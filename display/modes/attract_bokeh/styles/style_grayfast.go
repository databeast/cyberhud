package styles

import (
	"fmt"

	"github.com/databeast/cyberhud/display/modes/attract_bokeh/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// ─── GrayscaleFast ────────────────────────────────────────────────────────────

// BokehGrayscaleFastStyle targets grayscale fast-refresh panels.
type BokehGrayscaleFastStyle struct {
	Width  int
	Height int
}

func (s BokehGrayscaleFastStyle) Name() string {
	return fmt.Sprintf("grayscale-fast-%dx%d", s.Width, s.Height)
}

func (s BokehGrayscaleFastStyle) Requirements() style.SurfaceRequirements {
	return style.SurfaceRequirements{
		MinWidth:        s.Width,
		MinHeight:       s.Height,
		Capability:      style.GrayscaleFast,
		MinRows:         0,
		MinCharsPerLine: 0,
	}
}

func (s BokehGrayscaleFastStyle) Supports(hints textlayout.TextHints) style.Fitness {
	return style.EvaluateFitness(s.Requirements(), hints)
}

func (s BokehGrayscaleFastStyle) Build(_ source.BokehFrame, _ source.Policy, _ style.StyleContext) style.ViewData {
	return style.ViewData{}
}
