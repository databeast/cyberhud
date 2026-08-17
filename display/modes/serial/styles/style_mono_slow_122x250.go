package styles

import (
	"github.com/databeast/cyberhud/display/modes/serial/source"
	"github.com/databeast/cyberhud/display/style"
)

// MonoSlow122x250Style renders serial data on a 122×250 monochrome slow (e-paper) panel.
var MonoSlow122x250Style = def{
	name: "mono-slow-122x250",
	reqs: style.SurfaceRequirements{
		MinWidth:   122,
		MinHeight:  250,
		Capability: style.MonoSlow,
	},
	p: Params{BuildFn: buildMonoSlow122x250},
}

func buildMonoSlow122x250(snap source.Snapshot, p source.Policy, ctx style.StyleContext, d def) style.ViewData {
	return buildMonoSlowItemsOnly(snap, p, ctx, d, 14)
}
