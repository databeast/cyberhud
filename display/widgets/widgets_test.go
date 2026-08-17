package widgets

import (
	"fmt"
	"image"
	"testing"

	"github.com/databeast/cyberhud/display/style/layout"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// ---------------------------------------------------------------------------
// Source: cached_renderer_test.go
// ---------------------------------------------------------------------------

type testConfig struct {
	ID int
}

type testResult struct {
	Value string
}

// TestCachedRendererCacheHit verifies that calling Render twice with the same
// Config signature only invokes the underlying render function once, and both
// calls return the same result.
//

func TestCachedRendererCacheHit(t *testing.T) {
	callCount := 0
	render := func(cfg testConfig) *testResult {
		callCount++
		return &testResult{Value: fmt.Sprintf("rendered-%d", cfg.ID)}
	}
	sign := func(cfg testConfig) uint64 {
		return uint64(cfg.ID)
	}

	cr := newCachedRenderer(render, sign)

	cfg := testConfig{ID: 42}

	// First call — should invoke render
	r1 := cr.Render(cfg)
	if callCount != 1 {
		t.Fatalf("expected 1 render call after first Render, got %d", callCount)
	}
	if r1 == nil || r1.Value != "rendered-42" {
		t.Fatalf("unexpected first result: %v", r1)
	}

	// Second call with same config — should return cached result
	r2 := cr.Render(cfg)
	if callCount != 1 {
		t.Fatalf("expected 1 render call after second Render (cache hit), got %d", callCount)
	}
	if r2 != r1 {
		t.Fatalf("expected same pointer on cache hit, got different results")
	}
}

// TestCachedRendererCacheMiss verifies that calling Render with different Config
// signatures invokes the underlying render function each time and returns the
// new result.
//

func TestCachedRendererCacheMiss(t *testing.T) {
	callCount := 0
	render := func(cfg testConfig) *testResult {
		callCount++
		return &testResult{Value: fmt.Sprintf("rendered-%d", cfg.ID)}
	}
	sign := func(cfg testConfig) uint64 {
		return uint64(cfg.ID)
	}

	cr := newCachedRenderer(render, sign)

	// First call
	r1 := cr.Render(testConfig{ID: 1})
	if callCount != 1 {
		t.Fatalf("expected 1 render call, got %d", callCount)
	}
	if r1 == nil || r1.Value != "rendered-1" {
		t.Fatalf("unexpected first result: %v", r1)
	}

	// Second call with different signature — cache miss
	r2 := cr.Render(testConfig{ID: 2})
	if callCount != 2 {
		t.Fatalf("expected 2 render calls after cache miss, got %d", callCount)
	}
	if r2 == nil || r2.Value != "rendered-2" {
		t.Fatalf("unexpected second result: %v", r2)
	}

	// Results should be different pointers with different values
	if r1 == r2 {
		t.Fatalf("expected different pointers on cache miss")
	}

	// Third call with yet another signature — cache miss again
	r3 := cr.Render(testConfig{ID: 3})
	if callCount != 3 {
		t.Fatalf("expected 3 render calls after second cache miss, got %d", callCount)
	}
	if r3 == nil || r3.Value != "rendered-3" {
		t.Fatalf("unexpected third result: %v", r3)
	}
}

// TestCachedRendererValidFlag verifies that the first call to Render always
// invokes the underlying render function, even when the Config signature is 0.
// This ensures the `valid` flag correctly distinguishes "never rendered" from
// "signature is zero."
//

func TestCachedRendererValidFlag(t *testing.T) {
	callCount := 0
	render := func(cfg testConfig) *testResult {
		callCount++
		return &testResult{Value: fmt.Sprintf("rendered-%d", cfg.ID)}
	}
	// Sign function that always returns 0
	sign := func(cfg testConfig) uint64 {
		return 0
	}

	cr := newCachedRenderer(render, sign)

	cfg := testConfig{ID: 99}

	// First call — must invoke render even though signature is 0
	r1 := cr.Render(cfg)
	if callCount != 1 {
		t.Fatalf("expected 1 render call on first invocation (sig=0), got %d", callCount)
	}
	if r1 == nil || r1.Value != "rendered-99" {
		t.Fatalf("unexpected first result: %v", r1)
	}

	// Second call with same config (same signature = 0) — should use cache
	r2 := cr.Render(cfg)
	if callCount != 1 {
		t.Fatalf("expected 1 render call after second invocation (cache hit with sig=0), got %d", callCount)
	}
	if r2 != r1 {
		t.Fatalf("expected same pointer on cache hit, got different results")
	}
}

// ---------------------------------------------------------------------------
// Source: compositor_test.go
// ---------------------------------------------------------------------------

// mockRenderable is a simple Renderable that returns a fixed sprite.
type mockRenderable struct {
	sprite *Sprite
	called bool
}

func (m *mockRenderable) RenderFrame() *Sprite {
	m.called = true
	return m.sprite
}

// mockDescribedRenderable implements both Renderable and Described.
type mockDescribedRenderable struct {
	sprite     *Sprite
	descriptor Descriptor
	called     bool
}

func (m *mockDescribedRenderable) RenderFrame() *Sprite {
	m.called = true
	return m.sprite
}

func (m *mockDescribedRenderable) Describe() Descriptor {
	return m.descriptor
}

func testSprite(label string) *Sprite {
	return &Sprite{
		Image:  image.NewRGBA(image.Rect(0, 0, 10, 10)),
		Label:  label,
		Bounds: image.Rect(0, 0, 10, 10),
	}
}

func TestCompositor_Add_AppendsNonNilSprites(t *testing.T) {
	comp := NewCompositor(SuppressionContext{})
	r := &mockRenderable{sprite: testSprite("a")}

	comp.Add(r)

	if !r.called {
		t.Fatal("expected RenderFrame to be called")
	}
	sprites := comp.Sprites()
	if len(sprites) != 1 {
		t.Fatalf("expected 1 sprite, got %d", len(sprites))
	}
	if sprites[0].Label != "a" {
		t.Fatalf("expected label 'a', got %q", sprites[0].Label)
	}
}

func TestCompositor_Add_DiscardsNilResults(t *testing.T) {
	comp := NewCompositor(SuppressionContext{})
	r := &mockRenderable{sprite: nil}

	comp.Add(r)

	if !r.called {
		t.Fatal("expected RenderFrame to be called even when result is nil")
	}
	if len(comp.Sprites()) != 0 {
		t.Fatalf("expected 0 sprites, got %d", len(comp.Sprites()))
	}
}

func TestCompositor_AddIf_FalseSkipsInvocation(t *testing.T) {
	comp := NewCompositor(SuppressionContext{})
	r := &mockRenderable{sprite: testSprite("skipped")}

	comp.AddIf(false, r)

	if r.called {
		t.Fatal("expected RenderFrame NOT to be called when condition is false")
	}
	if len(comp.Sprites()) != 0 {
		t.Fatalf("expected 0 sprites, got %d", len(comp.Sprites()))
	}
}

func TestCompositor_AddIf_TrueInvokes(t *testing.T) {
	comp := NewCompositor(SuppressionContext{})
	r := &mockRenderable{sprite: testSprite("included")}

	comp.AddIf(true, r)

	if !r.called {
		t.Fatal("expected RenderFrame to be called when condition is true")
	}
	if len(comp.Sprites()) != 1 {
		t.Fatalf("expected 1 sprite, got %d", len(comp.Sprites()))
	}
}

func TestCompositor_InsertionOrder(t *testing.T) {
	comp := NewCompositor(SuppressionContext{})
	comp.Add(&mockRenderable{sprite: testSprite("first")})
	comp.Add(&mockRenderable{sprite: testSprite("second")})
	comp.Add(&mockRenderable{sprite: testSprite("third")})

	sprites := comp.Sprites()
	if len(sprites) != 3 {
		t.Fatalf("expected 3 sprites, got %d", len(sprites))
	}
	expected := []string{"first", "second", "third"}
	for i, s := range sprites {
		if s.Label != expected[i] {
			t.Errorf("sprite[%d]: expected label %q, got %q", i, expected[i], s.Label)
		}
	}
}

func TestCompositor_Suppression_SkipsMatchingWidgets(t *testing.T) {
	alwaysSuppress := func(_ Descriptor, _ SuppressionContext) bool { return true }
	comp := NewCompositor(SuppressionContext{}, alwaysSuppress)

	r := &mockDescribedRenderable{
		sprite:     testSprite("suppressed"),
		descriptor: Descriptor{Name: "test"},
	}

	comp.Add(r)

	if r.called {
		t.Fatal("expected RenderFrame NOT to be called when suppressed")
	}
	if len(comp.Sprites()) != 0 {
		t.Fatalf("expected 0 sprites, got %d", len(comp.Sprites()))
	}
}

func TestCompositor_Suppression_DoesNotAffectNonDescribed(t *testing.T) {
	alwaysSuppress := func(_ Descriptor, _ SuppressionContext) bool { return true }
	comp := NewCompositor(SuppressionContext{}, alwaysSuppress)

	// mockRenderable does NOT implement Described, so suppression doesn't apply
	r := &mockRenderable{sprite: testSprite("not-described")}

	comp.Add(r)

	if !r.called {
		t.Fatal("expected RenderFrame to be called for non-Described widget")
	}
	if len(comp.Sprites()) != 1 {
		t.Fatalf("expected 1 sprite, got %d", len(comp.Sprites()))
	}
}

func TestCompositor_Suppression_ShortCircuitOR(t *testing.T) {
	rule1Called := false
	rule2Called := false

	rule1 := func(_ Descriptor, _ SuppressionContext) bool {
		rule1Called = true
		return true // suppresses
	}
	rule2 := func(_ Descriptor, _ SuppressionContext) bool {
		rule2Called = true
		return true
	}

	comp := NewCompositor(SuppressionContext{}, rule1, rule2)
	r := &mockDescribedRenderable{
		sprite:     testSprite("test"),
		descriptor: Descriptor{Name: "test"},
	}

	comp.Add(r)

	if !rule1Called {
		t.Fatal("expected rule1 to be called")
	}
	if rule2Called {
		t.Fatal("expected rule2 NOT to be called (short-circuit)")
	}
}

func TestCompositor_NoRules_RendersAll(t *testing.T) {
	comp := NewCompositor(SuppressionContext{})
	r := &mockDescribedRenderable{
		sprite:     testSprite("rendered"),
		descriptor: Descriptor{Name: "test"},
	}

	comp.Add(r)

	if !r.called {
		t.Fatal("expected RenderFrame to be called with no rules")
	}
	if len(comp.Sprites()) != 1 {
		t.Fatalf("expected 1 sprite, got %d", len(comp.Sprites()))
	}
}

func TestCompositor_NoOverlapEnforcement(t *testing.T) {
	comp := NewCompositor(SuppressionContext{})

	// Two widgets with the same bounds (overlapping)
	s1 := &Sprite{Image: image.NewRGBA(image.Rect(0, 0, 50, 50)), Bounds: image.Rect(0, 0, 50, 50), Label: "a"}
	s2 := &Sprite{Image: image.NewRGBA(image.Rect(0, 0, 50, 50)), Bounds: image.Rect(0, 0, 50, 50), Label: "b"}

	comp.Add(&mockRenderable{sprite: s1})
	comp.Add(&mockRenderable{sprite: s2})

	sprites := comp.Sprites()
	if len(sprites) != 2 {
		t.Fatalf("expected 2 sprites (no collision detection), got %d", len(sprites))
	}
}

// ---------------------------------------------------------------------------
// Source: placement_test.go
// ---------------------------------------------------------------------------

// testBridge creates a LayoutBridge with known dimensions for placement tests.
// With PaddingPct=0 and 200x150 panel:
//   - contentOriginX = 0, contentOriginY = 0 (no title bar, no inset)
//   - availableContentWidth = 200
//   - availableContentHeight = 150
func testBridge() layout.LayoutCalculator {
	hints := textlayout.TextHints{
		PixelWidth:   200,
		PixelHeight:  150,
		RowHeight:    10,
		GlyphAdvance: 6,
		GlyphHeight:  8,
	}
	cfg := layout.BridgeConfig{
		PaddingPct: 0,
	}
	return layout.NewLayoutBridge(hints, cfg)
}

func TestPlaceTopRight(t *testing.T) {
	bridge := testBridge()

	// Normal case: element fits
	got := PlaceTopRight(bridge, 40, 20)
	want := image.Rectangle{
		Min: image.Point{X: 0 + 200 - 40, Y: 0},
		Max: image.Point{X: 0 + 200, Y: 0 + 20},
	}
	if got != want {
		t.Errorf("PlaceTopRight(40, 20) = %v, want %v", got, want)
	}

	// Overflow: element too wide
	got = PlaceTopRight(bridge, 210, 20)
	if got != (image.Rectangle{}) {
		t.Errorf("PlaceTopRight(210, 20) = %v, want empty rect", got)
	}

	// Overflow: element too tall
	got = PlaceTopRight(bridge, 40, 200)
	if got != (image.Rectangle{}) {
		t.Errorf("PlaceTopRight(40, 200) = %v, want empty rect", got)
	}
}

func TestPlaceBottom(t *testing.T) {
	bridge := testBridge()

	// Normal case: element fits
	got := PlaceBottom(bridge, 30)
	want := image.Rectangle{
		Min: image.Point{X: 0, Y: 0 + 150 - 30},
		Max: image.Point{X: 0 + 200, Y: 0 + 150},
	}
	if got != want {
		t.Errorf("PlaceBottom(30) = %v, want %v", got, want)
	}

	// Overflow: element too tall
	got = PlaceBottom(bridge, 200)
	if got != (image.Rectangle{}) {
		t.Errorf("PlaceBottom(200) = %v, want empty rect", got)
	}
}

func TestPlaceBottomRight(t *testing.T) {
	bridge := testBridge()

	// Normal case: element fits
	got := PlaceBottomRight(bridge, 50, 25)
	want := image.Rectangle{
		Min: image.Point{X: 0 + 200 - 50, Y: 0 + 150 - 25},
		Max: image.Point{X: 0 + 200, Y: 0 + 150},
	}
	if got != want {
		t.Errorf("PlaceBottomRight(50, 25) = %v, want %v", got, want)
	}

	// Overflow: element too wide
	got = PlaceBottomRight(bridge, 210, 25)
	if got != (image.Rectangle{}) {
		t.Errorf("PlaceBottomRight(210, 25) = %v, want empty rect", got)
	}

	// Overflow: element too tall
	got = PlaceBottomRight(bridge, 50, 200)
	if got != (image.Rectangle{}) {
		t.Errorf("PlaceBottomRight(50, 200) = %v, want empty rect", got)
	}
}

// ---------------------------------------------------------------------------
// Source: registry_test.go
// ---------------------------------------------------------------------------

// stubDescribed is a minimal Described implementation for registry tests.
type stubDescribed struct{}

func (s stubDescribed) Describe() Descriptor {
	return Descriptor{Name: "stub"}
}

func TestRegister_DuplicatePanics(t *testing.T) {
	// Save and restore registry state to avoid polluting other tests.
	orig := registry
	registry = map[string]func() Described{}
	defer func() { registry = orig }()

	factory := func() Described { return stubDescribed{} }
	Register("alpha", factory)

	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic on duplicate registration, got none")
		}
		msg, ok := r.(string)
		if !ok {
			t.Fatalf("expected string panic, got %T: %v", r, r)
		}
		if msg == "" {
			t.Fatal("panic message should be descriptive")
		}
	}()

	Register("alpha", factory) // should panic
}

func TestRegistered_ReturnsSortedNames(t *testing.T) {
	orig := registry
	registry = map[string]func() Described{}
	defer func() { registry = orig }()

	factory := func() Described { return stubDescribed{} }
	Register("zulu", factory)
	Register("alpha", factory)
	Register("mike", factory)

	names := Registered()
	if len(names) != 3 {
		t.Fatalf("expected 3 names, got %d", len(names))
	}
	expected := []string{"alpha", "mike", "zulu"}
	for i, name := range names {
		if name != expected[i] {
			t.Errorf("index %d: expected %q, got %q", i, expected[i], name)
		}
	}
}

func TestRegistered_EmptyRegistry(t *testing.T) {
	orig := registry
	registry = map[string]func() Described{}
	defer func() { registry = orig }()

	names := Registered()
	if len(names) != 0 {
		t.Fatalf("expected 0 names for empty registry, got %d", len(names))
	}
}

// ---------------------------------------------------------------------------
// Source: suppression_test.go
// ---------------------------------------------------------------------------

func TestSuppressOnEink_NonEinkContext(t *testing.T) {
	rule := SuppressOnEink()
	desc := Descriptor{Name: "test", Capabilities: []string{}}
	ctx := SuppressionContext{IsEink: false}

	if rule(desc, ctx) {
		t.Error("SuppressOnEink should not suppress when IsEink is false")
	}
}

func TestSuppressOnEink_EinkWithoutCapability(t *testing.T) {
	rule := SuppressOnEink()
	desc := Descriptor{Name: "test", Capabilities: []string{"animated"}}
	ctx := SuppressionContext{IsEink: true}

	if !rule(desc, ctx) {
		t.Error("SuppressOnEink should suppress when IsEink is true and widget lacks eink-safe")
	}
}

func TestSuppressOnEink_EinkWithCapability(t *testing.T) {
	rule := SuppressOnEink()
	desc := Descriptor{Name: "test", Capabilities: []string{"eink-safe"}}
	ctx := SuppressionContext{IsEink: true}

	if rule(desc, ctx) {
		t.Error("SuppressOnEink should not suppress when widget has eink-safe capability")
	}
}

func TestSuppressOnEink_EinkWithNilCapabilities(t *testing.T) {
	rule := SuppressOnEink()
	desc := Descriptor{Name: "test", Capabilities: nil}
	ctx := SuppressionContext{IsEink: true}

	if !rule(desc, ctx) {
		t.Error("SuppressOnEink should suppress when IsEink is true and Capabilities is nil")
	}
}

func TestSuppressBelow_WidgetFits(t *testing.T) {
	rule := SuppressBelow(100, 80)
	desc := Descriptor{Name: "test", MinWidth: 50, MinHeight: 40}
	ctx := SuppressionContext{}

	if rule(desc, ctx) {
		t.Error("SuppressBelow should not suppress when widget fits within available dimensions")
	}
}

func TestSuppressBelow_WidgetTooWide(t *testing.T) {
	rule := SuppressBelow(100, 80)
	desc := Descriptor{Name: "test", MinWidth: 101, MinHeight: 40}
	ctx := SuppressionContext{}

	if !rule(desc, ctx) {
		t.Error("SuppressBelow should suppress when widget MinWidth exceeds available width")
	}
}

func TestSuppressBelow_WidgetTooTall(t *testing.T) {
	rule := SuppressBelow(100, 80)
	desc := Descriptor{Name: "test", MinWidth: 50, MinHeight: 81}
	ctx := SuppressionContext{}

	if !rule(desc, ctx) {
		t.Error("SuppressBelow should suppress when widget MinHeight exceeds available height")
	}
}

func TestSuppressBelow_ExactFit(t *testing.T) {
	rule := SuppressBelow(100, 80)
	desc := Descriptor{Name: "test", MinWidth: 100, MinHeight: 80}
	ctx := SuppressionContext{}

	if rule(desc, ctx) {
		t.Error("SuppressBelow should not suppress when widget dimensions exactly equal available")
	}
}

func TestSuppressBelow_BothExceed(t *testing.T) {
	rule := SuppressBelow(100, 80)
	desc := Descriptor{Name: "test", MinWidth: 200, MinHeight: 160}
	ctx := SuppressionContext{}

	if !rule(desc, ctx) {
		t.Error("SuppressBelow should suppress when both dimensions exceed available")
	}
}
