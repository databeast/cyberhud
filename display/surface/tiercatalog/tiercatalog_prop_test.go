package tiercatalog

import (
	"fmt"
	"testing"

	"github.com/databeast/cyberhud/display/surface/fonts"
	"pgregory.net/rapid"
)

// stubCandidate implements the Candidate interface for property testing.
type stubCandidate struct {
	id       string
	metrics  font.Metrics
	scalable bool
}

func (s stubCandidate) ID() string                   { return s.id }
func (s stubCandidate) MetricsAt(_ int) font.Metrics { return s.metrics }
func (s stubCandidate) IsScalable() bool             { return s.scalable }

func TestProperty2_BestFitOptimality(t *testing.T) {
	families := []string{"spleen", "terminus", "cozette", "testfont"}

	rapid.Check(t, func(t *rapid.T) {
		// Generate a random number of candidates (2-20).
		n := rapid.IntRange(2, 20).Draw(t, "numCandidates")

		// Build a list of unique candidates.
		candidates := make([]Candidate, n)
		usedIDs := map[string]bool{}
		for i := 0; i < n; i++ {
			height := rapid.IntRange(5, 48).Draw(t, fmt.Sprintf("height_%d", i))
			family := rapid.SampledFrom(families).Draw(t, fmt.Sprintf("family_%d", i))
			// Generate a unique ID by appending the index.
			id := fmt.Sprintf("%s-%dx%d-%d", family, height/2+3, height, i)
			for usedIDs[id] {
				id = fmt.Sprintf("%s-%dx%d-%d-dup", family, height/2+3, height, i)
			}
			usedIDs[id] = true

			candidates[i] = stubCandidate{
				id: id,
				metrics: font.Metrics{
					GlyphWidth:   height/2 + 3,
					GlyphHeight:  height,
					GlyphAdvance: height/2 + 3,
					RowHeight:    height + 2,
				},
				scalable: false,
			}
		}

		// Generate a random target pixel height.
		targetPx := rapid.IntRange(1, 50).Draw(t, "targetPx")

		// Call bestFit.
		selected := bestFit(candidates, targetPx)

		// Compute distance of the selected candidate.
		selectedDist := abs(selected.MetricsAt(0).GlyphHeight - targetPx)
		selectedHeight := selected.MetricsAt(0).GlyphHeight
		selectedPri := familyPriority(selected.ID())

		// For EVERY other candidate, verify the selection is optimal.
		for i, other := range candidates {
			if other.ID() == selected.ID() {
				continue
			}

			otherDist := abs(other.MetricsAt(0).GlyphHeight - targetPx)
			otherHeight := other.MetricsAt(0).GlyphHeight
			otherPri := familyPriority(other.ID())

			// The selected candidate must have distance <= every other candidate.
			if selectedDist > otherDist {
				t.Fatalf("candidate %d (%s, height=%d, dist=%d) is closer to target %d than selected (%s, height=%d, dist=%d)",
					i, other.ID(), otherHeight, otherDist, targetPx, selected.ID(), selectedHeight, selectedDist)
			}

			// When distances are equal, check tie-breaking rules.
			if selectedDist == otherDist {
				// Tie-break 1: prefer smaller GlyphHeight.
				if selectedHeight > otherHeight {
					t.Fatalf("equidistant candidate %d (%s, height=%d) has smaller height than selected (%s, height=%d) at target %d",
						i, other.ID(), otherHeight, selected.ID(), selectedHeight, targetPx)
				}

				// Tie-break 2: when heights also equal, prefer higher family priority.
				if selectedHeight == otherHeight {
					if selectedPri < otherPri {
						t.Fatalf("same-height candidate %d (%s, priority=%d) has higher priority than selected (%s, priority=%d) at target %d",
							i, other.ID(), otherPri, selected.ID(), selectedPri, targetPx)
					}

					// Tie-break 3: when priority also ties, prefer lexicographically smaller FontID.
					if selectedPri == otherPri {
						if selected.ID() > other.ID() {
							t.Fatalf("same-priority candidate %d (%s) has lexicographically smaller ID than selected (%s) at target %d",
								i, other.ID(), selected.ID(), targetPx)
						}
					}
				}
			}
		}
	})
}

func TestProperty3_TierFullLargest(t *testing.T) {
	families := []string{"spleen", "terminus", "cozette", "testfont"}

	rapid.Check(t, func(t *rapid.T) {
		// Generate a random number of candidates (2-20).
		n := rapid.IntRange(2, 20).Draw(t, "numCandidates")

		// Generate a random advance budget.
		advanceBudget := rapid.IntRange(4, 32).Draw(t, "advanceBudget")

		// Build a list of unique candidates with varying metrics.
		candidates := make([]Candidate, n)
		usedIDs := map[string]bool{}
		for i := 0; i < n; i++ {
			height := rapid.IntRange(5, 48).Draw(t, fmt.Sprintf("height_%d", i))
			advance := rapid.IntRange(3, 32).Draw(t, fmt.Sprintf("advance_%d", i))
			family := rapid.SampledFrom(families).Draw(t, fmt.Sprintf("family_%d", i))
			id := fmt.Sprintf("%s-%dx%d-%d", family, advance, height, i)
			for usedIDs[id] {
				id = fmt.Sprintf("%s-%dx%d-%d-dup", family, advance, height, i)
			}
			usedIDs[id] = true

			candidates[i] = stubCandidate{
				id: id,
				metrics: font.Metrics{
					GlyphWidth:   advance - 1,
					GlyphHeight:  height,
					GlyphAdvance: advance,
					RowHeight:    height + 2,
				},
				scalable: false,
			}
		}

		// Filter candidates by advance budget (only those with GlyphAdvance <= budget qualify).
		var qualifying []Candidate
		for _, c := range candidates {
			if c.MetricsAt(0).GlyphAdvance <= advanceBudget {
				qualifying = append(qualifying, c)
			}
		}

		// If no candidates qualify, skip this iteration.
		if len(qualifying) == 0 {
			return
		}

		// Call selectLargest on the qualifying candidates.
		selected := selectLargest(qualifying)
		selectedHeight := selected.MetricsAt(0).GlyphHeight
		selectedPri := familyPriority(selected.ID())

		// Assert: selected.GlyphHeight >= every other qualifying candidate's GlyphHeight.
		for i, other := range qualifying {
			if other.ID() == selected.ID() {
				continue
			}
			otherHeight := other.MetricsAt(0).GlyphHeight
			otherPri := familyPriority(other.ID())

			if selectedHeight < otherHeight {
				t.Fatalf("candidate %d (%s, height=%d) has larger GlyphHeight than selected (%s, height=%d)",
					i, other.ID(), otherHeight, selected.ID(), selectedHeight)
			}

			// For same-height candidates, verify tie-breaking.
			if selectedHeight == otherHeight {
				// Tie-break 1: higher family priority wins.
				if selectedPri < otherPri {
					t.Fatalf("same-height candidate %d (%s, priority=%d) has higher family priority than selected (%s, priority=%d)",
						i, other.ID(), otherPri, selected.ID(), selectedPri)
				}
				// Tie-break 2: when priority ties, smaller ID wins.
				if selectedPri == otherPri && selected.ID() > other.ID() {
					t.Fatalf("same-priority candidate %d (%s) has lexicographically smaller ID than selected (%s)",
						i, other.ID(), selected.ID())
				}
			}
		}
	})
}

// propTestFace implements font.Face for property testing purposes.
// Unlike stubCandidate (which implements the Candidate interface),
// propTestFace implements font.Face so it can be registered in the font registry.
type propTestFace struct {
	id      string
	metrics font.Metrics
}

func (f propTestFace) ID() string                    { return f.id }
func (f propTestFace) Metrics() font.Metrics         { return f.metrics }
func (f propTestFace) GlyphRow(_ rune, _ int) uint32 { return 0 }

func TestProperty4_WidthSafety(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// 1. Snapshot and clear the font registry for a clean state.
		restore := font.SnapshotAndClear()
		defer restore()

		// 2. Generate 1-20 random fonts and register them.
		numFonts := rapid.IntRange(1, 20).Draw(t, "numFonts")
		for i := 0; i < numFonts; i++ {
			glyphWidth := rapid.IntRange(3, 32).Draw(t, fmt.Sprintf("glyphWidth_%d", i))
			glyphHeight := rapid.IntRange(5, 48).Draw(t, fmt.Sprintf("glyphHeight_%d", i))
			glyphAdvance := rapid.IntRange(3, 32).Draw(t, fmt.Sprintf("glyphAdvance_%d", i))
			rowHeight := glyphHeight + rapid.IntRange(0, 4).Draw(t, fmt.Sprintf("rowExtra_%d", i))

			families := []string{"spleen", "terminus", "cozette", "testfont"}
			family := rapid.SampledFrom(families).Draw(t, fmt.Sprintf("family_%d", i))
			id := fmt.Sprintf("%s-%dx%d-%d", family, glyphWidth, glyphHeight, i)

			face := propTestFace{
				id: id,
				metrics: font.Metrics{
					GlyphWidth:   glyphWidth,
					GlyphHeight:  glyphHeight,
					GlyphAdvance: glyphAdvance,
					RowHeight:    rowHeight,
				},
			}
			font.Register(face)
		}

		// 3. Generate random PixelWidth, PixelHeight, MinChars.
		pixelWidth := rapid.IntRange(1, 4096).Draw(t, "pixelWidth")
		pixelHeight := rapid.IntRange(1, 4096).Draw(t, "pixelHeight")
		minChars := rapid.IntRange(1, 100).Draw(t, "minChars")

		// 4. Call Build.
		catalog, err := Build(Params{
			PixelWidth:  pixelWidth,
			PixelHeight: pixelHeight,
			MinChars:    minChars,
		})

		// 5. If Build returns an error, the property is vacuously true.
		if err != nil {
			return
		}

		// 6. For every tier in the catalog, assert width safety.
		effectiveMinChars := minChars
		if effectiveMinChars <= 0 {
			effectiveMinChars = 10
		}

		for _, tier := range catalog.Tiers() {
			entry, ok := catalog.Get(tier)
			if !ok {
				continue
			}
			if entry.GlyphAdvance*effectiveMinChars > pixelWidth {
				t.Fatalf("width safety violated: tier %s has GlyphAdvance=%d, effectiveMinChars=%d, product=%d > PixelWidth=%d (FontID=%s)",
					tier, entry.GlyphAdvance, effectiveMinChars, entry.GlyphAdvance*effectiveMinChars, pixelWidth, entry.FontID)
			}
		}
	})
}

func TestProperty5_TierMonotonicity(t *testing.T) {
	families := []string{"spleen", "terminus", "cozette", "testfont"}

	rapid.Check(t, func(t *rapid.T) {
		restore := font.SnapshotAndClear()
		defer restore()

		// Generate 1-20 random fonts and register them.
		n := rapid.IntRange(1, 20).Draw(t, "numFonts")
		for i := 0; i < n; i++ {
			glyphHeight := rapid.IntRange(5, 48).Draw(t, fmt.Sprintf("gh_%d", i))
			glyphWidth := rapid.IntRange(3, 32).Draw(t, fmt.Sprintf("gw_%d", i))
			glyphAdvance := rapid.IntRange(3, 32).Draw(t, fmt.Sprintf("ga_%d", i))
			rowHeight := glyphHeight + rapid.IntRange(0, 4).Draw(t, fmt.Sprintf("rh_%d", i))
			family := rapid.SampledFrom(families).Draw(t, fmt.Sprintf("fam_%d", i))
			id := fmt.Sprintf("%s-%dx%d-prop5-%d", family, glyphWidth, glyphHeight, i)

			font.Register(propTestFace{
				id: id,
				metrics: font.Metrics{
					GlyphWidth:   glyphWidth,
					GlyphHeight:  glyphHeight,
					GlyphAdvance: glyphAdvance,
					RowHeight:    rowHeight,
				},
			})
		}

		// Generate random panel dimensions and MinChars.
		pixelWidth := rapid.IntRange(1, 4096).Draw(t, "pixelWidth")
		pixelHeight := rapid.IntRange(1, 4096).Draw(t, "pixelHeight")
		minChars := rapid.IntRange(1, 100).Draw(t, "minChars")

		cat, err := Build(Params{
			PixelWidth:  pixelWidth,
			PixelHeight: pixelHeight,
			MinChars:    minChars,
		})
		if err != nil {
			return // vacuously true when build fails
		}

		// Walk tierOrder: for each consecutive pair, assert non-decreasing GlyphHeight.
		for i := 1; i < len(tierOrder); i++ {
			prev, okPrev := cat.Get(tierOrder[i-1])
			curr, okCurr := cat.Get(tierOrder[i])
			if !okPrev || !okCurr {
				t.Fatalf("tier %q or %q missing from catalog", tierOrder[i-1], tierOrder[i])
			}
			if prev.GlyphHeight > curr.GlyphHeight {
				t.Fatalf("monotonicity violated: tier %q (GlyphHeight=%d) > tier %q (GlyphHeight=%d)",
					tierOrder[i-1], prev.GlyphHeight, tierOrder[i], curr.GlyphHeight)
			}
		}
	})
}

func TestProperty6_FontIDValidity(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// 1. Snapshot and clear the font registry for a clean state.
		restore := font.SnapshotAndClear()
		defer restore()

		// 2. Generate 1-20 random fonts and register them.
		numFonts := rapid.IntRange(1, 20).Draw(t, "numFonts")
		for i := 0; i < numFonts; i++ {
			glyphWidth := rapid.IntRange(3, 32).Draw(t, fmt.Sprintf("gw_%d", i))
			glyphHeight := rapid.IntRange(5, 48).Draw(t, fmt.Sprintf("gh_%d", i))
			glyphAdvance := rapid.IntRange(3, 32).Draw(t, fmt.Sprintf("ga_%d", i))
			rowHeight := glyphHeight + rapid.IntRange(0, 4).Draw(t, fmt.Sprintf("rh_%d", i))

			families := []string{"spleen", "terminus", "cozette", "testfont"}
			family := rapid.SampledFrom(families).Draw(t, fmt.Sprintf("fam_%d", i))
			id := fmt.Sprintf("%s-%dx%d-prop6-%d", family, glyphWidth, glyphHeight, i)

			font.Register(propTestFace{
				id: id,
				metrics: font.Metrics{
					GlyphWidth:   glyphWidth,
					GlyphHeight:  glyphHeight,
					GlyphAdvance: glyphAdvance,
					RowHeight:    rowHeight,
				},
			})
		}

		// 3. Generate random PixelWidth, PixelHeight, MinChars.
		pixelWidth := rapid.IntRange(1, 4096).Draw(t, "pixelWidth")
		pixelHeight := rapid.IntRange(1, 4096).Draw(t, "pixelHeight")
		minChars := rapid.IntRange(1, 100).Draw(t, "minChars")

		// 4. Call Build.
		catalog, err := Build(Params{
			PixelWidth:  pixelWidth,
			PixelHeight: pixelHeight,
			MinChars:    minChars,
		})

		// 5. If Build returns an error, the property is vacuously true.
		if err != nil {
			return
		}

		// 6. For every tier in the catalog, assert FontID validity.
		for _, tier := range catalog.Tiers() {
			entry, ok := catalog.Get(tier)
			if !ok {
				continue
			}

			// FontID must be non-empty.
			if entry.FontID == "" {
				t.Fatalf("tier %q has empty FontID", tier)
			}

			// FontID must resolve in the font registry.
			_, found := font.Get(entry.FontID)
			if !found {
				t.Fatalf("tier %q has FontID=%q which is not resolvable in the font registry", tier, entry.FontID)
			}
		}
	})
}

func TestProperty7_Determinism(t *testing.T) {
	families := []string{"spleen", "terminus", "cozette", "testfont"}

	rapid.Check(t, func(t *rapid.T) {
		// 1. Snapshot and clear the font registry for a clean state.
		restore := font.SnapshotAndClear()
		defer restore()

		// 2. Register 1-20 random propTestFace fonts.
		numFonts := rapid.IntRange(1, 20).Draw(t, "numFonts")
		for i := 0; i < numFonts; i++ {
			glyphWidth := rapid.IntRange(3, 32).Draw(t, fmt.Sprintf("gw_%d", i))
			glyphHeight := rapid.IntRange(5, 48).Draw(t, fmt.Sprintf("gh_%d", i))
			glyphAdvance := rapid.IntRange(3, 32).Draw(t, fmt.Sprintf("ga_%d", i))
			rowHeight := glyphHeight + rapid.IntRange(0, 4).Draw(t, fmt.Sprintf("rh_%d", i))
			family := rapid.SampledFrom(families).Draw(t, fmt.Sprintf("fam_%d", i))
			id := fmt.Sprintf("%s-%dx%d-prop7-%d", family, glyphWidth, glyphHeight, i)

			font.Register(propTestFace{
				id: id,
				metrics: font.Metrics{
					GlyphWidth:   glyphWidth,
					GlyphHeight:  glyphHeight,
					GlyphAdvance: glyphAdvance,
					RowHeight:    rowHeight,
				},
			})
		}

		// 3. Generate random Params.
		pixelWidth := rapid.IntRange(1, 4096).Draw(t, "pixelWidth")
		pixelHeight := rapid.IntRange(1, 4096).Draw(t, "pixelHeight")
		minChars := rapid.IntRange(0, 100).Draw(t, "minChars")
		ppi := rapid.Float64Range(-10.0, 300.0).Draw(t, "ppi")

		params := Params{
			PixelWidth:  pixelWidth,
			PixelHeight: pixelHeight,
			MinChars:    minChars,
			PPI:         ppi,
		}

		// 4. Call Build twice with identical params.
		cat1, err1 := Build(params)
		cat2, err2 := Build(params)

		// 5. Assert both return same error status.
		if (err1 == nil) != (err2 == nil) {
			t.Fatalf("determinism violated: first Build error=%v, second Build error=%v", err1, err2)
		}

		// If both errored, that's deterministic — done.
		if err1 != nil {
			return
		}

		// 6. Assert identical entries for all tiers.
		for _, tier := range tierOrder {
			e1, ok1 := cat1.Get(tier)
			e2, ok2 := cat2.Get(tier)

			if ok1 != ok2 {
				t.Fatalf("determinism violated: tier %q present in first=%v, second=%v", tier, ok1, ok2)
			}
			if !ok1 {
				continue
			}

			if e1.GlyphWidth != e2.GlyphWidth {
				t.Fatalf("determinism violated: tier %q GlyphWidth first=%d, second=%d", tier, e1.GlyphWidth, e2.GlyphWidth)
			}
			if e1.GlyphHeight != e2.GlyphHeight {
				t.Fatalf("determinism violated: tier %q GlyphHeight first=%d, second=%d", tier, e1.GlyphHeight, e2.GlyphHeight)
			}
			if e1.GlyphAdvance != e2.GlyphAdvance {
				t.Fatalf("determinism violated: tier %q GlyphAdvance first=%d, second=%d", tier, e1.GlyphAdvance, e2.GlyphAdvance)
			}
			if e1.RowHeight != e2.RowHeight {
				t.Fatalf("determinism violated: tier %q RowHeight first=%d, second=%d", tier, e1.RowHeight, e2.RowHeight)
			}
			if e1.FontID != e2.FontID {
				t.Fatalf("determinism violated: tier %q FontID first=%q, second=%q", tier, e1.FontID, e2.FontID)
			}
		}
	})
}

func TestProperty8_SuccessGuarantee(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// 1. Snapshot and clear the font registry for a clean state.
		restore := font.SnapshotAndClear()
		defer restore()

		// 2. Generate 1-20 random fonts and register them.
		numFonts := rapid.IntRange(1, 20).Draw(t, "numFonts")
		families := []string{"spleen", "terminus", "cozette", "testfont"}
		type fontInfo struct {
			id      string
			advance int
		}
		var registered []fontInfo

		for i := 0; i < numFonts; i++ {
			glyphWidth := rapid.IntRange(3, 32).Draw(t, fmt.Sprintf("gw_%d", i))
			glyphHeight := rapid.IntRange(5, 48).Draw(t, fmt.Sprintf("gh_%d", i))
			glyphAdvance := rapid.IntRange(3, 32).Draw(t, fmt.Sprintf("ga_%d", i))
			rowHeight := glyphHeight + rapid.IntRange(0, 4).Draw(t, fmt.Sprintf("rh_%d", i))
			family := rapid.SampledFrom(families).Draw(t, fmt.Sprintf("fam_%d", i))
			id := fmt.Sprintf("%s-%dx%d-prop8-%d", family, glyphWidth, glyphHeight, i)

			font.Register(propTestFace{
				id: id,
				metrics: font.Metrics{
					GlyphWidth:   glyphWidth,
					GlyphHeight:  glyphHeight,
					GlyphAdvance: glyphAdvance,
					RowHeight:    rowHeight,
				},
			})
			registered = append(registered, fontInfo{id: id, advance: glyphAdvance})
		}

		// 3. Generate random PixelWidth (1-4096), PixelHeight (1-4096), MinChars (1-100).
		pixelWidth := rapid.IntRange(1, 4096).Draw(t, "pixelWidth")
		pixelHeight := rapid.IntRange(1, 4096).Draw(t, "pixelHeight")
		minChars := rapid.IntRange(1, 100).Draw(t, "minChars")

		// 4. Compute advanceBudget = PixelWidth / MinChars.
		advanceBudget := pixelWidth / minChars

		// 5. Check if ANY registered font has GlyphAdvance <= advanceBudget.
		//    If none do, skip (vacuously true).
		anyQualifies := false
		for _, f := range registered {
			if f.advance <= advanceBudget {
				anyQualifies = true
				break
			}
		}
		if !anyQualifies {
			return
		}

		// 6. Call Build. Assert: err == nil (success guaranteed when at least one font qualifies).
		catalog, err := Build(Params{
			PixelWidth:  pixelWidth,
			PixelHeight: pixelHeight,
			MinChars:    minChars,
		})
		if err != nil {
			t.Fatalf("Build returned error despite qualifying fonts: %v (pixelWidth=%d, minChars=%d, advanceBudget=%d)",
				err, pixelWidth, minChars, advanceBudget)
		}

		// 7. Assert: all tiers in tierOrder have a valid entry (ok=true from Get).
		for _, tier := range tierOrder {
			_, ok := catalog.Get(tier)
			if !ok {
				t.Fatalf("tier %q is missing from catalog despite successful build", tier)
			}
		}
	})
}

func TestProperty11_NegativePPIEquivalence(t *testing.T) {
	families := []string{"spleen", "terminus", "cozette", "testfont"}

	rapid.Check(t, func(t *rapid.T) {
		// 1. Snapshot and clear the font registry for a clean state.
		restore := font.SnapshotAndClear()
		defer restore()

		// 2. Register 1-20 random fonts.
		numFonts := rapid.IntRange(1, 20).Draw(t, "numFonts")
		for i := 0; i < numFonts; i++ {
			glyphWidth := rapid.IntRange(3, 32).Draw(t, fmt.Sprintf("gw_%d", i))
			glyphHeight := rapid.IntRange(5, 48).Draw(t, fmt.Sprintf("gh_%d", i))
			glyphAdvance := rapid.IntRange(3, 32).Draw(t, fmt.Sprintf("ga_%d", i))
			rowHeight := glyphHeight + rapid.IntRange(0, 4).Draw(t, fmt.Sprintf("rh_%d", i))
			family := rapid.SampledFrom(families).Draw(t, fmt.Sprintf("fam_%d", i))
			id := fmt.Sprintf("%s-%dx%d-prop11-%d", family, glyphWidth, glyphHeight, i)

			font.Register(propTestFace{
				id: id,
				metrics: font.Metrics{
					GlyphWidth:   glyphWidth,
					GlyphHeight:  glyphHeight,
					GlyphAdvance: glyphAdvance,
					RowHeight:    rowHeight,
				},
			})
		}

		// 3. Generate random Params.
		pixelWidth := rapid.IntRange(1, 4096).Draw(t, "pixelWidth")
		pixelHeight := rapid.IntRange(1, 4096).Draw(t, "pixelHeight")
		minChars := rapid.IntRange(1, 100).Draw(t, "minChars")

		// 4. Generate a negative PPI (range -100.0 to -0.1).
		negativePPI := -rapid.Float64Range(0.1, 100.0).Draw(t, "negativePPI")

		// 5. Build with the negative PPI → catalog1.
		catalog1, err1 := Build(Params{
			PixelWidth:  pixelWidth,
			PixelHeight: pixelHeight,
			MinChars:    minChars,
			PPI:         negativePPI,
		})

		// 6. Build with PPI=0 (same other params) → catalog2.
		catalog2, err2 := Build(Params{
			PixelWidth:  pixelWidth,
			PixelHeight: pixelHeight,
			MinChars:    minChars,
			PPI:         0,
		})

		// 7. Assert: both produce same error state.
		if (err1 == nil) != (err2 == nil) {
			t.Fatalf("error state mismatch: negativePPI=%f err=%v, PPI=0 err=%v",
				negativePPI, err1, err2)
		}

		// If both errored, the property holds (same error state).
		if err1 != nil {
			return
		}

		// 8. Assert: if successful, same entries for all tiers.
		for _, tier := range catalog1.Tiers() {
			entry1, ok1 := catalog1.Get(tier)
			entry2, ok2 := catalog2.Get(tier)

			if ok1 != ok2 {
				t.Fatalf("tier %q presence mismatch: negativePPI=%f present=%v, PPI=0 present=%v",
					tier, negativePPI, ok1, ok2)
			}
			if !ok1 {
				continue
			}

			if entry1 != entry2 {
				t.Fatalf("tier %q entries differ: negativePPI=%f entry=%+v, PPI=0 entry=%+v",
					tier, negativePPI, entry1, entry2)
			}
		}
	})
}

func TestProperty10_NoPanics(t *testing.T) {
	families := []string{"spleen", "terminus", "cozette", "testfont"}

	rapid.Check(t, func(t *rapid.T) {
		// 1. Snapshot and clear the font registry for a clean state.
		restore := font.SnapshotAndClear()
		defer restore()

		// 2. Optionally register 0-20 fonts (including zero to test empty registry).
		numFonts := rapid.IntRange(0, 20).Draw(t, "numFonts")
		for i := 0; i < numFonts; i++ {
			glyphWidth := rapid.IntRange(3, 32).Draw(t, fmt.Sprintf("gw_%d", i))
			glyphHeight := rapid.IntRange(5, 48).Draw(t, fmt.Sprintf("gh_%d", i))
			glyphAdvance := rapid.IntRange(3, 32).Draw(t, fmt.Sprintf("ga_%d", i))
			rowHeight := glyphHeight + rapid.IntRange(0, 4).Draw(t, fmt.Sprintf("rh_%d", i))
			family := rapid.SampledFrom(families).Draw(t, fmt.Sprintf("fam_%d", i))
			id := fmt.Sprintf("%s-%dx%d-prop10-%d", family, glyphWidth, glyphHeight, i)

			font.Register(propTestFace{
				id: id,
				metrics: font.Metrics{
					GlyphWidth:   glyphWidth,
					GlyphHeight:  glyphHeight,
					GlyphAdvance: glyphAdvance,
					RowHeight:    rowHeight,
				},
			})
		}

		// 3. Generate extreme Params.
		pixelWidth := rapid.IntRange(-100, 4096).Draw(t, "pixelWidth")
		pixelHeight := rapid.IntRange(-100, 4096).Draw(t, "pixelHeight")
		minChars := rapid.IntRange(-50, 100).Draw(t, "minChars")
		ppi := rapid.Float64Range(-100.0, 500.0).Draw(t, "ppi")

		// 4. Call Build — if it panics, the test framework catches it as a failure.
		// We don't check the result; the property is simply "doesn't panic".
		Build(Params{
			PixelWidth:  pixelWidth,
			PixelHeight: pixelHeight,
			MinChars:    minChars,
			PPI:         ppi,
		})
	})
}

func TestProperty9_ErrorOnInsufficientFonts(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		restore := font.SnapshotAndClear()
		defer restore()

		// Draw which sub-case to test: 0 = empty registry, 1 = no font fits.
		subCase := rapid.IntRange(0, 1).Draw(t, "subCase")

		switch subCase {
		case 0:
			// (a) Empty registry: no fonts registered at all.
			pixelWidth := rapid.IntRange(1, 4096).Draw(t, "pixelWidth")
			pixelHeight := rapid.IntRange(1, 4096).Draw(t, "pixelHeight")
			minChars := rapid.IntRange(1, 100).Draw(t, "minChars")

			_, err := Build(Params{
				PixelWidth:  pixelWidth,
				PixelHeight: pixelHeight,
				MinChars:    minChars,
			})
			if err == nil {
				t.Fatal("expected error when no fonts are registered, got nil")
			}

		case 1:
			// (b) No font fits: all registered fonts have advance >= 20,
			// and the advance budget is forced below 20.
			numFonts := rapid.IntRange(1, 10).Draw(t, "numFonts")
			for i := 0; i < numFonts; i++ {
				advance := rapid.IntRange(20, 40).Draw(t, fmt.Sprintf("advance_%d", i))
				glyphHeight := rapid.IntRange(5, 48).Draw(t, fmt.Sprintf("gh_%d", i))
				glyphWidth := advance - 1
				rowHeight := glyphHeight + 2
				id := fmt.Sprintf("testfont-%dx%d-prop9-%d", glyphWidth, glyphHeight, i)

				font.Register(propTestFace{
					id: id,
					metrics: font.Metrics{
						GlyphWidth:   glyphWidth,
						GlyphHeight:  glyphHeight,
						GlyphAdvance: advance,
						RowHeight:    rowHeight,
					},
				})
			}

			// Choose PixelWidth and MinChars so advanceBudget = PixelWidth/MinChars < 20.
			// MinChars in [1,100], PixelWidth in [1, minChars*19] ensures budget <= 19.
			minChars := rapid.IntRange(1, 100).Draw(t, "minChars")
			maxPixelWidth := minChars * 19 // ensures PixelWidth / minChars <= 19 < 20
			if maxPixelWidth < 1 {
				maxPixelWidth = 1
			}
			pixelWidth := rapid.IntRange(1, maxPixelWidth).Draw(t, "pixelWidth")
			pixelHeight := rapid.IntRange(1, 4096).Draw(t, "pixelHeight")

			_, err := Build(Params{
				PixelWidth:  pixelWidth,
				PixelHeight: pixelHeight,
				MinChars:    minChars,
			})
			if err == nil {
				t.Fatalf("expected error when no font fits: pixelWidth=%d, minChars=%d, advanceBudget=%d (all fonts have advance >= 20)",
					pixelWidth, minChars, pixelWidth/minChars)
			}
		}
	})
}

func TestProperty12_MinCharsDefault(t *testing.T) {
	families := []string{"spleen", "terminus", "cozette", "testfont"}

	rapid.Check(t, func(t *rapid.T) {
		// 1. Snapshot and clear the font registry for a clean state.
		restore := font.SnapshotAndClear()
		defer restore()

		// 2. Register 1-20 random fonts.
		numFonts := rapid.IntRange(1, 20).Draw(t, "numFonts")
		for i := 0; i < numFonts; i++ {
			glyphWidth := rapid.IntRange(3, 32).Draw(t, fmt.Sprintf("gw_%d", i))
			glyphHeight := rapid.IntRange(5, 48).Draw(t, fmt.Sprintf("gh_%d", i))
			glyphAdvance := rapid.IntRange(3, 32).Draw(t, fmt.Sprintf("ga_%d", i))
			rowHeight := glyphHeight + rapid.IntRange(0, 4).Draw(t, fmt.Sprintf("rh_%d", i))
			family := rapid.SampledFrom(families).Draw(t, fmt.Sprintf("fam_%d", i))
			id := fmt.Sprintf("%s-%dx%d-prop12-%d", family, glyphWidth, glyphHeight, i)

			font.Register(propTestFace{
				id: id,
				metrics: font.Metrics{
					GlyphWidth:   glyphWidth,
					GlyphHeight:  glyphHeight,
					GlyphAdvance: glyphAdvance,
					RowHeight:    rowHeight,
				},
			})
		}

		// 3. Generate random PixelWidth, PixelHeight, PPI.
		pixelWidth := rapid.IntRange(1, 4096).Draw(t, "pixelWidth")
		pixelHeight := rapid.IntRange(1, 4096).Draw(t, "pixelHeight")
		ppi := rapid.Float64Range(-10.0, 300.0).Draw(t, "ppi")

		// 4. Build with MinChars=0 → catalog1.
		catalog1, err1 := Build(Params{
			PixelWidth:  pixelWidth,
			PixelHeight: pixelHeight,
			MinChars:    0,
			PPI:         ppi,
		})

		// 5. Build with MinChars=10 → catalog2.
		catalog2, err2 := Build(Params{
			PixelWidth:  pixelWidth,
			PixelHeight: pixelHeight,
			MinChars:    10,
			PPI:         ppi,
		})

		// 6. Assert: both produce same error status.
		if (err1 == nil) != (err2 == nil) {
			t.Fatalf("MinChars default violated: MinChars=0 error=%v, MinChars=10 error=%v",
				err1, err2)
		}

		// If both errored, the property holds.
		if err1 != nil {
			return
		}

		// 7. Assert: same entries for all tiers.
		for _, tier := range tierOrder {
			e1, ok1 := catalog1.Get(tier)
			e2, ok2 := catalog2.Get(tier)

			if ok1 != ok2 {
				t.Fatalf("MinChars default violated: tier %q present with MinChars=0: %v, MinChars=10: %v",
					tier, ok1, ok2)
			}
			if !ok1 {
				continue
			}

			if e1.GlyphWidth != e2.GlyphWidth ||
				e1.GlyphHeight != e2.GlyphHeight ||
				e1.GlyphAdvance != e2.GlyphAdvance ||
				e1.RowHeight != e2.RowHeight ||
				e1.FontID != e2.FontID {
				t.Fatalf("MinChars default violated: tier %q entry differs\n  MinChars=0:  %+v\n  MinChars=10: %+v",
					tier, e1, e2)
			}
		}
	})
}

func TestProperty14_PositiveEntryMetrics(t *testing.T) {
	families := []string{"spleen", "terminus", "cozette", "testfont"}

	rapid.Check(t, func(t *rapid.T) {
		// 1. Snapshot and clear the font registry for a clean state.
		restore := font.SnapshotAndClear()
		defer restore()

		// 2. Generate 1-20 random fonts with ALL positive metrics and register them.
		numFonts := rapid.IntRange(1, 20).Draw(t, "numFonts")
		for i := 0; i < numFonts; i++ {
			glyphWidth := rapid.IntRange(1, 32).Draw(t, fmt.Sprintf("gw_%d", i))
			glyphHeight := rapid.IntRange(1, 48).Draw(t, fmt.Sprintf("gh_%d", i))
			glyphAdvance := rapid.IntRange(1, 32).Draw(t, fmt.Sprintf("ga_%d", i))
			rowHeight := rapid.IntRange(1, 52).Draw(t, fmt.Sprintf("rh_%d", i))
			family := rapid.SampledFrom(families).Draw(t, fmt.Sprintf("fam_%d", i))
			id := fmt.Sprintf("%s-%dx%d-prop14-%d", family, glyphWidth, glyphHeight, i)

			font.Register(propTestFace{
				id: id,
				metrics: font.Metrics{
					GlyphWidth:   glyphWidth,
					GlyphHeight:  glyphHeight,
					GlyphAdvance: glyphAdvance,
					RowHeight:    rowHeight,
				},
			})
		}

		// 3. Generate random PixelWidth, PixelHeight, MinChars.
		pixelWidth := rapid.IntRange(1, 4096).Draw(t, "pixelWidth")
		pixelHeight := rapid.IntRange(1, 4096).Draw(t, "pixelHeight")
		minChars := rapid.IntRange(1, 100).Draw(t, "minChars")

		// 4. Call Build.
		catalog, err := Build(Params{
			PixelWidth:  pixelWidth,
			PixelHeight: pixelHeight,
			MinChars:    minChars,
		})

		// 5. If Build returns an error, the property is vacuously true.
		if err != nil {
			return
		}

		// 6. For every tier in the catalog, assert all Entry metrics are positive.
		for _, tier := range catalog.Tiers() {
			entry, ok := catalog.Get(tier)
			if !ok {
				continue
			}
			if entry.GlyphWidth <= 0 {
				t.Fatalf("tier %q has non-positive GlyphWidth=%d (FontID=%s)",
					tier, entry.GlyphWidth, entry.FontID)
			}
			if entry.GlyphHeight <= 0 {
				t.Fatalf("tier %q has non-positive GlyphHeight=%d (FontID=%s)",
					tier, entry.GlyphHeight, entry.FontID)
			}
			if entry.GlyphAdvance <= 0 {
				t.Fatalf("tier %q has non-positive GlyphAdvance=%d (FontID=%s)",
					tier, entry.GlyphAdvance, entry.FontID)
			}
			if entry.RowHeight <= 0 {
				t.Fatalf("tier %q has non-positive RowHeight=%d (FontID=%s)",
					tier, entry.RowHeight, entry.FontID)
			}
		}
	})
}

func TestProperty13_TierIndependence(t *testing.T) {
	families := []string{"spleen", "terminus", "cozette", "testfont"}
	// Tiers eligible for custom mm override (TierFull has no mm target).
	mmTiers := []Tier{TierSmall, TierNormal, TierLarge, TierHuge, TierColossal}

	rapid.Check(t, func(t *rapid.T) {
		// 1. Snapshot and clear the font registry for a clean state.
		restore := font.SnapshotAndClear()
		defer restore()

		// 2. Register 2-20 random fonts.
		numFonts := rapid.IntRange(2, 20).Draw(t, "numFonts")
		for i := 0; i < numFonts; i++ {
			glyphWidth := rapid.IntRange(3, 32).Draw(t, fmt.Sprintf("gw_%d", i))
			glyphHeight := rapid.IntRange(5, 48).Draw(t, fmt.Sprintf("gh_%d", i))
			glyphAdvance := rapid.IntRange(3, 32).Draw(t, fmt.Sprintf("ga_%d", i))
			rowHeight := glyphHeight + rapid.IntRange(0, 4).Draw(t, fmt.Sprintf("rh_%d", i))
			family := rapid.SampledFrom(families).Draw(t, fmt.Sprintf("fam_%d", i))
			id := fmt.Sprintf("%s-%dx%d-prop13-%d", family, glyphWidth, glyphHeight, i)

			font.Register(propTestFace{
				id: id,
				metrics: font.Metrics{
					GlyphWidth:   glyphWidth,
					GlyphHeight:  glyphHeight,
					GlyphAdvance: glyphAdvance,
					RowHeight:    rowHeight,
				},
			})
		}

		// 3. Generate random Params with PPI > 0 (so mm targets are used).
		pixelWidth := rapid.IntRange(1, 4096).Draw(t, "pixelWidth")
		pixelHeight := rapid.IntRange(1, 4096).Draw(t, "pixelHeight")
		minChars := rapid.IntRange(1, 100).Draw(t, "minChars")
		ppi := rapid.Float64Range(10.0, 300.0).Draw(t, "ppi")

		// 4. Pick a random tier T (not TierFull) to modify.
		tierT := rapid.SampledFrom(mmTiers).Draw(t, "tierT")

		// 5. Generate a custom mm target for tier T.
		customMM := rapid.Float64Range(1.0, 30.0).Draw(t, "customMM")

		// 6. Build candidate pool with default params (same for both builds).
		baseParams := Params{
			PixelWidth:  pixelWidth,
			PixelHeight: pixelHeight,
			MinChars:    minChars,
			PPI:         ppi,
		}

		candidates, _, err := buildCandidatePool(baseParams)
		if err != nil {
			// No qualifying fonts — property is vacuously true.
			return
		}

		// 7. Resolve pixel targets with default targets (no custom override).
		params1 := Params{
			PixelWidth:  pixelWidth,
			PixelHeight: pixelHeight,
			MinChars:    minChars,
			PPI:         ppi,
		}
		targets1 := resolvePixelTargets(params1)

		// 8. Resolve pixel targets with a custom mm target for tier T only.
		params2 := Params{
			PixelWidth:  pixelWidth,
			PixelHeight: pixelHeight,
			MinChars:    minChars,
			PPI:         ppi,
			TierTargetsMM: map[Tier]float64{
				tierT: customMM,
			},
		}
		targets2 := resolvePixelTargets(params2)

		// 9. For every tier OTHER than T: bestFit must select the same candidate.
		for _, tier := range mmTiers {
			if tier == tierT {
				continue
			}

			selected1 := bestFit(candidates, targets1[tier])
			selected2 := bestFit(candidates, targets2[tier])

			if selected1.ID() != selected2.ID() {
				t.Fatalf("tier independence violated: changing tier %q's mm target affected tier %q: "+
					"default target selected %q (target=%d), custom target for %q selected %q (target=%d)",
					tierT, tier, selected1.ID(), targets1[tier], tierT, selected2.ID(), targets2[tier])
			}
		}
	})
}
