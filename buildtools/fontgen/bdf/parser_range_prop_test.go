package bdf

import (
	"fmt"
	"strings"
	"testing"

	"pgregory.net/rapid"
)

// For any BDF font data containing glyphs at various codepoints, and for any
// set of CodepointRange values, ParseWithConfig SHALL retain only glyphs whose
// codepoints fall within at least one configured range, and discard all others.

// genCodepointRange generates a single valid CodepointRange with Low <= High.
func genCodepointRange(t *rapid.T, label string) CodepointRange {
	low := rune(rapid.IntRange(33, 0x1000).Draw(t, label+"_low"))
	high := rune(rapid.IntRange(int(low), int(low)+100).Draw(t, label+"_high"))
	return CodepointRange{Low: low, High: high}
}

// genCodepointRanges generates a slice of 1-4 non-overlapping CodepointRange values.
func genCodepointRanges(t *rapid.T) []CodepointRange {
	count := rapid.IntRange(1, 4).Draw(t, "rangeCount")
	ranges := make([]CodepointRange, count)
	// Use well-separated base offsets to avoid overlap.
	bases := []rune{33, 200, 400, 600}
	for i := 0; i < count; i++ {
		base := bases[i]
		span := rapid.IntRange(1, 20).Draw(t, fmt.Sprintf("span_%d", i))
		ranges[i] = CodepointRange{Low: base, High: base + rune(span)}
	}
	return ranges
}

// buildMinimalBDFGlyph creates BDF text for a single glyph at the given codepoint.
func buildMinimalBDFGlyph(cp rune) string {
	var b strings.Builder
	b.WriteString("STARTCHAR glyph\n")
	b.WriteString(fmt.Sprintf("ENCODING %d\n", cp))
	b.WriteString("SWIDTH 500 0\n")
	b.WriteString("DWIDTH 5 0\n")
	b.WriteString("BBX 4 4 0 0\n")
	b.WriteString("BITMAP\n")
	b.WriteString("F0\n")
	b.WriteString("90\n")
	b.WriteString("90\n")
	b.WriteString("F0\n")
	b.WriteString("ENDCHAR\n")
	return b.String()
}

// buildBDFFont creates a complete BDF font string with glyphs at the given codepoints.
func buildBDFFont(codepoints []rune) string {
	var b strings.Builder
	b.WriteString("STARTFONT 2.1\n")
	b.WriteString("FONT TestRangeFont\n")
	b.WriteString("SIZE 16 72 72\n")
	b.WriteString("FONTBOUNDINGBOX 4 4 0 0\n")
	b.WriteString(fmt.Sprintf("CHARS %d\n", len(codepoints)))
	for _, cp := range codepoints {
		b.WriteString(buildMinimalBDFGlyph(cp))
	}
	b.WriteString("ENDFONT\n")
	return b.String()
}

func TestProp_BDFParserRangeFiltering(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate random codepoint ranges for filtering.
		ranges := genCodepointRanges(t)

		// Generate a set of codepoints: some inside ranges, some outside.
		// First, collect codepoints guaranteed to be inside the ranges.
		var insideCPs []rune
		for _, r := range ranges {
			// Pick 1-3 codepoints within each range.
			count := rapid.IntRange(1, 3).Draw(t, fmt.Sprintf("inside_count_%d_%d", r.Low, r.High))
			for j := 0; j < count; j++ {
				cp := rune(rapid.IntRange(int(r.Low), int(r.High)).Draw(t, fmt.Sprintf("inside_cp_%d_%d_%d", r.Low, r.High, j)))
				insideCPs = append(insideCPs, cp)
			}
		}

		// Generate codepoints guaranteed to be outside all ranges.
		var outsideCPs []rune
		outsideCount := rapid.IntRange(1, 5).Draw(t, "outsideCount")
		for i := 0; i < outsideCount; i++ {
			// Pick codepoints well above the highest range to guarantee they're outside.
			cp := rune(rapid.IntRange(900, 1100).Draw(t, fmt.Sprintf("outside_cp_%d", i)))
			// Verify it's actually outside all ranges.
			if inRanges(cp, ranges) {
				continue // Skip if it accidentally falls in range.
			}
			outsideCPs = append(outsideCPs, cp)
		}

		// Combine all codepoints and build a BDF font.
		allCPs := append(insideCPs, outsideCPs...)

		// Deduplicate codepoints.
		seen := make(map[rune]bool)
		var uniqueCPs []rune
		for _, cp := range allCPs {
			if !seen[cp] {
				seen[cp] = true
				uniqueCPs = append(uniqueCPs, cp)
			}
		}

		if len(uniqueCPs) == 0 {
			return // Degenerate case, skip.
		}

		bdfData := buildBDFFont(uniqueCPs)

		// Parse with the generated ranges.
		font, err := ParseWithConfig(strings.NewReader(bdfData), ParseConfig{Ranges: ranges})

		// Determine expected glyphs (those in at least one range).
		var expectedInside []rune
		for _, cp := range uniqueCPs {
			if inRanges(cp, ranges) {
				expectedInside = append(expectedInside, cp)
			}
		}

		// If no glyphs match the ranges, ParseWithConfig should return an error.
		if len(expectedInside) == 0 {
			if err == nil {
				t.Fatalf("expected error when no glyphs match ranges, got nil")
			}
			return
		}

		if err != nil {
			t.Fatalf("ParseWithConfig failed: %v\nBDF:\n%s", err, bdfData)
		}

		// Property 1: Every retained glyph must be within at least one configured range.
		for cp := range font.Glyphs {
			if !inRanges(cp, ranges) {
				t.Errorf("glyph U+%04X was retained but is not within any configured range", cp)
			}
		}

		// Property 2: Every codepoint within a configured range that was present in the
		// input must appear in the output.
		for _, cp := range uniqueCPs {
			if inRanges(cp, ranges) {
				if _, ok := font.Glyphs[cp]; !ok {
					t.Errorf("glyph U+%04X is within a configured range but was not retained", cp)
				}
			}
		}

		// Property 3: No glyphs outside the configured ranges should be present.
		for _, cp := range outsideCPs {
			if _, ok := font.Glyphs[cp]; ok {
				t.Errorf("glyph U+%04X is outside all configured ranges but was retained", cp)
			}
		}
	})
}
