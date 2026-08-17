package catalog

import "sync"

// PolicySnapshotter is implemented by mode packages that support policy persistence.
// ALL registered modes MUST implement this interface — use StubSnapshotter for
// modes without configurable parameters.
type PolicySnapshotter interface {
	// SnapshotPolicy returns the current policy as a JSON-serializable map.
	SnapshotPolicy() map[string]interface{}
	// RestorePolicy applies policy values from a JSON map (with normalization).
	RestorePolicy(data map[string]interface{}) error
}

// StubSnapshotter is a no-op implementation for modes without configurable params.
type StubSnapshotter struct{}

func (StubSnapshotter) SnapshotPolicy() map[string]interface{}     { return map[string]interface{}{} }
func (StubSnapshotter) RestorePolicy(map[string]interface{}) error { return nil }

var (
	snapshottersMu sync.RWMutex
	snapshotters   = map[string]PolicySnapshotter{}
)

// RegisterSnapshotter associates a mode ID with its PolicySnapshotter implementation.
func RegisterSnapshotter(modeID string, s PolicySnapshotter) {
	if modeID == "" || s == nil {
		return
	}
	snapshottersMu.Lock()
	defer snapshottersMu.Unlock()
	snapshotters[modeID] = s
}

// Snapshotter returns the registered PolicySnapshotter for a mode, or nil.
func Snapshotter(modeID string) PolicySnapshotter {
	snapshottersMu.RLock()
	defer snapshottersMu.RUnlock()
	return snapshotters[modeID]
}
