package styles

import (
	"github.com/databeast/cyberhud/display/style"
)

// Clock style declarations for 480x320 panels.
//
// Each entry is a hand-tweakable declaration over the core layouts in
// core.go: adjust its Params (font tier, row toggles, Fast/Color layout
// selection) or attach a bespoke BuildFn to tune the clock's look for
// this specific resolution and PPI.

var MonoSlow480x320Style = styleDef{
	name: "mono-slow-480x320",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 320, Capability: style.MonoSlow},
}

var MonoSlow320x480Style = styleDef{
	name: "mono-slow-320x480",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 480, Capability: style.MonoSlow},
}

var GrayscaleSlow480x320Style = styleDef{
	name: "grayscale-slow-480x320",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 320, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow320x480Style = styleDef{
	name: "grayscale-slow-320x480",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 480, Capability: style.GrayscaleSlow},
}

var GrayscaleFast480x320Style = styleDef{
	name: "grayscale-fast-480x320",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 320, Capability: style.GrayscaleFast},
	p:    Params{Fast: true},
}

var GrayscaleFast320x480Style = styleDef{
	name: "grayscale-fast-320x480",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 480, Capability: style.GrayscaleFast},
	p:    Params{Fast: true},
}

var ColorSlow480x320Style = styleDef{
	name: "color-slow-480x320",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 320, Capability: style.ColorSlow},
}

var ColorSlow320x480Style = styleDef{
	name: "color-slow-320x480",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 480, Capability: style.ColorSlow},
}

var ColorLarge480x320Style = styleDef{
	name: "color-480x320",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 320, Capability: style.ColorFast},
	p:    Params{Fast: true, Color: true},
}

var ColorLarge320x480Style = styleDef{
	name: "color-320x480",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 480, Capability: style.ColorFast},
	p:    Params{Fast: true, Color: true},
}
