package styles

import "github.com/databeast/cyberhud/display/style"

// GPIO style declarations for this panel group.
//
// These are hand-tweakable declarations over the shared layouts in core.go.

var MonoFast200x200Style = def{
	name: "mono-fast-200x200",
	reqs: style.SurfaceRequirements{MinWidth: 200, MinHeight: 200, Capability: style.MonoFast},
}

var GrayscaleSlow200x200Style = def{
	name: "grayscale-slow-200x200",
	reqs: style.SurfaceRequirements{MinWidth: 200, MinHeight: 200, Capability: style.GrayscaleSlow},
}

var GrayscaleFast200x200Style = def{
	name: "grayscale-fast-200x200",
	reqs: style.SurfaceRequirements{MinWidth: 200, MinHeight: 200, Capability: style.GrayscaleFast},
}

var EinkStyle200x200 = def{
	name: "eink-200x200",
	reqs: style.SurfaceRequirements{MinWidth: 200, MinHeight: 200, Capability: style.MonoSlow},
	p:    Params{Layout: layoutEink},
}

var ColorSlow200x200Style = def{
	name: "color-slow-200x200",
	reqs: style.SurfaceRequirements{MinWidth: 200, MinHeight: 200, Capability: style.ColorSlow},
}

var ColorFast200x200Style = def{
	name: "color-fast-200x200",
	reqs: style.SurfaceRequirements{MinWidth: 200, MinHeight: 200, Capability: style.ColorFast},
}
