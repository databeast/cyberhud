package displaymode

import (
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// EinkSmall212x104Style targets the 212×104 e-ink panel.
type EinkSmall212x104Style struct{}

func (s EinkSmall212x104Style) Name() string { return "eink-212x104" }

// Requirements returns SurfaceRequirements for the 212×104 panel.
// Target: 212×104, 1-bit Mono packed, Deferred refresh.
func (s EinkSmall212x104Style) Requirements() style.SurfaceRequirements {
	return style.SurfaceRequirements{MinWidth: 212, MinHeight: 104, Capability: style.MonoSlow}
}

func (s EinkSmall212x104Style) Supports(hints textlayout.TextHints) style.Fitness {
	return style.EvaluateFitness(s.Requirements(), hints)
}

func (s EinkSmall212x104Style) Build(snapshot Snapshot, pol Policy, ctx style.StyleContext) style.ViewData {
	return style.ViewData{Items: []string{"(template)"}}
}
