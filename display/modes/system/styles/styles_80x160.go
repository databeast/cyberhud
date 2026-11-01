package styles

import (
	"github.com/databeast/cyberhud/display/style"
)

// System style declarations for 80x160 panels (portrait).
// Variant set: skeleton (adaptiveBuild) + default, compact, cores, top.

var MonoSlow80x160Style = def{
	name: "mono-slow-80x160",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.MonoSlow},
}

var MonoSlow80x160DefaultStyle = def{
	name: "mono-slow-80x160-default",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildDefaultStyle},
}

var MonoSlow80x160CompactStyle = def{
	name: "mono-slow-80x160-compact",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildCompactStyle},
}

var MonoSlow80x160CoresStyle = def{
	name: "mono-slow-80x160-cores",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildCoresStyle},
}

var MonoSlow80x160TopStyle = def{
	name: "mono-slow-80x160-top",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildTopStyle},
}

var GrayscaleSlow80x160Style = def{
	name: "grayscale-slow-80x160",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow80x160DefaultStyle = def{
	name: "grayscale-slow-80x160-default",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildDefaultStyle},
}

var GrayscaleSlow80x160CompactStyle = def{
	name: "grayscale-slow-80x160-compact",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildCompactStyle},
}

var GrayscaleSlow80x160CoresStyle = def{
	name: "grayscale-slow-80x160-cores",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildCoresStyle},
}

var GrayscaleSlow80x160TopStyle = def{
	name: "grayscale-slow-80x160-top",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildTopStyle},
}

var GrayscaleFast80x160Style = def{
	name: "grayscale-fast-80x160",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.GrayscaleFast},
}

var GrayscaleFast80x160DefaultStyle = def{
	name: "grayscale-fast-80x160-default",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildDefaultStyle},
}

var GrayscaleFast80x160CompactStyle = def{
	name: "grayscale-fast-80x160-compact",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildCompactStyle},
}

var GrayscaleFast80x160CoresStyle = def{
	name: "grayscale-fast-80x160-cores",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildCoresStyle},
}

var GrayscaleFast80x160TopStyle = def{
	name: "grayscale-fast-80x160-top",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildTopStyle},
}

var ColorSlow80x160Style = def{
	name: "color-slow-80x160",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.ColorSlow},
}

var ColorSlow80x160DefaultStyle = def{
	name: "color-slow-80x160-default",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildDefaultStyle},
}

var ColorSlow80x160CompactStyle = def{
	name: "color-slow-80x160-compact",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildCompactStyle},
}

var ColorSlow80x160CoresStyle = def{
	name: "color-slow-80x160-cores",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildCoresStyle},
}

var ColorSlow80x160TopStyle = def{
	name: "color-slow-80x160-top",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildTopStyle},
}

var ColorFast80x160Style = def{
	name: "color-80x160",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.ColorFast},
}

var ColorFast80x160DefaultStyle = def{
	name: "color-80x160-default",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.ColorFast},
	p:    Params{BuildFn: buildDefaultStyle},
}

var ColorFast80x160CompactStyle = def{
	name: "color-80x160-compact",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.ColorFast},
	p:    Params{BuildFn: buildCompactStyle},
}

var ColorFast80x160CoresStyle = def{
	name: "color-80x160-cores",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.ColorFast},
	p:    Params{BuildFn: buildCoresStyle},
}

var ColorFast80x160TopStyle = def{
	name: "color-80x160-top",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.ColorFast},
	p:    Params{BuildFn: buildTopStyle},
}

