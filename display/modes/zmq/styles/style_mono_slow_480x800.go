package styles

import (
	"github.com/databeast/cyberhud/display/modes/zmq/content"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
)

// MonoSlow480x800Style renders the ZMQ message feed on a 480x800 mono-slow panel.
var MonoSlow480x800Style = def{
	name: "mono-slow-480x800",
	reqs: style.SurfaceRequirements{
		MinWidth:   480,
		MinHeight:  800,
		Capability: style.MonoSlow,
	},
	p: Params{BuildFn: buildMonoSlow480x800},
}

func buildMonoSlow480x800(snap content.ZMQData, pol content.Policy, ctx style.StyleContext, _ def) style.ViewData {
	return buildItemsOnly(snap, pol, ctx, itemsConfig{
		tier:       tiercatalog.TierLarge,
		paddingPct: 2,
	})
}
