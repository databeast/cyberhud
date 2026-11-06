package styles

import (
	"github.com/databeast/cyberhud/display/style"
)

// System style declarations for 400x240 panels (landscape).
// Skeleton only (adaptiveBuild).

var MonoFast400x240Style = def{
	name: "mono-fast-400x240",
	reqs: style.SurfaceRequirements{MinWidth: 400, MinHeight: 240, Capability: style.MonoFast},
}
