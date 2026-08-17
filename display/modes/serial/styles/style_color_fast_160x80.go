package styles

import (
	"github.com/databeast/cyberhud/display/modes/serial/source"
	"github.com/databeast/cyberhud/display/style"
)

var ColorFast160x80Style = def{
	name: "color-160x80",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.ColorFast},
	p:    Params{BuildFn: buildColorFast160x80},
}

func buildColorFast160x80(snap source.Snapshot, p source.Policy, ctx style.StyleContext, d def) style.ViewData {
	return buildColorItemsOnly(snap, p, ctx, d, 5, 0, false)
}
