package styles

import "github.com/databeast/cyberhud/display/style"

// Dashboard style declarations for this resolution family.
//
// Each entry is a hand-tweakable declaration over the core layouts in core.go:
// adjust its Params layout selector or attach a bespoke BuildFn to tune the
// dashboard's look for this specific display.

var GrayscaleFast800x480 = styleDef{
	name: "grayscale-fast-800x480",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildGrayscaleFast240x135},
}

var EinkStyle800x480 = styleDef{
	name: "eink-800x480",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildEinkStyle800x480},
}

var GrayscaleFast480x800 = styleDef{
	name: "grayscale-fast-480x800",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildGrayscaleFast240x135},
}

var EinkStyle480x800 = styleDef{
	name: "eink-480x800",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildEinkStyle200x200},
}

var ColorStyle480x800 = styleDef{
	name: "color-480x800",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.ColorFast},
	p:    Params{BuildFn: buildColorStyle320x240},
}

var ColorStyle800x480 = styleDef{
	name: "color-800x480",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.ColorFast},
	p:    Params{BuildFn: buildColorStyle320x240},
}

var MonoFast480x800Style = styleDef{
	name: "mono-fast-480x800",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.MonoFast},
}

var MonoFast800x480Style = styleDef{
	name: "mono-fast-800x480",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.MonoFast},
}

var GrayscaleSlow800x480Style = styleDef{
	name: "grayscale-slow-800x480",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.GrayscaleSlow},
	p:    Params{Layout: LayoutGrayscaleDashboard},
}

var GrayscaleSlow480x800Style = styleDef{
	name: "grayscale-slow-480x800",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.GrayscaleSlow},
	p:    Params{Layout: LayoutGrayscaleDashboard},
}

var ColorSlow800x480Style = styleDef{
	name: "color-slow-800x480",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.ColorSlow},
	p:    Params{Layout: LayoutColorSkeleton},
}

var ColorSlow480x800Style = styleDef{
	name: "color-slow-480x800",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.ColorSlow},
	p:    Params{Layout: LayoutColorSkeleton},
}
