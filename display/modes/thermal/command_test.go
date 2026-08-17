package thermal

import (
	"strings"
	"testing"
)

func TestHandleCommand_ThresholdRejection(t *testing.T) {
	// Ensure default policy is in a valid state before test.
	if err := SetPolicy(DefaultPolicy()); err != nil {
		t.Fatalf("SetPolicy(default) failed: %v", err)
	}

	// warn_threshold=90 crit_threshold=70 should return error because 90 >= 70.
	result := HandleCommand([]string{"warn_threshold=90", "crit_threshold=70"})
	if !strings.HasPrefix(result, "ERR") {
		t.Fatalf("expected ERR response for warn>=crit, got: %s", result)
	}
	if !strings.Contains(result, "warn_threshold") {
		t.Fatalf("expected error to mention warn_threshold, got: %s", result)
	}

	// Verify policy was rolled back to previous valid state.
	p := GetPolicy()
	if p.WarnThreshold == 90 || p.CritThreshold == 70 {
		t.Fatalf("policy should have been rolled back, got warn=%d crit=%d", p.WarnThreshold, p.CritThreshold)
	}
}

func TestHandleCommand_ValidThresholds(t *testing.T) {
	// Reset to default.
	if err := SetPolicy(DefaultPolicy()); err != nil {
		t.Fatalf("SetPolicy(default) failed: %v", err)
	}

	// warn_threshold=60 crit_threshold=80 is valid (60 < 80).
	result := HandleCommand([]string{"warn_threshold=60", "crit_threshold=80"})
	if strings.HasPrefix(result, "ERR") {
		t.Fatalf("expected success for valid thresholds, got: %s", result)
	}

	p := GetPolicy()
	if p.WarnThreshold != 60 {
		t.Fatalf("expected WarnThreshold=60, got %d", p.WarnThreshold)
	}
	if p.CritThreshold != 80 {
		t.Fatalf("expected CritThreshold=80, got %d", p.CritThreshold)
	}
}
