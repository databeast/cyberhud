package styles

import (
	"fmt"

	"github.com/databeast/cyberhud/display/modes/attract_bokeh/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// ─── ColorFast ────────────────────────────────────────────────────────────────

// BokehColorFastStyle targets color fast-refresh TFT panels.
type BokehColorFastStyle struct {
	Width  int
	Height int
}

func (s BokehColorFastStyle) Name() string {
	return fmt.Sprintf("color-%dx%d", s.Width, s.Height)
}

func (s BokehColorFastStyle) Requirements() style.SurfaceRequirements {
	return style.SurfaceRequirements{
		MinWidth:        s.Width,
		MinHeight:       s.Height,
		Capability:      style.ColorFast,
		MinRows:         0,
		MinCharsPerLine: 0,
	}
}

func (s BokehColorFastStyle) Supports(hints textlayout.TextHints) style.Fitness {
	return style.EvaluateFitness(s.Requirements(), hints)
}

func (s BokehColorFastStyle) Build(_ source.BokehFrame, _ source.Policy, _ style.StyleContext) style.ViewData {
	return style.ViewData{}
}
