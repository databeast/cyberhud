package styles

import (
	"github.com/databeast/cyberhud/display/style"
)

// Clock style declarations for 160x128 panels.
//
// Each entry is a hand-tweakable declaration over the core layouts in
// core.go: adjust its Params (font tier, row toggles, Fast/Color layout
// selection) or attach a bespoke BuildFn to tune the clock's look for
// this specific resolution and PPI.

var MonoSlow160x128Style = styleDef{
	name: "mono-slow-160x128",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 128, Capability: style.MonoSlow},
}

var MonoSlow128x160Style = styleDef{
	name: "mono-slow-128x160",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 160, Capability: style.MonoSlow},
}

var GrayscaleSlow128x160Style = styleDef{
	name: "grayscale-slow-128x160",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 160, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow160x128Style = styleDef{
	name: "grayscale-slow-160x128",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 128, Capability: style.GrayscaleSlow},
}

var GrayscaleFast160x128Style = styleDef{
	name: "grayscale-fast-160x128",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 128, Capability: style.GrayscaleFast},
}

var GrayscaleFast128x160Style = styleDef{
	name: "grayscale-fast-128x160",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 160, Capability: style.GrayscaleFast},
}

var ColorSlow160x128Style = styleDef{
	name: "color-slow-160x128",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 128, Capability: style.ColorSlow},
}

var ColorSlow128x160Style = styleDef{
	name: "color-slow-128x160",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 160, Capability: style.ColorSlow},
}

var ColorSmall160x128Style = styleDef{
	name: "color-160x128",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 128, Capability: style.ColorFast},
	p:    Params{Fast: true, Color: true},
}

var ColorSmall128x160Style = styleDef{
	name: "color-128x160",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 160, Capability: style.ColorFast},
	p:    Params{Fast: true, Color: true},
}
