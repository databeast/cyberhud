// Adafruit SHARP Memory Display Breakout 1.3" 168x144 Monochrome (product
// 3502): a Sharp Memory LCD — a write-only 1-bit reflective display that pairs
// e-paper-like low power with LCD-like fast refresh, showing black pixels on a
// gray background. It is not backlit (daylight-readable; needs external light in
// the dark).
//
// The controller is the bespoke Sharp Memory LCD protocol (driver id
// "sharpmem"), not an ST/SSD/ILI controller. It needs only SPI (SCLK/MOSI/CS)
// plus the DISP/EXTMODE/EXTCOMIN control lines, which are handled in hardware;
// there is no data/command, reset, or backlight GPIO. All four pin fields are
// therefore left empty, and the panel is output-only.
package adafruit_sharp_1_3

import (
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/hardware/driver"
	"github.com/databeast/cyberhud/hardware/panels"
)

func init() {
	panels.Register(panels.Definition{
		Name:        "adafruit-sharp-1.3",
		Description: "Adafruit 1.3\" 168x144 Sharp Memory LCD (monochrome, black-on-gray)",
		Controller:  "sharpmem",
		Monochrome:  true,
		PPI:         panels.CalcPPI(168, 144, 1.3),
		Config: driver.DriverConfig{
			Width:  168,
			Height: 144,
		},
		// Sharp Memory LCD uses SPI only (SCLK/MOSI/CS); it has no DC, RST,
		// backlight, or busy line. The DISP/EXTMODE/EXTCOMIN control lines are
		// handled in hardware.
		DCPin:   "",
		RSTPin:  "",
		BusyPin: "",
		BLPin:   "",
	})

	// 168x144 has no exact style; the closest fitting mono-fast landscape style
	// is 160x128, which a 168x144 surface satisfies. Modes without a matching
	// alias fall through to the style system's fitness-based resolution.
	style.RegisterAliases("adafruit-sharp-1.3", map[string]string{
		"dashboard": "mono-fast-160x128",
		"thermal":   "overview",
		"wifi":      "mono-fast-160x128",
		"gpio":      "mono-fast-160x128",
	})
}
