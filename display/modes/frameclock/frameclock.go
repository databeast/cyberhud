// Package frameclock provides the time source that display modes read when
// deciding whether to advance animation state.
//
// # What problem this solves
//
// Modes that animate — the ticker's auto-scroll, the pager's smooth scroll and
// page transitions — gate advancement on elapsed wall-clock time. A snapshot
// test renders a handful of frames in a tight loop, which completes in
// microseconds, so no gate ever elapses and every captured frame shows the
// un-advanced initial state. Tests worked around this by reaching into mode
// internals and driving animation by hand, which meant the captured image
// reflected test simulation rather than production behaviour.
//
// Reading time through this package lets a test freeze the clock at a fixed
// instant and step it forward by a chosen interval before each render pass. The
// production tick path then advances on its own, exactly as it does on hardware,
// and the result is deterministic and independent of how fast the machine runs.
//
// # Production behaviour
//
// Until Freeze is called, Now returns time.Now. Nothing about the default path
// changes, and no mode needs to know whether it is under test.
package frameclock

import (
	"sync"
	"time"
)

// state holds the active time source.
//
// frozen is nil in production, in which case Now falls through to time.Now.
// The mutex guards both fields because modes read the clock from goroutines
// other than the one that froze it: the pager runs a background reader, and
// the ticker reads under its own feed-state lock.
var state struct {
	mu     sync.RWMutex
	frozen *time.Time
}

// snapshotMu serializes freeze/restore cycles.
//
// The frozen clock is process-wide, so two concurrent freezes would interleave
// their advances and produce nondeterministic output for both. Freeze blocks
// until any in-progress cycle has restored, which makes concurrent snapshot
// tests wait rather than corrupt each other.
var snapshotMu sync.Mutex

// Now returns the current time from the active source.
//
// Modes call this instead of time.Now so their animation gates can be stepped
// deterministically under test.
func Now() time.Time {
	state.mu.RLock()
	frozen := state.frozen
	state.mu.RUnlock()
	if frozen != nil {
		return *frozen
	}
	return time.Now()
}

// Frozen reports whether a fixed time source is currently installed.
func Frozen() bool {
	state.mu.RLock()
	defer state.mu.RUnlock()
	return state.frozen != nil
}

// Freeze installs start as the fixed time source and returns a function that
// restores the wall clock.
//
// Freeze blocks while another freeze cycle is active. The returned restore
// function is safe to call more than once; only the first call has an effect,
// so a deferred restore and an explicit one cannot double-release the lock.
func Freeze(start time.Time) (restore func()) {
	snapshotMu.Lock()

	state.mu.Lock()
	at := start
	state.frozen = &at
	state.mu.Unlock()

	var once sync.Once
	return func() {
		once.Do(func() {
			state.mu.Lock()
			state.frozen = nil
			state.mu.Unlock()
			snapshotMu.Unlock()
		})
	}
}

// Advance moves the frozen clock forward by d.
//
// Advance is a no-op when the clock is not frozen, so a mode or helper that
// calls it outside a snapshot cannot perturb production timing. Negative
// durations are ignored: time moving backwards would make elapsed-time
// arithmetic in the modes produce negative deltas, which they treat as "no
// advancement" and which no real clock produces.
func Advance(d time.Duration) {
	if d <= 0 {
		return
	}
	state.mu.Lock()
	defer state.mu.Unlock()
	if state.frozen == nil {
		return
	}
	next := state.frozen.Add(d)
	state.frozen = &next
}
