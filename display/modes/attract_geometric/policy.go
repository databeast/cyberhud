package attract_geometric

import (
	"sync"

	"github.com/databeast/cyberhud/display/modes/attract_geometric/source"
)

// policyState holds the mutex-protected active policy.
var policyState = struct {
	sync.RWMutex
	policy source.Policy
}{
	policy: source.DefaultPolicy(),
}

// GetPolicy returns the current geometric policy (thread-safe read under RWMutex).
func GetPolicy() source.Policy {
	policyState.RLock()
	defer policyState.RUnlock()
	return policyState.policy
}

// SetPolicy updates the geometric policy after normalization (thread-safe write).
func SetPolicy(p source.Policy) {
	policyState.Lock()
	defer policyState.Unlock()
	policyState.policy = source.NormalizePolicy(p)
}

func DefaultPolicy() source.Policy { return source.DefaultPolicy() }

func normalizePolicy(p source.Policy) source.Policy { return source.NormalizePolicy(p) }

func policyFingerprint(p source.Policy) string { return p.Fingerprint() }
