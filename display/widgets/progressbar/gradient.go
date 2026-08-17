package progressbar

import (
	"image/color"
	"math"
	"sort"

	"github.com/databeast/cyberhud/display/widgets/gradient"
)

// interpolateGradient finds the two bracketing stops for the given normalized position t
// and linearly interpolates each RGBA channel independently.
// Stops are sorted by position before interpolation, and positions are clamped to [0, 1].
// This reimplements the logic from the gradient package's unexported interpolateColor.
func interpolateGradient(stops []gradient.ColorStop, t float64) color.RGBA {
	// Make a working copy, clamp positions, and sort by position.
	work := make([]gradient.ColorStop, len(stops))
	copy(work, stops)

	for i := range work {
		if work[i].Position < 0.0 {
			work[i].Position = 0.0
		} else if work[i].Position > 1.0 {
			work[i].Position = 1.0
		}
	}

	sort.SliceStable(work, func(i, j int) bool {
		return work[i].Position < work[j].Position
	})

	n := len(work)

	// For t at or before the first stop, return the first stop's color
	// (handles duplicate positions at the start by using the last one at that position).
	if t <= work[0].Position {
		last := 0
		for i := 1; i < n; i++ {
			if work[i].Position == work[0].Position {
				last = i
			} else {
				break
			}
		}
		return work[last].Color
	}

	// For t >= last stop position: return last stop color.
	if t >= work[n-1].Position {
		return work[n-1].Color
	}

	// Find the two bracketing stops for t.
	// We want the last stop with Position <= t as the left bracket,
	// and the first stop with Position > t as the right bracket.
	left := 0
	for i := 1; i < n; i++ {
		if work[i].Position <= t {
			left = i
		} else {
			break
		}
	}

	right := left + 1

	// Guard against same-position brackets (shouldn't happen given the loop above).
	if work[left].Position == work[right].Position {
		return work[right].Color
	}

	// Compute interpolation fraction between the two bracketing stops.
	fraction := (t - work[left].Position) / (work[right].Position - work[left].Position)

	// Linearly interpolate each RGBA channel and round to nearest uint8.
	r := math.Round(float64(work[left].Color.R) + (float64(work[right].Color.R)-float64(work[left].Color.R))*fraction)
	g := math.Round(float64(work[left].Color.G) + (float64(work[right].Color.G)-float64(work[left].Color.G))*fraction)
	b := math.Round(float64(work[left].Color.B) + (float64(work[right].Color.B)-float64(work[left].Color.B))*fraction)
	a := math.Round(float64(work[left].Color.A) + (float64(work[right].Color.A)-float64(work[left].Color.A))*fraction)

	return color.RGBA{
		R: uint8(r),
		G: uint8(g),
		B: uint8(b),
		A: uint8(a),
	}
}
