package borderframe

import (
	"image"
	"image/color"
	"math"
)

// clampGlowRadius clamps radius to valid range [0, 32].
func clampGlowRadius(radius int) int {
	if radius < 0 {
		return 0
	}
	if radius > 32 {
		return 32
	}
	return radius
}

// renderGlow creates a glow layer image for the given border mask.
// The returned image has the same bounds as the border image.
// Glow pixels extend up to effectiveRadius pixels from any border pixel (alpha > 0 in borderMask).
// Alpha at distance d from nearest border pixel: GlowColor.A * (1 - d/effectiveRadius) * pulseIntensity
// Returns nil if effectiveRadius is 0.
func renderGlow(borderMask *image.RGBA, bounds image.Rectangle, glowRadius int, glowColor color.RGBA, pulseIntensity float64) *image.RGBA {
	radius := clampGlowRadius(glowRadius)
	if radius == 0 {
		return nil
	}

	// Compute effective radius scaled by pulse intensity.
	effectiveRadius := float64(radius) * pulseIntensity
	if effectiveRadius <= 0 {
		return nil
	}

	maskBounds := borderMask.Bounds()
	glowImg := image.NewRGBA(maskBounds)

	// Collect border pixel coordinates (pixels with alpha > 0 in the mask).
	var borderPixels []image.Point
	for y := maskBounds.Min.Y; y < maskBounds.Max.Y; y++ {
		rowStart := (y - maskBounds.Min.Y) * borderMask.Stride
		for x := maskBounds.Min.X; x < maskBounds.Max.X; x++ {
			i := rowStart + (x-maskBounds.Min.X)*4
			if borderMask.Pix[i+3] > 0 {
				borderPixels = append(borderPixels, image.Point{X: x, Y: y})
			}
		}
	}

	if len(borderPixels) == 0 {
		return nil
	}

	// Compute clamped bounds for glow iteration (intersection of expanded area with bounds).
	radiusCeil := int(math.Ceil(effectiveRadius))
	glowMinX := maskBounds.Min.X
	glowMinY := maskBounds.Min.Y
	glowMaxX := maskBounds.Max.X
	glowMaxY := maskBounds.Max.Y

	// Clamp to Config.Bounds rectangle (no pixels outside the frame boundary).
	if bounds.Min.X > glowMinX {
		glowMinX = bounds.Min.X
	}
	if bounds.Min.Y > glowMinY {
		glowMinY = bounds.Min.Y
	}
	if bounds.Max.X < glowMaxX {
		glowMaxX = bounds.Max.X
	}
	if bounds.Max.Y < glowMaxY {
		glowMaxY = bounds.Max.Y
	}

	// Build a distance field: for each pixel in the glow area, find distance to nearest border pixel.
	// Use brute-force approach: iterate each border pixel and update surrounding pixels.
	// This is efficient for small frames and max radius of 32.
	distField := make([]float64, (glowMaxX-glowMinX)*(glowMaxY-glowMinY))
	for i := range distField {
		distField[i] = math.MaxFloat64
	}

	fieldW := glowMaxX - glowMinX

	for _, bp := range borderPixels {
		// Determine bounding box around this border pixel within the radius.
		minX := bp.X - radiusCeil
		maxX := bp.X + radiusCeil
		minY := bp.Y - radiusCeil
		maxY := bp.Y + radiusCeil

		// Clamp to glow bounds.
		if minX < glowMinX {
			minX = glowMinX
		}
		if minY < glowMinY {
			minY = glowMinY
		}
		if maxX >= glowMaxX {
			maxX = glowMaxX - 1
		}
		if maxY >= glowMaxY {
			maxY = glowMaxY - 1
		}

		for y := minY; y <= maxY; y++ {
			for x := minX; x <= maxX; x++ {
				dx := float64(x - bp.X)
				dy := float64(y - bp.Y)
				dist := math.Sqrt(dx*dx + dy*dy)
				idx := (y-glowMinY)*fieldW + (x - glowMinX)
				if dist < distField[idx] {
					distField[idx] = dist
				}
			}
		}
	}

	// Render glow pixels based on distance field.
	baseAlpha := float64(glowColor.A)

	for y := glowMinY; y < glowMaxY; y++ {
		for x := glowMinX; x < glowMaxX; x++ {
			idx := (y-glowMinY)*fieldW + (x - glowMinX)
			dist := distField[idx]

			// Skip border pixels themselves (distance 0 means it IS a border pixel).
			if dist == 0 {
				continue
			}

			// Skip pixels beyond effective radius.
			if dist >= effectiveRadius {
				continue
			}

			// Linear alpha falloff: alpha = GlowColor.A * (1 - d/effectiveRadius) * pulseIntensity
			alpha := baseAlpha * (1.0 - dist/effectiveRadius) * pulseIntensity
			if alpha <= 0 {
				continue
			}
			if alpha > 255 {
				alpha = 255
			}

			// Write glow pixel.
			pixIdx := (y-maskBounds.Min.Y)*glowImg.Stride + (x-maskBounds.Min.X)*4
			glowImg.Pix[pixIdx+0] = glowColor.R
			glowImg.Pix[pixIdx+1] = glowColor.G
			glowImg.Pix[pixIdx+2] = glowColor.B
			glowImg.Pix[pixIdx+3] = uint8(math.Round(alpha))
		}
	}

	return glowImg
}
