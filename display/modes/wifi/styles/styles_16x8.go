package styles

import (
	"github.com/databeast/cyberhud/display/style"
)

// WiFi style declarations for 16x8 LED matrix panels (skeleton only).

var GrayscaleFast16x8Style = def{
	name: "grayscale-fast-16x8",
	reqs: style.SurfaceRequirements{MinWidth: 16, MinHeight: 8, Capability: style.GrayscaleFast},
	p:    Params{Fast: true},
}
