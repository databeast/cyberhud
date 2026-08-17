package source

import (
	"fmt"

	"github.com/databeast/cyberhud/display/catalog"
)

// Policy captures all runtime-configurable parameters for the particles mode.
// Fields are exposed through catalog registration and CLI command handling via
// the "attract_particles" command verb.
type Policy struct {
	Speed   float64 // [0.1, 10.0], default 1.0 — particle speed multiplier
	Density float64 // [0.1, 1.0], default 0.5 — particle count scaling
	Drift   float64 // [0.0, 1.0], default 0.3 — random directional wander per frame
	Glow    float64 // [0.1, 1.0], default 0.5 — hue cycling intensity
}

func (p Policy) Fingerprint() string {
	return policyFingerprint(p)
}

func (p Policy) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"speed":   p.Speed,
		"density": p.Density,
		"drift":   p.Drift,
		"glow":    p.Glow,
	}
}

func (p Policy) Options() []catalog.OptionDefinition {
	return []catalog.OptionDefinition{
		{Key: "speed", Type: "float", Summary: "Particle speed multiplier.", Default: "1.0", Allowed: []string{}},
		{Key: "density", Type: "float", Summary: "Particle count scaling factor from 0.1 to 1.0.", Default: "0.5", Allowed: []string{}},
		{Key: "drift", Type: "float", Summary: "Random directional wander per frame, 0.0 to 1.0.", Default: "0.3", Allowed: []string{}},
		{Key: "glow", Type: "float", Summary: "Hue cycling intensity from 0.1 to 1.0.", Default: "0.5", Allowed: []string{}},
	}
}

// DefaultPolicy returns the default particles policy with all fields initialized.
func DefaultPolicy() Policy {
	return Policy{
		Speed:   1.0,
		Density: 0.5,
		Drift:   0.3,
		Glow:    0.5,
	}
}

// NormalizePolicy ensures the policy fields contain valid values,
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

	// Clamp Drift to [0.0, 1.0]
	if p.Drift < 0.0 {
		p.Drift = 0.0
	}
	if p.Drift > 1.0 {
		p.Drift = 1.0
	}

	// Clamp Glow to [0.1, 1.0]
	if p.Glow < 0.1 {
		p.Glow = 0.1
	}
	if p.Glow > 1.0 {
		p.Glow = 1.0
	}

	return p
}

// policyFingerprint returns a pipe-delimited string representation of all
// policy fields, used for change detection in RenderCacheKey.
func policyFingerprint(p Policy) string {
	return fmt.Sprintf("%v|%v|%v|%v", p.Speed, p.Density, p.Drift, p.Glow)
}
