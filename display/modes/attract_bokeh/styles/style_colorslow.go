package styles

import (
	"fmt"

	"github.com/databeast/cyberhud/display/modes/attract_bokeh/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// ─── ColorSlow ────────────────────────────────────────────────────────────────

// BokehColorSlowStyle targets color slow-refresh e-ink panels.
type BokehColorSlowStyle struct {
	Width  int
	Height int
}

func (s BokehColorSlowStyle) Name() string {
	return fmt.Sprintf("color-slow-%dx%d", s.Width, s.Height)
}

func (s BokehColorSlowStyle) Requirements() style.SurfaceRequirements {
	return style.SurfaceRequirements{
		MinWidth:        s.Width,
		MinHeight:       s.Height,
		Capability:      style.ColorSlow,
		MinRows:         0,
		MinCharsPerLine: 0,
	}
}

func (s BokehColorSlowStyle) Supports(hints textlayout.TextHints) style.Fitness {
	return style.EvaluateFitness(s.Requirements(), hints)
}

func (s BokehColorSlowStyle) Build(_ source.BokehFrame, _ source.Policy, _ style.StyleContext) style.ViewData {
	return style.ViewData{Static: true}
}
