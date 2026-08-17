package bdf

import (
	"strings"
	"testing"
)

// minimalBDF is a minimal valid BDF file with one ASCII glyph (space, codepoint 32).
const minimalBDF = `STARTFONT 2.1
FONT -Test-Font
SIZE 8 72 72
FONTBOUNDINGBOX 5 8 0 -1
STARTCHAR space
ENCODING 32
SWIDTH 500 0
DWIDTH 5 0
BBX 5 8 0 -1
BITMAP
00
00
00
00
00
00
00
00
ENDCHAR
ENDFONT
`

func TestParseMinimalBDF(t *testing.T) {
	f, err := Parse(strings.NewReader(minimalBDF))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Name != "-Test-Font" {
		t.Errorf("expected font name '-Test-Font', got %q", f.Name)
	}
	if len(f.Glyphs) != 1 {
		t.Fatalf("expected 1 glyph, got %d", len(f.Glyphs))
	}
	g, ok := f.Glyphs[32]
	if !ok {
		t.Fatal("expected glyph at codepoint 32")
	}
	if g.Width != 5 || g.Height != 8 {
		t.Errorf("expected width=5 height=8, got width=%d height=%d", g.Width, g.Height)
	}
	if g.DWidth != 5 {
		t.Errorf("expected DWidth=5, got %d", g.DWidth)
	}
	for i, row := range g.Rows {
		if row != 0 {
			t.Errorf("row %d: expected 0, got 0x%08X", i, row)
		}
	}
}

// testBDFWithA has the 'A' character (codepoint 65) with known pixel data.
const testBDFWithA = `STARTFONT 2.1
FONT -Test-Font
SIZE 8 72 72
FONTBOUNDINGBOX 8 16 0 -2
STARTCHAR A
ENCODING 65
SWIDTH 500 0
DWIDTH 8 0
BBX 8 8 0 0
BITMAP
18
24
42
7E
42
42
42
00
ENDCHAR
ENDFONT
`

func TestParseBDFGlyphA(t *testing.T) {
	f, err := Parse(strings.NewReader(testBDFWithA))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	g, ok := f.Glyphs['A']
	if !ok {
		t.Fatal("expected glyph 'A'")
	}
	if g.Width != 8 || g.Height != 8 {
		t.Errorf("expected width=8 height=8, got width=%d height=%d", g.Width, g.Height)
	}
	if len(g.Rows) != 8 {
		t.Fatalf("expected 8 rows, got %d", len(g.Rows))
	}

	// With BBX 8 8 0 0 (xoff=0), hex "18" → 0x18 = 0001 1000
	// Shifted left by (32 - 8) = 24 bits → 0x18000000
	expected := []uint32{
		0x18000000, // 18
		0x24000000, // 24
		0x42000000, // 42
		0x7E000000, // 7E
		0x42000000, // 42
		0x42000000, // 42
		0x42000000, // 42
		0x00000000, // 00
	}
	for i, exp := range expected {
		if g.Rows[i] != exp {
			t.Errorf("row %d: expected 0x%08X, got 0x%08X", i, exp, g.Rows[i])
		}
	}
}

func TestParseSkipsNonASCII(t *testing.T) {
	bdf := `STARTFONT 2.1
FONT -Test-Font
SIZE 8 72 72
FONTBOUNDINGBOX 5 8 0 0
STARTCHAR space
ENCODING 32
DWIDTH 5 0
BBX 5 8 0 0
BITMAP
00
00
00
00
00
00
00
00
ENDCHAR
STARTCHAR nonascii
ENCODING 200
DWIDTH 5 0
BBX 5 8 0 0
BITMAP
FF
FF
FF
FF
FF
FF
FF
FF
ENDCHAR
STARTCHAR tilde
ENCODING 126
DWIDTH 5 0
BBX 5 8 0 0
BITMAP
00
00
00
00
00
00
00
00
ENDCHAR
STARTCHAR highuni
ENCODING 9999
DWIDTH 5 0
BBX 5 8 0 0
BITMAP
00
00
00
00
00
00
00
00
ENDCHAR
ENDFONT
`
	f, err := Parse(strings.NewReader(bdf))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(f.Glyphs) != 2 {
		t.Errorf("expected 2 glyphs (32 and 126), got %d", len(f.Glyphs))
	}
	if _, ok := f.Glyphs[32]; !ok {
		t.Error("expected glyph 32")
	}
	if _, ok := f.Glyphs[126]; !ok {
		t.Error("expected glyph 126")
	}
	if _, ok := f.Glyphs[200]; ok {
		t.Error("glyph 200 should have been skipped")
	}
}

func TestParseMissingSTARTFONT(t *testing.T) {
	bdf := `FONT -Test-Font
STARTCHAR space
ENCODING 32
BBX 5 8 0 0
BITMAP
00
ENDCHAR
`
	_, err := Parse(strings.NewReader(bdf))
	if err == nil {
		t.Fatal("expected error for missing STARTFONT")
	}
	if !strings.Contains(err.Error(), "not a valid BDF file (missing STARTFONT)") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestParseInvalidHex(t *testing.T) {
	bdf := `STARTFONT 2.1
FONT -Test-Font
STARTCHAR excl
ENCODING 33
DWIDTH 5 0
BBX 5 8 0 0
BITMAP
ZZ
00
00
00
00
00
00
00
ENDCHAR
ENDFONT
`
	_, err := Parse(strings.NewReader(bdf))
	if err == nil {
		t.Fatal("expected error for invalid hex")
	}
	if !strings.Contains(err.Error(), "invalid hex") {
		t.Errorf("unexpected error message: %v", err)
	}
	if !strings.Contains(err.Error(), "U+0021") {
		t.Errorf("expected error to contain glyph codepoint U+0021: %v", err)
	}
}

func TestParseWrongRowCount(t *testing.T) {
	bdf := `STARTFONT 2.1
FONT -Test-Font
STARTCHAR excl
ENCODING 33
DWIDTH 5 0
BBX 5 8 0 0
BITMAP
00
00
00
ENDCHAR
ENDFONT
`
	_, err := Parse(strings.NewReader(bdf))
	if err == nil {
		t.Fatal("expected error for wrong row count")
	}
	if !strings.Contains(err.Error(), "has 3 bitmap rows, expected 8") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestParseNoASCIIGlyphs(t *testing.T) {
	bdf := `STARTFONT 2.1
FONT -Test-Font
STARTCHAR nonascii
ENCODING 200
DWIDTH 5 0
BBX 5 8 0 0
BITMAP
00
00
00
00
00
00
00
00
ENDCHAR
ENDFONT
`
	_, err := Parse(strings.NewReader(bdf))
	if err == nil {
		t.Fatal("expected error for no ASCII glyphs")
	}
	if !strings.Contains(err.Error(), "no glyphs found in configured ranges") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestParseBBXWithXOffset(t *testing.T) {
	// A glyph with xoff=2: data should be shifted right by 2 from MSB.
	bdf := `STARTFONT 2.1
FONT -Test-Font
STARTCHAR excl
ENCODING 33
DWIDTH 8 0
BBX 6 4 2 0
BITMAP
FC
84
84
FC
ENDCHAR
ENDFONT
`
	f, err := Parse(strings.NewReader(bdf))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	g := f.Glyphs['!']

	// hex "FC" = 11111100 in 8 bits
	// Shift left by (32 - 8) = 24 → 0xFC000000
	// Then shift right by xoff=2 → 0x3F000000
	expected := uint32(0x3F000000)
	if g.Rows[0] != expected {
		t.Errorf("row 0: expected 0x%08X, got 0x%08X", expected, g.Rows[0])
	}
}

func TestParseEmptyInput(t *testing.T) {
	_, err := Parse(strings.NewReader(""))
	if err == nil {
		t.Fatal("expected error for empty input")
	}
	if !strings.Contains(err.Error(), "not a valid BDF file (missing STARTFONT)") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestParseNegativeXOffset(t *testing.T) {
	// A glyph with xoff=-1: data should be shifted left by 1 more from MSB position.
	bdf := `STARTFONT 2.1
FONT -Test-Font
STARTCHAR excl
ENCODING 33
DWIDTH 8 0
BBX 6 1 -1 0
BITMAP
7C
ENDCHAR
ENDFONT
`
	f, err := Parse(strings.NewReader(bdf))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	g := f.Glyphs['!']

	// hex "7C" = 0111 1100 in 8 bits
	// Shift left by (32 - 8) = 24 → 0x7C000000
	// Then shift left by 1 (xoff=-1) → 0xF8000000
	expected := uint32(0xF8000000)
	if g.Rows[0] != expected {
		t.Errorf("row 0: expected 0x%08X, got 0x%08X", expected, g.Rows[0])
	}
}

func TestParseWideGlyph16px(t *testing.T) {
	// A 16-pixel wide glyph uses 4 hex chars per row.
	bdf := `STARTFONT 2.1
FONT -Test-Font
STARTCHAR A
ENCODING 65
DWIDTH 16 0
BBX 16 2 0 0
BITMAP
FFFF
8001
ENDCHAR
ENDFONT
`
	f, err := Parse(strings.NewReader(bdf))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	g := f.Glyphs['A']

	// hex "FFFF" = 16 bits all set
	// Shift left by (32 - 16) = 16 → 0xFFFF0000
	if g.Rows[0] != 0xFFFF0000 {
		t.Errorf("row 0: expected 0xFFFF0000, got 0x%08X", g.Rows[0])
	}
	// hex "8001" = 1000 0000 0000 0001
	// Shift left by (32 - 16) = 16 → 0x80010000
	if g.Rows[1] != 0x80010000 {
		t.Errorf("row 1: expected 0x80010000, got 0x%08X", g.Rows[1])
	}
}
