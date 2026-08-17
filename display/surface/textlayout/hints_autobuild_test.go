package textlayout_test

import (
	"image"
	"testing"

	"github.com/databeast/cyberhud/display/surface/fonts"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

func TestDefaultTextHints_ValidBounds_PopulatesCatalog(t *testing.T) {
	restore := font.SnapshotAndClear()
	defer restore()

	// Register fonts that fit a 128px wide region (advance <= 12, since MinChars=10).
	font.Register(stubFace{
		id:      "spleen-5x8",
		metrics: font.Metrics{GlyphWidth: 5, GlyphHeight: 8, GlyphAdvance: 6, RowHeight: 10},
	})
	font.Register(stubFace{
		id:      "spleen-8x16",
		metrics: font.Metrics{GlyphWidth: 8, GlyphHeight: 16, GlyphAdvance: 9, RowHeight: 18},
	})

	bounds := image.Rect(0, 0, 128, 64)
	h := textlayout.DefaultTextHints(bounds)

	if h.Catalog.PixelWidth() != 128 {
		t.Errorf("Catalog.PixelWidth() = %d, want 128", h.Catalog.PixelWidth())
	}
	if h.Catalog.PixelHeight() != 64 {
		t.Errorf("Catalog.PixelHeight() = %d, want 64", h.Catalog.PixelHeight())
	}
	if h.Catalog.MinChars() != 10 {
		t.Errorf("Catalog.MinChars() = %d, want 10", h.Catalog.MinChars())
	}
}

func TestDefaultTextHints_ZeroBounds_ZeroCatalog(t *testing.T) {
	restore := font.SnapshotAndClear()
	defer restore()

	font.Register(stubFace{
		id:      "spleen-5x8",
		metrics: font.Metrics{GlyphWidth: 5, GlyphHeight: 8, GlyphAdvance: 6, RowHeight: 10},
	})

	bounds := image.Rect(0, 0, 0, 0)
	h := textlayout.DefaultTextHints(bounds)

	if h.Catalog.PixelWidth() != 0 {
		t.Errorf("Catalog.PixelWidth() = %d, want 0 for zero bounds", h.Catalog.PixelWidth())
	}
	if h.Catalog.PixelHeight() != 0 {
		t.Errorf("Catalog.PixelHeight() = %d, want 0 for zero bounds", h.Catalog.PixelHeight())
	}
}

func TestDefaultTextHints_NoFontFits_ZeroCatalog(t *testing.T) {
	restore := font.SnapshotAndClear()
	defer restore()

	// Register only a font with advance 20. For a 128px region with MinChars=10,
	// maxAdvance = 128/10 = 12, so advance 20 won't fit.
	font.Register(stubFace{
		id:      "large-20x30",
		metrics: font.Metrics{GlyphWidth: 20, GlyphHeight: 30, GlyphAdvance: 20, RowHeight: 32},
	})

	bounds := image.Rect(0, 0, 128, 64)
	h := textlayout.DefaultTextHints(bounds)

	if h.Catalog.PixelWidth() != 0 {
		t.Errorf("Catalog.PixelWidth() = %d, want 0 when no font fits", h.Catalog.PixelWidth())
	}
	if h.Catalog.PixelHeight() != 0 {
		t.Errorf("Catalog.PixelHeight() = %d, want 0 when no font fits", h.Catalog.PixelHeight())
	}
}
