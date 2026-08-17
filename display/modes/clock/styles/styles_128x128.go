package styles

import (
	"github.com/databeast/cyberhud/display/style"
)

// Clock style declarations for 128x128 panels.
//
// Each entry is a hand-tweakable declaration over the core layouts in
// core.go: adjust its Params (font tier, row toggles, Fast/Color layout
// selection) or attach a bespoke BuildFn to tune the clock's look for
// this specific resolution and PPI.

var MonoSlow128x128Style = styleDef{
	name: "mono-slow-128x128",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 128, Capability: style.MonoSlow},
}

var MonoSmall128x128Style = styleDef{
	name: "mono-128x128",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 128, Capability: style.MonoFast},
}

var GrayscaleSlow128x128Style = styleDef{
	name: "grayscale-slow-128x128",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 128, Capability: style.GrayscaleSlow},
}

var GrayscaleFast128x128Style = styleDef{
	name: "grayscale-fast-128x128",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 128, Capability: style.GrayscaleFast},
	p:    Params{Fast: true},
}

var ColorSlow128x128Style = styleDef{
	name: "color-slow-128x128",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 128, Capability: style.ColorSlow},
}

var ColorSmall128x128Style = styleDef{
	name: "color-128x128",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 128, Capability: style.ColorFast},
	p:    Params{Fast: true, Color: true},
}
