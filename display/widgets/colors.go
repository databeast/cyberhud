package widgets

import "image/color"

// ResolveColors applies default colors to zero-value inputs.
// Zero foreground → opaque white. Zero background → opaque black.
// Non-zero values pass through unchanged.
func ResolveColors(fg, bg color.RGBA) (color.RGBA, color.RGBA) {
	zero := color.RGBA{}
	if fg == zero {
		fg = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	}
	if bg == zero {
		bg = color.RGBA{R: 0, G: 0, B: 0, A: 255}
	}
	return fg, bg
}
