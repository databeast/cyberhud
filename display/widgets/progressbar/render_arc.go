package progressbar

import (
	"image"
	"math"

	"github.com/databeast/cyberhud/display/widgets"
)

// renderArc draws a partial circular arc (gauge) spanning the configured SweepAngle,
// filled proportional to Value from the start point determined by orientation.
//
// The arc is inscribed within cfg.Bounds using the smaller dimension as the outer
// diameter. Thickness and gradient logic are identical to the ring renderer.
//
// Horizontal orientation: fill begins from the left endpoint of the sweep.
// Vertical orientation: fill begins from the bottom endpoint of the sweep, progressing upward.
func renderArc(cfg Config, geom RenderGeometry) *image.RGBA {
	w := cfg.Bounds.Dx()
	h := cfg.Bounds.Dy()
	img := image.NewRGBA(image.Rect(0, 0, w, h))

	fg, bg := widgets.ResolveColors(cfg.Foreground, cfg.Background)

	// Center of the arc.
	cx := float64(w) / 2.0
	cy := float64(h) / 2.0

	// Outer radius inscribed within bounds.
	outerR := math.Min(float64(w), float64(h)) / 2.0
	// Inner radius from thickness.
	innerR := outerR - float64(cfg.Thickness)
	if innerR < 0 {
		innerR = 0
	}

	// Sweep angle in radians.
	sweepRad := cfg.SweepAngle * math.Pi / 180.0

	// The arc is centered symmetrically about the orientation reference angle.
	// arcStart is the angular beginning of the arc track (the "left endpoint").
	// geom.StartAngle is -π/2 for horizontal (12-o'clock) or π for vertical (9-o'clock).
	arcStart := geom.StartAngle - sweepRad/2.0

	// Fill extent based on value.
	fillExtent := cfg.Value * sweepRad

	// Determine if we have a valid gradient (≥2 stops).
	useGradient := cfg.Gradient != nil && len(cfg.Gradient.Stops) >= 2

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			// Use pixel center for calculations.
			dx := float64(x) + 0.5 - cx
			dy := float64(y) + 0.5 - cy
			dist := math.Sqrt(dx*dx + dy*dy)

			// Check if pixel is within the annulus.
			if dist > outerR || dist < innerR {
				img.SetRGBA(x, y, bg)
				continue
			}

			// Compute angle clockwise from 12-o'clock (standard screen coords).
			// atan2(dx, -dy) gives angle clockwise from 12-o'clock in [-π, π].
			angle := math.Atan2(dx, -dy)

			// Compute relative angle from arcStart, normalized to [0, 2π).
			relAngle := angle - arcStart
			// Normalize to [0, 2π).
			relAngle = math.Mod(relAngle, 2*math.Pi)
			if relAngle < 0 {
				relAngle += 2 * math.Pi
			}

			// Check if pixel is within the arc sweep.
			if relAngle > sweepRad {
				// Outside the arc track → background.
				img.SetRGBA(x, y, bg)
				continue
			}

			// Within the arc track: check if it's in the filled portion.
			if relAngle < fillExtent {
				// Filled portion.
				if useGradient {
					// Gradient position: 0.0 at arc start, 1.0 at fill completion.
					var t float64
					if fillExtent > 0 {
						t = relAngle / fillExtent
					}
					c := interpolateGradient(cfg.Gradient.Stops, t)
					img.SetRGBA(x, y, c)
				} else {
					img.SetRGBA(x, y, fg)
				}
			} else {
				// Unfilled portion of the arc track → background.
				img.SetRGBA(x, y, bg)
			}
		}
	}

	return img
}
