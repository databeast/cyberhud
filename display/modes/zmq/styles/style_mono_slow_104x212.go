package styles

import (
	"github.com/databeast/cyberhud/display/modes/zmq/content"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
)

// MonoSlow104x212Style renders the ZMQ message feed on a 104x212 mono-slow panel.
var MonoSlow104x212Style = def{
	name: "mono-slow-104x212",
	reqs: style.SurfaceRequirements{
		MinWidth:   104,
		MinHeight:  212,
		Capability: style.MonoSlow,
	},
	p: Params{BuildFn: buildMonoSlow104x212},
}

func buildMonoSlow104x212(snap content.ZMQData, pol content.Policy, ctx style.StyleContext, _ def) style.ViewData {
	return buildItemsOnly(snap, pol, ctx, itemsConfig{
		tier:       tiercatalog.TierNormal,
		paddingPct: 2,
	})
}
