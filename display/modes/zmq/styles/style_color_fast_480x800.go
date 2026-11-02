package styles

import (
	"github.com/databeast/cyberhud/display/modes/zmq/content"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
)

// ColorFast480x800Style renders the ZMQ message feed on a 480x800 color-fast panel.
var ColorFast480x800Style = def{
	name: "color-fast-480x800",
	reqs: style.SurfaceRequirements{
		MinWidth:   480,
		MinHeight:  800,
		Capability: style.ColorFast,
	},
	p: Params{BuildFn: buildColorFast480x800},
}

func buildColorFast480x800(snap content.ZMQData, pol content.Policy, ctx style.StyleContext, _ def) style.ViewData {
	return buildColorItemsOnly(snap, pol, ctx, itemsConfig{
		tier:       tiercatalog.TierLarge,
		paddingPct: 0,
	})
}
