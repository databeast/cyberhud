package styles

import (
	"fmt"

	"github.com/databeast/cyberhud/display/modes/attract_bokeh/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// ─── GrayscaleSlow ────────────────────────────────────────────────────────────

// BokehGrayscaleSlowStyle targets grayscale slow-refresh e-ink panels.
type BokehGrayscaleSlowStyle struct {
	Width  int
	Height int
}

func (s BokehGrayscaleSlowStyle) Name() string {
	return fmt.Sprintf("grayscale-slow-%dx%d", s.Width, s.Height)
}

func (s BokehGrayscaleSlowStyle) Requirements() style.SurfaceRequirements {
	return style.SurfaceRequirements{
		MinWidth:        s.Width,
		MinHeight:       s.Height,
		Capability:      style.GrayscaleSlow,
		MinRows:         0,
		MinCharsPerLine: 0,
	}
}

func (s BokehGrayscaleSlowStyle) Supports(hints textlayout.TextHints) style.Fitness {
	return style.EvaluateFitness(s.Requirements(), hints)
}

func (s BokehGrayscaleSlowStyle) Build(_ source.BokehFrame, _ source.Policy, _ style.StyleContext) style.ViewData {
	return style.ViewData{Static: true}
}
