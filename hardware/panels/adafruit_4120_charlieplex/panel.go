package adafruit_4120_charlieplex

import (
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/hardware/driver"
	"github.com/databeast/cyberhud/hardware/panels"
)

func init() {
	panels.Register(panels.Definition{
		Name:        "adafruit-4120-charlieplex",
		Description: "Adafruit 16x8 CharliePlex LED Matrix Bonnet (IS31FL3731, Green)",
		Controller:  "is31fl3731",
		Config: driver.DriverConfig{
			Width:  16,
			Height: 8,
			Layout: "charlie-bonnet",
		},
		Orientations: map[driver.Orientation]driver.OrientationConfig{
			driver.OrientationNormal: {Width: 16, Height: 8},
			// 90° clockwise mount: portrait 8x16; the is31fl3731 driver rotates
			// coordinates into the bonnet's native landscape wiring.
			driver.OrientationCW: {Width: 8, Height: 16},
		},
		PPI: panels.CalcPPI(16, 8, 0.8),
		// SPI pins left empty — this is an I2C device
		DCPin:   "",
		RSTPin:  "",
		BusyPin: "",
		BLPin:   "",
	})

	style.RegisterAliases("adafruit-4120-charlieplex", map[string]string{
		"serial": "compact",
	})
}
