package styles

// Programmatic generation of the pager style grid.
//
// The pager supports the same panel resolutions as the other text modes, across
// every capability class. Rather than hand-write one `def` literal per
// capability×resolution (the shape the old pager carried, and the shape that
// rots the moment a dimension is added), the grid is a dimension table per class
// run through the shared constructors in constructors.go. Each class's exported
// slice is consumed by the registry in styles.go.
//
// Adding a panel size is a one-line edit to the relevant table; the style, its
// canonical name, and its requirements all follow from the constructor.

import "github.com/databeast/cyberhud/display/style"

// slowDims are the e-paper resolutions shared by every slow capability class
// (both portrait and landscape orientations of each physical panel).
var slowDims = []dim{
	{122, 250}, {176, 264}, {200, 200}, {212, 104}, {296, 128},
	{400, 300}, {480, 800}, {800, 480}, {104, 212}, {250, 122},
	{128, 296}, {264, 176}, {300, 400},
}

// monoFastDims are the OLED resolutions for the mono-fast class.
var monoFastDims = []dim{
	{128, 32}, {128, 64}, {128, 128}, {32, 128}, {64, 128},
}

// fastDims are the TFT/LED-matrix resolutions shared by the grayscale-fast and
// color-fast classes.
var fastDims = []dim{
	{160, 80}, {160, 128}, {240, 135}, {240, 240}, {320, 240},
	{480, 320}, {800, 480}, {80, 160}, {128, 160}, {135, 240},
	{240, 320}, {320, 480}, {480, 800}, {128, 128},
}

// generateStyles maps a dimension table through a constructor into a style slice.
func generateStyles(dims []dim, ctor func(dim) def) []style.Style[snapshot, policy] {
	out := make([]style.Style[snapshot, policy], len(dims))
	for i, d := range dims {
		out[i] = ctor(d)
	}
	return out
}

// Per-class generated style sets, in the capability order the registry expects.
var (
	monoSlowStyles      = generateStyles(slowDims, monoSlow)
	monoFastStyles      = generateStyles(monoFastDims, monoFast)
	grayscaleSlowStyles = generateStyles(slowDims, grayscaleSlow)
	grayscaleFastStyles = generateStyles(fastDims, grayscaleFast)
	colorSlowStyles     = generateStyles(slowDims, colorSlow)
	colorFastStyles     = generateStyles(fastDims, colorFast)
)

// allStyles is every generated style, concatenated in capability order
// (MonoSlow → MonoFast → GrayscaleSlow → GrayscaleFast → ColorSlow → ColorFast),
// which also defines fitness tie-breaking (earlier wins).
func allStyles() []style.Style[snapshot, policy] {
	var out []style.Style[snapshot, policy]
	out = append(out, monoSlowStyles...)
	out = append(out, monoFastStyles...)
	out = append(out, grayscaleSlowStyles...)
	out = append(out, grayscaleFastStyles...)
	out = append(out, colorSlowStyles...)
	out = append(out, colorFastStyles...)
	return out
}
