package styles

import (
	"github.com/databeast/cyberhud/display/modes/zmq/content"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
)

// ColorSlow400x300Style renders the ZMQ message feed on a 400x300 color-slow panel.
var ColorSlow400x300Style = def{
	name: "color-slow-400x300",
	reqs: style.SurfaceRequirements{
		MinWidth:   400,
		MinHeight:  300,
		Capability: style.ColorSlow,
	},
	p: Params{BuildFn: buildColorSlow400x300},
}

func buildColorSlow400x300(snap content.ZMQData, pol content.Policy, ctx style.StyleContext, _ def) style.ViewData {
	return buildColorItemsOnly(snap, pol, ctx, itemsConfig{
		tier:       tiercatalog.TierNormal,
		paddingPct: 2,
	})
}
