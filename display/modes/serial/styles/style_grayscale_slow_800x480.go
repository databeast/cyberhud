package styles

import (
	"github.com/databeast/cyberhud/display/modes/serial/source"
	"github.com/databeast/cyberhud/display/style"
)

var GrayscaleSlow800x480Style = def{
	name: "grayscale-slow-800x480",
	reqs: style.SurfaceRequirements{
		MinWidth:   800,
		MinHeight:  480,
		Capability: style.GrayscaleSlow,
	},
	p: Params{BuildFn: buildGrayscaleSlow800x480},
}

func buildGrayscaleSlow800x480(snap source.Snapshot, p source.Policy, ctx style.StyleContext, d def) style.ViewData {
	return buildGrayscaleItemsOnly(snap, p, ctx, d, 30, true)
}
