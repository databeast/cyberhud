package styles

import (
	"image/color"

	"github.com/databeast/cyberhud/display/style"
)

// loadColor maps a normalized load value in [0,1] to a green→amber→red color,
// the system mode's signature heat encoding (see the doc.go # User Experience
// section): cool green at idle, amber at moderate load, urgent red near
// saturation. This lets a glance across the cores wall or the top leaderboard
// read the machine's effort by color alone.
//
// isColor gates the encoding: on color-capable panels the full ramp is used;
// on grayscale/mono panels color cannot carry the meaning, so the caller should
// fall back to fill level / ordering (loadColor returns a bright foreground so
// the element is still visible, but conveys no hue information there).
func loadColor(load float64, isColor bool) color.RGBA {
	if load < 0 {
		load = 0
	}
	if load > 1 {
		load = 1
	}
	if !isColor {
		// Grayscale/mono: a single bright foreground. Luminance/fill carries the
		// magnitude; hue cannot.
		return color.RGBA{R: 255, G: 255, B: 255, A: 255}
	}
	// Piecewise green→amber→red:
	//   0.0 → green (0,200,0)
	//   0.5 → amber (230,170,0)
	//   1.0 → red   (230,40,30)
	switch {
	case load <= 0.5:
		t := load / 0.5
		return lerpRGBA(color.RGBA{0, 200, 0, 255}, color.RGBA{230, 170, 0, 255}, t)
	default:
		t := (load - 0.5) / 0.5
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
