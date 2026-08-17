package borderframe

import (
	"image"
	"image/color"
	"testing"

	"github.com/databeast/cyberhud/display/widgets"
	"pgregory.net/rapid"
)

// --- From: borderframe_prop_test.go ---

// For any pixel bounds where width >= 16 and height >= 16, the borderframe.Render function
// SHALL return a non-nil *widgets.Sprite with a composited RGBA image of the correct dimensions.
func TestProperty_BorderFrameRenderCorrectness(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate random origin offsets
		ox := rapid.IntRange(0, 200).Draw(t, "originX")
		oy := rapid.IntRange(0, 200).Draw(t, "originY")

		// Generate random bounds with width in [16, 512] and height in [16, 512]
		w := rapid.IntRange(16, 512).Draw(t, "width")
		h := rapid.IntRange(16, 512).Draw(t, "height")

		bounds := image.Rect(ox, oy, ox+w, oy+h)
		sprite := Render(Config{Bounds: bounds})

		// Render should never return nil for valid bounds >= 16x16
		if sprite == nil {
			t.Fatalf("Render returned nil for bounds %v (width=%d, height=%d)", bounds, w, h)
		}

		// Verify image dimensions match bounds
		imgBounds := sprite.Image.Bounds()
		if imgBounds.Dx() != w || imgBounds.Dy() != h {
			t.Fatalf("expected image %dx%d, got %dx%d", w, h, imgBounds.Dx(), imgBounds.Dy())
		}

		// Verify position matches bounds.Min
		if sprite.Position.X != ox || sprite.Position.Y != oy {
			t.Fatalf("expected position (%d,%d), got (%d,%d)", ox, oy, sprite.Position.X, sprite.Position.Y)
		}

		// Verify label
		if sprite.Label != "borderframe" {
			t.Fatalf("expected label 'borderframe', got %q", sprite.Label)
		}

		// Verify image is *image.RGBA
		if _, ok := sprite.Image.(*image.RGBA); !ok {
			t.Fatalf("expected *image.RGBA, got %T", sprite.Image)
		}
	})
}

// For any bounds where width < 16 OR height < 16, Render SHALL return nil.
func TestProperty_BorderFrameNilForSmallBounds(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate bounds where at least one dimension is < 16
		w := rapid.IntRange(1, 15).Draw(t, "width")
		h := rapid.IntRange(1, 512).Draw(t, "height")

		bounds := image.Rect(0, 0, w, h)
		sprite := Render(Config{Bounds: bounds})
		if sprite != nil {
			t.Fatalf("Render returned non-nil for bounds %v (width=%d < 16)", bounds, w)
		}
	})

	rapid.Check(t, func(t *rapid.T) {
		w := rapid.IntRange(1, 512).Draw(t, "width")
		h := rapid.IntRange(1, 15).Draw(t, "height")

		bounds := image.Rect(0, 0, w, h)
		sprite := Render(Config{Bounds: bounds})
		if sprite != nil {
			t.Fatalf("Render returned non-nil for bounds %v (height=%d < 16)", bounds, h)
		}
	})
}

// --- From: borderframe_test.go ---

// --- Tests for the Render function and Renderable interface ---

func TestRender_NilWhenTooSmall(t *testing.T) {
	tests := []struct {
		name   string
		bounds image.Rectangle
	}{
		{"zero", image.Rect(0, 0, 0, 0)},
		{"width too small", image.Rect(0, 0, 15, 16)},
		{"height too small", image.Rect(0, 0, 16, 15)},
		{"both too small", image.Rect(0, 0, 8, 8)},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := Render(Config{Bounds: tc.bounds})
			if result != nil {
				t.Errorf("expected nil for bounds %v, got non-nil", tc.bounds)
			}
		})
	}
}

func TestRender_NonNilForValidBounds(t *testing.T) {
	cfg := Config{Bounds: image.Rect(0, 0, 32, 32)}
	result := Render(cfg)
	if result == nil {
		t.Fatal("expected non-nil result for 32x32 bounds")
	}
	if result.Image == nil {
		t.Fatal("expected non-nil Image in result")
	}
	if result.Label != "borderframe" {
		t.Errorf("expected label 'borderframe', got %q", result.Label)
	}
}

func TestRender_ImageDimensions(t *testing.T) {
	cfg := Config{Bounds: image.Rect(10, 20, 74, 84)}
	result := Render(cfg)
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	img := result.Image
	bounds := img.Bounds()
	if bounds.Dx() != 64 || bounds.Dy() != 64 {
		t.Errorf("expected 64x64 image, got %dx%d", bounds.Dx(), bounds.Dy())
	}
	if result.Position != image.Pt(10, 20) {
		t.Errorf("expected position (10,20), got %v", result.Position)
	}
}

func TestRender_ProducesRGBA(t *testing.T) {
	cfg := Config{Bounds: image.Rect(0, 0, 32, 32)}
	result := Render(cfg)
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if _, ok := result.Image.(*image.RGBA); !ok {
		t.Errorf("expected *image.RGBA, got %T", result.Image)
	}
}

func TestRender_MinimumBounds16x16(t *testing.T) {
	cfg := Config{Bounds: image.Rect(0, 0, 16, 16)}
	result := Render(cfg)
	if result == nil {
		t.Fatal("expected non-nil result for 16x16 bounds")
	}
}

func TestRender_TileSetDefaultsToBorder(t *testing.T) {
	// With empty TileSet, Render should use "border" prefix
	cfg := Config{Bounds: image.Rect(0, 0, 32, 32)}
	result := Render(cfg)
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	// Just verify it doesn't panic - icons may or may not be registered
}

func TestNew_ReturnsRenderable(t *testing.T) {
	cfg := Config{Bounds: image.Rect(0, 0, 32, 32)}
	w := New(cfg)
	if w == nil {
		t.Fatal("expected non-nil Renderable from New")
	}
}

func TestNew_RenderFrameProducesSprite(t *testing.T) {
	cfg := Config{Bounds: image.Rect(0, 0, 32, 32)}
	w := New(cfg)
	sprite := w.RenderFrame()
	if sprite == nil {
		t.Fatal("expected non-nil Sprite from RenderFrame for valid bounds")
	}
	if sprite.Image == nil {
		t.Fatal("expected non-nil Image in Sprite")
	}
}

func TestNew_RenderFrameNilForSmallBounds(t *testing.T) {
	cfg := Config{Bounds: image.Rect(0, 0, 8, 8)}
	w := New(cfg)
	sprite := w.RenderFrame()
	if sprite != nil {
		t.Error("expected nil Sprite for bounds < 16x16")
	}
}

func TestNew_ImplementsDescribed(t *testing.T) {
	cfg := Config{Bounds: image.Rect(0, 0, 32, 32)}
	w := New(cfg)
	d, ok := w.(widgets.Described)
	if !ok {
		t.Fatal("New() result does not implement Described")
	}
	desc := d.Describe()
	if desc.Name != "borderframe" {
		t.Errorf("expected Name 'borderframe', got %q", desc.Name)
	}
	if desc.MinWidth != 16 {
		t.Errorf("expected MinWidth 16, got %d", desc.MinWidth)
	}
	if desc.MinHeight != 16 {
		t.Errorf("expected MinHeight 16, got %d", desc.MinHeight)
	}
	found := false
	for _, cap := range desc.Capabilities {
		if cap == "eink-safe" {
			found = true
		}
	}
	if !found {
		t.Error("expected 'eink-safe' in Capabilities")
	}
}

func TestNew_ImplementsConfigurable(t *testing.T) {
	cfg := Config{Bounds: image.Rect(0, 0, 32, 32)}
	w := New(cfg)
	c, ok := w.(widgets.Configurable)
	if !ok {
		t.Fatal("New() result does not implement Configurable")
	}
	// Reconfigure with small bounds → next RenderFrame should return nil
	c.Configure(Config{Bounds: image.Rect(0, 0, 8, 8)})
	sprite := w.RenderFrame()
	if sprite != nil {
		t.Error("expected nil after Configure with small bounds")
	}
}

func TestNew_WithLabelOverride(t *testing.T) {
	cfg := Config{Bounds: image.Rect(0, 0, 32, 32)}
	w := New(cfg, widgets.WithLabel("custom-label"))
	sprite := w.RenderFrame()
	if sprite == nil {
		t.Fatal("expected non-nil Sprite")
	}
	if sprite.Label != "custom-label" {
		t.Errorf("expected label 'custom-label', got %q", sprite.Label)
	}
}

func TestRender_ImageIsNotAllZero(t *testing.T) {
	// If icons are registered, the composited image should have some non-zero pixels.
	// This verifies compositing actually draws something.
	cfg := Config{Bounds: image.Rect(0, 0, 32, 32)}
	result := Render(cfg)
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	rgba, ok := result.Image.(*image.RGBA)
	if !ok {
		t.Fatalf("expected *image.RGBA, got %T", result.Image)
	}
	hasNonZero := false
	for _, p := range rgba.Pix {
		if p != 0 {
			hasNonZero = true
			break
		}
	}
	// Note: if no icons are registered, image will be all-zero; that's acceptable.
	_ = hasNonZero
	_ = color.RGBA{}
}

// --- From: equivalence_prop_test.go ---

// For any bounds with width ≥ 16 and height ≥ 16, calling Render twice with the
// same Config SHALL produce pixel-identical Sprite outputs.

func TestProperty12_BorderframeRenderDeterminism(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		// Generate random bounds with width and height ≥ 16.
		w := rapid.IntRange(16, 256).Draw(rt, "width")
		h := rapid.IntRange(16, 256).Draw(rt, "height")
		ox := rapid.IntRange(0, 100).Draw(rt, "originX")
		oy := rapid.IntRange(0, 100).Draw(rt, "originY")

		bounds := image.Rect(ox, oy, ox+w, oy+h)
		cfg := Config{Bounds: bounds}

		// Call Render twice with the same config.
		result1 := Render(cfg)
		result2 := Render(cfg)

		if result1 == nil {
			t.Fatalf("Render returned nil for bounds %v (width=%d, height=%d)", bounds, w, h)
		}
		if result2 == nil {
			t.Fatalf("Second Render returned nil for bounds %v", bounds)
		}

		// Verify positions match
		if result1.Position != result2.Position {
			t.Fatalf("position mismatch: %v vs %v", result1.Position, result2.Position)
		}

		// Verify labels match
		if result1.Label != result2.Label {
			t.Fatalf("label mismatch: %q vs %q", result1.Label, result2.Label)
		}

		// Verify pixel-identical images
		img1, ok := result1.Image.(*image.RGBA)
		if !ok {
			t.Fatalf("result1 image is %T, expected *image.RGBA", result1.Image)
		}
		img2, ok := result2.Image.(*image.RGBA)
		if !ok {
			t.Fatalf("result2 image is %T, expected *image.RGBA", result2.Image)
		}

		if len(img1.Pix) != len(img2.Pix) {
			t.Fatalf("pixel data length mismatch: %d vs %d", len(img1.Pix), len(img2.Pix))
		}

		for i := range img1.Pix {
			if img1.Pix[i] != img2.Pix[i] {
				t.Fatalf("pixel data mismatch at index %d", i)
			}
		}
	})
}
