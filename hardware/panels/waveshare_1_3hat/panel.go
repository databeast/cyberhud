package waveshare_1_3hat

import (
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/hardware/driver"
	"github.com/databeast/cyberhud/hardware/panels"
)

func init() {
	panels.Register(panels.Definition{
		Name:        "waveshare-1.3hat",
		Description: "Waveshare 1.3-inch LCD HAT (240x240 ST7789) with buttons and joystick",
		Controller:  "st7789",
		PPI:         panels.CalcPPI(240, 240, 1.3),
		Config: driver.DriverConfig{
			Width:  240,
			Height: 240,
			MADCTL: driver.MadctlRGB | driver.MadctlMX,
		},
		DCPin:  panels.GPIO25,
		RSTPin: panels.GPIO27,
		BLPin:  panels.GPIO24,
		Inputs: panels.InputPins{
			Key1:       panels.GPIO5,
			Key2:       panels.GPIO6,
			Key3:       panels.GPIO13,
			JoyUp:      panels.GPIO19,
			JoyDown:    panels.GPIO21,
			JoyLeft:    panels.GPIO16,
			JoyRight:   panels.GPIO20,
			JoyPressed: panels.GPIO26,
		},
	})

	style.RegisterAliases("waveshare-1.3hat", map[string]string{
		"dashboard": "color-240x240",
		"clock":     "color-240x240",
		"thermal":   "overview",
		"ticker":    "color-240x240",
		"wifi":      "color-240x240",
		"gpio":      "color-240x240",
		"serial":    "color-240x240",
	})
}
