package led

import (
	"image"
	"image/color"
	"math"
)

// applyBorder draws a border ring of cfg.BorderWidth pixels around the LED body
// perimeter. The border follows the shape contour (circle, square, diamond, or
// rounded rectangle) and is drawn in the space between the body and the glow/outer edge.
//
// Border pixels are those that are:
//   - Inside the outer shape boundary (at full Diameter minus glow inset)
//   - Outside the inner body boundary (Diameter - 2×borderWidth)
//
// The border is at full configured border color regardless of LED state (no dimming
// for Off state — per Requirement 11.6).
//
// Layer ordering: border is drawn AFTER glow and BEFORE body fill.
func applyBorder(img *image.RGBA, cfg Config, glowRadius int) {
	if cfg.BorderWidth <= 0 {
		return
	}

	borderColor := cfg.BorderColor
	// Assert: zero-value border color should already be resolved to 50% gray by
	// resolveColors, but guard here as well.
	if borderColor == (color.RGBA{}) {
		borderColor = color.RGBA{R: 128, G: 128, B: 128, A: 255}
	}

	outputSize := img.Bounds().Dx()

	// The outer boundary of the border region starts at the glow inset.
	// outerInset = glowRadius (the border's outer edge aligns with the start of
	// the "Diameter" region within the output image).
	outerInset := glowRadius
	// The inner boundary of the border region is inset further by borderWidth.
	innerInset := glowRadius + cfg.BorderWidth

	// Outer shape bounding box (within which the border's outer edge is drawn)
	outerRect := image.Rect(outerInset, outerInset, outputSize-outerInset, outputSize-outerInset)
	// Inner shape bounding box (the body area — border doesn't extend into this)
	innerRect := image.Rect(innerInset, innerInset, outputSize-innerInset, outputSize-innerInset)

	switch cfg.Shape {
	case Square:
		applyBorderSquare(img, outerRect, innerRect, borderColor)
	case Diamond:
		applyBorderDiamond(img, outerRect, innerRect, borderColor)
	case RoundedSquare:
		applyBorderRoundedSquare(img, outerRect, innerRect, borderColor)
	default: // Circle
		applyBorderCircle(img, outerRect, innerRect, borderColor)
	}
}

// applyBorderCircle draws border pixels for a circular LED.
// A pixel is in the border if its distance from center is:
//   - <= outerRadius (inside the outer circle)
//   - > innerRadius (outside the inner/body circle)
func applyBorderCircle(img *image.RGBA, outerRect, innerRect image.Rectangle, borderColor color.RGBA) {
	outerW := outerRect.Dx()
	outerH := outerRect.Dy()
	outerCx := float64(outerW) / 2.0
	outerCy := float64(outerH) / 2.0
	outerRadius := float64(outerW) / 2.0
	if float64(outerH)/2.0 < outerRadius {
		outerRadius = float64(outerH) / 2.0
	}

	innerW := innerRect.Dx()
	innerH := innerRect.Dy()
	innerRadius := float64(innerW) / 2.0
	if float64(innerH)/2.0 < innerRadius {
		innerRadius = float64(innerH) / 2.0
	}

	// Iterate over the outer bounding box
	for py := outerRect.Min.Y; py < outerRect.Max.Y; py++ {
		for px := outerRect.Min.X; px < outerRect.Max.X; px++ {
			lx := px - outerRect.Min.X
			ly := py - outerRect.Min.Y

			dx := float64(lx) + 0.5 - outerCx
			dy := float64(ly) + 0.5 - outerCy
			dist := math.Sqrt(dx*dx + dy*dy)

			if dist <= outerRadius && dist > innerRadius {
				img.SetRGBA(px, py, borderColor)
			}
		}
	}
}

// applyBorderSquare draws border pixels for a square LED.
// A pixel is in the border if it's inside outerRect but outside innerRect.
func applyBorderSquare(img *image.RGBA, outerRect, innerRect image.Rectangle, borderColor color.RGBA) {
	for py := outerRect.Min.Y; py < outerRect.Max.Y; py++ {
		for px := outerRect.Min.X; px < outerRect.Max.X; px++ {
			// Inside outer but outside inner = border
			if !pointInRect(px, py, innerRect) {
				img.SetRGBA(px, py, borderColor)
			}
		}
	}
}

// applyBorderDiamond draws border pixels for a diamond-shaped LED.
// A pixel is in the border if its Manhattan distance from center is:
//   - <= outerRadius (inside the outer diamond)
//   - > innerRadius (outside the inner/body diamond)
func applyBorderDiamond(img *image.RGBA, outerRect, innerRect image.Rectangle, borderColor color.RGBA) {
	outerW := outerRect.Dx()
	outerH := outerRect.Dy()
	outerCx := float64(outerW) / 2.0
	outerCy := float64(outerH) / 2.0
	outerRadius := float64(outerW) / 2.0
	if float64(outerH)/2.0 < outerRadius {
		outerRadius = float64(outerH) / 2.0
	}

	innerW := innerRect.Dx()
	innerH := innerRect.Dy()
	innerRadius := float64(innerW) / 2.0
	if float64(innerH)/2.0 < innerRadius {
		innerRadius = float64(innerH) / 2.0
	}

	for py := outerRect.Min.Y; py < outerRect.Max.Y; py++ {
		for px := outerRect.Min.X; px < outerRect.Max.X; px++ {
			lx := px - outerRect.Min.X
			ly := py - outerRect.Min.Y

			pcx := float64(lx) + 0.5
			pcy := float64(ly) + 0.5

			// Manhattan distance from center
			dist := math.Abs(pcx-outerCx) + math.Abs(pcy-outerCy)

			if dist <= outerRadius && dist > innerRadius {
				img.SetRGBA(px, py, borderColor)
			}
		}
	}
}

// applyBorderRoundedSquare draws border pixels for a rounded-square LED.
// A pixel is in the border if it's inside the outer rounded rectangle but
// outside the inner rounded rectangle.
func applyBorderRoundedSquare(img *image.RGBA, outerRect, innerRect image.Rectangle, borderColor color.RGBA) {
	outerW := outerRect.Dx()
	outerH := outerRect.Dy()
	outerSide := outerW
	if outerH < outerSide {
		outerSide = outerH
	}
	outerCornerRadius := float64(outerSide / 4)

	innerW := innerRect.Dx()
	innerH := innerRect.Dy()
	innerSide := innerW
	if innerH < innerSide {
		innerSide = innerH
	}
	innerCornerRadius := float64(innerSide / 4)

	for py := outerRect.Min.Y; py < outerRect.Max.Y; py++ {
		for px := outerRect.Min.X; px < outerRect.Max.X; px++ {
			// Check if pixel is inside outer rounded rect
			outerLx := float64(px-outerRect.Min.X) + 0.5
			outerLy := float64(py-outerRect.Min.Y) + 0.5
			inOuter := isInsideRoundedRectGeneric(outerLx, outerLy, float64(outerW), float64(outerH), outerCornerRadius)

			if !inOuter {
				continue
			}

			// Check if pixel is inside inner rounded rect
			innerLx := float64(px-innerRect.Min.X) + 0.5
			innerLy := float64(py-innerRect.Min.Y) + 0.5
			inInner := false
			if px >= innerRect.Min.X && px < innerRect.Max.X && py >= innerRect.Min.Y && py < innerRect.Max.Y {
				inInner = isInsideRoundedRectGeneric(innerLx, innerLy, float64(innerW), float64(innerH), innerCornerRadius)
			}

			if !inInner {
				img.SetRGBA(px, py, borderColor)
			}
		}
	}
}

// isInsideRoundedRectGeneric checks if a point (px, py) in local coordinates is
// inside a rounded rectangle of dimensions w×h with the given corner radius.
func isInsideRoundedRectGeneric(px, py, w, h, cr float64) bool {
	if px < 0 || px > w || py < 0 || py > h {
		return false
	}

	inLeftCol := px < cr
	inRightCol := px > w-cr
	inTopRow := py < cr
	inBottomRow := py > h-cr

	if inLeftCol && inTopRow {
		dx := px - cr
		dy := py - cr
		return math.Sqrt(dx*dx+dy*dy) <= cr
	}
	if inRightCol && inTopRow {
		dx := px - (w - cr)
		dy := py - cr
		return math.Sqrt(dx*dx+dy*dy) <= cr
	}
	if inLeftCol && inBottomRow {
		dx := px - cr
		dy := py - (h - cr)
		return math.Sqrt(dx*dx+dy*dy) <= cr
	}
	if inRightCol && inBottomRow {
		dx := px - (w - cr)
		dy := py - (h - cr)
		return math.Sqrt(dx*dx+dy*dy) <= cr
	}

	return true
}

// pointInRect checks if a pixel coordinate is within the given rectangle.
func pointInRect(px, py int, r image.Rectangle) bool {
	return px >= r.Min.X && px < r.Max.X && py >= r.Min.Y && py < r.Max.Y
}
