package styles

import "github.com/databeast/cyberhud/display/style"

// GPIO style declarations for this panel group.
//
// These are hand-tweakable declarations over the shared layouts in core.go.

var MonoFast300x400Style = def{
	name: "mono-fast-300x400",
	reqs: style.SurfaceRequirements{MinWidth: 300, MinHeight: 400, Capability: style.MonoFast},
}

var MonoFast400x300Style = def{
	name: "mono-fast-400x300",
	reqs: style.SurfaceRequirements{MinWidth: 400, MinHeight: 300, Capability: style.MonoFast},
}

var GrayscaleSlow300x400Style = def{
	name: "grayscale-slow-300x400",
	reqs: style.SurfaceRequirements{MinWidth: 300, MinHeight: 400, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow400x300Style = def{
	name: "grayscale-slow-400x300",
	reqs: style.SurfaceRequirements{MinWidth: 400, MinHeight: 300, Capability: style.GrayscaleSlow},
}

var GrayscaleFast300x400Style = def{
	name: "grayscale-fast-300x400",
	reqs: style.SurfaceRequirements{MinWidth: 300, MinHeight: 400, Capability: style.GrayscaleFast},
}

var GrayscaleFast400x300Style = def{
	name: "grayscale-fast-400x300",
	reqs: style.SurfaceRequirements{MinWidth: 400, MinHeight: 300, Capability: style.GrayscaleFast},
}

var EinkStyle400x300 = def{
	name: "eink-400x300",
	reqs: style.SurfaceRequirements{MinWidth: 400, MinHeight: 300, Capability: style.MonoSlow},
	p:    Params{Layout: layoutEink},
}

var EinkStyle296x128 = def{
	name: "eink-296x128",
	reqs: style.SurfaceRequirements{MinWidth: 296, MinHeight: 128, Capability: style.MonoSlow},
	p:    Params{Layout: layoutEink},
}

var EinkStyle300x400 = def{
	name: "eink-300x400",
	reqs: style.SurfaceRequirements{MinWidth: 300, MinHeight: 400, Capability: style.MonoSlow},
	p:    Params{Layout: layoutEink},
}

var ColorSlow300x400Style = def{
	name: "color-slow-300x400",
	reqs: style.SurfaceRequirements{MinWidth: 300, MinHeight: 400, Capability: style.ColorSlow},
}

var ColorSlow400x300Style = def{
	name: "color-slow-400x300",
	reqs: style.SurfaceRequirements{MinWidth: 400, MinHeight: 300, Capability: style.ColorSlow},
}

var ColorFast300x400Style = def{
	name: "color-fast-300x400",
	reqs: style.SurfaceRequirements{MinWidth: 300, MinHeight: 400, Capability: style.ColorFast},
}

var ColorFast400x300Style = def{
	name: "color-fast-400x300",
	reqs: style.SurfaceRequirements{MinWidth: 400, MinHeight: 300, Capability: style.ColorFast},
}
