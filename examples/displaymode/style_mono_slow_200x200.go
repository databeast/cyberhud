package displaymode

import (
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// MonoSlow200x200Style targets the 200×200 e-ink panel.
type MonoSlow200x200Style struct{}

func (s MonoSlow200x200Style) Name() string { return "mono-slow-200x200" }

// Requirements returns SurfaceRequirements for the 200×200 panel.
// Target: 200×200, 1-bit Mono packed, Deferred refresh.
func (s MonoSlow200x200Style) Requirements() style.SurfaceRequirements {
	return style.SurfaceRequirements{MinWidth: 200, MinHeight: 200, Capability: style.MonoSlow}
}

func (s MonoSlow200x200Style) Supports(hints textlayout.TextHints) style.Fitness {
	return style.EvaluateFitness(s.Requirements(), hints)
}

func (s MonoSlow200x200Style) Build(snapshot Snapshot, pol Policy, ctx style.StyleContext) style.ViewData {
	return style.ViewData{Items: []string{"(template)"}}
}
