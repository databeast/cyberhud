package led

import (
	"image"
	"image/color"
	"math"
	"testing"

	"pgregory.net/rapid"
)

// TestProperty14_BorderOccupiesPerimeterBandAndInsetsBody verifies that for any
// valid Config with border width W in [1, 4] (after clamping to floor(Body_Radius/3)),
// the outermost W pixels following the shape perimeter SHALL be the border color,
// and the body fill area SHALL be inset by W pixels so that body + border fits
// within Diameter.

func TestProperty14_BorderOccupiesPerimeterBandAndInsetsBody(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Use a diameter large enough that border won't be fully clamped away
		diameter := rapid.IntRange(12, 64).Draw(t, "diameter")
		borderWidth := rapid.IntRange(1, 4).Draw(t, "borderWidth")

		// Clamp border width as validate does: min(borderWidth, floor(diameter/2 / 3))
		bodyRadius := diameter / 2
		maxBorder := bodyRadius / 3
		effectiveBorderWidth := borderWidth
		if effectiveBorderWidth > maxBorder {
			effectiveBorderWidth = maxBorder
		}
		if effectiveBorderWidth <= 0 {
			return // Skip: too small for any border
		}

		// Use distinct colors so border pixels are identifiable
		borderColor := color.RGBA{R: 255, G: 0, B: 0, A: 255}
		fg := color.RGBA{R: 0, G: 200, B: 0, A: 255}

		cfg := Config{
			Shape:       Circle,
			State:       On,
			Brightness:  -1.0,
			Diameter:    diameter,
			Bounds:      image.Rect(0, 0, diameter, diameter),
			Foreground:  fg,
			BorderWidth: borderWidth,
			BorderColor: borderColor,
			GlowEnabled: false,
		}

		result := Render(cfg)
		if result == nil {
			t.Fatal("expected non-nil sprite")
		}

		img := result.Image.(*image.RGBA)
		outputSize := img.Bounds().Dx()

		// With no glow, glowRadius = 0. The border occupies outerInset=0 to innerInset=effectiveBorderWidth.
		outerInset := 0
		innerInset := effectiveBorderWidth

		outerRadius := float64(outputSize-2*outerInset) / 2.0
		innerRadius := float64(outputSize-2*innerInset) / 2.0
		center := float64(outputSize) / 2.0

		// Verify: there exist border-colored pixels in the border band
		borderPixelCount := 0
		// Also verify: no body-fill pixels exist in the border band
		for py := 0; py < outputSize; py++ {
			for px := 0; px < outputSize; px++ {
				pcx := float64(px) + 0.5 - center
				pcy := float64(py) + 0.5 - center
				dist := math.Sqrt(pcx*pcx + pcy*pcy)

				c := img.RGBAAt(px, py)

				if dist <= outerRadius && dist > innerRadius {
					// This pixel should be in the border band
					if c == borderColor {
						borderPixelCount++
					}
				}
			}
		}

		if borderPixelCount == 0 {
			t.Fatalf("no border-colored pixels found in perimeter band [diameter=%d, borderWidth=%d, effectiveBorderWidth=%d]",
				diameter, borderWidth, effectiveBorderWidth)
		}

		// Verify body is inset: body pixels (fg-colored) should only appear inside innerRadius
		for py := 0; py < outputSize; py++ {
			for px := 0; px < outputSize; px++ {
				pcx := float64(px) + 0.5 - center
				pcy := float64(py) + 0.5 - center
				dist := math.Sqrt(pcx*pcx + pcy*pcy)

				c := img.RGBAAt(px, py)

				// If pixel is at full foreground color, it must be inside the inner radius
				if c == fg && dist > innerRadius+0.5 {
					t.Fatalf("body pixel at (%d,%d) outside inner boundary (dist=%.2f, innerRadius=%.2f) [diameter=%d, effectiveBorderWidth=%d]",
						px, py, dist, innerRadius, diameter, effectiveBorderWidth)
				}
			}
		}
	})
}

// TestProperty15_BorderSkippedWhenWidthIsZero verifies that for any Config with
// border width = 0, the output SHALL contain no border-colored pixels.

func TestProperty15_BorderSkippedWhenWidthIsZero(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		diameter := rapid.IntRange(3, 64).Draw(t, "diameter")
		shape := Shape(rapid.IntRange(0, 3).Draw(t, "shape"))

		// Use a distinct border color that won't appear naturally
		borderColor := color.RGBA{R: 128, G: 0, B: 128, A: 255}
		// Use foreground that is different from borderColor
		fg := color.RGBA{R: 0, G: 200, B: 0, A: 255}

		cfg := Config{
			Shape:       shape,
			State:       On,
			Brightness:  -1.0,
			Diameter:    diameter,
			Bounds:      image.Rect(0, 0, diameter, diameter),
			Foreground:  fg,
			BorderWidth: 0, // No border
			BorderColor: borderColor,
			GlowEnabled: false,
		}

		result := Render(cfg)
		if result == nil {
			t.Fatal("expected non-nil sprite")
		}

		img := result.Image.(*image.RGBA)
		bounds := img.Bounds()

		for py := bounds.Min.Y; py < bounds.Max.Y; py++ {
			for px := bounds.Min.X; px < bounds.Max.X; px++ {
				c := img.RGBAAt(px, py)
				if c == borderColor {
					t.Fatalf("found border-colored pixel at (%d,%d) with borderWidth=0 [shape=%s, diameter=%d]",
						px, py, shapeName(shape), diameter)
				}
			}
		}
	})
}

// TestProperty16_BorderWidthClampingIdempotence verifies that for any Config with
// border width exceeding floor(Body_Radius/3), the rendered output SHALL be
// pixel-identical to the same Config with border width set to floor(Body_Radius/3).

func TestProperty16_BorderWidthClampingIdempotence(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Use a small diameter so that floor(bodyRadius/3) < 4
		diameter := rapid.IntRange(6, 20).Draw(t, "diameter")
		shape := Shape(rapid.IntRange(0, 3).Draw(t, "shape"))

		bodyRadius := diameter / 2
		maxBorder := bodyRadius / 3

		if maxBorder < 1 {
			return // Skip: too small for any border
		}

		// Use a width that exceeds the max
		excessiveWidth := rapid.IntRange(maxBorder+1, 4).Draw(t, "excessiveWidth")
		if excessiveWidth <= maxBorder {
			return // Skip: random range didn't produce a valid excess
		}

		borderColor := color.RGBA{R: 200, G: 100, B: 50, A: 255}
		fg := color.RGBA{R: 0, G: 200, B: 0, A: 255}

		// Config with excessive border width
		cfgExcessive := Config{
			Shape:       shape,
			State:       On,
			Brightness:  -1.0,
			Diameter:    diameter,
			Bounds:      image.Rect(0, 0, diameter, diameter),
			Foreground:  fg,
			BorderWidth: excessiveWidth,
			BorderColor: borderColor,
			GlowEnabled: false,
		}

		// Config with clamped border width
		cfgClamped := Config{
			Shape:       shape,
			State:       On,
			Brightness:  -1.0,
			Diameter:    diameter,
			Bounds:      image.Rect(0, 0, diameter, diameter),
			Foreground:  fg,
			BorderWidth: maxBorder,
			BorderColor: borderColor,
			GlowEnabled: false,
		}

		resultExcessive := Render(cfgExcessive)
		resultClamped := Render(cfgClamped)

		if resultExcessive == nil || resultClamped == nil {
			t.Fatal("expected non-nil sprites")
		}

		imgExcessive := resultExcessive.Image.(*image.RGBA)
		imgClamped := resultClamped.Image.(*image.RGBA)

		boundsE := imgExcessive.Bounds()
		boundsC := imgClamped.Bounds()

		if boundsE.Dx() != boundsC.Dx() || boundsE.Dy() != boundsC.Dy() {
			t.Fatalf("image dimensions differ: excessive=%dx%d, clamped=%dx%d",
				boundsE.Dx(), boundsE.Dy(), boundsC.Dx(), boundsC.Dy())
		}

		for py := boundsE.Min.Y; py < boundsE.Max.Y; py++ {
			for px := boundsE.Min.X; px < boundsE.Max.X; px++ {
				ce := imgExcessive.RGBAAt(px, py)
				cc := imgClamped.RGBAAt(px, py)
				if ce != cc {
					t.Fatalf("pixel mismatch at (%d,%d): excessive=RGBA(%d,%d,%d,%d), clamped=RGBA(%d,%d,%d,%d) [diameter=%d, excessiveWidth=%d, maxBorder=%d]",
						px, py, ce.R, ce.G, ce.B, ce.A, cc.R, cc.G, cc.B, cc.A,
						diameter, excessiveWidth, maxBorder)
				}
			}
		}
	})
}

// TestProperty9_GlowOpacityLinearFalloffAndBrightnessModulation verifies that
// for any valid Config with glow enabled and effective brightness B > 0, every
// glow pixel at distance d from the body edge SHALL have alpha equal to
// floor(GlowColor.A × (1 − d / glowRadius) × B).

func TestProperty9_GlowOpacityLinearFalloffAndBrightnessModulation(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		diameter := rapid.IntRange(10, 40).Draw(t, "diameter")
		glowRadius := rapid.IntRange(3, 16).Draw(t, "glowRadius")
		brightness := rapid.Float64Range(0.1, 1.0).Draw(t, "brightness")

		// Use foreground with A=255 for predictable glow alpha calculations
		fg := color.RGBA{
			R: uint8(rapid.IntRange(1, 255).Draw(t, "fgR")),
			G: uint8(rapid.IntRange(1, 255).Draw(t, "fgG")),
			B: uint8(rapid.IntRange(1, 255).Draw(t, "fgB")),
			A: 255,
		}

		cfg := Config{
			Shape:       Circle,
			State:       On,
			Brightness:  brightness,
			Diameter:    diameter,
			Bounds:      image.Rect(0, 0, diameter, diameter),
			Foreground:  fg,
			GlowEnabled: true,
			GlowRadius:  glowRadius,
			BorderWidth: 0,
		}

		result := Render(cfg)
		if result == nil {
			t.Fatal("expected non-nil sprite")
		}

		img := result.Image.(*image.RGBA)
		outputSize := img.Bounds().Dx()
		center := float64(outputSize) / 2.0
		shapeRadius := float64(diameter) / 2.0
		glowRadiusF := float64(glowRadius)

		// Sample glow pixels along the horizontal axis to the right of the shape
		// These are on the glow region for a circle (distance from perimeter is clear)
		verified := 0
		for px := 0; px < outputSize; px++ {
			py := outputSize / 2 // Horizontal center line

			pcx := float64(px) + 0.5
			pcy := float64(py) + 0.5

			dx := pcx - center
			dy := pcy - center
			dist := math.Sqrt(dx*dx+dy*dy) - shapeRadius

			// Only check glow pixels (outside body, within glow radius)
			if dist > 0.5 && dist < glowRadiusF-0.5 {
				c := img.RGBAAt(px, py)
				if c.A == 0 {
					continue // Pixel might be just barely outside, skip
				}

				// Expected alpha = floor(glowBase.A × (1 − dist / glowRadius) × brightness)
				falloff := 1.0 - dist/glowRadiusF
				expectedAlpha := math.Floor(float64(fg.A) * falloff * brightness)
				if expectedAlpha < 0 {
					expectedAlpha = 0
				}
				if expectedAlpha > 255 {
					expectedAlpha = 255
				}

				// Allow ±1 tolerance for floating point rounding
				diff := math.Abs(float64(c.A) - expectedAlpha)
				if diff > 1.0 {
					t.Fatalf("glow pixel at (%d,%d) has alpha=%d, expected=%d (dist=%.2f, falloff=%.4f, brightness=%.4f) [diameter=%d, glowRadius=%d]",
						px, py, c.A, int(expectedAlpha), dist, falloff, brightness, diameter, glowRadius)
				}

				// Verify RGB matches glow base color
				if c.R != fg.R || c.G != fg.G || c.B != fg.B {
					t.Fatalf("glow pixel at (%d,%d) has RGB=(%d,%d,%d), expected=(%d,%d,%d)",
						px, py, c.R, c.G, c.B, fg.R, fg.G, fg.B)
				}

				verified++
			}
		}

		if verified == 0 {
			t.Fatalf("no glow pixels verified on horizontal center line [diameter=%d, glowRadius=%d]",
				diameter, glowRadius)
		}
	})
}

// TestProperty9_GlowSuppressedWhenBrightnessZero verifies that when B = 0,
// no glow pixels SHALL exist.

func TestProperty9_GlowSuppressedWhenBrightnessZero(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		diameter := rapid.IntRange(6, 40).Draw(t, "diameter")
		glowRadius := rapid.IntRange(2, 16).Draw(t, "glowRadius")

		fg := color.RGBA{R: 0, G: 200, B: 0, A: 255}

		cfg := Config{
			Shape:       Circle,
			State:       Off,
			Brightness:  0.0,
			Diameter:    diameter,
			Bounds:      image.Rect(0, 0, diameter, diameter),
			Foreground:  fg,
			GlowEnabled: true,
			GlowRadius:  glowRadius,
			BorderWidth: 0,
		}

		result := Render(cfg)
		if result == nil {
			t.Fatal("expected non-nil sprite")
		}

		img := result.Image.(*image.RGBA)
		outputSize := img.Bounds().Dx()
		center := float64(outputSize) / 2.0
		shapeRadius := float64(diameter) / 2.0

		// Check that no pixels in the glow region have non-zero alpha matching glow color
		for py := 0; py < outputSize; py++ {
			for px := 0; px < outputSize; px++ {
				pcx := float64(px) + 0.5
				pcy := float64(py) + 0.5
				dx := pcx - center
				dy := pcy - center
				dist := math.Sqrt(dx*dx + dy*dy)

				// Only check glow region (outside the body shape)
				if dist > shapeRadius+0.5 {
					c := img.RGBAAt(px, py)
					if c.A != 0 {
						t.Fatalf("found glow pixel at (%d,%d) with alpha=%d when brightness=0 [diameter=%d, glowRadius=%d]",
							px, py, c.A, diameter, glowRadius)
					}
				}
			}
		}
	})
}

// TestProperty10_GlowContourFollowsShapePerimeter verifies that for any valid
// Config with glow enabled and a non-Circle shape, glow pixel distance SHALL be
// computed from the nearest point on the body perimeter (not from center),
// producing a glow region that conforms to the shape contour.
//
// For a square shape, two glow pixels that have the same perimeter distance
// (as computed by distanceFromSquarePerimeter) must have the same alpha.
// We verify this by sampling pixels at matching Chebyshev distances on different
// sides of the square.

func TestProperty10_GlowContourFollowsShapePerimeter(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		diameter := rapid.IntRange(12, 40).Draw(t, "diameter")
		glowRadius := rapid.IntRange(4, 12).Draw(t, "glowRadius")

		// Test with Square shape - glow should follow square perimeter
		fg := color.RGBA{R: 200, G: 100, B: 50, A: 255}

		cfg := Config{
			Shape:       Square,
			State:       On,
			Brightness:  1.0,
			Diameter:    diameter,
			Bounds:      image.Rect(0, 0, diameter, diameter),
			Foreground:  fg,
			GlowEnabled: true,
			GlowRadius:  glowRadius,
			BorderWidth: 0,
		}

		result := Render(cfg)
		if result == nil {
			t.Fatal("expected non-nil sprite")
		}

		img := result.Image.(*image.RGBA)
		outputSize := img.Bounds().Dx()
		center := float64(outputSize) / 2.0
		shapeRadius := float64(diameter) / 2.0

		// Verify that glow pixels use shape-perimeter distance by checking that
		// pixels at the same computed distance from the square perimeter have
		// the same alpha value.
		//
		// For a square with Chebyshev distance: dist = max(|dx|,|dy|) - halfSide
		// A pixel on the top midline at (center, y) has dist = |y - center| - halfSide
		// A pixel on the right midline at (x, center) has dist = |x - center| - halfSide
		// At the same absolute distance from edge, these must match.

		// Sample glow pixels along the top midline (directly above center)
		// and along the right midline (directly right of center)
		verified := 0
		for offset := 1; offset < glowRadius; offset++ {
			// Pixel on the top midline: directly above the shape center
			topPx := int(center)
			topPy := int(center-shapeRadius) - offset
			if topPy < 0 || topPx >= outputSize {
				continue
			}

			// Pixel on the right midline: directly right of the shape center
			rightPx := int(center+shapeRadius) + offset - 1
			rightPy := int(center)
			if rightPx >= outputSize || rightPy >= outputSize {
				continue
			}

			// Both pixels should have the same distance from the square perimeter
			topDist := distanceFromPerimeter(float64(topPx)+0.5, float64(topPy)+0.5, center, center, shapeRadius, Square, diameter)
			rightDist := distanceFromPerimeter(float64(rightPx)+0.5, float64(rightPy)+0.5, center, center, shapeRadius, Square, diameter)

			// Only compare if distances are close (same integer offset should be close)
			if math.Abs(topDist-rightDist) > 0.5 {
				continue
			}

			topColor := img.RGBAAt(topPx, topPy)
			rightColor := img.RGBAAt(rightPx, rightPy)

			// Both glow pixels at the same distance should have the same alpha
			if topColor.A > 0 && rightColor.A > 0 {
				if absDiff(int(topColor.A), int(rightColor.A)) > 1 {
					t.Fatalf("glow alpha mismatch at equidistant points (offset=%d, topDist=%.2f, rightDist=%.2f): top(%d,%d).A=%d, right(%d,%d).A=%d [diameter=%d, glowRadius=%d]",
						offset, topDist, rightDist, topPx, topPy, topColor.A, rightPx, rightPy, rightColor.A, diameter, glowRadius)
				}
				verified++
			}
		}

		if verified == 0 {
			t.Fatalf("no equidistant glow pixel pairs verified [diameter=%d, glowRadius=%d]",
				diameter, glowRadius)
		}
	})
}

// absDiff returns the absolute difference between two integers.
func absDiff(a, b int) int {
	if a > b {
		return a - b
	}
	return b - a
}
