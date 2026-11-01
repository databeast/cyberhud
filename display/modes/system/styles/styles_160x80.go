package styles

import (
	"github.com/databeast/cyberhud/display/style"
)

// System style declarations for 160x80 panels (landscape).
// Variant set: skeleton (adaptiveBuild) + default, compact, cores, top.

var MonoSlow160x80Style = def{
	name: "mono-slow-160x80",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.MonoSlow},
}

var MonoSlow160x80DefaultStyle = def{
	name: "mono-slow-160x80-default",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildDefaultStyle},
}

var MonoSlow160x80CompactStyle = def{
	name: "mono-slow-160x80-compact",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildCompactStyle},
}

var MonoSlow160x80CoresStyle = def{
	name: "mono-slow-160x80-cores",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildCoresStyle},
}

var MonoSlow160x80TopStyle = def{
	name: "mono-slow-160x80-top",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildTopStyle},
}

var GrayscaleSlow160x80Style = def{
	name: "grayscale-slow-160x80",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow160x80DefaultStyle = def{
	name: "grayscale-slow-160x80-default",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildDefaultStyle},
}

var GrayscaleSlow160x80CompactStyle = def{
	name: "grayscale-slow-160x80-compact",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildCompactStyle},
}

var GrayscaleSlow160x80CoresStyle = def{
	name: "grayscale-slow-160x80-cores",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildCoresStyle},
}

var GrayscaleSlow160x80TopStyle = def{
	name: "grayscale-slow-160x80-top",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildTopStyle},
}

var GrayscaleFast160x80Style = def{
	name: "grayscale-fast-160x80",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.GrayscaleFast},
}

var GrayscaleFast160x80DefaultStyle = def{
	name: "grayscale-fast-160x80-default",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildDefaultStyle},
}

var GrayscaleFast160x80CompactStyle = def{
	name: "grayscale-fast-160x80-compact",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildCompactStyle},
}

var GrayscaleFast160x80CoresStyle = def{
	name: "grayscale-fast-160x80-cores",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildCoresStyle},
}

var GrayscaleFast160x80TopStyle = def{
	name: "grayscale-fast-160x80-top",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildTopStyle},
}

var ColorSlow160x80Style = def{
	name: "color-slow-160x80",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.ColorSlow},
}

var ColorSlow160x80DefaultStyle = def{
	name: "color-slow-160x80-default",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildDefaultStyle},
}

var ColorSlow160x80CompactStyle = def{
	name: "color-slow-160x80-compact",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildCompactStyle},
}

var ColorSlow160x80CoresStyle = def{
	name: "color-slow-160x80-cores",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildCoresStyle},
}

var ColorSlow160x80TopStyle = def{
	name: "color-slow-160x80-top",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildTopStyle},
}

var ColorFast160x80Style = def{
	name: "color-160x80",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.ColorFast},
}

var ColorFast160x80DefaultStyle = def{
	name: "color-160x80-default",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.ColorFast},
	p:    Params{BuildFn: buildDefaultStyle},
}

var ColorFast160x80CompactStyle = def{
	name: "color-160x80-compact",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.ColorFast},
	p:    Params{BuildFn: buildCompactStyle},
}

var ColorFast160x80CoresStyle = def{
	name: "color-160x80-cores",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.ColorFast},
	p:    Params{BuildFn: buildCoresStyle},
}

var ColorFast160x80TopStyle = def{
	name: "color-160x80-top",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.ColorFast},
	p:    Params{BuildFn: buildTopStyle},
}

