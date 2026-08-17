package styles

import "github.com/databeast/cyberhud/display/style"

// GPIO control style declarations for 128x128 panels.

var MonoSlow128x128Style = def{
	name: "mono-slow-128x128",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 128, Capability: style.MonoSlow},
}

var MonoFast128x128Style = def{
	name: "mono-fast-128x128",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 128, Capability: style.MonoFast},
}

var GrayscaleSlow128x128Style = def{
	name: "grayscale-slow-128x128",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 128, Capability: style.GrayscaleSlow},
}

var GrayscaleFast128x128Style = def{
	name: "grayscale-fast-128x128",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 128, Capability: style.GrayscaleFast},
}

var ColorSlow128x128Style = def{
	name: "color-slow-128x128",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 128, Capability: style.ColorSlow},
}

var ColorFast128x128Style = def{
	name: "color-fast-128x128",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 128, Capability: style.ColorFast},
}
