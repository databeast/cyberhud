package displaymode

import (
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// ColorLarge800x480Style targets the 800×480 color TFT/IPS panel.
type ColorLarge800x480Style struct{}

func (s ColorLarge800x480Style) Name() string { return "color-800x480" }

// Requirements returns SurfaceRequirements for the 800×480 color panel.
// Target: 800×480, RGB888/RGB666, Continuous refresh.
func (s ColorLarge800x480Style) Requirements() style.SurfaceRequirements {
	return style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.ColorFast}
}

func (s ColorLarge800x480Style) Supports(hints textlayout.TextHints) style.Fitness {
	return style.EvaluateFitness(s.Requirements(), hints)
}

func (s ColorLarge800x480Style) Build(snapshot Snapshot, pol Policy, ctx style.StyleContext) style.ViewData {
	return style.ViewData{Items: []string{"(template)"}}
}
