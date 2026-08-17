package styles

import "github.com/databeast/cyberhud/display/style"

// GPIO style declarations for this panel group.
//
// These are hand-tweakable declarations over the shared layouts in core.go.

var MonoFast104x212Style = def{
	name: "mono-fast-104x212",
	reqs: style.SurfaceRequirements{MinWidth: 104, MinHeight: 212, Capability: style.MonoFast},
}

var MonoFast212x104Style = def{
	name: "mono-fast-212x104",
	reqs: style.SurfaceRequirements{MinWidth: 212, MinHeight: 104, Capability: style.MonoFast},
}

var GrayscaleSlow104x212Style = def{
	name: "grayscale-slow-104x212",
	reqs: style.SurfaceRequirements{MinWidth: 104, MinHeight: 212, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow212x104Style = def{
	name: "grayscale-slow-212x104",
	reqs: style.SurfaceRequirements{MinWidth: 212, MinHeight: 104, Capability: style.GrayscaleSlow},
}

var GrayscaleFast104x212Style = def{
	name: "grayscale-fast-104x212",
	reqs: style.SurfaceRequirements{MinWidth: 104, MinHeight: 212, Capability: style.GrayscaleFast},
}

var GrayscaleFast212x104Style = def{
	name: "grayscale-fast-212x104",
	reqs: style.SurfaceRequirements{MinWidth: 212, MinHeight: 104, Capability: style.GrayscaleFast},
}

var EinkStyle212x104 = def{
	name: "eink-212x104",
	reqs: style.SurfaceRequirements{MinWidth: 212, MinHeight: 104, Capability: style.MonoSlow},
	p:    Params{Layout: layoutEink},
}

var EinkStyle104x212 = def{
	name: "eink-104x212",
	reqs: style.SurfaceRequirements{MinWidth: 104, MinHeight: 212, Capability: style.MonoSlow},
	p:    Params{Layout: layoutEink},
}

var ColorSlow104x212Style = def{
	name: "color-slow-104x212",
	reqs: style.SurfaceRequirements{MinWidth: 104, MinHeight: 212, Capability: style.ColorSlow},
}

var ColorSlow212x104Style = def{
	name: "color-slow-212x104",
	reqs: style.SurfaceRequirements{MinWidth: 212, MinHeight: 104, Capability: style.ColorSlow},
}

var ColorFast104x212Style = def{
	name: "color-fast-104x212",
	reqs: style.SurfaceRequirements{MinWidth: 104, MinHeight: 212, Capability: style.ColorFast},
}

var ColorFast212x104Style = def{
	name: "color-fast-212x104",
	reqs: style.SurfaceRequirements{MinWidth: 212, MinHeight: 104, Capability: style.ColorFast},
}
