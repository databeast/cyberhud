package styles

import (
	"github.com/databeast/cyberhud/display/style"
)

// Skeleton registrations that fill per-resolution/capability gaps. They reuse
// plainBuild until bespoke layouts exist, and deliberately omit MinRows and
// MinCharsPerLine so tiny panels are not rejected.

func skeleton(name string, w, h int, c style.Capability) def {
	return def{
		name: name,
		reqs: style.SurfaceRequirements{MinWidth: w, MinHeight: h, Capability: c},
		p:    Params{BuildFn: plainBuild},
	}
}

var (
	MonoFast128x32Style  = skeleton("mono-fast-128x32", 128, 32, style.MonoFast)
	MonoFast128x64Style  = skeleton("mono-fast-128x64", 128, 64, style.MonoFast)
	MonoFast168x144Style = skeleton("mono-fast-168x144", 168, 144, style.MonoFast)
	MonoFast400x240Style = skeleton("mono-fast-400x240", 400, 240, style.MonoFast)

	GrayscaleFast15x7Style   = skeleton("grayscale-fast-15x7", 15, 7, style.GrayscaleFast)
	GrayscaleFast16x8Style   = skeleton("grayscale-fast-16x8", 16, 8, style.GrayscaleFast)
	GrayscaleFast128x64Style = skeleton("grayscale-fast-128x64", 128, 64, style.GrayscaleFast)

	ColorSlow200x200Style = skeleton("color-slow-200x200", 200, 200, style.ColorSlow)

	ColorFast13x9Style    = skeleton("color-fast-13x9", 13, 9, style.ColorFast)
	ColorFast128x128Style = skeleton("color-fast-128x128", 128, 128, style.ColorFast)
	ColorFast128x160Style = skeleton("color-fast-128x160", 128, 160, style.ColorFast)
	ColorFast160x80Style  = skeleton("color-fast-160x80", 160, 80, style.ColorFast)
	ColorFast240x135Style = skeleton("color-fast-240x135", 240, 135, style.ColorFast)
	ColorFast240x240Style = skeleton("color-fast-240x240", 240, 240, style.ColorFast)
	ColorFast240x320Style = skeleton("color-fast-240x320", 240, 320, style.ColorFast)
	ColorFast320x240Style = skeleton("color-fast-320x240", 320, 240, style.ColorFast)
	ColorFast480x320Style = skeleton("color-fast-480x320", 480, 320, style.ColorFast)
)
