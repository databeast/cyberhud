package styles

import (
	"github.com/databeast/cyberhud/display/modes/zmq/content"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
)

// GrayscaleSlow320x240Style renders the ZMQ message feed on a 320x240 grayscale-slow panel.
var GrayscaleSlow320x240Style = def{
	name: "grayscale-slow-320x240",
	reqs: style.SurfaceRequirements{
		MinWidth:   320,
		MinHeight:  240,
		Capability: style.GrayscaleSlow,
	},
	p: Params{BuildFn: buildGrayscaleSlow320x240},
}

func buildGrayscaleSlow320x240(snap content.ZMQData, pol content.Policy, ctx style.StyleContext, _ def) style.ViewData {
	return buildItemsOnly(snap, pol, ctx, itemsConfig{
		tier:       tiercatalog.TierNormal,
		paddingPct: 2,
	})
}
