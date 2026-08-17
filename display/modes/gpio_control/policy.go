package gpio_control

import (
	"strings"
	"sync"

	"github.com/databeast/cyberhud/display/catalog"
	"github.com/databeast/cyberhud/display/modes/gpio_control/source"
	gpiomgr "github.com/databeast/cyberhud/hardware/gpio"
)

type Policy = source.Policy

func DefaultPolicy() Policy { return source.DefaultPolicy() }

// gpioControlSnapshotter implements catalog.PolicySnapshotter for the gpio_control mode.
type gpioControlSnapshotter struct{}

var state = struct {
	sync.RWMutex
	pins []gpiomgr.PinState
	mgr  *gpiomgr.Manager
}{
	pins: []gpiomgr.PinState{},
}

// SnapshotPolicy returns the current gpio_control policy as a JSON-serializable map
// using snake_case keys matching the protocol wire format.
func (gpioControlSnapshotter) SnapshotPolicy() map[string]interface{} {
	p := PolicySnapshot()
	return map[string]interface{}{
		"style": p.Style,
		"font":  p.Font,
	}
}

// RestorePolicy applies policy values from a JSON map, running them through
// the same normalization as SetPolicy. Keys use the protocol wire format.
func (gpioControlSnapshotter) RestorePolicy(data map[string]interface{}) error {
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
	catalog.RegisterSnapshotter("gpio-control", gpioControlSnapshotter{})
}

var (
	policyMu sync.RWMutex
	policy   = DefaultPolicy()
)

// PolicySnapshot returns a copy of the current policy.
func PolicySnapshot() Policy {
	policyMu.RLock()
	defer policyMu.RUnlock()
	return policy
}

// SetPolicy updates the current policy.
func SetPolicy(p Policy) {
	policyMu.Lock()
	defer policyMu.Unlock()
	policy = normalizePolicy(p)
}

// normalizePolicy validates and normalizes the policy fields.
// Empty style means auto-detect via fitness.
func normalizePolicy(p Policy) Policy {
	p.Style = strings.ToLower(strings.TrimSpace(p.Style))
	p.Font = strings.TrimSpace(p.Font)
	if p.Style != "" && gpioControlRegistry.Lookup(p.Style) == nil {
		p.Style = ""
	}
	if p.Font == "" {
		p.Font = "auto"
	}
	return p
}

// fontValidator accepts "auto" or any non-empty trimmed string for the font option.
// Invalid font IDs are handled at render time by falling back to auto.
func fontValidator(value string) string {
	if strings.TrimSpace(value) == "" {
		return "must be \"auto\" or a non-empty font ID"
	}
	return ""
}
