package styles

import (
	"github.com/databeast/cyberhud/display/modes/serial/source"
	"github.com/databeast/cyberhud/display/style"
)

var ColorFast240x320Style = def{
	name: "color-240x320",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 320, Capability: style.ColorFast},
	p:    Params{BuildFn: buildColorFast240x320},
}

func buildColorFast240x320(snap source.Snapshot, p source.Policy, ctx style.StyleContext, d def) style.ViewData {
	return buildColorItemsOnly(snap, p, ctx, d, 20, 0, false)
}
