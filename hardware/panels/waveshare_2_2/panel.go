package waveshare_2_2

import (
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/hardware/driver"
	"github.com/databeast/cyberhud/hardware/panels"
)

func init() {
	panels.Register(panels.Definition{
		Name:        "waveshare-2.2",
		Description: "Waveshare 2.2-inch SPI display (320x240 ST7789) without onboard inputs",
		Controller:  "st7789",
		PPI:         panels.CalcPPI(320, 240, 2.2),
		Config: driver.DriverConfig{
			Width:  320,
			Height: 240,
			MADCTL: driver.MadctlRGB | driver.MadctlMV,
		},
		DCPin:  panels.GPIO25,
		RSTPin: panels.GPIO27,
		BLPin:  panels.GPIO24,
	})

	style.RegisterAliases("waveshare-2.2", map[string]string{
		"dashboard": "color-320x240",
		"clock":     "color-320x240",
		"thermal":   "overview",
		"ticker":    "color-320x240",
		"wifi":      "color-320x240",
		"gpio":      "color-320x240",
		"serial":    "color-320x240",
	})
}
