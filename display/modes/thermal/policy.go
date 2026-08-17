package thermal

import (
	"fmt"
	"strings"
	"sync"

	"github.com/databeast/cyberhud/display/catalog"
	"github.com/databeast/cyberhud/display/modes/thermal/source"
)

type Policy = source.Policy
type ThermalSnapshot = source.ThermalSnapshot
type ZoneReading = source.ZoneReading
type TripPoint = source.TripPoint

func DefaultPolicy() Policy { return source.DefaultPolicy() }

// thermalSnapshotter implements catalog.PolicySnapshotter for the thermal mode.
type thermalSnapshotter struct{}

// SnapshotPolicy returns the current thermal policy as a JSON-serializable map
// using snake_case keys matching the protocol wire format.
func (thermalSnapshotter) SnapshotPolicy() map[string]interface{} {
	p := GetPolicy()
	return map[string]interface{}{
		"style":            p.Style,
		"font":             p.Font,
		"refresh_ms":       p.RefreshMS,
		"warn_threshold":   p.WarnThreshold,
		"crit_threshold":   p.CritThreshold,
		"show_border":      p.ShowBorder,
		"unit":             p.Unit,
		"fgcolor":          p.FGColor,
		"show_led":         p.ShowLED,
		"show_refresh_bar": p.ShowRefreshBar,
	}
}

// RestorePolicy applies policy values from a JSON map, running them through
// the same normalization as SetPolicy. Keys use the protocol wire format.
func (thermalSnapshotter) RestorePolicy(data map[string]interface{}) error {
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
	if v, ok := data["refresh_ms"]; ok {
		if n, ok := toInt(v); ok {
			p.RefreshMS = n
		}
	}
	if v, ok := data["warn_threshold"]; ok {
		if n, ok := toInt(v); ok {
			p.WarnThreshold = n
		}
	}
	if v, ok := data["crit_threshold"]; ok {
		if n, ok := toInt(v); ok {
			p.CritThreshold = n
		}
	}
	if v, ok := data["show_border"]; ok {
		if b, ok := toBool(v); ok {
			p.ShowBorder = b
		}
	}
	if v, ok := data["unit"]; ok {
		if s, ok := v.(string); ok {
			p.Unit = s
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
	if v, ok := data["show_refresh_bar"]; ok {
		if b, ok := toBool(v); ok {
			p.ShowRefreshBar = b
		}
	}

	return SetPolicy(p)
}

// toInt extracts an int from an interface value, handling both float64
// (JSON numbers decode as float64) and direct int types.
func toInt(v interface{}) (int, bool) {
	switch n := v.(type) {
	case float64:
		return int(n), true
	case int:
		return n, true
	case int64:
		return int(n), true
	default:
		return 0, false
	}
}

// toBool extracts a bool from an interface value.
func toBool(v interface{}) (bool, bool) {
	b, ok := v.(bool)
	return b, ok
}

func init() {
	catalog.RegisterSnapshotter("thermal", thermalSnapshotter{})
}

// policyState holds the mutex-protected active policy.
var policyState = struct {
	sync.RWMutex
	policy Policy
}{
	policy: DefaultPolicy(),
}

// GetPolicy returns the current thermal policy (thread-safe).
func GetPolicy() Policy {
	policyState.RLock()
	defer policyState.RUnlock()
	return policyState.policy
}

// SetPolicy updates the thermal policy after normalization (thread-safe).
// Returns an error if the policy fails validation (e.g. warn_threshold >= crit_threshold).
func SetPolicy(p Policy) error {
	normalized, err := normalizePolicy(p)
	if err != nil {
		return err
	}
	policyState.Lock()
	defer policyState.Unlock()
	policyState.policy = normalized
	return nil
}

// normalizePolicy ensures the policy fields contain valid values, with
// registry-aware style validation kept in the controller layer.
// Returns an error if warn_threshold >= crit_threshold after clamping.
func normalizePolicy(p Policy) (Policy, error) {
	d := source.DefaultPolicy()

	p.Style = strings.ToLower(strings.TrimSpace(p.Style))
	if p.Style != "" && thermalRegistry.Lookup(p.Style) == nil {
		p.Style = ""
	}

	p.Font = strings.ToLower(strings.TrimSpace(p.Font))
	if p.Font == "" {
		p.Font = d.Font
	}

	p.RefreshMS = clampRefreshMS(p.RefreshMS)

	if p.WarnThreshold < 0 {
		p.WarnThreshold = 0
	}
	if p.CritThreshold < 0 {
		p.CritThreshold = 0
	}

	// Reject invalid threshold pair (requirement 2.6).
	if p.WarnThreshold >= p.CritThreshold {
		return Policy{}, fmt.Errorf("warn_threshold (%d) must be < crit_threshold (%d)", p.WarnThreshold, p.CritThreshold)
	}

	p.Unit = strings.ToUpper(strings.TrimSpace(p.Unit))
	if !isAllowed(p.Unit, source.AllowedUnits) {
		p.Unit = d.Unit
	}

	p.FGColor = strings.ToLower(strings.TrimSpace(p.FGColor))
	if !isAllowed(p.FGColor, source.AllowedFGColors) {
		p.FGColor = d.FGColor
	}

	return p, nil
}
