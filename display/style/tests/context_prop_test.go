package tests_test

import (
	"reflect"
	"testing"

	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/style/layout"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"pgregory.net/rapid"
)

// For any random TextHints, constructing a StyleContext via NewStyleContext and
// then reading back via accessors returns values identical to the inputs.

func genContextTextHints(t *rapid.T) textlayout.TextHints {
	return textlayout.TextHints{
		PixelWidth:               rapid.IntRange(0, 1000).Draw(t, "pixelWidth"),
		PixelHeight:              rapid.IntRange(0, 1000).Draw(t, "pixelHeight"),
		GlyphWidth:               rapid.IntRange(1, 32).Draw(t, "glyphWidth"),
		GlyphHeight:              rapid.IntRange(1, 32).Draw(t, "glyphHeight"),
		GlyphAdvance:             rapid.IntRange(1, 32).Draw(t, "glyphAdvance"),
		RowHeight:                rapid.IntRange(1, 64).Draw(t, "rowHeight"),
		PreferEventRefresh:       rapid.Bool().Draw(t, "preferEventRefresh"),
		Capability:               rapid.IntRange(0, 5).Draw(t, "capability"),
		SupportsVerticalScroll:   rapid.Bool().Draw(t, "supportsVerticalScroll"),
		SupportsHorizontalScroll: rapid.Bool().Draw(t, "supportsHorizontalScroll"),
	}
}

func TestProperty1_ConstructionRoundTripPreservation(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		hints := genContextTextHints(t)
		paddingPct := rapid.IntRange(0, 50).Draw(t, "paddingPct")
		ctx := style.NewStyleContext(hints)

		// Layout round-trip: a bridge built via the context matches one built
		// directly from the same hints and padding.
		want := layout.NewLayoutBridge(hints, layout.BridgeConfig{PaddingPct: paddingPct})
		if ctx.Layout(paddingPct) != want {
			t.Fatalf("Layout(%d) mismatch: got %+v, want %+v", paddingPct, ctx.Layout(paddingPct), want)
		}

		// FontCatalog round-trip
		if !reflect.DeepEqual(ctx.FontCatalog(), hints.Catalog) {
			t.Fatalf("FontCatalog() mismatch: got %+v, want %+v", ctx.FontCatalog(), hints.Catalog)
		}

		// Capability round-trip
		if ctx.Cap() != style.Capability(hints.Capability) {
			t.Fatalf("Cap() = %v, want %v", ctx.Cap(), style.Capability(hints.Capability))
		}

		// VerticalScroll round-trip
		if ctx.VerticalScroll() != hints.SupportsVerticalScroll {
			t.Fatalf("VerticalScroll() = %v, want %v", ctx.VerticalScroll(), hints.SupportsVerticalScroll)
		}

		// HorizontalScroll round-trip
		if ctx.HorizontalScroll() != hints.SupportsHorizontalScroll {
			t.Fatalf("HorizontalScroll() = %v, want %v", ctx.HorizontalScroll(), hints.SupportsHorizontalScroll)
		}

		// Hints round-trip (compare individual comparable fields)
		gotHints := ctx.Hints()
		if gotHints.PixelWidth != hints.PixelWidth {
			t.Fatalf("Hints().PixelWidth mismatch: got %d, want %d", gotHints.PixelWidth, hints.PixelWidth)
		}
		if gotHints.PixelHeight != hints.PixelHeight {
			t.Fatalf("Hints().PixelHeight mismatch: got %d, want %d", gotHints.PixelHeight, hints.PixelHeight)
		}
		if gotHints.GlyphAdvance != hints.GlyphAdvance {
			t.Fatalf("Hints().GlyphAdvance mismatch: got %d, want %d", gotHints.GlyphAdvance, hints.GlyphAdvance)
		}
		if gotHints.RowHeight != hints.RowHeight {
			t.Fatalf("Hints().RowHeight mismatch: got %d, want %d", gotHints.RowHeight, hints.RowHeight)
		}
		if gotHints.Capability != hints.Capability {
			t.Fatalf("Hints().Capability mismatch: got %d, want %d", gotHints.Capability, hints.Capability)
		}
	})
}

// For any StyleContext constructed from arbitrary TextHints (with non-trivial
// PixelWidth/PixelHeight values in range 100-2000), there shall be no exported
// method on StyleContext whose single return value is a bare int that equals
// the original PixelWidth or PixelHeight.
//
// Note: Styles access raw pixel dimensions via ctx.Hints() (which returns a struct,
// not a bare int). This is architecturally correct — styles need hints to construct
// their own LayoutBridge and for panel-covering elements. This test validates that
// no *convenience accessor* leaks raw dimensions as bare int return values.

func TestProperty2_NoPixelDimensionLeakage(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate hints with non-trivial pixel dimensions (avoid small values
		// that could coincidentally match boolean/small-integer return values).
		hints := textlayout.TextHints{
			PixelWidth:               rapid.IntRange(100, 2000).Draw(t, "pixelWidth"),
			PixelHeight:              rapid.IntRange(100, 2000).Draw(t, "pixelHeight"),
			GlyphWidth:               rapid.IntRange(1, 32).Draw(t, "glyphWidth"),
			GlyphHeight:              rapid.IntRange(1, 32).Draw(t, "glyphHeight"),
			GlyphAdvance:             rapid.IntRange(1, 32).Draw(t, "glyphAdvance"),
			RowHeight:                rapid.IntRange(1, 64).Draw(t, "rowHeight"),
			Capability:               rapid.IntRange(0, 5).Draw(t, "capability"),
			SupportsVerticalScroll:   rapid.Bool().Draw(t, "supportsVerticalScroll"),
			SupportsHorizontalScroll: rapid.Bool().Draw(t, "supportsHorizontalScroll"),
		}

		ctx := style.NewStyleContext(hints)

		// Use reflection to enumerate all exported methods on StyleContext.
		ctxVal := reflect.ValueOf(ctx)
		ctxType := ctxVal.Type()

		for i := 0; i < ctxType.NumMethod(); i++ {
			method := ctxType.Method(i)
			// Only check methods with no arguments (beyond receiver) that return a single value.
			if method.Type.NumIn() != 1 || method.Type.NumOut() != 1 {
				continue
			}
			result := ctxVal.Method(i).Call(nil)
			// Check if the result is an int that matches PixelWidth or PixelHeight.
			if result[0].Kind() == reflect.Int {
				val := int(result[0].Int())
				if val == hints.PixelWidth {
					t.Fatalf("method %s returns PixelWidth value %d", method.Name, val)
				}
				if val == hints.PixelHeight {
					t.Fatalf("method %s returns PixelHeight value %d", method.Name, val)
				}
			}
		}
	})
}
