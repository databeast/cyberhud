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

// ToFloat64 extracts a float64 from an interface value, handling both
// float64 (native JSON number) and int/int64 conversions.
func ToFloat64(v interface{}) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case int:
		return float64(n), true
	case int64:
		return float64(n), true
	default:
		return 0, false
	}
}

// ToInt extracts an int from an interface value, handling both float64
// (JSON numbers decode as float64) and direct int types.
func ToInt(v interface{}) (int, bool) {
	switch n := v.(type) {
	case float64:
		return int(n), true
	case int:
		return n, true
	case int64:
		return int(n), true
	default:
		return 0, false
	}
}
