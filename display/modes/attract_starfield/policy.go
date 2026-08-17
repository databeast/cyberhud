package attract_starfield

import (
	"sync"

	"github.com/databeast/cyberhud/display/catalog"
	"github.com/databeast/cyberhud/display/modes/attract_starfield/source"
)

type Policy = source.Policy

func DefaultPolicy() source.Policy { return source.DefaultPolicy() }

// policyState holds the mutex-protected active policy.
var policyState = struct {
	sync.RWMutex
	policy source.Policy
}{
	policy: source.DefaultPolicy(),
}

// GetPolicy returns the current starfield policy (thread-safe read under RWMutex).
func GetPolicy() source.Policy {
	policyState.RLock()
	defer policyState.RUnlock()
	return policyState.policy
}

// SetPolicy updates the starfield policy after normalization (thread-safe write).
func SetPolicy(p source.Policy) {
	policyState.Lock()
	defer policyState.Unlock()
	policyState.policy = normalizePolicy(p)
}

// normalizePolicy clamps all policy fields to their valid ranges.
func normalizePolicy(p source.Policy) source.Policy {
	if p.Speed < 0.1 {
		p.Speed = 0.1
	}
	if p.Speed > 10.0 {
		p.Speed = 10.0
	}

	if p.Density < 0.1 {
		p.Density = 0.1
	}
	if p.Density > 1.0 {
		p.Density = 1.0
	}

	if p.Layers < 1 {
		p.Layers = 1
	}
	if p.Layers > 8 {
		p.Layers = 8
	}

	return p
}

func policyFingerprint(p source.Policy) string { return p.Fingerprint() }

// starfieldSnapshotter implements catalog.PolicySnapshotter for the attract_starfield mode.
type starfieldSnapshotter struct{}

// SnapshotPolicy returns the current starfield policy as a JSON-serializable map
// using snake_case keys matching the protocol wire format.
func (starfieldSnapshotter) SnapshotPolicy() map[string]interface{} {
	return GetPolicy().ToMap()
}

// RestorePolicy applies policy values from a JSON map, running them through
// the same normalization as SetPolicy. Keys use the protocol wire format.
func (starfieldSnapshotter) RestorePolicy(data map[string]interface{}) error {
	p := source.DefaultPolicy()

	if v, ok := data["speed"]; ok {
		if f, ok := source.ToFloat64(v); ok {
			p.Speed = f
		}
	}
	if v, ok := data["density"]; ok {
		if f, ok := source.ToFloat64(v); ok {
			p.Density = f
		}
	}
	if v, ok := data["layers"]; ok {
		if n, ok := source.ToInt(v); ok {
			p.Layers = n
		}
	}

	SetPolicy(p)
	return nil
}

func init() {
	catalog.RegisterSnapshotter("attract_starfield", starfieldSnapshotter{})
}
