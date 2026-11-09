// Adafruit SHARP Memory Display Breakout 2.7" 400x240 Monochrome (product
// 4694): the larger sibling of the 1.3" Sharp Memory LCD. Same write-only 1-bit
// reflective technology (e-paper-like low power, LCD-like fast refresh,
// black-on-gray, no backlight), same bespoke Sharp Memory LCD protocol, just a
// bigger 400x240 panel.
//
// Reuses the "sharpmem" driver. The controller needs only SPI (SCLK/MOSI/CS)
// plus the DISP/EXTMODE/EXTCOMIN control lines handled in hardware; there is no
// data/command, reset, or backlight GPIO, so all four pin fields are left empty
// and the panel is output-only.
package adafruit_sharp_2_7

import (
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/hardware/driver"
	"github.com/databeast/cyberhud/hardware/panels"
)

func init() {
	panels.Register(panels.Definition{
		Name:        "adafruit-sharp-2.7",
		Description: "Adafruit 2.7\" 400x240 Sharp Memory LCD (monochrome, black-on-gray)",
		Controller:  "sharpmem",
		Monochrome:  true,
		PPI:         panels.CalcPPI(400, 240, 2.7),
		Config: driver.DriverConfig{
			Width:  400,
			Height: 240,
		},
		// Sharp Memory LCD uses SPI only (SCLK/MOSI/CS); it has no DC, RST,
		// backlight, or busy line. The DISP/EXTMODE/EXTCOMIN control lines are
		// handled in hardware.
		DCPin:   "",
		RSTPin:  "",
		BusyPin: "",
		BLPin:   "",
	})

	// Dashboard, WiFi and GPIO use the closest fitting 320x240 styles.
	style.RegisterAliases("adafruit-sharp-2.7", map[string]string{
		"dashboard": "mono-fast-320x240",
		"thermal":   "mono-fast-400x240-overview",
		"wifi":      "mono-fast-320x240",
		"gpio":      "mono-fast-320x240",
	})
}
