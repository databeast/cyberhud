package attract_bokeh

import (
	"sync"

	"github.com/databeast/cyberhud/display/catalog"
	"github.com/databeast/cyberhud/display/modes/attract_bokeh/source"
)

// policyState holds the mutex-protected active policy.
var policyState = struct {
	sync.RWMutex
	policy source.Policy
}{
	policy: source.DefaultPolicy(),
}

// GetPolicy returns the current bokeh policy (thread-safe read under RWMutex).
func GetPolicy() source.Policy {
	policyState.RLock()
	defer policyState.RUnlock()
	return policyState.policy
}

// SetPolicy updates the bokeh policy after normalization (thread-safe write).
func SetPolicy(p source.Policy) {
	policyState.Lock()
	defer policyState.Unlock()
	policyState.policy = source.NormalizePolicy(p)
}

// bokehSnapshotter implements catalog.PolicySnapshotter for the attract_bokeh mode.
type bokehSnapshotter struct{}

// SnapshotPolicy returns the current bokeh policy as a JSON-serializable map
// with snake_case keys.
func (bokehSnapshotter) SnapshotPolicy() map[string]interface{} {
	p := GetPolicy()
	return map[string]interface{}{
		"speed":         p.Speed,
		"density":       p.Density,
		"size_variance": p.SizeVariance,
		"saturation":    p.Saturation,
	}
}

// RestorePolicy applies policy values from a JSON map, using the same
// normalization path as SetPolicy to ensure all constraints are enforced.
func (bokehSnapshotter) RestorePolicy(data map[string]interface{}) error {
	p := source.DefaultPolicy()

	if v, ok := data["speed"]; ok {
		if f, ok := v.(float64); ok {
			p.Speed = f
		}
	}
	if v, ok := data["density"]; ok {
		if f, ok := v.(float64); ok {
			p.Density = f
		}
	}
	if v, ok := data["size_variance"]; ok {
		if f, ok := v.(float64); ok {
			p.SizeVariance = f
		}
	}
	if v, ok := data["saturation"]; ok {
		if f, ok := v.(float64); ok {
			p.Saturation = f
		}
	}

	SetPolicy(p)
	return nil
}

func init() {
	catalog.RegisterSnapshotter("attract_bokeh", bokehSnapshotter{})
}
