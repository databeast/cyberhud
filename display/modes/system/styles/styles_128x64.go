package styles

import (
	"github.com/databeast/cyberhud/display/style"
)

// System style declarations for 128x64 panels (landscape).
// Variant set: skeleton (adaptiveBuild) + default, compact, cores, top.

var MonoFast128x64Style = def{
	name: "mono-128x64",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 64, Capability: style.MonoFast},
}

var MonoFast128x64DefaultStyle = def{
	name: "mono-128x64-default",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 64, Capability: style.MonoFast},
	p:    Params{BuildFn: buildDefaultStyle},
}

var MonoFast128x64CompactStyle = def{
	name: "mono-128x64-compact",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 64, Capability: style.MonoFast},
	p:    Params{BuildFn: buildCompactStyle},
}

var MonoFast128x64CoresStyle = def{
	name: "mono-128x64-cores",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 64, Capability: style.MonoFast},
	p:    Params{BuildFn: buildCoresStyle},
}

var MonoFast128x64TopStyle = def{
	name: "mono-128x64-top",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 64, Capability: style.MonoFast},
	p:    Params{BuildFn: buildTopStyle},
}

var GrayscaleSlow128x64Style = def{
	name: "grayscale-slow-128x64",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 64, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow128x64DefaultStyle = def{
	name: "grayscale-slow-128x64-default",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 64, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildDefaultStyle},
}

var GrayscaleSlow128x64CompactStyle = def{
	name: "grayscale-slow-128x64-compact",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 64, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildCompactStyle},
}

var GrayscaleSlow128x64CoresStyle = def{
	name: "grayscale-slow-128x64-cores",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 64, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildCoresStyle},
}

var GrayscaleSlow128x64TopStyle = def{
	name: "grayscale-slow-128x64-top",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 64, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildTopStyle},
}

var ColorSlow128x64Style = def{
	name: "color-slow-128x64",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 64, Capability: style.ColorSlow},
}

var ColorSlow128x64DefaultStyle = def{
	name: "color-slow-128x64-default",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 64, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildDefaultStyle},
}

var ColorSlow128x64CompactStyle = def{
	name: "color-slow-128x64-compact",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 64, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildCompactStyle},
}

var ColorSlow128x64CoresStyle = def{
	name: "color-slow-128x64-cores",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 64, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildCoresStyle},
}

var ColorSlow128x64TopStyle = def{
	name: "color-slow-128x64-top",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 64, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildTopStyle},
}

