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
// e-ink displays. The frame shows a frozen oscilloscope pattern representative
// of the mode's visual theme.
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

	baseCycles := 3.0 + p.Density*2.5
	traceCount := p.Traces
	if traceCount < 1 {
		traceCount = 1
	}
	if traceCount > 3 {
		traceCount = 3
	}

	halfHeight := float64(panelHeight) / 2.0
	baseAmplitude := p.Amplitude * halfHeight
	traceThickness := 2
	if panelWidth < 96 || panelHeight < 64 {
		traceThickness = 1
	}

	for t := 0; t < traceCount; t++ {
		phase := float64(t)*1.7 + 0.4
		shapeMix := math.Mod(float64(t)*0.31+0.17, 1.0)
		cycles := math.Round(baseCycles + 0.75*float64(t))
		if cycles < 2 {
			cycles = 2
		}
		ampScale := baseAmplitude * (0.9 + 0.18*float64(t))
		if ampScale <= 0 {
			ampScale = baseAmplitude
		}
		drawOscilloscopeTrace(img, panelWidth, panelHeight, false, cycles, phase, ampScale, shapeMix, color.RGBA{255, 255, 255, 255}, traceThickness)
	}

	return []widgets.Sprite{
		{Image: img, Position: image.Point{X: 0, Y: 0}, Label: "waveform-eink"},
	}
}
