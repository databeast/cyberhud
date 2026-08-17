package st7789

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
	cmdNORON   = 0x13
	cmdINVON   = 0x21
	cmdDISPON  = 0x29
	cmdCASET   = 0x2A
	cmdRASET   = 0x2B
	cmdRAMWR   = 0x2C
	cmdMADCTL  = 0x36
	cmdCOLMOD  = 0x3A
)

func init() {
	driver.Register(driver.Definition{
		ID:         "st7789",
		Title:      "ST7789",
		Summary:    "SPI color LCD controller used by 240x240, 240x135, and 320x240 TFT panels.",
		Monochrome: false,
		OptionDefs: []driver.OptionDefinition{
			{Key: "width", Type: "int", Summary: "Logical panel width in pixels.", Default: "240"},
			{Key: "height", Type: "int", Summary: "Logical panel height in pixels.", Default: "240"},
			{Key: "spi_hz", Type: "frequency", Summary: "SPI bus speed used for the panel transport.", Default: "40MHz"},
			{Key: "madctl", Type: "byte", Summary: "ST7789 MADCTL register value controlling orientation and RGB/BGR order.", Default: "0x40"},
			{Key: "x_offset", Type: "int", Summary: "Physical framebuffer X offset applied to window commands.", Default: "0"},
			{Key: "y_offset", Type: "int", Summary: "Physical framebuffer Y offset applied to window commands.", Default: "0"},
		},
		DefaultText: textlayout.DefaultTextHints(image.Rect(0, 0, 240, 240)),
		NewSPI: func(port spi.Port, dc, rst, bl gpio.PinOut, _ gpio.PinIn, cfg driver.DriverConfig) (driver.DrawTarget, error) {
			return New(port, dc, rst, bl, cfg)
		},
	})
}

// Dev is a handle to an ST7789 LCD.
type Dev struct {
	c      spi.Conn
	dc     gpio.PinOut
	rst    gpio.PinOut
	bl     gpio.PinOut
	width  int
	height int
	madctl byte
	xOff   int
	yOff   int
}

// New initialises the ST7789 and returns a ready-to-use Dev.
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
		cfg.MADCTL = driver.MadctlMX
	}
	conn, err := port.Connect(cfg.SPIHz, spi.Mode0, 8)
	if err != nil {
		return nil, err
	}
	d := &Dev{c: conn, dc: dc, rst: rst, bl: bl, width: cfg.Width, height: cfg.Height, madctl: cfg.MADCTL, xOff: cfg.XOffset, yOff: cfg.YOffset}
	if err := d.init(); err != nil {
		return nil, err
	}
	return d, nil
}

// Bounds returns the display rectangle.
func (d *Dev) Bounds() image.Rectangle { return image.Rect(0, 0, d.width, d.height) }

// TextHints reports text/scroll capabilities for ST7789 color panels.
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

// BacklightOn sets the backlight to on or off.
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
	time.Sleep(10 * time.Millisecond)
	if err := d.rst.Out(gpio.High); err != nil {
		return err
	}
	time.Sleep(120 * time.Millisecond)
	seq := []struct {
		cmd  byte
		data []byte
		wait time.Duration
	}{
		{cmdSWRESET, nil, 150 * time.Millisecond},
		{cmdSLPOUT, nil, 500 * time.Millisecond},
		{cmdCOLMOD, []byte{0x55}, 0},
		{cmdMADCTL, []byte{d.madctl}, 0},
		{cmdINVON, nil, 0},
		{cmdNORON, nil, 10 * time.Millisecond},
		{cmdDISPON, nil, 10 * time.Millisecond},
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
