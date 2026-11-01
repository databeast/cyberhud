package styles

import (
	"github.com/databeast/cyberhud/display/style"
)

// System style declarations for 122x250 panels (portrait).
// Variant set: skeleton (adaptiveBuild) + default, compact, cores, top.

var GrayscaleSlow122x250Style = def{
	name: "grayscale-slow-122x250",
	reqs: style.SurfaceRequirements{MinWidth: 122, MinHeight: 250, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow122x250DefaultStyle = def{
	name: "grayscale-slow-122x250-default",
	reqs: style.SurfaceRequirements{MinWidth: 122, MinHeight: 250, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildDefaultStyle},
}

var GrayscaleSlow122x250CompactStyle = def{
	name: "grayscale-slow-122x250-compact",
	reqs: style.SurfaceRequirements{MinWidth: 122, MinHeight: 250, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildCompactStyle},
}

var GrayscaleSlow122x250CoresStyle = def{
	name: "grayscale-slow-122x250-cores",
	reqs: style.SurfaceRequirements{MinWidth: 122, MinHeight: 250, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildCoresStyle},
}

var GrayscaleSlow122x250TopStyle = def{
	name: "grayscale-slow-122x250-top",
	reqs: style.SurfaceRequirements{MinWidth: 122, MinHeight: 250, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildTopStyle},
}

var ColorSlow122x250Style = def{
	name: "color-slow-122x250",
	reqs: style.SurfaceRequirements{MinWidth: 122, MinHeight: 250, Capability: style.ColorSlow},
}

var ColorSlow122x250DefaultStyle = def{
	name: "color-slow-122x250-default",
	reqs: style.SurfaceRequirements{MinWidth: 122, MinHeight: 250, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildDefaultStyle},
}

var ColorSlow122x250CompactStyle = def{
	name: "color-slow-122x250-compact",
	reqs: style.SurfaceRequirements{MinWidth: 122, MinHeight: 250, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildCompactStyle},
}

var ColorSlow122x250CoresStyle = def{
	name: "color-slow-122x250-cores",
	reqs: style.SurfaceRequirements{MinWidth: 122, MinHeight: 250, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildCoresStyle},
}

var ColorSlow122x250TopStyle = def{
	name: "color-slow-122x250-top",
	reqs: style.SurfaceRequirements{MinWidth: 122, MinHeight: 250, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildTopStyle},
}

