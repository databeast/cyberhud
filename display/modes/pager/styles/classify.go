package styles

import "github.com/databeast/cyberhud/display/surface/textlayout"

// SurfaceClass represents the classification of a display surface for
// rendering strategy selection.
type SurfaceClass int

const (
	// SurfaceFast indicates a fast-refresh display (OLED, TFT) that supports
	// smooth pixel-level scrolling at high frame rates.
	SurfaceFast SurfaceClass = iota

	// SurfaceSlow indicates a slow-refresh display (e-ink) that requires
	// full-page transitions rather than continuous pixel scrolling.
	SurfaceSlow
)

// classifySurface determines whether the display surface described by hints
// is a fast or slow display. The classification drives rendering strategy
// selection: surfaceFast uses smooth scrolling, surfaceSlow uses page
// transitions.
//
// Classification rules:
//   - Capability ∈ {CapMonoSlow, CapGrayscaleSlow, CapColorSlow} → surfaceSlow
//   - PreferEventRefresh = true → surfaceSlow (regardless of Capability)
//   - Otherwise → surfaceFast
//   - Zero-value / unavailable TextHints → surfaceFast (default)
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
