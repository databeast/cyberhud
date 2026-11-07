package pager

import (
	"time"

	"github.com/databeast/cyberhud/display/modes/pager/source"
	"github.com/databeast/cyberhud/display/region"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// buildView produces the pager's frame for one render pass.
//
// This is where all pager logic lives. It reads policy, classifies the panel,
// advances the scroll or page animation against the frame clock, snapshots the
// buffer, and resolves the visible window — then bakes that settled result into
// a source.PagerSnapshot and dispatches through the registry exactly like clock
// and thermal. The resolved style only places the frame; it sees no clock, no
// animation state, no buffer.
//
// elapsed is the time since the previous pass, capped for smooth scrolling.
// It is threaded in from the instance, which holds the per-frame clock, so this
// function stays a pure translation of (policy, hints, elapsed) into a frame.
func buildView(hints textlayout.TextHints, elapsed time.Duration, typer *source.TyperState) style.ViewData {
	pol := GetPolicy()

	// No source configured → neutral status, no style dispatch.
	if pol.Source == "" {
		return buildNoSourceView()
	}

	ctx := style.NewStyleContext(hints)
	layout := source.ComputeLayout(ctx)
	sc := source.ClassifySurface(hints)
	if !sc.IsSlow() && typer != nil && layout.GlyphAdvance > 0 && layout.RowHeight > 0 {
		bridge := ctx.Layout(source.PagerPaddingPct)
		layout.VisibleColumns = bridge.AvailableContentWidth() / layout.GlyphAdvance
		layout.VisibleRows = bridge.AvailableContentHeight() / layout.RowHeight
	}

	// Resolve this pass's frame: advance the active animation and select the
	// visible lines. The snapshot that comes back is fully settled.
	snap := resolveFrame(sc, layout, pol, elapsed, typer)

	// Registry-based dispatch: configured → alias → fitness. The style places
	// the frame; which plane it commits to follows from snap.IsSlow.
	s, reason := style.ResolveStyle(pagerRegistry, hints, "pager", pol.Style)
	vd := s.Build(snap, pol, ctx)
	vd.StyleReport = style.ReportFor(s, reason, hints)
	return vd
}

// resolveFrame advances the active scroll or page animation by elapsed and
// returns the settled frame for this pass. The scroll path adapts velocity to
// the backlog then advances the pixel offset; the page path ticks the fade state
// machine. Both select and truncate the visible lines against the layout so the
// style receives a ready-to-place window.
func resolveFrame(sc source.SurfaceClass, layout source.Layout, pol source.Policy, elapsed time.Duration, typer *source.TyperState) source.PagerSnapshot {
	base := source.PagerSnapshot{
		Policy:         pol,
		IsSlow:         sc.IsSlow(),
		RowHeight:      layout.RowHeight,
		VisibleColumns: layout.VisibleColumns,
	}

	if layout.VisibleRows <= 0 || layout.VisibleColumns <= 0 {
		return base
	}

	_, buf, scroll, page := source.ActiveStateSnapshot()
	if buf == nil {
		return base
	}
	if sc.IsSlow() {
		lines, firstLine := buf.SnapshotWithPosition()
		if page == nil {
			return base
		}
		if elapsed > 0 {
			page.Tick(elapsed, lines, firstLine, layout.VisibleRows, pol.FadeOutMS, pol.FadeInMS, pol.MaxWaitS)
		}
		base.Lines = source.PageFrame(page, layout)
		base.Phase = page.Phase()
		base.Alpha = page.FadeAlpha()
		return base
	}

	if typer != nil {
		lines, breaks, partial, firstLine := buf.SnapshotForTyper()
		typer.Advance(elapsed, pol.TypeSpeed, lines, breaks, partial, firstLine, layout.VisibleRows, layout.VisibleColumns)
		base.IsTyper = true
		base.Lines = typer.Frame(layout.VisibleRows)
		return base
	}

	if scroll == nil {
		return base
	}
	lines := buf.Snapshot()
	scroll.AdaptVelocity(len(lines), layout.VisibleRows)
	if elapsed > 0 {
		scroll.Advance(elapsed.Seconds(), layout.RowHeight)
	}
	visible, offsetY := source.ScrollFrame(scroll, lines, layout)
	base.Lines = visible
	base.OffsetY = offsetY
	return base
}

// buildNoSourceView is the neutral status shown when no data source is set.
func buildNoSourceView() style.ViewData {
	return style.ViewData{
		Items:  []string{"[pager] no source"},
		
	}
}

// RenderCacheKey returns a change-detection key incorporating the buffer
// sequence, the resolved scroll offset or page phase+alpha, and a policy
// fingerprint.
func renderCacheKey(typer *source.TyperState) uint32 {
	_, buf, scroll, page := source.ActiveStateSnapshot()

	pol := GetPolicy()
	fp := pol.Fingerprint()

	var seq uint64
	if buf != nil {
		seq = buf.Seq()
	}

	var token string
	switch {
	case typer != nil:
		return region.CalcRegionCacheKey("pager", seq, typer.Revision(), fp)
	case scroll != nil:
		token = itoa(scroll.OffsetPx())
	case page != nil:
		token = itoa(int(page.Phase())) + ":" + itoa(int(page.FadeAlpha()*100))
	default:
		token = "0"
	}

	return region.CalcRegionCacheKey("pager", seq, token, fp)
}

// itoa is a tiny strconv.Itoa avoiding an import churn in the cache-key path.
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var b [20]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		b[i] = '-'
	}
	return string(b[i:])
}
