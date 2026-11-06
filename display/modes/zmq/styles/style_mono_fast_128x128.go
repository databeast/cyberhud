package styles

import (
	"github.com/databeast/cyberhud/display/modes/zmq/content"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
)

// MonoFast128x128Style renders the ZMQ message feed on a 128x128 mono-fast panel.
var MonoFast128x128Style = def{
	name: "mono-fast-128x128",
	reqs: style.SurfaceRequirements{
		MinWidth:   128,
		MinHeight:  128,
		Capability: style.MonoFast,
	},
	p: Params{BuildFn: buildMonoFast128x128},
}

func buildMonoFast128x128(snap content.ZMQData, pol content.Policy, ctx style.StyleContext, _ def) style.ViewData {
	return buildItemsOnly(snap, pol, ctx, itemsConfig{
		tier:       tiercatalog.TierNormal,
		paddingPct: 0,
	})
}
