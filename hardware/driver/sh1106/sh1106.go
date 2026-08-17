package sh1106

import (
	"errors"
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
	sh1106Width  = 128
	sh1106Height = 64
)

func init() {
	driver.Register(driver.Definition{
		ID:         "sh1106",
		Title:      "SH1106",
		Summary:    "SPI monochrome OLED controller commonly used by 128x64 HAT panels.",
		Monochrome: true,
		OptionDefs: []driver.OptionDefinition{
			{Key: "width", Type: "int", Summary: "Logical panel width in pixels.", Default: "128"},
			{Key: "height", Type: "int", Summary: "Logical panel height in pixels.", Default: "64"},
			{Key: "spi_hz", Type: "frequency", Summary: "SPI bus speed used for the panel transport.", Default: "10MHz"},
			{Key: "col_offset", Type: "int", Summary: "Column offset applied when addressing SH1106 display RAM.", Default: "2"},
		},
		DefaultText: textlayout.DefaultTextHints(image.Rect(0, 0, sh1106Width, sh1106Height)),
		NewSPI: func(port spi.Port, dc, rst, bl gpio.PinOut, _ gpio.PinIn, cfg driver.DriverConfig) (driver.DrawTarget, error) {
			return New(port, dc, rst, bl, cfg)
		},
	})
}

type SH1106 struct {
	cfg driver.DriverConfig
	c   spi.Conn
	dc  gpio.PinOut
	rst gpio.PinOut
	bl  gpio.PinOut
}

func New(port spi.Port, dc, rst, bl gpio.PinOut, cfg driver.DriverConfig) (*SH1106, error) {
	if port == nil {
		return nil, errors.New("display: spi port must not be nil")
	}
	if dc == nil || rst == nil {
		return nil, errors.New("display: dc and rst pins must not be nil")
	}
	if cfg.Width <= 0 {
		cfg.Width = sh1106Width
	}
	if cfg.Height <= 0 {
		cfg.Height = sh1106Height
	}
	if cfg.SPIHz <= 0 {
		cfg.SPIHz = 10 * physic.MegaHertz
	}
	if cfg.ColOffset == 0 {
		cfg.ColOffset = 2
	}
	conn, err := port.Connect(cfg.SPIHz, spi.Mode0, 8)
	if err != nil {
		return nil, err
	}
	d := &SH1106{cfg: cfg, c: conn, dc: dc, rst: rst, bl: bl}
	if err := d.init(); err != nil {
		return nil, err
	}
	return d, nil
}

func (d *SH1106) Bounds() image.Rectangle { return image.Rect(0, 0, d.cfg.Width, d.cfg.Height) }

// TextHints reports what this panel knows about itself: its pixel dimensions, its
// colour/refresh capability and its scrolling behaviour.
//
// It deliberately does not report glyph metrics. A panel has no font — font choice
// depends on region size and which faces are registered, so it belongs to the
// Region, which fills those fields from its tier catalog (see
// region.applyBaselineGlyphMetrics). This driver previously asserted the
// textlayout 5x8/6/10 constants, which described a font it had no knowledge of.
//
// PPI is likewise left unset here because this driver serves several physical panel
// sizes. Supplying a real value where it is known (per panel product or per screen)
// is what makes the tier system's millimetre targets physically meaningful; see
// region.AssumedPPI.
func (d *SH1106) TextHints() textlayout.TextHints {
	b := d.Bounds()
	return textlayout.TextHints{
		PixelWidth: b.Dx(), PixelHeight: b.Dy(),
		SupportsVerticalScroll: true, SupportsHorizontalScroll: true, SupportsAutoScroll: true,
		PreferEventRefresh:     false,
		Capability:             textlayout.CapMonoFast,
		DefaultTickerDirection: textlayout.TickerDirectionVertical, DefaultLineMode: textlayout.LineModeTruncate,
	}
}

func (d *SH1106) DrawImage(src draw.Image) error {
	bounds := d.Bounds()
	pages := bounds.Dy() / 8
	buf := make([]byte, bounds.Dx())
	for page := 0; page < pages; page++ {
		for x := 0; x < bounds.Dx(); x++ {
			var col byte
			for bit := 0; bit < 8; bit++ {
				y := page*8 + bit
				r, g, b, _ := src.At(x, y).RGBA()
				if luminanceOn(r, g, b) {
					col |= 1 << bit
				}
			}
			buf[x] = col
		}
		if err := d.setPage(page); err != nil {
			return err
		}
		if err := d.data(buf); err != nil {
			return err
		}
	}
	return nil
}

func (d *SH1106) init() error {
	if err := d.rst.Out(gpio.Low); err != nil {
		return err
	}
	time.Sleep(10 * time.Millisecond)
	if err := d.rst.Out(gpio.High); err != nil {
		return err
	}
	time.Sleep(10 * time.Millisecond)
	initCmds := [][]byte{{0xAE}, {0xD5, 0x80}, {0xA8, 0x3F}, {0xD3, 0x00}, {0x40}, {0xAD, 0x8B}, {0xA1}, {0xC8}, {0xDA, 0x12}, {0x81, 0xCF}, {0xD9, 0xF1}, {0xDB, 0x40}, {0xA4}, {0xA6}}
	for _, c := range initCmds {
		if err := d.cmd(c[0]); err != nil {
			return err
		}
		for _, p := range c[1:] {
			if err := d.cmd(p); err != nil {
				return err
			}
		}
	}
	if err := d.clearGRAM(); err != nil {
		return err
	}
	if err := d.cmd(0xAF); err != nil {
		return err
	}
	if d.bl != nil {
		return d.bl.Out(gpio.High)
	}
	return nil
}

func (d *SH1106) clearGRAM() error {
	pages := d.cfg.Height / 8
	buf := make([]byte, d.cfg.Width)
	for page := 0; page < pages; page++ {
		if err := d.setPage(page); err != nil {
			return err
		}
		if err := d.data(buf); err != nil {
			return err
		}
	}
	return nil
}

func (d *SH1106) setPage(page int) error {
	col := d.cfg.ColOffset
	if err := d.cmd(byte(0xB0 | (page & 0x0F))); err != nil {
		return err
	}
	if err := d.cmd(byte(0x00 | (col & 0x0F))); err != nil {
		return err
	}
	return d.cmd(byte(0x10 | ((col >> 4) & 0x0F)))
}

func (d *SH1106) cmd(c byte) error {
	if err := d.dc.Out(gpio.Low); err != nil {
		return err
	}
	return d.c.Tx([]byte{c}, nil)
}

func (d *SH1106) data(buf []byte) error {
	if err := d.dc.Out(gpio.High); err != nil {
		return err
	}
	return d.c.Tx(buf, nil)
}

func luminanceOn(r, g, b uint32) bool { return (299*r+587*g+114*b)/1000 > 0x3000 }
