package styles

import (
	"github.com/databeast/cyberhud/display/modes/serial/source"
	"github.com/databeast/cyberhud/display/style"
)

var GrayscaleSlow240x135Style = def{
	name: "grayscale-slow-240x135",
	reqs: style.SurfaceRequirements{
		MinWidth:   240,
		MinHeight:  135,
		Capability: style.GrayscaleSlow,
	},
	p: Params{BuildFn: buildGrayscaleSlow240x135},
}

func buildGrayscaleSlow240x135(snap source.Snapshot, p source.Policy, ctx style.StyleContext, d def) style.ViewData {
	return buildGrayscaleItemsOnly(snap, p, ctx, d, 8, true)
}
