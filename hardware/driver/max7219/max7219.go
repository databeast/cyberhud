package max7219

import (
	"errors"
	"image"
	"image/draw"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/conn/v3/spi"

	"github.com/databeast/cyberhud/hardware/driver"
)

const (
	max7219Width  = 8
	max7219Height = 8

	regNoop        byte = 0x00
	regDigit0      byte = 0x01
	regDecodeMode  byte = 0x09
	regIntensity   byte = 0x0A
	regScanLimit   byte = 0x0B
	regShutdown    byte = 0x0C
	regDisplayTest byte = 0x0F
)

func init() {
	driver.Register(driver.Definition{
		ID:         "max7219",
		Title:      "MAX7219",
		Summary:    "SPI LED matrix controller (8x8 per device, chainable).",
		Monochrome: true,
		OptionDefs: []driver.OptionDefinition{
			{Key: "width", Type: "int", Summary: "Display width in pixels (multiple of 8 for chained devices).", Default: "8"},
			{Key: "height", Type: "int", Summary: "Display height in pixels.", Default: "8"},
			{Key: "spi_hz", Type: "frequency", Summary: "SPI bus speed.", Default: "10MHz"},
			{Key: "devices", Type: "int", Summary: "Number of cascaded MAX7219 devices.", Default: "1"},
		},
		NewSPI: func(port spi.Port, dc, rst, bl gpio.PinOut, _ gpio.PinIn, cfg driver.DriverConfig) (driver.DrawTarget, error) {
			return New(port, cfg)
		},
	})
}

// MAX7219 drives one or more cascaded MAX7219 LED matrix controllers over SPI.
type MAX7219 struct {
	cfg     driver.DriverConfig
	c       spi.Conn
	devices int
}

// New creates and initializes a MAX7219 chain.
func New(port spi.Port, cfg driver.DriverConfig) (*MAX7219, error) {
	if port == nil {
		return nil, errors.New("display: spi port must not be nil")
	}
	if cfg.Width <= 0 {
		cfg.Width = max7219Width
	}
	if cfg.Height <= 0 {
		cfg.Height = max7219Height
	}
	if cfg.SPIHz <= 0 {
		cfg.SPIHz = 10 * physic.MegaHertz
	}

	devices := cfg.Width / 8
	if devices < 1 {
		devices = 1
	}

	conn, err := port.Connect(cfg.SPIHz, spi.Mode0, 8)
	if err != nil {
		return nil, err
	}

	d := &MAX7219{cfg: cfg, c: conn, devices: devices}
	if err := d.init(); err != nil {
		return nil, err
	}
	return d, nil
}

func (d *MAX7219) Bounds() image.Rectangle {
	return image.Rect(0, 0, d.cfg.Width, d.cfg.Height)
}

// DrawImage writes pixel data to the cascaded MAX7219 chain.
// Each device handles 8 columns. Rows are sent as register writes (digit 0-7).
func (d *MAX7219) DrawImage(src draw.Image) error {
	for row := 0; row < 8; row++ {
		// Build a frame: for each device, a 2-byte (register, data) pair.
		buf := make([]byte, d.devices*2)
		for dev := 0; dev < d.devices; dev++ {
			var colByte byte
			for bit := 0; bit < 8; bit++ {
				x := dev*8 + bit
				r, g, b, _ := src.At(x, row).RGBA()
				if luminanceOn(r, g, b) {
					colByte |= 1 << uint(7-bit)
				}
			}
			// Devices are daisy-chained: last device in chain receives first.
			// We send in reverse order so device 0 is the leftmost.
			idx := (d.devices - 1 - dev) * 2
			buf[idx] = regDigit0 + byte(row)
			buf[idx+1] = colByte
		}
		if err := d.c.Tx(buf, nil); err != nil {
			return err
		}
	}
	return nil
}

func (d *MAX7219) init() error {
	// Send initialization commands to all devices in the chain.
	if err := d.writeAll(regDisplayTest, 0x00); err != nil {
		return err
	}
	if err := d.writeAll(regShutdown, 0x01); err != nil {
		return err
	}
	if err := d.writeAll(regDecodeMode, 0x00); err != nil {
		return err
	}
	if err := d.writeAll(regScanLimit, 0x07); err != nil {
		return err
	}
	if err := d.writeAll(regIntensity, 0x08); err != nil {
		return err
	}

	// Clear all digit registers.
	for row := byte(0); row < 8; row++ {
		if err := d.writeAll(regDigit0+row, 0x00); err != nil {
			return err
		}
	}
	return nil
}

// writeAll sends the same register+value pair to all devices in the chain.
func (d *MAX7219) writeAll(reg, val byte) error {
	buf := make([]byte, d.devices*2)
	for i := 0; i < d.devices; i++ {
		buf[i*2] = reg
		buf[i*2+1] = val
	}
	return d.c.Tx(buf, nil)
}

func luminanceOn(r, g, b uint32) bool {
	return (299*r+587*g+114*b)/1000 > 0x3000
}
