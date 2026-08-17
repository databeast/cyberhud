package tests_test

import (
	"image"
	"testing"

	"github.com/databeast/cyberhud/display/region"
	"github.com/databeast/cyberhud/display/surface"
	"github.com/databeast/cyberhud/display/surface/fonts"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// registerTestFonts registers fonts that fit a 128px wide region (maxAdvance = 128/10 = 12).
func registerTestFonts() {
	fonts := []stubFace{
		{id: "spleen-5x8", metrics: font.Metrics{GlyphWidth: 5, GlyphHeight: 8, GlyphAdvance: 6, RowHeight: 10}},
		{id: "spleen-8x16", metrics: font.Metrics{GlyphWidth: 8, GlyphHeight: 16, GlyphAdvance: 9, RowHeight: 18}},
	}
	for _, f := range fonts {
		font.Register(f)
	}
}

// TestResolveTextHints_Branch1_NoScreens verifies that when a Region has no screen
// positions, the catalog is populated from the Region's bounds dimensions.
func TestResolveTextHints_Branch1_NoScreens(t *testing.T) {
	restore := font.SnapshotAndClear()
	defer restore()
	registerTestFonts()

	bounds := image.Rect(0, 0, 128, 64)
	surf := surface.New(bounds)
	r := region.NewRegionWithScreens("test-branch1", bounds, surf, nil, "", 0, 0)

	hints := r.TextHints()
	if hints.Catalog.PixelWidth() != 128 {
		t.Errorf("Branch1: Catalog.PixelWidth() = %d, want 128", hints.Catalog.PixelWidth())
	}
	if hints.Catalog.PixelHeight() != 64 {
		t.Errorf("Branch1: Catalog.PixelHeight() = %d, want 64", hints.Catalog.PixelHeight())
	}
}

// TestResolveTextHints_Branch2_WithHintProvider verifies that when a Region is entirely
// within one Physical Screen that provides a HintProvider, the catalog uses the Region's
// resolved dimensions.
func TestResolveTextHints_Branch2_WithHintProvider(t *testing.T) {
	restore := font.SnapshotAndClear()
	defer restore()
	registerTestFonts()

	// Region is 128x64, screen is larger (256x128) so Region.In(Screen.Bounds) is true.
	regionBounds := image.Rect(0, 0, 128, 64)
	screenBounds := image.Rect(0, 0, 256, 128)

	screens := []region.ScreenPosition{
		{
			Index:  0,
			Name:   "main",
			Bounds: screenBounds,
			HintProvider: func() textlayout.TextHints {
				// Return hints without a catalog — auto-build should populate it.
				return textlayout.TextHints{
					PixelWidth:               256,
					PixelHeight:              128,
					GlyphWidth:               5,
					GlyphHeight:              7,
					GlyphAdvance:             6,
					RowHeight:                10,
					SupportsVerticalScroll:   true,
					SupportsHorizontalScroll: false, // hardware-specific flag
				}
			},
		},
	}

	surf := surface.New(regionBounds)
	r := region.NewRegionWithScreens("test-branch2-hint", regionBounds, surf, screens, "", 0, 0)

	hints := r.TextHints()
	// The catalog should be built with the Region's dimensions, not the screen's.
	if hints.Catalog.PixelWidth() != 128 {
		t.Errorf("Branch2 HintProvider: Catalog.PixelWidth() = %d, want 128", hints.Catalog.PixelWidth())
	}
	if hints.Catalog.PixelHeight() != 64 {
		t.Errorf("Branch2 HintProvider: Catalog.PixelHeight() = %d, want 64", hints.Catalog.PixelHeight())
	}
	// Verify hardware-specific flags were preserved from HintProvider.
	if hints.SupportsHorizontalScroll {
		t.Errorf("Branch2 HintProvider: SupportsHorizontalScroll should be false (from HintProvider)")
	}
}

// TestResolveTextHints_Branch2_NoHintProvider verifies that when a Region is entirely
// within one Physical Screen that lacks a HintProvider, the catalog is built from
// the Region's bounds dimensions using defaults.
func TestResolveTextHints_Branch2_NoHintProvider(t *testing.T) {
	restore := font.SnapshotAndClear()
	defer restore()
	registerTestFonts()

	regionBounds := image.Rect(0, 0, 128, 64)
	screenBounds := image.Rect(0, 0, 256, 128)

	screens := []region.ScreenPosition{
		{
			Index:        0,
			Name:         "main",
			Bounds:       screenBounds,
			HintProvider: nil, // no HintProvider
		},
	}

	surf := surface.New(regionBounds)
	r := region.NewRegionWithScreens("test-branch2-nohint", regionBounds, surf, screens, "", 0, 0)

	hints := r.TextHints()
	if hints.Catalog.PixelWidth() != 128 {
		t.Errorf("Branch2 NoHintProvider: Catalog.PixelWidth() = %d, want 128", hints.Catalog.PixelWidth())
	}
	if hints.Catalog.PixelHeight() != 64 {
		t.Errorf("Branch2 NoHintProvider: Catalog.PixelHeight() = %d, want 64", hints.Catalog.PixelHeight())
	}
}

// TestResolveTextHints_Branch3_MultiScreen verifies that when a Region spans multiple
// Physical Screens (neither fully contains the Region), the catalog is built from
// the Region's bounds dimensions.
func TestResolveTextHints_Branch3_MultiScreen(t *testing.T) {
	restore := font.SnapshotAndClear()
	defer restore()
	registerTestFonts()

	// Region spans from x=0 to x=200, but no single screen contains it fully.
	regionBounds := image.Rect(0, 0, 200, 64)

	// Screen 1 covers x=0..128, Screen 2 covers x=100..256. Both overlap
	// the region but neither fully contains it.
	screens := []region.ScreenPosition{
		{
			Index:  0,
			Name:   "left",
			Bounds: image.Rect(0, 0, 128, 64),
			HintProvider: func() textlayout.TextHints {
				return textlayout.TextHints{PixelWidth: 128, PixelHeight: 64}
			},
		},
		{
			Index:  1,
			Name:   "right",
			Bounds: image.Rect(100, 0, 256, 64),
			HintProvider: func() textlayout.TextHints {
				return textlayout.TextHints{PixelWidth: 156, PixelHeight: 64}
			},
		},
	}

	surf := surface.New(regionBounds)
	r := region.NewRegionWithScreens("test-branch3", regionBounds, surf, screens, "", 0, 0)

	hints := r.TextHints()
	// Catalog should use the Region's own bounds (200x64).
	if hints.Catalog.PixelWidth() != 200 {
		t.Errorf("Branch3 MultiScreen: Catalog.PixelWidth() = %d, want 200", hints.Catalog.PixelWidth())
	}
	if hints.Catalog.PixelHeight() != 64 {
		t.Errorf("Branch3 MultiScreen: Catalog.PixelHeight() = %d, want 64", hints.Catalog.PixelHeight())
	}
}

// TestResolveTextHints_ZeroBounds verifies that when a Region has 0×0 bounds,
// the catalog remains zero-valued (no panic, no error).
func TestResolveTextHints_ZeroBounds(t *testing.T) {
	restore := font.SnapshotAndClear()
	defer restore()
	registerTestFonts()

	bounds := image.Rect(0, 0, 0, 0)
	surf := surface.New(bounds)
	r := region.NewRegionWithScreens("test-zero", bounds, surf, nil, "", 0, 0)

	hints := r.TextHints()
	if hints.Catalog.PixelWidth() != 0 {
		t.Errorf("ZeroBounds: Catalog.PixelWidth() = %d, want 0", hints.Catalog.PixelWidth())
	}
	if hints.Catalog.PixelHeight() != 0 {
		t.Errorf("ZeroBounds: Catalog.PixelHeight() = %d, want 0", hints.Catalog.PixelHeight())
	}
}
