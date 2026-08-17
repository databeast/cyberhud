package styles

import (
	"github.com/databeast/cyberhud/display/style"
)

// Clock style declarations for 250x122 panels.
//
// Each entry is a hand-tweakable declaration over the core layouts in
// core.go: adjust its Params (font tier, row toggles, Fast/Color layout
// selection) or attach a bespoke BuildFn to tune the clock's look for
// this specific resolution and PPI.

var GrayscaleSlow250x122Style = styleDef{
	name: "grayscale-slow-250x122",
	reqs: style.SurfaceRequirements{MinWidth: 250, MinHeight: 122, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow122x250Style = styleDef{
	name: "grayscale-slow-122x250",
	reqs: style.SurfaceRequirements{MinWidth: 122, MinHeight: 250, Capability: style.GrayscaleSlow},
}

var ColorSlow250x122Style = styleDef{
	name: "color-slow-250x122",
	reqs: style.SurfaceRequirements{MinWidth: 250, MinHeight: 122, Capability: style.ColorSlow},
}

var ColorSlow122x250Style = styleDef{
	name: "color-slow-122x250",
	reqs: style.SurfaceRequirements{MinWidth: 122, MinHeight: 250, Capability: style.ColorSlow},
}
