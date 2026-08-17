package styles

import "github.com/databeast/cyberhud/display/style"

// Dashboard style declarations for this resolution family.
//
// Each entry is a hand-tweakable declaration over the core layouts in core.go:
// adjust its Params layout selector or attach a bespoke BuildFn to tune the
// dashboard's look for this specific display.

var EinkStyle264x176 = styleDef{
	name: "eink-264x176",
	reqs: style.SurfaceRequirements{MinWidth: 264, MinHeight: 176, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildEinkStyle200x200},
}

var EinkStyle176x264 = styleDef{
	name: "eink-176x264",
	reqs: style.SurfaceRequirements{MinWidth: 176, MinHeight: 264, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildEinkStyle200x200},
}

var MonoFast176x264Style = styleDef{
	name: "mono-fast-176x264",
	reqs: style.SurfaceRequirements{MinWidth: 176, MinHeight: 264, Capability: style.MonoFast},
}

var MonoFast264x176Style = styleDef{
	name: "mono-fast-264x176",
	reqs: style.SurfaceRequirements{MinWidth: 264, MinHeight: 176, Capability: style.MonoFast},
}

var GrayscaleSlow264x176Style = styleDef{
	name: "grayscale-slow-264x176",
	reqs: style.SurfaceRequirements{MinWidth: 264, MinHeight: 176, Capability: style.GrayscaleSlow},
	p:    Params{Layout: LayoutGrayscaleDashboard},
}

var GrayscaleSlow176x264Style = styleDef{
	name: "grayscale-slow-176x264",
	reqs: style.SurfaceRequirements{MinWidth: 176, MinHeight: 264, Capability: style.GrayscaleSlow},
	p:    Params{Layout: LayoutGrayscaleDashboard},
}

var GrayscaleFast264x176Style = styleDef{
	name: "grayscale-fast-264x176",
	reqs: style.SurfaceRequirements{MinWidth: 264, MinHeight: 176, Capability: style.GrayscaleFast},
	p:    Params{Layout: LayoutGrayscaleDashboard},
}

var GrayscaleFast176x264Style = styleDef{
	name: "grayscale-fast-176x264",
	reqs: style.SurfaceRequirements{MinWidth: 176, MinHeight: 264, Capability: style.GrayscaleFast},
	p:    Params{Layout: LayoutGrayscaleDashboard},
}

var ColorSlow264x176Style = styleDef{
	name: "color-slow-264x176",
	reqs: style.SurfaceRequirements{MinWidth: 264, MinHeight: 176, Capability: style.ColorSlow},
	p:    Params{Layout: LayoutColorSkeleton},
}

var ColorSlow176x264Style = styleDef{
	name: "color-slow-176x264",
	reqs: style.SurfaceRequirements{MinWidth: 176, MinHeight: 264, Capability: style.ColorSlow},
	p:    Params{Layout: LayoutColorSkeleton},
}

var ColorFast264x176Style = styleDef{
	name: "color-fast-264x176",
	reqs: style.SurfaceRequirements{MinWidth: 264, MinHeight: 176, Capability: style.ColorFast},
	p:    Params{Layout: LayoutColorFastSkeleton},
}

var ColorFast176x264Style = styleDef{
	name: "color-fast-176x264",
	reqs: style.SurfaceRequirements{MinWidth: 176, MinHeight: 264, Capability: style.ColorFast},
	p:    Params{Layout: LayoutColorFastSkeleton},
}
