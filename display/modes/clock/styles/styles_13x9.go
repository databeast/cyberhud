package styles

import (
	"github.com/databeast/cyberhud/display/style"
)

// Clock style declarations for 13x9 panels (color LED matrix).

var ColorSmall13x9Style = styleDef{
	name: "color-fast-13x9",
	reqs: style.SurfaceRequirements{MinWidth: 13, MinHeight: 9, Capability: style.ColorFast},
	p:    Params{Fast: true, Color: true},
}
