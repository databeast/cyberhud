package displaymode

import (
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// EinkSmall176x264Style targets the 176×264 e-ink panel.
type EinkSmall176x264Style struct{}

func (s EinkSmall176x264Style) Name() string { return "eink-176x264" }

// Requirements returns SurfaceRequirements for the 176×264 panel.
// Target: 176×264, 1-bit Mono packed, Deferred refresh.
func (s EinkSmall176x264Style) Requirements() style.SurfaceRequirements {
	return style.SurfaceRequirements{MinWidth: 176, MinHeight: 264, Capability: style.MonoSlow}
}

func (s EinkSmall176x264Style) Supports(hints textlayout.TextHints) style.Fitness {
	return style.EvaluateFitness(s.Requirements(), hints)
}

func (s EinkSmall176x264Style) Build(snapshot Snapshot, pol Policy, ctx style.StyleContext) style.ViewData {
	return style.ViewData{Items: []string{"(template)"}}
}
