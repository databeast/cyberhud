// Adafruit 128x64 OLED Bonnet for Raspberry Pi (product 3531): a 128x64
// monochrome SSD1306 OLED in a Pi HAT/bonnet form factor, with a 5-way joystick
// and two pushbuttons. The display speaks I2C (default address 0x3C); the
// joystick and buttons are on fixed Pi GPIOs.
//
// Like the 1.3" TFT Bonnet, this is a Pi bonnet with a FIXED pinout, so the
// input GPIOs are hardcoded rather than user-supplied. Being an I2C OLED, the
// display has no DC/RST/SPI control pins — those fields stay empty and the panel
// is reached over the system I2C bus.
//
// Input pinout (standard Adafruit bonnet layout):
//   - Buttons: A=GPIO5, B=GPIO6 (active-low, pulled up).
//   - Joystick: up=GPIO17, down=GPIO22, left=GPIO27, right=GPIO23, center=GPIO4.
package adafruit_128x64_oled_bonnet

import (
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/hardware/driver"
	"github.com/databeast/cyberhud/hardware/panels"
)

func init() {
	panels.Register(panels.Definition{
		Name:        "adafruit-128x64-oled-bonnet",
		Description: "Adafruit 128x64 OLED Bonnet (SSD1306 I2C) with joystick and two buttons",
		Controller:  "ssd1306",
		Monochrome:  true,
		PPI:         panels.CalcPPI(128, 64, 1.3),
		Config: driver.DriverConfig{
			Width:  128,
			Height: 64,
		},
		// I2C OLED — no DC/RST/SPI control pins. Bus defaults to the system I2C
		// bus; address defaults to the controller's 0x3C.
		DCPin:   "",
		RSTPin:  "",
		BusyPin: "",
		BLPin:   "",
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

	style.RegisterAliases("adafruit-128x64-oled-bonnet", map[string]string{
		"dashboard": "mono-fast-128x64",
		"clock":     "mono-fast-128x64",
		"thermal":   "mono-fast-128x64",
		"ticker":    "mono-fast-128x64",
		"wifi":      "mono-fast-128x64",
		"gpio":      "list",
		"serial":    "mono-fast-128x64",
		"usb":       "mono-slow-128x64",
	})
}
