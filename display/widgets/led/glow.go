package led

import (
	"image"
	"image/color"
	"math"
)

// applyGlow draws a radial glow region outside the LED body (or border outer edge)
// onto the output image. The glow is the first rendered layer after the transparent base.
//
// For each pixel outside the body/border but within glowRadius of the perimeter:
//
//	alpha = floor(glowBaseColor.A × (1 − distance / glowRadius) × effectiveBrightness)
//	color = glowBaseColor RGB (no brightness scaling on RGB — only alpha modulation)
//
// The glow base color is determined by:
//   - If gradient is configured with valid stops: use the outermost gradient stop color (highest position)
//   - Otherwise: use cfg.Foreground
func applyGlow(img *image.RGBA, cfg Config, glowRadius int, effectiveBrightness float64) {
	if glowRadius <= 0 || effectiveBrightness <= 0.0 {
		return
	}

	// Determine glow base color.
	glowBase := resolveGlowBaseColor(cfg)

	outputSize := img.Bounds().Dx()
	center := float64(outputSize) / 2.0

	// The shape boundary is at Diameter/2 from the center of the output image.
	// The output image is (Diameter + 2*glowRadius) × (Diameter + 2*glowRadius),
	// so the shape center is at (outputSize/2, outputSize/2).
	// The shape outer edge (including border) is at Diameter/2 from center.
	shapeRadius := float64(cfg.Diameter) / 2.0

	glowRadiusF := float64(glowRadius)

	for py := 0; py < outputSize; py++ {
		for px := 0; px < outputSize; px++ {
			// Pixel center coordinates
			pcx := float64(px) + 0.5
			pcy := float64(py) + 0.5

			// Compute distance from pixel to the shape perimeter
			dist := distanceFromPerimeter(pcx, pcy, center, center, shapeRadius, cfg.Shape, cfg.Diameter)

			// Only glow pixels: outside the shape (dist > 0) and within glow radius
			if dist > 0 && dist <= glowRadiusF {
				// Alpha = floor(glowBaseColor.A × (1 − dist / glowRadius) × effectiveBrightness)
				falloff := 1.0 - dist/glowRadiusF
				alpha := math.Floor(float64(glowBase.A) * falloff * effectiveBrightness)
				if alpha < 0 {
					alpha = 0
				}
				if alpha > 255 {
					alpha = 255
				}

				c := color.RGBA{
					R: glowBase.R,
					G: glowBase.G,
					B: glowBase.B,
					A: uint8(alpha),
				}
				img.SetRGBA(px, py, c)
			}
		}
	}
}

// resolveGlowBaseColor determines the glow base color from the config.
// If gradient is configured and has valid stops, uses the outermost stop (highest position).
// Otherwise uses cfg.Foreground.
func resolveGlowBaseColor(cfg Config) color.RGBA {
	if cfg.Gradient != nil && len(cfg.Gradient.Stops) >= 2 {
		// Find the stop with the highest position (outermost)
		outermost := cfg.Gradient.Stops[0]
		for _, stop := range cfg.Gradient.Stops[1:] {
			if stop.Position >= outermost.Position {
				outermost = stop
			}
		}
		return outermost.Color
	}
	return cfg.Foreground
}

// distanceFromPerimeter computes the signed distance from a pixel center (px, py) to
// the nearest point on the shape perimeter. Positive means outside the shape.
//
// The shape is centered at (cx, cy) with outer radius shapeRadius (= Diameter/2).
// This accounts for the full shape including border (border is inside Diameter).
func distanceFromPerimeter(px, py, cx, cy, shapeRadius float64, shape Shape, diameter int) float64 {
	switch shape {
	case Square:
		return distanceFromSquarePerimeter(px, py, cx, cy, shapeRadius)
	case Diamond:
		return distanceFromDiamondPerimeter(px, py, cx, cy, shapeRadius)
	case RoundedSquare:
		return distanceFromRoundedSquarePerimeter(px, py, cx, cy, shapeRadius, diameter)
	default: // Circle
		return distanceFromCirclePerimeter(px, py, cx, cy, shapeRadius)
	}
}

// distanceFromCirclePerimeter computes signed distance from point to circle edge.
// Positive = outside, negative = inside.
func distanceFromCirclePerimeter(px, py, cx, cy, radius float64) float64 {
	dx := px - cx
	dy := py - cy
	return math.Sqrt(dx*dx+dy*dy) - radius
}

// distanceFromSquarePerimeter computes signed distance from point to axis-aligned
// square edge using Chebyshev distance.
// The square has half-side = shapeRadius, centered at (cx, cy).
func distanceFromSquarePerimeter(px, py, cx, cy, halfSide float64) float64 {
	// Chebyshev distance from center minus half-side
	dx := math.Abs(px - cx)
	dy := math.Abs(py - cy)
	chebyshev := math.Max(dx, dy)
	return chebyshev - halfSide
}

// distanceFromDiamondPerimeter computes signed distance from point to diamond edge.
// The diamond has vertices at distance shapeRadius from center along axes.
// Manhattan distance from center - shapeRadius gives signed distance.
func distanceFromDiamondPerimeter(px, py, cx, cy, bodyRadius float64) float64 {
	dx := math.Abs(px - cx)
	dy := math.Abs(py - cy)
	// Manhattan distance
	manhattan := dx + dy
	return manhattan - bodyRadius
}

// distanceFromRoundedSquarePerimeter computes signed distance from point to rounded
// rectangle edge. Corner radius = 25% of body side length (Diameter).
func distanceFromRoundedSquarePerimeter(px, py, cx, cy, halfSide float64, diameter int) float64 {
	// Corner radius = 25% of diameter (the full side length), integer division floors
	cornerRadius := float64(diameter / 4)

	// Translate to first quadrant (exploit symmetry)
	dx := math.Abs(px - cx)
	dy := math.Abs(py - cy)

	// The rounded rectangle has half-extents of halfSide in both directions
	// with corner circles of radius cornerRadius at the corners.
	// Inner rectangle half-extents (where corners begin)
	innerHalfX := halfSide - cornerRadius
	innerHalfY := halfSide - cornerRadius

	if dx <= innerHalfX {
		// In the horizontal band - distance is just vertical distance to edge
		return dy - halfSide
	}
	if dy <= innerHalfY {
		// In the vertical band - distance is just horizontal distance to edge
		return dx - halfSide
	}

	// In the corner region - distance to the corner circle
	cornerDx := dx - innerHalfX
	cornerDy := dy - innerHalfY
	return math.Sqrt(cornerDx*cornerDx+cornerDy*cornerDy) - cornerRadius
}
