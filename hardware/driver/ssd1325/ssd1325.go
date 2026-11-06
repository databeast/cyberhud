// Package ssd1325 drives SSD1325-based monochrome-white OLED panels over SPI —
// notably the Adafruit 2.7" 128x64 OLED graphic module (product 2674).
//
// Although the glass is a single colour (white), the SSD1325 is a 4-bit
// GRAYSCALE controller: each pixel is driven with 16 PWM brightness levels, two
// pixels packed per RAM byte (high nibble = left pixel, low nibble = right).
// The driver is therefore advertised as grayscale-fast, not mono — declaring it
// mono would under-drive the panel and strip the luminance detail the hardware
// can show, exactly the under-declaration the capability-conformance system
// exists to prevent.
//
// Addressing. Unlike the page-based SSD1306/SSD1305/SH1106, the SSD1325 uses a
// column/row window (set-column 0x15, set-row 0x75) and streams nibble-packed
// grayscale. Columns are addressed in pairs (one column address = two pixels),
// so a 128-wide panel spans column addresses 0..63.
package ssd1325

import (
	"errors"
	"image"
	"image/color"
	"image/draw"
	"time"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/conn/v3/spi"

	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/hardware/driver"
)

const (
	ssd1325Width  = 128
	ssd1325Height = 64

	cmdSetColumnAddr byte = 0x15
	cmdSetRowAddr    byte = 0x75
	cmdSetContrast   byte = 0x81
	cmdMasterCurrent byte = 0x87
	cmdSetRemap      byte = 0xA0
	cmdSetStartLine  byte = 0xA1
	cmdSetOffset     byte = 0xA2
	cmdNormalDisplay byte = 0xA4
	cmdSetMultiplex  byte = 0xA8
	cmdDisplayOff    byte = 0xAE
	cmdDisplayOn     byte = 0xAF
	cmdSetPhaseLen   byte = 0xB1
	cmdSetClock      byte = 0xB3
	cmdGrayTable     byte = 0xB8
	cmdDefaultGray   byte = 0xB9
	cmdSetPrecharge  byte = 0xBC
	cmdSetVCOMH      byte = 0xBE
	cmdSetVSL        byte = 0xBF
)

func init() {
	driver.Register(driver.Definition{
		ID:         "ssd1325",
		Title:      "SSD1325",
		Summary:    "SPI 4-bit grayscale OLED controller used by 128x64 white OLED graphic modules.",
		Monochrome: false,
		ColorDepth: driver.DepthGray,
		Capability: textlayout.CapGrayscaleFast,
		OptionDefs: []driver.OptionDefinition{
			{Key: "width", Type: "int", Summary: "Display width in pixels.", Default: "128"},
			{Key: "height", Type: "int", Summary: "Display height in pixels.", Default: "64"},
			{Key: "spi_hz", Type: "frequency", Summary: "SPI bus speed.", Default: "8MHz"},
		},
		NewSPI: func(port spi.Port, dc, rst, bl gpio.PinOut, _ gpio.PinIn, cfg driver.DriverConfig) (driver.DrawTarget, error) {
			return New(port, dc, rst, bl, cfg)
		},
	})
}

// Dev is a handle to an SSD1325 grayscale OLED.
type Dev struct {
	c      spi.Conn
	dc     gpio.PinOut
	rst    gpio.PinOut
	bl     gpio.PinOut
	width  int
	height int
}

// New initialises the SSD1325 and returns a ready-to-use Dev.
func New(port spi.Port, dc, rst, bl gpio.PinOut, cfg driver.DriverConfig) (*Dev, error) {
	if port == nil {
		return nil, errors.New("display: spi port must not be nil")
	}
	if dc == nil || rst == nil {
		return nil, errors.New("display: dc and rst pins must not be nil")
	}
	if cfg.Width <= 0 {
		cfg.Width = ssd1325Width
	}
	if cfg.Height <= 0 {
		cfg.Height = ssd1325Height
	}
	if cfg.SPIHz <= 0 {
		cfg.SPIHz = 8 * physic.MegaHertz
	}
	conn, err := port.Connect(cfg.SPIHz, spi.Mode0, 8)
	if err != nil {
		return nil, err
	}
	d := &Dev{c: conn, dc: dc, rst: rst, bl: bl, width: cfg.Width, height: cfg.Height}
	if err := d.init(); err != nil {
		return nil, err
	}
	return d, nil
}

// Bounds returns the display rectangle.
func (d *Dev) Bounds() image.Rectangle { return image.Rect(0, 0, d.width, d.height) }

// TextHints reports this panel as a fast grayscale OLED (16 PWM levels).
func (d *Dev) TextHints() textlayout.TextHints {
	b := d.Bounds()
	return textlayout.TextHints{
		PixelWidth: b.Dx(), PixelHeight: b.Dy(),
		SupportsVerticalScroll: true, SupportsHorizontalScroll: true, SupportsAutoScroll: true,
		PreferEventRefresh:     false,
		Capability:             textlayout.CapGrayscaleFast,
		DefaultTickerDirection: textlayout.TickerDirectionVertical, DefaultLineMode: textlayout.LineModeTruncate,
	}
}

// DrawImage converts the source image to 4-bit grayscale and streams it. Two
// horizontally-adjacent pixels are packed per byte (high nibble = left pixel,
// low nibble = right pixel).
func (d *Dev) DrawImage(src draw.Image) error {
	// Address the full panel: columns in pairs (0..width/2-1), rows 0..height-1.
	if err := d.cmdData(cmdSetColumnAddr, []byte{0x00, byte(d.width/2 - 1)}); err != nil {
		return err
	}
	if err := d.cmdData(cmdSetRowAddr, []byte{0x00, byte(d.height - 1)}); err != nil {
		return err
	}

	lineBytes := d.width / 2
	buf := make([]byte, lineBytes*d.height)
	i := 0
	for y := 0; y < d.height; y++ {
		for x := 0; x < d.width; x += 2 {
			left := gray4(src.At(x, y))
			right := byte(0)
			if x+1 < d.width {
				right = gray4(src.At(x+1, y))
			}
			buf[i] = left<<4 | right
			i++
		}
	}
	return d.data(buf)
}

// gray4 maps an image colour to a 4-bit (0..15) grayscale level via BT.601
// luminance.
func gray4(c color.Color) byte {
	r, g, b, _ := c.RGBA()
	// driver.Luminance returns ~0..0xFFFF; scale to 0..15.
	lum := driver.Luminance(r, g, b) >> 12
	if lum > 15 {
		lum = 15
	}
	return byte(lum)
}

func (d *Dev) init() error {
	if err := d.rst.Out(gpio.High); err != nil {
		return err
	}
	time.Sleep(1 * time.Millisecond)
	if err := d.rst.Out(gpio.Low); err != nil {
		return err
	}
	time.Sleep(10 * time.Millisecond)
	if err := d.rst.Out(gpio.High); err != nil {
		return err
	}
	time.Sleep(10 * time.Millisecond)

	seq := []struct {
		cmd  byte
		data []byte
	}{
		{cmdDisplayOff, nil},
		{cmdSetClock, []byte{0xF1}},
		{cmdSetMultiplex, []byte{0x3F}},  // 64 rows
		{cmdSetOffset, []byte{0x4C}},     // display offset
		{cmdSetStartLine, []byte{0x00}},
		{cmdSetRemap, []byte{0x56}},      // column/nibble remap for this module
		{cmdSetContrast, []byte{0x7F}},
		{cmdMasterCurrent, []byte{0x0F}},
		{cmdSetPhaseLen, []byte{0x21}},
		{cmdSetPrecharge, []byte{0x10}},
		{cmdSetVCOMH, []byte{0x1C}},
		{cmdSetVSL, []byte{0x0D}},
		{cmdDefaultGray, nil}, // use the built-in linear grayscale table
		{cmdNormalDisplay, nil},
	}
	for _, s := range seq {
		if err := d.cmd(s.cmd); err != nil {
			return err
		}
		if len(s.data) > 0 {
			if err := d.data(s.data); err != nil {
				return err
			}
		}
	}
	if err := d.cmd(cmdDisplayOn); err != nil {
		return err
	}
	if d.bl != nil {
		return d.bl.Out(gpio.High)
	}
	return nil
}

func (d *Dev) cmdData(c byte, payload []byte) error {
	if err := d.cmd(c); err != nil {
		return err
	}
	if len(payload) == 0 {
		return nil
	}
	return d.data(payload)
}

func (d *Dev) cmd(c byte) error {
	if err := d.dc.Out(gpio.Low); err != nil {
		return err
	}
	return d.c.Tx([]byte{c}, nil)
}

func (d *Dev) data(buf []byte) error {
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
