package styles

import "github.com/databeast/cyberhud/display/style"

// Dashboard style declarations for this resolution family.
//
// Each entry is a hand-tweakable declaration over the core layouts in core.go:
// adjust its Params layout selector or attach a bespoke BuildFn to tune the
// dashboard's look for this specific display.

var ColorStyle80x160 = styleDef{
	name: "color-80x160",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.ColorFast},
	p:    Params{Layout: LayoutCompactPortrait},
}

var GrayscaleFast80x160 = styleDef{
	name: "grayscale-fast-80x160",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildGrayscaleFast80x160},
}

var GrayscaleFast160x80 = styleDef{
	name: "grayscale-fast-160x80",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildGrayscaleFast80x160},
}

var ColorStyle160x80 = styleDef{
	name: "color-160x80",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.ColorFast},
	p:    Params{BuildFn: buildColorStyle160x80},
}

var MonoSlow160x80Style = styleDef{
	name: "mono-slow-160x80",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.MonoSlow},
}

var MonoSlow80x160Style = styleDef{
	name: "mono-slow-80x160",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.MonoSlow},
}

var MonoFast160x80Style = styleDef{
	name: "mono-fast-160x80",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.MonoFast},
}

var MonoFast80x160Style = styleDef{
	name: "mono-fast-80x160",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.MonoFast},
}

var GrayscaleSlow80x160Style = styleDef{
	name: "grayscale-slow-80x160",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.GrayscaleSlow},
	p:    Params{Layout: LayoutGrayscaleDashboard},
}

var GrayscaleSlow160x80Style = styleDef{
	name: "grayscale-slow-160x80",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.GrayscaleSlow},
	p:    Params{Layout: LayoutGrayscaleDashboard},
}

var ColorSlow80x160Style = styleDef{
	name: "color-slow-80x160",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.ColorSlow},
	p:    Params{Layout: LayoutColorSkeleton},
}

var ColorSlow160x80Style = styleDef{
	name: "color-slow-160x80",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.ColorSlow},
	p:    Params{Layout: LayoutColorSkeleton},
}
