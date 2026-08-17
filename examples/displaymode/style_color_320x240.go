package displaymode

import (
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// ColorMedium320x240Style targets the 320×240 color TFT panel.
type ColorMedium320x240Style struct{}

func (s ColorMedium320x240Style) Name() string { return "color-320x240" }

// Requirements returns SurfaceRequirements for the 320×240 panel.
// Target: 320×240, RGB565, Continuous refresh.
func (s ColorMedium320x240Style) Requirements() style.SurfaceRequirements {
	return style.SurfaceRequirements{MinWidth: 320, MinHeight: 240, Capability: style.ColorFast}
}

func (s ColorMedium320x240Style) Supports(hints textlayout.TextHints) style.Fitness {
	return style.EvaluateFitness(s.Requirements(), hints)
}

func (s ColorMedium320x240Style) Build(snapshot Snapshot, pol Policy, ctx style.StyleContext) style.ViewData {
	return style.ViewData{Items: []string{"(template)"}}
}
