package styles

import "github.com/databeast/cyberhud/display/style"

// GPIO style declarations for this panel group.
//
// These are hand-tweakable declarations over the shared layouts in core.go.

var MonoFast176x264Style = def{
	name: "mono-fast-176x264",
	reqs: style.SurfaceRequirements{MinWidth: 176, MinHeight: 264, Capability: style.MonoFast},
}

var MonoFast264x176Style = def{
	name: "mono-fast-264x176",
	reqs: style.SurfaceRequirements{MinWidth: 264, MinHeight: 176, Capability: style.MonoFast},
}

var GrayscaleSlow176x264Style = def{
	name: "grayscale-slow-176x264",
	reqs: style.SurfaceRequirements{MinWidth: 176, MinHeight: 264, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow264x176Style = def{
	name: "grayscale-slow-264x176",
	reqs: style.SurfaceRequirements{MinWidth: 264, MinHeight: 176, Capability: style.GrayscaleSlow},
}

var GrayscaleFast176x264Style = def{
	name: "grayscale-fast-176x264",
	reqs: style.SurfaceRequirements{MinWidth: 176, MinHeight: 264, Capability: style.GrayscaleFast},
}

var GrayscaleFast264x176Style = def{
	name: "grayscale-fast-264x176",
	reqs: style.SurfaceRequirements{MinWidth: 264, MinHeight: 176, Capability: style.GrayscaleFast},
}

var EinkStyle264x176 = def{
	name: "eink-264x176",
	reqs: style.SurfaceRequirements{MinWidth: 264, MinHeight: 176, Capability: style.MonoSlow},
	p:    Params{Layout: layoutEink},
}

var EinkStyle176x264 = def{
	name: "eink-176x264",
	reqs: style.SurfaceRequirements{MinWidth: 176, MinHeight: 264, Capability: style.MonoSlow},
	p:    Params{Layout: layoutEink},
}

var ColorSlow176x264Style = def{
	name: "color-slow-176x264",
	reqs: style.SurfaceRequirements{MinWidth: 176, MinHeight: 264, Capability: style.ColorSlow},
}

var ColorSlow264x176Style = def{
	name: "color-slow-264x176",
	reqs: style.SurfaceRequirements{MinWidth: 264, MinHeight: 176, Capability: style.ColorSlow},
}

var ColorFast176x264Style = def{
	name: "color-fast-176x264",
	reqs: style.SurfaceRequirements{MinWidth: 176, MinHeight: 264, Capability: style.ColorFast},
}

var ColorFast264x176Style = def{
	name: "color-fast-264x176",
	reqs: style.SurfaceRequirements{MinWidth: 264, MinHeight: 176, Capability: style.ColorFast},
}
