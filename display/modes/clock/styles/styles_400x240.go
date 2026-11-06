package styles

import (
	"github.com/databeast/cyberhud/display/style"
)

// Clock style declarations for 400x240 panels (Sharp memory LCD).

var MonoSmall400x240Style = styleDef{
	name: "mono-fast-400x240",
	reqs: style.SurfaceRequirements{MinWidth: 400, MinHeight: 240, Capability: style.MonoFast},
}
