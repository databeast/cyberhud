package styles

import (
	"github.com/databeast/cyberhud/display/modes/zmq/content"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
)

// ColorFast160x80Style renders the ZMQ message feed on a 160x80 color-fast panel.
var ColorFast160x80Style = def{
	name: "color-fast-160x80",
	reqs: style.SurfaceRequirements{
		MinWidth:   160,
		MinHeight:  80,
		Capability: style.ColorFast,
	},
	p: Params{BuildFn: buildColorFast160x80},
}

func buildColorFast160x80(snap content.ZMQData, pol content.Policy, ctx style.StyleContext, _ def) style.ViewData {
	return buildColorItemsOnly(snap, pol, ctx, itemsConfig{
		tier:       tiercatalog.TierSmall,
		paddingPct: 0,
	})
}
