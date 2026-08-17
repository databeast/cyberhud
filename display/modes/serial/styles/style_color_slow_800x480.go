package styles

import (
	"github.com/databeast/cyberhud/display/modes/serial/source"
	"github.com/databeast/cyberhud/display/style"
)

var ColorSlow800x480Style = def{
	name: "color-slow-800x480",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildColorSlow800x480},
}

func buildColorSlow800x480(snap source.Snapshot, p source.Policy, ctx style.StyleContext, d def) style.ViewData {
	return buildColorItemsOnly(snap, p, ctx, d, 30, 2, true)
}
