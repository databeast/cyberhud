package styles

import (
	"github.com/databeast/cyberhud/display/modes/serial/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
)

var RawStyle = def{
	name: "raw",
	reqs: style.SurfaceRequirements{
		MinRows:         1,
		MinCharsPerLine: 8,
	},
	p: Params{BuildFn: buildRawStyle},
}

func buildRawStyle(snapshot source.Snapshot, pol source.Policy, ctx style.StyleContext, _ def) style.ViewData {
	p := pol

	// Get region constraints from the style context.
	hints := ctx.Hints()

	// Resolve face via tier catalog for sprite widget rendering.
	// Raw style doesn't produce sprites, but we still update hints metrics.
	face := ctx.Face("spleen", tiercatalog.TierNormal)
	hints = textlayout.WithFont(hints, face)

	vd := buildRaw(snapshot, p, hints, face)
	return vd
}
