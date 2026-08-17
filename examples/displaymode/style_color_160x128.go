package displaymode

import (
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// ColorSmall160x128Style targets the 160×128 color TFT panel.
type ColorSmall160x128Style struct{}

func (s ColorSmall160x128Style) Name() string { return "color-160x128" }

// Requirements returns SurfaceRequirements for the 160×128 panel.
// Target: 160×128, RGB565, Continuous refresh.
func (s ColorSmall160x128Style) Requirements() style.SurfaceRequirements {
	return style.SurfaceRequirements{MinWidth: 160, MinHeight: 128, Capability: style.ColorFast}
}

func (s ColorSmall160x128Style) Supports(hints textlayout.TextHints) style.Fitness {
	return style.EvaluateFitness(s.Requirements(), hints)
}

func (s ColorSmall160x128Style) Build(snapshot Snapshot, pol Policy, ctx style.StyleContext) style.ViewData {
	return style.ViewData{Items: []string{"(template)"}}
}
