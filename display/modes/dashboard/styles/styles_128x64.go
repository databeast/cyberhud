package styles

import "github.com/databeast/cyberhud/display/style"

// Dashboard style declarations for this resolution family.
//
// Each entry is a hand-tweakable declaration over the core layouts in core.go:
// adjust its Params layout selector or attach a bespoke BuildFn to tune the
// dashboard's look for this specific display.

var MonoStyle64x128 = styleDef{
	name: "mono-64x128",
	reqs: style.SurfaceRequirements{MinWidth: 64, MinHeight: 128, Capability: style.MonoFast},
	p:    Params{BuildFn: buildMonoStyle64x128},
}

var MonoStyle128x64 = styleDef{
	name: "mono-128x64",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 64, Capability: style.MonoFast},
	p:    Params{BuildFn: buildMonoStyle64x128},
}

var MonoSlow128x64Style = styleDef{
	name: "mono-slow-128x64",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 64, Capability: style.MonoSlow},
}

var MonoSlow64x128Style = styleDef{
	name: "mono-slow-64x128",
	reqs: style.SurfaceRequirements{MinWidth: 64, MinHeight: 128, Capability: style.MonoSlow},
}

var GrayscaleSlow128x64Style = styleDef{
	name: "grayscale-slow-128x64",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 64, Capability: style.GrayscaleSlow},
	p:    Params{Layout: LayoutGrayscaleDashboard},
}

var GrayscaleSlow64x128Style = styleDef{
	name: "grayscale-slow-64x128",
	reqs: style.SurfaceRequirements{MinWidth: 64, MinHeight: 128, Capability: style.GrayscaleSlow},
	p:    Params{Layout: LayoutGrayscaleDashboard},
}

var GrayscaleFast64x128Style = styleDef{
	name: "grayscale-fast-64x128",
	reqs: style.SurfaceRequirements{MinWidth: 64, MinHeight: 128, Capability: style.GrayscaleFast},
	p:    Params{Layout: LayoutGrayscaleDashboard},
}

var GrayscaleFast128x64Style = styleDef{
	name: "grayscale-fast-128x64",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 64, Capability: style.GrayscaleFast},
	p:    Params{Layout: LayoutGrayscaleDashboard},
}

var ColorSlow64x128Style = styleDef{
	name: "color-slow-64x128",
	reqs: style.SurfaceRequirements{MinWidth: 64, MinHeight: 128, Capability: style.ColorSlow},
	p:    Params{Layout: LayoutColorSkeleton},
}

var ColorSlow128x64Style = styleDef{
	name: "color-slow-128x64",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 64, Capability: style.ColorSlow},
	p:    Params{Layout: LayoutColorSkeleton},
}

var ColorFast64x128Style = styleDef{
	name: "color-fast-64x128",
	reqs: style.SurfaceRequirements{MinWidth: 64, MinHeight: 128, Capability: style.ColorFast},
	p:    Params{Layout: LayoutColorFastSkeleton},
}

var ColorFast128x64Style = styleDef{
	name: "color-fast-128x64",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 64, Capability: style.ColorFast},
	p:    Params{Layout: LayoutColorFastSkeleton},
}
