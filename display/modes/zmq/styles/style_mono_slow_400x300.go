package styles

import (
	"github.com/databeast/cyberhud/display/modes/zmq/content"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
)

// MonoSlow400x300Style renders the ZMQ message feed on a 400x300 mono-slow panel.
var MonoSlow400x300Style = def{
	name: "mono-slow-400x300",
	reqs: style.SurfaceRequirements{
		MinWidth:   400,
		MinHeight:  300,
		Capability: style.MonoSlow,
	},
	p: Params{BuildFn: buildMonoSlow400x300},
}

func buildMonoSlow400x300(snap content.ZMQData, pol content.Policy, ctx style.StyleContext, _ def) style.ViewData {
	return buildItemsOnly(snap, pol, ctx, itemsConfig{
		tier:       tiercatalog.TierNormal,
		paddingPct: 2,
	})
}
