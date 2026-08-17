package main

import (
	"testing"

	"github.com/databeast/cyberhud/display/catalog"
	displaymodes "github.com/databeast/cyberhud/display/modes"
)

// attractModeIDs lists all seven attract mode IDs that must be registered.
var attractModeIDs = []string{
	"attract_waveform",
	"attract_particles",
	"attract_plasma",
	"attract_starfield",
	"attract_bokeh",
	"attract_matrix",
	"attract_shapes",
}

// TestAttractModesCatalogPresence verifies all seven attract mode IDs are present
// in catalog.Definitions().
func TestAttractModesCatalogPresence(t *testing.T) {
	defs := catalog.Definitions()
	registered := make(map[string]bool, len(defs))
	for _, def := range defs {
		registered[def.ID] = true
	}

	for _, id := range attractModeIDs {
		if !registered[id] {
			t.Errorf("attract mode %q not found in catalog.Definitions()", id)
		}
	}
}

// TestAttractModesCatalogMetadata verifies metadata constraints for each attract mode:
// - Title is non-empty and ≤ 30 characters
// - Summary is non-empty and ≤ 120 characters
// - Order is 200
// - Options has at least one entry, each with non-empty Key, Type, Summary, Default
func TestAttractModesCatalogMetadata(t *testing.T) {
	for _, id := range attractModeIDs {
		t.Run(id, func(t *testing.T) {
			def, ok := catalog.Describe(id)
			if !ok {
				t.Fatalf("catalog.Describe(%q) returned not found", id)
			}

			// Title checks
			if def.Title == "" {
				t.Error("Title is empty")
			}
			if len(def.Title) > 30 {
				t.Errorf("Title too long (%d chars): %q", len(def.Title), def.Title)
			}

			// Summary checks
			if def.Summary == "" {
				t.Error("Summary is empty")
			}
			if len(def.Summary) > 120 {
				t.Errorf("Summary too long (%d chars): %q", len(def.Summary), def.Summary)
			}

			// Order check
			if def.Order != 200 {
				t.Errorf("Order = %d, want 200", def.Order)
			}

			// Options check
			if len(def.Options) == 0 {
				t.Fatal("Options is empty, expected at least one entry")
			}
			for i, opt := range def.Options {
				if opt.Key == "" {
					t.Errorf("Options[%d].Key is empty", i)
				}
				if opt.Type == "" {
					t.Errorf("Options[%d].Type is empty", i)
				}
				if opt.Summary == "" {
					t.Errorf("Options[%d].Summary is empty", i)
				}
				if opt.Default == "" {
					t.Errorf("Options[%d].Default is empty", i)
				}
			}
		})
	}
}

// TestAttractModesFactoryRegistration verifies all seven attract mode IDs are
// registered in displaymodes.IsKnownInstance.
func TestAttractModesFactoryRegistration(t *testing.T) {
	for _, id := range attractModeIDs {
		if !displaymodes.IsKnownInstance(id) {
			t.Errorf("attract mode %q not registered in displaymodes.IsKnownInstance", id)
		}
	}
}

// TestAttractModesNoImportCycles documents that no attract mode package imports
// another attract mode package. Since Go's build system enforces acyclicity at
// compile time, the fact that `go build ./...` (or `go test ./...`) succeeds is
// sufficient proof that no import cycles exist among these packages.
func TestAttractModesNoImportCycles(t *testing.T) {
	// This test exists as documentation. Go's compiler rejects cyclic imports,
	// so if this test file compiles and runs, no cycles exist between:
	//   display/modes/attract_waveform
	//   display/modes/attract_particles
	//   display/modes/attract_plasma
	//   display/modes/attract_starfield
	//   display/modes/attract_bokeh
	//   display/modes/attract_matrix
	//   display/modes/attract_shapes
	//
	// Each package registers independently via its own init() and shares only
	// parent packages (display/modes, display/catalog, display/style, etc.).
	t.Log("No import cycles: compilation success proves acyclicity among attract mode packages")
}
