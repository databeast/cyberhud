package regionid

import (
	"testing"

	"pgregory.net/rapid"
)

// For any valid region ID (surface name matching [a-z][a-z0-9-]* and non-negative
// index), parsing the canonical string representation <surface>.<index> and converting
// back to string should produce the original string.

// validSurfaceFirstChars contains all characters valid as the first character of a surface name.
const validSurfaceFirstChars = "abcdefghijklmnopqrstuvwxyz"

// validSurfaceTailChars contains all characters valid in the tail of a surface name.
const validSurfaceTailChars = "abcdefghijklmnopqrstuvwxyz0123456789-"

func TestProperty_RegionID_RoundTrip(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		// Generate a valid surface name matching [a-z][a-z0-9-]*
		// First character: lowercase letter [a-z]
		firstIdx := rapid.IntRange(0, len(validSurfaceFirstChars)-1).Draw(rt, "firstIdx")
		firstChar := validSurfaceFirstChars[firstIdx]

		// Remaining characters: lowercase letters, digits, or hyphens
		tailLen := rapid.IntRange(0, 20).Draw(rt, "tailLen")
		tail := make([]byte, tailLen)
		for i := range tail {
			charIdx := rapid.IntRange(0, len(validSurfaceTailChars)-1).Draw(rt, "tailCharIdx")
			tail[i] = validSurfaceTailChars[charIdx]
		}
		surface := string(firstChar) + string(tail)

		// Generate a non-negative index
		index := rapid.IntRange(0, 10000).Draw(rt, "index")

		// Construct the ID
		original := ID{Surface: surface, Index: index}

		// Round-trip: String() → Parse() should recover the original ID
		canonical := original.String()
		parsed, err := Parse(canonical)
		if err != nil {
			t.Fatalf("Parse(%q) returned error: %v (original: %+v)", canonical, err, original)
		}

		if parsed != original {
			t.Fatalf("round-trip failed: Parse(%q) = %+v, want %+v", canonical, parsed, original)
		}
	})
}
