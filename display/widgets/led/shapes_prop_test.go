package led

import (
	"image"
	"image/color"
	"math"
	"testing"

	"pgregory.net/rapid"
)

// TestProperty1_ShapePixelMembershipCorrectness verifies that for any valid Config
// with Diameter ≥ 3 and LED in On state (no border, no glow), every pixel within the
// output image SHALL be either foreground-colored (if inside the shape's geometric
// boundary) or fully transparent (if outside), where:
//   - Circle: pixel center distance from image center ≤ Diameter/2
//   - Square: all pixels within the body area
//   - Diamond: |px - cx| + |py - cy| ≤ Body_Radius
//   - RoundedSquare: within the rounded rectangle with 25% corner radius

func TestProperty1_ShapePixelMembershipCorrectness(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate random valid config parameters
		diameter := rapid.IntRange(3, 64).Draw(t, "diameter")
		shape := Shape(rapid.IntRange(0, 3).Draw(t, "shape"))

		// Use a non-zero, non-transparent foreground color so we can distinguish inside from outside
		fg := color.RGBA{
			R: uint8(rapid.IntRange(1, 255).Draw(t, "fgR")),
			G: uint8(rapid.IntRange(1, 255).Draw(t, "fgG")),
			B: uint8(rapid.IntRange(1, 255).Draw(t, "fgB")),
			A: 255, // Full alpha to clearly distinguish from transparent
		}

		cfg := Config{
			Shape:      shape,
			State:      On,
			Brightness: -1.0, // Use discrete state → On → effectiveBrightness = 1.0
			Diameter:   diameter,
			Bounds:     image.Rect(0, 0, diameter, diameter),
			Foreground: fg,
			// No border, no glow
			BorderWidth: 0,
			GlowEnabled: false,
		}

		result := Render(cfg)
		if result == nil {
			t.Fatalf("expected non-nil result for Diameter=%d, Shape=%d", diameter, shape)
		}

		img := result.Image
		bounds := img.Bounds()

		// With no border and no glow, the body area is the full image
		bodyRect := image.Rect(0, 0, diameter, diameter)

		// Compute geometry parameters for membership checks
		w := float64(bodyRect.Dx())
		h := float64(bodyRect.Dy())
		cx := w / 2.0
		cy := h / 2.0
		radius := w / 2.0
		if h/2.0 < radius {
			radius = h / 2.0
		}

		// For RoundedSquare: corner radius = floor(side * 0.25)
		side := bodyRect.Dx()
		if bodyRect.Dy() < side {
			side = bodyRect.Dy()
		}
		cornerRadius := float64(side / 4)

		transparent := color.RGBA{R: 0, G: 0, B: 0, A: 0}

		for py := bounds.Min.Y; py < bounds.Max.Y; py++ {
			for px := bounds.Min.X; px < bounds.Max.X; px++ {
				// Local coordinates within the body area
				lx := px - bodyRect.Min.X
				ly := py - bodyRect.Min.Y

				// Pixel center
				pcx := float64(lx) + 0.5
				pcy := float64(ly) + 0.5

				// Determine geometric membership
				insideGeometry := isInsideShape(shape, pcx, pcy, cx, cy, radius, w, h, cornerRadius)

				// Get actual pixel color
				actualColor := img.(*image.RGBA).RGBAAt(px, py)

				if insideGeometry {
					// Pixel should be foreground-colored (at full brightness, so exactly fg)
					if actualColor != fg {
						t.Fatalf("pixel (%d,%d) inside %s geometry but got RGBA(%d,%d,%d,%d), expected RGBA(%d,%d,%d,%d) [diameter=%d]",
							px, py, shapeName(shape),
							actualColor.R, actualColor.G, actualColor.B, actualColor.A,
							fg.R, fg.G, fg.B, fg.A, diameter)
					}
				} else {
					// Pixel should be fully transparent
					if actualColor != transparent {
						t.Fatalf("pixel (%d,%d) outside %s geometry but got RGBA(%d,%d,%d,%d), expected transparent [diameter=%d]",
							px, py, shapeName(shape),
							actualColor.R, actualColor.G, actualColor.B, actualColor.A, diameter)
					}
				}
			}
		}
	})
}

// isInsideShape computes geometric membership for a pixel center at (pcx, pcy)
// within a body of dimensions (w, h), centered at (cx, cy).
func isInsideShape(shape Shape, pcx, pcy, cx, cy, radius, w, h, cornerRadius float64) bool {
	switch shape {
	case Circle:
		dx := pcx - cx
		dy := pcy - cy
		return math.Sqrt(dx*dx+dy*dy) <= radius

	case Square:
		// All pixels within the body area
		return pcx >= 0 && pcx <= w && pcy >= 0 && pcy <= h

	case Diamond:
		// Manhattan distance from center ≤ body radius
		return math.Abs(pcx-cx)+math.Abs(pcy-cy) <= radius

	case RoundedSquare:
		return isInsideRoundedRectGeometry(pcx, pcy, w, h, cornerRadius)

	default:
		// Fallback to circle (unrecognized shapes are mapped to Circle by validate)
		dx := pcx - cx
		dy := pcy - cy
		return math.Sqrt(dx*dx+dy*dy) <= radius
	}
}

// isInsideRoundedRectGeometry checks if (pcx, pcy) is inside a rounded rectangle
// of dimensions (w, h) with the given corner radius.
func isInsideRoundedRectGeometry(pcx, pcy, w, h, cr float64) bool {
	if pcx < 0 || pcx > w || pcy < 0 || pcy > h {
		return false
	}

	inLeftCol := pcx < cr
	inRightCol := pcx > w-cr
	inTopRow := pcy < cr
	inBottomRow := pcy > h-cr

	// Top-left corner
	if inLeftCol && inTopRow {
		dx := pcx - cr
		dy := pcy - cr
		return math.Sqrt(dx*dx+dy*dy) <= cr
	}
	// Top-right corner
	if inRightCol && inTopRow {
		dx := pcx - (w - cr)
		dy := pcy - cr
		return math.Sqrt(dx*dx+dy*dy) <= cr
	}
	// Bottom-left corner
	if inLeftCol && inBottomRow {
		dx := pcx - cr
		dy := pcy - (h - cr)
		return math.Sqrt(dx*dx+dy*dy) <= cr
	}
	// Bottom-right corner
	if inRightCol && inBottomRow {
		dx := pcx - (w - cr)
		dy := pcy - (h - cr)
		return math.Sqrt(dx*dx+dy*dy) <= cr
	}

	// In the central cross — always inside
	return true
}

// shapeName returns a human-readable name for a shape enum value.
func shapeName(s Shape) string {
	switch s {
	case Circle:
		return "Circle"
	case Square:
		return "Square"
	case Diamond:
		return "Diamond"
	case RoundedSquare:
		return "RoundedSquare"
	default:
		return "Unknown"
	}
}
