package styles

import (
	"github.com/databeast/cyberhud/display/modes/zmq/content"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
)

// MonoSlow200x200Style renders the ZMQ message feed on a 200x200 mono-slow panel.
var MonoSlow200x200Style = def{
	name: "mono-slow-200x200",
	reqs: style.SurfaceRequirements{
		MinWidth:   200,
		MinHeight:  200,
		Capability: style.MonoSlow,
	},
	p: Params{BuildFn: buildMonoSlow200x200},
}

func buildMonoSlow200x200(snap content.ZMQData, pol content.Policy, ctx style.StyleContext, _ def) style.ViewData {
	return buildItemsOnly(snap, pol, ctx, itemsConfig{
		tier:       tiercatalog.TierNormal,
		paddingPct: 2,
	})
}
