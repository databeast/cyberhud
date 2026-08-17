package styles

import (
	"github.com/databeast/cyberhud/display/modes/serial/source"
	"github.com/databeast/cyberhud/display/style"
)

var GrayscaleSlow240x240Style = def{
	name: "grayscale-slow-240x240",
	reqs: style.SurfaceRequirements{
		MinWidth:   240,
		MinHeight:  240,
		Capability: style.GrayscaleSlow,
	},
	p: Params{BuildFn: buildGrayscaleSlow240x240},
}

func buildGrayscaleSlow240x240(snap source.Snapshot, p source.Policy, ctx style.StyleContext, d def) style.ViewData {
	return buildGrayscaleItemsOnly(snap, p, ctx, d, 15, true)
}
