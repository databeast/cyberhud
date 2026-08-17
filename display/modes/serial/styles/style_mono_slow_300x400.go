package styles

import (
	"github.com/databeast/cyberhud/display/modes/serial/source"
	"github.com/databeast/cyberhud/display/style"
)

// MonoSlow300x400Style renders serial data on a 300×400 monochrome slow (e-paper) panel.
var MonoSlow300x400Style = def{
	name: "mono-slow-300x400",
	reqs: style.SurfaceRequirements{
		MinWidth:   300,
		MinHeight:  400,
		Capability: style.MonoSlow,
	},
	p: Params{BuildFn: buildMonoSlow300x400},
}

func buildMonoSlow300x400(snap source.Snapshot, p source.Policy, ctx style.StyleContext, d def) style.ViewData {
	return buildMonoSlowItemsOnly(snap, p, ctx, d, 24)
}
