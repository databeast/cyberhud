package styles

import (
	"github.com/databeast/cyberhud/display/modes/serial/source"
	"github.com/databeast/cyberhud/display/style"
)

var GrayscaleSlow480x320Style = def{
	name: "grayscale-slow-480x320",
	reqs: style.SurfaceRequirements{
		MinWidth:   480,
		MinHeight:  320,
		Capability: style.GrayscaleSlow,
	},
	p: Params{BuildFn: buildGrayscaleSlow480x320},
}

func buildGrayscaleSlow480x320(snap source.Snapshot, p source.Policy, ctx style.StyleContext, d def) style.ViewData {
	return buildGrayscaleItemsOnly(snap, p, ctx, d, 20, true)
}
