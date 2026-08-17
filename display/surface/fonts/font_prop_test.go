package font_test

import (
	"image"
	"image/color"
	"strings"
	"testing"

	"github.com/databeast/cyberhud/display/surface/fonts"
	"pgregory.net/rapid"
)

// --- From: draw_prop_test.go ---

// For any valid font Face, character, position (x, y), foreground color, and
// target image, DrawGlyph sets exactly those pixels where the corresponding
// GlyphRow bit is 1 AND the pixel coordinate falls within image bounds,
// and does NOT modify any other pixels.

func genSmallImage(t *rapid.T) *image.RGBA {
	w := rapid.IntRange(5, 50).Draw(t, "imgWidth")
	h := rapid.IntRange(5, 50).Draw(t, "imgHeight")
	return image.NewRGBA(image.Rect(0, 0, w, h))
}

func genPosition(t *rapid.T) (int, int) {
	x := rapid.IntRange(-20, 60).Draw(t, "x")
	y := rapid.IntRange(-20, 60).Draw(t, "y")
	return x, y
}

func genRune(t *rapid.T) rune {
	return rune(rapid.IntRange(32, 126).Draw(t, "char"))
}

func genColor(t *rapid.T) color.RGBA {
	// Generate a non-black color so we can distinguish set pixels from background
	r := rapid.IntRange(1, 255).Draw(t, "colorR")
	g := rapid.IntRange(0, 255).Draw(t, "colorG")
	b := rapid.IntRange(0, 255).Draw(t, "colorB")
	a := rapid.IntRange(1, 255).Draw(t, "colorA")
	return color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
}

func TestProp_DrawGlyphPixelFidelity(t *testing.T) {
	face := font.Default()
	if face == nil {
		t.Skip("no default face registered")
	}

	rapid.Check(t, func(t *rapid.T) {
		img := genSmallImage(t)
		x, y := genPosition(t)
		ch := genRune(t)
		fg := genColor(t)

		bounds := img.Bounds()
		metrics := face.Metrics()

		// Save a copy of the image before DrawGlyph
		original := image.NewRGBA(bounds)
		copy(original.Pix, img.Pix)

		// Call DrawGlyph
		font.DrawGlyph(img, face, ch, x, y, fg)

		// Compute the expected set of pixels that should be modified.
		// A pixel (px, py) should be set if:
		//   1. It corresponds to a bit=1 in GlyphRow for the character
		//   2. It falls within the image bounds
		type pixel struct{ x, y int }
		expectedSet := make(map[pixel]bool)

		for row := 0; row < metrics.GlyphHeight; row++ {
			py := y + row
			if py < bounds.Min.Y || py >= bounds.Max.Y {
				continue
			}
			bits := face.GlyphRow(ch, row)
			for col := 0; col < metrics.GlyphWidth; col++ {
				px := x + col
				if px < bounds.Min.X || px >= bounds.Max.X {
					continue
				}
				if bits&(1<<uint(31-col)) != 0 {
					expectedSet[pixel{px, py}] = true
				}
			}
		}

		// Verify: every pixel in expectedSet is set to the expected color.
		// When fg.A == 255, the pixel should match fg exactly.
		// When fg.A < 255, the pixel is alpha-blended against the original dst.
		for p := range expectedSet {
			got := img.RGBAAt(p.x, p.y)
			var want color.RGBA
			if fg.A == 255 {
				want = fg
			} else {
				// Source-over blend against original pixel (zero on a fresh image).
				dst := original.RGBAAt(p.x, p.y)
				aa := uint32(fg.A)
				invA := 255 - aa
				want = color.RGBA{
					R: uint8((uint32(fg.R)*aa + uint32(dst.R)*invA) / 255),
					G: uint8((uint32(fg.G)*aa + uint32(dst.G)*invA) / 255),
					B: uint8((uint32(fg.B)*aa + uint32(dst.B)*invA) / 255),
					A: uint8(aa + (uint32(dst.A)*invA)/255),
				}
			}
			if got != want {
				t.Fatalf("pixel (%d,%d) should be %v but got %v (ch=%q, pos=(%d,%d), fg=%v)",
					p.x, p.y, want, got, ch, x, y, fg)
			}
		}

		// Verify: no other pixels were modified
		for py := bounds.Min.Y; py < bounds.Max.Y; py++ {
			for px := bounds.Min.X; px < bounds.Max.X; px++ {
				if expectedSet[pixel{px, py}] {
					continue
				}
				got := img.RGBAAt(px, py)
				orig := original.RGBAAt(px, py)
				if got != orig {
					t.Fatalf("pixel (%d,%d) was modified but shouldn't be: got %v, original %v (ch=%q, pos=(%d,%d))",
						px, py, got, orig, ch, x, y)
				}
			}
		}
	})
}

// --- From: font_prop_test.go ---

// For any requested pixel height (positive integer), the ByHeight function returns
// a registered font variant such that:
// (a) the variant's GlyphHeight is the largest among all registered variants that
//     does not exceed the requested height, AND
// (b) when multiple variants share that GlyphHeight, the variant from the highest-priority
//     family (Spleen > Terminus > Cozette) is selected.
// (c) If the requested height is smaller than all registered variants, the smallest
//     variant is returned.
// (d) If larger than all, the largest is returned.

func familyPriority(id string) int {
	lower := strings.ToLower(id)
	switch {
	case strings.HasPrefix(lower, "spleen-"):
		return 3
	case strings.HasPrefix(lower, "terminus-"):
		return 2
	case strings.HasPrefix(lower, "cozette-"):
		return 1
	default:
		return 0
	}
}

func TestByHeightSelectionInvariant(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		px := rapid.IntRange(1, 200).Draw(t, "px")

		result := font.ByHeight(px)
		if result == nil {
			t.Fatal("ByHeight returned nil")
		}

		all := font.List()
		if len(all) == 0 {
			t.Skip("no fonts registered")
		}

		resultHeight := result.Metrics().GlyphHeight
		resultPriority := familyPriority(result.ID())

		// Determine the set of faces with GlyphHeight <= px.
		var bestHeight int
		hasFitting := false
		for _, f := range all {
			h := f.Metrics().GlyphHeight
			if h <= px {
				if !hasFitting || h > bestHeight {
					bestHeight = h
					hasFitting = true
				}
			}
		}

		if hasFitting {
			// (a) The result's GlyphHeight must be the largest that does not exceed px.
			if resultHeight != bestHeight {
				t.Fatalf("px=%d: result height=%d, want largest fitting height=%d (id=%q)",
					px, resultHeight, bestHeight, result.ID())
			}

			// (b) Among all faces at bestHeight, the result has highest family priority.
			var maxPriorityAtHeight int
			for _, f := range all {
				if f.Metrics().GlyphHeight == bestHeight {
					p := familyPriority(f.ID())
					if p > maxPriorityAtHeight {
						maxPriorityAtHeight = p
					}
				}
			}
			if resultPriority < maxPriorityAtHeight {
				t.Fatalf("px=%d: result priority=%d (id=%q), want max priority=%d at height=%d",
					px, resultPriority, result.ID(), maxPriorityAtHeight, bestHeight)
			}
		} else {
			// (c) No face fits: result must be the smallest registered variant.
			smallestHeight := all[0].Metrics().GlyphHeight
			for _, f := range all {
				h := f.Metrics().GlyphHeight
				if h < smallestHeight {
					smallestHeight = h
				}
			}
			if resultHeight != smallestHeight {
				t.Fatalf("px=%d: no fitting face; result height=%d, want smallest=%d (id=%q)",
					px, resultHeight, smallestHeight, result.ID())
			}

			// Among all faces at the smallest height, highest priority wins.
			var maxPriorityAtSmallest int
			for _, f := range all {
				if f.Metrics().GlyphHeight == smallestHeight {
					p := familyPriority(f.ID())
					if p > maxPriorityAtSmallest {
						maxPriorityAtSmallest = p
					}
				}
			}
			if resultPriority < maxPriorityAtSmallest {
				t.Fatalf("px=%d: fallback priority=%d (id=%q), want max priority=%d at height=%d",
					px, resultPriority, result.ID(), maxPriorityAtSmallest, smallestHeight)
			}
		}
	})
}
