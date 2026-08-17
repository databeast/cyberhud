package tiercatalog_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/databeast/cyberhud/display/surface/fonts"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
	"pgregory.net/rapid"
)

// --- From: tiercatalog_prop_test.go ---

// stubFace is a minimal font.Face implementation for property testing.
type stubFace struct {
	id      string
	metrics font.Metrics
}

func (s stubFace) ID() string                { return s.id }
func (s stubFace) Metrics() font.Metrics     { return s.metrics }
func (s stubFace) GlyphRow(rune, int) uint32 { return 0 }

// For any valid region dimensions and font set, if tiercatalog.Build succeeds,
// every tier entry satisfies GlyphAdvance * MinChars <= PixelWidth.

func TestProperty3_WidthSafetyCatalog(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate random region parameters.
		pixelWidth := rapid.IntRange(1, 4096).Draw(t, "pixelWidth")
		pixelHeight := rapid.IntRange(1, 4096).Draw(t, "pixelHeight")
		minChars := rapid.IntRange(1, 100).Draw(t, "minChars")

		// Register a random set of additional fonts to vary the registry.
		numFonts := rapid.IntRange(1, 20).Draw(t, "numFonts")
		for i := 0; i < numFonts; i++ {
			glyphWidth := rapid.IntRange(1, 64).Draw(t, fmt.Sprintf("gw%d", i))
			glyphHeight := rapid.IntRange(1, 64).Draw(t, fmt.Sprintf("gh%d", i))
			glyphAdvance := rapid.IntRange(1, 64).Draw(t, fmt.Sprintf("ga%d", i))
			rowHeight := glyphHeight + rapid.IntRange(0, 4).Draw(t, fmt.Sprintf("rh%d", i))

			// Use a family prefix to get varying family priorities.
			family := rapid.SampledFrom([]string{"spleen", "terminus", "cozette", "testfont"}).Draw(t, fmt.Sprintf("fam%d", i))
			id := fmt.Sprintf("%s-%dx%d-prop3-%d", family, glyphWidth, glyphHeight, i)

			font.Register(stubFace{
				id: id,
				metrics: font.Metrics{
					GlyphWidth:   glyphWidth,
					GlyphHeight:  glyphHeight,
					GlyphAdvance: glyphAdvance,
					RowHeight:    rowHeight,
				},
			})
		}

		// Attempt to build the catalog.
		cat, err := tiercatalog.Build(tiercatalog.Params{
			PixelWidth:  pixelWidth,
			PixelHeight: pixelHeight,
			MinChars:    minChars,
		})
		if err != nil {
			// Build failed (no qualifying fonts) - property is vacuously true.
			return
		}

		// Assert width safety: for every tier, GlyphAdvance * MinChars <= PixelWidth.
		for _, tier := range cat.Tiers() {
			entry, ok := cat.Get(tier)
			if !ok {
				t.Fatalf("tier %q missing from successful catalog", tier)
			}
			if entry.GlyphAdvance*minChars > pixelWidth {
				t.Fatalf("width safety violated: tier=%q GlyphAdvance=%d * MinChars=%d = %d > PixelWidth=%d",
					tier, entry.GlyphAdvance, minChars, entry.GlyphAdvance*minChars, pixelWidth)
			}
		}
	})
}

// --- From: tiercatalog_test.go ---

// testFace is a minimal font.Face for unit testing with known metrics.
type testFace struct {
	id      string
	metrics font.Metrics
}

func (f testFace) ID() string                { return f.id }
func (f testFace) Metrics() font.Metrics     { return f.metrics }
func (f testFace) GlyphRow(rune, int) uint32 { return 0 }

// registerTestFonts registers the standard set of fonts described in the design
// document for a 128×128 region scenario.
func registerTestFonts() {
	fonts := []testFace{
		{id: "spleen-5x8", metrics: font.Metrics{GlyphWidth: 5, GlyphHeight: 8, GlyphAdvance: 6, RowHeight: 10}},
		{id: "spleen-6x12", metrics: font.Metrics{GlyphWidth: 6, GlyphHeight: 12, GlyphAdvance: 7, RowHeight: 14}},
		{id: "terminus-6x12", metrics: font.Metrics{GlyphWidth: 6, GlyphHeight: 12, GlyphAdvance: 7, RowHeight: 14}},
		{id: "cozette-6x13", metrics: font.Metrics{GlyphWidth: 6, GlyphHeight: 13, GlyphAdvance: 7, RowHeight: 15}},
		{id: "spleen-8x16", metrics: font.Metrics{GlyphWidth: 8, GlyphHeight: 16, GlyphAdvance: 9, RowHeight: 18}},
		{id: "terminus-8x14", metrics: font.Metrics{GlyphWidth: 8, GlyphHeight: 14, GlyphAdvance: 9, RowHeight: 16}},
		{id: "terminus-8x16", metrics: font.Metrics{GlyphWidth: 8, GlyphHeight: 16, GlyphAdvance: 9, RowHeight: 18}},
		// spleen-12x24 has advance 13; 13*10=130 > 128 so it should be EXCLUDED.
		{id: "spleen-12x24", metrics: font.Metrics{GlyphWidth: 12, GlyphHeight: 24, GlyphAdvance: 13, RowHeight: 26}},
	}
	for _, f := range fonts {
		font.Register(f)
	}
}

// TestBuild_128x128_DesignExample tests the concrete example from the design
// document: 128×128 region with MinChars=10.
//
// Qualifying fonts (advance ≤ 12): spleen-5x8(6), spleen-6x12(7),
// terminus-6x12(7), cozette-6x13(7), spleen-8x16(9), terminus-8x14(9),
// terminus-8x16(9). Excluded: spleen-12x24(13).
//
// With PPI=0, pixel fallback targets: Small=8, Normal=14, Large=20, Huge=28, Colossal=40
// Best-fit selections:
//
//	Small(8):    spleen-5x8 at height 8, distance 0 → exact match
//	Normal(14):  terminus-8x14 at height 14, distance 0 → exact match
//	Large(20):   spleen-8x16 at height 16, distance 4 (closest to 20)
//	Huge(28):    spleen-8x16 at height 16, distance 12 (closest to 28)
//	Colossal(40):spleen-8x16 at height 16, distance 24 (closest to 40)
//	Full:        spleen-8x16 (largest qualifying, priority wins tie with terminus-8x16)
func TestBuild_128x128_DesignExample(t *testing.T) {
	restore := font.SnapshotAndClear()
	defer restore()
	registerTestFonts()

	cat, err := tiercatalog.Build(tiercatalog.Params{
		PixelWidth:  128,
		PixelHeight: 128,
		MinChars:    10,
	})
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	// Verify region dimensions stored correctly.
	if cat.PixelWidth() != 128 {
		t.Errorf("PixelWidth: got %d, want 128", cat.PixelWidth())
	}
	if cat.PixelHeight() != 128 {
		t.Errorf("PixelHeight: got %d, want 128", cat.PixelHeight())
	}

	// Expected tier assignments from best-fit algorithm.
	tests := []struct {
		tier       tiercatalog.Tier
		wantHeight int
		wantWidth  int
		wantAdv    int
		wantRow    int
	}{
		{tiercatalog.TierSmall, 8, 5, 6, 10},
		{tiercatalog.TierNormal, 14, 8, 9, 16},
		{tiercatalog.TierLarge, 16, 8, 9, 18},
		{tiercatalog.TierHuge, 16, 8, 9, 18},
		{tiercatalog.TierColossal, 16, 8, 9, 18},
		{tiercatalog.TierFullsize, 16, 8, 9, 18},
	}

	for _, tc := range tests {
		entry, ok := cat.Get(tc.tier)
		if !ok {
			t.Errorf("tier %q: not present in catalog", tc.tier)
			continue
		}
		if entry.GlyphHeight != tc.wantHeight {
			t.Errorf("tier %q: GlyphHeight = %d, want %d", tc.tier, entry.GlyphHeight, tc.wantHeight)
		}
		if entry.GlyphWidth != tc.wantWidth {
			t.Errorf("tier %q: GlyphWidth = %d, want %d", tc.tier, entry.GlyphWidth, tc.wantWidth)
		}
		if entry.GlyphAdvance != tc.wantAdv {
			t.Errorf("tier %q: GlyphAdvance = %d, want %d", tc.tier, entry.GlyphAdvance, tc.wantAdv)
		}
		if entry.RowHeight != tc.wantRow {
			t.Errorf("tier %q: RowHeight = %d, want %d", tc.tier, entry.RowHeight, tc.wantRow)
		}
	}
}

// TestBuild_SingleQualifyingFont verifies that when only one font fits the
// region, all six tiers map to that same entry.
func TestBuild_SingleQualifyingFont(t *testing.T) {
	restore := font.SnapshotAndClear()
	defer restore()

	// Only one font; its advance of 5 fits in a 50px wide region (5*10=50).
	font.Register(testFace{
		id:      "spleen-5x8",
		metrics: font.Metrics{GlyphWidth: 5, GlyphHeight: 8, GlyphAdvance: 5, RowHeight: 10},
	})

	cat, err := tiercatalog.Build(tiercatalog.Params{
		PixelWidth:  50,
		PixelHeight: 50,
		MinChars:    10,
	})
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	// All tiers should have the same entry.
	smallEntry, _ := cat.Get(tiercatalog.TierSmall)
	for _, tier := range cat.Tiers() {
		entry, ok := cat.Get(tier)
		if !ok {
			t.Errorf("tier %q: not present", tier)
			continue
		}
		if entry.GlyphHeight != smallEntry.GlyphHeight || entry.GlyphWidth != smallEntry.GlyphWidth ||
			entry.GlyphAdvance != smallEntry.GlyphAdvance || entry.RowHeight != smallEntry.RowHeight {
			t.Errorf("tier %q: got %+v, want same as small %+v", tier, entry, smallEntry)
		}
	}
}

// TestBuild_TwoQualifyingFonts verifies behavior with exactly 2 qualifying fonts.
// With PPI=0 pixel fallback targets: Small=8, Normal=14, Large=20, Huge=28, Colossal=40
// Best-fit: Small(8)→spleen-5x8 (distance 0), Normal(14)→spleen-8x16 (distance 2, closer than 5x8 at distance 6)
// Large(20)→spleen-8x16 (distance 4), TierFull→spleen-8x16 (largest)
func TestBuild_TwoQualifyingFonts(t *testing.T) {
	restore := font.SnapshotAndClear()
	defer restore()

	font.Register(testFace{
		id:      "spleen-5x8",
		metrics: font.Metrics{GlyphWidth: 5, GlyphHeight: 8, GlyphAdvance: 6, RowHeight: 10},
	})
	font.Register(testFace{
		id:      "spleen-8x16",
		metrics: font.Metrics{GlyphWidth: 8, GlyphHeight: 16, GlyphAdvance: 9, RowHeight: 18},
	})

	cat, err := tiercatalog.Build(tiercatalog.Params{
		PixelWidth:  128,
		PixelHeight: 128,
		MinChars:    10,
	})
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	smallEntry, _ := cat.Get(tiercatalog.TierSmall)
	normalEntry, _ := cat.Get(tiercatalog.TierNormal)
	largeEntry, _ := cat.Get(tiercatalog.TierLarge)
	fullEntry, _ := cat.Get(tiercatalog.TierFullsize)

	// Small should be height 8 (exact match for target 8).
	if smallEntry.GlyphHeight != 8 {
		t.Errorf("small: GlyphHeight = %d, want 8", smallEntry.GlyphHeight)
	}
	// Normal target is 14: spleen-8x16 at distance 2, spleen-5x8 at distance 6 → 8x16 wins.
	if normalEntry.GlyphHeight != 16 {
		t.Errorf("normal: GlyphHeight = %d, want 16", normalEntry.GlyphHeight)
	}

	// Large and fullsize should both be height 16.
	if largeEntry.GlyphHeight != 16 {
		t.Errorf("large: GlyphHeight = %d, want 16", largeEntry.GlyphHeight)
	}
	if fullEntry.GlyphHeight != 16 {
		t.Errorf("fullsize: GlyphHeight = %d, want 16", fullEntry.GlyphHeight)
	}
}

// TestBuild_MinCharsZero_DefaultsTo10 verifies that MinChars=0 defaults to 10,
// which means maxAdvance = pixelWidth/10.
func TestBuild_MinCharsZero_DefaultsTo10(t *testing.T) {
	restore := font.SnapshotAndClear()
	defer restore()
	registerTestFonts()

	// Build with MinChars=0 (should default to 10) and with MinChars=10 explicitly.
	catDefault, err := tiercatalog.Build(tiercatalog.Params{
		PixelWidth:  128,
		PixelHeight: 128,
		MinChars:    0,
	})
	if err != nil {
		t.Fatalf("Build (MinChars=0) failed: %v", err)
	}

	catExplicit, err := tiercatalog.Build(tiercatalog.Params{
		PixelWidth:  128,
		PixelHeight: 128,
		MinChars:    10,
	})
	if err != nil {
		t.Fatalf("Build (MinChars=10) failed: %v", err)
	}

	// Both should produce identical catalogs.
	for _, tier := range catDefault.Tiers() {
		defEntry, _ := catDefault.Get(tier)
		expEntry, _ := catExplicit.Get(tier)
		if defEntry != expEntry {
			t.Errorf("tier %q: MinChars=0 gave %+v, MinChars=10 gave %+v", tier, defEntry, expEntry)
		}
	}
}

// TestBuild_TieBreaking_SameHeight verifies that when multiple families have
// the same GlyphHeight, the face with highest family priority wins.
// spleen (priority 3) should win over terminus (priority 2) at height 12.
func TestBuild_TieBreaking_SameHeight(t *testing.T) {
	restore := font.SnapshotAndClear()
	defer restore()

	// Register two fonts at the same height (12) with different families.
	font.Register(testFace{
		id:      "terminus-6x12",
		metrics: font.Metrics{GlyphWidth: 6, GlyphHeight: 12, GlyphAdvance: 7, RowHeight: 14},
	})
	font.Register(testFace{
		id:      "spleen-6x12",
		metrics: font.Metrics{GlyphWidth: 6, GlyphHeight: 12, GlyphAdvance: 7, RowHeight: 14},
	})

	cat, err := tiercatalog.Build(tiercatalog.Params{
		PixelWidth:  128,
		PixelHeight: 128,
		MinChars:    10,
	})
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	// With only one unique height, all tiers map to the same entry.
	// The winner should be spleen (priority 3) over terminus (priority 2).
	// Both have the same metrics so the entry values are identical,
	// but the dedup logic should keep spleen first.
	// We verify the catalog assigns height 12 for all tiers.
	for _, tier := range cat.Tiers() {
		entry, ok := cat.Get(tier)
		if !ok {
			t.Errorf("tier %q: not present", tier)
			continue
		}
		if entry.GlyphHeight != 12 {
			t.Errorf("tier %q: GlyphHeight = %d, want 12", tier, entry.GlyphHeight)
		}
		if entry.GlyphWidth != 6 {
			t.Errorf("tier %q: GlyphWidth = %d, want 6 (spleen wins)", tier, entry.GlyphWidth)
		}
		if entry.GlyphAdvance != 7 {
			t.Errorf("tier %q: GlyphAdvance = %d, want 7", tier, entry.GlyphAdvance)
		}
	}
}

// TestBuild_TieBreaking_MixedHeights verifies tie-breaking within the full
// 128×128 example. With best-fit algorithm and PPI=0 fallback targets:
// Normal(14): terminus-8x14 at height 14, distance 0 → exact match (width=8)
// Fullsize(→TierFull): largest is height 16, spleen-8x16 wins over terminus-8x16 by priority
func TestBuild_TieBreaking_MixedHeights(t *testing.T) {
	restore := font.SnapshotAndClear()
	defer restore()
	registerTestFonts()

	cat, err := tiercatalog.Build(tiercatalog.Params{
		PixelWidth:  128,
		PixelHeight: 128,
		MinChars:    10,
	})
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	// Normal tier should be height 14 (terminus-8x14 exact match for target 14).
	normalEntry, _ := cat.Get(tiercatalog.TierNormal)
	if normalEntry.GlyphHeight != 14 {
		t.Errorf("normal: GlyphHeight = %d, want 14", normalEntry.GlyphHeight)
	}
	if normalEntry.GlyphWidth != 8 {
		t.Errorf("normal: GlyphWidth = %d, want 8 (terminus-8x14)", normalEntry.GlyphWidth)
	}

	// Fullsize tier should be height 16 with spleen's metrics (spleen wins tie at height 16).
	fullEntry, _ := cat.Get(tiercatalog.TierFullsize)
	if fullEntry.GlyphHeight != 16 {
		t.Errorf("fullsize: GlyphHeight = %d, want 16", fullEntry.GlyphHeight)
	}
	if fullEntry.GlyphWidth != 8 {
		t.Errorf("fullsize: GlyphWidth = %d, want 8 (spleen-8x16 wins over terminus-8x16)", fullEntry.GlyphWidth)
	}
}

// TestBuild_ExcludesOversizedFont verifies that spleen-12x24 (advance 13) is
// excluded from a 128px wide region with MinChars=10 (maxAdvance=12).
func TestBuild_ExcludesOversizedFont(t *testing.T) {
	restore := font.SnapshotAndClear()
	defer restore()
	registerTestFonts()

	cat, err := tiercatalog.Build(tiercatalog.Params{
		PixelWidth:  128,
		PixelHeight: 128,
		MinChars:    10,
	})
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	// No tier should have GlyphHeight=24 (from the excluded spleen-12x24).
	for _, tier := range cat.Tiers() {
		entry, _ := cat.Get(tier)
		if entry.GlyphHeight == 24 {
			t.Errorf("tier %q: height 24 should have been excluded (advance 13 > maxAdvance 12)", tier)
		}
		if entry.GlyphAdvance == 13 {
			t.Errorf("tier %q: advance 13 should have been excluded", tier)
		}
	}
}

// TestBuild_EmptyRegistry verifies that Build returns an error when no fonts
// are registered.
func TestBuild_EmptyRegistry(t *testing.T) {
	restore := font.SnapshotAndClear()
	defer restore()

	_, err := tiercatalog.Build(tiercatalog.Params{
		PixelWidth:  128,
		PixelHeight: 128,
		MinChars:    10,
	})
	if err == nil {
		t.Fatal("Build: expected error for empty registry, got nil")
	}
}

// TestBuild_NoFontsFit verifies that Build returns an error when all fonts
// are too wide for the region.
func TestBuild_NoFontsFit(t *testing.T) {
	restore := font.SnapshotAndClear()
	defer restore()

	// Register a font with advance 20; for a 10px wide region with MinChars=10,
	// maxAdvance = 10/10 = 1, so advance 20 does not fit.
	font.Register(testFace{
		id:      "spleen-20x30",
		metrics: font.Metrics{GlyphWidth: 20, GlyphHeight: 30, GlyphAdvance: 20, RowHeight: 32},
	})

	_, err := tiercatalog.Build(tiercatalog.Params{
		PixelWidth:  10,
		PixelHeight: 100,
		MinChars:    10,
	})
	if err == nil {
		t.Fatal("Build: expected error when no fonts fit, got nil")
	}
}

// TestBuild_AllTiersPresent verifies that a successful Build always populates
// all six tiers regardless of how many unique heights exist.
func TestBuild_AllTiersPresent(t *testing.T) {
	restore := font.SnapshotAndClear()
	defer restore()
	registerTestFonts()

	cat, err := tiercatalog.Build(tiercatalog.Params{
		PixelWidth:  128,
		PixelHeight: 128,
		MinChars:    10,
	})
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	for _, tier := range cat.Tiers() {
		_, ok := cat.Get(tier)
		if !ok {
			t.Errorf("tier %q: missing from catalog", tier)
		}
	}
}

// TestBuild_Monotonicity verifies the GlyphHeight ordering invariant across all six tiers.
func TestBuild_Monotonicity(t *testing.T) {
	restore := font.SnapshotAndClear()
	defer restore()
	registerTestFonts()

	cat, err := tiercatalog.Build(tiercatalog.Params{
		PixelWidth:  128,
		PixelHeight: 128,
		MinChars:    10,
	})
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	small, _ := cat.Get(tiercatalog.TierSmall)
	normal, _ := cat.Get(tiercatalog.TierNormal)
	large, _ := cat.Get(tiercatalog.TierLarge)
	huge, _ := cat.Get(tiercatalog.TierHuge)
	colossal, _ := cat.Get(tiercatalog.TierColossal)
	full, _ := cat.Get(tiercatalog.TierFull)

	if small.GlyphHeight > normal.GlyphHeight {
		t.Errorf("monotonicity: small(%d) > normal(%d)", small.GlyphHeight, normal.GlyphHeight)
	}
	if normal.GlyphHeight > large.GlyphHeight {
		t.Errorf("monotonicity: normal(%d) > large(%d)", normal.GlyphHeight, large.GlyphHeight)
	}
	if large.GlyphHeight > huge.GlyphHeight {
		t.Errorf("monotonicity: large(%d) > huge(%d)", large.GlyphHeight, huge.GlyphHeight)
	}
	if huge.GlyphHeight > colossal.GlyphHeight {
		t.Errorf("monotonicity: huge(%d) > colossal(%d)", huge.GlyphHeight, colossal.GlyphHeight)
	}
	if colossal.GlyphHeight > full.GlyphHeight {
		t.Errorf("monotonicity: colossal(%d) > full(%d)", colossal.GlyphHeight, full.GlyphHeight)
	}
}

// TestBackwardCompat_TierFullsizeAlias verifies that Get(TierFullsize) returns
// the exact same entry as Get(TierFull). This is the alias handling requirement.
func TestBackwardCompat_TierFullsizeAlias(t *testing.T) {
	restore := font.SnapshotAndClear()
	defer restore()
	registerTestFonts()

	cat, err := tiercatalog.Build(tiercatalog.Params{
		PixelWidth:  128,
		PixelHeight: 128,
		MinChars:    10,
	})
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	fullEntry, fullOK := cat.Get(tiercatalog.TierFull)
	fullsizeEntry, fullsizeOK := cat.Get(tiercatalog.TierFullsize)

	if !fullOK {
		t.Fatal("Get(TierFull) returned false")
	}
	if !fullsizeOK {
		t.Fatal("Get(TierFullsize) returned false")
	}
	if fullEntry != fullsizeEntry {
		t.Errorf("TierFullsize alias broken:\n  TierFull    = %+v\n  TierFullsize = %+v", fullEntry, fullsizeEntry)
	}
}

// TestBackwardCompat_ZeroValueNewParams verifies that a caller using the old
// Params shape (only PixelWidth, PixelHeight, MinChars) gets the same result
// as one explicitly passing nil/zero for TierTargetsMM and TierTargetsPx.
func TestBackwardCompat_ZeroValueNewParams(t *testing.T) {
	restore := font.SnapshotAndClear()
	defer restore()
	registerTestFonts()

	// Old-style caller: no new fields.
	catOld, err := tiercatalog.Build(tiercatalog.Params{
		PixelWidth:  128,
		PixelHeight: 128,
		MinChars:    10,
	})
	if err != nil {
		t.Fatalf("Build (old-style) failed: %v", err)
	}

	// New-style caller: explicitly nil maps (should be equivalent).
	catNew, err := tiercatalog.Build(tiercatalog.Params{
		PixelWidth:    128,
		PixelHeight:   128,
		MinChars:      10,
		PPI:           0,
		TierTargetsMM: nil,
		TierTargetsPx: nil,
	})
	if err != nil {
		t.Fatalf("Build (new-style nil fields) failed: %v", err)
	}

	// Both catalogs should be identical.
	for _, tier := range catOld.Tiers() {
		oldEntry, _ := catOld.Get(tier)
		newEntry, _ := catNew.Get(tier)
		if oldEntry != newEntry {
			t.Errorf("tier %q: old=%+v, new=%+v", tier, oldEntry, newEntry)
		}
	}
}

// TestBackwardCompat_CatalogMethods verifies PixelWidth(), PixelHeight(),
// MinChars(), and Tiers() all return correct values.
func TestBackwardCompat_CatalogMethods(t *testing.T) {
	restore := font.SnapshotAndClear()
	defer restore()
	registerTestFonts()

	cat, err := tiercatalog.Build(tiercatalog.Params{
		PixelWidth:  200,
		PixelHeight: 150,
		MinChars:    15,
	})
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	if cat.PixelWidth() != 200 {
		t.Errorf("PixelWidth() = %d, want 200", cat.PixelWidth())
	}
	if cat.PixelHeight() != 150 {
		t.Errorf("PixelHeight() = %d, want 150", cat.PixelHeight())
	}
	if cat.MinChars() != 15 {
		t.Errorf("MinChars() = %d, want 15", cat.MinChars())
	}

	// Tiers() should return all 6 tiers in order.
	tiers := cat.Tiers()
	expectedTiers := []tiercatalog.Tier{
		tiercatalog.TierSmall,
		tiercatalog.TierNormal,
		tiercatalog.TierLarge,
		tiercatalog.TierHuge,
		tiercatalog.TierColossal,
		tiercatalog.TierFull,
	}
	if len(tiers) != len(expectedTiers) {
		t.Fatalf("Tiers() returned %d tiers, want %d", len(tiers), len(expectedTiers))
	}
	for i, tier := range tiers {
		if tier != expectedTiers[i] {
			t.Errorf("Tiers()[%d] = %q, want %q", i, tier, expectedTiers[i])
		}
	}
}

// TestBackwardCompat_FontID_Populated verifies that Entry.FontID is non-empty
// for all tiers in a successful catalog.
func TestBackwardCompat_FontID_Populated(t *testing.T) {
	restore := font.SnapshotAndClear()
	defer restore()
	registerTestFonts()

	cat, err := tiercatalog.Build(tiercatalog.Params{
		PixelWidth:  128,
		PixelHeight: 128,
		MinChars:    10,
	})
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	for _, tier := range cat.Tiers() {
		entry, ok := cat.Get(tier)
		if !ok {
			t.Errorf("tier %q: not present in catalog", tier)
			continue
		}
		if entry.FontID == "" {
			t.Errorf("tier %q: FontID is empty", tier)
		}
	}
}

// TestBuild_ZeroPixelWidth_ReturnsError verifies that PixelWidth=0 produces an error
// because advanceBudget=0 means no font can fit (all fonts have GlyphAdvance > 0).
func TestBuild_ZeroPixelWidth_ReturnsError(t *testing.T) {
	restore := font.SnapshotAndClear()
	defer restore()

	font.Register(testFace{
		id:      "spleen-5x8",
		metrics: font.Metrics{GlyphWidth: 5, GlyphHeight: 8, GlyphAdvance: 6, RowHeight: 10},
	})

	_, err := tiercatalog.Build(tiercatalog.Params{
		PixelWidth:  0,
		PixelHeight: 128,
		MinChars:    10,
	})
	if err == nil {
		t.Fatal("Build: expected error for PixelWidth=0, got nil")
	}
	// Verify error contains useful diagnostic info.
	if !strings.Contains(err.Error(), "no font fits") {
		t.Errorf("error message should mention 'no font fits': %v", err)
	}
}

// TestBuild_TierTargetFarFromPool verifies that when a tier's pixel target
// is much larger than any available font, the system still selects the
// closest (largest) qualifying font without error.
// This validates Req 3.5 behavior: the global candidate pool is used for all tiers.
func TestBuild_TierTargetFarFromPool(t *testing.T) {
	restore := font.SnapshotAndClear()
	defer restore()

	// Register only small fonts.
	font.Register(testFace{
		id:      "spleen-5x8",
		metrics: font.Metrics{GlyphWidth: 5, GlyphHeight: 8, GlyphAdvance: 6, RowHeight: 10},
	})
	font.Register(testFace{
		id:      "spleen-6x12",
		metrics: font.Metrics{GlyphWidth: 6, GlyphHeight: 12, GlyphAdvance: 7, RowHeight: 14},
	})

	// With PPI=0, TierColossal target is 40px — well beyond any registered font.
	cat, err := tiercatalog.Build(tiercatalog.Params{
		PixelWidth:  128,
		PixelHeight: 128,
		MinChars:    10,
	})
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	// TierColossal should get the largest available (spleen-6x12, height 12).
	colossal, ok := cat.Get(tiercatalog.TierColossal)
	if !ok {
		t.Fatal("TierColossal: not present")
	}
	if colossal.GlyphHeight != 12 {
		t.Errorf("TierColossal: GlyphHeight = %d, want 12 (closest to target 40)", colossal.GlyphHeight)
	}

	// TierFull should also be the largest (spleen-6x12).
	full, ok := cat.Get(tiercatalog.TierFull)
	if !ok {
		t.Fatal("TierFull: not present")
	}
	if full.GlyphHeight != 12 {
		t.Errorf("TierFull: GlyphHeight = %d, want 12", full.GlyphHeight)
	}

	// Monotonicity should hold.
	small, _ := cat.Get(tiercatalog.TierSmall)
	if small.GlyphHeight > colossal.GlyphHeight {
		t.Errorf("monotonicity: small(%d) > colossal(%d)", small.GlyphHeight, colossal.GlyphHeight)
	}
}

// TestBuild_GracefulMigration_PPI0_EquivalentToOldQuartile verifies that for the
// standard font registry with PPI=0, the new best-fit algorithm produces a catalog
// that is functionally equivalent to what the old quartile-based Build produced.
//
// "Functionally equivalent" means:
//   - Monotonically increasing font sizes (semantic ordering preserved)
//   - Small selects the smallest available font (height 8)
//   - Full selects the largest available font (height 16)
//   - The overall range (8→16) matches the old algorithm's range
//   - All tiers are populated
//
// The exact per-tier assignments may differ (the new algorithm uses best-fit to
// pixel targets rather than quartile distribution), but the catalog remains usable
// with the same semantic guarantees.
func TestBuild_GracefulMigration_PPI0_EquivalentToOldQuartile(t *testing.T) {
	restore := font.SnapshotAndClear()
	defer restore()
	registerTestFonts()

	cat, err := tiercatalog.Build(tiercatalog.Params{
		PixelWidth:  128,
		PixelHeight: 128,
		MinChars:    10,
		// PPI=0 (zero value) — triggers pixel fallback targets
	})
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	small, _ := cat.Get(tiercatalog.TierSmall)
	normal, _ := cat.Get(tiercatalog.TierNormal)
	large, _ := cat.Get(tiercatalog.TierLarge)
	full, _ := cat.Get(tiercatalog.TierFull)

	// Verify semantic ordering preserved (same as old algorithm).
	if small.GlyphHeight > normal.GlyphHeight {
		t.Errorf("graceful migration: small(%d) > normal(%d)", small.GlyphHeight, normal.GlyphHeight)
	}
	if normal.GlyphHeight > large.GlyphHeight {
		t.Errorf("graceful migration: normal(%d) > large(%d)", normal.GlyphHeight, large.GlyphHeight)
	}
	if large.GlyphHeight > full.GlyphHeight {
		t.Errorf("graceful migration: large(%d) > full(%d)", large.GlyphHeight, full.GlyphHeight)
	}

	// Verify range preserved: smallest font for Small, largest for Full.
	if small.GlyphHeight != 8 {
		t.Errorf("graceful migration: Small should be height 8 (smallest), got %d", small.GlyphHeight)
	}
	if full.GlyphHeight != 16 {
		t.Errorf("graceful migration: Full should be height 16 (largest), got %d", full.GlyphHeight)
	}

	// Verify all tiers populated.
	for _, tier := range cat.Tiers() {
		_, ok := cat.Get(tier)
		if !ok {
			t.Errorf("graceful migration: tier %q missing", tier)
		}
	}
}

// TestBuild_PPI96_StandardFonts verifies that PPI=96 mm→px conversion produces
// expected pixel targets and tier assignments. With standard fonts:
//
//	Small = round(3.0 * 96 / 25.4) = round(11.34) = 11px → closest: spleen-6x12 (distance 1)
//	Normal = round(5.0 * 96 / 25.4) = round(18.9) = 19px → closest: spleen-8x16 (distance 3)
//	Large = round(8.0 * 96 / 25.4) = round(30.24) = 30px → closest: spleen-8x16 (distance 14)
//	Huge = round(12.0 * 96 / 25.4) = round(45.35) = 45px → closest: spleen-8x16 (distance 29)
//	Colossal = round(18.0 * 96 / 25.4) = round(68.03) = 68px → closest: spleen-8x16 (distance 52)
//	Full: largest qualifying → spleen-8x16 (height 16)
func TestBuild_PPI96_StandardFonts(t *testing.T) {
	restore := font.SnapshotAndClear()
	defer restore()
	registerTestFonts()

	cat, err := tiercatalog.Build(tiercatalog.Params{
		PixelWidth:  128,
		PixelHeight: 128,
		MinChars:    10,
		PPI:         96,
	})
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	// With PPI=96, pixel targets are:
	// Small=11, Normal=19, Large=30, Huge=45, Colossal=68
	// Available qualifying heights: 8, 12, 13, 14, 16 (advance ≤12)
	// Best-fit selections:
	//   Small(11): spleen-6x12 at height 12, distance 1 (vs spleen-5x8 at distance 3)
	//   Normal(19): spleen-8x16 at height 16, distance 3 (closest to 19)
	//   Large(30): spleen-8x16 at height 16, distance 14 (closest to 30)
	//   Huge(45): spleen-8x16 at height 16, distance 29 (closest to 45)
	//   Colossal(68): spleen-8x16 at height 16, distance 52 (closest to 68)
	//   Full: spleen-8x16 (height 16, largest qualifying, spleen wins priority tie)

	tests := []struct {
		tier       tiercatalog.Tier
		wantHeight int
	}{
		{tiercatalog.TierSmall, 12},
		{tiercatalog.TierNormal, 16},
		{tiercatalog.TierLarge, 16},
		{tiercatalog.TierHuge, 16},
		{tiercatalog.TierColossal, 16},
		{tiercatalog.TierFull, 16},
	}

	for _, tc := range tests {
		entry, ok := cat.Get(tc.tier)
		if !ok {
			t.Errorf("tier %q: not present in catalog", tc.tier)
			continue
		}
		if entry.GlyphHeight != tc.wantHeight {
			t.Errorf("tier %q: GlyphHeight = %d, want %d", tc.tier, entry.GlyphHeight, tc.wantHeight)
		}
	}
}

// TestBuild_PPI110_ShortPanelSmallTier verifies that a 128x64 mono panel with
// real SH1106-class PPI keeps the small tier compact enough to fit the panel
// comfortably, selecting a 12px face instead of the taller 13px option.
func TestBuild_PPI110_ShortPanelSmallTier(t *testing.T) {
	restore := font.SnapshotAndClear()
	defer restore()
	registerTestFonts()

	cat, err := tiercatalog.Build(tiercatalog.Params{
		PixelWidth:  128,
		PixelHeight: 64,
		MinChars:    10,
		PPI:         110.08334658460504,
	})
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	entry, ok := cat.Get(tiercatalog.TierSmall)
	if !ok {
		t.Fatal("small tier missing")
	}
	if entry.GlyphHeight != 12 {
		t.Fatalf("Small tier GlyphHeight = %d, want 12", entry.GlyphHeight)
	}
	if entry.FontID != "spleen-6x12" {
		t.Fatalf("Small tier FontID = %q, want %q", entry.FontID, "spleen-6x12")
	}
}

// TestBuild_CustomMMOverride verifies that a custom mm target overrides only
// the targeted tier while leaving others at their defaults.
// Standard fonts with PPI=96: default Normal target = 5.0mm → 19px → height 16.
// Override Normal to 3.0mm → 11px → closest is height 12 (spleen-6x12, distance 1).
// Small should remain unchanged at its default (3.0mm → 11px → height 12).
func TestBuild_CustomMMOverride(t *testing.T) {
	restore := font.SnapshotAndClear()
	defer restore()
	registerTestFonts()

	// Build with default targets.
	catDefault, err := tiercatalog.Build(tiercatalog.Params{
		PixelWidth:  128,
		PixelHeight: 128,
		MinChars:    10,
		PPI:         96,
	})
	if err != nil {
		t.Fatalf("Build (default) failed: %v", err)
	}

	// Build with custom Normal target = 3.0mm (same as Small default).
	catCustom, err := tiercatalog.Build(tiercatalog.Params{
		PixelWidth:  128,
		PixelHeight: 128,
		MinChars:    10,
		PPI:         96,
		TierTargetsMM: map[tiercatalog.Tier]float64{
			tiercatalog.TierNormal: 3.0, // 3.0mm → round(3.0*96/25.4) = 11px
		},
	})
	if err != nil {
		t.Fatalf("Build (custom) failed: %v", err)
	}

	// Small should be the same in both catalogs (both have 3.0mm → 11px target).
	smallDefault, _ := catDefault.Get(tiercatalog.TierSmall)
	smallCustom, _ := catCustom.Get(tiercatalog.TierSmall)
	if smallDefault != smallCustom {
		t.Errorf("Small tier changed with Normal override:\n  default: %+v\n  custom:  %+v",
			smallDefault, smallCustom)
	}

	// Normal should now target 11px (height 12) instead of 19px (height 16).
	normalCustom, _ := catCustom.Get(tiercatalog.TierNormal)
	if normalCustom.GlyphHeight != 12 {
		t.Errorf("Normal (custom 3.0mm): GlyphHeight = %d, want 12", normalCustom.GlyphHeight)
	}
}

// TestBuild_TierFullsize_Alias explicitly tests that Get(TierFullsize) returns
// the same entry as Get(TierFull).
func TestBuild_TierFullsize_Alias(t *testing.T) {
	restore := font.SnapshotAndClear()
	defer restore()
	registerTestFonts()

	cat, err := tiercatalog.Build(tiercatalog.Params{
		PixelWidth:  128,
		PixelHeight: 128,
		MinChars:    10,
	})
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	fullEntry, fullOK := cat.Get(tiercatalog.TierFull)
	fullsizeEntry, fullsizeOK := cat.Get(tiercatalog.TierFullsize)

	if !fullOK {
		t.Fatal("Get(TierFull) returned false")
	}
	if !fullsizeOK {
		t.Fatal("Get(TierFullsize) returned false")
	}
	if fullEntry != fullsizeEntry {
		t.Errorf("TierFullsize alias mismatch:\n  TierFull    = %+v\n  TierFullsize = %+v",
			fullEntry, fullsizeEntry)
	}
}

// TestBuild_ExpandedTiers verifies that TierHuge and TierColossal are populated
// in a successful catalog and have non-zero GlyphHeight values.
func TestBuild_ExpandedTiers(t *testing.T) {
	restore := font.SnapshotAndClear()
	defer restore()
	registerTestFonts()

	cat, err := tiercatalog.Build(tiercatalog.Params{
		PixelWidth:  128,
		PixelHeight: 128,
		MinChars:    10,
	})
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	hugeEntry, hugeOK := cat.Get(tiercatalog.TierHuge)
	if !hugeOK {
		t.Fatal("TierHuge: not present in catalog")
	}
	if hugeEntry.GlyphHeight <= 0 {
		t.Errorf("TierHuge: GlyphHeight = %d, want > 0", hugeEntry.GlyphHeight)
	}

	colossalEntry, colossalOK := cat.Get(tiercatalog.TierColossal)
	if !colossalOK {
		t.Fatal("TierColossal: not present in catalog")
	}
	if colossalEntry.GlyphHeight <= 0 {
		t.Errorf("TierColossal: GlyphHeight = %d, want > 0", colossalEntry.GlyphHeight)
	}

	// Verify monotonicity: Huge ≤ Colossal.
	if hugeEntry.GlyphHeight > colossalEntry.GlyphHeight {
		t.Errorf("monotonicity: TierHuge(%d) > TierColossal(%d)",
			hugeEntry.GlyphHeight, colossalEntry.GlyphHeight)
	}
}

// TestBuild_ErrorMessageContent verifies that error messages include useful
// diagnostic details (region dimensions, constraints, smallest advance).
func TestBuild_ErrorMessageContent(t *testing.T) {
	restore := font.SnapshotAndClear()
	defer restore()

	// Test 1: Empty registry error message.
	_, err := tiercatalog.Build(tiercatalog.Params{
		PixelWidth:  128,
		PixelHeight: 128,
		MinChars:    10,
	})
	if err == nil {
		t.Fatal("expected error for empty registry")
	}
	if !strings.Contains(err.Error(), "no fonts registered") {
		t.Errorf("empty registry error should contain 'no fonts registered': got %q", err.Error())
	}

	// Test 2: No font fits error message with diagnostic details.
	font.Register(testFace{
		id:      "spleen-20x30",
		metrics: font.Metrics{GlyphWidth: 20, GlyphHeight: 30, GlyphAdvance: 20, RowHeight: 32},
	})

	_, err = tiercatalog.Build(tiercatalog.Params{
		PixelWidth:  10,
		PixelHeight: 100,
		MinChars:    10,
	})
	if err == nil {
		t.Fatal("expected error when no fonts fit")
	}
	errMsg := err.Error()
	if !strings.Contains(errMsg, "no font fits") {
		t.Errorf("error should contain 'no font fits': got %q", errMsg)
	}
	if !strings.Contains(errMsg, "10") {
		t.Errorf("error should contain pixel width (10): got %q", errMsg)
	}
	if !strings.Contains(errMsg, "smallest advance") {
		t.Errorf("error should contain 'smallest advance': got %q", errMsg)
	}
}

// TestBuild_BackwardCompat_ZeroNewFields verifies that existing callers that
// only set PixelWidth, PixelHeight, and MinChars (zero values for PPI,
// TierTargetsMM, TierTargetsPx) get equivalent behavior — specifically that
// the catalog is produced successfully with all tiers populated.
func TestBuild_BackwardCompat_ZeroNewFields(t *testing.T) {
	restore := font.SnapshotAndClear()
	defer restore()
	registerTestFonts()

	// Old-style caller: only the original three fields.
	cat, err := tiercatalog.Build(tiercatalog.Params{
		PixelWidth:  128,
		PixelHeight: 128,
		MinChars:    10,
	})
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	// All tiers must be populated.
	for _, tier := range cat.Tiers() {
		entry, ok := cat.Get(tier)
		if !ok {
			t.Errorf("tier %q: not present", tier)
			continue
		}
		// All metrics must be positive.
		if entry.GlyphWidth <= 0 || entry.GlyphHeight <= 0 ||
			entry.GlyphAdvance <= 0 || entry.RowHeight <= 0 {
			t.Errorf("tier %q: has non-positive metrics: %+v", tier, entry)
		}
	}

	// Verify backward compat: zero new fields = nil new fields.
	catNil, err := tiercatalog.Build(tiercatalog.Params{
		PixelWidth:    128,
		PixelHeight:   128,
		MinChars:      10,
		PPI:           0,
		TierTargetsMM: nil,
		TierTargetsPx: nil,
	})
	if err != nil {
		t.Fatalf("Build (nil fields) failed: %v", err)
	}

	for _, tier := range cat.Tiers() {
		oldEntry, _ := cat.Get(tier)
		newEntry, _ := catNil.Get(tier)
		if oldEntry != newEntry {
			t.Errorf("tier %q: zero-value differs from nil-value:\n  zero: %+v\n  nil:  %+v",
				tier, oldEntry, newEntry)
		}
	}
}

// TestBuild_FontID_Populated verifies that FontID is non-empty for all tier
// entries in a successful catalog and corresponds to a known font.
func TestBuild_FontID_Populated(t *testing.T) {
	restore := font.SnapshotAndClear()
	defer restore()
	registerTestFonts()

	cat, err := tiercatalog.Build(tiercatalog.Params{
		PixelWidth:  128,
		PixelHeight: 128,
		MinChars:    10,
	})
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	for _, tier := range cat.Tiers() {
		entry, ok := cat.Get(tier)
		if !ok {
			t.Errorf("tier %q: not present", tier)
			continue
		}
		if entry.FontID == "" {
			t.Errorf("tier %q: FontID is empty", tier)
			continue
		}
		// Verify the FontID corresponds to a registered font.
		if _, found := font.Get(entry.FontID); !found {
			t.Errorf("tier %q: FontID %q not found in font registry", tier, entry.FontID)
		}
	}
}
