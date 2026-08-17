package source

import "time"

// Snapshot captures the current USB bench state.
type Snapshot struct {
	Sequence        uint64
	Connected       bool
	HasLast         bool
	LastConnectedAt time.Time
	Device          DeviceInfo
}
