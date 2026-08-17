package wifi

import (
	"strings"
	"sync"

	"github.com/databeast/cyberhud/display/catalog"
	"github.com/databeast/cyberhud/display/modes/wifi/source"
)

// wifiSnapshotter implements catalog.PolicySnapshotter for the wifi mode.
type wifiSnapshotter struct{}

// SnapshotPolicy returns the current wifi policy as a JSON-serializable map
// using snake_case keys matching the protocol wire format.
func (wifiSnapshotter) SnapshotPolicy() map[string]interface{} {
	p := GetPolicy()
	return map[string]interface{}{
		"style":          p.Style,
		"fgcolor":        p.FGColor,
		"signal_display": p.SignalDisplay,
		"show_frequency": p.ShowFrequency,
		"show_interface": p.ShowInterface,
		"show_channel":   p.ShowChannel,
	}
}

// RestorePolicy applies policy values from a JSON map, running them through
// the same normalization as SetPolicy. Keys use the protocol wire format.
func (wifiSnapshotter) RestorePolicy(data map[string]interface{}) error {
	p := source.DefaultPolicy()

	if v, ok := data["style"]; ok {
		if s, ok := v.(string); ok {
			p.Style = s
		}
	}
	if v, ok := data["fgcolor"]; ok {
		if s, ok := v.(string); ok {
			p.FGColor = s
		}
	}
	if v, ok := data["signal_display"]; ok {
		if s, ok := v.(string); ok {
			p.SignalDisplay = s
		}
	}
	if v, ok := data["show_frequency"]; ok {
		if b, ok := toBool(v); ok {
			p.ShowFrequency = b
		}
	}
	if v, ok := data["show_interface"]; ok {
		if b, ok := toBool(v); ok {
			p.ShowInterface = b
		}
	}
	if v, ok := data["show_channel"]; ok {
		if b, ok := toBool(v); ok {
			p.ShowChannel = b
		}
	}

	SetPolicy(p)
	return nil
}

// toBool extracts a bool from an interface value.
func toBool(v interface{}) (bool, bool) {
	b, ok := v.(bool)
	return b, ok
}

func init() {
	catalog.RegisterSnapshotter("wifi", wifiSnapshotter{})
}

// policyState holds the mutex-protected active policy.
var policyState = struct {
	sync.RWMutex
	policy source.Policy
}{
	policy: source.DefaultPolicy(),
}

// GetPolicy returns the current WiFi policy (thread-safe read under RWMutex).
func GetPolicy() source.Policy {
	policyState.RLock()
	defer policyState.RUnlock()
	return policyState.policy
}

// SetPolicy updates the WiFi policy after normalization (thread-safe write under Mutex).
func SetPolicy(p source.Policy) {
	policyState.Lock()
	defer policyState.Unlock()
	policyState.policy = normalizePolicy(p)
}

// isAllowed checks if value is present in the allowed list.
func normalizePolicy(p source.Policy) source.Policy {
	d := source.DefaultPolicy()
	p.Style = strings.ToLower(strings.TrimSpace(p.Style))
	p.FGColor = strings.ToLower(strings.TrimSpace(p.FGColor))
	p.SignalDisplay = strings.ToLower(strings.TrimSpace(p.SignalDisplay))
	if p.Style != "" && wifiRegistry.Lookup(p.Style) == nil {
		p.Style = ""
	}
	if !isAllowed(p.FGColor, source.AllowedFGColors) {
		p.FGColor = d.FGColor
	}
	if !isAllowed(p.SignalDisplay, source.AllowedSignalDisplay) {
		p.SignalDisplay = d.SignalDisplay
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

// boolStr returns "true" or "false" for a bool value.
func boolStr(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
