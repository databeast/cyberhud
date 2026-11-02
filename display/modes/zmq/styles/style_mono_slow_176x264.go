package styles

import (
	"github.com/databeast/cyberhud/display/modes/zmq/content"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
)

// MonoSlow176x264Style renders the ZMQ message feed on a 176x264 mono-slow panel.
var MonoSlow176x264Style = def{
	name: "mono-slow-176x264",
	reqs: style.SurfaceRequirements{
		MinWidth:   176,
		MinHeight:  264,
		Capability: style.MonoSlow,
	},
	p: Params{BuildFn: buildMonoSlow176x264},
}

func buildMonoSlow176x264(snap content.ZMQData, pol content.Policy, ctx style.StyleContext, _ def) style.ViewData {
	return buildItemsOnly(snap, pol, ctx, itemsConfig{
		tier:       tiercatalog.TierNormal,
		paddingPct: 2,
	})
}
