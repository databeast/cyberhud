package styles

// Core style dispatch for the serial monitor.
//
// Serial has one visual intent: the console instrument panel — a connection
// LED, a port/baud header, a connection-state line, and the ANSI-colored output
// lines. There is exactly one look; it adapts to the panel rather than splitting
// into separate named styles. Every per-resolution style is a thin `def`
// declaration (name + surface requirements + a Params flag) generated in
// generate.go from the constructors in constructors.go; the look itself lives in
// the shared engine in engine.go.
//
// The Params flag selects which engine variant renders:
//   - Fast  → the live instrument panel (animated LED, throughput sparkline) for
//     OLED/TFT/LED-matrix panels.
//   - zero  → the e-paper-safe panel (static, no animated widgets) for slow
//     refresh displays.
//   - Color → the Fast variant with the accent palette engaged.
//
// A style only places an already-settled snapshot; it reads no clock and holds
// no per-frame state. The raw-vs-chrome choice is a policy toggle (Policy.Raw),
// not a style, so every panel size honours it uniformly.

import (
	"github.com/databeast/cyberhud/display/modes/serial/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// Params holds the hand-tweakable knobs for a single serial style.
type Params struct {
	// Fast selects the live instrument panel (animated LED + throughput
	// sparkline). Zero value selects the e-paper-safe static panel.
	Fast bool

	// Color engages the accent palette on the Fast panel.
	Color bool

	// BuildFn, when set, replaces the engine entirely for this style.
	BuildFn func(source.Snapshot, source.Policy, style.StyleContext, def) style.ViewData
}

// def is the declaration form of a serial style: name + surface requirements +
// engine selection. It implements style.Style[source.Snapshot, source.Policy].
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
	if d.p.Fast {
		return buildFastPanel(snapshot, pol, ctx, d)
	}
	return buildSlowPanel(snapshot, pol, ctx, d)
}
