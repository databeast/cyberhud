package systemd

import (
	"strings"
	"sync"

	"github.com/databeast/cyberhud/display/catalog"
	"github.com/databeast/cyberhud/display/modes/systemd/source"
)

// systemdSnapshotter implements catalog.PolicySnapshotter for the systemd mode.
type systemdSnapshotter struct{}

func (systemdSnapshotter) SnapshotPolicy() map[string]interface{} {
	return GetPolicy().ToMap()
}

func (systemdSnapshotter) RestorePolicy(data map[string]interface{}) error {
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

func init() {
	catalog.RegisterSnapshotter("systemd", systemdSnapshotter{})
}

var policyState = struct {
	sync.RWMutex
	policy source.Policy
}{
	policy: source.DefaultPolicy(),
}

func GetPolicy() source.Policy {
	policyState.RLock()
	defer policyState.RUnlock()
	return policyState.policy
}

func SetPolicy(p source.Policy) {
	policyState.Lock()
	defer policyState.Unlock()
	policyState.policy = normalizePolicy(p)
}

func normalizePolicy(p source.Policy) source.Policy {
	d := source.DefaultPolicy()
	p.Style = strings.ToLower(strings.TrimSpace(p.Style))
	if p.Style != "" && systemdRegistry.Lookup(p.Style) == nil {
		p.Style = ""
	}
	p.ColorAccent = strings.ToLower(strings.TrimSpace(p.ColorAccent))
	if !isAllowed(p.ColorAccent, source.AllowedAccents) {
		p.ColorAccent = d.ColorAccent
	}
	return p
}

func isAllowed(value string, allowed []string) bool {
	for _, a := range allowed {
		if value == a {
			return true
		}
	}
	return false
}
