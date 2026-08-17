package modehints_test

import (
	"reflect"
	"testing"

	"github.com/databeast/cyberhud/display/region/modehints"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"pgregory.net/rapid"
)

// --- From: idempotence_prop_test.go ---

// Property 3: PropagateHints is idempotent
//
// For any valid TextHints value, calling PropagateHints once or multiple times
// with the same value produces the same stored state in each registered receiver.

// mockIdempotenceReceiver records the hints value received via SetPanelHints.
type mockIdempotenceReceiver struct {
	hints textlayout.TextHints
	calls int
}

func (m *mockIdempotenceReceiver) SetPanelHints(hints textlayout.TextHints) {
	m.hints = hints
	m.calls++
}

// TestProperty3_PropagateHintsIdempotence verifies that calling PropagateHints
// multiple times with the same TextHints value results in an identical stored
// state in each receiver after each call.

func TestProperty3_PropagateHintsIdempotence(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		// 1. Reset receivers for test isolation.
		modehints.ResetHintsReceivers()

		// 2. Generate N receivers and register them.
		n := rapid.IntRange(1, 10).Draw(rt, "numReceivers")
		receivers := make([]*mockIdempotenceReceiver, n)
		for i := range receivers {
			receivers[i] = &mockIdempotenceReceiver{}
			modehints.RegisterHintsReceiver(receivers[i])
		}

		// 3. Generate arbitrary TextHints.
		hints := textlayout.TextHints{
			PixelWidth:               rapid.IntRange(1, 1920).Draw(rt, "pixelWidth"),
			PixelHeight:              rapid.IntRange(1, 1080).Draw(rt, "pixelHeight"),
			GlyphWidth:               rapid.IntRange(1, 32).Draw(rt, "glyphWidth"),
			GlyphHeight:              rapid.IntRange(1, 48).Draw(rt, "glyphHeight"),
			GlyphAdvance:             rapid.IntRange(1, 32).Draw(rt, "glyphAdvance"),
			RowHeight:                rapid.IntRange(1, 64).Draw(rt, "rowHeight"),
			SupportsVerticalScroll:   rapid.Bool().Draw(rt, "supportsVerticalScroll"),
			SupportsHorizontalScroll: rapid.Bool().Draw(rt, "supportsHorizontalScroll"),
			SupportsAutoScroll:       rapid.Bool().Draw(rt, "supportsAutoScroll"),
			PreferEventRefresh:       rapid.Bool().Draw(rt, "preferEventRefresh"),
			DefaultTickerDirection:   rapid.SampledFrom([]string{"vertical", "none"}).Draw(rt, "tickerDirection"),
			DefaultLineMode:          rapid.SampledFrom([]string{"truncate", "clip"}).Draw(rt, "lineMode"),
		}

		// 4. Call PropagateHints once — record state.
		modehints.PropagateHints(hints)

		stateAfterFirst := make([]textlayout.TextHints, n)
		for i, r := range receivers {
			stateAfterFirst[i] = r.hints
		}

		// 5. Call PropagateHints again with the same value — record state.
		modehints.PropagateHints(hints)

		stateAfterSecond := make([]textlayout.TextHints, n)
		for i, r := range receivers {
			stateAfterSecond[i] = r.hints
		}

		// 6. Assert states are identical after both calls.
		for i := range receivers {
			if !reflect.DeepEqual(stateAfterFirst[i], stateAfterSecond[i]) {
				t.Fatalf("receiver %d: state after first call %+v != state after second call %+v",
					i, stateAfterFirst[i], stateAfterSecond[i])
			}
			// Also verify the stored value matches the input hints.
			if !reflect.DeepEqual(stateAfterSecond[i], hints) {
				t.Fatalf("receiver %d: stored hints %+v != input hints %+v",
					i, stateAfterSecond[i], hints)
			}
		}

		// 7. Call PropagateHints a third time and assert same state.
		modehints.PropagateHints(hints)

		for i, r := range receivers {
			if !reflect.DeepEqual(r.hints, stateAfterSecond[i]) {
				t.Fatalf("receiver %d: state after third call %+v != state after second call %+v",
					i, r.hints, stateAfterSecond[i])
			}
		}
	})
}

// --- From: propagate_prop_test.go ---

// Property 2: PropagateHints distributes uniformly to all receivers.
// For any set of registered PanelHintsReceivers and any valid TextHints value,
// after PropagateHints is called, every receiver holds the exact same hints
// value that was passed in.

// mockReceiver captures the TextHints value delivered by PropagateHints.
type mockReceiver struct {
	received textlayout.TextHints
	called   bool
}

func (m *mockReceiver) SetPanelHints(hints textlayout.TextHints) {
	m.received = hints
	m.called = true
}

// genTextHints generates arbitrary TextHints values for property testing.
func genTextHints(t *rapid.T) textlayout.TextHints {
	return textlayout.TextHints{
		PixelWidth:               rapid.IntRange(1, 2000).Draw(t, "pixelWidth"),
		PixelHeight:              rapid.IntRange(1, 2000).Draw(t, "pixelHeight"),
		GlyphWidth:               rapid.IntRange(1, 32).Draw(t, "glyphWidth"),
		GlyphHeight:              rapid.IntRange(1, 32).Draw(t, "glyphHeight"),
		GlyphAdvance:             rapid.IntRange(1, 32).Draw(t, "glyphAdvance"),
		RowHeight:                rapid.IntRange(1, 64).Draw(t, "rowHeight"),
		SupportsVerticalScroll:   rapid.Bool().Draw(t, "supportsVerticalScroll"),
		SupportsHorizontalScroll: rapid.Bool().Draw(t, "supportsHorizontalScroll"),
		SupportsAutoScroll:       rapid.Bool().Draw(t, "supportsAutoScroll"),
		PreferEventRefresh:       rapid.Bool().Draw(t, "preferEventRefresh"),
		DefaultTickerDirection:   rapid.SampledFrom([]string{"vertical", "none"}).Draw(t, "tickerDirection"),
		DefaultLineMode:          rapid.SampledFrom([]string{"truncate", "clip"}).Draw(t, "lineMode"),
	}
}

// TestProperty_PropagateHints_UniformDistribution verifies that for any
// number of registered receivers (1-20) and any valid TextHints value,
// PropagateHints delivers the exact same hints to every receiver.
func TestProperty_PropagateHints_UniformDistribution(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		// Isolate test state
		modehints.ResetHintsReceivers()

		// Generate N receivers (1-20)
		n := rapid.IntRange(1, 20).Draw(rt, "receiverCount")
		receivers := make([]*mockReceiver, n)
		for i := range receivers {
			receivers[i] = &mockReceiver{}
			modehints.RegisterHintsReceiver(receivers[i])
		}

		// Generate arbitrary TextHints
		hints := genTextHints(rt)

		// Propagate
		modehints.PropagateHints(hints)

		// Assert all receivers got the exact same value
		for i, r := range receivers {
			if !r.called {
				t.Fatalf("receiver %d was not called", i)
			}
			if !reflect.DeepEqual(r.received, hints) {
				t.Fatalf("receiver %d received %+v, want %+v", i, r.received, hints)
			}
		}
	})
}
