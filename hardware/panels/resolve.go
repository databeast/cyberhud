package panels

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/databeast/cyberhud/hardware/driver"
)

// Overrides defines optional CLI/config overrides applied to a base panel definition.
type Overrides struct {
	Width              int
	Height             int
	MADCTL             string
	Rotate             string // "true", "false", or "" (unset)
	XOffset            int
	YOffset            int
	Orientation        string            // orientation: "normal", "flip", "cw", "ccw" — applies to all screens
	ScreenOrientations map[string]string // per-screen orientations keyed by screen name
	DCPin              string
	RSTPin             string
	BLPin              string
	BusyPin            string
	InputKey1          string
	InputKey2          string
	InputKey3          string
	InputUp            string
	InputDown          string
	InputLeft          string
	InputRight         string
	InputPress         string
}

// Resolve returns one panel definition with overrides applied.
func Resolve(name string, o Overrides) (Definition, error) {
	def, err := Get(name)
	if err != nil {
		return Definition{}, err
	}

	applyPositiveInt(&def.Config.Width, o.Width)
	applyPositiveInt(&def.Config.Height, o.Height)

	// Apply orientation before raw MADCTL so that MADCTL override still wins.
	if orient, ok := trimmedValue(o.Orientation); ok {
		parsed, err := driver.ParseOrientation(orient)
		if err != nil {
			return Definition{}, err
		}
		applyOrientation(&def, parsed)
	}

	// Apply per-screen orientations.
	for screenName, orient := range o.ScreenOrientations {
		parsed, err := driver.ParseOrientation(orient)
		if err != nil {
			return Definition{}, fmt.Errorf("screen %q: %w", screenName, err)
		}
		applyScreenOrientation(&def, screenName, parsed)
	}

	if madctl, ok := trimmedValue(o.MADCTL); ok {
		parsed, err := strconv.ParseUint(madctl, 0, 8)
		if err != nil {
			return Definition{}, err
		}
		def.Config.MADCTL = byte(parsed)
	}
	if rot, ok := trimmedValue(o.Rotate); ok {
		def.Config.Rotate180 = strings.EqualFold(rot, "true") || rot == "1"
	}
	applyNonNegativeInt(&def.Config.XOffset, o.XOffset)
	applyNonNegativeInt(&def.Config.YOffset, o.YOffset)
	applyTrimmedString(&def.DCPin, o.DCPin)
	applyTrimmedString(&def.RSTPin, o.RSTPin)
	applyBLPinOverride(&def.BLPin, o.BLPin)
	applyTrimmedString(&def.BusyPin, o.BusyPin)
	applyTrimmedString(&def.Inputs.Key1, o.InputKey1)
	applyTrimmedString(&def.Inputs.Key2, o.InputKey2)
	applyTrimmedString(&def.Inputs.Key3, o.InputKey3)
	applyTrimmedString(&def.Inputs.JoyUp, o.InputUp)
	applyTrimmedString(&def.Inputs.JoyDown, o.InputDown)
	applyTrimmedString(&def.Inputs.JoyLeft, o.InputLeft)
	applyTrimmedString(&def.Inputs.JoyRight, o.InputRight)
	applyTrimmedString(&def.Inputs.JoyPressed, o.InputPress)
	return def, nil
}

// ConfigForScreen merges a screen-specific config onto a base driver config.
func ConfigForScreen(screen Screen, base driver.DriverConfig) driver.DriverConfig {
	cfg := base
	if screen.Config.Width != 0 {
		cfg.Width = screen.Config.Width
	}
	if screen.Config.Height != 0 {
		cfg.Height = screen.Config.Height
	}
	if screen.Config.SPIHz != 0 {
		cfg.SPIHz = screen.Config.SPIHz
	}
	if screen.Config.MADCTL != 0 {
		cfg.MADCTL = screen.Config.MADCTL
	}
	if screen.Config.XOffset != 0 {
		cfg.XOffset = screen.Config.XOffset
	}
	if screen.Config.YOffset != 0 {
		cfg.YOffset = screen.Config.YOffset
	}
	if screen.Config.ColOffset != 0 {
		cfg.ColOffset = screen.Config.ColOffset
	}
	if screen.Config.BusyTimeout != 0 {
		cfg.BusyTimeout = screen.Config.BusyTimeout
	}
	if screen.Config.FullRefreshEvery != 0 {
		cfg.FullRefreshEvery = screen.Config.FullRefreshEvery
	}
	if screen.Config.BusyHigh {
		cfg.BusyHigh = true
	}
	return cfg
}

func trimmedValue(s string) (string, bool) {
	s = strings.TrimSpace(s)
	return s, s != ""
}

func applyTrimmedString(dst *string, value string) {
	if v, ok := trimmedValue(value); ok {
		*dst = v
	}
}

func applyBLPinOverride(dst *string, value string) {
	if v, ok := trimmedValue(value); ok {
		if strings.EqualFold(v, "none") || v == "-" {
			*dst = ""
			return
		}
		*dst = v
	}
}

func applyPositiveInt(dst *int, value int) {
	if value > 0 {
		*dst = value
	}
}

func applyNonNegativeInt(dst *int, value int) {
	if value >= 0 {
		*dst = value
	}
}
