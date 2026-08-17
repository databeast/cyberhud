package styles

import (
	"github.com/databeast/cyberhud/display/modes/menu/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

type Params struct {
	BuildFn func(source.MenuSnapshot, source.Policy, style.StyleContext, def) style.ViewData
}

type def struct {
	name string
	reqs style.SurfaceRequirements
	p    Params
}

var _ style.Style[source.MenuSnapshot, source.Policy] = def{}

func (d def) Name() string { return d.name }

func (d def) Requirements() style.SurfaceRequirements { return d.reqs }

func (d def) Supports(hints textlayout.TextHints) style.Fitness {
	return style.EvaluateFitness(d.reqs, hints)
}

func (d def) Build(data source.MenuSnapshot, pol source.Policy, ctx style.StyleContext) style.ViewData {
	if d.p.BuildFn != nil {
		return d.p.BuildFn(data, pol, ctx, d)
	}
	return style.ViewData{}
}
