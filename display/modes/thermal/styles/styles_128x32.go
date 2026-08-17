package styles

import (
	"github.com/databeast/cyberhud/display/style"
)

// Thermal style declarations for 128x32 panels (landscape narrow strip).
// Height < 64 → minimal variant only.

var MonoSlow128x32Style = def{
	name: "mono-slow-128x32",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 32, Capability: style.MonoSlow},
}

var MonoSlow128x32MinimalStyle = def{
	name: "mono-slow-128x32-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 32, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildMinimalStyle},
}

var MonoFast128x32Style = def{
	name: "mono-128x32",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 32, Capability: style.MonoFast},
}

var MonoFast128x32MinimalStyle = def{
	name: "mono-128x32-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 32, Capability: style.MonoFast},
	p:    Params{BuildFn: buildMinimalStyle},
}

var GrayscaleSlow128x32Style = def{
	name: "grayscale-slow-128x32",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 32, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow128x32MinimalStyle = def{
	name: "grayscale-slow-128x32-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 32, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildMinimalStyle},
}

var ColorSlow128x32Style = def{
	name: "color-slow-128x32",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 32, Capability: style.ColorSlow},
}

var ColorSlow128x32MinimalStyle = def{
	name: "color-slow-128x32-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 32, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildMinimalStyle},
}
