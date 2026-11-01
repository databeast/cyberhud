package styles

import (
	"github.com/databeast/cyberhud/display/style"
)

// System style declarations for 32x128 panels (tiny).
// Variant set: skeleton (adaptiveBuild) + compact only.

var MonoSlow32x128Style = def{
	name: "mono-slow-32x128",
	reqs: style.SurfaceRequirements{MinWidth: 32, MinHeight: 128, Capability: style.MonoSlow},
}

var MonoSlow32x128CompactStyle = def{
	name: "mono-slow-32x128-compact",
	reqs: style.SurfaceRequirements{MinWidth: 32, MinHeight: 128, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildCompactStyle},
}

var MonoFast32x128Style = def{
	name: "mono-32x128",
	reqs: style.SurfaceRequirements{MinWidth: 32, MinHeight: 128, Capability: style.MonoFast},
}

var MonoFast32x128CompactStyle = def{
	name: "mono-32x128-compact",
	reqs: style.SurfaceRequirements{MinWidth: 32, MinHeight: 128, Capability: style.MonoFast},
	p:    Params{BuildFn: buildCompactStyle},
}

var GrayscaleSlow32x128Style = def{
	name: "grayscale-slow-32x128",
	reqs: style.SurfaceRequirements{MinWidth: 32, MinHeight: 128, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow32x128CompactStyle = def{
	name: "grayscale-slow-32x128-compact",
	reqs: style.SurfaceRequirements{MinWidth: 32, MinHeight: 128, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildCompactStyle},
}

var ColorSlow32x128Style = def{
	name: "color-slow-32x128",
	reqs: style.SurfaceRequirements{MinWidth: 32, MinHeight: 128, Capability: style.ColorSlow},
}

var ColorSlow32x128CompactStyle = def{
	name: "color-slow-32x128-compact",
	reqs: style.SurfaceRequirements{MinWidth: 32, MinHeight: 128, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildCompactStyle},
}

