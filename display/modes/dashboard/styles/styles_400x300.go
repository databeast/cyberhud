package styles

import "github.com/databeast/cyberhud/display/style"

// Dashboard style declarations for this resolution family.
//
// Each entry is a hand-tweakable declaration over the core layouts in core.go:
// adjust its Params layout selector or attach a bespoke BuildFn to tune the
// dashboard's look for this specific display.

var EinkStyle300x400 = styleDef{
	name: "eink-300x400",
	reqs: style.SurfaceRequirements{MinWidth: 300, MinHeight: 400, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildEinkStyle200x200},
}

var EinkStyle400x300 = styleDef{
	name: "eink-400x300",
	reqs: style.SurfaceRequirements{MinWidth: 400, MinHeight: 300, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildEinkStyle200x200},
}

var MonoFast300x400Style = styleDef{
	name: "mono-fast-300x400",
	reqs: style.SurfaceRequirements{MinWidth: 300, MinHeight: 400, Capability: style.MonoFast},
}

var MonoFast400x300Style = styleDef{
	name: "mono-fast-400x300",
	reqs: style.SurfaceRequirements{MinWidth: 400, MinHeight: 300, Capability: style.MonoFast},
}

var GrayscaleSlow400x300Style = styleDef{
	name: "grayscale-slow-400x300",
	reqs: style.SurfaceRequirements{MinWidth: 400, MinHeight: 300, Capability: style.GrayscaleSlow},
	p:    Params{Layout: LayoutGrayscaleDashboard},
}

var GrayscaleSlow300x400Style = styleDef{
	name: "grayscale-slow-300x400",
	reqs: style.SurfaceRequirements{MinWidth: 300, MinHeight: 400, Capability: style.GrayscaleSlow},
	p:    Params{Layout: LayoutGrayscaleDashboard},
}

var GrayscaleFast400x300Style = styleDef{
	name: "grayscale-fast-400x300",
	reqs: style.SurfaceRequirements{MinWidth: 400, MinHeight: 300, Capability: style.GrayscaleFast},
	p:    Params{Layout: LayoutGrayscaleDashboard},
}

var GrayscaleFast300x400Style = styleDef{
	name: "grayscale-fast-300x400",
	reqs: style.SurfaceRequirements{MinWidth: 300, MinHeight: 400, Capability: style.GrayscaleFast},
	p:    Params{Layout: LayoutGrayscaleDashboard},
}

var ColorSlow400x300Style = styleDef{
	name: "color-slow-400x300",
	reqs: style.SurfaceRequirements{MinWidth: 400, MinHeight: 300, Capability: style.ColorSlow},
	p:    Params{Layout: LayoutColorSkeleton},
}

var ColorSlow300x400Style = styleDef{
	name: "color-slow-300x400",
	reqs: style.SurfaceRequirements{MinWidth: 300, MinHeight: 400, Capability: style.ColorSlow},
	p:    Params{Layout: LayoutColorSkeleton},
}

var ColorFast400x300Style = styleDef{
	name: "color-fast-400x300",
	reqs: style.SurfaceRequirements{MinWidth: 400, MinHeight: 300, Capability: style.ColorFast},
	p:    Params{Layout: LayoutColorFastSkeleton},
}

var ColorFast300x400Style = styleDef{
	name: "color-fast-300x400",
	reqs: style.SurfaceRequirements{MinWidth: 300, MinHeight: 400, Capability: style.ColorFast},
	p:    Params{Layout: LayoutColorFastSkeleton},
}
