package attract_waveform

import (
	"github.com/databeast/cyberhud/display/catalog"
	"github.com/databeast/cyberhud/display/modes/attract_waveform/source"
)

// waveformSnapshotter implements catalog.PolicySnapshotter for the attract_waveform mode.
type waveformSnapshotter struct{}

// SnapshotPolicy returns the current waveform policy as a JSON-serializable map
// using snake_case keys matching the protocol wire format.
func (waveformSnapshotter) SnapshotPolicy() map[string]interface{} {
	p := GetPolicy()
	return map[string]interface{}{
		"speed":       p.Speed,
		"density":     p.Density,
		"amplitude":   p.Amplitude,
		"traces":      p.Traces,
		"persistence": p.Persistence,
		"direction":   p.Direction,
	}
}

// RestorePolicy applies policy values from a JSON map, running them through
// the same normalization as SetPolicy. Keys use the protocol wire format.
func (waveformSnapshotter) RestorePolicy(data map[string]interface{}) error {
	p := source.DefaultPolicy()

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
	if v, ok := data["amplitude"]; ok {
		if f, ok := toFloat64(v); ok {
			p.Amplitude = f
		}
	}
	if v, ok := data["traces"]; ok {
		if n, ok := toInt(v); ok {
			p.Traces = n
		}
	}
	if v, ok := data["persistence"]; ok {
		if f, ok := toFloat64(v); ok {
			p.Persistence = f
		}
	}
	if v, ok := data["direction"]; ok {
		if s, ok := v.(string); ok {
			p.Direction = s
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

// toInt extracts an int from an interface value, handling both float64
// (JSON numbers decode as float64) and direct int types.
func toInt(v interface{}) (int, bool) {
	switch n := v.(type) {
	case float64:
		return int(n), true
	case int:
		return n, true
	case int64:
		return int(n), true
	default:
		return 0, false
	}
}

func init() {
	catalog.RegisterSnapshotter("attract_waveform", waveformSnapshotter{})
}
