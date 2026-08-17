package styles

import "github.com/databeast/cyberhud/display/style"

// Dashboard style declarations for this resolution family.
//
// Each entry is a hand-tweakable declaration over the core layouts in core.go:
// adjust its Params layout selector or attach a bespoke BuildFn to tune the
// dashboard's look for this specific display.

var EinkStyle122x250 = styleDef{
	name: "eink-122x250",
	reqs: style.SurfaceRequirements{MinWidth: 122, MinHeight: 250, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildEinkStyle200x200},
}

var EinkStyle250x122 = styleDef{
	name: "eink-250x122",
	reqs: style.SurfaceRequirements{MinWidth: 250, MinHeight: 122, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildEinkStyle200x200},
}

var MonoFast122x250Style = styleDef{
	name: "mono-fast-122x250",
	reqs: style.SurfaceRequirements{MinWidth: 122, MinHeight: 250, Capability: style.MonoFast},
}

var MonoFast250x122Style = styleDef{
	name: "mono-fast-250x122",
	reqs: style.SurfaceRequirements{MinWidth: 250, MinHeight: 122, Capability: style.MonoFast},
}

var GrayscaleSlow122x250Style = styleDef{
	name: "grayscale-slow-122x250",
	reqs: style.SurfaceRequirements{MinWidth: 122, MinHeight: 250, Capability: style.GrayscaleSlow},
	p:    Params{Layout: LayoutGrayscaleDashboard},
}

var GrayscaleSlow250x122Style = styleDef{
	name: "grayscale-slow-250x122",
	reqs: style.SurfaceRequirements{MinWidth: 250, MinHeight: 122, Capability: style.GrayscaleSlow},
	p:    Params{Layout: LayoutGrayscaleDashboard},
}

var GrayscaleFast250x122Style = styleDef{
	name: "grayscale-fast-250x122",
	reqs: style.SurfaceRequirements{MinWidth: 250, MinHeight: 122, Capability: style.GrayscaleFast},
	p:    Params{Layout: LayoutGrayscaleDashboard},
}

var GrayscaleFast122x250Style = styleDef{
	name: "grayscale-fast-122x250",
	reqs: style.SurfaceRequirements{MinWidth: 122, MinHeight: 250, Capability: style.GrayscaleFast},
	p:    Params{Layout: LayoutGrayscaleDashboard},
}

var ColorSlow250x122Style = styleDef{
	name: "color-slow-250x122",
	reqs: style.SurfaceRequirements{MinWidth: 250, MinHeight: 122, Capability: style.ColorSlow},
	p:    Params{Layout: LayoutColorSkeleton},
}

var ColorSlow122x250Style = styleDef{
	name: "color-slow-122x250",
	reqs: style.SurfaceRequirements{MinWidth: 122, MinHeight: 250, Capability: style.ColorSlow},
	p:    Params{Layout: LayoutColorSkeleton},
}

var ColorFast250x122Style = styleDef{
	name: "color-fast-250x122",
	reqs: style.SurfaceRequirements{MinWidth: 250, MinHeight: 122, Capability: style.ColorFast},
	p:    Params{Layout: LayoutColorFastSkeleton},
}

var ColorFast122x250Style = styleDef{
	name: "color-fast-122x250",
	reqs: style.SurfaceRequirements{MinWidth: 122, MinHeight: 250, Capability: style.ColorFast},
	p:    Params{Layout: LayoutColorFastSkeleton},
}
