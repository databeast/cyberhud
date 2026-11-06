package attract_geometric

import (
	"github.com/databeast/cyberhud/display/catalog"
	"github.com/databeast/cyberhud/display/catalog/cmdutil"
)

// geometricSnapshotter implements catalog.PolicySnapshotter for the attract_geometric mode.
type geometricSnapshotter struct{}

// SnapshotPolicy returns the current geometric policy as a JSON-serializable map
// using snake_case keys matching the protocol wire format.
func (geometricSnapshotter) SnapshotPolicy() map[string]interface{} {
	p := GetPolicy()
	return map[string]interface{}{
		"speed":          p.Speed,
		"density":        p.Density,
		"glow_intensity": p.GlowIntensity,
		"fragment_rate":  p.FragmentRate,
	}
}

// RestorePolicy applies policy values from a JSON map, running them through
// the same normalization as SetPolicy. Keys use the protocol wire format.
func (geometricSnapshotter) RestorePolicy(data map[string]interface{}) error {
	p := DefaultPolicy()

	if v, ok := data["speed"]; ok {
		if f, ok := cmdutil.ToFloat64(v); ok {
			p.Speed = f
		}
	}
	if v, ok := data["density"]; ok {
		if f, ok := cmdutil.ToFloat64(v); ok {
			p.Density = f
		}
	}
	if v, ok := data["glow_intensity"]; ok {
		if f, ok := cmdutil.ToFloat64(v); ok {
			p.GlowIntensity = f
		}
	}
	if v, ok := data["fragment_rate"]; ok {
		if f, ok := cmdutil.ToFloat64(v); ok {
			p.FragmentRate = f
		}
	}

	SetPolicy(p)
	return nil
}

func init() {
	catalog.RegisterSnapshotter("attract_geometric", geometricSnapshotter{})
}
