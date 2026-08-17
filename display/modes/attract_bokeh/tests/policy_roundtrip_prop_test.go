package tests

import (
	"encoding/json"
	"testing"

	"github.com/databeast/cyberhud/display/modes/attract_bokeh"
	"pgregory.net/rapid"
)

// TestProperty_PolicySerializationRoundTrip verifies that for any valid bokeh
// policy struct within normalized bounds, serializing to JSON via SnapshotPolicy()
// and deserializing back via RestorePolicy() produces an equivalent policy struct.
//
// Deserialize(Serialize(policy)) == policy
func TestProperty_PolicySerializationRoundTrip(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		// Generate a valid policy within normalized bounds.
		speed := rapid.Float64Range(0.1, 10.0).Draw(rt, "speed")
		density := rapid.Float64Range(0.1, 1.0).Draw(rt, "density")
		sizeVariance := rapid.Float64Range(0.0, 1.0).Draw(rt, "sizeVariance")
		saturation := rapid.Float64Range(0.0, 1.0).Draw(rt, "saturation")

		original := attract_bokeh.Policy{
			Speed:        speed,
			Density:      density,
			SizeVariance: sizeVariance,
			Saturation:   saturation,
		}

		// Apply normalization (ensures the policy is within valid bounds).
		original = attract_bokeh.NormalizePolicy(original)

		// Set the policy so SnapshotPolicy() reads it.
		attract_bokeh.SetPolicy(original)
		defer attract_bokeh.SetPolicy(attract_bokeh.DefaultPolicy())

		// Serialize via snapshotter.
		serialized := attract_bokeh.SnapshotPolicyForTest()

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
		if err := attract_bokeh.RestorePolicyForTest(decoded); err != nil {
			t.Fatalf("RestorePolicy failed: %v", err)
		}

		// Verify the restored policy matches the original.
		restored := attract_bokeh.GetPolicy()

		if restored.Speed != original.Speed {
			t.Fatalf("Speed mismatch: got %v, want %v", restored.Speed, original.Speed)
		}
		if restored.Density != original.Density {
			t.Fatalf("Density mismatch: got %v, want %v", restored.Density, original.Density)
		}
		if restored.SizeVariance != original.SizeVariance {
			t.Fatalf("SizeVariance mismatch: got %v, want %v", restored.SizeVariance, original.SizeVariance)
		}
		if restored.Saturation != original.Saturation {
			t.Fatalf("Saturation mismatch: got %v, want %v", restored.Saturation, original.Saturation)
		}
	})
}
