package styles

import (
	"github.com/databeast/cyberhud/display/style"
)

// System style declarations for 16x8 LED matrix panels.
// Skeleton only (adaptiveBuild).

var GrayscaleFast16x8Style = def{
	name: "grayscale-fast-16x8",
	reqs: style.SurfaceRequirements{MinWidth: 16, MinHeight: 8, Capability: style.GrayscaleFast},
}
