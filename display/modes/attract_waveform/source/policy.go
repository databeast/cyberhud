package source

import (
	"fmt"

	"github.com/databeast/cyberhud/display/catalog"
)

// Policy captures all runtime-configurable parameters for the waveform attract mode.
// Fields are exposed through catalog registration and CLI command handling via
// the "attract_waveform" command verb.
type Policy struct {
	Speed       float64 // animation rate multiplier, [0.1, 10.0], default 1.0
	Density     float64 // number of waveform cycles across the panel (maps [0.1,1.0] → [1.5,5.5] cycles), default 0.5
	Amplitude   float64 // fraction of half panel height for trace amplitude, [0.1, 1.0], default 0.8
	Traces      int     // number of simultaneous waveform traces, [1, 8], default 3
	Persistence float64 // phosphor trail decay as fraction of panel width, [0.1, 1.0], default 0.5
	Direction   string  // trace direction: "horizontal" or "vertical", default "auto" (picks based on aspect ratio)
}

func (p Policy) Options() []catalog.OptionDefinition {
	return []catalog.OptionDefinition{
		{Key: "speed", Type: "float", Summary: "Animation rate multiplier.", Default: "1.0", Allowed: []string{}},
		{Key: "density", Type: "float", Summary: "Waveform cycle density: maps [0.1,1.0] to [1.5,5.5] cycles across the panel. Higher values show more cycles and faster-apparent scrolling.", Default: "0.5", Allowed: []string{}},
		{Key: "amplitude", Type: "float", Summary: "Trace amplitude as fraction of panel height.", Default: "0.8", Allowed: []string{}},
		{Key: "traces", Type: "int", Summary: "Number of simultaneous waveform traces.", Default: "3", Allowed: []string{}},
		{Key: "persistence", Type: "float", Summary: "Phosphor trail decay as fraction of panel width.", Default: "0.5", Allowed: []string{}},
		{Key: "direction", Type: "string", Summary: "Trace direction: auto, horizontal, or vertical.", Default: "auto", Allowed: []string{"auto", "horizontal", "vertical"}},
	}
}

// DefaultPolicy returns the default waveform policy with all fields initialized
// to valid values that do not trigger normalization.
func DefaultPolicy() Policy {
	return Policy{
		Speed:       1.0,
		Density:     0.5,
		Amplitude:   0.8,
		Traces:      3,
		Persistence: 0.5,
		Direction:   "auto",
	}
}

func (p Policy) Fingerprint() string {
	return policyFingerprint(p)
}

func (p Policy) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"speed":       p.Speed,
		"density":     p.Density,
		"amplitude":   p.Amplitude,
		"traces":      p.Traces,
		"persistence": p.Persistence,
		"direction":   p.Direction,
	}
}

// policyFingerprint returns a pipe-delimited string representation of all
// policy fields, used for change detection in RenderCacheKey.
func policyFingerprint(p Policy) string {
	return fmt.Sprintf("%v|%v|%v|%v|%v|%s",
		p.Speed, p.Density, p.Amplitude, p.Traces, p.Persistence, p.Direction)
}
