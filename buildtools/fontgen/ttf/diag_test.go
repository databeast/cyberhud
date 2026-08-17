package ttf

import (
	"fmt"
	"os"
	"testing"
)

func TestDiagMatrixFont(t *testing.T) {
	f, err := os.Open("../vendor/matrix-code-font/Matrix Code Font.ttf")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	// Try various target heights to find one that produces visible glyphs.
	for _, targetHeight := range []int{12, 16, 20, 24, 32, 48} {
		cfg := ParseConfig{
			Ranges: []CodepointRange{
				{Low: 33, High: 126},
				{Low: 65382, High: 65437},
			},
			TargetHeight: targetHeight,
		}

		// Re-open file for each test
		f.Seek(0, 0)

		font, err := Parse(f, cfg)
		if err != nil {
			t.Logf("targetHeight=%d: error: %v", targetHeight, err)
			continue
		}

		nonZeroGlyphs := 0
		for _, gd := range font.Glyphs {
			hasInk := false
			for _, row := range gd.Rows {
				if row != 0 {
					hasInk = true
					break
				}
			}
			if hasInk {
				nonZeroGlyphs++
			}
		}

		fmt.Printf("targetHeight=%d: total glyphs=%d, non-zero glyphs=%d, maxWidth=%d\n",
			targetHeight, len(font.Glyphs), nonZeroGlyphs, font.GlyphWidth)

		// Print a sample glyph at this height
		if gd, ok := font.Glyphs['A']; ok {
			fmt.Printf("  'A' width=%d rows:\n", gd.Width)
			for i, row := range gd.Rows {
				fmt.Printf("    row %2d: 0x%08X\n", i, row)
			}
		}
	}
}
