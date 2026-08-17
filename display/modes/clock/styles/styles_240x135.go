package styles

import (
	"github.com/databeast/cyberhud/display/style"
)

// Clock style declarations for 240x135 panels.
//
// Each entry is a hand-tweakable declaration over the core layouts in
// core.go: adjust its Params (font tier, row toggles, Fast/Color layout
// selection) or attach a bespoke BuildFn to tune the clock's look for
// this specific resolution and PPI.

var MonoSlow240x135Style = styleDef{
	name: "mono-slow-240x135",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.MonoSlow},
}

var MonoSlow135x240Style = styleDef{
	name: "mono-slow-135x240",
	reqs: style.SurfaceRequirements{MinWidth: 135, MinHeight: 240, Capability: style.MonoSlow},
}

var GrayscaleSlow135x240Style = styleDef{
	name: "grayscale-slow-135x240",
	reqs: style.SurfaceRequirements{MinWidth: 135, MinHeight: 240, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow240x135Style = styleDef{
	name: "grayscale-slow-240x135",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.GrayscaleSlow},
}

var GrayscaleFast240x135Style = styleDef{
	name: "grayscale-fast-240x135",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.GrayscaleFast},
	p:    Params{Fast: true},
}

var GrayscaleFast135x240Style = styleDef{
	name: "grayscale-fast-135x240",
	reqs: style.SurfaceRequirements{MinWidth: 135, MinHeight: 240, Capability: style.GrayscaleFast},
	p:    Params{Fast: true},
}

var ColorSlow240x135Style = styleDef{
	name: "color-slow-240x135",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.ColorSlow},
}

var ColorSlow135x240Style = styleDef{
	name: "color-slow-135x240",
	reqs: style.SurfaceRequirements{MinWidth: 135, MinHeight: 240, Capability: style.ColorSlow},
}

var ColorMedium240x135Style = styleDef{
	name: "color-240x135",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.ColorFast},
	p:    Params{Fast: true, Color: true},
}

var ColorMedium135x240Style = styleDef{
	name: "color-135x240",
	reqs: style.SurfaceRequirements{MinWidth: 135, MinHeight: 240, Capability: style.ColorFast},
	p:    Params{Fast: true, Color: true},
}
