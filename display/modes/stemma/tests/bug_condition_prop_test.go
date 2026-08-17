package tests_test

import (
	"testing"
	"time"

	"github.com/databeast/cyberhud/display/modes/stemma"
	"github.com/databeast/cyberhud/display/modes/stemma/source"
	"github.com/databeast/cyberhud/display/modes/stemma/tests"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"pgregory.net/rapid"
)

// Bug Condition Exploration Property Test
//
// Property 1: Bug Condition - Scanner Runs Without Stemma Mode Active
//
// This test encodes the EXPECTED (correct) behavior:
//   - The scanner should NOT be running/accessible when the stemma mode is inactive
//   - Deactivate() should stop the scanner (GlobalScanner returns nil after deactivation)
//   - Activate() should be the ONLY path that starts the scanner
//   - Before any Activate() call, GlobalScanner should return nil
//
// On UNFIXED code this test will FAIL because:
//   - Activate() and Deactivate() are no-ops
//   - main.go unconditionally sets the global scanner at boot
//   - The global scanner persists regardless of mode lifecycle calls

// TestProperty1_DeactivateClearsGlobalScanner verifies that after calling
// instance.Deactivate(), the global scanner is nil (not accessible).
// On unfixed code, Deactivate() is a no-op so the scanner remains set.
func TestProperty1_DeactivateClearsGlobalScanner(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate random I2C bus configuration to simulate various daemon states
		busCount := rapid.IntRange(1, 3).Draw(t, "busCount")
		buses := make([]string, busCount)
		for i := range buses {
			buses[i] = rapid.SampledFrom([]string{
				"/dev/i2c-1", "/dev/i2c-3", "/dev/i2c-99",
			}).Draw(t, "bus")
		}
		interval := time.Duration(rapid.IntRange(1, 10).Draw(t, "intervalSec")) * time.Second

		// Simulate what main.go does: create scanner, start it, set global
		scanner := source.New(buses, interval)
		scanner.Start()
		source.SetGlobalScanner(scanner)

		// Create a stemma mode instance (mimics what the mode system does)
		hints := textlayout.TextHints{
			PixelWidth:   128,
			PixelHeight:  64,
			GlyphAdvance: 6,
			RowHeight:    10,
		}
		inst := tests.NewInstanceForTest(hints)

		// Activate then Deactivate — the expected behavior is that
		// Deactivate clears the global scanner
		inst.Activate()
		inst.Deactivate()

		// EXPECTED: After Deactivate(), the global scanner should be nil
		// BUG: On unfixed code, Deactivate() is a no-op so GlobalScanner() != nil
		if got := source.GlobalScanner(); got != nil {
			t.Fatalf("after Deactivate(), GlobalScanner() = %v, want nil "+
				"(scanner still running — Deactivate is a no-op)", got)
		}

		// Cleanup: stop the scanner we started
		scanner.Stop()
		source.SetGlobalScanner(nil)
	})
}

// TestProperty1_ActivateIsOnlyPathToStartScanner verifies that before
// Activate() is called, no scanner is available via GlobalScanner().
// On unfixed code, the scanner is set globally by main.go before any
// mode lifecycle call, so this fails.
func TestProperty1_ActivateIsOnlyPathToStartScanner(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Ensure clean state: no global scanner set
		source.SetGlobalScanner(nil)

		// Create instance — on correct code, the instance should NOT
		// have a running scanner until Activate() is called
		hints := textlayout.TextHints{
			PixelWidth:   rapid.IntRange(64, 320).Draw(t, "pixelWidth"),
			PixelHeight:  rapid.IntRange(32, 240).Draw(t, "pixelHeight"),
			GlyphAdvance: rapid.IntRange(4, 12).Draw(t, "glyphAdvance"),
			RowHeight:    rapid.IntRange(8, 16).Draw(t, "rowHeight"),
		}
		inst := tests.NewInstanceForTest(hints)

		// Before Activate(), GlobalScanner() should be nil
		// (no scanner running when stemma mode hasn't been activated)
		if got := source.GlobalScanner(); got != nil {
			t.Fatalf("before Activate(), GlobalScanner() = %v, want nil "+
				"(scanner should not exist until mode is activated)", got)
		}

		// Now simulate what main.go incorrectly does: set global scanner externally
		// This simulates the bug condition where the scanner exists before Activate
		buses := []string{"/dev/i2c-1"}
		externalScanner := source.New(buses, 2*time.Second)
		externalScanner.Start()
		source.SetGlobalScanner(externalScanner)

		// Even with an external scanner set, calling Activate() on the instance
		// should be what controls scanner lifecycle.
		// After deactivation, the scanner should be gone.
		inst.Activate()
		inst.Deactivate()

		// EXPECTED: scanner is nil after deactivate
		if got := source.GlobalScanner(); got != nil {
			t.Fatalf("after Activate()+Deactivate() with external scanner, "+
				"GlobalScanner() = %v, want nil (Deactivate should clear it)", got)
		}

		// Cleanup
		externalScanner.Stop()
		source.SetGlobalScanner(nil)
	})
}

// TestProperty1_ScannerNotRunningWhenModeInactive verifies the core bug condition:
// when the stemma mode is not active (no Activate() called), the scanner should
// NOT be accessible via GlobalScanner().
// On fixed code, main.go no longer sets the global scanner unconditionally — the
// scanner lifecycle is fully owned by Activate()/Deactivate().
// This test verifies:
//   - Without Activate(), GlobalScanner() returns nil (no scanner running)
//   - After Activate(), GlobalScanner() returns non-nil (scanner is running)
//   - After Deactivate(), GlobalScanner() returns nil again (scanner stopped)
func TestProperty1_ScannerNotRunningWhenModeInactive(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate a random non-stemma active mode name to contextualize the scenario
		activeMode := rapid.SampledFrom([]string{
			"clock", "dashboard", "serial", "thermal", "system",
			"gpio", "ticker", "usb", "wifi", "menu", "image",
		}).Draw(t, "activeMode")
		_ = activeMode // used for documentation/context in failure messages

		// Ensure clean state: no global scanner set (simulates fixed main.go)
		source.SetGlobalScanner(nil)

		// Configure scanner params so Activate() can create one
		buses := []string{"/dev/i2c-1"}
		interval := time.Duration(rapid.IntRange(1, 5).Draw(t, "intervalSec")) * time.Second
		stemma.SetScannerConfig(buses, interval)

		// Create stemma instance but do NOT activate it (simulates another mode being active)
		hints := textlayout.TextHints{
			PixelWidth:   128,
			PixelHeight:  64,
			GlyphAdvance: 6,
			RowHeight:    10,
		}
		inst := tests.NewInstanceForTest(hints)

		// EXPECTED: Without Activate(), GlobalScanner() should be nil
		if got := source.GlobalScanner(); got != nil {
			t.Fatalf("stemma mode not activated (active mode would be %q), "+
				"but GlobalScanner() != nil (scanner running without activation)", activeMode)
		}

		// Now activate stemma — scanner should become available
		inst.Activate()
		if got := source.GlobalScanner(); got == nil {
			t.Fatalf("after Activate(), GlobalScanner() is nil — scanner was not started")
		}

		// Deactivate stemma — scanner should be cleared
		inst.Deactivate()
		if got := source.GlobalScanner(); got != nil {
			t.Fatalf("after Deactivate(), GlobalScanner() != nil — scanner still running")
		}
	})
}
