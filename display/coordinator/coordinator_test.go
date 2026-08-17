package coordinator_test

import (
	"testing"

	"github.com/databeast/cyberhud/display/coordinator"
	_ "github.com/databeast/cyberhud/display/modes/clock"
	_ "github.com/databeast/cyberhud/display/modes/dashboard"
	_ "github.com/databeast/cyberhud/display/modes/menu"
	"pgregory.net/rapid"
)

func TestStateSetNextPrev(t *testing.T) {
	s := coordinator.NewState(coordinator.Region{
		Index:      1,
		Name:       "aux",
		Controller: "st7735s",
		Modes:      []string{"stemma", "gpio", "clock"},
		Default:    "stemma",
	})
	if got := s.CurrentMode(1); got != "stemma" {
		t.Fatalf("CurrentMode() = %q, want stemma", got)
	}
	if got, err := s.Next(1); err != nil || got != "gpio" {
		t.Fatalf("Next() got=%q err=%v", got, err)
	}
	if got, err := s.Prev(1); err != nil || got != "stemma" {
		t.Fatalf("Prev() got=%q err=%v", got, err)
	}
	if got, err := s.Set(1, "CLOCK"); err != nil || got != "clock" {
		t.Fatalf("Set() got=%q err=%v", got, err)
	}
}

func TestStateStatusSorted(t *testing.T) {
	s := coordinator.NewState(
		coordinator.Region{Index: 2, Name: "right", Modes: []string{"gpio"}},
		coordinator.Region{Index: 0, Name: "main", Modes: []string{"menu", "dashboard"}, Default: "dashboard"},
	)
	st := s.Status()
	if len(st) != 2 {
		t.Fatalf("Status() len=%d, want 2", len(st))
	}
	if st[0].Index != 0 || st[1].Index != 2 {
		t.Fatalf("Status() indexes=%v,%v, want 0,2", st[0].Index, st[1].Index)
	}
	if st[0].Current != "dashboard" {
		t.Fatalf("main current=%q, want dashboard", st[0].Current)
	}
}

func TestStateRegionLookup(t *testing.T) {
	s := coordinator.NewState(coordinator.Region{
		Index:      3,
		Name:       "aux",
		Controller: "st7735s",
		Modes:      []string{"stemma", "gpio"},
		Default:    "gpio",
	})
	p, ok := s.Region(3)
	if !ok {
		t.Fatal("Region(3) expected ok=true")
	}
	if p.Current != "gpio" {
		t.Fatalf("Region(3).Current=%q, want gpio", p.Current)
	}
	if !s.HasRegion(3) {
		t.Fatal("HasRegion(3)=false, want true")
	}
	if s.HasRegion(4) {
		t.Fatal("HasRegion(4)=true, want false")
	}
}

func TestStateDefinitions(t *testing.T) {
	s := coordinator.NewState(coordinator.Region{
		Index:      0,
		Name:       "main",
		Controller: "st7789",
		Modes:      []string{"menu", "dashboard", "custom-x"},
		Default:    "dashboard",
	})
	defs := s.Definitions()
	if len(defs) != 1 {
		t.Fatalf("Definitions() len=%d, want 1", len(defs))
	}
	if defs[0].Current != "dashboard" {
		t.Fatalf("Definitions()[0].Current=%q, want dashboard", defs[0].Current)
	}
	if len(defs[0].Modes) != 3 {
		t.Fatalf("Definitions()[0].Modes len=%d, want 3", len(defs[0].Modes))
	}
	if defs[0].Modes[0].Title != "Menu" {
		t.Fatalf("expected built-in title for menu, got %+v", defs[0].Modes[0])
	}
	if defs[0].Modes[2].Scope != "display mode" {
		t.Fatalf("expected display mode scope fallback, got %+v", defs[0].Modes[2])
	}
}

func TestProperty_StateFallback_EmptyModesUsesFallback(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		// Generate an arbitrary non-empty fallback mode string.
		fallback := rapid.StringMatching(`[a-z][a-z0-9\-]{0,19}`).Draw(rt, "fallbackMode")
		panelIndex := rapid.IntRange(0, 10).Draw(rt, "panelIndex")

		// Create a region with zero modes.
		regions := []coordinator.Region{
			{
				Index: panelIndex,
				Name:  "test-panel",
				Modes: []string{}, // zero configured modes
			},
		}

		s := coordinator.NewState()
		s.ResetWithFallback(regions, fallback)

		// After reset, the region should have exactly [fallback] as modes.
		ps, ok := s.Region(panelIndex)
		if !ok {
			rt.Fatalf("Region(%d) not found after ResetWithFallback", panelIndex)
		}
		if len(ps.Modes) != 1 {
			rt.Fatalf("expected exactly 1 mode, got %d: %v", len(ps.Modes), ps.Modes)
		}
		if ps.Modes[0] != fallback {
			rt.Fatalf("expected mode %q, got %q", fallback, ps.Modes[0])
		}
		// The current mode should also be the fallback.
		if ps.Current != fallback {
			rt.Fatalf("expected current=%q, got %q", fallback, ps.Current)
		}
	})
}

func TestProperty_StateFallback_NonEmptyModesUnchanged(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		// Generate an arbitrary fallback mode string.
		fallback := rapid.StringMatching(`[a-z][a-z0-9\-]{0,19}`).Draw(rt, "fallbackMode")
		panelIndex := rapid.IntRange(0, 10).Draw(rt, "panelIndex")

		// Generate a non-empty list of modes (1 to 5 modes).
		numModes := rapid.IntRange(1, 5).Draw(rt, "numModes")
		modes := make([]string, numModes)
		seen := map[string]bool{}
		for i := 0; i < numModes; i++ {
			for {
				m := rapid.StringMatching(`[a-z][a-z0-9\-]{0,9}`).Draw(rt, "mode")
				if !seen[m] {
					seen[m] = true
					modes[i] = m
					break
				}
			}
		}

		panels := []coordinator.Region{
			{
				Index:   panelIndex,
				Name:    "test-panel",
				Modes:   modes,
				Default: modes[0],
			},
		}

		s := coordinator.NewState()
		s.ResetWithFallback(panels, fallback)

		// After reset, the region should retain its original modes (not replaced by fallback).
		ps, ok := s.Region(panelIndex)
		if !ok {
			rt.Fatalf("Region(%d) not found after ResetWithFallback", panelIndex)
		}
		if len(ps.Modes) != numModes {
			rt.Fatalf("expected %d modes, got %d: %v", numModes, len(ps.Modes), ps.Modes)
		}
		for i, m := range modes {
			if ps.Modes[i] != m {
				rt.Fatalf("mode[%d]: expected %q, got %q", i, m, ps.Modes[i])
			}
		}
	})
}
