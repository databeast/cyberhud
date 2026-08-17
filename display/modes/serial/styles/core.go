package styles

import (
	"github.com/databeast/cyberhud/display/modes/serial/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// Params holds the hand-tweakable knobs for a serial style.
type Params struct {
	BuildFn func(source.Snapshot, source.Policy, style.StyleContext, def) style.ViewData
}

// def is the declaration form of a serial style: name + surface requirements
// + layout tweaks. It implements style.Style[source.Snapshot, source.Policy].
type def struct {
	name string
	reqs style.SurfaceRequirements
	p    Params
}

var _ style.Style[source.Snapshot, source.Policy] = def{}

func (d def) Name() string { return d.name }

func (d def) Requirements() style.SurfaceRequirements { return d.reqs }

func (d def) Supports(hints textlayout.TextHints) style.Fitness {
	return style.EvaluateFitness(d.reqs, hints)
}

func (d def) Build(snapshot source.Snapshot, pol source.Policy, ctx style.StyleContext) style.ViewData {
	if d.p.BuildFn != nil {
		return d.p.BuildFn(snapshot, pol, ctx, d)
	}
	return style.ViewData{}
}
