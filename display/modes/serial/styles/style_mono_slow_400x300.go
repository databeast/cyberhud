package styles

import (
	"github.com/databeast/cyberhud/display/modes/serial/source"
	"github.com/databeast/cyberhud/display/style"
)

// MonoSlow400x300Style renders serial data on a 400×300 monochrome slow (e-paper) panel.
var MonoSlow400x300Style = def{
	name: "mono-slow-400x300",
	reqs: style.SurfaceRequirements{
		MinWidth:   400,
		MinHeight:  300,
		Capability: style.MonoSlow,
	},
	p: Params{BuildFn: buildMonoSlow400x300},
}

func buildMonoSlow400x300(snap source.Snapshot, p source.Policy, ctx style.StyleContext, d def) style.ViewData {
	return buildMonoSlowItemsOnly(snap, p, ctx, d, 18)
}
