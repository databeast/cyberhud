package is31fl3731

import (
	"fmt"
	"image"
	"image/draw"

	"periph.io/x/conn/v3/i2c"

	"github.com/databeast/cyberhud/hardware/driver"
)

func init() {
	driver.Register(driver.Definition{
		ID:         "is31fl3731",
		Title:      "IS31FL3731",
		Summary:    "Charlieplexed LED matrix controller (I2C)",
		Monochrome: false,
		OptionDefs: []driver.OptionDefinition{
			{Key: "width", Type: "int", Summary: "Display width in pixels", Default: "15"},
			{Key: "height", Type: "int", Summary: "Display height in pixels", Default: "7"},
			{Key: "address", Type: "byte", Summary: "I2C device address", Default: "0x74"},
			{Key: "layout", Type: "string", Summary: "Pixel-address layout: linear, charlie-wing (15x7 FeatherWing), charlie-bonnet (16x8 Bonnet)", Default: "linear"},
		},
		NewI2C: func(bus i2c.Bus, cfg driver.DriverConfig) (driver.DrawTarget, error) {
			return New(bus, cfg)
		},
	})
}

const (
	// DefaultAddr is the default I2C address for the IS31FL3731 (Adafruit default
	// with address select pads open).
	DefaultAddr uint16 = 0x74

	// Register page select command register.
	regPageSelect byte = 0xFD

	// Function register page (page 8).
	pageFunctionReg byte = 0x0B

	// Frame 0 page.
	pageFrame0 byte = 0x00

	// Function registers.
	regShutdown     byte = 0x0A
	regMode         byte = 0x00
	regDisplayFrame byte = 0x01

	// LED enable register range within a frame page.
	ledEnableStart byte = 0x00
	ledEnableEnd   byte = 0x11

	// PWM register range within a frame page (144 registers, one per LED).
	pwmStart byte = 0x24
	pwmEnd   byte = 0xB3

	// pwmCount is the number of PWM registers (and addressable LEDs) per frame.
	pwmCount = int(pwmEnd) - int(pwmStart) + 1
)

// pixelAddrFunc maps a panel (x, y) coordinate to a PWM register offset (0–143).
type pixelAddrFunc func(x, y int) int

// layoutDef describes a named non-linear pixel-address layout and the native
// panel dimensions it was designed for.
type layoutDef struct {
	width  int
	height int
	addr   pixelAddrFunc
}

// Named layouts matching Adafruit's CircuitPython IS31FL3731 driver formulas.
var layouts = map[string]layoutDef{
	// Adafruit 15x7 CharliePlex LED Matrix FeatherWing (SKU 2965).
	"charlie-wing": {width: 15, height: 7, addr: func(x, y int) int {
		if x > 7 {
			x = 15 - x
			y += 8
		} else {
			y = 7 - y
		}
		return x*16 + y
	}},
	// Adafruit 16x8 CharliePlex LED Matrix Bonnet (SKUs 4120/4122/4127).
	"charlie-bonnet": {width: 16, height: 8, addr: func(x, y int) int {
		if x >= 8 {
			return (x-6)*16 - (y + 1)
		}
		return (x+1)*16 + (7 - y)
	}},
}

// resolvePixelAddr returns the pixel-address function for the configured layout
// and dimensions. The empty layout and "linear" use row-major addressing. Named
// layouts accept either their native dimensions or the transpose (panel mounted
// rotated 90°), in which case coordinates are rotated before mapping.
func resolvePixelAddr(cfg driver.DriverConfig) (pixelAddrFunc, error) {
	switch cfg.Layout {
	case "", "linear":
		w := cfg.Width
		return func(x, y int) int { return y*w + x }, nil
	}

	l, ok := layouts[cfg.Layout]
	if !ok {
		return nil, fmt.Errorf("unknown layout %q", cfg.Layout)
	}

	switch {
	case cfg.Width == l.width && cfg.Height == l.height:
		return l.addr, nil
	case cfg.Width == l.height && cfg.Height == l.width:
		// Panel mounted rotated 90°: rotate logical coordinates into the
		// layout's native landscape orientation.
		return func(x, y int) int {
			return l.addr(y, cfg.Width-1-x)
		}, nil
	default:
		return nil, fmt.Errorf("layout %q requires %dx%d (or %dx%d) dimensions, got %dx%d",
			cfg.Layout, l.width, l.height, l.height, l.width, cfg.Width, cfg.Height)
	}
}

// IS31FL3731 drives an IS31FL3731 charlieplexed LED matrix controller over I2C.
type IS31FL3731 struct {
	dev       *i2c.Dev
	cfg       driver.DriverConfig
	pixelAddr pixelAddrFunc
	bufLen    int // PWM buffer length: W*H for linear, full 144 for named layouts
}

// New opens the IS31FL3731 at the configured I2C address (default 0x74) and runs
// the hardware initialization sequence. It returns an error containing the target
// address (hex) and the underlying reason if any I2C write fails.
func New(bus i2c.Bus, cfg driver.DriverConfig) (*IS31FL3731, error) {
	addr := DefaultAddr
	if cfg.I2CAddr != 0 {
		addr = cfg.I2CAddr
	}

	pixelAddr, err := resolvePixelAddr(cfg)
	if err != nil {
		return nil, fmt.Errorf("is31fl3731 [0x%02X]: %w", addr, err)
	}

	bufLen := cfg.Width * cfg.Height
	if cfg.Layout != "" && cfg.Layout != "linear" {
		bufLen = pwmCount
	}

	dev := &i2c.Dev{Bus: bus, Addr: addr}
	d := &IS31FL3731{dev: dev, cfg: cfg, pixelAddr: pixelAddr, bufLen: bufLen}

	if err := d.init(); err != nil {
		return nil, fmt.Errorf("is31fl3731 [0x%02X]: %w", addr, err)
	}
	return d, nil
}

// init performs the hardware initialization sequence.
func (d *IS31FL3731) init() error {
	// 1. Select function register page (page 8).
	if err := d.writeReg(regPageSelect, pageFunctionReg); err != nil {
		return fmt.Errorf("select function page: %w", err)
	}

	// 2. Exit software shutdown.
	if err := d.writeReg(regShutdown, 0x01); err != nil {
		return fmt.Errorf("exit shutdown: %w", err)
	}

	// 3. Set picture mode.
	if err := d.writeReg(regMode, 0x00); err != nil {
		return fmt.Errorf("set picture mode: %w", err)
	}

	// 4. Set display frame to Frame 0.
	if err := d.writeReg(regDisplayFrame, 0x00); err != nil {
		return fmt.Errorf("set display frame: %w", err)
	}

	// 5. Select Frame 0 page.
	if err := d.writeReg(regPageSelect, pageFrame0); err != nil {
		return fmt.Errorf("select frame 0 page: %w", err)
	}

	// 6. Enable exactly the LEDs addressed by the configured layout
	// (registers 0x00–0x11, 18 bytes = 144 bits, bit n of byte m = LED m*8+n).
	enableBuf := make([]byte, 18)
	for y := 0; y < d.cfg.Height; y++ {
		for x := 0; x < d.cfg.Width; x++ {
			addr := d.pixelAddr(x, y)
			if addr < 0 || addr >= pwmCount {
				continue
			}
			enableBuf[addr/8] |= 1 << (addr % 8)
		}
	}

	if err := d.writeBulk(ledEnableStart, enableBuf); err != nil {
		return fmt.Errorf("enable LEDs: %w", err)
	}

	// 7. Zero all PWM registers (0x24–0xB3 = 144 bytes).
	pwmBuf := make([]byte, pwmCount)
	if err := d.writeBulk(pwmStart, pwmBuf); err != nil {
		return fmt.Errorf("zero PWM registers: %w", err)
	}

	return nil
}

// writeReg writes a single byte to an I2C register.
func (d *IS31FL3731) writeReg(reg, val byte) error {
	return d.dev.Tx([]byte{reg, val}, nil)
}

// writeBulk writes a slice of bytes starting at the given register address.
func (d *IS31FL3731) writeBulk(startReg byte, data []byte) error {
	buf := make([]byte, 1+len(data))
	buf[0] = startReg
	copy(buf[1:], data)
	return d.dev.Tx(buf, nil)
}

// Bounds returns the pixel dimensions of the LED matrix.
func (d *IS31FL3731) Bounds() image.Rectangle {
	return image.Rect(0, 0, d.cfg.Width, d.cfg.Height)
}

// DrawImage converts the source image pixels to grayscale brightness values using
// BT.601 luminance and writes them to the PWM registers, placing each pixel at
// the register offset given by the configured pixel-address layout.
// Only pixels within the configured bounds are processed; pixels outside are ignored.
// Any I2C write error is propagated to the caller.
func (d *IS31FL3731) DrawImage(src draw.Image) error {
	bounds := d.Bounds()
	width := bounds.Max.X
	height := bounds.Max.Y

	buf := make([]byte, d.bufLen)
	for row := 0; row < height; row++ {
		for col := 0; col < width; col++ {
			addr := d.pixelAddr(col, row)
			if addr < 0 || addr >= len(buf) {
				continue
			}
			r, g, b, _ := src.At(col, row).RGBA()
			// RGBA() returns values in 0–65535 range; shift to 0–255.
			buf[addr] = luminance(uint32(r>>8), uint32(g>>8), uint32(b>>8))
		}
	}

	return d.writeBulk(pwmStart, buf)
}

// luminance computes BT.601 grayscale brightness from 0–255 range RGB values.
func luminance(r, g, b uint32) byte {
	y := (299*r + 587*g + 114*b) / 1000
	if y > 255 {
		return 255
	}
	return byte(y)
}
