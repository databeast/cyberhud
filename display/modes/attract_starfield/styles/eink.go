package styles

import (
	"image"
	"image/color"
	"math"
	"math/rand"

	"github.com/databeast/cyberhud/display/modes/attract_starfield/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/widgets"
)

// BuildEinkView produces a static starfield scatter frame suitable for e-ink panels.
// The output is deterministic for a given policy (stable cache key).
func BuildEinkView(panelWidth, panelHeight int, p source.Policy) style.ViewData {
	sprite := renderEinkStarfield(panelWidth, panelHeight, p)

	compositor := widgets.NewCompositor(widgets.SuppressionContext{
		IsEink:          true,
		AvailableWidth:  panelWidth,
		AvailableHeight: panelHeight,
	})
	compositor.Add(source.NewStaticSprite(sprite))

	return style.ViewData{
		Static:  true,
		Sprites: compositor.Sprites(),
	}
}

// renderEinkStarfield draws a static scatter of stars for e-ink display.
func renderEinkStarfield(panelWidth, panelHeight int, p source.Policy) *widgets.Sprite {
	img := image.NewRGBA(image.Rect(0, 0, panelWidth, panelHeight))

	seed := int64(p.Speed*100) + int64(p.Density*1000) + int64(p.Layers*10000)
	rng := rand.New(rand.NewSource(seed))

	centerX := float64(panelWidth) / 2.0
	centerY := float64(panelHeight) / 2.0
	maxDist := math.Sqrt(centerX*centerX + centerY*centerY)

	area := panelWidth * panelHeight
	count := int(float64(area) * p.Density / 512.0)
	if count < 4 {
		count = 4
	}
	if count > 300 {
		count = 300
	}

	for i := 0; i < count; i++ {
		angle := rng.Float64() * 2 * math.Pi
		dist := rng.Float64() * maxDist * 0.9

		sx := int(centerX + math.Cos(angle)*dist)
		sy := int(centerY + math.Sin(angle)*dist)

		if sx < 0 || sx >= panelWidth || sy < 0 || sy >= panelHeight {
			continue
		}

		brightness := uint8(80 + int(175.0*(dist/maxDist)))
		img.SetRGBA(sx, sy, color.RGBA{R: brightness, G: brightness, B: brightness, A: 255})
	}

	return &widgets.Sprite{
		Image:    img,
		Position: image.Point{},
		Label:    "starfield-eink",
	}
}
