package styles

import (
	"github.com/databeast/cyberhud/display/modes/serial/source"
	"github.com/databeast/cyberhud/display/style"
)

var GrayscaleFast320x240Style = def{
	name: "grayscale-fast-320x240",
	reqs: style.SurfaceRequirements{
		MinWidth:   320,
		MinHeight:  240,
		Capability: style.GrayscaleFast,
	},
	p: Params{BuildFn: buildGrayscaleFast320x240},
}

func buildGrayscaleFast320x240(snap source.Snapshot, p source.Policy, ctx style.StyleContext, d def) style.ViewData {
	return buildGrayscaleItemsOnly(snap, p, ctx, d, 15, false)
}
