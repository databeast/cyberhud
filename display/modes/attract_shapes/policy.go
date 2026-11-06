package attract_shapes

import (
	"sync"

	"github.com/databeast/cyberhud/display/catalog"
	"github.com/databeast/cyberhud/display/catalog/cmdutil"
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
		if f, ok := cmdutil.ToFloat64(v); ok {
			p.Speed = f
		}
	}
	if v, ok := data["density"]; ok {
		if f, ok := cmdutil.ToFloat64(v); ok {
			p.Density = f
		}
	}
	if v, ok := data["shape_count"]; ok {
		if n, ok := cmdutil.ToInt(v); ok {
			p.ShapeCount = n
		}
	}
	if v, ok := data["pulse_rate"]; ok {
		if f, ok := cmdutil.ToFloat64(v); ok {
			p.PulseRate = f
		}
	}
	if v, ok := data["complexity"]; ok {
		if n, ok := cmdutil.ToInt(v); ok {
			p.Complexity = n
		}
	}

	SetPolicy(p)
	return nil
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
