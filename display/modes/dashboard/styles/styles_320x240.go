package styles

import "github.com/databeast/cyberhud/display/style"

// Dashboard style declarations for this resolution family.
//
// Each entry is a hand-tweakable declaration over the core layouts in core.go:
// adjust its Params layout selector or attach a bespoke BuildFn to tune the
// dashboard's look for this specific display.

var ColorStyle320x240 = styleDef{
	name: "color-320x240",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 240, Capability: style.ColorFast},
	p:    Params{BuildFn: buildColorStyle320x240},
}

var ColorStyle240x320 = styleDef{
	name: "color-240x320",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 320, Capability: style.ColorFast},
	p:    Params{BuildFn: buildColorStyle320x240},
}

var GrayscaleFast240x320 = styleDef{
	name: "grayscale-fast-240x320",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 320, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildGrayscaleFast240x135},
}

var MonoSlow240x320Style = styleDef{
	name: "mono-slow-240x320",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 320, Capability: style.MonoSlow},
}

var MonoSlow320x240Style = styleDef{
	name: "mono-slow-320x240",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 240, Capability: style.MonoSlow},
}

var MonoFast240x320Style = styleDef{
	name: "mono-fast-240x320",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 320, Capability: style.MonoFast},
}

var MonoFast320x240Style = styleDef{
	name: "mono-fast-320x240",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 240, Capability: style.MonoFast},
}

var GrayscaleSlow320x240Style = styleDef{
	name: "grayscale-slow-320x240",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 240, Capability: style.GrayscaleSlow},
	p:    Params{Layout: LayoutGrayscaleDashboard},
}

var GrayscaleSlow240x320Style = styleDef{
	name: "grayscale-slow-240x320",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 320, Capability: style.GrayscaleSlow},
	p:    Params{Layout: LayoutGrayscaleDashboard},
}

var ColorSlow320x240Style = styleDef{
	name: "color-slow-320x240",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 240, Capability: style.ColorSlow},
	p:    Params{Layout: LayoutColorSkeleton},
}

var ColorSlow240x320Style = styleDef{
	name: "color-slow-240x320",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 320, Capability: style.ColorSlow},
	p:    Params{Layout: LayoutColorSkeleton},
}
