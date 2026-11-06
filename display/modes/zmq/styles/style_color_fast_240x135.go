package styles

import (
	"github.com/databeast/cyberhud/display/modes/zmq/content"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
)

// ColorFast240x135Style renders the ZMQ message feed on a 240x135 color-fast panel.
var ColorFast240x135Style = def{
	name: "color-fast-240x135",
	reqs: style.SurfaceRequirements{
		MinWidth:   240,
		MinHeight:  135,
		Capability: style.ColorFast,
	},
	p: Params{BuildFn: buildColorFast240x135},
}

func buildColorFast240x135(snap content.ZMQData, pol content.Policy, ctx style.StyleContext, _ def) style.ViewData {
	return buildColorItemsOnly(snap, pol, ctx, itemsConfig{
		tier:       tiercatalog.TierNormal,
		paddingPct: 0,
	})
}
