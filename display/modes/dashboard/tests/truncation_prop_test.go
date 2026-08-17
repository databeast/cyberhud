package tests_test

import (
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/databeast/cyberhud/display/surface/textlayout"
	"pgregory.net/rapid"
)

// For any text row and any maxChars > 0, if the row text exceeds maxChars runes,
// the output SHALL be exactly maxChars runes long with the final rune being U+2026
// (ellipsis). If the row text does not exceed maxChars runes, the output SHALL equal
// the trimmed input.
func TestProperty_TextTruncationEllipsis(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate random text with varying lengths (0–100 chars).
		// Include a mix of ASCII and multi-byte characters to exercise rune-based logic.
		text := rapid.StringMatching(`[\x20-\x7E\x{00C0}-\x{00FF}\x{4E00}-\x{4E10}]{0,100}`).Draw(t, "text")

		// Generate random maxChars in range [1, 50].
		maxChars := rapid.IntRange(1, 50).Draw(t, "maxChars")

		// Call the Truncate function under test.
		result := textlayout.Truncate(text, maxChars)

		// The function first trims whitespace, so compute the trimmed input.
		trimmed := strings.TrimSpace(text)
		trimmedRuneCount := utf8.RuneCountInString(trimmed)
		resultRuneCount := utf8.RuneCountInString(result)

		if trimmedRuneCount > maxChars {
			// Property: output is exactly maxChars runes long.
			if resultRuneCount != maxChars {
				t.Fatalf("text exceeds maxChars: expected result rune count %d, got %d (text=%q, trimmed=%q, maxChars=%d, result=%q)",
					maxChars, resultRuneCount, text, trimmed, maxChars, result)
			}

			// Property: the last rune is U+2026 (ellipsis '…').
			runes := []rune(result)
			lastRune := runes[len(runes)-1]
			if lastRune != '\u2026' {
				t.Fatalf("expected last rune to be U+2026 (ellipsis), got U+%04X (text=%q, maxChars=%d, result=%q)",
					lastRune, text, maxChars, result)
			}

			// Additional check: the prefix before ellipsis matches the first maxChars-1 runes of trimmed.
			expectedPrefix := string([]rune(trimmed)[:maxChars-1])
			actualPrefix := string(runes[:len(runes)-1])
			if actualPrefix != expectedPrefix {
				t.Fatalf("prefix mismatch: expected %q, got %q (text=%q, maxChars=%d)",
					expectedPrefix, actualPrefix, text, maxChars)
			}
		} else {
			// Property: output equals the trimmed input (no truncation needed).
			if result != trimmed {
				t.Fatalf("text within maxChars: expected result=%q to equal trimmed=%q (text=%q, maxChars=%d)",
					result, trimmed, text, maxChars)
			}
		}
	})
}
