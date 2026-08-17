package styles

import (
	"github.com/databeast/cyberhud/display/modes/serial/source"
	"github.com/databeast/cyberhud/display/style"
)

var GrayscaleFast160x80Style = def{
	name: "grayscale-fast-160x80",
	reqs: style.SurfaceRequirements{
		MinWidth:   160,
		MinHeight:  80,
		Capability: style.GrayscaleFast,
	},
	p: Params{BuildFn: buildGrayscaleFast160x80},
}

func buildGrayscaleFast160x80(snap source.Snapshot, p source.Policy, ctx style.StyleContext, d def) style.ViewData {
	return buildGrayscaleItemsOnly(snap, p, ctx, d, 5, false)
}
