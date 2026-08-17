package led

import (
	"image"
	"image/color"
	"math"
	"sort"

	"github.com/databeast/cyberhud/display/widgets/gradient"
)

// applyGradientFill fills body pixels using radially interpolated colors from center
// to edge. It uses the shape membership function to determine which pixels are inside
// the body, then for each inside pixel computes a normalized distance from center
// (0.0 = center, 1.0 = edge) and interpolates color between gradient stops at that
// distance, applying brightness scaling to the result.
//
// This function replaces the shape body renderer entirely when gradient is active and
// the LED is On (effectiveBrightness > 0).
func applyGradientFill(img *image.RGBA, bodyRect image.Rectangle, cfg Config, effectiveBrightness float64) {
	stops := prepareStops(cfg.Gradient.Stops)
	if len(stops) < 2 {
		// Shouldn't reach here (caller checks), but fall back to nothing.
		return
	}

	w := bodyRect.Dx()
	h := bodyRect.Dy()
	cx := float64(w) / 2.0
	cy := float64(h) / 2.0

	// Compute shape-specific radius for normalization.
	radius := float64(w) / 2.0
	if float64(h)/2.0 < radius {
		radius = float64(h) / 2.0
	}

	// For RoundedSquare, precompute corner radius.
	side := w
	if h < side {
		side = h
	}
	cornerRadius := float64(side / 4)

	// Corner arc centers (local coordinates) for RoundedSquare membership check.
	tlx, tly := cornerRadius, cornerRadius
	trx, try_ := float64(w)-cornerRadius, cornerRadius
	blx, bly := cornerRadius, float64(h)-cornerRadius
	brx, bry := float64(w)-cornerRadius, float64(h)-cornerRadius

	for py := bodyRect.Min.Y; py < bodyRect.Max.Y; py++ {
		for px := bodyRect.Min.X; px < bodyRect.Max.X; px++ {
			lx := px - bodyRect.Min.X
			ly := py - bodyRect.Min.Y

			// Pixel center in local coordinates.
			pcx := float64(lx) + 0.5
			pcy := float64(ly) + 0.5

			// Check shape membership and compute normalized distance.
			var inside bool
			var normDist float64

			switch cfg.Shape {
			case Square:
				// All pixels within bodyRect are inside the square.
				inside = true
				// Chebyshev distance normalized by halfSide.
				halfSide := float64(w) / 2.0
				if float64(h)/2.0 < halfSide {
					halfSide = float64(h) / 2.0
				}
				dx := math.Abs(pcx - cx)
				dy := math.Abs(pcy - cy)
				normDist = math.Max(dx, dy) / halfSide

			case Diamond:
				// Manhattan distance membership: |px-cx| + |py-cy| <= bodyRadius
				dist := math.Abs(pcx-cx) + math.Abs(pcy-cy)
				inside = dist <= radius
				if inside && radius > 0 {
					normDist = dist / radius
				}

			case RoundedSquare:
				inside = isInsideRoundedRect(pcx, pcy, float64(w), float64(h), cornerRadius, tlx, tly, trx, try_, blx, bly, brx, bry)
				if inside {
					// Approximate as Chebyshev with corner adjustment.
					halfSide := float64(w) / 2.0
					if float64(h)/2.0 < halfSide {
						halfSide = float64(h) / 2.0
					}
					dx := math.Abs(pcx - cx)
					dy := math.Abs(pcy - cy)
					normDist = math.Max(dx, dy) / halfSide
				}

			default: // Circle
				dx := pcx - cx
				dy := pcy - cy
				dist := math.Sqrt(dx*dx + dy*dy)
				inside = dist <= radius
				if inside && radius > 0 {
					normDist = dist / radius
				}
			}

			if !inside {
				continue
			}

			// Clamp normalized distance to [0, 1].
			if normDist < 0 {
				normDist = 0
			}
			if normDist > 1.0 {
				normDist = 1.0
			}

			// Interpolate color at this normalized distance.
			c := interpolateGradientColor(stops, normDist)

			// Apply brightness scaling.
			c = applyBrightnessToColor(c, effectiveBrightness)

			img.SetRGBA(px, py, c)
		}
	}
}

// prepareStops clamps stop positions to [0.0, 1.0] and sorts by position.
// Returns the prepared slice (does not mutate the input).
func prepareStops(stops []gradient.ColorStop) []gradient.ColorStop {
	out := make([]gradient.ColorStop, len(stops))
	copy(out, stops)

	// Clamp positions to [0.0, 1.0].
	for i := range out {
		if out[i].Position < 0.0 {
			out[i].Position = 0.0
		} else if out[i].Position > 1.0 {
			out[i].Position = 1.0
		}
	}

	// Stable sort by position.
	sort.SliceStable(out, func(i, j int) bool {
		return out[i].Position < out[j].Position
	})

	return out
}

// interpolateGradientColor linearly interpolates between adjacent gradient stops
// at the given normalized position t in [0.0, 1.0].
// Position ≤ first stop → first stop color.
// Position ≥ last stop → last stop color.
func interpolateGradientColor(stops []gradient.ColorStop, t float64) color.RGBA {
	n := len(stops)
	if n == 0 {
		return color.RGBA{}
	}

	// At or before first stop.
	if t <= stops[0].Position {
		return stops[0].Color
	}

	// At or beyond last stop.
	if t >= stops[n-1].Position {
		return stops[n-1].Color
	}

	// Find bracketing stops.
	left := 0
	for i := 1; i < n; i++ {
		if stops[i].Position <= t {
			left = i
		} else {
			break
		}
	}
	right := left + 1

	// Same position — use the right stop color.
	if stops[left].Position == stops[right].Position {
		return stops[right].Color
	}

	// Linear interpolation fraction.
	fraction := (t - stops[left].Position) / (stops[right].Position - stops[left].Position)

	r := math.Round(float64(stops[left].Color.R) + (float64(stops[right].Color.R)-float64(stops[left].Color.R))*fraction)
	g := math.Round(float64(stops[left].Color.G) + (float64(stops[right].Color.G)-float64(stops[left].Color.G))*fraction)
	b := math.Round(float64(stops[left].Color.B) + (float64(stops[right].Color.B)-float64(stops[left].Color.B))*fraction)
	a := math.Round(float64(stops[left].Color.A) + (float64(stops[right].Color.A)-float64(stops[left].Color.A))*fraction)

	return color.RGBA{
		R: uint8(r),
		G: uint8(g),
		B: uint8(b),
		A: uint8(a),
	}
}

// applyBrightnessToColor scales RGB channels by brightness, preserving alpha.
func applyBrightnessToColor(c color.RGBA, brightness float64) color.RGBA {
	if brightness >= 1.0 {
		return c
	}
	if brightness <= 0.0 {
		return color.RGBA{R: 0, G: 0, B: 0, A: c.A}
	}
	return color.RGBA{
		R: uint8(math.Floor(float64(c.R) * brightness)),
		G: uint8(math.Floor(float64(c.G) * brightness)),
		B: uint8(math.Floor(float64(c.B) * brightness)),
		A: c.A,
	}
}

// shouldApplyGradient returns true when the gradient fill should replace
// the normal shape body renderer.
func shouldApplyGradient(cfg Config, isOff bool) bool {
	if isOff {
		return false
	}
	if cfg.Gradient == nil {
		return false
	}
	if len(cfg.Gradient.Stops) < 2 {
		return false
	}
	return true
}
