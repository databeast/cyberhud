package waveshare_1_3_oled_hat

import (
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/hardware/driver"
	"github.com/databeast/cyberhud/hardware/panels"
)

func init() {
	panels.Register(panels.Definition{
		Name:        "waveshare-1.3-oled-hat",
		Description: "Waveshare 1.3-inch OLED HAT (128x64 SH1106) with three buttons and D-pad",
		Controller:  "sh1106",
		Monochrome:  true,
		PPI:         panels.CalcPPI(128, 64, 1.3),
		Config: driver.DriverConfig{
			Width:     128,
			Height:    64,
			ColOffset: 2,
		},
		DCPin:  panels.GPIO24,
		RSTPin: panels.GPIO25,
		BLPin:  "",
		Inputs: panels.InputPins{
			Key1:       panels.GPIO21,
			Key2:       panels.GPIO20,
			Key3:       panels.GPIO16,
			JoyUp:      panels.GPIO19,
			JoyDown:    panels.GPIO6,
			JoyLeft:    panels.GPIO5,
			JoyRight:   panels.GPIO26,
			JoyPressed: panels.GPIO13,
		},
	})

	style.RegisterAliases("waveshare-1.3-oled-hat", map[string]string{
		"dashboard": "mono-fast-128x64",
		"clock":     "mono-128x64",
		"thermal":   "mono-oled-compact",
		"ticker":    "mono-128x64",
		"wifi":      "mono-fast-128x64",
		"gpio":      "list",
		"serial":    "mono-fast-128x64",
		"usb":       "mono-slow-128x64",
	})
}
