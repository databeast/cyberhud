package tests

import (
	"image"
	"image/color"
	"testing"

	fontpkg "github.com/databeast/cyberhud/display/surface/fonts"
	"github.com/databeast/cyberhud/display/widgets/textlabel"
	"pgregory.net/rapid"
)

// ============================================================================
//
// These tests verify that widget rendering produces identical output for
// same inputs and are independent of the old centralized font resolution path.
// ============================================================================

// --- Helpers ---

// preservationTestFace is a minimal font.Face for preservation property testing.
type preservationTestFace struct {
	id      string
	metrics fontpkg.Metrics
}

func (f preservationTestFace) ID() string                { return f.id }
func (f preservationTestFace) Metrics() fontpkg.Metrics  { return f.metrics }
func (f preservationTestFace) GlyphRow(rune, int) uint32 { return 0 }

// registerPreservationModeFonts registers a controlled set of fonts needed for
// mode-level preservation testing.
func registerPreservationModeFonts() func() {
	restore := fontpkg.SnapshotAndClear()
	fonts := []preservationTestFace{
		{id: "spleen-5x8", metrics: fontpkg.Metrics{GlyphWidth: 5, GlyphHeight: 8, GlyphAdvance: 6, RowHeight: 10}},
		{id: "spleen-6x12", metrics: fontpkg.Metrics{GlyphWidth: 6, GlyphHeight: 12, GlyphAdvance: 7, RowHeight: 14}},
		{id: "spleen-8x16", metrics: fontpkg.Metrics{GlyphWidth: 8, GlyphHeight: 16, GlyphAdvance: 9, RowHeight: 18}},
		{id: "terminus-6x12", metrics: fontpkg.Metrics{GlyphWidth: 6, GlyphHeight: 12, GlyphAdvance: 7, RowHeight: 14}},
		{id: "terminus-8x14", metrics: fontpkg.Metrics{GlyphWidth: 8, GlyphHeight: 14, GlyphAdvance: 9, RowHeight: 16}},
		{id: "cozette-6x13", metrics: fontpkg.Metrics{GlyphWidth: 6, GlyphHeight: 13, GlyphAdvance: 7, RowHeight: 15}},
	}
	for _, f := range fonts {
		fontpkg.Register(f)
	}
	return restore
}

// --- Property Tests ---

// TestPreservation_WidgetRenderingDeterminism verifies that textlabel.Render()
// produces identical bitmap output given the same face and config inputs.
// This confirms widget rendering mechanics are independent of font selection logic.

func TestPreservation_WidgetRenderingDeterminism(t *testing.T) {
	restore := registerPreservationModeFonts()
	defer restore()

	rapid.Check(t, func(rt *rapid.T) {
		// Generate a config with consistent inputs.
		text := rapid.StringMatching(`[A-Za-z0-9 ]{1,30}`).Draw(rt, "text")
		x := rapid.IntRange(0, 100).Draw(rt, "x")
		y := rapid.IntRange(0, 100).Draw(rt, "y")
		w := rapid.IntRange(10, 200).Draw(rt, "w")
		h := rapid.IntRange(8, 50).Draw(rt, "h")
		alignment := rapid.SampledFrom([]textlabel.Alignment{
			textlabel.Left, textlabel.Center, textlabel.Right,
		}).Draw(rt, "alignment")

		r := uint8(rapid.IntRange(0, 255).Draw(rt, "r"))
		g := uint8(rapid.IntRange(0, 255).Draw(rt, "g"))
		b := uint8(rapid.IntRange(0, 255).Draw(rt, "b"))

		// Use one of the registered fonts.
		fontID := rapid.SampledFrom([]string{
			"spleen-5x8", "spleen-6x12", "spleen-8x16",
			"terminus-6x12", "terminus-8x14", "cozette-6x13",
		}).Draw(rt, "fontID")
		face, ok := fontpkg.Get(fontID)
		if !ok {
			t.Fatalf("font %q not registered", fontID)
		}

		cfg := textlabel.Config{
			Text:       text,
			Bounds:     image.Rect(x, y, x+w, y+h),
			Font:       face,
			Alignment:  alignment,
			Foreground: color.RGBA{R: r, G: g, B: b, A: 255},
		}

		// Render twice with identical inputs.
		sprite1 := textlabel.Render(cfg)
		sprite2 := textlabel.Render(cfg)

		// Both must produce results (or both nil for invalid bounds).
		if (sprite1 == nil) != (sprite2 == nil) {
			t.Fatalf("render inconsistency: sprite1=%v, sprite2=%v", sprite1, sprite2)
		}
		if sprite1 == nil {
			return // both nil, consistent
		}

		// Position must be identical.
		if sprite1.Position != sprite2.Position {
			t.Fatalf("Position mismatch: %v vs %v", sprite1.Position, sprite2.Position)
		}

		// Image dimensions must be identical.
		b1 := sprite1.Image.Bounds()
		b2 := sprite2.Image.Bounds()
		if b1 != b2 {
			t.Fatalf("Image bounds mismatch: %v vs %v", b1, b2)
		}

		// Pixel data must be identical.
		for py := b1.Min.Y; py < b1.Max.Y; py++ {
			for px := b1.Min.X; px < b1.Max.X; px++ {
				r1, g1, b1v, a1 := sprite1.Image.At(px, py).RGBA()
				r2, g2, b2v, a2 := sprite2.Image.At(px, py).RGBA()
				if r1 != r2 || g1 != g2 || b1v != b2v || a1 != a2 {
					t.Fatalf("pixel (%d,%d) mismatch: (%d,%d,%d,%d) vs (%d,%d,%d,%d)",
						px, py, r1, g1, b1v, a1, r2, g2, b2v, a2)
				}
			}
		}
	})
}
