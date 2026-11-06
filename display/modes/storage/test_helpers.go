package storage

import (
	"github.com/databeast/cyberhud/display/modes/storage/source"
	"github.com/databeast/cyberhud/display/style"
)

// StorageRegistryEnumerate exposes registered storage styles for snapshot tests.
func StorageRegistryEnumerate() []style.Style[source.StorageSnapshot, source.Policy] {
	return storageRegistry.Enumerate()
}

// AllowedStyleNames returns the registered style names (for tests/tooling).
func AllowedStyleNames() []string { return allowedStyleNames() }

// NormalizePolicy exposes policy normalization for tests.
func NormalizePolicy(p Policy) Policy { return normalizePolicy(p) }

// UpdateSnapshot stores a snapshot in the source state (for tests/tooling).
func UpdateSnapshot(snap StorageSnapshot) { source.UpdateSnapshot(snap) }

// ResetSnapshotState clears the stored snapshot (for tests).
func ResetSnapshotState() { source.ResetSnapshotStateForTest() }
