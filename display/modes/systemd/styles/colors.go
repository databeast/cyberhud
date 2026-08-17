package styles

import (
	"image/color"

	"github.com/databeast/cyberhud/display/modes/systemd/source"
	sharedcolor "github.com/databeast/cyberhud/display/style/color"
)

// resolveBootAccent returns the appropriate accent color for boot progress rendering.
// When bootComplete is true, returns green {0x00, 0xFF, 0x00, 0xFF} regardless of accent.
// Otherwise delegates to sharedcolor.ResolveAccent which handles "none" → opaque white.
//
// Framework pattern demonstrated: ColorAccent via sharedcolor.ResolveAccent.
func resolveBootAccent(accent string, bootComplete bool) color.RGBA {
	if bootComplete {
		return color.RGBA{R: 0x00, G: 0xFF, B: 0x00, A: 0xFF}
	}
	return sharedcolor.ResolveAccent(accent)
}

// gradientAccent returns a desaturated and dimmed variant of the boot accent
// color, suitable for use as a gradient background stop. The full-intensity
// accent is too vivid for a background fill; this applies 50% desaturation
// (blend toward luminance gray) followed by 60% brightness reduction.
func gradientAccent(accent string, bootComplete bool) color.RGBA {
	c := resolveBootAccent(accent, bootComplete)
	return desaturateAndDim(c)
}

// desaturateAndDim applies 50% desaturation then dims to 60% brightness.
// Desaturation blends each channel toward the color's perceptual luminance.
// The result is a muted, dark color suitable for gradient backgrounds behind text.
func desaturateAndDim(c color.RGBA) color.RGBA {
	// Perceptual luminance (BT.601 coefficients).
	lum := 0.299*float64(c.R) + 0.587*float64(c.G) + 0.114*float64(c.B)

	// 50% desaturation: blend each channel halfway toward luminance.
	r := float64(c.R)*0.5 + lum*0.5
	g := float64(c.G)*0.5 + lum*0.5
	b := float64(c.B)*0.5 + lum*0.5

	// 60% dim (multiply by 0.6).
	r *= 0.6
	g *= 0.6
	b *= 0.6

	return color.RGBA{
		R: clampByte(r),
		G: clampByte(g),
		B: clampByte(b),
		A: 255,
	}
}

// clampByte clamps a float64 to [0, 255] and returns the uint8 value.
func clampByte(v float64) uint8 {
	if v < 0 {
		return 0
	}
	if v > 255 {
		return 255
	}
	return uint8(v)
}

// bootFraction returns the boot completion as a float64 in [0.0, 1.0].
// Uses active/total target counts when available; falls back to a simple
// state-machine heuristic when target enumeration returns zeros.
//
// Framework pattern demonstrated: Boot completion metric for gradient positioning.
func bootFraction(snap source.Snapshot) float64 {
	if snap.BootComplete {
		return 1.0
	}
	if snap.TotalTargets > 0 {
		frac := float64(snap.ActiveTargets) / float64(snap.TotalTargets)
		if frac > 1.0 {
			frac = 1.0
		}
		return frac
	}
	// Fallback heuristic: 0% = no loading, 50% = loading, 75% = loaded one
	if snap.Loaded != "" {
		return 0.75
	}
	if snap.Loading != "" {
		return 0.50
	}
	return 0.0
}
