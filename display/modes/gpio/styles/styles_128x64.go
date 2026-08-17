package styles

import "github.com/databeast/cyberhud/display/style"

// GPIO style declarations for this panel group.
//
// These are hand-tweakable declarations over the shared layouts in core.go.

var MonoSlow128x64Style = def{
	name: "mono-slow-128x64",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 64, Capability: style.MonoSlow},
}

var MonoSlow64x128Style = def{
	name: "mono-slow-64x128",
	reqs: style.SurfaceRequirements{MinWidth: 64, MinHeight: 128, Capability: style.MonoSlow},
}

var MonoStyle128x64 = def{
	name: "mono-128x64",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 64, Capability: style.MonoFast},
	p:    Params{Layout: layoutMonoRows},
}

var MonoStyle64x128 = def{
	name: "mono-64x128",
	reqs: style.SurfaceRequirements{MinWidth: 64, MinHeight: 128, Capability: style.MonoFast},
	p:    Params{Layout: layoutMonoRows},
}

var GrayscaleSlow128x64Style = def{
	name: "grayscale-slow-128x64",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 64, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow64x128Style = def{
	name: "grayscale-slow-64x128",
	reqs: style.SurfaceRequirements{MinWidth: 64, MinHeight: 128, Capability: style.GrayscaleSlow},
}

var GrayscaleFast128x64Style = def{
	name: "grayscale-fast-128x64",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 64, Capability: style.GrayscaleFast},
}

var GrayscaleFast64x128Style = def{
	name: "grayscale-fast-64x128",
	reqs: style.SurfaceRequirements{MinWidth: 64, MinHeight: 128, Capability: style.GrayscaleFast},
}

var ColorSlow128x64Style = def{
	name: "color-slow-128x64",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 64, Capability: style.ColorSlow},
}

var ColorSlow64x128Style = def{
	name: "color-slow-64x128",
	reqs: style.SurfaceRequirements{MinWidth: 64, MinHeight: 128, Capability: style.ColorSlow},
}

var ColorFast128x64Style = def{
	name: "color-fast-128x64",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 64, Capability: style.ColorFast},
}

var ColorFast64x128Style = def{
	name: "color-fast-64x128",
	reqs: style.SurfaceRequirements{MinWidth: 64, MinHeight: 128, Capability: style.ColorFast},
}
