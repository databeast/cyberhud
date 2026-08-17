package tests

import (
	"testing"

	"github.com/databeast/cyberhud/display/modes/attract_bokeh"
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
		// Generate random values including out-of-range: negative speed,
		// density outside [0.1,1.0], size_variance and saturation outside [0.0,1.0].
		speed := rapid.Float64Range(-5.0, 20.0).Draw(rt, "speed")
		density := rapid.Float64Range(-1.0, 5.0).Draw(rt, "density")
		sizeVariance := rapid.Float64Range(-2.0, 3.0).Draw(rt, "sizeVariance")
		saturation := rapid.Float64Range(-2.0, 3.0).Draw(rt, "saturation")

		// Path 1: RestorePolicy with a map (simulates JSON restore from disk).
		data := map[string]interface{}{
			"speed":         speed,
			"density":       density,
			"size_variance": sizeVariance,
			"saturation":    saturation,
		}

		if err := attract_bokeh.RestorePolicyForTest(data); err != nil {
			t.Fatalf("RestorePolicy failed: %v", err)
		}
		restoredPolicy := attract_bokeh.GetPolicy()

		// Path 2: SetPolicy directly with same raw values.
		directPolicy := attract_bokeh.Policy{
			Speed:        speed,
			Density:      density,
			SizeVariance: sizeVariance,
			Saturation:   saturation,
		}
		attract_bokeh.SetPolicy(directPolicy)
		setPolicyResult := attract_bokeh.GetPolicy()

		// Reset to default after test.
		defer attract_bokeh.SetPolicy(attract_bokeh.DefaultPolicy())

		// Both paths must produce identical results.
		if restoredPolicy.Speed != setPolicyResult.Speed {
			t.Fatalf("Speed mismatch: RestorePolicy got %v, SetPolicy got %v (raw: %v)",
				restoredPolicy.Speed, setPolicyResult.Speed, speed)
		}
		if restoredPolicy.Density != setPolicyResult.Density {
			t.Fatalf("Density mismatch: RestorePolicy got %v, SetPolicy got %v (raw: %v)",
				restoredPolicy.Density, setPolicyResult.Density, density)
		}
		if restoredPolicy.SizeVariance != setPolicyResult.SizeVariance {
			t.Fatalf("SizeVariance mismatch: RestorePolicy got %v, SetPolicy got %v (raw: %v)",
				restoredPolicy.SizeVariance, setPolicyResult.SizeVariance, sizeVariance)
		}
		if restoredPolicy.Saturation != setPolicyResult.Saturation {
			t.Fatalf("Saturation mismatch: RestorePolicy got %v, SetPolicy got %v (raw: %v)",
				restoredPolicy.Saturation, setPolicyResult.Saturation, saturation)
		}
	})
}
