package styles

import (
	"github.com/databeast/cyberhud/display/style"
)

// System style declarations for 264x176 panels (landscape).
// Variant set: skeleton (adaptiveBuild) + default, compact, cores, top.

var GrayscaleSlow264x176Style = def{
	name: "grayscale-slow-264x176",
	reqs: style.SurfaceRequirements{MinWidth: 264, MinHeight: 176, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow264x176DefaultStyle = def{
	name: "grayscale-slow-264x176-default",
	reqs: style.SurfaceRequirements{MinWidth: 264, MinHeight: 176, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildDefaultStyle},
}

var GrayscaleSlow264x176CompactStyle = def{
	name: "grayscale-slow-264x176-compact",
	reqs: style.SurfaceRequirements{MinWidth: 264, MinHeight: 176, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildCompactStyle},
}

var GrayscaleSlow264x176CoresStyle = def{
	name: "grayscale-slow-264x176-cores",
	reqs: style.SurfaceRequirements{MinWidth: 264, MinHeight: 176, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildCoresStyle},
}

var GrayscaleSlow264x176TopStyle = def{
	name: "grayscale-slow-264x176-top",
	reqs: style.SurfaceRequirements{MinWidth: 264, MinHeight: 176, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildTopStyle},
}

var ColorSlow264x176Style = def{
	name: "color-slow-264x176",
	reqs: style.SurfaceRequirements{MinWidth: 264, MinHeight: 176, Capability: style.ColorSlow},
}

var ColorSlow264x176DefaultStyle = def{
	name: "color-slow-264x176-default",
	reqs: style.SurfaceRequirements{MinWidth: 264, MinHeight: 176, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildDefaultStyle},
}

var ColorSlow264x176CompactStyle = def{
	name: "color-slow-264x176-compact",
	reqs: style.SurfaceRequirements{MinWidth: 264, MinHeight: 176, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildCompactStyle},
}

var ColorSlow264x176CoresStyle = def{
	name: "color-slow-264x176-cores",
	reqs: style.SurfaceRequirements{MinWidth: 264, MinHeight: 176, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildCoresStyle},
}

var ColorSlow264x176TopStyle = def{
	name: "color-slow-264x176-top",
	reqs: style.SurfaceRequirements{MinWidth: 264, MinHeight: 176, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildTopStyle},
}

