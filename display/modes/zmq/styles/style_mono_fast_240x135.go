package styles

import (
	"github.com/databeast/cyberhud/display/modes/zmq/content"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
)

// MonoFast240x135Style renders the ZMQ message feed on a 240x135 mono-fast panel.
var MonoFast240x135Style = def{
	name: "mono-fast-240x135",
	reqs: style.SurfaceRequirements{
		MinWidth:   240,
		MinHeight:  135,
		Capability: style.MonoFast,
	},
	p: Params{BuildFn: buildMonoFast240x135},
}

func buildMonoFast240x135(snap content.ZMQData, pol content.Policy, ctx style.StyleContext, _ def) style.ViewData {
	return buildItemsOnly(snap, pol, ctx, itemsConfig{
		tier:       tiercatalog.TierNormal,
		paddingPct: 0,
	})
}
