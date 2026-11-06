package styles

import (
	"github.com/databeast/cyberhud/display/modes/zmq/content"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
)

// ColorSlow122x250Style renders the ZMQ message feed on a 122x250 color-slow panel.
var ColorSlow122x250Style = def{
	name: "color-slow-122x250",
	reqs: style.SurfaceRequirements{
		MinWidth:   122,
		MinHeight:  250,
		Capability: style.ColorSlow,
	},
	p: Params{BuildFn: buildColorSlow122x250},
}

func buildColorSlow122x250(snap content.ZMQData, pol content.Policy, ctx style.StyleContext, _ def) style.ViewData {
	return buildColorItemsOnly(snap, pol, ctx, itemsConfig{
		tier:       tiercatalog.TierNormal,
		paddingPct: 2,
	})
}
