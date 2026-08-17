package styles

import (
	"github.com/databeast/cyberhud/display/modes/serial/source"
	"github.com/databeast/cyberhud/display/style"
)

// MonoSlow200x200Style renders serial data on a 200×200 monochrome slow (e-paper) panel.
var MonoSlow200x200Style = def{
	name: "mono-slow-200x200",
	reqs: style.SurfaceRequirements{
		MinWidth:   200,
		MinHeight:  200,
		Capability: style.MonoSlow,
	},
	p: Params{BuildFn: buildMonoSlow200x200},
}

func buildMonoSlow200x200(snap source.Snapshot, p source.Policy, ctx style.StyleContext, d def) style.ViewData {
	return buildMonoSlowItemsOnly(snap, p, ctx, d, 12)
}
