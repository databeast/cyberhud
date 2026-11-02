package styles

import (
	"github.com/databeast/cyberhud/display/modes/zmq/content"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
)

// GrayscaleFast240x135Style renders the ZMQ message feed on a 240x135 grayscale-fast panel.
var GrayscaleFast240x135Style = def{
	name: "grayscale-fast-240x135",
	reqs: style.SurfaceRequirements{
		MinWidth:   240,
		MinHeight:  135,
		Capability: style.GrayscaleFast,
	},
	p: Params{BuildFn: buildGrayscaleFast240x135},
}

func buildGrayscaleFast240x135(snap content.ZMQData, pol content.Policy, ctx style.StyleContext, _ def) style.ViewData {
	return buildItemsOnly(snap, pol, ctx, itemsConfig{
		tier:       tiercatalog.TierNormal,
		paddingPct: 0,
	})
}
