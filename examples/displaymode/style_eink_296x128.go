package displaymode

import (
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// EinkMedium296x128Style targets the 296×128 e-ink panel.
type EinkMedium296x128Style struct{}

func (s EinkMedium296x128Style) Name() string { return "eink-296x128" }

// Requirements returns SurfaceRequirements for the 296×128 panel.
// Target: 296×128, 1-bit/2-bit Mono packed, Deferred refresh.
func (s EinkMedium296x128Style) Requirements() style.SurfaceRequirements {
	return style.SurfaceRequirements{MinWidth: 296, MinHeight: 128, Capability: style.MonoSlow}
}

func (s EinkMedium296x128Style) Supports(hints textlayout.TextHints) style.Fitness {
	return style.EvaluateFitness(s.Requirements(), hints)
}

func (s EinkMedium296x128Style) Build(snapshot Snapshot, pol Policy, ctx style.StyleContext) style.ViewData {
	return style.ViewData{Items: []string{"(template)"}}
}
