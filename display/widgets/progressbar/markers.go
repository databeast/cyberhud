package progressbar

import (
	"image"
	"image/color"
	"math"
)

// clampMarkerValue restricts v to the range [0.0, 1.0].
func clampMarkerValue(v float64) float64 {
	if v < 0.0 {
		return 0.0
	}
	if v > 1.0 {
		return 1.0
	}
	return v
}

// drawMarkers draws threshold marker lines on top of the rendered bar image.
// Markers are 1-pixel-wide lines perpendicular to the bar axis at each marker's
// value position. This function is called LAST in the render pipeline (after
// animation) so markers remain visible over all other layers.
func drawMarkers(img *image.RGBA, cfg Config) {
	if len(cfg.Markers) == 0 {
		return
	}

	w := img.Bounds().Dx()
	h := img.Bounds().Dy()

	for _, m := range cfg.Markers {
		val := clampMarkerValue(m.Value)
		c := m.Color
		if c == (color.RGBA{}) {
			c = color.RGBA{255, 255, 255, 255}
		}

		switch cfg.Style {
		case Linear, Segmented:
			if cfg.Orientation == OrientVertical {
				// Vertical: draw horizontal line at y position spanning full width.
				y := h - 1 - int(val*float64(h-1))
				for x := 0; x < w; x++ {
					img.SetRGBA(x, y, c)
				}
			} else {
				// Horizontal: draw vertical line at x position spanning full height.
				x := int(val * float64(w-1))
				for y := 0; y < h; y++ {
					img.SetRGBA(x, y, c)
				}
			}
		case Ring, Pie:
			drawRadialMarker(img, cfg, val, c)
		case Arc:
			drawArcRadialMarker(img, cfg, val, c)
		}
	}
}

// drawRadialMarker draws a radial tick line for Ring/Pie style at the angular
// position determined by the marker value. The line extends from the inner edge
// to the outer edge of the annulus (Ring) or from center to outer edge (Pie).
func drawRadialMarker(img *image.RGBA, cfg Config, val float64, c color.RGBA) {
	w := img.Bounds().Dx()
	h := img.Bounds().Dy()

	cx := float64(w) / 2.0
	cy := float64(h) / 2.0

	outerR := math.Min(float64(w), float64(h)) / 2.0

	var innerR float64
	if cfg.Style == Ring {
		innerR = outerR - float64(cfg.Thickness)
		if innerR < 0 {
			innerR = 0
		}
	}
	// For Pie, innerR stays 0 (line from center to outer edge).

	// Compute the start offset (same logic as renderRing / renderPie).
	// geom.StartAngle for Ring/Pie:
	//   Horizontal → -π/2 (12-o'clock)
	//   Vertical   → π   (9-o'clock)
	// Convert to clockwise offset from 12-o'clock:
	//   offset = startAngle + π/2, normalized to [0, 2π)
	var startAngle float64
	if cfg.Orientation == OrientVertical {
		startAngle = math.Pi
	} else {
		startAngle = -math.Pi / 2.0
	}
	startOffset := startAngle + math.Pi/2.0
	startOffset = math.Mod(startOffset, 2.0*math.Pi)
	if startOffset < 0 {
		startOffset += 2.0 * math.Pi
	}

	// Angular position of the marker: val * full circle.
	markerAngle := startOffset + val*2.0*math.Pi
	// Normalize to [0, 2π)
	markerAngle = math.Mod(markerAngle, 2.0*math.Pi)
	if markerAngle < 0 {
		markerAngle += 2.0 * math.Pi
	}

	// Convert clockwise-from-north angle to standard math angle for sin/cos.
	// Clockwise from north: 0 = up (negative Y), angle increases clockwise.
	// dx = sin(angle), dy = -cos(angle)
	sinA := math.Sin(markerAngle)
	cosA := math.Cos(markerAngle)

	// Draw the radial line from innerR to outerR using Bresenham-like stepping.
	steps := int(outerR - innerR)
	if steps < 1 {
		steps = 1
	}
	for i := 0; i <= steps; i++ {
		r := innerR + float64(i)*(outerR-innerR)/float64(steps)
		px := int(cx + r*sinA)
		py := int(cy - r*cosA)
		if px >= 0 && px < w && py >= 0 && py < h {
			img.SetRGBA(px, py, c)
		}
	}
}

// drawArcRadialMarker draws a radial tick line for Arc style at the angular
// position determined by the marker value within the arc's sweep.
func drawArcRadialMarker(img *image.RGBA, cfg Config, val float64, c color.RGBA) {
	w := img.Bounds().Dx()
	h := img.Bounds().Dy()

	cx := float64(w) / 2.0
	cy := float64(h) / 2.0

	outerR := math.Min(float64(w), float64(h)) / 2.0
	innerR := outerR - float64(cfg.Thickness)
	if innerR < 0 {
		innerR = 0
	}

	// Sweep angle in radians.
	sweepRad := cfg.SweepAngle * math.Pi / 180.0

	// geom.StartAngle for Arc:
	//   Horizontal → -π/2 (12-o'clock)
	//   Vertical   → π   (9-o'clock)
	var geomStartAngle float64
	if cfg.Orientation == OrientVertical {
		geomStartAngle = math.Pi
	} else {
		geomStartAngle = -math.Pi / 2.0
	}

	// arcStart is the beginning of the arc (the "left endpoint"),
	// centered about the orientation reference angle.
	arcStart := geomStartAngle - sweepRad/2.0

	// The marker is at val fraction along the arc sweep.
	// markerAngle is in the same coordinate system as atan2(dx, -dy) [-π, π].
	markerStdAngle := arcStart + val*sweepRad

	// Convert standard math angle to screen coordinates for drawing.
	// In our coordinate system: angle measured from -Y axis (12-o'clock),
	// positive clockwise. We use atan2(dx, -dy).
	// To draw at angle θ (in the atan2(dx,-dy) system):
	//   dx = sin(θ), dy = -cos(θ)  (relative to center)
	// But our arcStart is already in this system (geom.StartAngle uses atan2(dx,-dy) convention).
	sinA := math.Sin(markerStdAngle)
	cosA := math.Cos(markerStdAngle)

	// dx_dir = sin(angle), dy_dir = -cos(angle) gives direction from center
	// in the atan2(dx, -dy) system.
	dxDir := sinA
	dyDir := -cosA

	// Draw the radial line from innerR to outerR.
	steps := int(outerR - innerR)
	if steps < 1 {
		steps = 1
	}
	for i := 0; i <= steps; i++ {
		r := innerR + float64(i)*(outerR-innerR)/float64(steps)
		px := int(cx + r*dxDir)
		py := int(cy + r*dyDir)
		if px >= 0 && px < w && py >= 0 && py < h {
			img.SetRGBA(px, py, c)
		}
	}
}
