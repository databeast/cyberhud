package tests

import (
	"strings"
	"testing"

	"github.com/databeast/cyberhud/hardware/panels"
	"pgregory.net/rapid"
)

// For any list of modes and excluded list, every mode NOT in the excluded list SHALL
// appear in the output in the same relative order as the input.
func TestProperty6_NonExcludedModesPreservedInOrder(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate a random list of mode names (0-20 items).
		modeCount := rapid.IntRange(0, 20).Draw(t, "modeCount")
		modes := make([]string, modeCount)
		for i := range modes {
			modes[i] = rapid.StringMatching(`[a-zA-Z0-9_\-]{1,15}`).Draw(t, "mode")
		}

		// Generate a random exclusion set (0-10 items).
		excludedCount := rapid.IntRange(0, 10).Draw(t, "excludedCount")
		excluded := make([]string, excludedCount)
		for i := range excluded {
			if modeCount > 0 && rapid.Bool().Draw(t, "pickFromModes") {
				idx := rapid.IntRange(0, modeCount-1).Draw(t, "idx")
				base := modes[idx]
				if rapid.Bool().Draw(t, "changeCase") {
					base = strings.ToUpper(base)
				}
				excluded[i] = base
			} else {
				excluded[i] = rapid.StringMatching(`[a-zA-Z0-9_\-]{1,15}`).Draw(t, "excludedMode")
			}
		}

		result := panels.TestFilterExcluded(modes, excluded)

		// Compute expected result manually, mirroring filterExcluded logic:
		// 1. Build the lowercase exclusion set.
		exSet := make(map[string]struct{}, len(excluded))
		for _, e := range excluded {
			e = strings.ToLower(strings.TrimSpace(e))
			if e != "" {
				exSet[e] = struct{}{}
			}
		}

		// 2. If exclusion set is effectively empty, expect a direct copy of modes.
		//    Otherwise, expect lowercased non-empty non-excluded modes in order.
		var expected []string
		if len(exSet) == 0 {
			expected = append([]string(nil), modes...)
		} else {
			for _, m := range modes {
				mLower := strings.ToLower(strings.TrimSpace(m))
				if mLower == "" {
					continue
				}
				if _, blocked := exSet[mLower]; blocked {
					continue
				}
				expected = append(expected, mLower)
			}
		}

		// Verify: result matches expected (same elements, same order).
		if len(result) != len(expected) {
			t.Fatalf("length mismatch: got %d, want %d\nmodes=%v\nexcluded=%v\nresult=%v\nexpected=%v",
				len(result), len(expected), modes, excluded, result, expected)
		}
		for i := range expected {
			if result[i] != expected[i] {
				t.Fatalf("order mismatch at index %d: got %q, want %q\nmodes=%v\nexcluded=%v\nresult=%v\nexpected=%v",
					i, result[i], expected[i], modes, excluded, result, expected)
			}
		}
	})
}

// For any list of modes and any list of excluded modes, the result of filterExcluded
// SHALL contain no element that appears in the excluded list (case-insensitive).
func TestProperty5_ExcludedModesNeverAppearInOutput(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate a random list of mode names (0-20 items, arbitrary strings).
		modeCount := rapid.IntRange(0, 20).Draw(t, "modeCount")
		modes := make([]string, modeCount)
		for i := range modes {
			modes[i] = rapid.StringMatching(`[a-zA-Z0-9_\-]{0,15}`).Draw(t, "mode")
		}

		// Generate a random exclusion set (0-10 items).
		// Use a mix of modes from the input list and arbitrary strings.
		excludedCount := rapid.IntRange(0, 10).Draw(t, "excludedCount")
		excluded := make([]string, excludedCount)
		for i := range excluded {
			if modeCount > 0 && rapid.Bool().Draw(t, "pickFromModes") {
				// Pick from existing modes (possibly with case variation)
				idx := rapid.IntRange(0, modeCount-1).Draw(t, "idx")
				base := modes[idx]
				if rapid.Bool().Draw(t, "changeCase") {
					base = strings.ToUpper(base)
				}
				excluded[i] = base
			} else {
				excluded[i] = rapid.StringMatching(`[a-zA-Z0-9_\-]{0,15}`).Draw(t, "excludedMode")
			}
		}

		result := panels.TestFilterExcluded(modes, excluded)

		// Build a lowercase set of excluded modes for checking.
		exSet := make(map[string]struct{}, len(excluded))
		for _, e := range excluded {
			e = strings.ToLower(strings.TrimSpace(e))
			if e != "" {
				exSet[e] = struct{}{}
			}
		}

		// Verify: no element in result is in the excluded set (case-insensitive).
		for _, r := range result {
			rLower := strings.ToLower(strings.TrimSpace(r))
			if _, found := exSet[rLower]; found {
				t.Fatalf("excluded mode %q (lowered: %q) appeared in output.\nmodes=%v\nexcluded=%v\nresult=%v",
					r, rLower, modes, excluded, result)
			}
		}
	})
}
