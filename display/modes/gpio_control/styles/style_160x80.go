package styles

import "github.com/databeast/cyberhud/display/style"

// GPIO control style declarations for 160x80 panels.

var MonoSlow160x80Style = def{
	name: "mono-slow-160x80",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.MonoSlow},
}

var MonoFast160x80Style = def{
	name: "mono-fast-160x80",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.MonoFast},
}

var GrayscaleSlow160x80Style = def{
	name: "grayscale-slow-160x80",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.GrayscaleSlow},
}

var GrayscaleFast160x80Style = def{
	name: "grayscale-fast-160x80",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.GrayscaleFast},
}

var ColorSlow160x80Style = def{
	name: "color-slow-160x80",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.ColorSlow},
}

var ColorFast160x80Style = def{
	name: "color-fast-160x80",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.ColorFast},
}
