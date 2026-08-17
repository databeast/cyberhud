package styles

import "github.com/databeast/cyberhud/display/style"

// WiFi style declarations for 160x128 panels.

var MonoSlow160x128Style = def{
	name: "mono-slow-160x128",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 128, Capability: style.MonoSlow},
}

var MonoFast160x128Style = def{
	name: "mono-fast-160x128",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 128, Capability: style.MonoFast},
	p:    Params{Fast: true},
}

var GrayscaleSlow160x128Style = def{
	name: "grayscale-slow-160x128",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 128, Capability: style.GrayscaleSlow},
}

var GrayscaleFast160x128Style = def{
	name: "grayscale-fast-160x128",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 128, Capability: style.GrayscaleFast},
	p:    Params{Fast: true},
}

var ColorSlow160x128Style = def{
	name: "color-slow-160x128",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 128, Capability: style.ColorSlow},
}

var ColorFast160x128Style = def{
	name: "color-fast-160x128",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 128, Capability: style.ColorFast},
	p:    Params{Fast: true, Color: true},
}

// Portrait variant (128×160)

var MonoSlow128x160Style = def{
	name: "mono-slow-128x160",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 160, Capability: style.MonoSlow},
}

var MonoFast128x160Style = def{
	name: "mono-fast-128x160",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 160, Capability: style.MonoFast},
	p:    Params{Fast: true},
}

var GrayscaleSlow128x160Style = def{
	name: "grayscale-slow-128x160",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 160, Capability: style.GrayscaleSlow},
}

var GrayscaleFast128x160Style = def{
	name: "grayscale-fast-128x160",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 160, Capability: style.GrayscaleFast},
	p:    Params{Fast: true},
}

var ColorSlow128x160Style = def{
	name: "color-slow-128x160",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 160, Capability: style.ColorSlow},
}

var ColorFast128x160Style = def{
	name: "color-fast-128x160",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 160, Capability: style.ColorFast},
	p:    Params{Fast: true, Color: true},
}
