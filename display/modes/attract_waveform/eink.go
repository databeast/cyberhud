package attract_waveform

import (
	"image"
	"image/color"
	"math"

	"github.com/databeast/cyberhud/display/modes/attract_waveform/source"
	"github.com/databeast/cyberhud/display/widgets"
)

// einkFrameCounter is NOT incremented on e-ink panels — the static frame
// never changes unless the policy changes.
// The RenderCacheKey for e-ink uses only the policy fingerprint.

// buildEinkSprites generates a static decorative waveform frame suitable for
// e-ink displays. It draws a single frozen snapshot of the living-signal model
// at a fixed instant — no animation, no phosphor accumulation (e-ink cannot
// show a fading trail), just one slice of the non-repeating signal rendered as
// crisp white traces.
func buildEinkSprites(p source.Policy, panelWidth, panelHeight int) []widgets.Sprite {
	if panelWidth <= 0 || panelHeight <= 0 {
		return []widgets.Sprite{{Image: image.NewRGBA(image.Rect(0, 0, 1, 1)), Position: image.Point{}}}
	}

	img := image.NewRGBA(image.Rect(0, 0, panelWidth, panelHeight))
	for y := 0; y < panelHeight; y++ {
		for x := 0; x < panelWidth; x++ {
			img.SetRGBA(x, y, color.RGBA{0, 0, 0, 255})
		}
	}

	traceCount := p.Traces
	if traceCount < 1 {
		traceCount = 1
	}
	if traceCount > 3 {
		traceCount = 3
	}

	baseFreq := 0.8 + 0.6*clamp01(p.Density)
	halfHeight := float64(panelHeight) / 2.0
	baseAmplitude := clamp01(p.Amplitude) * halfHeight * 0.92

	// E-ink has no glow/persistence to carry a thin beam, so draw the frozen
	// traces boldly: thin strokes look spindly on a crisp e-paper panel, and a
	// sparse frozen frame reads as weak. Scale thickness with panel size.
	thickness := 4
	if panelWidth < 300 || panelHeight < 200 {
		thickness = 3
	}
	if panelWidth < 128 || panelHeight < 96 {
		thickness = 2
	}
	if panelWidth < 96 || panelHeight < 64 {
		thickness = 1
	}

	// Fixed instant for the frozen frame.
	const frozenT = 1.37

	denom := panelWidth - 1
	if denom < 1 {
		denom = 1
	}

	white := color.RGBA{255, 255, 255, 255}

	for t := 0; t < traceCount; t++ {
		ft := float64(t)
		seed := 0.73*ft + 0.19
		sp := buildPartials(t, baseFreq)

		prevX, prevY := -1, -1
		for s := 0; s <= denom; s++ {
			u := float64(s) / float64(denom)
			env := macroEnvelope(u, frozenT, seed)
			v := sp.sample(u, frozenT) * env
			ay := int(math.Round(halfHeight - 0.5 + v*baseAmplitude))
			if ay < 0 {
				ay = 0
			}
			if ay >= panelHeight {
				ay = panelHeight - 1
			}
			drawEinkDot(img, s, ay, white, thickness)
			if prevX >= 0 {
				drawEinkLine(img, prevX, prevY, s, ay, white, thickness)
			}
			prevX, prevY = s, ay
		}
	}

	return []widgets.Sprite{
		{Image: img, Position: image.Point{X: 0, Y: 0}, Label: "waveform-eink"},
	}
}

func drawEinkDot(img *image.RGBA, x, y int, c color.RGBA, thickness int) {
	if thickness < 1 {
		thickness = 1
	}
	// radius grows with thickness: 1→0, 2→1, 3→1, 4→2 ... Use (thickness)/2 but
	// ensure thickness≥2 yields at least radius 1 so strokes read as solid.
	radius := thickness / 2
	b := img.Bounds()
	for dy := -radius; dy <= radius; dy++ {
		for dx := -radius; dx <= radius; dx++ {
			px, py := x+dx, y+dy
			if px < b.Min.X || px >= b.Max.X || py < b.Min.Y || py >= b.Max.Y {
				continue
			}
			img.SetRGBA(px, py, c)
		}
	}
}

func drawEinkLine(img *image.RGBA, x0, y0, x1, y1 int, c color.RGBA, thickness int) {
	dx := int(math.Abs(float64(x1 - x0)))
	dy := -int(math.Abs(float64(y1 - y0)))
	sx := -1
	if x0 < x1 {
		sx = 1
	}
	sy := -1
	if y0 < y1 {
		sy = 1
	}
	err := dx + dy
	for {
		drawEinkDot(img, x0, y0, c, thickness)
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * err
		if e2 >= dy {
			err += dy
			x0 += sx
		}
		if e2 <= dx {
			err += dx
			y0 += sy
		}
	}
}
