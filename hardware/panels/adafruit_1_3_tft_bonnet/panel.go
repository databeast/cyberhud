// Adafruit 1.3" Color TFT Bonnet for Raspberry Pi (product 4506): a 240x240 IPS
// TFT (ST7789, SPI) in a Pi HAT/bonnet form factor, with a 5-way joystick and
// two front buttons. Unlike the bare Adafruit breakouts, this is a Pi bonnet
// with a FIXED pinout, so the GPIO assignments are hardcoded here rather than
// supplied by the user.
//
// Pinout (from the Adafruit pinout guide):
//   - Display: SPI0 (SCK/MOSI/CE0), DC=GPIO25, backlight=GPIO26. There is no
//     reset GPIO on this bonnet — reset is not broken out — but the ST7789
//     driver requires a non-nil RST pin. GPIO24 is explicitly listed as unused
//     on this board, so it is assigned as a harmless RST stand-in: the pin is
//     not wired to the panel, and the driver's SPI software-reset (SWRESET)
//     still initialises the controller.
//   - Buttons: A=GPIO5, B=GPIO6 (active-low, pulled up).
//   - Joystick: up=GPIO17, down=GPIO22, left=GPIO27, right=GPIO23, center=GPIO4.
package adafruit_1_3_tft_bonnet

import (
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/hardware/driver"
	"github.com/databeast/cyberhud/hardware/panels"
)

func init() {
	panels.Register(panels.Definition{
		Name:        "adafruit-1.3-tft-bonnet",
		Description: "Adafruit 1.3\" Color TFT Bonnet (240x240 ST7789) with joystick and two buttons",
		Controller:  "st7789",
		PPI:         panels.CalcPPI(240, 240, 1.3),
		Config: driver.DriverConfig{
			Width:  240,
			Height: 240,
			MADCTL: driver.MadctlRGB | driver.MadctlMX,
		},
		DCPin:  panels.GPIO25,
		RSTPin: panels.GPIO24, // not wired on the bonnet; unused GPIO stand-in (SWRESET handles reset)
		BLPin:  panels.GPIO26,
		Inputs: panels.InputPins{
			Key1:       panels.GPIO5,
			Key2:       panels.GPIO6,
			JoyUp:      panels.GPIO17,
			JoyDown:    panels.GPIO22,
			JoyLeft:    panels.GPIO27,
			JoyRight:   panels.GPIO23,
			JoyPressed: panels.GPIO4,
		},
	})

	style.RegisterAliases("adafruit-1.3-tft-bonnet", map[string]string{
		"dashboard": "color-fast-240x240",
		"clock":     "color-fast-240x240",
		"thermal":   "overview",
		"ticker":    "color-fast-240x240",
		"wifi":      "color-fast-240x240",
		"gpio":      "color-fast-240x240",
		"serial":    "color-fast-240x240",
	})
}
