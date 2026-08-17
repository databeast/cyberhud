package stemma

import (
	"strings"
	"sync"

	"github.com/databeast/cyberhud/display/catalog"
	"github.com/databeast/cyberhud/display/modes/stemma/source"
)

// stemmaSnapshotter implements catalog.PolicySnapshotter for the stemma mode.
type stemmaSnapshotter struct{}

// SnapshotPolicy returns the current stemma policy as a JSON-serializable map
// using snake_case keys matching the protocol wire format.
func (stemmaSnapshotter) SnapshotPolicy() map[string]interface{} {
	return GetPolicy().ToMap()
}

// RestorePolicy applies policy values from a JSON map, running them through
// the same normalization as SetPolicy. Keys use the protocol wire format.
func (stemmaSnapshotter) RestorePolicy(data map[string]interface{}) error {
	p := source.DefaultPolicy()

	if v, ok := data["style"]; ok {
		if s, ok := v.(string); ok {
			p.Style = s
		}
	}

	SetPolicy(p)
	return nil
}

func init() {
	catalog.RegisterSnapshotter("stemma", stemmaSnapshotter{})
}

// policyState holds the mutex-protected active policy.
var policyState = struct {
	sync.RWMutex
	policy source.Policy
}{
	policy: source.DefaultPolicy(),
}

// GetPolicy returns the current Stemma policy (thread-safe).
func GetPolicy() source.Policy {
	policyState.RLock()
	defer policyState.RUnlock()
	return policyState.policy
}

// SetPolicy updates the Stemma policy after normalization (thread-safe).
func SetPolicy(p source.Policy) {
	policyState.Lock()
	defer policyState.Unlock()
	policyState.policy = normalizePolicy(p)
}

// DefaultPolicy returns the default Stemma policy.
func DefaultPolicy() source.Policy {
	return source.DefaultPolicy()
}

// normalizePolicy ensures registry-dependent policy fields contain valid values.
func normalizePolicy(p source.Policy) source.Policy {
	p.Style = strings.ToLower(strings.TrimSpace(p.Style))
	if p.Style != "" && stemmaRegistry.Lookup(p.Style) == nil {
		p.Style = ""
	}
	return p
}
