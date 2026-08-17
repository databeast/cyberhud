package main

import (
	"encoding/json"
	"testing"

	"pgregory.net/rapid"
)

// For any PolicyStore containing N mode entries, querying All() should return a
// map with exactly N keys. No fields should be omitted and no extra modes should appear.

func TestProperty_PolicyDump_Completeness(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		// Generate a random number of mode entries (0 to 20).
		n := rapid.IntRange(0, 20).Draw(rt, "numModes")

		// Generate unique mode IDs and their policy JSON values.
		modeIDs := make([]string, 0, n)
		seen := make(map[string]bool)
		policies := make(map[string]json.RawMessage)

		for i := 0; i < n; i++ {
			// Generate a mode ID matching [a-z][a-z0-9_]* pattern.
			modeID := rapid.StringMatching(`[a-z][a-z0-9_]{0,19}`).Draw(rt, "modeID")
			if seen[modeID] {
				continue // skip duplicates
			}
			seen[modeID] = true
			modeIDs = append(modeIDs, modeID)

			// Generate random JSON object for the policy.
			numFields := rapid.IntRange(0, 5).Draw(rt, "numFields")
			obj := make(map[string]interface{})
			for f := 0; f < numFields; f++ {
				key := rapid.StringMatching(`[a-z_]{1,10}`).Draw(rt, "fieldKey")
				val := rapid.Float64Range(0.0, 100.0).Draw(rt, "fieldVal")
				obj[key] = val
			}
			data, err := json.Marshal(obj)
			if err != nil {
				t.Fatalf("failed to marshal policy JSON: %v", err)
			}
			policies[modeID] = json.RawMessage(data)
		}

		// Create PolicyStore and populate it.
		ps := NewPolicyStore()
		for _, id := range modeIDs {
			ps.Set(id, policies[id])
		}

		// Query All() — simulates "policy dump".
		all := ps.All()

		// Property: exactly len(modeIDs) keys in the result.
		if len(all) != len(modeIDs) {
			t.Fatalf("All() returned %d entries, want %d (modes added: %v)",
				len(all), len(modeIDs), modeIDs)
		}

		// Property: every mode ID we added is present with correct value.
		for _, id := range modeIDs {
			got, ok := all[id]
			if !ok {
				t.Fatalf("All() missing mode %q", id)
			}
			if string(got) != string(policies[id]) {
				t.Fatalf("All()[%q] = %s, want %s", id, got, policies[id])
			}
		}

		// Property: no extra modes appear (covered by length check, but verify explicitly).
		for k := range all {
			if !seen[k] {
				t.Fatalf("All() contains unexpected mode %q", k)
			}
		}
	})
}

// For any mode ID and any policy values, after calling Set for that mode,
// querying the PolicyStore via Get should return the stored values.

func TestProperty_PolicyStore_CapturesSetPolicy(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		ps := NewPolicyStore()

		// Generate a random mode ID.
		modeID := rapid.StringMatching(`[a-z][a-z0-9_]{0,19}`).Draw(rt, "modeID")

		// Generate random JSON bytes as the policy value.
		numFields := rapid.IntRange(0, 8).Draw(rt, "numFields")
		obj := make(map[string]interface{})
		for f := 0; f < numFields; f++ {
			key := rapid.StringMatching(`[a-z_]{1,12}`).Draw(rt, "key")
			// Mix of value types for realistic policy data.
			valType := rapid.IntRange(0, 2).Draw(rt, "valType")
			switch valType {
			case 0:
				obj[key] = rapid.Float64Range(0.0, 100.0).Draw(rt, "floatVal")
			case 1:
				obj[key] = rapid.IntRange(0, 1000).Draw(rt, "intVal")
			case 2:
				obj[key] = rapid.Bool().Draw(rt, "boolVal")
			}
		}
		data, err := json.Marshal(obj)
		if err != nil {
			t.Fatalf("failed to marshal policy: %v", err)
		}
		policyJSON := json.RawMessage(data)

		// Call Set (simulates SetPolicy notification to PolicyStore).
		ps.Set(modeID, policyJSON)

		// Query via Get — must return stored value.
		got := ps.Get(modeID)
		if got == nil {
			t.Fatalf("Get(%q) returned nil after Set", modeID)
		}
		if string(got) != string(policyJSON) {
			t.Fatalf("Get(%q) = %s, want %s", modeID, got, policyJSON)
		}

		// Also verify via All() that the mode appears.
		all := ps.All()
		allGot, ok := all[modeID]
		if !ok {
			t.Fatalf("All() missing mode %q after Set", modeID)
		}
		if string(allGot) != string(policyJSON) {
			t.Fatalf("All()[%q] = %s, want %s", modeID, allGot, policyJSON)
		}
	})
}

// For any mode ID, after calling Set twice with different values, Get returns
// the most recent value (last-write-wins semantics).

func TestProperty_PolicyStore_SetOverwritesExisting(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		ps := NewPolicyStore()

		modeID := rapid.StringMatching(`[a-z][a-z0-9_]{0,15}`).Draw(rt, "modeID")

		// First Set.
		obj1 := map[string]interface{}{
			"speed": rapid.Float64Range(0.1, 50.0).Draw(rt, "speed1"),
		}
		data1, _ := json.Marshal(obj1)
		ps.Set(modeID, json.RawMessage(data1))

		// Second Set with different value.
		obj2 := map[string]interface{}{
			"speed": rapid.Float64Range(50.1, 100.0).Draw(rt, "speed2"),
		}
		data2, _ := json.Marshal(obj2)
		ps.Set(modeID, json.RawMessage(data2))

		// Get must return the second value.
		got := ps.Get(modeID)
		if string(got) != string(data2) {
			t.Fatalf("Get(%q) after overwrite = %s, want %s", modeID, got, data2)
		}
	})
}
