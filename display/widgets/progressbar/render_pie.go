package progressbar

import (
	"image"
	"math"

	"github.com/databeast/cyberhud/display/widgets"
)

// renderPie draws a solid pie (filled circle) swept clockwise from geom.StartAngle.
// Gradient is ignored per Req 1.6 — the pie always uses solid foreground color.
func renderPie(cfg Config, geom RenderGeometry) *image.RGBA {
	w := cfg.Bounds.Dx()
	h := cfg.Bounds.Dy()
	img := image.NewRGBA(image.Rect(0, 0, w, h))

	fg, bg := widgets.ResolveColors(cfg.Foreground, cfg.Background)

	// Center of the bounds (using pixel centers for sub-pixel accuracy).
	cx := float64(w) / 2.0
	cy := float64(h) / 2.0

	// Radius of the inscribed circle (smaller of the two dimensions).
	r := math.Min(float64(w), float64(h)) / 2.0

	// Convert geom.StartAngle to a clockwise offset from 12-o'clock.
	// geom.StartAngle is in standard math radians:
	//   -π/2 = 12-o'clock (top center)
	//   π    = 9-o'clock (left center)
	//
	// Our rendering uses a coordinate system where 0 = 12-o'clock measured clockwise.
	// To convert: offset = startAngle - (-π/2) = startAngle + π/2
	// Then normalize to [0, 2π).
	startOffset := geom.StartAngle + math.Pi/2.0
	// Normalize to [0, 2π)
	startOffset = math.Mod(startOffset, 2.0*math.Pi)
	if startOffset < 0 {
		startOffset += 2.0 * math.Pi
	}

	// Sweep angle proportional to value.
	sweepAngle := cfg.Value * 2.0 * math.Pi

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			// Use pixel center (x+0.5, y+0.5) for distance/angle calculations.
			dx := float64(x) + 0.5 - cx
			dy := float64(y) + 0.5 - cy
			dist := math.Sqrt(dx*dx + dy*dy)

			if dist > r {
				// Outside the circle → background.
				img.SetRGBA(x, y, bg)
			} else {
				// Compute angle clockwise from 12-o'clock (negative Y axis is 0°).
				angle := math.Atan2(dx, -dy)
				// Normalize to [0, 2π).
				if angle < 0 {
					angle += 2.0 * math.Pi
				}

				// Adjust for start offset.
				angle = angle - startOffset
				if angle < 0 {
					angle += 2.0 * math.Pi
				}

				if angle < sweepAngle {
					img.SetRGBA(x, y, fg)
				} else {
					img.SetRGBA(x, y, bg)
				}
			}
		}
	}

	return img
}
