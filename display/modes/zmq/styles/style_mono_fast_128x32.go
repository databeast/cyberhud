package styles

import (
	"github.com/databeast/cyberhud/display/modes/zmq/content"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
)

// MonoFast128x32Style renders the ZMQ message feed on a 128x32 mono-fast panel.
var MonoFast128x32Style = def{
	name: "mono-fast-128x32",
	reqs: style.SurfaceRequirements{
		MinWidth:   128,
		MinHeight:  32,
		Capability: style.MonoFast,
	},
	p: Params{BuildFn: buildMonoFast128x32},
}

func buildMonoFast128x32(snap content.ZMQData, pol content.Policy, ctx style.StyleContext, _ def) style.ViewData {
	return buildItemsOnly(snap, pol, ctx, itemsConfig{
		tier:       tiercatalog.TierSmall,
		paddingPct: 0,
	})
}
