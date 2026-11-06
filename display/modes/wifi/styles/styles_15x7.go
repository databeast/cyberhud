package styles

import (
	"github.com/databeast/cyberhud/display/style"
)

// WiFi style declarations for 15x7 LED matrix panels (skeleton only).

var GrayscaleFast15x7Style = def{
	name: "grayscale-fast-15x7",
	reqs: style.SurfaceRequirements{MinWidth: 15, MinHeight: 7, Capability: style.GrayscaleFast},
	p:    Params{Fast: true},
}
