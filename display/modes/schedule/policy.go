package schedule

import (
	"strings"

	"github.com/databeast/cyberhud/display/catalog"
	"github.com/databeast/cyberhud/display/catalog/cmdutil"
	"github.com/databeast/cyberhud/display/modes/schedule/source"
)

type Policy = source.Policy
type Snapshot = source.Snapshot

// DefaultPolicy returns the baseline schedule policy.
func DefaultPolicy() Policy { return source.DefaultPolicy() }

// normalizePolicy applies source normalization then drops an unknown style name
// so fitness resolution takes over.
func normalizePolicy(p Policy) Policy {
	p = source.NormalizePolicy(p)
	p.Style = strings.ToLower(strings.TrimSpace(p.Style))
	if p.Style != "" && scheduleRegistry.Lookup(p.Style) == nil {
		p.Style = ""
	}
	return p
}

// PolicySnapshot returns the current schedule policy.
func PolicySnapshot() Policy { return source.PolicySnapshot() }

// SetPolicy updates the schedule policy.
func SetPolicy(policy Policy) { source.SetPolicy(normalizePolicy(policy)) }

// scheduleSnapshotter implements catalog.PolicySnapshotter for persistence.
type scheduleSnapshotter struct{}

// SnapshotPolicy returns the current policy as a JSON-serializable map using
// snake_case wire keys.
func (scheduleSnapshotter) SnapshotPolicy() map[string]interface{} {
	return PolicySnapshot().ToMap()
}

// RestorePolicy applies policy values from a JSON map, running them through the
// same normalization as SetPolicy. Keys use the protocol wire format.
func (scheduleSnapshotter) RestorePolicy(data map[string]interface{}) error {
	p := DefaultPolicy()

	if v, ok := data["source"]; ok {
		if s, ok := v.(string); ok {
			p.Source = s
		}
	}
	if v, ok := data["refresh_sec"]; ok {
		if n, ok := cmdutil.ToInt(v); ok {
			p.RefreshSec = n
		}
	}
	if v, ok := data["time_format"]; ok {
		if s, ok := v.(string); ok {
			p.TimeFormat = s
		}
	}
	if v, ok := data["show_upcoming"]; ok {
		if b, ok := cmdutil.ToBool(v); ok {
			p.ShowUpcoming = b
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

	SetPolicy(p)
	return nil
}

func init() {
	catalog.RegisterSnapshotter("schedule", scheduleSnapshotter{})
}
