package styles

import (
	"github.com/databeast/cyberhud/display/modes/zmq/content"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
)

// ColorSlow176x264Style renders the ZMQ message feed on a 176x264 color-slow panel.
var ColorSlow176x264Style = def{
	name: "color-slow-176x264",
	reqs: style.SurfaceRequirements{
		MinWidth:   176,
		MinHeight:  264,
		Capability: style.ColorSlow,
	},
	p: Params{BuildFn: buildColorSlow176x264},
}

func buildColorSlow176x264(snap content.ZMQData, pol content.Policy, ctx style.StyleContext, _ def) style.ViewData {
	return buildColorItemsOnly(snap, pol, ctx, itemsConfig{
		tier:       tiercatalog.TierNormal,
		paddingPct: 2,
	})
}
