package styles

import (
	"github.com/databeast/cyberhud/display/modes/serial/source"
	"github.com/databeast/cyberhud/display/style"
)

// MonoSlow264x176Style renders serial data on a 264×176 monochrome slow (e-paper) panel.
var MonoSlow264x176Style = def{
	name: "mono-slow-264x176",
	reqs: style.SurfaceRequirements{
		MinWidth:   264,
		MinHeight:  176,
		Capability: style.MonoSlow,
	},
	p: Params{BuildFn: buildMonoSlow264x176},
}

func buildMonoSlow264x176(snap source.Snapshot, p source.Policy, ctx style.StyleContext, d def) style.ViewData {
	return buildMonoSlowItemsOnly(snap, p, ctx, d, 10)
}
