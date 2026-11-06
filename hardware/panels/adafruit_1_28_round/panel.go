// Adafruit 1.28" 240x240 Round TFT LCD with MicroSD (product 6178): a round IPS
// SPI TFT built around a GC9A01A controller, 240x240 RGB565. The panel is
// physically round; the framebuffer is a full 240x240 square whose corner
// pixels are hidden behind the bezel.
//
// This definition covers the DISPLAY path only. The onboard microSD is not
// used, and the board has no input.
//
// The board is a breakout (with an EYESPI FPC connector), not a Pi HAT, so
// there is no fixed GPIO mapping: the user wires the board's DC/RST lines (and
// the backlight, if driven from a GPIO) to whichever Raspberry Pi GPIOs they
// choose and supplies them as config overrides (display.dc / display.rst /
// display.bl, or the -display-dc / -display-rst / -display-bl flags). The pin
// fields are therefore left empty here; starting the panel without DC/RST
// overrides fails with a "pin not found" error. Native resolution is 240x240.
package adafruit_1_28_round

import (
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/hardware/driver"
	"github.com/databeast/cyberhud/hardware/panels"
)

func init() {
	panels.Register(panels.Definition{
		Name:        "adafruit-1.28-round",
		Description: "Adafruit 1.28\" 240x240 Round IPS TFT (GC9A01A SPI)",
		Controller:  "gc9a01a",
		PPI:         panels.CalcPPI(240, 240, 1.28),
		Config: driver.DriverConfig{
			Width:  240,
			Height: 240,
			// GC9A01A round-panel default MADCTL (vertical scan). May need
			// on-device adjustment for mounting orientation via -display-madctl.
			MADCTL: 0x08,
		},
		// No fixed GPIO mapping — the user supplies DC/RST (and optionally the
		// backlight) via config overrides for their own wiring.
		DCPin:  "",
		RSTPin: "",
		BLPin:  "",
	})

	style.RegisterAliases("adafruit-1.28-round", map[string]string{
		"dashboard": "color-fast-240x240",
		"clock":     "color-fast-240x240",
		"thermal":   "overview",
		"ticker":    "color-fast-240x240",
		"wifi":      "color-fast-240x240",
		"gpio":      "color-fast-240x240",
		"serial":    "color-fast-240x240",
	})
}
