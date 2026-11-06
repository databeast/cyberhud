package styles

import (
	"fmt"

	"github.com/databeast/cyberhud/display/style"
)

// stamp builds a serial style declaration: the canonical "<prefix>-WxH" name,
// the surface requirements for this panel (the uniform two-row / eight-char
// floor every serial style shares), and the engine-selection Params.
//
// The per-resolution styles_WxH.go files call this so each cell stays a short,
// hand-editable one-liner whose Params are the thing worth tweaking — rather
// than repeating the identical requirements boilerplate in every literal.
func stamp(prefix string, w, h int, cap style.Capability, p Params) def {
	return def{
		name: fmt.Sprintf("%s-%dx%d", prefix, w, h),
		reqs: style.SurfaceRequirements{
			MinWidth:        w,
			MinHeight:       h,
			Capability:      cap,
			MinRows:         2,
			MinCharsPerLine: 8,
		},
		p: p,
	}
}

// Per-capability stamps. Slow classes use the e-paper-safe static panel
// (zero-value Params); fast classes use the live panel, with color panels also
// engaging the accent palette.
func monoSlowStyle(w, h int) def { return stamp("mono-slow", w, h, style.MonoSlow, Params{}) }
func grayscaleSlowStyle(w, h int) def {
	return stamp("grayscale-slow", w, h, style.GrayscaleSlow, Params{})
}
func colorSlowStyle(w, h int) def { return stamp("color-slow", w, h, style.ColorSlow, Params{}) }

func monoFastStyle(w, h int) def { return stamp("mono-fast", w, h, style.MonoFast, Params{Fast: true}) }
func grayscaleFastStyle(w, h int) def {
	return stamp("grayscale-fast", w, h, style.GrayscaleFast, Params{Fast: true})
}
func colorFastStyle(w, h int) def {
	return stamp("color-fast", w, h, style.ColorFast, Params{Fast: true, Color: true})
}
