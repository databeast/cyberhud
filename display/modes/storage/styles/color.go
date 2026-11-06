package styles

import (
	"image/color"

	"github.com/databeast/cyberhud/display/style"
)

// usageColor maps a fill fraction in [0,1] to a green→amber→red heat ramp — the
// storage mode's signature encoding: cool green for plenty of room, amber as the
// filesystem fills, urgent red near capacity. This lets a glance across the
// mounts read disk pressure by color alone.
//
// isColor gates the encoding: on color-capable panels the full ramp is used; on
// grayscale/mono panels color cannot carry meaning, so a single bright
// foreground is returned and the bar's fill level carries the magnitude.
func usageColor(fill float64, isColor bool) color.RGBA {
	if fill < 0 {
		fill = 0
	}
	if fill > 1 {
		fill = 1
	}
	if !isColor {
		return color.RGBA{R: 255, G: 255, B: 255, A: 255}
	}
	// Piecewise green→amber→red, with red reserved for the last stretch so a
	// nearly-full disk reads as alarming:
	//   0.00 → green (0,200,0)
	//   0.80 → amber (230,170,0)
	//   1.00 → red   (230,40,30)
	switch {
	case fill <= 0.80:
		t := fill / 0.80
		return lerpRGBA(color.RGBA{0, 200, 0, 255}, color.RGBA{230, 170, 0, 255}, t)
	default:
		t := (fill - 0.80) / 0.20
		return lerpRGBA(color.RGBA{230, 170, 0, 255}, color.RGBA{230, 40, 30, 255}, t)
	}
}

// lerpRGBA linearly interpolates between two colors at t in [0,1].
func lerpRGBA(a, b color.RGBA, t float64) color.RGBA {
	if t < 0 {
		t = 0
	}
	if t > 1 {
		t = 1
	}
	return color.RGBA{
		R: uint8(float64(a.R) + (float64(b.R)-float64(a.R))*t),
		G: uint8(float64(a.G) + (float64(b.G)-float64(a.G))*t),
		B: uint8(float64(a.B) + (float64(b.B)-float64(a.B))*t),
		A: 255,
	}
}

// isColorPanel reports whether the panel can render color, from the style
// context capability. Grayscale and mono fall back to level/ordering encoding.
func isColorPanel(ctx style.StyleContext) bool {
	return ctx.Cap() >= style.ColorSlow
}

// barBackground returns the progress-bar track color for the panel type: a dark
// gray track on color panels, transparent on mono/grayscale where a filled track
// would smear the 1-bit output.
func barBackground(isColor bool) color.RGBA {
	if isColor {
		return color.RGBA{R: 40, G: 40, B: 40, A: 255}
	}
	return color.RGBA{R: 0, G: 0, B: 0, A: 0}
}
