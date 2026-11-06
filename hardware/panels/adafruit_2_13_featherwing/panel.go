// Adafruit 2.13" Monochrome eInk / ePaper Display FeatherWing (product 4195): a
// 250x122 black/white e-paper board built around an SSD1680 controller, with an
// onboard SRAM buffer and a microSD socket. It shares the same display and
// controller as the 2.13" e-ink Bonnet, but is a Feather-form add-on with no
// buttons.
//
// This definition covers the DISPLAY path only and is output-only (the
// FeatherWing has no onboard input). The onboard SRAM and microSD are not used
// — CyberHUD buffers frames on the host and streams them to the controller.
//
// The FeatherWing is driven from the Raspberry Pi's SPI bus plus DC/RST/BUSY
// control lines, but there is no fixed GPIO mapping: the user wires those lines
// to whichever Raspberry Pi GPIOs they choose and supplies them as config
// overrides (display.dc / display.rst / display.busy, or the -display-dc /
// -display-rst / -display-busy flags). The pin fields are therefore left empty
// here; starting the panel without pin overrides fails with a "pin not found"
// error. Native orientation is 250x122 landscape. The SSD1680 driver's RAM
// layout rotates the 250-wide logical frame into the controller's native
// 122-wide portrait RAM automatically.
package adafruit_2_13_featherwing

import (
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/hardware/driver"
	"github.com/databeast/cyberhud/hardware/panels"
)

func init() {
	panels.Register(panels.Definition{
		Name:        "adafruit-2.13-featherwing",
		Description: "Adafruit 2.13\" Monochrome eInk FeatherWing (250x122 SSD1680) without onboard inputs",
		Controller:  "ssd1680",
		Monochrome:  true,
		PPI:         panels.CalcPPI(250, 122, 2.13),
		Config: driver.DriverConfig{
			Width:            250,
			Height:           122,
			BusyHigh:         true,
			FullRefreshEvery: 20,
		},
		// No fixed GPIO mapping — the user supplies DC/RST/BUSY via config
		// overrides for their own wiring. Backlight is not used on e-paper.
		DCPin:   "",
		RSTPin:  "",
		BusyPin: "",
		BLPin:   "",
	})

	style.RegisterAliases("adafruit-2.13-featherwing", map[string]string{
		"dashboard": "mono-slow-250x122",
		"clock":     "mono-slow-240x135",
		"thermal":   "overview",
		"ticker":    "mono-slow-250x122",
		"gpio":      "mono-slow-250x122",
		"serial":    "mono-slow-250x122",
	})
}
