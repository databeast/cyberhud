package source

import "fmt"
import "github.com/databeast/cyberhud/display/catalog"

// Policy captures all runtime-configurable parameters for the geometric attract mode.
type Policy struct {
	Speed         float64 // [0.1, 10.0], default 1.0 — animation speed multiplier
	Density       float64 // [0.1, 1.0], default 0.5 — cluster count scaling
	GlowIntensity float64 // [0.0, 1.0], default 1.0 — glow effect strength
	FragmentRate  float64 // [0.0, 2.0], default 1.0 — pseudocode fragment spawn rate
}

// DefaultPolicy returns the default geometric policy with all fields initialized
// within their valid ranges.
func DefaultPolicy() Policy {
	return Policy{
		Speed:         1.0,
		Density:       0.5,
		GlowIntensity: 1.0,
		FragmentRate:  1.0,
	}
}

// NormalizePolicy ensures all policy fields are within their valid ranges,
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

	// GlowIntensity: [0.0, 1.0]
	if p.GlowIntensity < 0.0 {
		p.GlowIntensity = 0.0
	}
	if p.GlowIntensity > 1.0 {
		p.GlowIntensity = 1.0
	}

	// FragmentRate: [0.0, 2.0]
	if p.FragmentRate < 0.0 {
		p.FragmentRate = 0.0
	}
	if p.FragmentRate > 2.0 {
		p.FragmentRate = 2.0
	}

	return p
}

// Fingerprint returns a deterministic string representation of the policy.
func (p Policy) Fingerprint() string {
	return fmt.Sprintf("%v|%v|%v|%v", p.Speed, p.Density, p.GlowIntensity, p.FragmentRate)
}

// ToMap returns the policy as a snake_case map for snapshotting and wiring.
func (p Policy) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"speed":          p.Speed,
		"density":        p.Density,
		"glow_intensity": p.GlowIntensity,
		"fragment_rate":  p.FragmentRate,
	}
}

// Options returns the catalog metadata for the geometric policy.
func (Policy) Options() []catalog.OptionDefinition {
	return []catalog.OptionDefinition{
		{Key: "speed", Type: "float", Summary: "Animation speed multiplier.", Default: "1.0", Allowed: []string{}},
		{Key: "density", Type: "float", Summary: "Cluster count scaling from 0.1 to 1.0.", Default: "0.5", Allowed: []string{}},
		{Key: "glow_intensity", Type: "float", Summary: "Glow effect strength from 0.0 to 1.0.", Default: "1.0", Allowed: []string{}},
		{Key: "fragment_rate", Type: "float", Summary: "Pseudocode fragment spawn rate from 0.0 to 2.0.", Default: "1.0", Allowed: []string{}},
	}
}
