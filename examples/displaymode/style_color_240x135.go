package displaymode

import (
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// ColorMedium240x135Style targets the 240×135 color TFT panel.
type ColorMedium240x135Style struct{}

func (s ColorMedium240x135Style) Name() string { return "color-240x135" }

// Requirements returns SurfaceRequirements for the 240×135 panel.
// Target: 240×135, RGB565, Continuous refresh.
func (s ColorMedium240x135Style) Requirements() style.SurfaceRequirements {
	return style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.ColorFast}
}

func (s ColorMedium240x135Style) Supports(hints textlayout.TextHints) style.Fitness {
	return style.EvaluateFitness(s.Requirements(), hints)
}

func (s ColorMedium240x135Style) Build(snapshot Snapshot, pol Policy, ctx style.StyleContext) style.ViewData {
	return style.ViewData{Items: []string{"(template)"}}
}
