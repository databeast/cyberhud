package styles

import (
	"github.com/databeast/cyberhud/display/modes/zmq/content"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
)

// MonoSlow212x104Style renders the ZMQ message feed on a 212x104 mono-slow panel.
var MonoSlow212x104Style = def{
	name: "mono-slow-212x104",
	reqs: style.SurfaceRequirements{
		MinWidth:   212,
		MinHeight:  104,
		Capability: style.MonoSlow,
	},
	p: Params{BuildFn: buildMonoSlow212x104},
}

func buildMonoSlow212x104(snap content.ZMQData, pol content.Policy, ctx style.StyleContext, _ def) style.ViewData {
	return buildItemsOnly(snap, pol, ctx, itemsConfig{
		tier:       tiercatalog.TierNormal,
		paddingPct: 2,
	})
}
