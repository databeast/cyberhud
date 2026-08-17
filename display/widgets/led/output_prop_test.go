package led

import (
	"image"
	"image/color"
	"testing"

	"pgregory.net/rapid"
)

// TestProperty4_OutputDimensionsMatchDiameterAndGlowConfiguration verifies that for
// any valid Config with Diameter ≥ 3, the returned Sprite SHALL have Image pixel
// dimensions equal to (Diameter + 2 × effectiveGlowRadius) in both width and height
// when glow is enabled, or Diameter × Diameter when glow is disabled. Position SHALL
// equal Bounds.Min.

func TestProperty4_OutputDimensionsMatchDiameterAndGlowConfiguration(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		diameter := rapid.IntRange(3, 64).Draw(t, "diameter")
		glowEnabled := rapid.Bool().Draw(t, "glowEnabled")
		glowRadius := rapid.IntRange(0, 32).Draw(t, "glowRadius")
		shape := Shape(rapid.IntRange(0, 3).Draw(t, "shape"))
		state := State(rapid.IntRange(0, 2).Draw(t, "state"))

		// Random position for Bounds.Min
		posX := rapid.IntRange(-100, 500).Draw(t, "posX")
		posY := rapid.IntRange(-100, 500).Draw(t, "posY")

		fg := color.RGBA{
			R: uint8(rapid.IntRange(1, 255).Draw(t, "fgR")),
			G: uint8(rapid.IntRange(1, 255).Draw(t, "fgG")),
			B: uint8(rapid.IntRange(1, 255).Draw(t, "fgB")),
			A: 255,
		}

		cfg := Config{
			Shape:       shape,
			State:       state,
			Brightness:  -1.0,
			Diameter:    diameter,
			Bounds:      image.Rect(posX, posY, posX+diameter, posY+diameter),
			Foreground:  fg,
			GlowEnabled: glowEnabled,
			GlowRadius:  glowRadius,
		}

		result := Render(cfg)
		if result == nil {
			t.Fatal("expected non-nil sprite for valid config")
		}

		// Compute expected effective glow radius
		expectedGlowRadius := 0
		if glowEnabled {
			expectedGlowRadius = computeExpectedGlowRadius(diameter, glowRadius)
		}

		expectedSize := diameter + 2*expectedGlowRadius

		// Check image dimensions
		imgBounds := result.Image.Bounds()
		actualWidth := imgBounds.Dx()
		actualHeight := imgBounds.Dy()

		if actualWidth != expectedSize {
			t.Fatalf("width mismatch: got %d, want %d (diameter=%d, glowEnabled=%v, glowRadius=%d, effectiveGlowRadius=%d)",
				actualWidth, expectedSize, diameter, glowEnabled, glowRadius, expectedGlowRadius)
		}
		if actualHeight != expectedSize {
			t.Fatalf("height mismatch: got %d, want %d (diameter=%d, glowEnabled=%v, glowRadius=%d, effectiveGlowRadius=%d)",
				actualHeight, expectedSize, diameter, glowEnabled, glowRadius, expectedGlowRadius)
		}

		// Check position equals Bounds.Min
		if result.Position != cfg.Bounds.Min {
			t.Fatalf("position mismatch: got %v, want %v", result.Position, cfg.Bounds.Min)
		}
	})
}

// computeExpectedGlowRadius mirrors the effectiveGlowRadius logic from led.go:
// If GlowRadius == 0, default to 30% of body radius (Diameter/2), clamped to [1, 32].
// Otherwise use the configured value (already clamped to [0, 32] by validate).
func computeExpectedGlowRadius(diameter, glowRadius int) int {
	if glowRadius == 0 {
		bodyRadius := diameter / 2
		r := int(float64(bodyRadius) * 0.3)
		if r < 1 {
			r = 1
		}
		if r > 32 {
			r = 32
		}
		return r
	}
	// Clamp to [0, 32] as validate does
	if glowRadius > 32 {
		return 32
	}
	return glowRadius
}

// TestProperty27_LabelCorrectnessBasedOnBrightnessAndState verifies that for any valid
// Config, the Sprite label SHALL be:
//   - "led/on" when effective brightness > 0.5, or State is On with Brightness = -1.0
//   - "led/off" when effective brightness = 0.0, or State is Off with Brightness = -1.0
//   - "led/warning" when State is Warning (with Brightness = -1.0), or effective brightness ∈ (0.0, 0.5]

func TestProperty27_LabelCorrectnessBasedOnBrightnessAndState(t *testing.T) {
	t.Run("discrete_mode", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			diameter := rapid.IntRange(3, 64).Draw(t, "diameter")
			state := State(rapid.IntRange(0, 2).Draw(t, "state"))

			cfg := Config{
				Shape:      Circle,
				State:      state,
				Brightness: -1.0, // Discrete mode
				Diameter:   diameter,
				Bounds:     image.Rect(0, 0, diameter, diameter),
				Foreground: color.RGBA{R: 0, G: 200, B: 0, A: 255},
			}

			result := Render(cfg)
			if result == nil {
				t.Fatal("expected non-nil sprite")
			}

			var expectedLabel string
			switch state {
			case On:
				expectedLabel = "led/on"
			case Off:
				expectedLabel = "led/off"
			case Warning:
				expectedLabel = "led/warning"
			}

			if result.Label != expectedLabel {
				t.Fatalf("label mismatch in discrete mode: state=%d, got %q, want %q",
					state, result.Label, expectedLabel)
			}
		})
	})

	t.Run("continuous_brightness_on", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			diameter := rapid.IntRange(3, 64).Draw(t, "diameter")
			// Brightness > 0.5 → "led/on"
			brightness := rapid.Float64Range(0.501, 1.0).Draw(t, "brightness")

			cfg := Config{
				Shape:      Circle,
				State:      Off, // State should not matter in continuous mode
				Brightness: brightness,
				Diameter:   diameter,
				Bounds:     image.Rect(0, 0, diameter, diameter),
				Foreground: color.RGBA{R: 0, G: 200, B: 0, A: 255},
			}

			result := Render(cfg)
			if result == nil {
				t.Fatal("expected non-nil sprite")
			}

			if result.Label != "led/on" {
				t.Fatalf("label mismatch: brightness=%f (>0.5), got %q, want \"led/on\"",
					brightness, result.Label)
			}
		})
	})

	t.Run("continuous_brightness_off", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			diameter := rapid.IntRange(3, 64).Draw(t, "diameter")

			cfg := Config{
				Shape:      Circle,
				State:      On, // State should not matter in continuous mode
				Brightness: 0.0,
				Diameter:   diameter,
				Bounds:     image.Rect(0, 0, diameter, diameter),
				Foreground: color.RGBA{R: 0, G: 200, B: 0, A: 255},
			}

			result := Render(cfg)
			if result == nil {
				t.Fatal("expected non-nil sprite")
			}

			if result.Label != "led/off" {
				t.Fatalf("label mismatch: brightness=0.0, got %q, want \"led/off\"",
					result.Label)
			}
		})
	})

	t.Run("continuous_brightness_warning", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			diameter := rapid.IntRange(3, 64).Draw(t, "diameter")
			// Brightness ∈ (0.0, 0.5] → "led/warning"
			brightness := rapid.Float64Range(0.001, 0.5).Draw(t, "brightness")

			cfg := Config{
				Shape:      Circle,
				State:      On, // State should not matter in continuous mode
				Brightness: brightness,
				Diameter:   diameter,
				Bounds:     image.Rect(0, 0, diameter, diameter),
				Foreground: color.RGBA{R: 0, G: 200, B: 0, A: 255},
			}

			result := Render(cfg)
			if result == nil {
				t.Fatal("expected non-nil sprite")
			}

			if result.Label != "led/warning" {
				t.Fatalf("label mismatch: brightness=%f (in (0.0, 0.5]), got %q, want \"led/warning\"",
					brightness, result.Label)
			}
		})
	})
}
