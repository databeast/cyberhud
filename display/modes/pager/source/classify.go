package source

import "github.com/databeast/cyberhud/display/surface/textlayout"

// SurfaceClass is the pager's runtime classification of a panel, which decides
// the render strategy BuildView advances. It is mode logic, not style placement:
// a style never classifies a panel, it is handed a frame already resolved for
// one.
type SurfaceClass int

const (
	// SurfaceFast is a fast-refresh panel (OLED, TFT) that supports smooth
	// pixel-level scrolling.
	SurfaceFast SurfaceClass = iota

	// SurfaceSlow is a slow-refresh panel (e-ink) that requires full-page
	// transitions rather than continuous scrolling.
	SurfaceSlow
)

// ClassifySurface decides whether a panel is fast or slow.
//
//   - Capability ∈ {MonoSlow, GrayscaleSlow, ColorSlow} → SurfaceSlow
//   - PreferEventRefresh set → SurfaceSlow (regardless of capability)
//   - otherwise → SurfaceFast
//   - zero-value / unavailable hints → SurfaceFast
func ClassifySurface(hints textlayout.TextHints) SurfaceClass {
	switch hints.Capability {
	case textlayout.CapMonoSlow, textlayout.CapGrayscaleSlow, textlayout.CapColorSlow:
		return SurfaceSlow
	}
	if hints.PreferEventRefresh {
		return SurfaceSlow
	}
	return SurfaceFast
}

// IsSlow reports whether the class is the slow page-transition strategy.
func (c SurfaceClass) IsSlow() bool { return c == SurfaceSlow }
