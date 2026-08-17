package styles

import "github.com/databeast/cyberhud/display/style"

// WiFi style declarations for 200x200 panels.

var MonoSlow200x200Style = def{
	name: "mono-slow-200x200",
	reqs: style.SurfaceRequirements{MinWidth: 200, MinHeight: 200, Capability: style.MonoSlow},
}

var MonoFast200x200Style = def{
	name: "mono-fast-200x200",
	reqs: style.SurfaceRequirements{MinWidth: 200, MinHeight: 200, Capability: style.MonoFast},
	p:    Params{Fast: true},
}

var GrayscaleSlow200x200Style = def{
	name: "grayscale-slow-200x200",
	reqs: style.SurfaceRequirements{MinWidth: 200, MinHeight: 200, Capability: style.GrayscaleSlow},
}

var GrayscaleFast200x200Style = def{
	name: "grayscale-fast-200x200",
	reqs: style.SurfaceRequirements{MinWidth: 200, MinHeight: 200, Capability: style.GrayscaleFast},
	p:    Params{Fast: true},
}

var ColorSlow200x200Style = def{
	name: "color-slow-200x200",
	reqs: style.SurfaceRequirements{MinWidth: 200, MinHeight: 200, Capability: style.ColorSlow},
}

var ColorFast200x200Style = def{
	name: "color-fast-200x200",
	reqs: style.SurfaceRequirements{MinWidth: 200, MinHeight: 200, Capability: style.ColorFast},
	p:    Params{Fast: true, Color: true},
}
