package attract_popsquares

import (
	"image/color"
	"testing"

	"github.com/databeast/cyberhud/display/modes/attract_popsquares/source"
)

// TestTunePolicyForUltraLowResCapsActivity verifies the ultra-low-res tuning
// keeps tiny panels calm rather than noisy.
func TestTunePolicyForUltraLowResCapsActivity(t *testing.T) {
	p := tunePolicyForUltraLowRes(source.Policy{Speed: 2.0, PopRate: 0.9, GridScale: 0.5}, 16, 8)
	if p.PopRate > 0.3 || p.Speed > 0.8 || p.GridScale < 1.0 {
		t.Fatalf("ultra-low-res policy not calmed: %+v", p)
	}
}

// TestGridFillsPanelEdgeToEdge verifies the grid is full-bleed: every pixel is
// painted (no background showing through, no gaps).
func TestGridFillsPanelEdgeToEdge(t *testing.T) {
	p := source.DefaultPolicy()
	sprite := renderGrid(p, 320, 240, 1.5, false)
	img := sprite.Image
	b := img.Bounds()
	// Sample the four corners and the centre: all must be painted (opaque).
	pts := [][2]int{{0, 0}, {b.Dx() - 1, 0}, {0, b.Dy() - 1}, {b.Dx() - 1, b.Dy() - 1}, {b.Dx() / 2, b.Dy() / 2}}
	for _, pt := range pts {
		c := color.NRGBAModel.Convert(img.At(pt[0], pt[1])).(color.NRGBA)
		if c.A == 0 {
			t.Fatalf("pixel (%d,%d) is transparent — grid is not full-bleed", pt[0], pt[1])
		}
	}
}

// TestFieldIsCalm verifies the core "idling" invariant: at the default pop
// rate, only a MINORITY of cells are brightly lit. A light show (most cells
// bright) would break the supercomputer-idling vibe.
func TestFieldIsCalm(t *testing.T) {
	p := source.DefaultPolicy()
	sprite := renderGrid(p, 320, 240, 1.5, false)
	img := sprite.Image
	b := img.Bounds()

	var bright, total int
	for y := b.Min.Y; y < b.Max.Y; y += 4 {
		for x := b.Min.X; x < b.Max.X; x += 4 {
			r, g, bl, _ := img.At(x, y).RGBA()
			lum := 0.299*float64(r>>8) + 0.587*float64(g>>8) + 0.114*float64(bl>>8)
			if lum > 180 {
				bright++
			}
			total++
		}
	}
	frac := float64(bright) / float64(total)
	if frac > 0.4 {
		t.Fatalf("field is too busy: %.1f%% of samples are brightly lit — should be a sparse idle twinkle", frac*100)
	}
}

// TestPulseDoesNotRepeatAcrossCells guards against a synchronised pattern: two
// different cells should not pulse in lock-step. We check that a spread of cells
// produce a spread of pulse values at a fixed instant (not all identical).
func TestPulseDoesNotRepeatAcrossCells(t *testing.T) {
	seen := map[float64]int{}
	for gy := 0; gy < 8; gy++ {
		for gx := 0; gx < 8; gx++ {
			h := cellHash(gx, gy)
			// Quantise to buckets; a synchronised grid would collapse to few buckets.
			seen[float64(int(h*20))]++
		}
	}
	if len(seen) < 10 {
		t.Fatalf("cell phases are not well scattered: only %d distinct buckets across 64 cells", len(seen))
	}
}

// TestFrozenFrameIsNonBlank verifies the frozen-instant grid (used by the e-ink
// path) fills the panel and is neither all-black nor all-white — a calm field
// with some variation. Tests renderGrid directly to avoid fabricating panel
// hints, matching the plasma test's approach.
func TestFrozenFrameIsNonBlank(t *testing.T) {
	p := source.DefaultPolicy()
	sprite := renderGrid(p, 400, 300, 2.3, true) // mono, frozen instant — the e-ink content
	if sprite.Image == nil {
		t.Fatal("frozen frame produced no image")
	}
	b := sprite.Image.Bounds()
	var black, white, total int
	for y := b.Min.Y; y < b.Max.Y; y += 3 {
		for x := b.Min.X; x < b.Max.X; x += 3 {
			r, g, bl, _ := sprite.Image.At(x, y).RGBA()
			lum := (float64(r>>8) + float64(g>>8) + float64(bl>>8)) / 3.0
			switch {
			case lum < 8:
				black++
			case lum > 247:
				white++
			}
			total++
		}
	}
	if black == total {
		t.Fatal("frozen frame is entirely black")
	}
	if white == total {
		t.Fatal("frozen frame is entirely white")
	}
}
