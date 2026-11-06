// Package is31fl3741 drives ISSI IS31FL3741-based RGB LED matrices over I2C —
// notably the Adafruit 13x9 PWM RGB LED Matrix QT board (product 5201), 117 RGB
// LEDs with 8-bit PWM per channel for 24-bit colour.
//
// Unlike the charlieplexed IS31FL3731 (a grayscale single-colour matrix), the
// IS31FL3741 drives independent R/G/B channels, so this is a true colour panel.
//
// Protocol. The chip exposes 351 PWM channels split across two register pages
// (page 0: LEDs 0..179, page 1: LEDs 180..350). Page selection goes through the
// command register (0xFD) and must be preceded by writing the unlock value
// 0xC5 to the command-lock register (0xFE) each time. Page 4 holds the function
// registers (config/shutdown at 0x00, global current at 0x01, reset at 0x3F);
// pages 2 and 3 hold per-LED scaling. Initialisation resets the chip, sets all
// scaling to full, sets a global current, and clears the shutdown bit.
//
// Pixel mapping. The Adafruit 13x9 QT board uses a fixed row remap and a
// column-parity-dependent R/G/B sub-offset, replicated here from Adafruit's
// driver. Default colour order is BGR, matching the shipping QT matrices.
package is31fl3741

import (
	"fmt"
	"image"
	"image/draw"

	"periph.io/x/conn/v3/i2c"

	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/hardware/driver"
)

const (
	// DefaultAddr is the IS31FL3741 default I2C address (ADDR pads open). The
	// board can be strapped to 0x30/0x31/0x33/0x34 for tiling.
	DefaultAddr uint16 = 0x30

	regCommand     byte = 0xFD // page-select register
	regCommandLock byte = 0xFE // write unlockValue here before changing page
	regID          byte = 0xFC // reads 2*address

	unlockValue byte = 0xC5

	// Function-register page (page 4) registers.
	funcRegConfig   byte = 0x00 // bit0 = shutdown (0 = shutdown, 1 = normal)
	funcRegGCurrent byte = 0x01
	funcRegReset    byte = 0x3F
	resetValue      byte = 0xAE

	pagePWM0    byte = 0x00 // LEDs 0..179
	pagePWM1    byte = 0x01 // LEDs 180..350
	pageScale0  byte = 0x02 // scaling 0..179
	pageScale1  byte = 0x03 // scaling 180..350
	pageFunc    byte = 0x04
	maxPage     byte = 0x04
	pageUnset   byte = 0xFF // sentinel: no page selected yet
	ledCount         = 351  // total PWM channels
	page0Count       = 180  // channels addressable on page 0
	page1Count       = ledCount - page0Count

	// defaultGlobalCurrent balances brightness against the ~40mA/column sink
	// limit; Adafruit examples use a mid value. Tunable later if needed.
	defaultGlobalCurrent byte = 0xFF
)

func init() {
	driver.Register(driver.Definition{
		ID:         "is31fl3741",
		Title:      "IS31FL3741",
		Summary:    "I2C RGB LED matrix controller (24-bit color) used by the Adafruit 13x9 RGB matrix.",
		Monochrome: false,
		ColorDepth: driver.DepthColor,
		Capability: textlayout.CapColorFast,
		OptionDefs: []driver.OptionDefinition{
			{Key: "width", Type: "int", Summary: "Matrix width in pixels.", Default: "13"},
			{Key: "height", Type: "int", Summary: "Matrix height in pixels.", Default: "9"},
			{Key: "address", Type: "byte", Summary: "I2C device address (0x30 default; 0x31/0x33/0x34 for tiling).", Default: "0x30"},
			{Key: "layout", Type: "string", Summary: "Pixel-address layout: qt-13x9 (Adafruit 13x9 QT matrix).", Default: "qt-13x9"},
		},
		NewI2C: func(bus i2c.Bus, cfg driver.DriverConfig) (driver.DrawTarget, error) {
			return New(bus, cfg)
		},
	})
}

// rgbAddrFunc maps a panel (x, y) coordinate to its three PWM channel offsets
// (red, green, blue) in device order (before colour-order remap).
type rgbAddrFunc func(x, y int) (r, g, b int)

// qtRowmap reorders logical rows into the QT board's physical row wiring.
var qtRowmap = [9]int{8, 5, 4, 3, 2, 1, 0, 7, 6}

// qtPixelAddrs replicates Adafruit_RGBMatrixQT.pixel_addrs for the 13x9 board:
// a row remap, a split column base offset, and a column-parity-dependent sub-
// offset triple. The three returned channel addresses are in the board's
// positional order; the shipping QT matrix is wired BGR (see channelOrder), so
// the caller maps red/green/blue onto these positions via that order rather
// than assuming the first position is red.
func qtPixelAddrs(x, y int) (int, int, int) {
	y = qtRowmap[y]
	var base int
	if x < 10 {
		base = 3 * (x + y*10)
	} else {
		base = 3 * (x + 80 + y*3)
	}
	// Odd columns and the last column (12) use a swapped sub-offset triple.
	if x&1 == 1 || x == 12 {
		return base + 2, base + 0, base + 1
	}
	return base + 0, base + 1, base + 2
}

// IS31FL3741 drives an IS31FL3741 RGB LED matrix over I2C.
type IS31FL3741 struct {
	dev   *i2c.Dev
	cfg   driver.DriverConfig
	addrs rgbAddrFunc
	// rOff/gOff/bOff select which of the three positional channel addresses
	// returned by addrs carries red/green/blue. The shipping QT matrix is wired
	// BGR, so red→2, green→1, blue→0.
	rOff, gOff, bOff int
	page             byte
}

// New opens the IS31FL3741 at the configured I2C address (default 0x30) and runs
// the hardware initialisation sequence.
func New(bus i2c.Bus, cfg driver.DriverConfig) (*IS31FL3741, error) {
	addr := DefaultAddr
	if cfg.I2CAddr != 0 {
		addr = cfg.I2CAddr
	}

	addrs, err := resolveAddrFunc(cfg)
	if err != nil {
		return nil, fmt.Errorf("is31fl3741 [0x%02X]: %w", addr, err)
	}

	// Default BGR channel order, matching the shipping Adafruit QT matrices.
	d := &IS31FL3741{dev: &i2c.Dev{Bus: bus, Addr: addr}, cfg: cfg, addrs: addrs, rOff: 2, gOff: 1, bOff: 0, page: pageUnset}
	if err := d.init(); err != nil {
		return nil, fmt.Errorf("is31fl3741 [0x%02X]: %w", addr, err)
	}
	return d, nil
}

// resolveAddrFunc selects the per-pixel RGB address mapping for the configured
// layout. Only the Adafruit 13x9 QT board layout is supported today.
func resolveAddrFunc(cfg driver.DriverConfig) (rgbAddrFunc, error) {
	switch cfg.Layout {
	case "", "qt-13x9":
		if (cfg.Width != 13 || cfg.Height != 9) && cfg.Width != 0 && cfg.Height != 0 {
			return nil, fmt.Errorf("layout qt-13x9 requires 13x9 dimensions, got %dx%d", cfg.Width, cfg.Height)
		}
		return qtPixelAddrs, nil
	default:
		return nil, fmt.Errorf("unknown layout %q", cfg.Layout)
	}
}

// Bounds returns the pixel dimensions of the LED matrix.
func (d *IS31FL3741) Bounds() image.Rectangle {
	w, h := d.cfg.Width, d.cfg.Height
	if w == 0 {
		w = 13
	}
	if h == 0 {
		h = 9
	}
	return image.Rect(0, 0, w, h)
}

// TextHints reports this panel as a fast 24-bit colour matrix. DrawImage writes
// 8-bit PWM per channel, so emitted colour reaches the LEDs; declaring color-fast
// lets style selection offer colour styles.
func (d *IS31FL3741) TextHints() textlayout.TextHints {
	b := d.Bounds()
	return textlayout.TextHints{
		PixelWidth: b.Dx(), PixelHeight: b.Dy(),
		SupportsVerticalScroll: false, SupportsHorizontalScroll: false, SupportsAutoScroll: false,
		PreferEventRefresh:     false,
		Capability:             textlayout.CapColorFast,
		DefaultTickerDirection: textlayout.TickerDirectionNone, DefaultLineMode: textlayout.LineModeTruncate,
	}
}

func (d *IS31FL3741) init() error {
	// Verify the chip identity (ID register reads 2*address).
	if err := d.setPage(pageFunc); err != nil {
		return err
	}
	// Reset.
	if err := d.writeReg(funcRegReset, resetValue); err != nil {
		return fmt.Errorf("reset: %w", err)
	}
	// Reset returns the chip to page 4 but clears our cache intent; re-select.
	d.page = pageUnset

	// Set per-LED scaling to full on both scaling pages (0..179 and 180..350).
	full := make([]byte, page0Count)
	for i := range full {
		full[i] = 0xFF
	}
	if err := d.setPage(pageScale0); err != nil {
		return err
	}
	if err := d.writeBulk(0x00, full); err != nil {
		return fmt.Errorf("scaling page 0: %w", err)
	}
	if err := d.setPage(pageScale1); err != nil {
		return err
	}
	if err := d.writeBulk(0x00, full[:page1Count]); err != nil {
		return fmt.Errorf("scaling page 1: %w", err)
	}

	// Function page: set global current and clear the shutdown bit (enable).
	if err := d.setPage(pageFunc); err != nil {
		return err
	}
	if err := d.writeReg(funcRegGCurrent, defaultGlobalCurrent); err != nil {
		return fmt.Errorf("global current: %w", err)
	}
	if err := d.writeReg(funcRegConfig, 0x01); err != nil { // bit0 = normal operation
		return fmt.Errorf("enable: %w", err)
	}

	// Clear all PWM channels.
	zero := make([]byte, page0Count)
	if err := d.setPage(pagePWM0); err != nil {
		return err
	}
	if err := d.writeBulk(0x00, zero); err != nil {
		return fmt.Errorf("clear PWM page 0: %w", err)
	}
	if err := d.setPage(pagePWM1); err != nil {
		return err
	}
	if err := d.writeBulk(0x00, zero[:page1Count]); err != nil {
		return fmt.Errorf("clear PWM page 1: %w", err)
	}
	return nil
}

// DrawImage converts source pixels to 8-bit-per-channel RGB and writes the full
// PWM channel buffer across both pages. Pixels outside the configured bounds are
// ignored.
func (d *IS31FL3741) DrawImage(src draw.Image) error {
	b := d.Bounds()
	buf := make([]byte, ledCount)
	set := func(addr int, v byte) {
		if addr >= 0 && addr < ledCount {
			buf[addr] = v
		}
	}
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			// addrs returns three positional channel addresses; the colour
			// order (BGR on this board) decides which carries r/g/b.
			var pos [3]int
			pos[0], pos[1], pos[2] = d.addrs(x, y)
			r16, g16, b16, _ := src.At(x, y).RGBA()
			set(pos[d.rOff], byte(r16>>8))
			set(pos[d.gOff], byte(g16>>8))
			set(pos[d.bOff], byte(b16>>8))
		}
	}
	if err := d.setPage(pagePWM0); err != nil {
		return err
	}
	if err := d.writeBulk(0x00, buf[:page0Count]); err != nil {
		return err
	}
	if err := d.setPage(pagePWM1); err != nil {
		return err
	}
	return d.writeBulk(0x00, buf[page0Count:])
}

// setPage unlocks and selects a register page, caching the current page to
// avoid redundant writes.
func (d *IS31FL3741) setPage(page byte) error {
	if page > maxPage {
		return fmt.Errorf("invalid page %d", page)
	}
	if d.page == page {
		return nil
	}
	if err := d.writeReg(regCommandLock, unlockValue); err != nil {
		return fmt.Errorf("unlock: %w", err)
	}
	if err := d.writeReg(regCommand, page); err != nil {
		return fmt.Errorf("select page %d: %w", page, err)
	}
	d.page = page
	return nil
}

func (d *IS31FL3741) writeReg(reg, val byte) error {
	return d.dev.Tx([]byte{reg, val}, nil)
}

func (d *IS31FL3741) writeBulk(startReg byte, data []byte) error {
	buf := make([]byte, 1+len(data))
	buf[0] = startReg
	copy(buf[1:], data)
	return d.dev.Tx(buf, nil)
}
