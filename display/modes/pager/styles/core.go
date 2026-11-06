package styles

// Core style assembly for the pager mode.
//
// Every pager style is a `def` value — a name, the surface it targets, and a
// Params block of hand-tweak knobs. Its Build method places the resolved frame
// carried by the PagerSnapshot; the shared placement logic lives in
// placement.go, and the per-resolution def values are generated in generate.go
// from the constructors in constructors.go.
//
// A style only answers "where does this go on screen?". The snapshot it receives
// is already a settled frame (lines chosen, scroll offset or page phase/alpha
// resolved); nothing here reads a clock or advances animation. Which of the two
// render strategies applies is carried by snap.IsSlow, so a style never decides
// it and never mixes the two planes.

import (
	"github.com/databeast/cyberhud/display/modes/pager/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// snapshot and policy are the pager's concrete type arguments to the generic
// style machinery. Aliased here so the per-class generated slices in generate.go
// read without repeating the fully-qualified types.
type (
	snapshot = source.PagerSnapshot
	policy   = source.Policy
)

// Params holds the hand-tweakable knobs for a single pager style.
// The zero value produces the shared core placement.
type Params struct {
	// BuildFn, when set, replaces the core placement entirely for this style.
	// Attach it to tune a specific resolution/PPI without disturbing the shared
	// path the rest of the grid uses.
	BuildFn func(source.PagerSnapshot, source.Policy, style.StyleContext, def) style.ViewData
}

// def is the declaration form of a pager style: name + surface requirements +
// placement tweaks. It implements style.Style[source.PagerSnapshot, source.Policy].
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
	return coreBuild(snapshot, pol, ctx, d)
}

// coreBuild is the shared placement for every pager style. It commits the frame
// to a single plane based on snap.IsSlow: slow frames are placed on the sprite
// plane (page-fade), fast frames on the row plane (smooth scroll).
func coreBuild(snap source.PagerSnapshot, pol source.Policy, ctx style.StyleContext, d def) style.ViewData {
	if snap.IsSlow {
		return placePage(snap, pol, ctx, d)
	}
	return placeScroll(snap, pol, ctx, d)
}
