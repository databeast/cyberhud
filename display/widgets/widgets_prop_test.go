package widgets_test

import (
	"fmt"
	"image"
	"image/color"
	"reflect"
	"testing"
	"unsafe"

	"github.com/databeast/cyberhud/display/surface/fonts"
	"github.com/databeast/cyberhud/display/widgets"
	"github.com/databeast/cyberhud/display/widgets/borderframe"
	"github.com/databeast/cyberhud/display/widgets/led"
	"github.com/databeast/cyberhud/display/widgets/progressbar"
	"github.com/databeast/cyberhud/display/widgets/scrollbar"
	"github.com/databeast/cyberhud/display/widgets/sparkline"
	"github.com/databeast/cyberhud/display/widgets/textlabel"
	"pgregory.net/rapid"
)

// ---------------------------------------------------------------------------
// --- From: bugcondition_prop_test.go ---
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// Bug Condition Exploration Tests
// These tests encode EXPECTED behavior. They MUST FAIL on unfixed code,
// confirming the bugs exist. After the fixes are applied, they will pass.
// ---------------------------------------------------------------------------

// TestBugCondition_BorderframeNilImage tests that borderframe.Render produces
// a Sprite with a non-nil Image field.
// On unfixed code, Build unconditionally appended Sprites with nil Images for
// any icon name that is not registered. The new Render function composites
// into a single image, so this test verifies that the output is valid.
//
// NOTE: The icons package init() auto-registers all border icons, so this test
// currently passes.

func TestBugCondition_BorderframeNilImage(t *testing.T) {
	bounds := image.Rect(0, 0, 64, 64)
	sprite := borderframe.Render(borderframe.Config{Bounds: bounds})

	if sprite == nil {
		t.Skip("Render returned nil (bounds too small?), cannot test nil images")
	}

	if sprite.Image == nil {
		t.Fatalf("Bug 1 confirmed: borderframe.Render produced Sprite with nil Image (label=%q)", sprite.Label)
	}
}

// TestBugCondition_PieGuardSingleDimension tests that progressbar.Render
// returns nil when the Pie style is used with a single dimension < 3.
// On unfixed code, this WILL FAIL because the guard uses && instead of ||,
// so a 2×100 bound passes the guard and returns a non-nil result.

func TestBugCondition_PieGuardSingleDimension(t *testing.T) {
	// Case 1: width < 3, height >= 3 (2×100)
	cfg := progressbar.Config{
		Style:      progressbar.Pie,
		Value:      0.5,
		Bounds:     image.Rect(0, 0, 2, 100),
		Foreground: color.RGBA{R: 255, G: 255, B: 255, A: 255},
		Background: color.RGBA{R: 0, G: 0, B: 0, A: 255},
	}
	result := progressbar.Render(cfg)
	if result != nil {
		t.Fatalf("Bug 5 confirmed: progressbar.Render(Pie, 2×100) returned non-nil; expected nil when width < 3")
	}

	// Case 2: height < 3, width >= 3 (100×2)
	cfg2 := progressbar.Config{
		Style:      progressbar.Pie,
		Value:      0.5,
		Bounds:     image.Rect(0, 0, 100, 2),
		Foreground: color.RGBA{R: 255, G: 255, B: 255, A: 255},
		Background: color.RGBA{R: 0, G: 0, B: 0, A: 255},
	}
	result2 := progressbar.Render(cfg2)
	if result2 != nil {
		t.Fatalf("Bug 5 confirmed: progressbar.Render(Pie, 100×2) returned non-nil; expected nil when height < 3")
	}
}

// TestBugCondition_CachedRendererThrashing tests that CachedRenderer does not
// re-invoke the render function on every call when alternating between two
// distinct configs. On unfixed code, this WILL FAIL because the single-entry
// cache evicts on every alternating call, yielding 100 render invocations
// instead of the expected ≤ 2.

func TestBugCondition_CachedRendererThrashing(t *testing.T) {
	type TestConfig struct {
		ID int
	}
	type TestResult struct {
		Value string
	}

	renderCount := 0
	render := func(cfg TestConfig) *TestResult {
		renderCount++
		return &TestResult{Value: fmt.Sprintf("result-%d", cfg.ID)}
	}
	sign := func(cfg TestConfig) uint64 {
		return uint64(cfg.ID)
	}

	cache := widgets.NewCachedRendererForTest(render, sign)

	cfgA := TestConfig{ID: 1}
	cfgB := TestConfig{ID: 2}

	// Alternate between two configs for 100 iterations
	for i := 0; i < 100; i++ {
		if i%2 == 0 {
			cache.Render(cfgA)
		} else {
			cache.Render(cfgB)
		}
	}

	// With a 2-entry cache, render should be called at most 2 times (one cold miss per config).
	// With the unfixed single-entry cache, render is called on every alternating call (100 times).
	if renderCount > 2 {
		t.Fatalf("Bug 9 confirmed: CachedRenderer invoked render %d times for 100 alternating calls between 2 configs; expected ≤ 2", renderCount)
	}
}

// ---------------------------------------------------------------------------
// Property-Based Tests (rapid)
// ---------------------------------------------------------------------------

// TestProperty_BugCondition_PieGuardCompleteness verifies that for any
// progressbar.Config with Style == Pie where EITHER Bounds.Dx() < 3 OR
// Bounds.Dy() < 3, Render returns nil.
//
// On unfixed code, this WILL FAIL with a counterexample where one dimension
// is < 3 but the other is >= 3.

func TestProperty_BugCondition_PieGuardCompleteness(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate bounds where at least one dimension is < 3
		choice := rapid.IntRange(0, 2).Draw(t, "dimChoice")

		var w, h int
		switch choice {
		case 0: // width < 3 only
			w = rapid.IntRange(1, 2).Draw(t, "width")
			h = rapid.IntRange(3, 200).Draw(t, "height")
		case 1: // height < 3 only
			w = rapid.IntRange(3, 200).Draw(t, "width")
			h = rapid.IntRange(1, 2).Draw(t, "height")
		case 2: // both < 3
			w = rapid.IntRange(1, 2).Draw(t, "width")
			h = rapid.IntRange(1, 2).Draw(t, "height")
		}

		value := rapid.Float64Range(0.0, 1.0).Draw(t, "value")

		cfg := progressbar.Config{
			Style:      progressbar.Pie,
			Value:      value,
			Bounds:     image.Rect(0, 0, w, h),
			Foreground: color.RGBA{R: 255, G: 255, B: 255, A: 255},
			Background: color.RGBA{R: 0, G: 0, B: 0, A: 255},
		}

		result := progressbar.Render(cfg)
		if result != nil {
			t.Fatalf("Pie guard incomplete: Render(Pie, %d×%d) returned non-nil; "+
				"expected nil when either dimension < 3", w, h)
		}
	})
}

// TestProperty_BugCondition_CachedRendererAlternating verifies that for any
// sequence of CachedRenderer calls alternating between 2 distinct signatures,
// the render function is invoked at most 2 times (one per unique config).
//
// On unfixed code, this WILL FAIL because the single-entry cache yields
// a render call on every alternating miss.

func TestProperty_BugCondition_CachedRendererAlternating(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		type Cfg struct{ ID int }
		type Res struct{ V string }

		renderCount := 0
		render := func(cfg Cfg) *Res {
			renderCount++
			return &Res{V: fmt.Sprintf("r-%d", cfg.ID)}
		}
		sign := func(cfg Cfg) uint64 { return uint64(cfg.ID) }

		cache := widgets.NewCachedRendererForTest(render, sign)

		// Two distinct IDs
		idA := rapid.IntRange(0, 1000).Draw(t, "idA")
		idB := rapid.IntRange(1001, 2000).Draw(t, "idB")

		cfgA := Cfg{ID: idA}
		cfgB := Cfg{ID: idB}

		// Alternate for a random number of iterations (at least 4 to exercise the cache)
		iters := rapid.IntRange(4, 50).Draw(t, "iterations")
		for i := 0; i < iters; i++ {
			if i%2 == 0 {
				cache.Render(cfgA)
			} else {
				cache.Render(cfgB)
			}
		}

		// Property: render should be invoked at most 2 times
		if renderCount > 2 {
			t.Fatalf("CachedRenderer thrashing: render invoked %d times for %d alternating calls between sig %d and %d; expected ≤ 2",
				renderCount, iters, idA, idB)
		}
	})
}

// ---------------------------------------------------------------------------
// --- From: compositor_prop_test.go ---
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// Test helpers: mock Renderables for compositor property tests
// ---------------------------------------------------------------------------

// stubRenderable always returns a non-nil Sprite with a unique label.
type stubRenderable struct {
	label    string
	rendered bool
}

func (s *stubRenderable) RenderFrame() *widgets.Sprite {
	s.rendered = true
	return &widgets.Sprite{
		Image:    image.NewRGBA(image.Rect(0, 0, 1, 1)),
		Position: image.Point{},
		Label:    s.label,
	}
}

// nilRenderable always returns nil from RenderFrame.
type nilRenderable struct {
	rendered bool
}

func (n *nilRenderable) RenderFrame() *widgets.Sprite {
	n.rendered = true
	return nil
}

// describedRenderable implements both Renderable and Described.
type describedRenderable struct {
	label      string
	descriptor widgets.Descriptor
	rendered   bool
}

func (d *describedRenderable) RenderFrame() *widgets.Sprite {
	d.rendered = true
	return &widgets.Sprite{
		Image:    image.NewRGBA(image.Rect(0, 0, 1, 1)),
		Position: image.Point{},
		Label:    d.label,
	}
}

func (d *describedRenderable) Describe() widgets.Descriptor {
	return d.descriptor
}

// ---------------------------------------------------------------------------
// Property 5: Compositor Insertion Order
// For any sequence of N Add calls where all Renderables return non-nil,
// Sprites() has length N in insertion order.

// ---------------------------------------------------------------------------

func TestProperty5_CompositorInsertionOrder(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		n := rapid.IntRange(0, 50).Draw(t, "n")

		comp := widgets.NewCompositor(widgets.SuppressionContext{})

		renderables := make([]*stubRenderable, n)
		for i := 0; i < n; i++ {
			renderables[i] = &stubRenderable{label: rapid.StringMatching(`[a-z]{1,10}`).Draw(t, "label")}
			comp.Add(renderables[i])
		}

		sprites := comp.Sprites()

		// Property: length equals N
		if len(sprites) != n {
			t.Fatalf("expected %d sprites, got %d", n, len(sprites))
		}

		// Property: insertion order preserved (labels match)
		for i, s := range sprites {
			if s.Label != renderables[i].label {
				t.Fatalf("sprite[%d].Label = %q, want %q", i, s.Label, renderables[i].label)
			}
		}
	})
}

// ---------------------------------------------------------------------------
// Property 6: Compositor Nil Discard
// For any mix where K return non-nil and (N-K) return nil,
// Sprites() has length K in relative order.

// ---------------------------------------------------------------------------

func TestProperty6_CompositorNilDiscard(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		n := rapid.IntRange(1, 50).Draw(t, "n")

		comp := widgets.NewCompositor(widgets.SuppressionContext{})

		// Track which ones are non-nil, in order
		var expectedLabels []string

		for i := 0; i < n; i++ {
			returnsNil := rapid.Bool().Draw(t, "returnsNil")
			if returnsNil {
				comp.Add(&nilRenderable{})
			} else {
				label := rapid.StringMatching(`[a-z]{1,10}`).Draw(t, "label")
				expectedLabels = append(expectedLabels, label)
				comp.Add(&stubRenderable{label: label})
			}
		}

		sprites := comp.Sprites()

		// Property: length equals count of non-nil renderables
		if len(sprites) != len(expectedLabels) {
			t.Fatalf("expected %d sprites, got %d", len(expectedLabels), len(sprites))
		}

		// Property: relative order preserved
		for i, s := range sprites {
			if s.Label != expectedLabels[i] {
				t.Fatalf("sprite[%d].Label = %q, want %q", i, s.Label, expectedLabels[i])
			}
		}
	})
}

// ---------------------------------------------------------------------------
// Property 7: AddIf False Skips Invocation
// AddIf(false, r) does not invoke RenderFrame or increase Sprites count.

// ---------------------------------------------------------------------------

func TestProperty7_AddIfFalseSkipsInvocation(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		n := rapid.IntRange(1, 30).Draw(t, "n")

		comp := widgets.NewCompositor(widgets.SuppressionContext{})

		var skipped []*stubRenderable

		for i := 0; i < n; i++ {
			r := &stubRenderable{label: rapid.StringMatching(`[a-z]{1,10}`).Draw(t, "label")}
			skipped = append(skipped, r)
			comp.AddIf(false, r)
		}

		sprites := comp.Sprites()

		// Property: no sprites accumulated
		if len(sprites) != 0 {
			t.Fatalf("expected 0 sprites after AddIf(false, ...), got %d", len(sprites))
		}

		// Property: RenderFrame was never called
		for i, r := range skipped {
			if r.rendered {
				t.Fatalf("renderable[%d] was invoked despite AddIf(false, ...)", i)
			}
		}
	})
}

// ---------------------------------------------------------------------------
// Property 8: Suppression Rule Prevents Render
// When a rule matches, Add does not invoke RenderFrame.

// ---------------------------------------------------------------------------

func TestProperty8_SuppressionRulePreventsRender(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Create a rule that always suppresses
		alwaysSuppress := func(_ widgets.Descriptor, _ widgets.SuppressionContext) bool {
			return true
		}

		comp := widgets.NewCompositor(widgets.SuppressionContext{}, alwaysSuppress)

		n := rapid.IntRange(1, 20).Draw(t, "n")

		var renderables []*describedRenderable
		for i := 0; i < n; i++ {
			r := &describedRenderable{
				label: rapid.StringMatching(`[a-z]{1,10}`).Draw(t, "label"),
				descriptor: widgets.Descriptor{
					Name:      rapid.StringMatching(`[a-z]{3,8}`).Draw(t, "name"),
					MinWidth:  rapid.IntRange(1, 100).Draw(t, "minW"),
					MinHeight: rapid.IntRange(1, 100).Draw(t, "minH"),
				},
			}
			renderables = append(renderables, r)
			comp.Add(r)
		}

		sprites := comp.Sprites()

		// Property: no sprites accumulated when all are suppressed
		if len(sprites) != 0 {
			t.Fatalf("expected 0 sprites with always-suppress rule, got %d", len(sprites))
		}

		// Property: RenderFrame was never called
		for i, r := range renderables {
			if r.rendered {
				t.Fatalf("renderable[%d] was invoked despite suppression rule", i)
			}
		}
	})
}

// ---------------------------------------------------------------------------
// Property 11: Suppression Rule OR Composition
// Widget suppressed if ANY rule returns true (short-circuit).

// ---------------------------------------------------------------------------

func TestProperty11_SuppressionRuleORComposition(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate a random set of rule results (at least 2 rules)
		numRules := rapid.IntRange(2, 8).Draw(t, "numRules")

		ruleResults := make([]bool, numRules)
		for i := range ruleResults {
			ruleResults[i] = rapid.Bool().Draw(t, "ruleResult")
		}

		// Determine expected suppression: OR of all rule results
		expectedSuppressed := false
		for _, result := range ruleResults {
			if result {
				expectedSuppressed = true
				break
			}
		}

		// Track which rules were actually evaluated
		evaluated := make([]bool, numRules)

		// Build rules from results
		rules := make([]widgets.SuppressionRule, numRules)
		for i := range rules {
			idx := i // capture loop variable
			rules[i] = func(_ widgets.Descriptor, _ widgets.SuppressionContext) bool {
				evaluated[idx] = true
				return ruleResults[idx]
			}
		}

		comp := widgets.NewCompositor(widgets.SuppressionContext{}, rules...)

		r := &describedRenderable{
			label: "test-widget",
			descriptor: widgets.Descriptor{
				Name:      "test",
				MinWidth:  10,
				MinHeight: 10,
			},
		}

		comp.Add(r)

		sprites := comp.Sprites()

		// Property: widget is suppressed iff any rule returns true
		if expectedSuppressed {
			if len(sprites) != 0 {
				t.Fatalf("expected suppression (rules=%v), but got %d sprites", ruleResults, len(sprites))
			}
			if r.rendered {
				t.Fatal("RenderFrame was called despite suppression")
			}
		} else {
			if len(sprites) != 1 {
				t.Fatalf("expected 1 sprite (no suppression, rules=%v), got %d", ruleResults, len(sprites))
			}
			if !r.rendered {
				t.Fatal("RenderFrame was not called despite no suppression")
			}
		}

		// Property: short-circuit - if rule K returns true, rules K+1..N are not evaluated
		if expectedSuppressed {
			firstTrue := -1
			for i, result := range ruleResults {
				if result {
					firstTrue = i
					break
				}
			}
			// Rules after firstTrue should NOT have been evaluated
			for i := firstTrue + 1; i < numRules; i++ {
				if evaluated[i] {
					t.Fatalf("rule[%d] was evaluated after rule[%d] returned true (short-circuit violation)", i, firstTrue)
				}
			}
		}
	})
}

// ---------------------------------------------------------------------------
// --- From: options_prop_test.go ---
// ---------------------------------------------------------------------------

// For any widget constructed with WithCaching(), calling RenderFrame() twice with
// an unchanged configuration SHALL return the same *Sprite pointer on the second
// call (cache hit), without re-executing the underlying render logic.

func TestProperty13_WithCachingMemoization(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		// Generate a valid LED configuration (diameter >= 3 ensures non-nil render)
		diameter := rapid.IntRange(3, 64).Draw(rt, "diameter")
		state := led.State(rapid.IntRange(0, 1).Draw(rt, "state"))
		minX := rapid.IntRange(0, 200).Draw(rt, "minX")
		minY := rapid.IntRange(0, 200).Draw(rt, "minY")
		bounds := image.Rect(minX, minY, minX+diameter, minY+diameter)

		fg := color.RGBA{
			R: uint8(rapid.IntRange(1, 255).Draw(rt, "fgR")),
			G: uint8(rapid.IntRange(1, 255).Draw(rt, "fgG")),
			B: uint8(rapid.IntRange(1, 255).Draw(rt, "fgB")),
			A: 255,
		}
		bg := color.RGBA{
			R: uint8(rapid.IntRange(0, 255).Draw(rt, "bgR")),
			G: uint8(rapid.IntRange(0, 255).Draw(rt, "bgG")),
			B: uint8(rapid.IntRange(0, 255).Draw(rt, "bgB")),
			A: 255,
		}

		cfg := led.Config{
			State:      state,
			Brightness: -1.0,
			Diameter:   diameter,
			Bounds:     bounds,
			Foreground: fg,
			Background: bg,
		}

		// Create widget with caching enabled
		w := led.New(cfg, widgets.WithCaching())

		// First render — populates the cache
		sprite1 := w.RenderFrame()
		if sprite1 == nil {
			rt.Fatalf("first RenderFrame() returned nil for valid cfg (diameter=%d)", diameter)
		}

		// Second render — should return same pointer (cache hit)
		sprite2 := w.RenderFrame()
		if sprite2 == nil {
			rt.Fatalf("second RenderFrame() returned nil for unchanged cfg")
		}

		// Assert same pointer — cache must return the exact same *Sprite
		ptr1 := unsafe.Pointer(sprite1)
		ptr2 := unsafe.Pointer(sprite2)
		if ptr1 != ptr2 {
			rt.Fatalf("WithCaching: second call returned different pointer (%v != %v); expected cache hit",
				ptr1, ptr2)
		}
	})
}

// For any widget constructed with WithLabel(label) and for any valid configuration
// producing a non-nil Sprite, the Sprite's Label field SHALL equal the provided
// label string, overriding the widget's default label.

func TestProperty14_WithLabelOverride(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		// Generate a custom label (non-empty, printable ASCII)
		customLabel := rapid.StringMatching(`[a-zA-Z0-9_\-/]{1,64}`).Draw(rt, "customLabel")

		// Generate a valid LED configuration (diameter >= 3 ensures non-nil render)
		diameter := rapid.IntRange(3, 64).Draw(rt, "diameter")
		state := led.State(rapid.IntRange(0, 1).Draw(rt, "state"))
		bounds := image.Rect(0, 0, diameter, diameter)

		fg := color.RGBA{
			R: uint8(rapid.IntRange(1, 255).Draw(rt, "fgR")),
			G: uint8(rapid.IntRange(1, 255).Draw(rt, "fgG")),
			B: uint8(rapid.IntRange(1, 255).Draw(rt, "fgB")),
			A: 255,
		}
		bg := color.RGBA{R: 0, G: 0, B: 0, A: 255}

		cfg := led.Config{
			State:      state,
			Brightness: -1.0,
			Diameter:   diameter,
			Bounds:     bounds,
			Foreground: fg,
			Background: bg,
		}

		// Create widget with label override
		w := led.New(cfg, widgets.WithLabel(customLabel))

		sprite := w.RenderFrame()
		if sprite == nil {
			rt.Fatalf("RenderFrame() returned nil for valid cfg (diameter=%d)", diameter)
		}

		// Property: Sprite.Label must equal the provided custom label
		if sprite.Label != customLabel {
			rt.Fatalf("WithLabel: Sprite.Label = %q, want %q", sprite.Label, customLabel)
		}
	})
}

// For any widget implementing Configurable and for any two distinct valid configurations
// A and B producing visually different outputs, calling Configure(B) after initial
// construction with A SHALL cause the next RenderFrame() to produce output consistent
// with configuration B.

func TestProperty16_ConfigureUpdatesNextRender(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		// Generate two distinct progressbar configs that will produce different visuals.
		// Use Horizontal style for predictable pixel differences.
		w := rapid.IntRange(10, 100).Draw(rt, "width")
		h := rapid.IntRange(5, 50).Draw(rt, "height")
		boundsA := image.Rect(0, 0, w, h)

		// Config A: low value → few filled pixels
		valueA := rapid.Float64Range(0.0, 0.3).Draw(rt, "valueA")
		// Config B: high value → many filled pixels (guaranteed different)
		valueB := rapid.Float64Range(0.7, 1.0).Draw(rt, "valueB")

		fg := color.RGBA{R: 255, G: 255, B: 255, A: 255}
		bg := color.RGBA{R: 0, G: 0, B: 0, A: 255}

		cfgA := progressbar.Config{
			Style:      progressbar.Linear,
			Value:      valueA,
			Bounds:     boundsA,
			Foreground: fg,
			Background: bg,
		}

		cfgB := progressbar.Config{
			Style:      progressbar.Linear,
			Value:      valueB,
			Bounds:     boundsA,
			Foreground: fg,
			Background: bg,
		}

		// Create widget with config A
		renderable := progressbar.New(cfgA)

		// Render with config A
		spriteA := renderable.RenderFrame()
		if spriteA == nil {
			rt.Fatalf("RenderFrame() with config A returned nil (bounds=%v, value=%.2f)", boundsA, valueA)
		}

		// Configure with B
		configurable, ok := renderable.(widgets.Configurable)
		if !ok {
			rt.Fatal("progressbar widget does not implement Configurable")
		}
		configurable.Configure(cfgB)

		// Render with config B
		spriteB := renderable.RenderFrame()
		if spriteB == nil {
			rt.Fatalf("RenderFrame() with config B returned nil (bounds=%v, value=%.2f)", boundsA, valueB)
		}

		// Property: output must differ between A and B.
		// Compare pixel data — at least one pixel must differ.
		imgA, okA := spriteA.Image.(*image.RGBA)
		imgB, okB := spriteB.Image.(*image.RGBA)
		if !okA || !okB {
			rt.Fatal("expected *image.RGBA from both renders")
		}

		if len(imgA.Pix) != len(imgB.Pix) {
			// Different sizes → definitely different
			return
		}

		anyDiff := false
		for i := range imgA.Pix {
			if imgA.Pix[i] != imgB.Pix[i] {
				anyDiff = true
				break
			}
		}

		if !anyDiff {
			rt.Fatalf("Configure(B) did not change render output: valueA=%.4f, valueB=%.4f, bounds=%v",
				valueA, valueB, boundsA)
		}
	})
}

// ---------------------------------------------------------------------------
// --- From: preservation_prop_test.go ---
// ---------------------------------------------------------------------------

// ============================================================================
// Non-Buggy Widget Behavior Unchanged
//
// These tests capture baseline behaviors on UNFIXED code that must be preserved
// after the widget quality fixes are applied. They verify:
// - Progressbar Pie renders correctly when both dimensions ≥ 3
// - LED renders with non-zero foreground/background colors preserved
// - Textlabel uses the provided non-zero foreground unchanged
// - CachedRenderer caches same-config repeated calls (1 render invocation)

// ============================================================================

// TestPreservation_ProgressbarPieValidBoundsNonNil verifies that for any
// progressbar.Config with Style == Pie and Bounds.Dx() >= 3 && Bounds.Dy() >= 3,
// Render(cfg) returns non-nil with a valid image of matching dimensions.
//
// Property: For all Pie configs with adequate bounds, the result is non-nil and
// the image has the expected dimensions.

func TestPreservation_ProgressbarPieValidBoundsNonNil(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate bounds where both dimensions are >= 3
		w := rapid.IntRange(3, 100).Draw(t, "width")
		h := rapid.IntRange(3, 100).Draw(t, "height")
		minX := rapid.IntRange(0, 200).Draw(t, "minX")
		minY := rapid.IntRange(0, 200).Draw(t, "minY")
		bounds := image.Rect(minX, minY, minX+w, minY+h)

		value := rapid.Float64Range(0.0, 1.0).Draw(t, "value")

		fg := color.RGBA{
			R: uint8(rapid.IntRange(1, 255).Draw(t, "fgR")),
			G: uint8(rapid.IntRange(0, 255).Draw(t, "fgG")),
			B: uint8(rapid.IntRange(0, 255).Draw(t, "fgB")),
			A: uint8(rapid.IntRange(1, 255).Draw(t, "fgA")),
		}
		bg := color.RGBA{
			R: uint8(rapid.IntRange(0, 255).Draw(t, "bgR")),
			G: uint8(rapid.IntRange(0, 255).Draw(t, "bgG")),
			B: uint8(rapid.IntRange(0, 255).Draw(t, "bgB")),
			A: uint8(rapid.IntRange(1, 255).Draw(t, "bgA")),
		}

		cfg := progressbar.Config{
			Style:      progressbar.Pie,
			Value:      value,
			Bounds:     bounds,
			Foreground: fg,
			Background: bg,
		}

		result := progressbar.Render(cfg)

		if result == nil {
			t.Fatalf("Render returned nil for Pie with valid bounds %dx%d", w, h)
		}

		if result.Image == nil {
			t.Fatalf("Render returned non-nil result but nil Image for Pie with valid bounds %dx%d", w, h)
		}

		// Verify image dimensions match the bounds
		imgBounds := result.Image.Bounds()
		if imgBounds.Dx() != w || imgBounds.Dy() != h {
			t.Fatalf("Image dimensions %dx%d don't match bounds %dx%d",
				imgBounds.Dx(), imgBounds.Dy(), w, h)
		}
	})
}

// TestPreservation_LEDNonZeroColorsPreserved verifies that for any led.Config
// with non-zero Foreground and Background and valid diameter (>= 3),
// Render(cfg) output uses the provided foreground and background colors.
//
// For On state: pixels inside the circle should be foreground-colored.
// For Off state: pixels on the outline should be dimmed foreground.
// Background pixels outside the circle should match the provided background.
//
// Property: For all LED configs with non-zero fg/bg, the rendered colors
// match the expected foreground/background (no color substitution occurs).

func TestPreservation_LEDNonZeroColorsPreserved(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		diameter := rapid.IntRange(5, 32).Draw(t, "diameter")

		// Non-zero foreground - ensure at least one channel is non-zero and alpha is non-zero
		fg := color.RGBA{
			R: uint8(rapid.IntRange(1, 255).Draw(t, "fgR")),
			G: uint8(rapid.IntRange(1, 255).Draw(t, "fgG")),
			B: uint8(rapid.IntRange(1, 255).Draw(t, "fgB")),
			A: uint8(rapid.IntRange(1, 255).Draw(t, "fgA")),
		}
		// Non-zero background
		bg := color.RGBA{
			R: uint8(rapid.IntRange(0, 255).Draw(t, "bgR")),
			G: uint8(rapid.IntRange(0, 255).Draw(t, "bgG")),
			B: uint8(rapid.IntRange(0, 255).Draw(t, "bgB")),
			A: uint8(rapid.IntRange(1, 255).Draw(t, "bgA")),
		}
		// Ensure bg is non-zero: at least one of R/G/B must be non-zero or A is sufficient
		// to distinguish from zero value (color.RGBA{})
		if bg == (color.RGBA{}) {
			bg.A = 255
		}

		cfg := led.Config{
			State:      led.On,
			Brightness: -1.0,
			Diameter:   diameter,
			Bounds:     image.Rect(0, 0, diameter, diameter),
			Foreground: fg,
			Background: bg,
		}

		result := led.Render(cfg)
		if result == nil {
			t.Fatalf("Render returned nil for valid diameter %d", diameter)
		}

		img := result.Image
		if img == nil {
			t.Fatalf("Render returned nil image for valid diameter %d", diameter)
		}

		// For On state, the center pixel should be the foreground color.
		// Center of the LED circle is at (diameter/2, diameter/2).
		cx := diameter / 2
		cy := diameter / 2
		centerColor := img.At(cx, cy)
		r, g, b, a := centerColor.RGBA()
		// Convert to 8-bit
		centerRGBA := color.RGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: uint8(a >> 8)}

		// The center pixel should be the foreground (for On state, no shine at center for d>=5)
		// The shine pixel is at d/4, d/4, so center should be fg
		if centerRGBA != fg {
			t.Fatalf("center pixel (%d,%d) = %v, want foreground %v (diameter=%d)",
				cx, cy, centerRGBA, fg, diameter)
		}

		// A corner pixel (0,0) should be either the background color or transparent
		// (the LED sprite may not fill outside the circle if no background painting
		// occurs at that distance from the center).
		cornerColor := img.At(0, 0)
		r, g, b, a = cornerColor.RGBA()
		cornerRGBA := color.RGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: uint8(a >> 8)}

		transparent := color.RGBA{}
		if cornerRGBA != bg && cornerRGBA != transparent {
			t.Fatalf("corner pixel (0,0) = %v, want background %v or transparent (diameter=%d)",
				cornerRGBA, bg, diameter)
		}
	})
}

// TestPreservation_TextlabelNonZeroForegroundPreserved verifies that for any
// textlabel.Config with a non-zero Foreground and valid bounds, the rendered
// output uses the provided foreground color unchanged for text pixels.
//
// Property: For all textlabel configs with non-zero foreground and non-empty text,
// any non-transparent pixel in the output equals the provided Foreground.

func TestPreservation_TextlabelNonZeroForegroundPreserved(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		w := rapid.IntRange(20, 100).Draw(t, "width")
		h := rapid.IntRange(8, 50).Draw(t, "height")
		bounds := image.Rect(0, 0, w, h)

		// Non-zero foreground with all channels explicitly non-zero
		fg := color.RGBA{
			R: uint8(rapid.IntRange(1, 255).Draw(t, "fgR")),
			G: uint8(rapid.IntRange(1, 255).Draw(t, "fgG")),
			B: uint8(rapid.IntRange(1, 255).Draw(t, "fgB")),
			A: 255,
		}

		// Use a character that we know renders pixels with the default font
		text := "A"

		cfg := textlabel.Config{
			Text:       text,
			Bounds:     bounds,
			Font:       font.Default(),
			Alignment:  textlabel.Left,
			Foreground: fg,
		}

		result := textlabel.Render(cfg)
		if result == nil {
			t.Fatalf("Render returned nil for valid bounds %dx%d", w, h)
		}

		img := result.Image
		if img == nil {
			t.Fatalf("Render returned nil image for valid bounds %dx%d", w, h)
		}

		// Scan for any non-transparent pixel and verify it matches the foreground
		transparent := color.RGBA{}
		foundFgPixel := false

		imgBounds := img.Bounds()
		for y := imgBounds.Min.Y; y < imgBounds.Max.Y; y++ {
			for x := imgBounds.Min.X; x < imgBounds.Max.X; x++ {
				r, g, b, a := img.At(x, y).RGBA()
				px := color.RGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: uint8(a >> 8)}
				if px != transparent {
					foundFgPixel = true
					if px != fg {
						t.Fatalf("pixel (%d,%d) = %v, want foreground %v (text=%q)",
							x, y, px, fg, text)
					}
				}
			}
		}

		if !foundFgPixel {
			t.Fatalf("no foreground pixels found for text=%q with bounds %dx%d", text, w, h)
		}
	})
}

// TestPreservation_CachedRendererSameConfigSingleRender verifies that for any
// sequence of CachedRenderer.Render calls with the same config signature,
// the underlying render function is invoked exactly once.
//
// Property: For all same-signature call sequences of length N >= 1, the render
// function is called exactly 1 time (all subsequent calls are cache hits).

func TestPreservation_CachedRendererSameConfigSingleRender(t *testing.T) {
	type TestConfig struct {
		ID int
	}

	type TestResult struct {
		Value string
	}

	rapid.Check(t, func(t *rapid.T) {
		callCount := 0
		render := func(cfg TestConfig) *TestResult {
			callCount++
			return &TestResult{Value: fmt.Sprintf("rendered-%d", cfg.ID)}
		}
		sign := func(cfg TestConfig) uint64 {
			return uint64(cfg.ID)
		}

		cache := widgets.NewCachedRendererForTest(render, sign)

		// Generate a single config ID to repeat
		cfgID := rapid.IntRange(0, 1000).Draw(t, "cfgID")
		cfg := TestConfig{ID: cfgID}

		// Call Render multiple times with the same config
		numCalls := rapid.IntRange(2, 20).Draw(t, "numCalls")
		var firstResult *TestResult

		for i := 0; i < numCalls; i++ {
			result := cache.Render(cfg)
			if result == nil {
				t.Fatalf("call %d: Render returned nil for cfg %+v", i, cfg)
			}
			if i == 0 {
				firstResult = result
			} else {
				// All subsequent calls should return the same pointer (cached)
				if result != firstResult {
					t.Fatalf("call %d: expected cached result (same pointer), got different pointer", i)
				}
			}
		}

		// The render function should have been called exactly once
		if callCount != 1 {
			t.Fatalf("render function called %d times for %d same-config calls, want exactly 1",
				callCount, numCalls)
		}
	})
}

// ---------------------------------------------------------------------------
// --- From: widgets_prop_test.go ---
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// Source: border_prop_test.go
// ---------------------------------------------------------------------------

// TestPropertyDrawBorderEdgePixelInvariant verifies that for any image with dimensions
// width >= 1 and height >= 1, and for any foreground color, after calling DrawBorder,
// every pixel on the top row, bottom row, left column, and right column equals fg.

func TestPropertyDrawBorderEdgePixelInvariant(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		width := rapid.IntRange(1, 256).Draw(t, "width")
		height := rapid.IntRange(1, 256).Draw(t, "height")

		fg := color.RGBA{
			R: uint8(rapid.IntRange(0, 255).Draw(t, "fgR")),
			G: uint8(rapid.IntRange(0, 255).Draw(t, "fgG")),
			B: uint8(rapid.IntRange(0, 255).Draw(t, "fgB")),
			A: uint8(rapid.IntRange(0, 255).Draw(t, "fgA")),
		}

		img := image.NewRGBA(image.Rect(0, 0, width, height))
		widgets.DrawBorder(img, width, height, fg)

		// Verify top row (y=0)
		for x := 0; x < width; x++ {
			got := img.RGBAAt(x, 0)
			if got != fg {
				t.Fatalf("top row pixel (%d, 0): got %v, want %v", x, got, fg)
			}
		}

		// Verify bottom row (y=height-1)
		for x := 0; x < width; x++ {
			got := img.RGBAAt(x, height-1)
			if got != fg {
				t.Fatalf("bottom row pixel (%d, %d): got %v, want %v", x, height-1, got, fg)
			}
		}

		// Verify left column (x=0)
		for y := 0; y < height; y++ {
			got := img.RGBAAt(0, y)
			if got != fg {
				t.Fatalf("left column pixel (0, %d): got %v, want %v", y, got, fg)
			}
		}

		// Verify right column (x=width-1)
		for y := 0; y < height; y++ {
			got := img.RGBAAt(width-1, y)
			if got != fg {
				t.Fatalf("right column pixel (%d, %d): got %v, want %v", width-1, y, got, fg)
			}
		}
	})
}

// ---------------------------------------------------------------------------
// Source: cached_renderer_prop_test.go
// ---------------------------------------------------------------------------

// TestPropertyCachedRendererConsistency verifies that for any sequence of Configs
// and for any Config cfg in that sequence, CachedRenderer.Render(cfg) returns a
// result equal to directly calling the underlying render function with the same cfg.
// The cache must not serve stale results when the Config signature changes.

func TestPropertyCachedRendererConsistency(t *testing.T) {
	type TestConfig struct {
		ID int
	}

	type TestResult struct {
		Value string
	}

	render := func(cfg TestConfig) *TestResult {
		return &TestResult{Value: fmt.Sprintf("rendered-%d", cfg.ID)}
	}

	sign := func(cfg TestConfig) uint64 {
		return uint64(cfg.ID)
	}

	rapid.Check(t, func(t *rapid.T) {
		cache := widgets.NewCachedRendererForTest(render, sign)

		// Generate a random sequence of Config values (1 to 50 entries)
		seqLen := rapid.IntRange(1, 50).Draw(t, "seqLen")

		for i := 0; i < seqLen; i++ {
			cfg := TestConfig{
				ID: rapid.IntRange(-1000, 1000).Draw(t, fmt.Sprintf("cfg[%d].ID", i)),
			}

			// Get result from CachedRenderer
			cachedResult := cache.Render(cfg)

			// Get result directly from the underlying render function
			directResult := render(cfg)

			// Property: cached result must equal direct render result
			if cachedResult == nil {
				t.Fatalf("step %d: CachedRenderer.Render returned nil for cfg %+v", i, cfg)
			}
			if directResult == nil {
				t.Fatalf("step %d: direct render returned nil for cfg %+v", i, cfg)
			}
			if cachedResult.Value != directResult.Value {
				t.Fatalf("step %d: CachedRenderer.Render(cfg=%+v) = %q, want %q",
					i, cfg, cachedResult.Value, directResult.Value)
			}
		}
	})
}

// ---------------------------------------------------------------------------
// Source: colors_prop_test.go
// ---------------------------------------------------------------------------

// TestPropertyResolveColorsDefaultingAndPassthrough verifies that for any pair
// of color.RGBA values (fg, bg):
// - If fg is the zero value, the returned foreground is {255, 255, 255, 255}.
// - If bg is the zero value, the returned background is {0, 0, 0, 255}.
// - If fg is non-zero, the returned foreground equals fg.
// - If bg is non-zero, the returned background equals bg.

func TestPropertyResolveColorsDefaultingAndPassthrough(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate random RGBA values for fg and bg
		fg := color.RGBA{
			R: rapid.Byte().Draw(t, "fgR"),
			G: rapid.Byte().Draw(t, "fgG"),
			B: rapid.Byte().Draw(t, "fgB"),
			A: rapid.Byte().Draw(t, "fgA"),
		}
		bg := color.RGBA{
			R: rapid.Byte().Draw(t, "bgR"),
			G: rapid.Byte().Draw(t, "bgG"),
			B: rapid.Byte().Draw(t, "bgB"),
			A: rapid.Byte().Draw(t, "bgA"),
		}

		retFg, retBg := widgets.ResolveColors(fg, bg)

		zero := color.RGBA{}
		defaultWhite := color.RGBA{R: 255, G: 255, B: 255, A: 255}
		defaultBlack := color.RGBA{R: 0, G: 0, B: 0, A: 255}

		// Property: zero fg → default white foreground
		if fg == zero {
			if retFg != defaultWhite {
				t.Fatalf("zero fg: expected default white %v, got %v", defaultWhite, retFg)
			}
		} else {
			// Property: non-zero fg → passthrough
			if retFg != fg {
				t.Fatalf("non-zero fg: expected passthrough %v, got %v", fg, retFg)
			}
		}

		// Property: zero bg → default black background
		if bg == zero {
			if retBg != defaultBlack {
				t.Fatalf("zero bg: expected default black %v, got %v", defaultBlack, retBg)
			}
		} else {
			// Property: non-zero bg → passthrough
			if retBg != bg {
				t.Fatalf("non-zero bg: expected passthrough %v, got %v", bg, retBg)
			}
		}
	})
}

// ---------------------------------------------------------------------------
// Source: nil_guard_prop_test.go
// ---------------------------------------------------------------------------

// Property 19: Invalid Bounds/Size Produce Nil
// For any widget Config where bounds width < 1 or height < 1 (textlabel, scrollbar, sparkline),
// or where diameter < 3 (LED) or total < 1 (scrollbar), the Render function SHALL return nil.

// TestPropertyTextLabelNilForInvalidBounds verifies textlabel.Render returns nil
// when bounds have width ≤ 0 or height ≤ 0.
func TestPropertyTextLabelNilForInvalidBounds(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Decide which dimension(s) to make invalid: 0=width, 1=height, 2=both
		invalidDim := rapid.IntRange(0, 2).Draw(t, "invalidDim")

		var bounds image.Rectangle
		switch invalidDim {
		case 0: // width invalid
			w := rapid.IntRange(-10, 0).Draw(t, "width")
			h := rapid.IntRange(1, 50).Draw(t, "height")
			bounds = image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: w, Y: h}}
		case 1: // height invalid
			w := rapid.IntRange(1, 50).Draw(t, "width")
			h := rapid.IntRange(-10, 0).Draw(t, "height")
			bounds = image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: w, Y: h}}
		case 2: // both invalid
			w := rapid.IntRange(-10, 0).Draw(t, "width")
			h := rapid.IntRange(-10, 0).Draw(t, "height")
			bounds = image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: w, Y: h}}
		}

		text := rapid.StringMatching(`[a-zA-Z0-9 ]{0,20}`).Draw(t, "text")

		result := textlabel.Render(textlabel.Config{
			Text:       text,
			Bounds:     bounds,
			Font:       font.Default(),
			Alignment:  textlabel.Left,
			Foreground: color.RGBA{R: 255, G: 255, B: 255, A: 255},
		})

		if result != nil {
			t.Fatalf("expected nil for invalid bounds %v, got non-nil result", bounds)
		}
	})
}

// TestPropertyScrollbarNilForInvalidBounds verifies scrollbar.Render returns nil
// when bounds have width ≤ 0 or height ≤ 0.
func TestPropertyScrollbarNilForInvalidBounds(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		invalidDim := rapid.IntRange(0, 2).Draw(t, "invalidDim")

		var bounds image.Rectangle
		switch invalidDim {
		case 0: // width invalid
			w := rapid.IntRange(-10, 0).Draw(t, "width")
			h := rapid.IntRange(1, 50).Draw(t, "height")
			bounds = image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: w, Y: h}}
		case 1: // height invalid
			w := rapid.IntRange(1, 50).Draw(t, "width")
			h := rapid.IntRange(-10, 0).Draw(t, "height")
			bounds = image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: w, Y: h}}
		case 2: // both invalid
			w := rapid.IntRange(-10, 0).Draw(t, "width")
			h := rapid.IntRange(-10, 0).Draw(t, "height")
			bounds = image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: w, Y: h}}
		}

		// Use a valid total so that only bounds invalidity triggers nil
		total := rapid.IntRange(1, 100).Draw(t, "total")
		visible := rapid.IntRange(1, total).Draw(t, "visible")
		offset := rapid.IntRange(0, total-1).Draw(t, "offset")

		result := scrollbar.Render(scrollbar.Config{
			TotalItems:   total,
			VisibleItems: visible,
			ScrollOffset: offset,
			Bounds:       bounds,
			Foreground:   color.RGBA{R: 255, G: 255, B: 255, A: 255},
			Background:   color.RGBA{R: 0, G: 0, B: 0, A: 255},
		})

		if result != nil {
			t.Fatalf("expected nil for invalid bounds %v, got non-nil result", bounds)
		}
	})
}

// TestPropertyScrollbarNilForInvalidTotal verifies scrollbar.Render returns nil
// when TotalItems ≤ 0, even with valid bounds.
func TestPropertyScrollbarNilForInvalidTotal(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		w := rapid.IntRange(1, 50).Draw(t, "width")
		h := rapid.IntRange(1, 50).Draw(t, "height")
		bounds := image.Rectangle{
			Min: image.Point{X: 0, Y: 0},
			Max: image.Point{X: w, Y: h},
		}

		total := rapid.IntRange(-10, 0).Draw(t, "total")

		result := scrollbar.Render(scrollbar.Config{
			TotalItems:   total,
			VisibleItems: 5,
			ScrollOffset: 0,
			Bounds:       bounds,
			Foreground:   color.RGBA{R: 255, G: 255, B: 255, A: 255},
			Background:   color.RGBA{R: 0, G: 0, B: 0, A: 255},
		})

		if result != nil {
			t.Fatalf("expected nil for invalid total %d, got non-nil result", total)
		}
	})
}

// TestPropertySparklineNilForInvalidBounds verifies sparkline.Render returns nil
// when bounds have width ≤ 0 or height ≤ 0.
func TestPropertySparklineNilForInvalidBounds(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		invalidDim := rapid.IntRange(0, 2).Draw(t, "invalidDim")

		var bounds image.Rectangle
		switch invalidDim {
		case 0: // width invalid
			w := rapid.IntRange(-10, 0).Draw(t, "width")
			h := rapid.IntRange(1, 50).Draw(t, "height")
			bounds = image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: w, Y: h}}
		case 1: // height invalid
			w := rapid.IntRange(1, 50).Draw(t, "width")
			h := rapid.IntRange(-10, 0).Draw(t, "height")
			bounds = image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: w, Y: h}}
		case 2: // both invalid
			w := rapid.IntRange(-10, 0).Draw(t, "width")
			h := rapid.IntRange(-10, 0).Draw(t, "height")
			bounds = image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: w, Y: h}}
		}

		// Generate some data points
		dataLen := rapid.IntRange(0, 20).Draw(t, "dataLen")
		data := make([]float64, dataLen)
		for i := range data {
			data[i] = rapid.Float64Range(0.0, 1.0).Draw(t, "dataPoint")
		}

		style := sparkline.Style(rapid.IntRange(0, 1).Draw(t, "style"))

		result := sparkline.Render(sparkline.Config{
			Data:       data,
			Style:      style,
			Bounds:     bounds,
			Foreground: color.RGBA{R: 255, G: 255, B: 255, A: 255},
			Background: color.RGBA{R: 0, G: 0, B: 0, A: 255},
		})

		if result != nil {
			t.Fatalf("expected nil for invalid bounds %v, got non-nil result", bounds)
		}
	})
}

// TestPropertyLEDNilForInvalidDiameter verifies led.Render returns nil
// when diameter < 3.
func TestPropertyLEDNilForInvalidDiameter(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		diameter := rapid.IntRange(-10, 2).Draw(t, "diameter")
		state := led.State(rapid.IntRange(0, 1).Draw(t, "state"))

		result := led.Render(led.Config{
			State:      state,
			Brightness: -1.0,
			Diameter:   diameter,
			Bounds:     image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: 10, Y: 10}},
			Foreground: color.RGBA{R: 0, G: 200, B: 0, A: 255},
			Background: color.RGBA{R: 0, G: 0, B: 0, A: 255},
		})

		if result != nil {
			t.Fatalf("expected nil for invalid diameter %d, got non-nil result", diameter)
		}
	})
}

// ---------------------------------------------------------------------------
// Source: render_purity_prop_test.go
// ---------------------------------------------------------------------------

// Property 4: Render purity (Config immutability)
// For any widget subpackage and for any valid Config value, calling Render(cfg) SHALL NOT
// mutate the Config argument. A deep comparison of the Config before and after the call
// SHALL show equality.

// TestPropertyRenderPurity_Progressbar verifies that progressbar.Render does not mutate its Config.
func TestPropertyRenderPurity_Progressbar(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		w := rapid.IntRange(1, 100).Draw(t, "width")
		h := rapid.IntRange(1, 100).Draw(t, "height")
		bounds := image.Rectangle{
			Min: image.Point{X: rapid.IntRange(0, 200).Draw(t, "minX"), Y: rapid.IntRange(0, 200).Draw(t, "minY")},
			Max: image.Point{},
		}
		bounds.Max = image.Point{X: bounds.Min.X + w, Y: bounds.Min.Y + h}

		style := progressbar.Style(rapid.IntRange(0, 2).Draw(t, "style"))
		value := rapid.Float64Range(0.0, 1.0).Draw(t, "value")
		fg := color.RGBA{
			R: uint8(rapid.IntRange(0, 255).Draw(t, "fgR")),
			G: uint8(rapid.IntRange(0, 255).Draw(t, "fgG")),
			B: uint8(rapid.IntRange(0, 255).Draw(t, "fgB")),
			A: uint8(rapid.IntRange(0, 255).Draw(t, "fgA")),
		}
		bg := color.RGBA{
			R: uint8(rapid.IntRange(0, 255).Draw(t, "bgR")),
			G: uint8(rapid.IntRange(0, 255).Draw(t, "bgG")),
			B: uint8(rapid.IntRange(0, 255).Draw(t, "bgB")),
			A: uint8(rapid.IntRange(0, 255).Draw(t, "bgA")),
		}

		cfg := progressbar.Config{
			Style:      style,
			Value:      value,
			Bounds:     bounds,
			Foreground: fg,
			Background: bg,
		}

		before := cfg // struct copy (value type)
		progressbar.Render(cfg)

		if !reflect.DeepEqual(before, cfg) {
			t.Fatalf("progressbar.Render mutated Config:\nbefore: %+v\nafter:  %+v", before, cfg)
		}
	})
}

// TestPropertyRenderPurity_Scrollbar verifies that scrollbar.Render does not mutate its Config.
func TestPropertyRenderPurity_Scrollbar(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		w := rapid.IntRange(1, 50).Draw(t, "width")
		h := rapid.IntRange(1, 50).Draw(t, "height")
		bounds := image.Rectangle{
			Min: image.Point{X: rapid.IntRange(0, 200).Draw(t, "minX"), Y: rapid.IntRange(0, 200).Draw(t, "minY")},
			Max: image.Point{},
		}
		bounds.Max = image.Point{X: bounds.Min.X + w, Y: bounds.Min.Y + h}

		total := rapid.IntRange(1, 100).Draw(t, "total")
		visible := rapid.IntRange(1, total).Draw(t, "visible")
		offset := rapid.IntRange(0, total-1).Draw(t, "offset")
		fg := color.RGBA{
			R: uint8(rapid.IntRange(0, 255).Draw(t, "fgR")),
			G: uint8(rapid.IntRange(0, 255).Draw(t, "fgG")),
			B: uint8(rapid.IntRange(0, 255).Draw(t, "fgB")),
			A: uint8(rapid.IntRange(0, 255).Draw(t, "fgA")),
		}
		bg := color.RGBA{
			R: uint8(rapid.IntRange(0, 255).Draw(t, "bgR")),
			G: uint8(rapid.IntRange(0, 255).Draw(t, "bgG")),
			B: uint8(rapid.IntRange(0, 255).Draw(t, "bgB")),
			A: uint8(rapid.IntRange(0, 255).Draw(t, "bgA")),
		}

		cfg := scrollbar.Config{
			TotalItems:   total,
			VisibleItems: visible,
			ScrollOffset: offset,
			Bounds:       bounds,
			Foreground:   fg,
			Background:   bg,
		}

		before := cfg
		scrollbar.Render(cfg)

		if !reflect.DeepEqual(before, cfg) {
			t.Fatalf("scrollbar.Render mutated Config:\nbefore: %+v\nafter:  %+v", before, cfg)
		}
	})
}

// TestPropertyRenderPurity_Sparkline verifies that sparkline.Render does not mutate its Config.
func TestPropertyRenderPurity_Sparkline(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		w := rapid.IntRange(1, 100).Draw(t, "width")
		h := rapid.IntRange(1, 100).Draw(t, "height")
		bounds := image.Rectangle{
			Min: image.Point{X: rapid.IntRange(0, 200).Draw(t, "minX"), Y: rapid.IntRange(0, 200).Draw(t, "minY")},
			Max: image.Point{},
		}
		bounds.Max = image.Point{X: bounds.Min.X + w, Y: bounds.Min.Y + h}

		dataLen := rapid.IntRange(0, 30).Draw(t, "dataLen")
		data := make([]float64, dataLen)
		for i := range data {
			data[i] = rapid.Float64Range(0.0, 1.0).Draw(t, "dataPoint")
		}

		style := sparkline.Style(rapid.IntRange(0, 1).Draw(t, "style"))
		fg := color.RGBA{
			R: uint8(rapid.IntRange(0, 255).Draw(t, "fgR")),
			G: uint8(rapid.IntRange(0, 255).Draw(t, "fgG")),
			B: uint8(rapid.IntRange(0, 255).Draw(t, "fgB")),
			A: uint8(rapid.IntRange(0, 255).Draw(t, "fgA")),
		}
		bg := color.RGBA{
			R: uint8(rapid.IntRange(0, 255).Draw(t, "bgR")),
			G: uint8(rapid.IntRange(0, 255).Draw(t, "bgG")),
			B: uint8(rapid.IntRange(0, 255).Draw(t, "bgB")),
			A: uint8(rapid.IntRange(0, 255).Draw(t, "bgA")),
		}

		cfg := sparkline.Config{
			Data:       data,
			Style:      style,
			Bounds:     bounds,
			Foreground: fg,
			Background: bg,
		}

		// Deep-copy for sparkline since Data is a slice
		beforeData := make([]float64, len(data))
		copy(beforeData, data)
		before := sparkline.Config{
			Data:       beforeData,
			Style:      cfg.Style,
			Bounds:     cfg.Bounds,
			Foreground: cfg.Foreground,
			Background: cfg.Background,
		}

		sparkline.Render(cfg)

		if !reflect.DeepEqual(before, cfg) {
			t.Fatalf("sparkline.Render mutated Config:\nbefore: %+v\nafter:  %+v", before, cfg)
		}
	})
}

// TestPropertyRenderPurity_LED verifies that led.Render does not mutate its Config.
func TestPropertyRenderPurity_LED(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		diameter := rapid.IntRange(3, 64).Draw(t, "diameter")
		state := led.State(rapid.IntRange(0, 1).Draw(t, "state"))
		bounds := image.Rectangle{
			Min: image.Point{X: rapid.IntRange(0, 200).Draw(t, "minX"), Y: rapid.IntRange(0, 200).Draw(t, "minY")},
			Max: image.Point{X: rapid.IntRange(0, 200).Draw(t, "maxX"), Y: rapid.IntRange(0, 200).Draw(t, "maxY")},
		}
		fg := color.RGBA{
			R: uint8(rapid.IntRange(0, 255).Draw(t, "fgR")),
			G: uint8(rapid.IntRange(0, 255).Draw(t, "fgG")),
			B: uint8(rapid.IntRange(0, 255).Draw(t, "fgB")),
			A: uint8(rapid.IntRange(0, 255).Draw(t, "fgA")),
		}
		bg := color.RGBA{
			R: uint8(rapid.IntRange(0, 255).Draw(t, "bgR")),
			G: uint8(rapid.IntRange(0, 255).Draw(t, "bgG")),
			B: uint8(rapid.IntRange(0, 255).Draw(t, "bgB")),
			A: uint8(rapid.IntRange(0, 255).Draw(t, "bgA")),
		}

		cfg := led.Config{
			State:      state,
			Brightness: -1.0,
			Diameter:   diameter,
			Bounds:     bounds,
			Foreground: fg,
			Background: bg,
		}

		before := cfg
		led.Render(cfg)

		if !reflect.DeepEqual(before, cfg) {
			t.Fatalf("led.Render mutated Config:\nbefore: %+v\nafter:  %+v", before, cfg)
		}
	})
}

// TestPropertyRenderPurity_Textlabel verifies that textlabel.Render does not mutate its Config.
func TestPropertyRenderPurity_Textlabel(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		w := rapid.IntRange(1, 100).Draw(t, "width")
		h := rapid.IntRange(1, 100).Draw(t, "height")
		bounds := image.Rectangle{
			Min: image.Point{X: rapid.IntRange(0, 200).Draw(t, "minX"), Y: rapid.IntRange(0, 200).Draw(t, "minY")},
			Max: image.Point{},
		}
		bounds.Max = image.Point{X: bounds.Min.X + w, Y: bounds.Min.Y + h}

		text := rapid.StringMatching(`[a-zA-Z0-9 ]{0,20}`).Draw(t, "text")
		alignment := textlabel.Alignment(rapid.IntRange(0, 2).Draw(t, "alignment"))
		fg := color.RGBA{
			R: uint8(rapid.IntRange(0, 255).Draw(t, "fgR")),
			G: uint8(rapid.IntRange(0, 255).Draw(t, "fgG")),
			B: uint8(rapid.IntRange(0, 255).Draw(t, "fgB")),
			A: uint8(rapid.IntRange(0, 255).Draw(t, "fgA")),
		}

		cfg := textlabel.Config{
			Text:       text,
			Bounds:     bounds,
			Font:       font.Default(),
			Alignment:  alignment,
			Foreground: fg,
		}

		before := textlabel.Config{
			Text:       cfg.Text,
			Bounds:     cfg.Bounds,
			Font:       cfg.Font,
			Alignment:  cfg.Alignment,
			Foreground: cfg.Foreground,
		}

		textlabel.Render(cfg)

		if !reflect.DeepEqual(before, cfg) {
			t.Fatalf("textlabel.Render mutated Config:\nbefore: %+v\nafter:  %+v", before, cfg)
		}
	})
}
