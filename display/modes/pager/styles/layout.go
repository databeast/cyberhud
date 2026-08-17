package styles

import (
	"github.com/databeast/cyberhud/display/modes/pager/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
)

// Font configuration for the pager mode.
const pagerFontFamily = "spleen"

var pagerFontTier = tiercatalog.TierNormal

type Layout = source.Layout

// ComputeLayout calculates the visible rows and columns for the pager mode.
// It resolves font metrics via ctx.Entry (infallible tier-catalog API).
//
// If PixelWidth or PixelHeight is zero, the result is zero rows and columns.
func ComputeLayout(ctx style.StyleContext) Layout {
	hints := ctx.Hints()
	pixelWidth := hints.PixelWidth
	pixelHeight := hints.PixelHeight

	// Zero-dimension edge case: report zero rows/columns.
	if pixelWidth <= 0 || pixelHeight <= 0 {
		return Layout{
			PixelWidth:  pixelWidth,
			PixelHeight: pixelHeight,
		}
	}

	entry := ctx.Entry(pagerFontTier)
	glyphAdvance := entry.GlyphAdvance
	rowHeight := entry.RowHeight

	// Compute visible columns and rows via integer floor division.
	visibleColumns := 0
	if glyphAdvance > 0 {
		visibleColumns = pixelWidth / glyphAdvance
	}

	visibleRows := 0
	if rowHeight > 0 {
		visibleRows = pixelHeight / rowHeight
	}

	return Layout{
		VisibleColumns: visibleColumns,
		VisibleRows:    visibleRows,
		GlyphAdvance:   glyphAdvance,
		RowHeight:      rowHeight,
		PixelWidth:     pixelWidth,
		PixelHeight:    pixelHeight,
	}
}
