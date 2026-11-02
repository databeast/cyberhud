package styles

import (
	"github.com/databeast/cyberhud/display/modes/zmq/content"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
)

// ColorFast128x128Style renders the ZMQ message feed on a 128x128 color-fast panel.
var ColorFast128x128Style = def{
	name: "color-fast-128x128",
	reqs: style.SurfaceRequirements{
		MinWidth:   128,
		MinHeight:  128,
		Capability: style.ColorFast,
	},
	p: Params{BuildFn: buildColorFast128x128},
}

func buildColorFast128x128(snap content.ZMQData, pol content.Policy, ctx style.StyleContext, _ def) style.ViewData {
	return buildColorItemsOnly(snap, pol, ctx, itemsConfig{
		tier:       tiercatalog.TierNormal,
		paddingPct: 0,
	})
}
