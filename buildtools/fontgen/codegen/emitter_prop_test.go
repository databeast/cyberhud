package codegen

import (
	"bytes"
	"fmt"
	"testing"

	"pgregory.net/rapid"
)

// For any valid EmitConfig (random package name, font ID, struct/const/array names,
// glyph dimensions, and glyph row data), invoking Emit twice in sequence with the
// same configuration produces byte-for-byte identical output.

func TestEmitOutputIdempotence(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate valid Go identifiers using regex: starts with lowercase letter,
		// followed by 1-12 alphanumeric characters.
		packageName := rapid.StringMatching(`[a-z][a-z0-9]{1,12}`).Draw(t, "packageName")
		structName := rapid.StringMatching(`[A-Za-z][A-Za-z0-9]{1,12}`).Draw(t, "structName")
		constName := rapid.StringMatching(`[A-Z][A-Za-z0-9]{1,12}`).Draw(t, "constName")
		arrayName := rapid.StringMatching(`[a-z][a-z0-9]{1,12}`).Draw(t, "arrayName")

		// Generate glyph dimensions within valid ranges.
		glyphWidth := rapid.IntRange(1, 32).Draw(t, "glyphWidth")
		glyphHeight := rapid.IntRange(1, 64).Draw(t, "glyphHeight")
		glyphAdvance := rapid.IntRange(1, 64).Draw(t, "glyphAdvance")
		rowHeight := rapid.IntRange(1, 128).Draw(t, "rowHeight")

		// Generate a font ID like "font-NxM".
		fontID := fmt.Sprintf("font-%dx%d", glyphWidth, glyphHeight)

		// Generate random glyph map data with a subset of ASCII codepoints.
		numGlyphs := rapid.IntRange(1, 95).Draw(t, "numGlyphs")
		glyphMap := make(map[rune][]uint32, numGlyphs)
		for i := 0; i < numGlyphs; i++ {
			cp := rune(32 + i)
			rows := make([]uint32, glyphHeight)
			for r := 0; r < glyphHeight; r++ {
				rows[r] = rapid.Uint32().Draw(t, fmt.Sprintf("g%d_r%d", i, r))
			}
			glyphMap[cp] = rows
		}

		cfg := EmitConfig{
			PackageName:  packageName,
			FontID:       fontID,
			StructName:   structName,
			ConstName:    constName,
			ArrayName:    arrayName,
			GlyphWidth:   glyphWidth,
			GlyphHeight:  glyphHeight,
			GlyphAdvance: glyphAdvance,
			RowHeight:    rowHeight,
			GlyphMap:     glyphMap,
		}

		// Call Emit twice with the same config.
		var buf1, buf2 bytes.Buffer
		err1 := Emit(&buf1, cfg)
		err2 := Emit(&buf2, cfg)

		// Both calls must succeed.
		if err1 != nil {
			t.Fatalf("first Emit failed: %v", err1)
		}
		if err2 != nil {
			t.Fatalf("second Emit failed: %v", err2)
		}

		// Assert byte-for-byte identical output.
		if !bytes.Equal(buf1.Bytes(), buf2.Bytes()) {
			t.Fatalf("two sequential Emit calls produced different output (not idempotent)\n"+
				"first output length: %d\nsecond output length: %d",
				buf1.Len(), buf2.Len())
		}
	})
}
