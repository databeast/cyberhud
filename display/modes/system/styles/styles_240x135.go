package styles

import (
	"github.com/databeast/cyberhud/display/style"
)

// System style declarations for 240x135 panels (landscape).
// Variant set: skeleton (adaptiveBuild) + default, compact, cores, top.

var MonoSlow240x135Style = def{
	name: "mono-slow-240x135",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.MonoSlow},
}

var MonoSlow240x135DefaultStyle = def{
	name: "mono-slow-240x135-default",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildDefaultStyle},
}

var MonoSlow240x135CompactStyle = def{
	name: "mono-slow-240x135-compact",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildCompactStyle},
}

var MonoSlow240x135CoresStyle = def{
	name: "mono-slow-240x135-cores",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildCoresStyle},
}

var MonoSlow240x135TopStyle = def{
	name: "mono-slow-240x135-top",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildTopStyle},
}

var GrayscaleSlow240x135Style = def{
	name: "grayscale-slow-240x135",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow240x135DefaultStyle = def{
	name: "grayscale-slow-240x135-default",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildDefaultStyle},
}

var GrayscaleSlow240x135CompactStyle = def{
	name: "grayscale-slow-240x135-compact",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildCompactStyle},
}

var GrayscaleSlow240x135CoresStyle = def{
	name: "grayscale-slow-240x135-cores",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildCoresStyle},
}

var GrayscaleSlow240x135TopStyle = def{
	name: "grayscale-slow-240x135-top",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildTopStyle},
}

var GrayscaleFast240x135Style = def{
	name: "grayscale-fast-240x135",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.GrayscaleFast},
}

var GrayscaleFast240x135DefaultStyle = def{
	name: "grayscale-fast-240x135-default",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildDefaultStyle},
}

var GrayscaleFast240x135CompactStyle = def{
	name: "grayscale-fast-240x135-compact",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildCompactStyle},
}

var GrayscaleFast240x135CoresStyle = def{
	name: "grayscale-fast-240x135-cores",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildCoresStyle},
}

var GrayscaleFast240x135TopStyle = def{
	name: "grayscale-fast-240x135-top",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildTopStyle},
}

var ColorSlow240x135Style = def{
	name: "color-slow-240x135",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.ColorSlow},
}

var ColorSlow240x135DefaultStyle = def{
	name: "color-slow-240x135-default",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildDefaultStyle},
}

var ColorSlow240x135CompactStyle = def{
	name: "color-slow-240x135-compact",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildCompactStyle},
}

var ColorSlow240x135CoresStyle = def{
	name: "color-slow-240x135-cores",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildCoresStyle},
}

var ColorSlow240x135TopStyle = def{
	name: "color-slow-240x135-top",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildTopStyle},
}

var ColorFast240x135Style = def{
	name: "color-240x135",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.ColorFast},
}

var ColorFast240x135DefaultStyle = def{
	name: "color-240x135-default",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.ColorFast},
	p:    Params{BuildFn: buildDefaultStyle},
}

var ColorFast240x135CompactStyle = def{
	name: "color-240x135-compact",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.ColorFast},
	p:    Params{BuildFn: buildCompactStyle},
}

var ColorFast240x135CoresStyle = def{
	name: "color-240x135-cores",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.ColorFast},
	p:    Params{BuildFn: buildCoresStyle},
}

var ColorFast240x135TopStyle = def{
	name: "color-240x135-top",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.ColorFast},
	p:    Params{BuildFn: buildTopStyle},
}

