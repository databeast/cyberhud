package layout_test

import (
	"image"
	"testing"

	"github.com/databeast/cyberhud/display/style/layout"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"pgregory.net/rapid"
)

func TestProperty3_ContentOriginComputation(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate random TextHints
		pixelWidth := rapid.IntRange(0, 500).Draw(t, "pixelWidth")
		pixelHeight := rapid.IntRange(0, 500).Draw(t, "pixelHeight")
		glyphHeight := rapid.IntRange(-10, 50).Draw(t, "glyphHeight")

		hints := textlayout.TextHints{
			PixelWidth:  pixelWidth,
			PixelHeight: pixelHeight,
			GlyphHeight: glyphHeight,
		}

		// Generate random BridgeConfig
		paddingPct := rapid.IntRange(0, 50).Draw(t, "paddingPct")

		cfg := layout.BridgeConfig{
			PaddingPct: paddingPct,
		}

		lb := layout.NewLayoutBridge(hints, cfg)
		gotX, gotY := lb.ContentOrigin()

		// Compute expected padding insets
		padInsetX := paddingPct * pixelWidth / 100
		padInsetY := paddingPct * pixelHeight / 100

		// Verify ContentOrigin().X == padInsetX
		if gotX != padInsetX {
			t.Fatalf("ContentOrigin().X = %d, want %d (paddingPct=%d, pixelWidth=%d)",
				gotX, padInsetX, paddingPct, pixelWidth)
		}

		// Compute expected ContentOrigin().Y
		var expectedY int

		expectedY = padInsetY

		if gotY != expectedY {
			t.Fatalf("ContentOrigin().Y = %d, want %d (glyphHeight=%d, pixelHeight=%d, paddingPct=%d)",
				gotY, expectedY, glyphHeight, pixelHeight, paddingPct)
		}
	})
}

// deriveBarHeight reimplements the expected bar height formula for test verification.
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

// For any valid LayoutBridge and any row index N (including negative values) and any
// pixel offset, RowY(n, offset) SHALL equal ContentOrigin().Y + max(0, n) × RowHeight() + offset,
// with no clamping applied to the output.

func TestProperty1_RowYFormula(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate random TextHints
		pixelWidth := rapid.IntRange(1, 500).Draw(t, "pixelWidth")
		pixelHeight := rapid.IntRange(1, 500).Draw(t, "pixelHeight")
		glyphHeight := rapid.IntRange(-10, 50).Draw(t, "glyphHeight")
		rowHeight := rapid.IntRange(-5, 30).Draw(t, "rowHeight")

		hints := textlayout.TextHints{
			PixelWidth:  pixelWidth,
			PixelHeight: pixelHeight,
			GlyphHeight: glyphHeight,
			RowHeight:   rowHeight,
		}

		// Generate random BridgeConfig
		paddingPct := rapid.IntRange(0, 50).Draw(t, "paddingPct")

		cfg := layout.BridgeConfig{
			PaddingPct: paddingPct,
		}

		// Generate random row index and offset
		n := rapid.IntRange(-100, 100).Draw(t, "n")
		offset := rapid.IntRange(-200, 200).Draw(t, "offset")

		// Construct bridge and call RowY
		lb := layout.NewLayoutBridge(hints, cfg)
		got := lb.RowY(n, offset)

		// Compute expected: ContentOrigin().Y + max(0, n) * RowHeight() + offset
		_, originY := lb.ContentOrigin()
		effectiveN := n
		if effectiveN < 0 {
			effectiveN = 0
		}
		expected := originY + effectiveN*lb.RowHeight() + offset

		if got != expected {
			t.Fatalf("RowY(%d, %d) = %d, want %d (ContentOrigin().Y=%d, RowHeight()=%d)",
				n, offset, got, expected, originY, lb.RowHeight())
		}
	})
}

// For any valid LayoutBridge and any row index N (including negatives),
// RowBottomY(n) SHALL equal ContentOrigin().Y + (max(0, n) + 1) * RowHeight().
// For any item count (including negatives), ItemsBottomY(count) SHALL equal
// ContentOrigin().Y + max(0, count) * RowHeight().

func TestProperty2_BottomEdgeFormula(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate random TextHints
		pixelWidth := rapid.IntRange(1, 500).Draw(t, "pixelWidth")
		pixelHeight := rapid.IntRange(1, 500).Draw(t, "pixelHeight")
		rowHeight := rapid.IntRange(-5, 30).Draw(t, "rowHeight")
		glyphHeight := rapid.IntRange(-10, 50).Draw(t, "glyphHeight")

		hints := textlayout.TextHints{
			PixelWidth:  pixelWidth,
			PixelHeight: pixelHeight,
			RowHeight:   rowHeight,
			GlyphHeight: glyphHeight,
		}

		// Generate random BridgeConfig
		paddingPct := rapid.IntRange(0, 50).Draw(t, "paddingPct")

		cfg := layout.BridgeConfig{
			PaddingPct: paddingPct,
		}

		lb := layout.NewLayoutBridge(hints, cfg)

		// Generate random row index N in [-100, 100]
		n := rapid.IntRange(-100, 100).Draw(t, "rowIndex")
		// Generate random item count in [-100, 100]
		itemCount := rapid.IntRange(-100, 100).Draw(t, "itemCount")

		_, originY := lb.ContentOrigin()
		rh := lb.RowHeight()

		// Verify RowBottomY(n) == ContentOrigin().Y + (max(0, n) + 1) * RowHeight()
		clampedN := n
		if clampedN < 0 {
			clampedN = 0
		}
		expectedRowBottom := originY + (clampedN+1)*rh
		gotRowBottom := lb.RowBottomY(n)
		if gotRowBottom != expectedRowBottom {
			t.Fatalf("RowBottomY(%d) = %d, want %d (originY=%d, rowHeight=%d)",
				n, gotRowBottom, expectedRowBottom, originY, rh)
		}

		// Verify ItemsBottomY(count) == ContentOrigin().Y + max(0, count) * RowHeight()
		clampedCount := itemCount
		if clampedCount < 0 {
			clampedCount = 0
		}
		expectedItemsBottom := originY + clampedCount*rh
		gotItemsBottom := lb.ItemsBottomY(itemCount)
		if gotItemsBottom != expectedItemsBottom {
			t.Fatalf("ItemsBottomY(%d) = %d, want %d (originY=%d, rowHeight=%d)",
				itemCount, gotItemsBottom, expectedItemsBottom, originY, rh)
		}
	})
}

// For any valid LayoutBridge and any integer N, RemainingSpace(n) SHALL equal
// max(0, AvailableContentHeight() - max(0, n) * RowHeight()).

func TestProperty5_RemainingSpace(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate random TextHints
		pixelWidth := rapid.IntRange(1, 500).Draw(t, "pixelWidth")
		pixelHeight := rapid.IntRange(1, 500).Draw(t, "pixelHeight")
		rowHeight := rapid.IntRange(-5, 30).Draw(t, "rowHeight")

		hints := textlayout.TextHints{
			PixelWidth:  pixelWidth,
			PixelHeight: pixelHeight,
			RowHeight:   rowHeight,
		}

		// Generate random BridgeConfig
		paddingPct := rapid.IntRange(0, 50).Draw(t, "paddingPct")

		cfg := layout.BridgeConfig{
			PaddingPct: paddingPct,
		}

		lb := layout.NewLayoutBridge(hints, cfg)

		// Generate random n in [-100, 100]
		n := rapid.IntRange(-100, 100).Draw(t, "n")

		// Compute expected: max(0, AvailableContentHeight() - max(0, n) * RowHeight())
		clampedN := n
		if clampedN < 0 {
			clampedN = 0
		}
		expected := lb.AvailableContentHeight() - clampedN*lb.RowHeight()
		if expected < 0 {
			expected = 0
		}

		got := lb.RemainingSpace(n)
		if got != expected {
			t.Fatalf("RemainingSpace(%d) = %d, want %d (AvailableContentHeight=%d, RowHeight=%d)",
				n, got, expected, lb.AvailableContentHeight(), lb.RowHeight())
		}
	})
}

// For any valid LayoutBridge, RowX(charOffset) SHALL equal
// contentOriginX + max(0, charOffset) × GlyphAdvance() where contentOriginX
// is the padding inset. The result SHALL NOT be clamped to the panel width.

func TestProperty6_RowXFormula(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate random TextHints
		pixelWidth := rapid.IntRange(0, 500).Draw(t, "pixelWidth")
		pixelHeight := rapid.IntRange(0, 500).Draw(t, "pixelHeight")
		glyphAdvance := rapid.IntRange(-5, 20).Draw(t, "glyphAdvance")

		hints := textlayout.TextHints{
			PixelWidth:   pixelWidth,
			PixelHeight:  pixelHeight,
			GlyphAdvance: glyphAdvance,
		}

		// Generate random BridgeConfig
		paddingPct := rapid.IntRange(0, 50).Draw(t, "paddingPct")
		cfg := layout.BridgeConfig{
			PaddingPct: paddingPct,
		}

		// Generate random charOffset including negatives
		charOffset := rapid.IntRange(-50, 100).Draw(t, "charOffset")

		lb := layout.NewLayoutBridge(hints, cfg)
		got := lb.RowX(charOffset)

		// Compute expected padding inset (contentOriginX)
		padInsetX := paddingPct * pixelWidth / 100

		// Effective char offset: max(0, charOffset)
		effectiveCharOffset := charOffset
		if effectiveCharOffset < 0 {
			effectiveCharOffset = 0
		}

		expected := padInsetX + effectiveCharOffset*lb.GlyphAdvance()

		if got != expected {
			t.Fatalf("RowX(%d) = %d, want %d (padInsetX=%d, glyphAdvance=%d, "+
				"pixelWidth=%d, paddingPct=%d)",
				charOffset, got, expected, padInsetX, lb.GlyphAdvance(),
				pixelWidth, paddingPct)
		}
	})
}

// For any valid LayoutBridge, MaxVisibleRows() SHALL equal AvailableContentHeight() / RowHeight()
// (integer floor division), and the invariant MaxVisibleRows() * RowHeight() <= AvailableContentHeight()
// SHALL always hold. When AvailableContentHeight() is zero, MaxVisibleRows() SHALL be zero.

func TestProperty10_MaxVisibleRowsInvariant(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate random TextHints
		pixelWidth := rapid.IntRange(1, 500).Draw(t, "pixelWidth")
		pixelHeight := rapid.IntRange(1, 500).Draw(t, "pixelHeight")
		rowHeight := rapid.IntRange(-5, 30).Draw(t, "rowHeight")
		glyphHeight := rapid.IntRange(-10, 50).Draw(t, "glyphHeight")

		hints := textlayout.TextHints{
			PixelWidth:  pixelWidth,
			PixelHeight: pixelHeight,
			RowHeight:   rowHeight,
			GlyphHeight: glyphHeight,
		}

		// Generate random BridgeConfig
		paddingPct := rapid.IntRange(0, 50).Draw(t, "paddingPct")

		cfg := layout.BridgeConfig{
			PaddingPct: paddingPct,
		}

		lb := layout.NewLayoutBridge(hints, cfg)

		maxRows := lb.MaxVisibleRows()
		availHeight := lb.AvailableContentHeight()
		rh := lb.RowHeight()

		// Verify: if RowHeight() == 0, MaxVisibleRows() == 0
		// (shouldn't happen due to defaults, but check anyway)
		if rh == 0 {
			if maxRows != 0 {
				t.Fatalf("RowHeight()==0 but MaxVisibleRows()=%d, want 0", maxRows)
			}
			return
		}

		// Verify: when AvailableContentHeight() is zero, MaxVisibleRows() is zero
		if availHeight == 0 {
			if maxRows != 0 {
				t.Fatalf("AvailableContentHeight()==0 but MaxVisibleRows()=%d, want 0", maxRows)
			}
			return
		}

		// Verify: MaxVisibleRows() == AvailableContentHeight() / RowHeight() (integer floor division)
		expectedMaxRows := availHeight / rh
		if maxRows != expectedMaxRows {
			t.Fatalf("MaxVisibleRows()=%d, want AvailableContentHeight()/RowHeight()=%d/%d=%d",
				maxRows, availHeight, rh, expectedMaxRows)
		}

		// Verify invariant: MaxVisibleRows() * RowHeight() <= AvailableContentHeight()
		product := maxRows * rh
		if product > availHeight {
			t.Fatalf("MaxVisibleRows()*RowHeight()=%d*%d=%d > AvailableContentHeight()=%d",
				maxRows, rh, product, availHeight)
		}
	})
}

// For any valid LayoutBridge, ContentCenter() SHALL return
// (ContentOrigin().X + max(0, availableContentWidth)/2,
//
//	ContentOrigin().Y + max(0, AvailableContentHeight())/2)
//
// using integer floor division. When available width or height is zero, the center
// coordinate degenerates to the corresponding content origin coordinate.

func TestProperty11_ContentCenterFormula(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate random TextHints
		pixelWidth := rapid.IntRange(0, 500).Draw(t, "pixelWidth")
		pixelHeight := rapid.IntRange(0, 500).Draw(t, "pixelHeight")
		glyphHeight := rapid.IntRange(-10, 50).Draw(t, "glyphHeight")

		hints := textlayout.TextHints{
			PixelWidth:  pixelWidth,
			PixelHeight: pixelHeight,
			GlyphHeight: glyphHeight,
		}

		// Generate random BridgeConfig
		paddingPct := rapid.IntRange(0, 50).Draw(t, "paddingPct")

		cfg := layout.BridgeConfig{
			PaddingPct: paddingPct,
		}

		lb := layout.NewLayoutBridge(hints, cfg)

		// Get actual ContentCenter
		gotX, gotY := lb.ContentCenter()

		// Compute expected values
		originX, originY := lb.ContentOrigin()

		// Available content width from the bridge
		availableWidth := lb.AvailableContentWidth()

		// Available content height from the bridge
		availableHeight := lb.AvailableContentHeight()

		// Expected center using integer floor division
		expectedX := originX + availableWidth/2
		expectedY := originY + availableHeight/2

		if gotX != expectedX {
			t.Fatalf("ContentCenter().X = %d, want %d (originX=%d, availableWidth=%d, pixelWidth=%d, paddingPct=%d)",
				gotX, expectedX, originX, availableWidth, pixelWidth, paddingPct)
		}
		if gotY != expectedY {
			t.Fatalf("ContentCenter().Y = %d, want %d (originY=%d, availableHeight=%d, pixelHeight=%d, paddingPct=%d)",
				gotY, expectedY, originY, availableHeight, pixelHeight, paddingPct)
		}
	})
}

func TestProperty4_AvailableContentHeight(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate random TextHints
		pixelWidth := rapid.IntRange(0, 500).Draw(t, "pixelWidth")
		pixelHeight := rapid.IntRange(0, 500).Draw(t, "pixelHeight")
		glyphHeight := rapid.IntRange(-10, 50).Draw(t, "glyphHeight")

		hints := textlayout.TextHints{
			PixelWidth:  pixelWidth,
			PixelHeight: pixelHeight,
			GlyphHeight: glyphHeight,
		}

		// Generate random BridgeConfig
		paddingPct := rapid.IntRange(0, 50).Draw(t, "paddingPct")

		cfg := layout.BridgeConfig{
			PaddingPct: paddingPct,
		}

		lb := layout.NewLayoutBridge(hints, cfg)

		// Compute expected padding inset Y
		padInsetY := paddingPct * pixelHeight / 100

		// Get content origin Y from the bridge
		_, originY := lb.ContentOrigin()

		// Compute bottomChrome
		var bottomChrome int

		bottomChrome = padInsetY

		// Compute expected available content height
		expected := pixelHeight - originY - bottomChrome
		if expected < 0 {
			expected = 0
		}

		got := lb.AvailableContentHeight()
		if got != expected {
			t.Fatalf("AvailableContentHeight() = %d, want %d "+
				"(pixelWidth=%d, pixelHeight=%d, glyphHeight=%d, paddingPct=%d "+
				"originY=%d, bottomChrome=%d)",
				got, expected,
				pixelWidth, pixelHeight, glyphHeight, paddingPct,
				originY, bottomChrome)
		}
	})
}

// For any TextHints and BridgeConfig, constructing two LayoutBridge values from
// identical inputs SHALL produce identical results from all methods (RowY, RowBottomY,
// ItemsBottomY, ContentOrigin, AvailableContentHeight, RemainingSpace, RowX,
// MaxVisibleRows, ContentCenter, RowHeight, GlyphAdvance).

func TestProperty9_Determinism(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate random TextHints with all fields varied
		pixelWidth := rapid.IntRange(0, 500).Draw(t, "pixelWidth")
		pixelHeight := rapid.IntRange(0, 500).Draw(t, "pixelHeight")
		glyphHeight := rapid.IntRange(-10, 50).Draw(t, "glyphHeight")
		rowHeight := rapid.IntRange(-5, 30).Draw(t, "rowHeight")
		glyphAdvance := rapid.IntRange(-5, 20).Draw(t, "glyphAdvance")

		hints := textlayout.TextHints{
			PixelWidth:   pixelWidth,
			PixelHeight:  pixelHeight,
			GlyphHeight:  glyphHeight,
			RowHeight:    rowHeight,
			GlyphAdvance: glyphAdvance,
		}

		// Generate random BridgeConfig
		paddingPct := rapid.IntRange(0, 50).Draw(t, "paddingPct")

		cfg := layout.BridgeConfig{
			PaddingPct: paddingPct,
		}

		// Construct two LayoutBridge values from identical inputs
		lb1 := layout.NewLayoutBridge(hints, cfg)
		lb2 := layout.NewLayoutBridge(hints, cfg)

		// Generate random row indices, offsets, and char offsets for method calls
		rowIndex := rapid.IntRange(-50, 100).Draw(t, "rowIndex")
		offset := rapid.IntRange(-200, 200).Draw(t, "offset")
		charOffset := rapid.IntRange(-50, 100).Draw(t, "charOffset")
		itemCount := rapid.IntRange(-50, 100).Draw(t, "itemCount")

		// Verify RowY determinism
		if lb1.RowY(rowIndex, offset) != lb2.RowY(rowIndex, offset) {
			t.Fatalf("RowY(%d, %d) not deterministic: lb1=%d, lb2=%d",
				rowIndex, offset, lb1.RowY(rowIndex, offset), lb2.RowY(rowIndex, offset))
		}

		// Verify RowBottomY determinism
		if lb1.RowBottomY(rowIndex) != lb2.RowBottomY(rowIndex) {
			t.Fatalf("RowBottomY(%d) not deterministic: lb1=%d, lb2=%d",
				rowIndex, lb1.RowBottomY(rowIndex), lb2.RowBottomY(rowIndex))
		}

		// Verify ItemsBottomY determinism
		if lb1.ItemsBottomY(itemCount) != lb2.ItemsBottomY(itemCount) {
			t.Fatalf("ItemsBottomY(%d) not deterministic: lb1=%d, lb2=%d",
				itemCount, lb1.ItemsBottomY(itemCount), lb2.ItemsBottomY(itemCount))
		}

		// Verify ContentOrigin determinism
		x1, y1 := lb1.ContentOrigin()
		x2, y2 := lb2.ContentOrigin()
		if x1 != x2 || y1 != y2 {
			t.Fatalf("ContentOrigin() not deterministic: lb1=(%d,%d), lb2=(%d,%d)",
				x1, y1, x2, y2)
		}

		// Verify AvailableContentHeight determinism
		if lb1.AvailableContentHeight() != lb2.AvailableContentHeight() {
			t.Fatalf("AvailableContentHeight() not deterministic: lb1=%d, lb2=%d",
				lb1.AvailableContentHeight(), lb2.AvailableContentHeight())
		}

		// Verify RemainingSpace determinism
		if lb1.RemainingSpace(rowIndex) != lb2.RemainingSpace(rowIndex) {
			t.Fatalf("RemainingSpace(%d) not deterministic: lb1=%d, lb2=%d",
				rowIndex, lb1.RemainingSpace(rowIndex), lb2.RemainingSpace(rowIndex))
		}

		// Verify RowX determinism
		if lb1.RowX(charOffset) != lb2.RowX(charOffset) {
			t.Fatalf("RowX(%d) not deterministic: lb1=%d, lb2=%d",
				charOffset, lb1.RowX(charOffset), lb2.RowX(charOffset))
		}

		// Verify MaxVisibleRows determinism
		if lb1.MaxVisibleRows() != lb2.MaxVisibleRows() {
			t.Fatalf("MaxVisibleRows() not deterministic: lb1=%d, lb2=%d",
				lb1.MaxVisibleRows(), lb2.MaxVisibleRows())
		}

		// Verify ContentCenter determinism
		cx1, cy1 := lb1.ContentCenter()
		cx2, cy2 := lb2.ContentCenter()
		if cx1 != cx2 || cy1 != cy2 {
			t.Fatalf("ContentCenter() not deterministic: lb1=(%d,%d), lb2=(%d,%d)",
				cx1, cy1, cx2, cy2)
		}

		// Verify RowHeight determinism
		if lb1.RowHeight() != lb2.RowHeight() {
			t.Fatalf("RowHeight() not deterministic: lb1=%d, lb2=%d",
				lb1.RowHeight(), lb2.RowHeight())
		}

		// Verify GlyphAdvance determinism
		if lb1.GlyphAdvance() != lb2.GlyphAdvance() {
			t.Fatalf("GlyphAdvance() not deterministic: lb1=%d, lb2=%d",
				lb1.GlyphAdvance(), lb2.GlyphAdvance())
		}
	})
}

// For any valid LayoutBridge, any integer charCount, and any integer glyphAdvance parameter:
// - When charCount <= 0, both TextPixelWidth(charCount) and TextPixelWidthWith(charCount, glyphAdvance) return 0.
// - When charCount > 0 and glyphAdvance <= 0, TextPixelWidthWith falls back to the bridge's GlyphAdvance() and returns GlyphAdvance() × charCount.
// - When charCount > 0 and glyphAdvance > 0, TextPixelWidthWith returns glyphAdvance × charCount.
// - TextPixelWidth(charCount) equals TextPixelWidthWith(charCount, GlyphAdvance()) for all inputs.

func TestExtProperty1_TextPixelWidth(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate random TextHints
		pixelWidth := rapid.IntRange(0, 500).Draw(t, "pixelWidth")
		pixelHeight := rapid.IntRange(0, 500).Draw(t, "pixelHeight")
		glyphAdvance := rapid.IntRange(-5, 20).Draw(t, "glyphAdvance")
		glyphHeight := rapid.IntRange(-10, 50).Draw(t, "glyphHeight")

		hints := textlayout.TextHints{
			PixelWidth:   pixelWidth,
			PixelHeight:  pixelHeight,
			GlyphAdvance: glyphAdvance,
			GlyphHeight:  glyphHeight,
		}

		// Generate random BridgeConfig
		paddingPct := rapid.IntRange(0, 50).Draw(t, "paddingPct")

		cfg := layout.BridgeConfig{
			PaddingPct: paddingPct,
		}

		lb := layout.NewLayoutBridge(hints, cfg)

		// Generate random charCount and explicit glyph advance
		charCount := rapid.IntRange(-50, 200).Draw(t, "charCount")
		explicitGlyphAdvance := rapid.IntRange(-10, 30).Draw(t, "explicitGlyphAdvance")

		// Property: When charCount <= 0, both TextPixelWidth and TextPixelWidthWith return 0
		if charCount <= 0 {
			got := lb.TextPixelWidth(charCount)
			if got != 0 {
				t.Fatalf("TextPixelWidth(%d) = %d, want 0 (charCount <= 0)", charCount, got)
			}
			gotWith := lb.TextPixelWidthWith(charCount, explicitGlyphAdvance)
			if gotWith != 0 {
				t.Fatalf("TextPixelWidthWith(%d, %d) = %d, want 0 (charCount <= 0)",
					charCount, explicitGlyphAdvance, gotWith)
			}
			return
		}

		// charCount > 0 from here

		// Property: When glyphAdvance <= 0, TextPixelWidthWith falls back to bridge's GlyphAdvance()
		if explicitGlyphAdvance <= 0 {
			gotWith := lb.TextPixelWidthWith(charCount, explicitGlyphAdvance)
			expected := lb.GlyphAdvance() * charCount
			if gotWith != expected {
				t.Fatalf("TextPixelWidthWith(%d, %d) = %d, want %d (fallback to GlyphAdvance()=%d)",
					charCount, explicitGlyphAdvance, gotWith, expected, lb.GlyphAdvance())
			}
		}

		// Property: When charCount > 0 and glyphAdvance > 0, TextPixelWidthWith returns glyphAdvance × charCount
		if explicitGlyphAdvance > 0 {
			gotWith := lb.TextPixelWidthWith(charCount, explicitGlyphAdvance)
			expected := explicitGlyphAdvance * charCount
			if gotWith != expected {
				t.Fatalf("TextPixelWidthWith(%d, %d) = %d, want %d",
					charCount, explicitGlyphAdvance, gotWith, expected)
			}
		}

		// Property: TextPixelWidth(charCount) equals TextPixelWidthWith(charCount, GlyphAdvance())
		got := lb.TextPixelWidth(charCount)
		gotWith := lb.TextPixelWidthWith(charCount, lb.GlyphAdvance())
		if got != gotWith {
			t.Fatalf("TextPixelWidth(%d) = %d != TextPixelWidthWith(%d, %d) = %d",
				charCount, got, charCount, lb.GlyphAdvance(), gotWith)
		}

		// Also verify TextPixelWidth uses the bridge's configured advance
		expectedWidth := lb.GlyphAdvance() * charCount
		if got != expectedWidth {
			t.Fatalf("TextPixelWidth(%d) = %d, want GlyphAdvance()×charCount = %d×%d = %d",
				charCount, got, lb.GlyphAdvance(), charCount, expectedWidth)
		}
	})
}

// For any valid LayoutBridge, any integer charCount, and any integer glyphAdvance parameter:
// - CenterX(charCount) equals max(0, (AvailableContentWidth() - TextPixelWidth(charCount)) / 2) using integer division.
// - CenterXWith(charCount, glyphAdvance) equals max(0, (AvailableContentWidth() - TextPixelWidthWith(charCount, glyphAdvance)) / 2).
// - When charCount <= 0, text pixel width is 0, so result is AvailableContentWidth() / 2.

func TestExtProperty2_CenterX(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate random TextHints
		pixelWidth := rapid.IntRange(0, 500).Draw(t, "pixelWidth")
		pixelHeight := rapid.IntRange(0, 500).Draw(t, "pixelHeight")
		glyphAdvance := rapid.IntRange(-5, 20).Draw(t, "glyphAdvance")
		glyphHeight := rapid.IntRange(-10, 50).Draw(t, "glyphHeight")

		hints := textlayout.TextHints{
			PixelWidth:   pixelWidth,
			PixelHeight:  pixelHeight,
			GlyphAdvance: glyphAdvance,
			GlyphHeight:  glyphHeight,
		}

		// Generate random BridgeConfig
		paddingPct := rapid.IntRange(0, 50).Draw(t, "paddingPct")

		cfg := layout.BridgeConfig{
			PaddingPct: paddingPct,
		}

		lb := layout.NewLayoutBridge(hints, cfg)

		// Generate random charCount and explicit glyph advance
		charCount := rapid.IntRange(-50, 200).Draw(t, "charCount")
		explicitGlyphAdvance := rapid.IntRange(-10, 30).Draw(t, "explicitGlyphAdvance")

		acw := lb.AvailableContentWidth()

		// Property: CenterX(charCount) == max(0, (ACW - TextPixelWidth(charCount)) / 2)
		{
			got := lb.CenterX(charCount)
			textWidth := lb.TextPixelWidth(charCount)
			expected := (acw - textWidth) / 2
			if expected < 0 {
				expected = 0
			}
			if got != expected {
				t.Fatalf("CenterX(%d) = %d, want max(0, (ACW=%d - TextPixelWidth=%d) / 2) = %d",
					charCount, got, acw, textWidth, expected)
			}
		}

		// Property: CenterXWith(charCount, glyphAdvance) == max(0, (ACW - TextPixelWidthWith(charCount, glyphAdvance)) / 2)
		{
			got := lb.CenterXWith(charCount, explicitGlyphAdvance)
			textWidth := lb.TextPixelWidthWith(charCount, explicitGlyphAdvance)
			expected := (acw - textWidth) / 2
			if expected < 0 {
				expected = 0
			}
			if got != expected {
				t.Fatalf("CenterXWith(%d, %d) = %d, want max(0, (ACW=%d - TextPixelWidthWith=%d) / 2) = %d",
					charCount, explicitGlyphAdvance, got, acw, textWidth, expected)
			}
		}

		// Property: When charCount <= 0, result is ACW / 2
		if charCount <= 0 {
			got := lb.CenterX(charCount)
			expected := acw / 2
			if got != expected {
				t.Fatalf("CenterX(%d) = %d, want ACW/2 = %d/2 = %d (charCount <= 0)",
					charCount, got, acw, expected)
			}
		}
	})
}

// For any valid LayoutBridge constructed from arbitrary TextHints and BridgeConfig,
// AvailableContentWidth() equals max(0, PixelWidth - 2 × padInsetX) where
// padInsetX = PaddingPct * PixelWidth / 100.

func TestExtProperty4_AvailableContentWidth(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate random TextHints
		pixelWidth := rapid.IntRange(0, 500).Draw(t, "pixelWidth")
		pixelHeight := rapid.IntRange(0, 500).Draw(t, "pixelHeight")
		glyphHeight := rapid.IntRange(-10, 50).Draw(t, "glyphHeight")

		hints := textlayout.TextHints{
			PixelWidth:  pixelWidth,
			PixelHeight: pixelHeight,
			GlyphHeight: glyphHeight,
		}

		// Generate random BridgeConfig
		paddingPct := rapid.IntRange(0, 50).Draw(t, "paddingPct")

		cfg := layout.BridgeConfig{
			PaddingPct: paddingPct,
		}

		lb := layout.NewLayoutBridge(hints, cfg)

		// Compute expected padInsetX
		padInsetX := paddingPct * pixelWidth / 100

		// Expected: max(0, PixelWidth - 2 × padInsetX)
		expected := pixelWidth - 2*padInsetX
		if expected < 0 {
			expected = 0
		}

		got := lb.AvailableContentWidth()
		if got != expected {
			t.Fatalf("AvailableContentWidth() = %d, want max(0, %d - 2×%d) = %d "+
				"(pixelWidth=%d, paddingPct=%d, padInsetX=%d)",
				got, pixelWidth, padInsetX, expected,
				pixelWidth, paddingPct, padInsetX)
		}

		// Also verify consistency: AvailableContentWidth() is what CenterX uses internally
		// This verifies Requirement 4.2: same value used by horizontal centering
		charCount := rapid.IntRange(1, 50).Draw(t, "charCountForConsistency")
		centerX := lb.CenterX(charCount)
		textWidth := lb.TextPixelWidth(charCount)
		expectedCenter := (got - textWidth) / 2
		if expectedCenter < 0 {
			expectedCenter = 0
		}
		if centerX != expectedCenter {
			t.Fatalf("CenterX(%d) = %d inconsistent with AvailableContentWidth() = %d: "+
				"expected max(0, (%d - %d) / 2) = %d",
				charCount, centerX, got, got, textWidth, expectedCenter)
		}
	})
}

// For any valid LayoutBridge, any slice of row heights, and any integer spacing value:
// - When the slice is empty, CenterBlockY returns 0.
// - When spacing is negative, it is treated as 0.
// - Content block height = sum(rowHeights) + max(0, spacing) × (len(rowHeights) - 1).
// - CenterBlockY returns max(0, (AvailableContentHeight() - contentBlockHeight) / 2).

func TestExtProperty3_CenterBlockY(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate random TextHints for a realistic panel
		pixelWidth := rapid.IntRange(32, 500).Draw(t, "pixelWidth")
		pixelHeight := rapid.IntRange(32, 500).Draw(t, "pixelHeight")
		glyphHeight := rapid.IntRange(-10, 50).Draw(t, "glyphHeight")

		hints := textlayout.TextHints{
			PixelWidth:  pixelWidth,
			PixelHeight: pixelHeight,
			GlyphHeight: glyphHeight,
		}

		// Generate random BridgeConfig
		paddingPct := rapid.IntRange(0, 50).Draw(t, "paddingPct")

		cfg := layout.BridgeConfig{
			PaddingPct: paddingPct,
		}

		lb := layout.NewLayoutBridge(hints, cfg)

		// Generate row heights slice (0 to 5 rows, heights in [1, 60])
		numRows := rapid.IntRange(0, 5).Draw(t, "numRows")
		rowHeights := make([]int, numRows)
		for i := range rowHeights {
			rowHeights[i] = rapid.IntRange(1, 60).Draw(t, "rowHeight")
		}

		// Generate spacing (including negatives to test the clamp)
		spacing := rapid.IntRange(-10, 30).Draw(t, "spacing")

		got := lb.CenterBlockY(rowHeights, spacing)

		// Property: When the slice is empty, CenterBlockY returns 0
		if numRows == 0 {
			if got != 0 {
				t.Fatalf("CenterBlockY(empty, %d) = %d, want 0", spacing, got)
			}
			return
		}

		// Compute expected block height
		effectiveSpacing := spacing
		if effectiveSpacing < 0 {
			effectiveSpacing = 0
		}

		blockHeight := 0
		for _, h := range rowHeights {
			blockHeight += h
		}
		blockHeight += effectiveSpacing * (len(rowHeights) - 1)

		// Expected: max(0, (ACH - blockHeight) / 2)
		ach := lb.AvailableContentHeight()
		expected := (ach - blockHeight) / 2
		if expected < 0 {
			expected = 0
		}

		if got != expected {
			t.Fatalf("CenterBlockY(%v, %d) = %d, want max(0, (ACH=%d - blockH=%d) / 2) = %d",
				rowHeights, spacing, got, ach, blockHeight, expected)
		}
	})
}

// For any valid LayoutBridge and any non-empty slice of row heights:
// 1. Minimum visibility: visibleCount >= 1.
// 2. Fit guarantee: sum of first visibleCount row heights + spacing × (visibleCount - 1) <= ACH, OR visibleCount == 1.
// 3. Spacing bounds (small panel): When ACH <= 135 and visibleCount > 1, spacing in [1, 2].
// 4. Spacing bounds (large panel): When ACH > 135 and visibleCount > 1, spacing <= 8.
// 5. Single-row spacing: When visibleCount <= 1, spacing == 0.
// 6. Vertical centering: offsetY == max(0, (ACH - contentBlockH) / 2).
// 7. Empty input: all three values are 0.

func TestExtProperty5_FitRows(t *testing.T) {
	// Test the empty input case explicitly first.
	t.Run("empty_input", func(t *testing.T) {
		hints := textlayout.TextHints{
			PixelWidth:  128,
			PixelHeight: 128,
		}
		cfg := layout.BridgeConfig{}
		lb := layout.NewLayoutBridge(hints, cfg)

		spacing, offsetY, visibleCount := lb.FitRows([]int{})
		if spacing != 0 || offsetY != 0 || visibleCount != 0 {
			t.Fatalf("FitRows(empty) = (%d, %d, %d), want (0, 0, 0)",
				spacing, offsetY, visibleCount)
		}
	})

	// Property-based test for non-empty inputs.
	rapid.Check(t, func(t *rapid.T) {
		// Use realistic panel sizes (64-240px height)
		pixelWidth := rapid.IntRange(64, 240).Draw(t, "pixelWidth")
		pixelHeight := rapid.IntRange(64, 240).Draw(t, "pixelHeight")
		glyphHeight := rapid.IntRange(0, 20).Draw(t, "glyphHeight")

		hints := textlayout.TextHints{
			PixelWidth:  pixelWidth,
			PixelHeight: pixelHeight,
			GlyphHeight: glyphHeight,
		}

		// Generate BridgeConfig
		paddingPct := rapid.IntRange(0, 20).Draw(t, "paddingPct")

		cfg := layout.BridgeConfig{
			PaddingPct: paddingPct,
		}

		lb := layout.NewLayoutBridge(hints, cfg)
		ach := lb.AvailableContentHeight()

		// Generate row heights: 1-5 rows, individual heights in [8, 40]
		numRows := rapid.IntRange(1, 5).Draw(t, "numRows")
		rowHeights := make([]int, numRows)
		for i := range rowHeights {
			rowHeights[i] = rapid.IntRange(8, 40).Draw(t, "rowH")
		}

		spacing, offsetY, visibleCount := lb.FitRows(rowHeights)

		// Invariant 1: Minimum visibility - visibleCount >= 1
		if visibleCount < 1 {
			t.Fatalf("FitRows: visibleCount=%d < 1 (rowHeights=%v, ACH=%d)",
				visibleCount, rowHeights, ach)
		}

		// Invariant 2: Fit guarantee
		visibleSum := 0
		for i := 0; i < visibleCount; i++ {
			visibleSum += rowHeights[i]
		}
		contentWithSpacing := visibleSum + spacing*(visibleCount-1)
		if contentWithSpacing > ach && visibleCount != 1 {
			t.Fatalf("FitRows: content block %d > ACH %d with visibleCount=%d (not single-row preservation); "+
				"spacing=%d, rowHeights=%v",
				contentWithSpacing, ach, visibleCount, spacing, rowHeights[:visibleCount])
		}

		// Invariant 3: Spacing bounds (small panel): When ACH <= 135 and visibleCount > 1, spacing in [1, 2]
		if ach <= 135 && visibleCount > 1 {
			if spacing < 1 || spacing > 2 {
				t.Fatalf("FitRows: small panel (ACH=%d) with visibleCount=%d, spacing=%d not in [1, 2]; rowHeights=%v",
					ach, visibleCount, spacing, rowHeights[:visibleCount])
			}
		}

		// Invariant 4: Spacing bounds (large panel): When ACH > 135 and visibleCount > 1, spacing <= 8
		if ach > 135 && visibleCount > 1 {
			if spacing > 8 {
				t.Fatalf("FitRows: large panel (ACH=%d) with visibleCount=%d, spacing=%d > 8; rowHeights=%v",
					ach, visibleCount, spacing, rowHeights[:visibleCount])
			}
		}

		// Invariant 5: Single-row spacing: When visibleCount <= 1, spacing == 0
		if visibleCount <= 1 {
			if spacing != 0 {
				t.Fatalf("FitRows: visibleCount=%d but spacing=%d != 0; rowHeights=%v, ACH=%d",
					visibleCount, spacing, rowHeights, ach)
			}
		}

		// Invariant 6: Vertical centering: offsetY == max(0, (ACH - contentBlockH) / 2)
		contentBlockH := visibleSum + spacing*(visibleCount-1)
		expectedOffsetY := (ach - contentBlockH) / 2
		if expectedOffsetY < 0 {
			expectedOffsetY = 0
		}
		if offsetY != expectedOffsetY {
			t.Fatalf("FitRows: offsetY=%d, want max(0, (ACH=%d - contentBlockH=%d) / 2) = %d; "+
				"spacing=%d, visibleCount=%d, rowHeights=%v",
				offsetY, ach, contentBlockH, expectedOffsetY, spacing, visibleCount, rowHeights[:visibleCount])
		}
	})
}

// For any valid LayoutBridge and any integer elementHeight:
// - When elementHeight <= 0, the result is ContentOrigin().Y + AvailableContentHeight().
// - When 0 < elementHeight <= AvailableContentHeight(), the result is ContentOrigin().Y + AvailableContentHeight() - elementHeight.
// - When elementHeight > AvailableContentHeight(), the result is ContentOrigin().Y.

func TestExtProperty6_BottomAnchorY(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate random TextHints
		pixelWidth := rapid.IntRange(1, 500).Draw(t, "pixelWidth")
		pixelHeight := rapid.IntRange(1, 500).Draw(t, "pixelHeight")
		glyphHeight := rapid.IntRange(-10, 50).Draw(t, "glyphHeight")

		hints := textlayout.TextHints{
			PixelWidth:  pixelWidth,
			PixelHeight: pixelHeight,
			GlyphHeight: glyphHeight,
		}

		// Generate random BridgeConfig
		paddingPct := rapid.IntRange(0, 50).Draw(t, "paddingPct")

		cfg := layout.BridgeConfig{
			PaddingPct: paddingPct,
		}

		lb := layout.NewLayoutBridge(hints, cfg)

		// Generate random elementHeight spanning negative, zero, within range, and overflow cases
		elementHeight := rapid.IntRange(-50, 600).Draw(t, "elementHeight")

		got := lb.BottomAnchorY(elementHeight)
		_, originY := lb.ContentOrigin()
		ach := lb.AvailableContentHeight()

		var expected int
		if elementHeight <= 0 {
			expected = originY + ach
		} else if elementHeight > ach {
			expected = originY
		} else {
			expected = originY + ach - elementHeight
		}

		if got != expected {
			t.Fatalf("BottomAnchorY(%d) = %d, want %d (ContentOrigin().Y=%d, AvailableContentHeight()=%d)",
				elementHeight, got, expected, originY, ach)
		}
	})
}

// For any valid LayoutBridge and any integer elementWidth:
// - When elementWidth <= 0, the result is ContentOrigin().X + AvailableContentWidth().
// - When 0 < elementWidth <= AvailableContentWidth(), the result is ContentOrigin().X + AvailableContentWidth() - elementWidth.
// - When elementWidth > AvailableContentWidth(), the result is ContentOrigin().X.

func TestExtProperty7_TopRightAnchorX(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate random TextHints
		pixelWidth := rapid.IntRange(1, 500).Draw(t, "pixelWidth")
		pixelHeight := rapid.IntRange(1, 500).Draw(t, "pixelHeight")
		glyphHeight := rapid.IntRange(-10, 50).Draw(t, "glyphHeight")

		hints := textlayout.TextHints{
			PixelWidth:  pixelWidth,
			PixelHeight: pixelHeight,
			GlyphHeight: glyphHeight,
		}

		// Generate random BridgeConfig
		paddingPct := rapid.IntRange(0, 50).Draw(t, "paddingPct")

		cfg := layout.BridgeConfig{
			PaddingPct: paddingPct,
		}

		lb := layout.NewLayoutBridge(hints, cfg)

		// Generate random elementWidth spanning negative, zero, within range, and overflow cases
		elementWidth := rapid.IntRange(-50, 600).Draw(t, "elementWidth")

		got := lb.TopRightAnchorX(elementWidth)
		originX, _ := lb.ContentOrigin()
		acw := lb.AvailableContentWidth()

		var expected int
		if elementWidth <= 0 {
			expected = originX + acw
		} else if elementWidth > acw {
			expected = originX
		} else {
			expected = originX + acw - elementWidth
		}

		if got != expected {
			t.Fatalf("TopRightAnchorX(%d) = %d, want %d (ContentOrigin().X=%d, AvailableContentWidth()=%d)",
				elementWidth, got, expected, originX, acw)
		}
	})
}

// For any valid LayoutBridge, any integers rowIndex, charOffset, rowHeight, and optionally glyphAdvance:
// - Row index and character offset are clamped to 0 when negative.
// - effectiveGlyphAdvance: the explicit param when > 0, else bridge's GlyphAdvance().
// - Min.X = ContentOrigin().X + max(0, charOffset) × effectiveGlyphAdvance
// - Min.Y = ContentOrigin().Y + max(0, rowIndex) × rowHeight
// - Max.X = ContentOrigin().X + AvailableContentWidth()
// - Max.Y = Min.Y + rowHeight
// - When Min.X > Max.X, rectangle is empty (Min == Max, i.e., image.Rectangle{}).

func TestExtProperty8_InlineWidgetBounds(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate random TextHints
		pixelWidth := rapid.IntRange(1, 500).Draw(t, "pixelWidth")
		pixelHeight := rapid.IntRange(1, 500).Draw(t, "pixelHeight")
		glyphAdvance := rapid.IntRange(-5, 20).Draw(t, "glyphAdvance")
		glyphHeight := rapid.IntRange(-10, 50).Draw(t, "glyphHeight")

		hints := textlayout.TextHints{
			PixelWidth:   pixelWidth,
			PixelHeight:  pixelHeight,
			GlyphAdvance: glyphAdvance,
			GlyphHeight:  glyphHeight,
		}

		// Generate random BridgeConfig
		paddingPct := rapid.IntRange(0, 50).Draw(t, "paddingPct")

		cfg := layout.BridgeConfig{
			PaddingPct: paddingPct,
		}

		lb := layout.NewLayoutBridge(hints, cfg)

		// Generate random parameters for InlineWidgetBounds
		rowIndex := rapid.IntRange(-20, 50).Draw(t, "rowIndex")
		charOffset := rapid.IntRange(-20, 100).Draw(t, "charOffset")
		rowHeight := rapid.IntRange(0, 40).Draw(t, "rowHeight")
		explicitGlyphAdvance := rapid.IntRange(-10, 30).Draw(t, "explicitGlyphAdvance")

		// Test InlineWidgetBoundsWith (the parameterized form)
		got := lb.InlineWidgetBoundsWith(rowIndex, charOffset, rowHeight, explicitGlyphAdvance)

		originX, originY := lb.ContentOrigin()
		acw := lb.AvailableContentWidth()

		// Compute effective glyph advance
		effectiveGA := explicitGlyphAdvance
		if effectiveGA <= 0 {
			effectiveGA = lb.GlyphAdvance()
		}

		// Clamp rowIndex and charOffset
		clampedRowIndex := rowIndex
		if clampedRowIndex < 0 {
			clampedRowIndex = 0
		}
		clampedCharOffset := charOffset
		if clampedCharOffset < 0 {
			clampedCharOffset = 0
		}

		// Compute expected bounds
		expectedMinX := originX + clampedCharOffset*effectiveGA
		expectedMinY := originY + clampedRowIndex*rowHeight
		expectedMaxX := originX + acw
		expectedMaxY := expectedMinY + rowHeight

		if expectedMinX > expectedMaxX {
			// When Min.X > Max.X, rectangle should be empty
			expectedRect := image.Rectangle{}
			if got != expectedRect {
				t.Fatalf("InlineWidgetBoundsWith(%d, %d, %d, %d) = %v, want empty rect %v "+
					"(Min.X=%d > Max.X=%d)",
					rowIndex, charOffset, rowHeight, explicitGlyphAdvance, got, expectedRect,
					expectedMinX, expectedMaxX)
			}
		} else {
			expectedRect := image.Rectangle{
				Min: image.Point{X: expectedMinX, Y: expectedMinY},
				Max: image.Point{X: expectedMaxX, Y: expectedMaxY},
			}
			if got != expectedRect {
				t.Fatalf("InlineWidgetBoundsWith(%d, %d, %d, %d) = %v, want %v "+
					"(originX=%d, originY=%d, ACW=%d, effectiveGA=%d)",
					rowIndex, charOffset, rowHeight, explicitGlyphAdvance, got, expectedRect,
					originX, originY, acw, effectiveGA)
			}
		}

		// Also verify InlineWidgetBounds (non-parameterized) uses bridge's GlyphAdvance
		gotSimple := lb.InlineWidgetBounds(rowIndex, charOffset, rowHeight)
		gotWithBridgeGA := lb.InlineWidgetBoundsWith(rowIndex, charOffset, rowHeight, 0)
		if gotSimple != gotWithBridgeGA {
			t.Fatalf("InlineWidgetBounds(%d, %d, %d) = %v != InlineWidgetBoundsWith(%d, %d, %d, 0) = %v",
				rowIndex, charOffset, rowHeight, gotSimple,
				rowIndex, charOffset, rowHeight, gotWithBridgeGA)
		}
	})
}

// For any TextHints and BridgeConfig, constructing two LayoutBridge values from identical
// inputs produces identical results from ALL new methods when called with identical arguments:
// TextPixelWidth, TextPixelWidthWith, CenterX, CenterXWith, AvailableContentWidth,
// CenterBlockY, FitRows, BottomAnchorY, TopRightAnchorX, InlineWidgetBounds, InlineWidgetBoundsWith.

func TestExtProperty9_Determinism(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate all TextHints fields
		pixelWidth := rapid.IntRange(0, 500).Draw(t, "pixelWidth")
		pixelHeight := rapid.IntRange(0, 500).Draw(t, "pixelHeight")
		glyphHeight := rapid.IntRange(-10, 50).Draw(t, "glyphHeight")
		rowHeight := rapid.IntRange(-5, 30).Draw(t, "rowHeight")
		glyphAdvance := rapid.IntRange(-5, 20).Draw(t, "glyphAdvance")

		hints := textlayout.TextHints{
			PixelWidth:   pixelWidth,
			PixelHeight:  pixelHeight,
			GlyphHeight:  glyphHeight,
			RowHeight:    rowHeight,
			GlyphAdvance: glyphAdvance,
		}

		// Generate BridgeConfig
		paddingPct := rapid.IntRange(0, 50).Draw(t, "paddingPct")

		cfg := layout.BridgeConfig{
			PaddingPct: paddingPct,
		}

		// Construct two LayoutBridge values from identical inputs
		lb1 := layout.NewLayoutBridge(hints, cfg)
		lb2 := layout.NewLayoutBridge(hints, cfg)

		// Generate test parameters
		charCount := rapid.IntRange(-50, 200).Draw(t, "charCount")
		explicitGlyphAdvance := rapid.IntRange(-10, 30).Draw(t, "explicitGlyphAdvance")
		spacing := rapid.IntRange(-10, 30).Draw(t, "spacing")
		elementHeight := rapid.IntRange(-10, 300).Draw(t, "elementHeight")
		elementWidth := rapid.IntRange(-10, 300).Draw(t, "elementWidth")
		rowIndex := rapid.IntRange(-10, 50).Draw(t, "rowIndex")
		charOffset := rapid.IntRange(-10, 100).Draw(t, "charOffset")
		rh := rapid.IntRange(5, 30).Draw(t, "rowHeightParam")

		// Generate a small rowHeights slice (1-4 entries, values in [5, 30])
		numRows := rapid.IntRange(1, 4).Draw(t, "numRows")
		rowHeights := make([]int, numRows)
		for i := range rowHeights {
			rowHeights[i] = rapid.IntRange(5, 30).Draw(t, "rowH")
		}

		// Verify TextPixelWidth determinism
		if lb1.TextPixelWidth(charCount) != lb2.TextPixelWidth(charCount) {
			t.Fatalf("TextPixelWidth(%d) not deterministic: lb1=%d, lb2=%d",
				charCount, lb1.TextPixelWidth(charCount), lb2.TextPixelWidth(charCount))
		}

		// Verify TextPixelWidthWith determinism
		if lb1.TextPixelWidthWith(charCount, explicitGlyphAdvance) != lb2.TextPixelWidthWith(charCount, explicitGlyphAdvance) {
			t.Fatalf("TextPixelWidthWith(%d, %d) not deterministic: lb1=%d, lb2=%d",
				charCount, explicitGlyphAdvance,
				lb1.TextPixelWidthWith(charCount, explicitGlyphAdvance),
				lb2.TextPixelWidthWith(charCount, explicitGlyphAdvance))
		}

		// Verify CenterX determinism
		if lb1.CenterX(charCount) != lb2.CenterX(charCount) {
			t.Fatalf("CenterX(%d) not deterministic: lb1=%d, lb2=%d",
				charCount, lb1.CenterX(charCount), lb2.CenterX(charCount))
		}

		// Verify CenterXWith determinism
		if lb1.CenterXWith(charCount, explicitGlyphAdvance) != lb2.CenterXWith(charCount, explicitGlyphAdvance) {
			t.Fatalf("CenterXWith(%d, %d) not deterministic: lb1=%d, lb2=%d",
				charCount, explicitGlyphAdvance,
				lb1.CenterXWith(charCount, explicitGlyphAdvance),
				lb2.CenterXWith(charCount, explicitGlyphAdvance))
		}

		// Verify AvailableContentWidth determinism
		if lb1.AvailableContentWidth() != lb2.AvailableContentWidth() {
			t.Fatalf("AvailableContentWidth() not deterministic: lb1=%d, lb2=%d",
				lb1.AvailableContentWidth(), lb2.AvailableContentWidth())
		}

		// Verify CenterBlockY determinism
		if lb1.CenterBlockY(rowHeights, spacing) != lb2.CenterBlockY(rowHeights, spacing) {
			t.Fatalf("CenterBlockY(%v, %d) not deterministic: lb1=%d, lb2=%d",
				rowHeights, spacing,
				lb1.CenterBlockY(rowHeights, spacing),
				lb2.CenterBlockY(rowHeights, spacing))
		}

		// Verify FitRows determinism
		s1, o1, v1 := lb1.FitRows(rowHeights)
		s2, o2, v2 := lb2.FitRows(rowHeights)
		if s1 != s2 || o1 != o2 || v1 != v2 {
			t.Fatalf("FitRows(%v) not deterministic: lb1=(%d,%d,%d), lb2=(%d,%d,%d)",
				rowHeights, s1, o1, v1, s2, o2, v2)
		}

		// Verify BottomAnchorY determinism
		if lb1.BottomAnchorY(elementHeight) != lb2.BottomAnchorY(elementHeight) {
			t.Fatalf("BottomAnchorY(%d) not deterministic: lb1=%d, lb2=%d",
				elementHeight, lb1.BottomAnchorY(elementHeight), lb2.BottomAnchorY(elementHeight))
		}

		// Verify TopRightAnchorX determinism
		if lb1.TopRightAnchorX(elementWidth) != lb2.TopRightAnchorX(elementWidth) {
			t.Fatalf("TopRightAnchorX(%d) not deterministic: lb1=%d, lb2=%d",
				elementWidth, lb1.TopRightAnchorX(elementWidth), lb2.TopRightAnchorX(elementWidth))
		}

		// Verify InlineWidgetBounds determinism
		r1 := lb1.InlineWidgetBounds(rowIndex, charOffset, rh)
		r2 := lb2.InlineWidgetBounds(rowIndex, charOffset, rh)
		if r1 != r2 {
			t.Fatalf("InlineWidgetBounds(%d, %d, %d) not deterministic: lb1=%v, lb2=%v",
				rowIndex, charOffset, rh, r1, r2)
		}

		// Verify InlineWidgetBoundsWith determinism
		r1w := lb1.InlineWidgetBoundsWith(rowIndex, charOffset, rh, explicitGlyphAdvance)
		r2w := lb2.InlineWidgetBoundsWith(rowIndex, charOffset, rh, explicitGlyphAdvance)
		if r1w != r2w {
			t.Fatalf("InlineWidgetBoundsWith(%d, %d, %d, %d) not deterministic: lb1=%v, lb2=%v",
				rowIndex, charOffset, rh, explicitGlyphAdvance, r1w, r2w)
		}
	})
}

// --- From: layout_test.go ---

// Backward-compatibility unit tests verifying exact pixel equality with known
// values from existing style implementations using PaddingPct.

func TestBackwardCompatibility_PaddedStyle(t *testing.T) {
	// Style with PaddingPct that produces inset of ~8px on a 240px panel:
	// PaddingPct=3 on 240px => padInsetX = 3*240/100 = 7, padInsetY = 7
	// Use PaddingPct=4 for 240px => 4*240/100 = 9
	// For exact 0 inset, use PaddingPct=0
	hints := textlayout.TextHints{
		PixelWidth:   240,
		PixelHeight:  240,
		GlyphHeight:  7,
		GlyphAdvance: 6,
		RowHeight:    10,
	}
	cfg := layout.BridgeConfig{
		PaddingPct: 5,
	}
	lb := layout.NewLayoutBridge(hints, cfg)

	// PaddingPct=5, PixelWidth=240 => padInsetX = 5*240/100 = 12
	// PaddingPct=5, PixelHeight=240 => padInsetY = 5*240/100 = 12
	_, originY := lb.ContentOrigin()
	if originY != 12 {
		t.Fatalf("ContentOrigin().Y = %d, want 12", originY)
	}

	originX, _ := lb.ContentOrigin()
	if originX != 12 {
		t.Fatalf("ContentOrigin().X = %d, want 12", originX)
	}

	// RowY(n, 0) == 12 + n*10 for n=0..4
	for n := 0; n < 5; n++ {
		expected := 12 + n*10
		if got := lb.RowY(n, 0); got != expected {
			t.Fatalf("RowY(%d, 0) = %d, want %d", n, got, expected)
		}
	}
}

func TestBackwardCompatibility_NoPaddingStyle(t *testing.T) {
	// Style with PaddingPct=0, no title, no hint
	// Panel: 240x240, default font metrics
	hints := textlayout.TextHints{
		PixelWidth:   240,
		PixelHeight:  240,
		GlyphHeight:  7,
		GlyphAdvance: 6,
		RowHeight:    10,
	}
	cfg := layout.BridgeConfig{
		PaddingPct: 0,
	}
	lb := layout.NewLayoutBridge(hints, cfg)

	// ContentOrigin.Y should be 0 (no padding, no title)
	_, originY := lb.ContentOrigin()
	if originY != 0 {
		t.Fatalf("ContentOrigin().Y = %d, want 0", originY)
	}

	// RowY(n, 0) == n*10 for n=0..4
	for n := 0; n < 5; n++ {
		expected := n * 10
		if got := lb.RowY(n, 0); got != expected {
			t.Fatalf("RowY(%d, 0) = %d, want %d", n, got, expected)
		}
	}
}

func TestBackwardCompatibility_SerialWithTitleStyle(t *testing.T) {
	// Serial framed style: PaddingPct=0, no title bar (layout doesn't handle title)
	// Panel: 240x240, default font metrics
	hints := textlayout.TextHints{
		PixelWidth:   240,
		PixelHeight:  240,
		GlyphHeight:  7,
		GlyphAdvance: 6,
		RowHeight:    10,
	}
	cfg := layout.BridgeConfig{
		PaddingPct: 0,
	}
	lb := layout.NewLayoutBridge(hints, cfg)

	// ContentOrigin.Y should be 0 (no padding, no title in layout)
	_, originY := lb.ContentOrigin()
	if originY != 0 {
		t.Fatalf("ContentOrigin().Y = %d, want 0", originY)
	}

	// RowY(n, 0) == 0 + n*10 for n=0..4
	for n := 0; n < 5; n++ {
		expected := n * 10
		if got := lb.RowY(n, 0); got != expected {
			t.Fatalf("RowY(%d, 0) = %d, want %d", n, got, expected)
		}
	}
}

func TestBackwardCompatibility_PaddedWithTitle(t *testing.T) {
	// PaddingPct=5, no title bar in Layout Calculator
	// ContentOrigin.Y = padding inset (5*240/100 = 12)
	hints := textlayout.TextHints{
		PixelWidth:   240,
		PixelHeight:  240,
		GlyphHeight:  7,
		GlyphAdvance: 6,
		RowHeight:    10,
	}
	cfg := layout.BridgeConfig{
		PaddingPct: 5,
	}
	lb := layout.NewLayoutBridge(hints, cfg)

	// ContentOrigin.X should be 12 (padding inset: 5*240/100 = 12)
	originX, originY := lb.ContentOrigin()
	if originX != 12 {
		t.Fatalf("ContentOrigin().X = %d, want 12", originX)
	}
	// ContentOrigin.Y should be 12 (same padding inset applied vertically)
	if originY != 12 {
		t.Fatalf("ContentOrigin().Y = %d, want 12", originY)
	}

	// RowY(n, 0) == 12 + n*10 for n=0..4
	for n := 0; n < 5; n++ {
		expected := 12 + n*10
		if got := lb.RowY(n, 0); got != expected {
			t.Fatalf("RowY(%d, 0) = %d, want %d", n, got, expected)
		}
	}
}

// For any TextHints with PixelWidth and any valid BridgeConfig with PaddingPct,
// the LayoutBridge's AvailableContentWidth SHALL equal
// max(0, PixelWidth - 2 * (PaddingPct * PixelWidth / 100)).

func TestRendererPaddingProperty10_AvailableContentWidth(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		pw := rapid.IntRange(1, 800).Draw(t, "pixelWidth")
		ph := rapid.IntRange(1, 600).Draw(t, "pixelHeight")
		paddingPct := rapid.IntRange(0, 50).Draw(t, "paddingPct")

		hints := textlayout.TextHints{
			PixelWidth:   pw,
			PixelHeight:  ph,
			GlyphAdvance: 6,
			GlyphHeight:  7,
			RowHeight:    10,
		}
		cfg := layout.BridgeConfig{
			PaddingPct: paddingPct,
		}
		lb := layout.NewLayoutBridge(hints, cfg)

		padInsetX := paddingPct * pw / 100
		expected := pw - 2*padInsetX
		if expected < 0 {
			expected = 0
		}

		got := lb.AvailableContentWidth()
		if got != expected {
			t.Fatalf("AvailableContentWidth() = %d, want max(0, %d - 2*%d) = %d "+
				"(pixelWidth=%d, pixelHeight=%d, paddingPct=%d, padInsetX=%d)",
				got, pw, padInsetX, expected, pw, ph, paddingPct, padInsetX)
		}
	})
}

// For any LayoutBridge constructed from TextHints and a character count whose
// text pixel width is less than AvailableContentWidth, CenterX SHALL produce an offset such
// that offset + textPixelWidth/2 is within 1px of AvailableContentWidth/2. When text pixel
// width exceeds AvailableContentWidth, CenterX SHALL return 0.

func TestRendererPaddingProperty11_CenterXMidpointAccuracy(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		pw := rapid.IntRange(1, 800).Draw(t, "pixelWidth")
		ph := rapid.IntRange(1, 600).Draw(t, "pixelHeight")
		paddingPct := rapid.IntRange(0, 50).Draw(t, "paddingPct")
		ga := rapid.IntRange(1, 20).Draw(t, "glyphAdvance")

		hints := textlayout.TextHints{
			PixelWidth:   pw,
			PixelHeight:  ph,
			GlyphAdvance: ga,
			GlyphHeight:  7,
			RowHeight:    10,
		}
		cfg := layout.BridgeConfig{
			PaddingPct: paddingPct,
		}
		lb := layout.NewLayoutBridge(hints, cfg)

		acw := lb.AvailableContentWidth()

		// Generate char counts that produce text widths both less than and greater than ACW
		charCount := rapid.IntRange(1, 100).Draw(t, "charCount")
		textPixelWidth := lb.TextPixelWidth(charCount)
		centerX := lb.CenterX(charCount)

		if textPixelWidth >= acw {
			// When text pixel width exceeds AvailableContentWidth, CenterX SHALL return 0
			if centerX != 0 {
				t.Fatalf("CenterX(%d) = %d, want 0 (textPixelWidth=%d >= ACW=%d)",
					charCount, centerX, textPixelWidth, acw)
			}
		} else {
			// When text pixel width < ACW, verify midpoint accuracy:
			// offset + textPixelWidth/2 should be within 1px of ACW/2
			textMidpoint := centerX + textPixelWidth/2
			areaMidpoint := acw / 2

			diff := textMidpoint - areaMidpoint
			if diff < 0 {
				diff = -diff
			}
			if diff > 1 {
				t.Fatalf("CenterX(%d) midpoint inaccuracy: textMidpoint=%d (offset=%d + textWidth/2=%d), "+
					"areaMidpoint=%d (ACW=%d / 2), diff=%d > 1px "+
					"(pixelWidth=%d, pixelHeight=%d, paddingPct=%d, glyphAdvance=%d)",
					charCount, textMidpoint, centerX, textPixelWidth/2,
					areaMidpoint, acw, diff,
					pw, ph, paddingPct, ga)
			}
		}
	})
}
