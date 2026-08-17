package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"sort"
	"strings"
	"testing"

	"pgregory.net/rapid"
)

// *For any* non-empty set of icon entries with unique post-conversion constant names,
// EmitConstants SHALL produce output that is parseable by go/parser as valid Go source,
// contains exactly one const declaration per entry, and the constant declarations appear
// in lexicographic order by constant name.

func TestProp_ConstantsFileValidityAndOrdering(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate 1-20 unique snake_case names ensuring no collisions after SnakeToPascal.
		numEntries := rapid.IntRange(1, 20).Draw(t, "numEntries")

		var entries []IconEntry
		seenSnake := make(map[string]bool)
		seenPascal := make(map[string]bool)

		for i := 0; i < numEntries; i++ {
			for {
				name := genUniqueSnakeName(t, i)
				pascal := SnakeToPascal(name)
				if !seenSnake[name] && !seenPascal[pascal] {
					seenSnake[name] = true
					seenPascal[pascal] = true
					cp := rune(rapid.Uint32Range(0x0000, 0x10FFFF).Draw(t, fmt.Sprintf("cp_%d", i)))
					entries = append(entries, IconEntry{Name: name, Codepoint: cp})
					break
				}
			}
		}

		// Call EmitConstants to a buffer.
		var buf bytes.Buffer
		err := EmitConstants(&buf, "icons", entries)
		if err != nil {
			t.Fatalf("EmitConstants returned error: %v", err)
		}

		src := buf.String()

		// Assert: output is parseable by go/parser.
		fset := token.NewFileSet()
		file, parseErr := parser.ParseFile(fset, "constants.go", src, 0)
		if parseErr != nil {
			t.Fatalf("go/parser failed to parse output:\n%s\nerror: %v", src, parseErr)
		}

		// Extract constant names from the AST.
		var constNames []string
		for _, decl := range file.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok || genDecl.Tok != token.CONST {
				continue
			}
			for _, spec := range genDecl.Specs {
				vs, ok := spec.(*ast.ValueSpec)
				if !ok {
					continue
				}
				for _, ident := range vs.Names {
					constNames = append(constNames, ident.Name)
				}
			}
		}

		// Assert: count of constants equals number of entries.
		if len(constNames) != numEntries {
			t.Fatalf("expected %d constants, got %d", numEntries, len(constNames))
		}

		// Assert: constants appear in lexicographic order.
		if !sort.StringsAreSorted(constNames) {
			t.Fatalf("constants are not in lexicographic order: %v", constNames)
		}

		// Assert: all expected constant names are present.
		expectedNames := make([]string, numEntries)
		for i, e := range entries {
			expectedNames[i] = SnakeToPascal(e.Name)
		}
		sort.Strings(expectedNames)

		for i, expected := range expectedNames {
			if constNames[i] != expected {
				t.Fatalf("constant[%d]: expected %q, got %q", i, expected, constNames[i])
			}
		}
	})
}

// genUniqueSnakeName generates a valid snake_case icon name for uniqueness testing.
// Each name consists of 1-3 lowercase segments separated by underscores.
func genUniqueSnakeName(t *rapid.T, idx int) string {
	numSegments := rapid.IntRange(1, 3).Draw(t, fmt.Sprintf("numSeg_%d", idx))
	segments := make([]string, numSegments)
	for i := range segments {
		segLen := rapid.IntRange(2, 5).Draw(t, fmt.Sprintf("segLen_%d_%d", idx, i))
		seg := make([]byte, segLen)
		for j := range seg {
			seg[j] = byte(rapid.IntRange('a', 'z').Draw(t, fmt.Sprintf("ch_%d_%d_%d", idx, i, j)))
		}
		segments[i] = string(seg)
	}
	return strings.Join(segments, "_")
}
