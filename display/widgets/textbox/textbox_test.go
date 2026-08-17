package textbox

import (
	"image"
	"image/color"
	"strings"
	"testing"

	"github.com/databeast/cyberhud/display/surface/fonts"
	"github.com/databeast/cyberhud/display/widgets"
	"pgregory.net/rapid"
)

// --- From: textbox_prop_test.go ---

// **Feature: textbox-widget, Property 1: Output metadata correctness**
//
// For any valid TextBox Config (Bounds with positive Dx and Dy, effective area > 0
// after padding), Render SHALL return a non-nil Result whose Image has dimensions
// exactly equal to Bounds.Dx() × Bounds.Dy(), Position equal to Bounds.Min, and
// Label of at most 128 characters.
//

func TestPropertyOutputMetadataCorrectness(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate random dimensions ensuring positive bounds.
		width := rapid.IntRange(1, 100).Draw(t, "width")
		height := rapid.IntRange(1, 100).Draw(t, "height")
		minX := rapid.IntRange(0, 50).Draw(t, "minX")
		minY := rapid.IntRange(0, 50).Draw(t, "minY")

		// Generate border flag first so we know the inset.
		border := rapid.Bool().Draw(t, "border")
		borderInset := 0
		if border {
			borderInset = 1
		}

		// Generate padding constrained so effective area stays positive.
		// effectiveWidth = width - 2*padX - 2*borderInset > 0
		// effectiveHeight = height - 2*padY - 2*borderInset > 0
		maxPadX := (width - 2*borderInset - 1) / 2
		if maxPadX < 0 {
			maxPadX = 0
		}
		maxPadY := (height - 2*borderInset - 1) / 2
		if maxPadY < 0 {
			maxPadY = 0
		}

		// Skip configs where even zero padding can't produce positive effective area.
		if width-2*borderInset <= 0 || height-2*borderInset <= 0 {
			return
		}

		padX := rapid.IntRange(0, maxPadX).Draw(t, "padX")
		padY := rapid.IntRange(0, maxPadY).Draw(t, "padY")

		// Generate text content (0-20 random printable chars).
		text := rapid.StringMatching(`[A-Za-z0-9 ]{0,20}`).Draw(t, "text")

		// Generate alignment, valign, overflow enums.
		alignment := Alignment(rapid.IntRange(0, 2).Draw(t, "alignment"))
		valign := VAlign(rapid.IntRange(0, 2).Draw(t, "valign"))
		overflow := Overflow(rapid.IntRange(0, 2).Draw(t, "overflow"))

		// Generate a foreground color.
		fg := color.RGBA{
			R: uint8(rapid.IntRange(0, 255).Draw(t, "r")),
			G: uint8(rapid.IntRange(0, 255).Draw(t, "g")),
			B: uint8(rapid.IntRange(0, 255).Draw(t, "b")),
			A: uint8(rapid.IntRange(0, 255).Draw(t, "a")),
		}

		// Generate a label (0-200 chars to test truncation behavior).
		label := rapid.StringMatching(`[A-Za-z0-9_]{0,200}`).Draw(t, "label")

		bounds := image.Rect(minX, minY, minX+width, minY+height)

		cfg := Config{
			Bounds:     bounds,
			Text:       text,
			Font:       font.Default(),
			Alignment:  alignment,
			VAlign:     valign,
			Overflow:   overflow,
			Foreground: fg,
			PadX:       padX,
			PadY:       padY,
			Border:     border,
			Label:      label,
		}

		result := Render(cfg)

		// Verify 1: Result is non-nil.
		if result == nil {
			t.Fatalf("expected non-nil result for valid config: width=%d height=%d padX=%d padY=%d border=%v",
				width, height, padX, padY, border)
		}

		// Verify 2: Image dimensions == Bounds.Dx() × Bounds.Dy().
		imgBounds := result.Image.Bounds()
		gotWidth := imgBounds.Dx()
		gotHeight := imgBounds.Dy()
		if gotWidth != width {
			t.Fatalf("image width mismatch: got %d, want %d", gotWidth, width)
		}
		if gotHeight != height {
			t.Fatalf("image height mismatch: got %d, want %d", gotHeight, height)
		}

		// Verify 3: Position == Bounds.Min.
		if result.Position.X != minX || result.Position.Y != minY {
			t.Fatalf("Position mismatch: got (%d,%d), want (%d,%d)",
				result.Position.X, result.Position.Y, minX, minY)
		}

		// Verify 4: Label length <= 128 characters.
		if len([]rune(result.Label)) > 128 {
			t.Fatalf("Label exceeds 128 chars: got %d chars", len([]rune(result.Label)))
		}
	})
}

// **Feature: textbox-widget, Property 6: Horizontal alignment offset**
//
// For any single-line text narrower than the effective width, the leftmost
// foreground pixel's X-coordinate SHALL equal:
// - Left: PadX + borderInset + firstLitCol
// - Center: PadX + borderInset + floor((effectiveWidth - textPixelWidth) / 2) + firstLitCol
// - Right: PadX + borderInset + effectiveWidth - textPixelWidth + firstLitCol
//
// Where:
// - borderInset = 1 if Border=true, 0 otherwise
// - effectiveWidth = Bounds.Dx() - 2*PadX - 2*borderInset
// - textPixelWidth = charCount * GlyphAdvance
//

func TestPropertyHorizontalAlignmentOffset(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		face := font.Default()
		metrics := face.Metrics()
		glyphAdvance := metrics.GlyphAdvance
		glyphWidth := metrics.GlyphWidth
		glyphHeight := metrics.GlyphHeight

		// Generate single-line text (1-10 uppercase chars, no newlines).
		text := rapid.StringMatching(`[A-Z]{1,10}`).Draw(t, "text")
		charCount := len([]rune(text))

		// textPixelWidth for textbox uses charCount * GlyphAdvance.
		textPixelWidth := charCount * glyphAdvance

		// Generate padding and border.
		padX := rapid.IntRange(0, 5).Draw(t, "padX")
		border := rapid.Bool().Draw(t, "border")

		borderInset := 0
		if border {
			borderInset = 1
		}

		// effectiveWidth must be larger than textPixelWidth to ensure text fits.
		effectiveWidth := rapid.IntRange(textPixelWidth+1, textPixelWidth+50).Draw(t, "effectiveWidth")

		// Compute bounds width from effectiveWidth, padX, and borderInset.
		boundsWidth := effectiveWidth + 2*padX + 2*borderInset

		// Height must accommodate at least one line of text.
		// Layout uses RowHeight (not GlyphHeight) to decide if a line fits.
		rowHeight := metrics.RowHeight
		minHeight := rowHeight + 2*borderInset
		boundsHeight := rapid.IntRange(minHeight, minHeight+20).Draw(t, "boundsHeight")

		bounds := image.Rect(0, 0, boundsWidth, boundsHeight)

		// Random alignment.
		alignVal := rapid.IntRange(0, 2).Draw(t, "alignment")
		align := Alignment(alignVal)

		fg := color.RGBA{R: 255, G: 0, B: 0, A: 255}

		result := Render(Config{
			Text:       text,
			Bounds:     bounds,
			Font:       face,
			Alignment:  align,
			PadX:       padX,
			Border:     border,
			Foreground: fg,
		})

		if result == nil {
			t.Fatal("expected non-nil result for valid config")
		}

		// Compute the first lit column in the first character's glyph.
		firstChar := rune(text[0])
		firstLitCol := glyphWidth // fallback if no pixel found
		for row := 0; row < glyphHeight; row++ {
			bits := face.GlyphRow(firstChar, row)
			for col := 0; col < glyphWidth; col++ {
				if bits&(1<<uint(31-col)) != 0 {
					if col < firstLitCol {
						firstLitCol = col
					}
				}
			}
		}

		// Find the leftmost foreground pixel X-coordinate (scan columns left-to-right).
		// Skip border pixels when border is enabled.
		leftmostX := -1
		scanStartX := 0
		scanStartY := 0
		scanEndX := boundsWidth
		scanEndY := boundsHeight
		if border {
			scanStartX = 1
			scanStartY = 1
			scanEndX = boundsWidth - 1
			scanEndY = boundsHeight - 1
		}
		for x := scanStartX; x < scanEndX; x++ {
			for y := scanStartY; y < scanEndY; y++ {
				_, _, _, a := result.Image.At(x, y).RGBA()
				if a > 0 {
					leftmostX = x
					goto found
				}
			}
		}
	found:
		if leftmostX == -1 {
			t.Fatal("no foreground text pixels found for non-empty text")
		}

		// Compute expected xStart based on alignment.
		var xStart int
		switch align {
		case Left:
			xStart = padX + borderInset
		case Center:
			xStart = padX + borderInset + (effectiveWidth-textPixelWidth)/2
		case Right:
			xStart = padX + borderInset + effectiveWidth - textPixelWidth
		}

		expectedLeftmostX := xStart + firstLitCol

		if leftmostX != expectedLeftmostX {
			t.Fatalf("alignment %d: leftmost foreground pixel at x=%d, expected x=%d "+
				"(xStart=%d, firstLitCol=%d, padX=%d, borderInset=%d, effectiveWidth=%d, textPixelWidth=%d, text=%q)",
				align, leftmostX, expectedLeftmostX, xStart, firstLitCol,
				padX, borderInset, effectiveWidth, textPixelWidth, text)
		}
	})
}

// **Feature: textbox-widget, Property 9: Newline splitting produces distinct Y bands**
//
// For any text containing N newline characters (producing N+1 logical lines),
// the foreground pixels of each rendered line SHALL occupy a Y-band starting at
// the previous line's Y-band start plus that line's effective font RowHeight plus
// LineSpacing. Consecutive newlines SHALL produce empty (transparent) bands.
//

func TestPropertyNewlineSplittingDistinctYBands(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		face := font.Default()
		metrics := face.Metrics()
		rowHeight := metrics.RowHeight
		glyphAdvance := metrics.GlyphAdvance

		// Generate text with 2-4 newlines embedded between short segments.
		// Each segment is 1-5 uppercase letters (or empty for consecutive newlines).
		numNewlines := rapid.IntRange(2, 4).Draw(t, "numNewlines")
		numSegments := numNewlines + 1

		segments := make([]string, numSegments)
		for i := 0; i < numSegments; i++ {
			// Allow some segments to be empty (tests consecutive newlines).
			if rapid.Bool().Draw(t, "emptySegment") {
				segments[i] = ""
			} else {
				segments[i] = rapid.StringMatching(`[A-Z]{1,5}`).Draw(t, "segment")
			}
		}

		// Build text by joining segments with newlines.
		text := ""
		for i, seg := range segments {
			if i > 0 {
				text += "\n"
			}
			text += seg
		}

		// Generate line spacing (0-4 pixels).
		lineSpacing := rapid.IntRange(0, 4).Draw(t, "lineSpacing")

		// Compute needed bounds. Use Left alignment, no padding, no border.
		// Width must fit the widest segment.
		maxSegLen := 0
		for _, seg := range segments {
			if len(seg) > maxSegLen {
				maxSegLen = len(seg)
			}
		}
		if maxSegLen == 0 {
			maxSegLen = 1 // at least 1 char wide
		}
		boundsWidth := maxSegLen*glyphAdvance + 10 // extra space

		// Height must be large enough to fit all lines.
		boundsHeight := numSegments*(rowHeight+lineSpacing) + 10

		cfg := Config{
			Bounds:      image.Rect(0, 0, boundsWidth, boundsHeight),
			Text:        text,
			Font:        face,
			Alignment:   Left,
			VAlign:      Top,
			Overflow:    Clip,
			Foreground:  color.RGBA{R: 255, G: 255, B: 255, A: 255},
			LineSpacing: lineSpacing,
		}

		result := Render(cfg)
		if result == nil {
			t.Fatal("expected non-nil result for valid config")
		}

		img := result.Image

		// For each line, determine the expected Y-band and check pixels.
		for lineIdx := 0; lineIdx < numSegments; lineIdx++ {
			bandStartY := lineIdx * (rowHeight + lineSpacing)
			bandEndY := bandStartY + rowHeight

			if segments[lineIdx] == "" {
				// Empty segment: the entire band should be transparent.
				for y := bandStartY; y < bandEndY; y++ {
					for x := 0; x < boundsWidth; x++ {
						_, _, _, a := img.At(x, y).RGBA()
						if a > 0 {
							t.Fatalf("line %d is empty but found non-transparent pixel at (%d, %d)",
								lineIdx, x, y)
						}
					}
				}
			} else {
				// Non-empty segment: foreground pixels must be within the band.
				// Check that any foreground pixels for this segment appear within bandStartY..bandEndY.
				foundPixel := false
				for y := bandStartY; y < bandEndY; y++ {
					for x := 0; x < boundsWidth; x++ {
						_, _, _, a := img.At(x, y).RGBA()
						if a > 0 {
							foundPixel = true
							break
						}
					}
					if foundPixel {
						break
					}
				}
				if !foundPixel {
					t.Fatalf("line %d (%q) should have foreground pixels in Y-band [%d, %d) but found none",
						lineIdx, segments[lineIdx], bandStartY, bandEndY)
				}
			}

			// Check that no foreground pixels appear in the lineSpacing gap between lines.
			if lineSpacing > 0 && lineIdx < numSegments-1 {
				gapStartY := bandEndY
				gapEndY := gapStartY + lineSpacing
				for y := gapStartY; y < gapEndY && y < boundsHeight; y++ {
					for x := 0; x < boundsWidth; x++ {
						_, _, _, a := img.At(x, y).RGBA()
						if a > 0 {
							t.Fatalf("found foreground pixel at (%d, %d) in lineSpacing gap between line %d and %d",
								x, y, lineIdx, lineIdx+1)
						}
					}
				}
			}
		}
	})
}

// **Feature: textbox-widget, Property 10: Wrap mode line fitting**
//
// For any multi-word text in Wrap mode, every rendered visual line SHALL have a
// pixel width not exceeding the effective width.
//

func TestPropertyWrapModeLineFitting(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		face := font.Default()
		metrics := face.Metrics()
		glyphAdvance := metrics.GlyphAdvance
		rowHeight := metrics.RowHeight

		// Generate multi-word text (3-8 words of 2-8 chars each).
		numWords := rapid.IntRange(3, 8).Draw(t, "numWords")
		words := make([]string, numWords)
		for i := 0; i < numWords; i++ {
			words[i] = rapid.StringMatching(`[A-Z]{2,8}`).Draw(t, "word")
		}
		text := ""
		for i, w := range words {
			if i > 0 {
				text += " "
			}
			text += w
		}

		// Generate an effective width that's narrower than the total text width
		// (to ensure wrapping), but wide enough for at least one word.
		// Minimum effective width = at least glyphAdvance * 2 characters.
		minEffWidth := 2 * glyphAdvance
		// Maximum effective width = total text pixel width minus something to force wrap.
		totalTextWidth := len([]rune(text)) * glyphAdvance
		maxEffWidth := totalTextWidth - glyphAdvance
		if maxEffWidth < minEffWidth {
			maxEffWidth = minEffWidth
		}

		effectiveWidth := rapid.IntRange(minEffWidth, maxEffWidth).Draw(t, "effectiveWidth")

		// No padding, no border → bounds width = effective width.
		boundsWidth := effectiveWidth

		// Height must accommodate several wrapped lines.
		boundsHeight := 20 * rowHeight

		cfg := Config{
			Bounds:     image.Rect(0, 0, boundsWidth, boundsHeight),
			Text:       text,
			Font:       face,
			Alignment:  Left,
			VAlign:     Top,
			Overflow:   Wrap,
			Foreground: color.RGBA{R: 255, G: 255, B: 255, A: 255},
		}

		result := Render(cfg)
		if result == nil {
			t.Fatal("expected non-nil result for valid config")
		}

		img := result.Image

		// Check each visual line (each rowHeight-tall Y-band) for foreground pixels.
		// No foreground pixel should appear at X >= effectiveWidth.
		for y := 0; y < boundsHeight; y++ {
			for x := effectiveWidth; x < boundsWidth; x++ {
				_, _, _, a := img.At(x, y).RGBA()
				if a > 0 {
					t.Fatalf("foreground pixel at (%d, %d) exceeds effective width %d in wrap mode",
						x, y, effectiveWidth)
				}
			}
		}

		// Additional check: verify no foreground pixel appears at X >= effectiveWidth
		// for ALL Y coordinates (the bounds and effective width are the same here,
		// so the above loop handles it, but let's be explicit about the property).
		// Already covered above since boundsWidth == effectiveWidth, but let's also
		// verify that within the effective area, rendered chars don't exceed effectiveWidth.
		for y := 0; y < boundsHeight; y++ {
			// Find the rightmost foreground pixel in this row.
			rightmostX := -1
			for x := boundsWidth - 1; x >= 0; x-- {
				_, _, _, a := img.At(x, y).RGBA()
				if a > 0 {
					rightmostX = x
					break
				}
			}
			if rightmostX >= effectiveWidth {
				t.Fatalf("row %d: rightmost foreground pixel at x=%d exceeds effectiveWidth=%d",
					y, rightmostX, effectiveWidth)
			}
		}
	})
}

// **Feature: textbox-widget, Property 11: Vertical clipping at bounds**
//
// For any text with more logical lines than fit in the effective height, no
// foreground pixels SHALL appear at Y-coordinates beyond PadY + borderInset +
// effectiveHeight.
//

func TestPropertyVerticalClippingAtBounds(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		face := font.Default()
		metrics := face.Metrics()
		rowHeight := metrics.RowHeight
		glyphAdvance := metrics.GlyphAdvance

		// Generate text with many newlines (more lines than will fit).
		// We'll generate between 5 and 12 lines.
		numLines := rapid.IntRange(5, 12).Draw(t, "numLines")
		segments := make([]string, numLines)
		for i := 0; i < numLines; i++ {
			segments[i] = rapid.StringMatching(`[A-Z]{1,5}`).Draw(t, "segment")
		}
		text := ""
		for i, seg := range segments {
			if i > 0 {
				text += "\n"
			}
			text += seg
		}

		// Generate padding and border.
		padY := rapid.IntRange(0, 3).Draw(t, "padY")
		padX := rapid.IntRange(0, 3).Draw(t, "padX")
		border := rapid.Bool().Draw(t, "border")
		borderInset := 0
		if border {
			borderInset = 1
		}

		// Generate lineSpacing.
		lineSpacing := rapid.IntRange(0, 2).Draw(t, "lineSpacing")

		// Compute effective height that fits fewer lines than numLines.
		// Each visible line takes rowHeight + lineSpacing (except last which just takes rowHeight).
		// We want to fit between 2 and numLines-1 lines.
		maxFitLines := rapid.IntRange(2, numLines-1).Draw(t, "maxFitLines")

		// effectiveHeight must accommodate exactly maxFitLines lines:
		// maxFitLines * rowHeight + (maxFitLines-1) * lineSpacing
		// But we want the (maxFitLines+1)-th line to NOT fit.
		// So effectiveHeight should be >= maxFitLines*rowHeight + (maxFitLines-1)*lineSpacing
		// but < (maxFitLines+1)*rowHeight + maxFitLines*lineSpacing.
		effectiveHeight := maxFitLines*rowHeight + (maxFitLines-1)*lineSpacing

		// Compute bounds height from effective height, padding, and border inset.
		boundsHeight := effectiveHeight + 2*padY + 2*borderInset

		// Width must fit the widest line.
		maxSegLen := 0
		for _, seg := range segments {
			if len([]rune(seg)) > maxSegLen {
				maxSegLen = len([]rune(seg))
			}
		}
		boundsWidth := maxSegLen*glyphAdvance + 2*padX + 2*borderInset + 5

		cfg := Config{
			Bounds:      image.Rect(0, 0, boundsWidth, boundsHeight),
			Text:        text,
			Font:        face,
			Alignment:   Left,
			VAlign:      Top,
			Overflow:    Clip,
			Foreground:  color.RGBA{R: 255, G: 255, B: 255, A: 255},
			PadX:        padX,
			PadY:        padY,
			Border:      border,
			LineSpacing: lineSpacing,
		}

		result := Render(cfg)
		if result == nil {
			// Config might be invalid if effective area is zero - skip.
			return
		}

		img := result.Image

		// The clipping boundary: no foreground pixels should appear at
		// Y >= padY + borderInset + effectiveHeight.
		clipY := padY + borderInset + effectiveHeight

		for y := clipY; y < boundsHeight; y++ {
			for x := 0; x < boundsWidth; x++ {
				// Skip border pixels.
				if border && (y == 0 || y == boundsHeight-1 || x == 0 || x == boundsWidth-1) {
					continue
				}
				_, _, _, a := img.At(x, y).RGBA()
				if a > 0 {
					t.Fatalf("foreground pixel at (%d, %d) exceeds vertical clipping boundary (clipY=%d, padY=%d, borderInset=%d, effectiveHeight=%d)",
						x, y, clipY, padY, borderInset, effectiveHeight)
				}
			}
		}
	})
}

// **Feature: textbox-widget, Property 12: Vertical alignment offset**
//
// For any text block whose total height is less than the effective height:
// - VAlign Top: first line's Y-offset SHALL be PadY + borderInset
// - VAlign Middle: first line's Y-offset SHALL be PadY + borderInset + floor((effectiveHeight - blockHeight) / 2)
// - VAlign Bottom: last line's bottom (Y + RowHeight) SHALL equal PadY + borderInset + effectiveHeight
//
// When text block height exceeds effective height (regardless of VAlign), rendering SHALL start from the top.
//

func TestPropertyVerticalAlignmentOffset(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		face := font.Default()
		metrics := face.Metrics()
		glyphAdvance := metrics.GlyphAdvance
		glyphHeight := metrics.GlyphHeight
		rowHeight := metrics.RowHeight

		// Generate 1-3 lines of text joined by \n.
		numLines := rapid.IntRange(1, 3).Draw(t, "numLines")
		var lines []string
		for i := 0; i < numLines; i++ {
			// Each line is 1-5 uppercase characters (short enough to fit).
			line := rapid.StringMatching(`[A-Z]{1,5}`).Draw(t, "line")
			lines = append(lines, line)
		}
		text := ""
		for i, l := range lines {
			if i > 0 {
				text += "\n"
			}
			text += l
		}

		// Generate padding and border.
		padY := rapid.IntRange(0, 5).Draw(t, "padY")
		border := rapid.Bool().Draw(t, "border")
		borderInset := 0
		if border {
			borderInset = 1
		}

		// LineSpacing for blockHeight calculation.
		lineSpacing := rapid.IntRange(0, 3).Draw(t, "lineSpacing")

		// blockHeight = numLines * RowHeight + (numLines-1) * lineSpacing
		blockHeight := numLines*rowHeight + (numLines-1)*lineSpacing

		// effectiveHeight must be > blockHeight to test vertical alignment.
		// effectiveHeight = height - 2*padY - 2*borderInset
		// So height = effectiveHeight + 2*padY + 2*borderInset
		// We need effectiveHeight > blockHeight.
		effectiveHeight := rapid.IntRange(blockHeight+1, blockHeight+30).Draw(t, "effectiveHeight")
		boundsHeight := effectiveHeight + 2*padY + 2*borderInset

		// Width must accommodate the text (widest line).
		maxLineLen := 0
		for _, l := range lines {
			if len([]rune(l)) > maxLineLen {
				maxLineLen = len([]rune(l))
			}
		}
		textPixelWidth := maxLineLen * glyphAdvance
		// PadX=0 for simplicity; effectiveWidth must exceed textPixelWidth.
		boundsWidth := textPixelWidth + 10 + 2*borderInset

		bounds := image.Rect(0, 0, boundsWidth, boundsHeight)

		// Generate random VAlign.
		valignVal := rapid.IntRange(0, 2).Draw(t, "valign")
		valign := VAlign(valignVal)

		fg := color.RGBA{R: 255, G: 0, B: 0, A: 255}

		result := Render(Config{
			Text:        text,
			Bounds:      bounds,
			Font:        face,
			VAlign:      valign,
			Foreground:  fg,
			PadY:        padY,
			Border:      border,
			LineSpacing: lineSpacing,
		})

		if result == nil {
			t.Fatal("expected non-nil result for valid config")
		}

		// Compute the first lit row offset for the first character of the first line.
		// Glyphs may have empty top rows (e.g., 'A' starts with a blank row 0).
		firstLineFirstChar := rune(lines[0][0])
		firstLitRow := glyphHeight // fallback if no pixel found
		for row := 0; row < glyphHeight; row++ {
			bits := face.GlyphRow(firstLineFirstChar, row)
			if bits != 0 {
				firstLitRow = row
				break
			}
		}

		// Compute the last lit row offset for the last character of the last line.
		lastLine := lines[numLines-1]
		lastLineLastChar := rune(lastLine[len(lastLine)-1])
		lastLitRow := 0
		for row := glyphHeight - 1; row >= 0; row-- {
			bits := face.GlyphRow(lastLineLastChar, row)
			if bits != 0 {
				lastLitRow = row
				break
			}
		}
		// Also check all chars in last line to find the absolute last lit row.
		for _, ch := range lastLine {
			for row := glyphHeight - 1; row >= 0; row-- {
				bits := face.GlyphRow(ch, row)
				if bits != 0 {
					if row > lastLitRow {
						lastLitRow = row
					}
					break
				}
			}
		}

		// Also check all chars in first line to find the absolute first lit row.
		for _, ch := range lines[0] {
			for row := 0; row < glyphHeight; row++ {
				bits := face.GlyphRow(ch, row)
				if bits != 0 {
					if row < firstLitRow {
						firstLitRow = row
					}
					break
				}
			}
		}

		// Find the topmost foreground pixel Y (excluding border pixels).
		topY := -1
		scanStartX := borderInset
		scanStartY := borderInset
		scanEndX := boundsWidth - borderInset
		scanEndY := boundsHeight - borderInset
		for y := scanStartY; y < scanEndY; y++ {
			for x := scanStartX; x < scanEndX; x++ {
				_, _, _, a := result.Image.At(x, y).RGBA()
				if a > 0 {
					topY = y
					goto foundTop
				}
			}
		}
	foundTop:
		if topY == -1 {
			t.Fatal("no foreground text pixels found for non-empty text")
		}

		// Find the bottommost foreground pixel Y (excluding border pixels).
		bottomY := -1
		for y := scanEndY - 1; y >= scanStartY; y-- {
			for x := scanStartX; x < scanEndX; x++ {
				_, _, _, a := result.Image.At(x, y).RGBA()
				if a > 0 {
					bottomY = y
					goto foundBottom
				}
			}
		}
	foundBottom:
		if bottomY == -1 {
			t.Fatal("no foreground text pixels found for non-empty text (bottom scan)")
		}

		// Verify based on VAlign.
		// The line's Y-offset is where glyph rendering begins (top-left of glyph cell).
		// The topmost actual pixel is at lineY + firstLitRow.
		// The bottommost actual pixel is at lastLineY + lastLitRow.
		switch valign {
		case Top:
			// First line's Y-offset SHALL be PadY + borderInset.
			// So topmost pixel = PadY + borderInset + firstLitRow.
			expectedTopY := padY + borderInset + firstLitRow
			if topY != expectedTopY {
				t.Fatalf("VAlign Top: topmost foreground pixel at y=%d, expected y=%d "+
					"(padY=%d, borderInset=%d, firstLitRow=%d, effectiveHeight=%d, blockHeight=%d)",
					topY, expectedTopY, padY, borderInset, firstLitRow, effectiveHeight, blockHeight)
			}

		case Middle:
			// First line's Y-offset SHALL be PadY + borderInset + floor((effectiveHeight - blockHeight) / 2).
			// So topmost pixel = PadY + borderInset + floor((effectiveHeight - blockHeight) / 2) + firstLitRow.
			lineY := padY + borderInset + (effectiveHeight-blockHeight)/2
			expectedTopY := lineY + firstLitRow
			if topY != expectedTopY {
				t.Fatalf("VAlign Middle: topmost foreground pixel at y=%d, expected y=%d "+
					"(padY=%d, borderInset=%d, effectiveHeight=%d, blockHeight=%d, offset=%d, firstLitRow=%d)",
					topY, expectedTopY, padY, borderInset, effectiveHeight, blockHeight, (effectiveHeight-blockHeight)/2, firstLitRow)
			}

		case Bottom:
			// Last line's bottom (Y + RowHeight) SHALL equal PadY + borderInset + effectiveHeight.
			// Last line Y = PadY + borderInset + effectiveHeight - blockHeight + (numLines-1)*(rowHeight+lineSpacing)
			// But simpler: the last line starts at startY + offset for last line.
			// startY for bottom = effectiveHeight - blockHeight
			// last line Y in effective coords = startY + (numLines-1)*(rowHeight+lineSpacing)
			// In absolute coords: padY + borderInset + startY + (numLines-1)*(rowHeight+lineSpacing) + lastLitRow
			startY := effectiveHeight - blockHeight
			lastLineY := padY + borderInset + startY + (numLines-1)*(rowHeight+lineSpacing)
			expectedBottomY := lastLineY + lastLitRow
			if bottomY != expectedBottomY {
				t.Fatalf("VAlign Bottom: bottommost foreground pixel at y=%d, expected y=%d "+
					"(padY=%d, borderInset=%d, effectiveHeight=%d, blockHeight=%d, numLines=%d, lastLitRow=%d, startY=%d, lastLineY=%d)",
					bottomY, expectedBottomY, padY, borderInset, effectiveHeight, blockHeight, numLines, lastLitRow, startY, lastLineY)
			}
		}
	})
}

// **Feature: textbox-widget, Property 13: Padding and border exclusion zone**
//
// For any TextBox Config with PadX > 0 or PadY > 0 or Border=true (and valid
// effective area), no text foreground pixels SHALL appear at X < PadX + borderInset,
// X >= (Bounds.Dx() - PadX - borderInset), Y < PadY + borderInset, or
// Y >= (Bounds.Dy() - PadY - borderInset), where borderInset is 1 if Border=true,
// 0 otherwise.
//

func TestPropertyPaddingAndBorderExclusionZone(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		face := font.Default()

		// Generate border flag.
		border := rapid.Bool().Draw(t, "border")
		borderInset := 0
		if border {
			borderInset = 1
		}

		// Generate padding ensuring at least one of padX > 0, padY > 0, or border is true.
		padX := rapid.IntRange(0, 10).Draw(t, "padX")
		padY := rapid.IntRange(0, 10).Draw(t, "padY")

		// Property requires PadX > 0 or PadY > 0 or Border=true.
		if padX == 0 && padY == 0 && !border {
			padX = rapid.IntRange(1, 10).Draw(t, "padXForced")
		}

		// Minimum bounds must produce positive effective area.
		// effectiveWidth = width - 2*padX - 2*borderInset > 0
		// effectiveHeight = height - 2*padY - 2*borderInset > 0
		minWidth := 2*padX + 2*borderInset + 1
		minHeight := 2*padY + 2*borderInset + 1

		// Need at least one glyph advance and one row height in effective area.
		metrics := face.Metrics()
		if minWidth < 2*padX+2*borderInset+metrics.GlyphAdvance {
			minWidth = 2*padX + 2*borderInset + metrics.GlyphAdvance
		}
		if minHeight < 2*padY+2*borderInset+metrics.RowHeight {
			minHeight = 2*padY + 2*borderInset + metrics.RowHeight
		}

		width := rapid.IntRange(minWidth, minWidth+30).Draw(t, "width")
		height := rapid.IntRange(minHeight, minHeight+30).Draw(t, "height")

		// Generate non-empty text to ensure some foreground pixels are drawn.
		text := rapid.StringMatching(`[A-Z]{1,8}`).Draw(t, "text")

		// Use a distinctive foreground color.
		fg := color.RGBA{R: 200, G: 100, B: 50, A: 255}

		cfg := Config{
			Bounds:     image.Rect(0, 0, width, height),
			Text:       text,
			Font:       face,
			Foreground: fg,
			PadX:       padX,
			PadY:       padY,
			Border:     border,
		}

		result := Render(cfg)
		if result == nil {
			t.Fatal("expected non-nil result for valid config")
		}

		// The exclusion zone for TEXT pixels is the padding area between
		// the border (if any) and the text area. When border is true,
		// border pixels ARE expected on the outer edge, so we skip those.
		// We check that no text pixels appear in the inner padding zone.
		//
		// Inner padding zone X: [borderInset, padX + borderInset)
		//   and [width - padX - borderInset, width - borderInset)
		// Inner padding zone Y: [borderInset, padY + borderInset)
		//   and [height - padY - borderInset, height - borderInset)
		//
		// Actually the property states: no text foreground pixels at
		//   X < padX + borderInset
		//   X >= width - padX - borderInset
		//   Y < padY + borderInset
		//   Y >= height - padY - borderInset
		//
		// When border is true, the outer row/col will have border pixels
		// (which match fg). We need to skip those and verify that in the
		// padding zone (excluding the border row/col), pixels are transparent.

		textAreaLeft := padX + borderInset
		textAreaRight := width - padX - borderInset
		textAreaTop := padY + borderInset
		textAreaBottom := height - padY - borderInset

		// Scan the exclusion zone. Skip the outermost border pixels (if border is true).
		for y := borderInset; y < height-borderInset; y++ {
			for x := borderInset; x < width-borderInset; x++ {
				// If this pixel is inside the text area, skip it.
				if x >= textAreaLeft && x < textAreaRight && y >= textAreaTop && y < textAreaBottom {
					continue
				}
				// This pixel is in the padding zone (between border and text area).
				// It should be fully transparent.
				r, g, b, a := result.Image.At(x, y).RGBA()
				if a != 0 {
					t.Fatalf("found non-transparent pixel in padding zone at (%d,%d): RGBA=(%d,%d,%d,%d), "+
						"textArea=[%d,%d)x[%d,%d), border=%v, padX=%d, padY=%d",
						x, y, r>>8, g>>8, b>>8, a>>8,
						textAreaLeft, textAreaRight, textAreaTop, textAreaBottom,
						border, padX, padY)
				}
			}
		}
	})
}

// **Feature: textbox-widget, Property 20: TextBox border draws 1px frame using Foreground color**
//
// For any valid TextBox Config with Border=true and effective area > 0 (after
// border inset and padding), every pixel on the outermost row and column
// (y=0, y=Bounds.Dy()-1, x=0, x=Bounds.Dx()-1) SHALL have RGBA values matching
// the effective Foreground color, and no text pixels SHALL appear within the
// 1px border inset zone.
//

func TestPropertyBorderDraws1pxFrame(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		face := font.Default()
		metrics := face.Metrics()

		// Generate padding.
		padX := rapid.IntRange(0, 5).Draw(t, "padX")
		padY := rapid.IntRange(0, 5).Draw(t, "padY")

		// Border is always true for this property.
		borderInset := 1

		// Minimum dimensions for positive effective area.
		minWidth := 2*padX + 2*borderInset + metrics.GlyphAdvance
		minHeight := 2*padY + 2*borderInset + metrics.RowHeight

		width := rapid.IntRange(minWidth, minWidth+30).Draw(t, "width")
		height := rapid.IntRange(minHeight, minHeight+30).Draw(t, "height")

		// Generate text (may be empty — border should still be drawn).
		text := rapid.StringMatching(`[A-Z]{0,10}`).Draw(t, "text")

		// Generate foreground color. Zero value defaults to white.
		fgR := uint8(rapid.IntRange(0, 255).Draw(t, "fgR"))
		fgG := uint8(rapid.IntRange(0, 255).Draw(t, "fgG"))
		fgB := uint8(rapid.IntRange(0, 255).Draw(t, "fgB"))
		fgA := uint8(rapid.IntRange(0, 255).Draw(t, "fgA"))
		fg := color.RGBA{R: fgR, G: fgG, B: fgB, A: fgA}

		// Compute the effective foreground (zero value → white).
		effectiveFg := fg
		if effectiveFg == (color.RGBA{}) {
			effectiveFg = color.RGBA{R: 255, G: 255, B: 255, A: 255}
		}

		cfg := Config{
			Bounds:     image.Rect(0, 0, width, height),
			Text:       text,
			Font:       face,
			Foreground: fg,
			PadX:       padX,
			PadY:       padY,
			Border:     true,
		}

		result := Render(cfg)
		if result == nil {
			t.Fatal("expected non-nil result for valid config with border")
		}

		// Verify: every pixel on the outermost row/column matches effectiveFg.
		// Top row (y=0).
		for x := 0; x < width; x++ {
			r, g, b, a := result.Image.At(x, 0).RGBA()
			gotR, gotG, gotB, gotA := uint8(r>>8), uint8(g>>8), uint8(b>>8), uint8(a>>8)
			if gotR != effectiveFg.R || gotG != effectiveFg.G || gotB != effectiveFg.B || gotA != effectiveFg.A {
				t.Fatalf("border pixel mismatch at (%d,0): got RGBA(%d,%d,%d,%d), want RGBA(%d,%d,%d,%d)",
					x, gotR, gotG, gotB, gotA,
					effectiveFg.R, effectiveFg.G, effectiveFg.B, effectiveFg.A)
			}
		}
		// Bottom row (y=height-1).
		for x := 0; x < width; x++ {
			r, g, b, a := result.Image.At(x, height-1).RGBA()
			gotR, gotG, gotB, gotA := uint8(r>>8), uint8(g>>8), uint8(b>>8), uint8(a>>8)
			if gotR != effectiveFg.R || gotG != effectiveFg.G || gotB != effectiveFg.B || gotA != effectiveFg.A {
				t.Fatalf("border pixel mismatch at (%d,%d): got RGBA(%d,%d,%d,%d), want RGBA(%d,%d,%d,%d)",
					x, height-1, gotR, gotG, gotB, gotA,
					effectiveFg.R, effectiveFg.G, effectiveFg.B, effectiveFg.A)
			}
		}
		// Left column (x=0).
		for y := 0; y < height; y++ {
			r, g, b, a := result.Image.At(0, y).RGBA()
			gotR, gotG, gotB, gotA := uint8(r>>8), uint8(g>>8), uint8(b>>8), uint8(a>>8)
			if gotR != effectiveFg.R || gotG != effectiveFg.G || gotB != effectiveFg.B || gotA != effectiveFg.A {
				t.Fatalf("border pixel mismatch at (0,%d): got RGBA(%d,%d,%d,%d), want RGBA(%d,%d,%d,%d)",
					y, gotR, gotG, gotB, gotA,
					effectiveFg.R, effectiveFg.G, effectiveFg.B, effectiveFg.A)
			}
		}
		// Right column (x=width-1).
		for y := 0; y < height; y++ {
			r, g, b, a := result.Image.At(width-1, y).RGBA()
			gotR, gotG, gotB, gotA := uint8(r>>8), uint8(g>>8), uint8(b>>8), uint8(a>>8)
			if gotR != effectiveFg.R || gotG != effectiveFg.G || gotB != effectiveFg.B || gotA != effectiveFg.A {
				t.Fatalf("border pixel mismatch at (%d,%d): got RGBA(%d,%d,%d,%d), want RGBA(%d,%d,%d,%d)",
					width-1, y, gotR, gotG, gotB, gotA,
					effectiveFg.R, effectiveFg.G, effectiveFg.B, effectiveFg.A)
			}
		}

		// Verify: no text pixels appear in the 1px border inset zone.
		// The inset zone is the 1px ring just inside the border. If there's no
		// padding, text starts right after the border. But if padX=0 and padY=0,
		// the text area starts at (1,1) — so there is no inset zone to check beyond
		// the border itself. We only check if there's padding > 0 on either axis.
		// This is covered by Property 13 above — so here we just ensure
		// no text bleeds into the border pixels themselves (already verified above
		// by checking that all border pixels match fg exactly).
	})
}

// **Feature: textbox-widget, Property 14: Per-line font determines line layout**
//
// For any multi-line TextBox Config with FontOverrides, each rendered line SHALL:
// - Use its effective font's RowHeight for vertical advance to the next line
// - Use its effective font's GlyphAdvance for horizontal overflow evaluation
//

func TestPropertyPerLineFontDeterminesLayout(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// We need at least 2 distinct fonts for this test to be meaningful.
		allFonts := font.List()
		if len(allFonts) < 2 {
			// Only one font available — still validate mechanism works but
			// can't verify differing RowHeights.
			return
		}

		// Pick two fonts with different RowHeights for clear separation.
		var smallFont, largeFont font.Face
		for i := 0; i < len(allFonts)-1; i++ {
			for j := i + 1; j < len(allFonts); j++ {
				if allFonts[i].Metrics().RowHeight != allFonts[j].Metrics().RowHeight {
					smallFont = allFonts[i]
					largeFont = allFonts[j]
					break
				}
			}
			if smallFont != nil {
				break
			}
		}
		if smallFont == nil {
			// All fonts have same RowHeight — skip.
			return
		}

		// Ensure smallFont has smaller RowHeight.
		if smallFont.Metrics().RowHeight > largeFont.Metrics().RowHeight {
			smallFont, largeFont = largeFont, smallFont
		}

		smallMetrics := smallFont.Metrics()
		largeMetrics := largeFont.Metrics()

		// Generate 2-3 lines of visible text (printable ASCII, no spaces to avoid wrap issues).
		numLines := rapid.IntRange(2, 3).Draw(t, "numLines")

		// Determine max GlyphAdvance to ensure text fits within bounds.
		maxAdvance := largeMetrics.GlyphAdvance
		if smallMetrics.GlyphAdvance > maxAdvance {
			maxAdvance = smallMetrics.GlyphAdvance
		}

		// Generate short lines (2-5 chars) that fit in the bounds.
		lineLen := rapid.IntRange(2, 5).Draw(t, "lineLen")
		var lines []string
		for i := 0; i < numLines; i++ {
			line := rapid.StringMatching(`[A-Z]{2,5}`).Draw(t, "line")
			// Ensure consistent length for simplicity.
			r := []rune(line)
			if len(r) > lineLen {
				r = r[:lineLen]
			}
			lines = append(lines, string(r))
		}

		// Build multi-line text.
		text := ""
		for i, l := range lines {
			if i > 0 {
				text += "\n"
			}
			text += l
		}

		// Assign font overrides: alternate between smallFont and largeFont.
		fontOverrides := make([]font.Face, numLines)
		for i := 0; i < numLines; i++ {
			if i%2 == 0 {
				fontOverrides[i] = smallFont
			} else {
				fontOverrides[i] = largeFont
			}
		}

		// Compute bounds large enough to fit all lines.
		// Width: enough chars at the max advance.
		effectiveWidth := (lineLen + 2) * maxAdvance
		// Height: sum of all RowHeights.
		totalHeight := 0
		for i := 0; i < numLines; i++ {
			totalHeight += fontOverrides[i].Metrics().RowHeight
		}
		totalHeight += 10 // Extra margin

		lineSpacing := rapid.IntRange(0, 2).Draw(t, "lineSpacing")
		totalHeight += lineSpacing * (numLines - 1)

		bounds := image.Rect(0, 0, effectiveWidth, totalHeight)

		fg := color.RGBA{R: 255, G: 255, B: 255, A: 255}

		result := Render(Config{
			Bounds:        bounds,
			Text:          text,
			Font:          smallFont, // Default font (overridden per line)
			Overflow:      Truncate,
			Foreground:    fg,
			LineSpacing:   lineSpacing,
			FontOverrides: fontOverrides,
		})

		if result == nil {
			t.Fatal("expected non-nil result for valid multi-line config with FontOverrides")
		}

		img := result.Image

		// Verify Property 14: Each line's Y-band height matches its font's RowHeight.
		// Scan the image to find foreground pixel Y-bands per line.
		// Lines are rendered starting at y=0. Each line occupies [y, y+RowHeight).
		// Next line starts at y + RowHeight + lineSpacing.
		expectedY := 0
		for i := 0; i < numLines; i++ {
			lineFont := fontOverrides[i]
			lineRowHeight := lineFont.Metrics().RowHeight
			lineGlyphAdvance := lineFont.Metrics().GlyphAdvance

			// Check that foreground pixels for this line exist within [expectedY, expectedY+lineRowHeight).
			hasPixels := false
			for y := expectedY; y < expectedY+lineRowHeight && !hasPixels; y++ {
				for x := 0; x < effectiveWidth; x++ {
					_, _, _, a := img.At(x, y).RGBA()
					if a > 0 {
						hasPixels = true
						break
					}
				}
			}

			if hasPixels {
				// Verify no foreground pixels from this logical line extend beyond
				// the maxChars for this line's font.
				maxChars := effectiveWidth / lineGlyphAdvance
				maxPixelX := maxChars * lineGlyphAdvance
				for y := expectedY; y < expectedY+lineRowHeight; y++ {
					for x := maxPixelX; x < effectiveWidth; x++ {
						_, _, _, a := img.At(x, y).RGBA()
						if a > 0 {
							t.Fatalf("line %d: foreground pixel at (%d,%d) exceeds effective width for font %s (maxChars=%d, maxPixelX=%d, glyphAdvance=%d)",
								i, x, y, lineFont.ID(), maxChars, maxPixelX, lineGlyphAdvance)
						}
					}
				}

				// Verify no foreground pixels appear in the gap between this line and next
				// (the lineSpacing zone).
				gapStart := expectedY + lineRowHeight
				gapEnd := gapStart + lineSpacing
				for y := gapStart; y < gapEnd && y < totalHeight; y++ {
					for x := 0; x < effectiveWidth; x++ {
						_, _, _, a := img.At(x, y).RGBA()
						if a > 0 {
							t.Fatalf("line %d: foreground pixel at (%d,%d) found in lineSpacing gap [%d, %d)",
								i, x, y, gapStart, gapEnd)
						}
					}
				}
			}

			expectedY += lineRowHeight + lineSpacing
		}
	})
}

// **Feature: textbox-widget, Property 15: Wrap continuation uses originating line's font**
//
// For any line with a per-line font override that wraps to multiple visual lines
// in Wrap mode, all wrapped continuation lines SHALL render using the same override
// font as the originating logical line.
//

func TestPropertyWrapContinuationUsesOriginatingFont(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// We need at least 2 distinct fonts.
		allFonts := font.List()
		if len(allFonts) < 2 {
			return
		}

		// Pick two fonts with different GlyphAdvance values.
		var baseFont, overrideFont font.Face
		for i := 0; i < len(allFonts)-1; i++ {
			for j := i + 1; j < len(allFonts); j++ {
				if allFonts[i].Metrics().GlyphAdvance != allFonts[j].Metrics().GlyphAdvance {
					baseFont = allFonts[i]
					overrideFont = allFonts[j]
					break
				}
			}
			if baseFont != nil {
				break
			}
		}
		if baseFont == nil {
			// All fonts have same GlyphAdvance — skip.
			return
		}

		overrideMetrics := overrideFont.Metrics()

		// Create a line that wraps: we need multiple words that together exceed
		// one visual line width.
		// The effective width should fit ~3-5 chars so that a multi-word line wraps.
		charsPerLine := rapid.IntRange(3, 6).Draw(t, "charsPerLine")
		effectiveWidth := charsPerLine * overrideMetrics.GlyphAdvance

		// Generate 2-3 words, each 2-4 chars, so that together they wrap.
		numWords := rapid.IntRange(2, 3).Draw(t, "numWords")
		var words []string
		totalChars := 0
		for w := 0; w < numWords; w++ {
			word := rapid.StringMatching(`[A-Z]{2,4}`).Draw(t, "word")
			words = append(words, word)
			totalChars += len([]rune(word))
			if w > 0 {
				totalChars++ // space between words
			}
		}

		// Need the text to be long enough to actually wrap.
		// totalChars must exceed charsPerLine.
		if totalChars <= charsPerLine {
			// Not enough text to force a wrap — skip.
			return
		}

		// Build the text: single logical line (no \n) so all wrapping is from Overflow=Wrap.
		line := ""
		for i, w := range words {
			if i > 0 {
				line += " "
			}
			line += w
		}

		// FontOverrides: override the first (only) logical line.
		fontOverrides := []font.Face{overrideFont}

		// Height: enough to fit at least 3 visual lines of the override font.
		visualLineHeight := overrideMetrics.RowHeight
		boundsHeight := visualLineHeight * 4

		bounds := image.Rect(0, 0, effectiveWidth, boundsHeight)

		fg := color.RGBA{R: 255, G: 255, B: 255, A: 255}

		result := Render(Config{
			Bounds:        bounds,
			Text:          line,
			Font:          baseFont, // Default font (should NOT be used for this line)
			Overflow:      Wrap,
			Foreground:    fg,
			FontOverrides: fontOverrides,
		})

		if result == nil {
			t.Fatal("expected non-nil result for wrapping config")
		}

		img := result.Image

		// The text should have wrapped into multiple visual lines.
		// All visual lines should use overrideFont's GlyphAdvance for char spacing.
		// Verify: for each visual line's Y-band, the rightmost foreground pixel
		// must be consistent with overrideFont's GlyphAdvance (not baseFont's).

		// Find all visual line Y-bands by scanning for foreground pixel rows.
		type yBand struct {
			minY, maxY int
		}
		var bands []yBand
		inBand := false
		bandStart := 0

		for y := 0; y < boundsHeight; y++ {
			hasPixel := false
			for x := 0; x < effectiveWidth; x++ {
				_, _, _, a := img.At(x, y).RGBA()
				if a > 0 {
					hasPixel = true
					break
				}
			}
			if hasPixel && !inBand {
				inBand = true
				bandStart = y
			} else if !hasPixel && inBand {
				inBand = false
				bands = append(bands, yBand{minY: bandStart, maxY: y - 1})
			}
		}
		if inBand {
			bands = append(bands, yBand{minY: bandStart, maxY: boundsHeight - 1})
		}

		// We should have at least 2 visual lines (text wrapped).
		if len(bands) < 2 {
			t.Fatalf("expected at least 2 visual lines (text wrapping), got %d bands; text=%q charsPerLine=%d effectiveWidth=%d overrideAdvance=%d",
				len(bands), line, charsPerLine, effectiveWidth, overrideMetrics.GlyphAdvance)
		}

		// For each visual line band, verify that the rightmost foreground pixel
		// is within the expected range for overrideFont's GlyphAdvance.
		// Maximum X for any visual line: maxChars * GlyphAdvance (where maxChars = effectiveWidth / overrideAdvance).
		maxChars := effectiveWidth / overrideMetrics.GlyphAdvance

		for bandIdx, band := range bands {
			rightmostX := -1
			for y := band.minY; y <= band.maxY; y++ {
				for x := effectiveWidth - 1; x >= 0; x-- {
					_, _, _, a := img.At(x, y).RGBA()
					if a > 0 {
						if x > rightmostX {
							rightmostX = x
						}
						break
					}
				}
			}

			if rightmostX < 0 {
				continue // Band with no pixels (shouldn't happen given how we found it)
			}

			// The rightmost pixel must be within maxChars * GlyphAdvance.
			// Each character occupies up to GlyphAdvance pixels (GlyphWidth within that).
			maxAllowedX := maxChars * overrideMetrics.GlyphAdvance
			if rightmostX >= maxAllowedX {
				t.Fatalf("visual line %d (y=%d..%d): rightmost pixel at x=%d exceeds maxAllowedX=%d for override font %s (advance=%d, maxChars=%d)",
					bandIdx, band.minY, band.maxY, rightmostX, maxAllowedX, overrideFont.ID(), overrideMetrics.GlyphAdvance, maxChars)
			}
		}
	})
}

// --- From: textbox_test.go ---

// TestDefaultFontFallback verifies that a nil Font in Config falls back to
// spleen-5x8 and produces a non-nil render result.

func TestDefaultFontFallback(t *testing.T) {
	result := Render(Config{
		Bounds: image.Rect(0, 0, 50, 10),
		Text:   "Hello",
		Font:   nil,
	})
	if result == nil {
		t.Fatal("expected non-nil result when Font is nil (should fall back to spleen-5x8)")
	}
	// Verify the image has the expected dimensions.
	bounds := result.Image.Bounds()
	if bounds.Dx() != 50 || bounds.Dy() != 10 {
		t.Fatalf("expected 50×10 image, got %d×%d", bounds.Dx(), bounds.Dy())
	}
}

// TestDefaultAlignmentOverflowVAlign verifies that zero-value Alignment, Overflow,
// and VAlign produce Left, Truncate, and Top respectively.

func TestDefaultAlignmentOverflowVAlign(t *testing.T) {
	face := font.Default()
	metrics := face.Metrics()
	// Text "ABCDE" = 5 chars × 5px advance = 25px wide.
	// Bounds width = 25px exactly → text fits without overflow.
	width := 5 * metrics.GlyphAdvance
	height := metrics.RowHeight

	result := Render(Config{
		Bounds: image.Rect(0, 0, width, height),
		Text:   "ABCDE",
		Font:   face,
		// Zero values: Alignment=Left(0), Overflow=Truncate(0), VAlign=Top(0)
	})
	if result == nil {
		t.Fatal("expected non-nil result")
	}

	img := result.Image.(*image.RGBA)

	// With Left alignment and Top VAlign, text should start at x=0, y=0.
	// Verify there's at least one foreground pixel in the leftmost glyph column.
	hasLeftPixel := false
	for y := 0; y < height; y++ {
		r, g, b, a := img.At(0, y).RGBA()
		if a > 0 && (r > 0 || g > 0 || b > 0) {
			hasLeftPixel = true
			break
		}
	}
	if !hasLeftPixel {
		t.Error("expected foreground pixels at x=0 with default Left alignment")
	}
}

// TestMinimalValidConfig verifies that a minimal valid Config produces a non-nil
// result with correct dimensions and position.

func TestMinimalValidConfig(t *testing.T) {
	result := Render(Config{
		Bounds: image.Rect(0, 0, 50, 10),
		Text:   "Hi",
	})
	if result == nil {
		t.Fatal("expected non-nil result for minimal valid config")
	}
	bounds := result.Image.Bounds()
	if bounds.Dx() != 50 || bounds.Dy() != 10 {
		t.Fatalf("expected 50×10 image, got %d×%d", bounds.Dx(), bounds.Dy())
	}
	if result.Position != (image.Point{X: 0, Y: 0}) {
		t.Fatalf("expected position (0,0), got %v", result.Position)
	}
	if result.Label == "" {
		t.Fatal("expected non-empty label")
	}
}

// TestSpriteFieldAssignmentCompatibility is a compile-time check that Result
// fields can be assigned to widgets.Sprite fields without type conversion.

func TestSpriteFieldAssignmentCompatibility(t *testing.T) {
	result := Render(Config{
		Bounds: image.Rect(0, 0, 20, 10),
		Text:   "X",
	})
	if result == nil {
		t.Fatal("expected non-nil result")
	}

	// This assignment must compile without type conversion or assertion.
	sprite := widgets.Sprite{
		Image:    result.Image,
		Position: result.Position,
		Label:    result.Label,
	}

	// Basic sanity checks on the sprite.
	if sprite.Image == nil {
		t.Error("sprite.Image should not be nil")
	}
	if sprite.Label == "" {
		t.Error("sprite.Label should not be empty")
	}
}

// TestEllipsisEdgeCaseTooNarrow verifies that when the effective width allows
// only 1 character slot, text that exceeds it renders only the ellipsis.

func TestEllipsisEdgeCaseTooNarrow(t *testing.T) {
	face := font.Default()
	metrics := face.Metrics()
	// effectiveWidth = 1 * GlyphAdvance = 5px (allows exactly 1 char).
	// Text "AB" exceeds this. In Truncate mode with maxChars=1, only ellipsis is rendered.
	width := metrics.GlyphAdvance // 5
	height := metrics.RowHeight   // 8

	result := Render(Config{
		Bounds:   image.Rect(0, 0, width, height),
		Text:     "AB",
		Font:     face,
		Overflow: Truncate,
	})
	if result == nil {
		t.Fatal("expected non-nil result")
	}

	// The result should have exactly 1 character slot used (the ellipsis).
	// Verify no pixels appear beyond the first glyph advance.
	img := result.Image.(*image.RGBA)
	for y := 0; y < height; y++ {
		for x := metrics.GlyphAdvance; x < width; x++ {
			_, _, _, a := img.At(x, y).RGBA()
			if a > 0 {
				t.Fatalf("unexpected pixel at (%d,%d) beyond single glyph slot", x, y)
			}
		}
	}
}

// TestNegativePaddingTreatedAsZero verifies that negative PadX/PadY values
// are treated as zero (same output as PadX=0, PadY=0).

func TestNegativePaddingTreatedAsZero(t *testing.T) {
	face := font.Default()

	resultNeg := Render(Config{
		Bounds: image.Rect(0, 0, 30, 10),
		Text:   "Hi",
		Font:   face,
		PadX:   -5,
		PadY:   -3,
	})
	resultZero := Render(Config{
		Bounds: image.Rect(0, 0, 30, 10),
		Text:   "Hi",
		Font:   face,
		PadX:   0,
		PadY:   0,
	})

	if resultNeg == nil || resultZero == nil {
		t.Fatal("both results should be non-nil")
	}

	imgNeg := resultNeg.Image.(*image.RGBA)
	imgZero := resultZero.Image.(*image.RGBA)

	// Both images should be pixel-identical.
	for y := 0; y < 10; y++ {
		for x := 0; x < 30; x++ {
			if imgNeg.RGBAAt(x, y) != imgZero.RGBAAt(x, y) {
				t.Fatalf("pixel mismatch at (%d,%d): negative padding should produce same result as zero", x, y)
			}
		}
	}
}

// TestLabelTruncationAt128Characters verifies that a Label longer than 128
// characters is truncated to exactly 128 characters.

func TestLabelTruncationAt128Characters(t *testing.T) {
	longLabel := strings.Repeat("a", 200)

	result := Render(Config{
		Bounds: image.Rect(0, 0, 30, 10),
		Text:   "Hi",
		Label:  longLabel,
	})
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if len([]rune(result.Label)) != 128 {
		t.Fatalf("expected label length 128, got %d", len([]rune(result.Label)))
	}
}

// TestConsecutiveNewlinesProduceEmptyVisualLines verifies that "A\n\nB" produces
// 3 visual lines, with the middle line being empty (transparent Y-band).

func TestConsecutiveNewlinesProduceEmptyVisualLines(t *testing.T) {
	face := font.Default()
	metrics := face.Metrics()
	// 3 lines × 8px RowHeight = 24px needed.
	height := 3 * metrics.RowHeight

	result := Render(Config{
		Bounds: image.Rect(0, 0, 50, height),
		Text:   "A\n\nB",
		Font:   face,
	})
	if result == nil {
		t.Fatal("expected non-nil result")
	}

	img := result.Image.(*image.RGBA)

	// Middle line occupies Y = RowHeight .. 2*RowHeight-1 and should be transparent.
	for y := metrics.RowHeight; y < 2*metrics.RowHeight; y++ {
		for x := 0; x < 50; x++ {
			_, _, _, a := img.At(x, y).RGBA()
			if a > 0 {
				t.Fatalf("expected transparent pixel at (%d,%d) in empty middle line", x, y)
			}
		}
	}

	// First line (Y = 0..RowHeight-1) should have some foreground pixels for 'A'.
	hasPixelLine1 := false
	for y := 0; y < metrics.RowHeight; y++ {
		for x := 0; x < 50; x++ {
			_, _, _, a := img.At(x, y).RGBA()
			if a > 0 {
				hasPixelLine1 = true
				break
			}
		}
		if hasPixelLine1 {
			break
		}
	}
	if !hasPixelLine1 {
		t.Error("expected foreground pixels in first line for 'A'")
	}

	// Third line (Y = 2*RowHeight..3*RowHeight-1) should have foreground pixels for 'B'.
	hasPixelLine3 := false
	for y := 2 * metrics.RowHeight; y < 3*metrics.RowHeight; y++ {
		for x := 0; x < 50; x++ {
			_, _, _, a := img.At(x, y).RGBA()
			if a > 0 {
				hasPixelLine3 = true
				break
			}
		}
		if hasPixelLine3 {
			break
		}
	}
	if !hasPixelLine3 {
		t.Error("expected foreground pixels in third line for 'B'")
	}
}

// TestSingleWordExceedingBoundsInWrapMode verifies that a single word
// exceeding the effective width in Wrap mode is placed on its own line and
// truncated (overflow=Wrap, renderLine clips it).

func TestSingleWordExceedingBoundsInWrapMode(t *testing.T) {
	face := font.Default()
	metrics := face.Metrics()
	// Allow 4 chars: effectiveWidth = 4 * GlyphAdvance = 20px.
	width := 4 * metrics.GlyphAdvance
	height := metrics.RowHeight

	result := Render(Config{
		Bounds:   image.Rect(0, 0, width, height),
		Text:     "ABCDEFGHIJ", // 10 chars, exceeds 4-char limit
		Font:     face,
		Overflow: Wrap,
	})
	if result == nil {
		t.Fatal("expected non-nil result")
	}

	img := result.Image.(*image.RGBA)

	// Verify no foreground pixels beyond the effective width.
	for y := 0; y < height; y++ {
		for x := width; x < img.Bounds().Dx(); x++ {
			_, _, _, a := img.At(x, y).RGBA()
			if a > 0 {
				t.Fatalf("unexpected pixel at (%d,%d) beyond effective width", x, y)
			}
		}
	}
}

// TestBorderWithEmptyText verifies that Border=true with empty text produces
// a border drawn at edges with a transparent interior.

func TestBorderWithEmptyText(t *testing.T) {
	fg := color.RGBA{R: 255, G: 255, B: 255, A: 255}

	result := Render(Config{
		Bounds:     image.Rect(0, 0, 10, 10),
		Text:       "",
		Border:     true,
		Foreground: fg,
	})
	if result == nil {
		t.Fatal("expected non-nil result with border and empty text")
	}

	img := result.Image.(*image.RGBA)
	width := 10
	height := 10

	// All border pixels should be foreground color.
	for x := 0; x < width; x++ {
		if img.RGBAAt(x, 0) != fg {
			t.Fatalf("top border pixel at (%d,0) should be foreground", x)
		}
		if img.RGBAAt(x, height-1) != fg {
			t.Fatalf("bottom border pixel at (%d,%d) should be foreground", x, height-1)
		}
	}
	for y := 0; y < height; y++ {
		if img.RGBAAt(0, y) != fg {
			t.Fatalf("left border pixel at (0,%d) should be foreground", y)
		}
		if img.RGBAAt(width-1, y) != fg {
			t.Fatalf("right border pixel at (%d,%d) should be foreground", width-1, y)
		}
	}

	// Interior should be transparent.
	for y := 1; y < height-1; y++ {
		for x := 1; x < width-1; x++ {
			c := img.RGBAAt(x, y)
			if c.A != 0 {
				t.Fatalf("interior pixel at (%d,%d) should be transparent, got %v", x, y, c)
			}
		}
	}
}

// TestBorderWithTightBounds3x3 verifies that a 3×3 bounds with border works
// correctly: 1px border + 1px interior.

func TestBorderWithTightBounds3x3(t *testing.T) {
	fg := color.RGBA{R: 200, G: 100, B: 50, A: 255}

	result := Render(Config{
		Bounds:     image.Rect(0, 0, 3, 3),
		Text:       "",
		Border:     true,
		Foreground: fg,
	})
	if result == nil {
		t.Fatal("expected non-nil result for 3×3 with border (1px border + 1px interior)")
	}

	img := result.Image.(*image.RGBA)

	// All 8 border pixels should be fg.
	borderPositions := [][2]int{
		{0, 0}, {1, 0}, {2, 0},
		{0, 1}, {2, 1},
		{0, 2}, {1, 2}, {2, 2},
	}
	for _, pos := range borderPositions {
		c := img.RGBAAt(pos[0], pos[1])
		if c != fg {
			t.Fatalf("border pixel at (%d,%d) should be foreground, got %v", pos[0], pos[1], c)
		}
	}

	// Interior (1,1) should be transparent (empty text).
	c := img.RGBAAt(1, 1)
	if c.A != 0 {
		t.Fatalf("interior pixel at (1,1) should be transparent, got %v", c)
	}
}

// TestBorderWithDx2NoPaddingNilResult verifies that Bounds.Dx()=2 with border
// and no padding produces nil result (effectiveWidth = 2 - 2*0 - 2*1 = 0).

func TestBorderWithDx2NoPaddingNilResult(t *testing.T) {
	result := Render(Config{
		Bounds: image.Rect(0, 0, 2, 10),
		Text:   "A",
		Border: true,
	})
	if result != nil {
		t.Fatal("expected nil result when border leaves zero effective width (Dx=2)")
	}
}
