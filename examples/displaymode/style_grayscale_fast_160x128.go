package displaymode

import (
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// GrayscaleFast160x128Style provides grayscale-fast rendering for the 160×128
// color-capable panel when color or rapid refresh is unavailable.
type GrayscaleFast160x128Style struct{}

func (s GrayscaleFast160x128Style) Name() string { return "grayscale-fast-160x128" }

// Requirements returns SurfaceRequirements for grayscale-fast 160×128 rendering.
// Monochrome fallback for color-capable panel: 160×128, NeedsColor=false, NeedsRapidRefresh=false.
func (s GrayscaleFast160x128Style) Requirements() style.SurfaceRequirements {
	return style.SurfaceRequirements{MinWidth: 160, MinHeight: 128, Capability: style.GrayscaleFast}
}

func (s GrayscaleFast160x128Style) Supports(hints textlayout.TextHints) style.Fitness {
	return style.EvaluateFitness(s.Requirements(), hints)
}

func (s GrayscaleFast160x128Style) Build(snapshot Snapshot, pol Policy, ctx style.StyleContext) style.ViewData {
	return style.ViewData{Items: []string{"(template)"}}
}
