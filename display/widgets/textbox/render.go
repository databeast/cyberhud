package textbox

import (
	"image"
	"image/color"

	"github.com/databeast/cyberhud/display/surface/fonts"
)

// ellipsis is the Unicode character used as a truncation indicator.
const ellipsis = '…' // U+2026

// renderLine renders a single line of text onto img at the given Y position.
// It handles horizontal alignment and overflow (truncate/clip) within the
// effective width. The effectiveX is the left edge of the text area, and
// effectiveWidth is the available horizontal space.
func renderLine(img *image.RGBA, face font.Face, text string, effectiveX, y, effectiveWidth int, align Alignment, overflow Overflow, fg color.RGBA) {
	if effectiveWidth <= 0 || len(text) == 0 {
		return
	}

	if face == nil {
		return
	}
	metrics := face.Metrics()
	glyphAdvance := metrics.GlyphAdvance
	if glyphAdvance <= 0 {
		return
	}

	// Maximum number of characters that fit in the effective width.
	maxChars := effectiveWidth / glyphAdvance
	if maxChars <= 0 {
		// No characters fit at all — render nothing.
		return
	}

	// Determine what text (and how many chars) to actually render.
	runes := []rune(text)
	textLen := len(runes)

	var renderRunes []rune

	if textLen <= maxChars {
		// Text fits entirely — no overflow handling needed.
		renderRunes = runes
	} else {
		// Text exceeds available width — apply overflow strategy.
		switch overflow {
		case Truncate:
			renderRunes = truncateLine(runes, maxChars)
		case Clip:
			renderRunes = runes[:maxChars]
		default:
			// Wrap mode should not reach renderLine (layout.go handles it),
			// but if it does, treat like clip.
			renderRunes = runes[:maxChars]
		}
	}

	if len(renderRunes) == 0 {
		return
	}

	// Compute pixel width of the text to render.
	textPixelWidth := len(renderRunes) * glyphAdvance

	// Compute horizontal start position based on alignment.
	var xStart int
	switch align {
	case Center:
		xStart = effectiveX + (effectiveWidth-textPixelWidth)/2
	case Right:
		xStart = effectiveX + effectiveWidth - textPixelWidth
	default: // Left
		xStart = effectiveX
	}

	// Render each character using drawGlyph.
	maxX := effectiveX + effectiveWidth
	xCursor := xStart
	for _, ch := range renderRunes {
		if xCursor >= maxX {
			break
		}
		drawGlyph(img, face, ch, xCursor, y, fg)
		xCursor += glyphAdvance
	}
}

// truncateLine returns the runes to render when text exceeds maxChars in
// Truncate mode. It renders leading characters plus an ellipsis (U+2026)
// such that the total count does not exceed maxChars.
//
// Edge cases:
//   - maxChars >= 2: render (maxChars-1) leading chars + ellipsis
//   - maxChars == 1: render only the ellipsis
//   - maxChars == 0: render nothing (caller handles this before calling)
func truncateLine(runes []rune, maxChars int) []rune {
	if maxChars <= 0 {
		return nil
	}
	if maxChars == 1 {
		// Only the ellipsis fits.
		return []rune{ellipsis}
	}
	// Leading chars + ellipsis.
	result := make([]rune, maxChars)
	copy(result, runes[:maxChars-1])
	result[maxChars-1] = ellipsis
	return result
}
