package styles

import (
	"github.com/databeast/cyberhud/display/style"
)

// System style declarations for 64x128 panels (portrait).
// Variant set: skeleton (adaptiveBuild) + default, compact, cores, top.

var MonoSlow64x128Style = def{
	name: "mono-slow-64x128",
	reqs: style.SurfaceRequirements{MinWidth: 64, MinHeight: 128, Capability: style.MonoSlow},
}

var MonoSlow64x128DefaultStyle = def{
	name: "mono-slow-64x128-default",
	reqs: style.SurfaceRequirements{MinWidth: 64, MinHeight: 128, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildDefaultStyle},
}

var MonoSlow64x128CompactStyle = def{
	name: "mono-slow-64x128-compact",
	reqs: style.SurfaceRequirements{MinWidth: 64, MinHeight: 128, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildCompactStyle},
}

var MonoSlow64x128CoresStyle = def{
	name: "mono-slow-64x128-cores",
	reqs: style.SurfaceRequirements{MinWidth: 64, MinHeight: 128, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildCoresStyle},
}

var MonoSlow64x128TopStyle = def{
	name: "mono-slow-64x128-top",
	reqs: style.SurfaceRequirements{MinWidth: 64, MinHeight: 128, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildTopStyle},
}

var MonoFast64x128Style = def{
	name: "mono-64x128",
	reqs: style.SurfaceRequirements{MinWidth: 64, MinHeight: 128, Capability: style.MonoFast},
}

var MonoFast64x128DefaultStyle = def{
	name: "mono-64x128-default",
	reqs: style.SurfaceRequirements{MinWidth: 64, MinHeight: 128, Capability: style.MonoFast},
	p:    Params{BuildFn: buildDefaultStyle},
}

var MonoFast64x128CompactStyle = def{
	name: "mono-64x128-compact",
	reqs: style.SurfaceRequirements{MinWidth: 64, MinHeight: 128, Capability: style.MonoFast},
	p:    Params{BuildFn: buildCompactStyle},
}

var MonoFast64x128CoresStyle = def{
	name: "mono-64x128-cores",
	reqs: style.SurfaceRequirements{MinWidth: 64, MinHeight: 128, Capability: style.MonoFast},
	p:    Params{BuildFn: buildCoresStyle},
}

var MonoFast64x128TopStyle = def{
	name: "mono-64x128-top",
	reqs: style.SurfaceRequirements{MinWidth: 64, MinHeight: 128, Capability: style.MonoFast},
	p:    Params{BuildFn: buildTopStyle},
}

var GrayscaleSlow64x128Style = def{
	name: "grayscale-slow-64x128",
	reqs: style.SurfaceRequirements{MinWidth: 64, MinHeight: 128, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow64x128DefaultStyle = def{
	name: "grayscale-slow-64x128-default",
	reqs: style.SurfaceRequirements{MinWidth: 64, MinHeight: 128, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildDefaultStyle},
}

var GrayscaleSlow64x128CompactStyle = def{
	name: "grayscale-slow-64x128-compact",
	reqs: style.SurfaceRequirements{MinWidth: 64, MinHeight: 128, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildCompactStyle},
}

var GrayscaleSlow64x128CoresStyle = def{
	name: "grayscale-slow-64x128-cores",
	reqs: style.SurfaceRequirements{MinWidth: 64, MinHeight: 128, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildCoresStyle},
}

var GrayscaleSlow64x128TopStyle = def{
	name: "grayscale-slow-64x128-top",
	reqs: style.SurfaceRequirements{MinWidth: 64, MinHeight: 128, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildTopStyle},
}

var ColorSlow64x128Style = def{
	name: "color-slow-64x128",
	reqs: style.SurfaceRequirements{MinWidth: 64, MinHeight: 128, Capability: style.ColorSlow},
}

var ColorSlow64x128DefaultStyle = def{
	name: "color-slow-64x128-default",
	reqs: style.SurfaceRequirements{MinWidth: 64, MinHeight: 128, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildDefaultStyle},
}

var ColorSlow64x128CompactStyle = def{
	name: "color-slow-64x128-compact",
	reqs: style.SurfaceRequirements{MinWidth: 64, MinHeight: 128, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildCompactStyle},
}

var ColorSlow64x128CoresStyle = def{
	name: "color-slow-64x128-cores",
	reqs: style.SurfaceRequirements{MinWidth: 64, MinHeight: 128, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildCoresStyle},
}

var ColorSlow64x128TopStyle = def{
	name: "color-slow-64x128-top",
	reqs: style.SurfaceRequirements{MinWidth: 64, MinHeight: 128, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildTopStyle},
}

