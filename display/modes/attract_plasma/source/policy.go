package source

import (
	"fmt"

	"github.com/databeast/cyberhud/display/catalog"
)

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

	// Clamp CycleRate to [0.1, 5.0]
	if p.CycleRate < 0.1 {
		p.CycleRate = 0.1
	}
	if p.CycleRate > 5.0 {
		p.CycleRate = 5.0
	}

	// Clamp BlobScale to [0.5, 4.0]
	if p.BlobScale < 0.5 {
		p.BlobScale = 0.5
	}
	if p.BlobScale > 4.0 {
		p.BlobScale = 4.0
	}

	return p
}

// Policy captures all runtime-configurable parameters for the plasma mode.
// Fields are exposed through catalog registration and CLI command handling via
// the "attract_plasma" command verb.
type Policy struct {
	Speed     float64 // [0.1, 10.0], default 1.0 — morph speed multiplier
	Density   float64 // [0.1, 1.0], default 0.5 — unused, kept for uniformity
	CycleRate float64 // [0.1, 5.0], default 1.0 — color palette cycles/second
	BlobScale float64 // [0.5, 4.0], default 1.0 — spatial frequency multiplier
}

func (p Policy) Fingerprint() string {
	return fmt.Sprintf("%v|%v|%v|%v", p.Speed, p.Density, p.CycleRate, p.BlobScale)
}

func (p Policy) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"speed":      p.Speed,
		"density":    p.Density,
		"cycle_rate": p.CycleRate,
		"blob_scale": p.BlobScale,
	}
}

func (p Policy) Options() []catalog.OptionDefinition {
	return []catalog.OptionDefinition{
		{Key: "speed", Type: "float", Summary: "Morph speed multiplier.", Default: "1.0", Allowed: []string{}},
		{Key: "density", Type: "float", Summary: "Unused, kept for uniformity.", Default: "0.5", Allowed: []string{}},
		{Key: "cycle_rate", Type: "float", Summary: "Color palette cycles per second.", Default: "1.0", Allowed: []string{}},
		{Key: "blob_scale", Type: "float", Summary: "Spatial frequency multiplier for blob size.", Default: "1.0", Allowed: []string{}},
	}
}

// DefaultPolicy returns the default plasma policy with all fields initialized.
func DefaultPolicy() Policy {
	return Policy{
		Speed:     1.0,
		Density:   0.5,
		CycleRate: 1.0,
		BlobScale: 1.0,
	}
}
