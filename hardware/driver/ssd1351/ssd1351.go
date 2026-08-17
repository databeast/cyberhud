package ssd1351

import (
	"errors"
	"image"
	"image/draw"
	"time"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/conn/v3/spi"

	"github.com/databeast/cyberhud/hardware/driver"
)

const (
	ssd1351Width  = 128
	ssd1351Height = 128

	cmdSetColumn      byte = 0x15
	cmdSetRow         byte = 0x75
	cmdWriteRAM       byte = 0x5C
	cmdSetRemap       byte = 0xA0
	cmdSetStartLine   byte = 0xA1
	cmdSetOffset      byte = 0xA2
	cmdFunctionSel    byte = 0xAB
	cmdDisplayOff     byte = 0xAE
	cmdDisplayOn      byte = 0xAF
	cmdSetPrechargeV  byte = 0xBB
	cmdSetVCOMH       byte = 0xBE
	cmdSetContrast    byte = 0xC1
	cmdMasterContrast byte = 0xC7
	cmdSetMuxRatio    byte = 0xCA
	cmdSetPrecharge   byte = 0xB1
	cmdClockDiv       byte = 0xB3
	cmdSetPrecharge2  byte = 0xB6
	cmdCmdLock        byte = 0xFD
)

func init() {
	driver.Register(driver.Definition{
		ID:         "ssd1351",
		Title:      "SSD1351",
		Summary:    "SPI 128x128 RGB color OLED controller.",
		Monochrome: false,
		OptionDefs: []driver.OptionDefinition{
			{Key: "width", Type: "int", Summary: "Display width in pixels.", Default: "128"},
			{Key: "height", Type: "int", Summary: "Display height in pixels.", Default: "128"},
			{Key: "spi_hz", Type: "frequency", Summary: "SPI bus speed.", Default: "8MHz"},
		},
		NewSPI: func(port spi.Port, dc, rst, bl gpio.PinOut, _ gpio.PinIn, cfg driver.DriverConfig) (driver.DrawTarget, error) {
			return New(port, dc, rst, bl, cfg)
		},
	})
}

// SSD1351 drives an SSD1351 color OLED over SPI.
type SSD1351 struct {
	cfg driver.DriverConfig
	c   spi.Conn
	dc  gpio.PinOut
	rst gpio.PinOut
	bl  gpio.PinOut
}

// New creates and initializes an SSD1351 display.
func New(port spi.Port, dc, rst, bl gpio.PinOut, cfg driver.DriverConfig) (*SSD1351, error) {
	if port == nil {
		return nil, errors.New("display: spi port must not be nil")
	}
	if dc == nil || rst == nil {
		return nil, errors.New("display: dc and rst pins must not be nil")
	}
	if cfg.Width <= 0 {
		cfg.Width = ssd1351Width
	}
	if cfg.Height <= 0 {
		cfg.Height = ssd1351Height
	}
	if cfg.SPIHz <= 0 {
		cfg.SPIHz = 8 * physic.MegaHertz
	}

	conn, err := port.Connect(cfg.SPIHz, spi.Mode0, 8)
	if err != nil {
		return nil, err
	}

	d := &SSD1351{cfg: cfg, c: conn, dc: dc, rst: rst, bl: bl}
	if err := d.init(); err != nil {
		return nil, err
	}
	return d, nil
}

func (d *SSD1351) Bounds() image.Rectangle {
	return image.Rect(0, 0, d.cfg.Width, d.cfg.Height)
}

// DrawImage writes 16-bit RGB565 pixel data to the OLED GDDRAM.
func (d *SSD1351) DrawImage(src draw.Image) error {
	w := d.cfg.Width
	h := d.cfg.Height

	// Set column address window.
	if err := d.cmdData(cmdSetColumn, []byte{0x00, byte(w - 1)}); err != nil {
		return err
	}
	// Set row address window.
	if err := d.cmdData(cmdSetRow, []byte{0x00, byte(h - 1)}); err != nil {
		return err
	}
	// Begin write.
	if err := d.cmd(cmdWriteRAM); err != nil {
		return err
	}

	// Pack pixels as RGB565 (big-endian, 2 bytes per pixel).
	buf := make([]byte, w*h*2)
	idx := 0
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, g, b, _ := src.At(x, y).RGBA()
			// RGBA() returns 0–65535; scale to 5-6-5 bits.
			r5 := uint16(r>>11) & 0x1F
			g6 := uint16(g>>10) & 0x3F
			b5 := uint16(b>>11) & 0x1F
			rgb565 := (r5 << 11) | (g6 << 5) | b5
			buf[idx] = byte(rgb565 >> 8)
			buf[idx+1] = byte(rgb565)
			idx += 2
		}
	}

	return d.data(buf)
}

func (d *SSD1351) init() error {
	// Hardware reset.
	if err := d.rst.Out(gpio.Low); err != nil {
		return err
	}
	time.Sleep(10 * time.Millisecond)
	if err := d.rst.Out(gpio.High); err != nil {
		return err
	}
	time.Sleep(10 * time.Millisecond)

	// Unlock commands.
	if err := d.cmdData(cmdCmdLock, []byte{0x12}); err != nil {
		return err
	}
	if err := d.cmdData(cmdCmdLock, []byte{0xB1}); err != nil {
		return err
	}

	if err := d.cmd(cmdDisplayOff); err != nil {
		return err
	}
	if err := d.cmdData(cmdClockDiv, []byte{0xF1}); err != nil {
		return err
	}
	if err := d.cmdData(cmdSetMuxRatio, []byte{byte(d.cfg.Height - 1)}); err != nil {
		return err
	}
	// Remap: 65k color, COM split odd/even, color order RGB.
	if err := d.cmdData(cmdSetRemap, []byte{0x74}); err != nil {
		return err
	}
	if err := d.cmdData(cmdSetStartLine, []byte{0x00}); err != nil {
		return err
	}
	if err := d.cmdData(cmdSetOffset, []byte{0x00}); err != nil {
		return err
	}
	if err := d.cmdData(cmdFunctionSel, []byte{0x01}); err != nil {
		return err
	}
	if err := d.cmdData(cmdSetPrecharge, []byte{0x32}); err != nil {
		return err
	}
	if err := d.cmdData(cmdSetVCOMH, []byte{0x05}); err != nil {
		return err
	}
	if err := d.cmdData(cmdSetPrechargeV, []byte{0x17}); err != nil {
		return err
	}
	// Contrast per color channel.
	if err := d.cmdData(cmdSetContrast, []byte{0xC8, 0x80, 0xC8}); err != nil {
		return err
	}
	if err := d.cmdData(cmdMasterContrast, []byte{0x0F}); err != nil {
		return err
	}
	if err := d.cmdData(cmdSetPrecharge2, []byte{0x01}); err != nil {
		return err
	}
	if err := d.cmd(cmdDisplayOn); err != nil {
		return err
	}

	if d.bl != nil {
		return d.bl.Out(gpio.High)
	}
	return nil
}

func (d *SSD1351) cmdData(c byte, payload []byte) error {
	if err := d.cmd(c); err != nil {
		return err
	}
	if len(payload) == 0 {
		return nil
	}
	return d.data(payload)
}

func (d *SSD1351) cmd(c byte) error {
	if err := d.dc.Out(gpio.Low); err != nil {
		return err
	}
	return d.c.Tx([]byte{c}, nil)
}

func (d *SSD1351) data(buf []byte) error {
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
