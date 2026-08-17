package styles

import "github.com/databeast/cyberhud/display/style"

// WiFi style declarations for 264x176 panels.

var MonoSlow264x176Style = def{
	name: "mono-slow-264x176",
	reqs: style.SurfaceRequirements{MinWidth: 264, MinHeight: 176, Capability: style.MonoSlow},
}

var MonoFast264x176Style = def{
	name: "mono-fast-264x176",
	reqs: style.SurfaceRequirements{MinWidth: 264, MinHeight: 176, Capability: style.MonoFast},
	p:    Params{Fast: true},
}

var GrayscaleSlow264x176Style = def{
	name: "grayscale-slow-264x176",
	reqs: style.SurfaceRequirements{MinWidth: 264, MinHeight: 176, Capability: style.GrayscaleSlow},
}

var GrayscaleFast264x176Style = def{
	name: "grayscale-fast-264x176",
	reqs: style.SurfaceRequirements{MinWidth: 264, MinHeight: 176, Capability: style.GrayscaleFast},
	p:    Params{Fast: true},
}

var ColorSlow264x176Style = def{
	name: "color-slow-264x176",
	reqs: style.SurfaceRequirements{MinWidth: 264, MinHeight: 176, Capability: style.ColorSlow},
}

var ColorFast264x176Style = def{
	name: "color-fast-264x176",
	reqs: style.SurfaceRequirements{MinWidth: 264, MinHeight: 176, Capability: style.ColorFast},
	p:    Params{Fast: true, Color: true},
}
