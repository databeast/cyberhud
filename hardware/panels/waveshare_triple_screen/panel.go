package waveshare_triple_screen

import (
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/hardware/driver"
	"github.com/databeast/cyberhud/hardware/panels"
)

func init() {
	panels.Register(panels.Definition{
		Name:        "waveshare-triple-screen",
		Description: "Waveshare Zero LCD HAT A triple-screen (main 1.3in ST7789 + dual 0.96in ST7735S)",
		Controller:  "virtual",
		Virtual: []panels.Screen{
			{
				Index:       0,
				Name:        "main",
				Controller:  "st7789",
				SPI:         "SPI1.0",
				DefaultMode: "menu",
				Config:      driver.DriverConfig{Width: 240, Height: 240, MADCTL: driver.MadctlMX},
				DCPin:       panels.GPIO22,
				RSTPin:      panels.GPIO27,
				BLPin:       panels.GPIO19,
				XPosition:   80,
				YPosition:   0,
				MirrorX:     true,
				PPI:         panels.CalcPPI(240, 240, 1.3),
				Orientations: map[driver.Orientation]driver.OrientationConfig{
					driver.OrientationNormal: {MADCTL: driver.MadctlMX, MirrorX: true},
					driver.OrientationFlip:   {MADCTL: driver.MadctlMY, YOffset: 80, MirrorX: true},
					driver.OrientationCW:     {MADCTL: driver.MadctlMV | driver.MadctlMX},
					driver.OrientationCCW:    {MADCTL: driver.MadctlMV | driver.MadctlMY},
				},
			},
			{
				Index:       1,
				Name:        "left-aux",
				Controller:  "st7735s",
				SPI:         "SPI0.0",
				DefaultMode: "stemma",
				Config:      driver.DriverConfig{Width: 160, Height: 80, MADCTL: driver.MadctlMY | driver.MadctlMV | driver.MadctlBGR, XOffset: 1, YOffset: 26, InvertColors: true},
				DCPin:       panels.GPIO4,
				RSTPin:      panels.GPIO24,
				BLPin:       panels.GPIO13,
				XPosition:   0,
				YPosition:   0,
				Rotation:    90,
				PPI:         panels.CalcPPI(160, 80, 0.96),
				Orientations: map[driver.Orientation]driver.OrientationConfig{
					driver.OrientationNormal: {MADCTL: driver.MadctlMY | driver.MadctlMV | driver.MadctlBGR, XOffset: 1, YOffset: 26},
					driver.OrientationFlip:   {MADCTL: driver.MadctlMX | driver.MadctlMV | driver.MadctlBGR, XOffset: 1, YOffset: 26},
				},
			},
			{
				Index:       2,
				Name:        "right-aux",
				Controller:  "st7735s",
				SPI:         "SPI0.1",
				DefaultMode: "gpio",
				Config:      driver.DriverConfig{Width: 160, Height: 80, MADCTL: driver.MadctlMY | driver.MadctlMV | driver.MadctlBGR, XOffset: 1, YOffset: 26, InvertColors: true},
				DCPin:       panels.GPIO5,
				RSTPin:      panels.GPIO23,
				BLPin:       panels.GPIO12,
				XPosition:   320,
				YPosition:   0,
				Rotation:    90,
				PPI:         panels.CalcPPI(160, 80, 0.96),
				Orientations: map[driver.Orientation]driver.OrientationConfig{
					driver.OrientationNormal: {MADCTL: driver.MadctlMY | driver.MadctlMV | driver.MadctlBGR, XOffset: 1, YOffset: 26},
					driver.OrientationFlip:   {MADCTL: driver.MadctlMX | driver.MadctlMV | driver.MadctlBGR, XOffset: 1, YOffset: 26},
				},
			},
		},
		Inputs: panels.InputPins{Key1: panels.GPIO25, Key2: panels.GPIO26},
	})

	// Per-screen aliases (compound keys)
	style.RegisterAliases("waveshare-triple-screen/main", map[string]string{
		"dashboard": "color-240x240",
		"clock":     "color-240x240",
		"thermal":   "overview",
		"ticker":    "color-240x240",
		"wifi":      "color-240x240",
		"gpio":      "color-240x240",
		"serial":    "color-240x240",
	})
	style.RegisterAliases("waveshare-triple-screen/left-aux", map[string]string{
		"stemma":  "compact",
		"gpio":    "color-160x80",
		"thermal": "minimal",
		"ticker":  "color-160x80",
		"serial":  "color-160x80",
	})
	style.RegisterAliases("waveshare-triple-screen/right-aux", map[string]string{
		"gpio":    "color-160x80",
		"stemma":  "compact",
		"thermal": "minimal",
		"ticker":  "color-160x80",
		"serial":  "color-160x80",
	})

	// Parent-level fallback aliases
	style.RegisterAliases("waveshare-triple-screen", map[string]string{
		"clock":  "color-240x240",
		"serial": "color-240x240",
	})
}
