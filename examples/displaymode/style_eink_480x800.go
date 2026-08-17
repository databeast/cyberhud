package displaymode

import (
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// EinkLarge480x800Style targets the 480×800 e-ink panel.
type EinkLarge480x800Style struct{}

func (s EinkLarge480x800Style) Name() string { return "eink-480x800" }

// Requirements returns SurfaceRequirements for the 480×800 panel.
// Target: 480×800, 1-bit/2-bit Mono packed, Deferred refresh.
func (s EinkLarge480x800Style) Requirements() style.SurfaceRequirements {
	return style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.MonoSlow}
}

func (s EinkLarge480x800Style) Supports(hints textlayout.TextHints) style.Fitness {
	return style.EvaluateFitness(s.Requirements(), hints)
}

func (s EinkLarge480x800Style) Build(snapshot Snapshot, pol Policy, ctx style.StyleContext) style.ViewData {
	return style.ViewData{Items: []string{"(template)"}}
}
