// Adafruit 2.4" TFT LCD with Touchscreen Breakout (product 2478): a 240x320
// portrait SPI TFT built around an ILI9341 display controller, with a resistive
// touchscreen and a microSD socket on the same board.
//
// This definition covers the DISPLAY path only — the resistive touch overlay is
// not yet supported, so the panel is registered output-only (InputEnabled stays
// false).
//
// The board is a breakout, not a Pi HAT, so there is no fixed GPIO mapping: the
// user wires the board's DC/RST lines (and the LITE backlight line, if driven
// from a GPIO) to whichever Raspberry Pi GPIOs they choose and supplies them as
// config overrides (display.dc / display.rst / display.bl, or the -display-dc /
// -display-rst / -display-bl flags). The pin fields are therefore left empty
// here; starting the panel without DC/RST overrides fails with a "pin not
// found" error. Native orientation is 240x320 portrait.
package adafruit_2_4_ili9341

import (
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/hardware/driver"
	"github.com/databeast/cyberhud/hardware/panels"
)

func init() {
	panels.Register(panels.Definition{
		Name:        "adafruit-2.4-ili9341",
		Description: "Adafruit 2.4\" TFT LCD Touchscreen Breakout, ILI9341 240x320 SPI TFT (resistive touch not yet wired)",
		Controller:  "ili9341",
		PPI:         panels.CalcPPI(240, 320, 2.4),
		Config: driver.DriverConfig{
			Width:  240,
			Height: 320,
			// 0x48: MX set with BGR order — native 240x320 portrait, the common
			// Adafruit configuration. Orientation/colour order may need on-device
			// adjustment across board revisions (tune via -display-madctl).
			MADCTL: driver.MadctlMX | driver.MadctlBGR,
		},
		// No fixed GPIO mapping — the user supplies DC/RST (and optionally the
		// backlight) via config overrides for their own wiring.
		DCPin:  "",
		RSTPin: "",
		BLPin:  "",
	})

	style.RegisterAliases("adafruit-2.4-ili9341", map[string]string{
		"dashboard": "color-fast-240x320",
		"clock":     "color-fast-240x320",
		"thermal":   "color-fast-240x320-overview",
		"ticker":    "color-fast-240x320",
		"wifi":      "color-fast-240x320",
		"gpio":      "color-fast-240x320",
		"serial":    "color-fast-240x320",
	})
}
