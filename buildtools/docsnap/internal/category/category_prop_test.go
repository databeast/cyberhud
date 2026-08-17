package category

import (
	"testing"

	"pgregory.net/rapid"
)

// For any filename prefixed with "grayscale-fast-", Match SHALL return (Grayscale, true).
// For "degraded-" prefix, Match SHALL return ("", false) since no PrefixMapping for "degraded-" exists after the removal.

func TestProperty4_MatchGrayscaleFastPrefix(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		suffix := rapid.StringMatching(`[a-z0-9]{1,20}`).Draw(rt, "suffix")
		filename := "grayscale-fast-" + suffix

		cat, ok := Match(filename)
		if !ok {
			t.Fatalf("Match(%q) returned ok=false, want (Grayscale, true)", filename)
		}
		if cat != Grayscale {
			t.Fatalf("Match(%q) returned category=%q, want %q", filename, cat, Grayscale)
		}
	})
}

func TestProperty4_MatchDegradedPrefixNoMatch(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		suffix := rapid.StringMatching(`[a-z0-9]{1,20}`).Draw(rt, "suffix")
		filename := "degraded-" + suffix

		cat, ok := Match(filename)
		if ok {
			t.Fatalf("Match(%q) returned ok=true with category=%q, want (\"\", false)", filename, cat)
		}
		if cat != "" {
			t.Fatalf("Match(%q) returned category=%q, want \"\"", filename, cat)
		}
	})
}
