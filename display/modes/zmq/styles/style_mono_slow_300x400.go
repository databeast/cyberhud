package styles

import (
	"github.com/databeast/cyberhud/display/modes/zmq/content"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
)

// MonoSlow300x400Style renders the ZMQ message feed on a 300x400 mono-slow panel.
var MonoSlow300x400Style = def{
	name: "mono-slow-300x400",
	reqs: style.SurfaceRequirements{
		MinWidth:   300,
		MinHeight:  400,
		Capability: style.MonoSlow,
	},
	p: Params{BuildFn: buildMonoSlow300x400},
}

func buildMonoSlow300x400(snap content.ZMQData, pol content.Policy, ctx style.StyleContext, _ def) style.ViewData {
	return buildItemsOnly(snap, pol, ctx, itemsConfig{
		tier:       tiercatalog.TierLarge,
		paddingPct: 2,
	})
}
