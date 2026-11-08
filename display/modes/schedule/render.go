package schedule

import (
	"fmt"
	"time"

	"github.com/databeast/cyberhud/display/modes/schedule/source"
	"github.com/databeast/cyberhud/display/region"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// BuildView resolves the schedule style for the panel and lets it place the
// current snapshot, stamped with the render instant so the active-event decision
// tracks the clock between fetches. All fetching and parsing lives in the source
// package; this only assembles the snapshot and dispatches through the registry.
func BuildView(now time.Time, hints textlayout.TextHints) style.ViewData {
	snap := source.SnapshotNow()
	snap.Now = now
	p := normalizePolicy(PolicySnapshot())

	s, reason := style.ResolveStyle(scheduleRegistry, hints, "schedule", p.Style)
	ctx := style.NewStyleContext(hints)
	vd := s.Build(snap, p, ctx)
	vd.StyleReport = style.ReportFor(s, reason, hints)
	vd.Cursor = -1
	return vd
}

// RenderCacheKey returns a change-detection token. It folds the snapshot
// sequence (bumped on every fetch/policy change), the current active-event
// identity at minute precision (so the view flips when an event starts or ends
// even without a new fetch), and the policy fingerprint.
func RenderCacheKey(now time.Time) uint32 {
	snap := source.SnapshotNow()
	p := PolicySnapshot()

	activeID := "idle"
	if ev, ok := snap.ActiveEvent(now); ok {
		activeID = ev.Summary + "@" + ev.Start.Format(time.RFC3339)
	} else if ev, ok := snap.NextEvent(now); ok {
		activeID = "next:" + ev.Summary + "@" + ev.Start.Format(time.RFC3339)
	}

	key := fmt.Sprintf("sched|%d|%t|%s|%s|%s",
		snap.Sequence, snap.Loaded, now.Format("200601021504"), activeID, p.Fingerprint())
	return region.CalcRegionCacheKey(key)
}
