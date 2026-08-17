package gauges

import (
	"github.com/databeast/cyberhud/display/region"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// BuildView selects a style and renders the current gauges snapshot.
func BuildView(hints textlayout.TextHints) style.ViewData {
	pol := sourcePolicySnapshot()
	snap := sourceSnapshot()

	resolved, reason := style.ResolveStyle(gaugesRegistry, hints, "gauges", pol.Style)
	ctx := style.NewStyleContext(hints)

	vd := resolved.Build(snap, pol, ctx)
	vd.StyleReport = style.StyleReport{Name: resolved.Name(), Reason: reason}
	vd.Static = true
	if vd.PaddingPct == 0 {
		vd.PaddingPct = pol.PaddingPct
	}
	return vd
}

// RenderCacheKey returns a fingerprint that changes when either the gauges
// payload or the policy changes.
func RenderCacheKey() uint32 {
	pol := sourcePolicySnapshot()
	snap := sourceSnapshot()
	return region.CalcRegionCacheKey("gauges", sourceVersion(), pol.Fingerprint(), snap.Fingerprint())
}
