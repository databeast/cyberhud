package usb

import (
	"strings"
	"time"

	"github.com/databeast/cyberhud/display/catalog"
	"github.com/databeast/cyberhud/display/modes/usb/source"
)

type Policy = source.Policy
type Snapshot = source.Snapshot

func DefaultPolicy() Policy { return source.DefaultPolicy() }

// usbSnapshotter implements catalog.PolicySnapshotter for the usb mode.
type usbSnapshotter struct{}

// SnapshotPolicy returns the current usb policy as a JSON-serializable map
// using snake_case keys matching the protocol wire format.
func (usbSnapshotter) SnapshotPolicy() map[string]interface{} {
	p := PolicySnapshot()
	return map[string]interface{}{
		"poll_ms":           p.PollMS,
		"hold_unplugged_ms": p.HoldUnpluggedMS,
		"hide_root_hubs":    p.HideRootHubs,
		"style":             p.Style,
	}
}

// RestorePolicy applies policy values from a JSON map, running them through
// the same normalization as SetPolicy. Keys use the protocol wire format.
func (usbSnapshotter) RestorePolicy(data map[string]interface{}) error {
	p := DefaultPolicy()

	if v, ok := data["poll_ms"]; ok {
		if n, ok := toInt(v); ok {
			p.PollMS = n
		}
	}
	if v, ok := data["hold_unplugged_ms"]; ok {
		if n, ok := toInt(v); ok {
			p.HoldUnpluggedMS = n
		}
	}
	if v, ok := data["hide_root_hubs"]; ok {
		if b, ok := toBool(v); ok {
			p.HideRootHubs = b
		}
	}
	if v, ok := data["style"]; ok {
		if s, ok := v.(string); ok {
			p.Style = s
		}
	}

	SetPolicy(p)
	return nil
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
	catalog.RegisterSnapshotter("usb", usbSnapshotter{})
}

func normalizePolicy(policy source.Policy) source.Policy {
	d := source.DefaultPolicy()
	if policy.PollMS < 100 {
		policy.PollMS = d.PollMS
	}
	if policy.HoldUnpluggedMS < 0 {
		policy.HoldUnpluggedMS = d.HoldUnpluggedMS
	}
	// Normalize style using registry lookup; empty means auto-detect via fitness.
	policy.Style = strings.ToLower(strings.TrimSpace(policy.Style))
	if policy.Style != "" && usbRegistry.Lookup(policy.Style) == nil {
		policy.Style = ""
	}
	return policy
}

func applyHoldPolicy(snapshot *Snapshot, now time.Time, policy Policy) bool {
	if snapshot == nil || !snapshot.HasLast || snapshot.Connected || policy.HoldUnpluggedMS <= 0 {
		return false
	}
	if now.Sub(snapshot.LastConnectedAt) < time.Duration(policy.HoldUnpluggedMS)*time.Millisecond {
		return false
	}
	snapshot.HasLast = false
	snapshot.Connected = false
	snapshot.Device = source.DeviceInfo{}
	snapshot.LastConnectedAt = time.Time{}
	return true
}
