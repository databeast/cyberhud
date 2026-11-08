package styles

import (
	"fmt"

	"github.com/databeast/cyberhud/display/style"
)

// stamp builds a schedule style declaration: the canonical "<capability>-WxH"
// name, the surface requirements for this panel (the uniform one-row / eight-char
// floor every schedule style shares — the active-event title must fit), and the
// capability flag.
//
// Schedule has one layout (adaptiveBuild); these declarations exist only so the
// registry can pick the best-fitting capability/size cell for a given panel.
func stamp(prefix string, w, h int, cap style.Capability, color bool) def {
	return def{
		name: fmt.Sprintf("%s-%dx%d", prefix, w, h),
		reqs: style.SurfaceRequirements{
			MinWidth:        w,
			MinHeight:       h,
			Capability:      cap,
			MinRows:         1,
			MinCharsPerLine: 8,
		},
		p: Params{Color: color},
	}
}

func monoSlowStyle(w, h int) def      { return stamp("mono-slow", w, h, style.MonoSlow, false) }
func monoFastStyle(w, h int) def      { return stamp("mono-fast", w, h, style.MonoFast, false) }
func grayscaleSlowStyle(w, h int) def { return stamp("grayscale-slow", w, h, style.GrayscaleSlow, false) }
func grayscaleFastStyle(w, h int) def { return stamp("grayscale-fast", w, h, style.GrayscaleFast, false) }
func colorSlowStyle(w, h int) def     { return stamp("color-slow", w, h, style.ColorSlow, true) }
func colorFastStyle(w, h int) def     { return stamp("color-fast", w, h, style.ColorFast, true) }
