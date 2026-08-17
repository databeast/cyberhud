package styles

import "github.com/databeast/cyberhud/display/style"

// GPIO style declarations for this panel group.
//
// These are hand-tweakable declarations over the shared layouts in core.go.

var MonoSlow160x128Style = def{
	name: "mono-slow-160x128",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 128, Capability: style.MonoSlow},
}

var MonoSlow128x160Style = def{
	name: "mono-slow-128x160",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 160, Capability: style.MonoSlow},
}

var MonoFast160x128Style = def{
	name: "mono-fast-160x128",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 128, Capability: style.MonoFast},
}

var MonoFast128x160Style = def{
	name: "mono-fast-128x160",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 160, Capability: style.MonoFast},
}

var GrayscaleSlow128x160Style = def{
	name: "grayscale-slow-128x160",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 160, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow160x128Style = def{
	name: "grayscale-slow-160x128",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 128, Capability: style.GrayscaleSlow},
}

var GrayscaleFast128x160Style = def{
	name: "grayscale-fast-128x160",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 160, Capability: style.GrayscaleFast},
	p:    Params{Layout: layoutGrayscaleFast},
}

var ColorSlow128x160Style = def{
	name: "color-slow-128x160",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 160, Capability: style.ColorSlow},
}

var ColorSlow160x128Style = def{
	name: "color-slow-160x128",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 128, Capability: style.ColorSlow},
}

var ColorStyle160x128 = def{
	name: "color-160x128",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 128, Capability: style.ColorFast},
	p:    Params{Layout: layoutColorFastRows},
}

var ColorStyle128x160 = def{
	name: "color-128x160",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 160, Capability: style.ColorFast},
	p:    Params{Layout: layoutColorFastRows},
}
