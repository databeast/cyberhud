package styles

import (
	"github.com/databeast/cyberhud/display/modes/zmq/content"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
)

// ColorSlow800x480Style renders the ZMQ message feed on a 800x480 color-slow panel.
var ColorSlow800x480Style = def{
	name: "color-slow-800x480",
	reqs: style.SurfaceRequirements{
		MinWidth:   800,
		MinHeight:  480,
		Capability: style.ColorSlow,
	},
	p: Params{BuildFn: buildColorSlow800x480},
}

func buildColorSlow800x480(snap content.ZMQData, pol content.Policy, ctx style.StyleContext, _ def) style.ViewData {
	return buildColorItemsOnly(snap, pol, ctx, itemsConfig{
		tier:       tiercatalog.TierLarge,
		paddingPct: 2,
	})
}
