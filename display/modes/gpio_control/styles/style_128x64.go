package styles

import "github.com/databeast/cyberhud/display/style"

// GPIO control style declarations for 128x64 panels.

var MonoSlow128x64Style = def{
	name: "mono-slow-128x64",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 64, Capability: style.MonoSlow},
}

var MonoFast128x64Style = def{
	name: "mono-fast-128x64",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 64, Capability: style.MonoFast},
}

var GrayscaleSlow128x64Style = def{
	name: "grayscale-slow-128x64",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 64, Capability: style.GrayscaleSlow},
}

var GrayscaleFast128x64Style = def{
	name: "grayscale-fast-128x64",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 64, Capability: style.GrayscaleFast},
}

var ColorSlow128x64Style = def{
	name: "color-slow-128x64",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 64, Capability: style.ColorSlow},
}

var ColorFast128x64Style = def{
	name: "color-fast-128x64",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 64, Capability: style.ColorFast},
}
