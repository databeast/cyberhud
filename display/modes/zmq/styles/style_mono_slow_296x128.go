package styles

import (
	"github.com/databeast/cyberhud/display/modes/zmq/content"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
)

// MonoSlow296x128Style renders the ZMQ message feed on a 296x128 mono-slow panel.
var MonoSlow296x128Style = def{
	name: "mono-slow-296x128",
	reqs: style.SurfaceRequirements{
		MinWidth:   296,
		MinHeight:  128,
		Capability: style.MonoSlow,
	},
	p: Params{BuildFn: buildMonoSlow296x128},
}

func buildMonoSlow296x128(snap content.ZMQData, pol content.Policy, ctx style.StyleContext, _ def) style.ViewData {
	return buildItemsOnly(snap, pol, ctx, itemsConfig{
		tier:       tiercatalog.TierNormal,
		paddingPct: 2,
	})
}
