// Package ili9486 drives ILI9486-based SPI colour TFT panels — notably the
// ubiquitous 3.5" 480x320 Raspberry Pi LCDs (e.g. Waveshare 3.5" RPi LCD (A)),
// which pair an ILI9486 display controller with an XPT2046 resistive touch
// controller on SPI.
//
// Pixel format. These Waveshare 3.5" SPI boards consume 16-bit colour (RGB565),
// two bytes per pixel, like the ST7789. We set COLMOD=0x55 and emit 2
// bytes/pixel. (Sending 18-bit/3-byte RGB666 instead overruns each row — the
// right-hand columns go unaddressed — and desyncs the row pointer so alternate
// scanlines are skipped.)
//
// Orientation (MADCTL) and RGB/BGR order on these panels vary by board
// revision; the defaults here target 480x320 landscape with BGR order, which is
// the common Waveshare (A) configuration, but may need adjustment on device.
package ili9486

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
	cmdSWRESET = 0x01
	cmdSLPOUT  = 0x11
	cmdINVOFF  = 0x20
	cmdINVON   = 0x21
	cmdDISPON  = 0x29
	cmdCASET   = 0x2A
	cmdRASET   = 0x2B
	cmdRAMWR   = 0x2C
	cmdMADCTL  = 0x36
	cmdCOLMOD  = 0x3A

	// ILI9486-specific power/gamma registers used by the init sequence.
	cmdPWCTRL1 = 0xC0
	cmdPWCTRL2 = 0xC1
	cmdVMCTRL  = 0xC5
	cmdPGAMMA  = 0xE0
	cmdNGAMMA  = 0xE1

	// COLMOD value selecting 16-bit/pixel (RGB565), 2 bytes per pixel over SPI.
	// These Waveshare SPI boards consume 16-bit pixels; sending 18-bit/3-byte
	// overruns each row (right-hand columns go unaddressed) and desyncs the row
	// pointer (alternate scanlines skipped).
	colmod16bit = 0x55
)

func init() {
	driver.Register(driver.Definition{
		ID:         "ili9486",
		Title:      "ILI9486",
		Summary:    "SPI color LCD controller used by 480x320 3.5\" Raspberry Pi TFT panels (18-bit/pixel).",
		Monochrome: false,
		OptionDefs: []driver.OptionDefinition{
			{Key: "width", Type: "int", Summary: "Logical panel width in pixels.", Default: "480"},
			{Key: "height", Type: "int", Summary: "Logical panel height in pixels.", Default: "320"},
			{Key: "spi_hz", Type: "frequency", Summary: "SPI bus speed used for the panel transport.", Default: "16MHz"},
			{Key: "madctl", Type: "byte", Summary: "ILI9486 MADCTL register value controlling orientation and RGB/BGR order.", Default: "0x28"},
		},
		DefaultText: textlayout.DefaultTextHints(image.Rect(0, 0, 480, 320)),
		NewSPI: func(port spi.Port, dc, rst, bl gpio.PinOut, _ gpio.PinIn, cfg driver.DriverConfig) (driver.DrawTarget, error) {
			return New(port, dc, rst, bl, cfg)
		},
	})
}

// Dev is a handle to an ILI9486 LCD.
type Dev struct {
	c      spi.Conn
	dc     gpio.PinOut
	rst    gpio.PinOut
	bl     gpio.PinOut
	width  int
	height int
	madctl byte
	invert bool
}

// New initialises the ILI9486 and returns a ready-to-use Dev.
func New(port spi.Port, dc, rst, bl gpio.PinOut, cfg driver.DriverConfig) (*Dev, error) {
	if port == nil {
		return nil, errors.New("display: spi port must not be nil")
	}
	if dc == nil || rst == nil {
		return nil, errors.New("display: dc and rst pins must not be nil")
	}
	if cfg.Width <= 0 {
		cfg.Width = 480
	}
	if cfg.Height <= 0 {
		cfg.Height = 320
	}
	if cfg.SPIHz <= 0 {
		// These boards are conservative on SPI clock; 16MHz is a safe default.
		cfg.SPIHz = 16 * physic.MegaHertz
	}
	if cfg.MADCTL == 0 {
		// 0x28: MX off, MV set (row/column exchange → landscape), BGR order.
		// Landscape 480x320 with BGR is the common Waveshare (A) setup.
		cfg.MADCTL = 0x28
	}
	conn, err := port.Connect(cfg.SPIHz, spi.Mode0, 8)
	if err != nil {
		return nil, err
	}
	d := &Dev{
		c:      conn,
		dc:     dc,
		rst:    rst,
		bl:     bl,
		width:  cfg.Width,
		height: cfg.Height,
		madctl: cfg.MADCTL,
		invert: cfg.InvertColors,
	}
	if err := d.init(); err != nil {
		return nil, err
	}
	return d, nil
}

// Bounds returns the display rectangle.
func (d *Dev) Bounds() image.Rectangle { return image.Rect(0, 0, d.width, d.height) }

// TextHints reports text/scroll capabilities for ILI9486 color panels.
func (d *Dev) TextHints() textlayout.TextHints {
	b := d.Bounds()
	return textlayout.TextHints{
		PixelWidth: b.Dx(), PixelHeight: b.Dy(),
		SupportsVerticalScroll: true, SupportsHorizontalScroll: true, SupportsAutoScroll: true,
		PreferEventRefresh:     false,
		Capability:             textlayout.CapColorFast,
		DefaultTickerDirection: textlayout.TickerDirectionVertical, DefaultLineMode: textlayout.LineModeTruncate,
	}
}

// Draw renders src into the display window defined by r. Pixels are emitted as
// 18-bit RGB666, three bytes per pixel (see package doc).
func (d *Dev) Draw(r image.Rectangle, src image.Image, sp image.Point) error {
	r = r.Intersect(d.Bounds())
	if r.Empty() {
		return nil
	}
	if err := d.setWindow(r.Min.X, r.Min.Y, r.Max.X-1, r.Max.Y-1); err != nil {
		return err
	}
	buf := make([]byte, r.Dx()*r.Dy()*2)
	i := 0
	for y := r.Min.Y; y < r.Max.Y; y++ {
		for x := r.Min.X; x < r.Max.X; x++ {
			c := src.At(x+sp.X-r.Min.X, y+sp.Y-r.Min.Y)
			r16, g16, b16, _ := c.RGBA()
			// RGB565: 5 bits red, 6 bits green, 5 bits blue, 2 bytes/pixel.
			rgb := (r16>>11)<<11 | (g16>>10)<<5 | (b16 >> 11)
			buf[i] = byte(rgb >> 8)
			buf[i+1] = byte(rgb)
			i += 2
		}
	}
	return d.data(buf)
}

// Fill fills the entire display with a single colour.
func (d *Dev) Fill(c color.Color) error {
	return d.Draw(d.Bounds(), image.NewUniform(c), image.Point{})
}

// Clear turns all pixels black.
func (d *Dev) Clear() error { return d.Fill(color.Black) }

// DrawImage draws an RGBA image starting at the top-left of the panel.
func (d *Dev) DrawImage(img draw.Image) error { return d.Draw(d.Bounds(), img, image.Point{}) }

// BacklightOn sets the backlight to on or off. The Waveshare (A) ties the
// backlight on in hardware (no GPIO control), in which case bl is nil and this
// is a no-op.
func (d *Dev) BacklightOn(on bool) error {
	if d.bl == nil {
		return nil
	}
	return d.bl.Out(gpio.Level(on))
}

func (d *Dev) init() error {
	if err := d.rst.Out(gpio.Low); err != nil {
		return err
	}
	time.Sleep(20 * time.Millisecond)
	if err := d.rst.Out(gpio.High); err != nil {
		return err
	}
	time.Sleep(150 * time.Millisecond)

	invCmd := byte(cmdINVOFF)
	if d.invert {
		invCmd = cmdINVON
	}

	seq := []struct {
		cmd  byte
		data []byte
		wait time.Duration
	}{
		{cmdSWRESET, nil, 150 * time.Millisecond},
		{cmdSLPOUT, nil, 150 * time.Millisecond},
		// 16-bit colour (RGB565), 2 bytes/pixel.
		{cmdCOLMOD, []byte{colmod16bit}, 0},
		// Power control — standard ILI9486 values (as used by fbtft fb_ili9486).
		{cmdPWCTRL1, []byte{0x0E, 0x0E}, 0},
		{cmdPWCTRL2, []byte{0x41, 0x00}, 0},
		{cmdVMCTRL, []byte{0x00, 0x22, 0x80, 0x40}, 0},
		// Positive and negative gamma correction.
		{cmdPGAMMA, []byte{0x0F, 0x1F, 0x1C, 0x0C, 0x0F, 0x08, 0x48, 0x98, 0x37, 0x0A, 0x13, 0x04, 0x11, 0x0D, 0x00}, 0},
		{cmdNGAMMA, []byte{0x0F, 0x32, 0x2E, 0x0B, 0x0D, 0x05, 0x47, 0x75, 0x37, 0x06, 0x10, 0x03, 0x24, 0x20, 0x00}, 0},
		{cmdMADCTL, []byte{d.madctl}, 0},
		{invCmd, nil, 0},
		{cmdDISPON, nil, 150 * time.Millisecond},
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
		if s.wait > 0 {
			time.Sleep(s.wait)
		}
	}
	if d.bl != nil {
		return d.bl.Out(gpio.High)
	}
	return nil
}

func (d *Dev) setWindow(x0, y0, x1, y1 int) error {
	if err := d.cmd(cmdCASET); err != nil {
		return err
	}
	if err := d.data([]byte{byte(x0 >> 8), byte(x0), byte(x1 >> 8), byte(x1)}); err != nil {
		return err
	}
	if err := d.cmd(cmdRASET); err != nil {
		return err
	}
	if err := d.data([]byte{byte(y0 >> 8), byte(y0), byte(y1 >> 8), byte(y1)}); err != nil {
		return err
	}
	return d.cmd(cmdRAMWR)
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
