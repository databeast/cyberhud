package main

import (
	"encoding/json"
	"maps"
	"testing"

	"pgregory.net/rapid"
)

// orientationFieldEqual reports whether two *orientationField values are
// equivalent, treating nil and comparing both the Single and PerScreen forms.
func orientationFieldEqual(a, b *orientationField) bool {
	if a == nil || b == nil {
		return a == b
	}
	return a.Single == b.Single && maps.Equal(a.PerScreen, b.PerScreen)
}

// For any existing config file (with hardware settings and possibly old policies) and
// any new policy snapshot map, merging the new policies into the config should:
// 1. Preserve all non-policy fields (socket, i2c, scan, display) unchanged
// 2. For mode IDs present only in the old config, keep old values
// 3. For mode IDs present in the new snapshot, use new values (overwrite)

// genModeID generates a plausible mode ID string (lowercase with underscores).
func genModeID(rt *rapid.T, label string) string {
	const chars = "abcdefghijklmnopqrstuvwxyz_"
	length := rapid.IntRange(3, 20).Draw(rt, label+"_len")
	b := make([]byte, length)
	for i := range b {
		b[i] = chars[rapid.IntRange(0, len(chars)-1).Draw(rt, label+"_char")]
	}
	return string(b)
}

// genRawJSON generates a random JSON object as json.RawMessage.
func genRawJSON(rt *rapid.T, label string) json.RawMessage {
	m := map[string]interface{}{
		"key_" + rapid.StringMatching(`[a-z]{3,8}`).Draw(rt, label+"_key"): rapid.Float64Range(0, 100).Draw(rt, label+"_val"),
	}
	data, _ := json.Marshal(m)
	return json.RawMessage(data)
}

// genPolicyMap generates a random map[string]json.RawMessage with 0-5 entries.
func genPolicyMap(rt *rapid.T, label string) map[string]json.RawMessage {
	count := rapid.IntRange(0, 5).Draw(rt, label+"_count")
	if count == 0 {
		return nil
	}
	m := make(map[string]json.RawMessage, count)
	for i := 0; i < count; i++ {
		id := genModeID(rt, label+"_id")
		m[id] = genRawJSON(rt, label+"_val")
	}
	return m
}

// genFileDisplayConfig generates a random fileDisplayConfig.
func genFileDisplayConfig(rt *rapid.T) fileDisplayConfig {
	disabled := rapid.Bool().Draw(rt, "disabled")
	disableInput := rapid.Bool().Draw(rt, "disableInput")
	width := rapid.IntRange(64, 1920).Draw(rt, "width")
	height := rapid.IntRange(64, 1080).Draw(rt, "height")
	rotate := rapid.Bool().Draw(rt, "rotate")
	xOffset := rapid.IntRange(0, 100).Draw(rt, "xOffset")
	yOffset := rapid.IntRange(0, 100).Draw(rt, "yOffset")

	return fileDisplayConfig{
		Disabled:     &disabled,
		Profile:      rapid.StringMatching(`[a-z0-9-]{3,20}`).Draw(rt, "profile"),
		DisableInput: &disableInput,
		Orientation:  &orientationField{PerScreen: map[string]string{"main": rapid.SampledFrom([]string{"normal", "flip", "cw", "ccw"}).Draw(rt, "orientation")}},
		Width:        &width,
		Height:       &height,
		MADCTL:       rapid.StringMatching(`[A-Z0-9]{0,4}`).Draw(rt, "madctl"),
		Rotate:       &rotate,
		XOffset:      &xOffset,
		YOffset:      &yOffset,
		DC:           rapid.StringMatching(`[A-Z0-9]{0,8}`).Draw(rt, "dc"),
		RST:          rapid.StringMatching(`[A-Z0-9]{0,8}`).Draw(rt, "rst"),
		BL:           rapid.StringMatching(`[A-Z0-9]{0,8}`).Draw(rt, "bl"),
		Busy:         rapid.StringMatching(`[A-Z0-9]{0,8}`).Draw(rt, "busy"),
		InputKey1:    rapid.StringMatching(`[a-z0-9]{0,8}`).Draw(rt, "key1"),
		InputKey2:    rapid.StringMatching(`[a-z0-9]{0,8}`).Draw(rt, "key2"),
		InputKey3:    rapid.StringMatching(`[a-z0-9]{0,8}`).Draw(rt, "key3"),
		InputUp:      rapid.StringMatching(`[a-z0-9]{0,8}`).Draw(rt, "up"),
		InputDown:    rapid.StringMatching(`[a-z0-9]{0,8}`).Draw(rt, "down"),
		InputLeft:    rapid.StringMatching(`[a-z0-9]{0,8}`).Draw(rt, "left"),
		InputRight:   rapid.StringMatching(`[a-z0-9]{0,8}`).Draw(rt, "right"),
		InputPress:   rapid.StringMatching(`[a-z0-9]{0,8}`).Draw(rt, "press"),
	}
}

// genFileConfig generates a random fileConfig with optional policies.
func genFileConfig(rt *rapid.T) fileConfig {
	return fileConfig{
		Socket:   rapid.StringMatching(`/[a-z/]{3,30}\.sock`).Draw(rt, "socket"),
		I2C:      rapid.StringMatching(`/dev/i2c-[0-9]`).Draw(rt, "i2c"),
		Scan:     rapid.SampledFrom([]string{"1s", "5s", "10s", "30s", "1m"}).Draw(rt, "scan"),
		Display:  genFileDisplayConfig(rt),
		Policies: genPolicyMap(rt, "oldPolicies"),
	}
}

// mergePolicies simulates the merge logic: new policies overwrite existing entries,
// old-only entries are preserved.
func mergePolicies(cfg *fileConfig, newPolicies map[string]json.RawMessage) {
	if cfg.Policies == nil {
		cfg.Policies = make(map[string]json.RawMessage)
	}
	for modeID, data := range newPolicies {
		cfg.Policies[modeID] = data
	}
}

func TestProperty_PolicyConfigMergePreservation(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		// Generate a random existing config
		cfg := genFileConfig(rt)

		// Snapshot non-policy fields before merge
		origSocket := cfg.Socket
		origI2C := cfg.I2C
		origScan := cfg.Scan
		origDisplay := cfg.Display

		// Deep copy the old policies for later comparison
		oldPolicies := make(map[string]json.RawMessage, len(cfg.Policies))
		for k, v := range cfg.Policies {
			cp := make(json.RawMessage, len(v))
			copy(cp, v)
			oldPolicies[k] = cp
		}

		// Generate new policies to merge
		newPolicies := genPolicyMap(rt, "newPolicies")

		// Perform the merge
		mergePolicies(&cfg, newPolicies)

		// 1. Non-policy fields must be unchanged
		if cfg.Socket != origSocket {
			t.Fatalf("Socket changed: got %q, want %q", cfg.Socket, origSocket)
		}
		if cfg.I2C != origI2C {
			t.Fatalf("I2C changed: got %q, want %q", cfg.I2C, origI2C)
		}
		if cfg.Scan != origScan {
			t.Fatalf("Scan changed: got %q, want %q", cfg.Scan, origScan)
		}
		if cfg.Display.Profile != origDisplay.Profile {
			t.Fatalf("Display.Profile changed: got %q, want %q", cfg.Display.Profile, origDisplay.Profile)
		}
		if !orientationFieldEqual(cfg.Display.Orientation, origDisplay.Orientation) {
			t.Fatalf("Display.Orientation changed: got %v, want %v", cfg.Display.Orientation, origDisplay.Orientation)
		}
		if cfg.Display.MADCTL != origDisplay.MADCTL {
			t.Fatalf("Display.MADCTL changed: got %q, want %q", cfg.Display.MADCTL, origDisplay.MADCTL)
		}
		if cfg.Display.DC != origDisplay.DC {
			t.Fatalf("Display.DC changed: got %q, want %q", cfg.Display.DC, origDisplay.DC)
		}
		if cfg.Display.RST != origDisplay.RST {
			t.Fatalf("Display.RST changed: got %q, want %q", cfg.Display.RST, origDisplay.RST)
		}
		if cfg.Display.BL != origDisplay.BL {
			t.Fatalf("Display.BL changed: got %q, want %q", cfg.Display.BL, origDisplay.BL)
		}
		if cfg.Display.Width != nil && origDisplay.Width != nil && *cfg.Display.Width != *origDisplay.Width {
			t.Fatalf("Display.Width changed: got %d, want %d", *cfg.Display.Width, *origDisplay.Width)
		}
		if cfg.Display.Height != nil && origDisplay.Height != nil && *cfg.Display.Height != *origDisplay.Height {
			t.Fatalf("Display.Height changed: got %d, want %d", *cfg.Display.Height, *origDisplay.Height)
		}

		// 2. Old-only policies (keys present in old but not in new) must be preserved
		for modeID, oldData := range oldPolicies {
			if _, inNew := newPolicies[modeID]; inNew {
				continue // overwritten by new — checked in step 3
			}
			got, exists := cfg.Policies[modeID]
			if !exists {
				t.Fatalf("old-only policy %q was removed after merge", modeID)
			}
			if string(got) != string(oldData) {
				t.Fatalf("old-only policy %q was modified: got %s, want %s", modeID, got, oldData)
			}
		}

		// 3. New policies must overwrite (use new values)
		for modeID, newData := range newPolicies {
			got, exists := cfg.Policies[modeID]
			if !exists {
				t.Fatalf("new policy %q missing after merge", modeID)
			}
			if string(got) != string(newData) {
				t.Fatalf("new policy %q not applied: got %s, want %s", modeID, got, newData)
			}
		}
	})
}

func TestProperty_PolicyConfigMergeNilMaps(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		// Generate config with nil policies
		cfg := fileConfig{
			Socket:   rapid.StringMatching(`/[a-z/]{3,30}\.sock`).Draw(rt, "socket"),
			I2C:      rapid.StringMatching(`/dev/i2c-[0-9]`).Draw(rt, "i2c"),
			Scan:     rapid.SampledFrom([]string{"1s", "5s", "10s"}).Draw(rt, "scan"),
			Display:  genFileDisplayConfig(rt),
			Policies: nil, // explicitly nil
		}

		origSocket := cfg.Socket
		origI2C := cfg.I2C
		origScan := cfg.Scan

		// Generate new policies (may also be nil/empty)
		newPolicies := genPolicyMap(rt, "newPolicies")

		// Perform the merge
		mergePolicies(&cfg, newPolicies)

		// Non-policy fields preserved
		if cfg.Socket != origSocket {
			t.Fatalf("Socket changed: got %q, want %q", cfg.Socket, origSocket)
		}
		if cfg.I2C != origI2C {
			t.Fatalf("I2C changed: got %q, want %q", cfg.I2C, origI2C)
		}
		if cfg.Scan != origScan {
			t.Fatalf("Scan changed: got %q, want %q", cfg.Scan, origScan)
		}

		// All new policies present
		for modeID, newData := range newPolicies {
			got, exists := cfg.Policies[modeID]
			if !exists {
				t.Fatalf("new policy %q missing after merge from nil map", modeID)
			}
			if string(got) != string(newData) {
				t.Fatalf("new policy %q not applied: got %s, want %s", modeID, got, newData)
			}
		}

		// If newPolicies was nil/empty, cfg.Policies should still be initialized (not nil)
		if cfg.Policies == nil {
			t.Fatal("cfg.Policies is nil after merge (should be initialized)")
		}
	})
}

func TestProperty_PolicyConfigMergeEmptyNewPolicies(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		// Generate config with existing policies
		cfg := genFileConfig(rt)

		// Deep copy old policies
		oldPolicies := make(map[string]json.RawMessage, len(cfg.Policies))
		for k, v := range cfg.Policies {
			cp := make(json.RawMessage, len(v))
			copy(cp, v)
			oldPolicies[k] = cp
		}

		origSocket := cfg.Socket

		// Merge with empty map (no new policies)
		mergePolicies(&cfg, map[string]json.RawMessage{})

		// Non-policy fields unchanged
		if cfg.Socket != origSocket {
			t.Fatalf("Socket changed: got %q, want %q", cfg.Socket, origSocket)
		}

		// All old policies preserved exactly
		for modeID, oldData := range oldPolicies {
			got, exists := cfg.Policies[modeID]
			if !exists {
				t.Fatalf("policy %q was removed by empty merge", modeID)
			}
			if string(got) != string(oldData) {
				t.Fatalf("policy %q was modified by empty merge: got %s, want %s", modeID, got, oldData)
			}
		}
	})
}
