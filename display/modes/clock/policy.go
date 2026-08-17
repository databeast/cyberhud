package clock

import (
	"strings"
	"sync"

	"github.com/databeast/cyberhud/display/catalog"
	"github.com/databeast/cyberhud/display/modes/clock/source"
)

type Policy = source.Policy

func DefaultPolicy() Policy { return source.DefaultPolicy() }

// policyState holds the mutex-protected active policy.
var policyState = struct {
	sync.RWMutex
	policy source.Policy
}{
	policy: source.DefaultPolicy(),
}

// GetPolicy returns the current clock policy (thread-safe read under RWMutex).
//
// Framework pattern demonstrated: policy definition — concurrent-safe read access
// to the shared policy state, enabling render loops and CLI queries to observe
// consistent snapshots without blocking writers.
func GetPolicy() source.Policy {
	policyState.RLock()
	defer policyState.RUnlock()
	return policyState.policy
}

// SetPolicy updates the clock policy after normalization (thread-safe write under Mutex).
//
// Framework pattern demonstrated: policy normalization — validates and normalizes
// incoming policy values before persisting them, ensuring the active policy always
// contains well-formed data regardless of caller input.
func SetPolicy(p source.Policy) {
	policyState.Lock()
	defer policyState.Unlock()
	policyState.policy = normalizePolicy(p)
}

// normalizePolicy ensures the policy fields contain valid values,
// falling back to defaults for invalid entries.
func normalizePolicy(p source.Policy) source.Policy {
	d := source.DefaultPolicy()

	p.Style = strings.ToLower(strings.TrimSpace(p.Style))
	p.FGColor = strings.ToLower(strings.TrimSpace(p.FGColor))
	p.SecondsBar = strings.ToLower(strings.TrimSpace(p.SecondsBar))
	p.BorderColor = strings.ToLower(strings.TrimSpace(p.BorderColor))

	if p.Style != "" && clockRegistry.Lookup(p.Style) == nil {
		p.Style = ""
	}

	if !isAllowed(p.TimeFormat, source.AllowedTimeFormats) {
		p.TimeFormat = d.TimeFormat
	}

	if !isAllowed(p.DateFormat, source.AllowedDateFormats) {
		p.DateFormat = d.DateFormat
	}

	if !isAllowed(p.SecondsBar, source.AllowedSecondsBar) {
		p.SecondsBar = d.SecondsBar
	}

	if !isAllowed(p.BorderColor, source.AllowedBorderColors) {
		p.BorderColor = "auto"
	}

	if strings.TrimSpace(p.Timezone) == "" {
		p.Timezone = "local"
	}

	return p
}

// isAllowed checks if value is present in the allowed list.
func isAllowed(value string, allowed []string) bool {
	for _, a := range allowed {
		if value == a {
			return true
		}
	}
	return false
}

// clockSnapshotter implements catalog.PolicySnapshotter for the clock mode.
type clockSnapshotter struct{}

// SnapshotPolicy returns the current clock policy as a JSON-serializable map
// using snake_case keys matching the protocol wire format.
func (clockSnapshotter) SnapshotPolicy() map[string]interface{} {
	p := GetPolicy()
	return p.ToMap()
}

// RestorePolicy applies policy values from a JSON map, running them through
// the same normalization as SetPolicy. Keys use the protocol wire format.
func (clockSnapshotter) RestorePolicy(data map[string]interface{}) error {
	p := source.DefaultPolicy()

	if v, ok := data["style"]; ok {
		if s, ok := v.(string); ok {
			p.Style = s
		}
	}
	if v, ok := data["show_seconds"]; ok {
		if b, ok := toBool(v); ok {
			p.ShowSeconds = b
		}
	}
	if v, ok := data["time_format"]; ok {
		if s, ok := v.(string); ok {
			p.TimeFormat = s
		}
	}
	if v, ok := data["date_format"]; ok {
		if s, ok := v.(string); ok {
			p.DateFormat = s
		}
	}
	if v, ok := data["timezone"]; ok {
		if s, ok := v.(string); ok {
			p.Timezone = s
		}
	}
	if v, ok := data["show_weekday"]; ok {
		if b, ok := toBool(v); ok {
			p.ShowWeekday = b
		}
	}
	if v, ok := data["blink_colon"]; ok {
		if b, ok := toBool(v); ok {
			p.BlinkColon = b
		}
	}
	if v, ok := data["fgcolor"]; ok {
		if s, ok := v.(string); ok {
			p.FGColor = s
		}
	}
	if v, ok := data["show_led"]; ok {
		if b, ok := toBool(v); ok {
			p.ShowLED = b
		}
	}
	if v, ok := data["seconds_bar"]; ok {
		if s, ok := v.(string); ok {
			p.SecondsBar = s
		}
	}
	if v, ok := data["show_daybar"]; ok {
		if b, ok := toBool(v); ok {
			p.ShowDaybar = b
		}
	}
	if v, ok := data["show_border"]; ok {
		if b, ok := toBool(v); ok {
			p.ShowBorder = b
		}
	}
	if v, ok := data["border_color"]; ok {
		if s, ok := v.(string); ok {
			p.BorderColor = s
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
	catalog.RegisterSnapshotter("clock", clockSnapshotter{})
}
