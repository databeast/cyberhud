// Adafruit 1.54" Monochrome eInk / ePaper Display with SRAM (product 4196): a
// 200x200 black/white e-paper breakout built around an SSD1681 controller, with
// an onboard SRAM buffer, a microSD socket, and an 18-pin EYESPI FPC connector.
//
// This definition covers the DISPLAY path only. The board has no buttons, so it
// is output-only. The onboard SRAM and microSD are not used — CyberHUD buffers
// frames on the host and streams them to the controller directly.
//
// The board is a breakout, not a Pi HAT, so there is no fixed GPIO mapping: the
// user wires the board's DC/RST/BUSY lines to whichever Raspberry Pi GPIOs they
// choose and supplies them as config overrides (display.dc / display.rst /
// display.busy, or the -display-dc / -display-rst / -display-busy flags). The
// pin fields are therefore left empty here; starting the panel without pin
// overrides fails with a "pin not found" error. Native orientation is 200x200
// square.
package adafruit_1_54_ssd1681

import (
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/hardware/driver"
	"github.com/databeast/cyberhud/hardware/panels"
)

func init() {
	panels.Register(panels.Definition{
		Name:        "adafruit-1.54-ssd1681",
		Description: "Adafruit 1.54\" Monochrome eInk Breakout (200x200 SSD1681) without onboard inputs",
		Controller:  "ssd1681",
		Monochrome:  true,
		PPI:         panels.CalcPPI(200, 200, 1.54),
		Config: driver.DriverConfig{
			Width:            200,
			Height:           200,
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

	style.RegisterAliases("adafruit-1.54-ssd1681", map[string]string{
		"dashboard": "mono-slow-200x200",
		"clock":     "mono-slow-200x200",
		"thermal":   "overview",
		"gpio":      "mono-slow-200x200",
		"serial":    "mono-slow-200x200",
	})
}
