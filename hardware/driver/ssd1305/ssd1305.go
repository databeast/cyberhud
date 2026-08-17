package ssd1305

import (
	"errors"
	"fmt"
	"image"
	"image/draw"
	"time"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/conn/v3/spi"

	"github.com/databeast/cyberhud/hardware/driver"
)

const (
	defaultWidth  = 128
	defaultHeight = 32
	defaultAddr   = 0x3C
	defaultSPIHz  = 1 * physic.MegaHertz
)

func init() {
	driver.Register(driver.Definition{
		ID:         "ssd1305",
		Title:      "SSD1305",
		Summary:    "SPI/I2C monochrome OLED controller (128x32).",
		Monochrome: true,
		OptionDefs: []driver.OptionDefinition{
			{Key: "width", Type: "int", Summary: "Display width in pixels.", Default: "128"},
			{Key: "height", Type: "int", Summary: "Display height in pixels.", Default: "32"},
			{Key: "spi_hz", Type: "frequency", Summary: "SPI bus speed.", Default: "1MHz"},
			{Key: "address", Type: "byte", Summary: "I2C device address.", Default: "0x3C"},
		},
		NewSPI: func(port spi.Port, dc, rst, bl gpio.PinOut, _ gpio.PinIn, cfg driver.DriverConfig) (driver.DrawTarget, error) {
			return NewSPI(port, dc, rst, bl, cfg)
		},
		NewI2C: func(bus i2c.Bus, cfg driver.DriverConfig) (driver.DrawTarget, error) {
			return NewI2C(bus, cfg)
		},
	})
}

// SSD1305 drives an SSD1305 OLED display over SPI or I2C.
type SSD1305 struct {
	cfg driver.DriverConfig
	// SPI transport (non-nil when created via NewSPI)
	spiConn spi.Conn
	dc      gpio.PinOut
	rst     gpio.PinOut
	bl      gpio.PinOut
	// I2C transport (non-nil when created via NewI2C)
	i2cDev *i2c.Dev
}

// NewSPI constructs an SSD1305 instance using SPI transport.
func NewSPI(port spi.Port, dc, rst, bl gpio.PinOut, cfg driver.DriverConfig) (*SSD1305, error) {
	if port == nil {
		return nil, errors.New("ssd1305: spi port must not be nil")
	}
	if dc == nil {
		return nil, errors.New("ssd1305: dc pin must not be nil")
	}
	if rst == nil {
		return nil, errors.New("ssd1305: rst pin must not be nil")
	}

	applyDefaults(&cfg)

	if cfg.SPIHz <= 0 {
		cfg.SPIHz = defaultSPIHz
	}

	conn, err := port.Connect(cfg.SPIHz, spi.Mode3, 8)
	if err != nil {
		return nil, err
	}

	d := &SSD1305{
		cfg:     cfg,
		spiConn: conn,
		dc:      dc,
		rst:     rst,
		bl:      bl,
	}

	if err := d.initSPI(); err != nil {
		return nil, err
	}
	return d, nil
}

// NewI2C constructs an SSD1305 instance using I2C transport.
func NewI2C(bus i2c.Bus, cfg driver.DriverConfig) (*SSD1305, error) {
	if bus == nil {
		return nil, errors.New("ssd1305: i2c bus must not be nil")
	}

	applyDefaults(&cfg)

	addr := uint16(defaultAddr)
	if cfg.I2CAddr != 0 {
		addr = cfg.I2CAddr
	}

	dev := &i2c.Dev{Bus: bus, Addr: addr}
	d := &SSD1305{
		cfg:    cfg,
		i2cDev: dev,
	}

	if err := d.initI2C(); err != nil {
		return nil, fmt.Errorf("ssd1305 [0x%02X]: %w", addr, err)
	}
	return d, nil
}

// Bounds returns the display rectangle.
func (d *SSD1305) Bounds() image.Rectangle {
	return image.Rect(0, 0, d.cfg.Width, d.cfg.Height)
}

// DrawImage converts pixel data to page-column format and transmits to the display.
// Uses page addressing mode with column offset 4 (SSD1305 has 132-col internal buffer,
// display starts at column 4).
func (d *SSD1305) DrawImage(src draw.Image) error {
	w := d.cfg.Width
	h := d.cfg.Height
	pages := h / 8

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

	// Write using page addressing: set page, then column, then data for each page.
	for page := 0; page < pages; page++ {
		// Set page address
		if err := d.cmd(0xB0 + byte(page)); err != nil {
			return err
		}
		// Set lower column address (offset 4 for 132-col buffer)
		if err := d.cmd(0x04); err != nil {
			return err
		}
		// Set higher column address
		if err := d.cmd(0x10); err != nil {
			return err
		}
		// Write page data
		pageData := buf[page*w : (page+1)*w]
		if d.spiConn != nil {
			if err := d.spiData(pageData); err != nil {
				return err
			}
		} else {
			if err := d.i2cData(pageData); err != nil {
				return err
			}
		}
	}
	return nil
}

// cmd sends a command byte via the active transport.
func (d *SSD1305) cmd(c byte) error {
	if d.spiConn != nil {
		return d.spiCmd(c)
	}
	return d.i2cCmd(c)
}

// initCommands returns the SSD1305 initialization command sequence.
// When cfg.Rotate180 is true, segment remap and COM scan are flipped for 180° rotation.
func (d *SSD1305) initCommands() []byte {
	segRemap := byte(0xA1) // Normal: col 127 → SEG0
	comScan := byte(0xC8)  // Normal: remapped (bottom-to-top)
	if d.cfg.Rotate180 {
		segRemap = 0xA0 // Rotated: col 0 → SEG0
		comScan = 0xC0  // Rotated: normal (top-to-bottom)
	}

	return []byte{
		0xAE,       // Display OFF
		0x04,       // Set lower column start address (offset 4 for 132-col buffer)
		0x10,       // Set higher column start address
		0x40,       // Set display start line to 0
		0x81, 0x80, // Set contrast
		segRemap,   // Set segment remap
		0xA6,       // Set normal display (not inverted)
		0xA8, 0x1F, // Set multiplex ratio to 31 (32 rows)
		comScan,    // Set COM output scan direction
		0xD3, 0x00, // Set display offset to 0
		0xD5, 0xF0, // Set clock divide ratio/oscillator frequency
		0xD8, 0x05, // Set area color mode / low power display mode
		0xD9, 0xC2, // Set pre-charge period
		0xDA, 0x12, // Set COM pins hardware configuration
		0xDB, 0x08, // Set VCOMH deselect level
		0xAF, // Display ON
	}
}

// luminanceOn returns true if the pixel luminance exceeds the threshold.
func luminanceOn(r, g, b uint32) bool {
	return (299*r+587*g+114*b)/1000 > 0x3000
}

// applyDefaults sets width and height to defaults when either is not specified.
// If either dimension is invalid (≤ 0), both default to 128×32.
func applyDefaults(cfg *driver.DriverConfig) {
	if cfg.Width <= 0 || cfg.Height <= 0 {
		cfg.Width = defaultWidth
		cfg.Height = defaultHeight
	}
}

// initSPI performs the SSD1305 initialization sequence over SPI.
func (d *SSD1305) initSPI() error {
	// Hardware reset: RST high → low → high.
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

	for _, c := range d.initCommands() {
		if err := d.spiCmd(c); err != nil {
			return err
		}
	}

	// Set backlight pin HIGH if non-nil.
	if d.bl != nil {
		if err := d.bl.Out(gpio.High); err != nil {
			return err
		}
	}

	return nil
}

// spiCmd sends a single command byte over SPI with DC pin LOW.
func (d *SSD1305) spiCmd(c byte) error {
	if err := d.dc.Out(gpio.Low); err != nil {
		return err
	}
	return d.spiConn.Tx([]byte{c}, nil)
}

// spiData sends data bytes over SPI with DC pin HIGH.
func (d *SSD1305) spiData(buf []byte) error {
	if err := d.dc.Out(gpio.High); err != nil {
		return err
	}
	return d.spiConn.Tx(buf, nil)
}

// initI2C performs the SSD1305 initialization sequence over I2C.
func (d *SSD1305) initI2C() error {
	for _, c := range d.initCommands() {
		if err := d.i2cCmd(c); err != nil {
			return err
		}
	}

	return nil
}

// i2cCmd sends a single command byte over I2C with control byte 0x00.
func (d *SSD1305) i2cCmd(c byte) error {
	return d.i2cDev.Tx([]byte{0x00, c}, nil)
}

// i2cData sends data bytes over I2C with control byte 0x40.
func (d *SSD1305) i2cData(buf []byte) error {
	out := make([]byte, 1+len(buf))
	out[0] = 0x40
	copy(out[1:], buf)
	return d.i2cDev.Tx(out, nil)
}
