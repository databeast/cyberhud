package displaymode

import (
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// MonoSlow800x480Style targets the 800×480 e-ink panel.
type MonoSlow800x480Style struct{}

func (s MonoSlow800x480Style) Name() string { return "mono-slow-800x480" }

// Requirements returns SurfaceRequirements for the 800×480 e-ink panel.
// Target: 800×480, 1-bit/2-bit Mono packed, Deferred refresh.
func (s MonoSlow800x480Style) Requirements() style.SurfaceRequirements {
	return style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.MonoSlow}
}

func (s MonoSlow800x480Style) Supports(hints textlayout.TextHints) style.Fitness {
	return style.EvaluateFitness(s.Requirements(), hints)
}

func (s MonoSlow800x480Style) Build(snapshot Snapshot, pol Policy, ctx style.StyleContext) style.ViewData {
	return style.ViewData{Items: []string{"(template)"}}
}
