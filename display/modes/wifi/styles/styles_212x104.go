package styles

import "github.com/databeast/cyberhud/display/style"

// WiFi style declarations for 212x104 panels.

var MonoSlow212x104Style = def{
	name: "mono-slow-212x104",
	reqs: style.SurfaceRequirements{MinWidth: 212, MinHeight: 104, Capability: style.MonoSlow},
}

var MonoFast212x104Style = def{
	name: "mono-fast-212x104",
	reqs: style.SurfaceRequirements{MinWidth: 212, MinHeight: 104, Capability: style.MonoFast},
	p:    Params{Fast: true},
}

var GrayscaleSlow212x104Style = def{
	name: "grayscale-slow-212x104",
	reqs: style.SurfaceRequirements{MinWidth: 212, MinHeight: 104, Capability: style.GrayscaleSlow},
}

var GrayscaleFast212x104Style = def{
	name: "grayscale-fast-212x104",
	reqs: style.SurfaceRequirements{MinWidth: 212, MinHeight: 104, Capability: style.GrayscaleFast},
	p:    Params{Fast: true},
}

var ColorSlow212x104Style = def{
	name: "color-slow-212x104",
	reqs: style.SurfaceRequirements{MinWidth: 212, MinHeight: 104, Capability: style.ColorSlow},
}

var ColorFast212x104Style = def{
	name: "color-fast-212x104",
	reqs: style.SurfaceRequirements{MinWidth: 212, MinHeight: 104, Capability: style.ColorFast},
	p:    Params{Fast: true, Color: true},
}
