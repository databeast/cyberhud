package ht16k33

import (
	"fmt"
	"image"
	"image/draw"

	"periph.io/x/conn/v3/i2c"

	"github.com/databeast/cyberhud/hardware/driver"
)

const (
	ht16k33Width  = 16
	ht16k33Height = 8

	DefaultAddr uint16 = 0x70

	// HT16K33 command registers.
	cmdSystemSetup  byte = 0x20
	cmdOscillatorOn byte = 0x01
	cmdDisplaySetup byte = 0x80
	cmdDisplayOn    byte = 0x01
	cmdBrightness   byte = 0xE0
	cmdRAMAddr      byte = 0x00
)

func init() {
	driver.Register(driver.Definition{
		ID:         "ht16k33",
		Title:      "HT16K33",
		Summary:    "I2C LED matrix/segment controller (up to 16x8 LEDs).",
		Monochrome: true,
		OptionDefs: []driver.OptionDefinition{
			{Key: "width", Type: "int", Summary: "Display width in pixels/segments.", Default: "16"},
			{Key: "height", Type: "int", Summary: "Display height in pixels/rows.", Default: "8"},
			{Key: "address", Type: "byte", Summary: "I2C device address.", Default: "0x70"},
			{Key: "brightness", Type: "int", Summary: "LED brightness (0-15).", Default: "15"},
		},
		NewI2C: func(bus i2c.Bus, cfg driver.DriverConfig) (driver.DrawTarget, error) {
			return New(bus, cfg)
		},
	})
}

// HT16K33 drives an HT16K33 LED matrix controller over I2C.
type HT16K33 struct {
	dev *i2c.Dev
	cfg driver.DriverConfig
}

// New opens the HT16K33 at the configured I2C address and initializes the display.
func New(bus i2c.Bus, cfg driver.DriverConfig) (*HT16K33, error) {
	if cfg.Width <= 0 {
		cfg.Width = ht16k33Width
	}
	if cfg.Height <= 0 {
		cfg.Height = ht16k33Height
	}
	addr := DefaultAddr
	if cfg.I2CAddr != 0 {
		addr = cfg.I2CAddr
	}

	dev := &i2c.Dev{Bus: bus, Addr: addr}
	d := &HT16K33{dev: dev, cfg: cfg}

	if err := d.init(); err != nil {
		return nil, fmt.Errorf("ht16k33 [0x%02X]: %w", addr, err)
	}
	return d, nil
}

func (d *HT16K33) Bounds() image.Rectangle {
	return image.Rect(0, 0, d.cfg.Width, d.cfg.Height)
}

// DrawImage packs image pixels into the HT16K33 display RAM format.
// Each row is stored as a 16-bit word (little-endian) in 2 bytes.
// A pixel is "on" if its luminance exceeds a threshold.
func (d *HT16K33) DrawImage(src draw.Image) error {
	w := d.cfg.Width
	h := d.cfg.Height

	// HT16K33 has 16 bytes of display RAM (8 rows × 16-bit).
	buf := make([]byte, 16)
	for row := 0; row < h && row < 8; row++ {
		var word uint16
		for col := 0; col < w && col < 16; col++ {
			r, g, b, _ := src.At(col, row).RGBA()
			if luminanceOn(r, g, b) {
				word |= 1 << uint(col)
			}
		}
		buf[row*2] = byte(word & 0xFF)
		buf[row*2+1] = byte(word >> 8)
	}

	// Write display RAM starting at address 0x00.
	out := make([]byte, 1+len(buf))
	out[0] = cmdRAMAddr
	copy(out[1:], buf)
	return d.dev.Tx(out, nil)
}

func (d *HT16K33) init() error {
	// Turn on oscillator.
	if err := d.dev.Tx([]byte{cmdSystemSetup | cmdOscillatorOn}, nil); err != nil {
		return fmt.Errorf("enable oscillator: %w", err)
	}

	// Turn on display, no blinking.
	if err := d.dev.Tx([]byte{cmdDisplaySetup | cmdDisplayOn}, nil); err != nil {
		return fmt.Errorf("display on: %w", err)
	}

	// Set brightness to max (or configured value).
	brightness := byte(15)
	if err := d.dev.Tx([]byte{cmdBrightness | (brightness & 0x0F)}, nil); err != nil {
		return fmt.Errorf("set brightness: %w", err)
	}

	// Clear display RAM.
	clear := make([]byte, 17) // 1 byte addr + 16 bytes data (zeros)
	clear[0] = cmdRAMAddr
	if err := d.dev.Tx(clear, nil); err != nil {
		return fmt.Errorf("clear RAM: %w", err)
	}

	return nil
}

func luminanceOn(r, g, b uint32) bool {
	return (299*r+587*g+114*b)/1000 > 0x3000
}
