package styles

import (
	"github.com/databeast/cyberhud/display/style"
)

// Thermal style declarations for 32x128 panels (portrait narrow strip).
// Width < 64 → minimal variant only.

var MonoSlow32x128Style = def{
	name: "mono-slow-32x128",
	reqs: style.SurfaceRequirements{MinWidth: 32, MinHeight: 128, Capability: style.MonoSlow},
}

var MonoSlow32x128MinimalStyle = def{
	name: "mono-slow-32x128-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 32, MinHeight: 128, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildMinimalStyle},
}

var MonoFast32x128Style = def{
	name: "mono-32x128",
	reqs: style.SurfaceRequirements{MinWidth: 32, MinHeight: 128, Capability: style.MonoFast},
}

var MonoFast32x128MinimalStyle = def{
	name: "mono-32x128-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 32, MinHeight: 128, Capability: style.MonoFast},
	p:    Params{BuildFn: buildMinimalStyle},
}

var GrayscaleSlow32x128Style = def{
	name: "grayscale-slow-32x128",
	reqs: style.SurfaceRequirements{MinWidth: 32, MinHeight: 128, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow32x128MinimalStyle = def{
	name: "grayscale-slow-32x128-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 32, MinHeight: 128, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildMinimalStyle},
}

var ColorSlow32x128Style = def{
	name: "color-slow-32x128",
	reqs: style.SurfaceRequirements{MinWidth: 32, MinHeight: 128, Capability: style.ColorSlow},
}

var ColorSlow32x128MinimalStyle = def{
	name: "color-slow-32x128-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 32, MinHeight: 128, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildMinimalStyle},
}
