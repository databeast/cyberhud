package styles

import (
	"github.com/databeast/cyberhud/display/modes/serial/source"
	"github.com/databeast/cyberhud/display/style"
)

// MonoSlow176x264Style renders serial data on a 176×264 monochrome slow (e-paper) panel.
var MonoSlow176x264Style = def{
	name: "mono-slow-176x264",
	reqs: style.SurfaceRequirements{
		MinWidth:   176,
		MinHeight:  264,
		Capability: style.MonoSlow,
	},
	p: Params{BuildFn: buildMonoSlow176x264},
}

func buildMonoSlow176x264(snap source.Snapshot, p source.Policy, ctx style.StyleContext, d def) style.ViewData {
	return buildMonoSlowItemsOnly(snap, p, ctx, d, 16)
}
