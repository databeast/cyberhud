package styles

import (
	"strconv"

	"github.com/databeast/cyberhud/display/modes/storage/source"
	"github.com/databeast/cyberhud/display/style"
)

// dim renders a pixel dimension as a decimal string for style-name composition.
func dim(n int) string { return strconv.Itoa(n) }

// Variant suffixes for storage style names. The full name is
// "<capability-prefix>-<WxH>[-<variant>]"; a name with no suffix is the
// skeleton, which routes to adaptiveBuild.
const (
	VariantBars  = "bars"
	VariantFocus = "focus"
)

// mkStyles returns the three defs for one resolution+capability: a skeleton
// (adaptive), the explicit bars variant, and the explicit focus variant. The
// names follow the pattern enforced by styleNamePattern in the parent package.
func mkStyles(prefix string, w, h int, cap style.Capability) []def {
	reqs := style.SurfaceRequirements{MinWidth: w, MinHeight: h, Capability: cap}
	base := prefix // e.g. "color-fast-240x240"
	return []def{
		{name: base, reqs: reqs},
		{name: base + "-" + VariantBars, reqs: reqs, p: Params{BuildFn: buildBarsStyle}},
		{name: base + "-" + VariantFocus, reqs: reqs, p: Params{BuildFn: buildFocusStyle}},
	}
}

// All returns every registered storage style, grouped by capability class and
// then by resolution, with the skeleton first within each group so fitness
// tie-breaking prefers the adaptive layout for a bare resolution match.
//
// The resolution set covers the panels the project targets: small mono OLEDs
// (128x64, 128x32), square color/e-paper TFTs (128x128, 240x240, 200x200),
// landscape TFTs (160x128, 320x240, 480x320, 800x480), portrait e-paper
// (122x250, 128x296), and the common 135x240/240x320 ST7789 displays. A bare
// MinWidth/MinHeight requirement means any larger panel in the same capability
// class still fits via the fitness scorer, so uncommon sizes degrade sensibly
// rather than failing to resolve.
func All() []style.Style[source.StorageSnapshot, source.Policy] {
	type entry struct {
		prefix string
		w, h   int
		cap    style.Capability
	}

	// Resolutions shared across capability classes.
	resolutions := []struct {
		w, h int
	}{
		{128, 32},
		{128, 64},
		{128, 128},
		{128, 296},
		{122, 250},
		{135, 240},
		{160, 128},
		{200, 200},
		{240, 240},
		{240, 320},
		{320, 240},
		{480, 320},
		{800, 480},
	}

	caps := []struct {
		prefix string
		cap    style.Capability
	}{
		{"mono-fast", style.MonoFast},
		{"mono-slow", style.MonoSlow},
		{"grayscale-fast", style.GrayscaleFast},
		{"grayscale-slow", style.GrayscaleSlow},
		{"color-fast", style.ColorFast},
		{"color-slow", style.ColorSlow},
	}

	var entries []entry
	for _, c := range caps {
		for _, r := range resolutions {
			prefix := c.prefix + "-" + dim(r.w) + "x" + dim(r.h)
			entries = append(entries, entry{prefix: prefix, w: r.w, h: r.h, cap: c.cap})
		}
	}

	out := make([]style.Style[source.StorageSnapshot, source.Policy], 0, len(entries)*3)
	for _, e := range entries {
		for _, d := range mkStyles(e.prefix, e.w, e.h, e.cap) {
			out = append(out, d)
		}
	}
	return out
}
