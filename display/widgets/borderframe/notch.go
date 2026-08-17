package borderframe

import (
	"image"
	"image/color"
)

// renderNotches draws ticker notch decorations along all four edges of the border.
// Notches are 1px-wide perpendicular lines extending NotchLength pixels inward from
// the outer boundary. They are placed at NotchInterval pixel intervals starting from
// the corner tile boundary (8px from each edge start), skipping the 8×8 corner
// exclusion zones and omitting any notch within 1px of the far corner boundary.
func renderNotches(img *image.RGBA, cfg Config) {
	if cfg.NotchInterval <= 0 {
		return
	}

	w := img.Bounds().Dx()
	h := img.Bounds().Dy()

	notchLen := cfg.NotchLength
	if notchLen <= 0 {
		notchLen = 2
	}
	if notchLen > 8 {
		notchLen = 8
	}

	tint := effectiveTint(cfg.ColorTint)
	c := color.RGBA{R: tint.R, G: tint.G, B: tint.B, A: 255}

	// Corner exclusion zone is 8 pixels from each edge start/end.
	const cornerZone = tileSize // 8

	// Top edge: vertical notches drawn from y=0 downward for notchLen pixels.
	// x positions start at cornerZone + NotchInterval, repeating every NotchInterval.
	// Stop when x >= (w - cornerZone - 1) (far corner exclusion + 1px margin).
	farLimitX := w - cornerZone - 1
	for x := cornerZone + cfg.NotchInterval; x < farLimitX; x += cfg.NotchInterval {
		for dy := 0; dy < notchLen; dy++ {
			setPixel(img, x, dy, c)
		}
	}

	// Bottom edge: vertical notches drawn from y=(h-1) upward for notchLen pixels.
	for x := cornerZone + cfg.NotchInterval; x < farLimitX; x += cfg.NotchInterval {
		for dy := 0; dy < notchLen; dy++ {
			setPixel(img, x, h-1-dy, c)
		}
	}

	// Left edge: horizontal notches drawn from x=0 rightward for notchLen pixels.
	// y positions start at cornerZone + NotchInterval, repeating every NotchInterval.
	// Stop when y >= (h - cornerZone - 1) (far corner exclusion + 1px margin).
	farLimitY := h - cornerZone - 1
	for y := cornerZone + cfg.NotchInterval; y < farLimitY; y += cfg.NotchInterval {
		for dx := 0; dx < notchLen; dx++ {
			setPixel(img, dx, y, c)
		}
	}

	// Right edge: horizontal notches drawn from x=(w-1) leftward for notchLen pixels.
	for y := cornerZone + cfg.NotchInterval; y < farLimitY; y += cfg.NotchInterval {
		for dx := 0; dx < notchLen; dx++ {
			setPixel(img, w-1-dx, y, c)
		}
	}
}

// setPixel writes a color.RGBA directly into an *image.RGBA at (x, y).
func setPixel(img *image.RGBA, x, y int, c color.RGBA) {
	if x < img.Bounds().Min.X || x >= img.Bounds().Max.X ||
		y < img.Bounds().Min.Y || y >= img.Bounds().Max.Y {
		return
	}
	i := (y-img.Bounds().Min.Y)*img.Stride + (x-img.Bounds().Min.X)*4
	img.Pix[i+0] = c.R
	img.Pix[i+1] = c.G
	img.Pix[i+2] = c.B
	img.Pix[i+3] = c.A
}
