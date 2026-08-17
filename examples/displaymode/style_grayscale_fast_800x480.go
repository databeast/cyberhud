package displaymode

import (
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// GrayscaleFast800x480Style provides grayscale-fast rendering for the 800×480
// color-capable panel when color or rapid refresh is unavailable.
type GrayscaleFast800x480Style struct{}

func (s GrayscaleFast800x480Style) Name() string { return "grayscale-fast-800x480" }

// Requirements returns SurfaceRequirements for grayscale-fast 800×480 rendering.
// Monochrome fallback for color-capable panel: 800×480, NeedsColor=false, NeedsRapidRefresh=false.
func (s GrayscaleFast800x480Style) Requirements() style.SurfaceRequirements {
	return style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.GrayscaleFast}
}

func (s GrayscaleFast800x480Style) Supports(hints textlayout.TextHints) style.Fitness {
	return style.EvaluateFitness(s.Requirements(), hints)
}

func (s GrayscaleFast800x480Style) Build(snapshot Snapshot, pol Policy, ctx style.StyleContext) style.ViewData {
	return style.ViewData{Items: []string{"(template)"}}
}
