package styles

import "github.com/databeast/cyberhud/display/style"

// WiFi style declarations for 800x600 panels.

var MonoSlow800x600Style = def{
	name: "mono-slow-800x600",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 600, Capability: style.MonoSlow},
}

var MonoFast800x600Style = def{
	name: "mono-fast-800x600",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 600, Capability: style.MonoFast},
	p:    Params{Fast: true},
}

var GrayscaleSlow800x600Style = def{
	name: "grayscale-slow-800x600",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 600, Capability: style.GrayscaleSlow},
}

var GrayscaleFast800x600Style = def{
	name: "grayscale-fast-800x600",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 600, Capability: style.GrayscaleFast},
	p:    Params{Fast: true},
}

var ColorSlow800x600Style = def{
	name: "color-slow-800x600",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 600, Capability: style.ColorSlow},
}

var ColorFast800x600Style = def{
	name: "color-fast-800x600",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 600, Capability: style.ColorFast},
	p:    Params{Fast: true, Color: true},
}
