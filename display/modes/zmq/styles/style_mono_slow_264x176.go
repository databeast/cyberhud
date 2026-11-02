package styles

import (
	"github.com/databeast/cyberhud/display/modes/zmq/content"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
)

// MonoSlow264x176Style renders the ZMQ message feed on a 264x176 mono-slow panel.
var MonoSlow264x176Style = def{
	name: "mono-slow-264x176",
	reqs: style.SurfaceRequirements{
		MinWidth:   264,
		MinHeight:  176,
		Capability: style.MonoSlow,
	},
	p: Params{BuildFn: buildMonoSlow264x176},
}

func buildMonoSlow264x176(snap content.ZMQData, pol content.Policy, ctx style.StyleContext, _ def) style.ViewData {
	return buildItemsOnly(snap, pol, ctx, itemsConfig{
		tier:       tiercatalog.TierNormal,
		paddingPct: 2,
	})
}
