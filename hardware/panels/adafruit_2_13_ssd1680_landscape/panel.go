package adafruit_2_13_ssd1680_landscape

import (
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/hardware/driver"
	"github.com/databeast/cyberhud/hardware/panels"
)

func init() {
	panels.Register(panels.Definition{
		Name:        "adafruit-2.13-ssd1680",
		Description: "Adafruit 2.13-inch Monochrome E-Ink Bonnet (122x250 SSD1680) with two buttons",
		Controller:  "ssd1680",
		Monochrome:  true,
		PPI:         panels.CalcPPI(122, 250, 2.13),
		Config: driver.DriverConfig{
			Width:            122,
			Height:           250,
			BusyHigh:         true,
			FullRefreshEvery: 20,
		},
		Orientations: map[driver.Orientation]driver.OrientationConfig{
			driver.OrientationNormal: {},
			driver.OrientationFlip:   {Width: 122, Height: 250, Rotate180: true},
			driver.OrientationCW:     {Width: 250, Height: 122},
			driver.OrientationCCW:    {Width: 250, Height: 122, Rotate180: true},
		},
		DCPin:   panels.GPIO22,
		RSTPin:  panels.GPIO27,
		BusyPin: panels.GPIO17,
		BLPin:   "",
		Inputs: panels.InputPins{
			Key1: panels.GPIO5,
			Key2: panels.GPIO6,
		},
	})

	style.RegisterAliases("adafruit-2.13-ssd1680", map[string]string{
		"dashboard": "eink-250x122",
		"clock":     "mono-slow-240x135",
		"thermal":   "overview",
		"ticker":    "mono-slow-122x250",
		"gpio":      "eink-250x122",
		"serial":    "mono-slow-122x250",
	})
}
