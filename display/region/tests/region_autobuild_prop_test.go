package tests_test

import (
	"fmt"
	"image"
	"testing"

	"github.com/databeast/cyberhud/display/region"
	"github.com/databeast/cyberhud/display/surface"
	font "github.com/databeast/cyberhud/display/surface/fonts"
	"pgregory.net/rapid"
)

// stubFace is a minimal font.Face implementation for property testing.
type stubFace struct {
	id      string
	metrics font.Metrics
}

func (s stubFace) ID() string                { return s.id }
func (s stubFace) Metrics() font.Metrics     { return s.metrics }
func (s stubFace) GlyphRow(rune, int) uint32 { return 0 }

// For any Region with valid bounds (Dx > 0, Dy > 0) and a font registry containing
// at least one font that fits the width constraint, after resolveTextHints() completes,
// the resulting TextHints.Catalog.PixelWidth() SHALL equal TextHints.PixelWidth and
// TextHints.Catalog.PixelHeight() SHALL equal TextHints.PixelHeight.
func TestProperty_AutoBuildMatchesDimensions(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		restore := font.SnapshotAndClear()
		defer restore()
		// Generate random valid bounds large enough for fonts to fit.
		// MinChars=10, so maxAdvance = pixelWidth/10. We need at least one font
		// with advance <= maxAdvance.
		pixelWidth := rapid.IntRange(60, 1024).Draw(t, "pixelWidth")
		pixelHeight := rapid.IntRange(10, 1024).Draw(t, "pixelHeight")

		// Compute the maximum advance that will satisfy the width constraint.
		maxAdvance := pixelWidth / 10

		// Register at least one font that fits the width constraint.
		numFonts := rapid.IntRange(1, 10).Draw(t, "numFonts")
		registeredFitting := false
		for i := 0; i < numFonts; i++ {
			glyphWidth := rapid.IntRange(1, 32).Draw(t, fmt.Sprintf("gw%d", i))
			glyphHeight := rapid.IntRange(4, 64).Draw(t, fmt.Sprintf("gh%d", i))
			// Ensure at least one font has a fitting advance.
			var glyphAdvance int
			if !registeredFitting || rapid.Bool().Draw(t, fmt.Sprintf("fitAdv%d", i)) {
				// Generate an advance that fits.
				glyphAdvance = rapid.IntRange(1, maxAdvance).Draw(t, fmt.Sprintf("ga%d", i))
				registeredFitting = true
			} else {
				// Generate an arbitrary advance (may or may not fit).
				glyphAdvance = rapid.IntRange(1, 64).Draw(t, fmt.Sprintf("ga%d", i))
			}
			rowHeight := glyphHeight + rapid.IntRange(0, 4).Draw(t, fmt.Sprintf("rh%d", i))

			family := rapid.SampledFrom([]string{"spleen", "terminus", "cozette", "testfont"}).Draw(t, fmt.Sprintf("fam%d", i))
			id := fmt.Sprintf("%s-%dx%d-autobuild-%d", family, glyphWidth, glyphHeight, i)

			font.Register(stubFace{
				id: id,
				metrics: font.Metrics{
					GlyphWidth:   glyphWidth,
					GlyphHeight:  glyphHeight,
					GlyphAdvance: glyphAdvance,
					RowHeight:    rowHeight,
				},
			})
		}

		// Create bounds and surface.
		bounds := image.Rect(0, 0, pixelWidth, pixelHeight)
		surf := surface.New(bounds)

		// Branch 1: No screens — verify auto-build via DefaultTextHints path.
		r := region.NewRegionWithScreens("prop-test", bounds, surf, nil, "", 0, 0)

		hints := r.TextHints()

		// Verify catalog dimensions match the resolved text hints dimensions.
		if hints.Catalog.PixelWidth() != hints.PixelWidth {
			t.Fatalf("Catalog.PixelWidth()=%d, want TextHints.PixelWidth=%d",
				hints.Catalog.PixelWidth(), hints.PixelWidth)
		}
		if hints.Catalog.PixelHeight() != hints.PixelHeight {
			t.Fatalf("Catalog.PixelHeight()=%d, want TextHints.PixelHeight=%d",
				hints.Catalog.PixelHeight(), hints.PixelHeight)
		}
	})
}
