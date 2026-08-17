package styles

import (
	"github.com/databeast/cyberhud/display/modes/serial/source"
	"github.com/databeast/cyberhud/display/style"
)

var ColorFast480x800Style = def{
	name: "color-480x800",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.ColorFast},
	p:    Params{BuildFn: buildColorFast480x800},
}

func buildColorFast480x800(snap source.Snapshot, p source.Policy, ctx style.StyleContext, d def) style.ViewData {
	return buildColorItemsOnly(snap, p, ctx, d, 30, 0, false)
}
