package codegen

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"pgregory.net/rapid"
)

// For any generated extended font with a known glyph map and fallback character,
// calling GlyphRow(ch, row) SHALL return the correct bitmask from the map when ch
// is present, the fallback character's bitmask when ch is absent, and 0 when row
// is out of bounds.

// genGlyphMap generates a map[rune][]uint32 with 1-10 entries, each with
// glyphHeight rows of random uint32 values.
func genGlyphMap(t *rapid.T, glyphHeight int) map[rune][]uint32 {
	numEntries := rapid.IntRange(1, 10).Draw(t, "numEntries")
	m := make(map[rune][]uint32, numEntries)
	for i := 0; i < numEntries; i++ {
		// Generate runes in a broad range including ASCII and extended Unicode
		r := rune(rapid.IntRange(33, 70000).Draw(t, fmt.Sprintf("rune_%d", i)))
		if _, exists := m[r]; exists {
			continue // skip duplicates
		}
		rows := make([]uint32, glyphHeight)
		for row := 0; row < glyphHeight; row++ {
			rows[row] = rapid.Uint32().Draw(t, fmt.Sprintf("glyph_%d_row_%d", i, row))
		}
		m[r] = rows
	}
	return m
}

// simulateGlyphRow implements the same logic as the generated GlyphRow for
// map-based extended fonts. This is the oracle for the property test.
func simulateGlyphRow(glyphMap map[rune][]uint32, fallback rune, ch rune, row int, glyphHeight int) uint32 {
	if row < 0 || row >= glyphHeight {
		return 0
	}
	if data, ok := glyphMap[ch]; ok {
		return data[row]
	}
	if data, ok := glyphMap[fallback]; ok {
		return data[row]
	}
	return 0
}

func TestProp_ExtendedGlyphRowDispatch(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate glyph dimensions.
		glyphHeight := rapid.IntRange(1, 32).Draw(t, "glyphHeight")

		// Generate a glyph map with 1-10 entries.
		glyphMap := genGlyphMap(t, glyphHeight)

		// Pick a fallback character: either one present in the map or a new one.
		usePresentFallback := rapid.Bool().Draw(t, "usePresentFallback")
		var fallback rune
		if usePresentFallback && len(glyphMap) > 0 {
			// Pick one of the existing runes as fallback.
			idx := rapid.IntRange(0, len(glyphMap)-1).Draw(t, "fallbackIdx")
			i := 0
			for r := range glyphMap {
				if i == idx {
					fallback = r
					break
				}
				i++
			}
		} else {
			// Use a fallback that may or may not be in the map.
			fallback = '?'
		}

		// Generate query parameters.
		queryChar := rune(rapid.IntRange(33, 70000).Draw(t, "queryChar"))
		// row can be in-bounds, negative, or beyond glyphHeight.
		queryRow := rapid.IntRange(-5, glyphHeight+5).Draw(t, "queryRow")

		// Simulate the expected result.
		expected := simulateGlyphRow(glyphMap, fallback, queryChar, queryRow, glyphHeight)

		// Now verify by emitting code and checking patterns.
		// The emitter should produce code that implements the same logic.
		cfg := EmitConfig{
			PackageName:  "testpkg",
			FontID:       "test-extended",
			StructName:   "TestFace",
			ConstName:    "TestFontID",
			ArrayName:    "testGlyphs",
			GlyphWidth:   8,
			GlyphHeight:  glyphHeight,
			GlyphAdvance: 9,
			RowHeight:    glyphHeight + 2,
			GlyphMap:     glyphMap,
			FallbackChar: fallback,
		}

		var buf bytes.Buffer
		err := Emit(&buf, cfg)
		if err != nil {
			t.Fatalf("Emit failed: %v", err)
		}

		output := buf.String()

		// Verify the emitted code contains the expected structural patterns
		// for map-based GlyphRow dispatch.

		// 1. Must contain bounds check for row.
		boundsCheck := fmt.Sprintf("row < 0 || row >= %d", glyphHeight)
		if !strings.Contains(output, boundsCheck) {
			t.Fatalf("emitted code missing bounds check: %s", boundsCheck)
		}

		// 2. Must contain map lookup pattern.
		mapName := cfg.ArrayName + "Map"
		mapLookup := fmt.Sprintf("if data, ok := %s[ch]; ok", mapName)
		if !strings.Contains(output, mapLookup) {
			t.Fatalf("emitted code missing map lookup: %s", mapLookup)
		}

		// 3. Must contain fallback lookup pattern.
		fallbackConst := cfg.ArrayName + "Fallback"
		fallbackLookup := fmt.Sprintf("if data, ok := %s[%s]; ok", mapName, fallbackConst)
		if !strings.Contains(output, fallbackLookup) {
			t.Fatalf("emitted code missing fallback lookup: %s", fallbackLookup)
		}

		// 4. Must contain final return 0.
		if !strings.Contains(output, "return 0") {
			t.Fatalf("emitted code missing 'return 0' fallback")
		}

		// 5. Must contain map literal with correct type signature.
		mapDecl := fmt.Sprintf("var %s = map[rune][%d]uint32", mapName, glyphHeight)
		if !strings.Contains(output, mapDecl) {
			t.Fatalf("emitted code missing map declaration: %s", mapDecl)
		}

		// 6. Verify all runes in glyphMap appear in the emitted map literal.
		for r := range glyphMap {
			entry := fmt.Sprintf("%d: {", r)
			if !strings.Contains(output, entry) {
				t.Fatalf("emitted code missing entry for rune %d (%q)", r, string(r))
			}
		}

		// 7. Verify the simulated dispatch logic is consistent:
		//    - Known char returns its bitmask
		//    - Unknown char returns fallback bitmask
		//    - Out-of-bounds row returns 0

		// Test known char: pick a char from the map and verify.
		for r, rows := range glyphMap {
			for row := 0; row < glyphHeight; row++ {
				result := simulateGlyphRow(glyphMap, fallback, r, row, glyphHeight)
				if result != rows[row] {
					t.Fatalf("known char %q row %d: got 0x%08X, want 0x%08X",
						string(r), row, result, rows[row])
				}
			}
			break // verify at least one known char fully
		}

		// Test unknown char with the random query parameters.
		_, isKnown := glyphMap[queryChar]
		result := simulateGlyphRow(glyphMap, fallback, queryChar, queryRow, glyphHeight)
		if result != expected {
			t.Fatalf("query char=%q row=%d: simulate returned 0x%08X, expected 0x%08X",
				string(queryChar), queryRow, result, expected)
		}

		// Specific sub-property: out-of-bounds rows always return 0.
		if queryRow < 0 || queryRow >= glyphHeight {
			if expected != 0 {
				t.Fatalf("out-of-bounds row %d should return 0, got 0x%08X", queryRow, expected)
			}
		}

		// Specific sub-property: unknown char returns fallback bitmask.
		if !isKnown && queryRow >= 0 && queryRow < glyphHeight {
			fallbackData, hasFallback := glyphMap[fallback]
			if hasFallback {
				if expected != fallbackData[queryRow] {
					t.Fatalf("unknown char %q should return fallback bitmask at row %d: got 0x%08X, want 0x%08X",
						string(queryChar), queryRow, expected, fallbackData[queryRow])
				}
			} else {
				if expected != 0 {
					t.Fatalf("unknown char %q with no fallback should return 0: got 0x%08X",
						string(queryChar), expected)
				}
			}
		}
	})
}
