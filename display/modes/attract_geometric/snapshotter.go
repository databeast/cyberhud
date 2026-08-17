package attract_geometric

import (
	"github.com/databeast/cyberhud/display/catalog"
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
		if f, ok := toFloat64(v); ok {
			p.Speed = f
		}
	}
	if v, ok := data["density"]; ok {
		if f, ok := toFloat64(v); ok {
			p.Density = f
		}
	}
	if v, ok := data["glow_intensity"]; ok {
		if f, ok := toFloat64(v); ok {
			p.GlowIntensity = f
		}
	}
	if v, ok := data["fragment_rate"]; ok {
		if f, ok := toFloat64(v); ok {
			p.FragmentRate = f
		}
	}

	SetPolicy(p)
	return nil
}

// toFloat64 extracts a float64 from an interface value, handling both
// float64 (native JSON number) and int/int64 conversions.
func toFloat64(v interface{}) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case int:
		return float64(n), true
	case int64:
		return float64(n), true
	default:
		return 0, false
	}
}

func init() {
	catalog.RegisterSnapshotter("attract_geometric", geometricSnapshotter{})
}
