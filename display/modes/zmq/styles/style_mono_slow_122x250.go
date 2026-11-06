package styles

import (
	"github.com/databeast/cyberhud/display/modes/zmq/content"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
)

// MonoSlow122x250Style renders the ZMQ message feed on a 122x250 mono-slow panel.
var MonoSlow122x250Style = def{
	name: "mono-slow-122x250",
	reqs: style.SurfaceRequirements{
		MinWidth:   122,
		MinHeight:  250,
		Capability: style.MonoSlow,
	},
	p: Params{BuildFn: buildMonoSlow122x250},
}

func buildMonoSlow122x250(snap content.ZMQData, pol content.Policy, ctx style.StyleContext, _ def) style.ViewData {
	return buildItemsOnly(snap, pol, ctx, itemsConfig{
		tier:       tiercatalog.TierNormal,
		paddingPct: 2,
	})
}
