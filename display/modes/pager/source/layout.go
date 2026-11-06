package source

import (
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
)

// pagerFontTier is the tier the pager measures and renders text with. The pager
// is a text tailer: one uniform tier across the whole frame, sized by the
// panel's catalog rather than per-row.
const pagerFontTier = tiercatalog.TierNormal

const PagerPaddingPct = 2

// ComputeLayout resolves the visible row/column budget and text metrics for a
// panel. This is mode logic: BuildView uses it to size the visible window and to
// derive the scroll/page cadence before handing a resolved frame to a style.
//
// Zero PixelWidth or PixelHeight yields a zero-row/zero-column layout.
func ComputeLayout(ctx style.StyleContext) Layout {
	hints := ctx.Hints()
	pw := hints.PixelWidth
	ph := hints.PixelHeight

	if pw <= 0 || ph <= 0 {
		return Layout{PixelWidth: pw, PixelHeight: ph}
	}

	entry := ctx.Entry(pagerFontTier)
	glyphAdvance := entry.GlyphAdvance
	rowHeight := entry.RowHeight

	visibleColumns := 0
	if glyphAdvance > 0 {
		visibleColumns = pw / glyphAdvance
	}
	visibleRows := 0
	if rowHeight > 0 {
		visibleRows = ph / rowHeight
	}

	return Layout{
		VisibleColumns: visibleColumns,
		VisibleRows:    visibleRows,
		GlyphAdvance:   glyphAdvance,
		RowHeight:      rowHeight,
		PixelWidth:     pw,
		PixelHeight:    ph,
	}
}
