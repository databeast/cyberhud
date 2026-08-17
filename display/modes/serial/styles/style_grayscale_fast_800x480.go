package styles

import (
	"github.com/databeast/cyberhud/display/modes/serial/source"
	"github.com/databeast/cyberhud/display/style"
)

var GrayscaleFast800x480Style = def{
	name: "grayscale-fast-800x480",
	reqs: style.SurfaceRequirements{
		MinWidth:   800,
		MinHeight:  480,
		Capability: style.GrayscaleFast,
	},
	p: Params{BuildFn: buildGrayscaleFast800x480},
}

func buildGrayscaleFast800x480(snap source.Snapshot, p source.Policy, ctx style.StyleContext, d def) style.ViewData {
	return buildGrayscaleItemsOnly(snap, p, ctx, d, 30, false)
}
