// Adafruit Monochrome 2.7" 128x64 OLED Graphic Display Module (product 2674): a
// 128x64 white OLED driven by an SSD1325 controller over SPI. The glass is a
// single colour, but the SSD1325 drives it with 4-bit (16-level) grayscale, so
// this panel is a grayscale panel, not 1-bit mono.
//
// This definition covers the DISPLAY path only; being an OLED it is self-lit
// with no backlight, and the module has no input.
//
// The module is a breakout, not a Pi HAT, so there is no fixed GPIO mapping: the
// user wires DC/RST to whichever Raspberry Pi GPIOs they choose and supplies
// them as config overrides (display.dc / display.rst, or the -display-dc /
// -display-rst flags). The pin fields are therefore left empty here; starting
// the panel without DC/RST overrides fails with a "pin not found" error.
package adafruit_2_7_oled_ssd1325

import (
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/hardware/driver"
	"github.com/databeast/cyberhud/hardware/panels"
)

func init() {
	panels.Register(panels.Definition{
		Name:        "adafruit-2.7-oled-ssd1325",
		Description: "Adafruit 2.7\" 128x64 Grayscale OLED (SSD1325 SPI)",
		Controller:  "ssd1325",
		PPI:         panels.CalcPPI(128, 64, 2.7),
		Config: driver.DriverConfig{
			Width:  128,
			Height: 64,
		},
		// No fixed GPIO mapping — the user supplies DC/RST via config overrides.
		// OLED is self-lit, so there is no backlight pin.
		DCPin:  "",
		RSTPin: "",
		BLPin:  "",
	})

	// Alias only the modes that register a grayscale-fast-128x64 style; other
	// modes fall through to fitness-based resolution (a grayscale-fast panel
	// also satisfies mono-fast styles).
	style.RegisterAliases("adafruit-2.7-oled-ssd1325", map[string]string{
		"dashboard": "grayscale-fast-128x64",
		"thermal":   "overview",
		"wifi":      "grayscale-fast-128x64",
		"gpio":      "grayscale-fast-128x64",
		"usb":       "grayscale-fast-128x64",
	})
}
