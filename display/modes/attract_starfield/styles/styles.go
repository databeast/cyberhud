package styles

import (
	"fmt"

	"github.com/databeast/cyberhud/display/modes/attract_starfield/source"

	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/style/layout"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// ─── ColorFast ────────────────────────────────────────────────────────────────

// ColorFastStyle targets color TFT panels.
type ColorFastStyle struct {
	Width  int
	Height int
}

func (s ColorFastStyle) Name() string { return fmt.Sprintf("color-fast-%dx%d", s.Width, s.Height) }

func (s ColorFastStyle) Requirements() style.SurfaceRequirements {
	return style.SurfaceRequirements{
		MinWidth:        s.Width,
		MinHeight:       s.Height,
		Capability:      style.ColorFast,
		MinRows:         0,
		MinCharsPerLine: 0,
	}
}

func (s ColorFastStyle) Supports(hints textlayout.TextHints) style.Fitness {
	return style.EvaluateFitness(s.Requirements(), hints)
}

func (s ColorFastStyle) Build(snapshot source.Snapshot, pol source.Policy, ctx style.StyleContext) style.ViewData {
	hints := ctx.Hints()
	_ = layout.NewLayoutBridge(hints, layout.BridgeConfig{PaddingPct: 0})
	return style.ViewData{Static: snapshot.IsEink}
}

// ─── MonoFast ─────────────────────────────────────────────────────────────────

// MonoFastStyle targets monochrome OLED panels.
type MonoFastStyle struct {
	Width  int
	Height int
}

func (s MonoFastStyle) Name() string { return fmt.Sprintf("mono-fast-%dx%d", s.Width, s.Height) }

func (s MonoFastStyle) Requirements() style.SurfaceRequirements {
	return style.SurfaceRequirements{
		MinWidth:        s.Width,
		MinHeight:       s.Height,
		Capability:      style.MonoFast,
		MinRows:         0,
		MinCharsPerLine: 0,
	}
}

func (s MonoFastStyle) Supports(hints textlayout.TextHints) style.Fitness {
	return style.EvaluateFitness(s.Requirements(), hints)
}

func (s MonoFastStyle) Build(snapshot source.Snapshot, pol source.Policy, ctx style.StyleContext) style.ViewData {
	hints := ctx.Hints()
	_ = layout.NewLayoutBridge(hints, layout.BridgeConfig{PaddingPct: 0})
	return style.ViewData{Static: snapshot.IsEink}
}

// ─── MonoSlow (e-ink) ─────────────────────────────────────────────────────────

// MonoSlowStyle targets mono e-paper panels.
type MonoSlowStyle struct {
	Width  int
	Height int
}

func (s MonoSlowStyle) Name() string { return fmt.Sprintf("mono-slow-%dx%d", s.Width, s.Height) }

func (s MonoSlowStyle) Requirements() style.SurfaceRequirements {
	return style.SurfaceRequirements{
		MinWidth:        s.Width,
		MinHeight:       s.Height,
		Capability:      style.MonoSlow,
		MinRows:         0,
		MinCharsPerLine: 0,
	}
}

func (s MonoSlowStyle) Supports(hints textlayout.TextHints) style.Fitness {
	return style.EvaluateFitness(s.Requirements(), hints)
}

func (s MonoSlowStyle) Build(snapshot source.Snapshot, pol source.Policy, ctx style.StyleContext) style.ViewData {
	hints := ctx.Hints()
	_ = layout.NewLayoutBridge(hints, layout.BridgeConfig{PaddingPct: 0})
	return style.ViewData{Static: true}
}

// ─── GrayscaleSlow ────────────────────────────────────────────────────────────

// GrayscaleSlowStyle targets grayscale e-paper panels.
type GrayscaleSlowStyle struct {
	Width  int
	Height int
}

func (s GrayscaleSlowStyle) Name() string {
	return fmt.Sprintf("grayscale-slow-%dx%d", s.Width, s.Height)
}

func (s GrayscaleSlowStyle) Requirements() style.SurfaceRequirements {
	return style.SurfaceRequirements{
		MinWidth:        s.Width,
		MinHeight:       s.Height,
		Capability:      style.GrayscaleSlow,
		MinRows:         0,
		MinCharsPerLine: 0,
	}
}

func (s GrayscaleSlowStyle) Supports(hints textlayout.TextHints) style.Fitness {
	return style.EvaluateFitness(s.Requirements(), hints)
}

func (s GrayscaleSlowStyle) Build(snapshot source.Snapshot, pol source.Policy, ctx style.StyleContext) style.ViewData {
	hints := ctx.Hints()
	_ = layout.NewLayoutBridge(hints, layout.BridgeConfig{PaddingPct: 0})
	return style.ViewData{Static: true}
}

// ─── GrayscaleFast ────────────────────────────────────────────────────────────

// GrayscaleFastStyle targets grayscale fast-refresh panels.
type GrayscaleFastStyle struct {
	Width  int
	Height int
}

func (s GrayscaleFastStyle) Name() string {
	return fmt.Sprintf("grayscale-fast-%dx%d", s.Width, s.Height)
}

func (s GrayscaleFastStyle) Requirements() style.SurfaceRequirements {
	return style.SurfaceRequirements{
		MinWidth:        s.Width,
		MinHeight:       s.Height,
		Capability:      style.GrayscaleFast,
		MinRows:         0,
		MinCharsPerLine: 0,
	}
}

func (s GrayscaleFastStyle) Supports(hints textlayout.TextHints) style.Fitness {
	return style.EvaluateFitness(s.Requirements(), hints)
}

func (s GrayscaleFastStyle) Build(snapshot source.Snapshot, pol source.Policy, ctx style.StyleContext) style.ViewData {
	hints := ctx.Hints()
	_ = layout.NewLayoutBridge(hints, layout.BridgeConfig{PaddingPct: 0})
	return style.ViewData{Static: snapshot.IsEink}
}

// ─── ColorSlow ────────────────────────────────────────────────────────────────

// ColorSlowStyle targets color e-paper panels.
type ColorSlowStyle struct {
	Width  int
	Height int
}

func (s ColorSlowStyle) Name() string { return fmt.Sprintf("color-slow-%dx%d", s.Width, s.Height) }

func (s ColorSlowStyle) Requirements() style.SurfaceRequirements {
	return style.SurfaceRequirements{
		MinWidth:        s.Width,
		MinHeight:       s.Height,
		Capability:      style.ColorSlow,
		MinRows:         0,
		MinCharsPerLine: 0,
	}
}

func (s ColorSlowStyle) Supports(hints textlayout.TextHints) style.Fitness {
	return style.EvaluateFitness(s.Requirements(), hints)
}

func (s ColorSlowStyle) Build(snapshot source.Snapshot, pol source.Policy, ctx style.StyleContext) style.ViewData {
	hints := ctx.Hints()
	_ = layout.NewLayoutBridge(hints, layout.BridgeConfig{PaddingPct: 0})
	return style.ViewData{Static: true}
}
