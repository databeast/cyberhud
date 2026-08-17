package styles

import "github.com/databeast/cyberhud/display/style"

// GPIO style declarations for this panel group.
//
// These are hand-tweakable declarations over the shared layouts in core.go.

var MonoSlow480x320Style = def{
	name: "mono-slow-480x320",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 320, Capability: style.MonoSlow},
}

var MonoSlow320x480Style = def{
	name: "mono-slow-320x480",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 480, Capability: style.MonoSlow},
}

var MonoFast480x320Style = def{
	name: "mono-fast-480x320",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 320, Capability: style.MonoFast},
}

var MonoFast320x480Style = def{
	name: "mono-fast-320x480",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 480, Capability: style.MonoFast},
}

var GrayscaleSlow320x480Style = def{
	name: "grayscale-slow-320x480",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 480, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow480x320Style = def{
	name: "grayscale-slow-480x320",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 320, Capability: style.GrayscaleSlow},
}

var GrayscaleFast480x320Style = def{
	name: "grayscale-fast-480x320",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 320, Capability: style.GrayscaleFast},
	p:    Params{Layout: layoutGrayscaleFast},
}

var GrayscaleFast320x480Style = def{
	name: "grayscale-fast-320x480",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 480, Capability: style.GrayscaleFast},
	p:    Params{Layout: layoutGrayscaleFast},
}

var ColorSlow320x480Style = def{
	name: "color-slow-320x480",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 480, Capability: style.ColorSlow},
}

var ColorSlow480x320Style = def{
	name: "color-slow-480x320",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 320, Capability: style.ColorSlow},
}

var ColorStyle480x320 = def{
	name: "color-480x320",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 320, Capability: style.ColorFast},
	p:    Params{Layout: layoutColorFastRows},
}

var ColorStyle320x480 = def{
	name: "color-320x480",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 480, Capability: style.ColorFast},
	p:    Params{Layout: layoutColorFastRows},
}
