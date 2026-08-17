package styles

import "github.com/databeast/cyberhud/display/style"

// GPIO control style declarations for 128x32 panels.

var MonoSlow128x32Style = def{
	name: "mono-slow-128x32",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 32, Capability: style.MonoSlow},
}

var MonoFast128x32Style = def{
	name: "mono-fast-128x32",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 32, Capability: style.MonoFast},
}

var GrayscaleSlow128x32Style = def{
	name: "grayscale-slow-128x32",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 32, Capability: style.GrayscaleSlow},
}

var GrayscaleFast128x32Style = def{
	name: "grayscale-fast-128x32",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 32, Capability: style.GrayscaleFast},
}

var ColorSlow128x32Style = def{
	name: "color-slow-128x32",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 32, Capability: style.ColorSlow},
}

var ColorFast128x32Style = def{
	name: "color-fast-128x32",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 32, Capability: style.ColorFast},
}
