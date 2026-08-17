package bdf

import (
	"fmt"
	"strings"
	"testing"

	"pgregory.net/rapid"
)

// For any valid BDF glyph in the ASCII 32–126 range with arbitrary width (1-32),
// height (1-16), and x-offset (≥0), parsing the BDF bitmap hex data and producing
// a uint32 row value yields a bitmask where bit positions correspond exactly to the
// pixel positions defined in the source BDF data, with MSB representing the leftmost pixel.

func TestBDFParseRoundTrip(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate random valid glyph parameters.
		width := rapid.IntRange(1, 32).Draw(t, "width")
		height := rapid.IntRange(1, 16).Draw(t, "height")
		xoff := rapid.IntRange(0, 4).Draw(t, "xoff")
		yoff := rapid.IntRange(-2, 2).Draw(t, "yoff")
		codepoint := rapid.IntRange(32, 126).Draw(t, "codepoint")
		dwidth := width + 1

		// Number of hex chars per row: ceil(width / 8) * 2 (byte-padded).
		hexBytes := (width + 7) / 8
		hexChars := hexBytes * 2
		hexBits := hexChars * 4

		// Generate random bitmap rows as uint32 values constrained to hexBits bits.
		// In BDF, the leftmost pixel is the MSB of the first byte.
		rows := make([]uint32, height)
		hexRows := make([]string, height)
		for i := 0; i < height; i++ {
			var maxVal uint32
			if hexBits >= 32 {
				maxVal = 0xFFFFFFFF
			} else {
				maxVal = (1 << uint(hexBits)) - 1
			}
			val := uint32(rapid.Uint32Range(0, maxVal).Draw(t, fmt.Sprintf("row%d", i)))
			rows[i] = val
			hexRows[i] = fmt.Sprintf("%0*X", hexChars, val)
		}

		// Construct a minimal valid BDF file.
		var bdf strings.Builder
		bdf.WriteString("STARTFONT 2.1\n")
		bdf.WriteString("FONT TestFont\n")
		bdf.WriteString("SIZE 16 72 72\n")
		bdf.WriteString(fmt.Sprintf("FONTBOUNDINGBOX %d %d %d %d\n", width, height, xoff, yoff))
		bdf.WriteString("CHARS 1\n")
		bdf.WriteString("STARTCHAR testglyph\n")
		bdf.WriteString(fmt.Sprintf("ENCODING %d\n", codepoint))
		bdf.WriteString(fmt.Sprintf("SWIDTH %d 0\n", 500))
		bdf.WriteString(fmt.Sprintf("DWIDTH %d 0\n", dwidth))
		bdf.WriteString(fmt.Sprintf("BBX %d %d %d %d\n", width, height, xoff, yoff))
		bdf.WriteString("BITMAP\n")
		for _, hr := range hexRows {
			bdf.WriteString(hr + "\n")
		}
		bdf.WriteString("ENDCHAR\n")
		bdf.WriteString("ENDFONT\n")

		// Parse the BDF.
		font, err := Parse(strings.NewReader(bdf.String()))
		if err != nil {
			t.Fatalf("Parse failed: %v\nBDF:\n%s", err, bdf.String())
		}

		// Verify the glyph was parsed.
		glyph, ok := font.Glyphs[rune(codepoint)]
		if !ok {
			t.Fatalf("glyph for codepoint %d not found in parsed font", codepoint)
		}

		// Verify glyph metadata.
		if glyph.Width != width {
			t.Errorf("width: got %d, want %d", glyph.Width, width)
		}
		if glyph.Height != height {
			t.Errorf("height: got %d, want %d", glyph.Height, height)
		}
		if glyph.XOff != xoff {
			t.Errorf("xoff: got %d, want %d", glyph.XOff, xoff)
		}
		if glyph.YOff != yoff {
			t.Errorf("yoff: got %d, want %d", glyph.YOff, yoff)
		}
		if glyph.DWidth != dwidth {
			t.Errorf("dwidth: got %d, want %d", glyph.DWidth, dwidth)
		}

		// Verify parsed row values match expected bit positions.
		// The parser does: shifted = val << (32 - hexBits), then shifted >>= xoff (if xoff > 0).
		if len(glyph.Rows) != height {
			t.Fatalf("row count: got %d, want %d", len(glyph.Rows), height)
		}
		for i, val := range rows {
			// Replicate the parser's alignment logic.
			expected := val << uint(32-hexBits)
			if xoff > 0 {
				expected >>= uint(xoff)
			}

			if glyph.Rows[i] != expected {
				t.Errorf("row %d: got 0x%08X, want 0x%08X (source hex val=0x%X, hexBits=%d, xoff=%d)",
					i, glyph.Rows[i], expected, val, hexBits, xoff)
			}
		}

		// Additionally verify that specific pixel bit positions are correct.
		// For the MSB-left convention: bit (hexBits-1-col) in the source value
		// corresponds to pixel column 'col' (0 = leftmost).
		// After alignment: pixel col is at bit (31 - xoff - col) in the parsed uint32.
		for i, val := range rows {
			for col := 0; col < width; col++ {
				// Source bit for pixel column 'col'.
				srcBit := (val >> uint(hexBits-1-col)) & 1

				// Parsed bit for pixel column 'col'.
				parsedBitPos := 31 - xoff - col
				if parsedBitPos < 0 || parsedBitPos > 31 {
					continue // Pixel shifted out of uint32 range
				}
				parsedBit := (glyph.Rows[i] >> uint(parsedBitPos)) & 1

				if srcBit != parsedBit {
					t.Errorf("row %d, pixel col %d: source bit=%d, parsed bit=%d "+
						"(val=0x%X, parsed=0x%08X, hexBits=%d, xoff=%d)",
						i, col, srcBit, parsedBit, val, glyph.Rows[i], hexBits, xoff)
				}
			}
		}
	})
}
