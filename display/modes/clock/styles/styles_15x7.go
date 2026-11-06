package styles

import (
	"github.com/databeast/cyberhud/display/style"
)

// Clock style declarations for 15x7 panels (grayscale LED matrix).

var GrayscaleFast15x7Style = styleDef{
	name: "grayscale-fast-15x7",
	reqs: style.SurfaceRequirements{MinWidth: 15, MinHeight: 7, Capability: style.GrayscaleFast},
	p:    Params{Fast: true},
}
