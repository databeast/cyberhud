package displaymode

import (
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// ColorSmall160x80Style targets the 160×80 color TFT panel.
type ColorSmall160x80Style struct{}

func (s ColorSmall160x80Style) Name() string { return "color-160x80" }

// Requirements returns SurfaceRequirements for the 160×80 panel.
// Target: 160×80, RGB565, Continuous refresh.
func (s ColorSmall160x80Style) Requirements() style.SurfaceRequirements {
	return style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.ColorFast}
}

func (s ColorSmall160x80Style) Supports(hints textlayout.TextHints) style.Fitness {
	return style.EvaluateFitness(s.Requirements(), hints)
}

func (s ColorSmall160x80Style) Build(snapshot Snapshot, pol Policy, ctx style.StyleContext) style.ViewData {
	return style.ViewData{Items: []string{"(template)"}}
}
