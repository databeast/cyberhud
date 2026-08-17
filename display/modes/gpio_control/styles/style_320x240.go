package styles

import "github.com/databeast/cyberhud/display/style"

// GPIO control style declarations for 320x240 panels.

var MonoSlow320x240Style = def{
	name: "mono-slow-320x240",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 240, Capability: style.MonoSlow},
}

var MonoFast320x240Style = def{
	name: "mono-fast-320x240",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 240, Capability: style.MonoFast},
}

var GrayscaleSlow320x240Style = def{
	name: "grayscale-slow-320x240",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 240, Capability: style.GrayscaleSlow},
}

var GrayscaleFast320x240Style = def{
	name: "grayscale-fast-320x240",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 240, Capability: style.GrayscaleFast},
}

var ColorSlow320x240Style = def{
	name: "color-slow-320x240",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 240, Capability: style.ColorSlow},
}

var ColorFast320x240Style = def{
	name: "color-fast-320x240",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 240, Capability: style.ColorFast},
}
