package progressbar

import (
	"image"
	"image/color"
)

// applyBorder draws a border outline around the bar perimeter. It overwrites
// the outermost borderWidth pixels on each edge with the configured BorderColor.
//
// For Linear style with RoundedCaps enabled, the border follows the stadium
// (capsule) curvature. For all other cases (Linear without caps, Segmented),
// the border is rectangular.
//
// Border is skipped when:
//   - BorderWidth <= 0
//   - BorderColor is the zero value (fully transparent black)
//   - Minor axis < 2*borderWidth + 2 (bar too thin for meaningful interior)
//
// For Segmented bars, the border wraps the entire bar perimeter (encompassing
// all cells and gaps), not individual cells.
func applyBorder(img *image.RGBA, cfg Config) {
	bw := effectiveBorderWidth(cfg)
	if bw <= 0 {
		return
	}
	if cfg.BorderColor == (color.RGBA{}) {
		return
	}

	w := img.Bounds().Dx()
	h := img.Bounds().Dy()

	// Determine minor axis based on orientation (same logic for Linear and Segmented).
	minorAxis := h
	if cfg.Orientation == OrientVertical {
		minorAxis = w
	}

	// Skip if too thin for border + at least 2 interior pixels.
	if minorAxis < 2*bw+2 {
		return
	}

	bc := cfg.BorderColor

	// Only Linear style with RoundedCaps gets capsule-shaped border.
	if cfg.Style == Linear && cfg.RoundedCaps && minorAxis >= 2 {
		drawCapsuleBorder(img, cfg, bw, bc)
	} else {
		drawRectBorder(img, w, h, bw, bc)
	}
}

func effectiveBorderWidth(cfg Config) int {
	return cfg.BorderWidth + cfg.BorderWall
}

// drawRectBorder draws a rectangular border of width bw around the image perimeter.
// Top band: rows [0, bw)
// Bottom band: rows [h-bw, h)
// Left band: rows [bw, h-bw), cols [0, bw)
// Right band: rows [bw, h-bw), cols [w-bw, w)
func drawRectBorder(img *image.RGBA, w, h, bw int, bc color.RGBA) {
	// Top band
	for y := 0; y < bw; y++ {
		for x := 0; x < w; x++ {
			img.SetRGBA(x, y, bc)
		}
	}
	// Bottom band
	for y := h - bw; y < h; y++ {
		for x := 0; x < w; x++ {
			img.SetRGBA(x, y, bc)
		}
	}
	// Left band (between top and bottom bands)
	for y := bw; y < h-bw; y++ {
		for x := 0; x < bw; x++ {
			img.SetRGBA(x, y, bc)
		}
	}
	// Right band (between top and bottom bands)
	for y := bw; y < h-bw; y++ {
		for x := w - bw; x < w; x++ {
			img.SetRGBA(x, y, bc)
		}
	}
}

// drawCapsuleBorder draws a border following the stadium/capsule shape for
// Linear bars with rounded caps enabled. The border is bw pixels thick and
// follows the semicircular curvature at both ends.
//
// A pixel is in the border if:
//  1. It is inside the capsule shape (outer boundary), AND
//  2. It is NOT inside the inset capsule (shrunk by bw on all sides)
//
// This produces a uniform-width border that follows the curved ends.
func drawCapsuleBorder(img *image.RGBA, cfg Config, bw int, bc color.RGBA) {
	w := img.Bounds().Dx()
	h := img.Bounds().Dy()

	// Outer capsule parameters
	var outerMinor int
	if cfg.Orientation == OrientVertical {
		outerMinor = w
	} else {
		outerMinor = h
	}
	outerRadius := float64(outerMinor) / 2.0
	outerR2 := outerRadius * outerRadius

	// Inner capsule parameters (inset by bw on each side)
	innerRadius := outerRadius - float64(bw)
	innerR2 := innerRadius * innerRadius

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if !insideCapsule(x, y, w, h, outerRadius, outerR2, cfg.Orientation) {
				// Outside the capsule entirely — skip (already transparent from caps pass)
				continue
			}
			if innerRadius > 0 && insideInsetCapsule(x, y, w, h, bw, innerRadius, innerR2, cfg.Orientation) {
				// Inside the inner capsule — this is interior, not border
				continue
			}
			// Between outer and inner capsule — this is the border band
			img.SetRGBA(x, y, bc)
		}
	}
}

// insideInsetCapsule returns true if the pixel at (x, y) is inside the inset
// capsule (the capsule shrunk by bw pixels on all sides).
func insideInsetCapsule(x, y, w, h, bw int, innerRadius, innerR2 float64, orient Orientation) bool {
	px := float64(x) + 0.5
	py := float64(y) + 0.5

	bwf := float64(bw)

	if orient == OrientVertical {
		// Vertical capsule: width is minor axis.
		// Outer capsule: radius = w/2, centers at (w/2, w/2) and (w/2, h - w/2)
		// Inner capsule: radius = w/2 - bw, centers at (w/2, w/2) and (w/2, h - w/2)
		// But the inner capsule is inset, so the straight sides are at x in [bw, w-bw].
		outerRadius := float64(w) / 2.0

		// Check x is within the inset straight sides
		if px < bwf || px > float64(w)-bwf {
			return false
		}

		centerX := float64(w) / 2.0

		if py < outerRadius {
			// In the top cap region of outer capsule
			// Inner cap center is at (centerX, outerRadius)
			dx := px - centerX
			dy := py - outerRadius
			return dx*dx+dy*dy <= innerR2
		} else if py > float64(h)-outerRadius {
			// In the bottom cap region of outer capsule
			// Inner cap center is at (centerX, h - outerRadius)
			dx := px - centerX
			dy := py - (float64(h) - outerRadius)
			return dx*dx+dy*dy <= innerR2
		}
		// In the middle straight section — check y is within inset
		return py >= bwf && py <= float64(h)-bwf
	}

	// Horizontal capsule: height is minor axis.
	// Outer capsule: radius = h/2, centers at (h/2, h/2) and (w - h/2, h/2)
	// Inner capsule: inset by bw on all sides.
	outerRadius := float64(h) / 2.0

	// Check y is within the inset straight sides
	if py < bwf || py > float64(h)-bwf {
		return false
	}

	centerY := float64(h) / 2.0

	if px < outerRadius {
		// In the left cap region of outer capsule
		// Inner cap center is at (outerRadius, centerY)
		dx := px - outerRadius
		dy := py - centerY
		return dx*dx+dy*dy <= innerR2
	} else if px > float64(w)-outerRadius {
		// In the right cap region of outer capsule
		// Inner cap center is at (w - outerRadius, centerY)
		dx := px - (float64(w) - outerRadius)
		dy := py - centerY
		return dx*dx+dy*dy <= innerR2
	}
	// In the middle straight section — check y is within inset
	return true
}
