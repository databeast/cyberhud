package st7789_240x135

import (
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/hardware/driver"
	"github.com/databeast/cyberhud/hardware/panels"
)

func init() {
	panels.Register(panels.Definition{
		Name:        "st7789-240x135",
		Description: "Generic ST7789 SPI panel (240x135) without input controls",
		Controller:  "st7789",
		Config: driver.DriverConfig{
			Width:   240,
			Height:  135,
			MADCTL:  driver.MadctlRGB | driver.MadctlMV,
			XOffset: 40,
			YOffset: 53,
		},
		DCPin:  panels.GPIO25,
		RSTPin: panels.GPIO27,
		BLPin:  panels.GPIO24,
	})

	style.RegisterAliases("st7789-240x135", map[string]string{
		"dashboard": "color-240x240",
		"clock":     "color-240x135",
		"thermal":   "minimal",
		"ticker":    "color-240x135",
		"wifi":      "color-240x240",
		"gpio":      "color-240x135",
		"serial":    "color-240x135",
	})
}
