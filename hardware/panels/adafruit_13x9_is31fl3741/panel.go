// Adafruit IS31FL3741 13x9 PWM RGB LED Matrix (product 5201): 117 RGB LEDs in a
// 13x9 grid driven by an ISSI IS31FL3741 over I2C, with STEMMA QT / Qwiic
// connectors. Each LED has 8-bit PWM per channel for 24-bit colour.
//
// This is an I2C device, so it uses no GPIO control pins — the DC/RST/BUSY/BL
// fields are left empty and the panel is reached over the system I2C bus at the
// controller's default address (0x30, or a strapped address for tiling). The
// matrix has no onboard input.
package adafruit_13x9_is31fl3741

import (
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/hardware/driver"
	"github.com/databeast/cyberhud/hardware/panels"
)

func init() {
	panels.Register(panels.Definition{
		Name:        "adafruit-13x9-is31fl3741",
		Description: "Adafruit 13x9 PWM RGB LED Matrix (IS31FL3741, I2C/STEMMA QT)",
		Controller:  "is31fl3741",
		Config: driver.DriverConfig{
			Width:  13,
			Height: 9,
			Layout: "qt-13x9",
		},
		PPI: panels.CalcPPI(13, 9, 1.70), // 13x9 LEDs on 3mm pitch ≈ 36x24mm → ~1.70" diagonal grid
		// I2C device — no GPIO control pins. Bus defaults to the system I2C bus;
		// address defaults to the controller's 0x30.
		DCPin:   "",
		RSTPin:  "",
		BusyPin: "",
		BLPin:   "",
	})

	style.RegisterAliases("adafruit-13x9-is31fl3741", map[string]string{
		"serial": "compact",
	})
}
