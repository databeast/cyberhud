package styles

import (
	"github.com/databeast/cyberhud/display/modes/serial/source"
	"github.com/databeast/cyberhud/display/style"
)

var ColorFast320x480Style = def{
	name: "color-320x480",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 480, Capability: style.ColorFast},
	p:    Params{BuildFn: buildColorFast320x480},
}

func buildColorFast320x480(snap source.Snapshot, p source.Policy, ctx style.StyleContext, d def) style.ViewData {
	return buildColorItemsOnly(snap, p, ctx, d, 20, 0, false)
}
