package styles

import (
	"github.com/databeast/cyberhud/display/modes/zmq/content"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
)

// MonoSlow250x122Style renders the ZMQ message feed on a 250x122 mono-slow panel.
var MonoSlow250x122Style = def{
	name: "mono-slow-250x122",
	reqs: style.SurfaceRequirements{
		MinWidth:   250,
		MinHeight:  122,
		Capability: style.MonoSlow,
	},
	p: Params{BuildFn: buildMonoSlow250x122},
}

func buildMonoSlow250x122(snap content.ZMQData, pol content.Policy, ctx style.StyleContext, _ def) style.ViewData {
	return buildItemsOnly(snap, pol, ctx, itemsConfig{
		tier:       tiercatalog.TierNormal,
		paddingPct: 2,
	})
}
