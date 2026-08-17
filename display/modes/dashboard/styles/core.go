package styles

// Core layout assembly for dashboard styles.
//
// Every per-resolution style in this package is declared as a `def` value in
// its styles_WxH.go file: a name, the surface requirements it targets, and a
// Params block of hand-tweakable layout selection. Edit declarations there for
// display-specific tuning, or set Params.BuildFn for a fully bespoke layout.

import (
	"github.com/databeast/cyberhud/display/modes/dashboard/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// Layout selects the shared dashboard builder family for a declaration.
type Layout uint8

const (
	LayoutGrayscaleDashboard Layout = iota
	LayoutColorSkeleton
	LayoutColorFastSkeleton
	LayoutCompactLandscape
	LayoutCompactPortrait
)

// Params holds the hand-tweakable knobs for a single per-resolution style.
// The zero value selects the mono skeleton layout.
type Params struct {
	Layout  Layout
	BuildFn func(data source.DashboardContent, pol source.Policy, ctx style.StyleContext, d styleDef) style.ViewData
}

// styleDef is the declaration form of a dashboard style: name + surface requirements
// + layout tweaks. It implements style.Style[source.DashboardContent, source.Policy].
type styleDef struct {
	name string
	reqs style.SurfaceRequirements
	p    Params
}

var _ style.Style[source.DashboardContent, source.Policy] = styleDef{}

func (d styleDef) Name() string { return d.name }

func (d styleDef) Requirements() style.SurfaceRequirements { return d.reqs }

func (d styleDef) Supports(hints textlayout.TextHints) style.Fitness {
	return style.EvaluateFitness(d.reqs, hints)
}

func (d styleDef) Build(data source.DashboardContent, pol source.Policy, ctx style.StyleContext) style.ViewData {
	if d.p.BuildFn != nil {
		return d.p.BuildFn(data, pol, ctx, d)
	}
	switch d.p.Layout {
	case LayoutGrayscaleDashboard:
		return buildGrayscaleDashboard(data, ctx)
	case LayoutColorSkeleton:
		return buildColorSkeleton(data, ctx)
	case LayoutColorFastSkeleton:
		return buildColorFastSkeleton(data, ctx)
	case LayoutCompactLandscape:
		return compactLandscapeLayout(data, pol, ctx)
	case LayoutCompactPortrait:
		return compactPorttraitLayout(data, pol, ctx)
	default:
		return buildMonoSkeleton(data, ctx)
	}
}
