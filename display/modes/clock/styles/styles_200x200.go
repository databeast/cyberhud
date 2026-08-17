package styles

import (
	"github.com/databeast/cyberhud/display/style"
)

// Clock style declarations for 200x200 panels.
//
// Each entry is a hand-tweakable declaration over the core layouts in
// core.go: adjust its Params (font tier, row toggles, Fast/Color layout
// selection) or attach a bespoke BuildFn to tune the clock's look for
// this specific resolution and PPI.

var GrayscaleSlow200x200Style = styleDef{
	name: "grayscale-slow-200x200",
	reqs: style.SurfaceRequirements{MinWidth: 200, MinHeight: 200, Capability: style.GrayscaleSlow},
}

var ColorSlow200x200Style = styleDef{
	name: "color-slow-200x200",
	reqs: style.SurfaceRequirements{MinWidth: 200, MinHeight: 200, Capability: style.ColorSlow},
}
