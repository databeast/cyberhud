package storage

import (
	"sync"

	"github.com/databeast/cyberhud/display/catalog"
	"github.com/databeast/cyberhud/display/catalog/cmdutil"
	"github.com/databeast/cyberhud/display/modes/storage/source"
)

type Policy = source.Policy
type StorageSnapshot = source.StorageSnapshot
type MountUsage = source.MountUsage

// DefaultPolicy returns the default storage policy.
func DefaultPolicy() Policy { return source.DefaultPolicy() }

// storageSnapshotter implements catalog.PolicySnapshotter for the storage mode.
type storageSnapshotter struct{}

// SnapshotPolicy returns the current storage policy as a JSON-serializable map
// using snake_case keys matching the protocol wire format.
func (storageSnapshotter) SnapshotPolicy() map[string]interface{} {
	return GetPolicy().ToMap()
}

// RestorePolicy applies policy values from a JSON map, running them through the
// same normalization as SetPolicy. Keys use the protocol wire format.
func (storageSnapshotter) RestorePolicy(data map[string]interface{}) error {
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
		if n, ok := cmdutil.ToInt(v); ok {
			p.RefreshMS = n
		}
	}
	if v, ok := data["unit"]; ok {
		if s, ok := v.(string); ok {
			p.Unit = s
		}
	}
	if v, ok := data["show_all"]; ok {
		if b, ok := v.(bool); ok {
			p.ShowAll = b
		}
	}

	SetPolicy(p)
	return nil
}

func init() {
	catalog.RegisterSnapshotter("storage", storageSnapshotter{})
}

// policyState holds the mutex-protected active policy.
var policyState = struct {
	sync.RWMutex
	policy Policy
}{
	policy: DefaultPolicy(),
}

// GetPolicy returns the current storage policy (thread-safe).
func GetPolicy() Policy {
	policyState.RLock()
	defer policyState.RUnlock()
	return policyState.policy
}

// SetPolicy updates the storage policy after normalization (thread-safe).
func SetPolicy(p Policy) {
	policyState.Lock()
	defer policyState.Unlock()
	policyState.policy = normalizePolicy(p)
}

// normalizePolicy performs registry-aware normalization: it runs the
// source-level normalization, then nulls an unknown style so resolution falls
// back to fitness.
func normalizePolicy(p Policy) Policy {
	p = source.NormalizePolicy(p)
	if p.Style != "" && storageRegistry.Lookup(p.Style) == nil {
		p.Style = ""
	}
	return p
}
