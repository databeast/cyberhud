package tests

import (
	"testing"

	"github.com/databeast/cyberhud/display/modes/attract_matrix"
	"pgregory.net/rapid"
)

// TestProperty_RestoreNormalizationEquivalence verifies that for any policy values
// (including out-of-range values), restoring a policy from a JSON map via
// RestorePolicy() produces the same result as calling SetPolicy() directly with
// the same raw values. Both paths must apply identical normalization.
//
// RestorePolicy(values) == normalizePolicy(fromRawValues(values))
func TestProperty_RestoreNormalizationEquivalence(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		// Generate random values including out-of-range: negative speeds,
		// trail lengths outside [4,128], density outside [0.1,1.0].
		minSpeed := rapid.Float64Range(-10.0, 100.0).Draw(rt, "minSpeed")
		maxSpeed := rapid.Float64Range(-10.0, 200.0).Draw(rt, "maxSpeed")
		trailLength := rapid.IntRange(-50, 256).Draw(rt, "trailLength")
		density := rapid.Float64Range(-1.0, 5.0).Draw(rt, "density")
		showBackground := rapid.Bool().Draw(rt, "showBackground")

		// Path 1: RestorePolicy with a map (simulates JSON restore from disk).
		data := map[string]interface{}{
			"min_speed":       minSpeed,
			"max_speed":       maxSpeed,
			"trail_length":    trailLength,
			"density":         density,
			"show_background": showBackground,
		}

		if err := attract_matrix.RestorePolicyForTest(data); err != nil {
			t.Fatalf("RestorePolicy failed: %v", err)
		}
		restoredPolicy := attract_matrix.GetPolicy()

		// Path 2: SetPolicy directly with same raw values.
		directPolicy := attract_matrix.Policy{
			MinSpeed:       minSpeed,
			MaxSpeed:       maxSpeed,
			TrailLength:    trailLength,
			Density:        density,
			ShowBackground: showBackground,
		}
		attract_matrix.SetPolicy(directPolicy)
		setPolicyResult := attract_matrix.GetPolicy()

		// Reset to default after test.
		defer attract_matrix.SetPolicy(attract_matrix.DefaultPolicy())

		// Both paths must produce identical results.
		if restoredPolicy.MinSpeed != setPolicyResult.MinSpeed {
			t.Fatalf("MinSpeed mismatch: RestorePolicy got %v, SetPolicy got %v (raw: %v)",
				restoredPolicy.MinSpeed, setPolicyResult.MinSpeed, minSpeed)
		}
		if restoredPolicy.MaxSpeed != setPolicyResult.MaxSpeed {
			t.Fatalf("MaxSpeed mismatch: RestorePolicy got %v, SetPolicy got %v (raw: %v)",
				restoredPolicy.MaxSpeed, setPolicyResult.MaxSpeed, maxSpeed)
		}
		if restoredPolicy.TrailLength != setPolicyResult.TrailLength {
			t.Fatalf("TrailLength mismatch: RestorePolicy got %v, SetPolicy got %v (raw: %v)",
				restoredPolicy.TrailLength, setPolicyResult.TrailLength, trailLength)
		}
		if restoredPolicy.Density != setPolicyResult.Density {
			t.Fatalf("Density mismatch: RestorePolicy got %v, SetPolicy got %v (raw: %v)",
				restoredPolicy.Density, setPolicyResult.Density, density)
		}
		if restoredPolicy.ShowBackground != setPolicyResult.ShowBackground {
			t.Fatalf("ShowBackground mismatch: RestorePolicy got %v, SetPolicy got %v (raw: %v)",
				restoredPolicy.ShowBackground, setPolicyResult.ShowBackground, showBackground)
		}
	})
}
