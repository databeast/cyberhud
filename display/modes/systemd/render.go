package systemd

import (
	"github.com/databeast/cyberhud/display/modes/systemd/source"
	"github.com/databeast/cyberhud/display/region"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// BuildView returns the systemd mode view data by delegating to the appropriate
// style builder based on the active policy.
func BuildView(hints textlayout.TextHints) style.ViewData {
	snap := source.PollSnapshot()
	p := GetPolicy()
	s, reason := style.ResolveStyle(systemdRegistry, hints, "systemd", p.Style)
	ctx := style.NewStyleContext(hints)
	svd := s.Build(snap, p, ctx)
	svd.StyleReport = style.StyleReport{Name: s.Name(), Reason: reason}
	return svd
}

func RenderCacheKey() uint32 {
	return region.CalcRegionCacheKey(source.Signature(GetPolicy()))
}
