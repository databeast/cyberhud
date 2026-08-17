package textlayout_test

import (
	"reflect"
	"testing"

	"github.com/databeast/cyberhud/display/surface/fonts"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"pgregory.net/rapid"
)

// --- From: hints_test.go ---

func TestMaxCharsPerRow(t *testing.T) {
	h := textlayout.TextHints{PixelWidth: 240, GlyphAdvance: 6}
	if got := textlayout.MaxCharsPerRow(h, 0); got != 40 {
		t.Fatalf("MaxCharsPerRow()=%d, want 40", got)
	}
	if got := textlayout.MaxCharsPerRow(h, 5); got != 38 {
		t.Fatalf("MaxCharsPerRow() with padding=%d, want 38", got)
	}
}

func TestMaxVisibleRows(t *testing.T) {
	h := textlayout.TextHints{PixelHeight: 240, RowHeight: 10}
	if got := textlayout.MaxVisibleRows(h, 0); got != 24 {
		t.Fatalf("MaxVisibleRows()=%d, want 24", got)
	}
	if got := textlayout.MaxVisibleRows(h, 2); got != 23 {
		t.Fatalf("MaxVisibleRows() with padding=%d, want 23", got)
	}
}

func TestHelpersGuardInvalidValues(t *testing.T) {
	if got := textlayout.MaxCharsPerRow(textlayout.TextHints{PixelWidth: 100, GlyphAdvance: 0}, 0); got != 0 {
		t.Fatalf("MaxCharsPerRow() with invalid advance=%d, want 0", got)
	}
	if got := textlayout.MaxVisibleRows(textlayout.TextHints{PixelHeight: 100, RowHeight: 0}, 0); got != 0 {
		t.Fatalf("MaxVisibleRows() with invalid row height=%d, want 0", got)
	}
}

// --- From: textlayout_prop_test.go ---

// mockFace implements font.Face for property-based testing.
type mockFace struct {
	id      string
	metrics font.Metrics
}

func (f mockFace) ID() string                       { return f.id }
func (f mockFace) Metrics() font.Metrics            { return f.metrics }
func (f mockFace) GlyphRow(ch rune, row int) uint32 { return 0 }

// genTextHints generates an arbitrary TextHints with random field values.
func genTextHints(t *rapid.T) textlayout.TextHints {
	tickerDirs := []string{textlayout.TickerDirectionVertical, textlayout.TickerDirectionNone, "horizontal"}
	lineModes := []string{textlayout.LineModeTruncate, textlayout.LineModeClip, "wrap"}

	return textlayout.TextHints{
		PixelWidth:               rapid.IntRange(1, 1920).Draw(t, "pixelWidth"),
		PixelHeight:              rapid.IntRange(1, 1080).Draw(t, "pixelHeight"),
		GlyphWidth:               rapid.IntRange(1, 32).Draw(t, "glyphWidth"),
		GlyphHeight:              rapid.IntRange(1, 32).Draw(t, "glyphHeight"),
		GlyphAdvance:             rapid.IntRange(1, 32).Draw(t, "glyphAdvance"),
		RowHeight:                rapid.IntRange(1, 64).Draw(t, "rowHeight"),
		SupportsVerticalScroll:   rapid.Bool().Draw(t, "supportsVerticalScroll"),
		SupportsHorizontalScroll: rapid.Bool().Draw(t, "supportsHorizontalScroll"),
		SupportsAutoScroll:       rapid.Bool().Draw(t, "supportsAutoScroll"),
		PreferEventRefresh:       rapid.Bool().Draw(t, "preferEventRefresh"),
		DefaultTickerDirection:   rapid.SampledFrom(tickerDirs).Draw(t, "defaultTickerDirection"),
		DefaultLineMode:          rapid.SampledFrom(lineModes).Draw(t, "defaultLineMode"),
	}
}

// genMockFace generates an arbitrary font.Face with random metrics.
func genMockFace(t *rapid.T) font.Face {
	return mockFace{
		id: rapid.StringMatching(`[a-z]{3,12}-[0-9]{1,2}px`).Draw(t, "fontID"),
		metrics: font.Metrics{
			GlyphWidth:   rapid.IntRange(1, 32).Draw(t, "fontGlyphWidth"),
			GlyphHeight:  rapid.IntRange(1, 32).Draw(t, "fontGlyphHeight"),
			GlyphAdvance: rapid.IntRange(1, 32).Draw(t, "fontGlyphAdvance"),
			RowHeight:    rapid.IntRange(1, 64).Draw(t, "fontRowHeight"),
		},
	}
}

// For any TextHints h and font.Face f where the font's metrics differ from the hints'
// glyph metrics (isBugCondition is true), applying textlayout.WithFont(h, f) SHALL return
// a new TextHints where RowHeight, GlyphWidth, GlyphHeight, and GlyphAdvance match the
// font's Metrics().

func TestProperty1_WithFont_LayoutMetricsMatchSelectedFont(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		h := genTextHints(t)
		f := genMockFace(t)

		// Ensure bug condition: at least one metric differs between hints and font.
		fm := f.Metrics()
		if h.GlyphWidth == fm.GlyphWidth && h.GlyphHeight == fm.GlyphHeight &&
			h.GlyphAdvance == fm.GlyphAdvance && h.RowHeight == fm.RowHeight {
			// Force a difference so we always test the bug condition path.
			h.RowHeight = fm.RowHeight + 1
		}

		result := textlayout.WithFont(h, f)

		if result.RowHeight != fm.RowHeight {
			t.Fatalf("RowHeight mismatch: got %d, want %d", result.RowHeight, fm.RowHeight)
		}
		if result.GlyphWidth != fm.GlyphWidth {
			t.Fatalf("GlyphWidth mismatch: got %d, want %d", result.GlyphWidth, fm.GlyphWidth)
		}
		if result.GlyphHeight != fm.GlyphHeight {
			t.Fatalf("GlyphHeight mismatch: got %d, want %d", result.GlyphHeight, fm.GlyphHeight)
		}
		if result.GlyphAdvance != fm.GlyphAdvance {
			t.Fatalf("GlyphAdvance mismatch: got %d, want %d", result.GlyphAdvance, fm.GlyphAdvance)
		}
	})
}

// For any TextHints h and font.Face f, applying WithFont(h, f) SHALL preserve
// h.PixelWidth, h.PixelHeight, h.SupportsVerticalScroll, h.SupportsHorizontalScroll,
// h.SupportsAutoScroll, h.PreferEventRefresh, h.DefaultTickerDirection, and
// h.DefaultLineMode unchanged.

func TestProperty2_WithFont_PreservesNonMetricFields(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		h := genTextHints(t)
		f := genMockFace(t)

		result := textlayout.WithFont(h, f)

		if result.PixelWidth != h.PixelWidth {
			t.Fatalf("PixelWidth changed: got %d, want %d", result.PixelWidth, h.PixelWidth)
		}
		if result.PixelHeight != h.PixelHeight {
			t.Fatalf("PixelHeight changed: got %d, want %d", result.PixelHeight, h.PixelHeight)
		}
		if result.SupportsVerticalScroll != h.SupportsVerticalScroll {
			t.Fatalf("SupportsVerticalScroll changed: got %v, want %v", result.SupportsVerticalScroll, h.SupportsVerticalScroll)
		}
		if result.SupportsHorizontalScroll != h.SupportsHorizontalScroll {
			t.Fatalf("SupportsHorizontalScroll changed: got %v, want %v", result.SupportsHorizontalScroll, h.SupportsHorizontalScroll)
		}
		if result.SupportsAutoScroll != h.SupportsAutoScroll {
			t.Fatalf("SupportsAutoScroll changed: got %v, want %v", result.SupportsAutoScroll, h.SupportsAutoScroll)
		}
		if result.PreferEventRefresh != h.PreferEventRefresh {
			t.Fatalf("PreferEventRefresh changed: got %v, want %v", result.PreferEventRefresh, h.PreferEventRefresh)
		}
		if result.DefaultTickerDirection != h.DefaultTickerDirection {
			t.Fatalf("DefaultTickerDirection changed: got %q, want %q", result.DefaultTickerDirection, h.DefaultTickerDirection)
		}
		if result.DefaultLineMode != h.DefaultLineMode {
			t.Fatalf("DefaultLineMode changed: got %q, want %q", result.DefaultLineMode, h.DefaultLineMode)
		}
	})
}

// For any TextHints h, applying WithFont(h, nil) SHALL return h completely unchanged
// (identity operation / nil-safety).

func TestProperty2_WithFont_NilFaceReturnsUnchanged(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		h := genTextHints(t)

		result := textlayout.WithFont(h, nil)

		if !reflect.DeepEqual(result, h) {
			t.Fatalf("WithFont(h, nil) != h\ngot:  %+v\nwant: %+v", result, h)
		}
	})
}
