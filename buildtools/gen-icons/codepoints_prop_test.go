package main

import (
	"fmt"
	"strings"
	"testing"

	"pgregory.net/rapid"
)

// *For any* codepoints input containing two or more lines with the same icon name
// but different codepoint values, ParseCodepoints SHALL return exactly one entry per
// unique name, and that entry's codepoint SHALL equal the value from the last occurrence
// of that name in the input.

func TestProp_DuplicateNameLastWins(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate 2-5 unique names.
		numNames := rapid.IntRange(2, 5).Draw(t, "numNames")
		names := make([]string, numNames)
		nameSet := make(map[string]bool)
		for i := 0; i < numNames; i++ {
			for {
				name := genValidIconName(t)
				if !nameSet[name] {
					names[i] = name
					nameSet[name] = true
					break
				}
			}
		}

		// For each name, generate 1-3 codepoints (occurrences).
		// At least one name must have >1 occurrence to test duplicates.
		type occurrence struct {
			name      string
			codepoint rune
		}

		var allOccurrences []occurrence
		lastCodepoint := make(map[string]rune)

		for _, name := range names {
			numOccurrences := rapid.IntRange(1, 3).Draw(t, fmt.Sprintf("occurrences_%s", name))
			for j := 0; j < numOccurrences; j++ {
				cp := rune(rapid.Uint32Range(0x0000, 0x10FFFF).Draw(t, fmt.Sprintf("cp_%s_%d", name, j)))
				allOccurrences = append(allOccurrences, occurrence{name: name, codepoint: cp})
				lastCodepoint[name] = cp
			}
		}

		// Ensure at least one name has duplicates (numOccurrences > 1).
		// If by chance all names got exactly 1 occurrence, add an extra occurrence for the first name.
		hasDuplicate := false
		countPerName := make(map[string]int)
		for _, occ := range allOccurrences {
			countPerName[occ.name]++
			if countPerName[occ.name] > 1 {
				hasDuplicate = true
			}
		}
		if !hasDuplicate {
			extraCp := rune(rapid.Uint32Range(0x0000, 0x10FFFF).Draw(t, "extraCp"))
			allOccurrences = append(allOccurrences, occurrence{name: names[0], codepoint: extraCp})
			lastCodepoint[names[0]] = extraCp
		}

		// Build the codepoints input text.
		var sb strings.Builder
		for _, occ := range allOccurrences {
			fmt.Fprintf(&sb, "%s %x\n", occ.name, occ.codepoint)
		}

		// Parse.
		entries, err := ParseCodepoints(strings.NewReader(sb.String()))
		if err != nil {
			t.Fatalf("ParseCodepoints returned error: %v", err)
		}

		// Assert: len(result) equals number of unique names.
		if len(entries) != numNames {
			t.Fatalf("expected %d unique entries, got %d", numNames, len(entries))
		}

		// Assert: for each unique name, the result entry's codepoint matches the last occurrence.
		entryMap := make(map[string]rune)
		for _, e := range entries {
			entryMap[e.Name] = e.Codepoint
		}

		for name, expectedCp := range lastCodepoint {
			got, ok := entryMap[name]
			if !ok {
				t.Fatalf("name %q not found in parsed entries", name)
			}
			if got != expectedCp {
				t.Fatalf("name %q: expected codepoint %U, got %U", name, expectedCp, got)
			}
		}
	})
}

// genValidIconName generates a valid snake_case icon name (lowercase ASCII + underscores,
// non-empty, no leading/trailing/consecutive underscores).
func genValidIconName(t *rapid.T) string {
	numSegments := rapid.IntRange(1, 3).Draw(t, "numSegments")
	segments := make([]string, numSegments)
	for i := range segments {
		segLen := rapid.IntRange(1, 6).Draw(t, fmt.Sprintf("segLen_%d", i))
		seg := make([]byte, segLen)
		for j := range seg {
			seg[j] = byte(rapid.IntRange('a', 'z').Draw(t, fmt.Sprintf("char_%d_%d", i, j)))
		}
		segments[i] = string(seg)
	}
	return strings.Join(segments, "_")
}
