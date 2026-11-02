package waveshare_3_5_lcd_a

import (
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/hardware/driver"
	"github.com/databeast/cyberhud/hardware/panels"
)

// Waveshare 3.5" RPi LCD (A): a 480x320 SPI TFT built around an ILI9486 display
// controller and an XPT2046 resistive touch controller. This definition covers
// the DISPLAY path only — touch (XPT2046) input is not yet supported, so the
// panel is registered output-only (InputEnabled stays false).
//
// Pin mapping follows the standard `waveshare35a` device-tree overlay used by
// this board: display on SPI0.0 (CE0), DC=GPIO24, RST=GPIO25, backlight tied on
// in hardware (no GPIO control). Native orientation is 480x320 landscape.
func init() {
	panels.Register(panels.Definition{
		Name:        "waveshare-3.5-lcd-a",
		Description: "Waveshare 3.5\" RPi LCD (A), ILI9486 480x320 SPI TFT (resistive touch not yet wired)",
		Controller:  "ili9486",
		Config: driver.DriverConfig{
			Width:  480,
			Height: 320,
			// 0x28: MV set (row/column exchange → landscape) with BGR order,
			// the common configuration for this board. Orientation/colour order
			// may need on-device adjustment across board revisions.
			MADCTL: 0x28,
		},
		DCPin:  panels.GPIO24,
		RSTPin: panels.GPIO25,
		// Backlight is hardwired on for the (A) revision — no BLPin.
	})

	style.RegisterAliases("waveshare-3.5-lcd-a", map[string]string{
		"dashboard":          "color-480x320",
		"clock":              "color-480x320",
		"thermal":            "overview",
		"ticker":             "color-480x320",
		"wifi":               "color-480x320",
		"gpio":               "color-480x320",
		"serial":             "color-480x320",
		"attract_popsquares": "color-480x320",
		"attract_waveform":   "color-480x320",
		"attract_plasma":     "color-480x320",
	})
}
