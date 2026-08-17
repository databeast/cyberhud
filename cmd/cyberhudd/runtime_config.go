package main

import (
	"time"
)

// buildRuntimeConfig returns a closure that captures the effective runtime
// configuration (after flag/config merge) and produces a *fileConfig suitable
// for JSON serialization via the "config dump" protocol command.
func buildRuntimeConfig(
	socketPath string,
	i2cBuses string,
	scanInterval time.Duration,
	noDisplay bool,
	noInput bool,
	profileName string,
	width, height int,
	madctl string,
	xOffset, yOffset int,
	dc, rst, bl, busy string,
	key1, key2, key3 string,
	up, down, left, right, press string,
) func() *fileConfig {
	return func() *fileConfig {
		cfg := &fileConfig{
			Socket: socketPath,
			I2C:    i2cBuses,
			Scan:   scanInterval.String(),
			Display: fileDisplayConfig{
				Profile:    profileName,
				MADCTL:     madctl,
				DC:         dc,
				RST:        rst,
				BL:         bl,
				Busy:       busy,
				InputKey1:  key1,
				InputKey2:  key2,
				InputKey3:  key3,
				InputUp:    up,
				InputDown:  down,
				InputLeft:  left,
				InputRight: right,
				InputPress: press,
			},
		}

		// *bool: only set when true (false is the default/zero)
		if noDisplay {
			cfg.Display.Disabled = ptrBool(true)
		}
		if noInput {
			cfg.Display.DisableInput = ptrBool(true)
		}

		// *int: -1 is the sentinel for "not set"; any other value (including 0) is valid
		if width != -1 {
			cfg.Display.Width = ptrInt(width)
		}
		if height != -1 {
			cfg.Display.Height = ptrInt(height)
		}
		if xOffset != -1 {
			cfg.Display.XOffset = ptrInt(xOffset)
		}
		if yOffset != -1 {
			cfg.Display.YOffset = ptrInt(yOffset)
		}

		return cfg
	}
}

func ptrBool(v bool) *bool { return &v }
func ptrInt(v int) *int    { return &v }
