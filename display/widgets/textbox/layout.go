package textbox

import (
	"strings"

	"github.com/databeast/cyberhud/display/surface/fonts"
)

// visualLine represents a single rendered line after layout.
type visualLine struct {
	text string
	face font.Face
	y    int // Y offset relative to text area origin
}

// computeLayout takes the text, effective dimensions, font, overflow mode,
// VAlign, LineSpacing, and FontOverrides, and returns the visual lines to render.
// The Y positions in each visualLine are relative to the text area origin
// (the caller adds padding/border offset).
func computeLayout(text string, effectiveWidth, effectiveHeight int, defaultFont font.Face, overflow Overflow, valign VAlign, lineSpacing int, fontOverrides []font.Face) []visualLine {
	if effectiveWidth <= 0 || effectiveHeight <= 0 {
		return nil
	}

	// Normalize lineSpacing: negative treated as 0.
	if lineSpacing < 0 {
		lineSpacing = 0
	}

	// Step 1: Split text at newline characters into logical lines.
	logicalLines := strings.Split(text, "\n")

	// Step 2: Resolve per-line font and expand wraps into visual lines.
	type pending struct {
		text string
		face font.Face
	}
	var pendingLines []pending

	for i, line := range logicalLines {
		face := effectiveFont(fontOverrides, i, defaultFont)

		if overflow == Wrap {
			wrapped := wrapLine(line, face, effectiveWidth)
			for _, wl := range wrapped {
				pendingLines = append(pendingLines, pending{text: wl, face: face})
			}
		} else {
			pendingLines = append(pendingLines, pending{text: line, face: face})
		}
	}

	// Step 3: Compute total text block height for vertical alignment.
	blockHeight := 0
	for i, pl := range pendingLines {
		blockHeight += pl.face.Metrics().RowHeight
		if i < len(pendingLines)-1 {
			blockHeight += lineSpacing
		}
	}

	// Step 4: Determine vertical start offset based on VAlign.
	startY := computeStartY(blockHeight, effectiveHeight, valign)

	// Step 5: Assign Y positions and clip lines that exceed effective height.
	var result []visualLine
	y := startY
	for _, pl := range pendingLines {
		rowHeight := pl.face.Metrics().RowHeight

		// Stop if this line would exceed the effective vertical bounds.
		if y+rowHeight > effectiveHeight {
			break
		}

		// Skip lines that start above the visible area (shouldn't happen
		// with current logic, but guard against negative startY).
		if y < 0 {
			y += rowHeight + lineSpacing
			continue
		}

		result = append(result, visualLine{
			text: pl.text,
			face: pl.face,
			y:    y,
		})

		y += rowHeight + lineSpacing
	}

	return result
}

// effectiveFont returns the font to use for the given logical line index.
// If fontOverrides has a non-nil entry at index i, that is used; otherwise
// the defaultFont is returned.
func effectiveFont(fontOverrides []font.Face, i int, defaultFont font.Face) font.Face {
	if i < len(fontOverrides) && fontOverrides[i] != nil {
		return fontOverrides[i]
	}
	return defaultFont
}

// wrapLine breaks a single logical line into visual lines that fit within
// effectiveWidth using word-wrap at whitespace boundaries. If a single word
// exceeds the available width, it is placed on its own line (renderLine will
// truncate it). An empty input line produces a single empty visual line.
func wrapLine(line string, face font.Face, effectiveWidth int) []string {
	glyphAdvance := face.Metrics().GlyphAdvance
	if glyphAdvance <= 0 {
		return []string{line}
	}

	maxChars := effectiveWidth / glyphAdvance
	if maxChars <= 0 {
		// Nothing fits — still produce one line (renderLine will handle it).
		return []string{line}
	}

	// Empty line → produce one empty visual line.
	if len(line) == 0 {
		return []string{""}
	}

	words := splitWords(line)
	var lines []string
	var currentLine strings.Builder

	currentLen := 0 // character count of current line

	for _, word := range words {
		wordLen := len([]rune(word))

		if currentLen == 0 {
			// First word on this line.
			if wordLen > maxChars {
				// Word exceeds line width — put it on its own line.
				lines = append(lines, word)
				// currentLine stays empty, currentLen stays 0.
			} else {
				currentLine.WriteString(word)
				currentLen = wordLen
			}
		} else {
			// Adding space + word to current line.
			needed := 1 + wordLen // space + word
			if currentLen+needed <= maxChars {
				currentLine.WriteRune(' ')
				currentLine.WriteString(word)
				currentLen += needed
			} else {
				// Flush current line and start a new one.
				lines = append(lines, currentLine.String())
				currentLine.Reset()

				if wordLen > maxChars {
					// Word exceeds line width — put it on its own line.
					lines = append(lines, word)
					currentLen = 0
				} else {
					currentLine.WriteString(word)
					currentLen = wordLen
				}
			}
		}
	}

	// Flush any remaining content.
	if currentLen > 0 {
		lines = append(lines, currentLine.String())
	}

	// If the input was all whitespace, splitWords may have produced nothing,
	// but we still need at least one visual line.
	if len(lines) == 0 {
		lines = []string{""}
	}

	return lines
}

// splitWords splits a string into words by whitespace. Consecutive whitespace
// is treated as a single separator. Leading/trailing whitespace is ignored.
func splitWords(s string) []string {
	return strings.Fields(s)
}

// computeStartY determines the Y offset for the first line based on
// vertical alignment, block height, and effective height.
func computeStartY(blockHeight, effectiveHeight int, valign VAlign) int {
	// If text block exceeds effective height, always start from top.
	if blockHeight > effectiveHeight {
		return 0
	}

	switch valign {
	case Middle:
		return (effectiveHeight - blockHeight) / 2
	case Bottom:
		return effectiveHeight - blockHeight
	default: // Top
		return 0
	}
}
