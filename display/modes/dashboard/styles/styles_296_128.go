package styles

import "github.com/databeast/cyberhud/display/style"

// Dashboard style declarations for this resolution family.
//
// Each entry is a hand-tweakable declaration over the core layouts in core.go:
// adjust its Params layout selector or attach a bespoke BuildFn to tune the
// dashboard's look for this specific display.

var EinkStyle128x296 = styleDef{
	name: "eink-128x296",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 296, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildEinkStyle200x200},
}

var EinkStyle296x128 = styleDef{
	name: "eink-296x128",
	reqs: style.SurfaceRequirements{MinWidth: 296, MinHeight: 128, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildEinkStyle200x200},
}

var MonoFast128x296Style = styleDef{
	name: "mono-fast-128x296",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 296, Capability: style.MonoFast},
}

var MonoFast296x128Style = styleDef{
	name: "mono-fast-296x128",
	reqs: style.SurfaceRequirements{MinWidth: 296, MinHeight: 128, Capability: style.MonoFast},
}

var GrayscaleSlow128x296Style = styleDef{
	name: "grayscale-slow-128x296",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 296, Capability: style.GrayscaleSlow},
	p:    Params{Layout: LayoutGrayscaleDashboard},
}

var GrayscaleSlow296x128Style = styleDef{
	name: "grayscale-slow-296x128",
	reqs: style.SurfaceRequirements{MinWidth: 296, MinHeight: 128, Capability: style.GrayscaleSlow},
	p:    Params{Layout: LayoutGrayscaleDashboard},
}

var GrayscaleFast296x128Style = styleDef{
	name: "grayscale-fast-296x128",
	reqs: style.SurfaceRequirements{MinWidth: 296, MinHeight: 128, Capability: style.GrayscaleFast},
	p:    Params{Layout: LayoutGrayscaleDashboard},
}

var GrayscaleFast128x296Style = styleDef{
	name: "grayscale-fast-128x296",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 296, Capability: style.GrayscaleFast},
	p:    Params{Layout: LayoutGrayscaleDashboard},
}

var ColorSlow296x128Style = styleDef{
	name: "color-slow-296x128",
	reqs: style.SurfaceRequirements{MinWidth: 296, MinHeight: 128, Capability: style.ColorSlow},
	p:    Params{Layout: LayoutColorSkeleton},
}

var ColorSlow128x296Style = styleDef{
	name: "color-slow-128x296",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 296, Capability: style.ColorSlow},
	p:    Params{Layout: LayoutColorSkeleton},
}

var ColorFast296x128Style = styleDef{
	name: "color-fast-296x128",
	reqs: style.SurfaceRequirements{MinWidth: 296, MinHeight: 128, Capability: style.ColorFast},
	p:    Params{Layout: LayoutColorFastSkeleton},
}

var ColorFast128x296Style = styleDef{
	name: "color-fast-128x296",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 296, Capability: style.ColorFast},
	p:    Params{Layout: LayoutColorFastSkeleton},
}
