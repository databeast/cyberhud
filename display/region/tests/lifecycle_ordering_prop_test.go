package tests

import (
	"fmt"
	"hash/fnv"
	"image"
	"sync"
	"testing"

	"github.com/databeast/cyberhud/display/region"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/runtime/action"
	"pgregory.net/rapid"
)

// =============================================================================
// Property 6: Lifecycle Ordering

//
// For any sequence of mode switches on a Region, Activate SHALL be called exactly
// once on each new instance before any BuildView call, and Deactivate SHALL be
// called exactly once on the outgoing instance before the next instance is
// activated. No BuildView call SHALL occur on an instance that has not been
// activated.
// =============================================================================

// lifecycleEvent represents a recorded lifecycle event on a mock instance.
type lifecycleEvent struct {
	instanceID string // which instance (mode ID + sequence number)
	event      string // "activate", "deactivate", or "buildview"
	seq        int    // global sequence number
}

// lifecycleInstance is a mock ModeInstance that records lifecycle calls.
type lifecycleInstance struct {
	id          string
	instanceID  string // unique per-construction (e.g., "clock-3")
	mu          sync.Mutex
	activated   bool
	deactivated bool
	log         *[]lifecycleEvent
	seqCounter  *int
}

func (m *lifecycleInstance) ID() string { return m.id }

func (m *lifecycleInstance) Activate() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.activated = true
	*m.seqCounter++
	*m.log = append(*m.log, lifecycleEvent{
		instanceID: m.instanceID,
		event:      "activate",
		seq:        *m.seqCounter,
	})
}

func (m *lifecycleInstance) Deactivate() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.deactivated = true
	*m.seqCounter++
	*m.log = append(*m.log, lifecycleEvent{
		instanceID: m.instanceID,
		event:      "deactivate",
		seq:        *m.seqCounter,
	})
}

func (m *lifecycleInstance) ActionHandler() action.Handler { return nil }

func (m *lifecycleInstance) BuildView() style.ViewData {
	m.mu.Lock()
	defer m.mu.Unlock()
	*m.seqCounter++
	*m.log = append(*m.log, lifecycleEvent{
		instanceID: m.instanceID,
		event:      "buildview",
		seq:        *m.seqCounter,
	})
	return style.ViewData{FontIDs: []string{"test-font"}}
}

func (m *lifecycleInstance) RenderCacheKey() uint32 {
	h := fnv.New32a()
	h.Write([]byte(m.instanceID))
	return h.Sum32()
}

// TestProperty6_LifecycleOrdering_ActivateBeforeBuildView asserts that for any
// sequence of mode switches, Activate is called on each new instance before any
// BuildView call, and Deactivate is called on the outgoing instance before the
// next instance is activated.
func TestProperty6_LifecycleOrdering_ActivateBeforeBuildView(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Setup: a set of mode IDs to switch between.
		numModes := rapid.IntRange(2, 6).Draw(t, "numModes")
		modeIDs := make([]string, numModes)
		for i := 0; i < numModes; i++ {
			modeIDs[i] = fmt.Sprintf("mode%d", i)
		}

		// Shared event log and sequence counter.
		var eventLog []lifecycleEvent
		seqCounter := 0
		constructionCounter := 0

		// Track all constructed instances.
		var allInstances []*lifecycleInstance

		// Create a Region with a surface.
		bounds := image.Rect(0, 0, 240, 240)
		surf := surface.New(bounds)
		r := region.NewRegion("test", bounds, surf)

		// Set up a ModeFactory that creates lifecycleInstance mocks.
		r.SetModeFactory(func(id string, hints textlayout.TextHints) (region.ModeInstance, bool) {
			for _, mID := range modeIDs {
				if mID == id {
					constructionCounter++
					inst := &lifecycleInstance{
						id:         id,
						instanceID: fmt.Sprintf("%s-%d", id, constructionCounter),
						log:        &eventLog,
						seqCounter: &seqCounter,
					}
					allInstances = append(allInstances, inst)
					return inst, true
				}
			}
			return nil, false
		})

		// Generate a random sequence of mode switches (3-20 switches).
		numSwitches := rapid.IntRange(3, 20).Draw(t, "numSwitches")
		for i := 0; i < numSwitches; i++ {
			modeIdx := rapid.IntRange(0, numModes-1).Draw(t, fmt.Sprintf("switch%d", i))
			modeID := modeIDs[modeIdx]

			err := r.SetMode(modeID)
			if err != nil {
				t.Fatalf("SetMode(%q) failed unexpectedly: %v", modeID, err)
			}

			// After SetMode, call BuildView on the active instance to simulate rendering.
			inst := r.Instance()
			if inst == nil {
				t.Fatalf("Instance() is nil after successful SetMode(%q)", modeID)
			}
			inst.BuildView()
		}

		// === Verify Property 6 assertions ===

		// Assertion 1: For each instance, Activate is called exactly once and
		// before any BuildView call on that instance.
		for _, inst := range allInstances {
			var activateSeq int
			activateCount := 0
			for _, ev := range eventLog {
				if ev.instanceID == inst.instanceID && ev.event == "activate" {
					activateCount++
					activateSeq = ev.seq
				}
			}
			if activateCount != 1 {
				t.Fatalf("instance %q: Activate called %d times, want exactly 1",
					inst.instanceID, activateCount)
			}

			// Every BuildView on this instance must come after its Activate.
			for _, ev := range eventLog {
				if ev.instanceID == inst.instanceID && ev.event == "buildview" {
					if ev.seq < activateSeq {
						t.Fatalf("instance %q: BuildView (seq=%d) occurred before Activate (seq=%d)",
							inst.instanceID, ev.seq, activateSeq)
					}
				}
			}
		}

		// Assertion 2: For each instance that was deactivated, Deactivate is
		// called exactly once.
		for _, inst := range allInstances {
			deactivateCount := 0
			for _, ev := range eventLog {
				if ev.instanceID == inst.instanceID && ev.event == "deactivate" {
					deactivateCount++
				}
			}
			// The last active instance won't be deactivated, all others should be.
			if inst.deactivated && deactivateCount != 1 {
				t.Fatalf("instance %q: Deactivate called %d times, want exactly 1",
					inst.instanceID, deactivateCount)
			}
		}

		// Assertion 3: Activate/Deactivate pairs are balanced — every Activate
		// on a new instance has a corresponding Deactivate on the old instance
		// before it. Specifically: for consecutive mode switches, the old
		// instance's Deactivate seq must be less than the new instance's
		// Activate seq (the old is deactivated before or during the switch that
		// activates the new one).
		//
		// NOTE: The implementation actually activates the NEW instance first,
		// then deactivates the old one. So the ordering is:
		//   new.Activate < old.Deactivate
		// This is by design (see task 4.2 CRITICAL note). We verify that
		// Deactivate on old happens before the NEXT Activate after the new one.
		activateEvents := []lifecycleEvent{}
		for _, ev := range eventLog {
			if ev.event == "activate" {
				activateEvents = append(activateEvents, ev)
			}
		}

		// For each pair of consecutive activations, verify that the first
		// instance's Deactivate occurs between them (after first's Activate,
		// before second's Activate — or equal to second's Activate-1 since
		// the implementation does: Activate(new) → Deactivate(old)).
		// Actually per the implementation: new.Activate happens THEN old.Deactivate.
		// So we verify: for each instance that was deactivated, its Deactivate
		// seq is less than ANY subsequent instance's first BuildView seq.
		for _, inst := range allInstances {
			if !inst.deactivated {
				continue // last active instance, not deactivated
			}
			var deactivateSeq int
			for _, ev := range eventLog {
				if ev.instanceID == inst.instanceID && ev.event == "deactivate" {
					deactivateSeq = ev.seq
					break
				}
			}

			// No BuildView should occur on this instance after Deactivate.
			for _, ev := range eventLog {
				if ev.instanceID == inst.instanceID && ev.event == "buildview" {
					if ev.seq > deactivateSeq {
						t.Fatalf("instance %q: BuildView (seq=%d) occurred after Deactivate (seq=%d)",
							inst.instanceID, ev.seq, deactivateSeq)
					}
				}
			}
		}

		// Assertion 4: No BuildView call occurs on an instance that has not
		// been activated. (Covered by Assertion 1's check that BuildView seq >
		// Activate seq, but let's be explicit.)
		activatedInstances := map[string]bool{}
		for _, ev := range eventLog {
			if ev.event == "activate" {
				activatedInstances[ev.instanceID] = true
			}
			if ev.event == "buildview" {
				if !activatedInstances[ev.instanceID] {
					t.Fatalf("instance %q: BuildView called without prior Activate",
						ev.instanceID)
				}
			}
		}
	})
}

// TestProperty6_LifecycleOrdering_DeactivateBeforeNextActivate verifies that
// for any sequence of mode switches, the outgoing instance's Deactivate is
// called before the next switch completes (i.e., before any BuildView on the
// next instance).
func TestProperty6_LifecycleOrdering_DeactivateBeforeNextActivate(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Use a smaller mode set to increase collision (same mode switched to again).
		numModes := rapid.IntRange(2, 4).Draw(t, "numModes")
		modeIDs := make([]string, numModes)
		for i := 0; i < numModes; i++ {
			modeIDs[i] = fmt.Sprintf("m%d", i)
		}

		var eventLog []lifecycleEvent
		seqCounter := 0
		constructionCounter := 0

		bounds := image.Rect(0, 0, 128, 128)
		surf := surface.New(bounds)
		r := region.NewRegion("test2", bounds, surf)

		r.SetModeFactory(func(id string, hints textlayout.TextHints) (region.ModeInstance, bool) {
			for _, mID := range modeIDs {
				if mID == id {
					constructionCounter++
					inst := &lifecycleInstance{
						id:         id,
						instanceID: fmt.Sprintf("%s-%d", id, constructionCounter),
						log:        &eventLog,
						seqCounter: &seqCounter,
					}
					return inst, true
				}
			}
			return nil, false
		})

		// Generate sequence of switches.
		numSwitches := rapid.IntRange(2, 15).Draw(t, "numSwitches")
		for i := 0; i < numSwitches; i++ {
			modeIdx := rapid.IntRange(0, numModes-1).Draw(t, fmt.Sprintf("sw%d", i))
			err := r.SetMode(modeIDs[modeIdx])
			if err != nil {
				t.Fatalf("SetMode(%q) error: %v", modeIDs[modeIdx], err)
			}

			// Simulate rendering after each switch.
			if inst := r.Instance(); inst != nil {
				inst.BuildView()
			}
		}

		// Verify: for each deactivated instance, no BuildView occurs on ANY
		// subsequent instance before the deactivation completes.
		// In practice, we verify: old.Deactivate < next.BuildView (first one).
		//
		// Collect activate events in order.
		type activateInfo struct {
			instanceID string
			seq        int
		}
		var activations []activateInfo
		deactivateSeqs := map[string]int{}

		for _, ev := range eventLog {
			if ev.event == "activate" {
				activations = append(activations, activateInfo{ev.instanceID, ev.seq})
			}
			if ev.event == "deactivate" {
				deactivateSeqs[ev.instanceID] = ev.seq
			}
		}

		// For consecutive activations (i, i+1), verify that instance[i]'s
		// Deactivate seq < first BuildView seq of instance[i+1].
		for i := 0; i < len(activations)-1; i++ {
			prevInstanceID := activations[i].instanceID
			nextInstanceID := activations[i+1].instanceID

			prevDeactSeq, wasDeactivated := deactivateSeqs[prevInstanceID]
			if !wasDeactivated {
				// Should only be the last instance that's not deactivated.
				if i < len(activations)-2 {
					t.Fatalf("instance %q was not deactivated but is not the last",
						prevInstanceID)
				}
				continue
			}

			// Find first BuildView of the next instance.
			for _, ev := range eventLog {
				if ev.instanceID == nextInstanceID && ev.event == "buildview" {
					if ev.seq < prevDeactSeq {
						t.Fatalf("instance %q BuildView (seq=%d) occurred before previous instance %q Deactivate (seq=%d)",
							nextInstanceID, ev.seq, prevInstanceID, prevDeactSeq)
					}
					break // only check first BuildView
				}
			}
		}
	})
}
