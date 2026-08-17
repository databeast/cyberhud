package styles

import (
	"github.com/databeast/cyberhud/display/modes/serial/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
)

var CompactStyle = def{
	name: "compact",
	reqs: style.SurfaceRequirements{
		MinRows:         1,
		MinCharsPerLine: 10,
	},
	p: Params{BuildFn: buildCompactStyle},
}

func buildCompactStyle(snapshot source.Snapshot, pol source.Policy, ctx style.StyleContext, _ def) style.ViewData {
	p := pol

	// Get region constraints from the style context.
	hints := ctx.Hints()

	// Resolve face via tier catalog for sprite widget rendering.
	face := ctx.Face("spleen", tiercatalog.TierNormal)
	hints = textlayout.WithFont(hints, face)

	vd := buildCompact(snapshot, p, hints, face)
	return vd
}
