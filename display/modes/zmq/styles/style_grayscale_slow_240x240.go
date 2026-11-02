package styles

import (
	"github.com/databeast/cyberhud/display/modes/zmq/content"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
)

// GrayscaleSlow240x240Style renders the ZMQ message feed on a 240x240 grayscale-slow panel.
var GrayscaleSlow240x240Style = def{
	name: "grayscale-slow-240x240",
	reqs: style.SurfaceRequirements{
		MinWidth:   240,
		MinHeight:  240,
		Capability: style.GrayscaleSlow,
	},
	p: Params{BuildFn: buildGrayscaleSlow240x240},
}

func buildGrayscaleSlow240x240(snap content.ZMQData, pol content.Policy, ctx style.StyleContext, _ def) style.ViewData {
	return buildItemsOnly(snap, pol, ctx, itemsConfig{
		tier:       tiercatalog.TierNormal,
		paddingPct: 2,
	})
}
