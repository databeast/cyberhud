package source

import (
	"fmt"

	"github.com/databeast/cyberhud/display/catalog"
)

// Policy captures all runtime-configurable parameters for the popsquares mode.
// Fields are exposed through catalog registration and CLI command handling via
// the "attract_popsquares" command verb.
//
// The defaults are deliberately calm: a modest grid, a low pop rate, and a cool
// palette, so out of the box the mode reads as a machine idling rather than a
// busy light show. Users can wind it up with the knobs if they want more.
type Policy struct {
	Speed     float64 // [0.1, 10.0], default 1.0 — overall animation rate (pop-and-fade cadence)
	GridScale float64 // [0.5, 4.0], default 1.0 — cell size multiplier (smaller = more, finer squares)
	PopRate   float64 // [0.1, 1.0], default 0.35 — fraction of cells gently active at once (kept low for the idle feel)
	Warmth    float64 // [0.0, 1.0], default 0.15 — palette temperature: 0 = cool blue/teal (idle), 1 = warm amber (busier)
}

// NormalizePolicy ensures the policy fields contain valid values, clamping
// out-of-range entries to their nearest valid bound.
func NormalizePolicy(p Policy) Policy {
	if p.Speed < 0.1 {
		p.Speed = 0.1
	}
	if p.Speed > 10.0 {
		p.Speed = 10.0
	}
	if p.GridScale < 0.5 {
		p.GridScale = 0.5
	}
	if p.GridScale > 4.0 {
		p.GridScale = 4.0
	}
	if p.PopRate < 0.1 {
		p.PopRate = 0.1
	}
	if p.PopRate > 1.0 {
		p.PopRate = 1.0
	}
	if p.Warmth < 0.0 {
		p.Warmth = 0.0
	}
	if p.Warmth > 1.0 {
		p.Warmth = 1.0
	}
	return p
}

func (p Policy) Fingerprint() string {
	return fmt.Sprintf("%v|%v|%v|%v", p.Speed, p.GridScale, p.PopRate, p.Warmth)
}

func (p Policy) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"speed":      p.Speed,
		"grid_scale": p.GridScale,
		"pop_rate":   p.PopRate,
		"warmth":     p.Warmth,
	}
}

func (p Policy) Options() []catalog.OptionDefinition {
	return []catalog.OptionDefinition{
		{Key: "speed", Type: "float", Summary: "Animation rate: pop-and-fade cadence.", Default: "1.0", Allowed: []string{}},
		{Key: "grid_scale", Type: "float", Summary: "Cell size multiplier. Smaller values pack more, finer squares.", Default: "1.0", Allowed: []string{}},
		{Key: "pop_rate", Type: "float", Summary: "Fraction of cells gently active at once. Low values read as idle.", Default: "0.35", Allowed: []string{}},
		{Key: "warmth", Type: "float", Summary: "Palette temperature: 0 cool/idle (blue-teal), 1 warm/busy (amber).", Default: "0.15", Allowed: []string{}},
	}
}

// DefaultPolicy returns the default popsquares policy.
func DefaultPolicy() Policy {
	return Policy{
		Speed:     1.0,
		GridScale: 1.0,
		PopRate:   0.35,
		Warmth:    0.15,
	}
}
