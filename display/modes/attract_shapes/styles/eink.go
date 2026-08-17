package styles

import (
	"image"
	"image/color"
	"math"
	"math/rand"

	"github.com/databeast/cyberhud/display/modes/attract_shapes/source"
	"github.com/databeast/cyberhud/display/widgets"
)

// buildEinkSprites generates a static decorative shapes frame suitable for
// e-ink displays. The frame shows scattered geometric shapes representative
// of the mode's visual theme without animation.
func BuildEinkSprites(p source.Policy, panelWidth, panelHeight int) []widgets.Sprite {
	if panelWidth <= 0 || panelHeight <= 0 {
		return []widgets.Sprite{{Image: image.NewRGBA(image.Rect(0, 0, 1, 1)), Position: image.Point{}}}
	}

	img := image.NewRGBA(image.Rect(0, 0, panelWidth, panelHeight))

	// Background is black (zero-value RGBA). Modes render with standard
	// "light on dark" convention — the panel driver handles any inversion
	// needed for the physical display technology.

	count := p.ShapeCount
	if count < 1 {
		count = 1
	}
	if count > 50 {
		count = 50
	}

	sides := p.Complexity
	if sides < 3 {
		sides = 3
	}
	if sides > 8 {
		sides = 8
	}

	minDim := panelWidth
	if panelHeight < minDim {
		minDim = panelHeight
	}
	baseRadius := float64(minDim) * 0.08

	// Use a fixed seed for deterministic e-ink rendering.
	rng := rand.New(rand.NewSource(42))

	for i := 0; i < count; i++ {
		cx := rng.Float64() * float64(panelWidth)
		cy := rng.Float64() * float64(panelHeight)
		radius := baseRadius * (0.5 + rng.Float64())
		angle := rng.Float64() * 2 * math.Pi

		// Draw white outlines on dark background.
		source.DrawPolygonOutline(img, cx, cy, radius, sides, angle, color.RGBA{255, 255, 255, 255})
	}

	return []widgets.Sprite{
		{Image: img, Position: image.Point{X: 0, Y: 0}, Label: "shapes-eink"},
	}
}
