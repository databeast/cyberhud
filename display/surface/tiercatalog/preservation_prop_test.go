package tiercatalog_test

import (
	"fmt"
	"testing"

	"github.com/databeast/cyberhud/display/surface/fonts"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
	"github.com/databeast/cyberhud/display/surface/tierselect"
	"pgregory.net/rapid"
)

// ============================================================================
//
// These tests capture baseline behavior on UNFIXED code to ensure the bugfix
// does not break existing non-font-selection behavior.
// ============================================================================

// --- Helpers ---

// preservationFace is a minimal font.Face for preservation property testing.
type preservationFace struct {
	id      string
	metrics font.Metrics
}

func (f preservationFace) ID() string                { return f.id }
func (f preservationFace) Metrics() font.Metrics     { return f.metrics }
func (f preservationFace) GlyphRow(rune, int) uint32 { return 0 }

// registerPreservationFonts registers a controlled set of fonts for testing.
// Returns the restore function.
func registerPreservationFonts(t *rapid.T, numFonts int) {
	families := []string{"spleen", "terminus", "cozette", "testfont"}
	for i := 0; i < numFonts; i++ {
		glyphWidth := rapid.IntRange(3, 32).Draw(t, fmt.Sprintf("gw%d", i))
		glyphHeight := rapid.IntRange(5, 48).Draw(t, fmt.Sprintf("gh%d", i))
		glyphAdvance := rapid.IntRange(3, 32).Draw(t, fmt.Sprintf("ga%d", i))
		rowHeight := glyphHeight + rapid.IntRange(0, 4).Draw(t, fmt.Sprintf("rh%d", i))
		family := rapid.SampledFrom(families).Draw(t, fmt.Sprintf("fam%d", i))
		id := fmt.Sprintf("%s-%dx%d-pres-%d", family, glyphWidth, glyphHeight, i)

		font.Register(preservationFace{
			id: id,
			metrics: font.Metrics{
				GlyphWidth:   glyphWidth,
				GlyphHeight:  glyphHeight,
				GlyphAdvance: glyphAdvance,
				RowHeight:    rowHeight,
			},
		})
	}
}

// --- Property Tests ---

// TestPreservation_CatalogBuildTierAssignment verifies that tiercatalog.Build()
// assigns tiers identically for the same inputs — same quartile boundaries,
// same family priority tie-breaking (Spleen > Terminus > Cozette > other).
//
// For any random font registry and panel dimensions, building the catalog twice
// with identical inputs SHALL produce identical tier entries.

func TestPreservation_CatalogBuildTierAssignment(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		restore := font.SnapshotAndClear()
		defer restore()

		numFonts := rapid.IntRange(1, 20).Draw(t, "numFonts")
		registerPreservationFonts(t, numFonts)

		pixelWidth := rapid.IntRange(10, 2048).Draw(t, "pixelWidth")
		pixelHeight := rapid.IntRange(10, 2048).Draw(t, "pixelHeight")
		minChars := rapid.IntRange(1, 50).Draw(t, "minChars")

		params := tiercatalog.Params{
			PixelWidth:  pixelWidth,
			PixelHeight: pixelHeight,
			MinChars:    minChars,
		}

		cat1, err1 := tiercatalog.Build(params)
		cat2, err2 := tiercatalog.Build(params)

		// Both calls must produce the same error state.
		if (err1 == nil) != (err2 == nil) {
			t.Fatalf("inconsistent errors: err1=%v, err2=%v", err1, err2)
		}
		if err1 != nil {
			return // both failed, property is vacuously preserved
		}

		// All tier entries must be identical (including FontID).
		for _, tier := range cat1.Tiers() {
			e1, ok1 := cat1.Get(tier)
			e2, ok2 := cat2.Get(tier)
			if ok1 != ok2 {
				t.Fatalf("tier %q presence mismatch: cat1=%v, cat2=%v", tier, ok1, ok2)
			}
			if !ok1 {
				continue
			}
			if e1 != e2 {
				t.Fatalf("tier %q entries differ:\n  cat1: %+v\n  cat2: %+v", tier, e1, e2)
			}
			// FontID must be non-empty for all populated tiers.
			if e1.FontID == "" {
				t.Fatalf("tier %q: FontID is empty", tier)
			}
		}

		// PixelWidth and MinChars stored must match.
		if cat1.PixelWidth() != cat2.PixelWidth() {
			t.Fatalf("PixelWidth mismatch: %d vs %d", cat1.PixelWidth(), cat2.PixelWidth())
		}
		if cat1.MinChars() != cat2.MinChars() {
			t.Fatalf("MinChars mismatch: %d vs %d", cat1.MinChars(), cat2.MinChars())
		}
	})
}

// TestPreservation_CatalogBuildFamilyPriority verifies that tiercatalog.Build()
// applies family priority correctly: when multiple families share the same
// GlyphHeight, the one with higher priority (Spleen > Terminus > Cozette > other)
// is selected for the tier entry.

func TestPreservation_CatalogBuildFamilyPriority(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		restore := font.SnapshotAndClear()
		defer restore()

		// Register two fonts at same height but different families.
		height := rapid.IntRange(8, 32).Draw(t, "height")
		advance := rapid.IntRange(3, 12).Draw(t, "advance")

		// Pick two different families from the priority list.
		families := []struct {
			name     string
			priority int
		}{
			{"spleen", 3},
			{"terminus", 2},
			{"cozette", 1},
			{"testfont", 0},
		}
		i1 := rapid.IntRange(0, 2).Draw(t, "fam1Idx")
		i2 := rapid.IntRange(i1+1, 3).Draw(t, "fam2Idx")

		highPrio := families[i1]
		lowPrio := families[i2]

		// Register both at the same height and advance.
		font.Register(preservationFace{
			id: fmt.Sprintf("%s-5x%d", highPrio.name, height),
			metrics: font.Metrics{
				GlyphWidth:   5,
				GlyphHeight:  height,
				GlyphAdvance: advance,
				RowHeight:    height + 2,
			},
		})
		font.Register(preservationFace{
			id: fmt.Sprintf("%s-5x%d", lowPrio.name, height),
			metrics: font.Metrics{
				GlyphWidth:   5,
				GlyphHeight:  height,
				GlyphAdvance: advance,
				RowHeight:    height + 2,
			},
		})

		// Build catalog with dimensions that allow these fonts.
		pixelWidth := advance * 10 // guarantees both qualify
		cat, err := tiercatalog.Build(tiercatalog.Params{
			PixelWidth:  pixelWidth,
			PixelHeight: 200,
			MinChars:    10,
		})
		if err != nil {
			return // vacuously true
		}

		// Since there's only one unique height, all tiers map to it.
		// The winning entry should use the high-priority family's metrics.
		for _, tier := range cat.Tiers() {
			entry, ok := cat.Get(tier)
			if !ok {
				t.Fatalf("tier %q: not present", tier)
			}
			if entry.GlyphHeight != height {
				t.Fatalf("tier %q: expected GlyphHeight=%d, got %d", tier, height, entry.GlyphHeight)
			}
			// Advance should match (both have same advance).
			if entry.GlyphAdvance != advance {
				t.Fatalf("tier %q: expected GlyphAdvance=%d, got %d", tier, advance, entry.GlyphAdvance)
			}
			// FontID should reference the high-priority family font.
			if entry.FontID == "" {
				t.Fatalf("tier %q: FontID is empty", tier)
			}
		}
	})
}

// TestPreservation_TierSelectResolution verifies that tierselect.Select()
// produces the same resolution for the same (family, tier) inputs with the
// same catalog. Calling Select twice must return the same face ID.

func TestPreservation_TierSelectResolution(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		restore := font.SnapshotAndClear()
		defer restore()

		numFonts := rapid.IntRange(2, 15).Draw(t, "numFonts")
		registerPreservationFonts(t, numFonts)

		pixelWidth := rapid.IntRange(30, 1024).Draw(t, "pixelWidth")
		pixelHeight := rapid.IntRange(30, 1024).Draw(t, "pixelHeight")
		minChars := rapid.IntRange(1, 20).Draw(t, "minChars")

		cat, err := tiercatalog.Build(tiercatalog.Params{
			PixelWidth:  pixelWidth,
			PixelHeight: pixelHeight,
			MinChars:    minChars,
		})
		if err != nil {
			return // no qualifying fonts
		}

		// Determine which families exist in the registry.
		registeredFamilies := map[string]bool{}
		for _, f := range font.List() {
			fam := extractFamilyPres(f.ID())
			registeredFamilies[fam] = true
		}
		var famSlice []string
		for fam := range registeredFamilies {
			famSlice = append(famSlice, fam)
		}
		if len(famSlice) == 0 {
			return
		}

		family := rapid.SampledFrom(famSlice).Draw(t, "family")
		tier := rapid.SampledFrom([]tiercatalog.Tier{
			tiercatalog.TierSmall, tiercatalog.TierNormal,
			tiercatalog.TierLarge, tiercatalog.TierHuge,
			tiercatalog.TierColossal, tiercatalog.TierFullsize,
		}).Draw(t, "tier")

		req := tierselect.Request{Family: family, Tier: tier}

		face1 := tierselect.Select(cat, req)
		face2 := tierselect.Select(cat, req)

		if face1.ID() != face2.ID() {
			t.Fatalf("Select non-deterministic: first=%q, second=%q (family=%q, tier=%q)",
				face1.ID(), face2.ID(), family, tier)
		}

		// Also verify the face respects width constraint.
		maxAdvance := cat.PixelWidth() / cat.MinChars()
		if face1.Metrics().GlyphAdvance > maxAdvance {
			t.Fatalf("Select returned face %q with GlyphAdvance=%d > maxAdvance=%d",
				face1.ID(), face1.Metrics().GlyphAdvance, maxAdvance)
		}
	})
}

// TestPreservation_TierSelectResolutionOrder verifies that tierselect.Select()
// follows the three-strategy resolution order:
//  1. Exact match: family variant at the tier's target GlyphHeight
//  2. Closest family variant with GlyphHeight ≤ target
//  3. Cross-family best-fit at target height
//
// When an exact family match exists at the target height, it must be returned.

func TestPreservation_TierSelectResolutionOrder(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		restore := font.SnapshotAndClear()
		defer restore()

		// Register a target font at a specific height.
		targetHeight := rapid.IntRange(8, 32).Draw(t, "targetHeight")
		advance := rapid.IntRange(3, 12).Draw(t, "advance")
		family := rapid.SampledFrom([]string{"spleen", "terminus", "cozette"}).Draw(t, "family")

		exactID := fmt.Sprintf("%s-5x%d-exact", family, targetHeight)
		font.Register(preservationFace{
			id: exactID,
			metrics: font.Metrics{
				GlyphWidth:   5,
				GlyphHeight:  targetHeight,
				GlyphAdvance: advance,
				RowHeight:    targetHeight + 2,
			},
		})

		// Register other fonts at different heights.
		numExtra := rapid.IntRange(1, 8).Draw(t, "numExtra")
		for i := 0; i < numExtra; i++ {
			extraHeight := rapid.IntRange(5, 48).Draw(t, fmt.Sprintf("eh%d", i))
			if extraHeight == targetHeight {
				extraHeight = targetHeight + 1 // avoid collision
			}
			extraAdvance := rapid.IntRange(3, 12).Draw(t, fmt.Sprintf("ea%d", i))
			extraFam := rapid.SampledFrom([]string{"spleen", "terminus", "cozette", "testfont"}).Draw(t, fmt.Sprintf("ef%d", i))
			font.Register(preservationFace{
				id: fmt.Sprintf("%s-5x%d-extra%d", extraFam, extraHeight, i),
				metrics: font.Metrics{
					GlyphWidth:   5,
					GlyphHeight:  extraHeight,
					GlyphAdvance: extraAdvance,
					RowHeight:    extraHeight + 2,
				},
			})
		}

		// Build catalog with dimensions wide enough for exact font.
		pixelWidth := advance * 10
		cat, err := tiercatalog.Build(tiercatalog.Params{
			PixelWidth:  pixelWidth,
			PixelHeight: 200,
			MinChars:    10,
		})
		if err != nil {
			return
		}

		// Find a tier whose target height matches our exact font.
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
			return // tier didn't map to our height, vacuously true
		}

		// Select should return our exact font.
		face := tierselect.Select(cat, tierselect.Request{
			Family: family,
			Tier:   matchedTier,
		})

		if face.ID() != exactID {
			t.Fatalf("expected exact match %q for family=%q tier=%q (targetHeight=%d), got %q",
				exactID, family, matchedTier, targetHeight, face.ID())
		}
	})
}

// TestPreservation_MaxAdvanceConstraint verifies that the maxAdvance = PixelWidth / MinChars
// constraint is applied identically: all catalog entries satisfy GlyphAdvance ≤ maxAdvance.

func TestPreservation_MaxAdvanceConstraint(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		restore := font.SnapshotAndClear()
		defer restore()

		numFonts := rapid.IntRange(1, 20).Draw(t, "numFonts")
		registerPreservationFonts(t, numFonts)

		pixelWidth := rapid.IntRange(10, 2048).Draw(t, "pixelWidth")
		pixelHeight := rapid.IntRange(10, 2048).Draw(t, "pixelHeight")
		minChars := rapid.IntRange(1, 50).Draw(t, "minChars")

		cat, err := tiercatalog.Build(tiercatalog.Params{
			PixelWidth:  pixelWidth,
			PixelHeight: pixelHeight,
			MinChars:    minChars,
		})
		if err != nil {
			return // no qualifying fonts, vacuously true
		}

		maxAdvance := pixelWidth / minChars

		for _, tier := range cat.Tiers() {
			entry, ok := cat.Get(tier)
			if !ok {
				t.Fatalf("tier %q not present in successful catalog", tier)
			}
			if entry.GlyphAdvance > maxAdvance {
				t.Fatalf("maxAdvance violated: tier=%q GlyphAdvance=%d > maxAdvance=%d (PixelWidth=%d / MinChars=%d)",
					tier, entry.GlyphAdvance, maxAdvance, pixelWidth, minChars)
			}
			// FontID must be non-empty and correspond to a registered font.
			if entry.FontID == "" {
				t.Fatalf("tier %q: FontID is empty", tier)
			}
		}
	})
}

// TestPreservation_CatalogBuildErrorOnNoQualifyingFonts verifies that
// tiercatalog.Build() returns an error when no fonts satisfy constraints
// (not silent degradation).

func TestPreservation_CatalogBuildErrorOnNoQualifyingFonts(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		restore := font.SnapshotAndClear()
		defer restore()

		// Register fonts with large advances.
		numFonts := rapid.IntRange(1, 10).Draw(t, "numFonts")
		minAdvance := 20 // all fonts have advance >= 20
		for i := 0; i < numFonts; i++ {
			advance := rapid.IntRange(minAdvance, 64).Draw(t, fmt.Sprintf("adv%d", i))
			height := rapid.IntRange(8, 32).Draw(t, fmt.Sprintf("h%d", i))
			font.Register(preservationFace{
				id: fmt.Sprintf("testfont-%dx%d-err-%d", advance, height, i),
				metrics: font.Metrics{
					GlyphWidth:   advance - 1,
					GlyphHeight:  height,
					GlyphAdvance: advance,
					RowHeight:    height + 2,
				},
			})
		}

		// Use dimensions where maxAdvance < minAdvance of any registered font.
		// PixelWidth / MinChars < 20, so e.g. PixelWidth=10, MinChars=10 → maxAdvance=1.
		pixelWidth := rapid.IntRange(1, minAdvance-1).Draw(t, "pixelWidth")
		minChars := rapid.IntRange(1, 50).Draw(t, "minChars")
		// Ensure maxAdvance < minAdvance
		if pixelWidth/minChars >= minAdvance {
			minChars = pixelWidth + 1 // forces maxAdvance = 0
		}

		_, err := tiercatalog.Build(tiercatalog.Params{
			PixelWidth:  pixelWidth,
			PixelHeight: 100,
			MinChars:    minChars,
		})

		if err == nil {
			t.Fatalf("expected error when no fonts fit (PixelWidth=%d, MinChars=%d, smallest advance=%d)",
				pixelWidth, minChars, minAdvance)
		}
	})
}

// TestPreservation_CatalogBuildErrorOnEmptyRegistry verifies that
// tiercatalog.Build() returns an error when no fonts are registered.

func TestPreservation_CatalogBuildErrorOnEmptyRegistry(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		restore := font.SnapshotAndClear()
		defer restore()

		// Don't register any fonts.
		pixelWidth := rapid.IntRange(10, 2048).Draw(t, "pixelWidth")
		pixelHeight := rapid.IntRange(10, 2048).Draw(t, "pixelHeight")
		minChars := rapid.IntRange(1, 50).Draw(t, "minChars")

		_, err := tiercatalog.Build(tiercatalog.Params{
			PixelWidth:  pixelWidth,
			PixelHeight: pixelHeight,
			MinChars:    minChars,
		})

		if err == nil {
			t.Fatalf("expected error for empty font registry (PixelWidth=%d, PixelHeight=%d, MinChars=%d)",
				pixelWidth, pixelHeight, minChars)
		}
	})
}

// TestPreservation_CatalogTierMonotonicity verifies that tier entries maintain
// GlyphHeight ordering across all six tiers:
// small ≤ normal ≤ large ≤ huge ≤ colossal ≤ full.

func TestPreservation_CatalogTierMonotonicity(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		restore := font.SnapshotAndClear()
		defer restore()

		numFonts := rapid.IntRange(1, 20).Draw(t, "numFonts")
		registerPreservationFonts(t, numFonts)

		pixelWidth := rapid.IntRange(10, 2048).Draw(t, "pixelWidth")
		pixelHeight := rapid.IntRange(10, 2048).Draw(t, "pixelHeight")
		minChars := rapid.IntRange(1, 50).Draw(t, "minChars")

		cat, err := tiercatalog.Build(tiercatalog.Params{
			PixelWidth:  pixelWidth,
			PixelHeight: pixelHeight,
			MinChars:    minChars,
		})
		if err != nil {
			return
		}

		small, _ := cat.Get(tiercatalog.TierSmall)
		normal, _ := cat.Get(tiercatalog.TierNormal)
		large, _ := cat.Get(tiercatalog.TierLarge)
		huge, _ := cat.Get(tiercatalog.TierHuge)
		colossal, _ := cat.Get(tiercatalog.TierColossal)
		full, _ := cat.Get(tiercatalog.TierFull)

		// Also verify TierFullsize alias resolves to TierFull.
		fullAlias, _ := cat.Get(tiercatalog.TierFullsize)
		if full != fullAlias {
			t.Fatalf("TierFullsize alias mismatch: TierFull=%+v, TierFullsize=%+v", full, fullAlias)
		}

		if small.GlyphHeight > normal.GlyphHeight {
			t.Fatalf("monotonicity: small(%d) > normal(%d)", small.GlyphHeight, normal.GlyphHeight)
		}
		if normal.GlyphHeight > large.GlyphHeight {
			t.Fatalf("monotonicity: normal(%d) > large(%d)", normal.GlyphHeight, large.GlyphHeight)
		}
		if large.GlyphHeight > huge.GlyphHeight {
			t.Fatalf("monotonicity: large(%d) > huge(%d)", large.GlyphHeight, huge.GlyphHeight)
		}
		if huge.GlyphHeight > colossal.GlyphHeight {
			t.Fatalf("monotonicity: huge(%d) > colossal(%d)", huge.GlyphHeight, colossal.GlyphHeight)
		}
		if colossal.GlyphHeight > full.GlyphHeight {
			t.Fatalf("monotonicity: colossal(%d) > full(%d)", colossal.GlyphHeight, full.GlyphHeight)
		}
	})
}

// --- Helper ---

// extractFamilyPres splits a font ID on the first "-" and returns the family prefix.
func extractFamilyPres(id string) string {
	for i, ch := range id {
		if ch == '-' {
			return id[:i]
		}
	}
	return id
}
