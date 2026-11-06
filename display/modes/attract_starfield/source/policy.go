package source

import (
	"fmt"

	"github.com/databeast/cyberhud/display/catalog"
)

// Policy captures all runtime-configurable parameters for the starfield mode.
type Policy struct {
	Speed   float64 // [0.1, 10.0], default 1.0 — travel speed multiplier
	Density float64 // [0.1, 1.0], default 0.5 — star density scaling
	Layers  int     // [1, 8], default 4 — depth layers
}

func (p Policy) Options() []catalog.OptionDefinition {
	return []catalog.OptionDefinition{
		{Key: "speed", Type: "float", Summary: "Travel speed multiplier (0.1-10.0).", Default: "1.0", Allowed: []string{}},
		{Key: "density", Type: "float", Summary: "Star density scaling (0.1-1.0).", Default: "0.5", Allowed: []string{}},
		{Key: "layers", Type: "int", Summary: "Number of depth layers (1-8).", Default: "4", Allowed: []string{}},
	}
}

func (p Policy) Fingerprint() string {
	return fmt.Sprintf("%v|%v|%v", p.Speed, p.Density, p.Layers)
}

func (p Policy) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"speed":   p.Speed,
		"density": p.Density,
		"layers":  p.Layers,
	}
}

// DefaultPolicy returns the default starfield policy with all fields in valid range.
func DefaultPolicy() Policy {
	return Policy{
		Speed:   1.0,
		Density: 0.5,
		Layers:  4,
	}
}
