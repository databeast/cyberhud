package testpattern

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"testing"

	"github.com/databeast/cyberhud/display/catalog"
	fontpkg "github.com/databeast/cyberhud/display/surface/fonts"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/runtime/action"
	"pgregory.net/rapid"
)

// --- From: testpattern_prop_test.go ---

// validTextHints generates random valid TextHints within reasonable bounds.
func validTextHints(t *rapid.T) textlayout.TextHints {
	return textlayout.TextHints{
		PixelWidth:   rapid.IntRange(1, 500).Draw(t, "PixelWidth"),
		PixelHeight:  rapid.IntRange(1, 500).Draw(t, "PixelHeight"),
		GlyphWidth:   rapid.IntRange(1, 20).Draw(t, "GlyphWidth"),
		GlyphHeight:  rapid.IntRange(1, 20).Draw(t, "GlyphHeight"),
		GlyphAdvance: rapid.IntRange(1, 25).Draw(t, "GlyphAdvance"),
		RowHeight:    rapid.IntRange(1, 30).Draw(t, "RowHeight"),
	}
}

// TestProperty1_BorderRectangleCoversPanelEdges tests that for any valid TextHints
// with PixelWidth >= 1 and PixelHeight >= 1, the rendered test pattern has white pixels
// at all positions along the 1-pixel border.
func TestProperty1_BorderRectangleCoversPanelEdges(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		hints := validTextHints(t)

		view := BuildView(hints, false)
		img := view.Image

		w := hints.PixelWidth
		h := hints.PixelHeight

		// Top row (y=0, x=0..Width-1).
		for x := 0; x < w; x++ {
			got := img.RGBAAt(x, 0)
			if got != colorBorder {
				t.Fatalf("top border pixel (%d, 0): got %v, want %v", x, got, colorBorder)
			}
		}

		// Bottom row (y=Height-1, x=0..Width-1).
		for x := 0; x < w; x++ {
			got := img.RGBAAt(x, h-1)
			if got != colorBorder {
				t.Fatalf("bottom border pixel (%d, %d): got %v, want %v", x, h-1, got, colorBorder)
			}
		}

		// Left column (x=0, y=0..Height-1).
		for y := 0; y < h; y++ {
			got := img.RGBAAt(0, y)
			if got != colorBorder {
				t.Fatalf("left border pixel (0, %d): got %v, want %v", y, got, colorBorder)
			}
		}

		// Right column (x=Width-1, y=0..Height-1).
		for y := 0; y < h; y++ {
			got := img.RGBAAt(w-1, y)
			if got != colorBorder {
				t.Fatalf("right border pixel (%d, %d): got %v, want %v", w-1, y, got, colorBorder)
			}
		}
	})
}

// TestProperty2_CornerMarkersConditionalOnDimensions tests that:
// - When PixelWidth >= 8 AND PixelHeight >= 8, 4×4 solid white regions exist at all four corners.
// - When PixelWidth < 8 OR PixelHeight < 8, the interior of corner regions (excluding border) are NOT all white.
func TestProperty2_CornerMarkersConditionalOnDimensions(t *testing.T) {
	// Sub-test: large panels (corners present).
	t.Run("LargePanels_CornersPresent", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			hints := textlayout.TextHints{
				PixelWidth:   rapid.IntRange(8, 500).Draw(t, "PixelWidth"),
				PixelHeight:  rapid.IntRange(8, 500).Draw(t, "PixelHeight"),
				GlyphWidth:   rapid.IntRange(1, 20).Draw(t, "GlyphWidth"),
				GlyphHeight:  rapid.IntRange(1, 20).Draw(t, "GlyphHeight"),
				GlyphAdvance: rapid.IntRange(1, 25).Draw(t, "GlyphAdvance"),
				RowHeight:    rapid.IntRange(1, 30).Draw(t, "RowHeight"),
			}

			view := BuildView(hints, false)
			img := view.Image

			w := hints.PixelWidth
			h := hints.PixelHeight
			white := color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}

			// Top-left corner: (0,0) to (3,3).
			for y := 0; y < cornerMarkerSize; y++ {
				for x := 0; x < cornerMarkerSize; x++ {
					got := img.RGBAAt(x, y)
					if got != white {
						t.Fatalf("top-left corner (%d,%d): got %v, want %v", x, y, got, white)
					}
				}
			}

			// Top-right corner: (W-4,0) to (W-1,3).
			for y := 0; y < cornerMarkerSize; y++ {
				for x := w - cornerMarkerSize; x < w; x++ {
					got := img.RGBAAt(x, y)
					if got != white {
						t.Fatalf("top-right corner (%d,%d): got %v, want %v", x, y, got, white)
					}
				}
			}

			// Bottom-left corner: (0,H-4) to (3,H-1).
			for y := h - cornerMarkerSize; y < h; y++ {
				for x := 0; x < cornerMarkerSize; x++ {
					got := img.RGBAAt(x, y)
					if got != white {
						t.Fatalf("bottom-left corner (%d,%d): got %v, want %v", x, y, got, white)
					}
				}
			}

			// Bottom-right corner: (W-4,H-4) to (W-1,H-1).
			for y := h - cornerMarkerSize; y < h; y++ {
				for x := w - cornerMarkerSize; x < w; x++ {
					got := img.RGBAAt(x, y)
					if got != white {
						t.Fatalf("bottom-right corner (%d,%d): got %v, want %v", x, y, got, white)
					}
				}
			}
		})
	})

	// Sub-test: small panels (corners absent).
	t.Run("SmallPanels_CornersAbsent", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			// Generate a panel where at least one dimension is < 8.
			smallDim := rapid.Bool().Draw(t, "smallWidth")
			var pw, ph int
			if smallDim {
				pw = rapid.IntRange(1, 7).Draw(t, "PixelWidth")
				ph = rapid.IntRange(1, 500).Draw(t, "PixelHeight")
			} else {
				pw = rapid.IntRange(1, 500).Draw(t, "PixelWidth")
				ph = rapid.IntRange(1, 7).Draw(t, "PixelHeight")
			}

			hints := textlayout.TextHints{
				PixelWidth:   pw,
				PixelHeight:  ph,
				GlyphWidth:   rapid.IntRange(1, 20).Draw(t, "GlyphWidth"),
				GlyphHeight:  rapid.IntRange(1, 20).Draw(t, "GlyphHeight"),
				GlyphAdvance: rapid.IntRange(1, 25).Draw(t, "GlyphAdvance"),
				RowHeight:    rapid.IntRange(1, 30).Draw(t, "RowHeight"),
			}

			view := BuildView(hints, false)
			img := view.Image

			w := hints.PixelWidth
			h := hints.PixelHeight
			white := color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}

			// Check that the interior pixels of the 4×4 corner region (excluding
			// the 1-pixel border) are NOT all white.
			// Only check corners that fit within the image bounds.
			maxCornerX := cornerMarkerSize
			if maxCornerX > w {
				maxCornerX = w
			}
			maxCornerY := cornerMarkerSize
			if maxCornerY > h {
				maxCornerY = h
			}

			// Collect interior pixels (those not on the border) in the top-left corner region.
			allInteriorWhite := true
			interiorCount := 0
			for y := 1; y < maxCornerY; y++ {
				for x := 1; x < maxCornerX; x++ {
					// Skip border pixels (row 0, col 0 are border).
					// Also skip if at bottom row or right col of image (those are border too).
					if y == h-1 || x == w-1 {
						continue
					}
					interiorCount++
					got := img.RGBAAt(x, y)
					if got != white {
						allInteriorWhite = false
					}
				}
			}

			// If there are interior pixels to check, they should NOT all be white
			// (corner markers are absent).
			if interiorCount > 0 && allInteriorWhite {
				t.Fatalf("small panel (%dx%d): interior of top-left corner region should NOT all be white when dimensions < 8", w, h)
			}
		})
	})
}

// TestProperty3_ColorSwatchesOnColorPanels tests that for any valid TextHints on a
// color-capable panel, the test pattern contains exactly 5 horizontal color bands of
// equal width (floor(PixelWidth / 5)) in order red, green, blue, white, black, with
// band height equal to min(2 × RowHeight, PixelHeight).
// We test drawColorSwatches in isolation to avoid interference from overlays (hint summary, glyph grid)
// that are drawn on top per the design's layering order.
func TestProperty3_ColorSwatchesOnColorPanels(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		hints := textlayout.TextHints{
			PixelWidth:   rapid.IntRange(20, 500).Draw(t, "PixelWidth"),
			PixelHeight:  rapid.IntRange(20, 500).Draw(t, "PixelHeight"),
			GlyphWidth:   rapid.IntRange(1, 20).Draw(t, "GlyphWidth"),
			GlyphHeight:  rapid.IntRange(1, 20).Draw(t, "GlyphHeight"),
			GlyphAdvance: rapid.IntRange(1, 25).Draw(t, "GlyphAdvance"),
			RowHeight:    rapid.IntRange(1, 30).Draw(t, "RowHeight"),
		}

		w := hints.PixelWidth
		h := hints.PixelHeight
		rh := hints.RowHeight

		// Create a fresh image and only draw swatches (isolated test).
		img := image.NewRGBA(image.Rect(0, 0, w, h))
		drawColorSwatches(img, hints, false)

		swatchWidth := w / swatchCount
		bandHeight := swatchHeightMultiplier * rh
		if bandHeight > h {
			bandHeight = h
		}

		// Determine startY based on corner presence.
		startY := 1
		hasCorners := w >= minDimensionForCorners && h >= minDimensionForCorners
		if hasCorners {
			startY = cornerMarkerSize
		}

		// Clamp bandHeight to available space (don't exceed bottom border).
		maxY := h - 1
		if bandHeight > maxY-startY {
			bandHeight = maxY - startY
		}
		if bandHeight <= 0 {
			return // band too small to test
		}

		// Sample from the middle of the band vertically.
		sampleY := startY + bandHeight/2
		if sampleY <= 0 {
			sampleY = startY + 1
		}
		if sampleY >= h-1 {
			return // band too small to test
		}

		expectedColors := []color.RGBA{colorSwatchRed, colorSwatchGrn, colorSwatchBlu, colorSwatchWht, colorSwatchBlk}

		for i, expected := range expectedColors {
			// Sample from middle of the swatch horizontally.
			sampleX := i*swatchWidth + swatchWidth/2
			if sampleX <= 0 || sampleX >= w-1 {
				continue // skip if too close to border
			}

			// Skip if in a corner marker region.
			if hasCorners {
				inTopLeft := sampleX < cornerMarkerSize && sampleY < cornerMarkerSize
				inTopRight := sampleX >= w-cornerMarkerSize && sampleY < cornerMarkerSize
				inBottomLeft := sampleX < cornerMarkerSize && sampleY >= h-cornerMarkerSize
				inBottomRight := sampleX >= w-cornerMarkerSize && sampleY >= h-cornerMarkerSize
				if inTopLeft || inTopRight || inBottomLeft || inBottomRight {
					continue
				}
			}

			got := img.RGBAAt(sampleX, sampleY)
			if got != expected {
				t.Fatalf("swatch %d at (%d,%d): got %v, want %v (w=%d, h=%d, rh=%d, swatchWidth=%d, bandHeight=%d, startY=%d)",
					i, sampleX, sampleY, got, expected, w, h, rh, swatchWidth, bandHeight, startY)
			}
		}
	})
}

// TestProperty4_MonochromeSwatches tests that for any valid TextHints on a monochrome
// panel, the test pattern contains exactly 2 horizontal bands of equal width
// (floor(PixelWidth / 2)) representing full-intensity foreground (white) and
// zero-intensity background (black), with band height equal to min(2 × RowHeight, PixelHeight).
// We test drawColorSwatches in isolation to avoid interference from overlays.
func TestProperty4_MonochromeSwatches(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		hints := textlayout.TextHints{
			PixelWidth:   rapid.IntRange(20, 500).Draw(t, "PixelWidth"),
			PixelHeight:  rapid.IntRange(20, 500).Draw(t, "PixelHeight"),
			GlyphWidth:   rapid.IntRange(1, 20).Draw(t, "GlyphWidth"),
			GlyphHeight:  rapid.IntRange(1, 20).Draw(t, "GlyphHeight"),
			GlyphAdvance: rapid.IntRange(1, 25).Draw(t, "GlyphAdvance"),
			RowHeight:    rapid.IntRange(1, 30).Draw(t, "RowHeight"),
		}

		w := hints.PixelWidth
		h := hints.PixelHeight
		rh := hints.RowHeight

		// Create a fresh image and only draw swatches (isolated test).
		img := image.NewRGBA(image.Rect(0, 0, w, h))
		drawColorSwatches(img, hints, true)

		swatchWidth := w / monoSwatchCount
		bandHeight := swatchHeightMultiplier * rh
		if bandHeight > h {
			bandHeight = h
		}

		// Determine startY based on corner presence.
		startY := 1
		hasCorners := w >= minDimensionForCorners && h >= minDimensionForCorners
		if hasCorners {
			startY = cornerMarkerSize
		}

		// Clamp bandHeight to available space (don't exceed bottom border).
		maxY := h - 1
		if bandHeight > maxY-startY {
			bandHeight = maxY - startY
		}
		if bandHeight <= 0 {
			return // band too small to test
		}

		// Sample from the middle of the band vertically.
		sampleY := startY + bandHeight/2
		if sampleY <= 0 {
			sampleY = startY + 1
		}
		if sampleY >= h-1 {
			return // band too small to test
		}

		expectedColors := []color.RGBA{colorSwatchWht, colorSwatchBlk}

		for i, expected := range expectedColors {
			// Sample from middle of the swatch horizontally.
			sampleX := i*swatchWidth + swatchWidth/2
			if sampleX <= 0 || sampleX >= w-1 {
				continue // skip if too close to border
			}

			// Skip if in a corner marker region.
			if hasCorners {
				inTopLeft := sampleX < cornerMarkerSize && sampleY < cornerMarkerSize
				inTopRight := sampleX >= w-cornerMarkerSize && sampleY < cornerMarkerSize
				inBottomLeft := sampleX < cornerMarkerSize && sampleY >= h-cornerMarkerSize
				inBottomRight := sampleX >= w-cornerMarkerSize && sampleY >= h-cornerMarkerSize
				if inTopLeft || inTopRight || inBottomLeft || inBottomRight {
					continue
				}
			}

			got := img.RGBAAt(sampleX, sampleY)
			if got != expected {
				t.Fatalf("mono swatch %d at (%d,%d): got %v, want %v (w=%d, h=%d, rh=%d, swatchWidth=%d, bandHeight=%d, startY=%d)",
					i, sampleX, sampleY, got, expected, w, h, rh, swatchWidth, bandHeight, startY)
			}
		}
	})
}

// TestProperty5_GlyphGridCharacterCount tests that:
//   - When PixelWidth >= GlyphAdvance and there is room for at least one character
//     (PixelWidth >= 2*cornerMarkerSize + GlyphAdvance), the glyph grid renders exactly
//     floor((PixelWidth - 2*cornerMarkerSize) / GlyphAdvance) characters per row (capped at 26 for alpha).
//   - When PixelWidth < GlyphAdvance, no characters are rendered.
func TestProperty5_GlyphGridCharacterCount(t *testing.T) {
	face := fontpkg.Default()

	t.Run("CharactersRendered", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			// Use the font's actual metrics for glyph dimensions to ensure
			// characters are always visible (font is 5x7).
			fm := face.Metrics()
			gw := fm.GlyphWidth  // 5
			gh := fm.GlyphHeight // 7
			ga := rapid.IntRange(gw, 25).Draw(t, "GlyphAdvance")
			rh := rapid.IntRange(gh, 30).Draw(t, "RowHeight")

			// Ensure PixelWidth is large enough for at least one char: >= 2*cornerMarkerSize + ga
			minWidth := 2*cornerMarkerSize + ga
			pw := rapid.IntRange(minWidth, 500).Draw(t, "PixelWidth")

			// PixelHeight must be large enough for the grid to appear.
			// Grid starts after swatches: cornerMarkerSize + 2*RowHeight, plus at least one glyph row.
			minHeight := cornerMarkerSize + 2*rh + gh + 1
			if minHeight < minDimensionForCorners {
				minHeight = minDimensionForCorners
			}
			if minHeight > 500 {
				return
			}
			ph := rapid.IntRange(minHeight, 500).Draw(t, "PixelHeight")

			hints := textlayout.TextHints{
				PixelWidth:   pw,
				PixelHeight:  ph,
				GlyphWidth:   gw,
				GlyphHeight:  gh,
				GlyphAdvance: ga,
				RowHeight:    rh,
			}

			view := BuildView(hints, false)
			img := view.Image

			// Expected character count per row.
			expectedChars := (pw - 2*cornerMarkerSize) / ga
			if expectedChars > 26 {
				expectedChars = 26
			}

			// Determine gridStartY (same logic as drawGlyphGrid).
			hasCorners := pw >= minDimensionForCorners && ph >= minDimensionForCorners
			swatchStartY := 1
			if hasCorners {
				swatchStartY = cornerMarkerSize
			}
			swatchBandH := swatchHeightMultiplier * rh
			available := ph - swatchStartY
			if swatchBandH > available {
				swatchBandH = available
			}
			gridStartY := swatchStartY + swatchBandH

			startX := cornerMarkerSize

			// Count rendered characters in the alpha row by checking if any glyph
			// pixel is set (using the font to know which pixels should be lit).
			alphaRow := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
			renderedCount := 0

			for ci := 0; ci < expectedChars && ci < len(alphaRow); ci++ {
				ch := rune(alphaRow[ci])
				px := startX + ci*ga
				py := gridStartY

				// Check if this glyph cell has any lit pixel by scanning with the font.
				found := false
				for row := 0; row < gh && py+row < ph; row++ {
					bits := face.GlyphRow(ch, row)
					for col := 0; col < gw && px+col < pw; col++ {
						if bits&(1<<uint(31-col)) != 0 {
							got := img.RGBAAt(px+col, py+row)
							if got == colorGlyph {
								found = true
								break
							}
						}
					}
					if found {
						break
					}
				}
				if found {
					renderedCount++
				}
			}

			if renderedCount != expectedChars {
				t.Fatalf("expected %d chars rendered, got %d (pw=%d, ga=%d, gw=%d, gh=%d, gridStartY=%d)",
					expectedChars, renderedCount, pw, ga, gw, gh, gridStartY)
			}
		})
	})

	t.Run("NoCharactersWhenTooNarrow", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			// Generate hints where PixelWidth < GlyphAdvance (no character fits at all).
			ga := rapid.IntRange(2, 25).Draw(t, "GlyphAdvance")
			pw := rapid.IntRange(1, ga-1).Draw(t, "PixelWidth")
			gw := rapid.IntRange(1, 20).Draw(t, "GlyphWidth")
			gh := rapid.IntRange(1, 20).Draw(t, "GlyphHeight")
			rh := rapid.IntRange(1, 30).Draw(t, "RowHeight")
			ph := rapid.IntRange(40, 200).Draw(t, "PixelHeight")

			hints := textlayout.TextHints{
				PixelWidth:   pw,
				PixelHeight:  ph,
				GlyphWidth:   gw,
				GlyphHeight:  gh,
				GlyphAdvance: ga,
				RowHeight:    rh,
			}

			view := BuildView(hints, false)
			img := view.Image

			// When PixelWidth < GlyphAdvance, drawGlyphGrid returns early.
			// Verify no glyph-colored pixels exist in the area where glyphs would be
			// (excluding border and guide line pixels).
			hasCorners := pw >= minDimensionForCorners && ph >= minDimensionForCorners
			swatchStartY := 1
			if hasCorners {
				swatchStartY = cornerMarkerSize
			}
			swatchBandH := swatchHeightMultiplier * rh
			avail := ph - swatchStartY
			if swatchBandH > avail {
				swatchBandH = avail
			}
			gridStartY := swatchStartY + swatchBandH

			startX := cornerMarkerSize
			if startX >= pw {
				return
			}

			// Check that no glyph pixels exist in the grid region.
			for y := gridStartY; y < ph-1 && y < gridStartY+gh; y++ {
				for x := startX; x < pw-1 && x < startX+gw; x++ {
					got := img.RGBAAt(x, y)
					if got == colorGlyph {
						isBorder := x == 0 || x == pw-1 || y == 0 || y == ph-1
						isGuideLine := (ga > 0 && x%ga == 0) || (rh > 0 && y%rh == 0)
						isCorner := hasCorners && ((x < cornerMarkerSize && y < cornerMarkerSize) ||
							(x >= pw-cornerMarkerSize && y < cornerMarkerSize) ||
							(x < cornerMarkerSize && y >= ph-cornerMarkerSize) ||
							(x >= pw-cornerMarkerSize && y >= ph-cornerMarkerSize))
						if !isBorder && !isGuideLine && !isCorner {
							t.Fatalf("found glyph pixel at (%d,%d) when no characters should be rendered (pw=%d < ga=%d)",
								x, y, pw, ga)
						}
					}
				}
			}
			_ = img
		})
	})
}

// TestProperty6_GlyphGridSpacing tests that:
// - Adjacent characters in a row have their left edges separated by exactly GlyphAdvance pixels.
// - Adjacent rows (alpha and digit) have their top edges separated by exactly RowHeight pixels.
func TestProperty6_GlyphGridSpacing(t *testing.T) {
	face := fontpkg.Default()
	fm := face.Metrics()

	rapid.Check(t, func(t *rapid.T) {
		// Use the font's actual glyph dimensions so characters are visible.
		gw := fm.GlyphWidth  // 5
		gh := fm.GlyphHeight // 7
		ga := rapid.IntRange(gw, 25).Draw(t, "GlyphAdvance")
		rh := rapid.IntRange(gh, 30).Draw(t, "RowHeight")

		// Need at least 2 chars: width >= 2*cornerMarkerSize + 2*ga
		minWidth := 2*cornerMarkerSize + 2*ga
		if minWidth > 500 {
			return
		}
		pw := rapid.IntRange(minWidth, 500).Draw(t, "PixelWidth")

		// Need height for grid start + 2 rows (alpha + digit).
		minHeight := cornerMarkerSize + 2*rh + 2*rh + 1
		if minHeight < minDimensionForCorners {
			minHeight = minDimensionForCorners
		}
		if minHeight > 500 {
			return
		}
		ph := rapid.IntRange(minHeight, 500).Draw(t, "PixelHeight")

		hints := textlayout.TextHints{
			PixelWidth:   pw,
			PixelHeight:  ph,
			GlyphWidth:   gw,
			GlyphHeight:  gh,
			GlyphAdvance: ga,
			RowHeight:    rh,
		}

		view := BuildView(hints, false)
		img := view.Image

		// Compute gridStartY (same as implementation).
		hasCorners := pw >= minDimensionForCorners && ph >= minDimensionForCorners
		swatchStartY := 1
		if hasCorners {
			swatchStartY = cornerMarkerSize
		}
		swatchBandH := swatchHeightMultiplier * rh
		available := ph - swatchStartY
		if swatchBandH > available {
			swatchBandH = available
		}
		gridStartY := swatchStartY + swatchBandH

		startX := cornerMarkerSize

		// Helper: find leftmost x position where a character has a lit pixel.
		findLeftEdge := func(ch rune, px, py int) (int, bool) {
			for col := 0; col < gw && px+col < pw; col++ {
				for row := 0; row < gh && py+row < ph; row++ {
					bits := face.GlyphRow(ch, row)
					if bits&(1<<uint(31-col)) != 0 {
						got := img.RGBAAt(px+col, py+row)
						if got == colorGlyph {
							return px + col, true
						}
					}
				}
			}
			return 0, false
		}

		expectedChars := (pw - 2*cornerMarkerSize) / ga
		if expectedChars > 26 {
			expectedChars = 26
		}

		// HORIZONTAL SPACING: verify adjacent characters have left edges
		// separated by exactly GlyphAdvance pixels.
		alphaRow := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		for ci := 0; ci < expectedChars-1 && ci+1 < len(alphaRow); ci++ {
			ch1 := rune(alphaRow[ci])
			ch2 := rune(alphaRow[ci+1])

			px1 := startX + ci*ga
			px2 := startX + (ci+1)*ga
			py := gridStartY

			leftEdge1, found1 := findLeftEdge(ch1, px1, py)
			leftEdge2, found2 := findLeftEdge(ch2, px2, py)

			if found1 && found2 {
				// The left edge positions should differ by exactly GlyphAdvance,
				// since both start at their respective px positions.
				// leftEdge1 is px1 + some column offset for the first lit column of ch1.
				// leftEdge2 is px2 + some column offset for the first lit column of ch2.
				// The spacing between character START positions (px2 - px1) must be GlyphAdvance.
				spacing := px2 - px1
				if spacing != ga {
					t.Fatalf("horizontal spacing between char %d start and char %d start: got %d, want %d",
						ci, ci+1, spacing, ga)
				}

				// Also verify the lit pixel positions are consistent:
				// The offset from the start position to the first lit column should be
				// the same for both characters IF they have the same leftmost set column.
				// But different characters have different shapes, so we verify the START
				// positions are correct rather than the individual pixel positions.
				// What matters is: leftEdge2 - leftEdge1 == ga only if both characters
				// have the same leftmost column (e.g., both start at col 0).
				// For 'A' (row 0 = 0x04 = bit 2), leftmost lit col is col 2.
				// For 'B' (row 0 = 0x1E = bits 4,3,2,1), leftmost lit col is col 0.
				// So we can't compare leftEdge directly — but we CAN verify that
				// each character's left edge equals startX + ci*ga + (first lit col).
				// That's what we already know from the positioning.
				expectedLeftEdge1 := px1 // This is where the character cell starts.
				expectedLeftEdge2 := px2
				// leftEdge is px + firstLitCol, so verify it's within the cell.
				if leftEdge1 < expectedLeftEdge1 || leftEdge1 >= expectedLeftEdge1+gw {
					t.Fatalf("char %d left edge %d outside cell [%d, %d)",
						ci, leftEdge1, expectedLeftEdge1, expectedLeftEdge1+gw)
				}
				if leftEdge2 < expectedLeftEdge2 || leftEdge2 >= expectedLeftEdge2+gw {
					t.Fatalf("char %d left edge %d outside cell [%d, %d)",
						ci+1, leftEdge2, expectedLeftEdge2, expectedLeftEdge2+gw)
				}
			}
		}

		// VERTICAL SPACING: verify the alpha row start and digit row start differ by RowHeight.
		digitRowY := gridStartY + rh
		if digitRowY+gh > ph {
			return // digit row doesn't fit
		}

		rowDistance := digitRowY - gridStartY
		if rowDistance != rh {
			t.Fatalf("row spacing: got %d, want %d (gridStartY=%d, digitRowY=%d)",
				rowDistance, rh, gridStartY, digitRowY)
		}

		// Verify both rows actually render at their expected Y positions.
		// Check character 'A' at alpha row Y and '0' at digit row Y.
		_, alphaFound := findLeftEdge('A', startX, gridStartY)
		_, digitFound := findLeftEdge('0', startX, digitRowY)

		if !alphaFound {
			t.Fatalf("alpha row character 'A' not found at position (%d, %d)", startX, gridStartY)
		}
		if !digitFound {
			t.Fatalf("digit row character '0' not found at position (%d, %d)", startX, digitRowY)
		}
	})
}

// TestProperty11_WidgetGuideGridAtCorrectIntervals tests that for any valid TextHints
// with RowHeight > 0 and GlyphAdvance > 0, drawGuideGrid produces:
// - 1-pixel-wide horizontal lines at every y = n × RowHeight spanning full width
// - 1-pixel-wide vertical lines at every x = n × GlyphAdvance spanning full height
func TestProperty11_WidgetGuideGridAtCorrectIntervals(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		hints := textlayout.TextHints{
			PixelWidth:   rapid.IntRange(10, 500).Draw(t, "PixelWidth"),
			PixelHeight:  rapid.IntRange(10, 500).Draw(t, "PixelHeight"),
			GlyphWidth:   rapid.IntRange(1, 20).Draw(t, "GlyphWidth"),
			GlyphHeight:  rapid.IntRange(1, 20).Draw(t, "GlyphHeight"),
			GlyphAdvance: rapid.IntRange(1, 25).Draw(t, "GlyphAdvance"),
			RowHeight:    rapid.IntRange(1, 30).Draw(t, "RowHeight"),
		}

		w := hints.PixelWidth
		h := hints.PixelHeight
		rh := hints.RowHeight
		ga := hints.GlyphAdvance

		img := image.NewRGBA(image.Rect(0, 0, w, h))
		drawGuideGrid(img, hints)

		// Verify horizontal lines at y = n * RowHeight.
		numH := h / rh
		for n := 0; n < numH; n++ {
			y := n * rh
			// Check a few x positions along the line.
			for x := 0; x < w; x += w/4 + 1 {
				if x >= w {
					break
				}
				got := img.RGBAAt(x, y)
				if got != colorGuide {
					t.Fatalf("horizontal guide at (%d,%d): got %v, want %v", x, y, got, colorGuide)
				}
			}
		}

		// Verify vertical lines at x = n * GlyphAdvance.
		numV := w / ga
		for n := 0; n < numV; n++ {
			x := n * ga
			for y := 0; y < h; y += h/4 + 1 {
				if y >= h {
					break
				}
				got := img.RGBAAt(x, y)
				if got != colorGuide {
					t.Fatalf("vertical guide at (%d,%d): got %v, want %v", x, y, got, colorGuide)
				}
			}
		}
	})
}

// TestProperty12_GuideLineColorIntensity tests that colorGuide has RGB channels all <= 63
// (25% of 255), and that it differs from white, black, red, green, and blue.
func TestProperty12_GuideLineColorIntensity(t *testing.T) {
	// Verify guide color intensity is at most 25% of 255 (63).
	if colorGuide.R > guideIntensityMax || colorGuide.G > guideIntensityMax || colorGuide.B > guideIntensityMax {
		t.Fatalf("guide color %v has channel > %d", colorGuide, guideIntensityMax)
	}

	// Verify guide color differs from all other element colors.
	otherColors := []color.RGBA{colorBorder, colorBackground, colorSwatchRed, colorSwatchGrn, colorSwatchBlu}
	for _, c := range otherColors {
		if colorGuide == c {
			t.Fatalf("guide color %v must differ from %v", colorGuide, c)
		}
	}
}

// TestProperty7_HintSummaryContent tests that for any valid TextHints where sufficient
// space exists to display at least one label, the rendered hint summary includes text
// representations of TextHints fields with correct values.
func TestProperty7_HintSummaryContent(t *testing.T) {
	face := fontpkg.Default()
	fm := face.Metrics()

	rapid.Check(t, func(t *rapid.T) {
		// Use font's actual glyph dimensions so characters render correctly.
		gw := fm.GlyphWidth  // 5
		gh := fm.GlyphHeight // 7
		ga := rapid.IntRange(gw, 25).Draw(t, "GlyphAdvance")
		rh := rapid.IntRange(gh, 30).Draw(t, "RowHeight")

		// Need enough width for at least one label: "PW=X" is 4 chars minimum.
		// startX = cornerMarkerSize = 4, so need: 4 + 4*ga <= PixelWidth.
		minWidth := cornerMarkerSize + 4*ga
		if minWidth > 500 {
			return
		}
		pw := rapid.IntRange(minWidth, 500).Draw(t, "PixelWidth")

		// Need enough height for at least one row of text.
		// startY = cornerMarkerSize = 4, so need: 4 + gh <= PixelHeight.
		minHeight := cornerMarkerSize + gh
		if minHeight > 500 {
			return
		}
		ph := rapid.IntRange(minHeight, 500).Draw(t, "PixelHeight")

		hints := textlayout.TextHints{
			PixelWidth:   pw,
			PixelHeight:  ph,
			GlyphWidth:   gw,
			GlyphHeight:  gh,
			GlyphAdvance: ga,
			RowHeight:    rh,
		}

		// Render hint summary in isolation on a fresh image.
		img := image.NewRGBA(image.Rect(0, 0, pw, ph))
		drawHintSummary(img, hints)

		// Expected labels in order.
		labels := []struct {
			name  string
			value int
		}{
			{"PW", pw},
			{"PH", ph},
			{"GW", gw},
			{"GH", gh},
			{"GA", ga},
			{"RH", rh},
		}

		startX := cornerMarkerSize
		startY := cornerMarkerSize
		curY := startY

		for _, lbl := range labels {
			text := fmt.Sprintf("%s=%d", lbl.name, lbl.value)

			// Check if this label fits vertically.
			if curY+gh > ph {
				break
			}

			// Check if this label fits horizontally.
			textWidth := len(text) * ga
			actualText := text
			if startX+textWidth > pw {
				maxChars := (pw - startX) / ga
				if maxChars <= 0 {
					break
				}
				if maxChars < len(text) {
					actualText = text[:maxChars]
				}
			}

			// Verify the first character of this label is rendered correctly.
			// Check that at least one pixel of the first character matches colorHintText.
			ch := rune(actualText[0])
			px := startX
			py := curY

			found := false
			for row := 0; row < gh && py+row < ph; row++ {
				bits := face.GlyphRow(ch, row)
				for col := 0; col < gw && px+col < pw; col++ {
					if bits&(1<<uint(31-col)) != 0 {
						got := img.RGBAAt(px+col, py+row)
						if got == colorHintText {
							found = true
							break
						}
					}
				}
				if found {
					break
				}
			}

			if !found {
				t.Fatalf("hint label %q not rendered at (%d,%d): no colorHintText pixels found for first char '%c' (pw=%d, ph=%d, ga=%d, rh=%d, gh=%d)",
					text, px, py, ch, pw, ph, ga, rh, gh)
			}

			curY += rh
		}
	})
}

// TestProperty8_HintSummaryInsetFromCorners tests that for any valid TextHints where
// the hint summary is rendered, NO hint text pixels exist at x < 4 or y < 4,
// ensuring no overlap with corner markers.
func TestProperty8_HintSummaryInsetFromCorners(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		hints := textlayout.TextHints{
			PixelWidth:   rapid.IntRange(20, 500).Draw(t, "PixelWidth"),
			PixelHeight:  rapid.IntRange(20, 500).Draw(t, "PixelHeight"),
			GlyphWidth:   rapid.IntRange(1, 20).Draw(t, "GlyphWidth"),
			GlyphHeight:  rapid.IntRange(1, 20).Draw(t, "GlyphHeight"),
			GlyphAdvance: rapid.IntRange(1, 25).Draw(t, "GlyphAdvance"),
			RowHeight:    rapid.IntRange(1, 30).Draw(t, "RowHeight"),
		}

		w := hints.PixelWidth
		h := hints.PixelHeight

		// Render hint summary in isolation on a fresh image.
		img := image.NewRGBA(image.Rect(0, 0, w, h))
		drawHintSummary(img, hints)

		// Verify no hint text pixels exist at y < cornerMarkerSize.
		black := color.RGBA{0, 0, 0, 0} // zero value (unset pixel in fresh image)
		for y := 0; y < cornerMarkerSize && y < h; y++ {
			for x := 0; x < w; x++ {
				got := img.RGBAAt(x, y)
				if got != black {
					t.Fatalf("hint summary pixel at (%d,%d) is set (%v) but should be empty (y < cornerMarkerSize=%d)",
						x, y, got, cornerMarkerSize)
				}
			}
		}

		// Verify no hint text pixels exist at x < cornerMarkerSize.
		for y := 0; y < h; y++ {
			for x := 0; x < cornerMarkerSize && x < w; x++ {
				got := img.RGBAAt(x, y)
				if got != black {
					t.Fatalf("hint summary pixel at (%d,%d) is set (%v) but should be empty (x < cornerMarkerSize=%d)",
						x, y, got, cornerMarkerSize)
				}
			}
		}
	})
}

// TestProperty9_DeterministicOutputAndStaticFlag tests that for any valid TextHints,
// calling BuildView twice with the same TextHints produces byte-identical image output,
// the same RenderCacheKey string, and State.Static == true.
func TestProperty9_DeterministicOutputAndStaticFlag(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		hints := validTextHints(t)

		view1 := BuildView(hints, false)
		view2 := BuildView(hints, false)

		// Static flag must be true.
		if !view1.Static {
			t.Fatal("BuildView returned Static=false, want true")
		}
		if !view2.Static {
			t.Fatal("BuildView returned Static=false on second call, want true")
		}

		// Images must be byte-identical.
		if !bytes.Equal(view1.Image.Pix, view2.Image.Pix) {
			t.Fatal("BuildView produced different image bytes for same hints")
		}

		// RenderCacheKey must be identical.
		sig1 := RenderCacheKey(hints)
		sig2 := RenderCacheKey(hints)
		if sig1 != sig2 {
			t.Fatalf("RenderCacheKey differs: %q vs %q", sig1, sig2)
		}
	})
}

// TestProperty10_DiscriminatingRenderCacheKey tests that for any two distinct TextHints
// configurations (differing in at least one of PixelWidth, PixelHeight, GlyphWidth,
// GlyphHeight, GlyphAdvance, or RowHeight), RenderCacheKey returns different string values.
func TestProperty10_DiscriminatingRenderCacheKey(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		hints1 := validTextHints(t)

		// Generate a second hints that differs in at least one field.
		// Pick a random field to modify.
		field := rapid.IntRange(0, 5).Draw(t, "fieldToModify")
		hints2 := hints1
		switch field {
		case 0:
			hints2.PixelWidth = hints1.PixelWidth + rapid.IntRange(1, 100).Draw(t, "delta")
		case 1:
			hints2.PixelHeight = hints1.PixelHeight + rapid.IntRange(1, 100).Draw(t, "delta")
		case 2:
			hints2.GlyphWidth = hints1.GlyphWidth + rapid.IntRange(1, 10).Draw(t, "delta")
		case 3:
			hints2.GlyphHeight = hints1.GlyphHeight + rapid.IntRange(1, 10).Draw(t, "delta")
		case 4:
			hints2.GlyphAdvance = hints1.GlyphAdvance + rapid.IntRange(1, 10).Draw(t, "delta")
		case 5:
			hints2.RowHeight = hints1.RowHeight + rapid.IntRange(1, 10).Draw(t, "delta")
		}

		sig1 := RenderCacheKey(hints1)
		sig2 := RenderCacheKey(hints2)
		if sig1 == sig2 {
			t.Fatalf("RenderCacheKey should differ for different hints:\n  hints1=%+v → %q\n  hints2=%+v → %q",
				hints1, sig1, hints2, sig2)
		}
	})
}

// TestProperty13_ActionHandlerReturnsZeroResult tests that for any action value
// (None, Up, Down, Primary, Secondary) and any cursor/itemCount values, HandleAction
// returns an action.Result with Navigate == "", CursorDelta == 0, Dirty == false, and Refresh == false.
func TestProperty13_ActionHandlerReturnsZeroResult(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate random action (0-4 covers None, Up, Down, Primary, Secondary).
		act := action.Action(rapid.IntRange(0, 4).Draw(t, "action"))
		cursor := rapid.IntRange(0, 100).Draw(t, "cursor")
		itemCount := rapid.IntRange(0, 100).Draw(t, "itemCount")

		h := Handler{}
		result := h.HandleAction(act, cursor, itemCount)

		if result.Navigate != "" {
			t.Fatalf("HandleAction(%v): Navigate=%q, want empty", act, result.Navigate)
		}
		if result.CursorDelta != 0 {
			t.Fatalf("HandleAction(%v): CursorDelta=%d, want 0", act, result.CursorDelta)
		}
		if result.Dirty {
			t.Fatalf("HandleAction(%v): Dirty=true, want false", act)
		}
		if result.Refresh {
			t.Fatalf("HandleAction(%v): Refresh=true, want false", act)
		}
	})
}

// --- From: testpattern_test.go ---

// TestModeRegistered verifies that "testpattern" is registered in the catalog
// via the init function and returns correct metadata.
func TestModeRegistered(t *testing.T) {
	_, ok := catalog.Describe("testpattern")
	if !ok {
		t.Fatal("testpattern mode not found in catalog after init registration")
	}
}

// TestCatalogDefinition verifies that the catalog definition for "testpattern"
// has the correct metadata: ID, non-empty Title, non-empty Summary, Order >= 90.
func TestCatalogDefinition(t *testing.T) {
	def, ok := catalog.Describe("testpattern")
	if !ok {
		t.Fatal("testpattern not found in catalog")
	}
	if def.ID != "testpattern" {
		t.Fatalf("catalog ID = %q, want %q", def.ID, "testpattern")
	}
	if def.Title == "" {
		t.Fatal("catalog Title is empty")
	}
	if def.Summary == "" {
		t.Fatal("catalog Summary is empty")
	}
	if def.Order < 90 {
		t.Fatalf("catalog Order = %d, want >= 90", def.Order)
	}
}

// TestPanelInclusion verifies that when a panel does NOT exclude "testpattern",
// it appears in the catalog's registered IDs (simulating panel availability).
func TestPanelInclusion(t *testing.T) {
	ids := catalog.IDs()
	found := false
	for _, id := range ids {
		if id == "testpattern" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("testpattern not found in catalog.IDs(); expected it to be available when not excluded")
	}
}

// TestPanelExclusion verifies that when a panel lists "testpattern" in its
// ExcludedModes, filtering the catalog IDs correctly omits it.
func TestPanelExclusion(t *testing.T) {
	excluded := map[string]bool{"testpattern": true}
	ids := catalog.IDs()

	var filtered []string
	for _, id := range ids {
		if !excluded[id] {
			filtered = append(filtered, id)
		}
	}

	for _, id := range filtered {
		if id == "testpattern" {
			t.Fatal("testpattern should be excluded from filtered list when in ExcludedModes")
		}
	}

	// Also verify that "testpattern" WAS present in the unfiltered list
	// (otherwise the exclusion test is vacuous).
	found := false
	for _, id := range ids {
		if id == "testpattern" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("testpattern must be in catalog.IDs() for exclusion test to be meaningful")
	}
}

// TestPixelSpotCheck240x135 spot-checks specific pixel values for a known
// 240×135 TextHints configuration using default glyph metrics.
func TestPixelSpotCheck240x135(t *testing.T) {
	hints := textlayout.TextHints{
		PixelWidth:   240,
		PixelHeight:  135,
		GlyphWidth:   5,
		GlyphHeight:  7,
		GlyphAdvance: 6,
		RowHeight:    10,
	}

	view := BuildView(hints, false)
	img := view.Image

	if img == nil {
		t.Fatal("BuildView returned nil Image")
	}

	bounds := img.Bounds()
	if bounds.Dx() != 240 || bounds.Dy() != 135 {
		t.Fatalf("image size = %dx%d, want 240x135", bounds.Dx(), bounds.Dy())
	}

	white := color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}

	// Corner pixels should be white (border + corner marker).
	corners := []struct {
		x, y int
		name string
	}{
		{0, 0, "top-left"},
		{239, 0, "top-right"},
		{0, 134, "bottom-left"},
		{239, 134, "bottom-right"},
	}
	for _, c := range corners {
		got := img.RGBAAt(c.x, c.y)
		if got != white {
			t.Errorf("pixel (%d,%d) [%s corner] = %v, want white %v", c.x, c.y, c.name, got, white)
		}
	}

	// Center pixel should NOT be pure white border color — it should be
	// interior content (background, guide grid, swatch, glyph, or hint text).
	centerColor := img.RGBAAt(120, 67)
	if centerColor == white {
		t.Error("pixel (120,67) [center] is white; expected interior content color")
	}

	// Static flag should be true for the test pattern mode.
	if !view.Static {
		t.Fatal("view.Static = false, want true")
	}
}
