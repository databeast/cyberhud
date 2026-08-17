package waveshare_2_23_oled_hat

import (
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/hardware/driver"
	"github.com/databeast/cyberhud/hardware/panels"
)

func init() {
	// I2C panel variant
	panels.Register(panels.Definition{
		Name:        "waveshare-2.23-oled-hat-i2c",
		Description: "Waveshare 2.23-inch OLED HAT (128x32 SSD1305) via I2C",
		Controller:  "ssd1305",
		Monochrome:  true,
		PPI:         panels.CalcPPI(128, 32, 2.23),
		Config: driver.DriverConfig{
			Width:   128,
			Height:  32,
			I2CAddr: 0x3C,
		},
		Orientations: map[driver.Orientation]driver.OrientationConfig{
			driver.OrientationNormal: {},
			driver.OrientationFlip:   {Rotate180: true},
		},
		I2CBus: "/dev/i2c-1",
	})

	// SPI panel variant
	panels.Register(panels.Definition{
		Name:        "waveshare-2.23-oled-hat-spi",
		Description: "Waveshare 2.23-inch OLED HAT (128x32 SSD1305) via SPI",
		Controller:  "ssd1305",
		Monochrome:  true,
		PPI:         panels.CalcPPI(128, 32, 2.23),
		Config: driver.DriverConfig{
			Width:  128,
			Height: 32,
		},
		Orientations: map[driver.Orientation]driver.OrientationConfig{
			driver.OrientationNormal: {},
			driver.OrientationFlip:   {Rotate180: true},
		},
		DCPin:  panels.GPIO24,
		RSTPin: panels.GPIO25,
		BLPin:  panels.GPIO18,
	})

	style.RegisterAliases("waveshare-2.23-oled-hat-i2c", map[string]string{
		"dashboard": "mono-fast-128x32",
		"clock":     "mono-128x32",
		"thermal":   "minimal",
		"ticker":    "mono-128x32",
		"gpio":      "list",
		"serial":    "mono-fast-128x32",
	})

	style.RegisterAliases("waveshare-2.23-oled-hat-spi", map[string]string{
		"dashboard": "mono-fast-128x32",
		"clock":     "mono-128x32",
		"thermal":   "minimal",
		"ticker":    "mono-128x32",
		"gpio":      "list",
		"serial":    "mono-fast-128x32",
	})
}
