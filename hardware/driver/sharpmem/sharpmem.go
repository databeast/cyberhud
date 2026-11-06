// Package sharpmem drives Sharp Memory LCD panels over SPI — notably the
// Adafruit 1.3" 168x144 Monochrome Memory Display breakout (product 3502) and
// its siblings (96x96, 400x240, etc.).
//
// Sharp Memory LCDs are a distinct device class from the TFT/OLED/e-paper
// controllers: there is no controller command/register set, no init sequence,
// and no gamma/power programming. The panel is a write-only 1-bit framebuffer
// addressed one line at a time with a tiny bespoke protocol:
//
//   - A write frame begins with a mode byte [VCOM | WRITECMD(0x80)], then for
//     each line a 1-based line-address byte, the packed 1bpp line data, and a
//     trailing 0x00; the whole frame ends with one more trailing 0x00.
//   - A clear is the two bytes [VCOM | CLEAR(0x20), 0x00].
//   - The VCOM bit (0x40) must toggle on every transfer (and at least ~once per
//     second even when idle) to prevent a DC bias building up in the panel.
//
// Two hardware quirks matter on a Raspberry Pi:
//   - CS is ACTIVE HIGH on these panels, the opposite of normal SPI. The host's
//     SPI chip-select is normally active-low, so the panel's CS usually needs an
//     inverter or a separate GPIO driven high-during-transfer. This driver drives
//     the standard SPI transfers; CS polarity is a wiring concern.
//   - Bit order. Adafruit's driver clocks data LSB-first with the leftmost pixel
//     in the least-significant bit. periph transmits MSB-first, so we pack the
//     leftmost pixel into the most-significant bit — which puts the same pixel on
//     the wire first. No per-byte bit reversal is needed.
//
// The panel is monochrome with a fast LCD-like refresh (not e-paper), so it is
// advertised as mono-fast.
package sharpmem

import (
	"errors"
	"image"
	"image/draw"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/conn/v3/spi"

	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/hardware/driver"
)

const (
	sharpWidth  = 168
	sharpHeight = 144

	bitWriteCmd byte = 0x80 // MLCD_WR: write one or more lines
	bitVCOM     byte = 0x40 // MLCD_VCOM: VCOM inversion bit (toggled each transfer)
	bitClear    byte = 0x20 // MLCD_CM: clear all memory
)

func init() {
	driver.Register(driver.Definition{
		ID:         "sharpmem",
		Title:      "Sharp Memory LCD",
		Summary:    "SPI write-only 1-bit memory LCD (black-on-gray), fast refresh, no backlight.",
		Monochrome: true,
		ColorDepth: driver.DepthMono1,
		Capability: textlayout.CapMonoFast,
		OptionDefs: []driver.OptionDefinition{
			{Key: "width", Type: "int", Summary: "Display width in pixels.", Default: "168"},
			{Key: "height", Type: "int", Summary: "Display height in pixels.", Default: "144"},
			{Key: "spi_hz", Type: "frequency", Summary: "SPI bus speed (these panels top out around 2MHz).", Default: "2MHz"},
		},
		NewSPI: func(port spi.Port, dc, rst, bl gpio.PinOut, _ gpio.PinIn, cfg driver.DriverConfig) (driver.DrawTarget, error) {
			return New(port, cfg)
		},
	})
}

// Dev is a handle to a Sharp Memory LCD.
type Dev struct {
	c      spi.Conn
	width  int
	height int
	vcom   byte
}

// New initialises the Sharp Memory LCD. There is no controller init sequence;
// the panel is cleared to white and left ready for line writes.
//
// Note the DC/RST/BL pins from the generic factory are intentionally unused:
// Sharp Memory LCDs have no data/command line, no reset line, and no backlight.
// Only SCLK/MOSI/CS (and the DISP/EXTMODE/EXTCOMIN control lines handled in
// hardware) are required.
func New(port spi.Port, cfg driver.DriverConfig) (*Dev, error) {
	if port == nil {
		return nil, errors.New("display: spi port must not be nil")
	}
	if cfg.Width <= 0 {
		cfg.Width = sharpWidth
	}
	if cfg.Height <= 0 {
		cfg.Height = sharpHeight
	}
	if cfg.SPIHz <= 0 {
		// Sharp memory LCDs are specified around 1–2MHz SPI.
		cfg.SPIHz = 2 * physic.MegaHertz
	}
	conn, err := port.Connect(cfg.SPIHz, spi.Mode0, 8)
	if err != nil {
		return nil, err
	}
	d := &Dev{c: conn, width: cfg.Width, height: cfg.Height, vcom: bitVCOM}
	if err := d.clear(); err != nil {
		return nil, err
	}
	return d, nil
}

// Bounds returns the display rectangle.
func (d *Dev) Bounds() image.Rectangle { return image.Rect(0, 0, d.width, d.height) }

// TextHints reports this panel as a fast monochrome LCD. Unlike e-paper it
// refreshes quickly, so it is advertised as mono-fast and supports scrolling.
func (d *Dev) TextHints() textlayout.TextHints {
	b := d.Bounds()
	return textlayout.TextHints{
		PixelWidth: b.Dx(), PixelHeight: b.Dy(),
		SupportsVerticalScroll: true, SupportsHorizontalScroll: true, SupportsAutoScroll: true,
		PreferEventRefresh:     false,
		Capability:             textlayout.CapMonoFast,
		DefaultTickerDirection: textlayout.TickerDirectionVertical, DefaultLineMode: textlayout.LineModeTruncate,
	}
}

// DrawImage packs the source image into 1bpp lines and writes the whole frame.
func (d *Dev) DrawImage(src draw.Image) error {
	b := d.Bounds()
	lineBytes := (d.width + 7) / 8

	// Mode byte: write command with the current VCOM state, then toggle VCOM.
	frame := make([]byte, 0, 1+d.height*(lineBytes+2)+1)
	frame = append(frame, d.vcom|bitWriteCmd)
	d.toggleVCOM()

	for y := 0; y < d.height; y++ {
		// 1-based line address.
		frame = append(frame, byte(y+1))
		line := make([]byte, lineBytes)
		for x := 0; x < d.width; x++ {
			r, g, bl, _ := src.At(b.Min.X+x, b.Min.Y+y).RGBA()
			// White pixel (luminance above threshold) sets the bit; black clears
			// it. Leftmost pixel goes in the most-significant bit so MSB-first SPI
			// transmits it first, matching the Sharp line order.
			if driver.Luminance(r, g, bl) > driver.MonoThreshold {
				line[x/8] |= 0x80 >> uint(x%8)
			}
		}
		frame = append(frame, line...)
		frame = append(frame, 0x00) // end-of-line
	}
	frame = append(frame, 0x00) // trailing byte after the final line

	return d.tx(frame)
}

// clear issues the dedicated clear command (faster than writing a blank frame).
func (d *Dev) clear() error {
	data := []byte{d.vcom | bitClear, 0x00}
	d.toggleVCOM()
	return d.tx(data)
}

// toggleVCOM flips the VCOM inversion bit, which must alternate on every
// transfer to avoid a DC bias in the panel.
func (d *Dev) toggleVCOM() {
	if d.vcom != 0 {
		d.vcom = 0
	} else {
		d.vcom = bitVCOM
	}
}

func (d *Dev) tx(buf []byte) error {
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
