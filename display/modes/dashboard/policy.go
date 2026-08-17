package dashboard

import (
	"strings"
	"sync"

	"github.com/databeast/cyberhud/display/catalog"
	"github.com/databeast/cyberhud/display/modes/dashboard/source"
)

var (
	policyMu sync.RWMutex
	policy   = source.DefaultPolicy()
)

func init() {
	catalog.RegisterSnapshotter("dashboard", dashboardSnapshotter{})
}

// dashboardSnapshotter implements catalog.PolicySnapshotter for the dashboard mode.
type dashboardSnapshotter struct{}

// SnapshotPolicy returns the current dashboard policy as a JSON-serializable map
// using snake_case keys matching the protocol wire format.
func (dashboardSnapshotter) SnapshotPolicy() map[string]interface{} {
	return GetPolicy().ToMap()

}

// RestorePolicy applies policy values from a JSON map, running them through
// the same normalization as SetPolicy. Keys use the protocol wire format.
func (dashboardSnapshotter) RestorePolicy(data map[string]interface{}) error {
	p := source.DefaultPolicy()

	if v, ok := data["style"]; ok {
		if s, ok := v.(string); ok {
			p.Style = s
		}
	}
	if v, ok := data["color_accent"]; ok {
		if s, ok := v.(string); ok {
			p.ColorAccent = s
		}
	}

	SetPolicy(p)
	return nil
}

// SetPolicy updates the dashboard policy with normalization (thread-safe).
func SetPolicy(p source.Policy) {
	p = normalizePolicy(p)
	policyMu.Lock()
	defer policyMu.Unlock()
	policy = p
}

// GetPolicy returns a copy of the current policy (thread-safe).
func GetPolicy() source.Policy {
	policyMu.RLock()
	defer policyMu.RUnlock()
	return policy
}

// ResetPolicy resets the dashboard policy to defaults (non-explicit).
// Used for testing.
func ResetPolicy() {
	policyMu.Lock()
	defer policyMu.Unlock()
	policy = source.DefaultPolicy()
}

// normalizePolicy validates and corrects all Policy fields.
// Empty Style means auto-detect via fitness; only validate non-empty values.
// Invalid ColorAccent → "cyan".
func normalizePolicy(p source.Policy) source.Policy {
	// Validate Style against the registry; empty means auto-detect via fitness.
	p.Style = strings.ToLower(strings.TrimSpace(p.Style))
	if p.Style != "" && dashboardRegistry.Lookup(p.Style) == nil {
		p.Style = ""
	}

	// Validate ColorAccent against allowed accents.
	validAccent := false
	for _, a := range source.AllowedAccents {
		if p.ColorAccent == a {
			validAccent = true
			break
		}
	}
	if !validAccent {
		p.ColorAccent = "cyan"
	}

	return p
}
