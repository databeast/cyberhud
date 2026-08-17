package styles

import "github.com/databeast/cyberhud/display/style"

// WiFi style declarations for 160x80 panels.

var MonoSlow160x80Style = def{
	name: "mono-slow-160x80",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.MonoSlow},
}

var MonoFast160x80Style = def{
	name: "mono-fast-160x80",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.MonoFast},
	p:    Params{Fast: true},
}

var GrayscaleSlow160x80Style = def{
	name: "grayscale-slow-160x80",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.GrayscaleSlow},
}

var GrayscaleFast160x80Style = def{
	name: "grayscale-fast-160x80",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.GrayscaleFast},
	p:    Params{Fast: true},
}

var ColorSlow160x80Style = def{
	name: "color-slow-160x80",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.ColorSlow},
}

var ColorFast160x80Style = def{
	name: "color-fast-160x80",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.ColorFast},
	p:    Params{Fast: true, Color: true},
}

// Portrait variant: 80x160

var MonoSlow80x160Style = def{
	name: "mono-slow-80x160",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.MonoSlow},
}

var MonoFast80x160Style = def{
	name: "mono-fast-80x160",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.MonoFast},
	p:    Params{Fast: true},
}

var GrayscaleSlow80x160Style = def{
	name: "grayscale-slow-80x160",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.GrayscaleSlow},
}

var GrayscaleFast80x160Style = def{
	name: "grayscale-fast-80x160",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.GrayscaleFast},
	p:    Params{Fast: true},
}

var ColorSlow80x160Style = def{
	name: "color-slow-80x160",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.ColorSlow},
}

var ColorFast80x160Style = def{
	name: "color-fast-80x160",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.ColorFast},
	p:    Params{Fast: true, Color: true},
}
