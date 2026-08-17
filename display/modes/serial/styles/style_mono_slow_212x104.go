package styles

import (
	"github.com/databeast/cyberhud/display/modes/serial/source"
	"github.com/databeast/cyberhud/display/style"
)

// MonoSlow212x104Style renders serial data on a 212×104 monochrome slow (e-paper) panel.
var MonoSlow212x104Style = def{
	name: "mono-slow-212x104",
	reqs: style.SurfaceRequirements{
		MinWidth:   212,
		MinHeight:  104,
		Capability: style.MonoSlow,
	},
	p: Params{BuildFn: buildMonoSlow212x104},
}

func buildMonoSlow212x104(snap source.Snapshot, p source.Policy, ctx style.StyleContext, d def) style.ViewData {
	return buildMonoSlowItemsOnly(snap, p, ctx, d, 6)
}
