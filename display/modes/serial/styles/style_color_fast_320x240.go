package styles

import (
	"github.com/databeast/cyberhud/display/modes/serial/source"
	"github.com/databeast/cyberhud/display/style"
)

var ColorFast320x240Style = def{
	name: "color-320x240",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 240, Capability: style.ColorFast},
	p:    Params{BuildFn: buildColorFast320x240},
}

func buildColorFast320x240(snap source.Snapshot, p source.Policy, ctx style.StyleContext, d def) style.ViewData {
	return buildColorItemsOnly(snap, p, ctx, d, 15, 0, false)
}
