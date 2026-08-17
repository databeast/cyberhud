package adafruit_15x7_charlieplex

import (
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/hardware/driver"
	"github.com/databeast/cyberhud/hardware/panels"
)

func init() {
	panels.Register(panels.Definition{
		Name:        "adafruit-15x7-charlieplex",
		Description: "Adafruit 15x7 CharliePlex LED Matrix FeatherWing (IS31FL3731)",
		Controller:  "is31fl3731",
		Config: driver.DriverConfig{
			Width:  15,
			Height: 7,
			Layout: "charlie-wing",
		},
		// SPI pins left empty — this is an I2C device
		DCPin:   "",
		RSTPin:  "",
		BusyPin: "",
		BLPin:   "",
	})

	style.RegisterAliases("adafruit-15x7-charlieplex", map[string]string{
		"serial": "compact",
	})
}
