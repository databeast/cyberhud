package displaymode

import (
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// GrayscaleFast240x135Style provides grayscale-fast rendering for the 240×135
// color-capable panel when color or rapid refresh is unavailable.
type GrayscaleFast240x135Style struct{}

func (s GrayscaleFast240x135Style) Name() string { return "grayscale-fast-240x135" }

// Requirements returns SurfaceRequirements for grayscale-fast 240×135 rendering.
// Monochrome fallback for color-capable panel: 240×135, NeedsColor=false, NeedsRapidRefresh=false.
func (s GrayscaleFast240x135Style) Requirements() style.SurfaceRequirements {
	return style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.GrayscaleFast}
}

func (s GrayscaleFast240x135Style) Supports(hints textlayout.TextHints) style.Fitness {
	return style.EvaluateFitness(s.Requirements(), hints)
}

func (s GrayscaleFast240x135Style) Build(snapshot Snapshot, pol Policy, ctx style.StyleContext) style.ViewData {
	return style.ViewData{Items: []string{"(template)"}}
}
