package styles

import (
	"github.com/databeast/cyberhud/display/modes/zmq/content"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
)

// MonoFast160x80Style renders the ZMQ message feed on a 160x80 mono-fast panel.
var MonoFast160x80Style = def{
	name: "mono-fast-160x80",
	reqs: style.SurfaceRequirements{
		MinWidth:   160,
		MinHeight:  80,
		Capability: style.MonoFast,
	},
	p: Params{BuildFn: buildMonoFast160x80},
}

func buildMonoFast160x80(snap content.ZMQData, pol content.Policy, ctx style.StyleContext, _ def) style.ViewData {
	return buildItemsOnly(snap, pol, ctx, itemsConfig{
		tier:       tiercatalog.TierSmall,
		paddingPct: 0,
	})
}
