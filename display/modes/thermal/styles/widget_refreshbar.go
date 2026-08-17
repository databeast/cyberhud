package styles

import (
	"image"
	"image/color"

	"github.com/databeast/cyberhud/display/modes/thermal/source"
	"github.com/databeast/cyberhud/display/widgets"
	"github.com/databeast/cyberhud/display/widgets/progressbar"
)

// buildRefreshBarSprite generates the horizontal progress bar sprite showing
// elapsed time since the last sample as a proportion of refresh_ms. The bar is
// positioned at the bottom of the effective content area, spanning the full
// Effective_Width with a height of 4 pixels.
//
// Returns nil when:
//   - pol.ShowRefreshBar is false
//   - effectiveHeight < 48 (panel too short for a useful bar)
//
// When active, the bar reduces the available content height by 4 pixels;
// callers must account for this before computing row visibility.
func buildRefreshBarSprite(pol source.Policy, elapsedMS int, effectiveWidth int, effectiveHeight int, contentOffsetX int, contentOffsetY int, isColor bool, nativeFG color.RGBA) *widgets.Sprite {
	// Suppress when policy disables the bar.
	if !pol.ShowRefreshBar {
		return nil
	}

	// Suppress when effective height is too small for a useful layout.
	if effectiveHeight < 48 {
		return nil
	}

	// Compute progress value, clamped to [0.0, 1.0].
	progress := 0.0
	if pol.RefreshMS > 0 {
		progress = float64(elapsedMS) / float64(pol.RefreshMS)
	}
	if progress < 0.0 {
		progress = 0.0
	}
	if progress > 1.0 {
		progress = 1.0
	}

	// Position bar at bottom of effective content area.
	const barHeight = 4
	barX := contentOffsetX
	barY := contentOffsetY + effectiveHeight - barHeight

	// Determine foreground and background colors based on panel type.
	var fg, bg color.RGBA
	if isColor {
		// Color panel: foreground is accent color.
		switch pol.FGColor {
		case "thermal":
			// Severity green for the "thermal" accent.
			fg = color.RGBA{R: 0, G: 255, B: 0, A: 255}
		case "none":
			// Opaque white for the "none" accent.
			fg = color.RGBA{R: 255, G: 255, B: 255, A: 255}
		default:
			// Named accent color.
			fg = accentColor(pol.FGColor)
		}
		// Background: dark gray.
		bg = color.RGBA{R: 40, G: 40, B: 40, A: 255}
	} else {
		// Monochrome panel: native foreground, transparent background.
		fg = nativeFG
		bg = color.RGBA{R: 0, G: 0, B: 0, A: 0}
	}

	// Render the progress bar widget.
	bounds := image.Rect(barX, barY, barX+effectiveWidth, barY+barHeight)
	result := progressbar.Render(progressbar.Config{
		Style:      progressbar.Linear,
		Value:      progress,
		Bounds:     bounds,
		Foreground: fg,
		Background: bg,
	})
	if result == nil {
		return nil
	}

	return result
}
