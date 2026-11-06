package serial

import (
	"strings"

	"github.com/databeast/cyberhud/display/catalog"
	"github.com/databeast/cyberhud/display/catalog/cmdutil"
	"github.com/databeast/cyberhud/display/modes/serial/source"
)

type Policy = source.Policy
type Snapshot = source.Snapshot

// serialRegistry and registeredStyleNames live in styles.go.

func DefaultPolicy() Policy { return source.DefaultPolicy() }

// normalizePolicy ensures policy fields contain valid values and registered styles.
func normalizePolicy(p Policy) Policy {
	p = source.NormalizePolicy(p)
	p.Style = strings.ToLower(strings.TrimSpace(p.Style))
	if p.Style != "" && serialRegistry.Lookup(p.Style) == nil {
		p.Style = ""
	}
	return p
}

// serialSnapshotter implements catalog.PolicySnapshotter for the serial mode.
type serialSnapshotter struct{}

// SnapshotPolicy returns the current serial policy as a JSON-serializable map
// using snake_case keys matching the protocol wire format.
func (serialSnapshotter) SnapshotPolicy() map[string]interface{} {
	return PolicySnapshot().ToMap()
}

// RestorePolicy applies policy values from a JSON map, running them through
// the same normalization as SetPolicy. Keys use the protocol wire format.
func (serialSnapshotter) RestorePolicy(data map[string]interface{}) error {
	p := DefaultPolicy()

	if v, ok := data["port"]; ok {
		if s, ok := v.(string); ok {
			p.Port = s
		}
	}
	if v, ok := data["baud"]; ok {
		if n, ok := cmdutil.ToInt(v); ok {
			p.Baud = n
		}
	}
	if v, ok := data["lines"]; ok {
		if n, ok := cmdutil.ToInt(v); ok {
			p.MaxLines = n
		}
	}
	if v, ok := data["autoselect"]; ok {
		if b, ok := cmdutil.ToBool(v); ok {
			p.AutoSelect = b
		}
	}
	if v, ok := data["scan_ms"]; ok {
		if n, ok := cmdutil.ToInt(v); ok {
			p.ScanMS = n
		}
	}
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
	if v, ok := data["raw"]; ok {
		if b, ok := cmdutil.ToBool(v); ok {
			p.Raw = b
		}
	}

	SetPolicy(p)
	return nil
}

func init() {
	catalog.RegisterSnapshotter("serial", serialSnapshotter{})
}

// PolicySnapshot returns the current serial monitor policy.
func PolicySnapshot() Policy {
	return source.PolicySnapshot()
}

// SetPolicy updates the serial monitor policy.
func SetPolicy(policy Policy) {
	source.SetPolicy(normalizePolicy(policy))
}
