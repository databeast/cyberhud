package displaymode

import (
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// GrayscaleFast320x240Style provides grayscale-fast rendering for the 320×240
// color-capable panel when color or rapid refresh is unavailable.
type GrayscaleFast320x240Style struct{}

func (s GrayscaleFast320x240Style) Name() string { return "grayscale-fast-320x240" }

// Requirements returns SurfaceRequirements for grayscale-fast 320×240 rendering.
// Monochrome fallback for color-capable panel: 320×240, NeedsColor=false, NeedsRapidRefresh=false.
func (s GrayscaleFast320x240Style) Requirements() style.SurfaceRequirements {
	return style.SurfaceRequirements{MinWidth: 320, MinHeight: 240, Capability: style.GrayscaleFast}
}

func (s GrayscaleFast320x240Style) Supports(hints textlayout.TextHints) style.Fitness {
	return style.EvaluateFitness(s.Requirements(), hints)
}

func (s GrayscaleFast320x240Style) Build(snapshot Snapshot, pol Policy, ctx style.StyleContext) style.ViewData {
	return style.ViewData{Items: []string{"(template)"}}
}
