package styles

import (
	"github.com/databeast/cyberhud/display/modes/zmq/content"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
)

// ColorFast320x240Style renders the ZMQ message feed on a 320x240 color-fast panel.
var ColorFast320x240Style = def{
	name: "color-fast-320x240",
	reqs: style.SurfaceRequirements{
		MinWidth:   320,
		MinHeight:  240,
		Capability: style.ColorFast,
	},
	p: Params{BuildFn: buildColorFast320x240},
}

func buildColorFast320x240(snap content.ZMQData, pol content.Policy, ctx style.StyleContext, _ def) style.ViewData {
	return buildColorItemsOnly(snap, pol, ctx, itemsConfig{
		tier:       tiercatalog.TierNormal,
		paddingPct: 0,
	})
}
