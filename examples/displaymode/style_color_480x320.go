package displaymode

import (
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// ColorLarge480x320Style targets the 480×320 color TFT panel.
type ColorLarge480x320Style struct{}

func (s ColorLarge480x320Style) Name() string { return "color-480x320" }

// Requirements returns SurfaceRequirements for the 480×320 panel.
// Target: 480×320, RGB565/RGB666, Continuous refresh.
func (s ColorLarge480x320Style) Requirements() style.SurfaceRequirements {
	return style.SurfaceRequirements{MinWidth: 480, MinHeight: 320, Capability: style.ColorFast}
}

func (s ColorLarge480x320Style) Supports(hints textlayout.TextHints) style.Fitness {
	return style.EvaluateFitness(s.Requirements(), hints)
}

func (s ColorLarge480x320Style) Build(snapshot Snapshot, pol Policy, ctx style.StyleContext) style.ViewData {
	return style.ViewData{Items: []string{"(template)"}}
}
