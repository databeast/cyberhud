package main

import (
	"encoding/json"
	"sync"
	"testing"
)

func TestNewPolicyStore(t *testing.T) {
	ps := NewPolicyStore()
	if ps == nil {
		t.Fatal("NewPolicyStore returned nil")
	}
	if ps.policies == nil {
		t.Fatal("policies map not initialized")
	}
}

func TestPolicyStore_GetNil(t *testing.T) {
	ps := NewPolicyStore()
	got := ps.Get("nonexistent")
	if got != nil {
		t.Fatalf("expected nil for missing mode, got %v", got)
	}
}

func TestPolicyStore_SetAndGet(t *testing.T) {
	ps := NewPolicyStore()
	data := json.RawMessage(`{"speed":1.5,"density":0.8}`)
	ps.Set("attract_bokeh", data)

	got := ps.Get("attract_bokeh")
	if got == nil {
		t.Fatal("expected stored policy, got nil")
	}
	if string(got) != string(data) {
		t.Fatalf("expected %s, got %s", data, got)
	}
}

func TestPolicyStore_SetOverwrite(t *testing.T) {
	ps := NewPolicyStore()
	ps.Set("clock", json.RawMessage(`{"style":"digital"}`))
	ps.Set("clock", json.RawMessage(`{"style":"analog"}`))

	got := ps.Get("clock")
	if string(got) != `{"style":"analog"}` {
		t.Fatalf("expected overwritten value, got %s", got)
	}
}

func TestPolicyStore_SetOnNilMap(t *testing.T) {
	// Test that Set creates the map if nil (defensive coding path)
	ps := &PolicyStore{}
	ps.Set("test_mode", json.RawMessage(`{"key":"value"}`))

	got := ps.Get("test_mode")
	if string(got) != `{"key":"value"}` {
		t.Fatalf("expected stored value after nil-map init, got %s", got)
	}
}

func TestPolicyStore_All(t *testing.T) {
	ps := NewPolicyStore()
	ps.Set("mode_a", json.RawMessage(`{"a":1}`))
	ps.Set("mode_b", json.RawMessage(`{"b":2}`))

	all := ps.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(all))
	}
	if string(all["mode_a"]) != `{"a":1}` {
		t.Fatalf("mode_a mismatch: %s", all["mode_a"])
	}
	if string(all["mode_b"]) != `{"b":2}` {
		t.Fatalf("mode_b mismatch: %s", all["mode_b"])
	}
}

func TestPolicyStore_AllReturnsDefensiveCopy(t *testing.T) {
	ps := NewPolicyStore()
	ps.Set("mode_x", json.RawMessage(`{"x":true}`))

	cp := ps.All()
	// Mutating the copy should not affect the store
	cp["mode_x"] = json.RawMessage(`{"x":false}`)
	cp["injected"] = json.RawMessage(`{}`)

	got := ps.Get("mode_x")
	if string(got) != `{"x":true}` {
		t.Fatalf("store was mutated via All() copy: %s", got)
	}
	if ps.Get("injected") != nil {
		t.Fatal("injected key visible in store via All() copy mutation")
	}
}

func TestPolicyStore_LoadFromConfig(t *testing.T) {
	ps := NewPolicyStore()
	cfg := &fileConfig{
		Policies: map[string]json.RawMessage{
			"attract_matrix": json.RawMessage(`{"density":0.9}`),
			"clock":          json.RawMessage(`{"style":"led"}`),
		},
	}
	ps.LoadFromConfig(cfg)

	got := ps.Get("attract_matrix")
	if string(got) != `{"density":0.9}` {
		t.Fatalf("attract_matrix policy mismatch: %s", got)
	}
	got = ps.Get("clock")
	if string(got) != `{"style":"led"}` {
		t.Fatalf("clock policy mismatch: %s", got)
	}
}

func TestPolicyStore_LoadFromConfigNil(t *testing.T) {
	ps := NewPolicyStore()
	// Should not panic
	ps.LoadFromConfig(nil)
	all := ps.All()
	if len(all) != 0 {
		t.Fatalf("expected empty store after nil config, got %d entries", len(all))
	}
}

func TestPolicyStore_LoadFromConfigMerges(t *testing.T) {
	ps := NewPolicyStore()
	ps.Set("existing", json.RawMessage(`{"keep":"me"}`))

	cfg := &fileConfig{
		Policies: map[string]json.RawMessage{
			"new_mode": json.RawMessage(`{"added":true}`),
		},
	}
	ps.LoadFromConfig(cfg)

	// Both should exist
	if ps.Get("existing") == nil {
		t.Fatal("existing entry lost after LoadFromConfig")
	}
	if ps.Get("new_mode") == nil {
		t.Fatal("new entry not loaded from config")
	}
}

func TestPolicyStore_ConcurrentAccess(t *testing.T) {
	ps := NewPolicyStore()
	var wg sync.WaitGroup

	// Concurrent writes
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			ps.Set("mode", json.RawMessage(`{"n":`+string(rune('0'+n%10))+`}`))
		}(i)
	}

	// Concurrent reads
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = ps.Get("mode")
			_ = ps.All()
		}()
	}

	wg.Wait()
	// No race detected = success (run with -race flag)
}
