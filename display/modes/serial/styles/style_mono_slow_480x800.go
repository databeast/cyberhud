package styles

import (
	"github.com/databeast/cyberhud/display/modes/serial/source"
	"github.com/databeast/cyberhud/display/style"
)

// MonoSlow480x800Style renders serial data on a 480×800 monochrome slow (e-paper) panel.
var MonoSlow480x800Style = def{
	name: "mono-slow-480x800",
	reqs: style.SurfaceRequirements{
		MinWidth:   480,
		MinHeight:  800,
		Capability: style.MonoSlow,
	},
	p: Params{BuildFn: buildMonoSlow480x800},
}

func buildMonoSlow480x800(snap source.Snapshot, p source.Policy, ctx style.StyleContext, d def) style.ViewData {
	return buildMonoSlowItemsOnly(snap, p, ctx, d, 30)
}
