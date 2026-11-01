package styles

import (
	"github.com/databeast/cyberhud/display/style"
)

// System style declarations for 176x264 panels (portrait).
// Variant set: skeleton (adaptiveBuild) + default, compact, cores, top.

var GrayscaleSlow176x264Style = def{
	name: "grayscale-slow-176x264",
	reqs: style.SurfaceRequirements{MinWidth: 176, MinHeight: 264, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow176x264DefaultStyle = def{
	name: "grayscale-slow-176x264-default",
	reqs: style.SurfaceRequirements{MinWidth: 176, MinHeight: 264, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildDefaultStyle},
}

var GrayscaleSlow176x264CompactStyle = def{
	name: "grayscale-slow-176x264-compact",
	reqs: style.SurfaceRequirements{MinWidth: 176, MinHeight: 264, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildCompactStyle},
}

var GrayscaleSlow176x264CoresStyle = def{
	name: "grayscale-slow-176x264-cores",
	reqs: style.SurfaceRequirements{MinWidth: 176, MinHeight: 264, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildCoresStyle},
}

var GrayscaleSlow176x264TopStyle = def{
	name: "grayscale-slow-176x264-top",
	reqs: style.SurfaceRequirements{MinWidth: 176, MinHeight: 264, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildTopStyle},
}

var ColorSlow176x264Style = def{
	name: "color-slow-176x264",
	reqs: style.SurfaceRequirements{MinWidth: 176, MinHeight: 264, Capability: style.ColorSlow},
}

var ColorSlow176x264DefaultStyle = def{
	name: "color-slow-176x264-default",
	reqs: style.SurfaceRequirements{MinWidth: 176, MinHeight: 264, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildDefaultStyle},
}

var ColorSlow176x264CompactStyle = def{
	name: "color-slow-176x264-compact",
	reqs: style.SurfaceRequirements{MinWidth: 176, MinHeight: 264, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildCompactStyle},
}

var ColorSlow176x264CoresStyle = def{
	name: "color-slow-176x264-cores",
	reqs: style.SurfaceRequirements{MinWidth: 176, MinHeight: 264, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildCoresStyle},
}

var ColorSlow176x264TopStyle = def{
	name: "color-slow-176x264-top",
	reqs: style.SurfaceRequirements{MinWidth: 176, MinHeight: 264, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildTopStyle},
}

