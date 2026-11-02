package styles

import (
	"github.com/databeast/cyberhud/display/modes/zmq/content"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
)

// GrayscaleSlow480x320Style renders the ZMQ message feed on a 480x320 grayscale-slow panel.
var GrayscaleSlow480x320Style = def{
	name: "grayscale-slow-480x320",
	reqs: style.SurfaceRequirements{
		MinWidth:   480,
		MinHeight:  320,
		Capability: style.GrayscaleSlow,
	},
	p: Params{BuildFn: buildGrayscaleSlow480x320},
}

func buildGrayscaleSlow480x320(snap content.ZMQData, pol content.Policy, ctx style.StyleContext, _ def) style.ViewData {
	return buildItemsOnly(snap, pol, ctx, itemsConfig{
		tier:       tiercatalog.TierNormal,
		paddingPct: 2,
	})
}
