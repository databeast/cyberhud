package source

import (
	"fmt"

	"github.com/databeast/cyberhud/display/catalog"
)

// Policy captures all runtime-configurable parameters for the shapes attract mode.
// Fields are exposed through catalog registration and CLI command handling via
// the "attract_shapes" command verb.
type Policy struct {
	Speed      float64 // rotation speed multiplier, [0.1, 10.0], default 1.0
	Density    float64 // visual density (unused for shapes, kept for uniformity), [0.1, 1.0], default 0.5
	ShapeCount int     // number of shapes to display, [1, 50], default 8
	PulseRate  float64 // scale oscillation rate in oscillations/second, [0.1, 5.0], default 1.0
	Complexity int     // number of polygon sides, [3, 8], default 6
}

func (p Policy) Fingerprint() string {
	return policyFingerprint(p)
}

func (p Policy) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"speed":       p.Speed,
		"density":     p.Density,
		"shape_count": p.ShapeCount,
		"pulse_rate":  p.PulseRate,
		"complexity":  p.Complexity,
	}
}

func (p Policy) Options() []catalog.OptionDefinition {
	return []catalog.OptionDefinition{
		{Key: "speed", Type: "float", Summary: "Rotation speed multiplier.", Default: "1.0", Allowed: []string{}},
		{Key: "density", Type: "float", Summary: "Visual density (0.1 to 1.0).", Default: "0.5", Allowed: []string{}},
		{Key: "shape_count", Type: "int", Summary: "Number of shapes to display.", Default: "8", Allowed: []string{}},
		{Key: "pulse_rate", Type: "float", Summary: "Scale oscillation rate in oscillations/second.", Default: "1.0", Allowed: []string{}},
		{Key: "complexity", Type: "int", Summary: "Number of polygon sides (3-8).", Default: "6", Allowed: []string{}},
	}
}

// DefaultPolicy returns the default shapes policy with all fields initialized
// to valid values that do not trigger normalization.
func DefaultPolicy() Policy {
	return Policy{
		Speed:      1.0,
		Density:    0.5,
		ShapeCount: 8,
		PulseRate:  1.0,
		Complexity: 6,
	}
}

// policyFingerprint returns a pipe-delimited string representation of all
// policy fields, used for change detection in RenderCacheKey.
func policyFingerprint(p Policy) string {
	return fmt.Sprintf("%v|%v|%v|%v|%v",
		p.Speed, p.Density, p.ShapeCount, p.PulseRate, p.Complexity)
}

// normalizePolicy ensures the policy fields contain valid values,
// clamping out-of-range entries to their nearest valid bound.
func NormalizePolicy(p Policy) Policy {
	// Clamp Speed to [0.1, 10.0]
	if p.Speed < 0.1 {
		p.Speed = 0.1
	}
	if p.Speed > 10.0 {
		p.Speed = 10.0
	}

	// Clamp Density to [0.1, 1.0]
	if p.Density < 0.1 {
		p.Density = 0.1
	}
	if p.Density > 1.0 {
		p.Density = 1.0
	}

	// Clamp ShapeCount to [1, 50]
	if p.ShapeCount < 1 {
		p.ShapeCount = 1
	}
	if p.ShapeCount > 50 {
		p.ShapeCount = 50
	}

	// Clamp PulseRate to [0.1, 5.0]
	if p.PulseRate < 0.1 {
		p.PulseRate = 0.1
	}
	if p.PulseRate > 5.0 {
		p.PulseRate = 5.0
	}

	// Clamp Complexity to [3, 8]
	if p.Complexity < 3 {
		p.Complexity = 3
	}
	if p.Complexity > 8 {
		p.Complexity = 8
	}

	return p
}
