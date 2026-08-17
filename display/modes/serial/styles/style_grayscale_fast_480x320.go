package styles

import (
	"github.com/databeast/cyberhud/display/modes/serial/source"
	"github.com/databeast/cyberhud/display/style"
)

var GrayscaleFast480x320Style = def{
	name: "grayscale-fast-480x320",
	reqs: style.SurfaceRequirements{
		MinWidth:   480,
		MinHeight:  320,
		Capability: style.GrayscaleFast,
	},
	p: Params{BuildFn: buildGrayscaleFast480x320},
}

func buildGrayscaleFast480x320(snap source.Snapshot, p source.Policy, ctx style.StyleContext, d def) style.ViewData {
	return buildGrayscaleItemsOnly(snap, p, ctx, d, 20, false)
}
