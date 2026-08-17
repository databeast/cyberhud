package styles

import "github.com/databeast/cyberhud/display/style"

// Dashboard style declarations for this resolution family.
//
// Each entry is a hand-tweakable declaration over the core layouts in core.go:
// adjust its Params layout selector or attach a bespoke BuildFn to tune the
// dashboard's look for this specific display.

var ColorStyle135x240 = styleDef{
	name: "color-135x240",
	reqs: style.SurfaceRequirements{MinWidth: 135, MinHeight: 240, Capability: style.ColorFast},
	p:    Params{BuildFn: buildColorStyle135x240},
}

var ColorStyle240x135 = styleDef{
	name: "color-240x135",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.ColorFast},
	p:    Params{BuildFn: buildColorStyle160x128},
}

var GrayscaleFast240x135 = styleDef{
	name: "grayscale-fast-240x135",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildGrayscaleFast240x135},
}

var GrayscaleFast135x240 = styleDef{
	name: "grayscale-fast-135x240",
	reqs: style.SurfaceRequirements{MinWidth: 135, MinHeight: 240, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildGrayscaleFast240x135},
}

var MonoSlow240x135Style = styleDef{
	name: "mono-slow-240x135",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.MonoSlow},
}

var MonoSlow135x240Style = styleDef{
	name: "mono-slow-135x240",
	reqs: style.SurfaceRequirements{MinWidth: 135, MinHeight: 240, Capability: style.MonoSlow},
}

var MonoFast240x135Style = styleDef{
	name: "mono-fast-240x135",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.MonoFast},
}

var MonoFast135x240Style = styleDef{
	name: "mono-fast-135x240",
	reqs: style.SurfaceRequirements{MinWidth: 135, MinHeight: 240, Capability: style.MonoFast},
}

var GrayscaleSlow135x240Style = styleDef{
	name: "grayscale-slow-135x240",
	reqs: style.SurfaceRequirements{MinWidth: 135, MinHeight: 240, Capability: style.GrayscaleSlow},
	p:    Params{Layout: LayoutGrayscaleDashboard},
}

var GrayscaleSlow240x135Style = styleDef{
	name: "grayscale-slow-240x135",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.GrayscaleSlow},
	p:    Params{Layout: LayoutGrayscaleDashboard},
}

var ColorSlow135x240Style = styleDef{
	name: "color-slow-135x240",
	reqs: style.SurfaceRequirements{MinWidth: 135, MinHeight: 240, Capability: style.ColorSlow},
	p:    Params{Layout: LayoutColorSkeleton},
}

var ColorSlow240x135Style = styleDef{
	name: "color-slow-240x135",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.ColorSlow},
	p:    Params{Layout: LayoutColorSkeleton},
}
