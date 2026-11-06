package styles

import (
	"github.com/databeast/cyberhud/display/style"
)

// Clock style declarations for 16x8 panels (grayscale LED matrix).

var GrayscaleFast16x8Style = styleDef{
	name: "grayscale-fast-16x8",
	reqs: style.SurfaceRequirements{MinWidth: 16, MinHeight: 8, Capability: style.GrayscaleFast},
	p:    Params{Fast: true},
}
