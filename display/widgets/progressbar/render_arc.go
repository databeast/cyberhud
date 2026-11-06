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

	// Convert geom.StartAngle into this renderer's angle convention.
	//
	// The pixel test below measures angle with atan2(dx, -dy), which is clockwise
	// from 12-o'clock: top = 0, right = +π/2, bottom = ±π, left = -π/2. But
	// geom.StartAngle is expressed in the shared standard-math convention used by
	// the Ring and Pie renderers (-π/2 = 12-o'clock, π = 9-o'clock). The two are
	// offset by 90°, so StartAngle must be converted with the same +π/2 shift that
	// render_pie.go and render_ring.go apply before use. Omitting this conversion
	// (the previous behavior) centered the arc on 9-o'clock instead of 12-o'clock,
	// which put the 270° gauge's gap on the wrong side and made the fill grow the
	// wrong way — the arc appeared rotated/flipped from the configured orientation.
	startOffset := geom.StartAngle + math.Pi/2.0

	// The arc track is centered symmetrically about the reference angle (12-o'clock
	// for horizontal), so for a 270° sweep the 90° gap sits at the bottom — the
	// familiar car-speedometer layout. arcStart is the clockwise beginning of the
	// track (its lower-left endpoint), from which the fill grows clockwise.
	arcStart := startOffset - sweepRad/2.0

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
