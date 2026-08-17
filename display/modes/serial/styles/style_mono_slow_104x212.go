package styles

import (
	"github.com/databeast/cyberhud/display/modes/serial/source"
	"github.com/databeast/cyberhud/display/style"
)

// MonoSlow104x212Style renders serial data on a 104×212 monochrome slow (e-paper) panel.
var MonoSlow104x212Style = def{
	name: "mono-slow-104x212",
	reqs: style.SurfaceRequirements{
		MinWidth:   104,
		MinHeight:  212,
		Capability: style.MonoSlow,
	},
	p: Params{BuildFn: buildMonoSlow104x212},
}

func buildMonoSlow104x212(snap source.Snapshot, p source.Policy, ctx style.StyleContext, d def) style.ViewData {
	return buildMonoSlowItemsOnly(snap, p, ctx, d, 12)
}
