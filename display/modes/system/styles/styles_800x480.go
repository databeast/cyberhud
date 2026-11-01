package styles

import (
	"github.com/databeast/cyberhud/display/style"
)

// System style declarations for 800x480 panels (landscape).
// Variant set: skeleton (adaptiveBuild) + default, compact, cores, top.

var MonoSlow800x480Style = def{
	name: "mono-slow-800x480",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.MonoSlow},
}

var MonoSlow800x480DefaultStyle = def{
	name: "mono-slow-800x480-default",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildDefaultStyle},
}

var MonoSlow800x480CompactStyle = def{
	name: "mono-slow-800x480-compact",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildCompactStyle},
}

var MonoSlow800x480CoresStyle = def{
	name: "mono-slow-800x480-cores",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildCoresStyle},
}

var MonoSlow800x480TopStyle = def{
	name: "mono-slow-800x480-top",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildTopStyle},
}

var GrayscaleSlow800x480Style = def{
	name: "grayscale-slow-800x480",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow800x480DefaultStyle = def{
	name: "grayscale-slow-800x480-default",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildDefaultStyle},
}

var GrayscaleSlow800x480CompactStyle = def{
	name: "grayscale-slow-800x480-compact",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildCompactStyle},
}

var GrayscaleSlow800x480CoresStyle = def{
	name: "grayscale-slow-800x480-cores",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildCoresStyle},
}

var GrayscaleSlow800x480TopStyle = def{
	name: "grayscale-slow-800x480-top",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildTopStyle},
}

var GrayscaleFast800x480Style = def{
	name: "grayscale-fast-800x480",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.GrayscaleFast},
}

var GrayscaleFast800x480DefaultStyle = def{
	name: "grayscale-fast-800x480-default",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildDefaultStyle},
}

var GrayscaleFast800x480CompactStyle = def{
	name: "grayscale-fast-800x480-compact",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildCompactStyle},
}

var GrayscaleFast800x480CoresStyle = def{
	name: "grayscale-fast-800x480-cores",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildCoresStyle},
}

var GrayscaleFast800x480TopStyle = def{
	name: "grayscale-fast-800x480-top",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildTopStyle},
}

var ColorSlow800x480Style = def{
	name: "color-slow-800x480",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.ColorSlow},
}

var ColorSlow800x480DefaultStyle = def{
	name: "color-slow-800x480-default",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildDefaultStyle},
}

var ColorSlow800x480CompactStyle = def{
	name: "color-slow-800x480-compact",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildCompactStyle},
}

var ColorSlow800x480CoresStyle = def{
	name: "color-slow-800x480-cores",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildCoresStyle},
}

var ColorSlow800x480TopStyle = def{
	name: "color-slow-800x480-top",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildTopStyle},
}

var ColorFast800x480Style = def{
	name: "color-800x480",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.ColorFast},
}

var ColorFast800x480DefaultStyle = def{
	name: "color-800x480-default",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.ColorFast},
	p:    Params{BuildFn: buildDefaultStyle},
}

var ColorFast800x480CompactStyle = def{
	name: "color-800x480-compact",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.ColorFast},
	p:    Params{BuildFn: buildCompactStyle},
}

var ColorFast800x480CoresStyle = def{
	name: "color-800x480-cores",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.ColorFast},
	p:    Params{BuildFn: buildCoresStyle},
}

var ColorFast800x480TopStyle = def{
	name: "color-800x480-top",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.ColorFast},
	p:    Params{BuildFn: buildTopStyle},
}

