package styles

import (
	"github.com/databeast/cyberhud/display/modes/zmq/content"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
)

// ColorFast240x320Style renders the ZMQ message feed on a 240x320 color-fast panel.
var ColorFast240x320Style = def{
	name: "color-fast-240x320",
	reqs: style.SurfaceRequirements{
		MinWidth:   240,
		MinHeight:  320,
		Capability: style.ColorFast,
	},
	p: Params{BuildFn: buildColorFast240x320},
}

func buildColorFast240x320(snap content.ZMQData, pol content.Policy, ctx style.StyleContext, _ def) style.ViewData {
	return buildColorItemsOnly(snap, pol, ctx, itemsConfig{
		tier:       tiercatalog.TierNormal,
		paddingPct: 0,
	})
}
