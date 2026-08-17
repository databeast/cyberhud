package styles

import "github.com/databeast/cyberhud/display/style"

// WiFi style declarations for 400x300 panels.

var MonoSlow400x300Style = def{
	name: "mono-slow-400x300",
	reqs: style.SurfaceRequirements{MinWidth: 400, MinHeight: 300, Capability: style.MonoSlow},
}

var MonoFast400x300Style = def{
	name: "mono-fast-400x300",
	reqs: style.SurfaceRequirements{MinWidth: 400, MinHeight: 300, Capability: style.MonoFast},
	p:    Params{Fast: true},
}

var GrayscaleSlow400x300Style = def{
	name: "grayscale-slow-400x300",
	reqs: style.SurfaceRequirements{MinWidth: 400, MinHeight: 300, Capability: style.GrayscaleSlow},
}

var GrayscaleFast400x300Style = def{
	name: "grayscale-fast-400x300",
	reqs: style.SurfaceRequirements{MinWidth: 400, MinHeight: 300, Capability: style.GrayscaleFast},
	p:    Params{Fast: true},
}

var ColorSlow400x300Style = def{
	name: "color-slow-400x300",
	reqs: style.SurfaceRequirements{MinWidth: 400, MinHeight: 300, Capability: style.ColorSlow},
}

var ColorFast400x300Style = def{
	name: "color-fast-400x300",
	reqs: style.SurfaceRequirements{MinWidth: 400, MinHeight: 300, Capability: style.ColorFast},
	p:    Params{Fast: true, Color: true},
}
