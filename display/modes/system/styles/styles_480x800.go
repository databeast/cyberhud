package styles

import (
	"github.com/databeast/cyberhud/display/style"
)

// System style declarations for 480x800 panels (portrait).
// Variant set: skeleton (adaptiveBuild) + default, compact, cores, top.

var GrayscaleSlow480x800Style = def{
	name: "grayscale-slow-480x800",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow480x800DefaultStyle = def{
	name: "grayscale-slow-480x800-default",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildDefaultStyle},
}

var GrayscaleSlow480x800CompactStyle = def{
	name: "grayscale-slow-480x800-compact",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildCompactStyle},
}

var GrayscaleSlow480x800CoresStyle = def{
	name: "grayscale-slow-480x800-cores",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildCoresStyle},
}

var GrayscaleSlow480x800TopStyle = def{
	name: "grayscale-slow-480x800-top",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildTopStyle},
}

var GrayscaleFast480x800Style = def{
	name: "grayscale-fast-480x800",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.GrayscaleFast},
}

var GrayscaleFast480x800DefaultStyle = def{
	name: "grayscale-fast-480x800-default",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildDefaultStyle},
}

var GrayscaleFast480x800CompactStyle = def{
	name: "grayscale-fast-480x800-compact",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildCompactStyle},
}

var GrayscaleFast480x800CoresStyle = def{
	name: "grayscale-fast-480x800-cores",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildCoresStyle},
}

var GrayscaleFast480x800TopStyle = def{
	name: "grayscale-fast-480x800-top",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildTopStyle},
}

var ColorSlow480x800Style = def{
	name: "color-slow-480x800",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.ColorSlow},
}

var ColorSlow480x800DefaultStyle = def{
	name: "color-slow-480x800-default",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildDefaultStyle},
}

var ColorSlow480x800CompactStyle = def{
	name: "color-slow-480x800-compact",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildCompactStyle},
}

var ColorSlow480x800CoresStyle = def{
	name: "color-slow-480x800-cores",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildCoresStyle},
}

var ColorSlow480x800TopStyle = def{
	name: "color-slow-480x800-top",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildTopStyle},
}

var ColorFast480x800Style = def{
	name: "color-480x800",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.ColorFast},
}

var ColorFast480x800DefaultStyle = def{
	name: "color-480x800-default",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.ColorFast},
	p:    Params{BuildFn: buildDefaultStyle},
}

var ColorFast480x800CompactStyle = def{
	name: "color-480x800-compact",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.ColorFast},
	p:    Params{BuildFn: buildCompactStyle},
}

var ColorFast480x800CoresStyle = def{
	name: "color-480x800-cores",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.ColorFast},
	p:    Params{BuildFn: buildCoresStyle},
}

var ColorFast480x800TopStyle = def{
	name: "color-480x800-top",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.ColorFast},
	p:    Params{BuildFn: buildTopStyle},
}

