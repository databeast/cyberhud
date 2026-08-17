package ssd1306

import (
	"fmt"
	"image"
	"image/draw"

	"periph.io/x/conn/v3/i2c"

	"github.com/databeast/cyberhud/hardware/driver"
)

const (
	ssd1306Width  = 128
	ssd1306Height = 64

	DefaultAddr uint16 = 0x3C

	cmdDisplayOff       byte = 0xAE
	cmdDisplayOn        byte = 0xAF
	cmdSetDisplayClock  byte = 0xD5
	cmdSetMultiplex     byte = 0xA8
	cmdSetDisplayOffset byte = 0xD3
	cmdSetStartLine     byte = 0x40
	cmdChargePump       byte = 0x8D
	cmdMemoryMode       byte = 0x20
	cmdSegRemap         byte = 0xA1
	cmdCOMScanDec       byte = 0xC8
	cmdSetCOMPins       byte = 0xDA
	cmdSetContrast      byte = 0x81
	cmdSetPrecharge     byte = 0xD9
	cmdSetVCOMDetect    byte = 0xDB
	cmdEntireDisplayRAM byte = 0xA4
	cmdNormalDisplay    byte = 0xA6
	cmdColumnAddr       byte = 0x21
	cmdPageAddr         byte = 0x22
)

func init() {
	driver.Register(driver.Definition{
		ID:         "ssd1306",
		Title:      "SSD1306",
		Summary:    "I2C monochrome OLED controller (128x64 or 128x32).",
		Monochrome: true,
		OptionDefs: []driver.OptionDefinition{
			{Key: "width", Type: "int", Summary: "Display width in pixels.", Default: "128"},
			{Key: "height", Type: "int", Summary: "Display height in pixels.", Default: "64"},
			{Key: "address", Type: "byte", Summary: "I2C device address.", Default: "0x3C"},
		},
		NewI2C: func(bus i2c.Bus, cfg driver.DriverConfig) (driver.DrawTarget, error) {
			return New(bus, cfg)
		},
	})
}

// SSD1306 drives an SSD1306 OLED display over I2C.
type SSD1306 struct {
	dev *i2c.Dev
	cfg driver.DriverConfig
}

// New opens the SSD1306 at the configured I2C address and runs initialization.
func New(bus i2c.Bus, cfg driver.DriverConfig) (*SSD1306, error) {
	if cfg.Width <= 0 {
		cfg.Width = ssd1306Width
	}
	if cfg.Height <= 0 {
		cfg.Height = ssd1306Height
	}
	addr := DefaultAddr
	if cfg.I2CAddr != 0 {
		addr = cfg.I2CAddr
	}

	dev := &i2c.Dev{Bus: bus, Addr: addr}
	d := &SSD1306{dev: dev, cfg: cfg}

	if err := d.init(); err != nil {
		return nil, fmt.Errorf("ssd1306 [0x%02X]: %w", addr, err)
	}
	return d, nil
}

func (d *SSD1306) Bounds() image.Rectangle {
	return image.Rect(0, 0, d.cfg.Width, d.cfg.Height)
}

// DrawImage converts pixels to the SSD1306 page-column format and writes to GDDRAM.
// Each byte represents 8 vertical pixels (LSB = top). Row-major page order.
func (d *SSD1306) DrawImage(src draw.Image) error {
	w := d.cfg.Width
	h := d.cfg.Height
	pages := h / 8

	// Set column and page address window to cover entire display.
	if err := d.cmdSeq(cmdColumnAddr, 0x00, byte(w-1)); err != nil {
		return err
	}
	if err := d.cmdSeq(cmdPageAddr, 0x00, byte(pages-1)); err != nil {
		return err
	}

	buf := make([]byte, w*pages)
	for page := 0; page < pages; page++ {
		for col := 0; col < w; col++ {
			var b byte
			for bit := 0; bit < 8; bit++ {
				y := page*8 + bit
				r, g, bl, _ := src.At(col, y).RGBA()
				if luminanceOn(r, g, bl) {
					b |= 1 << uint(bit)
				}
			}
			buf[page*w+col] = b
		}
	}

	return d.data(buf)
}

func (d *SSD1306) init() error {
	comPins := byte(0x12) // 128x64 default
	if d.cfg.Height == 32 {
		comPins = 0x02
	}

	cmds := []byte{
		cmdDisplayOff,
		cmdSetDisplayClock, 0x80,
		cmdSetMultiplex, byte(d.cfg.Height - 1),
		cmdSetDisplayOffset, 0x00,
		cmdSetStartLine,
		cmdChargePump, 0x14, // enable charge pump
		cmdMemoryMode, 0x00, // horizontal addressing mode
		cmdSegRemap,
		cmdCOMScanDec,
		cmdSetCOMPins, comPins,
		cmdSetContrast, 0xCF,
		cmdSetPrecharge, 0xF1,
		cmdSetVCOMDetect, 0x40,
		cmdEntireDisplayRAM,
		cmdNormalDisplay,
		cmdDisplayOn,
	}

	for _, c := range cmds {
		if err := d.cmd(c); err != nil {
			return err
		}
	}
	return nil
}

// cmd sends a single command byte (Co=1, D/C#=0 → control byte 0x00).
func (d *SSD1306) cmd(c byte) error {
	return d.dev.Tx([]byte{0x00, c}, nil)
}

// cmdSeq sends a command followed by parameter bytes.
func (d *SSD1306) cmdSeq(cmd byte, params ...byte) error {
	buf := make([]byte, 0, 2+len(params)*2)
	buf = append(buf, 0x00, cmd)
	for _, p := range params {
		buf = append(buf, 0x00, p)
	}
	return d.dev.Tx(buf, nil)
}

// data sends pixel data (Co=0, D/C#=1 → control byte 0x40).
func (d *SSD1306) data(buf []byte) error {
	out := make([]byte, 1+len(buf))
	out[0] = 0x40
	copy(out[1:], buf)
	return d.dev.Tx(out, nil)
}

func luminanceOn(r, g, b uint32) bool {
	return (299*r+587*g+114*b)/1000 > 0x3000
}
