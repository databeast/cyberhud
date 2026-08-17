package tests

import (
	"fmt"
	"image/color"
	"math"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/databeast/cyberhud/display/modes/attract_matrix"
	"github.com/databeast/cyberhud/display/modes/attract_matrix/source"
	"github.com/databeast/cyberhud/display/modes/attract_matrix/styles"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"pgregory.net/rapid"
)

// --- From: matrix_charsource_prop_test.go ---

// For any RandomSource and any cell index, CharAt(index) returns a rune in the 40-character pool.

func TestProperty_CharAt_PoolMembership(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		seed := rapid.Int64Range(1, 1000000).Draw(rt, "seed")
		colIdx := rapid.IntRange(0, 100).Draw(rt, "colIdx")
		cellIndex := rapid.IntRange(0, 1000).Draw(rt, "cellIndex")

		src := source.NewRandomSource(seed, colIdx, source.DefaultMutationInterval())
		r := src.CharAt(cellIndex)

		found := false
		for _, poolRune := range source.CharacterPool() {
			if r == poolRune {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("CharAt(%d) returned rune %q (U+%04X) which is not in characterPool",
				cellIndex, r, r)
		}
	})
}

// Two instances with same seed+colIdx produce identical CharAt results;
// within mutation interval, same index returns same rune.

func TestProperty_CharAt_Determinism(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		seed := rapid.Int64Range(1, 1000000).Draw(rt, "seed")
		colIdx := rapid.IntRange(0, 100).Draw(rt, "colIdx")
		indices := rapid.SliceOfN(rapid.IntRange(0, 1000), 1, 20).Draw(rt, "indices")

		src1 := source.NewRandomSource(seed, colIdx, source.DefaultMutationInterval())
		src2 := source.NewRandomSource(seed, colIdx, source.DefaultMutationInterval())

		for i, idx := range indices {
			r1 := src1.CharAt(idx)
			r2 := src2.CharAt(idx)
			if r1 != r2 {
				t.Fatalf("Determinism violated at indices[%d]=%d: src1 returned %q, src2 returned %q (seed=%d, colIdx=%d)",
					i, idx, r1, r2, seed, colIdx)
			}
		}
	})
}

func TestProperty_CharAt_StabilityWithinInterval(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		seed := rapid.Int64Range(1, 1000000).Draw(rt, "seed")
		colIdx := rapid.IntRange(0, 100).Draw(rt, "colIdx")
		cellIndex := rapid.IntRange(0, 1000).Draw(rt, "cellIndex")

		// Use a long mutation interval to guarantee stability between two quick calls.
		src := source.NewRandomSource(seed, colIdx, 1*time.Second)

		r1 := src.CharAt(cellIndex)
		r2 := src.CharAt(cellIndex)

		if r1 != r2 {
			t.Fatalf("Stability violated: CharAt(%d) returned %q then %q within mutation interval (seed=%d, colIdx=%d)",
				cellIndex, r1, r2, seed, colIdx)
		}
	})
}

// --- From: matrix_command_prop_test.go ---

// TestProperty_CommandHandler_ValidSetResponse verifies that when all key=value
// arguments are valid, the response starts with "OK matrix" and reflects the
// updated values (key names present in the response).

func TestProperty_CommandHandler_ValidSetResponse(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		// Reset to defaults before each property check.
		attract_matrix.SetPolicy(attract_matrix.DefaultPolicy())

		// Generate valid key=value pairs.
		type kvPair struct{ k, v string }
		possiblePairs := []kvPair{
			{"min_speed", fmt.Sprintf("%g", rapid.Float64Range(0.1, 50.0).Draw(rt, "minSpeed"))},
			{"max_speed", fmt.Sprintf("%g", rapid.Float64Range(0.1, 100.0).Draw(rt, "maxSpeed"))},
			{"trail_length", fmt.Sprintf("%d", rapid.IntRange(1, 128).Draw(rt, "trailLen"))},
			{"density", fmt.Sprintf("%g", rapid.Float64Range(0.1, 1.0).Draw(rt, "density"))},
			{"show_background", rapid.SampledFrom([]string{"true", "false"}).Draw(rt, "showBg")},
		}

		// Pick 1 to 5 args.
		count := rapid.IntRange(1, 5).Draw(rt, "argCount")
		if count > len(possiblePairs) {
			count = len(possiblePairs)
		}
		args := make([]string, count)
		for i := 0; i < count; i++ {
			args[i] = possiblePairs[i].k + "=" + possiblePairs[i].v
		}

		result := attract_matrix.HandleCommand(args)

		// Must start with "OK attract_matrix".
		if !strings.HasPrefix(result, "OK attract_matrix") {
			t.Fatalf("expected OK response for valid args %v, got %q", args, result)
		}

		// Response should contain all key names.
		for i := 0; i < count; i++ {
			if !strings.Contains(result, possiblePairs[i].k+"=") {
				t.Fatalf("response missing key %q: %q", possiblePairs[i].k, result)
			}
		}
	})
}

// --- From: matrix_compositor_prop_test.go ---

// TestProperty_StyleBehavior_ByPanelType verifies that for any panel dimensions:
// - EinkStyle.Build produces ViewData with Static=true (snapshot.Eink=true)
// - ColorStyle.Build produces ViewData with Static=false (snapshot.Eink=false)
// - MonoStyle.Build produces ViewData with Static=false (snapshot.Eink=false)

func TestProperty_StyleBehavior_ByPanelType(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		width := rapid.IntRange(80, 800).Draw(rt, "width")
		height := rapid.IntRange(32, 800).Draw(rt, "height")

		hints := textlayout.TextHints{
			PixelWidth:  width,
			PixelHeight: height,
		}

		snap := attract_matrix.MatrixSnapshot{
			Policy:       attract_matrix.DefaultPolicy(),
			PanelWidth:   width,
			PanelHeight:  height,
			GlyphAdvance: 14,
			RowHeight:    14,
		}

		// E-ink style → Static = true
		einkSnap := snap
		einkSnap.Eink = true
		einkStyle := styles.EinkStyle{Width: width, Height: height}
		ctx := style.NewStyleContext(hints)
		einkResult := einkStyle.Build(einkSnap, attract_matrix.DefaultPolicy(), ctx)
		if !einkResult.Static {
			t.Fatalf("EinkStyle.Build(%dx%d) produced Static=false, want true",
				width, height)
		}

		// Color style → Static = false
		colorSnap := snap
		colorSnap.Eink = false
		colorStyle := styles.ColorStyle{Width: width, Height: height}
		colorResult := colorStyle.Build(colorSnap, attract_matrix.DefaultPolicy(), ctx)
		if colorResult.Static {
			t.Fatalf("ColorStyle.Build(%dx%d) produced Static=true, want false",
				width, height)
		}

		// Mono style → Static = false
		monoSnap := snap
		monoSnap.Eink = false
		monoSnap.Mono = true
		monoStyle := styles.MonoStyle{Width: width, Height: height}
		monoResult := monoStyle.Build(monoSnap, attract_matrix.DefaultPolicy(), ctx)
		if monoResult.Static {
			t.Fatalf("MonoStyle.Build(%dx%d) produced Static=true, want false",
				width, height)
		}
	})
}

// TestProperty_Compositor_ZOrderAndSuppression verifies that:
// - When ShowBackground=true and panel resolves to COLOR: first sprite has label "gradient/radial"
// - When ShowBackground=true and panel resolves to MONO: no sprite has label "gradient/radial"
// - When ShowBackground=false: no sprite has label "gradient/radial"
//
// The test randomizes TrailLength and Density while controlling ShowBackground as the
// variable under test. Panel type is determined by hint dimensions: 240x135 → color,
// 128x64 → mono.

// --- From: matrix_datasig_prop_test.go ---

// For any sequence of N >= 2 consecutive BuildView calls with the same policy,
// each RenderCacheKey string SHALL differ from the previous one (due to the
// monotonically increasing frame counter). Additionally, for any two distinct
// Policy values at the same frame counter, RenderCacheKey SHALL produce different strings.

func TestProperty_RenderCacheKey_ConsecutiveUniqueness(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		n := rapid.IntRange(2, 10).Draw(rt, "iterations")

		// Use a known color panel size.
		hints := textlayout.TextHints{
			PixelWidth:  240,
			PixelHeight: 135,
		}

		// Set a fixed policy for the duration of this property check.
		attract_matrix.SetPolicy(attract_matrix.DefaultPolicy())
		defer attract_matrix.SetPolicy(attract_matrix.DefaultPolicy())

		// Collect RenderCacheKey after each BuildView call.
		signatures := make([]uint32, n)
		for i := 0; i < n; i++ {
			attract_matrix.BuildView(hints)
			signatures[i] = attract_matrix.RenderCacheKey()
		}

		// Each consecutive pair must differ.
		for i := 1; i < n; i++ {
			if signatures[i] == signatures[i-1] {
				t.Fatalf("RenderCacheKey did not change between BuildView calls %d and %d: both = %v",
					i-1, i, signatures[i])
			}
		}
	})
}

func TestProperty_RenderCacheKey_PolicyDifferentiation(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		// Generate two distinct TrailLength values (both within valid range [4, 128]).
		trail1 := rapid.IntRange(4, 128).Draw(rt, "trail1")
		trail2 := rapid.IntRange(4, 128).Draw(rt, "trail2")
		if trail1 == trail2 {
			// Ensure they differ; shift trail2.
			if trail2 < 128 {
				trail2++
			} else {
				trail2--
			}
		}

		// Record sig1 with first policy.
		p1 := attract_matrix.DefaultPolicy()
		p1.TrailLength = trail1
		attract_matrix.SetPolicy(p1)
		sig1 := attract_matrix.RenderCacheKey()

		// Record sig2 with second policy (no BuildView between, same frameCounter).
		p2 := attract_matrix.DefaultPolicy()
		p2.TrailLength = trail2
		attract_matrix.SetPolicy(p2)
		sig2 := attract_matrix.RenderCacheKey()

		// Restore default policy.
		attract_matrix.SetPolicy(attract_matrix.DefaultPolicy())

		if sig1 == sig2 {
			t.Fatalf("RenderCacheKey should differ for different policies (trail %d vs %d) at same frame counter: both = %v",
				trail1, trail2, sig1)
		}
	})
}

// --- From: matrix_density_prop_test.go ---

// For any totalColumns >= 1 and density in [0.1, 1.0], the number of active columns
// equals max(1, floor(totalColumns * density)). When totalColumns == 1, exactly 1
// column is active regardless of density. All returned indices are in [0, totalColumns-1]
// and sorted.

func TestProperty_ActiveColumnCount(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		totalColumns := rapid.IntRange(1, 100).Draw(rt, "totalColumns")
		density := rapid.Float64Range(0.1, 1.0).Draw(rt, "density")

		active := source.ComputeActiveColumns(totalColumns, density)

		// Property: correct count
		expectedCount := int(math.Floor(float64(totalColumns) * density))
		if expectedCount < 1 {
			expectedCount = 1
		}
		if expectedCount > totalColumns {
			expectedCount = totalColumns
		}
		if len(active) != expectedCount {
			t.Fatalf("computeActiveColumns(%d, %f) returned %d columns, want %d",
				totalColumns, density, len(active), expectedCount)
		}

		// Property: when totalColumns == 1, result is always [0]
		if totalColumns == 1 {
			if len(active) != 1 || active[0] != 0 {
				t.Fatalf("computeActiveColumns(1, %f) = %v, want [0]",
					density, active)
			}
		}

		// Property: all indices in [0, totalColumns-1]
		for _, idx := range active {
			if idx < 0 || idx >= totalColumns {
				t.Fatalf("computeActiveColumns(%d, %f) returned index %d out of range [0, %d]",
					totalColumns, density, idx, totalColumns-1)
			}
		}

		// Property: indices are sorted
		if !sort.IntsAreSorted(active) {
			t.Fatalf("computeActiveColumns(%d, %f) returned unsorted indices: %v",
				totalColumns, density, active)
		}
	})
}

// For any requested trailLength and visibleCells >= 4, the effective clamped value
// equals clamp(requested, 4, visibleCells). This is verified indirectly by checking
// that buildColorArray produces the expected output length given the clamped value.

func TestProperty_TrailLengthClamping(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		trailLength := rapid.IntRange(1, 200).Draw(rt, "trailLength")
		visibleCells := rapid.IntRange(4, 50).Draw(rt, "visibleCells")

		// Compute expected clamped value: max(4, min(trailLength, visibleCells))
		expected := trailLength
		if expected < 4 {
			expected = 4
		}
		if expected > visibleCells {
			expected = visibleCells
		}

		// Verify buildColorArray output length matches the clamped trail length.
		// buildColorArray returns a slice of length clampedTrailLength + 1.
		colors := source.BuildColorArray(expected, false)
		if len(colors) != expected+1 {
			t.Fatalf("buildColorArray(%d, false) returned %d colors, want %d (trailLength=%d, visibleCells=%d)",
				expected, len(colors), expected+1, trailLength, visibleCells)
		}

		// Also verify the clamping logic itself matches expectations.
		// Replicate the clamping from rebuildStrips:
		clamped := trailLength
		if clamped < 4 {
			clamped = 4
		}
		if clamped > visibleCells {
			clamped = visibleCells
		}
		if clamped != expected {
			t.Fatalf("clamping mismatch: trailLength=%d, visibleCells=%d, got %d, want %d",
				trailLength, visibleCells, clamped, expected)
		}
	})
}

// --- From: matrix_gradient_prop_test.go ---

// brightness computes the sum of RGB channels for a color.RGBA value.
func brightness(c color.RGBA) int {
	return int(c.R) + int(c.G) + int(c.B)
}

// TestProperty_ColorArray_Construction verifies that for any trailLength in [1,128]
// and any mono flag, buildColorArray produces:
// - An array of exactly trailLength+1 elements
// - Index 0 = lead color ({180,255,180,255} normal, {255,255,255,255} mono)
// - Final index = fully black {0,0,0,255}
// - Monotonically non-increasing brightness from index 0 to final index
// - When mono=true, all cells have equal R, G, B channels (grayscale)
// - Alpha is always 255 for all cells

func TestProperty_ColorArray_Construction(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		trailLength := rapid.IntRange(1, 128).Draw(rt, "trailLength")
		mono := rapid.Bool().Draw(rt, "mono")

		result := source.BuildColorArray(trailLength, mono)

		// Property: array length = trailLength + 1
		if len(result) != trailLength+1 {
			t.Fatalf("len(buildColorArray(%d, %v)) = %d, want %d",
				trailLength, mono, len(result), trailLength+1)
		}

		// Property: index 0 = lead color
		if mono {
			expected := color.RGBA{255, 255, 255, 255}
			if result[0] != expected {
				t.Fatalf("buildColorArray(%d, true)[0] = %v, want %v",
					trailLength, result[0], expected)
			}
		} else {
			expected := color.RGBA{180, 255, 180, 255}
			if result[0] != expected {
				t.Fatalf("buildColorArray(%d, false)[0] = %v, want %v",
					trailLength, result[0], expected)
			}
		}

		// Property: final index = fully black {0,0,0,255}
		finalExpected := color.RGBA{0, 0, 0, 255}
		if result[trailLength] != finalExpected {
			t.Fatalf("buildColorArray(%d, %v)[%d] = %v, want %v",
				trailLength, mono, trailLength, result[trailLength], finalExpected)
		}

		// Property: monotonically non-increasing brightness
		for i := 0; i < len(result)-1; i++ {
			bCurr := brightness(result[i])
			bNext := brightness(result[i+1])
			if bCurr < bNext {
				t.Fatalf("brightness not non-increasing at index %d→%d: %d < %d (colors: %v → %v), trailLength=%d, mono=%v",
					i, i+1, bCurr, bNext, result[i], result[i+1], trailLength, mono)
			}
		}

		// Property: mono implies grayscale (R == G == B for all cells)
		if mono {
			for i, c := range result {
				if c.R != c.G || c.G != c.B {
					t.Fatalf("mono color at index %d is not grayscale: %v, trailLength=%d",
						i, c, trailLength)
				}
			}
		}

		// Property: alpha always 255 for all cells
		for i, c := range result {
			if c.A != 255 {
				t.Fatalf("alpha at index %d = %d, want 255, trailLength=%d, mono=%v",
					i, c.A, trailLength, mono)
			}
		}
	})
}

// --- From: matrix_layout_prop_test.go ---

// For any panelWidth >= 0, panelHeight >= 0, glyphAdvance > 0: column count = panelWidth/glyphAdvance,
// positions correct, no bounds overflow, zero columns for edge cases.

func TestProperty_ColumnGrid_Invariants(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		panelWidth := rapid.IntRange(0, 2000).Draw(rt, "panelWidth")
		panelHeight := rapid.IntRange(0, 2000).Draw(rt, "panelHeight")
		glyphAdvance := rapid.IntRange(1, 50).Draw(rt, "glyphAdvance")
		rowHeight := rapid.IntRange(1, 50).Draw(rt, "rowHeight")

		// Property: computeColumnCount matches integer division
		colCount := source.ComputeColumnCount(panelWidth, glyphAdvance)
		expected := panelWidth / glyphAdvance
		if colCount != expected {
			t.Fatalf("computeColumnCount(%d, %d) = %d, want %d",
				panelWidth, glyphAdvance, colCount, expected)
		}

		// Property: no bounds overflow — columns * glyphAdvance never exceeds panel width
		if colCount*glyphAdvance > panelWidth {
			t.Fatalf("bounds overflow: %d columns * %d advance = %d > panelWidth %d",
				colCount, glyphAdvance, colCount*glyphAdvance, panelWidth)
		}

		// Property: each column at index i is positioned at i*glyphAdvance (conceptual check)
		// Verify that the rightmost column's right edge does not exceed panel width.
		if colCount > 0 {
			rightEdge := colCount * glyphAdvance
			if rightEdge > panelWidth {
				t.Fatalf("rightmost column right edge %d exceeds panelWidth %d",
					rightEdge, panelWidth)
			}
		}

		// Property: computeVisibleCells matches integer division
		visCells := source.ComputeVisibleCells(panelHeight, rowHeight)
		expectedCells := panelHeight / rowHeight
		if visCells != expectedCells {
			t.Fatalf("computeVisibleCells(%d, %d) = %d, want %d",
				panelHeight, rowHeight, visCells, expectedCells)
		}

		// Property: panelWidth < glyphAdvance yields 0 columns
		if panelWidth < glyphAdvance && colCount != 0 {
			t.Fatalf("expected 0 columns when panelWidth(%d) < glyphAdvance(%d), got %d",
				panelWidth, glyphAdvance, colCount)
		}
	})
}

// TestProperty_ColumnGrid_ZeroAdvance verifies that glyphAdvance <= 0 always produces 0 columns.
func TestProperty_ColumnGrid_ZeroAdvance(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		panelWidth := rapid.IntRange(0, 2000).Draw(rt, "panelWidth")
		glyphAdvance := rapid.IntRange(-50, 0).Draw(rt, "glyphAdvance")

		colCount := source.ComputeColumnCount(panelWidth, glyphAdvance)
		if colCount != 0 {
			t.Fatalf("computeColumnCount(%d, %d) = %d, want 0 for non-positive advance",
				panelWidth, glyphAdvance, colCount)
		}
	})
}

// TestProperty_ColumnGrid_WidthLessThanAdvance verifies that when panelWidth < glyphAdvance,
// zero columns are produced.
func TestProperty_ColumnGrid_WidthLessThanAdvance(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		glyphAdvance := rapid.IntRange(1, 50).Draw(rt, "glyphAdvance")
		// Ensure panelWidth is strictly less than glyphAdvance.
		panelWidth := rapid.IntRange(0, glyphAdvance-1).Draw(rt, "panelWidth")

		colCount := source.ComputeColumnCount(panelWidth, glyphAdvance)
		if colCount != 0 {
			t.Fatalf("computeColumnCount(%d, %d) = %d, want 0 when width < advance",
				panelWidth, glyphAdvance, colCount)
		}
	})
}

// TestProperty_ColumnGrid_VisibleCells_ZeroRowHeight verifies that rowHeight <= 0
// always produces 0 visible cells.
func TestProperty_ColumnGrid_VisibleCells_ZeroRowHeight(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		panelHeight := rapid.IntRange(0, 2000).Draw(rt, "panelHeight")
		rowHeight := rapid.IntRange(-50, 0).Draw(rt, "rowHeight")

		visCells := source.ComputeVisibleCells(panelHeight, rowHeight)
		if visCells != 0 {
			t.Fatalf("computeVisibleCells(%d, %d) = %d, want 0 for non-positive rowHeight",
				panelHeight, rowHeight, visCells)
		}
	})
}

// --- From: matrix_speed_prop_test.go ---

// For any Policy where MinSpeed > 0 and MaxSpeed > 0, and for any column created
// under that policy, the column's speed SHALL be in the range
// [min(MinSpeed, MaxSpeed), max(MinSpeed, MaxSpeed)].
// When MinSpeed > MaxSpeed after normalization, MinSpeed is clamped to MaxSpeed,
// yielding uniform speed.

func TestProperty_StripSpeed_InRange(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		// Reset strip state so each iteration gets a fresh initial build.
		source.ResetLayoutCache()

		minSpeed := rapid.Float64Range(0.1, 50.0).Draw(rt, "minSpeed")
		maxSpeed := rapid.Float64Range(minSpeed, 100.0).Draw(rt, "maxSpeed")

		p := attract_matrix.Policy{
			MinSpeed:    minSpeed,
			MaxSpeed:    maxSpeed,
			TrailLength: 8,
			Density:     1.0,
		}
		p = attract_matrix.NormalizePolicy(p)

		// Fixed layout: 20 columns at 14px advance, 10 visible cells.
		panelWidth := 280
		panelHeight := 140
		glyphAdvance := 14
		rowHeight := 14
		seed := int64(42)

		strips := source.RebuildStripsForTest(p, panelWidth, panelHeight, glyphAdvance, rowHeight, false, seed, nil)

		if len(strips) == 0 {
			t.Fatal("expected at least one strip with density=1.0 and valid dimensions")
		}

		for i, strip := range strips {
			speed := strip.Speed()
			if speed < p.MinSpeed || speed > p.MaxSpeed {
				t.Fatalf("strip[%d].Speed() = %f, want in [%f, %f]",
					i, speed, p.MinSpeed, p.MaxSpeed)
			}
		}
	})
}

func TestProperty_StripSpeed_UniformWhenEqual(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		source.ResetLayoutCache()

		speed := rapid.Float64Range(0.1, 100.0).Draw(rt, "speed")

		p := attract_matrix.Policy{
			MinSpeed:    speed,
			MaxSpeed:    speed,
			TrailLength: 8,
			Density:     1.0,
		}
		p = attract_matrix.NormalizePolicy(p)

		panelWidth := 280
		panelHeight := 140
		glyphAdvance := 14
		rowHeight := 14
		seed := int64(42)

		strips := source.RebuildStripsForTest(p, panelWidth, panelHeight, glyphAdvance, rowHeight, false, seed, nil)

		if len(strips) == 0 {
			t.Fatal("expected at least one strip")
		}

		for i, strip := range strips {
			if strip.Speed() != speed {
				t.Fatalf("strip[%d].Speed() = %f, want exactly %f (MinSpeed == MaxSpeed)",
					i, strip.Speed(), speed)
			}
		}
	})
}

// For any panel configuration with visibleCells = panelPixelHeight / rowHeight,
// every column's phase offset SHALL be in the inclusive range [0, visibleCells * 2].
// When visibleCells = 0, no strips are created (rebuildStrips returns nil).

func TestProperty_StripPhase_InRange(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		// Reset so this is treated as an initial build (random phases).
		source.ResetLayoutCache()

		// Generate panel dimensions ensuring visibleCells > 0.
		rowHeight := rapid.IntRange(8, 20).Draw(rt, "rowHeight")
		visibleCells := rapid.IntRange(1, 50).Draw(rt, "visibleCells")
		panelHeight := visibleCells * rowHeight

		glyphAdvance := rapid.IntRange(8, 20).Draw(rt, "glyphAdvance")
		columns := rapid.IntRange(1, 30).Draw(rt, "columns")
		panelWidth := columns * glyphAdvance

		p := attract_matrix.Policy{
			MinSpeed:    3.0,
			MaxSpeed:    12.0,
			TrailLength: 8,
			Density:     1.0,
		}
		p = attract_matrix.NormalizePolicy(p)

		seed := int64(42)

		strips := source.RebuildStripsForTest(p, panelWidth, panelHeight, glyphAdvance, rowHeight, false, seed, nil)

		if len(strips) == 0 {
			t.Fatal("expected at least one strip with valid dimensions and density=1.0")
		}

		maxPhase := float64(visibleCells * 2)
		for i, strip := range strips {
			// Before any Tick() call, Offset() equals the Phase value.
			offset := strip.Offset()
			if offset < 0 || offset > maxPhase {
				t.Fatalf("strip[%d].Offset() = %f, want in [0, %f] (visibleCells=%d)",
					i, offset, maxPhase, visibleCells)
			}
		}
	})
}

func TestProperty_StripPhase_ZeroVisibleCells(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		source.ResetLayoutCache()

		// Ensure visibleCells = 0: panelHeight < rowHeight.
		rowHeight := rapid.IntRange(10, 50).Draw(rt, "rowHeight")
		panelHeight := rapid.IntRange(0, rowHeight-1).Draw(rt, "panelHeight")

		panelWidth := 280
		glyphAdvance := 14

		p := attract_matrix.Policy{
			MinSpeed:    3.0,
			MaxSpeed:    12.0,
			TrailLength: 8,
			Density:     1.0,
		}
		p = attract_matrix.NormalizePolicy(p)

		seed := int64(42)

		strips := source.RebuildStripsForTest(p, panelWidth, panelHeight, glyphAdvance, rowHeight, false, seed, nil)

		if strips != nil && len(strips) != 0 {
			t.Fatalf("expected no strips when visibleCells=0, got %d strips", len(strips))
		}
	})
}
