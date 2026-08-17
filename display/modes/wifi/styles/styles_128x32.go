package styles

import "github.com/databeast/cyberhud/display/style"

// WiFi style declarations for 128x32 panels.

var MonoSlow128x32Style = def{
	name: "mono-slow-128x32",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 32, Capability: style.MonoSlow},
}

var MonoFast128x32Style = def{
	name: "mono-fast-128x32",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 32, Capability: style.MonoFast},
	p:    Params{Fast: true},
}

var GrayscaleSlow128x32Style = def{
	name: "grayscale-slow-128x32",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 32, Capability: style.GrayscaleSlow},
}

var GrayscaleFast128x32Style = def{
	name: "grayscale-fast-128x32",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 32, Capability: style.GrayscaleFast},
	p:    Params{Fast: true},
}

var ColorSlow128x32Style = def{
	name: "color-slow-128x32",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 32, Capability: style.ColorSlow},
}

var ColorFast128x32Style = def{
	name: "color-fast-128x32",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 32, Capability: style.ColorFast},
	p:    Params{Fast: true, Color: true},
}

// Portrait variant: 32x128

var MonoSlow32x128Style = def{
	name: "mono-slow-32x128",
	reqs: style.SurfaceRequirements{MinWidth: 32, MinHeight: 128, Capability: style.MonoSlow},
}

var MonoFast32x128Style = def{
	name: "mono-fast-32x128",
	reqs: style.SurfaceRequirements{MinWidth: 32, MinHeight: 128, Capability: style.MonoFast},
	p:    Params{Fast: true},
}

var GrayscaleSlow32x128Style = def{
	name: "grayscale-slow-32x128",
	reqs: style.SurfaceRequirements{MinWidth: 32, MinHeight: 128, Capability: style.GrayscaleSlow},
}

var GrayscaleFast32x128Style = def{
	name: "grayscale-fast-32x128",
	reqs: style.SurfaceRequirements{MinWidth: 32, MinHeight: 128, Capability: style.GrayscaleFast},
	p:    Params{Fast: true},
}

var ColorSlow32x128Style = def{
	name: "color-slow-32x128",
	reqs: style.SurfaceRequirements{MinWidth: 32, MinHeight: 128, Capability: style.ColorSlow},
}

var ColorFast32x128Style = def{
	name: "color-fast-32x128",
	reqs: style.SurfaceRequirements{MinWidth: 32, MinHeight: 128, Capability: style.ColorFast},
	p:    Params{Fast: true, Color: true},
}
