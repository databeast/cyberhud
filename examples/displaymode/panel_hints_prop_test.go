package displaymode

import (
	"reflect"
	"testing"

	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/runtime/action"
	"pgregory.net/rapid"
)

// =============================================================================
//
// For any TextHints value, after calling SetPanelHints(hints), calling
// getPanelHints() shall return (hints, true) with the hints value equal to the
// stored input.

// =============================================================================

func TestPropertyPanelHintsStorageRoundTrip(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		hints := textlayout.TextHints{
			PixelWidth:         rapid.IntRange(1, 2048).Draw(t, "pixelWidth"),
			PixelHeight:        rapid.IntRange(1, 2048).Draw(t, "pixelHeight"),
			GlyphWidth:         rapid.IntRange(1, 20).Draw(t, "glyphWidth"),
			GlyphHeight:        rapid.IntRange(1, 20).Draw(t, "glyphHeight"),
			GlyphAdvance:       rapid.IntRange(1, 24).Draw(t, "glyphAdvance"),
			RowHeight:          rapid.IntRange(1, 32).Draw(t, "rowHeight"),
			PreferEventRefresh: rapid.Bool().Draw(t, "preferEventRefresh"),
		}

		SetPanelHints(hints)
		got, ok := getPanelHints()

		if !ok {
			t.Fatal("getPanelHints() returned set=false after SetPanelHints")
		}
		if !reflect.DeepEqual(got, hints) {
			t.Fatalf("round-trip failed: stored %+v, got %+v", hints, got)
		}
	})
}

// =============================================================================
//
// For any fixed Policy state, calling RenderCacheKey() multiple times without
// changing the Policy shall return the same string on every call.

// =============================================================================

func TestPropertyRenderCacheKeyDeterminism(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		names := registeredStyleNames()
		idx := rapid.IntRange(0, len(names)-1).Draw(t, "nameIdx")
		SetPolicy(Policy{Style: names[idx]})

		sig1 := RenderCacheKey()
		sig2 := RenderCacheKey()
		sig3 := RenderCacheKey()

		if sig1 != sig2 || sig2 != sig3 {
			t.Fatalf("RenderCacheKey not deterministic: %q, %q, %q", sig1, sig2, sig3)
		}
	})
}

// =============================================================================
//
// For any two distinct Policy.Style values, the strings returned by
// RenderCacheKey() shall be different.

// =============================================================================

func TestPropertyRenderCacheKeySensitivity(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		names := registeredStyleNames()
		if len(names) < 2 {
			t.Skip("need at least 2 styles")
		}
		idx1 := rapid.IntRange(0, len(names)-1).Draw(t, "idx1")
		idx2 := rapid.IntRange(0, len(names)-1).Filter(func(i int) bool {
			return i != idx1
		}).Draw(t, "idx2")

		SetPolicy(Policy{Style: names[idx1]})
		sig1 := RenderCacheKey()

		SetPolicy(Policy{Style: names[idx2]})
		sig2 := RenderCacheKey()

		if sig1 == sig2 {
			t.Fatalf("RenderCacheKey same for different styles %q and %q: %q",
				names[idx1], names[idx2], sig1)
		}
	})
}

// =============================================================================
//
// For any Policy state, len(RenderCacheKey()) shall be at most 512 bytes.

// =============================================================================

func TestPropertyRenderCacheKeyLengthInvariant(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		names := registeredStyleNames()
		idx := rapid.IntRange(0, len(names)-1).Draw(t, "nameIdx")
		SetPolicy(Policy{Style: names[idx]})

		sig := RenderCacheKey()
		if len(sig) > 512 {
			t.Fatalf("RenderCacheKey length %d exceeds 512 bytes", len(sig))
		}
		if len(sig) == 0 {
			t.Fatal("RenderCacheKey returned empty string")
		}
	})
}

// =============================================================================
//
// For any action.Action value not equal to action.Left or action.Right,
// Handler{}.HandleAction(act, cursor, itemCount) shall return a zero-value
// action.Result{}.

// =============================================================================

func TestPropertyUnhandledActionsReturnZeroResult(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate any action that is NOT Left or Right
		act := action.Action(rapid.Uint8().Filter(func(a uint8) bool {
			return action.Action(a) != action.Left && action.Action(a) != action.Right
		}).Draw(t, "action"))

		cursor := rapid.IntRange(0, 100).Draw(t, "cursor")
		itemCount := rapid.IntRange(0, 100).Draw(t, "itemCount")

		result := Handler{}.HandleAction(act, cursor, itemCount)

		zero := action.Result{}
		if result != zero {
			t.Fatalf("action %d returned non-zero result: %+v", act, result)
		}
	})
}
