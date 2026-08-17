package styles

import (
	"github.com/databeast/cyberhud/display/modes/serial/source"
	"github.com/databeast/cyberhud/display/style"
)

var GrayscaleFast240x240Style = def{
	name: "grayscale-fast-240x240",
	reqs: style.SurfaceRequirements{
		MinWidth:   240,
		MinHeight:  240,
		Capability: style.GrayscaleFast,
	},
	p: Params{BuildFn: buildGrayscaleFast240x240},
}

func buildGrayscaleFast240x240(snap source.Snapshot, p source.Policy, ctx style.StyleContext, d def) style.ViewData {
	return buildGrayscaleItemsOnly(snap, p, ctx, d, 15, false)
}
