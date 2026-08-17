package attract_particles

import (
	"sync"

	"github.com/databeast/cyberhud/display/catalog"
	"github.com/databeast/cyberhud/display/modes/attract_particles/source"
)

type Policy = source.Policy

// policyState holds the mutex-protected active policy.
var policyState = struct {
	sync.RWMutex
	policy source.Policy
}{
	policy: source.DefaultPolicy(),
}

func DefaultPolicy() source.Policy { return source.DefaultPolicy() }

// GetPolicy returns the current particles policy (thread-safe read under RWMutex).
func GetPolicy() source.Policy {
	policyState.RLock()
	defer policyState.RUnlock()
	return policyState.policy
}

// SetPolicy updates the particles policy after normalization (thread-safe write under Mutex).
func SetPolicy(p source.Policy) {
	policyState.Lock()
	defer policyState.Unlock()
	policyState.policy = normalizePolicy(p)
}

// normalizePolicy keeps the controller's active policy within valid bounds.
func normalizePolicy(p source.Policy) source.Policy {
	return source.NormalizePolicy(p)
}

// particlesSnapshotter implements catalog.PolicySnapshotter for the attract_particles mode.
type particlesSnapshotter struct{}

// SnapshotPolicy returns the current particles policy as a JSON-serializable map
// using snake_case keys matching the protocol wire format.
func (particlesSnapshotter) SnapshotPolicy() map[string]interface{} {
	p := GetPolicy()
	return map[string]interface{}{
		"speed":   p.Speed,
		"density": p.Density,
		"drift":   p.Drift,
		"glow":    p.Glow,
	}
}

// RestorePolicy applies policy values from a JSON map, running them through
// the same normalization as SetPolicy. Keys use the protocol wire format.
func (particlesSnapshotter) RestorePolicy(data map[string]interface{}) error {
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
	if v, ok := data["drift"]; ok {
		if f, ok := toFloat64(v); ok {
			p.Drift = f
		}
	}
	if v, ok := data["glow"]; ok {
		if f, ok := toFloat64(v); ok {
			p.Glow = f
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
	catalog.RegisterSnapshotter("attract_particles", particlesSnapshotter{})
}
