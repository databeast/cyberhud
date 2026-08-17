package styles

// Core layout assembly for pager styles.
//
// Every concrete pager style in this package is declared as a `def` value: a
// name, surface requirements, and a Params block of hand-tweakable knobs. The
// def's Build method assembles the shared pager layout, adjusted by those
// Params. Set Params.BuildFn for a fully bespoke layout.

import (
	"github.com/databeast/cyberhud/display/modes/pager/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// Params holds the hand-tweakable knobs for a single pager style.
// The zero value produces the static slow-refresh pager layout.
type Params struct {
	// Fast selects the fast-refresh pager layout: use the snapshot's pixel
	// OffsetY and keep the output non-static for continuous scroll refresh.
	Fast bool

	// BuildFn, when set, replaces the core layout entirely for this style.
	BuildFn func(source.PagerSnapshot, source.Policy, style.StyleContext, def) style.ViewData
}

// def is the declaration form of a pager style: name + surface requirements +
// layout tweaks. It implements style.Style[source.PagerSnapshot, source.Policy].
type def struct {
	name string
	reqs style.SurfaceRequirements
	p    Params
}

var _ style.Style[source.PagerSnapshot, source.Policy] = def{}

func (d def) Name() string { return d.name }

func (d def) Requirements() style.SurfaceRequirements { return d.reqs }

func (d def) Supports(hints textlayout.TextHints) style.Fitness {
	return style.EvaluateFitness(d.reqs, hints)
}

func (d def) Build(snapshot source.PagerSnapshot, pol source.Policy, ctx style.StyleContext) style.ViewData {
	if d.p.BuildFn != nil {
		return d.p.BuildFn(snapshot, pol, ctx, d)
	}
	if d.p.Fast {
		return style.ViewData{
			Items:   snapshot.Lines,
			OffsetY: snapshot.OffsetY,
			Static:  false,
		}
	}
	return style.ViewData{
		Items:  snapshot.Lines,
		Static: true,
	}
}
