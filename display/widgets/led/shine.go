package led

import (
	"image"
	"image/color"
	"math"
)

// applyShine draws a specular highlight on the LED body.
// Must be called LAST in the layer stack (after body fill).
// Caller is responsible for checking preconditions:
//   - cfg.ShineStyle != ShineNone
//   - effectiveBrightness > 0.0
//   - cfg.Diameter >= 5
func applyShine(img *image.RGBA, cfg Config, bodyRect image.Rectangle, effectiveBrightness float64) {
	// Resolve shine opacity: 0 means fully opaque (255).
	resolvedOpacity := cfg.ShineOpacity
	if resolvedOpacity == 0 {
		resolvedOpacity = 255
	}

	// Compute final alpha: opacity modulated by brightness.
	alpha := uint8(math.Floor(float64(resolvedOpacity) * effectiveBrightness))
	if alpha == 0 {
		return
	}

	shineColor := color.RGBA{R: 255, G: 255, B: 255, A: alpha}

	bodyRadius := bodyRect.Dx() / 2

	switch cfg.ShineStyle {
	case ShineDot:
		drawShineDot(img, bodyRect, bodyRadius, shineColor)
	case ShineCrescent:
		drawShineCrescent(img, bodyRect, bodyRadius, shineColor)
	}
}

// drawShineDot draws a small white filled circle in the upper-left quadrant.
// Dot radius = floor(bodyRadius * 0.15), minimum 1.
// Center offset = (-25%, -25%) from body center (upper-left).
func drawShineDot(img *image.RGBA, bodyRect image.Rectangle, bodyRadius int, shineColor color.RGBA) {
	dotRadius := int(math.Floor(float64(bodyRadius) * 0.15))
	if dotRadius < 1 {
		dotRadius = 1
	}

	// Body center in image coordinates.
	cx := bodyRect.Min.X + bodyRect.Dx()/2
	cy := bodyRect.Min.Y + bodyRect.Dy()/2

	// Dot center: upper-left offset.
	offsetX := int(math.Floor(float64(bodyRadius) * 0.25))
	offsetY := int(math.Floor(float64(bodyRadius) * 0.25))
	dotCx := cx - offsetX
	dotCy := cy - offsetY

	// Draw filled circle at (dotCx, dotCy) with dotRadius.
	r := float64(dotRadius)
	for py := dotCy - dotRadius; py <= dotCy+dotRadius; py++ {
		for px := dotCx - dotRadius; px <= dotCx+dotRadius; px++ {
			// Check bounds.
			if px < img.Rect.Min.X || px >= img.Rect.Max.X || py < img.Rect.Min.Y || py >= img.Rect.Max.Y {
				continue
			}
			dx := float64(px) + 0.5 - (float64(dotCx) + 0.5)
			dy := float64(py) + 0.5 - (float64(dotCy) + 0.5)
			if math.Sqrt(dx*dx+dy*dy) <= r {
				img.SetRGBA(px, py, shineColor)
			}
		}
	}
}

// drawShineCrescent draws a white arc-shaped highlight along the upper-left edge.
// The crescent occupies an annular region:
//   - innerRadius = bodyRadius * 0.80
//   - outerRadius = bodyRadius * 0.90
//
// A pixel is in the crescent if:
//   - distance from body center is in [innerRadius, outerRadius]
//   - angle from center is in [180°, 270°] (upper-left quadrant: dx < 0 and dy < 0)
func drawShineCrescent(img *image.RGBA, bodyRect image.Rectangle, bodyRadius int, shineColor color.RGBA) {
	innerRadius := float64(bodyRadius) * 0.80
	outerRadius := float64(bodyRadius) * 0.90

	// Ensure minimum thickness of 1 pixel.
	if outerRadius-innerRadius < 1.0 {
		outerRadius = innerRadius + 1.0
	}

	// Body center in image coordinates.
	cx := float64(bodyRect.Min.X) + float64(bodyRect.Dx())/2.0
	cy := float64(bodyRect.Min.Y) + float64(bodyRect.Dy())/2.0

	for py := bodyRect.Min.Y; py < bodyRect.Max.Y; py++ {
		for px := bodyRect.Min.X; px < bodyRect.Max.X; px++ {
			// Pixel center relative to body center.
			dx := float64(px) + 0.5 - cx
			dy := float64(py) + 0.5 - cy

			dist := math.Sqrt(dx*dx + dy*dy)

			// Check annular region.
			if dist < innerRadius || dist > outerRadius {
				continue
			}

			// Check angular region: upper-left quadrant (180° to 270°).
			// Using standard math convention: angle measured from positive X axis.
			// 180° to 270° means dx <= 0 and dy <= 0 (upper-left in screen coords where Y increases downward).
			if dx > 0 || dy > 0 {
				continue
			}

			// Bounds check.
			if px < img.Rect.Min.X || px >= img.Rect.Max.X || py < img.Rect.Min.Y || py >= img.Rect.Max.Y {
				continue
			}

			img.SetRGBA(px, py, shineColor)
		}
	}
}
