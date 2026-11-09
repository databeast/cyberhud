// Adafruit 1.5" 16-bit Color OLED Breakout with microSD (product 1431): a
// 128x128 RGB OLED built around an SSD1351 controller, driven over 4-wire SPI.
// Later revisions add an EYESPI FPC connector but are otherwise identical.
//
// This definition covers the DISPLAY path only. The onboard microSD is not
// used, and the board has no input. Being an OLED, it is self-lit — there is no
// backlight line (the BL pin stays empty).
//
// The board is a breakout, not a Pi HAT, so there is no fixed GPIO mapping: the
// user wires the board's DC/RST lines to whichever Raspberry Pi GPIOs they
// choose and supplies them as config overrides (display.dc / display.rst, or
// the -display-dc / -display-rst flags). The pin fields are therefore left
// empty here; starting the panel without DC/RST overrides fails with a "pin not
// found" error. Native resolution is 128x128.
package adafruit_1_5_oled_ssd1351

import (
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/hardware/driver"
	"github.com/databeast/cyberhud/hardware/panels"
)

func init() {
	panels.Register(panels.Definition{
		Name:        "adafruit-1.5-oled-ssd1351",
		Description: "Adafruit 1.5\" 128x128 RGB Color OLED (SSD1351 SPI)",
		Controller:  "ssd1351",
		PPI:         panels.CalcPPI(128, 128, 1.5),
		Config: driver.DriverConfig{
			Width:  128,
			Height: 128,
		},
		// No fixed GPIO mapping — the user supplies DC/RST via config overrides
		// for their own wiring. OLED is self-lit, so there is no backlight pin.
		DCPin:  "",
		RSTPin: "",
		BLPin:  "",
	})

	style.RegisterAliases("adafruit-1.5-oled-ssd1351", map[string]string{
		"dashboard": "color-fast-128x128",
		"clock":     "color-fast-128x128",
		"thermal":   "color-fast-128x128-overview",
		"ticker":    "color-fast-128x128",
		"wifi":      "color-fast-128x128",
		"gpio":      "color-fast-128x128",
		"serial":    "color-fast-128x128",
	})
}
