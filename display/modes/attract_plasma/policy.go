package attract_plasma

import (
	"sync"

	"github.com/databeast/cyberhud/display/catalog"
	"github.com/databeast/cyberhud/display/catalog/cmdutil"
	"github.com/databeast/cyberhud/display/modes/attract_plasma/source"
)

type Policy = source.Policy

// policyState holds the mutex-protected active policy.
var policyState = struct {
	sync.RWMutex
	policy Policy
}{
	policy: DefaultPolicy(),
}

func DefaultPolicy() source.Policy {
	return source.DefaultPolicy()
}

// GetPolicy returns the current plasma policy (thread-safe read under RWMutex).
func GetPolicy() Policy {
	policyState.RLock()
	defer policyState.RUnlock()
	return policyState.policy
}

// SetPolicy updates the plasma policy after normalization (thread-safe write under Mutex).
func SetPolicy(p Policy) {
	policyState.Lock()
	defer policyState.Unlock()
	policyState.policy = normalizePolicy(p)
}

func normalizePolicy(p source.Policy) source.Policy {
	return source.NormalizePolicy(p)
}

// plasmaSnapshotter implements catalog.PolicySnapshotter for the attract_plasma mode.
type plasmaSnapshotter struct{}

// SnapshotPolicy returns the current plasma policy as a JSON-serializable map
// using snake_case keys matching the protocol wire format.
func (plasmaSnapshotter) SnapshotPolicy() map[string]interface{} {
	p := GetPolicy()
	return map[string]interface{}{
		"speed":      p.Speed,
		"density":    p.Density,
		"cycle_rate": p.CycleRate,
		"blob_scale": p.BlobScale,
	}
}

// RestorePolicy applies policy values from a JSON map, running them through
// the same normalization as SetPolicy. Keys use the protocol wire format.
func (plasmaSnapshotter) RestorePolicy(data map[string]interface{}) error {
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
	if v, ok := data["cycle_rate"]; ok {
		if f, ok := cmdutil.ToFloat64(v); ok {
			p.CycleRate = f
		}
	}
	if v, ok := data["blob_scale"]; ok {
		if f, ok := cmdutil.ToFloat64(v); ok {
			p.BlobScale = f
		}
	}

	SetPolicy(p)
	return nil
}

func init() {
	catalog.RegisterSnapshotter("attract_plasma", plasmaSnapshotter{})
}
