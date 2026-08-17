package progressbar

import (
	"image"
	"image/color"
)

// applyRoundedCaps masks pixels outside the "stadium" (capsule) shape to the
// background color, giving the bar semicircular ends. This is a post-processing
// step applied ONLY to Linear style bars when RoundedCaps is true and the minor
// axis is at least 2 pixels.
//
// For Horizontal orientation the capsule has:
//   - Left semicircle centered at (radius, radius) with radius = height/2
//   - Right semicircle centered at (width-radius, radius) with radius = height/2
//   - Rectangle between the semicircles spanning full height
//
// For Vertical orientation the capsule has:
//   - Top semicircle centered at (radius, radius) with radius = width/2
//   - Bottom semicircle centered at (radius, height-radius) with radius = width/2
//   - Rectangle between the semicircles spanning full width
//
// Pixels outside the capsule are set to transparent (zero RGBA).
func applyRoundedCaps(img *image.RGBA, cfg Config) {
	w := img.Bounds().Dx()
	h := img.Bounds().Dy()

	// Determine minor axis and radius based on orientation.
	var minorAxis int
	if cfg.Orientation == OrientVertical {
		minorAxis = w
	} else {
		minorAxis = h
	}

	// Skip if minor axis < 2 (render rectangular).
	if minorAxis < 2 {
		return
	}

	radius := float64(minorAxis) / 2.0
	bg := color.RGBA{0, 0, 0, 0} // transparent for masked pixels

	// Use squared distance to avoid sqrt for performance.
	r2 := radius * radius

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if !insideCapsule(x, y, w, h, radius, r2, cfg.Orientation) {
				img.SetRGBA(x, y, bg)
			}
		}
	}
}

// insideCapsule returns true if the pixel at (x, y) is inside the stadium shape.
func insideCapsule(x, y, w, h int, radius, r2 float64, orient Orientation) bool {
	// Use pixel center (x+0.5, y+0.5) for smoother edges.
	px := float64(x) + 0.5
	py := float64(y) + 0.5

	if orient == OrientVertical {
		// Vertical: capsule runs top-to-bottom, semicircles at top and bottom.
		// radius = width / 2.0
		// Top semicircle center: (radius, radius)
		// Bottom semicircle center: (radius, height - radius)
		// Middle rectangle: y in [radius, height-radius], x in [0, width]

		if py < radius {
			// In the top cap region — check distance from top center.
			dx := px - radius
			dy := py - radius
			return dx*dx+dy*dy <= r2
		} else if py > float64(h)-radius {
			// In the bottom cap region — check distance from bottom center.
			dx := px - radius
			dy := py - (float64(h) - radius)
			return dx*dx+dy*dy <= r2
		}
		// In the middle rectangle — always inside.
		return true
	}

	// Horizontal: capsule runs left-to-right, semicircles at left and right.
	// radius = height / 2.0
	// Left semicircle center: (radius, radius)
	// Right semicircle center: (width - radius, radius)
	// Middle rectangle: x in [radius, width-radius], y in [0, height]

	if px < radius {
		// In the left cap region — check distance from left center.
		dx := px - radius
		dy := py - radius
		return dx*dx+dy*dy <= r2
	} else if px > float64(w)-radius {
		// In the right cap region — check distance from right center.
		dx := px - (float64(w) - radius)
		dy := py - radius
		return dx*dx+dy*dy <= r2
	}
	// In the middle rectangle — always inside.
	return true
}
