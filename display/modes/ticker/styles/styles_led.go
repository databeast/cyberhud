package styles

import (
	"github.com/databeast/cyberhud/display/style"
)

// Skeleton registrations for 128x64 OLED and LED matrix panels.

var TickerGrayscaleFast128x64Style = def{
	name: "grayscale-fast-128x64",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 64, Capability: style.GrayscaleFast},
}

var TickerGrayscaleFast15x7Style = def{
	name: "grayscale-fast-15x7",
	reqs: style.SurfaceRequirements{MinWidth: 15, MinHeight: 7, Capability: style.GrayscaleFast},
}

var TickerGrayscaleFast16x8Style = def{
	name: "grayscale-fast-16x8",
	reqs: style.SurfaceRequirements{MinWidth: 16, MinHeight: 8, Capability: style.GrayscaleFast},
}

var TickerColorFast13x9Style = def{
	name: "color-fast-13x9",
	reqs: style.SurfaceRequirements{MinWidth: 13, MinHeight: 9, Capability: style.ColorFast},
}
