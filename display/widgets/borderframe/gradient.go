package borderframe

import (
	"image"
	"image/color"
	"sort"

	"github.com/databeast/cyberhud/display/widgets/gradient"
)

// maxGradientStops is the maximum number of gradient stops supported.
const maxGradientStops = 64

// perimeterTileCount returns the number of tiles in a clockwise perimeter traversal.
// cols and rows are the number of tile columns and rows (width/8 and height/8).
func perimeterTileCount(cols, rows int) int {
	if cols <= 0 || rows <= 0 {
		return 0
	}
	if cols == 1 && rows == 1 {
		return 1
	}
	if cols == 1 {
		return rows
	}
	if rows == 1 {
		return cols
	}
	return 2 * (cols + rows - 2)
}

// perimeterTilePositions returns the (x, y) pixel positions of all tiles in clockwise
// perimeter order. Each position is the top-left pixel of the 8×8 tile.
//
// Traversal order:
//   - top-left corner
//   - top edge tiles (left to right)
//   - top-right corner
//   - right edge tiles (top to bottom)
//   - bottom-right corner
//   - bottom edge tiles (right to left)
//   - bottom-left corner
//   - left edge tiles (bottom to top)
func perimeterTilePositions(cols, rows int) []image.Point {
	total := perimeterTileCount(cols, rows)
	if total == 0 {
		return nil
	}

	positions := make([]image.Point, 0, total)

	// Top-left corner.
	positions = append(positions, image.Pt(0, 0))

	if cols == 1 && rows == 1 {
		return positions
	}

	// Handle degenerate single-column case.
	if cols == 1 {
		for row := 1; row < rows; row++ {
			positions = append(positions, image.Pt(0, row*tileSize))
		}
		return positions
	}

	// Handle degenerate single-row case.
	if rows == 1 {
		for col := 1; col < cols; col++ {
			positions = append(positions, image.Pt(col*tileSize, 0))
		}
		return positions
	}

	// Top edge tiles (left to right, excluding corners).
	for col := 1; col < cols-1; col++ {
		positions = append(positions, image.Pt(col*tileSize, 0))
	}

	// Top-right corner.
	positions = append(positions, image.Pt((cols-1)*tileSize, 0))

	// Right edge tiles (top to bottom, excluding corners).
	for row := 1; row < rows-1; row++ {
		positions = append(positions, image.Pt((cols-1)*tileSize, row*tileSize))
	}

	// Bottom-right corner.
	positions = append(positions, image.Pt((cols-1)*tileSize, (rows-1)*tileSize))

	// Bottom edge tiles (right to left, excluding corners).
	for col := cols - 2; col >= 1; col-- {
		positions = append(positions, image.Pt(col*tileSize, (rows-1)*tileSize))
	}

	// Bottom-left corner.
	positions = append(positions, image.Pt(0, (rows-1)*tileSize))

	// Left edge tiles (bottom to top, excluding corners).
	for row := rows - 2; row >= 1; row-- {
		positions = append(positions, image.Pt(0, row*tileSize))
	}

	return positions
}

// interpolateColor interpolates between sorted ColorStops at a given normalized
// position in [0.0, 1.0]. Stops must be sorted by Position before calling.
// If position is at or before the first stop, returns the first stop's color.
// If position is at or after the last stop, returns the last stop's color.
func interpolateColor(stops []gradient.ColorStop, position float64) color.RGBA {
	if len(stops) == 0 {
		return color.RGBA{}
	}
	if len(stops) == 1 {
		return stops[0].Color
	}

	// Clamp position to [0.0, 1.0].
	if position <= 0.0 {
		return stops[0].Color
	}
	if position >= 1.0 {
		return stops[len(stops)-1].Color
	}

	// Find the two stops that bracket the position.
	for i := 0; i < len(stops)-1; i++ {
		if position >= stops[i].Position && position <= stops[i+1].Position {
			// Linear interpolation between stops[i] and stops[i+1].
			span := stops[i+1].Position - stops[i].Position
			if span <= 0 {
				return stops[i].Color
			}
			t := (position - stops[i].Position) / span

			c0 := stops[i].Color
			c1 := stops[i+1].Color

			return color.RGBA{
				R: lerpU8(c0.R, c1.R, t),
				G: lerpU8(c0.G, c1.G, t),
				B: lerpU8(c0.B, c1.B, t),
				A: lerpU8(c0.A, c1.A, t),
			}
		}
	}

	// Fallback: return last stop color.
	return stops[len(stops)-1].Color
}

// lerpU8 linearly interpolates between two uint8 values at parameter t in [0, 1].
func lerpU8(a, b uint8, t float64) uint8 {
	return uint8(float64(a)*(1.0-t) + float64(b)*t + 0.5)
}

// computePerimeterGradient returns a color for each tile in perimeter order,
// or nil if fewer than 2 stops are provided.
// Stops are sorted by Position before interpolation.
// If more than 64 stops are provided, only the first 64 are used.
func computePerimeterGradient(cols, rows int, stops []gradient.ColorStop) []color.RGBA {
	if len(stops) < 2 {
		return nil
	}

	// Clamp to maximum stops.
	if len(stops) > maxGradientStops {
		stops = stops[:maxGradientStops]
	}

	// Sort stops by position.
	sorted := make([]gradient.ColorStop, len(stops))
	copy(sorted, stops)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Position < sorted[j].Position
	})

	total := perimeterTileCount(cols, rows)
	if total == 0 {
		return nil
	}

	colors := make([]color.RGBA, total)

	for i := 0; i < total; i++ {
		var pos float64
		if total == 1 {
			pos = 0.0
		} else {
			pos = float64(i) / float64(total-1)
		}
		colors[i] = interpolateColor(sorted, pos)
	}

	return colors
}
