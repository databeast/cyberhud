package displaymode

import (
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// MonoSlow176x264Style targets the 176×264 e-ink panel.
type MonoSlow176x264Style struct{}

func (s MonoSlow176x264Style) Name() string { return "mono-slow-176x264" }

// Requirements returns SurfaceRequirements for the 176×264 panel.
// Target: 176×264, 1-bit Mono packed, Deferred refresh.
func (s MonoSlow176x264Style) Requirements() style.SurfaceRequirements {
	return style.SurfaceRequirements{MinWidth: 176, MinHeight: 264, Capability: style.MonoSlow}
}

func (s MonoSlow176x264Style) Supports(hints textlayout.TextHints) style.Fitness {
	return style.EvaluateFitness(s.Requirements(), hints)
}

func (s MonoSlow176x264Style) Build(snapshot Snapshot, pol Policy, ctx style.StyleContext) style.ViewData {
	return style.ViewData{Items: []string{"(template)"}}
}
