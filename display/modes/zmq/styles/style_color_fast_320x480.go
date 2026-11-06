package styles

import (
	"github.com/databeast/cyberhud/display/modes/zmq/content"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
)

// ColorFast320x480Style renders the ZMQ message feed on a 320x480 color-fast panel.
var ColorFast320x480Style = def{
	name: "color-fast-320x480",
	reqs: style.SurfaceRequirements{
		MinWidth:   320,
		MinHeight:  480,
		Capability: style.ColorFast,
	},
	p: Params{BuildFn: buildColorFast320x480},
}

func buildColorFast320x480(snap content.ZMQData, pol content.Policy, ctx style.StyleContext, _ def) style.ViewData {
	return buildColorItemsOnly(snap, pol, ctx, itemsConfig{
		tier:       tiercatalog.TierLarge,
		paddingPct: 0,
	})
}
