package styles

import (
	"github.com/databeast/cyberhud/display/style"
)

// Clock style declarations for 800x480 panels.
//
// Each entry is a hand-tweakable declaration over the core layouts in
// core.go: adjust its Params (font tier, row toggles, Fast/Color layout
// selection) or attach a bespoke BuildFn to tune the clock's look for
// this specific resolution and PPI.

var GrayscaleSlow480x800Style = styleDef{
	name: "grayscale-slow-480x800",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.GrayscaleSlow},
}

var GrayscaleFast480x800Style = styleDef{
	name: "grayscale-fast-480x800",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.GrayscaleFast},
	p:    Params{Fast: true},
}

var ColorSlow480x800Style = styleDef{
	name: "color-slow-480x800",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.ColorSlow},
}

var ColorLarge480x800Style = styleDef{
	name: "color-480x800",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.ColorFast},
}
