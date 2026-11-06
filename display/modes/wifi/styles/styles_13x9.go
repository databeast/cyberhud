package styles

import (
	"github.com/databeast/cyberhud/display/style"
)

// WiFi style declarations for 13x9 LED matrix panels (skeleton only).

var ColorFast13x9Style = def{
	name: "color-fast-13x9",
	reqs: style.SurfaceRequirements{MinWidth: 13, MinHeight: 9, Capability: style.ColorFast},
	p:    Params{Fast: true, Color: true},
}
