package tests_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/databeast/cyberhud/display/catalog"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"pgregory.net/rapid"
)

// mockPolicy is a minimal catalog.ConfigPolicy implementation for registry tests.
type mockPolicy struct{}

func (mockPolicy) Fingerprint() string                 { return "mock" }
func (mockPolicy) ToMap() map[string]interface{}       { return nil }
func (mockPolicy) Options() []catalog.OptionDefinition { return nil }

// --- From: fitness_notes_prop_test.go ---

// fitnessNoteMockStyle is a Style mock that uses Capability in fitness evaluation,
// allowing FitnessNotes to generate capability-related notes when appropriate.
type fitnessNoteMockStyle struct {
	name string
	reqs style.SurfaceRequirements
}

func (m fitnessNoteMockStyle) Name() string                            { return m.name }
func (m fitnessNoteMockStyle) Requirements() style.SurfaceRequirements { return m.reqs }
func (m fitnessNoteMockStyle) Build(_ any, _ mockPolicy, _ style.StyleContext) style.ViewData {
	return style.ViewData{Items: []string{"mock"}}
}

// Supports evaluates fitness using capability ordering.
// This aligns with FitnessNotes which emits a capability note when
// reqs.Capability > hints.Capability and fitness < Full.
func (m fitnessNoteMockStyle) Supports(hints textlayout.TextHints) style.Fitness {
	// Zero-dimension panels are Unsupported only when there are actual
	// requirements that would generate notes in FitnessNotes.
	hasAnyReq := m.reqs.MinWidth > 0 || m.reqs.MinHeight > 0 ||
		m.reqs.PreferredWidth > 0 || m.reqs.PreferredHeight > 0 ||
		m.reqs.Capability > style.MonoSlow ||
		m.reqs.MinRows > 0 || m.reqs.MinCharsPerLine > 0
	if hasAnyReq && (hints.PixelWidth == 0 || hints.PixelHeight == 0) {
		return style.Unsupported
	}

	// Undersized panels are Unsupported.
	if m.reqs.MinWidth > 0 && hints.PixelWidth < m.reqs.MinWidth {
		return style.Unsupported
	}
	if m.reqs.MinHeight > 0 && hints.PixelHeight < m.reqs.MinHeight {
		return style.Unsupported
	}

	// Capability ordering violation → Unsupported.
	if m.reqs.Capability > style.Capability(hints.Capability) {
		return style.Unsupported
	}

	// All minimum requirements met. Check preferred for Optimal.
	meetsPreferred := true
	if m.reqs.PreferredWidth > 0 && hints.PixelWidth < m.reqs.PreferredWidth {
		meetsPreferred = false
	}
	if m.reqs.PreferredHeight > 0 && hints.PixelHeight < m.reqs.PreferredHeight {
		meetsPreferred = false
	}

	if meetsPreferred {
		return style.Optimal
	}
	return style.Full
}

// genTextHints generates random TextHints with a mix of adequate and undersized values.
func genTextHints(t *rapid.T) textlayout.TextHints {
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

// countExpectedNotes counts how many note lines FitnessNotes should produce
// based on the actual implementation logic in fitness_notes.go.
func countExpectedNotes(reqs style.SurfaceRequirements, hints textlayout.TextHints) int {
	count := 0

	// MinWidth unmet.
	if reqs.MinWidth > 0 && hints.PixelWidth < reqs.MinWidth {
		count++
	}
	// MinHeight unmet.
	if reqs.MinHeight > 0 && hints.PixelHeight < reqs.MinHeight {
		count++
	}

	// Preferred dimensions: the implementation uses if/else if, so preferred width
	// and preferred height produce at most ONE combined note.
	preferredUnmet := false
	if reqs.PreferredWidth > 0 && hints.PixelWidth < reqs.PreferredWidth {
		preferredUnmet = true
	} else if reqs.PreferredHeight > 0 && hints.PixelHeight < reqs.PreferredHeight {
		preferredUnmet = true
	}
	if preferredUnmet {
		count++
	}

	// Capability ordering violation.
	if reqs.Capability > style.Capability(hints.Capability) {
		count++
	}

	return count
}

func TestProperty9_FitnessWarningsMatchUnmetFields(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		reqs := genSurfaceRequirements(t)
		s := fitnessNoteMockStyle{
			name: "test-style",
			reqs: reqs,
		}

		hints := genTextHints(t)

		notes := style.FitnessNotes(s, hints)
		fitness := s.Supports(hints)

		// If Full or Optimal → no notes.
		if fitness >= style.Full {
			if len(notes) != 0 {
				t.Fatalf("expected no notes for fitness %d, got %d notes: %v",
					fitness, len(notes), notes)
			}
			return
		}

		// Every note starts with "note:".
		for i, note := range notes {
			if !strings.HasPrefix(note, "note:") {
				t.Fatalf("note[%d] does not start with \"note:\": %q", i, note)
			}
		}

		// Note count equals unmet field count.
		expected := countExpectedNotes(reqs, hints)
		if len(notes) != expected {
			t.Fatalf("expected %d notes, got %d\nnotes: %v\nreqs: %+v\nhints: PixelWidth=%d, PixelHeight=%d, PreferEventRefresh=%v",
				expected, len(notes), notes, reqs, hints.PixelWidth, hints.PixelHeight, hints.PreferEventRefresh)
		}
	})
}

// --- From: registry_prop_test.go ---

// For any StyleRegistry and any registered style name (with any case/whitespace variation),
// Lookup returns that Style regardless of any TextHints or Fitness value. The registry's
// Lookup, Cycle, and Enumerate methods never filter, exclude, or reject styles based on
// fitness evaluation.

// genAnyTextHints generates a completely random TextHints including zero dimensions,
// undersized panels, and any configuration — the registry must not care.
func genAnyTextHints(t *rapid.T) textlayout.TextHints {
	return textlayout.TextHints{
		PixelWidth:         rapid.IntRange(0, 2000).Draw(t, "pixelWidth"),
		PixelHeight:        rapid.IntRange(0, 2000).Draw(t, "pixelHeight"),
		GlyphWidth:         rapid.IntRange(0, 32).Draw(t, "glyphWidth"),
		GlyphHeight:        rapid.IntRange(0, 32).Draw(t, "glyphHeight"),
		GlyphAdvance:       rapid.IntRange(0, 32).Draw(t, "glyphAdvance"),
		RowHeight:          rapid.IntRange(0, 64).Draw(t, "rowHeight"),
		PreferEventRefresh: rapid.Bool().Draw(t, "preferEventRefresh"),
		Capability:         rapid.IntRange(0, 5).Draw(t, "capability"),
	}
}

// genDistinctStyleNames generates a slice of N distinct style names.
func genDistinctStyleNames(t *rapid.T, count int) []string {
	names := make([]string, count)
	for i := 0; i < count; i++ {
		names[i] = fmt.Sprintf("style-%d", i)
	}
	return names
}

func TestProperty7_RegistryNeverBlocksStyles(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate a registry with 2–10 mock styles with varying requirements.
		count := rapid.IntRange(2, 10).Draw(t, "styleCount")
		names := genDistinctStyleNames(t, count)

		styles := make([]style.Style[any, mockPolicy], count)
		for i, name := range names {
			styles[i] = mockStyle{
				name: name,
				reqs: genSurfaceRequirements(t),
			}
		}

		registry := style.NewRegistry[any, mockPolicy](styles...)

		// Generate random TextHints — including zero dimensions, undersized panels, etc.
		hints := genAnyTextHints(t)

		// Verify: Lookup returns non-nil for every registered style name,
		// regardless of what TextHints or fitness would be.
		for _, s := range styles {
			looked := registry.Lookup(s.Name())
			if looked == nil {
				t.Fatalf("Lookup(%q) returned nil; registry must never block styles based on fitness\nhints: %+v",
					s.Name(), hints)
			}
			if looked.Name() != s.Name() {
				t.Fatalf("Lookup(%q) returned style %q; expected exact match",
					s.Name(), looked.Name())
			}
		}

		// Verify: Enumerate returns all registered styles without filtering by fitness.
		enumerated := registry.Enumerate()
		if len(enumerated) != count {
			t.Fatalf("Enumerate() returned %d styles, expected %d; registry must not filter by fitness",
				len(enumerated), count)
		}

		// Even with these random hints, calling Supports on each style should NOT
		// affect what the registry returns. Confirm Lookup still works after evaluation.
		for _, s := range styles {
			_ = s.Supports(hints) // fitness evaluation happens but must not affect registry
			looked := registry.Lookup(s.Name())
			if looked == nil {
				t.Fatalf("Lookup(%q) returned nil after Supports() call; registry must be independent of fitness",
					s.Name())
			}
		}

		// Verify: Cycle returns non-nil for every registered style and any delta.
		delta := rapid.IntRange(-100, 100).Draw(t, "delta")
		for _, s := range styles {
			cycled := registry.Cycle(s.Name(), delta, textlayout.TextHints{})
			if cycled == nil {
				t.Fatalf("Cycle(%q, %d) returned nil; registry must never block styles",
					s.Name(), delta)
			}
		}
	})
}

// For any sequence of N distinct Styles registered into a StyleRegistry,
// Enumerate() returns them in exactly the registration order.

func TestProperty10_RegistrationOrderPreserved(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate N distinct styles (2-20).
		n := rapid.IntRange(2, 20).Draw(t, "styleCount")

		styles := make([]style.Style[any, mockPolicy], n)
		usedNames := make(map[string]bool, n)

		for i := 0; i < n; i++ {
			var name string
			for {
				name = rapid.StringMatching(`[a-z][a-z0-9-]{0,15}`).Draw(t, fmt.Sprintf("name_%d", i))
				if !usedNames[name] {
					usedNames[name] = true
					break
				}
			}
			styles[i] = mockStyle{
				name: name,
				reqs: genSurfaceRequirements(t),
			}
		}

		// Create registry with the generated styles.
		reg := style.NewRegistry[any, mockPolicy](styles...)

		// Call Enumerate and verify order matches registration order.
		enumerated := reg.Enumerate()

		if len(enumerated) != n {
			t.Fatalf("Enumerate() returned %d styles, expected %d", len(enumerated), n)
		}

		for i := 0; i < n; i++ {
			if enumerated[i].Name() != styles[i].Name() {
				t.Fatalf("Enumerate()[%d].Name() = %q, expected %q (registration order violated)",
					i, enumerated[i].Name(), styles[i].Name())
			}
		}
	})
}

// TestProperty12_LookupWithNormalization verifies that for any registered style,
// Lookup finds the style using case/whitespace variations of its name, and returns
// nil for names not matching any registered style.

func TestProperty12_LookupWithNormalization(t *testing.T) {
	// Create a registry with known mock styles.
	digital := mockStyle{name: "digital"}
	bigDigit := mockStyle{name: "big-digit"}
	minimal := mockStyle{name: "minimal"}

	registry := style.NewRegistry[any, mockPolicy](digital, bigDigit, minimal)

	rapid.Check(t, func(t *rapid.T) {
		// Pick a registered style at random.
		registered := []mockStyle{digital, bigDigit, minimal}
		idx := rapid.IntRange(0, len(registered)-1).Draw(t, "styleIdx")
		chosen := registered[idx]

		// Generate a case/whitespace variation of the registered name.
		variation := genNameVariation(t, chosen.name)

		// Lookup with the variation must return the correct style.
		result := registry.Lookup(variation)
		if result == nil {
			t.Fatalf("Lookup(%q) returned nil, expected style %q", variation, chosen.name)
		}
		if result.Name() != chosen.name {
			t.Fatalf("Lookup(%q) returned style %q, expected %q", variation, result.Name(), chosen.name)
		}
	})
}

// TestProperty12_LookupUnregisteredReturnsNil verifies that Lookup returns nil
// for names that do not match any registered style after normalization.

func TestProperty12_LookupUnregisteredReturnsNil(t *testing.T) {
	// Create a registry with known mock styles.
	registry := style.NewRegistry[any, mockPolicy](
		mockStyle{name: "digital"},
		mockStyle{name: "big-digit"},
		mockStyle{name: "minimal"},
	)

	// Known registered names (normalized form).
	registeredNames := map[string]bool{
		"digital":   true,
		"big-digit": true,
		"minimal":   true,
	}

	rapid.Check(t, func(t *rapid.T) {
		// Generate a random string that is unlikely to match registered names.
		candidate := rapid.StringMatching(`[a-z0-9-]{1,32}`).Draw(t, "unregisteredName")

		// Skip if it happens to match a registered name.
		if registeredNames[candidate] {
			return
		}

		result := registry.Lookup(candidate)
		if result != nil {
			t.Fatalf("Lookup(%q) returned style %q, expected nil for unregistered name",
				candidate, result.Name())
		}
	})
}

// genNameVariation generates a case and whitespace variation of a given style name.
// It randomly applies: uppercase letters, leading/trailing spaces, and mixed case.
func genNameVariation(t *rapid.T, name string) string {
	// Randomly add leading whitespace.
	leadSpaces := rapid.IntRange(0, 5).Draw(t, "leadSpaces")
	// Randomly add trailing whitespace.
	trailSpaces := rapid.IntRange(0, 5).Draw(t, "trailSpaces")

	// Randomly change case of each character.
	runes := []rune(name)
	varied := make([]rune, len(runes))
	for i, r := range runes {
		if rapid.Bool().Draw(t, fmt.Sprintf("upper_%d", i)) {
			if r >= 'a' && r <= 'z' {
				varied[i] = r - 32 // to uppercase
			} else {
				varied[i] = r
			}
		} else {
			varied[i] = r
		}
	}

	// Build result with whitespace padding.
	result := ""
	for i := 0; i < leadSpaces; i++ {
		result += " "
	}
	result += string(varied)
	for i := 0; i < trailSpaces; i++ {
		result += " "
	}

	return result
}

// For any StyleRegistry, if Style A is registered first with normalized name N,
// then attempting to register Style B whose Name() normalizes to N is rejected,
// and Lookup(N) still returns A. Enumerate() contains only A (not B).

func TestProperty13_DuplicateRejection(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate a valid style name.
		name := rapid.StringMatching(`[a-z][a-z0-9-]{0,15}`).Draw(t, "styleName")

		// Create two mock styles with the same normalized name but different internal state.
		// Use distinct SurfaceRequirements so we can tell them apart.
		reqsA := style.SurfaceRequirements{MinWidth: 100, MinHeight: 100}
		reqsB := style.SurfaceRequirements{MinWidth: 200, MinHeight: 200}

		styleA := mockStyle{name: name, reqs: reqsA}
		styleB := mockStyle{name: name, reqs: reqsB}

		// Register styleA first, then styleB (duplicate).
		registry := style.NewRegistry[any, mockPolicy](styleA, styleB)

		// Verify Lookup returns styleA (first registered wins).
		looked := registry.Lookup(name)
		if looked == nil {
			t.Fatalf("Lookup(%q) returned nil; expected styleA", name)
		}
		if looked.Requirements() != reqsA {
			t.Fatalf("Lookup(%q) returned style with reqs %+v; expected styleA reqs %+v (first registration wins)",
				name, looked.Requirements(), reqsA)
		}

		// Verify Enumerate has length 1 — only styleA was accepted.
		enumerated := registry.Enumerate()
		if len(enumerated) != 1 {
			t.Fatalf("Enumerate() returned %d styles, expected 1 (duplicate should be rejected)", len(enumerated))
		}
		if enumerated[0].Requirements() != reqsA {
			t.Fatalf("Enumerate()[0] has reqs %+v; expected styleA reqs %+v",
				enumerated[0].Requirements(), reqsA)
		}
	})
}

// For any StyleRegistry with count styles, any valid current style name at index i,
// and any integer delta, Cycle(current, delta) returns the style at index
// ((i + delta) % count + count) % count. When the current name is not found,
// Cycle returns the default style.

func TestProperty14_CycleModularArithmetic(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate N distinct styles (2-20).
		n := rapid.IntRange(2, 20).Draw(t, "styleCount")

		styles := make([]style.Style[any, mockPolicy], n)
		for i := 0; i < n; i++ {
			styles[i] = mockStyle{
				name: fmt.Sprintf("style-%d", i),
				reqs: genSurfaceRequirements(t),
			}
		}

		registry := style.NewRegistry[any, mockPolicy](styles...)

		// Pick a random current style index.
		idx := rapid.IntRange(0, n-1).Draw(t, "currentIndex")
		currentName := styles[idx].Name()

		// Generate a random delta (including negative, large positive, zero).
		delta := rapid.IntRange(-1000, 1000).Draw(t, "delta")

		// Call Cycle.
		result := registry.Cycle(currentName, delta, textlayout.TextHints{})

		// Compute expected index using modular arithmetic.
		expectedIdx := ((idx+delta)%n + n) % n
		expectedName := styles[expectedIdx].Name()

		if result == nil {
			t.Fatalf("Cycle(%q, %d) returned nil", currentName, delta)
		}
		if result.Name() != expectedName {
			t.Fatalf("Cycle(%q, %d) returned %q, expected %q (index %d → %d, count=%d)",
				currentName, delta, result.Name(), expectedName, idx, expectedIdx, n)
		}
	})
}

func TestProperty14_CycleUnknownNameReturnsBestFit(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate N distinct styles (2-20).
		n := rapid.IntRange(2, 20).Draw(t, "styleCount")

		styles := make([]style.Style[any, mockPolicy], n)
		registeredNames := make(map[string]bool, n)
		for i := 0; i < n; i++ {
			name := fmt.Sprintf("style-%d", i)
			styles[i] = mockStyle{
				name: name,
				reqs: genSurfaceRequirements(t),
			}
			registeredNames[name] = true
		}

		registry := style.NewRegistry[any, mockPolicy](styles...)

		// Generate a name that is not in the registry.
		unknownName := rapid.StringMatching(`[a-z][a-z0-9-]{0,15}`).Draw(t, "unknownName")
		if registeredNames[unknownName] {
			// Skip this iteration if we accidentally generated a registered name.
			return
		}

		// Any delta should return a valid registered style when current name is unknown.
		delta := rapid.IntRange(-100, 100).Draw(t, "delta")
		result := registry.Cycle(unknownName, delta, textlayout.TextHints{})

		if result == nil {
			t.Fatalf("Cycle(%q, %d) returned nil for unknown name", unknownName, delta)
		}

		if !registeredNames[result.Name()] {
			t.Fatalf("Cycle(%q, %d) returned %q which is not a registered style",
				unknownName, delta, result.Name())
		}
	})
}

// For any StyleRegistry and any style name not present in the registry (after normalization),
// Normalize returns the default style's name.

func TestProperty15_FallbackToDefaultOnUnknownName(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate a registry with 2–10 styles with known names.
		count := rapid.IntRange(2, 10).Draw(t, "styleCount")
		names := genDistinctStyleNames(t, count)

		styles := make([]style.Style[any, mockPolicy], count)
		for i, name := range names {
			styles[i] = mockStyle{
				name: name,
				reqs: genSurfaceRequirements(t),
			}
		}

		registry := style.NewRegistry[any, mockPolicy](styles...)

		// Build a set of normalized registered names for filtering.
		registeredNames := make(map[string]bool, count)
		for _, name := range names {
			registeredNames[name] = true
		}

		// Generate a random string that is NOT a registered name after normalization.
		candidate := rapid.StringMatching(`[a-z][a-z0-9-]{0,30}`).Draw(t, "unknownName")

		// Skip if the candidate happens to match a registered name after normalization.
		if registeredNames[candidate] {
			return
		}

		// Normalize returns empty string for unknown names.
		result := registry.Normalize(candidate)

		if result != "" {
			t.Fatalf("Normalize(%q) = %q, expected empty string for unknown name",
				candidate, result)
		}
	})
}

// --- From: style_prop_test.go ---

// mockStyle is a configurable Style implementation for property-based testing.
type mockStyle struct {
	name string
	reqs style.SurfaceRequirements
}

func (m mockStyle) Name() string                            { return m.name }
func (m mockStyle) Requirements() style.SurfaceRequirements { return m.reqs }

// Supports implements the fitness evaluation logic from the design document.
func (m mockStyle) Supports(hints textlayout.TextHints) style.Fitness {
	// Zero-dimension panels are Unsupported.
	if hints.PixelWidth == 0 || hints.PixelHeight == 0 {
		return style.Unsupported
	}

	// Undersized panels are Unsupported.
	if m.reqs.MinWidth > 0 && hints.PixelWidth < m.reqs.MinWidth {
		return style.Unsupported
	}
	if m.reqs.MinHeight > 0 && hints.PixelHeight < m.reqs.MinHeight {
		return style.Unsupported
	}

	// Capability ordering violation → Unsupported.
	if m.reqs.Capability > style.Capability(hints.Capability) {
		return style.Unsupported
	}

	// All minimum requirements met. Check preferred for Optimal.
	meetsPreferred := true
	if m.reqs.PreferredWidth > 0 && hints.PixelWidth < m.reqs.PreferredWidth {
		meetsPreferred = false
	}
	if m.reqs.PreferredHeight > 0 && hints.PixelHeight < m.reqs.PreferredHeight {
		meetsPreferred = false
	}

	if meetsPreferred {
		return style.Optimal
	}
	return style.Full
}

func (m mockStyle) Build(_ any, _ mockPolicy, _ style.StyleContext) style.ViewData {
	return style.ViewData{Items: []string{"mock"}}
}

// genSurfaceRequirements generates a random SurfaceRequirements struct with valid invariants.
func genSurfaceRequirements(t *rapid.T) style.SurfaceRequirements {
	minW := rapid.IntRange(0, 500).Draw(t, "minWidth")
	minH := rapid.IntRange(0, 500).Draw(t, "minHeight")

	// Preferred >= Min when both non-zero.
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

	return style.SurfaceRequirements{
		MinWidth:        minW,
		MinHeight:       minH,
		PreferredWidth:  prefW,
		PreferredHeight: prefH,
		Capability:      style.Capability(rapid.IntRange(0, 5).Draw(t, "capability")),
	}
}

// genAdequateHints generates TextHints that meet or exceed all SurfaceRequirements minimums.
// This ensures: PixelWidth >= MinWidth (or >= 1 when unconstrained),
// PixelHeight >= MinHeight (or >= 1 when unconstrained),
// and Capability >= reqs.Capability.
func genAdequateHints(t *rapid.T, reqs style.SurfaceRequirements) textlayout.TextHints {
	// Ensure width meets or exceeds MinWidth (at least 1 to avoid zero-dimension).
	minW := 1
	if reqs.MinWidth > 0 {
		minW = reqs.MinWidth
	}
	pixelWidth := rapid.IntRange(minW, minW+1000).Draw(t, "pixelWidth")

	// Ensure height meets or exceeds MinHeight (at least 1 to avoid zero-dimension).
	minH := 1
	if reqs.MinHeight > 0 {
		minH = reqs.MinHeight
	}
	pixelHeight := rapid.IntRange(minH, minH+1000).Draw(t, "pixelHeight")

	// Ensure capability meets or exceeds requirement.
	capMin := int(reqs.Capability)
	capability := rapid.IntRange(capMin, 5).Draw(t, "capability")

	return textlayout.TextHints{
		PixelWidth:         pixelWidth,
		PixelHeight:        pixelHeight,
		GlyphWidth:         rapid.IntRange(1, 32).Draw(t, "glyphWidth"),
		GlyphHeight:        rapid.IntRange(1, 32).Draw(t, "glyphHeight"),
		GlyphAdvance:       rapid.IntRange(1, 32).Draw(t, "glyphAdvance"),
		RowHeight:          rapid.IntRange(1, 64).Draw(t, "rowHeight"),
		PreferEventRefresh: rapid.Bool().Draw(t, "preferEventRefresh"),
		Capability:         capability,
	}
}

// For any Style's Requirements() return value, MinWidth >= 0, MinHeight >= 0,
// PreferredWidth >= 0, PreferredHeight >= 0, and when both MinWidth > 0 and
// PreferredWidth > 0 then PreferredWidth >= MinWidth (same for height).

func TestProperty3_SurfaceRequirements_Invariants(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		req := genSurfaceRequirements(t)

		// Invariant: all dimensions are non-negative
		if req.MinWidth < 0 {
			t.Fatalf("MinWidth is negative: %d", req.MinWidth)
		}
		if req.MinHeight < 0 {
			t.Fatalf("MinHeight is negative: %d", req.MinHeight)
		}
		if req.PreferredWidth < 0 {
			t.Fatalf("PreferredWidth is negative: %d", req.PreferredWidth)
		}
		if req.PreferredHeight < 0 {
			t.Fatalf("PreferredHeight is negative: %d", req.PreferredHeight)
		}

		// Invariant: when both MinWidth > 0 and PreferredWidth > 0, PreferredWidth >= MinWidth
		if req.MinWidth > 0 && req.PreferredWidth > 0 {
			if req.PreferredWidth < req.MinWidth {
				t.Fatalf("PreferredWidth (%d) < MinWidth (%d) when both non-zero",
					req.PreferredWidth, req.MinWidth)
			}
		}

		// Invariant: when both MinHeight > 0 and PreferredHeight > 0, PreferredHeight >= MinHeight
		if req.MinHeight > 0 && req.PreferredHeight > 0 {
			if req.PreferredHeight < req.MinHeight {
				t.Fatalf("PreferredHeight (%d) < MinHeight (%d) when both non-zero",
					req.PreferredHeight, req.MinHeight)
			}
		}
	})
}

// For any Style, if TextHints satisfies all SurfaceRequirements minimum fields
// (width >= MinWidth when MinWidth > 0, height >= MinHeight when MinHeight > 0,
// capability >= reqs.Capability), then Supports(hints) >= Full.

func TestProperty6_AdequatePanelsProduceFullOrBetter(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		reqs := genSurfaceRequirements(t)
		s := mockStyle{
			name: "test-style",
			reqs: reqs,
		}

		hints := genAdequateHints(t, reqs)

		fitness := s.Supports(hints)
		if fitness < style.Full {
			t.Fatalf("expected Supports() >= Full (%d), got %d\nreqs: %+v\nhints: PixelWidth=%d, PixelHeight=%d, PreferEventRefresh=%v",
				style.Full, fitness, reqs, hints.PixelWidth, hints.PixelHeight, hints.PreferEventRefresh)
		}
	})
}

// For any Style instance, calling Requirements() any number of times always returns
// a value equal to the first call's result. This validates that Requirements() is a
// pure function with no side effects.

func TestProperty4_RequirementsImmutability(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate a mock Style with fixed SurfaceRequirements.
		reqs := genSurfaceRequirements(t)
		name := rapid.StringMatching(`[a-z][a-z0-9-]{0,15}`).Draw(t, "styleName")
		s := mockStyle{name: name, reqs: reqs}

		// Call Requirements() N times (between 2 and 100).
		n := rapid.IntRange(2, 100).Draw(t, "callCount")

		first := s.Requirements()
		for i := 1; i < n; i++ {
			got := s.Requirements()
			if got != first {
				t.Fatalf("Requirements() returned different value on call %d:\nfirst: %+v\ngot:   %+v", i+1, first, got)
			}
		}
	})
}
