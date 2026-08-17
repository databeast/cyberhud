package styles

import "github.com/databeast/cyberhud/display/style"

// Dashboard style declarations for this resolution family.
//
// Each entry is a hand-tweakable declaration over the core layouts in core.go:
// adjust its Params layout selector or attach a bespoke BuildFn to tune the
// dashboard's look for this specific display.

var ColorStyle320x480 = styleDef{
	name: "color-320x480",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 480, Capability: style.ColorFast},
	p:    Params{BuildFn: buildColorStyle320x240},
}

var ColorStyle480x320 = styleDef{
	name: "color-480x320",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 320, Capability: style.ColorFast},
	p:    Params{BuildFn: buildColorStyle320x240},
}

var GrayscaleFast320x240 = styleDef{
	name: "grayscale-fast-320x240",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 240, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildGrayscaleFast240x135},
}

var GrayscaleFast320x480 = styleDef{
	name: "grayscale-fast-320x480",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 480, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildGrayscaleFast240x135},
}

var GrayscaleFast480x320 = styleDef{
	name: "grayscale-fast-480x320",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 320, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildGrayscaleFast240x135},
}

var MonoSlow320x480Style = styleDef{
	name: "mono-slow-320x480",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 480, Capability: style.MonoSlow},
}

var MonoSlow480x320Style = styleDef{
	name: "mono-slow-480x320",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 320, Capability: style.MonoSlow},
}

var MonoFast320x480Style = styleDef{
	name: "mono-fast-320x480",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 480, Capability: style.MonoFast},
}

var MonoFast480x320Style = styleDef{
	name: "mono-fast-480x320",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 320, Capability: style.MonoFast},
}

var GrayscaleSlow480x320Style = styleDef{
	name: "grayscale-slow-480x320",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 320, Capability: style.GrayscaleSlow},
	p:    Params{Layout: LayoutGrayscaleDashboard},
}

var GrayscaleSlow320x480Style = styleDef{
	name: "grayscale-slow-320x480",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 480, Capability: style.GrayscaleSlow},
	p:    Params{Layout: LayoutGrayscaleDashboard},
}

var ColorSlow480x320Style = styleDef{
	name: "color-slow-480x320",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 320, Capability: style.ColorSlow},
	p:    Params{Layout: LayoutColorSkeleton},
}

var ColorSlow320x480Style = styleDef{
	name: "color-slow-320x480",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 480, Capability: style.ColorSlow},
	p:    Params{Layout: LayoutColorSkeleton},
}
