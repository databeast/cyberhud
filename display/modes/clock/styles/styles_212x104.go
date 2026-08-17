package styles

import (
	"github.com/databeast/cyberhud/display/style"
)

// Clock style declarations for 212x104 panels.
//
// Each entry is a hand-tweakable declaration over the core layouts in
// core.go: adjust its Params (font tier, row toggles, Fast/Color layout
// selection) or attach a bespoke BuildFn to tune the clock's look for
// this specific resolution and PPI.

var GrayscaleSlow212x104Style = styleDef{
	name: "grayscale-slow-212x104",
	reqs: style.SurfaceRequirements{MinWidth: 212, MinHeight: 104, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow104x212Style = styleDef{
	name: "grayscale-slow-104x212",
	reqs: style.SurfaceRequirements{MinWidth: 104, MinHeight: 212, Capability: style.GrayscaleSlow},
}

var ColorSlow212x104Style = styleDef{
	name: "color-slow-212x104",
	reqs: style.SurfaceRequirements{MinWidth: 212, MinHeight: 104, Capability: style.ColorSlow},
}

var ColorSlow104x212Style = styleDef{
	name: "color-slow-104x212",
	reqs: style.SurfaceRequirements{MinWidth: 104, MinHeight: 212, Capability: style.ColorSlow},
}
