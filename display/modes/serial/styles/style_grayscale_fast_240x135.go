package styles

import (
	"github.com/databeast/cyberhud/display/modes/serial/source"
	"github.com/databeast/cyberhud/display/style"
)

var GrayscaleFast240x135Style = def{
	name: "grayscale-fast-240x135",
	reqs: style.SurfaceRequirements{
		MinWidth:   240,
		MinHeight:  135,
		Capability: style.GrayscaleFast,
	},
	p: Params{BuildFn: buildGrayscaleFast240x135},
}

func buildGrayscaleFast240x135(snap source.Snapshot, p source.Policy, ctx style.StyleContext, d def) style.ViewData {
	return buildGrayscaleItemsOnly(snap, p, ctx, d, 8, false)
}
