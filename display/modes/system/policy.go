package system

import (
	"fmt"
	"strings"
	"sync"

	"github.com/databeast/cyberhud/display/catalog"
	"github.com/databeast/cyberhud/display/modes/system/source"
)

type Policy = source.Policy

// systemSnapshotter implements catalog.PolicySnapshotter for the system mode.
type systemSnapshotter struct{}

// SnapshotPolicy returns the current system policy as a JSON-serializable map
// using snake_case keys matching the protocol wire format.
func (systemSnapshotter) SnapshotPolicy() map[string]interface{} {
	p := GetPolicy()
	return map[string]interface{}{
		"style": p.Style,
		"font":  p.Font,
	}
}

// RestorePolicy applies policy values from a JSON map, running them through
// the same normalization as SetPolicy. Keys use the protocol wire format.
func (systemSnapshotter) RestorePolicy(data map[string]interface{}) error {
	p := DefaultPolicy()

	if v, ok := data["style"]; ok {
		if s, ok := v.(string); ok {
			p.Style = s
		}
	}
	if v, ok := data["font"]; ok {
		if s, ok := v.(string); ok {
			p.Font = s
		}
	}

	SetPolicy(p)
	return nil
}

func init() {
	catalog.RegisterSnapshotter("system", systemSnapshotter{})
}

// GetPolicy returns the current system policy (thread-safe).
func GetPolicy() Policy {
	policyState.RLock()
	defer policyState.RUnlock()
	return policyState.policy
}

// SetPolicy updates the system policy after normalization (thread-safe).
func SetPolicy(p Policy) {
	policyState.Lock()
	defer policyState.Unlock()
	policyState.policy = normalizePolicy(p)
}

func DefaultPolicy() Policy {
	return source.DefaultPolicy()
}

// normalizePolicy ensures policy fields contain valid values and validates styles
// against the root-owned registry.
func normalizePolicy(p Policy) Policy {
	p = source.NormalizePolicy(p)
	if p.Style != "" && systemRegistry.Lookup(p.Style) == nil {
		p.Style = ""
	}
	return p
}

// fontValidator accepts "auto" or any non-empty string for the font option.
func fontValidator(value string) string {
	if strings.TrimSpace(value) == "" {
		return "must be \"auto\" or a registered font ID"
	}
	return ""
}

// queryResponse builds the "OK system style=... font=..." response.
func queryResponse() string {
	p := GetPolicy()
	return fmt.Sprintf("OK system style=%s font=%s", p.Style, p.Font)
}

// policyState holds the mutex-protected active policy.
var policyState = struct {
	sync.RWMutex
	policy Policy
}{
	policy: DefaultPolicy(),
}
