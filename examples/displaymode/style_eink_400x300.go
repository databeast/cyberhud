package displaymode

import (
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// EinkMedium400x300Style targets the 400×300 e-ink panel.
type EinkMedium400x300Style struct{}

func (s EinkMedium400x300Style) Name() string { return "eink-400x300" }

// Requirements returns SurfaceRequirements for the 400×300 panel.
// Target: 400×300, 1-bit/2-bit Mono packed, Deferred refresh.
func (s EinkMedium400x300Style) Requirements() style.SurfaceRequirements {
	return style.SurfaceRequirements{MinWidth: 400, MinHeight: 300, Capability: style.MonoSlow}
}

func (s EinkMedium400x300Style) Supports(hints textlayout.TextHints) style.Fitness {
	return style.EvaluateFitness(s.Requirements(), hints)
}

func (s EinkMedium400x300Style) Build(snapshot Snapshot, pol Policy, ctx style.StyleContext) style.ViewData {
	return style.ViewData{Items: []string{"(template)"}}
}
