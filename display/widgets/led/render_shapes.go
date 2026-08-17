package led

import (
	"image"
	"image/color"
	"math"
)

// dimColor computes the dimmed foreground color for off-state LED outline.
// Returns (floor(R×0.3), floor(G×0.3), floor(B×0.3), A).
func dimColor(c color.RGBA) color.RGBA {
	return color.RGBA{
		R: uint8(math.Floor(float64(c.R) * 0.3)),
		G: uint8(math.Floor(float64(c.G) * 0.3)),
		B: uint8(math.Floor(float64(c.B) * 0.3)),
		A: c.A,
	}
}

// scaleBrightness applies brightness scaling to an RGBA color.
// Each RGB channel = floor(channel × brightness), alpha unchanged.
func scaleBrightness(c color.RGBA, brightness float64) color.RGBA {
	if brightness >= 1.0 {
		return c
	}
	if brightness <= 0.0 {
		return color.RGBA{R: 0, G: 0, B: 0, A: c.A}
	}
	return color.RGBA{
		R: uint8(math.Floor(float64(c.R) * brightness)),
		G: uint8(math.Floor(float64(c.G) * brightness)),
		B: uint8(math.Floor(float64(c.B) * brightness)),
		A: c.A,
	}
}

// renderCircleBody renders a filled circle inscribed within bodyRect.
// For On state: fill with brightness-scaled foreground color.
// For Off state: 1px outline at dimmed color, interior at background.
func renderCircleBody(img *image.RGBA, bodyRect image.Rectangle, fillColor, bgColor color.RGBA, brightness float64, isOff bool) {
	w := bodyRect.Dx()
	h := bodyRect.Dy()
	cx := float64(w) / 2.0
	cy := float64(h) / 2.0
	radius := float64(w) / 2.0 // inscribed circle, use min side / 2
	if float64(h)/2.0 < radius {
		radius = float64(h) / 2.0
	}

	scaledFill := scaleBrightness(fillColor, brightness)
	dimmed := dimColor(fillColor)

	for py := bodyRect.Min.Y; py < bodyRect.Max.Y; py++ {
		for px := bodyRect.Min.X; px < bodyRect.Max.X; px++ {
			lx := px - bodyRect.Min.X
			ly := py - bodyRect.Min.Y

			dx := float64(lx) + 0.5 - cx
			dy := float64(ly) + 0.5 - cy
			dist := math.Sqrt(dx*dx + dy*dy)

			if dist <= radius {
				if isOff {
					// Check if outline pixel (inside but has neighbor outside)
					if isCircleOutlinePixel(lx, ly, cx, cy, radius) {
						img.SetRGBA(px, py, dimmed)
					} else {
						img.SetRGBA(px, py, bgColor)
					}
				} else {
					img.SetRGBA(px, py, scaledFill)
				}
			}
		}
	}
}

// isCircleOutlinePixel checks if pixel at (lx, ly) in local coords is on the
// circle outline (inside but has at least one 4-neighbor outside).
func isCircleOutlinePixel(lx, ly int, cx, cy, radius float64) bool {
	neighbors := [4][2]int{
		{lx - 1, ly},
		{lx + 1, ly},
		{lx, ly - 1},
		{lx, ly + 1},
	}
	for _, n := range neighbors {
		dx := float64(n[0]) + 0.5 - cx
		dy := float64(n[1]) + 0.5 - cy
		if math.Sqrt(dx*dx+dy*dy) > radius {
			return true
		}
	}
	return false
}

// renderSquareBody renders a filled square occupying the full bodyRect.
// For On state: fill with brightness-scaled foreground color.
// For Off state: 1px outline at dimmed color, interior at background.
func renderSquareBody(img *image.RGBA, bodyRect image.Rectangle, fillColor, bgColor color.RGBA, brightness float64, isOff bool) {
	scaledFill := scaleBrightness(fillColor, brightness)
	dimmed := dimColor(fillColor)

	for py := bodyRect.Min.Y; py < bodyRect.Max.Y; py++ {
		for px := bodyRect.Min.X; px < bodyRect.Max.X; px++ {
			if isOff {
				// Outline: pixel on the edge of bodyRect (1px border)
				if px == bodyRect.Min.X || px == bodyRect.Max.X-1 ||
					py == bodyRect.Min.Y || py == bodyRect.Max.Y-1 {
					img.SetRGBA(px, py, dimmed)
				} else {
					img.SetRGBA(px, py, bgColor)
				}
			} else {
				img.SetRGBA(px, py, scaledFill)
			}
		}
	}
}

// renderDiamondBody renders a filled diamond (45° rotated square) inscribed
// within bodyRect. Pixel membership uses Manhattan distance:
// |px - cx| + |py - cy| ≤ bodyRadius
// For On state: fill with brightness-scaled foreground color.
// For Off state: 1px outline at dimmed color, interior at background.
func renderDiamondBody(img *image.RGBA, bodyRect image.Rectangle, fillColor, bgColor color.RGBA, brightness float64, isOff bool) {
	w := bodyRect.Dx()
	h := bodyRect.Dy()
	cx := float64(w) / 2.0
	cy := float64(h) / 2.0
	// bodyRadius is half the body side length (use the smaller dimension)
	bodyRadius := float64(w) / 2.0
	if float64(h)/2.0 < bodyRadius {
		bodyRadius = float64(h) / 2.0
	}

	scaledFill := scaleBrightness(fillColor, brightness)
	dimmed := dimColor(fillColor)

	for py := bodyRect.Min.Y; py < bodyRect.Max.Y; py++ {
		for px := bodyRect.Min.X; px < bodyRect.Max.X; px++ {
			lx := px - bodyRect.Min.X
			ly := py - bodyRect.Min.Y

			// Pixel center
			pcx := float64(lx) + 0.5
			pcy := float64(ly) + 0.5

			// Manhattan distance from center
			dist := math.Abs(pcx-cx) + math.Abs(pcy-cy)

			if dist <= bodyRadius {
				if isOff {
					if isDiamondOutlinePixel(lx, ly, cx, cy, bodyRadius) {
						img.SetRGBA(px, py, dimmed)
					} else {
						img.SetRGBA(px, py, bgColor)
					}
				} else {
					img.SetRGBA(px, py, scaledFill)
				}
			}
		}
	}
}

// isDiamondOutlinePixel checks if a pixel inside the diamond has at least one
// 4-connected neighbor outside the diamond (Manhattan distance > bodyRadius).
func isDiamondOutlinePixel(lx, ly int, cx, cy, bodyRadius float64) bool {
	neighbors := [4][2]int{
		{lx - 1, ly},
		{lx + 1, ly},
		{lx, ly - 1},
		{lx, ly + 1},
	}
	for _, n := range neighbors {
		ncx := float64(n[0]) + 0.5
		ncy := float64(n[1]) + 0.5
		if math.Abs(ncx-cx)+math.Abs(ncy-cy) > bodyRadius {
			return true
		}
	}
	return false
}

// renderRoundedSquareBody renders a filled rounded rectangle within bodyRect.
// Corner radius = floor(bodySide * 0.25).
// A pixel is inside if:
//   - It is within the central cross (not in a corner zone), OR
//   - It is within a corner zone AND its distance from the corner arc center ≤ cornerRadius
//
// For On state: fill with brightness-scaled foreground color.
// For Off state: 1px outline at dimmed color, interior at background.
func renderRoundedSquareBody(img *image.RGBA, bodyRect image.Rectangle, fillColor, bgColor color.RGBA, brightness float64, isOff bool) {
	w := bodyRect.Dx()
	h := bodyRect.Dy()
	// Corner radius = 25% of body side length, rounded down
	side := w
	if h < side {
		side = h
	}
	cornerRadius := side / 4 // integer division floors

	scaledFill := scaleBrightness(fillColor, brightness)
	dimmed := dimColor(fillColor)

	// Corner arc centers (in local coordinates)
	cr := float64(cornerRadius)
	// Top-left corner center
	tlx, tly := cr, cr
	// Top-right corner center
	trx, try_ := float64(w)-cr, cr
	// Bottom-left corner center
	blx, bly := cr, float64(h)-cr
	// Bottom-right corner center
	brx, bry := float64(w)-cr, float64(h)-cr

	for py := bodyRect.Min.Y; py < bodyRect.Max.Y; py++ {
		for px := bodyRect.Min.X; px < bodyRect.Max.X; px++ {
			lx := px - bodyRect.Min.X
			ly := py - bodyRect.Min.Y

			// Pixel center in local coords
			pcx := float64(lx) + 0.5
			pcy := float64(ly) + 0.5

			inside := isInsideRoundedRect(pcx, pcy, float64(w), float64(h), cr, tlx, tly, trx, try_, blx, bly, brx, bry)

			if inside {
				if isOff {
					if isRoundedRectOutlinePixel(lx, ly, w, h, cr, tlx, tly, trx, try_, blx, bly, brx, bry) {
						img.SetRGBA(px, py, dimmed)
					} else {
						img.SetRGBA(px, py, bgColor)
					}
				} else {
					img.SetRGBA(px, py, scaledFill)
				}
			}
		}
	}
}

// isInsideRoundedRect checks if pixel center (pcx, pcy) is inside the rounded rectangle.
func isInsideRoundedRect(pcx, pcy, w, h, cr, tlx, tly, trx, try_, blx, bly, brx, bry float64) bool {
	// Must be within the bounding box first
	if pcx < 0 || pcx > w || pcy < 0 || pcy > h {
		return false
	}

	// Check if in a corner zone
	inLeftCol := pcx < cr
	inRightCol := pcx > w-cr
	inTopRow := pcy < cr
	inBottomRow := pcy > h-cr

	// If in a corner zone, check distance from corner arc center
	if inLeftCol && inTopRow {
		// Top-left corner
		dx := pcx - tlx
		dy := pcy - tly
		return math.Sqrt(dx*dx+dy*dy) <= cr
	}
	if inRightCol && inTopRow {
		// Top-right corner
		dx := pcx - trx
		dy := pcy - try_
		return math.Sqrt(dx*dx+dy*dy) <= cr
	}
	if inLeftCol && inBottomRow {
		// Bottom-left corner
		dx := pcx - blx
		dy := pcy - bly
		return math.Sqrt(dx*dx+dy*dy) <= cr
	}
	if inRightCol && inBottomRow {
		// Bottom-right corner
		dx := pcx - brx
		dy := pcy - bry
		return math.Sqrt(dx*dx+dy*dy) <= cr
	}

	// In the central cross (not in any corner zone) — always inside
	return true
}

// isRoundedRectOutlinePixel checks if a pixel inside the rounded rectangle has at least
// one 4-connected neighbor outside the rounded rectangle.
func isRoundedRectOutlinePixel(lx, ly, w, h int, cr, tlx, tly, trx, try_, blx, bly, brx, bry float64) bool {
	neighbors := [4][2]int{
		{lx - 1, ly},
		{lx + 1, ly},
		{lx, ly - 1},
		{lx, ly + 1},
	}
	fw := float64(w)
	fh := float64(h)
	for _, n := range neighbors {
		ncx := float64(n[0]) + 0.5
		ncy := float64(n[1]) + 0.5
		if !isInsideRoundedRect(ncx, ncy, fw, fh, cr, tlx, tly, trx, try_, blx, bly, brx, bry) {
			return true
		}
	}
	return false
}
