package styles

import (
	"github.com/databeast/cyberhud/display/modes/serial/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
)

var DashboardStyle = def{
	name: "dashboard",
	reqs: style.SurfaceRequirements{
		MinWidth:  64,
		MinHeight: 64,
	},
	p: Params{BuildFn: buildDashboardStyle},
}

func buildDashboardStyle(snapshot source.Snapshot, pol source.Policy, ctx style.StyleContext, _ def) style.ViewData {
	p := pol

	// Get region constraints from the style context.
	hints := ctx.Hints()

	// Resolve face via tier catalog for sprite widget rendering.
	face := ctx.Face("spleen", tiercatalog.TierNormal)
	hints = textlayout.WithFont(hints, face)

	vd := buildDashboard(snapshot, p, hints, face)
	return vd
}
