package menu

import (
	"github.com/databeast/cyberhud/display/modes/menu/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// BuildView returns the menu mode view data using registry-based style dispatch.
func BuildView(hints textlayout.TextHints, borderBuilder BorderBuilder) style.ViewData {
	p := GetPolicy()

	s, reason := style.ResolveStyle(menuRegistry, hints, "menu", p.Style)

	snap := source.BuildData(0, 0)
	ctx := style.NewStyleContext(hints)
	svd := s.Build(snap, p, ctx)
	svd.StyleReport = style.StyleReport{Name: s.Name(), Reason: reason}
	svd.Static = false

	return svd
}
