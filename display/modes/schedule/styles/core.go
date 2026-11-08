// Package styles holds the hand-crafted presentation layer for the schedule
// display mode.
//
// Schedule has one visual intent: answer "what is on right now". The primary
// row is the active event's title, sized as large as the panel allows; a
// secondary row carries its time window; a third, optional row previews the
// next upcoming event. When nothing is active the view states that plainly and,
// when configured, counts down to the next event instead.
//
// There is exactly one look. It adapts to the panel through the font tier it
// picks for the title rather than splitting into separate named layouts, so
// every per-resolution style here is a thin declaration (name + surface
// requirements + a capability flag); the layout itself lives in adaptiveBuild.
// A style only places an already-settled snapshot — it reads no clock of its
// own; the active-event selection happens in the mode's render path against the
// current instant and arrives pre-resolved in the view model.
package styles

import (
	"github.com/databeast/cyberhud/display/modes/schedule/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// schedulePaddingPct is the inset the layout is solved against, matching the
// value carried back to the renderer in ViewData.PaddingPct.
const schedulePaddingPct = 6

// View is the resolved presentation model a style draws. The mode's render path
// builds it from the snapshot and the current instant so the style never does
// time arithmetic or event selection.
type View struct {
	// Title is the headline row: the active event's summary, a next-event label,
	// or a status line such as a load error.
	Title string
	// Subtitle is the active event's time window, or empty.
	Subtitle string
	// Upcoming is the optional next-event preview row, or empty.
	Upcoming string
	// Accent engages the accent palette on color panels (there is an active
	// event). Off while idle or erroring keeps the idle/error look restrained.
	Accent bool
}

// Params holds the per-style knobs.
type Params struct {
	// Color engages the accent palette on capable panels.
	Color bool
	// BuildFn, when set, replaces adaptiveBuild for this style.
	BuildFn func(source.Snapshot, source.Policy, style.StyleContext, def) style.ViewData
}

// def is the declaration form of a schedule style: name + surface requirements +
// Params. It implements style.Style[source.Snapshot, source.Policy].
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

func (d def) Build(snap source.Snapshot, pol source.Policy, ctx style.StyleContext) style.ViewData {
	if d.p.BuildFn != nil {
		return d.p.BuildFn(snap, pol, ctx, d)
	}
	return adaptiveBuild(snap, pol, ctx, d)
}
