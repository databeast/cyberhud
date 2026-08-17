package styles

import (
	"github.com/databeast/cyberhud/display/style"
)

// Clock style declarations for 160x80 panels.
//
// Each entry is a hand-tweakable declaration over the core layouts in
// core.go: adjust its Params (font tier, row toggles, Fast/Color layout
// selection) or attach a bespoke BuildFn to tune the clock's look for
// this specific resolution and PPI.

var MonoSlow160x80Style = styleDef{
	name: "mono-slow-160x80",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.MonoSlow},
}

var MonoSlow80x160Style = styleDef{
	name: "mono-slow-80x160",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.MonoSlow},
}

var GrayscaleSlow160x80Style = styleDef{
	name: "grayscale-slow-160x80",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow80x160Style = styleDef{
	name: "grayscale-slow-80x160",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.GrayscaleSlow},
}

var GrayscaleFast160x80Style = styleDef{
	name: "grayscale-fast-160x80",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.GrayscaleFast},
	p:    Params{Fast: true},
}

var GrayscaleFast80x160Style = styleDef{
	name: "grayscale-fast-80x160",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.GrayscaleFast},
	p:    Params{Fast: true},
}

var ColorSlow160x80Style = styleDef{
	name: "color-slow-160x80",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.ColorSlow},
}

var ColorSlow80x160Style = styleDef{
	name: "color-slow-80x160",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.ColorSlow},
}

var ColorSmall160x80Style = styleDef{
	name: "color-160x80",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.ColorFast},
	p:    Params{Fast: true, Color: true},
}

// ColorSmall80x160Style targets the 80×160 (portrait) color TFT panel.
var ColorSmall80x160Style = styleDef{
	name: "color-80x160",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.ColorFast},
}
