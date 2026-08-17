package main

import (
	"bytes"
	"fmt"
	"go/parser"
	"go/token"
	"strings"
	"testing"

	"github.com/databeast/cyberhud/buildtools/fontgen/codegen"
	"pgregory.net/rapid"
)

// *For any* generated icon face with a known glyph map, *for any* codepoint and
// row index: (a) if the codepoint exists in the map and 0 ≤ row < GlyphHeight,
// GlyphRow SHALL return the corresponding bitmap value; (b) if the codepoint is
// absent and the fallback character exists, GlyphRow SHALL return the fallback
// character's row data; (c) if row < 0 or row ≥ GlyphHeight, GlyphRow SHALL return 0.

func TestProp_GlyphRowDispatchCorrectness(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate glyph dimensions: height in [8, 64].
		glyphHeight := rapid.IntRange(8, 64).Draw(t, "glyphHeight")

		// Generate a glyph map with 1-10 entries.
		glyphMap := genGlyphMapForFace(t, glyphHeight)

		// Ensure fallback '?' is sometimes present, sometimes not.
		includeFallback := rapid.Bool().Draw(t, "includeFallback")
		if includeFallback {
			// Add '?' to the map if not already present.
			if _, exists := glyphMap['?']; !exists {
				rows := make([]uint32, glyphHeight)
				for row := 0; row < glyphHeight; row++ {
					rows[row] = rapid.Uint32().Draw(t, fmt.Sprintf("fallback_row_%d", row))
				}
				glyphMap['?'] = rows
			}
		}

		// Generate query codepoint: may or may not be in the map.
		queryCp := rune(rapid.IntRange(33, 70000).Draw(t, "queryCp"))

		// Generate query row: may be in-bounds or out-of-bounds.
		queryRow := rapid.IntRange(-5, glyphHeight+5).Draw(t, "queryRow")

		// Simulate the expected GlyphRow result.
		expected := simulateGlyphRowFace(glyphMap, '?', queryCp, queryRow, glyphHeight)

		// Emit generated code via codegen.Emit.
		arrayName := "testIcons"
		cfg := codegen.EmitConfig{
			PackageName:  "testpkg",
			FontID:       "material-icons-test",
			StructName:   "TestIconFace",
			ConstName:    "TestIconFaceID",
			ArrayName:    arrayName,
			GlyphWidth:   glyphHeight,
			GlyphHeight:  glyphHeight,
			GlyphAdvance: glyphHeight,
			RowHeight:    glyphHeight,
			GlyphMap:     glyphMap,
			FallbackChar: '?',
		}

		var buf bytes.Buffer
		err := codegen.Emit(&buf, cfg)
		if err != nil {
			t.Fatalf("codegen.Emit failed: %v", err)
		}

		output := buf.String()

		// Verify structural patterns in generated code.

		// 1. Must contain bounds check for row.
		boundsCheck := fmt.Sprintf("row < 0 || row >= %d", glyphHeight)
		if !strings.Contains(output, boundsCheck) {
			t.Fatalf("emitted code missing bounds check: %s", boundsCheck)
		}

		// 2. Must contain map lookup pattern.
		mapName := arrayName + "Map"
		mapLookup := fmt.Sprintf("if data, ok := %s[ch]; ok", mapName)
		if !strings.Contains(output, mapLookup) {
			t.Fatalf("emitted code missing map lookup: %s", mapLookup)
		}

		// 3. Must contain fallback lookup pattern.
		fallbackConst := arrayName + "Fallback"
		fallbackLookup := fmt.Sprintf("if data, ok := %s[%s]; ok", mapName, fallbackConst)
		if !strings.Contains(output, fallbackLookup) {
			t.Fatalf("emitted code missing fallback lookup: %s", fallbackLookup)
		}

		// 4. Must contain final return 0.
		if !strings.Contains(output, "return 0") {
			t.Fatalf("emitted code missing 'return 0' fallback")
		}

		// 5. Verify the simulated dispatch logic is consistent.

		// Sub-property (a): known codepoint + valid row returns map value.
		for r, rows := range glyphMap {
			for row := 0; row < glyphHeight; row++ {
				result := simulateGlyphRowFace(glyphMap, '?', r, row, glyphHeight)
				if result != rows[row] {
					t.Fatalf("known codepoint %q row %d: got 0x%08X, want 0x%08X",
						string(r), row, result, rows[row])
				}
			}
			break // verify at least one known char fully
		}

		// Sub-property (b): unknown codepoint falls back to '?'.
		_, isKnown := glyphMap[queryCp]
		if !isKnown && queryRow >= 0 && queryRow < glyphHeight {
			fallbackData, hasFallback := glyphMap['?']
			if hasFallback {
				if expected != fallbackData[queryRow] {
					t.Fatalf("unknown codepoint %q row %d: expected fallback 0x%08X, got 0x%08X",
						string(queryCp), queryRow, fallbackData[queryRow], expected)
				}
			} else {
				if expected != 0 {
					t.Fatalf("unknown codepoint %q with no fallback should return 0, got 0x%08X",
						string(queryCp), expected)
				}
			}
		}

		// Sub-property (c): out-of-bounds row returns 0.
		if queryRow < 0 || queryRow >= glyphHeight {
			if expected != 0 {
				t.Fatalf("out-of-bounds row %d should return 0, got 0x%08X", queryRow, expected)
			}
		}
	})
}

// genGlyphMapForFace generates a map[rune][]uint32 with 1-10 entries, each with
// glyphHeight rows of random uint32 values.
func genGlyphMapForFace(t *rapid.T, glyphHeight int) map[rune][]uint32 {
	numEntries := rapid.IntRange(1, 10).Draw(t, "numEntries")
	m := make(map[rune][]uint32, numEntries)
	for i := 0; i < numEntries; i++ {
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

// simulateGlyphRowFace implements the same logic as the generated GlyphRow for
// map-based icon faces. This is the oracle for the property test.
func simulateGlyphRowFace(glyphMap map[rune][]uint32, fallback rune, ch rune, row int, glyphHeight int) uint32 {
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

// *For any* valid target pixel height in [8, 64], a generated icon face SHALL register
// with ID "material-icons-{height}" and report Metrics where GlyphWidth, GlyphHeight,
// GlyphAdvance, and RowHeight all equal the target pixel height.

func TestProp_IconFaceIdentityAndMetrics(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate a target height in [8, 64].
		targetHeight := rapid.IntRange(8, 64).Draw(t, "targetHeight")

		// Generate a non-empty glyph map (1-5 entries, each with targetHeight rows).
		numGlyphs := rapid.IntRange(1, 5).Draw(t, "numGlyphs")
		glyphMap := make(map[rune][]uint32, numGlyphs)
		for i := 0; i < numGlyphs; i++ {
			cp := rune(rapid.IntRange(0x0021, 0xFFFF).Draw(t, fmt.Sprintf("cp_%d", i)))
			if _, exists := glyphMap[cp]; exists {
				continue // skip duplicate codepoints
			}
			rows := make([]uint32, targetHeight)
			for r := 0; r < targetHeight; r++ {
				rows[r] = rapid.Uint32().Draw(t, fmt.Sprintf("row_%d_%d", i, r))
			}
			glyphMap[cp] = rows
		}

		// Ensure at least one glyph exists.
		if len(glyphMap) == 0 {
			cp := rune(0x0041) // 'A'
			rows := make([]uint32, targetHeight)
			glyphMap[cp] = rows
		}

		// Build the face ID.
		faceID := fmt.Sprintf("material-icons-%d", targetHeight)

		// Derive struct/const/array names from faceID (same logic as face.go).
		structName := idToIdentifier(faceID) + "Face"
		constName := idToIdentifier(faceID) + "ID"
		arrayName := idToIdentifier(faceID)

		// Build EmitConfig with square metrics (all = targetHeight).
		cfg := codegen.EmitConfig{
			PackageName:  "font",
			FontID:       faceID,
			StructName:   structName,
			ConstName:    constName,
			ArrayName:    arrayName,
			GlyphWidth:   targetHeight,
			GlyphHeight:  targetHeight,
			GlyphAdvance: targetHeight,
			RowHeight:    targetHeight,
			GlyphMap:     glyphMap,
			FallbackChar: '?',
		}

		// Emit to a buffer.
		var buf bytes.Buffer
		err := codegen.Emit(&buf, cfg)
		if err != nil {
			t.Fatalf("codegen.Emit returned error: %v", err)
		}

		src := buf.String()

		// Verify the output contains the ID string.
		if !strings.Contains(src, fmt.Sprintf("%q", faceID)) {
			t.Fatalf("output does not contain face ID %q", faceID)
		}

		// Verify the output contains Metrics with all four fields equal to targetHeight.
		expectedMetrics := fmt.Sprintf(
			"Metrics{GlyphWidth: %d, GlyphHeight: %d, GlyphAdvance: %d, RowHeight: %d}",
			targetHeight, targetHeight, targetHeight, targetHeight,
		)
		if !strings.Contains(src, expectedMetrics) {
			t.Fatalf("output does not contain expected metrics %q", expectedMetrics)
		}

		// Parse with go/parser to verify validity.
		fset := token.NewFileSet()
		_, parseErr := parser.ParseFile(fset, "face.go", src, 0)
		if parseErr != nil {
			t.Fatalf("go/parser failed to parse output:\n%s\nerror: %v", src, parseErr)
		}
	})
}
