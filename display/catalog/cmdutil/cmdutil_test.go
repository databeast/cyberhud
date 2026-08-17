package cmdutil

import (
	"strings"
	"testing"
	"testing/quick"
)

// referenceParseBool is the original local implementation duplicated in
// clock, gpio, serial, and usb packages. It serves as the oracle for the
// property test.
func referenceParseBool(s string) (bool, bool) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "true", "yes", "on", "1":
		return true, true
	case "false", "no", "off", "0":
		return false, true
	default:
		return false, false
	}
}

// TestProperty_ParseBoolEquivalence verifies that for any arbitrary string input,
// cmdutil.ParseBool returns the same (bool, bool) pair as the old local
// implementations that were duplicated across mode packages.

func TestProperty_ParseBoolEquivalence(t *testing.T) {
	f := func(s string) bool {
		gotVal, gotOk := ParseBool(s)
		wantVal, wantOk := referenceParseBool(s)
		return gotVal == wantVal && gotOk == wantOk
	}

	cfg := &quick.Config{MaxCount: 100}
	if err := quick.Check(f, cfg); err != nil {
		t.Error(err)
	}
}
