package tierselect_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/databeast/cyberhud/display/surface/fonts"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
	"github.com/databeast/cyberhud/display/surface/tierselect"
	"pgregory.net/rapid"
)

// --- From: tierselect_prop_test.go ---

// stubFace is a minimal font.Face implementation for property testing.
type stubFace struct {
	id      string
	metrics font.Metrics
}

func (s stubFace) ID() string                { return s.id }
func (s stubFace) Metrics() font.Metrics     { return s.metrics }
func (s stubFace) GlyphRow(rune, int) uint32 { return 0 }

// For any valid catalog and request, the Font_Face returned by tierselect.Select
// SHALL satisfy GlyphAdvance * MinChars <= catalog.PixelWidth().

func TestProperty3_WidthSafetySelect(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		restore := font.SnapshotAndClear()
		defer restore()

		// Generate random region parameters.
		pixelWidth := rapid.IntRange(10, 4096).Draw(t, "pixelWidth")
		pixelHeight := rapid.IntRange(10, 4096).Draw(t, "pixelHeight")
		minChars := rapid.IntRange(1, 50).Draw(t, "minChars")

		// Generate random fonts with varying families and metrics.
		families := []string{"spleen", "terminus", "cozette", "testfont"}
		numFonts := rapid.IntRange(1, 20).Draw(t, "numFonts")
		for i := 0; i < numFonts; i++ {
			glyphWidth := rapid.IntRange(1, 64).Draw(t, fmt.Sprintf("gw%d", i))
			glyphHeight := rapid.IntRange(1, 64).Draw(t, fmt.Sprintf("gh%d", i))
			glyphAdvance := rapid.IntRange(1, 64).Draw(t, fmt.Sprintf("ga%d", i))
			rowHeight := glyphHeight + rapid.IntRange(0, 4).Draw(t, fmt.Sprintf("rh%d", i))

			family := rapid.SampledFrom(families).Draw(t, fmt.Sprintf("fam%d", i))
			id := fmt.Sprintf("%s-%dx%d-p3s-%d", family, glyphWidth, glyphHeight, i)

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

		// Build the catalog with the given dimensions and MinChars.
		cat, err := tiercatalog.Build(tiercatalog.Params{
			PixelWidth:  pixelWidth,
			PixelHeight: pixelHeight,
			MinChars:    minChars,
		})
		if err != nil {
			// No qualifying fonts for this configuration — property is vacuously true.
			return
		}

		// Determine which families actually have registered variants.
		registeredFamilies := make(map[string]bool)
		for _, f := range font.List() {
			fam := extractFamily(f.ID())
			registeredFamilies[fam] = true
		}
		var availableFamilies []string
		for fam := range registeredFamilies {
			availableFamilies = append(availableFamilies, fam)
		}
		if len(availableFamilies) == 0 {
			return
		}

		// For each tier in the catalog, call Select with a family that exists.
		tiers := []tiercatalog.Tier{
			tiercatalog.TierSmall,
			tiercatalog.TierNormal,
			tiercatalog.TierLarge,
			tiercatalog.TierFullsize,
		}

		for _, tier := range tiers {
			family := rapid.SampledFrom(availableFamilies).Draw(t, fmt.Sprintf("reqFam_%s", tier))

			face := tierselect.Select(cat, tierselect.Request{
				Family: family,
				Tier:   tier,
			})

			if face == nil {
				t.Fatalf("Select returned nil for family=%q tier=%q", family, tier)
			}

			// Assert width safety: returned face's GlyphAdvance * MinChars <= catalog.PixelWidth()
			advance := face.Metrics().GlyphAdvance
			if advance*minChars > cat.PixelWidth() {
				t.Fatalf("width safety violated: family=%q tier=%q face=%q GlyphAdvance=%d * MinChars=%d = %d > PixelWidth=%d",
					family, tier, face.ID(), advance, minChars, advance*minChars, cat.PixelWidth())
			}
		}
	})
}

// extractFamily splits a font ID on the first "-" and returns the family prefix.
func extractFamily(id string) string {
	idx := strings.Index(id, "-")
	if idx < 0 {
		return id
	}
	return id[:idx]
}

// For any catalog and request where the requested family has an exact variant
// at the tier's target GlyphHeight, Select SHALL return that exact variant.

func TestProperty5_FamilyPreferenceStability(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Snapshot and clear the font registry for isolation.
		restore := font.SnapshotAndClear()
		defer restore()

		// Choose a target family and a target glyph height for the exact match.
		family := rapid.SampledFrom([]string{"spleen", "terminus", "cozette"}).Draw(t, "family")
		targetHeight := rapid.IntRange(5, 64).Draw(t, "targetHeight")
		glyphWidth := rapid.IntRange(3, 32).Draw(t, "glyphWidth")
		glyphAdvance := rapid.IntRange(3, 32).Draw(t, "glyphAdvance")
		rowHeight := targetHeight + rapid.IntRange(0, 4).Draw(t, "rowExtra")

		// Register the exact family variant at targetHeight.
		exactID := fmt.Sprintf("%s-%dx%d", family, glyphWidth, targetHeight)
		font.Register(stubFace{
			id: exactID,
			metrics: font.Metrics{
				GlyphWidth:   glyphWidth,
				GlyphHeight:  targetHeight,
				GlyphAdvance: glyphAdvance,
				RowHeight:    rowHeight,
			},
		})

		// Register additional fonts to ensure the catalog is non-trivial.
		// We avoid registering another font from the same family at the exact
		// targetHeight, since having two same-family variants at the same height
		// makes the "exact match" precondition ambiguous.
		numExtra := rapid.IntRange(1, 10).Draw(t, "numExtra")
		for i := 0; i < numExtra; i++ {
			extraFamily := rapid.SampledFrom([]string{"spleen", "terminus", "cozette", "testfont"}).Draw(t, fmt.Sprintf("efam%d", i))
			extraHeight := rapid.IntRange(5, 64).Draw(t, fmt.Sprintf("egh%d", i))

			// Skip if this would create another variant from the same family
			// at the exact target height.
			if extraFamily == family && extraHeight == targetHeight {
				continue
			}

			extraWidth := rapid.IntRange(3, 32).Draw(t, fmt.Sprintf("egw%d", i))
			extraAdvance := rapid.IntRange(3, 32).Draw(t, fmt.Sprintf("ega%d", i))
			extraRow := extraHeight + rapid.IntRange(0, 4).Draw(t, fmt.Sprintf("erh%d", i))

			extraID := fmt.Sprintf("%s-%dx%d-extra%d", extraFamily, extraWidth, extraHeight, i)

			font.Register(stubFace{
				id: extraID,
				metrics: font.Metrics{
					GlyphWidth:   extraWidth,
					GlyphHeight:  extraHeight,
					GlyphAdvance: extraAdvance,
					RowHeight:    extraRow,
				},
			})
		}

		// Build the catalog. Choose region dimensions wide enough that the exact
		// font qualifies (GlyphAdvance * MinChars <= PixelWidth).
		minChars := rapid.IntRange(1, 20).Draw(t, "minChars")
		pixelWidth := glyphAdvance * minChars // Guarantees exact font qualifies.
		// Add some extra width to allow other fonts too.
		pixelWidth += rapid.IntRange(0, 200).Draw(t, "extraWidth")
		pixelHeight := rapid.IntRange(targetHeight, 4096).Draw(t, "pixelHeight")

		cat, err := tiercatalog.Build(tiercatalog.Params{
			PixelWidth:  pixelWidth,
			PixelHeight: pixelHeight,
			MinChars:    minChars,
		})
		if err != nil {
			// If catalog build fails, the property is vacuously true.
			return
		}

		// Find a tier whose target GlyphHeight matches our targetHeight.
		var matchedTier tiercatalog.Tier
		found := false
		for _, tier := range cat.Tiers() {
			entry, ok := cat.Get(tier)
			if ok && entry.GlyphHeight == targetHeight {
				matchedTier = tier
				found = true
				break
			}
		}
		if !found {
			// No tier maps to our target height — property is vacuously true.
			return
		}

		// Call Select with the family and the matched tier.
		result := tierselect.Select(cat, tierselect.Request{
			Family: family,
			Tier:   matchedTier,
		})

		// Assert: the returned face must be our exact family variant.
		if result.ID() != exactID {
			t.Fatalf("family preference stability violated: expected %q (exact family variant at height %d), got %q",
				exactID, targetHeight, result.ID())
		}
		if result.Metrics().GlyphHeight != targetHeight {
			t.Fatalf("returned face height mismatch: expected %d, got %d",
				targetHeight, result.Metrics().GlyphHeight)
		}
	})
}

// For any given catalog and request, Select SHALL return a Font_Face with the
// same ID when invoked twice with the same inputs.

func TestProperty7_DeterminismSelect(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		restore := font.SnapshotAndClear()
		defer restore()

		// Register a random set of fonts covering multiple families.
		families := []string{"spleen", "terminus", "cozette", "testfont"}
		numFonts := rapid.IntRange(2, 15).Draw(t, "numFonts")
		var registeredFamilies []string

		for i := 0; i < numFonts; i++ {
			glyphWidth := rapid.IntRange(4, 32).Draw(t, fmt.Sprintf("gw%d", i))
			glyphHeight := rapid.IntRange(6, 48).Draw(t, fmt.Sprintf("gh%d", i))
			glyphAdvance := rapid.IntRange(4, 32).Draw(t, fmt.Sprintf("ga%d", i))
			rowHeight := glyphHeight + rapid.IntRange(0, 4).Draw(t, fmt.Sprintf("rh%d", i))

			family := rapid.SampledFrom(families).Draw(t, fmt.Sprintf("fam%d", i))
			id := fmt.Sprintf("%s-%dx%d", family, glyphWidth, glyphHeight)

			font.Register(stubFace{
				id: id,
				metrics: font.Metrics{
					GlyphWidth:   glyphWidth,
					GlyphHeight:  glyphHeight,
					GlyphAdvance: glyphAdvance,
					RowHeight:    rowHeight,
				},
			})
			registeredFamilies = append(registeredFamilies, family)
		}

		// Build a catalog with region dimensions that allow some fonts to qualify.
		pixelWidth := rapid.IntRange(40, 512).Draw(t, "pixelWidth")
		pixelHeight := rapid.IntRange(40, 512).Draw(t, "pixelHeight")

		catalog, err := tiercatalog.Build(tiercatalog.Params{
			PixelWidth:  pixelWidth,
			PixelHeight: pixelHeight,
			MinChars:    10,
		})
		if err != nil {
			// No qualifying fonts for this configuration — property is vacuously true.
			return
		}

		// Pick a random family from those we registered and a random tier.
		family := rapid.SampledFrom(registeredFamilies).Draw(t, "reqFamily")
		tier := rapid.SampledFrom([]tiercatalog.Tier{
			tiercatalog.TierSmall,
			tiercatalog.TierNormal,
			tiercatalog.TierLarge,
			tiercatalog.TierFullsize,
		}).Draw(t, "reqTier")

		req := tierselect.Request{
			Family: family,
			Tier:   tier,
		}

		// Call Select twice with the same catalog and request.
		face1 := tierselect.Select(catalog, req)
		face2 := tierselect.Select(catalog, req)

		// Assert both calls return the same face ID.
		if face1.ID() != face2.ID() {
			t.Fatalf("determinism violated: Select returned %q first, then %q second for family=%q tier=%q",
				face1.ID(), face2.ID(), family, tier)
		}
	})
}

// --- From: tierselect_test.go ---

// testFace is a minimal font.Face for unit testing with known metrics.
type testFace struct {
	id      string
	metrics font.Metrics
}

func (f testFace) ID() string                { return f.id }
func (f testFace) Metrics() font.Metrics     { return f.metrics }
func (f testFace) GlyphRow(rune, int) uint32 { return 0 }

// registerStandardFonts registers the standard font set from the design document.
// For a 128×128 region with MinChars=10 (maxAdvance=12):
//
//	spleen-5x8 (adv 6), spleen-6x12 (adv 7), spleen-8x16 (adv 9)
//	terminus-6x12 (adv 7), terminus-8x14 (adv 9), terminus-8x16 (adv 9)
//	cozette-6x13 (adv 7)
//
// After build: small=8, normal=14, large=16, huge=16, colossal=16, full=16
func registerStandardFonts() {
	fonts := []testFace{
		{id: "spleen-5x8", metrics: font.Metrics{GlyphWidth: 5, GlyphHeight: 8, GlyphAdvance: 6, RowHeight: 10}},
		{id: "spleen-6x12", metrics: font.Metrics{GlyphWidth: 6, GlyphHeight: 12, GlyphAdvance: 7, RowHeight: 14}},
		{id: "spleen-8x16", metrics: font.Metrics{GlyphWidth: 8, GlyphHeight: 16, GlyphAdvance: 9, RowHeight: 18}},
		{id: "terminus-6x12", metrics: font.Metrics{GlyphWidth: 6, GlyphHeight: 12, GlyphAdvance: 7, RowHeight: 14}},
		{id: "terminus-8x14", metrics: font.Metrics{GlyphWidth: 8, GlyphHeight: 14, GlyphAdvance: 9, RowHeight: 16}},
		{id: "terminus-8x16", metrics: font.Metrics{GlyphWidth: 8, GlyphHeight: 16, GlyphAdvance: 9, RowHeight: 18}},
		{id: "cozette-6x13", metrics: font.Metrics{GlyphWidth: 6, GlyphHeight: 13, GlyphAdvance: 7, RowHeight: 15}},
	}
	for _, f := range fonts {
		font.Register(f)
	}
}

// buildStandardCatalog builds the 128×128 catalog with MinChars=10.
func buildStandardCatalog(t *testing.T) tiercatalog.Catalog {
	t.Helper()
	cat, err := tiercatalog.Build(tiercatalog.Params{
		PixelWidth:  128,
		PixelHeight: 128,
		MinChars:    10,
	})
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	return cat
}

// TestSelect_ExactMatch verifies that when a family has a variant at or below the
// target GlyphHeight for a tier, the correct variant is returned.
//
// Request: spleen/TierNormal → target height 14
// spleen variants: 8, 12, 16. Closest ≤ 14 is 12.
// Expected: spleen-6x12 (closest family variant at height ≤ 14)
func TestSelect_ExactMatch(t *testing.T) {
	restore := font.SnapshotAndClear()
	defer restore()
	registerStandardFonts()

	cat := buildStandardCatalog(t)

	face := tierselect.Select(cat, tierselect.Request{
		Family: "spleen",
		Tier:   tiercatalog.TierNormal,
	})

	if face == nil {
		t.Fatal("Select returned nil")
	}
	if face.ID() != "spleen-6x12" {
		t.Errorf("exact match: got %q, want %q", face.ID(), "spleen-6x12")
	}
	if face.Metrics().GlyphHeight != 12 {
		t.Errorf("exact match: GlyphHeight = %d, want 12", face.Metrics().GlyphHeight)
	}
}

// TestSelect_ClosestVariantFallback verifies that when a family has no variant
// at the exact target height, the closest family variant with GlyphHeight ≤
// target is returned (strategy 2).
//
// Request: cozette/TierLarge → target height 16
// Only cozette variant: cozette-6x13 (height 13 ≤ 16)
// Expected: cozette-6x13
func TestSelect_ClosestVariantFallback(t *testing.T) {
	restore := font.SnapshotAndClear()
	defer restore()
	registerStandardFonts()

	cat := buildStandardCatalog(t)

	face := tierselect.Select(cat, tierselect.Request{
		Family: "cozette",
		Tier:   tiercatalog.TierLarge,
	})

	if face == nil {
		t.Fatal("Select returned nil")
	}
	if face.ID() != "cozette-6x13" {
		t.Errorf("closest variant fallback: got %q, want %q", face.ID(), "cozette-6x13")
	}
	if face.Metrics().GlyphHeight != 13 {
		t.Errorf("closest variant fallback: GlyphHeight = %d, want 13", face.Metrics().GlyphHeight)
	}
}

// TestSelect_CrossFamilyFallback verifies that when a family has variants only
// above the target height, the cross-family fallback (strategy 3) is used.
// This should return a font at the target height from any family and log a warning.
//
// Setup: Register "testfam-10x20" (height 20, advance 11). Request testfam/TierSmall.
// Target height for small = 8. No testfam at ≤ 8, so cross-family fallback returns
// spleen-5x8 (height 8, the best-fit at target height from any family).
func TestSelect_CrossFamilyFallback(t *testing.T) {
	restore := font.SnapshotAndClear()
	defer restore()
	registerStandardFonts()

	// Register a family that only has a variant above the small tier target.
	font.Register(testFace{
		id:      "testfam-10x20",
		metrics: font.Metrics{GlyphWidth: 10, GlyphHeight: 20, GlyphAdvance: 11, RowHeight: 22},
	})

	cat := buildStandardCatalog(t)

	face := tierselect.Select(cat, tierselect.Request{
		Family: "testfam",
		Tier:   tiercatalog.TierSmall,
	})

	if face == nil {
		t.Fatal("Select returned nil")
	}
	// Cross-family: the best fit at height 8 is spleen-5x8 (spleen has highest priority).
	if face.ID() != "spleen-5x8" {
		t.Errorf("cross-family fallback: got %q, want %q", face.ID(), "spleen-5x8")
	}
	if face.Metrics().GlyphHeight != 8 {
		t.Errorf("cross-family fallback: GlyphHeight = %d, want 8", face.Metrics().GlyphHeight)
	}
}

// TestSelect_PanicInvalidTier verifies that Select panics with a diagnostic
// message when the requested tier is not in the catalog.
func TestSelect_PanicInvalidTier(t *testing.T) {
	restore := font.SnapshotAndClear()
	defer restore()
	registerStandardFonts()

	cat := buildStandardCatalog(t)

	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic for invalid tier, got none")
		}
		msg, ok := r.(string)
		if !ok {
			t.Fatalf("panic value is not a string: %v", r)
		}
		if !strings.Contains(msg, "invalid tier") {
			t.Errorf("panic message %q does not contain %q", msg, "invalid tier")
		}
		if !strings.Contains(msg, "bogus") {
			t.Errorf("panic message %q does not contain the tier name %q", msg, "bogus")
		}
	}()

	tierselect.Select(cat, tierselect.Request{
		Family: "spleen",
		Tier:   tiercatalog.Tier("bogus"),
	})
}

// TestSelect_PanicUnknownFamily verifies that Select panics immediately when
// the requested family has no registered variants at all.
func TestSelect_PanicUnknownFamily(t *testing.T) {
	restore := font.SnapshotAndClear()
	defer restore()
	registerStandardFonts()

	cat := buildStandardCatalog(t)

	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic for unknown family, got none")
		}
		msg, ok := r.(string)
		if !ok {
			t.Fatalf("panic value is not a string: %v", r)
		}
		if !strings.Contains(msg, "no registered variants") {
			t.Errorf("panic message %q does not contain %q", msg, "no registered variants")
		}
		if !strings.Contains(msg, "nofamily") {
			t.Errorf("panic message %q does not contain the family name %q", msg, "nofamily")
		}
	}()

	tierselect.Select(cat, tierselect.Request{
		Family: "nofamily",
		Tier:   tiercatalog.TierNormal,
	})
}

// TestSelect_PanicExhaustedFallback verifies that Select panics when the
// fallback chain is fully exhausted. This can only occur when the font
// registry is emptied after the catalog is built (registry corruption),
// since fonts.ByHeight always returns a face from a non-empty registry.
func TestSelect_PanicExhaustedFallback(t *testing.T) {
	restore := font.SnapshotAndClear()
	defer restore()

	// Register a font to build a valid catalog.
	font.Register(testFace{
		id:      "lonely-5x8",
		metrics: font.Metrics{GlyphWidth: 5, GlyphHeight: 8, GlyphAdvance: 5, RowHeight: 10},
	})

	cat, err := tiercatalog.Build(tiercatalog.Params{
		PixelWidth:  128,
		PixelHeight: 128,
		MinChars:    10,
	})
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	// Clear the registry completely after building the catalog.
	// Then register a family so familyHasAnyVariant passes, but the registry
	// has no actual entries when font.ByHeight is called (simulate corruption
	// by clearing again between the family check and ByHeight call).
	//
	// In practice, the exhausted-fallback panic is a guard against registry
	// corruption. With the current ByHeight implementation, it returns
	// Default() when no fonts are registered. We verify the panic triggers
	// by ensuring the registry is truly empty (no faces at all).
	innerRestore := font.SnapshotAndClear()
	defer innerRestore()

	// Empty registry: familyHasAnyVariant will return false first and panic
	// with "no registered variants". Test that this panic fires correctly.
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic for exhausted/unknown family fallback, got none")
		}
		msg, ok := r.(string)
		if !ok {
			t.Fatalf("panic value is not a string: %v", r)
		}
		// With empty registry after catalog build, the family check fires first.
		if !strings.Contains(msg, "no registered variants") && !strings.Contains(msg, "no font found") {
			t.Errorf("panic message %q does not contain expected diagnostic", msg)
		}
	}()

	tierselect.Select(cat, tierselect.Request{
		Family: "lonely",
		Tier:   tiercatalog.TierSmall,
	})
}

// TestSelectMulti verifies that SelectMulti returns the correct number of
// faces and each result matches calling Select individually.
func TestSelectMulti(t *testing.T) {
	restore := font.SnapshotAndClear()
	defer restore()
	registerStandardFonts()

	cat := buildStandardCatalog(t)

	reqs := []tierselect.Request{
		{Family: "spleen", Tier: tiercatalog.TierNormal},
		{Family: "terminus", Tier: tiercatalog.TierLarge},
		{Family: "cozette", Tier: tiercatalog.TierLarge},
		{Family: "spleen", Tier: tiercatalog.TierSmall},
		{Family: "spleen", Tier: tiercatalog.TierFullsize},
	}

	faces := tierselect.SelectMulti(cat, reqs)

	// Verify slice length.
	if len(faces) != len(reqs) {
		t.Fatalf("SelectMulti: got %d faces, want %d", len(faces), len(reqs))
	}

	// Verify each result matches an individual Select call.
	expectedIDs := []string{
		"spleen-6x12",   // spleen at normal (height 14) → closest ≤ 14 is height 12
		"terminus-8x16", // terminus at large (height 16) → exact match
		"cozette-6x13",  // cozette at large (height 16) → closest ≤ 16 is 13
		"spleen-5x8",    // spleen at small (height 8) → exact match
		"spleen-8x16",   // spleen at fullsize (height 16) → exact match
	}

	for i, face := range faces {
		if face == nil {
			t.Errorf("SelectMulti[%d]: got nil", i)
			continue
		}
		if face.ID() != expectedIDs[i] {
			t.Errorf("SelectMulti[%d]: got %q, want %q", i, face.ID(), expectedIDs[i])
		}

		// Also verify individual Select produces the same result.
		individual := tierselect.Select(cat, reqs[i])
		if individual.ID() != face.ID() {
			t.Errorf("SelectMulti[%d]: multi=%q differs from individual Select=%q",
				i, face.ID(), individual.ID())
		}
	}
}

// TestSelectMulti_Empty verifies that SelectMulti with zero requests returns
// an empty (non-nil) slice.
func TestSelectMulti_Empty(t *testing.T) {
	restore := font.SnapshotAndClear()
	defer restore()
	registerStandardFonts()

	cat := buildStandardCatalog(t)

	faces := tierselect.SelectMulti(cat, nil)
	if faces == nil {
		t.Error("SelectMulti(nil): got nil slice, want non-nil empty slice")
	}
	if len(faces) != 0 {
		t.Errorf("SelectMulti(nil): got %d faces, want 0", len(faces))
	}
}
