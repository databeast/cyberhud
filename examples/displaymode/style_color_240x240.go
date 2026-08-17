package displaymode

import (
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// ColorMedium240x240Style targets the 240×240 color TFT panel.
type ColorMedium240x240Style struct{}

func (s ColorMedium240x240Style) Name() string { return "color-240x240" }

// Requirements returns SurfaceRequirements for the 240×240 panel.
// Target: 240×240, RGB565, Continuous refresh.
func (s ColorMedium240x240Style) Requirements() style.SurfaceRequirements {
	return style.SurfaceRequirements{MinWidth: 240, MinHeight: 240, Capability: style.ColorFast}
}

func (s ColorMedium240x240Style) Supports(hints textlayout.TextHints) style.Fitness {
	return style.EvaluateFitness(s.Requirements(), hints)
}

func (s ColorMedium240x240Style) Build(snapshot Snapshot, pol Policy, ctx style.StyleContext) style.ViewData {
	return style.ViewData{Items: []string{"(template)"}}
}
