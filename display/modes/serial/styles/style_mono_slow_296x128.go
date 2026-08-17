package styles

import (
	"github.com/databeast/cyberhud/display/modes/serial/source"
	"github.com/databeast/cyberhud/display/style"
)

// MonoSlow296x128Style renders serial data on a 296×128 monochrome slow (e-paper) panel.
var MonoSlow296x128Style = def{
	name: "mono-slow-296x128",
	reqs: style.SurfaceRequirements{
		MinWidth:   296,
		MinHeight:  128,
		Capability: style.MonoSlow,
	},
	p: Params{BuildFn: buildMonoSlow296x128},
}

func buildMonoSlow296x128(snap source.Snapshot, p source.Policy, ctx style.StyleContext, d def) style.ViewData {
	return buildMonoSlowItemsOnly(snap, p, ctx, d, 7)
}
