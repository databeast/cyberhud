package source

import (
	"image"
	"image/color"
	"math"
	"math/rand"
)

// renderEinkFrame creates a static grayscale bokeh scatter image.
func renderEinkFrame(panelWidth, panelHeight int, p Policy) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, panelWidth, panelHeight))

	shortDim := panelWidth
	if panelHeight < shortDim {
		shortDim = panelHeight
	}

	count := CircleCount(panelWidth, panelHeight, p.Density)

	// Use a seeded RNG for deterministic e-ink frames.
	rng := rand.New(rand.NewSource(42))

	for i := 0; i < count; i++ {
		minRadius := 0.02 * float64(shortDim)
		maxRadius := 0.15 * float64(shortDim)
		midRadius := (minRadius + maxRadius) / 2.0
		spread := (maxRadius - minRadius) / 2.0
		radius := midRadius + (rng.Float64()*2.0-1.0)*spread*p.SizeVariance
		if radius < minRadius {
			radius = minRadius
		}
		if radius > maxRadius {
			radius = maxRadius
		}

		x := rng.Float64() * float64(panelWidth)
		y := rng.Float64() * float64(panelHeight)
		opacity := 0.15 + rng.Float64()*0.45

		// Draw grayscale circle with soft edges.
		drawEinkCircle(img, x, y, radius, opacity)
	}

	return img
}

// drawEinkCircle renders a single grayscale circle with radial gradient onto the image.
func drawEinkCircle(img *image.RGBA, cx, cy, radius, opacity float64) {
	bounds := img.Bounds()
	r := radius
	if r < 1 {
		r = 1
	}

	minX := int(math.Floor(cx - r))
	maxX := int(math.Ceil(cx + r))
	minY := int(math.Floor(cy - r))
	maxY := int(math.Ceil(cy + r))

	if minX < bounds.Min.X {
		minX = bounds.Min.X
	}
	if maxX > bounds.Max.X {
		maxX = bounds.Max.X
	}
	if minY < bounds.Min.Y {
		minY = bounds.Min.Y
	}
	if maxY > bounds.Max.Y {
		maxY = bounds.Max.Y
	}

	for py := minY; py < maxY; py++ {
		for px := minX; px < maxX; px++ {
			dx := float64(px) + 0.5 - cx
			dy := float64(py) + 0.5 - cy
			dist := math.Sqrt(dx*dx + dy*dy)
			if dist >= r {
				continue
			}

			alpha := (1.0 - dist/r) * opacity
			lum := uint8(alpha * 255.0)

			src := color.RGBA{R: lum, G: lum, B: lum, A: lum}
			existing := img.RGBAAt(px, py)
			img.SetRGBA(px, py, alphaBlend(existing, src))
		}
	}
}
