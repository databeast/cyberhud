// Package ssd1681_tricolor drives SSD1681-based tri-color (black/red/white) SPI
// e-paper panels — notably the Adafruit 1.54" tri-color eInk breakout (product
// 4868), a 200x200 board that pairs the SSD1681 controller with an onboard SRAM
// buffer and a microSD socket.
//
// This is the colour sibling of the monochrome ssd1681 driver. The SSD1681
// exposes two 1-bpp image planes: the black/white plane (command 0x24) and a
// second colour plane (command 0x26). The monochrome driver repurposes 0x26 as
// a previous-frame buffer for partial refreshes; here 0x26 carries the RED ink,
// so partial refresh is not available and every update is a full refresh — the
// normal mode of operation for tri-color e-paper.
//
// Pixel classification follows the uc8179 tri-color convention: a pixel whose
// red channel clearly dominates is red; otherwise a low-luminance pixel is
// black and everything else is white.
package ssd1681_tricolor

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
	ssd1681Width  = 200
	ssd1681Height = 200

	ssd1681DriverOutputControl = 0x01
	ssd1681DataEntryMode       = 0x11
	ssd1681SWReset             = 0x12
	ssd1681MasterActivation    = 0x20
	ssd1681DisplayUpdateCtrl2  = 0x22
	ssd1681WriteRAMBlack       = 0x24 // black/white plane
	ssd1681WriteRAMRed         = 0x26 // red (second colour) plane
	ssd1681BorderWaveform      = 0x3C
	ssd1681SetRAMXRange        = 0x44
	ssd1681SetRAMYRange        = 0x45
	ssd1681SetRAMXCounter      = 0x4E
	ssd1681SetRAMYCounter      = 0x4F
	ssd1681TempSensorControl   = 0x18

	// Full-refresh update sequence for tri-color panels (load temperature, enable
	// clock/analog, display with the colour-aware waveform, disable).
	ssd1681RefreshFull = 0xF7
)

func init() {
	driver.Register(driver.Definition{
		ID:         "ssd1681-tricolor",
		Title:      "SSD1681 (tri-color)",
		Summary:    "SPI tri-color e-ink controller (black/white/red) used by 1.54-inch 200x200 panels.",
		Monochrome: false,
		// Tri-color (black/white/red) carries hue, so the native depth is color.
		// Slow refresh → color-slow. DrawImage drives the red plane separately
		// from the B/W plane, so a style's hue survives to the panel.
		ColorDepth: driver.DepthColor,
		Capability: textlayout.CapColorSlow,
		IsEPaper:   true,
		OptionDefs: []driver.OptionDefinition{
			{Key: "width", Type: "int", Summary: "Logical panel width in pixels.", Default: "200"},
			{Key: "height", Type: "int", Summary: "Logical panel height in pixels.", Default: "200"},
			{Key: "spi_hz", Type: "frequency", Summary: "SPI bus speed used for the panel transport.", Default: "4MHz"},
			{Key: "busy_high", Type: "bool", Summary: "Whether the busy pin reads high while the controller is refreshing.", Default: "true"},
			{Key: "busy_timeout", Type: "duration", Summary: "Maximum time to wait for busy-pin refresh completion.", Default: "30s"},
		},
		NewSPI: func(port spi.Port, dc, rst, _ gpio.PinOut, busy gpio.PinIn, cfg driver.DriverConfig) (driver.DrawTarget, error) {
			return New(port, dc, rst, busy, cfg)
		},
	})
}

// SSD1681TriColor is a handle to a tri-color SSD1681 e-paper panel.
type SSD1681TriColor struct {
	cfg       driver.DriverConfig
	c         spi.Conn
	dc        gpio.PinOut
	rst       gpio.PinOut
	bsy       gpio.PinIn
	ramWidth  int
	ramHeight int
	rotate90  bool
	prevBlack []byte
	prevRed   []byte
	firstDraw bool
}

// New initialises the tri-color SSD1681 and returns a ready-to-use handle.
func New(port spi.Port, dc, rst gpio.PinOut, busy gpio.PinIn, cfg driver.DriverConfig) (*SSD1681TriColor, error) {
	if port == nil {
		return nil, errors.New("display: spi port must not be nil")
	}
	if dc == nil || rst == nil || busy == nil {
		return nil, errors.New("display: dc, rst, and busy pins must not be nil")
	}
	if cfg.Width <= 0 {
		cfg.Width = ssd1681Width
	}
	if cfg.Height <= 0 {
		cfg.Height = ssd1681Height
	}
	if cfg.SPIHz <= 0 {
		cfg.SPIHz = 4 * physic.MegaHertz
	}
	if cfg.BusyTimeout <= 0 {
		// Tri-color full refreshes are slow — several seconds to ~15s.
		cfg.BusyTimeout = 30 * time.Second
	}
	conn, err := port.Connect(cfg.SPIHz, spi.Mode0, 8)
	if err != nil {
		return nil, err
	}
	if err := busy.In(gpio.PullNoChange, gpio.NoEdge); err != nil {
		return nil, err
	}
	ramWidth, ramHeight, rotate90 := ssd1681RAMLayout(cfg)
	d := &SSD1681TriColor{cfg: cfg, c: conn, dc: dc, rst: rst, bsy: busy, ramWidth: ramWidth, ramHeight: ramHeight, rotate90: rotate90, firstDraw: true}
	if err := d.init(); err != nil {
		return nil, err
	}
	return d, nil
}

// Bounds returns the display rectangle.
func (d *SSD1681TriColor) Bounds() image.Rectangle { return image.Rect(0, 0, d.cfg.Width, d.cfg.Height) }

// TextHints advertises this panel as a slow color e-ink display. Its DrawImage
// drives a separate red plane alongside black/white, so hue reaches the panel;
// declaring color-slow lets selection offer color styles rather than treating
// the panel as mono. E-ink semantics: no scrolling, event-driven refresh.
func (d *SSD1681TriColor) TextHints() textlayout.TextHints {
	b := d.Bounds()
	return textlayout.TextHints{
		PixelWidth: b.Dx(), PixelHeight: b.Dy(),
		SupportsVerticalScroll: false, SupportsHorizontalScroll: false, SupportsAutoScroll: false,
		PreferEventRefresh:     true,
		Capability:             textlayout.CapColorSlow,
		DefaultTickerDirection: textlayout.TickerDirectionNone, DefaultLineMode: textlayout.LineModeTruncate,
	}
}

// DrawImage classifies each pixel into black/red/white, packs the black and red
// planes, and performs a full refresh. Tri-color panels do not support partial
// refresh, so every changed frame triggers a full update; an unchanged frame is
// skipped.
func (d *SSD1681TriColor) DrawImage(src draw.Image) error {
	black, red := d.packFrames(src)
	if !d.firstDraw && bytes.Equal(black, d.prevBlack) && bytes.Equal(red, d.prevRed) {
		return nil
	}
	if err := d.setWindowAndCursor(); err != nil {
		return err
	}
	if err := d.cmd(ssd1681WriteRAMBlack); err != nil {
		return err
	}
	if err := d.data(black); err != nil {
		return err
	}
	if err := d.setWindowAndCursor(); err != nil {
		return err
	}
	if err := d.cmd(ssd1681WriteRAMRed); err != nil {
		return err
	}
	if err := d.data(red); err != nil {
		return err
	}
	if err := d.refresh(); err != nil {
		return err
	}
	d.prevBlack = black
	d.prevRed = red
	d.firstDraw = false
	return nil
}

func (d *SSD1681TriColor) init() error {
	if err := d.rst.Out(gpio.Low); err != nil {
		return err
	}
	time.Sleep(10 * time.Millisecond)
	if err := d.rst.Out(gpio.High); err != nil {
		return err
	}
	time.Sleep(10 * time.Millisecond)
	if err := d.cmd(ssd1681SWReset); err != nil {
		return err
	}
	if err := d.waitBusy(); err != nil {
		return err
	}
	if err := d.cmdData(ssd1681DriverOutputControl, []byte{byte((d.ramHeight - 1) & 0xFF), byte(((d.ramHeight - 1) >> 8) & 0xFF), 0x00}); err != nil {
		return err
	}
	if err := d.cmdData(ssd1681DataEntryMode, []byte{0x03}); err != nil {
		return err
	}
	// Border waveform 0x05 keeps the border white on tri-color panels.
	if err := d.cmdData(ssd1681BorderWaveform, []byte{0x05}); err != nil {
		return err
	}
	if err := d.cmdData(ssd1681TempSensorControl, []byte{0x80}); err != nil {
		return err
	}
	return d.setWindowAndCursor()
}

// packFrames classifies the source image into two 1-bpp planes. In the black/
// white plane a set bit is white and a cleared bit is black (so the buffer
// starts all-white, 0xFF). In the red plane a set bit is red and a cleared bit
// leaves the pixel to the black/white plane (buffer starts all-zero).
func (d *SSD1681TriColor) packFrames(src draw.Image) (black []byte, red []byte) {
	b := d.Bounds()
	ramWidth, ramHeight := d.ramWidth, d.ramHeight
	if ramWidth <= 0 || ramHeight <= 0 {
		ramWidth, ramHeight, _ = ssd1681RAMLayout(d.cfg)
	}
	lineBytes := (ramWidth + 7) / 8
	black = make([]byte, lineBytes*ramHeight)
	red = make([]byte, lineBytes*ramHeight)
	for i := range black {
		black[i] = 0xFF // white
	}
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			r, g, bl, _ := src.At(x, y).RGBA()
			isRed := isRedPixel(r, g, bl)
			isBlack := !isRed && driver.Luminance(r, g, bl) < driver.MonoThreshold

			lx := x - b.Min.X
			ly := y - b.Min.Y
			rx, ry := lx, ly
			if d.rotate90 {
				rx = (ramWidth - 1) - ly
				ry = lx
			}
			idx := ry*lineBytes + rx/8
			bit := uint(7 - (rx % 8))
			if isBlack {
				black[idx] &^= 1 << bit // clear bit = black
			}
			if isRed {
				red[idx] |= 1 << bit // set bit = red
			}
		}
	}
	return black, red
}

func (d *SSD1681TriColor) refresh() error {
	if err := d.cmdData(ssd1681DisplayUpdateCtrl2, []byte{ssd1681RefreshFull}); err != nil {
		return err
	}
	if err := d.cmd(ssd1681MasterActivation); err != nil {
		return err
	}
	return d.waitBusy()
}

func (d *SSD1681TriColor) setWindowAndCursor() error {
	lineBytes := (d.ramWidth + 7) / 8
	if err := d.cmdData(ssd1681SetRAMXRange, []byte{0x00, byte(lineBytes - 1)}); err != nil {
		return err
	}
	yMax := d.ramHeight - 1
	if err := d.cmdData(ssd1681SetRAMYRange, []byte{0x00, 0x00, byte(yMax & 0xFF), byte((yMax >> 8) & 0xFF)}); err != nil {
		return err
	}
	if err := d.cmdData(ssd1681SetRAMXCounter, []byte{0x00}); err != nil {
		return err
	}
	return d.cmdData(ssd1681SetRAMYCounter, []byte{0x00, 0x00})
}

func (d *SSD1681TriColor) waitBusy() error {
	deadline := time.Now().Add(d.cfg.BusyTimeout)
	for {
		lvl := d.bsy.Read()
		busy := (d.cfg.BusyHigh && lvl == gpio.High) || (!d.cfg.BusyHigh && lvl == gpio.Low)
		if !busy {
			return nil
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("display: ssd1681-tricolor busy timeout after %s", d.cfg.BusyTimeout)
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func (d *SSD1681TriColor) cmdData(c byte, payload []byte) error {
	if err := d.cmd(c); err != nil {
		return err
	}
	if len(payload) == 0 {
		return nil
	}
	return d.data(payload)
}

func (d *SSD1681TriColor) cmd(c byte) error {
	if err := d.dc.Out(gpio.Low); err != nil {
		return err
	}
	return d.c.Tx([]byte{c}, nil)
}

func (d *SSD1681TriColor) data(buf []byte) error {
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

// isRedPixel reports whether the red channel clearly dominates green and blue,
// matching the uc8179 tri-color convention. Inputs are 16-bit channel values as
// returned by color.RGBA().
func isRedPixel(r, g, b uint32) bool {
	return r > 0x8000 && g < 0x4000 && b < 0x4000
}

// ssd1681RAMLayout normalises the configured dimensions to the controller's
// native portrait RAM orientation, reporting whether a 90° rotation is needed.
// For the square 200x200 panel the dimensions are equal, so rotate90 is false.
func ssd1681RAMLayout(cfg driver.DriverConfig) (width int, height int, rotate90 bool) {
	width = cfg.Width
	height = cfg.Height
	if cfg.Width > cfg.Height {
		width = cfg.Height
		height = cfg.Width
		rotate90 = true
	}
	return width, height, rotate90
}
