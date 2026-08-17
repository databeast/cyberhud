package styles

import "github.com/databeast/cyberhud/display/style"

// GPIO style declarations for this panel group.
//
// These are hand-tweakable declarations over the shared layouts in core.go.

var MonoFast800x480Style = def{
	name: "mono-fast-800x480",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.MonoFast},
}

var MonoFast480x800Style = def{
	name: "mono-fast-480x800",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.MonoFast},
}

var GrayscaleSlow480x800Style = def{
	name: "grayscale-slow-480x800",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow800x480Style = def{
	name: "grayscale-slow-800x480",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.GrayscaleSlow},
}

var GrayscaleFast800x480Style = def{
	name: "grayscale-fast-800x480",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.GrayscaleFast},
	p:    Params{Layout: layoutGrayscaleFast},
}

var GrayscaleFast480x800Style = def{
	name: "grayscale-fast-480x800",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.GrayscaleFast},
	p:    Params{Layout: layoutGrayscaleFast},
}

var EinkStyle800x480 = def{
	name: "eink-800x480",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.MonoSlow},
	p:    Params{BuildFn: einkFramed800Build},
}

var EinkStyle480x800 = def{
	name: "eink-480x800",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.MonoSlow},
	p:    Params{Layout: layoutEink},
}

var ColorSlow480x800Style = def{
	name: "color-slow-480x800",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.ColorSlow},
}

var ColorSlow800x480Style = def{
	name: "color-slow-800x480",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.ColorSlow},
}

var ColorStyle800x480 = def{
	name: "color-800x480",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.ColorFast},
	p:    Params{Layout: layoutColorFastRows},
}

var ColorStyle480x800 = def{
	name: "color-480x800",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.ColorFast},
	p:    Params{Layout: layoutColorFastRows},
}
