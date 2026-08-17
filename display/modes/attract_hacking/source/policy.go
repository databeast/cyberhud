package source

import (
	"fmt"

	"github.com/databeast/cyberhud/display/catalog"
)

// Policy captures runtime tuning for the Hollywood-hacking attract mode.
type Policy struct {
	Speed     float64 // animation pacing multiplier, [0.1, 3.0], default 1.0
	Density   float64 // visual density, [0.1, 1.0], default 0.7
	Glitch    float64 // glitch intensity, [0.0, 1.0], default 0.5
	Intensity float64 // glow intensity, [0.1, 1.0], default 0.8
	Pulse     float64 // pulse speed, [0.1, 1.5], default 0.7
}

func (p Policy) Fingerprint() string {
	return fmt.Sprintf("%v|%v|%v|%v|%v", p.Speed, p.Density, p.Glitch, p.Intensity, p.Pulse)
}

func (p Policy) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"speed":     p.Speed,
		"density":   p.Density,
		"glitch":    p.Glitch,
		"intensity": p.Intensity,
		"pulse":     p.Pulse,
	}
}

func (p Policy) Options() []catalog.OptionDefinition {
	return []catalog.OptionDefinition{
		{Key: "speed", Type: "float", Summary: "Animation pacing multiplier.", Default: "1.0", Allowed: []string{}},
		{Key: "density", Type: "float", Summary: "Overall clutter and log density.", Default: "0.7", Allowed: []string{}},
		{Key: "glitch", Type: "float", Summary: "Amount of visual glitching and offset jitter.", Default: "0.5", Allowed: []string{}},
		{Key: "intensity", Type: "float", Summary: "Glow and neon intensity.", Default: "0.8", Allowed: []string{}},
		{Key: "pulse", Type: "float", Summary: "Scanline and pulse cadence.", Default: "0.7", Allowed: []string{}},
	}
}

func DefaultPolicy() Policy {
	return Policy{Speed: 1.0, Density: 0.7, Glitch: 0.5, Intensity: 0.8, Pulse: 0.7}
}
