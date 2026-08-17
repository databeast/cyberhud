package textlayout_test

import (
	"fmt"
	"image"
	"testing"

	font "github.com/databeast/cyberhud/display/surface/fonts"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"pgregory.net/rapid"
)

type stubFace struct {
	id      string
	metrics font.Metrics
}

func (s stubFace) ID() string                       { return s.id }
func (s stubFace) Metrics() font.Metrics            { return s.metrics }
func (s stubFace) GlyphRow(ch rune, row int) uint32 { return 0 }

// For any bounds rectangle with Dx() > 0 and Dy() > 0 and a font registry containing
// at least one font that fits the width constraint, textlayout.DefaultTextHints(bounds)
// SHALL return TextHints where Catalog.PixelWidth() == bounds.Dx() and
// Catalog.MinChars() == 10.
func TestProperty_DefaultTextHintsAutoBuild(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		restore := font.SnapshotAndClear()
		defer restore()
		// Generate valid bounds ensuring pixelWidth >= 60 so maxAdvance >= 6
		// (MinChars=10, maxAdvance = pixelWidth/10)
		pixelWidth := rapid.IntRange(60, 1024).Draw(t, "pixelWidth")
		pixelHeight := rapid.IntRange(10, 1024).Draw(t, "pixelHeight")

		// Register at least one font that fits: advance <= pixelWidth/10
		maxAdvance := pixelWidth / 10
		numFonts := rapid.IntRange(1, 10).Draw(t, "numFonts")
		for i := 0; i < numFonts; i++ {
			advance := rapid.IntRange(1, maxAdvance).Draw(t, fmt.Sprintf("adv%d", i))
			height := rapid.IntRange(1, 64).Draw(t, fmt.Sprintf("h%d", i))
			glyphWidth := advance - 1
			if glyphWidth < 1 {
				glyphWidth = 1
			}
			font.Register(stubFace{
				id: fmt.Sprintf("testfont-%dx%d-%d", advance, height, i),
				metrics: font.Metrics{
					GlyphWidth:   glyphWidth,
					GlyphHeight:  height,
					GlyphAdvance: advance,
					RowHeight:    height + 2,
				},
			})
		}

		bounds := image.Rect(0, 0, pixelWidth, pixelHeight)
		h := textlayout.DefaultTextHints(bounds)

		if h.Catalog.PixelWidth() != pixelWidth {
			t.Fatalf("Catalog.PixelWidth() = %d, want %d", h.Catalog.PixelWidth(), pixelWidth)
		}
		if h.Catalog.PixelHeight() != pixelHeight {
			t.Fatalf("Catalog.PixelHeight() = %d, want %d", h.Catalog.PixelHeight(), pixelHeight)
		}
		if h.Catalog.MinChars() != 10 {
			t.Fatalf("Catalog.MinChars() = %d, want 10", h.Catalog.MinChars())
		}
	})
}
