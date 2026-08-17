package st7735s

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
	cmdINVOFF  = 0x20
	cmdINVON   = 0x21
	cmdDISPON  = 0x29
	cmdCASET   = 0x2A
	cmdRASET   = 0x2B
	cmdRAMWR   = 0x2C
	cmdMADCTL  = 0x36
	cmdCOLMOD  = 0x3A
	cmdFRMCTR1 = 0xB1
	cmdFRMCTR2 = 0xB2
	cmdFRMCTR3 = 0xB3
	cmdINVCTR  = 0xB4
	cmdPWCTR1  = 0xC0
	cmdPWCTR2  = 0xC1
	cmdPWCTR3  = 0xC2
	cmdPWCTR4  = 0xC3
	cmdPWCTR5  = 0xC4
	cmdVMCTR1  = 0xC5
	cmdGMCTRP1 = 0xE0
	cmdGMCTRN1 = 0xE1
)

func init() {
	driver.Register(driver.Definition{
		ID:         "st7735s",
		Title:      "ST7735S",
		Summary:    "SPI color LCD controller used by compact TFT panels such as 128x128 and 160x80 displays.",
		Monochrome: false,
		OptionDefs: []driver.OptionDefinition{
			{Key: "width", Type: "int", Summary: "Logical panel width in pixels.", Default: "128"},
			{Key: "height", Type: "int", Summary: "Logical panel height in pixels.", Default: "128"},
			{Key: "spi_hz", Type: "frequency", Summary: "SPI bus speed used for the panel transport.", Default: "40MHz"},
			{Key: "madctl", Type: "byte", Summary: "ST7735S MADCTL register value controlling orientation and RGB/BGR order.", Default: "0xA8"},
			{Key: "x_offset", Type: "int", Summary: "Physical framebuffer X offset applied to window commands.", Default: "2"},
			{Key: "y_offset", Type: "int", Summary: "Physical framebuffer Y offset applied to window commands.", Default: "1"},
		},
		DefaultText: textlayout.DefaultTextHints(image.Rect(0, 0, 128, 128)),
		NewSPI: func(port spi.Port, dc, rst, bl gpio.PinOut, _ gpio.PinIn, cfg driver.DriverConfig) (driver.DrawTarget, error) {
			return New(port, dc, rst, bl, cfg)
		},
	})
}

type ST7735S struct {
	cfg driver.DriverConfig
	c   spi.Conn
	dc  gpio.PinOut
	rst gpio.PinOut
	bl  gpio.PinOut
}

// invCmd returns the inversion command byte based on the InvertColors config flag.
func (d *ST7735S) invCmd() byte {
	if d.cfg.InvertColors {
		return cmdINVON
	}
	return cmdINVOFF
}

func New(port spi.Port, dc, rst, bl gpio.PinOut, cfg driver.DriverConfig) (*ST7735S, error) {
	if port == nil {
		return nil, errors.New("display: spi port must not be nil")
	}
	if dc == nil || rst == nil {
		return nil, errors.New("display: dc and rst pins must not be nil")
	}
	if cfg.Width <= 0 {
		cfg.Width = 128
	}
	if cfg.Height <= 0 {
		cfg.Height = 128
	}
	if cfg.SPIHz <= 0 {
		cfg.SPIHz = 40 * physic.MegaHertz
	}
	if cfg.MADCTL == 0 {
		cfg.MADCTL = driver.MadctlMY | driver.MadctlMV | driver.MadctlBGR
	}
	if cfg.XOffset == 0 {
		cfg.XOffset = 2
	}
	if cfg.YOffset == 0 {
		cfg.YOffset = 1
	}
	conn, err := port.Connect(cfg.SPIHz, spi.Mode0, 8)
	if err != nil {
		return nil, err
	}
	d := &ST7735S{cfg: cfg, c: conn, dc: dc, rst: rst, bl: bl}
	if err := d.init(); err != nil {
		return nil, err
	}
	return d, nil
}

func (d *ST7735S) Bounds() image.Rectangle { return image.Rect(0, 0, d.cfg.Width, d.cfg.Height) }

func (d *ST7735S) TextHints() textlayout.TextHints {
	b := d.Bounds()
	return textlayout.TextHints{
		PixelWidth: b.Dx(), PixelHeight: b.Dy(),
		SupportsVerticalScroll: true, SupportsHorizontalScroll: true, SupportsAutoScroll: true,
		PreferEventRefresh:     false,
		Capability:             textlayout.CapColorFast,
		DefaultTickerDirection: textlayout.TickerDirectionVertical, DefaultLineMode: textlayout.LineModeTruncate,
	}
}

func (d *ST7735S) DrawImage(src draw.Image) error {
	return d.Draw(d.Bounds(), src, image.Point{})
}

func (d *ST7735S) Draw(r image.Rectangle, src image.Image, sp image.Point) error {
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
			r16, g16, b16, _ := src.At(x+sp.X-r.Min.X, y+sp.Y-r.Min.Y).RGBA()
			rgb := (r16>>11)<<11 | (g16>>10)<<5 | (b16 >> 11)
			buf[i], buf[i+1] = byte(rgb>>8), byte(rgb)
			i += 2
		}
	}
	return d.data(buf)
}

func (d *ST7735S) Fill(c color.Color) error {
	return d.Draw(d.Bounds(), image.NewUniform(c), image.Point{})
}

func (d *ST7735S) BacklightOn(on bool) error {
	if d.bl == nil {
		return nil
	}
	return d.bl.Out(gpio.Level(on))
}

func (d *ST7735S) init() error {
	if err := d.rst.Out(gpio.Low); err != nil {
		return err
	}
	time.Sleep(10 * time.Millisecond)
	if err := d.rst.Out(gpio.High); err != nil {
		return err
	}
	time.Sleep(120 * time.Millisecond)

	type entry struct {
		cmd  byte
		data []byte
		wait time.Duration
	}
	x0 := byte(d.cfg.XOffset)
	x1 := byte(d.cfg.XOffset + d.cfg.Width - 1)
	y0 := byte(d.cfg.YOffset)
	y1 := byte(d.cfg.YOffset + d.cfg.Height - 1)
	seq := []entry{
		{cmdSWRESET, nil, 150 * time.Millisecond},
		{cmdSLPOUT, nil, 500 * time.Millisecond},
		{cmdFRMCTR1, []byte{0x01, 0x2C, 0x2D}, 0},
		{cmdFRMCTR2, []byte{0x01, 0x2C, 0x2D}, 0},
		{cmdFRMCTR3, []byte{0x01, 0x2C, 0x2D, 0x01, 0x2C, 0x2D}, 0},
		{cmdINVCTR, []byte{0x07}, 0},
		{cmdPWCTR1, []byte{0xA2, 0x02, 0x84}, 0},
		{cmdPWCTR2, []byte{0xC5}, 0},
		{cmdPWCTR3, []byte{0x0A, 0x00}, 0},
		{cmdPWCTR4, []byte{0x8A, 0x2A}, 0},
		{cmdPWCTR5, []byte{0x8A, 0xEE}, 0},
		{cmdVMCTR1, []byte{0x0E}, 0},
		{d.invCmd(), nil, 0},
		{cmdMADCTL, []byte{d.cfg.MADCTL}, 0},
		{cmdCOLMOD, []byte{0x05}, 0},
		{cmdCASET, []byte{0x00, x0, 0x00, x1}, 0},
		{cmdRASET, []byte{0x00, y0, 0x00, y1}, 0},
		{cmdGMCTRP1, []byte{0x02, 0x1C, 0x07, 0x12, 0x37, 0x32, 0x29, 0x2D, 0x29, 0x25, 0x2B, 0x39, 0x00, 0x01, 0x03, 0x10}, 0},
		{cmdGMCTRN1, []byte{0x03, 0x1D, 0x07, 0x06, 0x2E, 0x2C, 0x29, 0x2D, 0x2E, 0x2E, 0x37, 0x3F, 0x00, 0x00, 0x02, 0x10}, 0},
		{cmdNORON, nil, 10 * time.Millisecond},
		{cmdDISPON, nil, 100 * time.Millisecond},
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

func (d *ST7735S) setWindow(x0, y0, x1, y1 int) error {
	x0 += d.cfg.XOffset
	x1 += d.cfg.XOffset
	y0 += d.cfg.YOffset
	y1 += d.cfg.YOffset
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

func (d *ST7735S) cmd(c byte) error {
	if err := d.dc.Out(gpio.Low); err != nil {
		return err
	}
	return d.c.Tx([]byte{c}, nil)
}

func (d *ST7735S) data(buf []byte) error {
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
