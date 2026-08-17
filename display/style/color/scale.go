package color

import "image/color"

// Scale multiplies each RGB channel of c by factor, clamping factor to [0.0, 1.0].
// Alpha is always forced to 255 regardless of input.
func Scale(c color.RGBA, factor float64) color.RGBA {
	if factor <= 0.0 {
		return color.RGBA{R: 0, G: 0, B: 0, A: 255}
	}
	if factor >= 1.0 {
		return color.RGBA{R: c.R, G: c.G, B: c.B, A: 255}
	}
	return color.RGBA{
		R: uint8(float64(c.R) * factor),
		G: uint8(float64(c.G) * factor),
		B: uint8(float64(c.B) * factor),
		A: 255,
	}
}

// Gradient generates a brightness ramp of n steps from c (index 0, full brightness)
// down to black (index n-1). Each step i has factor = (n-1-i) / (n-1).
// Returns nil if n < 1. Returns []color.RGBA{c with A=255} if n == 1.
func Gradient(c color.RGBA, n int) []color.RGBA {
	if n < 1 {
		return nil
	}
	if n == 1 {
		return []color.RGBA{{R: c.R, G: c.G, B: c.B, A: 255}}
	}
	result := make([]color.RGBA, n)
	for i := 0; i < n; i++ {
		factor := float64(n-1-i) / float64(n-1)
		result[i] = Scale(c, factor)
	}
	return result
}
