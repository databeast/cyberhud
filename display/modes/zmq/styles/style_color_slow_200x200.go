package styles

import (
	"github.com/databeast/cyberhud/display/modes/zmq/content"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
)

// ColorSlow200x200Style renders the ZMQ message feed on a 200x200 color-slow panel.
var ColorSlow200x200Style = def{
	name: "color-slow-200x200",
	reqs: style.SurfaceRequirements{
		MinWidth:   200,
		MinHeight:  200,
		Capability: style.ColorSlow,
	},
	p: Params{BuildFn: buildColorSlow200x200},
}

func buildColorSlow200x200(snap content.ZMQData, pol content.Policy, ctx style.StyleContext, _ def) style.ViewData {
	return buildColorItemsOnly(snap, pol, ctx, itemsConfig{
		tier:       tiercatalog.TierNormal,
		paddingPct: 2,
	})
}
