// Adafruit Monochrome 1.3" 128x64 SH1106G SPI OLED (product 5228): the bare OLED
// module used in the Adafruit MacroPad RP2040. It is sold as a raw display with
// no PCB or supporting circuitry — just the panel and its 16-pin 0.5mm FPC tail
// — intended to be integrated into a custom board.
//
// The controller is the SH1106 (the "G" is the module variant), driven over
// 4-wire SPI, so this panel reuses the existing "sh1106" driver. Being an OLED
// it is self-lit, with no backlight. It is a bare module with no onboard input.
//
// There is no fixed GPIO mapping: the user wires DC/RST to whichever Raspberry
// Pi GPIOs they choose and supplies them as config overrides (display.dc /
// display.rst, or the -display-dc / -display-rst flags). The pin fields are
// therefore left empty here; starting the panel without DC/RST overrides fails
// with a "pin not found" error.
package adafruit_1_3_oled_sh1106

import (
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/hardware/driver"
	"github.com/databeast/cyberhud/hardware/panels"
)

func init() {
	panels.Register(panels.Definition{
		Name:        "adafruit-1.3-oled-sh1106",
		Description: "Adafruit 1.3\" 128x64 SH1106 Monochrome OLED (bare MacroPad module, SPI)",
		Controller:  "sh1106",
		Monochrome:  true,
		PPI:         panels.CalcPPI(128, 64, 1.3),
		Config: driver.DriverConfig{
			Width:  128,
			Height: 64,
			// SH1106 has 132 columns of RAM for a 128-wide panel; the visible
			// window starts 2 columns in.
			ColOffset: 2,
		},
		// No fixed GPIO mapping — the user supplies DC/RST via config overrides
		// for their own wiring. OLED is self-lit, so there is no backlight pin.
		DCPin:  "",
		RSTPin: "",
		BLPin:  "",
	})

	style.RegisterAliases("adafruit-1.3-oled-sh1106", map[string]string{
		"dashboard": "mono-fast-128x64",
		"clock":     "mono-fast-128x64",
		"thermal":   "mono-fast-128x64",
		"ticker":    "mono-fast-128x64",
		"wifi":      "mono-fast-128x64",
		"gpio":      "mono-fast-128x64",
		"serial":    "mono-fast-128x64",
		"usb":       "mono-slow-128x64",
	})
}
