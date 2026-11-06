// Adafruit 1.8" Color TFT LCD with microSD Breakout (product 358): a 128x160
// RGB TFT built around an ST7735R controller, driven over 4-wire SPI, with a
// backlight and an onboard microSD slot. Later revisions add an EYESPI FPC
// connector.
//
// The ST7735R shares the ST7735 "red tab" power/gamma initialisation and
// command set with the ST7735S, so this panel reuses the existing "st7735s"
// driver. Native orientation is 128x160 portrait.
//
// This definition covers the DISPLAY path only; the microSD is not used and the
// board has no input.
//
// The board is a breakout, not a Pi HAT, so there is no fixed GPIO mapping: the
// user wires the board's DC/RST lines (and the LITE backlight line) to whichever
// Raspberry Pi GPIOs they choose and supplies them as config overrides
// (display.dc / display.rst / display.bl, or the -display-dc / -display-rst /
// -display-bl flags). The pin fields are therefore left empty here; starting the
// panel without DC/RST overrides fails with a "pin not found" error.
package adafruit_1_8_tft_st7735r

import (
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/hardware/driver"
	"github.com/databeast/cyberhud/hardware/panels"
)

func init() {
	panels.Register(panels.Definition{
		Name:        "adafruit-1.8-tft",
		Description: "Adafruit 1.8\" 128x160 Color TFT (ST7735R SPI)",
		Controller:  "st7735s",
		PPI:         panels.CalcPPI(128, 160, 1.8),
		Config: driver.DriverConfig{
			Width:  128,
			Height: 160,
			// ST7735 red-tab framebuffer margins. These 1.8" panels vary slightly
			// by batch; if the image is shifted, tune with display.x_offset /
			// display.y_offset.
			XOffset: 2,
			YOffset: 1,
		},
		// No fixed GPIO mapping — the user supplies DC/RST (and optionally the
		// backlight) via config overrides for their own wiring.
		DCPin:  "",
		RSTPin: "",
		BLPin:  "",
	})

	style.RegisterAliases("adafruit-1.8-tft", map[string]string{
		"dashboard": "color-fast-128x160",
		"clock":     "color-fast-128x160",
		"thermal":   "overview",
		"ticker":    "color-fast-128x160",
		"wifi":      "color-fast-128x160",
		"gpio":      "color-fast-128x160",
		"serial":    "color-fast-128x160",
	})
}
