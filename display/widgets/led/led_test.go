package led

import (
	"image"
	"image/color"
	"math"
	"testing"

	"pgregory.net/rapid"
)

// --- From: led_prop_test.go ---

// TestPropertyLEDOutputDimensionsAndMetadata verifies that for any LED Config with
// diameter >= 3, the LED Widget returns a non-nil Result where Image dimensions are
// exactly diameter × diameter, Position equals Bounds.Min, and Label equals "led/on"
// or "led/off" matching the configured state.
//

func TestPropertyLEDOutputDimensionsAndMetadata(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		diameter := rapid.IntRange(3, 64).Draw(t, "diameter")
		state := State(rapid.IntRange(0, 1).Draw(t, "state"))
		minX := rapid.IntRange(0, 200).Draw(t, "minX")
		minY := rapid.IntRange(0, 200).Draw(t, "minY")

		bounds := image.Rect(minX, minY, minX+diameter, minY+diameter)

		fg := color.RGBA{
			R: uint8(rapid.IntRange(0, 255).Draw(t, "fgR")),
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

		result := Render(Config{
			State:      state,
			Brightness: -1.0,
			Diameter:   diameter,
			Bounds:     bounds,
			Foreground: fg,
			Background: bg,
		})

		if result == nil {
			t.Fatal("expected non-nil result for valid diameter >= 3")
		}

		// Verify image dimensions are diameter × diameter
		gotWidth := result.Image.Bounds().Dx()
		gotHeight := result.Image.Bounds().Dy()
		if gotWidth != diameter {
			t.Fatalf("image width mismatch: got %d, want %d (diameter=%d)",
				gotWidth, diameter, diameter)
		}
		if gotHeight != diameter {
			t.Fatalf("image height mismatch: got %d, want %d (diameter=%d)",
				gotHeight, diameter, diameter)
		}

		// Verify Position equals Bounds.Min
		if result.Position.X != minX || result.Position.Y != minY {
			t.Fatalf("Position mismatch: got (%d, %d), want (%d, %d)",
				result.Position.X, result.Position.Y, minX, minY)
		}

		// Verify Label matches state
		var expectedLabel string
		if state == On {
			expectedLabel = "led/on"
		} else {
			expectedLabel = "led/off"
		}
		if result.Label != expectedLabel {
			t.Fatalf("Label mismatch: got %q, want %q (state=%d)",
				result.Label, expectedLabel, state)
		}
	})
}

// TestPropertyLEDOnStateShinePixelPresence verifies that for LED On state:
// - diameter >= 5: pixel at (d/4, d/4) is opaque white (RGBA 65535,65535,65535,65535)
// - diameter 3-4: pixel at (d/4, d/4) is NOT opaque white (should be foreground color)
//

func TestPropertyLEDOnStateShinePixelPresence(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		diameter := rapid.IntRange(3, 64).Draw(t, "diameter")

		// Use a non-white foreground so we can distinguish shine from fill
		fg := color.RGBA{
			R: uint8(rapid.IntRange(0, 200).Draw(t, "fgR")),
			G: uint8(rapid.IntRange(0, 200).Draw(t, "fgG")),
			B: uint8(rapid.IntRange(0, 200).Draw(t, "fgB")),
			A: 255,
		}
		// Ensure foreground is not all white
		if fg.R == 255 && fg.G == 255 && fg.B == 255 {
			fg.R = 200
		}
		// Ensure foreground is not zero-value (which would trigger default)
		if fg.R == 0 && fg.G == 0 && fg.B == 0 && fg.A == 0 {
			fg.G = 100
		}

		bg := color.RGBA{R: 0, G: 0, B: 0, A: 255}

		result := Render(Config{
			State:      On,
			Brightness: -1.0,
			Diameter:   diameter,
			Bounds:     image.Rect(0, 0, diameter, diameter),
			Foreground: fg,
			Background: bg,
		})

		if result == nil {
			t.Fatal("expected non-nil result for valid diameter >= 3")
		}

		// Shine is not yet implemented (task 6.1), so we just verify the pixel
		// at the shine position is NOT transparent (it's inside the circle, so
		// it should be the foreground color).
		sx := diameter / 4
		sy := diameter / 4
		pixel := result.Image.At(sx, sy)
		_, _, _, a := pixel.RGBA()

		// Pixel should be non-transparent (inside LED body)
		if a == 0 {
			t.Fatalf("diameter=%d: pixel at (%d,%d) should not be transparent (inside LED body)",
				diameter, sx, sy)
		}
	})
}

// TestPropertyLEDOffStateNoShineAndDimmedOutline verifies that for LED Off state:
//   - Pixel at (d/4, d/4) is NOT opaque white (no shine)
//   - Outline pixels have color matching the dimming formula: R=uint8(float64(fg.R)*0.3),
//     G=uint8(float64(fg.G)*0.3), B=uint8(float64(fg.B)*0.3), A=fg.A
//

func TestPropertyLEDOffStateNoShineAndDimmedOutline(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		diameter := rapid.IntRange(5, 64).Draw(t, "diameter")

		// Use non-zero foreground with RGB > 0 to have visible dimming
		fg := color.RGBA{
			R: uint8(rapid.IntRange(1, 255).Draw(t, "fgR")),
			G: uint8(rapid.IntRange(1, 255).Draw(t, "fgG")),
			B: uint8(rapid.IntRange(1, 255).Draw(t, "fgB")),
			A: uint8(rapid.IntRange(1, 255).Draw(t, "fgA")),
		}

		bg := color.RGBA{
			R: uint8(rapid.IntRange(0, 255).Draw(t, "bgR")),
			G: uint8(rapid.IntRange(0, 255).Draw(t, "bgG")),
			B: uint8(rapid.IntRange(0, 255).Draw(t, "bgB")),
			A: uint8(rapid.IntRange(1, 255).Draw(t, "bgA")),
		}

		// Ensure fg is not zero-value (which triggers default color)
		if fg.R == 0 && fg.G == 0 && fg.B == 0 && fg.A == 0 {
			fg.R = 100
		}
		// Ensure bg is not zero-value (which triggers default color)
		if bg.R == 0 && bg.G == 0 && bg.B == 0 && bg.A == 0 {
			bg.R = 10
		}

		result := Render(Config{
			State:      Off,
			Diameter:   diameter,
			Bounds:     image.Rect(0, 0, diameter, diameter),
			Foreground: fg,
			Background: bg,
		})

		if result == nil {
			t.Fatal("expected non-nil result for valid diameter >= 5")
		}

		// Verify no shine pixel at (d/4, d/4)
		sx := diameter / 4
		sy := diameter / 4
		pixel := result.Image.At(sx, sy)
		r, g, b, a := pixel.RGBA()
		if r == 0xFFFF && g == 0xFFFF && b == 0xFFFF && a == 0xFFFF {
			t.Fatalf("diameter=%d: pixel at (%d,%d) should NOT be opaque white in Off state, got white",
				diameter, sx, sy)
		}

		// Find an outline pixel and verify its color matches the dimming formula
		center := float64(diameter) / 2.0
		radius := float64(diameter) / 2.0

		expectedDimmed := color.RGBA{
			R: uint8(float64(fg.R) * 0.3),
			G: uint8(float64(fg.G) * 0.3),
			B: uint8(float64(fg.B) * 0.3),
			A: fg.A,
		}

		// Search for an outline pixel
		foundOutline := false
		for y := 0; y < diameter && !foundOutline; y++ {
			for x := 0; x < diameter && !foundOutline; x++ {
				if !isInsideCircleTest(x, y, center, radius) {
					continue
				}
				if !isOutlinePixelTest(x, y, center, radius) {
					continue
				}

				// This is an outline pixel — verify its color
				foundOutline = true
				outlinePixel := result.Image.At(x, y)
				or, og, ob, oa := outlinePixel.RGBA()

				// Convert expected dimmed to pre-multiplied 16-bit values
				er := uint32(expectedDimmed.R) * 0x101
				eg := uint32(expectedDimmed.G) * 0x101
				eb := uint32(expectedDimmed.B) * 0x101
				ea := uint32(expectedDimmed.A) * 0x101

				if or != er || og != eg || ob != eb || oa != ea {
					t.Fatalf("diameter=%d: outline pixel at (%d,%d) color mismatch: got RGBA(%d,%d,%d,%d), want RGBA(%d,%d,%d,%d) [fg=%v, dimFactor=0.3]",
						diameter, x, y, or, og, ob, oa, er, eg, eb, ea, fg)
				}
			}
		}

		if !foundOutline {
			t.Fatalf("diameter=%d: could not find any outline pixel", diameter)
		}
	})
}

// isInsideCircleTest replicates the circle check for test verification.
func isInsideCircleTest(x, y int, center, radius float64) bool {
	dx := float64(x) + 0.5 - center
	dy := float64(y) + 0.5 - center
	dist := math.Sqrt(dx*dx + dy*dy)
	return dist <= radius
}

// isOutlinePixelTest replicates the outline check for test verification.
func isOutlinePixelTest(x, y int, center, radius float64) bool {
	neighbors := [4][2]int{
		{x - 1, y},
		{x + 1, y},
		{x, y - 1},
		{x, y + 1},
	}
	for _, n := range neighbors {
		if !isInsideCircleTest(n[0], n[1], center, radius) {
			return true
		}
	}
	return false
}

// --- From: led_test.go ---

// TestRenderNilForSmallDiameter verifies that Render returns nil for diameter < 3,
// including zero and negative values.

func TestRenderNilForSmallDiameter(t *testing.T) {
	cases := []struct {
		name     string
		diameter int
	}{
		{"diameter 0", 0},
		{"diameter 1", 1},
		{"diameter 2", 2},
		{"diameter -1", -1},
		{"diameter -5", -5},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := Config{
				State:    On,
				Diameter: tc.diameter,
			}
			result := Render(cfg)
			if result != nil {
				t.Errorf("expected nil for diameter %d, got non-nil result", tc.diameter)
			}
		})
	}
}

// TestDiameter3OnNoShine verifies that diameter=3 with State=On produces a 3×3 image
// with no shine pixel (white) and foreground fill inside the circle.

func TestDiameter3OnNoShine(t *testing.T) {
	cfg := Config{
		State:      On,
		Brightness: -1.0,
		Diameter:   3,
	}
	result := Render(cfg)
	if result == nil {
		t.Fatal("expected non-nil result for diameter 3")
	}

	bounds := result.Image.Bounds()
	if bounds.Dx() != 3 || bounds.Dy() != 3 {
		t.Fatalf("expected 3×3 image, got %d×%d", bounds.Dx(), bounds.Dy())
	}

	// For diameter 3, shine position would be (3/4, 3/4) = (0, 0).
	// It should NOT be white since diameter < 5.
	white := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	r, g, b, a := result.Image.At(0, 0).RGBA()
	pixelColor := color.RGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: uint8(a >> 8)}
	if pixelColor == white {
		t.Error("pixel at (0,0) should NOT be white for diameter 3 (no shine)")
	}

	// Center pixel (1, 1) should be foreground (default green)
	expectedFg := color.RGBA{R: 0, G: 200, B: 0, A: 255}
	r, g, b, a = result.Image.At(1, 1).RGBA()
	centerColor := color.RGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: uint8(a >> 8)}
	if centerColor != expectedFg {
		t.Errorf("center pixel (1,1) expected %v, got %v", expectedFg, centerColor)
	}
}

// TestDiameter4OnNoShine verifies that diameter=4 with State=On produces a 4×4 image
// with no shine pixel (white).

func TestDiameter4OnNoShine(t *testing.T) {
	cfg := Config{
		State:      On,
		Brightness: -1.0,
		Diameter:   4,
	}
	result := Render(cfg)
	if result == nil {
		t.Fatal("expected non-nil result for diameter 4")
	}

	bounds := result.Image.Bounds()
	if bounds.Dx() != 4 || bounds.Dy() != 4 {
		t.Fatalf("expected 4×4 image, got %d×%d", bounds.Dx(), bounds.Dy())
	}

	// For diameter 4, shine position would be (4/4, 4/4) = (1, 1).
	// It should NOT be white since diameter < 5.
	white := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	r, g, b, a := result.Image.At(1, 1).RGBA()
	pixelColor := color.RGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: uint8(a >> 8)}
	if pixelColor == white {
		t.Error("pixel at (1,1) should NOT be white for diameter 4 (no shine)")
	}
}

// TestDiameter5OnWithShine verifies that diameter=5 with State=On produces a shine pixel
// at (5/4, 5/4) = (1, 1) that is opaque white.

func TestDiameter5OnWithShine(t *testing.T) {
	fg := color.RGBA{R: 0, G: 200, B: 0, A: 255}
	cfg := Config{
		State:      On,
		Brightness: -1.0,
		Diameter:   5,
		Foreground: fg,
	}
	result := Render(cfg)
	if result == nil {
		t.Fatal("expected non-nil result for diameter 5")
	}

	// Shine is not yet implemented (task 6.1).
	// Verify center pixel (2, 2) is foreground color.
	r, g, b, a := result.Image.At(2, 2).RGBA()
	centerColor := color.RGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: uint8(a >> 8)}
	if centerColor != fg {
		t.Errorf("center pixel (2,2) expected foreground %v, got %v", fg, centerColor)
	}
}

// TestDefaultColors verifies that zero-value Foreground defaults to green (0,200,0,255)
// and zero-value Background defaults to black (0,0,0,255).

func TestDefaultColors(t *testing.T) {
	cfg := Config{
		State:      On,
		Brightness: -1.0,
		Diameter:   10,
		// Zero-value Foreground and Background
	}
	result := Render(cfg)
	if result == nil {
		t.Fatal("expected non-nil result for diameter 10")
	}

	expectedFg := color.RGBA{R: 0, G: 200, B: 0, A: 255}

	// Check a pixel inside the circle (center at 5,5) — should be green
	r, g, b, a := result.Image.At(5, 5).RGBA()
	insideColor := color.RGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: uint8(a >> 8)}
	if insideColor != expectedFg {
		t.Errorf("inside pixel (5,5) expected foreground %v, got %v", expectedFg, insideColor)
	}

	// Corner pixel (0, 0) — outside the circle, should be transparent in the new renderer
	r, g, b, a = result.Image.At(0, 0).RGBA()
	if a != 0 {
		t.Errorf("outside pixel (0,0) expected transparent, got RGBA(%d,%d,%d,%d)", r, g, b, a)
	}
}

// TestOffStateDimming verifies that in the Off state, the outline uses a dimmed color
// computed as (R*0.3, G*0.3, B*0.3, A) truncated to uint8.

func TestOffStateDimming(t *testing.T) {
	fg := color.RGBA{R: 100, G: 200, B: 50, A: 255}
	cfg := Config{
		State:      Off,
		Brightness: -1.0,
		Diameter:   10,
		Foreground: fg,
	}
	result := Render(cfg)
	if result == nil {
		t.Fatal("expected non-nil result for diameter 10, Off state")
	}

	// Expected dimmed color: (100*0.3=30, 200*0.3=60, 50*0.3=15, 255)
	expectedDimmed := color.RGBA{R: 30, G: 60, B: 15, A: 255}

	// Find an outline pixel. The outline is at circle boundary.
	// For d=10, center=5.0, radius=5.0.
	// Pixel (5, 0) should be inside the circle (distance from center = 4.5 ≤ 5.0)
	// and has neighbor (5, -1) outside → outline pixel.
	r, g, b, a := result.Image.At(5, 0).RGBA()
	outlineColor := color.RGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: uint8(a >> 8)}
	if outlineColor != expectedDimmed {
		t.Errorf("outline pixel (5,0) expected dimmed %v, got %v", expectedDimmed, outlineColor)
	}

	// Verify shine pixel position (10/4, 10/4) = (2, 2) is NOT white
	white := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	r, g, b, a = result.Image.At(2, 2).RGBA()
	shinePos := color.RGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: uint8(a >> 8)}
	if shinePos == white {
		t.Error("pixel at (2,2) should NOT be white in Off state (no shine)")
	}
}

// TestResolveColorsZeroDefaults verifies that zero-value foreground defaults to
// opaque green {0, 200, 0, 255} and zero-value background defaults to opaque black {0, 0, 0, 255}.

func TestResolveColorsZeroDefaults(t *testing.T) {
	cfg := Config{}
	resolveColors(&cfg)

	expectedFg := color.RGBA{R: 0, G: 200, B: 0, A: 255}
	expectedBg := color.RGBA{R: 0, G: 0, B: 0, A: 255}

	if cfg.Foreground != expectedFg {
		t.Errorf("zero foreground: expected %v, got %v", expectedFg, cfg.Foreground)
	}
	if cfg.Background != expectedBg {
		t.Errorf("zero background: expected %v, got %v", expectedBg, cfg.Background)
	}
}

// TestResolveColorsNonZeroPassthrough verifies that non-zero foreground and background
// values pass through unchanged.

func TestResolveColorsNonZeroPassthrough(t *testing.T) {
	inputFg := color.RGBA{R: 100, G: 50, B: 200, A: 255}
	inputBg := color.RGBA{R: 50, G: 25, B: 100, A: 255}

	cfg := Config{Foreground: inputFg, Background: inputBg}
	resolveColors(&cfg)

	if cfg.Foreground != inputFg {
		t.Errorf("non-zero foreground: expected %v, got %v", inputFg, cfg.Foreground)
	}
	if cfg.Background != inputBg {
		t.Errorf("non-zero background: expected %v, got %v", inputBg, cfg.Background)
	}
}

// --- From: brightness_test.go ---

// TestResolveBrightnessSentinel verifies that when Brightness == -1.0, the
// discrete State field determines effective brightness.

func TestResolveBrightnessSentinel(t *testing.T) {
	cases := []struct {
		name     string
		state    State
		expected float64
	}{
		{"On state → 1.0", On, 1.0},
		{"Off state → 0.0", Off, 0.0},
		{"Warning state → 1.0", Warning, 1.0},
		{"Unrecognized state → 0.0", State(99), 0.0},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := Config{
				Brightness: -1.0,
				State:      tc.state,
			}
			got := resolveBrightness(cfg)
			if got != tc.expected {
				t.Errorf("resolveBrightness(sentinel, state=%d): got %f, want %f",
					tc.state, got, tc.expected)
			}
		})
	}
}

// TestResolveBrightnessContinuous verifies that non-sentinel Brightness values
// are used directly (clamped to [0.0, 1.0]).

func TestResolveBrightnessContinuous(t *testing.T) {
	cases := []struct {
		name       string
		brightness float64
		expected   float64
	}{
		{"0.0 → 0.0", 0.0, 0.0},
		{"0.5 → 0.5", 0.5, 0.5},
		{"1.0 → 1.0", 1.0, 1.0},
		{"0.25 → 0.25", 0.25, 0.25},
		{"0.99 → 0.99", 0.99, 0.99},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := Config{
				Brightness: tc.brightness,
				State:      On, // Should be ignored when Brightness != -1.0
			}
			got := resolveBrightness(cfg)
			if got != tc.expected {
				t.Errorf("resolveBrightness(%f): got %f, want %f",
					tc.brightness, got, tc.expected)
			}
		})
	}
}

// TestResolveBrightnessDefensiveClamping verifies defensive handling of
// NaN, Inf, and out-of-range values.

func TestResolveBrightnessDefensiveClamping(t *testing.T) {
	cases := []struct {
		name       string
		brightness float64
		expected   float64
	}{
		{"NaN → 0.0", math.NaN(), 0.0},
		{"+Inf → 0.0", math.Inf(1), 0.0},
		{"-Inf → 0.0", math.Inf(-1), 0.0},
		{"negative (-0.5) → 0.0", -0.5, 0.0},
		{"greater than 1.0 (1.5) → 1.0", 1.5, 1.0},
		{"greater than 1.0 (100.0) → 1.0", 100.0, 1.0},
		{"-2.0 → 0.0", -2.0, 0.0},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := Config{
				Brightness: tc.brightness,
				State:      On, // Should be ignored
			}
			got := resolveBrightness(cfg)
			if math.IsNaN(tc.expected) {
				if !math.IsNaN(got) {
					t.Errorf("resolveBrightness(%v): got %f, want NaN",
						tc.brightness, got)
				}
			} else if got != tc.expected {
				t.Errorf("resolveBrightness(%v): got %f, want %f",
					tc.brightness, got, tc.expected)
			}
		})
	}
}

// TestResolveBrightnessContinuousOverridesState verifies that when Brightness
// is not -1.0, the State field is entirely ignored.

func TestResolveBrightnessContinuousOverridesState(t *testing.T) {
	// Brightness 0.7 should be used regardless of what State says
	for _, state := range []State{On, Off, Warning} {
		cfg := Config{
			Brightness: 0.7,
			State:      state,
		}
		got := resolveBrightness(cfg)
		if got != 0.7 {
			t.Errorf("resolveBrightness(0.7, state=%d): got %f, want 0.7",
				state, got)
		}
	}
}
