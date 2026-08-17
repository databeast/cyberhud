package progressbar

import (
	"image"
	"image/color"

	"github.com/databeast/cyberhud/display/widgets"
)

// renderLinear draws a rectangular progress bar for both horizontal and vertical
// orientations. It uses gradient fill when a valid GradientFill (≥2 stops) is
// provided, otherwise falls back to solid foreground color.
//
// For Horizontal orientation: fill progresses left-to-right.
// For Vertical orientation: fill progresses bottom-to-top.
func renderLinear(cfg Config, geom RenderGeometry) *image.RGBA {
	w := cfg.Bounds.Dx()
	h := cfg.Bounds.Dy()
	img := image.NewRGBA(image.Rect(0, 0, w, h))

	fg, bg := widgets.ResolveColors(cfg.Foreground, cfg.Background)

	// Determine if we have a valid gradient (≥2 stops).
	useGradient := cfg.Gradient != nil && len(cfg.Gradient.Stops) >= 2

	switch geom.FillDirection {
	case FillLeftToRight:
		fillCols := int(float64(w) * cfg.Value)

		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				if x < fillCols {
					if useGradient {
						// Map pixel x to gradient t within the fill region.
						var t float64
						if fillCols > 1 {
							t = float64(x) / float64(fillCols-1)
						}
						c := interpolateGradient(cfg.Gradient.Stops, t)
						img.SetRGBA(x, y, c)
					} else {
						img.SetRGBA(x, y, fg)
					}
				} else {
					img.SetRGBA(x, y, bg)
				}
			}
		}

	case FillBottomToTop:
		fillRows := int(float64(h) * cfg.Value)

		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				if y >= h-fillRows {
					if useGradient {
						// Map pixel y to gradient t within the fill region.
						// Bottom of fill region is y = h-1, top of fill region is y = h-fillRows.
						// t = 0.0 at bottom (y = h-1), t = 1.0 at top (y = h-fillRows).
						var t float64
						if fillRows > 1 {
							t = float64((h-1)-y) / float64(fillRows-1)
						}
						c := interpolateGradient(cfg.Gradient.Stops, t)
						img.SetRGBA(x, y, c)
					} else {
						img.SetRGBA(x, y, fg)
					}
				} else {
					img.SetRGBA(x, y, bg)
				}
			}
		}

	default:
		// Fallback: fill horizontal (shouldn't reach here for Linear style).
		fillCols := int(float64(w) * cfg.Value)
		fillHorizontalSolid(img, w, h, fillCols, fg, bg)
	}

	return img
}

// fillHorizontalSolid is a helper to fill a horizontal bar with solid colors.
func fillHorizontalSolid(img *image.RGBA, w, h, fillCols int, fg, bg color.RGBA) {
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if x < fillCols {
				img.SetRGBA(x, y, fg)
			} else {
				img.SetRGBA(x, y, bg)
			}
		}
	}
}
