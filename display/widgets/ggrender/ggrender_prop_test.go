package ggrender

import (
	"image"
	"image/color"
	"math"
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/image/font/gofont/goregular"
	"pgregory.net/rapid"
)

// --- From: canvas_prop_test.go ---

// TestCanvasNilGuard validates Property 2: Invalid Dimensions Produce Nil Canvas.
//
// For any width ≤ 0 or height ≤ 0, NewCanvas(width, height) SHALL return nil.

func TestCanvasNilGuard(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate at least one non-positive dimension (width non-positive).
		w := rapid.IntRange(-100, 0).Draw(t, "width")
		h := rapid.IntRange(-100, 100).Draw(t, "height")
		c := NewCanvas(w, h)
		if c != nil {
			t.Fatalf("NewCanvas(%d, %d) should be nil for non-positive width", w, h)
		}
	})
}

// TestCanvasNilGuardHeightNonPositive validates the complementary case where
// width is positive but height is non-positive.

func TestCanvasNilGuardHeightNonPositive(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Width is positive, height is non-positive.
		w := rapid.IntRange(1, 100).Draw(t, "width")
		h := rapid.IntRange(-100, 0).Draw(t, "height")
		c := NewCanvas(w, h)
		if c != nil {
			t.Fatalf("NewCanvas(%d, %d) should be nil for non-positive height", w, h)
		}
	})
}

func TestCanvasTransparentBackground(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		w := rapid.IntRange(1, 128).Draw(t, "width")
		h := rapid.IntRange(1, 128).Draw(t, "height")
		c := NewCanvas(w, h)
		if c == nil {
			t.Fatalf("NewCanvas(%d, %d) returned nil", w, h)
		}
		img := c.Image()
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				r, g, b, a := img.At(x, y).RGBA()
				if r != 0 || g != 0 || b != 0 || a != 0 {
					t.Fatalf("pixel (%d,%d) not transparent: RGBA(%d,%d,%d,%d)", x, y, r, g, b, a)
				}
			}
		}
	})
}

// TestToResultLabelTruncation validates Property 7: Label Truncation.
//
// For any label string with length > 128 runes, ToResult SHALL produce a
// result whose Label has exactly 128 runes.

func TestToResultLabelTruncation(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate a label with 129-300 alphanumeric runes (always > 128)
		label := rapid.StringMatching(`[a-zA-Z0-9]{129,300}`).Draw(t, "label")
		if len([]rune(label)) <= 128 {
			t.Fatalf("precondition: generated label should be > 128 runes, got %d", len([]rune(label)))
		}
		c := NewCanvas(32, 32)
		if c == nil {
			t.Fatal("NewCanvas(32,32) returned nil")
		}
		result := c.ToResult(image.Point{}, label)
		got := len([]rune(result.Label))
		if got != 128 {
			t.Fatalf("expected label length 128, got %d (input length %d)", got, len([]rune(label)))
		}
	})
}

// TestToResultFieldFidelity validates Property 8: ToResult Field Fidelity.
//
// For random positions and labels (≤ 128 chars), verify the result fields match inputs.

func TestToResultFieldFidelity(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		x := rapid.IntRange(-1000, 1000).Draw(t, "x")
		y := rapid.IntRange(-1000, 1000).Draw(t, "y")
		label := rapid.StringMatching(`[a-zA-Z0-9_\-/]{0,128}`).Draw(t, "label")
		pos := image.Point{X: x, Y: y}

		c := NewCanvas(64, 64)
		if c == nil {
			t.Fatal("NewCanvas(64,64) returned nil")
		}
		result := c.ToResult(pos, label)
		if result.Position != pos {
			t.Fatalf("position %v != expected %v", result.Position, pos)
		}
		if result.Label != label {
			t.Fatalf("label %q != expected %q", result.Label, label)
		}
		if result.Image == nil {
			t.Fatal("result.Image is nil")
		}
	})
}

// --- From: font_prop_test.go ---

// setupTestFont writes the Go Regular TTF font to a temporary file and returns its path.
func setupTestFont(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	fontPath := filepath.Join(dir, "GoRegular.ttf")
	if err := os.WriteFile(fontPath, goregular.TTF, 0644); err != nil {
		t.Fatalf("failed to write test font: %v", err)
	}
	return fontPath
}

// For any valid font file path and point size, calling LoadFont N times with
// the same arguments SHALL return the same *Font pointer every time.

func TestFontCacheIdempotence(t *testing.T) {
	fontPath := setupTestFont(t)

	rapid.Check(t, func(t *rapid.T) {
		size := rapid.Float64Range(8.0, 72.0).Draw(t, "size")

		f1, err := LoadFont(fontPath, size)
		if err != nil {
			t.Fatalf("LoadFont failed: %v", err)
		}

		n := rapid.IntRange(2, 10).Draw(t, "n")
		for i := 0; i < n; i++ {
			f2, err := LoadFont(fontPath, size)
			if err != nil {
				t.Fatalf("LoadFont (iteration %d) failed: %v", i, err)
			}
			if f1 != f2 {
				t.Fatalf("LoadFont returned different pointer on iteration %d for size %.2f", i, size)
			}
		}
	})
}

// --- From: render_prop_test.go ---

// TestRenderNilForInvalidBounds validates Property 9: Render Nil for Invalid Bounds.
//
// For random Configs with zero or negative dimensions, Render SHALL return nil.

func TestRenderNilForInvalidBounds(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate bounds with zero or negative dimensions.
		// Use image.Rectangle literal to avoid image.Rect normalization.
		x := rapid.IntRange(0, 100).Draw(t, "x")
		y := rapid.IntRange(0, 100).Draw(t, "y")
		// Make width zero or negative (w <= 0)
		w := rapid.IntRange(-50, 0).Draw(t, "w")
		h := rapid.IntRange(-50, 100).Draw(t, "h")

		bounds := image.Rectangle{
			Min: image.Point{X: x, Y: y},
			Max: image.Point{X: x + w, Y: y + h},
		}

		cfg := Config{
			Bounds: bounds,
			Label:  "test",
			Color:  color.RGBA{R: 255, A: 255},
		}
		result := Render(cfg)
		if result != nil {
			t.Fatalf("Render should return nil for invalid bounds Dx=%d Dy=%d", cfg.Bounds.Dx(), cfg.Bounds.Dy())
		}
	})
}

// TestMinimumBoundsGuard validates Property 14: Minimum Bounds Guard.
//
// For random Configs with width < 16, verify Render returns nil.

func TestMinimumBoundsGuard(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate bounds where width is < 16
		w := rapid.IntRange(1, 15).Draw(t, "w")
		h := rapid.IntRange(1, 512).Draw(t, "h")

		cfg := Config{
			Bounds: image.Rect(0, 0, w, h),
			Label:  "test",
			Color:  color.RGBA{R: 255, A: 255},
		}
		result := Render(cfg)
		if result != nil {
			t.Fatalf("Render should return nil for bounds %dx%d (width < 16)", w, h)
		}
	})
}

// TestMinimumBoundsGuardHeight validates Property 14: Minimum Bounds Guard.
//
// For random Configs with height < 16 (and width >= 16), verify Render returns nil.

func TestMinimumBoundsGuardHeight(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate bounds where width is valid but height is < 16
		w := rapid.IntRange(16, 512).Draw(t, "w")
		h := rapid.IntRange(1, 15).Draw(t, "h")

		cfg := Config{
			Bounds: image.Rect(0, 0, w, h),
			Label:  "test",
			Color:  color.RGBA{R: 255, A: 255},
		}
		result := Render(cfg)
		if result != nil {
			t.Fatalf("Render should return nil for bounds %dx%d (height < 16)", w, h)
		}
	})
}

// --- From: shapes_prop_test.go ---

// TestShapeClippingSafety validates Property 5 from the design: Out-of-Bounds Shapes Clip Without Error.
//
// For any Canvas with valid dimensions and for any shape whose dimensions or position
// exceed the Canvas bounds, drawing the shape SHALL complete without panic and the
// resulting image bounds SHALL remain unchanged (equal to the original Canvas dimensions).

func TestShapeClippingSafety(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		w := rapid.IntRange(16, 128).Draw(t, "width")
		h := rapid.IntRange(16, 128).Draw(t, "height")
		c := NewCanvas(w, h)
		if c == nil {
			t.Fatalf("NewCanvas(%d,%d) returned nil", w, h)
		}

		// Draw shapes with extreme out-of-bounds coordinates
		col := color.RGBA{R: 255, A: 255}
		bigX := rapid.Float64Range(-10000, 10000).Draw(t, "x")
		bigY := rapid.Float64Range(-10000, 10000).Draw(t, "y")
		bigR := rapid.Float64Range(-1000, 1000).Draw(t, "r")

		c.FillRect(bigX, bigY, bigR, bigR, col)
		c.FillRoundedRect(bigX, bigY, bigR, bigR, bigR, col)
		c.FillCircle(bigX, bigY, bigR, col)
		c.StrokeLine(bigX, bigY, bigR, bigR, 5.0, col)
		c.FillArc(bigX, bigY, bigR, 0, 6.28, col)

		// Verify bounds unchanged
		bounds := c.Image().Bounds()
		expected := image.Rect(0, 0, w, h)
		if bounds != expected {
			t.Fatalf("bounds changed: %v != %v", bounds, expected)
		}
	})
}

// For canvas >= 32x32 with circle at center, verify center pixel is non-transparent.

func TestCircleCenterPixel(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		w := rapid.IntRange(32, 256).Draw(t, "width")
		h := rapid.IntRange(32, 256).Draw(t, "height")
		c := NewCanvas(w, h)
		if c == nil {
			t.Fatalf("NewCanvas(%d,%d) returned nil", w, h)
		}

		cx := float64(w) / 2
		cy := float64(h) / 2
		radius := float64(rapid.IntRange(4, min(w, h)/2).Draw(t, "radius"))
		col := color.RGBA{R: 255, G: 0, B: 0, A: 255}

		c.FillCircle(cx, cy, radius, col)

		// Check center pixel is non-transparent
		_, _, _, a := c.Image().At(w/2, h/2).RGBA()
		if a == 0 {
			t.Fatalf("center pixel (%d,%d) is transparent after drawing circle at center with radius %.0f", w/2, h/2, radius)
		}
	})
}

// --- From: sign_prop_test.go ---

// TestSignSensitivity validates Property 12: Sign Sensitivity.
//
// For random Config pairs differing in exactly one field, Sign SHALL produce
// different uint64 values (collision rate < 0.1%).

func TestSignSensitivity(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		cfg1 := Config{
			Bounds: image.Rect(
				rapid.IntRange(0, 100).Draw(t, "minX"),
				rapid.IntRange(0, 100).Draw(t, "minY"),
				rapid.IntRange(100, 500).Draw(t, "maxX"),
				rapid.IntRange(100, 500).Draw(t, "maxY"),
			),
			Label: rapid.StringMatching(`[a-zA-Z0-9]{1,64}`).Draw(t, "label"),
			Color: color.RGBA{
				R: rapid.Byte().Draw(t, "r"),
				G: rapid.Byte().Draw(t, "g"),
				B: rapid.Byte().Draw(t, "b"),
				A: rapid.Byte().Draw(t, "a"),
			},
		}

		// Create cfg2 differing in exactly one field
		cfg2 := cfg1
		field := rapid.IntRange(0, 2).Draw(t, "field")
		switch field {
		case 0:
			cfg2.Bounds.Max.X += rapid.IntRange(1, 100).Draw(t, "delta")
		case 1:
			cfg2.Label += "x"
		case 2:
			cfg2.Color.R = cfg1.Color.R + 1
		}

		s1 := Sign(cfg1)
		s2 := Sign(cfg2)
		if s1 == s2 {
			t.Fatalf("Sign collision: both %d for configs differing in field %d", s1, field)
		}
	})
}

// --- From: text_prop_test.go ---

// For any non-empty text string, MeasureText produces a consistent, deterministic
// advance width that closely matches the actual rendered ink extent. The advance
// width (sum of glyph advances) and the ink extent (rightmost rendered pixel) may
// differ slightly due to glyph overhang or right-side bearing, but the difference
// must remain small and bounded.
//
// Properties validated:
//   - Determinism: MeasureText called twice returns identical values
//   - Closeness: |advance - ink| is small (within a font-size-relative tolerance)

func TestTextMeasurementConsistency(t *testing.T) {
	fontPath := setupTestFont(t)
	f, err := LoadFont(fontPath, 24.0)
	if err != nil {
		t.Fatalf("LoadFont: %v", err)
	}

	rapid.Check(t, func(t *rapid.T) {
		// Generate strings that start and end with visible characters to avoid
		// trailing-space edge cases (trailing spaces have advance but no ink).
		text := rapid.StringMatching(`[a-zA-Z0-9][a-zA-Z0-9 ]{0,48}[a-zA-Z0-9]`).Draw(t, "text")

		c := NewCanvas(1024, 128)
		if c == nil {
			t.Fatal("NewCanvas returned nil")
		}

		// Property: Determinism — two calls with same input produce identical results
		w1, h1, err := c.MeasureText(text, f)
		if err != nil {
			t.Fatalf("MeasureText (first call): %v", err)
		}
		w2, h2, err := c.MeasureText(text, f)
		if err != nil {
			t.Fatalf("MeasureText (second call): %v", err)
		}
		if w1 != w2 || h1 != h2 {
			t.Fatalf("MeasureText not deterministic: first=(%.2f, %.2f) second=(%.2f, %.2f)", w1, h1, w2, h2)
		}

		if w1 <= 0 {
			t.Fatalf("measured width %.2f should be > 0", w1)
		}
		if h1 <= 0 {
			t.Fatalf("measured height %.2f should be > 0", h1)
		}

		// Draw the text and find the rightmost non-transparent pixel (ink extent)
		col := color.RGBA{R: 255, G: 255, B: 255, A: 255}
		err = c.DrawText(text, 0, 64, f, col, AlignLeft)
		if err != nil {
			t.Fatalf("DrawText: %v", err)
		}

		img := c.Image()
		rightmost := -1
		for x := 0; x < 1024; x++ {
			for y := 0; y < 128; y++ {
				_, _, _, a := img.At(x, y).RGBA()
				if a > 0 && x > rightmost {
					rightmost = x
				}
			}
		}

		if rightmost < 0 {
			t.Fatalf("no ink pixels found for text %q", text)
		}

		inkWidth := float64(rightmost + 1)

		// Property: Closeness — the advance width and ink extent should be
		// very close. They can differ because:
		//   - Right-side bearing makes advance > ink
		//   - Glyph overhang (e.g., italic, "Q" tail) makes ink > advance
		// For a 24pt font, a tolerance of 3px handles both cases.
		diff := math.Abs(w1 - inkWidth)
		if diff > 3 {
			t.Fatalf("|advance - ink| = %.2f > 3px for text %q (advance=%.2f, ink=%.2f)", diff, text, w1, inkWidth)
		}
	})
}
