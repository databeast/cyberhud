package catalog_test

import (
	"testing"

	"github.com/databeast/cyberhud/display/catalog"
	_ "github.com/databeast/cyberhud/display/modes/clock"
	_ "github.com/databeast/cyberhud/display/modes/dashboard"
	_ "github.com/databeast/cyberhud/display/modes/gpio"
	_ "github.com/databeast/cyberhud/display/modes/gpio_control"
	_ "github.com/databeast/cyberhud/display/modes/image"
	_ "github.com/databeast/cyberhud/display/modes/menu"
	_ "github.com/databeast/cyberhud/display/modes/serial"
	_ "github.com/databeast/cyberhud/display/modes/stemma"
	_ "github.com/databeast/cyberhud/display/modes/system"
	_ "github.com/databeast/cyberhud/display/modes/systemd"
	_ "github.com/databeast/cyberhud/display/modes/ticker"
	_ "github.com/databeast/cyberhud/display/modes/usb"
	"pgregory.net/rapid"
)

// For any panelIndex > 0, SystemdBootPolicy.BootMode returns an empty string,
// ensuring only the primary panel (index 0) receives a boot mode override.

func TestProperty5_BootModeReturnsEmptyForNonPrimaryPanels(t *testing.T) {
	policy := catalog.SystemdBootPolicy{}

	rapid.Check(t, func(t *rapid.T) {
		// Generate any panel index > 0
		panelIndex := rapid.IntRange(1, 10000).Draw(t, "panelIndex")

		result := policy.BootMode(panelIndex)

		if result != "" {
			t.Fatalf("BootMode(%d) = %q, want empty string for non-primary panel",
				panelIndex, result)
		}
	})

	// Also verify the complementary property: panel 0 returns "systemd"
	if got := policy.BootMode(0); got != "systemd" {
		t.Fatalf("BootMode(0) = %q, want \"systemd\"", got)
	}
}

// Property 4: StandardPolicy always returns a valid mode
// *For any* non-empty list of available modes and any inputEnabled boolean,
// StandardPolicy.ResolveDefault returns a mode that is a member of the available list.
// Specifically: if inputEnabled and "menu" ∈ available, returns "menu"; if ¬inputEnabled
// and "dashboard" ∈ available, returns "dashboard"; otherwise returns the first available mode.

// genNonEmptyModeList generates a non-empty list of mode name strings.
func genNonEmptyModeList(t *rapid.T) []string {
	// Use a pool of realistic mode names plus some arbitrary strings.
	pool := []string{
		"menu", "dashboard", "clock", "stemma", "gpio", "system",
		"systemd", "ticker", "image", "usb", "serial", "thermal",
		"cycle", "testpattern", "testfonts", "gpio-control",
	}
	n := rapid.IntRange(1, 8).Draw(t, "listLen")
	modes := make([]string, n)
	for i := range modes {
		modes[i] = rapid.SampledFrom(pool).Draw(t, "mode")
	}
	return modes
}

// contains checks whether s is in the slice.
func contains(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}

// TestProperty_StandardPolicy_ResultInAvailable verifies that for any non-empty
// available list and any inputEnabled boolean, the result is always a member of available.
func TestProperty_StandardPolicy_ResultInAvailable(t *testing.T) {
	policy := catalog.StandardPolicy{}

	rapid.Check(t, func(rt *rapid.T) {
		available := genNonEmptyModeList(rt)
		inputEnabled := rapid.Bool().Draw(rt, "inputEnabled")

		result := policy.ResolveDefault(available, inputEnabled)

		if !contains(available, result) {
			t.Fatalf("ResolveDefault(%v, %v) = %q, which is not in available",
				available, inputEnabled, result)
		}
	})
}

// TestProperty_StandardPolicy_InputEnabledPrefersMenu verifies that when
// inputEnabled=true and "menu" is in the available list, the result is "menu".
func TestProperty_StandardPolicy_InputEnabledPrefersMenu(t *testing.T) {
	policy := catalog.StandardPolicy{}

	rapid.Check(t, func(rt *rapid.T) {
		// Generate a list that always contains "menu".
		others := rapid.SliceOfN(rapid.SampledFrom([]string{
			"clock", "stemma", "gpio", "system", "dashboard",
			"ticker", "image", "usb", "serial", "thermal",
		}), 0, 6).Draw(rt, "others")

		// Insert "menu" at a random position.
		pos := rapid.IntRange(0, len(others)).Draw(rt, "menuPos")
		available := make([]string, 0, len(others)+1)
		available = append(available, others[:pos]...)
		available = append(available, "menu")
		available = append(available, others[pos:]...)

		result := policy.ResolveDefault(available, true)

		if result != "menu" {
			t.Fatalf("ResolveDefault(%v, true) = %q, want \"menu\"",
				available, result)
		}
	})
}

// TestProperty_StandardPolicy_InputDisabledPrefersDashboard verifies that when
// inputEnabled=false and "dashboard" is in the available list, the result is "dashboard".
func TestProperty_StandardPolicy_InputDisabledPrefersDashboard(t *testing.T) {
	policy := catalog.StandardPolicy{}

	rapid.Check(t, func(rt *rapid.T) {
		// Generate a list that always contains "dashboard".
		others := rapid.SliceOfN(rapid.SampledFrom([]string{
			"clock", "stemma", "gpio", "system", "menu",
			"ticker", "image", "usb", "serial", "thermal",
		}), 0, 6).Draw(rt, "others")

		// Insert "dashboard" at a random position.
		pos := rapid.IntRange(0, len(others)).Draw(rt, "dashPos")
		available := make([]string, 0, len(others)+1)
		available = append(available, others[:pos]...)
		available = append(available, "dashboard")
		available = append(available, others[pos:]...)

		result := policy.ResolveDefault(available, false)

		if result != "dashboard" {
			t.Fatalf("ResolveDefault(%v, false) = %q, want \"dashboard\"",
				available, result)
		}
	})
}

// TestProperty_StandardPolicy_FallbackToFirst verifies that when the preferred
// mode is not in the available list, the result is always available[0].
func TestProperty_StandardPolicy_FallbackToFirst(t *testing.T) {
	policy := catalog.StandardPolicy{}

	rapid.Check(t, func(rt *rapid.T) {
		inputEnabled := rapid.Bool().Draw(rt, "inputEnabled")

		// Build a list that does NOT contain the preferred mode.
		// If inputEnabled, preferred is "menu"; else preferred is "dashboard".
		var excluded string
		if inputEnabled {
			excluded = "menu"
		} else {
			excluded = "dashboard"
		}

		// Pool without the excluded mode.
		pool := make([]string, 0, 10)
		for _, m := range []string{
			"clock", "stemma", "gpio", "system", "ticker",
			"image", "usb", "serial", "thermal", "cycle",
			"testpattern", "testfonts", "gpio-control",
		} {
			if m != excluded {
				pool = append(pool, m)
			}
		}
		// Also add the OTHER preferred mode to pool (it shouldn't be selected).
		if inputEnabled {
			pool = append(pool, "dashboard")
		} else {
			pool = append(pool, "menu")
		}

		n := rapid.IntRange(1, 6).Draw(rt, "listLen")
		available := make([]string, n)
		for i := range available {
			available[i] = rapid.SampledFrom(pool).Draw(rt, "mode")
		}

		result := policy.ResolveDefault(available, inputEnabled)

		if result != available[0] {
			t.Fatalf("ResolveDefault(%v, %v) = %q, want first element %q (preferred %q not in list)",
				available, inputEnabled, result, available[0], excluded)
		}
	})
}

func TestDescribeRegisteredModeIncludesOptions(t *testing.T) {
	def, ok := catalog.Describe("ticker")
	if !ok {
		t.Fatal("Describe(ticker) expected ok=true")
	}
	if def.Title != "Ticker" {
		t.Fatalf("Describe(ticker)=%+v", def)
	}
	if len(def.Options) == 0 {
		t.Fatalf("Describe(ticker).Options expected published options")
	}
	if def.Options[0].Key == "" {
		t.Fatalf("Describe(ticker).Options[0] missing key: %+v", def.Options[0])
	}
}

func TestIDs(t *testing.T) {
	// IDs returns all registered mode IDs in priority order.
	allModes := catalog.IDs()
	if len(allModes) == 0 {
		t.Fatal("IDs() expected non-empty")
	}

	// All modes should be available
	seen := map[string]bool{}
	for _, id := range allModes {
		seen[id] = true
	}
	for _, want := range []string{"systemd", "menu", "dashboard", "stemma", "gpio", "gpio-control", "system", "ticker", "image", "usb", "serial", "clock"} {
		if !seen[want] {
			t.Fatalf("IDs() missing %q: %v", want, allModes)
		}
	}
}

func TestModeCommandsRegistered(t *testing.T) {
	for _, verb := range []string{"ticker", "image", "usb", "serial"} {
		cmd, ok := catalog.Command(verb)
		if !ok {
			t.Fatalf("Command(%q) missing", verb)
		}
		if cmd.Verb != verb {
			t.Fatalf("Command(%q).Verb=%q", verb, cmd.Verb)
		}
		if cmd.Handle == nil {
			t.Fatalf("Command(%q).Handle is nil", verb)
		}
	}
}
