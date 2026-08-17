package styles

import "github.com/databeast/cyberhud/display/style"

// Dashboard style declarations for this resolution family.
//
// Each entry is a hand-tweakable declaration over the core layouts in core.go:
// adjust its Params layout selector or attach a bespoke BuildFn to tune the
// dashboard's look for this specific display.

var ColorStyle240x240 = styleDef{
	name: "color-240x240",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 240, Capability: style.ColorFast},
	p:    Params{BuildFn: buildColorStyle240x240},
}

var GrayscaleFast240x240 = styleDef{
	name: "grayscale-fast-240x240",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 240, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildGrayscaleFast240x135},
}

var MonoSlow240x240Style = styleDef{
	name: "mono-slow-240x240",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 240, Capability: style.MonoSlow},
}

var MonoFast240x240Style = styleDef{
	name: "mono-fast-240x240",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 240, Capability: style.MonoFast},
}

var GrayscaleSlow240x240Style = styleDef{
	name: "grayscale-slow-240x240",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 240, Capability: style.GrayscaleSlow},
	p:    Params{Layout: LayoutGrayscaleDashboard},
}

var ColorSlow240x240Style = styleDef{
	name: "color-slow-240x240",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 240, Capability: style.ColorSlow},
	p:    Params{Layout: LayoutColorSkeleton},
}
