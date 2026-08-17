package ggrender

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"strings"
	"testing"

	"github.com/databeast/cyberhud/display/widgets"
)

// --- From: font_test.go ---

func TestLoadFontNonExistentPath(t *testing.T) {
	_, err := LoadFont("/nonexistent/path/font.ttf", 12.0)
	if err == nil {
		t.Fatal("expected error for non-existent font path")
	}
	if !strings.Contains(err.Error(), "ggrender: load font") {
		t.Fatalf("error should contain 'ggrender: load font', got: %v", err)
	}
	if !strings.Contains(err.Error(), "/nonexistent/path/font.ttf") {
		t.Fatalf("error should contain the font path, got: %v", err)
	}
}

func TestListFontsAfterLoading(t *testing.T) {
	fontPath := setupTestFont(t)
	size := 16.0
	_, err := LoadFont(fontPath, size)
	if err != nil {
		t.Fatalf("LoadFont failed: %v", err)
	}
	fonts := ListFonts()
	expected := fmt.Sprintf("%s:%.2f", fontPath, size)
	found := false
	for _, id := range fonts {
		if id == expected {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("ListFonts() = %v, want to contain %q", fonts, expected)
	}
}

// --- From: integration_test.go ---

func TestFullPipelineIntegration(t *testing.T) {
	// 1. Create Canvas
	c := NewCanvas(128, 128)
	if c == nil {
		t.Fatal("NewCanvas returned nil")
	}

	// 2. Draw shapes
	red := color.RGBA{R: 255, A: 255}
	c.FillRect(0, 0, 64, 64, red)
	c.FillCircle(96, 96, 20, color.RGBA{G: 255, A: 255})
	c.StrokeLine(0, 0, 127, 127, 2, color.RGBA{B: 255, A: 255})

	// 3. Draw text (load font first)
	fontPath := setupTestFont(t)
	f, err := LoadFont(fontPath, 16.0)
	if err != nil {
		t.Fatalf("LoadFont: %v", err)
	}
	err = c.DrawText("Hello", 10, 64, f, color.RGBA{R: 255, G: 255, B: 255, A: 255}, AlignLeft)
	if err != nil {
		t.Fatalf("DrawText: %v", err)
	}

	// 4. Convert to Sprite
	pos := image.Point{X: 10, Y: 20}
	result := c.ToResult(pos, "integration-test")
	if result == nil {
		t.Fatal("ToResult returned nil")
	}
	if result.Image == nil {
		t.Fatal("result.Image is nil")
	}
	if result.Position != pos {
		t.Fatalf("position mismatch: %v != %v", result.Position, pos)
	}
	if result.Label != "integration-test" {
		t.Fatalf("label mismatch: %q", result.Label)
	}
}

// --- From: shapes_test.go ---

// TestFillRect verifies that FillRect draws a filled rectangle at the expected
// location and leaves other regions transparent.

func TestFillRect(t *testing.T) {
	c := NewCanvas(64, 64)
	if c == nil {
		t.Fatal("NewCanvas returned nil for valid dimensions")
	}

	red := color.RGBA{R: 255, G: 0, B: 0, A: 255}
	c.FillRect(0, 0, 32, 32, red)

	img := c.Image()

	// Pixel (16,16) should be red (inside the rectangle)
	r, g, b, a := img.At(16, 16).RGBA()
	if a == 0 {
		t.Fatal("pixel (16,16) should be non-transparent after FillRect")
	}
	// Compare against red (RGBA values are pre-multiplied and scaled to 16-bit)
	if r != 0xffff || g != 0 || b != 0 {
		t.Fatalf("pixel (16,16) expected red, got RGBA(%d,%d,%d,%d)", r, g, b, a)
	}

	// Pixel (48,48) should be transparent (outside the rectangle)
	_, _, _, a = img.At(48, 48).RGBA()
	if a != 0 {
		t.Fatalf("pixel (48,48) should be transparent, got alpha=%d", a)
	}
}

// TestFillRoundedRect verifies that FillRoundedRect fills the interior region
// of a rounded rectangle with the specified color.

func TestFillRoundedRect(t *testing.T) {
	c := NewCanvas(64, 64)
	if c == nil {
		t.Fatal("NewCanvas returned nil for valid dimensions")
	}

	green := color.RGBA{R: 0, G: 255, B: 0, A: 255}
	c.FillRoundedRect(8, 8, 48, 48, 8, green)

	img := c.Image()

	// Pixel (32,32) is in the center of the rounded rect and should be green
	r, g, b, a := img.At(32, 32).RGBA()
	if a == 0 {
		t.Fatal("pixel (32,32) should be non-transparent after FillRoundedRect")
	}
	if r != 0 || g != 0xffff || b != 0 {
		t.Fatalf("pixel (32,32) expected green, got RGBA(%d,%d,%d,%d)", r, g, b, a)
	}
}

// TestStrokeLine verifies that StrokeLine draws a visible line along the
// expected diagonal path.

func TestStrokeLine(t *testing.T) {
	c := NewCanvas(64, 64)
	if c == nil {
		t.Fatal("NewCanvas returned nil for valid dimensions")
	}

	white := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	c.StrokeLine(0, 0, 63, 63, 3, white)

	img := c.Image()

	// Pixel (32,32) is on the diagonal and should be non-transparent
	_, _, _, a := img.At(32, 32).RGBA()
	if a == 0 {
		t.Fatal("pixel (32,32) should be non-transparent after StrokeLine along diagonal")
	}
}

// TestFillArc verifies that FillArc draws a filled pie wedge and produces
// non-transparent pixels in the expected region.

func TestFillArc(t *testing.T) {
	c := NewCanvas(64, 64)
	if c == nil {
		t.Fatal("NewCanvas returned nil for valid dimensions")
	}

	blue := color.RGBA{R: 0, G: 0, B: 255, A: 255}
	// Draw arc at center (32,32), radius 20, from 0 to π/2 (first quadrant)
	c.FillArc(32, 32, 20, 0, math.Pi/2, blue)

	img := c.Image()

	// The center of the arc wedge in the first quadrant (0 to π/2) should be
	// roughly at an angle of π/4 from center, at about half the radius.
	// That's approximately (32 + 10*cos(π/4), 32 + 10*sin(π/4)) ≈ (39, 39)
	_, _, _, a := img.At(39, 39).RGBA()
	if a == 0 {
		t.Fatal("pixel (39,39) should be non-transparent after FillArc in first quadrant")
	}
}

// --- From: sign_test.go ---

func TestSignRenderCacheCompatibility(t *testing.T) {
	// This test verifies at compile time that Render and Sign are
	// type-compatible with widgets.NewRenderCache[Config, widgets.Sprite].
	cr := widgets.NewRenderCache[Config, widgets.Sprite](Render, Sign)
	if cr == nil {
		t.Fatal("NewRenderCache returned nil")
	}
}

// --- From: text_test.go ---

func TestDrawTextNilFontReturnsError(t *testing.T) {
	c := NewCanvas(64, 64)
	if c == nil {
		t.Fatal("NewCanvas returned nil")
	}
	col := color.RGBA{R: 255, A: 255}
	err := c.DrawText("hello", 10, 10, nil, col, AlignLeft)
	if err == nil {
		t.Fatal("expected error for nil font")
	}
	if !strings.Contains(err.Error(), "no font set") {
		t.Fatalf("expected 'no font set' in error, got: %v", err)
	}
}

func TestMeasureTextPositive(t *testing.T) {
	fontPath := setupTestFont(t)
	f, err := LoadFont(fontPath, 24.0)
	if err != nil {
		t.Fatalf("LoadFont: %v", err)
	}
	c := NewCanvas(256, 64)
	if c == nil {
		t.Fatal("NewCanvas returned nil")
	}
	w, h, err := c.MeasureText("Hello World", f)
	if err != nil {
		t.Fatalf("MeasureText: %v", err)
	}
	if w <= 0 {
		t.Fatalf("expected width > 0, got %f", w)
	}
	if h <= 0 {
		t.Fatalf("expected height > 0, got %f", h)
	}
}
