package displaymode

import (
	"strings"
	"testing"

	"pgregory.net/rapid"
)

// =============================================================================
//
// For any valid style name present in registeredStyleNames(), applying it via
// the CmdHandler's Apply callback and then reading it via the Get callback
// shall return the same value (after normalization).

// =============================================================================

func TestPropertyCmdHandlerGetApplyRoundTrip(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		names := registeredStyleNames()
		idx := rapid.IntRange(0, len(names)-1).Draw(t, "nameIndex")
		name := names[idx]

		// Apply via HandleCommand
		result := HandleCommand([]string{"style=" + name})
		if strings.HasPrefix(result, "ERR") {
			t.Fatalf("valid style %q rejected: %s", name, result)
		}

		// Get via GetPolicy
		got := GetPolicy().Style
		if got != name {
			t.Fatalf("round-trip failed: applied %q, got %q", name, got)
		}
	})
}

// =============================================================================
//
// For any string that is not a recognized key name or whose value fails
// validation, HandleCommand shall return a response prefixed with "ERR " and
// the Policy state shall be unchanged from before the call.

// =============================================================================

func TestPropertyInvalidCommandRejection(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Save current policy
		before := GetPolicy()

		// Generate invalid commands
		choice := rapid.IntRange(0, 2).Draw(t, "choice")
		var args []string
		switch choice {
		case 0:
			// Unrecognized key
			key := rapid.StringMatching(`[a-z]{3,10}`).Draw(t, "badKey")
			if key == "style" {
				key = "badkey"
			}
			args = []string{key + "=somevalue"}
		case 1:
			// Valid key, invalid value
			value := rapid.StringMatching(`[A-Z]{5,15}`).Draw(t, "badValue")
			// Make sure value is not a valid style name
			names := registeredStyleNames()
			isValid := false
			for _, n := range names {
				if strings.EqualFold(n, value) {
					isValid = true
					break
				}
			}
			if isValid {
				value = "DEFINITELY_INVALID_STYLE_XYZ123"
			}
			args = []string{"style=" + value}
		case 2:
			// Bare token (no = sign)
			token := rapid.StringMatching(`[a-z]{3,8}`).Draw(t, "bareToken")
			args = []string{token}
		}

		result := HandleCommand(args)

		if !strings.HasPrefix(result, "ERR") {
			t.Fatalf("expected ERR prefix for args %v, got: %s", args, result)
		}

		// Policy should be unchanged
		after := GetPolicy()
		if after.Style != before.Style {
			t.Fatalf("policy mutated on invalid command: before=%q, after=%q", before.Style, after.Style)
		}
	})
}

// =============================================================================
//
// For any Policy value, after SetPolicy(p) followed by GetPolicy(), the
// returned Policy's Style field shall be lowercase, whitespace-trimmed, and
// either empty (meaning best-fit) or present in the registry's style name list.
// If the input Style was a valid registry name (after trimming and lowercasing),
// the output Style shall equal that name. Otherwise the output Style shall be
// empty (best-fit).

// =============================================================================

func TestPropertyPolicyNormalizationCorrectness(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		names := registeredStyleNames()

		// Mix of valid names (with case/whitespace variations) and invalid strings
		choice := rapid.IntRange(0, 2).Draw(t, "choice")
		var inputStyle string
		switch choice {
		case 0:
			// Valid name with random case and whitespace
			idx := rapid.IntRange(0, len(names)-1).Draw(t, "nameIdx")
			base := names[idx]
			inputStyle = rapid.SampledFrom([]string{
				base,
				strings.ToUpper(base),
				"  " + base + "  ",
				strings.ToUpper(base[:1]) + base[1:],
			}).Draw(t, "caseVariant")
		case 1:
			// Completely random string (likely invalid)
			inputStyle = rapid.String().Draw(t, "randomStyle")
		case 2:
			// Empty or whitespace only
			inputStyle = rapid.SampledFrom([]string{"", "   ", "\t"}).Draw(t, "emptyStyle")
		}

		SetPolicy(Policy{Style: inputStyle})
		got := GetPolicy()

		// Style should be lowercase and trimmed
		if got.Style != strings.ToLower(strings.TrimSpace(got.Style)) {
			t.Fatalf("policy Style not normalized: %q", got.Style)
		}

		// Style should be either empty (best-fit) or a registered name
		if got.Style != "" {
			found := false
			for _, n := range names {
				if n == got.Style {
					found = true
					break
				}
			}
			if !found {
				t.Fatalf("policy Style %q not in registry and not empty", got.Style)
			}
		}

		// If input was valid (after trim+lower), output should match
		normalized := strings.ToLower(strings.TrimSpace(inputStyle))
		isValidInput := false
		for _, n := range names {
			if n == normalized {
				isValidInput = true
				break
			}
		}
		if isValidInput {
			if got.Style != normalized {
				t.Fatalf("valid input %q normalized to %q instead of %q", inputStyle, got.Style, normalized)
			}
		} else {
			// Invalid input should fall back to default
			if got.Style != DefaultPolicy().Style {
				t.Fatalf("invalid input %q normalized to %q instead of default %q", inputStyle, got.Style, DefaultPolicy().Style)
			}
		}
	})
}
