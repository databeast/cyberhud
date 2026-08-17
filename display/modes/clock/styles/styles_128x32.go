package styles

import (
	"github.com/databeast/cyberhud/display/style"
)

// Clock style declarations for 128x32 panels.
//
// Each entry is a hand-tweakable declaration over the core layouts in
// core.go: adjust its Params (font tier, row toggles, Fast/Color layout
// selection) or attach a bespoke BuildFn to tune the clock's look for
// this specific resolution and PPI.

var MonoSlow128x32Style = styleDef{
	name: "mono-slow-128x32",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 32, Capability: style.MonoSlow},
}

var MonoSlow32x128Style = styleDef{
	name: "mono-slow-32x128",
	reqs: style.SurfaceRequirements{MinWidth: 32, MinHeight: 128, Capability: style.MonoSlow},
}

var MonoSmall128x32Style = styleDef{
	name: "mono-128x32",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 32, Capability: style.MonoFast},
}

var MonoSmall32x128Style = styleDef{
	name: "mono-32x128",
	reqs: style.SurfaceRequirements{MinWidth: 32, MinHeight: 128, Capability: style.MonoFast},
}

var GrayscaleSlow128x32Style = styleDef{
	name: "grayscale-slow-128x32",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 32, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow32x128Style = styleDef{
	name: "grayscale-slow-32x128",
	reqs: style.SurfaceRequirements{MinWidth: 32, MinHeight: 128, Capability: style.GrayscaleSlow},
}

var ColorSlow128x32Style = styleDef{
	name: "color-slow-128x32",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 32, Capability: style.ColorSlow},
}

var ColorSlow32x128Style = styleDef{
	name: "color-slow-32x128",
	reqs: style.SurfaceRequirements{MinWidth: 32, MinHeight: 128, Capability: style.ColorSlow},
}
