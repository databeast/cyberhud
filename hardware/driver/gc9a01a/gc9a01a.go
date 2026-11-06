// Package gc9a01a drives GC9A01A-based round IPS SPI TFT panels — notably the
// Adafruit 1.28" 240x240 Round TFT (product 6178) and the Waveshare RP2350/
// ESP32-S3 1.28" round LCDs, all 240x240 RGB565 over 4-wire SPI.
//
// The GC9A01A is ST7789-adjacent (shared CASET/RASET/RAMWR windowing and
// RGB565/COLMOD=0x55 pixel format) but requires its own long vendor power/gamma
// initialisation sequence and a different default MADCTL. The init table and
// MADCTL/rotation values here are ported from a known-good TinyGo driver that
// was verified on real GC9A01A hardware (the ratzbadge project's
// ratzbadgedrivers/gc9a01a), adapted to this project's periph-based SPI
// DrawTarget interface.
//
// The panel is physically round, but the framebuffer is a full 240x240 square;
// the corner pixels are addressable, they are simply hidden behind the bezel.
package gc9a01a

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
	cmdCASET  = 0x2A
	cmdRASET  = 0x2B
	cmdRAMWR  = 0x2C
	cmdMADCTL = 0x36

	// madctlDefault is the GC9A01A round-panel default (vertical scan), matching
	// the verified Waveshare/Adafruit init. It is NOT the ST7789 default.
	madctlDefault = 0x08
)

// initCmd encodes one command in the GC9A01A power-on sequence.
type initCmd struct {
	cmd   byte
	data  []byte
	delay time.Duration
}

// initSequence is the full GC9A01A power-on initialisation sequence, ported
// verbatim from the hardware-verified ratzbadge TinyGo driver (itself derived
// from the Waveshare RP2350-LCD-1.28 reference C implementation). The trailing
// commands set RGB565 pixel format (0x3A=0x55), tearing effect on (0x35),
// display inversion on (0x21), sleep out (0x11), and display on (0x29).
var initSequence = []initCmd{
	{0xEF, nil, 0},
	{0xEB, []byte{0x14}, 0},
	{0xFE, nil, 0},
	{0xEF, nil, 0},
	{0xEB, []byte{0x14}, 0},
	{0x84, []byte{0x40}, 0},
	{0x85, []byte{0xFF}, 0},
	{0x86, []byte{0xFF}, 0},
	{0x87, []byte{0xFF}, 0},
	{0x88, []byte{0x0A}, 0},
	{0x89, []byte{0x21}, 0},
	{0x8A, []byte{0x00}, 0},
	{0x8B, []byte{0x80}, 0},
	{0x8C, []byte{0x01}, 0},
	{0x8D, []byte{0x01}, 0},
	{0x8E, []byte{0xFF}, 0},
	{0x8F, []byte{0xFF}, 0},
	{0xB6, []byte{0x00, 0x20}, 0},
	{cmdMADCTL, []byte{madctlDefault}, 0},
	{0x3A, []byte{0x55}, 0}, // RGB565, 16-bit
	{0x90, []byte{0x08, 0x08, 0x08, 0x08}, 0},
	{0xBD, []byte{0x06}, 0},
	{0xBC, []byte{0x00}, 0},
	{0xFF, []byte{0x60, 0x01, 0x04}, 0},
	{0xC3, []byte{0x13}, 0},
	{0xC4, []byte{0x13}, 0},
	{0xC9, []byte{0x22}, 0},
	{0xBE, []byte{0x11}, 0},
	{0xE1, []byte{0x10, 0x0E}, 0},
	{0xDF, []byte{0x21, 0x0C, 0x02}, 0},
	{0xF0, []byte{0x45, 0x09, 0x08, 0x08, 0x26, 0x2A}, 0},
	{0xF1, []byte{0x43, 0x70, 0x72, 0x36, 0x37, 0x6F}, 0},
	{0xF2, []byte{0x45, 0x09, 0x08, 0x08, 0x26, 0x2A}, 0},
	{0xF3, []byte{0x43, 0x70, 0x72, 0x36, 0x37, 0x6F}, 0},
	{0xED, []byte{0x1B, 0x0B}, 0},
	{0xAE, []byte{0x77}, 0},
	{0xCD, []byte{0x63}, 0},
	{0x70, []byte{0x07, 0x07, 0x04, 0x0E, 0x0F, 0x09, 0x07, 0x08, 0x03}, 0},
	{0xE8, []byte{0x34}, 0},
	{0x62, []byte{0x18, 0x0D, 0x71, 0xED, 0x70, 0x70, 0x18, 0x0F, 0x71, 0xEF, 0x70, 0x70}, 0},
	{0x63, []byte{0x18, 0x11, 0x71, 0xF1, 0x70, 0x70, 0x18, 0x13, 0x71, 0xF3, 0x70, 0x70}, 0},
	{0x64, []byte{0x28, 0x29, 0xF1, 0x01, 0xF1, 0x00, 0x07}, 0},
	{0x66, []byte{0x3C, 0x00, 0xCD, 0x67, 0x45, 0x45, 0x10, 0x00, 0x00, 0x00}, 0},
	{0x67, []byte{0x00, 0x3C, 0x00, 0x00, 0x00, 0x01, 0x54, 0x10, 0x32, 0x98}, 0},
	{0x74, []byte{0x10, 0x85, 0x80, 0x00, 0x00, 0x4E, 0x00}, 0},
	{0x98, []byte{0x3E, 0x07}, 0},
	{0x35, nil, 0},
	{0x21, nil, 0},
	{0x11, nil, 120 * time.Millisecond},
	{0x29, nil, 20 * time.Millisecond},
}

func init() {
	driver.Register(driver.Definition{
		ID:         "gc9a01a",
		Title:      "GC9A01A",
		Summary:    "SPI color LCD controller used by 240x240 1.28\" round IPS TFT panels (16-bit/pixel).",
		Monochrome: false,
		ColorDepth: driver.DepthColor,
		Capability: textlayout.CapColorFast,
		OptionDefs: []driver.OptionDefinition{
			{Key: "width", Type: "int", Summary: "Logical panel width in pixels.", Default: "240"},
			{Key: "height", Type: "int", Summary: "Logical panel height in pixels.", Default: "240"},
			{Key: "spi_hz", Type: "frequency", Summary: "SPI bus speed used for the panel transport.", Default: "40MHz"},
			{Key: "madctl", Type: "byte", Summary: "GC9A01A MADCTL register value controlling orientation and RGB/BGR order.", Default: "0x08"},
		},
		NewSPI: func(port spi.Port, dc, rst, bl gpio.PinOut, _ gpio.PinIn, cfg driver.DriverConfig) (driver.DrawTarget, error) {
			return New(port, dc, rst, bl, cfg)
		},
	})
}

// Dev is a handle to a GC9A01A round LCD.
type Dev struct {
	c      spi.Conn
	dc     gpio.PinOut
	rst    gpio.PinOut
	bl     gpio.PinOut
	width  int
	height int
	madctl byte
}

// New initialises the GC9A01A and returns a ready-to-use Dev.
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
		cfg.Height = 240
	}
	if cfg.SPIHz <= 0 {
		cfg.SPIHz = 40 * physic.MegaHertz
	}
	if cfg.MADCTL == 0 {
		cfg.MADCTL = madctlDefault
	}
	conn, err := port.Connect(cfg.SPIHz, spi.Mode0, 8)
	if err != nil {
		return nil, err
	}
	d := &Dev{c: conn, dc: dc, rst: rst, bl: bl, width: cfg.Width, height: cfg.Height, madctl: cfg.MADCTL}
	if err := d.init(); err != nil {
		return nil, err
	}
	return d, nil
}

// Bounds returns the display rectangle. The panel is round but the framebuffer
// is a full square; corner pixels are addressable but hidden by the bezel.
func (d *Dev) Bounds() image.Rectangle { return image.Rect(0, 0, d.width, d.height) }

// TextHints reports text/scroll capabilities for GC9A01A color panels.
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

// Draw renders src into the display window defined by r.
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
			// RGB565, 2 bytes/pixel, big-endian (high byte first).
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
	// Hardware reset: high 5ms → low 10ms → high 120ms (ported timing).
	if err := d.rst.Out(gpio.High); err != nil {
		return err
	}
	time.Sleep(5 * time.Millisecond)
	if err := d.rst.Out(gpio.Low); err != nil {
		return err
	}
	time.Sleep(10 * time.Millisecond)
	if err := d.rst.Out(gpio.High); err != nil {
		return err
	}
	time.Sleep(120 * time.Millisecond)

	for _, s := range initSequence {
		cmd := s.cmd
		// Honour a MADCTL override from config in place of the table default.
		if cmd == cmdMADCTL {
			if err := d.cmd(cmd); err != nil {
				return err
			}
			if err := d.data([]byte{d.madctl}); err != nil {
				return err
			}
		} else {
			if err := d.cmd(cmd); err != nil {
				return err
			}
			if len(s.data) > 0 {
				if err := d.data(s.data); err != nil {
					return err
				}
			}
		}
		if s.delay > 0 {
			time.Sleep(s.delay)
		}
	}

	if d.bl != nil {
		return d.bl.Out(gpio.High)
	}
	return nil
}

func (d *Dev) setWindow(x0, y0, x1, y1 int) error {
	// GC9A01A CASET/RASET take 16-bit coordinates (high byte then low byte).
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
