package source

import (
	"fmt"

	"github.com/databeast/cyberhud/display/catalog"
)

// Policy captures all runtime-configurable parameters for the bokeh attract mode.
type Policy struct {
	Speed        float64 // [0.1, 10.0], default 1.0 — drift speed multiplier
	Density      float64 // [0.1, 1.0], default 0.5 — circle count scaling
	SizeVariance float64 // [0.0, 1.0], default 0.5 — radius spread (0=uniform, 1=max spread)
	Saturation   float64 // [0.0, 1.0], default 0.7 — color saturation
}

func (p Policy) Fingerprint() string {
	return fmt.Sprintf("%v|%v|%v|%v", p.Speed, p.Density, p.SizeVariance, p.Saturation)
}

func (p Policy) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"speed":         p.Speed,
		"density":       p.Density,
		"size_variance": p.SizeVariance,
		"saturation":    p.Saturation,
	}
}

func (p Policy) Options() []catalog.OptionDefinition {
	return []catalog.OptionDefinition{
		{Key: "speed", Type: "float", Summary: "Drift speed multiplier.", Default: "1.0", Allowed: []string{}},
		{Key: "density", Type: "float", Summary: "Circle count scaling from 0.1 to 1.0.", Default: "0.5", Allowed: []string{}},
		{Key: "size_variance", Type: "float", Summary: "Radius spread, 0.0 uniform to 1.0 maximum.", Default: "0.5", Allowed: []string{}},
		{Key: "saturation", Type: "float", Summary: "Color saturation from 0.0 to 1.0.", Default: "0.7", Allowed: []string{}},
	}
}

// normalizePolicy ensures all policy fields are within their valid ranges,
// clamping out-of-range entries to their nearest valid bound.
func NormalizePolicy(p Policy) Policy {
	// Speed: [0.1, 10.0]
	if p.Speed < 0.1 {
		p.Speed = 0.1
	}
	if p.Speed > 10.0 {
		p.Speed = 10.0
	}

	// Density: [0.1, 1.0]
	if p.Density < 0.1 {
		p.Density = 0.1
	}
	if p.Density > 1.0 {
		p.Density = 1.0
	}

	// SizeVariance: [0.0, 1.0]
	if p.SizeVariance < 0.0 {
		p.SizeVariance = 0.0
	}
	if p.SizeVariance > 1.0 {
		p.SizeVariance = 1.0
	}

	// Saturation: [0.0, 1.0]
	if p.Saturation < 0.0 {
		p.Saturation = 0.0
	}
	if p.Saturation > 1.0 {
		p.Saturation = 1.0
	}

	return p
}

// policyFingerprint returns a deterministic string representation of all policy
// fields, used for change detection in RenderCacheKey.
func policyFingerprint(p Policy) string {
	return fmt.Sprintf("%v|%v|%v|%v", p.Speed, p.Density, p.SizeVariance, p.Saturation)
}

// DefaultPolicy returns the default bokeh policy with all fields initialized
// within their valid ranges.
func DefaultPolicy() Policy {
	return Policy{
		Speed:        1.0,
		Density:      0.5,
		SizeVariance: 0.5,
		Saturation:   0.7,
	}
}
