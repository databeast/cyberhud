package widgets

import (
	"image"
	"image/draw"
	"testing"

	"github.com/databeast/cyberhud/display/style/layout"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"pgregory.net/rapid"
)

// ---------------------------------------------------------------------------
// --- From: placement_prop_test.go ---
// ---------------------------------------------------------------------------

// For any LayoutBridge with content area (originX, originY, width, height) and for any
// element dimensions (elemW, elemH):
// - PlaceTopRight SHALL return a rectangle with Max.X = originX + width and Min.Y = originY.
// - PlaceBottom SHALL return a rectangle with Min.X = originX, Max.X = originX + width, and Max.Y = originY + height.
// - PlaceBottomRight SHALL return a rectangle with Max.X = originX + width and Max.Y = originY + height.
// - When elemW > width or elemH > height, all placement functions SHALL return image.Rectangle{}.

// genBridge generates a LayoutBridge with random panel dimensions.
func genBridge(t *rapid.T) layout.LayoutCalculator {
	pixelWidth := rapid.IntRange(20, 800).Draw(t, "pixelWidth")
	pixelHeight := rapid.IntRange(20, 480).Draw(t, "pixelHeight")

	hints := textlayout.TextHints{
		PixelWidth:   pixelWidth,
		PixelHeight:  pixelHeight,
		RowHeight:    10,
		GlyphAdvance: 6,
		GlyphHeight:  8,
	}
	cfg := layout.BridgeConfig{}
	return layout.NewLayoutBridge(hints, cfg)
}

// TestProperty15_PlaceTopRight_Anchors verifies that PlaceTopRight returns a rectangle
// with Max.X = originX + width and Min.Y = originY for any fitting element.
func TestProperty15_PlaceTopRight_Anchors(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		bridge := genBridge(rt)

		w := bridge.AvailableContentWidth()
		h := bridge.AvailableContentHeight()
		if w <= 0 || h <= 0 {
			return // skip degenerate panels with no content area
		}

		elemW := rapid.IntRange(1, w).Draw(rt, "elemW")
		elemH := rapid.IntRange(1, h).Draw(rt, "elemH")

		originX, originY := bridge.ContentOrigin()
		got := PlaceTopRight(bridge, elemW, elemH)

		if got.Max.X != originX+w {
			t.Fatalf("PlaceTopRight(%d, %d): Max.X = %d, want %d (originX=%d, width=%d)",
				elemW, elemH, got.Max.X, originX+w, originX, w)
		}
		if got.Min.Y != originY {
			t.Fatalf("PlaceTopRight(%d, %d): Min.Y = %d, want %d (originY=%d)",
				elemW, elemH, got.Min.Y, originY, originY)
		}
	})
}

// TestProperty15_PlaceBottom_Anchors verifies that PlaceBottom returns a rectangle
// with Min.X = originX, Max.X = originX + width, and Max.Y = originY + height.
func TestProperty15_PlaceBottom_Anchors(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		bridge := genBridge(rt)

		w := bridge.AvailableContentWidth()
		h := bridge.AvailableContentHeight()
		if w <= 0 || h <= 0 {
			return // skip degenerate panels
		}

		elemH := rapid.IntRange(1, h).Draw(rt, "elemH")

		originX, originY := bridge.ContentOrigin()
		got := PlaceBottom(bridge, elemH)

		if got.Min.X != originX {
			t.Fatalf("PlaceBottom(%d): Min.X = %d, want %d (originX=%d)",
				elemH, got.Min.X, originX, originX)
		}
		if got.Max.X != originX+w {
			t.Fatalf("PlaceBottom(%d): Max.X = %d, want %d (originX=%d, width=%d)",
				elemH, got.Max.X, originX+w, originX, w)
		}
		if got.Max.Y != originY+h {
			t.Fatalf("PlaceBottom(%d): Max.Y = %d, want %d (originY=%d, height=%d)",
				elemH, got.Max.Y, originY+h, originY, h)
		}
	})
}

// TestProperty15_PlaceBottomRight_Anchors verifies that PlaceBottomRight returns a
// rectangle with Max.X = originX + width and Max.Y = originY + height.
func TestProperty15_PlaceBottomRight_Anchors(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		bridge := genBridge(rt)

		w := bridge.AvailableContentWidth()
		h := bridge.AvailableContentHeight()
		if w <= 0 || h <= 0 {
			return // skip degenerate panels
		}

		elemW := rapid.IntRange(1, w).Draw(rt, "elemW")
		elemH := rapid.IntRange(1, h).Draw(rt, "elemH")

		originX, originY := bridge.ContentOrigin()
		got := PlaceBottomRight(bridge, elemW, elemH)

		if got.Max.X != originX+w {
			t.Fatalf("PlaceBottomRight(%d, %d): Max.X = %d, want %d (originX=%d, width=%d)",
				elemW, elemH, got.Max.X, originX+w, originX, w)
		}
		if got.Max.Y != originY+h {
			t.Fatalf("PlaceBottomRight(%d, %d): Max.Y = %d, want %d (originY=%d, height=%d)",
				elemW, elemH, got.Max.Y, originY+h, originY, h)
		}
	})
}

// TestProperty15_Overflow_ReturnsZero verifies that when element dimensions exceed
// the content area, all placement functions return image.Rectangle{}.
func TestProperty15_Overflow_ReturnsZero(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		bridge := genBridge(rt)

		w := bridge.AvailableContentWidth()
		h := bridge.AvailableContentHeight()
		if w <= 0 || h <= 0 {
			return // skip degenerate panels
		}

		// Generate element dimensions where at least one exceeds the content area.
		overflowWidth := rapid.Bool().Draw(rt, "overflowWidth")
		var elemW, elemH int
		if overflowWidth {
			elemW = rapid.IntRange(w+1, w+200).Draw(rt, "elemW")
			elemH = rapid.IntRange(1, h+200).Draw(rt, "elemH")
		} else {
			elemW = rapid.IntRange(1, w+200).Draw(rt, "elemW")
			elemH = rapid.IntRange(h+1, h+200).Draw(rt, "elemH")
		}

		zero := image.Rectangle{}

		gotTopRight := PlaceTopRight(bridge, elemW, elemH)
		if gotTopRight != zero {
			t.Fatalf("PlaceTopRight(%d, %d) with content %dx%d = %v, want empty rect",
				elemW, elemH, w, h, gotTopRight)
		}

		gotBottom := PlaceBottom(bridge, elemH)
		// PlaceBottom only checks height overflow
		if elemH > h {
			if gotBottom != zero {
				t.Fatalf("PlaceBottom(%d) with content height %d = %v, want empty rect",
					elemH, h, gotBottom)
			}
		}

		gotBottomRight := PlaceBottomRight(bridge, elemW, elemH)
		if gotBottomRight != zero {
			t.Fatalf("PlaceBottomRight(%d, %d) with content %dx%d = %v, want empty rect",
				elemW, elemH, w, h, gotBottomRight)
		}
	})
}

// ---------------------------------------------------------------------------
// --- From: registry_prop_test.go ---
// ---------------------------------------------------------------------------

// For any registered widget type, constructing an instance via its factory and
// calling RenderFrame() with valid bounds (≥ Descriptor.MinWidth × MinHeight)
// returns non-nil *Sprite with non-nil Image.

// testWidget is a minimal widget implementing Renderable, Described, and Configurable.
// It is used to validate the property test machinery before real widgets are registered.
type testWidget struct {
	bounds image.Rectangle
}

func (w *testWidget) RenderFrame() *Sprite {
	width := w.bounds.Dx()
	height := w.bounds.Dy()
	if width < 4 || height < 4 {
		return nil
	}
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	return &Sprite{
		Image:  img,
		Bounds: w.bounds,
		Label:  "test-widget",
	}
}

func (w *testWidget) Describe() Descriptor {
	return Descriptor{
		Name:      "test-widget",
		MinWidth:  4,
		MinHeight: 4,
	}
}

func (w *testWidget) Configure(cfg interface{}) {
	if r, ok := cfg.(image.Rectangle); ok {
		w.bounds = r
	}
}

func TestProperty1_RegistryRenderContract(t *testing.T) {
	// Save and restore registry state to avoid polluting other tests.
	orig := registry
	registry = map[string]func() Described{}
	defer func() { registry = orig }()

	// Register a test widget so the property can exercise the full path.
	Register("test-widget", func() Described {
		return &testWidget{}
	})

	rapid.Check(t, func(rt *rapid.T) {
		names := Registered()
		if len(names) == 0 {
			t.Skip("no widgets registered")
		}

		// Pick a registered widget type.
		name := rapid.SampledFrom(names).Draw(rt, "widgetType")
		factory := registry[name]
		instance := factory()

		desc := instance.Describe()

		// Generate valid bounds: width ≥ MinWidth, height ≥ MinHeight.
		// Clamp minimums to at least 1 to avoid degenerate cases.
		minW := desc.MinWidth
		if minW < 1 {
			minW = 1
		}
		minH := desc.MinHeight
		if minH < 1 {
			minH = 1
		}

		width := rapid.IntRange(minW, minW+200).Draw(rt, "width")
		height := rapid.IntRange(minH, minH+200).Draw(rt, "height")
		bounds := image.Rect(0, 0, width, height)

		// If the widget implements Configurable, configure it with valid bounds.
		if c, ok := instance.(Configurable); ok {
			c.Configure(bounds)
		}

		// The widget must also implement Renderable to call RenderFrame.
		r, ok := instance.(Renderable)
		if !ok {
			t.Fatalf("widget %q factory returned a Described that does not implement Renderable", name)
		}

		sprite := r.RenderFrame()
		if sprite == nil {
			t.Fatalf("widget %q: RenderFrame() returned nil with valid bounds %v (min %dx%d)",
				name, bounds, desc.MinWidth, desc.MinHeight)
		}
		if sprite.Image == nil {
			t.Fatalf("widget %q: RenderFrame() returned Sprite with nil Image for bounds %v",
				name, bounds)
		}
	})
}

// For any registered widget type and for any configuration where bounds width ≤ 0,
// bounds height ≤ 0, or bounds are below the widget's minimum dimensions,
// RenderFrame() SHALL return nil.

func TestProperty2_InvalidBoundsNilGuard(t *testing.T) {
	// Save and restore registry state to avoid polluting other tests.
	orig := registry
	registry = map[string]func() Described{}
	defer func() { registry = orig }()

	// Register the test widget so the property can exercise the full path.
	Register("test-widget", func() Described {
		return &testWidget{}
	})

	rapid.Check(t, func(rt *rapid.T) {
		names := Registered()
		if len(names) == 0 {
			t.Skip("no widgets registered")
		}

		// Pick a registered widget type.
		name := rapid.SampledFrom(names).Draw(rt, "widgetType")
		factory := registry[name]
		instance := factory()

		desc := instance.Describe()

		// Generate invalid bounds: one of three cases:
		// 0 = zero/negative width, 1 = zero/negative height, 2 = below minimum dimensions
		invalidCase := rapid.IntRange(0, 2).Draw(rt, "invalidCase")

		var bounds image.Rectangle
		switch invalidCase {
		case 0: // zero or negative width
			w := rapid.IntRange(-10, 0).Draw(rt, "width")
			h := rapid.IntRange(desc.MinHeight, desc.MinHeight+100).Draw(rt, "height")
			bounds = image.Rectangle{
				Min: image.Point{X: 0, Y: 0},
				Max: image.Point{X: w, Y: h},
			}
		case 1: // zero or negative height
			w := rapid.IntRange(desc.MinWidth, desc.MinWidth+100).Draw(rt, "width")
			h := rapid.IntRange(-10, 0).Draw(rt, "height")
			bounds = image.Rectangle{
				Min: image.Point{X: 0, Y: 0},
				Max: image.Point{X: w, Y: h},
			}
		case 2: // below minimum dimensions (positive but too small)
			// At least one dimension must be below the minimum.
			// Pick which dimension is below minimum.
			dimChoice := rapid.IntRange(0, 2).Draw(rt, "dimChoice")
			var w, h int
			switch dimChoice {
			case 0: // width below min
				if desc.MinWidth > 1 {
					w = rapid.IntRange(1, desc.MinWidth-1).Draw(rt, "width")
				} else {
					// MinWidth is 1, so we can't go below min with positive values.
					// Fall back to zero width.
					w = 0
				}
				h = rapid.IntRange(desc.MinHeight, desc.MinHeight+100).Draw(rt, "height")
			case 1: // height below min
				w = rapid.IntRange(desc.MinWidth, desc.MinWidth+100).Draw(rt, "width")
				if desc.MinHeight > 1 {
					h = rapid.IntRange(1, desc.MinHeight-1).Draw(rt, "height")
				} else {
					// MinHeight is 1, so we can't go below min with positive values.
					// Fall back to zero height.
					h = 0
				}
			default: // both below min
				if desc.MinWidth > 1 {
					w = rapid.IntRange(1, desc.MinWidth-1).Draw(rt, "width")
				} else {
					w = 0
				}
				if desc.MinHeight > 1 {
					h = rapid.IntRange(1, desc.MinHeight-1).Draw(rt, "height")
				} else {
					h = 0
				}
			}
			bounds = image.Rectangle{
				Min: image.Point{X: 0, Y: 0},
				Max: image.Point{X: w, Y: h},
			}
		}

		// Configure the widget with invalid bounds.
		if c, ok := instance.(Configurable); ok {
			c.Configure(bounds)
		}

		// The widget must implement Renderable.
		r, ok := instance.(Renderable)
		if !ok {
			t.Fatalf("widget %q factory returned a Described that does not implement Renderable", name)
		}

		sprite := r.RenderFrame()
		if sprite != nil {
			t.Fatalf("widget %q: RenderFrame() returned non-nil with invalid bounds %v (min %dx%d)",
				name, bounds, desc.MinWidth, desc.MinHeight)
		}
	})
}

// For any registered non-animated widget and for any valid configuration,
// calling RenderFrame() twice with the same configuration SHALL produce
// Sprite values with pixel-identical Images, equal Positions, equal Bounds,
// and equal Labels.

func TestProperty3_RenderDeterminism(t *testing.T) {
	// Save and restore registry state to avoid polluting other tests.
	orig := registry
	registry = map[string]func() Described{}
	defer func() { registry = orig }()

	// Register the test widget (not animated) so the property can exercise the full path.
	Register("test-widget", func() Described {
		return &testWidget{}
	})

	rapid.Check(t, func(rt *rapid.T) {
		names := Registered()
		if len(names) == 0 {
			t.Skip("no widgets registered")
		}

		// Pick a registered widget type.
		name := rapid.SampledFrom(names).Draw(rt, "widgetType")
		factory := registry[name]
		instance := factory()

		// Skip animated widgets — they may produce different output each call.
		if _, animated := instance.(Animated); animated {
			rt.Skip("skipping animated widget")
		}

		desc := instance.Describe()

		// Generate valid bounds: width ≥ MinWidth, height ≥ MinHeight.
		minW := desc.MinWidth
		if minW < 1 {
			minW = 1
		}
		minH := desc.MinHeight
		if minH < 1 {
			minH = 1
		}

		width := rapid.IntRange(minW, minW+200).Draw(rt, "width")
		height := rapid.IntRange(minH, minH+200).Draw(rt, "height")
		bounds := image.Rect(0, 0, width, height)

		// Configure the widget with valid bounds.
		if c, ok := instance.(Configurable); ok {
			c.Configure(bounds)
		}

		r, ok := instance.(Renderable)
		if !ok {
			t.Fatalf("widget %q factory returned a Described that does not implement Renderable", name)
		}

		// First render.
		sprite1 := r.RenderFrame()
		// Second render with same config.
		sprite2 := r.RenderFrame()

		if sprite1 == nil && sprite2 == nil {
			return // Both nil is deterministic.
		}
		if (sprite1 == nil) != (sprite2 == nil) {
			t.Fatalf("widget %q: RenderFrame() non-deterministic nil: first=%v, second=%v",
				name, sprite1 == nil, sprite2 == nil)
		}

		// Compare Labels.
		if sprite1.Label != sprite2.Label {
			t.Fatalf("widget %q: Label mismatch: %q vs %q", name, sprite1.Label, sprite2.Label)
		}

		// Compare Positions.
		if sprite1.Position != sprite2.Position {
			t.Fatalf("widget %q: Position mismatch: %v vs %v", name, sprite1.Position, sprite2.Position)
		}

		// Compare Bounds.
		if sprite1.Bounds != sprite2.Bounds {
			t.Fatalf("widget %q: Bounds mismatch: %v vs %v", name, sprite1.Bounds, sprite2.Bounds)
		}

		// Compare Images pixel-by-pixel.
		if sprite1.Image == nil && sprite2.Image == nil {
			return
		}
		if (sprite1.Image == nil) != (sprite2.Image == nil) {
			t.Fatalf("widget %q: Image nil mismatch", name)
		}

		b1 := sprite1.Image.Bounds()
		b2 := sprite2.Image.Bounds()
		if b1 != b2 {
			t.Fatalf("widget %q: Image bounds differ: %v vs %v", name, b1, b2)
		}

		// Convert both to RGBA for reliable pixel comparison.
		rgba1 := image.NewRGBA(b1)
		draw.Draw(rgba1, b1, sprite1.Image, b1.Min, draw.Src)
		rgba2 := image.NewRGBA(b2)
		draw.Draw(rgba2, b2, sprite2.Image, b2.Min, draw.Src)

		if len(rgba1.Pix) != len(rgba2.Pix) {
			t.Fatalf("widget %q: pixel buffer lengths differ: %d vs %d",
				name, len(rgba1.Pix), len(rgba2.Pix))
		}
		for i := range rgba1.Pix {
			if rgba1.Pix[i] != rgba2.Pix[i] {
				t.Fatalf("widget %q: pixel data differs at byte offset %d", name, i)
			}
		}
	})
}

// For any registered widget type and for any valid configuration producing a non-nil
// Sprite, the Sprite's Image SHALL be of type *image.RGBA or *image.NRGBA.

func TestProperty4_RGBAOutputFormat(t *testing.T) {
	// Save and restore registry state to avoid polluting other tests.
	orig := registry
	registry = map[string]func() Described{}
	defer func() { registry = orig }()

	// Register the test widget so the property can exercise the full path.
	Register("test-widget", func() Described {
		return &testWidget{}
	})

	rapid.Check(t, func(rt *rapid.T) {
		names := Registered()
		if len(names) == 0 {
			t.Skip("no widgets registered")
		}

		// Pick a registered widget type.
		name := rapid.SampledFrom(names).Draw(rt, "widgetType")
		factory := registry[name]
		instance := factory()

		desc := instance.Describe()

		// Generate valid bounds: width ≥ MinWidth, height ≥ MinHeight.
		minW := desc.MinWidth
		if minW < 1 {
			minW = 1
		}
		minH := desc.MinHeight
		if minH < 1 {
			minH = 1
		}

		width := rapid.IntRange(minW, minW+200).Draw(rt, "width")
		height := rapid.IntRange(minH, minH+200).Draw(rt, "height")
		bounds := image.Rect(0, 0, width, height)

		// Configure the widget with valid bounds.
		if c, ok := instance.(Configurable); ok {
			c.Configure(bounds)
		}

		r, ok := instance.(Renderable)
		if !ok {
			t.Fatalf("widget %q factory returned a Described that does not implement Renderable", name)
		}

		sprite := r.RenderFrame()
		if sprite == nil {
			// Widget returned nil (valid for this config) — nothing to assert.
			return
		}
		if sprite.Image == nil {
			t.Fatalf("widget %q: non-nil Sprite has nil Image", name)
		}

		// Assert the Image is *image.RGBA or *image.NRGBA.
		switch sprite.Image.(type) {
		case *image.RGBA:
			// OK
		case *image.NRGBA:
			// OK
		default:
			t.Fatalf("widget %q: Image type is %T, expected *image.RGBA or *image.NRGBA",
				name, sprite.Image)
		}
	})
}

// ---------------------------------------------------------------------------
// --- From: suppression_prop_test.go ---
// ---------------------------------------------------------------------------

// For any Descriptor and SuppressionContext where IsEink is true,
// SuppressOnEink() returns true iff Capabilities do not contain "eink-safe".

// genCapabilities generates a random capability slice that may or may not include "eink-safe".
func genCapabilities(t *rapid.T) []string {
	otherCaps := []string{"animated", "transparent", "high-contrast", "vector", "scalable"}
	n := rapid.IntRange(0, 4).Draw(t, "numOtherCaps")
	caps := make([]string, 0, n+1)
	for i := 0; i < n; i++ {
		idx := rapid.IntRange(0, len(otherCaps)-1).Draw(t, "capIdx")
		caps = append(caps, otherCaps[idx])
	}
	includeEinkSafe := rapid.Bool().Draw(t, "includeEinkSafe")
	if includeEinkSafe {
		caps = append(caps, "eink-safe")
	}
	return caps
}

// hasEinkSafe checks whether a capabilities slice contains "eink-safe".
func hasEinkSafe(caps []string) bool {
	for _, c := range caps {
		if c == "eink-safe" {
			return true
		}
	}
	return false
}

// TestProperty9_SuppressOnEink_Correctness verifies that for any Descriptor and
// SuppressionContext, SuppressOnEink() returns true iff (IsEink AND "eink-safe" not in Capabilities).
func TestProperty9_SuppressOnEink_Correctness(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		caps := genCapabilities(rt)
		desc := Descriptor{
			Name:         rapid.StringMatching(`[a-z]{3,10}`).Draw(rt, "name"),
			MinWidth:     rapid.IntRange(1, 500).Draw(rt, "minWidth"),
			MinHeight:    rapid.IntRange(1, 500).Draw(rt, "minHeight"),
			Capabilities: caps,
		}
		isEink := rapid.Bool().Draw(rt, "isEink")
		ctx := SuppressionContext{
			IsEink:          isEink,
			AvailableWidth:  rapid.IntRange(10, 800).Draw(rt, "availWidth"),
			AvailableHeight: rapid.IntRange(10, 480).Draw(rt, "availHeight"),
		}

		rule := SuppressOnEink()
		got := rule(desc, ctx)

		// Expected: suppress iff IsEink AND widget lacks "eink-safe"
		want := isEink && !hasEinkSafe(caps)

		if got != want {
			t.Fatalf("SuppressOnEink()(desc{Capabilities: %v}, ctx{IsEink: %v}) = %v, want %v",
				caps, isEink, got, want)
		}
	})
}

// TestProperty9_SuppressOnEink_NilCapabilities verifies the rule also works correctly
// when Capabilities is nil (should suppress on e-ink since "eink-safe" cannot be present).
func TestProperty9_SuppressOnEink_NilCapabilities(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		desc := Descriptor{
			Name:         rapid.StringMatching(`[a-z]{3,10}`).Draw(rt, "name"),
			MinWidth:     rapid.IntRange(1, 500).Draw(rt, "minWidth"),
			MinHeight:    rapid.IntRange(1, 500).Draw(rt, "minHeight"),
			Capabilities: nil,
		}
		isEink := rapid.Bool().Draw(rt, "isEink")
		ctx := SuppressionContext{IsEink: isEink}

		rule := SuppressOnEink()
		got := rule(desc, ctx)

		// With nil capabilities, "eink-safe" is never present, so suppress iff IsEink.
		want := isEink

		if got != want {
			t.Fatalf("SuppressOnEink()(desc{Capabilities: nil}, ctx{IsEink: %v}) = %v, want %v",
				isEink, got, want)
		}
	})
}

// For any Descriptor with MinWidth W and MinHeight H, SuppressBelow(AW, AH)
// returns true iff W > AW or H > AH.

// TestProperty10_SuppressBelow_Correctness verifies that for any Descriptor with
// MinWidth W and MinHeight H, SuppressBelow(AW, AH) returns true iff W > AW or H > AH.
func TestProperty10_SuppressBelow_Correctness(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		minW := rapid.IntRange(0, 1000).Draw(rt, "minWidth")
		minH := rapid.IntRange(0, 1000).Draw(rt, "minHeight")
		desc := Descriptor{
			Name:         rapid.StringMatching(`[a-z]{3,10}`).Draw(rt, "name"),
			MinWidth:     minW,
			MinHeight:    minH,
			Capabilities: nil,
		}

		availW := rapid.IntRange(0, 1000).Draw(rt, "availWidth")
		availH := rapid.IntRange(0, 1000).Draw(rt, "availHeight")

		rule := SuppressBelow(availW, availH)
		ctx := SuppressionContext{}
		got := rule(desc, ctx)

		// Expected: suppress iff MinWidth > availW OR MinHeight > availH
		want := minW > availW || minH > availH

		if got != want {
			t.Fatalf("SuppressBelow(%d, %d)(desc{MinWidth: %d, MinHeight: %d}) = %v, want %v",
				availW, availH, minW, minH, got, want)
		}
	})
}
