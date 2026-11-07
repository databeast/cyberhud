package styles

// Core dispatch for GPIO styles.
//
// Every per-resolution style in this package is declared as a def value in its
// styles_WxH.go file: a name and surface requirements. All of them render the
// single live-functions view (see buildListInternal in styles.go) via
// adaptiveBuild. Params.BuildFn remains as an escape hatch for a future bespoke
// layout, but nothing currently sets it.

import (
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// Params holds the optional per-style override. The common case is the zero
// value (no override), which routes every style through adaptiveBuild — the one
// live-functions view this mode delivers.
type Params struct {
	BuildFn func(data GpioSnapshot, pol Policy, ctx style.StyleContext, d def) style.ViewData
}

type def struct {
	name string
	reqs style.SurfaceRequirements
	p    Params
}

var _ style.Style[GpioSnapshot, Policy] = def{}

func (d def) Name() string { return d.name }

func (d def) Requirements() style.SurfaceRequirements { return d.reqs }

func (d def) Supports(hints textlayout.TextHints) style.Fitness {
	return style.EvaluateFitness(d.reqs, hints)
}

func (d def) Build(data GpioSnapshot, pol Policy, ctx style.StyleContext) style.ViewData {
	if d.p.BuildFn != nil {
		return d.p.BuildFn(data, pol, ctx, d)
	}
	// Every style renders the single live-functions view. The mode deliberately
	// no longer forks into list/icons/detail/dashboard/activity layouts — those
	// were the old sprawl pulling in three directions. The legacy style names
	// survive as historical aliases (a user who set style=detail still gets the
	// good view), but all resolve to the one idea: "what is the header doing,
	// right now?"
	return adaptiveBuild(data, pol, ctx, d)
}

// adaptiveBuild renders the live-functions view at every panel size, so there
// is no "not implemented" fall-through and no aggregate-counter degradation (a
// tally like "Pins:N Hi:K" throws away pin identity, the one thing that
// matters). buildListInternal distributes its rows across the available height
// and truncates text to the width, so a tiny strip shows fewer, tighter rows of
// the SAME per-function unit rather than a different representation. Color vs
// mono is handled inside the builder.
func adaptiveBuild(data GpioSnapshot, pol Policy, ctx style.StyleContext, d def) style.ViewData {
	hints := ctx.Hints()
	if hints.PixelWidth == 0 || hints.PixelHeight == 0 {
		return style.ViewData{}
	}
	// On a roomy panel, draw the physical connector as a live map: position
	// carries identity (and "where"), color/shape carry state — the vivid view
	// that earns a color panel. Where a legible 20×2 pad grid will not fit, fall
	// back to the compact live-functions text list.
	bridge := ctx.Layout(0)
	if boardMapFits(bridge.AvailableContentWidth(), bridge.AvailableContentHeight()) {
		return buildBoardMap(data, ctx)
	}
	return buildListInternal(data, ctx)
}
