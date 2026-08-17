package styles

import "github.com/databeast/cyberhud/display/style"

// WiFi style declarations for 800x480 panels.

var MonoSlow800x480Style = def{
	name: "mono-slow-800x480",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.MonoSlow},
}

var MonoFast800x480Style = def{
	name: "mono-fast-800x480",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.MonoFast},
	p:    Params{Fast: true},
}

var GrayscaleSlow800x480Style = def{
	name: "grayscale-slow-800x480",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.GrayscaleSlow},
}

var GrayscaleFast800x480Style = def{
	name: "grayscale-fast-800x480",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.GrayscaleFast},
	p:    Params{Fast: true},
}

var ColorSlow800x480Style = def{
	name: "color-slow-800x480",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.ColorSlow},
}

var ColorFast800x480Style = def{
	name: "color-fast-800x480",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.ColorFast},
	p:    Params{Fast: true, Color: true},
}

// Portrait variant: 480x800

var MonoSlow480x800Style = def{
	name: "mono-slow-480x800",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.MonoSlow},
}

var MonoFast480x800Style = def{
	name: "mono-fast-480x800",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.MonoFast},
	p:    Params{Fast: true},
}

var GrayscaleSlow480x800Style = def{
	name: "grayscale-slow-480x800",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.GrayscaleSlow},
}

var GrayscaleFast480x800Style = def{
	name: "grayscale-fast-480x800",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.GrayscaleFast},
	p:    Params{Fast: true},
}

var ColorSlow480x800Style = def{
	name: "color-slow-480x800",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.ColorSlow},
}

var ColorFast480x800Style = def{
	name: "color-fast-480x800",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.ColorFast},
	p:    Params{Fast: true, Color: true},
}
