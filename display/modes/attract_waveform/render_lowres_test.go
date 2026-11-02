package attract_waveform

import (
	"image/color"
	"math"
	"testing"

	"github.com/databeast/cyberhud/display/modes/attract_waveform/source"
)

func TestTunePolicyForUltraLowResCapsTraceLoad(t *testing.T) {
	p := tunePolicyForUltraLowRes(source.Policy{Speed: 1.3, Traces: 4, Persistence: 0.8, Amplitude: 0.9}, 16, 8)
	if p.Traces > 1 || p.Persistence > 0.25 || p.Amplitude > 0.65 || p.Speed > 0.85 {
		t.Fatalf("unexpected tuned policy: %+v", p)
	}
}

func TestBuildEinkSpritesFillsThePanel(t *testing.T) {
	p := source.DefaultPolicy()
	p.Traces = 3
	p.Amplitude = 0.8
	p.Density = 0.6

	sprites := buildEinkSprites(p, 400, 300)
	if len(sprites) != 1 {
		t.Fatalf("expected 1 sprite, got %d", len(sprites))
	}
	img := sprites[0].Image
	if img == nil {
		t.Fatal("expected image to be generated")
	}

	bounds := img.Bounds()
	totalPixels := bounds.Dx() * bounds.Dy()
	whitePixels := 0
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := color.NRGBAModel.Convert(img.At(x, y)).(color.NRGBA)
			if c.A > 0 && c.R > 200 && c.G > 200 && c.B > 200 {
				whitePixels++
			}
		}
	}

	if whitePixels <= totalPixels/20 {
		t.Fatalf("e-ink waveform was too sparse: %d white pixels out of %d (%.2f%%)", whitePixels, totalPixels, float64(whitePixels)/float64(totalPixels)*100)
	}
}

func TestBuildSpritesTouchesBothEdges(t *testing.T) {
	ResetSnapshotState()
	p := source.DefaultPolicy()
	p.Traces = 3
	p.Amplitude = 0.75
	p.Density = 0.6

	sprites := buildSprites(p, 240, 135, false)
	if len(sprites) != 1 || sprites[0].Image == nil {
		t.Fatal("expected one rendered sprite")
	}

	img := sprites[0].Image
	bounds := img.Bounds()
	left := 0
	right := 0
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		if c := color.NRGBAModel.Convert(img.At(bounds.Min.X, y)).(color.NRGBA); c.A > 0 && (c.R > 0 || c.G > 0 || c.B > 0) {
			left++
		}
		if c := color.NRGBAModel.Convert(img.At(bounds.Max.X-1, y)).(color.NRGBA); c.A > 0 && (c.R > 0 || c.G > 0 || c.B > 0) {
			right++
		}
	}

	if left == 0 || right == 0 {
		t.Fatalf("waveform did not wrap to both edges: left=%d right=%d", left, right)
	}
}

func TestTraceColorsAreDistinctPerTrace(t *testing.T) {
	// Each trace must get a visibly different phosphor hue so the traces read as
	// separate signals rather than one tangle of the same colour.
	r0, g0, b0 := traceColor(0, 1.0)
	r1, g1, b1 := traceColor(1, 1.0)
	r2, g2, b2 := traceColor(2, 1.0)
	if r0 == r1 && g0 == g1 && b0 == b1 {
		t.Fatalf("traces 0 and 1 share a colour: (%v,%v,%v)", r0, g0, b0)
	}
	if r1 == r2 && g1 == g2 && b1 == b2 {
		t.Fatalf("traces 1 and 2 share a colour: (%v,%v,%v)", r1, g1, b1)
	}
}

// TestSignalDoesNotRepeatAcrossSweep guards the core aesthetic promise: within a
// single frame the trace must not be a short periodic tile (the old bug was ~5
// identical humps). We sample the first half and second half of the sweep and
// require they differ substantially — a repeating tile would make them match.
func TestSignalDoesNotRepeatAcrossSweep(t *testing.T) {
	sp := buildPartials(0, 0.8+0.6*0.5)
	const tFixed = 2.0
	var diff float64
	const n = 200
	for i := 0; i < n; i++ {
		u1 := 0.25 * float64(i) / float64(n) // first quarter
		u2 := u1 + 0.5                        // mirrored into third quarter
		d := sp.sample(u1, tFixed) - sp.sample(u2, tFixed)
		diff += math.Abs(d)
	}
	avg := diff / float64(n)
	if avg < 0.05 {
		t.Fatalf("signal looks periodic across the sweep (avg half-period diff %.4f) — it is repeating within the frame", avg)
	}
}

// TestSignalEvolvesOverTime guards that successive frames are not identical —
// the signal must never look twice-the-same.
func TestSignalEvolvesOverTime(t *testing.T) {
	sp := buildPartials(0, 0.8+0.6*0.5)
	var diff float64
	const n = 200
	for i := 0; i < n; i++ {
		u := float64(i) / float64(n)
		d := sp.sample(u, 1.0) - sp.sample(u, 2.5)
		diff += math.Abs(d)
	}
	avg := diff / float64(n)
	if avg < 0.02 {
		t.Fatalf("signal barely changes over time (avg diff %.4f) — it looks static/repeating", avg)
	}
}

func TestBuildSpritesUsesLargeVerticalRange(t *testing.T) {
	ResetSnapshotState()
	p := source.DefaultPolicy()
	p.Traces = 3
	p.Amplitude = 1.0
	p.Density = 0.6

	sprites := buildSprites(p, 240, 135, false)
	if len(sprites) != 1 || sprites[0].Image == nil {
		t.Fatal("expected one rendered sprite")
	}

	img := sprites[0].Image
	bounds := img.Bounds()
	minY := bounds.Max.Y
	maxY := bounds.Min.Y
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := color.NRGBAModel.Convert(img.At(x, y)).(color.NRGBA)
			if c.A > 0 && (c.R > 0 || c.G > 0 || c.B > 0) {
				if y < minY {
					minY = y
				}
				if y > maxY {
					maxY = y
				}
			}
		}
	}

	if maxY-minY < (bounds.Dy()*2)/3 {
		t.Fatalf("vertical excursion still too small: minY=%d maxY=%d height=%d", minY, maxY, bounds.Dy())
	}
}
