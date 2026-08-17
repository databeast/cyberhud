package epd4in26

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/draw"
	"time"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/conn/v3/spi"

	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/hardware/driver"
)

const (
	epd4in26Width  = 800
	epd4in26Height = 480

	cmdSWReset            byte = 0x12
	cmdDriverOutputCtrl   byte = 0x01
	cmdDataEntryMode      byte = 0x11
	cmdBorderWaveform     byte = 0x3C
	cmdTempSensor         byte = 0x18
	cmdSoftStart          byte = 0x0C
	cmdDisplayUpdateCtrl2 byte = 0x22
	cmdMasterActivation   byte = 0x20
	cmdWriteRAMBW         byte = 0x24
	cmdWriteRAMRed        byte = 0x26
	cmdSetRAMXRange       byte = 0x44
	cmdSetRAMYRange       byte = 0x45
	cmdSetRAMXCounter     byte = 0x4E
	cmdSetRAMYCounter     byte = 0x4F
	cmdDeepSleep          byte = 0x10

	refreshFull byte = 0xF7
)

func init() {
	driver.Register(driver.Definition{
		ID:         "epd4in26",
		Title:      "EPD 4.26\"",
		Summary:    "SPI monochrome e-ink controller for Waveshare 4.26-inch 800x480 B/W display.",
		Monochrome: true,
		IsEPaper:   true,
		OptionDefs: []driver.OptionDefinition{
			{Key: "width", Type: "int", Summary: "Display width in pixels.", Default: "800"},
			{Key: "height", Type: "int", Summary: "Display height in pixels.", Default: "480"},
			{Key: "spi_hz", Type: "frequency", Summary: "SPI bus speed.", Default: "4MHz"},
			{Key: "busy_timeout", Type: "duration", Summary: "Maximum time to wait for refresh.", Default: "30s"},
		},
		DefaultText: textlayout.TextHints{
			PixelWidth: epd4in26Width, PixelHeight: epd4in26Height,
			SupportsVerticalScroll: false, SupportsHorizontalScroll: false, SupportsAutoScroll: false,
			PreferEventRefresh:     true,
			Capability:             textlayout.CapMonoSlow,
			DefaultTickerDirection: textlayout.TickerDirectionNone, DefaultLineMode: textlayout.LineModeTruncate,
		},
		NewSPI: func(port spi.Port, dc, rst, _ gpio.PinOut, busy gpio.PinIn, cfg driver.DriverConfig) (driver.DrawTarget, error) {
			return New(port, dc, rst, busy, cfg)
		},
	})
}

// EPD4in26 is a handle to a Waveshare 4.26" 800x480 B/W e-ink display.
type EPD4in26 struct {
	cfg  driver.DriverConfig
	c    spi.Conn
	dc   gpio.PinOut
	rst  gpio.PinOut
	bsy  gpio.PinIn
	prev []byte
}

func New(port spi.Port, dc, rst gpio.PinOut, busy gpio.PinIn, cfg driver.DriverConfig) (*EPD4in26, error) {
	if port == nil {
		return nil, errors.New("display: spi port must not be nil")
	}
	if dc == nil || rst == nil || busy == nil {
		return nil, errors.New("display: dc, rst, and busy pins must not be nil")
	}
	if cfg.Width <= 0 {
		cfg.Width = epd4in26Width
	}
	if cfg.Height <= 0 {
		cfg.Height = epd4in26Height
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
	bufSize := (cfg.Width / 8) * cfg.Height
	prev := make([]byte, bufSize)
	for i := range prev {
		prev[i] = 0xFF
	}
	d := &EPD4in26{cfg: cfg, c: conn, dc: dc, rst: rst, bsy: busy, prev: prev}
	if err := d.init(); err != nil {
		return nil, err
	}
	return d, nil
}

func (d *EPD4in26) Bounds() image.Rectangle {
	return image.Rect(0, 0, d.cfg.Width, d.cfg.Height)
}

func (d *EPD4in26) TextHints() textlayout.TextHints {
	b := d.Bounds()
	return textlayout.TextHints{
		PixelWidth: b.Dx(), PixelHeight: b.Dy(),
		SupportsVerticalScroll: false, SupportsHorizontalScroll: false, SupportsAutoScroll: false,
		PreferEventRefresh:     true,
		Capability:             textlayout.CapMonoSlow,
		DefaultTickerDirection: textlayout.TickerDirectionNone, DefaultLineMode: textlayout.LineModeTruncate,
	}
}

func (d *EPD4in26) DrawImage(src draw.Image) error {
	next := d.packFrame(src)
	if bytes.Equal(next, d.prev) {
		return nil
	}
	if err := d.cmd(cmdWriteRAMBW); err != nil {
		return err
	}
	if err := d.data(next); err != nil {
		return err
	}
	if err := d.refresh(); err != nil {
		return err
	}
	copy(d.prev, next)
	return nil
}

func (d *EPD4in26) init() error {
	// Hardware reset
	if err := d.rst.Out(gpio.High); err != nil {
		return err
	}
	time.Sleep(20 * time.Millisecond)
	if err := d.rst.Out(gpio.Low); err != nil {
		return err
	}
	time.Sleep(2 * time.Millisecond)
	if err := d.rst.Out(gpio.High); err != nil {
		return err
	}
	time.Sleep(20 * time.Millisecond)

	if err := d.waitBusy(); err != nil {
		return err
	}

	// Software reset
	if err := d.cmd(cmdSWReset); err != nil {
		return err
	}
	if err := d.waitBusy(); err != nil {
		return err
	}

	// Internal temperature sensor
	if err := d.cmdData(cmdTempSensor, []byte{0x80}); err != nil {
		return err
	}

	// Soft start
	if err := d.cmdData(cmdSoftStart, []byte{0xAE, 0xC7, 0xC3, 0xC0, 0x80}); err != nil {
		return err
	}

	// Driver output control: height - 1
	h := d.cfg.Height
	if err := d.cmdData(cmdDriverOutputCtrl, []byte{byte((h - 1) & 0xFF), byte(((h - 1) >> 8) & 0xFF), 0x02}); err != nil {
		return err
	}

	// Border waveform
	if err := d.cmdData(cmdBorderWaveform, []byte{0x01}); err != nil {
		return err
	}

	// Data entry mode: X+, Y-
	if err := d.cmdData(cmdDataEntryMode, []byte{0x01}); err != nil {
		return err
	}

	// Set RAM window and cursor
	if err := d.setWindow(0, d.cfg.Height-1, d.cfg.Width-1, 0); err != nil {
		return err
	}
	if err := d.setCursor(0, 0); err != nil {
		return err
	}

	return d.waitBusy()
}

func (d *EPD4in26) packFrame(src draw.Image) []byte {
	w := d.cfg.Width
	h := d.cfg.Height
	lineBytes := w / 8
	buf := make([]byte, lineBytes*h)
	for i := range buf {
		buf[i] = 0xFF
	}
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, g, b, _ := src.At(x, y).RGBA()
			// Convention: modes draw white/bright on black background.
			// E-paper: bit set (1) = white, bit cleared (0) = black.
			// So bright source pixels should CLEAR bits (become black on e-paper)
			// to produce the inverted look: bright content → dark ink on white paper.
			luma := (299*r + 587*g + 114*b) / 1000
			if luma >= 0x8000 {
				buf[y*lineBytes+x/8] &^= 0x80 >> uint(x%8)
			}
		}
	}
	return buf
}

func (d *EPD4in26) refresh() error {
	if err := d.cmdData(cmdDisplayUpdateCtrl2, []byte{refreshFull}); err != nil {
		return err
	}
	if err := d.cmd(cmdMasterActivation); err != nil {
		return err
	}
	return d.waitBusy()
}

func (d *EPD4in26) setWindow(xStart, yStart, xEnd, yEnd int) error {
	if err := d.cmdData(cmdSetRAMXRange, []byte{
		byte(xStart & 0xFF), byte((xStart >> 8) & 0x03),
		byte(xEnd & 0xFF), byte((xEnd >> 8) & 0x03),
	}); err != nil {
		return err
	}
	return d.cmdData(cmdSetRAMYRange, []byte{
		byte(yStart & 0xFF), byte((yStart >> 8) & 0xFF),
		byte(yEnd & 0xFF), byte((yEnd >> 8) & 0xFF),
	})
}

func (d *EPD4in26) setCursor(x, y int) error {
	if err := d.cmdData(cmdSetRAMXCounter, []byte{
		byte(x & 0xFF), byte((x >> 8) & 0x03),
	}); err != nil {
		return err
	}
	return d.cmdData(cmdSetRAMYCounter, []byte{
		byte(y & 0xFF), byte((y >> 8) & 0xFF),
	})
}

func (d *EPD4in26) waitBusy() error {
	deadline := time.Now().Add(d.cfg.BusyTimeout)
	for {
		// Busy is active-HIGH on this panel (busy while pin reads HIGH)
		if d.bsy.Read() == gpio.Low {
			time.Sleep(20 * time.Millisecond)
			return nil
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("display: epd4in26 busy timeout after %s", d.cfg.BusyTimeout)
		}
		time.Sleep(20 * time.Millisecond)
	}
}

func (d *EPD4in26) cmdData(c byte, payload []byte) error {
	if err := d.cmd(c); err != nil {
		return err
	}
	if len(payload) == 0 {
		return nil
	}
	return d.data(payload)
}

func (d *EPD4in26) cmd(c byte) error {
	if err := d.dc.Out(gpio.Low); err != nil {
		return err
	}
	return d.c.Tx([]byte{c}, nil)
}

func (d *EPD4in26) data(buf []byte) error {
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
