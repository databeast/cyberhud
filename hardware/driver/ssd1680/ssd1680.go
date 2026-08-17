package ssd1680

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
	ssd1680Width  = 250
	ssd1680Height = 122

	ssd1680DriverOutputControl = 0x01
	ssd1680DataEntryMode       = 0x11
	ssd1680SWReset             = 0x12
	ssd1680MasterActivation    = 0x20
	ssd1680DisplayUpdateCtrl2  = 0x22
	ssd1680WriteRAM            = 0x24
	ssd1680WriteRAM2           = 0x26
	ssd1680BorderWaveform      = 0x3C
	ssd1680SetRAMXRange        = 0x44
	ssd1680SetRAMYRange        = 0x45
	ssd1680SetRAMXCounter      = 0x4E
	ssd1680SetRAMYCounter      = 0x4F
	ssd1680TempSensorControl   = 0x18

	ssd1680RefreshFull    = 0xF7
	ssd1680RefreshPartial = 0xFF
)

func init() {
	driver.Register(driver.Definition{
		ID:         "ssd1680",
		Title:      "SSD1680",
		Summary:    "SPI monochrome e-ink controller used by 2.13-inch panels with busy-pin driven refresh cycles.",
		Monochrome: true,
		IsEPaper:   true,
		OptionDefs: []driver.OptionDefinition{
			{Key: "width", Type: "int", Summary: "Logical panel width in pixels.", Default: "250"},
			{Key: "height", Type: "int", Summary: "Logical panel height in pixels.", Default: "122"},
			{Key: "spi_hz", Type: "frequency", Summary: "SPI bus speed used for the panel transport.", Default: "4MHz"},
			{Key: "busy_high", Type: "bool", Summary: "Whether the busy pin reads high while the controller is refreshing.", Default: "true"},
			{Key: "busy_timeout", Type: "duration", Summary: "Maximum time to wait for busy-pin refresh completion.", Default: "8s"},
			{Key: "full_refresh_every", Type: "int", Summary: "Number of partial refreshes before a maintenance full refresh.", Default: "20"},
		},
		DefaultText: textlayout.TextHints{
			PixelWidth: ssd1680Width, PixelHeight: ssd1680Height,
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

type SSD1680 struct {
	cfg          driver.DriverConfig
	c            spi.Conn
	dc           gpio.PinOut
	rst          gpio.PinOut
	bsy          gpio.PinIn
	ramWidth     int
	ramHeight    int
	rotate90     bool
	prev         []byte
	partialCount int
	firstRefresh bool
}

func New(port spi.Port, dc, rst gpio.PinOut, busy gpio.PinIn, cfg driver.DriverConfig) (*SSD1680, error) {
	if port == nil {
		return nil, errors.New("display: spi port must not be nil")
	}
	if dc == nil || rst == nil || busy == nil {
		return nil, errors.New("display: dc, rst, and busy pins must not be nil")
	}
	if cfg.Width <= 0 {
		cfg.Width = ssd1680Width
	}
	if cfg.Height <= 0 {
		cfg.Height = ssd1680Height
	}
	if cfg.SPIHz <= 0 {
		cfg.SPIHz = 4 * physic.MegaHertz
	}
	if cfg.BusyTimeout <= 0 {
		cfg.BusyTimeout = 8 * time.Second
	}
	if cfg.FullRefreshEvery <= 0 {
		cfg.FullRefreshEvery = 20
	}
	conn, err := port.Connect(cfg.SPIHz, spi.Mode0, 8)
	if err != nil {
		return nil, err
	}
	if err := busy.In(gpio.PullNoChange, gpio.NoEdge); err != nil {
		return nil, err
	}
	ramWidth, ramHeight, rotate90 := ssd1680RAMLayout(cfg)
	lineBytes := (ramWidth + 7) / 8
	prev := make([]byte, lineBytes*ramHeight)
	for i := range prev {
		prev[i] = 0xFF
	}
	d := &SSD1680{cfg: cfg, c: conn, dc: dc, rst: rst, bsy: busy, ramWidth: ramWidth, ramHeight: ramHeight, rotate90: rotate90, prev: prev, firstRefresh: true}
	if err := d.init(); err != nil {
		return nil, err
	}
	return d, nil
}

func (d *SSD1680) Bounds() image.Rectangle { return image.Rect(0, 0, d.cfg.Width, d.cfg.Height) }

func (d *SSD1680) TextHints() textlayout.TextHints {
	b := d.Bounds()
	return textlayout.TextHints{
		PixelWidth: b.Dx(), PixelHeight: b.Dy(),
		SupportsVerticalScroll: false, SupportsHorizontalScroll: false, SupportsAutoScroll: false,
		PreferEventRefresh:     true,
		Capability:             textlayout.CapMonoSlow,
		DefaultTickerDirection: textlayout.TickerDirectionNone, DefaultLineMode: textlayout.LineModeTruncate,
	}
}

func (d *SSD1680) DrawImage(src draw.Image) error {
	next := d.packFrame(src)
	if bytes.Equal(next, d.prev) {
		return nil
	}
	full := d.firstRefresh || d.cfg.FullRefreshEvery <= 0 || d.partialCount >= d.cfg.FullRefreshEvery
	if full {
		if err := d.writeFrame(next, next); err != nil {
			return err
		}
		if err := d.refresh(ssd1680RefreshFull); err != nil {
			return err
		}
		d.partialCount = 0
		d.firstRefresh = false
	} else {
		if err := d.writeFrame(next, d.prev); err != nil {
			return err
		}
		if err := d.refresh(ssd1680RefreshPartial); err != nil {
			return err
		}
		d.partialCount++
	}
	copy(d.prev, next)
	return nil
}

func (d *SSD1680) init() error {
	if err := d.rst.Out(gpio.Low); err != nil {
		return err
	}
	time.Sleep(10 * time.Millisecond)
	if err := d.rst.Out(gpio.High); err != nil {
		return err
	}
	time.Sleep(10 * time.Millisecond)
	if err := d.cmd(ssd1680SWReset); err != nil {
		return err
	}
	if err := d.waitBusy(); err != nil {
		return err
	}
	if err := d.cmdData(ssd1680DriverOutputControl, []byte{byte((d.ramHeight - 1) & 0xFF), byte(((d.ramHeight - 1) >> 8) & 0xFF), 0x00}); err != nil {
		return err
	}
	if err := d.cmdData(ssd1680DataEntryMode, []byte{0x03}); err != nil {
		return err
	}
	if err := d.cmdData(ssd1680BorderWaveform, []byte{0x05}); err != nil {
		return err
	}
	if err := d.cmdData(ssd1680TempSensorControl, []byte{0x80}); err != nil {
		return err
	}
	if err := d.setWindowAndCursor(); err != nil {
		return err
	}
	if err := d.writeFrame(d.prev, d.prev); err != nil {
		return err
	}
	return d.refresh(ssd1680RefreshFull)
}

func (d *SSD1680) packFrame(src draw.Image) []byte {
	b := d.Bounds()
	ramWidth, ramHeight := d.ramWidth, d.ramHeight
	if ramWidth <= 0 || ramHeight <= 0 {
		ramWidth, ramHeight, _ = ssd1680RAMLayout(d.cfg)
	}
	lineBytes := (ramWidth + 7) / 8
	buf := make([]byte, lineBytes*ramHeight)
	for i := range buf {
		buf[i] = 0xFF
	}
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			r, g, bl, _ := src.At(x, y).RGBA()
			on := luminanceOn(r, g, bl)
			lx := x - b.Min.X
			ly := y - b.Min.Y
			rx, ry := lx, ly
			if d.rotate90 {
				rx = (ramWidth - 1) - ly
				ry = lx
			}
			idx := ry*lineBytes + rx/8
			bit := uint(7 - (rx % 8))
			if on {
				buf[idx] &^= 1 << bit
			}
		}
	}
	return buf
}

func (d *SSD1680) writeFrame(curr []byte, prev []byte) error {
	if err := d.setWindowAndCursor(); err != nil {
		return err
	}
	if err := d.cmd(ssd1680WriteRAM); err != nil {
		return err
	}
	if err := d.data(curr); err != nil {
		return err
	}
	if err := d.setWindowAndCursor(); err != nil {
		return err
	}
	if err := d.cmd(ssd1680WriteRAM2); err != nil {
		return err
	}
	return d.data(prev)
}

func (d *SSD1680) refresh(mode byte) error {
	if err := d.cmdData(ssd1680DisplayUpdateCtrl2, []byte{mode}); err != nil {
		return err
	}
	if err := d.cmd(ssd1680MasterActivation); err != nil {
		return err
	}
	return d.waitBusy()
}

func (d *SSD1680) setWindowAndCursor() error {
	lineBytes := (d.ramWidth + 7) / 8
	if err := d.cmdData(ssd1680SetRAMXRange, []byte{0x00, byte(lineBytes - 1)}); err != nil {
		return err
	}
	yMax := d.ramHeight - 1
	if err := d.cmdData(ssd1680SetRAMYRange, []byte{0x00, 0x00, byte(yMax & 0xFF), byte((yMax >> 8) & 0xFF)}); err != nil {
		return err
	}
	if err := d.cmdData(ssd1680SetRAMXCounter, []byte{0x00}); err != nil {
		return err
	}
	return d.cmdData(ssd1680SetRAMYCounter, []byte{0x00, 0x00})
}

func (d *SSD1680) waitBusy() error {
	deadline := time.Now().Add(d.cfg.BusyTimeout)
	for {
		lvl := d.bsy.Read()
		busy := (d.cfg.BusyHigh && lvl == gpio.High) || (!d.cfg.BusyHigh && lvl == gpio.Low)
		if !busy {
			return nil
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("display: ssd1680 busy timeout after %s", d.cfg.BusyTimeout)
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func (d *SSD1680) cmdData(c byte, payload []byte) error {
	if err := d.cmd(c); err != nil {
		return err
	}
	if len(payload) == 0 {
		return nil
	}
	return d.data(payload)
}

func (d *SSD1680) cmd(c byte) error {
	if err := d.dc.Out(gpio.Low); err != nil {
		return err
	}
	return d.c.Tx([]byte{c}, nil)
}

func (d *SSD1680) data(buf []byte) error {
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

func ssd1680RAMLayout(cfg driver.DriverConfig) (width int, height int, rotate90 bool) {
	width = cfg.Width
	height = cfg.Height
	if cfg.Width > cfg.Height {
		width = cfg.Height
		height = cfg.Width
		rotate90 = true
	}
	return width, height, rotate90
}

func luminanceOn(r, g, b uint32) bool { luma := (299*r + 587*g + 114*b) / 1000; return luma > 0x3000 }
