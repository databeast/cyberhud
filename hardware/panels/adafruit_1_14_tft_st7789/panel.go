// Adafruit 1.14" 240x135 Color "Newxie" TFT Display (product 6113): a tiny
// 240x135 IPS colour TFT (ST7789, SPI) on a minimal breakout designed for
// vertical "nixie-tube" stacking. It has a backlight but — unlike the larger
// 1.14" breakout — no microSD and no reset line broken out; only 7 pins are
// needed (SPI + DC + backlight + power/ground).
//
// Reuses the existing "st7789" driver. Native orientation is 240x135 landscape,
// using the standard ST7789 240x135 framebuffer window offsets.
//
// The board is a breakout, not a Pi HAT, so there is no fixed GPIO mapping: the
// user wires DC (and the backlight) to whichever Raspberry Pi GPIOs they choose
// and supplies them as config overrides (display.dc / display.bl, or the
// -display-dc / -display-bl flags). This board does NOT break out a reset line,
// but the st7789 driver requires a reset pin to be configured; point display.rst
// at any spare GPIO (it is harmless and unused — the ST7789's SPI software reset
// initialises the controller). Starting the panel without DC/RST overrides fails
// with a "pin not found" error.
package adafruit_1_14_tft_st7789

import (
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/hardware/driver"
	"github.com/databeast/cyberhud/hardware/panels"
)

func init() {
	panels.Register(panels.Definition{
		Name:        "adafruit-1.14-tft",
		Description: "Adafruit 1.14\" 240x135 Color TFT (ST7789 SPI, Newxie)",
		Controller:  "st7789",
		PPI:         panels.CalcPPI(240, 135, 1.14),
		Config: driver.DriverConfig{
			Width:   240,
			Height:  135,
			MADCTL:  driver.MadctlRGB | driver.MadctlMV,
			XOffset: 40,
			YOffset: 53,
		},
		// No fixed GPIO mapping — the user supplies DC (and optionally the
		// backlight) via config overrides. This board has no reset line broken
		// out; point RST at any spare GPIO to satisfy the driver (SWRESET handles
		// the actual reset).
		DCPin:  "",
		RSTPin: "",
		BLPin:  "",
	})

	style.RegisterAliases("adafruit-1.14-tft", map[string]string{
		"dashboard": "color-fast-240x135",
		"clock":     "color-fast-240x135",
		"thermal":   "minimal",
		"ticker":    "color-fast-240x135",
		"wifi":      "color-fast-240x135",
		"gpio":      "color-fast-240x135",
		"serial":    "color-fast-240x135",
	})
}
