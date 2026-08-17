package styles

import "github.com/databeast/cyberhud/display/style"

// GPIO control style declarations for 240x240 panels.

var MonoSlow240x240Style = def{
	name: "mono-slow-240x240",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 240, Capability: style.MonoSlow},
}

var MonoFast240x240Style = def{
	name: "mono-fast-240x240",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 240, Capability: style.MonoFast},
}

var GrayscaleSlow240x240Style = def{
	name: "grayscale-slow-240x240",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 240, Capability: style.GrayscaleSlow},
}

var GrayscaleFast240x240Style = def{
	name: "grayscale-fast-240x240",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 240, Capability: style.GrayscaleFast},
}

var ColorSlow240x240Style = def{
	name: "color-slow-240x240",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 240, Capability: style.ColorSlow},
}

var ColorFast240x240Style = def{
	name: "color-fast-240x240",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 240, Capability: style.ColorFast},
}
