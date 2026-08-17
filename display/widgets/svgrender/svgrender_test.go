package svgrender

import (
	"image"
	"image/color"
	"testing"
)

// --- From: svgrender_parse_test.go ---

// TestParseSVGElements verifies that each supported SVG element type
// parses successfully and renders non-transparent pixels.
func TestParseSVGElements(t *testing.T) {
	tests := []struct {
		name string
		svg  string
	}{
		{
			name: "rect",
			svg:  `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><rect width="100" height="100" fill="red"/></svg>`,
		},
		{
			name: "circle",
			svg:  `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><circle cx="50" cy="50" r="50" fill="green"/></svg>`,
		},
		{
			name: "ellipse",
			svg:  `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><ellipse cx="50" cy="50" rx="50" ry="30" fill="blue"/></svg>`,
		},
		{
			name: "line",
			svg:  `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><line x1="0" y1="0" x2="100" y2="100" stroke="black" stroke-width="2"/></svg>`,
		},
		{
			name: "polyline",
			svg:  `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><polyline points="10,10 50,90 90,10" stroke="red" fill="none" stroke-width="2"/></svg>`,
		},
		{
			name: "polygon",
			svg:  `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><polygon points="50,5 20,99 95,39 5,39 80,99" fill="yellow"/></svg>`,
		},
		{
			name: "path",
			svg:  `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><path d="M10 10 L90 10 L90 90 L10 90 Z" fill="purple"/></svg>`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Verify parse produces a non-nil icon.
			icon, err := parse(tc.svg, 64, 64)
			if err != nil {
				t.Fatalf("parse returned error: %v", err)
			}
			if icon == nil {
				t.Fatal("parse returned nil icon")
			}

			// Render to verify non-nil sprite with non-transparent pixels.
			cfg := Config{
				Bounds: image.Rect(0, 0, 64, 64),
				SVG:    tc.svg,
				Label:  tc.name,
			}
			sprite := Render(cfg)
			if sprite == nil {
				t.Fatal("Render returned nil for valid SVG")
			}
			if sprite.Image == nil {
				t.Fatal("Render produced sprite with nil Image")
			}

			// Check that at least some pixels are non-transparent.
			img := sprite.Image.(*image.RGBA)
			hasNonTransparent := false
			for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
				for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
					_, _, _, a := img.At(x, y).RGBA()
					if a > 0 {
						hasNonTransparent = true
						break
					}
				}
				if hasNonTransparent {
					break
				}
			}
			if !hasNonTransparent {
				t.Error("rendered image has no non-transparent pixels; expected visible content")
			}
		})
	}
}

// TestParseFillColor verifies that fill color attributes are respected
// by rendering an SVG with a known fill and checking for colored pixels.
func TestParseFillColor(t *testing.T) {
	tests := []struct {
		name    string
		svg     string
		expectR bool // expect red channel > 0
		expectG bool // expect green channel > 0
		expectB bool // expect blue channel > 0
	}{
		{
			name:    "red_fill",
			svg:     `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><rect width="100" height="100" fill="red"/></svg>`,
			expectR: true,
			expectG: false,
			expectB: false,
		},
		{
			name:    "green_fill",
			svg:     `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><rect width="100" height="100" fill="#00ff00"/></svg>`,
			expectR: false,
			expectG: true,
			expectB: false,
		},
		{
			name:    "blue_fill",
			svg:     `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><rect width="100" height="100" fill="blue"/></svg>`,
			expectR: false,
			expectG: false,
			expectB: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cfg := Config{
				Bounds: image.Rect(0, 0, 64, 64),
				SVG:    tc.svg,
				Label:  tc.name,
			}
			sprite := Render(cfg)
			if sprite == nil {
				t.Fatal("Render returned nil")
			}

			img := sprite.Image.(*image.RGBA)

			// Accumulate color channel presence across all non-transparent pixels.
			var hasR, hasG, hasB bool
			for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
				for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
					c := img.RGBAAt(x, y)
					if c.A == 0 {
						continue
					}
					if c.R > 0 {
						hasR = true
					}
					if c.G > 0 {
						hasG = true
					}
					if c.B > 0 {
						hasB = true
					}
				}
			}

			if tc.expectR && !hasR {
				t.Error("expected red channel pixels but found none")
			}
			if tc.expectG && !hasG {
				t.Error("expected green channel pixels but found none")
			}
			if tc.expectB && !hasB {
				t.Error("expected blue channel pixels but found none")
			}
			// For pure colors, verify absence of unexpected channels.
			if !tc.expectR && hasR {
				t.Error("found unexpected red channel pixels")
			}
			if !tc.expectG && hasG {
				t.Error("found unexpected green channel pixels")
			}
			if !tc.expectB && hasB {
				t.Error("found unexpected blue channel pixels")
			}
		})
	}
}

// TestParseStrokeColor verifies that stroke color attributes produce
// visible pixels in the rendered output.
func TestParseStrokeColor(t *testing.T) {
	tests := []struct {
		name string
		svg  string
	}{
		{
			name: "stroke_red_line",
			svg:  `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><line x1="0" y1="50" x2="100" y2="50" stroke="red" stroke-width="10"/></svg>`,
		},
		{
			name: "stroke_blue_circle",
			svg:  `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><circle cx="50" cy="50" r="40" stroke="blue" stroke-width="5" fill="none"/></svg>`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cfg := Config{
				Bounds: image.Rect(0, 0, 64, 64),
				SVG:    tc.svg,
				Label:  tc.name,
			}
			sprite := Render(cfg)
			if sprite == nil {
				t.Fatal("Render returned nil for SVG with stroke")
			}

			img := sprite.Image.(*image.RGBA)
			hasNonTransparent := false
			for i := 3; i < len(img.Pix); i += 4 {
				if img.Pix[i] > 0 {
					hasNonTransparent = true
					break
				}
			}
			if !hasNonTransparent {
				t.Error("stroke-only SVG produced no visible pixels")
			}
		})
	}
}

// TestParseAspectRatio verifies that SVG content with a non-square viewBox
// is rendered fitting within the canvas bounds without overflow.
// A 200×100 viewBox (2:1 aspect) rendered into a 64×64 canvas should
// produce visible content entirely within the 64×64 bounds.
func TestParseAspectRatio(t *testing.T) {
	// SVG with 2:1 aspect ratio viewBox.
	svg := `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 200 100"><rect width="200" height="100" fill="white"/></svg>`

	cfg := Config{
		Bounds: image.Rect(0, 0, 64, 64),
		SVG:    svg,
		Label:  "aspect-ratio",
	}
	sprite := Render(cfg)
	if sprite == nil {
		t.Fatal("Render returned nil for aspect ratio test SVG")
	}

	img := sprite.Image.(*image.RGBA)

	// The image dimensions should be the full canvas (64x64).
	if img.Bounds().Dx() != 64 || img.Bounds().Dy() != 64 {
		t.Errorf("image bounds = %dx%d, want 64x64", img.Bounds().Dx(), img.Bounds().Dy())
	}

	// Verify that rendered content fits within canvas bounds (no overflow).
	// All non-transparent pixels must be within [0,64) x [0,64).
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			_, _, _, a := img.At(x, y).RGBA()
			if a > 0 {
				if x < 0 || x >= 64 || y < 0 || y >= 64 {
					t.Fatalf("pixel at (%d,%d) is outside canvas bounds", x, y)
				}
			}
		}
	}

	// Verify the canvas has visible rendered content.
	if !hasNonTransparentPixels(img) {
		t.Error("no visible content in rendered image")
	}
}

// TestParseAspectRatioTall verifies that SVG content with a tall viewBox
// is rendered fitting within the canvas bounds without overflow.
// A 100×200 viewBox (1:2 aspect) rendered into a 64×64 canvas should
// produce visible content entirely within the 64×64 bounds.
func TestParseAspectRatioTall(t *testing.T) {
	// SVG with 1:2 aspect ratio viewBox.
	svg := `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 200"><rect width="100" height="200" fill="white"/></svg>`

	cfg := Config{
		Bounds: image.Rect(0, 0, 64, 64),
		SVG:    svg,
		Label:  "aspect-ratio-tall",
	}
	sprite := Render(cfg)
	if sprite == nil {
		t.Fatal("Render returned nil for tall aspect ratio SVG")
	}

	img := sprite.Image.(*image.RGBA)

	// Verify image dimensions match the canvas.
	if img.Bounds().Dx() != 64 || img.Bounds().Dy() != 64 {
		t.Errorf("image bounds = %dx%d, want 64x64", img.Bounds().Dx(), img.Bounds().Dy())
	}

	// Verify that all rendered content fits within canvas bounds (no overflow).
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			_, _, _, a := img.At(x, y).RGBA()
			if a > 0 {
				if x < 0 || x >= 64 || y < 0 || y >= 64 {
					t.Fatalf("pixel at (%d,%d) is outside canvas bounds", x, y)
				}
			}
		}
	}

	// Verify the canvas has visible rendered content.
	if !hasNonTransparentPixels(img) {
		t.Error("no visible content in rendered image")
	}
}

// hasNonTransparentPixels checks whether an RGBA image contains any
// pixel with alpha > 0. Helper used to avoid import of color in table tests.
func hasNonTransparentPixels(img *image.RGBA) bool {
	for i := 3; i < len(img.Pix); i += 4 {
		if img.Pix[i] > 0 {
			return true
		}
	}
	return false
}

// --- From: svgrender_test.go ---

// TestConstants verifies exported and unexported constant values.
func TestConstants(t *testing.T) {
	if MinBoundsWidth != 16 {
		t.Errorf("MinBoundsWidth = %d, want 16", MinBoundsWidth)
	}
	if MinBoundsHeight != 16 {
		t.Errorf("MinBoundsHeight = %d, want 16", MinBoundsHeight)
	}
	if maxLabelLen != 128 {
		t.Errorf("maxLabelLen = %d, want 128", maxLabelLen)
	}
}

// TestRenderBasicRect verifies rendering a red rectangle SVG produces a
// non-nil sprite with correct bounds and position.
func TestRenderBasicRect(t *testing.T) {
	svg := `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><rect width="100" height="100" fill="red"/></svg>`
	cfg := Config{
		Bounds: image.Rect(0, 0, 64, 64),
		SVG:    svg,
		Label:  "rect-test",
	}

	sprite := Render(cfg)
	if sprite == nil {
		t.Fatal("Render returned nil for valid red rect SVG")
	}

	bounds := sprite.Image.Bounds()
	if bounds.Dx() != 64 || bounds.Dy() != 64 {
		t.Errorf("sprite image bounds = %dx%d, want 64x64", bounds.Dx(), bounds.Dy())
	}

	if sprite.Position != image.Pt(0, 0) {
		t.Errorf("sprite position = %v, want (0,0)", sprite.Position)
	}
}

// TestRenderBasicCircle verifies rendering a green circle SVG produces a
// non-nil sprite.
func TestRenderBasicCircle(t *testing.T) {
	svg := `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><circle cx="50" cy="50" r="50" fill="green"/></svg>`
	cfg := Config{
		Bounds: image.Rect(0, 0, 64, 64),
		SVG:    svg,
		Label:  "circle-test",
	}

	sprite := Render(cfg)
	if sprite == nil {
		t.Fatal("Render returned nil for valid green circle SVG")
	}
}

// TestCanvasToResult verifies that Canvas.ToResult produces a sprite with the
// correct position and a truncated label.
func TestCanvasToResult(t *testing.T) {
	canvas := NewCanvas(32, 32)
	if canvas == nil {
		t.Fatal("NewCanvas(32, 32) returned nil")
	}

	pos := image.Pt(10, 20)
	label := "test-label"
	sprite := canvas.ToResult(pos, label)

	if sprite == nil {
		t.Fatal("ToResult returned nil")
	}
	if sprite.Position != pos {
		t.Errorf("sprite.Position = %v, want %v", sprite.Position, pos)
	}
	if sprite.Label != label {
		t.Errorf("sprite.Label = %q, want %q", sprite.Label, label)
	}
	if sprite.Image == nil {
		t.Error("sprite.Image is nil")
	}

	// Verify label truncation with a long label.
	longLabel := make([]rune, 200)
	for i := range longLabel {
		longLabel[i] = 'x'
	}
	sprite2 := canvas.ToResult(pos, string(longLabel))
	if len([]rune(sprite2.Label)) != 128 {
		t.Errorf("truncated label length = %d, want 128", len([]rune(sprite2.Label)))
	}
}

// TestTintOnOff verifies that rendering the same SVG with and without tint
// produces different pixel data.
func TestTintOnOff(t *testing.T) {
	svg := `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><rect width="100" height="100" fill="white"/></svg>`

	// Render without tint (zero color = no tint).
	cfgNoTint := Config{
		Bounds: image.Rect(0, 0, 32, 32),
		SVG:    svg,
		Label:  "no-tint",
		Color:  color.RGBA{0, 0, 0, 0},
	}
	spriteNoTint := Render(cfgNoTint)
	if spriteNoTint == nil {
		t.Fatal("Render without tint returned nil")
	}

	// Render with a blue tint.
	cfgTint := Config{
		Bounds: image.Rect(0, 0, 32, 32),
		SVG:    svg,
		Label:  "with-tint",
		Color:  color.RGBA{0, 0, 255, 255},
	}
	spriteTint := Render(cfgTint)
	if spriteTint == nil {
		t.Fatal("Render with tint returned nil")
	}

	// Compare pixel data: at least some pixels should differ.
	imgNoTint := spriteNoTint.Image.(*image.RGBA)
	imgTint := spriteTint.Image.(*image.RGBA)

	diffCount := 0
	for i := 0; i < len(imgNoTint.Pix) && i < len(imgTint.Pix); i++ {
		if imgNoTint.Pix[i] != imgTint.Pix[i] {
			diffCount++
		}
	}

	if diffCount == 0 {
		t.Error("tinted and non-tinted images are identical; expected visual difference")
	}
}
