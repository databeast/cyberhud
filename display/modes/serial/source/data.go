package source

import "time"

// Snapshot is a read-only view of the live serial monitor state.
type Snapshot struct {
	Sequence      uint64
	Connected     bool
	Port          string
	Baud          int
	AutoSelect    bool
	LastError     string
	ErrorCategory ErrorCategory // Classified error type
	Lines         []string
	LineColors    [][]ColorSegment // Per-line ANSI color segments
	Throughput    [32]int          // Bytes-per-second sliding window (oldest first)
	ScrollOffset  int              // Current scroll position
	UpdatedAt     time.Time
}
