package styles

import (
	"fmt"

	"github.com/databeast/cyberhud/display/modes/attract_bokeh/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// ─── MonoFast ─────────────────────────────────────────────────────────────────

// BokehMonoFastStyle targets monochrome fast-refresh OLED panels.
type BokehMonoFastStyle struct {
	Width  int
	Height int
}

func (s BokehMonoFastStyle) Name() string {
	return fmt.Sprintf("mono-%dx%d", s.Width, s.Height)
}

func (s BokehMonoFastStyle) Requirements() style.SurfaceRequirements {
	return style.SurfaceRequirements{
		MinWidth:        s.Width,
		MinHeight:       s.Height,
		Capability:      style.MonoFast,
		MinRows:         0,
		MinCharsPerLine: 0,
	}
}

func (s BokehMonoFastStyle) Supports(hints textlayout.TextHints) style.Fitness {
	return style.EvaluateFitness(s.Requirements(), hints)
}

func (s BokehMonoFastStyle) Build(_ source.BokehFrame, _ source.Policy, _ style.StyleContext) style.ViewData {
	return style.ViewData{}
}
