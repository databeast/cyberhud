package styles

import (
	"fmt"

	"github.com/databeast/cyberhud/display/modes/attract_bokeh/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// ─── MonoSlow ─────────────────────────────────────────────────────────────────

// BokehMonoSlowStyle targets monochrome slow-refresh e-ink panels.
type BokehMonoSlowStyle struct {
	Width  int
	Height int
}

func (s BokehMonoSlowStyle) Name() string {
	return fmt.Sprintf("mono-slow-%dx%d", s.Width, s.Height)
}

func (s BokehMonoSlowStyle) Requirements() style.SurfaceRequirements {
	return style.SurfaceRequirements{
		MinWidth:        s.Width,
		MinHeight:       s.Height,
		Capability:      style.MonoSlow,
		MinRows:         0,
		MinCharsPerLine: 0,
	}
}

func (s BokehMonoSlowStyle) Supports(hints textlayout.TextHints) style.Fitness {
	return style.EvaluateFitness(s.Requirements(), hints)
}

func (s BokehMonoSlowStyle) Build(_ source.BokehFrame, _ source.Policy, _ style.StyleContext) style.ViewData {
	return style.ViewData{Static: true}
}
