package ttf

import (
	"os"
	"strings"
	"testing"
)

func TestParse_ErrorCases(t *testing.T) {
	t.Run("targetHeight <= 0", func(t *testing.T) {
		r := strings.NewReader("not a font")
		_, err := Parse(r, ParseConfig{
			TargetHeight: 0,
			Ranges:       []CodepointRange{{Low: 32, High: 126}},
		})
		if err == nil {
			t.Fatal("expected error for targetHeight <= 0")
		}
		if !strings.Contains(err.Error(), "targetHeight must be > 0") {
			t.Fatalf("unexpected error message: %s", err)
		}
	})

	t.Run("negative targetHeight", func(t *testing.T) {
		r := strings.NewReader("not a font")
		_, err := Parse(r, ParseConfig{
			TargetHeight: -5,
			Ranges:       []CodepointRange{{Low: 32, High: 126}},
		})
		if err == nil {
			t.Fatal("expected error for negative targetHeight")
		}
	})

	t.Run("no ranges specified", func(t *testing.T) {
		r := strings.NewReader("not a font")
		_, err := Parse(r, ParseConfig{
			TargetHeight: 12,
			Ranges:       nil,
		})
		if err == nil {
			t.Fatal("expected error for no ranges")
		}
		if !strings.Contains(err.Error(), "no codepoint ranges") {
			t.Fatalf("unexpected error message: %s", err)
		}
	})

	t.Run("invalid font data", func(t *testing.T) {
		r := strings.NewReader("this is not a valid TTF file")
		_, err := Parse(r, ParseConfig{
			TargetHeight: 12,
			Ranges:       []CodepointRange{{Low: 32, High: 126}},
		})
		if err == nil {
			t.Fatal("expected error for invalid font data")
		}
		if !strings.Contains(err.Error(), "parsing font data") {
			t.Fatalf("unexpected error message: %s", err)
		}
	})

	t.Run("no glyphs in range", func(t *testing.T) {
		// Use the real font but request codepoints that don't exist in it.
		f, err := os.Open("../vendor/matrix-code-font/Matrix Code Font.ttf")
		if err != nil {
			t.Skip("Matrix Code Font.ttf not available:", err)
		}
		defer f.Close()

		// Request a range where no glyphs exist (private use area).
		_, err = Parse(f, ParseConfig{
			TargetHeight: 12,
			Ranges:       []CodepointRange{{Low: 0xF0000, High: 0xF000F}},
		})
		if err == nil {
			t.Fatal("expected error for no glyphs in range")
		}
		if !strings.Contains(err.Error(), "no glyphs found") {
			t.Fatalf("unexpected error message: %s", err)
		}
	})
}

func TestParse_BasicFunctionality(t *testing.T) {
	f, err := os.Open("../vendor/matrix-code-font/Matrix Code Font.ttf")
	if err != nil {
		t.Skip("Matrix Code Font.ttf not available:", err)
	}
	defer f.Close()

	font, err := Parse(f, ParseConfig{
		TargetHeight: 12,
		Ranges:       []CodepointRange{{Low: 'A', High: 'Z'}},
	})
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Verify basic properties.
	if font.GlyphHeight != 12 {
		t.Errorf("GlyphHeight = %d, want 12", font.GlyphHeight)
	}

	if font.GlyphWidth <= 0 {
		t.Error("GlyphWidth should be > 0")
	}

	if font.GlyphWidth > 32 {
		t.Errorf("GlyphWidth = %d, exceeds uint32 bit width", font.GlyphWidth)
	}

	// Should have at least some glyphs for A-Z.
	if len(font.Glyphs) == 0 {
		t.Fatal("expected at least some glyphs for A-Z")
	}

	// Verify each glyph has correct structure.
	for cp, gd := range font.Glyphs {
		if gd.Codepoint != cp {
			t.Errorf("glyph map key %c doesn't match GlyphData.Codepoint %c", cp, gd.Codepoint)
		}
		if len(gd.Rows) != 12 {
			t.Errorf("glyph %c has %d rows, want 12", cp, len(gd.Rows))
		}
		if gd.Width <= 0 {
			t.Errorf("glyph %c has width %d, want > 0", cp, gd.Width)
		}

		// At least one row should have ink for visible characters.
		hasInk := false
		for _, row := range gd.Rows {
			if row != 0 {
				hasInk = true
				break
			}
		}
		if !hasInk && cp != ' ' {
			t.Errorf("glyph %c has no ink in any row", cp)
		}
	}
}

func TestParse_MultipleRanges(t *testing.T) {
	f, err := os.Open("../vendor/matrix-code-font/Matrix Code Font.ttf")
	if err != nil {
		t.Skip("Matrix Code Font.ttf not available:", err)
	}
	defer f.Close()

	font, err := Parse(f, ParseConfig{
		TargetHeight: 12,
		Ranges: []CodepointRange{
			{Low: 'A', High: 'C'},
			{Low: '0', High: '2'},
		},
	})
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Should have glyphs from both ranges.
	for _, cp := range []rune{'A', 'B', 'C', '0', '1', '2'} {
		if _, ok := font.Glyphs[cp]; !ok {
			t.Errorf("expected glyph for %c", cp)
		}
	}

	// Should NOT have glyphs outside the ranges.
	for _, cp := range []rune{'D', 'E', '3', '4'} {
		if _, ok := font.Glyphs[cp]; ok {
			t.Errorf("did not expect glyph for %c", cp)
		}
	}
}
