package displaymode

import (
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// MonoSmall128x32Style targets the 128×32 monochrome OLED panel.
type MonoSmall128x32Style struct{}

func (s MonoSmall128x32Style) Name() string { return "mono-128x32" }

// Requirements returns SurfaceRequirements for the 128×32 panel.
// Target: 128×32, 1-bit Mono packed, Deferred refresh.
func (s MonoSmall128x32Style) Requirements() style.SurfaceRequirements {
	return style.SurfaceRequirements{MinWidth: 128, MinHeight: 32, Capability: style.MonoFast}
}

func (s MonoSmall128x32Style) Supports(hints textlayout.TextHints) style.Fitness {
	return style.EvaluateFitness(s.Requirements(), hints)
}

func (s MonoSmall128x32Style) Build(snapshot Snapshot, pol Policy, ctx style.StyleContext) style.ViewData {
	return style.ViewData{Items: []string{"(template)"}}
}
