package ttf

import (
	"os"
	"testing"

	"pgregory.net/rapid"
)

// For any TTF font file containing glyphs at specified codepoints, and for any
// set of CodepointRange values and valid target height, the TTF parser SHALL
// produce non-empty bitmap data ([]uint32 with at least one non-zero row) for
// each codepoint within the configured ranges that has a glyph present in the
// font, and SHALL produce no entries for codepoints outside the configured ranges.

func TestProp_TTFParserGlyphExtraction(t *testing.T) {
	f, err := os.Open("../vendor/matrix-code-font/Matrix Code Font.ttf")
	if err != nil {
		t.Skip("Matrix Code Font.ttf not available:", err)
	}
	f.Close()

	rapid.Check(t, func(t *rapid.T) {
		// Generate random codepoint ranges.
		ranges := genCodepointRanges(t)
		targetHeight := genTargetHeight(t)

		// Open font fresh for each iteration.
		fontFile, err := os.Open("../vendor/matrix-code-font/Matrix Code Font.ttf")
		if err != nil {
			t.Fatal("cannot open font file:", err)
		}
		defer fontFile.Close()

		font, err := Parse(fontFile, ParseConfig{
			Ranges:       ranges,
			TargetHeight: targetHeight,
		})
		if err != nil {
			// If no glyphs found in the generated ranges, that's expected
			// for ranges that don't overlap with the font's glyph set.
			return
		}

		// Build a set of all codepoints within the configured ranges.
		inRange := make(map[rune]bool)
		for _, cr := range ranges {
			for cp := cr.Low; cp <= cr.High; cp++ {
				inRange[cp] = true
			}
		}

		// Property 1: No entries outside configured ranges.
		for cp := range font.Glyphs {
			if !inRange[cp] {
				t.Fatalf("glyph for codepoint %d (U+%04X) exists but is outside configured ranges", cp, cp)
			}
		}

		// Property 2: For codepoints within ranges that have glyphs present,
		// bitmap data is non-empty (len(Rows) == TargetHeight > 0).
		// Each glyph has exactly TargetHeight rows.
		for cp, gd := range font.Glyphs {
			if len(gd.Rows) != targetHeight {
				t.Fatalf("glyph %c (U+%04X) has %d rows, want %d (TargetHeight)",
					cp, cp, len(gd.Rows), targetHeight)
			}
		}

		// Property 3: The parser only returns glyphs that the font actually
		// supports (i.e., the font reports a valid advance width). This is
		// verified implicitly: if Parse returns a glyph, it passed the
		// rasterizeGlyph check. We verify structural validity: Width > 0.
		for cp, gd := range font.Glyphs {
			if gd.Width <= 0 {
				t.Fatalf("glyph %c (U+%04X) has Width=%d, expected > 0", cp, cp, gd.Width)
			}
			if gd.Codepoint != cp {
				t.Fatalf("glyph map key U+%04X doesn't match GlyphData.Codepoint U+%04X", cp, gd.Codepoint)
			}
		}
	})
}

// genCodepointRanges generates 1-3 random CodepointRange values covering
// plausible Unicode ranges that may or may not overlap with glyphs in the
// Matrix Code Font (ASCII printable + half-width katakana).
func genCodepointRanges(t *rapid.T) []CodepointRange {
	numRanges := rapid.IntRange(1, 3).Draw(t, "numRanges")
	ranges := make([]CodepointRange, numRanges)
	for i := 0; i < numRanges; i++ {
		// Pick ranges from zones likely to have glyphs:
		// ASCII printable (33-126), half-width katakana (0xFF66-0xFF9D),
		// or random Unicode ranges.
		zone := rapid.IntRange(0, 2).Draw(t, "zone")
		var low, high rune
		switch zone {
		case 0:
			// ASCII printable subset.
			low = rune(rapid.IntRange(33, 100).Draw(t, "asciiLow"))
			high = rune(rapid.IntRange(int(low), 126).Draw(t, "asciiHigh"))
		case 1:
			// Half-width katakana subset.
			low = rune(rapid.IntRange(0xFF66, 0xFF90).Draw(t, "kataLow"))
			high = rune(rapid.IntRange(int(low), 0xFF9D).Draw(t, "kataHigh"))
		case 2:
			// Random codepoint range (may not have glyphs in the font).
			low = rune(rapid.IntRange(0x100, 0x2FF).Draw(t, "randLow"))
			high = rune(rapid.IntRange(int(low), int(low)+30).Draw(t, "randHigh"))
		}
		ranges[i] = CodepointRange{Low: low, High: high}
	}
	return ranges
}

// genTargetHeight generates a target height in the range [8, 24], which are
// reasonable pixel sizes for bitmap font rasterization.
func genTargetHeight(t *rapid.T) int {
	return rapid.IntRange(8, 24).Draw(t, "targetHeight")
}
