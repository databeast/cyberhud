package tests_test

import (
	"testing"

	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"pgregory.net/rapid"
)

// Capability Ordering Determines Fitness
//
// This test validates that the capability ordering check in EvaluateFitness
// correctly returns Unsupported when reqs.Capability > hints.Capability,
// and does NOT return Unsupported when hints.Capability >= reqs.Capability
// (assuming dimensions are adequate).

// bugfixMockStyle is a mock FitnessReporter for bug condition exploration.
type bugfixMockStyle struct {
	reqs style.SurfaceRequirements
}

func (m bugfixMockStyle) Name() string                            { return "bugfix-test" }
func (m bugfixMockStyle) Requirements() style.SurfaceRequirements { return m.reqs }
func (m bugfixMockStyle) Supports(hints textlayout.TextHints) style.Fitness {
	return style.EvaluateFitness(m.reqs, hints)
}

// Test 1a: Capability Ordering - Panel meets requirement
//
// Generate SurfaceRequirements with a given Capability level and TextHints
// with Capability >= that level, with adequate dimensions.
// Assert EvaluateFitness does NOT return Unsupported.
func TestBugCondition_RapidRefreshInversion(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate adequate dimensions (non-zero, meeting minimums)
		minW := rapid.IntRange(1, 200).Draw(t, "minWidth")
		minH := rapid.IntRange(1, 200).Draw(t, "minHeight")

		// Style requires ColorFast (highest level)
		reqs := style.SurfaceRequirements{
			Capability: style.ColorFast,
			MinWidth:   minW,
			MinHeight:  minH,
		}

		// Panel provides ColorFast capability (meets requirement)
		pixelW := rapid.IntRange(minW, minW+500).Draw(t, "pixelWidth")
		pixelH := rapid.IntRange(minH, minH+500).Draw(t, "pixelHeight")

		hints := textlayout.TextHints{
			PixelWidth:  pixelW,
			PixelHeight: pixelH,
			Capability:  int(style.ColorFast),
		}

		fitness := style.EvaluateFitness(reqs, hints)

		if fitness == style.Unsupported {
			t.Fatalf("EvaluateFitness returned Unsupported for panel "+
				"with adequate capability (ColorFast).\n"+
				"reqs: {Capability: ColorFast, MinWidth: %d, MinHeight: %d}\n"+
				"hints: {PixelWidth: %d, PixelHeight: %d, Capability: ColorFast}",
				minW, minH, pixelW, pixelH)
		}
	})
}

// Test 1b: Capability Ordering - Panel below requirement returns Unsupported
//
// Generate SurfaceRequirements with a Capability > panel's Capability.
// Assert EvaluateFitness returns Unsupported.
func TestBugCondition_ColorCapabilityMissing(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Style requires ColorFast but panel only provides MonoFast
		minW := rapid.IntRange(1, 200).Draw(t, "minWidth")
		minH := rapid.IntRange(1, 200).Draw(t, "minHeight")

		reqs := style.SurfaceRequirements{
			Capability: style.ColorFast,
			MinWidth:   minW,
			MinHeight:  minH,
		}

		// Panel has lower capability (MonoFast < ColorFast)
		pixelW := rapid.IntRange(minW, minW+500).Draw(t, "pixelWidth")
		pixelH := rapid.IntRange(minH, minH+500).Draw(t, "pixelHeight")

		hints := textlayout.TextHints{
			PixelWidth:  pixelW,
			PixelHeight: pixelH,
			Capability:  int(style.MonoFast),
		}

		fitness := style.EvaluateFitness(reqs, hints)

		if fitness != style.Unsupported {
			t.Fatalf("EvaluateFitness should return Unsupported when panel "+
				"capability (MonoFast) is below style requirement (ColorFast).\n"+
				"reqs: {Capability: ColorFast, MinWidth: %d, MinHeight: %d}\n"+
				"hints: {PixelWidth: %d, PixelHeight: %d, Capability: MonoFast}\n"+
				"got fitness: %d",
				minW, minH, pixelW, pixelH, fitness)
		}
	})
}

// ============================================================================
// Non-Capability Fitness Evaluation Unchanged
//
// For all inputs where Capability == MonoSlow (minimum requirement),
// EvaluateFitness produces results based purely on dimension checks.
// The capability ordering check is trivially satisfied (MonoSlow <= anything).

// ============================================================================

// genPreservationReqs generates SurfaceRequirements with Capability=MonoSlow.
// These inputs are completely unaffected by the capability ordering check.
func genPreservationReqs(t *rapid.T) style.SurfaceRequirements {
	minW := rapid.IntRange(0, 500).Draw(t, "minWidth")
	minH := rapid.IntRange(0, 500).Draw(t, "minHeight")

	prefW := 0
	if minW > 0 {
		prefW = rapid.IntRange(minW, minW+500).Draw(t, "preferredWidth")
	} else {
		prefW = rapid.IntRange(0, 1000).Draw(t, "preferredWidth")
	}
	prefH := 0
	if minH > 0 {
		prefH = rapid.IntRange(minH, minH+500).Draw(t, "preferredHeight")
	} else {
		prefH = rapid.IntRange(0, 1000).Draw(t, "preferredHeight")
	}

	minRows := rapid.IntRange(0, 20).Draw(t, "minRows")
	minChars := rapid.IntRange(0, 40).Draw(t, "minCharsPerLine")

	return style.SurfaceRequirements{
		MinWidth:        minW,
		MinHeight:       minH,
		PreferredWidth:  prefW,
		PreferredHeight: prefH,
		Capability:      style.MonoSlow, // Minimum requirement — always satisfied
		MinRows:         minRows,
		MinCharsPerLine: minChars,
	}
}

// genPreservationHints generates TextHints covering the full input space.
func genPreservationHints(t *rapid.T) textlayout.TextHints {
	return textlayout.TextHints{
		PixelWidth:         rapid.IntRange(0, 1000).Draw(t, "pixelWidth"),
		PixelHeight:        rapid.IntRange(0, 1000).Draw(t, "pixelHeight"),
		GlyphWidth:         rapid.IntRange(1, 32).Draw(t, "glyphWidth"),
		GlyphHeight:        rapid.IntRange(1, 32).Draw(t, "glyphHeight"),
		GlyphAdvance:       rapid.IntRange(1, 32).Draw(t, "glyphAdvance"),
		RowHeight:          rapid.IntRange(1, 64).Draw(t, "rowHeight"),
		PreferEventRefresh: rapid.Bool().Draw(t, "preferEventRefresh"),
		Capability:         rapid.IntRange(0, 5).Draw(t, "capability"),
	}
}

// TestPreservation_ZeroDimensionUnsupported verifies that zero-dimension hints
// always return Unsupported.

func TestPreservation_ZeroDimensionUnsupported(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		reqs := genPreservationReqs(t)

		zeroWidth := rapid.Bool().Draw(t, "zeroWidth")
		var hints textlayout.TextHints
		if zeroWidth {
			hints = textlayout.TextHints{
				PixelWidth:         0,
				PixelHeight:        rapid.IntRange(0, 1000).Draw(t, "pixelHeight"),
				GlyphWidth:         rapid.IntRange(1, 32).Draw(t, "glyphWidth"),
				GlyphHeight:        rapid.IntRange(1, 32).Draw(t, "glyphHeight"),
				GlyphAdvance:       rapid.IntRange(1, 32).Draw(t, "glyphAdvance"),
				RowHeight:          rapid.IntRange(1, 64).Draw(t, "rowHeight"),
				PreferEventRefresh: rapid.Bool().Draw(t, "preferEventRefresh"),
				Capability:         rapid.IntRange(0, 5).Draw(t, "capability"),
			}
		} else {
			hints = textlayout.TextHints{
				PixelWidth:         rapid.IntRange(0, 1000).Draw(t, "pixelWidth"),
				PixelHeight:        0,
				GlyphWidth:         rapid.IntRange(1, 32).Draw(t, "glyphWidth"),
				GlyphHeight:        rapid.IntRange(1, 32).Draw(t, "glyphHeight"),
				GlyphAdvance:       rapid.IntRange(1, 32).Draw(t, "glyphAdvance"),
				RowHeight:          rapid.IntRange(1, 64).Draw(t, "rowHeight"),
				PreferEventRefresh: rapid.Bool().Draw(t, "preferEventRefresh"),
				Capability:         rapid.IntRange(0, 5).Draw(t, "capability"),
			}
		}

		fitness := style.EvaluateFitness(reqs, hints)
		if fitness != style.Unsupported {
			t.Fatalf("EvaluateFitness should return Unsupported for zero-dimension hints.\n"+
				"reqs: %+v\nhints: {PixelWidth: %d, PixelHeight: %d}\nfitness: %d",
				reqs, hints.PixelWidth, hints.PixelHeight, fitness)
		}
	})
}

// TestPreservation_UndersizedDimensionsUnsupported verifies undersized panels
// return Unsupported.

func TestPreservation_UndersizedDimensionsUnsupported(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		minW := rapid.IntRange(2, 500).Draw(t, "minWidth")
		minH := rapid.IntRange(2, 500).Draw(t, "minHeight")

		reqs := style.SurfaceRequirements{
			MinWidth:   minW,
			MinHeight:  minH,
			Capability: style.MonoSlow,
		}

		violateWidth := rapid.Bool().Draw(t, "violateWidth")
		var pixW, pixH int
		if violateWidth {
			pixW = rapid.IntRange(1, minW-1).Draw(t, "pixelWidth")
			pixH = rapid.IntRange(1, 1000).Draw(t, "pixelHeight")
		} else {
			pixW = rapid.IntRange(1, 1000).Draw(t, "pixelWidth")
			pixH = rapid.IntRange(1, minH-1).Draw(t, "pixelHeight")
		}

		hints := textlayout.TextHints{
			PixelWidth:         pixW,
			PixelHeight:        pixH,
			GlyphWidth:         rapid.IntRange(1, 32).Draw(t, "glyphWidth"),
			GlyphHeight:        rapid.IntRange(1, 32).Draw(t, "glyphHeight"),
			GlyphAdvance:       rapid.IntRange(1, 32).Draw(t, "glyphAdvance"),
			RowHeight:          rapid.IntRange(1, 64).Draw(t, "rowHeight"),
			PreferEventRefresh: rapid.Bool().Draw(t, "preferEventRefresh"),
			Capability:         rapid.IntRange(0, 5).Draw(t, "capability"),
		}

		fitness := style.EvaluateFitness(reqs, hints)
		if fitness != style.Unsupported {
			t.Fatalf("EvaluateFitness should return Unsupported for undersized panel.\n"+
				"reqs: {MinWidth: %d, MinHeight: %d}\n"+
				"hints: {PixelWidth: %d, PixelHeight: %d}\nfitness: %d",
				minW, minH, pixW, pixH, fitness)
		}
	})
}

// TestPreservation_AdequateDimensionsOptimal verifies panels meeting all dimensions
// return Optimal.

func TestPreservation_AdequateDimensionsOptimal(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		minW := rapid.IntRange(1, 200).Draw(t, "minWidth")
		minH := rapid.IntRange(1, 200).Draw(t, "minHeight")
		prefW := rapid.IntRange(minW, minW+200).Draw(t, "preferredWidth")
		prefH := rapid.IntRange(minH, minH+200).Draw(t, "preferredHeight")

		reqs := style.SurfaceRequirements{
			MinWidth:        minW,
			MinHeight:       minH,
			PreferredWidth:  prefW,
			PreferredHeight: prefH,
			Capability:      style.MonoSlow,
			MinRows:         0,
			MinCharsPerLine: 0,
		}

		pixW := rapid.IntRange(prefW, prefW+500).Draw(t, "pixelWidth")
		pixH := rapid.IntRange(prefH, prefH+500).Draw(t, "pixelHeight")

		hints := textlayout.TextHints{
			PixelWidth:         pixW,
			PixelHeight:        pixH,
			GlyphWidth:         rapid.IntRange(1, 32).Draw(t, "glyphWidth"),
			GlyphHeight:        rapid.IntRange(1, 32).Draw(t, "glyphHeight"),
			GlyphAdvance:       rapid.IntRange(1, 32).Draw(t, "glyphAdvance"),
			RowHeight:          rapid.IntRange(1, 64).Draw(t, "rowHeight"),
			PreferEventRefresh: rapid.Bool().Draw(t, "preferEventRefresh"),
			Capability:         rapid.IntRange(0, 5).Draw(t, "capability"),
		}

		fitness := style.EvaluateFitness(reqs, hints)
		if fitness < style.Optimal {
			t.Fatalf("EvaluateFitness should return Optimal tier when all preferred dims are met.\n"+
				"reqs: {MinWidth: %d, MinHeight: %d, PreferredWidth: %d, PreferredHeight: %d}\n"+
				"hints: {PixelWidth: %d, PixelHeight: %d}\nfitness: %d",
				minW, minH, prefW, prefH, pixW, pixH, fitness)
		}
	})
}

// TestPreservation_MeetsMinButNotPreferredFull verifies panels meeting min but
// not preferred return Full.

func TestPreservation_MeetsMinButNotPreferredFull(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		minW := rapid.IntRange(1, 200).Draw(t, "minWidth")
		minH := rapid.IntRange(1, 200).Draw(t, "minHeight")
		prefW := rapid.IntRange(minW+50, minW+500).Draw(t, "preferredWidth")
		prefH := rapid.IntRange(minH+50, minH+500).Draw(t, "preferredHeight")

		reqs := style.SurfaceRequirements{
			MinWidth:        minW,
			MinHeight:       minH,
			PreferredWidth:  prefW,
			PreferredHeight: prefH,
			Capability:      style.MonoSlow,
			MinRows:         0,
			MinCharsPerLine: 0,
		}

		// Meet min but below preferred width.
		pixW := rapid.IntRange(minW, prefW-1).Draw(t, "pixelWidth")
		pixH := rapid.IntRange(prefH, prefH+500).Draw(t, "pixelHeight")

		hints := textlayout.TextHints{
			PixelWidth:         pixW,
			PixelHeight:        pixH,
			GlyphWidth:         rapid.IntRange(1, 32).Draw(t, "glyphWidth"),
			GlyphHeight:        rapid.IntRange(1, 32).Draw(t, "glyphHeight"),
			GlyphAdvance:       rapid.IntRange(1, 32).Draw(t, "glyphAdvance"),
			RowHeight:          rapid.IntRange(1, 64).Draw(t, "rowHeight"),
			PreferEventRefresh: rapid.Bool().Draw(t, "preferEventRefresh"),
			Capability:         rapid.IntRange(0, 5).Draw(t, "capability"),
		}

		fitness := style.EvaluateFitness(reqs, hints)
		if fitness < style.Full || fitness >= style.Optimal {
			t.Fatalf("EvaluateFitness should return Full tier when min is met but preferred is not.\n"+
				"reqs: {MinWidth: %d, MinHeight: %d, PreferredWidth: %d, PreferredHeight: %d}\n"+
				"hints: {PixelWidth: %d, PixelHeight: %d}\nfitness: %d",
				minW, minH, prefW, prefH, pixW, pixH, fitness)
		}
	})
}

// TestPreservation_TextConstraintsUnsupported verifies that when text constraints
// cannot be satisfied, EvaluateFitness returns Unsupported.

func TestPreservation_TextConstraintsUnsupported(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		minRows := rapid.IntRange(50, 200).Draw(t, "minRows")
		minChars := rapid.IntRange(50, 200).Draw(t, "minCharsPerLine")

		reqs := style.SurfaceRequirements{
			MinWidth:        1,
			MinHeight:       1,
			Capability:      style.MonoSlow,
			MinRows:         minRows,
			MinCharsPerLine: minChars,
		}

		hints := textlayout.TextHints{
			PixelWidth:         rapid.IntRange(10, 30).Draw(t, "pixelWidth"),
			PixelHeight:        rapid.IntRange(10, 30).Draw(t, "pixelHeight"),
			GlyphWidth:         5,
			GlyphHeight:        7,
			GlyphAdvance:       6,
			RowHeight:          10,
			PreferEventRefresh: rapid.Bool().Draw(t, "preferEventRefresh"),
			Capability:         rapid.IntRange(0, 5).Draw(t, "capability"),
		}

		fitness := style.EvaluateFitness(reqs, hints)
		if fitness != style.Unsupported {
			t.Fatalf("EvaluateFitness should return Unsupported when text constraints cannot be met.\n"+
				"reqs: {MinRows: %d, MinCharsPerLine: %d}\n"+
				"hints: {PixelWidth: %d, PixelHeight: %d}\nfitness: %d",
				minRows, minChars, hints.PixelWidth, hints.PixelHeight, fitness)
		}
	})
}

// TestPreservation_FullBehaviorAcrossAllInputs is the comprehensive preservation property.

func TestPreservation_FullBehaviorAcrossAllInputs(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		reqs := genPreservationReqs(t)
		hints := genPreservationHints(t)

		actual := style.EvaluateFitness(reqs, hints)
		expectedTier := expectedPreservationFitness(reqs, hints)

		actualTier := fitnessTier(actual)
		if actualTier != expectedTier {
			t.Fatalf("EvaluateFitness tier mismatch for non-capability inputs.\n"+
				"reqs: {MinWidth: %d, MinHeight: %d, PreferredWidth: %d, PreferredHeight: %d, "+
				"MinRows: %d, MinCharsPerLine: %d}\n"+
				"hints: {PixelWidth: %d, PixelHeight: %d, GlyphAdvance: %d, RowHeight: %d}\n"+
				"expected tier: %d, actual: %d (tier: %d)",
				reqs.MinWidth, reqs.MinHeight, reqs.PreferredWidth, reqs.PreferredHeight,
				reqs.MinRows, reqs.MinCharsPerLine,
				hints.PixelWidth, hints.PixelHeight, hints.GlyphAdvance, hints.RowHeight,
				expectedTier, actual, actualTier)
		}
	})
}

// fitnessTier extracts the base tier from a scored Fitness value.
func fitnessTier(f style.Fitness) style.Fitness {
	switch {
	case f >= style.Optimal:
		return style.Optimal
	case f >= style.Full:
		return style.Full
	case f >= style.Degraded:
		return style.Degraded
	default:
		return style.Unsupported
	}
}

// expectedPreservationFitness computes expected fitness for non-capability inputs.
func expectedPreservationFitness(reqs style.SurfaceRequirements, hints textlayout.TextHints) style.Fitness {
	if hints.PixelWidth == 0 || hints.PixelHeight == 0 {
		return style.Unsupported
	}
	if reqs.MinWidth > 0 && hints.PixelWidth < reqs.MinWidth {
		return style.Unsupported
	}
	if reqs.MinHeight > 0 && hints.PixelHeight < reqs.MinHeight {
		return style.Unsupported
	}
	// Capability check: MonoSlow <= any capability, so always passes.
	if reqs.MinRows > 0 || reqs.MinCharsPerLine > 0 {
		if style.TextFitness(reqs, hints) == style.Unsupported {
			return style.Unsupported
		}
	}
	meetsPreferred := true
	if reqs.PreferredWidth > 0 && hints.PixelWidth < reqs.PreferredWidth {
		meetsPreferred = false
	}
	if reqs.PreferredHeight > 0 && hints.PixelHeight < reqs.PreferredHeight {
		meetsPreferred = false
	}
	if meetsPreferred {
		return style.Optimal
	}
	return style.Full
}
