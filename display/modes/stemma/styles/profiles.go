package styles

import (
	"github.com/databeast/cyberhud/display/modes/stemma/source"
	"github.com/databeast/cyberhud/display/style"
)

func listProfile(name string, width, height int, capability style.Capability) def {
	return def{
		name: name,
		reqs: style.SurfaceRequirements{
			MinWidth:        width,
			MinHeight:       height,
			Capability:      capability,
			MinRows:         2,
			MinCharsPerLine: 10,
		},
		p: Params{Layout: layoutList, RowFormatter: defaultDeviceRow},
	}
}

// AllStyles returns the full STEMMA style matrix: generic layouts plus the
// explicit size/capability profiles used for panel-specific resolution.
func AllStyles() []style.Style[source.StemmaSnapshot, source.Policy] {
	return []style.Style[source.StemmaSnapshot, source.Policy]{
		ListStyle,
		CompactStyle,

		// MonoSlow
		listProfile("mono-slow-32x128", 32, 128, style.MonoSlow),
		listProfile("mono-slow-64x128", 64, 128, style.MonoSlow),
		listProfile("mono-slow-80x160", 80, 160, style.MonoSlow),
		listProfile("mono-slow-128x32", 128, 32, style.MonoSlow),
		listProfile("mono-slow-128x64", 128, 64, style.MonoSlow),
		listProfile("mono-slow-128x128", 128, 128, style.MonoSlow),
		listProfile("mono-slow-128x160", 128, 160, style.MonoSlow),
		listProfile("mono-slow-135x240", 135, 240, style.MonoSlow),
		listProfile("mono-slow-160x80", 160, 80, style.MonoSlow),
		listProfile("mono-slow-160x128", 160, 128, style.MonoSlow),
		listProfile("mono-slow-240x135", 240, 135, style.MonoSlow),
		listProfile("mono-slow-240x240", 240, 240, style.MonoSlow),
		listProfile("mono-slow-240x320", 240, 320, style.MonoSlow),
		listProfile("mono-slow-320x240", 320, 240, style.MonoSlow),
		listProfile("mono-slow-320x480", 320, 480, style.MonoSlow),
		listProfile("mono-slow-480x320", 480, 320, style.MonoSlow),
		MonoSlow800x480Style,

		// MonoFast
		listProfile("mono-128x32", 128, 32, style.MonoFast),
		listProfile("mono-128x64", 128, 64, style.MonoFast),
		listProfile("mono-128x128", 128, 128, style.MonoFast),
		listProfile("mono-32x128", 32, 128, style.MonoFast),
		listProfile("mono-64x128", 64, 128, style.MonoFast),

		// GrayscaleSlow
		listProfile("grayscale-slow-32x128", 32, 128, style.GrayscaleSlow),
		listProfile("grayscale-slow-64x128", 64, 128, style.GrayscaleSlow),
		listProfile("grayscale-slow-80x160", 80, 160, style.GrayscaleSlow),
		listProfile("grayscale-slow-104x212", 104, 212, style.GrayscaleSlow),
		listProfile("grayscale-slow-122x250", 122, 250, style.GrayscaleSlow),
		listProfile("grayscale-slow-128x32", 128, 32, style.GrayscaleSlow),
		listProfile("grayscale-slow-128x64", 128, 64, style.GrayscaleSlow),
		listProfile("grayscale-slow-128x128", 128, 128, style.GrayscaleSlow),
		listProfile("grayscale-slow-128x160", 128, 160, style.GrayscaleSlow),
		listProfile("grayscale-slow-128x296", 128, 296, style.GrayscaleSlow),
		listProfile("grayscale-slow-135x240", 135, 240, style.GrayscaleSlow),
		listProfile("grayscale-slow-160x80", 160, 80, style.GrayscaleSlow),
		listProfile("grayscale-slow-160x128", 160, 128, style.GrayscaleSlow),
		listProfile("grayscale-slow-176x264", 176, 264, style.GrayscaleSlow),
		listProfile("grayscale-slow-200x200", 200, 200, style.GrayscaleSlow),
		listProfile("grayscale-slow-212x104", 212, 104, style.GrayscaleSlow),
		listProfile("grayscale-slow-240x135", 240, 135, style.GrayscaleSlow),
		listProfile("grayscale-slow-240x240", 240, 240, style.GrayscaleSlow),
		listProfile("grayscale-slow-240x320", 240, 320, style.GrayscaleSlow),
		listProfile("grayscale-slow-250x122", 250, 122, style.GrayscaleSlow),
		listProfile("grayscale-slow-264x176", 264, 176, style.GrayscaleSlow),
		listProfile("grayscale-slow-296x128", 296, 128, style.GrayscaleSlow),
		listProfile("grayscale-slow-300x400", 300, 400, style.GrayscaleSlow),
		listProfile("grayscale-slow-320x240", 320, 240, style.GrayscaleSlow),
		listProfile("grayscale-slow-320x480", 320, 480, style.GrayscaleSlow),
		listProfile("grayscale-slow-400x300", 400, 300, style.GrayscaleSlow),
		listProfile("grayscale-slow-480x320", 480, 320, style.GrayscaleSlow),
		listProfile("grayscale-slow-480x800", 480, 800, style.GrayscaleSlow),
		listProfile("grayscale-slow-800x480", 800, 480, style.GrayscaleSlow),

		// GrayscaleFast
		listProfile("grayscale-fast-80x160", 80, 160, style.GrayscaleFast),
		listProfile("grayscale-fast-128x128", 128, 128, style.GrayscaleFast),
		listProfile("grayscale-fast-128x160", 128, 160, style.GrayscaleFast),
		listProfile("grayscale-fast-135x240", 135, 240, style.GrayscaleFast),
		listProfile("grayscale-fast-160x80", 160, 80, style.GrayscaleFast),
		listProfile("grayscale-fast-160x128", 160, 128, style.GrayscaleFast),
		listProfile("grayscale-fast-240x135", 240, 135, style.GrayscaleFast),
		listProfile("grayscale-fast-240x240", 240, 240, style.GrayscaleFast),
		listProfile("grayscale-fast-240x320", 240, 320, style.GrayscaleFast),
		listProfile("grayscale-fast-320x240", 320, 240, style.GrayscaleFast),
		listProfile("grayscale-fast-320x480", 320, 480, style.GrayscaleFast),
		listProfile("grayscale-fast-480x320", 480, 320, style.GrayscaleFast),
		listProfile("grayscale-fast-480x800", 480, 800, style.GrayscaleFast),
		listProfile("grayscale-fast-800x480", 800, 480, style.GrayscaleFast),

		// ColorSlow
		listProfile("color-slow-32x128", 32, 128, style.ColorSlow),
		listProfile("color-slow-64x128", 64, 128, style.ColorSlow),
		listProfile("color-slow-80x160", 80, 160, style.ColorSlow),
		listProfile("color-slow-104x212", 104, 212, style.ColorSlow),
		listProfile("color-slow-122x250", 122, 250, style.ColorSlow),
		listProfile("color-slow-128x32", 128, 32, style.ColorSlow),
		listProfile("color-slow-128x64", 128, 64, style.ColorSlow),
		listProfile("color-slow-128x128", 128, 128, style.ColorSlow),
		listProfile("color-slow-128x160", 128, 160, style.ColorSlow),
		listProfile("color-slow-128x296", 128, 296, style.ColorSlow),
		listProfile("color-slow-135x240", 135, 240, style.ColorSlow),
		listProfile("color-slow-160x80", 160, 80, style.ColorSlow),
		listProfile("color-slow-160x128", 160, 128, style.ColorSlow),
		listProfile("color-slow-176x264", 176, 264, style.ColorSlow),
		listProfile("color-slow-200x200", 200, 200, style.ColorSlow),
		listProfile("color-slow-212x104", 212, 104, style.ColorSlow),
		listProfile("color-slow-240x135", 240, 135, style.ColorSlow),
		listProfile("color-slow-240x240", 240, 240, style.ColorSlow),
		listProfile("color-slow-240x320", 240, 320, style.ColorSlow),
		listProfile("color-slow-250x122", 250, 122, style.ColorSlow),
		listProfile("color-slow-264x176", 264, 176, style.ColorSlow),
		listProfile("color-slow-296x128", 296, 128, style.ColorSlow),
		listProfile("color-slow-300x400", 300, 400, style.ColorSlow),
		listProfile("color-slow-320x240", 320, 240, style.ColorSlow),
		listProfile("color-slow-320x480", 320, 480, style.ColorSlow),
		listProfile("color-slow-400x300", 400, 300, style.ColorSlow),
		listProfile("color-slow-480x320", 480, 320, style.ColorSlow),
		listProfile("color-slow-480x800", 480, 800, style.ColorSlow),
		listProfile("color-slow-800x480", 800, 480, style.ColorSlow),

		// ColorFast
		listProfile("color-80x160", 80, 160, style.ColorFast),
		listProfile("color-128x128", 128, 128, style.ColorFast),
		listProfile("color-128x160", 128, 160, style.ColorFast),
		listProfile("color-medium-135x240", 135, 240, style.ColorFast),
		listProfile("color-160x80", 160, 80, style.ColorFast),
		listProfile("color-160x128", 160, 128, style.ColorFast),
		listProfile("color-medium-240x135", 240, 135, style.ColorFast),
		listProfile("color-medium-240x240", 240, 240, style.ColorFast),
		listProfile("color-medium-240x320", 240, 320, style.ColorFast),
		listProfile("color-medium-320x240", 320, 240, style.ColorFast),
		listProfile("color-large-320x480", 320, 480, style.ColorFast),
		listProfile("color-large-480x320", 480, 320, style.ColorFast),
		listProfile("color-large-480x800", 480, 800, style.ColorFast),
		listProfile("color-large-800x480", 800, 480, style.ColorFast),
	}
}
