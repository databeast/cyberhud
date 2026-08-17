package styles

import "github.com/databeast/cyberhud/display/style"

// WiFi style declarations for 296x128 panels.

var MonoSlow296x128Style = def{
	name: "mono-slow-296x128",
	reqs: style.SurfaceRequirements{MinWidth: 296, MinHeight: 128, Capability: style.MonoSlow},
}

var MonoFast296x128Style = def{
	name: "mono-fast-296x128",
	reqs: style.SurfaceRequirements{MinWidth: 296, MinHeight: 128, Capability: style.MonoFast},
	p:    Params{Fast: true},
}

var GrayscaleSlow296x128Style = def{
	name: "grayscale-slow-296x128",
	reqs: style.SurfaceRequirements{MinWidth: 296, MinHeight: 128, Capability: style.GrayscaleSlow},
}

var GrayscaleFast296x128Style = def{
	name: "grayscale-fast-296x128",
	reqs: style.SurfaceRequirements{MinWidth: 296, MinHeight: 128, Capability: style.GrayscaleFast},
	p:    Params{Fast: true},
}

var ColorSlow296x128Style = def{
	name: "color-slow-296x128",
	reqs: style.SurfaceRequirements{MinWidth: 296, MinHeight: 128, Capability: style.ColorSlow},
}

var ColorFast296x128Style = def{
	name: "color-fast-296x128",
	reqs: style.SurfaceRequirements{MinWidth: 296, MinHeight: 128, Capability: style.ColorFast},
	p:    Params{Fast: true, Color: true},
}

// Portrait variant (128×296)

var MonoSlow128x296Style = def{
	name: "mono-slow-128x296",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 296, Capability: style.MonoSlow},
}

var MonoFast128x296Style = def{
	name: "mono-fast-128x296",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 296, Capability: style.MonoFast},
	p:    Params{Fast: true},
}

var GrayscaleSlow128x296Style = def{
	name: "grayscale-slow-128x296",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 296, Capability: style.GrayscaleSlow},
}

var GrayscaleFast128x296Style = def{
	name: "grayscale-fast-128x296",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 296, Capability: style.GrayscaleFast},
	p:    Params{Fast: true},
}

var ColorSlow128x296Style = def{
	name: "color-slow-128x296",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 296, Capability: style.ColorSlow},
}

var ColorFast128x296Style = def{
	name: "color-fast-128x296",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 296, Capability: style.ColorFast},
	p:    Params{Fast: true, Color: true},
}
