package displaymode

import (
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// MonoSlow122x250Style targets the 122×250 e-ink panel.
type MonoSlow122x250Style struct{}

func (s MonoSlow122x250Style) Name() string { return "mono-slow-122x250" }

// Requirements returns SurfaceRequirements for the 122×250 panel.
// Target: 122×250, 1-bit Mono packed, Deferred refresh.
func (s MonoSlow122x250Style) Requirements() style.SurfaceRequirements {
	return style.SurfaceRequirements{MinWidth: 122, MinHeight: 250, Capability: style.MonoSlow}
}

func (s MonoSlow122x250Style) Supports(hints textlayout.TextHints) style.Fitness {
	return style.EvaluateFitness(s.Requirements(), hints)
}

func (s MonoSlow122x250Style) Build(snapshot Snapshot, pol Policy, ctx style.StyleContext) style.ViewData {
	return style.ViewData{Items: []string{"(template)"}}
}
