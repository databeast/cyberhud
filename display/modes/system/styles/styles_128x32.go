package styles

import (
	"github.com/databeast/cyberhud/display/style"
)

// System style declarations for 128x32 panels (tiny).
// Variant set: skeleton (adaptiveBuild) + compact only.

var MonoSlow128x32Style = def{
	name: "mono-slow-128x32",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 32, Capability: style.MonoSlow},
}

var MonoSlow128x32CompactStyle = def{
	name: "mono-slow-128x32-compact",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 32, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildCompactStyle},
}

var MonoFast128x32Style = def{
	name: "mono-128x32",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 32, Capability: style.MonoFast},
}

var MonoFast128x32CompactStyle = def{
	name: "mono-128x32-compact",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 32, Capability: style.MonoFast},
	p:    Params{BuildFn: buildCompactStyle},
}

var GrayscaleSlow128x32Style = def{
	name: "grayscale-slow-128x32",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 32, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow128x32CompactStyle = def{
	name: "grayscale-slow-128x32-compact",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 32, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildCompactStyle},
}

var ColorSlow128x32Style = def{
	name: "color-slow-128x32",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 32, Capability: style.ColorSlow},
}

var ColorSlow128x32CompactStyle = def{
	name: "color-slow-128x32-compact",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 32, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildCompactStyle},
}

