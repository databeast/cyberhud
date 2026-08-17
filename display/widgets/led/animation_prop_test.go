package led

import (
	"image"
	"image/color"
	"math"
	"testing"
	"time"

	"pgregory.net/rapid"
)

// TestProperty11_AnimationBrightnessBounds verifies that for any Config with
// animation configured and any elapsed time, the effective brightness stays within
// the defined bounds:
//   - Pulse: result in [MinBrightness, 1.0]
//   - Blink: result is exactly 0.0 or 1.0
//   - Fade: result in [0.0, 1.0]

func TestProperty11_AnimationBrightnessBounds(t *testing.T) {
	t.Run("Pulse", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			// Random period in valid range [100ms, 5000ms]
			periodMs := rapid.IntRange(100, 5000).Draw(t, "periodMs")
			period := time.Duration(periodMs) * time.Millisecond

			// Random MinBrightness in [0.0, 0.99]
			minBrightMillis := rapid.IntRange(0, 990).Draw(t, "minBrightMillis")
			minBrightness := float64(minBrightMillis) / 1000.0

			// Random elapsed time (up to several periods)
			elapsedMs := rapid.IntRange(0, 20000).Draw(t, "elapsedMs")
			elapsed := time.Duration(elapsedMs) * time.Millisecond

			// Base brightness (Pulse modulates from scratch, so use 1.0 as base)
			baseBrightness := 1.0

			cfg := Config{
				Animation: AnimationConfig{
					Type:          Pulse,
					Period:        period,
					MinBrightness: minBrightness,
				},
				animElapsed: elapsed,
			}

			result := resolveAnimation(cfg, baseBrightness)

			if result < minBrightness-1e-9 {
				t.Fatalf("Pulse brightness %f below MinBrightness %f [period=%v, elapsed=%v]",
					result, minBrightness, period, elapsed)
			}
			if result > 1.0+1e-9 {
				t.Fatalf("Pulse brightness %f exceeds 1.0 [period=%v, elapsed=%v, minBrightness=%f]",
					result, period, elapsed, minBrightness)
			}
		})
	})

	t.Run("Blink", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			// Random period in valid range [100ms, 5000ms]
			periodMs := rapid.IntRange(100, 5000).Draw(t, "periodMs")
			period := time.Duration(periodMs) * time.Millisecond

			// Random elapsed time
			elapsedMs := rapid.IntRange(0, 20000).Draw(t, "elapsedMs")
			elapsed := time.Duration(elapsedMs) * time.Millisecond

			baseBrightness := 1.0

			cfg := Config{
				Animation: AnimationConfig{
					Type:   Blink,
					Period: period,
				},
				animElapsed: elapsed,
			}

			result := resolveAnimation(cfg, baseBrightness)

			if result != 0.0 && result != 1.0 {
				t.Fatalf("Blink brightness %f is not exactly 0.0 or 1.0 [period=%v, elapsed=%v]",
					result, period, elapsed)
			}
		})
	})

	t.Run("Fade", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			// Random period in valid range [100ms, 5000ms]
			periodMs := rapid.IntRange(100, 5000).Draw(t, "periodMs")
			period := time.Duration(periodMs) * time.Millisecond

			// Random elapsed time
			elapsedMs := rapid.IntRange(0, 20000).Draw(t, "elapsedMs")
			elapsed := time.Duration(elapsedMs) * time.Millisecond

			baseBrightness := 1.0

			cfg := Config{
				Animation: AnimationConfig{
					Type:   Fade,
					Period: period,
				},
				animElapsed: elapsed,
			}

			result := resolveAnimation(cfg, baseBrightness)

			if result < -1e-9 {
				t.Fatalf("Fade brightness %f below 0.0 [period=%v, elapsed=%v]",
					result, period, elapsed)
			}
			if result > 1.0+1e-9 {
				t.Fatalf("Fade brightness %f exceeds 1.0 [period=%v, elapsed=%v]",
					result, period, elapsed)
			}
		})
	})
}

// TestProperty12_NoAnimationEquivalence verifies that for any Config where
// Animation.Type is NoAnimation, or the Period is ≤ 0, or an unrecognized
// Animation type is used, the rendered output SHALL be equivalent to the same
// Config rendered without any animation effects applied.
//
// We test this by calling resolveAnimation directly and verifying baseBrightness
// is returned unchanged.

func TestProperty12_NoAnimationEquivalence(t *testing.T) {
	t.Run("NoAnimation_type", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			// Random base brightness
			baseBright := rapid.Float64Range(0.0, 1.0).Draw(t, "baseBrightness")
			// Random elapsed time
			elapsedMs := rapid.IntRange(0, 10000).Draw(t, "elapsedMs")
			elapsed := time.Duration(elapsedMs) * time.Millisecond
			// Random period (shouldn't matter since type is NoAnimation)
			periodMs := rapid.IntRange(100, 5000).Draw(t, "periodMs")
			period := time.Duration(periodMs) * time.Millisecond

			cfg := Config{
				Animation: AnimationConfig{
					Type:   NoAnimation,
					Period: period,
				},
				animElapsed: elapsed,
			}

			result := resolveAnimation(cfg, baseBright)

			if result != baseBright {
				t.Fatalf("NoAnimation type changed brightness from %f to %f [period=%v, elapsed=%v]",
					baseBright, result, period, elapsed)
			}
		})
	})

	t.Run("Zero_or_negative_period", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			// Animation type that would normally animate
			animType := Animation(rapid.IntRange(1, 3).Draw(t, "animType")) // Pulse, Blink, or Fade

			baseBright := rapid.Float64Range(0.0, 1.0).Draw(t, "baseBrightness")
			elapsedMs := rapid.IntRange(0, 10000).Draw(t, "elapsedMs")
			elapsed := time.Duration(elapsedMs) * time.Millisecond

			// Zero or negative period
			periodMs := rapid.IntRange(-5000, 0).Draw(t, "periodMs")
			period := time.Duration(periodMs) * time.Millisecond

			cfg := Config{
				Animation: AnimationConfig{
					Type:          animType,
					Period:        period,
					MinBrightness: 0.3,
				},
				animElapsed: elapsed,
			}

			result := resolveAnimation(cfg, baseBright)

			if result != baseBright {
				t.Fatalf("Zero/negative period changed brightness from %f to %f [animType=%d, period=%v, elapsed=%v]",
					baseBright, result, animType, period, elapsed)
			}
		})
	})

	t.Run("Unrecognized_animation_type", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			// Unrecognized animation type (values outside [0, 3])
			animType := Animation(rapid.IntRange(4, 100).Draw(t, "animType"))

			baseBright := rapid.Float64Range(0.0, 1.0).Draw(t, "baseBrightness")
			elapsedMs := rapid.IntRange(0, 10000).Draw(t, "elapsedMs")
			elapsed := time.Duration(elapsedMs) * time.Millisecond
			periodMs := rapid.IntRange(100, 5000).Draw(t, "periodMs")
			period := time.Duration(periodMs) * time.Millisecond

			cfg := Config{
				Animation: AnimationConfig{
					Type:          animType,
					Period:        period,
					MinBrightness: 0.3,
				},
				animElapsed: elapsed,
			}

			result := resolveAnimation(cfg, baseBright)

			if result != baseBright {
				t.Fatalf("Unrecognized animation type %d changed brightness from %f to %f [period=%v, elapsed=%v]",
					animType, baseBright, result, period, elapsed)
			}
		})
	})
}

// TestProperty13_AnimationModulatesGlowInSyncWithBrightness verifies that for any
// Config with Pulse animation and glow enabled, the glow opacity at the body edge
// SHALL equal GlowColor.A × effectiveBrightness at any elapsed time.
//
// We render a full LED with Pulse animation + glow, then check a glow pixel near
// the body edge to verify its alpha ≈ floor(GlowColor.A × effectiveBrightness).

func TestProperty13_AnimationModulatesGlowInSyncWithBrightness(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		diameter := rapid.IntRange(12, 40).Draw(t, "diameter")
		glowRadius := rapid.IntRange(4, 12).Draw(t, "glowRadius")

		// Pulse animation parameters
		periodMs := rapid.IntRange(200, 3000).Draw(t, "periodMs")
		period := time.Duration(periodMs) * time.Millisecond
		minBrightMillis := rapid.IntRange(100, 800).Draw(t, "minBrightMillis")
		minBrightness := float64(minBrightMillis) / 1000.0

		// Random elapsed time within a few periods
		elapsedMs := rapid.IntRange(0, periodMs*3).Draw(t, "elapsedMs")
		elapsed := time.Duration(elapsedMs) * time.Millisecond

		// Use foreground with A=255 for predictable glow calculations
		fg := color.RGBA{R: 0, G: 200, B: 0, A: 255}

		cfg := Config{
			Shape:       Circle,
			State:       On,
			Brightness:  -1.0,
			Diameter:    diameter,
			Bounds:      image.Rect(0, 0, diameter, diameter),
			Foreground:  fg,
			GlowEnabled: true,
			GlowRadius:  glowRadius,
			BorderWidth: 0,
			Animation: AnimationConfig{
				Type:          Pulse,
				Period:        period,
				MinBrightness: minBrightness,
			},
			animElapsed: elapsed,
		}

		// Compute expected effective brightness by calling resolveAnimation directly
		baseBrightness := resolveBrightness(cfg)
		effectiveBrightness := resolveAnimation(cfg, baseBrightness)

		// Skip edge case where brightness rounds to zero (no glow rendered)
		if effectiveBrightness < 0.01 {
			return
		}

		result := Render(cfg)
		if result == nil {
			t.Fatal("expected non-nil sprite")
		}

		img := result.Image.(*image.RGBA)
		outputSize := img.Bounds().Dx()
		center := float64(outputSize) / 2.0
		shapeRadius := float64(diameter) / 2.0
		glowRadiusF := float64(glowRadius)

		// Sample a glow pixel on the horizontal midline, just outside the body
		// Find the first pixel that's clearly in the glow region (1-2 pixels outside body edge)
		verified := false
		midY := outputSize / 2
		for px := 0; px < outputSize; px++ {
			pcx := float64(px) + 0.5
			pcy := float64(midY) + 0.5

			dx := pcx - center
			dy := pcy - center
			dist := math.Sqrt(dx*dx+dy*dy) - shapeRadius

			// Check pixels just outside the body (small distance into glow)
			if dist > 0.5 && dist < glowRadiusF*0.5 {
				c := img.RGBAAt(px, midY)
				if c.A == 0 {
					continue
				}

				// Expected alpha = floor(glowBase.A × (1 − dist / glowRadius) × effectiveBrightness)
				falloff := 1.0 - dist/glowRadiusF
				expectedAlpha := math.Floor(float64(fg.A) * falloff * effectiveBrightness)
				if expectedAlpha < 0 {
					expectedAlpha = 0
				}
				if expectedAlpha > 255 {
					expectedAlpha = 255
				}

				// Allow ±1 tolerance for floating point rounding
				diff := math.Abs(float64(c.A) - expectedAlpha)
				if diff > 1.0 {
					t.Fatalf("glow pixel at (%d,%d) has alpha=%d, expected=%d "+
						"(dist=%.2f, falloff=%.4f, effectiveBrightness=%.4f) "+
						"[diameter=%d, glowRadius=%d, period=%v, elapsed=%v, minBrightness=%.3f]",
						px, midY, c.A, int(expectedAlpha), dist, falloff, effectiveBrightness,
						diameter, glowRadius, period, elapsed, minBrightness)
				}

				verified = true
				break
			}
		}

		if !verified {
			t.Fatalf("no glow pixel verified on horizontal midline "+
				"[diameter=%d, glowRadius=%d, effectiveBrightness=%.4f]",
				diameter, glowRadius, effectiveBrightness)
		}
	})
}
