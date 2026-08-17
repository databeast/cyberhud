package tests

import (
	"image"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/databeast/cyberhud/display/region"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/runtime/action"
)

// --- Mock ModeInstance for lifecycle testing ---

// lifecycleMock records Activate/Deactivate calls with sequence numbers
// to verify ordering guarantees.
type lifecycleMock struct {
	id              string
	activateCount   int
	deactivateCount int
	activateSeq     int // global sequence number when Activate was called
	deactivateSeq   int // global sequence number when Deactivate was called

	// Optional: block Deactivate for timeout testing.
	deactivateBlock time.Duration

	// Optional: panic in Activate.
	activatePanic interface{}

	seq *int32 // shared atomic counter for ordering
}

func newLifecycleMock(id string, seq *int32) *lifecycleMock {
	return &lifecycleMock{id: id, seq: seq}
}

func (m *lifecycleMock) ID() string { return m.id }

func (m *lifecycleMock) Activate() {
	if m.activatePanic != nil {
		panic(m.activatePanic)
	}
	m.activateCount++
	m.activateSeq = int(atomic.AddInt32(m.seq, 1))
}

func (m *lifecycleMock) Deactivate() {
	if m.deactivateBlock > 0 {
		time.Sleep(m.deactivateBlock)
	}
	m.deactivateCount++
	m.deactivateSeq = int(atomic.AddInt32(m.seq, 1))
}

func (m *lifecycleMock) ActionHandler() action.Handler { return nil }

func (m *lifecycleMock) BuildView() style.ViewData {
	return style.ViewData{}
}

func (m *lifecycleMock) RenderCacheKey() uint32 { return 0 }

// --- Helper to create a Region with a ModeFactory ---

func newTestSetModeRegion() (*region.Region, *int32) {
	bounds := image.Rect(0, 0, 100, 100)
	surf := surface.New(bounds)
	r := region.NewRegion("test", bounds, surf)
	seq := new(int32)
	return r, seq
}

// --- Tests ---

// TestSetMode_SuccessfulSwitch_OldDeactivatedAfterNewActivated verifies that on a
// successful mode switch, the old instance's Deactivate is called AFTER the new
// instance's Activate succeeds, and the new instance is stored.

func TestSetMode_SuccessfulSwitch_OldDeactivatedAfterNewActivated(t *testing.T) {
	r, seq := newTestSetModeRegion()

	oldInstance := newLifecycleMock("old-mode", seq)
	newInstance := newLifecycleMock("new-mode", seq)

	// Pre-set the old instance as the active one.
	r.TestSetInstance(oldInstance)
	r.TestSetMode("old-mode")

	// Configure ModeFactory to return newInstance for "new-mode".
	r.SetModeFactory(func(id string, hints textlayout.TextHints) (region.ModeInstance, bool) {
		if id == "new-mode" {
			return newInstance, true
		}
		return nil, false
	})

	err := r.SetMode("new-mode")
	if err != nil {
		t.Fatalf("SetMode() returned unexpected error: %v", err)
	}

	// Verify new instance is now stored.
	if r.Instance() != newInstance {
		t.Fatal("expected new instance to be stored as active")
	}

	// Verify mode field updated.
	if r.CurrentMode() != "new-mode" {
		t.Fatalf("CurrentMode()=%q, want %q", r.CurrentMode(), "new-mode")
	}

	// Verify Activate was called on new instance.
	if newInstance.activateCount != 1 {
		t.Fatalf("new instance Activate called %d times, want 1", newInstance.activateCount)
	}

	// Verify Deactivate was called on old instance.
	if oldInstance.deactivateCount != 1 {
		t.Fatalf("old instance Deactivate called %d times, want 1", oldInstance.deactivateCount)
	}

	// CRITICAL: Verify ordering — new Activate BEFORE old Deactivate.
	if newInstance.activateSeq >= oldInstance.deactivateSeq {
		t.Fatalf("ordering violation: new.Activate (seq=%d) should be before old.Deactivate (seq=%d)",
			newInstance.activateSeq, oldInstance.deactivateSeq)
	}
}

// TestSetMode_UnknownMode_ErrorReturnedPreviousUnchanged verifies that when SetMode
// is called with an unknown mode ID, it returns an error, leaves the previous
// instance unchanged, and does NOT call Deactivate on the previous instance.

func TestSetMode_UnknownMode_ErrorReturnedPreviousUnchanged(t *testing.T) {
	r, seq := newTestSetModeRegion()

	oldInstance := newLifecycleMock("current", seq)
	r.TestSetInstance(oldInstance)
	r.TestSetMode("current")

	// Configure ModeFactory to return false for unknown modes.
	r.SetModeFactory(func(id string, hints textlayout.TextHints) (region.ModeInstance, bool) {
		return nil, false
	})

	err := r.SetMode("nonexistent")
	if err == nil {
		t.Fatal("SetMode() should return error for unknown mode")
	}

	// Previous instance should remain unchanged.
	if r.Instance() != oldInstance {
		t.Fatal("previous instance should remain active after unknown mode error")
	}

	// Mode should remain unchanged.
	if r.CurrentMode() != "current" {
		t.Fatalf("CurrentMode()=%q, want %q", r.CurrentMode(), "current")
	}

	// Deactivate should NOT have been called on the previous instance.
	if oldInstance.deactivateCount != 0 {
		t.Fatalf("old instance Deactivate called %d times, want 0", oldInstance.deactivateCount)
	}
}

// TestSetMode_FactoryPanicRecovery verifies that when the ModeFactory (factory) panics,
// SetMode recovers, returns an error, and the previous instance remains unchanged.

func TestSetMode_FactoryPanicRecovery(t *testing.T) {
	r, seq := newTestSetModeRegion()

	oldInstance := newLifecycleMock("stable", seq)
	r.TestSetInstance(oldInstance)
	r.TestSetMode("stable")

	// Configure ModeFactory to panic.
	r.SetModeFactory(func(id string, hints textlayout.TextHints) (region.ModeInstance, bool) {
		panic("factory exploded")
	})

	err := r.SetMode("broken")
	if err == nil {
		t.Fatal("SetMode() should return error when factory panics")
	}

	// Previous instance should remain unchanged.
	if r.Instance() != oldInstance {
		t.Fatal("previous instance should remain active after factory panic")
	}

	// Mode should remain unchanged.
	if r.CurrentMode() != "stable" {
		t.Fatalf("CurrentMode()=%q, want %q", r.CurrentMode(), "stable")
	}

	// Deactivate should NOT have been called on the previous instance.
	if oldInstance.deactivateCount != 0 {
		t.Fatalf("old instance Deactivate called %d times, want 0", oldInstance.deactivateCount)
	}
}

// TestSetMode_ActivatePanicRecovery verifies that when Activate panics on the new
// instance, SetMode recovers, calls Deactivate on the new instance (cleanup), and
// retains the previous instance.

func TestSetMode_ActivatePanicRecovery(t *testing.T) {
	r, seq := newTestSetModeRegion()

	oldInstance := newLifecycleMock("stable", seq)
	r.TestSetInstance(oldInstance)
	r.TestSetMode("stable")

	panicInstance := newLifecycleMock("panicky", seq)
	panicInstance.activatePanic = "activate boom"

	// Configure ModeFactory to return the panicky instance.
	r.SetModeFactory(func(id string, hints textlayout.TextHints) (region.ModeInstance, bool) {
		if id == "panicky" {
			return panicInstance, true
		}
		return nil, false
	})

	err := r.SetMode("panicky")
	if err == nil {
		t.Fatal("SetMode() should return error when Activate panics")
	}

	// Previous instance should remain unchanged.
	if r.Instance() != oldInstance {
		t.Fatal("previous instance should remain active after Activate panic")
	}

	// Mode should remain unchanged.
	if r.CurrentMode() != "stable" {
		t.Fatalf("CurrentMode()=%q, want %q", r.CurrentMode(), "stable")
	}

	// Deactivate should have been called on the new (panicky) instance as cleanup.
	if panicInstance.deactivateCount != 1 {
		t.Fatalf("panicky instance Deactivate called %d times, want 1 (cleanup)", panicInstance.deactivateCount)
	}

	// Deactivate should NOT have been called on the old instance.
	if oldInstance.deactivateCount != 0 {
		t.Fatalf("old instance Deactivate called %d times, want 0", oldInstance.deactivateCount)
	}
}

// TestSetMode_DeactivateTimeout verifies that when the old instance's Deactivate
// blocks for longer than 5 seconds, SetMode still completes and the instance is
// dereferenced regardless.

func TestSetMode_DeactivateTimeout(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping timeout test in short mode")
	}

	r, seq := newTestSetModeRegion()

	// Old instance whose Deactivate blocks for 7 seconds (exceeds 5s timeout).
	oldInstance := newLifecycleMock("blocking", seq)
	oldInstance.deactivateBlock = 7 * time.Second
	r.TestSetInstance(oldInstance)
	r.TestSetMode("blocking")

	newInstance := newLifecycleMock("fast", seq)

	r.SetModeFactory(func(id string, hints textlayout.TextHints) (region.ModeInstance, bool) {
		if id == "fast" {
			return newInstance, true
		}
		return nil, false
	})

	// Track that SetMode completes in approximately 5 seconds (not 7).
	var wg sync.WaitGroup
	var elapsed time.Duration
	var setModeErr error

	wg.Add(1)
	go func() {
		defer wg.Done()
		start := time.Now()
		setModeErr = r.SetMode("fast")
		elapsed = time.Since(start)
	}()

	wg.Wait()

	if setModeErr != nil {
		t.Fatalf("SetMode() returned unexpected error: %v", setModeErr)
	}

	// SetMode should complete around 5 seconds (the timeout), not 7 seconds.
	if elapsed < 4*time.Second {
		t.Fatalf("SetMode completed too fast (%v), expected ~5s timeout", elapsed)
	}
	if elapsed > 6*time.Second {
		t.Fatalf("SetMode took too long (%v), expected ~5s timeout not full 7s block", elapsed)
	}

	// New instance should be stored regardless of timeout.
	if r.Instance() != newInstance {
		t.Fatal("new instance should be stored after deactivate timeout")
	}

	// New instance should have been activated.
	if newInstance.activateCount != 1 {
		t.Fatalf("new instance Activate called %d times, want 1", newInstance.activateCount)
	}

	// Mode should be updated.
	if r.CurrentMode() != "fast" {
		t.Fatalf("CurrentMode()=%q, want %q", r.CurrentMode(), "fast")
	}
}
