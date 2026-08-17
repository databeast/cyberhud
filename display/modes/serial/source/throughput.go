package source

import (
	"sync"
	"time"
)

// ThroughputTracker maintains a sliding window of 32 bytes-per-second measurements.
type ThroughputTracker struct {
	mu       sync.Mutex
	window   [32]int   // Circular buffer of bytes-per-second values
	head     int       // Next write position
	count    int       // Number of populated entries (max 32)
	current  int       // Accumulated bytes in the current interval
	lastTick time.Time // Start of the current 1-second interval
}

// Add records received bytes into the current interval accumulator.
func (t *ThroughputTracker) Add(n int) {
	t.mu.Lock()
	t.current += n
	t.mu.Unlock()
}

// Tick advances the sampling window if 1 second has elapsed.
// Called from the read loop or a periodic timer.
func (t *ThroughputTracker) Tick(now time.Time) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.lastTick.IsZero() {
		t.lastTick = now
		return
	}

	// Advance one slot for each full second elapsed.
	for now.Sub(t.lastTick) >= time.Second {
		t.window[t.head] = t.current
		t.head = (t.head + 1) % 32
		if t.count < 32 {
			t.count++
		}
		t.current = 0
		t.lastTick = t.lastTick.Add(time.Second)
	}
}

// History returns the 32-entry sliding window as an array (oldest first).
func (t *ThroughputTracker) History() [32]int {
	t.mu.Lock()
	defer t.mu.Unlock()

	var out [32]int
	if t.count == 0 {
		return out
	}

	// The oldest entry is at (head - count) mod 32.
	start := (t.head - t.count + 32) % 32
	for i := 0; i < 32; i++ {
		if i < t.count {
			out[i] = t.window[(start+i)%32]
		}
		// Unpopulated entries stay 0.
	}
	return out
}

// Reset clears all 32 entries to zero (called on disconnect/port change).
func (t *ThroughputTracker) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.window = [32]int{}
	t.head = 0
	t.count = 0
	t.current = 0
	t.lastTick = time.Time{}
}

// Normalize returns the throughput history scaled to [0.0, 1.0]
// relative to the configured baud rate (maxBPS = baud / 10 for 8N1).
func (t *ThroughputTracker) Normalize(maxBPS int) []float64 {
	hist := t.History()
	out := make([]float64, 32)
	if maxBPS <= 0 {
		return out
	}
	for i, v := range hist {
		ratio := float64(v) / float64(maxBPS)
		if ratio > 1.0 {
			ratio = 1.0
		}
		if ratio < 0.0 {
			ratio = 0.0
		}
		out[i] = ratio
	}
	return out
}
