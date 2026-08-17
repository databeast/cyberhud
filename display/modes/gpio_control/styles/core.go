package styles

import (
	"github.com/databeast/cyberhud/display/modes/gpio_control/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

type Params struct {
	BuildFn func(data source.Data, pol source.Policy, ctx style.StyleContext, d def) style.ViewData
}

type def struct {
	name string
	reqs style.SurfaceRequirements
	p    Params
}

var _ style.Style[source.Data, source.Policy] = def{}

func (d def) Name() string { return d.name }

func (d def) Requirements() style.SurfaceRequirements { return d.reqs }

func (d def) Supports(hints textlayout.TextHints) style.Fitness {
	return style.EvaluateFitness(d.reqs, hints)
}

func (d def) Build(data source.Data, pol source.Policy, ctx style.StyleContext) style.ViewData {
	if d.p.BuildFn != nil {
		return d.p.BuildFn(data, pol, ctx, d)
	}
	// Shared layout dispatch based on capability tier and dimensions.
	if d.reqs.Capability >= style.ColorSlow && d.reqs.MinWidth >= 64 && d.reqs.MinHeight >= 64 {
		return sharedGridBuild(data, pol, ctx, d)
	}
	return sharedListBuild(data, pol, ctx, d)
}
