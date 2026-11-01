package styles

import (
	"github.com/databeast/cyberhud/display/style"
)

// System style declarations for 104x212 panels (portrait).
// Variant set: skeleton (adaptiveBuild) + default, compact, cores, top.

var GrayscaleSlow104x212Style = def{
	name: "grayscale-slow-104x212",
	reqs: style.SurfaceRequirements{MinWidth: 104, MinHeight: 212, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow104x212DefaultStyle = def{
	name: "grayscale-slow-104x212-default",
	reqs: style.SurfaceRequirements{MinWidth: 104, MinHeight: 212, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildDefaultStyle},
}

var GrayscaleSlow104x212CompactStyle = def{
	name: "grayscale-slow-104x212-compact",
	reqs: style.SurfaceRequirements{MinWidth: 104, MinHeight: 212, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildCompactStyle},
}

var GrayscaleSlow104x212CoresStyle = def{
	name: "grayscale-slow-104x212-cores",
	reqs: style.SurfaceRequirements{MinWidth: 104, MinHeight: 212, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildCoresStyle},
}

var GrayscaleSlow104x212TopStyle = def{
	name: "grayscale-slow-104x212-top",
	reqs: style.SurfaceRequirements{MinWidth: 104, MinHeight: 212, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildTopStyle},
}

var ColorSlow104x212Style = def{
	name: "color-slow-104x212",
	reqs: style.SurfaceRequirements{MinWidth: 104, MinHeight: 212, Capability: style.ColorSlow},
}

var ColorSlow104x212DefaultStyle = def{
	name: "color-slow-104x212-default",
	reqs: style.SurfaceRequirements{MinWidth: 104, MinHeight: 212, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildDefaultStyle},
}

var ColorSlow104x212CompactStyle = def{
	name: "color-slow-104x212-compact",
	reqs: style.SurfaceRequirements{MinWidth: 104, MinHeight: 212, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildCompactStyle},
}

var ColorSlow104x212CoresStyle = def{
	name: "color-slow-104x212-cores",
	reqs: style.SurfaceRequirements{MinWidth: 104, MinHeight: 212, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildCoresStyle},
}

var ColorSlow104x212TopStyle = def{
	name: "color-slow-104x212-top",
	reqs: style.SurfaceRequirements{MinWidth: 104, MinHeight: 212, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildTopStyle},
}

