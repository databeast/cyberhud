package attract_waveform

import (
	"image/color"
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

func TestWaveformTraceColorVariesWithTraceState(t *testing.T) {
	a := waveformTraceColor(3, 4.0, 0.25, 10.0)
	b := waveformTraceColor(6, 8.0, 3.25, 10.0)
	if a == b {
		t.Fatalf("expected color to vary with frequency/amplitude/phase, got identical colors: %+v", a)
	}
}

func TestBuildSpritesUsesLargeVerticalRange(t *testing.T) {
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
