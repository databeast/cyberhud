package styles

import (
	"github.com/databeast/cyberhud/display/style"
)

// System style declarations for 168x144 panels (landscape).
// Skeleton only (adaptiveBuild).

var MonoFast168x144Style = def{
	name: "mono-fast-168x144",
	reqs: style.SurfaceRequirements{MinWidth: 168, MinHeight: 144, Capability: style.MonoFast},
}
