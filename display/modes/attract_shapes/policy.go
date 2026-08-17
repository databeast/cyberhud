package attract_shapes

import (
	"sync"

	"github.com/databeast/cyberhud/display/catalog"
	"github.com/databeast/cyberhud/display/modes/attract_shapes/source"
)

// shapesSnapshotter implements catalog.PolicySnapshotter for the attract_shapes mode.
type shapesSnapshotter struct{}

// SnapshotPolicy returns the current shapes policy as a JSON-serializable map
// using snake_case keys matching the protocol wire format.
func (shapesSnapshotter) SnapshotPolicy() map[string]interface{} {
	p := GetPolicy()
	return map[string]interface{}{
		"speed":       p.Speed,
		"density":     p.Density,
		"shape_count": p.ShapeCount,
		"pulse_rate":  p.PulseRate,
		"complexity":  p.Complexity,
	}
}

// RestorePolicy applies policy values from a JSON map, running them through
// the same normalization as SetPolicy. Keys use the protocol wire format.
func (shapesSnapshotter) RestorePolicy(data map[string]interface{}) error {
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
	if v, ok := data["shape_count"]; ok {
		if n, ok := toInt(v); ok {
			p.ShapeCount = n
		}
	}
	if v, ok := data["pulse_rate"]; ok {
		if f, ok := toFloat64(v); ok {
			p.PulseRate = f
		}
	}
	if v, ok := data["complexity"]; ok {
		if n, ok := toInt(v); ok {
			p.Complexity = n
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
	catalog.RegisterSnapshotter("attract_shapes", shapesSnapshotter{})
}

// policyState holds the mutex-protected active policy.
var policyState = struct {
	sync.RWMutex
	policy source.Policy
}{
	policy: source.DefaultPolicy(),
}

// GetPolicy returns the current shapes policy (thread-safe read under RWMutex).
func GetPolicy() source.Policy {
	policyState.RLock()
	defer policyState.RUnlock()
	return policyState.policy
}

// SetPolicy updates the shapes policy after normalization (thread-safe write under Mutex).
func SetPolicy(p source.Policy) {
	policyState.Lock()
	defer policyState.Unlock()
	policyState.policy = normalizePolicy(p)
}

func DefaultPolicy() source.Policy {
	return source.DefaultPolicy()
}

func normalizePolicy(p source.Policy) source.Policy {
	return source.NormalizePolicy(p)
}
