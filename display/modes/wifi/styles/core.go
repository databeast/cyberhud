package styles

// Core layout dispatch for WiFi styles.
//
// Every per-resolution style in this package is declared as a `def` value in
// its styles_WxH.go file: a name, the surface requirements it targets, and a
// Params block of hand-tweakable knobs. The def's Build method dispatches to
// the shared WiFi layouts, adjusted by those Params.
//
// To hand-tweak WiFi for a specific display, edit that resolution's Params -
// or set Params.BuildFn for a fully bespoke layout.

import (
	"github.com/databeast/cyberhud/display/modes/wifi/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// Params holds the hand-tweakable knobs for a single per-resolution WiFi style.
// The zero value produces the shared placeholder layout used by unimplemented
// legacy skeleton styles.
type Params struct {
	// MonoSlow selects the shared polished mono slow-refresh WiFi layout.
	// Retained for backward compatibility with existing styles that set it.
	MonoSlow bool

	// Fast selects the shared fast-refresh layout engine (icons, bars, progress).
	Fast bool

	// Color enables accent color rendering in the Fast layout engine.
	Color bool

	// BuildFn, when set, replaces the core layout entirely for this style.
	BuildFn func(data source.WifiData, pol source.Policy, ctx style.StyleContext, d def) style.ViewData
}

// def is the declaration form of a WiFi style: name + surface requirements +
// layout tweaks. It implements style.Style[source.WifiData, source.Policy].
type def struct {
	name string
	reqs style.SurfaceRequirements
	p    Params
}

var _ style.Style[source.WifiData, source.Policy] = def{}

func (d def) Name() string { return d.name }

func (d def) Requirements() style.SurfaceRequirements { return d.reqs }

func (d def) Supports(hints textlayout.TextHints) style.Fitness {
	return style.EvaluateFitness(d.reqs, hints)
}

func (d def) Build(data source.WifiData, pol source.Policy, ctx style.StyleContext) style.ViewData {
	if d.p.BuildFn != nil {
		return d.p.BuildFn(data, pol, ctx, d)
	}
	if d.p.Fast {
		return buildFastWifi(data, pol, ctx, d)
	}
	// MonoSlow flag or zero-value Params both dispatch to Slow engine.
	return buildSlowWifi(data, ctx, d.reqs)
}
