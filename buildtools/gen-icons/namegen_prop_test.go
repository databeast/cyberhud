package main

import (
	"go/token"
	"strings"
	"testing"
	"unicode"

	"pgregory.net/rapid"
)

// *For any* valid snake_case icon name (composed of one or more lowercase segments
// separated by underscores), SnakeToPascal SHALL produce a string that starts with
// "Icon", followed by each segment with its first character capitalized, contains no
// underscores, and is a valid Go identifier.

func TestProp_SnakeToPascalConversion(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate a valid snake_case name: 1-5 segments joined by underscores.
		// Each segment: first char a-z, rest 1-9 chars from a-z or 0-9.
		numSegments := rapid.IntRange(1, 5).Draw(t, "numSegments")
		segments := make([]string, numSegments)
		for i := range segments {
			// First char must be a letter (a-z)
			firstChar := rapid.ByteRange('a', 'z').Draw(t, "firstChar")
			// Remaining chars: 0-9 more chars from a-z or 0-9
			restLen := rapid.IntRange(0, 9).Draw(t, "restLen")
			rest := make([]byte, restLen)
			for j := range rest {
				if rapid.Bool().Draw(t, "isDigit") {
					rest[j] = rapid.ByteRange('0', '9').Draw(t, "digit")
				} else {
					rest[j] = rapid.ByteRange('a', 'z').Draw(t, "letter")
				}
			}
			segments[i] = string(firstChar) + string(rest)
		}
		name := strings.Join(segments, "_")

		result := SnakeToPascal(name)

		// 1. Result starts with "Icon"
		if !strings.HasPrefix(result, "Icon") {
			t.Fatalf("SnakeToPascal(%q) = %q: does not start with \"Icon\"", name, result)
		}

		// 2. Result contains no underscores
		if strings.Contains(result, "_") {
			t.Fatalf("SnakeToPascal(%q) = %q: contains underscore", name, result)
		}

		// 3. Character at position 4 (after "Icon") is uppercase
		afterIcon := result[4:]
		if len(afterIcon) == 0 {
			t.Fatalf("SnakeToPascal(%q) = %q: nothing after \"Icon\"", name, result)
		}
		firstRune := rune(afterIcon[0])
		if !unicode.IsUpper(firstRune) {
			t.Fatalf("SnakeToPascal(%q) = %q: char after \"Icon\" is not uppercase (got %c)", name, result, firstRune)
		}

		// 4. Result is a valid Go identifier
		if !token.IsIdentifier(result) {
			t.Fatalf("SnakeToPascal(%q) = %q: not a valid Go identifier", name, result)
		}

		// 5. For each segment, the corresponding position in result has uppercase first char
		remaining := afterIcon
		for _, seg := range segments {
			if len(remaining) == 0 {
				t.Fatalf("SnakeToPascal(%q) = %q: ran out of result chars while checking segments", name, result)
			}
			r := rune(remaining[0])
			if !unicode.IsUpper(r) && !unicode.IsDigit(r) {
				// If the segment starts with a letter, it must be uppercase in the result
				if seg[0] >= 'a' && seg[0] <= 'z' {
					t.Fatalf("SnakeToPascal(%q) = %q: segment %q first char not capitalized in result", name, result, seg)
				}
			}
			// Advance past this segment in the result
			remaining = remaining[len(seg):]
		}
	})
}
