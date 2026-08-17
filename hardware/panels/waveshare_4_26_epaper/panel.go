package waveshare_4_26_epaper

import (
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/hardware/driver"
	"github.com/databeast/cyberhud/hardware/panels"
)

func init() {
	panels.Register(panels.Definition{
		Name:        "waveshare-4.26-epaper",
		Description: "Waveshare 4.26-inch E-Paper HAT (800x480 B/W) with SPI interface",
		Controller:  "epd4in26",
		Monochrome:  true,
		PPI:         panels.CalcPPI(800, 480, 4.26),
		Config: driver.DriverConfig{
			Width:  800,
			Height: 480,
		},
		Orientations: map[driver.Orientation]driver.OrientationConfig{
			driver.OrientationNormal: {},
			driver.OrientationFlip:   {Rotate180: true},
		},
		DCPin:   panels.GPIO25,
		RSTPin:  panels.GPIO17,
		BusyPin: panels.GPIO24,
		BLPin:   "",
	})

	style.RegisterAliases("waveshare-4.26-epaper", map[string]string{
		"dashboard": "eink-800x480",
		"clock":     "eink-800x480",
		"thermal":   "overview",
		"ticker":    "eink-800x480",
		"gpio":      "eink-800x480",
		"serial":    "mono-slow-800x480",
	})
}
