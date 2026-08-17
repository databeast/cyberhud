package styles

import (
	"github.com/databeast/cyberhud/display/modes/serial/source"
	"github.com/databeast/cyberhud/display/style"
)

// MonoSlow250x122Style renders serial data on a 250×122 monochrome slow (e-paper) panel.
var MonoSlow250x122Style = def{
	name: "mono-slow-250x122",
	reqs: style.SurfaceRequirements{
		MinWidth:   250,
		MinHeight:  122,
		Capability: style.MonoSlow,
	},
	p: Params{BuildFn: buildMonoSlow250x122},
}

func buildMonoSlow250x122(snap source.Snapshot, p source.Policy, ctx style.StyleContext, d def) style.ViewData {
	return buildMonoSlowItemsOnly(snap, p, ctx, d, 7)
}
