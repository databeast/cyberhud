package displaymode

import (
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// GrayscaleFast480x320Style provides grayscale-fast rendering for the 480×320
// color-capable panel when color or rapid refresh is unavailable.
type GrayscaleFast480x320Style struct{}

func (s GrayscaleFast480x320Style) Name() string { return "grayscale-fast-480x320" }

// Requirements returns SurfaceRequirements for grayscale-fast 480×320 rendering.
// Monochrome fallback for color-capable panel: 480×320, NeedsColor=false, NeedsRapidRefresh=false.
func (s GrayscaleFast480x320Style) Requirements() style.SurfaceRequirements {
	return style.SurfaceRequirements{MinWidth: 480, MinHeight: 320, Capability: style.GrayscaleFast}
}

func (s GrayscaleFast480x320Style) Supports(hints textlayout.TextHints) style.Fitness {
	return style.EvaluateFitness(s.Requirements(), hints)
}

func (s GrayscaleFast480x320Style) Build(snapshot Snapshot, pol Policy, ctx style.StyleContext) style.ViewData {
	return style.ViewData{Items: []string{"(template)"}}
}
