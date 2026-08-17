package waveshare_1_44

import (
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/hardware/driver"
	"github.com/databeast/cyberhud/hardware/panels"
)

func init() {
	panels.Register(panels.Definition{
		Name:        "waveshare-1.44",
		Description: "Waveshare 1.44-inch LCD HAT (128x128 ST7735S) with three buttons and joystick",
		Controller:  "st7735s",
		PPI:         panels.CalcPPI(128, 128, 1.44),
		Config: driver.DriverConfig{
			Width:   128,
			Height:  128,
			MADCTL:  driver.MadctlMY | driver.MadctlMV | driver.MadctlBGR,
			XOffset: 2,
			YOffset: 1,
		},
		Orientations: map[driver.Orientation]driver.OrientationConfig{
			driver.OrientationNormal: {MADCTL: driver.MadctlMY | driver.MadctlMV | driver.MadctlBGR, XOffset: 2, YOffset: 1},
			driver.OrientationFlip:   {MADCTL: driver.MadctlMX | driver.MadctlMV | driver.MadctlBGR, XOffset: 2, YOffset: 1},
		},
		DCPin:  panels.GPIO25,
		RSTPin: panels.GPIO27,
		BLPin:  panels.GPIO24,
		Inputs: panels.InputPins{
			Key1:       panels.GPIO21,
			Key2:       panels.GPIO20,
			Key3:       panels.GPIO16,
			JoyUp:      panels.GPIO6,
			JoyDown:    panels.GPIO19,
			JoyLeft:    panels.GPIO5,
			JoyRight:   panels.GPIO26,
			JoyPressed: panels.GPIO13,
		},
	})

	style.RegisterAliases("waveshare-1.44", map[string]string{
		"dashboard": "color-128x160",
		"clock":     "color-128x128",
		"thermal":   "minimal",
		"ticker":    "color-128x128",
		"wifi":      "color-128x128",
		"gpio":      "color-128x128",
		"serial":    "color-128x128",
	})
}
