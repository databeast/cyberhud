// Package layout provides the Style Layout Calculator, a pure-computation utility
// that centralizes text-row-to-pixel coordinate translation for the display
// system's style implementations.
//
// Styles need to position sprites relative to text rows (e.g., "3px below the
// last text row"). Without the bridge, each style duplicates the renderer's
// layout math (rowHeight, content origin, title bar offset, padding inset),
// leading to subtle bugs when one copy diverges from another.
//
// The bridge is constructed from [textlayout.TextHints] and a [BridgeConfig],
// performs defensive clamping at construction time, and exposes query methods
// for pixel positioning. All arithmetic is integer-only with no side effects,
// making the bridge safe for concurrent use without synchronization.
package layout

import (
	"image"

	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// BridgeConfig holds the inset configuration
// needed to construct a LayoutCalculator.
type BridgeConfig struct {
	PaddingPct int // Percentage inset 0–50 of panel dimensions per edge
}

// LayoutCalculator encapsulates all layout geometry derived from TextHints and
// BridgeConfig. It is an immutable value type: all fields are set at construction
// and never mutated.
type LayoutCalculator struct {
	contentOriginX         int
	contentOriginY         int
	rowHeight              int
	glyphAdvance           int
	pixelWidth             int
	pixelHeight            int
	availableContentHeight int
	availableContentWidth  int
}

// AvailableContentHeight returns the number of vertical pixels available for
// content rows (between title bar and hint bar, minus padding insets).
func (lb LayoutCalculator) AvailableContentHeight() int {
	return lb.availableContentHeight
}

// AvailableContentWidth returns the horizontal pixels available for content,
// which is PixelWidth minus twice the padding inset, clamped to zero.
func (lb LayoutCalculator) AvailableContentWidth() int {
	return lb.availableContentWidth
}

// RowY returns the Y pixel coordinate where text row n begins, plus an
// optional vertical pixel offset. Negative row indices are clamped to zero.
// The result is not clamped to the panel height.
func (lb LayoutCalculator) RowY(n, offsetPx int) int {
	if n < 0 {
		n = 0
	}
	return lb.contentOriginY + n*lb.rowHeight + offsetPx
}

// RowHeight returns the effective row height in pixels.
func (lb LayoutCalculator) RowHeight() int {
	return lb.rowHeight
}

// GlyphAdvance returns the effective glyph advance in pixels.
func (lb LayoutCalculator) GlyphAdvance() int {
	return lb.glyphAdvance
}

func (lb LayoutCalculator) MaxChars() int {
	if lb.glyphAdvance > 0 {
		return lb.AvailableContentWidth() / lb.GlyphAdvance()
	}
	return 0
}

// RowBottomY returns the Y pixel coordinate of the bottom edge of text row n.
// Negative row indices are clamped to zero.
func (lb LayoutCalculator) RowBottomY(n int) int {
	if n < 0 {
		n = 0
	}
	return lb.contentOriginY + (n+1)*lb.rowHeight
}

// RowX returns the X pixel coordinate for a given character offset within a
// text row. The base is always contentOriginX (the padding inset). Negative
// char offsets are treated as zero. The result is not clamped to the panel width.
func (lb LayoutCalculator) RowX(charOffset int) int {
	if charOffset < 0 {
		charOffset = 0
	}
	return lb.contentOriginX + charOffset*lb.glyphAdvance
}

// ContentOrigin returns the (X, Y) pixel coordinates of the content area's
// top-left corner.
func (lb LayoutCalculator) ContentOrigin() (int, int) {
	return lb.contentOriginX, lb.contentOriginY
}

// ContentCenter returns the center pixel coordinates (X, Y) of the available
// content area. When available width or height is zero, the center degenerates
// to the content origin coordinate for that axis.
func (lb LayoutCalculator) ContentCenter() (int, int) {
	w := lb.availableContentWidth
	if w < 0 {
		w = 0
	}
	h := lb.availableContentHeight
	if h < 0 {
		h = 0
	}
	centerX := lb.contentOriginX + w/2
	centerY := lb.contentOriginY + h/2
	return centerX, centerY
}

// RemainingSpace returns how many vertical pixels remain after n text rows
// have been allocated. Negative n is treated as zero. The result is never
// negative.
func (lb LayoutCalculator) RemainingSpace(n int) int {
	if n < 0 {
		n = 0
	}
	if lb.rowHeight == 0 {
		return 0
	}
	remaining := lb.availableContentHeight - n*lb.rowHeight
	if remaining < 0 {
		return 0
	}
	return remaining
}

// MaxVisibleRows returns the number of text rows that fit in the available
// content area. Returns zero if row height is zero or there is no available
// content height.
func (lb LayoutCalculator) MaxVisibleRows() int {
	if lb.rowHeight == 0 || lb.availableContentHeight == 0 {
		return 0
	}
	return lb.availableContentHeight / lb.rowHeight
}

// ItemsBottomY returns the Y pixel coordinate immediately below the last of
// itemCount text rows. Negative item counts are treated as zero.
func (lb LayoutCalculator) ItemsBottomY(itemCount int) int {
	if itemCount < 0 {
		itemCount = 0
	}
	return lb.contentOriginY + itemCount*lb.rowHeight
}

// TextPixelWidth returns GlyphAdvance × charCount using the bridge's configured advance.
// Returns 0 when charCount <= 0.
func (lb LayoutCalculator) TextPixelWidth(charCount int) int {
	if charCount <= 0 {
		return 0
	}
	return lb.glyphAdvance * charCount
}

// TextPixelWidthWith returns glyphAdvance × charCount using the provided advance.
// Falls back to the bridge's configured advance when glyphAdvance <= 0.
// Returns 0 when charCount <= 0.
func (lb LayoutCalculator) TextPixelWidthWith(charCount, glyphAdvance int) int {
	if charCount <= 0 {
		return 0
	}
	if glyphAdvance <= 0 {
		glyphAdvance = lb.glyphAdvance
	}
	return glyphAdvance * charCount
}

// CenterX returns the horizontal offset to center text of charCount characters.
// Result is max(0, (AvailableContentWidth - TextPixelWidth(charCount)) / 2).
func (lb LayoutCalculator) CenterX(charCount int) int {
	offset := (lb.availableContentWidth - lb.TextPixelWidth(charCount)) / 2
	if offset < 0 {
		return 0
	}
	return offset
}

// CenterXWith returns the horizontal centering offset using an explicit glyph advance.
// Falls back to bridge's configured advance when glyphAdvance <= 0.
func (lb LayoutCalculator) CenterXWith(charCount, glyphAdvance int) int {
	offset := (lb.availableContentWidth - lb.TextPixelWidthWith(charCount, glyphAdvance)) / 2
	if offset < 0 {
		return 0
	}
	return offset
}

// CenterBlockY returns the vertical offset to center a content block defined by
// rowHeights and inter-row spacing within the available content height.
// Returns 0 for empty slices or when the block exceeds available height.
func (lb LayoutCalculator) CenterBlockY(rowHeights []int, spacing int) int {
	if len(rowHeights) == 0 {
		return 0
	}
	if spacing < 0 {
		spacing = 0
	}

	blockHeight := 0
	for _, h := range rowHeights {
		blockHeight += h
	}
	blockHeight += spacing * (len(rowHeights) - 1)

	offset := (lb.availableContentHeight - blockHeight) / 2
	if offset < 0 {
		return 0
	}
	return offset
}

// FitRows determines how many rows fit, their spacing, and the vertical centering
// offset. Returns (spacing, offsetY, visibleCount).
// Implements the adaptive row fitting algorithm from clock.ComputeLayout.
//
// Row omission proceeds from the end of the slice until the remaining rows fit
// within the available content height using a minimum inter-row spacing of 1px.
// At least one row is always preserved.
//
// For panels with ACH ≤ 135px, inter-row spacing is clamped to [1, 2].
// For panels with ACH > 135px, spacing is min(8, remainingSpace/(visibleCount+1)).
// When only one row is visible, spacing is 0.
// The content block is vertically centered in the available space.
func (lb LayoutCalculator) FitRows(rowHeights []int) (spacing int, offsetY int, visibleCount int) {
	if len(rowHeights) == 0 {
		return 0, 0, 0
	}

	ach := lb.availableContentHeight

	// Determine visible rows (drop from end if they don't fit with min 1px spacing).
	visibleCount = len(rowHeights)
	for visibleCount > 1 {
		totalH := sumRowHeights(rowHeights[:visibleCount]) + 1*(visibleCount-1)
		if totalH <= ach {
			break
		}
		visibleCount--
	}
	if visibleCount < 1 {
		visibleCount = 1
	}

	totalRowH := sumRowHeights(rowHeights[:visibleCount])
	remainingSpace := ach - totalRowH
	if remainingSpace < 0 {
		remainingSpace = 0
	}

	// Calculate spacing.
	if visibleCount <= 1 {
		spacing = 0
	} else {
		divisor := visibleCount + 1
		if ach <= 135 {
			calculated := remainingSpace / divisor
			if calculated < 1 {
				calculated = 1
			}
			if calculated > 2 {
				calculated = 2
			}
			spacing = calculated
		} else {
			calculated := remainingSpace / divisor
			if calculated > 8 {
				calculated = 8
			}
			spacing = calculated
		}
	}

	// Vertical centering.
	contentBlockH := totalRowH + spacing*(visibleCount-1)
	offsetY = (ach - contentBlockH) / 2
	if offsetY < 0 {
		offsetY = 0
	}

	return spacing, offsetY, visibleCount
}

// sumRowHeights returns the sum of all values in heights.
func sumRowHeights(heights []int) int {
	total := 0
	for _, h := range heights {
		total += h
	}
	return total
}

// BottomAnchorY returns the Y coordinate for an element anchored at the bottom
// of the content area. Returns ContentOrigin.Y when element overflows.
func (lb LayoutCalculator) BottomAnchorY(elementHeight int) int {
	if elementHeight <= 0 {
		return lb.contentOriginY + lb.availableContentHeight
	}
	if elementHeight > lb.availableContentHeight {
		return lb.contentOriginY
	}
	return lb.contentOriginY + lb.availableContentHeight - elementHeight
}

// TopRightAnchorX returns the X coordinate for an element anchored at the right
// edge of the content area. Returns ContentOrigin.X when element overflows.
func (lb LayoutCalculator) TopRightAnchorX(elementWidth int) int {
	if elementWidth <= 0 {
		return lb.contentOriginX + lb.availableContentWidth
	}
	if elementWidth > lb.availableContentWidth {
		return lb.contentOriginX
	}
	return lb.contentOriginX + lb.availableContentWidth - elementWidth
}

// NewLayoutBridge constructs a LayoutCalculator from the given TextHints and
// BridgeConfig. It performs all defensive clamping and derivation at
// construction time so that query methods can be trivial arithmetic.
func NewLayoutBridge(hints textlayout.TextHints, cfg BridgeConfig) LayoutCalculator {
	lb := LayoutCalculator{
		pixelWidth:  hints.PixelWidth,
		pixelHeight: hints.PixelHeight,
	}

	// Derive rowHeight.
	if hints.RowHeight > 0 {
		lb.rowHeight = hints.RowHeight
	} else {
		lb.rowHeight = textlayout.RowHeight // 10
	}

	// Derive glyphAdvance.
	if hints.GlyphAdvance > 0 {
		lb.glyphAdvance = hints.GlyphAdvance
	} else {
		lb.glyphAdvance = textlayout.GlyphAdvance // 6
	}

	// Clamp PaddingPct to [0, 50].
	paddingPct := cfg.PaddingPct
	if paddingPct < 0 {
		paddingPct = 0
	}
	if paddingPct > 50 {
		paddingPct = 50
	}

	// Compute padding insets (sole inset mechanism).
	padInsetX := paddingPct * hints.PixelWidth / 100
	padInsetY := paddingPct * hints.PixelHeight / 100

	// Content origin.
	lb.contentOriginX = padInsetX
	lb.contentOriginY = padInsetY

	// Available content dimensions.
	lb.availableContentWidth = hints.PixelWidth - 2*padInsetX
	if lb.availableContentWidth < 0 {
		lb.availableContentWidth = 0
	}

	bottomChrome := padInsetY
	lb.availableContentHeight = hints.PixelHeight - lb.contentOriginY - bottomChrome
	if lb.availableContentHeight < 0 {
		lb.availableContentHeight = 0
	}

	return lb
}

// deriveBarHeight computes the pixel height of a title or hint bar based on
// visibility, the panel's pixel height (which determines the base size), and
// the glyph height (which may override the base when glyphs are tall).
func deriveBarHeight(show bool, glyphHeight, pixelHeight int) int {
	if !show {
		return 0
	}
	base := 20
	if pixelHeight < 80 {
		base = 10
	}
	if glyphHeight > 0 {
		floor := glyphHeight + 4
		if floor > base {
			return floor
		}
	}
	return base
}

// RowInput describes a single text row's geometry for multi-tier layout computation.
type RowInput struct {
	TextLen      int // Number of characters in the row.
	GlyphAdvance int // Glyph advance for this row's font tier; ≤0 uses bridge default.
	RowHeight    int // Row height for this row's font tier; ≤0 uses bridge default.
}

// RowsResult holds the computed layout geometry for a set of tiered rows.
//
// This is the complete solution to the row layout problem for one frame. Every
// field must reach the renderer intact. Anything the renderer re-derives instead
// of reading from here will eventually disagree with OffsetY, because OffsetY was
// computed from these exact values — see the VisibleCount doc below for the
// concrete failure that motivated carrying the full result.
type RowsResult struct {
	Offsets    []int // Per-row horizontal centering offsets.
	RowHeights []int // Per-row pixel heights (normalized; never ≤0).
	OffsetY    int   // Vertical centering offset from FitRows.
	Spacing    int   // Inter-row spacing from FitRows.

	// VisibleCount is how many rows actually fit the region, as decided by
	// [LayoutCalculator.FitRows]. It may be fewer than len(Offsets): FitRows drops
	// rows from the end until the remainder fits.
	//
	// This field previously did not exist, and LayoutRows discarded FitRows'
	// third return value with `spacing, offsetY, _ := lb.FitRows(rowHeights)`.
	// The consequence was a silent inconsistency: on a region where only 2 of 3
	// rows fit, OffsetY was centred for a 2-row block while the caller still
	// emitted 3 items, and the renderer drew as many as its own unrelated
	// heuristic allowed. The block was then neither centred nor correctly clipped.
	//
	// Renderers must truncate to VisibleCount rather than computing their own row
	// budget. The style has already solved the fit against the real per-row font
	// metrics; a renderer dividing region height by a single row height cannot
	// reproduce that answer when rows use different fonts.
	VisibleCount int
}

// LayoutRows computes horizontal centering offsets and row heights for rows
// that may each have different font metrics (multi-tier layout). It normalizes
// zero/negative glyph advance and row height values to the bridge defaults,
// centers each row individually, and delegates vertical fitting to FitRows.
func (lb LayoutCalculator) LayoutRows(rows []RowInput) RowsResult {
	offsets := make([]int, len(rows))
	rowHeights := make([]int, len(rows))

	for i, r := range rows {
		ga := r.GlyphAdvance
		if ga <= 0 {
			ga = lb.glyphAdvance
		}
		if ga <= 0 {
			ga = 1
		}

		rh := r.RowHeight
		if rh <= 0 {
			rh = lb.rowHeight
		}
		if rh <= 0 {
			rh = 1
		}

		offset := (lb.availableContentWidth - r.TextLen*ga) / 2
		if offset < 0 {
			offset = 0
		}
		offsets[i] = offset
		rowHeights[i] = rh
	}

	spacing, offsetY, visibleCount := lb.FitRows(rowHeights)
	return RowsResult{
		Offsets:      offsets,
		RowHeights:   rowHeights,
		OffsetY:      offsetY,
		Spacing:      spacing,
		VisibleCount: visibleCount,
	}
}

// InlineWidgetBounds returns the pixel rectangle for a widget placed to the right
// of text at the given character offset on a specific row. Uses the bridge's
// configured glyph advance for the Min.X computation.
// Negative rowIndex and charOffset are clamped to 0.
// When Min.X > Max.X, returns an empty rectangle.
func (lb LayoutCalculator) InlineWidgetBounds(rowIndex, charOffset, rowHeight int) image.Rectangle {
	return lb.InlineWidgetBoundsWith(rowIndex, charOffset, rowHeight, 0)
}

// InlineWidgetBoundsWith returns the pixel rectangle for a widget placed to the
// right of text at the given character offset on a specific row, using an explicit
// glyph advance. Falls back to the bridge's configured advance when glyphAdvance <= 0.
// Negative rowIndex and charOffset are clamped to 0.
// When Min.X > Max.X, returns an empty rectangle.
func (lb LayoutCalculator) InlineWidgetBoundsWith(rowIndex, charOffset, rowHeight, glyphAdvance int) image.Rectangle {
	if rowIndex < 0 {
		rowIndex = 0
	}
	if charOffset < 0 {
		charOffset = 0
	}

	effectiveGlyphAdvance := glyphAdvance
	if effectiveGlyphAdvance <= 0 {
		effectiveGlyphAdvance = lb.glyphAdvance
	}

	minX := lb.contentOriginX + charOffset*effectiveGlyphAdvance
	minY := lb.contentOriginY + rowIndex*rowHeight
	maxX := lb.contentOriginX + lb.availableContentWidth
	maxY := minY + rowHeight

	if minX > maxX {
		return image.Rectangle{}
	}

	return image.Rectangle{
		Min: image.Point{X: minX, Y: minY},
		Max: image.Point{X: maxX, Y: maxY},
	}
}
