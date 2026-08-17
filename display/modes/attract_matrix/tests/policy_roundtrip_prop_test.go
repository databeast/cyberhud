package tests

import (
	"encoding/json"
	"testing"

	"github.com/databeast/cyberhud/display/modes/attract_matrix"
	"pgregory.net/rapid"
)

// TestProperty_PolicySerializationRoundTrip verifies that for any valid matrix
// policy struct within normalized bounds, serializing to JSON via SnapshotPolicy()
// and deserializing back via RestorePolicy() produces an equivalent policy struct.
//
// Deserialize(Serialize(policy)) == policy
func TestProperty_PolicySerializationRoundTrip(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		// Generate a valid policy within normalized bounds.
		minSpeed := rapid.Float64Range(0.1, 50.0).Draw(rt, "minSpeed")
		maxSpeed := rapid.Float64Range(minSpeed, 100.0).Draw(rt, "maxSpeed")
		trailLength := rapid.IntRange(4, 128).Draw(rt, "trailLength")
		density := rapid.Float64Range(0.1, 1.0).Draw(rt, "density")
		showBackground := rapid.Bool().Draw(rt, "showBackground")

		original := attract_matrix.Policy{
			MinSpeed:       minSpeed,
			MaxSpeed:       maxSpeed,
			TrailLength:    trailLength,
			Density:        density,
			ShowBackground: showBackground,
		}

		// Apply normalization (ensures the policy is within valid bounds).
		original = attract_matrix.NormalizePolicy(original)

		// Set the policy so SnapshotPolicy() reads it.
		attract_matrix.SetPolicy(original)
		defer attract_matrix.SetPolicy(attract_matrix.DefaultPolicy())

		// Serialize via snapshotter.
		serialized := attract_matrix.SnapshotPolicyForTest()

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
		if err := attract_matrix.RestorePolicyForTest(decoded); err != nil {
			t.Fatalf("RestorePolicy failed: %v", err)
		}

		// Verify the restored policy matches the original.
		restored := attract_matrix.GetPolicy()

		if restored.MinSpeed != original.MinSpeed {
			t.Fatalf("MinSpeed mismatch: got %v, want %v", restored.MinSpeed, original.MinSpeed)
		}
		if restored.MaxSpeed != original.MaxSpeed {
			t.Fatalf("MaxSpeed mismatch: got %v, want %v", restored.MaxSpeed, original.MaxSpeed)
		}
		if restored.TrailLength != original.TrailLength {
			t.Fatalf("TrailLength mismatch: got %v, want %v", restored.TrailLength, original.TrailLength)
		}
		if restored.Density != original.Density {
			t.Fatalf("Density mismatch: got %v, want %v", restored.Density, original.Density)
		}
		if restored.ShowBackground != original.ShowBackground {
			t.Fatalf("ShowBackground mismatch: got %v, want %v", restored.ShowBackground, original.ShowBackground)
		}
	})
}
