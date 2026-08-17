package progressbar

import (
	"image"
	"image/color"

	"github.com/databeast/cyberhud/display/widgets"
)

// renderSegmented draws a segmented (chunked) bar along the primary axis.
// Cells are arranged and filled according to geom.FillDirection:
//   - FillLeftToRight: cells left-to-right, fill progresses left-to-right
//   - FillBottomToTop: cells bottom-to-top, fill progresses bottom-to-top
//
// If the bar is too short for at least 2 cells + 1 gap, the function falls
// back to an unsegmented solid fill (identical to Linear style rendering).
func renderSegmented(cfg Config, geom RenderGeometry) *image.RGBA {
	fg, bg := widgets.ResolveColors(cfg.Foreground, cfg.Background)

	primaryAxis := geom.PrimaryAxis
	minorAxis := geom.MinorAxis
	gap := cfg.SegmentGap

	// Compute segment count.
	segmentCount := cfg.SegmentCount
	if segmentCount <= 0 {
		segmentCount = primaryAxis / (4 + gap)
	}
	if segmentCount < 1 {
		segmentCount = 1
	}

	// Compute cell width: distribute available space (minus gaps) evenly.
	cellWidth := (primaryAxis - (segmentCount-1)*gap) / segmentCount
	if cellWidth < 1 {
		cellWidth = 1
	}

	// Too-short bar fallback: if primary axis can't fit 2 cells + 1 gap,
	// render as an unsegmented linear fill.
	if primaryAxis < 2*cellWidth+gap {
		return renderSegmentedFallback(cfg, geom, fg, bg)
	}

	// Fill extent: how many pixels along primary axis are considered "filled".
	fillExtent := int(float64(primaryAxis) * cfg.Value)

	// Create the output image.
	var imgW, imgH int
	if geom.FillDirection == FillLeftToRight {
		imgW = cfg.Bounds.Dx()
		imgH = cfg.Bounds.Dy()
	} else {
		// Vertical: primary axis = height, minor axis = width
		imgW = cfg.Bounds.Dx()
		imgH = cfg.Bounds.Dy()
	}
	img := image.NewRGBA(image.Rect(0, 0, imgW, imgH))

	// Fill entire image with background first.
	for y := 0; y < imgH; y++ {
		for x := 0; x < imgW; x++ {
			img.SetRGBA(x, y, bg)
		}
	}

	// Determine whether to use gradient coloring.
	useGradient := cfg.Gradient != nil && len(cfg.Gradient.Stops) >= 2

	// Draw cells along the primary axis.
	for i := 0; i < segmentCount; i++ {
		cellStart := i * (cellWidth + gap)
		cellEnd := cellStart + cellWidth
		cellCenter := cellStart + cellWidth/2

		// Only fill if cell center falls within the fill region.
		if cellCenter >= fillExtent {
			continue // Leave as background (already set).
		}

		// Determine the fill color for this cell.
		cellColor := fg
		if useGradient {
			// Compute normalized position t along the primary axis at cell center.
			var t float64
			if primaryAxis > 1 {
				t = float64(cellCenter) / float64(primaryAxis-1)
			}
			cellColor = interpolateGradient(cfg.Gradient.Stops, t)
		}

		// Fill this cell with the computed color.
		fillCell(img, geom, cellStart, cellEnd, minorAxis, cellColor)
	}

	return img
}

// fillCell fills a single segment cell in the image between cellStart and
// cellEnd along the primary axis, spanning the full minor axis.
func fillCell(img *image.RGBA, geom RenderGeometry, cellStart, cellEnd, minorAxis int, c color.RGBA) {
	if geom.FillDirection == FillLeftToRight {
		// Horizontal: primary axis = X, minor axis = Y
		for y := 0; y < minorAxis; y++ {
			for x := cellStart; x < cellEnd; x++ {
				img.SetRGBA(x, y, c)
			}
		}
	} else {
		// Vertical: primary axis = Y (bottom-to-top), minor axis = X
		// Position 0 along primary axis maps to pixel y = height-1.
		height := img.Bounds().Dy()
		for x := 0; x < minorAxis; x++ {
			for pos := cellStart; pos < cellEnd; pos++ {
				y := height - 1 - pos
				if y >= 0 && y < height {
					img.SetRGBA(x, y, c)
				}
			}
		}
	}
}

// renderSegmentedFallback renders the bar as a simple unsegmented solid fill,
// identical to a Linear style bar with the same orientation and value.
func renderSegmentedFallback(cfg Config, geom RenderGeometry, fg, bg color.RGBA) *image.RGBA {
	if geom.FillDirection == FillBottomToTop {
		return renderVertical(cfg.Bounds, cfg.Value, fg, bg)
	}
	return renderHorizontal(cfg.Bounds, cfg.Value, fg, bg)
}
