package tests

import (
	"encoding/json"
	"testing"

	"github.com/databeast/cyberhud/display/modes/ticker"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"pgregory.net/rapid"
)

// TestProperty_PolicySerializationRoundTrip verifies that for any valid ticker
// policy struct within normalized bounds, serializing to JSON via SnapshotPolicy()
// and deserializing back via RestorePolicy() produces an equivalent policy struct.
//
// Deserialize(Serialize(policy)) == policy
func TestProperty_PolicySerializationRoundTrip(t *testing.T) {
	// Collect valid style names from the registry.
	styles := ticker.TickerRegistryEnumerate()
	if len(styles) == 0 {
		t.Fatal("tickerRegistry contains zero styles")
	}
	styleNames := make([]string, len(styles))
	for i, s := range styles {
		styleNames[i] = s.Name()
	}

	rapid.Check(t, func(rt *rapid.T) {
		// Generate valid policy values from allowed sets.
		style := rapid.SampledFrom(styleNames).Draw(rt, "style")
		font := rapid.SampledFrom([]string{"auto", "spleen-8x16", "spleen-16x32"}).Draw(rt, "font")
		fontTier := rapid.SampledFrom(ticker.AllowedFontTiers()).Draw(rt, "fontTier")
		lineMode := rapid.SampledFrom([]string{textlayout.LineModeTruncate, textlayout.LineModeClip}).Draw(rt, "lineMode")
		direction := rapid.SampledFrom([]string{textlayout.TickerDirectionVertical, "horizontal", textlayout.TickerDirectionNone}).Draw(rt, "direction")
		autoScrollMS := rapid.IntRange(0, 5000).Draw(rt, "autoScrollMS")
		accent := rapid.SampledFrom([]string{"cyan", "green", "amber", "red", "white", "none"}).Draw(rt, "accent")

		original := ticker.Policy{
			Style:        style,
			Font:         font,
			FontTier:     fontTier,
			LineMode:     lineMode,
			Direction:    direction,
			AutoScrollMS: autoScrollMS,
			Accent:       accent,
		}

		// Apply normalization (ensures the policy is within valid bounds).
		original = ticker.NormalizePolicy(original)

		// Set the policy so PolicySnapshot() reads it.
		ticker.SetPolicy(original)
		defer ticker.SetPolicy(ticker.DefaultPolicy())

		// Serialize via snapshotter.
		snap := ticker.NewSnapshotterForTest()
		serialized := snap.SnapshotPolicy()

		// Encode to JSON (simulates wire/disk format).
		jsonBytes, err := json.Marshal(serialized)
		if err != nil {
			t.Fatalf("json.Marshal failed: %v", err)
		}

		// Decode from JSON back to map.
		var decoded map[string]interface{}
		if err := json.Unmarshal(jsonBytes, &decoded); err != nil {
			t.Fatalf("json.Unmarshal failed: %v", err)
		}

		// Restore via snapshotter.
		if err := snap.RestorePolicy(decoded); err != nil {
			t.Fatalf("RestorePolicy failed: %v", err)
		}

		// Verify the restored policy matches the original.
		restored := ticker.PolicySnapshot()

		if restored.Style != original.Style {
			t.Fatalf("Style mismatch: got %q, want %q", restored.Style, original.Style)
		}
		if restored.Font != original.Font {
			t.Fatalf("Font mismatch: got %q, want %q", restored.Font, original.Font)
		}
		if restored.FontTier != original.FontTier {
			t.Fatalf("FontTier mismatch: got %q, want %q", restored.FontTier, original.FontTier)
		}
		if restored.LineMode != original.LineMode {
			t.Fatalf("LineMode mismatch: got %q, want %q", restored.LineMode, original.LineMode)
		}
		if restored.Direction != original.Direction {
			t.Fatalf("Direction mismatch: got %q, want %q", restored.Direction, original.Direction)
		}
		if restored.AutoScrollMS != original.AutoScrollMS {
			t.Fatalf("AutoScrollMS mismatch: got %v, want %v", restored.AutoScrollMS, original.AutoScrollMS)
		}
		if restored.Accent != original.Accent {
			t.Fatalf("Accent mismatch: got %q, want %q", restored.Accent, original.Accent)
		}
	})
}
