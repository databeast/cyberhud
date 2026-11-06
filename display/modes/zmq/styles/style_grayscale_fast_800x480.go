package styles

import (
	"github.com/databeast/cyberhud/display/modes/zmq/content"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
)

// GrayscaleFast800x480Style renders the ZMQ message feed on a 800x480 grayscale-fast panel.
var GrayscaleFast800x480Style = def{
	name: "grayscale-fast-800x480",
	reqs: style.SurfaceRequirements{
		MinWidth:   800,
		MinHeight:  480,
		Capability: style.GrayscaleFast,
	},
	p: Params{BuildFn: buildGrayscaleFast800x480},
}

func buildGrayscaleFast800x480(snap content.ZMQData, pol content.Policy, ctx style.StyleContext, _ def) style.ViewData {
	return buildItemsOnly(snap, pol, ctx, itemsConfig{
		tier:       tiercatalog.TierLarge,
		paddingPct: 0,
	})
}
