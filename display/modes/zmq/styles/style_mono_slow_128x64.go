package styles

import (
	"github.com/databeast/cyberhud/display/modes/zmq/content"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
)

// MonoSlow128x64Style renders the ZMQ message feed on a 128x64 mono-slow panel.
var MonoSlow128x64Style = def{
	name: "mono-slow-128x64",
	reqs: style.SurfaceRequirements{
		MinWidth:   128,
		MinHeight:  64,
		Capability: style.MonoSlow,
	},
	p: Params{BuildFn: buildMonoSlow128x64},
}

func buildMonoSlow128x64(snap content.ZMQData, pol content.Policy, ctx style.StyleContext, _ def) style.ViewData {
	return buildItemsOnly(snap, pol, ctx, itemsConfig{
		tier:       tiercatalog.TierSmall,
		paddingPct: 2,
	})
}
