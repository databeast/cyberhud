package styles

import (
	"fmt"
	"strings"

	"github.com/databeast/cyberhud/display/modes/attract_geometric/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// Compile-time interface compliance checks.
var (
	_ style.Style[source.GeometricFrame, source.Policy] = GeometricColorFastStyle{}
	_ style.Style[source.GeometricFrame, source.Policy] = GeometricMonoFastStyle{}
	_ style.Style[source.GeometricFrame, source.Policy] = GeometricMonoSlowStyle{}
	_ style.Style[source.GeometricFrame, source.Policy] = GeometricGrayscaleFastStyle{}
	_ style.Style[source.GeometricFrame, source.Policy] = GeometricGrayscaleSlowStyle{}
	_ style.Style[source.GeometricFrame, source.Policy] = GeometricColorSlowStyle{}
)

// ─── ColorFast ────────────────────────────────────────────────────────────────

// GeometricColorFastStyle targets color fast-refresh TFT panels.
type GeometricColorFastStyle struct {
	Width  int
	Height int
}

func (s GeometricColorFastStyle) Name() string {
	return fmt.Sprintf("color-%dx%d", s.Width, s.Height)
}

func (s GeometricColorFastStyle) Requirements() style.SurfaceRequirements {
	return style.SurfaceRequirements{
		MinWidth:        s.Width,
		MinHeight:       s.Height,
		Capability:      style.ColorFast,
		MinRows:         0,
		MinCharsPerLine: 0,
	}
}

func (s GeometricColorFastStyle) Supports(hints textlayout.TextHints) style.Fitness {
	return style.EvaluateFitness(s.Requirements(), hints)
}

func (s GeometricColorFastStyle) Build(_ source.GeometricFrame, _ source.Policy, _ style.StyleContext) style.ViewData {
	return style.ViewData{}
}

// ─── MonoFast ─────────────────────────────────────────────────────────────────

// GeometricMonoFastStyle targets monochrome fast-refresh OLED panels.
type GeometricMonoFastStyle struct {
	Width  int
	Height int
}

func (s GeometricMonoFastStyle) Name() string {
	return fmt.Sprintf("mono-%dx%d", s.Width, s.Height)
}

func (s GeometricMonoFastStyle) Requirements() style.SurfaceRequirements {
	return style.SurfaceRequirements{
		MinWidth:        s.Width,
		MinHeight:       s.Height,
		Capability:      style.MonoFast,
		MinRows:         0,
		MinCharsPerLine: 0,
	}
}

func (s GeometricMonoFastStyle) Supports(hints textlayout.TextHints) style.Fitness {
	return style.EvaluateFitness(s.Requirements(), hints)
}

func (s GeometricMonoFastStyle) Build(_ source.GeometricFrame, _ source.Policy, _ style.StyleContext) style.ViewData {
	return style.ViewData{}
}

// ─── MonoSlow ─────────────────────────────────────────────────────────────────

// GeometricMonoSlowStyle targets monochrome slow-refresh panels.
type GeometricMonoSlowStyle struct {
	Width  int
	Height int
}

func (s GeometricMonoSlowStyle) Name() string {
	return fmt.Sprintf("mono-slow-%dx%d", s.Width, s.Height)
}

func (s GeometricMonoSlowStyle) Requirements() style.SurfaceRequirements {
	return style.SurfaceRequirements{
		MinWidth:        s.Width,
		MinHeight:       s.Height,
		Capability:      style.MonoSlow,
		MinRows:         0,
		MinCharsPerLine: 0,
	}
}

func (s GeometricMonoSlowStyle) Supports(hints textlayout.TextHints) style.Fitness {
	return style.EvaluateFitness(s.Requirements(), hints)
}

func (s GeometricMonoSlowStyle) Build(_ source.GeometricFrame, _ source.Policy, _ style.StyleContext) style.ViewData {
	return style.ViewData{Static: true}
}

// ─── GrayscaleFast ────────────────────────────────────────────────────────────

// GeometricGrayscaleFastStyle targets grayscale fast-refresh panels.
type GeometricGrayscaleFastStyle struct {
	Width  int
	Height int
}

func (s GeometricGrayscaleFastStyle) Name() string {
	return fmt.Sprintf("grayscale-fast-%dx%d", s.Width, s.Height)
}

func (s GeometricGrayscaleFastStyle) Requirements() style.SurfaceRequirements {
	return style.SurfaceRequirements{
		MinWidth:        s.Width,
		MinHeight:       s.Height,
		Capability:      style.GrayscaleFast,
		MinRows:         0,
		MinCharsPerLine: 0,
	}
}

func (s GeometricGrayscaleFastStyle) Supports(hints textlayout.TextHints) style.Fitness {
	return style.EvaluateFitness(s.Requirements(), hints)
}

func (s GeometricGrayscaleFastStyle) Build(_ source.GeometricFrame, _ source.Policy, _ style.StyleContext) style.ViewData {
	return style.ViewData{}
}

// ─── GrayscaleSlow ────────────────────────────────────────────────────────────

// GeometricGrayscaleSlowStyle targets grayscale slow-refresh panels.
type GeometricGrayscaleSlowStyle struct {
	Width  int
	Height int
}

func (s GeometricGrayscaleSlowStyle) Name() string {
	return fmt.Sprintf("grayscale-slow-%dx%d", s.Width, s.Height)
}

func (s GeometricGrayscaleSlowStyle) Requirements() style.SurfaceRequirements {
	return style.SurfaceRequirements{
		MinWidth:        s.Width,
		MinHeight:       s.Height,
		Capability:      style.GrayscaleSlow,
		MinRows:         0,
		MinCharsPerLine: 0,
	}
}

func (s GeometricGrayscaleSlowStyle) Supports(hints textlayout.TextHints) style.Fitness {
	return style.EvaluateFitness(s.Requirements(), hints)
}

func (s GeometricGrayscaleSlowStyle) Build(_ source.GeometricFrame, _ source.Policy, _ style.StyleContext) style.ViewData {
	return style.ViewData{Static: true}
}

// ─── ColorSlow ────────────────────────────────────────────────────────────────

// GeometricColorSlowStyle targets color slow-refresh panels.
type GeometricColorSlowStyle struct {
	Width  int
	Height int
}

func (s GeometricColorSlowStyle) Name() string {
	return fmt.Sprintf("color-slow-%dx%d", s.Width, s.Height)
}

func (s GeometricColorSlowStyle) Requirements() style.SurfaceRequirements {
	return style.SurfaceRequirements{
		MinWidth:        s.Width,
		MinHeight:       s.Height,
		Capability:      style.ColorSlow,
		MinRows:         0,
		MinCharsPerLine: 0,
	}
}

func (s GeometricColorSlowStyle) Supports(hints textlayout.TextHints) style.Fitness {
	return style.EvaluateFitness(s.Requirements(), hints)
}

func (s GeometricColorSlowStyle) Build(_ source.GeometricFrame, _ source.Policy, _ style.StyleContext) style.ViewData {
	return style.ViewData{Static: true}
}

// Registry returns the mode's style registry.
func Registry() *style.StyleRegistry[source.GeometricFrame, source.Policy] {
	return geometricRegistry
}

// ResolveBestStyleName returns the name of the best-fit style for the given hints.
func ResolveBestStyleName(hints textlayout.TextHints) string {
	return resolveBestStyleName(hints)
}

// ─── Registry ─────────────────────────────────────────────────────────────────

// geometricRegistry is the per-mode StyleRegistry for the geometric attract display mode.
// Registration order follows capability ordering: MonoSlow → MonoFast →
// GrayscaleSlow → GrayscaleFast → ColorSlow → ColorFast.
var geometricRegistry = func() *style.StyleRegistry[source.GeometricFrame, source.Policy] {
	r := style.NewRegistry[source.GeometricFrame, source.Policy](
		// MonoSlow (e-paper mono)
		GeometricMonoSlowStyle{Width: 122, Height: 250},
		GeometricMonoSlowStyle{Width: 176, Height: 264},
		GeometricMonoSlowStyle{Width: 200, Height: 200},
		GeometricMonoSlowStyle{Width: 212, Height: 104},
		GeometricMonoSlowStyle{Width: 296, Height: 128},
		GeometricMonoSlowStyle{Width: 400, Height: 300},
		GeometricMonoSlowStyle{Width: 480, Height: 800},
		GeometricMonoSlowStyle{Width: 800, Height: 480},
		GeometricMonoSlowStyle{Width: 104, Height: 212},
		GeometricMonoSlowStyle{Width: 250, Height: 122},
		GeometricMonoSlowStyle{Width: 128, Height: 296},
		GeometricMonoSlowStyle{Width: 264, Height: 176},
		GeometricMonoSlowStyle{Width: 300, Height: 400},

		// MonoFast (OLED mono)
		GeometricMonoFastStyle{Width: 16, Height: 8},
		GeometricMonoFastStyle{Width: 8, Height: 16},
		GeometricMonoFastStyle{Width: 128, Height: 32},
		GeometricMonoFastStyle{Width: 128, Height: 64},
		GeometricMonoFastStyle{Width: 128, Height: 128},
		GeometricMonoFastStyle{Width: 32, Height: 128},
		GeometricMonoFastStyle{Width: 64, Height: 128},

		// GrayscaleSlow (grayscale e-paper)
		GeometricGrayscaleSlowStyle{Width: 122, Height: 250},
		GeometricGrayscaleSlowStyle{Width: 176, Height: 264},
		GeometricGrayscaleSlowStyle{Width: 200, Height: 200},
		GeometricGrayscaleSlowStyle{Width: 212, Height: 104},
		GeometricGrayscaleSlowStyle{Width: 296, Height: 128},
		GeometricGrayscaleSlowStyle{Width: 400, Height: 300},
		GeometricGrayscaleSlowStyle{Width: 480, Height: 800},
		GeometricGrayscaleSlowStyle{Width: 800, Height: 480},
		GeometricGrayscaleSlowStyle{Width: 104, Height: 212},
		GeometricGrayscaleSlowStyle{Width: 250, Height: 122},
		GeometricGrayscaleSlowStyle{Width: 128, Height: 296},
		GeometricGrayscaleSlowStyle{Width: 264, Height: 176},
		GeometricGrayscaleSlowStyle{Width: 300, Height: 400},

		// GrayscaleFast (grayscale LED matrix)
		GeometricGrayscaleFastStyle{Width: 16, Height: 8},
		GeometricGrayscaleFastStyle{Width: 8, Height: 16},
		GeometricGrayscaleFastStyle{Width: 160, Height: 80},
		GeometricGrayscaleFastStyle{Width: 160, Height: 128},
		GeometricGrayscaleFastStyle{Width: 240, Height: 135},
		GeometricGrayscaleFastStyle{Width: 240, Height: 240},
		GeometricGrayscaleFastStyle{Width: 320, Height: 240},
		GeometricGrayscaleFastStyle{Width: 480, Height: 320},
		GeometricGrayscaleFastStyle{Width: 800, Height: 480},
		GeometricGrayscaleFastStyle{Width: 80, Height: 160},
		GeometricGrayscaleFastStyle{Width: 128, Height: 160},
		GeometricGrayscaleFastStyle{Width: 135, Height: 240},
		GeometricGrayscaleFastStyle{Width: 240, Height: 320},
		GeometricGrayscaleFastStyle{Width: 320, Height: 480},
		GeometricGrayscaleFastStyle{Width: 480, Height: 800},
		GeometricGrayscaleFastStyle{Width: 128, Height: 128},

		// ColorSlow (color e-paper)
		GeometricColorSlowStyle{Width: 122, Height: 250},
		GeometricColorSlowStyle{Width: 176, Height: 264},
		GeometricColorSlowStyle{Width: 200, Height: 200},
		GeometricColorSlowStyle{Width: 212, Height: 104},
		GeometricColorSlowStyle{Width: 296, Height: 128},
		GeometricColorSlowStyle{Width: 400, Height: 300},
		GeometricColorSlowStyle{Width: 480, Height: 800},
		GeometricColorSlowStyle{Width: 800, Height: 480},
		GeometricColorSlowStyle{Width: 104, Height: 212},
		GeometricColorSlowStyle{Width: 250, Height: 122},
		GeometricColorSlowStyle{Width: 128, Height: 296},
		GeometricColorSlowStyle{Width: 264, Height: 176},
		GeometricColorSlowStyle{Width: 300, Height: 400},

		// ColorFast (color TFT)
		GeometricColorFastStyle{Width: 16, Height: 8},
		GeometricColorFastStyle{Width: 8, Height: 16},
		GeometricColorFastStyle{Width: 160, Height: 80},
		GeometricColorFastStyle{Width: 160, Height: 128},
		GeometricColorFastStyle{Width: 240, Height: 135},
		GeometricColorFastStyle{Width: 240, Height: 240},
		GeometricColorFastStyle{Width: 320, Height: 240},
		GeometricColorFastStyle{Width: 480, Height: 320},
		GeometricColorFastStyle{Width: 800, Height: 480},
		GeometricColorFastStyle{Width: 80, Height: 160},
		GeometricColorFastStyle{Width: 128, Height: 160},
		GeometricColorFastStyle{Width: 135, Height: 240},
		GeometricColorFastStyle{Width: 240, Height: 320},
		GeometricColorFastStyle{Width: 320, Height: 480},
		GeometricColorFastStyle{Width: 480, Height: 800},
		GeometricColorFastStyle{Width: 128, Height: 128},
	)
	return r
}()

// ─── Helpers ──────────────────────────────────────────────────────────────────

// resolveBestStyleName returns the name of the best-fit style for the given hints.
func resolveBestStyleName(hints textlayout.TextHints) string {
	s, _ := style.ResolveStyle(geometricRegistry, hints, "attract_geometric", "")
	return s.Name()
}

// isSlowRefresh returns true when the best-fit style for the given hints
// is a slow-refresh variant (e-paper / e-ink).
func isSlowRefresh(hints textlayout.TextHints) bool {
	name := resolveBestStyleName(hints)
	return strings.Contains(name, "-slow-")
}
