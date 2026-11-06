package styles

// Shared constructors for pager styles.
//
// Every per-resolution style in the pager grid is produced by one of these
// constructors rather than hand-written. A constructor stamps the canonical
// name for its capability class, the surface requirements for the given
// dimensions, and the shared core placement (via the zero-value Params, whose
// Build dispatches to coreBuild). This is the single place a pager style's name
// format and requirements shape are defined, so the whole grid stays uniform.

import (
	"fmt"

	"github.com/databeast/cyberhud/display/style"
)

// dim is a panel resolution in pixels.
type dim struct {
	w, h int
}

// pagerReqs builds the surface requirements shared by every pager style:
// minimum dimensions from the panel size, the capability class, and a floor of
// one row / one char per line so a style never resolves onto a panel too small
// to show anything.
func pagerReqs(d dim, cap style.Capability) style.SurfaceRequirements {
	return style.SurfaceRequirements{
		MinWidth:        d.w,
		MinHeight:       d.h,
		Capability:      cap,
		MinRows:         1,
		MinCharsPerLine: 1,
	}
}

// newStyle is the single stamp: canonical name + requirements + shared core
// placement. The name encodes capability prefix and WxH, matching the house
// convention used by clock and thermal.
func newStyle(prefix string, d dim, cap style.Capability) def {
	return def{
		name: fmt.Sprintf("%s-%dx%d", prefix, d.w, d.h),
		reqs: pagerReqs(d, cap),
	}
}

func monoSlow(d dim) def      { return newStyle("mono-slow", d, style.MonoSlow) }
func monoFast(d dim) def      { return newStyle("mono-fast", d, style.MonoFast) }
func grayscaleSlow(d dim) def { return newStyle("grayscale-slow", d, style.GrayscaleSlow) }
func grayscaleFast(d dim) def { return newStyle("grayscale-fast", d, style.GrayscaleFast) }
func colorSlow(d dim) def     { return newStyle("color-slow", d, style.ColorSlow) }
func colorFast(d dim) def     { return newStyle("color-fast", d, style.ColorFast) }
