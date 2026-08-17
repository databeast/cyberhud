package uc8179

import (
	"errors"
	"fmt"
	"image"
	"image/draw"
	"time"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/conn/v3/spi"

	"github.com/databeast/cyberhud/hardware/driver"
)

const (
	uc8179Width  = 400
	uc8179Height = 300

	cmdPanelSetting     byte = 0x00
	cmdPowerSetting     byte = 0x01
	cmdPowerOn          byte = 0x04
	cmdBoosterSoftStart byte = 0x06
	cmdDataStartTx1     byte = 0x10
	cmdDataStartTx2     byte = 0x13
	cmdDisplayRefresh   byte = 0x12
	cmdPLLControl       byte = 0x30
	cmdVCOMDataInterval byte = 0x50
	cmdTCONResolution   byte = 0x61
	cmdSPISetting       byte = 0x62
)

func init() {
	driver.Register(driver.Definition{
		ID:         "uc8179",
		Title:      "UC8179",
		Summary:    "SPI tri-color e-ink controller (black/white/red, 400x300).",
		Monochrome: false,
		IsEPaper:   true,
		OptionDefs: []driver.OptionDefinition{
			{Key: "width", Type: "int", Summary: "Display width in pixels.", Default: "400"},
			{Key: "height", Type: "int", Summary: "Display height in pixels.", Default: "300"},
			{Key: "spi_hz", Type: "frequency", Summary: "SPI bus speed.", Default: "4MHz"},
			{Key: "busy_high", Type: "bool", Summary: "Busy pin active level.", Default: "true"},
			{Key: "busy_timeout", Type: "duration", Summary: "Max wait for refresh.", Default: "30s"},
		},
		NewSPI: func(port spi.Port, dc, rst, _ gpio.PinOut, busy gpio.PinIn, cfg driver.DriverConfig) (driver.DrawTarget, error) {
			return New(port, dc, rst, busy, cfg)
		},
	})
}

// UC8179 drives a UC8179 tri-color e-ink controller over SPI.
type UC8179 struct {
	cfg driver.DriverConfig
	c   spi.Conn
	dc  gpio.PinOut
	rst gpio.PinOut
	bsy gpio.PinIn
}

// New creates and initializes a UC8179 e-ink display.
func New(port spi.Port, dc, rst gpio.PinOut, busy gpio.PinIn, cfg driver.DriverConfig) (*UC8179, error) {
	if port == nil {
		return nil, errors.New("display: spi port must not be nil")
	}
	if dc == nil || rst == nil || busy == nil {
		return nil, errors.New("display: dc, rst, and busy pins must not be nil")
	}
	if cfg.Width <= 0 {
		cfg.Width = uc8179Width
	}
	if cfg.Height <= 0 {
		cfg.Height = uc8179Height
	}
	if cfg.SPIHz <= 0 {
		cfg.SPIHz = 4 * physic.MegaHertz
	}
	if cfg.BusyTimeout <= 0 {
		cfg.BusyTimeout = 30 * time.Second
	}

	conn, err := port.Connect(cfg.SPIHz, spi.Mode0, 8)
	if err != nil {
		return nil, err
	}
	if err := busy.In(gpio.PullNoChange, gpio.NoEdge); err != nil {
		return nil, err
	}

	d := &UC8179{cfg: cfg, c: conn, dc: dc, rst: rst, bsy: busy}
	if err := d.init(); err != nil {
		return nil, err
	}
	return d, nil
}

func (d *UC8179) Bounds() image.Rectangle {
	return image.Rect(0, 0, d.cfg.Width, d.cfg.Height)
}

// DrawImage converts pixels to B/W and red channel data for the tri-color panel.
// Black/white is determined by luminance threshold. Red channel is extracted when
// the red component dominates green and blue.
func (d *UC8179) DrawImage(src draw.Image) error {
	w := d.cfg.Width
	h := d.cfg.Height
	lineBytes := (w + 7) / 8
	bwBuf := make([]byte, lineBytes*h)
	redBuf := make([]byte, lineBytes*h)

	// Initialize B/W buffer to white (1=white, 0=black).
	for i := range bwBuf {
		bwBuf[i] = 0xFF
	}

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, g, b, _ := src.At(x, y).RGBA()
			idx := y*lineBytes + x/8
			bit := uint(7 - (x % 8))

			isRed := isRedPixel(r, g, b)
			isBlack := !isRed && luminanceOn(r, g, b)

			if isBlack {
				bwBuf[idx] &^= 1 << bit // set black
			}
			if isRed {
				redBuf[idx] |= 1 << bit // set red
			}
		}
	}

	// Write B/W data.
	if err := d.cmd(cmdDataStartTx1); err != nil {
		return err
	}
	if err := d.data(bwBuf); err != nil {
		return err
	}

	// Write red channel data.
	if err := d.cmd(cmdDataStartTx2); err != nil {
		return err
	}
	if err := d.data(redBuf); err != nil {
		return err
	}

	// Trigger refresh and wait for completion.
	if err := d.cmd(cmdDisplayRefresh); err != nil {
		return err
	}
	return d.waitBusy()
}

func (d *UC8179) init() error {
	// Hardware reset.
	if err := d.rst.Out(gpio.Low); err != nil {
		return err
	}
	time.Sleep(10 * time.Millisecond)
	if err := d.rst.Out(gpio.High); err != nil {
		return err
	}
	time.Sleep(100 * time.Millisecond)
	if err := d.waitBusy(); err != nil {
		return err
	}

	// Panel setting: tri-color, LUT from register.
	if err := d.cmdData(cmdPanelSetting, []byte{0x0F}); err != nil {
		return err
	}

	// Power setting.
	if err := d.cmdData(cmdPowerSetting, []byte{0x07, 0x07, 0x3F, 0x3F}); err != nil {
		return err
	}
	// Booster soft start.
	if err := d.cmdData(cmdBoosterSoftStart, []byte{0x17, 0x17, 0x28, 0x17}); err != nil {
		return err
	}
	// PLL control.
	if err := d.cmdData(cmdPLLControl, []byte{0x06}); err != nil {
		return err
	}
	// Power on.
	if err := d.cmd(cmdPowerOn); err != nil {
		return err
	}
	if err := d.waitBusy(); err != nil {
		return err
	}
	// Resolution.
	wH := byte((d.cfg.Width >> 8) & 0xFF)
	wL := byte(d.cfg.Width & 0xFF)
	hH := byte((d.cfg.Height >> 8) & 0xFF)
	hL := byte(d.cfg.Height & 0xFF)
	if err := d.cmdData(cmdTCONResolution, []byte{wH, wL, hH, hL}); err != nil {
		return err
	}
	// VCOM and data interval.
	if err := d.cmdData(cmdVCOMDataInterval, []byte{0x10, 0x07}); err != nil {
		return err
	}
	return nil
}

func (d *UC8179) waitBusy() error {
	deadline := time.Now().Add(d.cfg.BusyTimeout)
	for {
		lvl := d.bsy.Read()
		busy := (d.cfg.BusyHigh && lvl == gpio.High) || (!d.cfg.BusyHigh && lvl == gpio.Low)
		if !busy {
			return nil
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("display: uc8179 busy timeout after %s", d.cfg.BusyTimeout)
		}
		time.Sleep(50 * time.Millisecond)
	}
}

func (d *UC8179) cmdData(c byte, payload []byte) error {
	if err := d.cmd(c); err != nil {
		return err
	}
	if len(payload) == 0 {
		return nil
	}
	return d.data(payload)
}

func (d *UC8179) cmd(c byte) error {
	if err := d.dc.Out(gpio.Low); err != nil {
		return err
	}
	return d.c.Tx([]byte{c}, nil)
}

func (d *UC8179) data(buf []byte) error {
	if err := d.dc.Out(gpio.High); err != nil {
		return err
	}
	const chunk = 4096
	for len(buf) > 0 {
		n := len(buf)
		if n > chunk {
			n = chunk
		}
		if err := d.c.Tx(buf[:n], nil); err != nil {
			return err
		}
		buf = buf[n:]
	}
	return nil
}

// isRedPixel returns true when the red channel significantly dominates.
func isRedPixel(r, g, b uint32) bool {
	return r > 0x8000 && g < 0x4000 && b < 0x4000
}

// luminanceOn returns true if the pixel is dark (for black marking).
func luminanceOn(r, g, b uint32) bool {
	return (299*r+587*g+114*b)/1000 < 0x3000
}
