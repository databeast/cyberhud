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

// TTYX and TTYY are pixel coordinates measured within the TTY text-block's own
// frame — the laid-out, centered, padding-inset block of text rows a style hands
// the renderer. They are NOT panel-absolute (world) pixels: the block is a
// subregion floating inside the region, so a TTY coordinate only becomes a world
// pixel after its origin offset is added.
//
// World-relative pixels stay plain int / image.Point. Object-relative (TTY) pixels
// are these distinct types so a TTY value cannot be silently assigned into a
// world-space field (e.g. widgets.Sprite.Position) — that assignment is the
// recurring footgun. Cross the boundary deliberately via LayoutCalculator.TTYToPixel.
type TTYX int
type TTYY int

// BridgeConfig holds the inset configuration
// needed to construct a LayoutCalculator.
type BridgeConfig struct {
	PaddingPct int // Percentage inset 0–50 of panel dimensions per edge
}

// LayoutCalculator encapsulates all layout geometry derived from TextHints and
// BridgeConfig. It is an immutable value type: all fields are set at construction
// and never mutated.
// Metrics is the glyph advance and row height of a SPECIFIC resolved font,
// bound together so a measuring call cannot silently use metrics from a
// different font than the renderer will draw with.
//
// The fields are unexported and there is no bare-int path into them: the only
// way to obtain a usable Metrics is MetricsFrom (from the tiercatalog.Entry the
// renderer resolves the same tier to) or DefaultMetrics (the deliberate,
// greppable "I am drawing in the region default face" deviation). A style that
// measures with ctx.Entry(tier).Metrics() and emits that same tier is therefore
// guaranteed that the advance it centred with is the advance drawn — the
// measure-with-one-font-draw-with-another bug class becomes unrepresentable.
//
// The zero Metrics is unusable (advance/rowHeight 0); measuring methods treat it
// as "no font" and return zero rather than falling back to a hidden default.
type Metrics struct {
	advance    int
	glyphWidth int // width of a single glyph's ink box; <= advance (advance adds inter-glyph spacing)
	rowHeight  int
	fontID     string // the face these metrics belong to; for diagnostics
}

// Advance returns the glyph advance in pixels (one glyph's horizontal step,
// including the inter-glyph gap).
func (m Metrics) Advance() int { return m.advance }

// GlyphWidth returns the ink-box width of a single glyph in pixels. It is at most
// Advance: the last glyph on a row occupies GlyphWidth, not a full Advance, which
// is why TextWidth is exact rather than advance*count.
func (m Metrics) GlyphWidth() int { return m.glyphWidth }

// RowHeight returns the row height in pixels.
func (m Metrics) RowHeight() int { return m.rowHeight }

// FontID returns the identifier of the face these metrics describe.
func (m Metrics) FontID() string { return m.fontID }

// Usable reports whether these metrics carry positive advance and row height.
func (m Metrics) Usable() bool { return m.advance > 0 && m.rowHeight > 0 }

// NewMetrics builds a Metrics from raw values. It is the low-level constructor
// used by the catalog/context bridges (MetricsFrom, DefaultMetrics); styles
// should obtain Metrics through StyleContext, not call this directly. Kept
// exported because the tiercatalog→layout binding lives in the style package,
// which cannot reach into layout's unexported fields.
//
// glyphWidth should be the face's single-glyph ink width; pass the same value as
// advance if a caller genuinely only has the advance, which makes TextWidth
// conservative (it will over-estimate the final glyph by the advance gap) rather
// than wrong.
func NewMetrics(advance, glyphWidth, rowHeight int, fontID string) Metrics {
	return Metrics{advance: advance, glyphWidth: glyphWidth, rowHeight: rowHeight, fontID: fontID}
}

// TextWidth returns the exact horizontal pixels the string needs to render fully
// in this font, by the same rule the surface and textlabel draw with: each glyph
// steps one Advance, except the last, which occupies only GlyphWidth. An empty
// string needs 0; an unusable Metrics (no advance) needs 0.
//
// This is the system's single answer to "how wide is this line?". Compare it to a
// content width (see Fits) to decide whether a row renders whole, needs a smaller
// font, or must be dropped — never truncate text to force a fit.
func (m Metrics) TextWidth(s string) int {
	if m.advance <= 0 {
		return 0
	}
	n := len([]rune(s))
	if n <= 0 {
		return 0
	}
	gw := m.glyphWidth
	if gw <= 0 {
		gw = m.advance
	}
	return (n-1)*m.advance + gw
}

// Fits reports whether the string renders fully within widthPx in this font.
// A non-positive widthPx fits only the empty string.
func (m Metrics) Fits(s string, widthPx int) bool {
	return m.TextWidth(s) <= widthPx
}

type LayoutCalculator struct {
	contentOriginX         int
	contentOriginY         int
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

// RowY returns the Y pixel coordinate where text row n begins, plus an optional
// vertical pixel offset, for rows of the given font's height. Negative row
// indices are clamped to zero. The result is not clamped to the panel height.
func (lb LayoutCalculator) RowY(n, offsetPx int, m Metrics) int {
	if n < 0 {
		n = 0
	}
	return lb.contentOriginY + n*m.rowHeight + offsetPx
}

// MaxChars returns how many glyphs of the given font fit across the content
// width. Returns 0 for an unusable Metrics.
func (lb LayoutCalculator) MaxChars(m Metrics) int {
	if m.advance > 0 {
		return lb.availableContentWidth / m.advance
	}
	return 0
}

// RowBottomY returns the Y pixel coordinate of the bottom edge of text row n for
// rows of the given font's height. Negative row indices are clamped to zero.
func (lb LayoutCalculator) RowBottomY(n int, m Metrics) int {
	if n < 0 {
		n = 0
	}
	return lb.contentOriginY + (n+1)*m.rowHeight
}

// RowX returns the X pixel coordinate for a given character offset within a text
// row, for the given font. The base is always contentOriginX (the padding
// inset). Negative char offsets are treated as zero. Not clamped to panel width.
func (lb LayoutCalculator) RowX(charOffset int, m Metrics) int {
	if charOffset < 0 {
		charOffset = 0
	}
	return lb.contentOriginX + charOffset*m.advance
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

// RemainingSpace returns how many vertical pixels remain after n text rows of
// the given font have been allocated. Negative n is treated as zero. The result
// is never negative.
func (lb LayoutCalculator) RemainingSpace(n int, m Metrics) int {
	if n < 0 {
		n = 0
	}
	if m.rowHeight == 0 {
		return 0
	}
	remaining := lb.availableContentHeight - n*m.rowHeight
	if remaining < 0 {
		return 0
	}
	return remaining
}

// MaxVisibleRows returns the number of text rows of the given font that fit in
// the available content area. Returns zero for an unusable Metrics or no height.
func (lb LayoutCalculator) MaxVisibleRows(m Metrics) int {
	if m.rowHeight == 0 || lb.availableContentHeight == 0 {
		return 0
	}
	return lb.availableContentHeight / m.rowHeight
}

// ItemsBottomY returns the Y pixel coordinate immediately below the last of
// itemCount text rows of the given font. Negative item counts are treated as zero.
func (lb LayoutCalculator) ItemsBottomY(itemCount int, m Metrics) int {
	if itemCount < 0 {
		itemCount = 0
	}
	return lb.contentOriginY + itemCount*m.rowHeight
}

// TextPixelWidth returns the pixel width of charCount glyphs of the given font.
// Returns 0 when charCount <= 0 or the Metrics is unusable.
func (lb LayoutCalculator) TextPixelWidth(charCount int, m Metrics) int {
	if charCount <= 0 {
		return 0
	}
	return m.advance * charCount
}

// CenterX returns the horizontal offset to center charCount glyphs of the given
// font within the content width. Result is max(0, (ACW - advance*charCount)/2).
// Because m comes only from the font the renderer draws with, the centring is
// guaranteed to match the drawn text.
func (lb LayoutCalculator) CenterX(charCount int, m Metrics) TTYX {
	offset := (lb.availableContentWidth - lb.TextPixelWidth(charCount, m)) / 2
	if offset < 0 {
		return 0
	}
	return TTYX(offset)
}

// CenterBlockY returns the vertical offset to center a content block defined by
// rowHeights and inter-row spacing within the available content height.
// Returns 0 for empty slices or when the block exceeds available height.
func (lb LayoutCalculator) CenterBlockY(rowHeights []int, spacing int) TTYY {
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
	return TTYY(offset)
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

	// The calculator no longer holds a glyph advance or row height. Text
	// measurement requires a layout.Metrics (from the font the renderer will
	// draw with); the calculator provides only pixel geometry (content area,
	// padding origin) that is font-independent. This is deliberate: a baked-in
	// default advance was the frictionless path to measuring with the wrong font.

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

// RowInput describes a single text row's geometry for multi-tier layout computation.
type RowInput struct {
	TextLen int     // Number of characters in the row.
	Font    Metrics // The font this row is measured and drawn with. Must be usable.
}

// RowsResult holds the computed layout geometry for a set of tiered rows.
//
// This is the complete solution to the row layout problem for one frame. Every
// field must reach the renderer intact. Anything the renderer re-derives instead
// of reading from here will eventually disagree with OffsetY, because OffsetY was
// computed from these exact values — see the VisibleCount doc below for the
// concrete failure that motivated carrying the full result.
type RowsResult struct {
	Offsets    []TTYX // Per-row horizontal centering offsets (TTY-block frame).
	RowHeights []int  // Per-row pixel heights (normalized; never ≤0). A magnitude, not a coordinate.
	OffsetY    TTYY   // Vertical centering offset from FitRows (TTY-block frame).
	Spacing    int    // Inter-row spacing from FitRows (a gap, frame-agnostic).

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
	offsets := make([]TTYX, len(rows))
	rowHeights := make([]int, len(rows))

	for i, r := range rows {
		ga := r.Font.advance
		if ga <= 0 {
			ga = 1
		}
		rh := r.Font.rowHeight
		if rh <= 0 {
			rh = 1
		}

		offset := (lb.availableContentWidth - r.TextLen*ga) / 2
		if offset < 0 {
			offset = 0
		}
		offsets[i] = TTYX(offset)
		rowHeights[i] = rh
	}

	spacing, offsetY, visibleCount := lb.FitRows(rowHeights)

	return RowsResult{
		Offsets:      offsets,
		RowHeights:   rowHeights,
		OffsetY:      TTYY(offsetY),
		Spacing:      spacing,
		VisibleCount: visibleCount,
	}
}

// TTYToPixel projects a coordinate from the TTY text-block frame into
// panel-absolute (world) pixel space by adding the block's content origin
// (padding inset). This is the sanctioned, deliberate crossing from
// object-relative to world-relative coordinates — e.g. to pin a widgets.Sprite
// beside a specific text row. The returned image.Point is world-space and may be
// assigned directly to Sprite.Position.
func (lb LayoutCalculator) TTYToPixel(x TTYX, y TTYY) image.Point {
	return image.Point{
		X: lb.contentOriginX + int(x),
		Y: lb.contentOriginY + int(y),
	}
}

// CellToPixel returns the world pixel of the top-left corner of the character
// cell at (row, col), measured with font m. Negative row/col clamp to 0.
// Not clamped to the panel.
func (lb LayoutCalculator) CellToPixel(row, col int, m Metrics) image.Point {
	return image.Point{X: lb.RowX(col, m), Y: lb.RowY(row, 0, m)}
}

// CellRect returns the world pixel rectangle covered by the character cell at
// (row, col), measured with font m. Negative row/col clamp to 0.
func (lb LayoutCalculator) CellRect(row, col int, m Metrics) image.Rectangle {
	min := lb.CellToPixel(row, col, m)
	return image.Rectangle{
		Min: min,
		Max: image.Point{X: min.X + m.advance, Y: min.Y + m.rowHeight},
	}
}

// PixelToCell is the inverse of CellToPixel: it maps a world pixel to the
// (row, col) of the character cell containing it. ok is false when m is
// unusable or p lies left of or above the content origin, or at or beyond the
// right edge of the content width.
func (lb LayoutCalculator) PixelToCell(p image.Point, m Metrics) (row, col int, ok bool) {
	if !m.Usable() {
		return 0, 0, false
	}
	dx, dy := p.X-lb.contentOriginX, p.Y-lb.contentOriginY
	if dx < 0 || dy < 0 || dx >= lb.availableContentWidth {
		return 0, 0, false
	}
	return dy / m.rowHeight, dx / m.advance, true
}

// InlineWidgetBounds returns the pixel rectangle for a widget placed to the
// right of text at the given character offset on a specific row, measured with
// the row's font m. Negative rowIndex and charOffset are clamped to 0.
// When Min.X > Max.X, returns an empty rectangle.
func (lb LayoutCalculator) InlineWidgetBounds(rowIndex, charOffset, rowHeight int, m Metrics) image.Rectangle {
	if rowIndex < 0 {
		rowIndex = 0
	}
	if charOffset < 0 {
		charOffset = 0
	}

	minX := lb.contentOriginX + charOffset*m.advance
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
