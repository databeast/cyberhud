package st7789_320x240

import (
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/hardware/driver"
	"github.com/databeast/cyberhud/hardware/panels"
)

func init() {
	panels.Register(panels.Definition{
		Name:        "st7789-320x240",
		Description: "Generic ST7789 SPI panel (320x240) without input controls",
		Controller:  "st7789",
		Config: driver.DriverConfig{
			Width:  320,
			Height: 240,
			MADCTL: driver.MadctlRGB | driver.MadctlMV,
		},
		DCPin:  panels.GPIO25,
		RSTPin: panels.GPIO27,
		BLPin:  panels.GPIO24,
	})

	style.RegisterAliases("st7789-320x240", map[string]string{
		"dashboard": "color-320x240",
		"clock":     "color-320x240",
		"thermal":   "overview",
		"ticker":    "color-320x240",
		"wifi":      "color-320x240",
		"gpio":      "color-320x240",
		"serial":    "color-320x240",
	})
}
