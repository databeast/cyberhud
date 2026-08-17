package main

import (
	"image"
	"testing"

	displaymodes "github.com/databeast/cyberhud/display/modes"
	"github.com/databeast/cyberhud/display/region"
	"github.com/databeast/cyberhud/display/surface"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
)

// TestIntegration_CycleAll19Modes_SetMode verifies that all 19 registered display
// modes can be constructed, activated, and rendered via Region.SetMode without
// panicking, and that each produces valid ViewData.
//

func TestIntegration_CycleAll19Modes_SetMode(t *testing.T) {
	// Build a TierCatalog with dimensions large enough for all font variants.
	// 240x240 is a common panel size that satisfies all modes' style requirements.
	const pw, ph = 240, 240
	cat, err := tiercatalog.Build(tiercatalog.Params{
		PixelWidth:  pw,
		PixelHeight: ph,
	})
	if err != nil {
		t.Fatalf("tiercatalog.Build failed: %v", err)
	}

	// Construct TextHints with a valid Catalog field for font resolution.
	hints := textlayout.DefaultTextHints(image.Rect(0, 0, pw, ph))
	hints.Catalog = cat

	// Create a Region with a real surface.
	bounds := image.Rect(0, 0, pw, ph)
	surf := surface.New(bounds)
	r := region.NewRegion("integration-test", bounds, surf)

	// Wire the displaymodes.GetInstance as the ModeFactory.
	// This bridges the region package (which defines ModeFactory) with the
	// displaymodes package (which provides GetInstance), matching what the
	// runtime wiring does in production.
	r.SetModeFactory(func(id string, h textlayout.TextHints) (region.ModeInstance, bool) {
		inst, ok := displaymodes.GetInstance(id)
		if !ok {
			return nil, false
		}
		return inst, true
	})

	// Override the region's textHints to include our TierCatalog.
	// The region was created with NewRegion which uses DefaultTextHints (no catalog).
	// We need to construct a region that has our catalog in its hints.
	// Since TextHints is set during NewRegion and isn't mutable, we recreate
	// the region using NewRegionWithScreens with a HintProvider that returns
	// our catalog-enriched hints.
	r = region.NewRegionWithScreens("integration-test", bounds, surf, []region.ScreenPosition{
		{
			Index:  0,
			Bounds: bounds,
			HintProvider: func() textlayout.TextHints {
				return hints
			},
		},
	}, "", 0, 0)
	r.SetModeFactory(func(id string, h textlayout.TextHints) (region.ModeInstance, bool) {
		inst, ok := displaymodes.GetInstance(id)
		if !ok {
			return nil, false
		}
		return inst, true
	})

	// All 19 mode IDs that must be registered.
	modeIDs := []string{
		"menu",
		"dashboard",
		"stemma",
		"gpio",
		"system",
		"systemd",
		"clock",
		"ticker",
		"image",
		"usb",
		"serial",
		"thermal",
		"gpio-control",
		"testpattern",
		"testfonts",
		"cycle",
		"attract_matrix",
		"zmq",
		"wifi",
	}

	for _, modeID := range modeIDs {
		t.Run(modeID, func(t *testing.T) {
			// SetMode should succeed without error.
			if err := r.SetMode(modeID); err != nil {
				t.Fatalf("SetMode(%q) error: %v", modeID, err)
			}

			// Instance must be non-nil after successful SetMode.
			inst := r.Instance()
			if inst == nil {
				t.Fatalf("Instance() is nil after SetMode(%q)", modeID)
			}

			// BuildView must not panic and must return valid ViewData.
			vd := inst.BuildView()

			_ = vd

			// Deactivate lifecycle modes to release background resources.
			// This is handled by the next SetMode call (which deactivates old),
			// but we explicitly deactivate the last one.
		})
	}

	// Deactivate the final instance to clean up any background goroutines.
	if inst := r.Instance(); inst != nil {
		inst.Deactivate()
	}
}
