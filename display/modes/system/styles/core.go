package styles

import (
	"github.com/databeast/cyberhud/display/modes/system/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

type Params struct {
	BuildFn func(source.SystemSnapshot, source.Policy, style.StyleContext, def) style.ViewData
}

type def struct {
	name string
	reqs style.SurfaceRequirements
	p    Params
}

var _ style.Style[source.SystemSnapshot, source.Policy] = def{}

func (d def) Name() string { return d.name }

func (d def) Requirements() style.SurfaceRequirements { return d.reqs }

func (d def) Supports(hints textlayout.TextHints) style.Fitness {
	return style.EvaluateFitness(d.reqs, hints)
}

func (d def) Build(data source.SystemSnapshot, pol source.Policy, ctx style.StyleContext) style.ViewData {
	if d.p.BuildFn != nil {
		return d.p.BuildFn(data, pol, ctx, d)
	}
	return adaptiveBuild(data, pol, ctx, d)
}

// adaptiveBuild is the fallback layout for skeleton styles (nil BuildFn),
// selected by fitness when the user has not picked a specific variant. It
// chooses a sensible content variant based on panel dimensions:
//   - Tiny panels (width < 64 or height < 64): the compact one-line name tag.
//   - Everything else: the default identity card.
//
// The user reaches the data-dense cores/top variants by selecting them
// explicitly; the skeleton defaults to an always-sensible readout.
func adaptiveBuild(data source.SystemSnapshot, pol source.Policy, ctx style.StyleContext, d def) style.ViewData {
	hints := ctx.Hints()
	if hints.PixelWidth < 64 || hints.PixelHeight < 64 {
		return buildCompactStyle(data, pol, ctx, d)
	}
	return buildDefaultStyle(data, pol, ctx, d)
}
