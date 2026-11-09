// Package ili9341 drives ILI9341-based SPI colour TFT panels — notably the
// Adafruit 2.4" / 2.8" / 3.2" TFT breakouts and the many generic 240x320 SPI
// modules built around this controller.
//
// Pixel format. The ILI9341 is driven here in 16-bit colour (RGB565), two bytes
// per pixel (COLMOD=0x55), the same transport ST7789 uses. The CASET/RASET/RAMWR
// windowing command set is shared with ST7789/ILI9486; what differs is the
// power and gamma init sequence, which uses the ILI9341's own register map
// (0xC0/0xC1/0xC5/0xC7 power, 0xE0/0xE1 gamma) — hence a dedicated driver rather
// than reusing st7789.
//
// Native orientation is 240x320 portrait. MADCTL defaults to 0x48 (MX set, BGR
// order), the common Adafruit configuration; orientation and RGB/BGR order vary
// by module and may need on-device adjustment via the -display-madctl flag.
package ili9341

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

	// ILI9341-specific power/gamma registers used by the init sequence.
	cmdPWCTRL1  = 0xC0 // Power control 1
	cmdPWCTRL2  = 0xC1 // Power control 2
	cmdVMCTRL1  = 0xC5 // VCOM control 1
	cmdVMCTRL2  = 0xC7 // VCOM control 2
	cmdFRMCTR1  = 0xB1 // Frame rate control (normal mode)
	cmdDFUNCTR  = 0xB6 // Display function control
	cmdPGAMMA   = 0xE0 // Positive gamma correction
	cmdNGAMMA   = 0xE1 // Negative gamma correction
	cmdGAMMASET = 0x26 // Gamma set

	// COLMOD value selecting 16-bit/pixel (RGB565), 2 bytes per pixel over SPI.
	colmod16bit = 0x55
)

func init() {
	driver.Register(driver.Definition{
		ID:         "ili9341",
		Title:      "ILI9341",
		Summary:    "SPI color LCD controller used by 240x320 2.4\"/2.8\"/3.2\" TFT panels (16-bit/pixel).",
		Monochrome: false,
		ColorDepth: driver.DepthColor,
		Capability: textlayout.CapColorFast,
		OptionDefs: []driver.OptionDefinition{
			{Key: "width", Type: "int", Summary: "Logical panel width in pixels.", Default: "240"},
			{Key: "height", Type: "int", Summary: "Logical panel height in pixels.", Default: "320"},
			{Key: "spi_hz", Type: "frequency", Summary: "SPI bus speed used for the panel transport.", Default: "32MHz"},
			{Key: "madctl", Type: "byte", Summary: "ILI9341 MADCTL register value controlling orientation and RGB/BGR order.", Default: "0x48"},
			{Key: "x_offset", Type: "int", Summary: "Physical framebuffer X offset applied to window commands.", Default: "0"},
			{Key: "y_offset", Type: "int", Summary: "Physical framebuffer Y offset applied to window commands.", Default: "0"},
		},
		NewSPI: func(port spi.Port, dc, rst, bl gpio.PinOut, _ gpio.PinIn, cfg driver.DriverConfig) (driver.DrawTarget, error) {
			return New(port, dc, rst, bl, cfg)
		},
	})
}

// Dev is a handle to an ILI9341 LCD.
type Dev struct {
	c      spi.Conn
	dc     gpio.PinOut
	rst    gpio.PinOut
	bl     gpio.PinOut
	width  int
	height int
	xOff   int
	yOff   int
	madctl byte
	invert bool
}

// New initialises the ILI9341 and returns a ready-to-use Dev.
func New(port spi.Port, dc, rst, bl gpio.PinOut, cfg driver.DriverConfig) (*Dev, error) {
	if port == nil {
		return nil, errors.New("display: spi port must not be nil")
	}
	if dc == nil || rst == nil {
		return nil, errors.New("display: dc and rst pins must not be nil")
	}
	if cfg.Width <= 0 {
		cfg.Width = 240
	}
	if cfg.Height <= 0 {
		cfg.Height = 320
	}
	if cfg.SPIHz <= 0 {
		// The ILI9341 happily runs at 32MHz over SPI on a Pi.
		cfg.SPIHz = 32 * physic.MegaHertz
	}
	if cfg.MADCTL == 0 {
		// 0x48: MX set with BGR order — the common Adafruit portrait setup.
		cfg.MADCTL = driver.MadctlMX | driver.MadctlBGR
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
		xOff:   cfg.XOffset,
		yOff:   cfg.YOffset,
		invert: cfg.InvertColors,
	}
	if err := d.init(); err != nil {
		return nil, err
	}
	return d, nil
}

// Bounds returns the display rectangle.
func (d *Dev) Bounds() image.Rectangle { return image.Rect(0, 0, d.width, d.height) }

// TextHints reports text/scroll capabilities for ILI9341 color panels.
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
// 16-bit RGB565, two bytes per pixel.
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

// BacklightOn sets the backlight to on or off. When the board ties the
// backlight on in hardware (no GPIO control), bl is nil and this is a no-op.
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
		// Power/VCOM control — standard ILI9341 values (as used by the Adafruit
		// and fbtft fb_ili9341 drivers).
		{cmdPWCTRL1, []byte{0x23}, 0},
		{cmdPWCTRL2, []byte{0x10}, 0},
		{cmdVMCTRL1, []byte{0x3E, 0x28}, 0},
		{cmdVMCTRL2, []byte{0x86}, 0},
		// Frame rate ~79Hz, display function control.
		{cmdFRMCTR1, []byte{0x00, 0x18}, 0},
		{cmdDFUNCTR, []byte{0x08, 0x82, 0x27}, 0},
		{cmdGAMMASET, []byte{0x01}, 0},
		// Positive and negative gamma correction.
		{cmdPGAMMA, []byte{0x0F, 0x31, 0x2B, 0x0C, 0x0E, 0x08, 0x4E, 0xF1, 0x37, 0x07, 0x10, 0x03, 0x0E, 0x09, 0x00}, 0},
		{cmdNGAMMA, []byte{0x00, 0x0E, 0x14, 0x03, 0x11, 0x07, 0x31, 0xC1, 0x48, 0x08, 0x0F, 0x0C, 0x31, 0x36, 0x0F}, 0},
		// 16-bit colour (RGB565), 2 bytes/pixel.
		{cmdCOLMOD, []byte{colmod16bit}, 0},
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
	x0 += d.xOff
	x1 += d.xOff
	y0 += d.yOff
	y1 += d.yOff
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
