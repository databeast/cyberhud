package styles

import "github.com/databeast/cyberhud/display/style"

// WiFi style declarations for 250x122 panels.

var MonoSlow250x122Style = def{
	name: "mono-slow-250x122",
	reqs: style.SurfaceRequirements{MinWidth: 250, MinHeight: 122, Capability: style.MonoSlow},
}

var MonoFast250x122Style = def{
	name: "mono-fast-250x122",
	reqs: style.SurfaceRequirements{MinWidth: 250, MinHeight: 122, Capability: style.MonoFast},
	p:    Params{Fast: true},
}

var GrayscaleSlow250x122Style = def{
	name: "grayscale-slow-250x122",
	reqs: style.SurfaceRequirements{MinWidth: 250, MinHeight: 122, Capability: style.GrayscaleSlow},
}

var GrayscaleFast250x122Style = def{
	name: "grayscale-fast-250x122",
	reqs: style.SurfaceRequirements{MinWidth: 250, MinHeight: 122, Capability: style.GrayscaleFast},
	p:    Params{Fast: true},
}

var ColorSlow250x122Style = def{
	name: "color-slow-250x122",
	reqs: style.SurfaceRequirements{MinWidth: 250, MinHeight: 122, Capability: style.ColorSlow},
}

var ColorFast250x122Style = def{
	name: "color-fast-250x122",
	reqs: style.SurfaceRequirements{MinWidth: 250, MinHeight: 122, Capability: style.ColorFast},
	p:    Params{Fast: true, Color: true},
}
