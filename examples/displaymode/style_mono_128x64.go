package displaymode

import (
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// MonoSmall128x64Style targets the 128×64 monochrome OLED panel.
type MonoSmall128x64Style struct{}

func (s MonoSmall128x64Style) Name() string { return "mono-128x64" }

// Requirements returns SurfaceRequirements for the 128×64 panel.
// Target: 128×64, 1-bit Mono packed, Deferred refresh.
func (s MonoSmall128x64Style) Requirements() style.SurfaceRequirements {
	return style.SurfaceRequirements{MinWidth: 128, MinHeight: 64, Capability: style.MonoFast}
}

func (s MonoSmall128x64Style) Supports(hints textlayout.TextHints) style.Fitness {
	return style.EvaluateFitness(s.Requirements(), hints)
}

func (s MonoSmall128x64Style) Build(snapshot Snapshot, pol Policy, ctx style.StyleContext) style.ViewData {
	return style.ViewData{Items: []string{"(template)"}}
}
