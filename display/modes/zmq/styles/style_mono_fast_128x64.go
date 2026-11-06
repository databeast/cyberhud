package styles

import (
	"github.com/databeast/cyberhud/display/modes/zmq/content"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
)

// MonoFast128x64Style renders the ZMQ message feed on a 128x64 mono-fast panel.
var MonoFast128x64Style = def{
	name: "mono-fast-128x64",
	reqs: style.SurfaceRequirements{
		MinWidth:   128,
		MinHeight:  64,
		Capability: style.MonoFast,
	},
	p: Params{BuildFn: buildMonoFast128x64},
}

func buildMonoFast128x64(snap content.ZMQData, pol content.Policy, ctx style.StyleContext, _ def) style.ViewData {
	return buildItemsOnly(snap, pol, ctx, itemsConfig{
		tier:       tiercatalog.TierSmall,
		paddingPct: 0,
	})
}
