package styles

import (
	"fmt"

	"github.com/databeast/cyberhud/display/modes/attract_matrix/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/style/layout"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// Compile-time interface compliance checks.
var (
	// MonoSlow
	_ style.Style[source.MatrixSnapshot, source.Policy] = MonoSlowStyle{}

	// MonoFast
	_ style.Style[source.MatrixSnapshot, source.Policy] = MonoStyle{}

	// GrayscaleSlow
	_ style.Style[source.MatrixSnapshot, source.Policy] = GrayscaleSlowStyle{}

	// GrayscaleFast
	_ style.Style[source.MatrixSnapshot, source.Policy] = GrayscaleFastStyle{}

	// ColorSlow
	_ style.Style[source.MatrixSnapshot, source.Policy] = ColorSlowStyle{}

	// ColorFast
	_ style.Style[source.MatrixSnapshot, source.Policy] = ColorStyle{}

	// EinkStyle (MonoSlow, legacy parameterized type)
	_ style.Style[source.MatrixSnapshot, source.Policy] = EinkStyle{}
)

// ─── ColorFast ────────────────────────────────────────────────────────────────

// ColorStyle targets color TFT panels at the specified resolution.
type ColorStyle struct {
	Width  int
	Height int
}

func (s ColorStyle) Name() string { return fmt.Sprintf("color-%dx%d", s.Width, s.Height) }

func (s ColorStyle) Requirements() style.SurfaceRequirements {
	return style.SurfaceRequirements{
		MinWidth:   s.Width,
		MinHeight:  s.Height,
		Capability: style.ColorFast,
	}
}

func (s ColorStyle) Supports(hints textlayout.TextHints) style.Fitness {
	return style.EvaluateFitness(s.Requirements(), hints)
}

func (s ColorStyle) Build(snapshot source.MatrixSnapshot, _ source.Policy, ctx style.StyleContext) style.ViewData {
	hints := ctx.Hints()
	_ = layout.NewLayoutBridge(hints, layout.BridgeConfig{PaddingPct: 0})
	return buildFrame(snapshot)
}

// ─── MonoFast ─────────────────────────────────────────────────────────────────

// MonoStyle targets monochrome OLED panels at the specified resolution.
type MonoStyle struct {
	Width  int
	Height int
}

func (s MonoStyle) Name() string { return fmt.Sprintf("mono-%dx%d", s.Width, s.Height) }

func (s MonoStyle) Requirements() style.SurfaceRequirements {
	return style.SurfaceRequirements{
		MinWidth:   s.Width,
		MinHeight:  s.Height,
		Capability: style.MonoFast,
	}
}

func (s MonoStyle) Supports(hints textlayout.TextHints) style.Fitness {
	return style.EvaluateFitness(s.Requirements(), hints)
}

func (s MonoStyle) Build(snapshot source.MatrixSnapshot, _ source.Policy, ctx style.StyleContext) style.ViewData {
	hints := ctx.Hints()
	_ = layout.NewLayoutBridge(hints, layout.BridgeConfig{PaddingPct: 0})
	return buildFrame(snapshot)
}

// ─── MonoSlow (e-ink) ─────────────────────────────────────────────────────────

// EinkStyle targets e-ink panels at the specified resolution.
type EinkStyle struct {
	Width  int
	Height int
}

func (s EinkStyle) Name() string { return fmt.Sprintf("eink-%dx%d", s.Width, s.Height) }

func (s EinkStyle) Requirements() style.SurfaceRequirements {
	return style.SurfaceRequirements{
		MinWidth:   s.Width,
		MinHeight:  s.Height,
		Capability: style.MonoSlow,
	}
}

func (s EinkStyle) Supports(hints textlayout.TextHints) style.Fitness {
	return style.EvaluateFitness(s.Requirements(), hints)
}

func (s EinkStyle) Build(snapshot source.MatrixSnapshot, _ source.Policy, ctx style.StyleContext) style.ViewData {
	hints := ctx.Hints()
	_ = layout.NewLayoutBridge(hints, layout.BridgeConfig{PaddingPct: 0})
	return buildFrame(snapshot)
}

// ─── MonoSlowStyle (parameterized skeleton) ───────────────────────────────────

// MonoSlowStyle is a skeleton style for mono slow-refresh (e-paper mono) panels.
// Capability: MonoSlow (1-bit, slow refresh).
type MonoSlowStyle struct {
	Width  int
	Height int
}

func (s MonoSlowStyle) Name() string { return fmt.Sprintf("mono-slow-%dx%d", s.Width, s.Height) }

func (s MonoSlowStyle) Requirements() style.SurfaceRequirements {
	return style.SurfaceRequirements{
		MinWidth:   s.Width,
		MinHeight:  s.Height,
		Capability: style.MonoSlow,
	}
}

func (s MonoSlowStyle) Supports(hints textlayout.TextHints) style.Fitness {
	return style.EvaluateFitness(s.Requirements(), hints)
}

func (s MonoSlowStyle) Build(_ source.MatrixSnapshot, _ source.Policy, ctx style.StyleContext) style.ViewData {
	hints := ctx.Hints()
	_ = layout.NewLayoutBridge(hints, layout.BridgeConfig{PaddingPct: 0})
	return style.ViewData{}
}

// ─── GrayscaleSlowStyle (parameterized skeleton) ──────────────────────────────

// GrayscaleSlowStyle is a skeleton style for grayscale slow-refresh (grayscale e-paper) panels.
// Capability: GrayscaleSlow (multi-level luminance, slow refresh).
type GrayscaleSlowStyle struct {
	Width  int
	Height int
}

func (s GrayscaleSlowStyle) Name() string {
	return fmt.Sprintf("grayscale-slow-%dx%d", s.Width, s.Height)
}

func (s GrayscaleSlowStyle) Requirements() style.SurfaceRequirements {
	return style.SurfaceRequirements{
		MinWidth:   s.Width,
		MinHeight:  s.Height,
		Capability: style.GrayscaleSlow,
	}
}

func (s GrayscaleSlowStyle) Supports(hints textlayout.TextHints) style.Fitness {
	return style.EvaluateFitness(s.Requirements(), hints)
}

func (s GrayscaleSlowStyle) Build(_ source.MatrixSnapshot, _ source.Policy, ctx style.StyleContext) style.ViewData {
	hints := ctx.Hints()
	_ = layout.NewLayoutBridge(hints, layout.BridgeConfig{PaddingPct: 0})
	return style.ViewData{}
}

// ─── GrayscaleFastStyle (parameterized skeleton) ──────────────────────────────

// GrayscaleFastStyle is a skeleton style for grayscale fast-refresh (grayscale LED matrix) panels.
// Capability: GrayscaleFast (multi-level luminance, fast refresh).
type GrayscaleFastStyle struct {
	Width  int
	Height int
}

func (s GrayscaleFastStyle) Name() string {
	return fmt.Sprintf("grayscale-fast-%dx%d", s.Width, s.Height)
}

func (s GrayscaleFastStyle) Requirements() style.SurfaceRequirements {
	return style.SurfaceRequirements{
		MinWidth:   s.Width,
		MinHeight:  s.Height,
		Capability: style.GrayscaleFast,
	}
}

func (s GrayscaleFastStyle) Supports(hints textlayout.TextHints) style.Fitness {
	return style.EvaluateFitness(s.Requirements(), hints)
}

func (s GrayscaleFastStyle) Build(_ source.MatrixSnapshot, _ source.Policy, ctx style.StyleContext) style.ViewData {
	hints := ctx.Hints()
	_ = layout.NewLayoutBridge(hints, layout.BridgeConfig{PaddingPct: 0})
	return style.ViewData{}
}

// ─── ColorSlowStyle (parameterized skeleton) ──────────────────────────────────

// ColorSlowStyle is a skeleton style for color slow-refresh (color e-paper) panels.
// Capability: ColorSlow (RGB, slow refresh).
type ColorSlowStyle struct {
	Width  int
	Height int
}

func (s ColorSlowStyle) Name() string { return fmt.Sprintf("color-slow-%dx%d", s.Width, s.Height) }

func (s ColorSlowStyle) Requirements() style.SurfaceRequirements {
	return style.SurfaceRequirements{
		MinWidth:   s.Width,
		MinHeight:  s.Height,
		Capability: style.ColorSlow,
	}
}

func (s ColorSlowStyle) Supports(hints textlayout.TextHints) style.Fitness {
	return style.EvaluateFitness(s.Requirements(), hints)
}

func (s ColorSlowStyle) Build(_ source.MatrixSnapshot, _ source.Policy, ctx style.StyleContext) style.ViewData {
	hints := ctx.Hints()
	_ = layout.NewLayoutBridge(hints, layout.BridgeConfig{PaddingPct: 0})
	return style.ViewData{}
}

// ─── Shared Rendering ─────────────────────────────────────────────────────────

// buildFrame is the shared rendering function for all matrix styles.
// It produces a single frame of matrix rain output. Currently returns
// a minimal ViewData; the full column-drop rendering pipeline is not yet implemented.
func buildFrame(snapshot source.MatrixSnapshot) style.ViewData {
	return style.ViewData{
		Static: snapshot.Eink,
	}
}
