package styles

import (
	"github.com/databeast/cyberhud/display/style"
)

// Clock style declarations for 400x300 panels.
//
// Each entry is a hand-tweakable declaration over the core layouts in
// core.go: adjust its Params (font tier, row toggles, Fast/Color layout
// selection) or attach a bespoke BuildFn to tune the clock's look for
// this specific resolution and PPI.

var GrayscaleSlow400x300Style = styleDef{
	name: "grayscale-slow-400x300",
	reqs: style.SurfaceRequirements{MinWidth: 400, MinHeight: 300, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow300x400Style = styleDef{
	name: "grayscale-slow-300x400",
	reqs: style.SurfaceRequirements{MinWidth: 300, MinHeight: 400, Capability: style.GrayscaleSlow},
}

var ColorSlow400x300Style = styleDef{
	name: "color-slow-400x300",
	reqs: style.SurfaceRequirements{MinWidth: 400, MinHeight: 300, Capability: style.ColorSlow},
}

var ColorSlow300x400Style = styleDef{
	name: "color-slow-300x400",
	reqs: style.SurfaceRequirements{MinWidth: 300, MinHeight: 400, Capability: style.ColorSlow},
}
