package progressbar

import (
	"image"
	"math"

	"github.com/databeast/cyberhud/display/widgets"
)

// renderRing draws a hollow circular ring (donut) progress indicator.
// The ring is inscribed within cfg.Bounds using the smaller dimension as the
// outer diameter. Fill progresses clockwise from geom.StartAngle (12-o'clock
// for horizontal, 9-o'clock for vertical) proportional to cfg.Value.
//
// Unfilled annulus pixels (the "track") are drawn in the background color.
// Pixels outside the annulus are also background.
// When a valid GradientFill (≥2 stops) is provided, the fill color is
// interpolated along the angular extent of the fill arc.
func renderRing(cfg Config, geom RenderGeometry) *image.RGBA {
	w := cfg.Bounds.Dx()
	h := cfg.Bounds.Dy()
	img := image.NewRGBA(image.Rect(0, 0, w, h))

	fg, bg := widgets.ResolveColors(cfg.Foreground, cfg.Background)

	// Determine if we have a valid gradient (≥2 stops).
	useGradient := cfg.Gradient != nil && len(cfg.Gradient.Stops) >= 2

	// Ring geometry.
	cx := float64(w) / 2.0
	cy := float64(h) / 2.0
	outerR := math.Min(float64(w), float64(h)) / 2.0
	innerR := outerR - float64(cfg.Thickness)
	if innerR < 0 {
		innerR = 0
	}

	// Sweep angle: the angular extent of the fill.
	sweepAngle := cfg.Value * 2.0 * math.Pi

	// Compute the start angle offset for clockwise-from-north calculations.
	// atan2(dx, -dy) gives clockwise angle from north (12-o'clock) in [−π, π].
	// We normalize to [0, 2π), then subtract startAngleOffset to get the
	// relative angle from the configured start position.
	//
	// For Horizontal: startAngle = -π/2, offset from north = 0
	//   (north IS 12-o'clock, so no offset needed)
	// For Vertical: startAngle = π (9-o'clock), offset from north = 3π/2
	//   (9-o'clock is 270° clockwise from 12-o'clock)
	var startOffset float64
	if cfg.Orientation == OrientVertical {
		startOffset = 3.0 * math.Pi / 2.0
	}
	// For Horizontal, startOffset = 0 (12-o'clock is the natural zero of atan2(dx, -dy))

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			// Use pixel center for sub-pixel accuracy.
			dx := float64(x) + 0.5 - cx
			dy := float64(y) + 0.5 - cy
			dist := math.Sqrt(dx*dx + dy*dy)

			if dist < innerR || dist > outerR {
				// Outside the annulus → background.
				img.SetRGBA(x, y, bg)
			} else {
				// Within the annulus — determine if this pixel is in the fill arc.
				// Compute clockwise angle from 12-o'clock (north).
				angle := math.Atan2(dx, -dy)
				// Normalize to [0, 2π).
				if angle < 0 {
					angle += 2.0 * math.Pi
				}

				// Adjust for the start offset (orientation).
				relative := angle - startOffset
				if relative < 0 {
					relative += 2.0 * math.Pi
				}

				if relative < sweepAngle {
					// Filled portion.
					if useGradient {
						// Gradient position: 0.0 at start angle, 1.0 at completion.
						var t float64
						if sweepAngle > 0 {
							t = relative / sweepAngle
						}
						c := interpolateGradient(cfg.Gradient.Stops, t)
						img.SetRGBA(x, y, c)
					} else {
						img.SetRGBA(x, y, fg)
					}
				} else {
					// Track (unfilled annulus) → background.
					img.SetRGBA(x, y, bg)
				}
			}
		}
	}

	return img
}
