package source

import (
	"fmt"

	"github.com/databeast/cyberhud/display/catalog"
)

// Policy captures all runtime-configurable parameters for the matrix rain mode.
// Fields are exposed through catalog registration and CLI command handling via
// the "matrix" command verb.
//
// as the single source of truth for mode behavior, referenced by catalog defaults,
// normalization, CLI handlers, and the rendering pipeline.
type Policy struct {
	MinSpeed       float64 // cells/second, default 3.0, must be > 0
	MaxSpeed       float64 // cells/second, default 12.0, must be > 0
	TrailLength    int     // default 16, clamped to [4, 128]
	Density        float64 // [0.1, 1.0], default 1.0
	ShowBackground bool    // default false
}

func (p Policy) Options() []catalog.OptionDefinition {
	return []catalog.OptionDefinition{
		{Key: "min_speed", Type: "float", Summary: "Minimum column scroll speed in cells/second.", Default: "1.5", Allowed: []string{}},
		{Key: "max_speed", Type: "float", Summary: "Maximum column scroll speed in cells/second.", Default: "6", Allowed: []string{}},
		{Key: "trail_length", Type: "int", Summary: "Number of trailing cells behind the lead character.", Default: "16", Allowed: []string{}},
		{Key: "density", Type: "float", Summary: "Fraction of columns active, from 0.1 to 1.0.", Default: "1.0", Allowed: []string{}},
		{Key: "show_background", Type: "bool", Summary: "Enable radial gradient background on color panels.", Default: "false", Allowed: []string{"true", "false"}},
	}
}

func (p Policy) Fingerprint() string {
	return PolicyFingerprint(p)
}

func (p Policy) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"min_speed":       p.MinSpeed,
		"max_speed":       p.MaxSpeed,
		"trail_length":    p.TrailLength,
		"density":         p.Density,
		"show_background": p.ShowBackground,
	}
}

// DefaultPolicy returns the default matrix policy with all fields initialized.
//
// default state that catalog registration and normalization reference as the baseline.
func DefaultPolicy() Policy {
	return Policy{
		MinSpeed:       1.5,
		MaxSpeed:       6.0,
		TrailLength:    16,
		Density:        1.0,
		ShowBackground: false,
	}
}

// normalizePolicy ensures the policy fields contain valid values,
// clamping out-of-range entries to their nearest valid bound.
func NormalizePolicy(p Policy) Policy {
	d := DefaultPolicy()

	// Clamp MinSpeed > 0
	if p.MinSpeed <= 0 {
		p.MinSpeed = d.MinSpeed
	}

	// Clamp MaxSpeed > 0
	if p.MaxSpeed <= 0 {
		p.MaxSpeed = d.MaxSpeed
	}

	// If MinSpeed > MaxSpeed, set MinSpeed = MaxSpeed
	if p.MinSpeed > p.MaxSpeed {
		p.MinSpeed = p.MaxSpeed
	}

	// Clamp TrailLength to [4, 128]
	if p.TrailLength < 4 {
		p.TrailLength = 4
	}
	if p.TrailLength > 128 {
		p.TrailLength = 128
	}

	// Clamp Density to [0.1, 1.0]
	if p.Density < 0.1 {
		p.Density = 0.1
	}
	if p.Density > 1.0 {
		p.Density = 1.0
	}

	return p
}

// policyFingerprint returns a pipe-delimited string representation of all
// policy fields, used for change detection in RenderCacheKey and strip cache
// invalidation.
func PolicyFingerprint(p Policy) string {
	return fmt.Sprintf("%v|%v|%v|%v|%v",
		p.MinSpeed, p.MaxSpeed, p.TrailLength, p.Density, p.ShowBackground)
}
