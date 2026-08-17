package led

import (
	"image"
	"image/color"
	"math"
	"testing"

	"pgregory.net/rapid"
)

// TestProperty3_DiameterBelowMinimumReturnsNil verifies that for any Config with
// Diameter < 3 (including negative and zero values), Render SHALL return nil
// regardless of Shape, State, Group, or any other field.

func TestProperty3_DiameterBelowMinimumReturnsNil(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate diameter in [-100, 2] (always below minimum of 3)
		diameter := rapid.IntRange(-100, 2).Draw(t, "diameter")

		// Generate random shape (including invalid values)
		shape := Shape(rapid.IntRange(-5, 10).Draw(t, "shape"))

		// Generate random state
		state := State(rapid.IntRange(0, 2).Draw(t, "state"))

		// Generate random brightness
		brightness := rapid.Float64Range(-2.0, 2.0).Draw(t, "brightness")

		// Generate random foreground color
		fg := color.RGBA{
			R: uint8(rapid.IntRange(0, 255).Draw(t, "fgR")),
			G: uint8(rapid.IntRange(0, 255).Draw(t, "fgG")),
			B: uint8(rapid.IntRange(0, 255).Draw(t, "fgB")),
			A: uint8(rapid.IntRange(0, 255).Draw(t, "fgA")),
		}

		// Optionally include a group to verify groups also return nil
		var group []GroupEntry
		hasGroup := rapid.Bool().Draw(t, "hasGroup")
		if hasGroup {
			numEntries := rapid.IntRange(1, 5).Draw(t, "numEntries")
			group = make([]GroupEntry, numEntries)
			for i := range group {
				group[i] = GroupEntry{
					State: State(rapid.IntRange(0, 2).Draw(t, "entryState")),
				}
			}
		}

		cfg := Config{
			Shape:      shape,
			State:      state,
			Brightness: brightness,
			Diameter:   diameter,
			Bounds:     image.Rect(0, 0, 100, 100),
			Foreground: fg,
			Group:      group,
		}

		result := Render(cfg)
		if result != nil {
			t.Fatalf("expected nil for Diameter=%d, got non-nil sprite (shape=%d, state=%d, brightness=%f, group=%v)",
				diameter, shape, state, brightness, hasGroup)
		}
	})
}

// TestProperty2_UnrecognizedShapeFallsBackToCircle verifies that for any Config
// with an unrecognized Shape enum value (outside [0, 3]), the rendered output SHALL
// be pixel-identical to the same Config with Shape set to Circle.

func TestProperty2_UnrecognizedShapeFallsBackToCircle(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate an invalid shape value (outside [0, 3])
		invalidShape := Shape(rapid.OneOf(
			rapid.IntRange(-100, -1),
			rapid.IntRange(4, 100),
		).Draw(t, "invalidShape"))

		diameter := rapid.IntRange(3, 64).Draw(t, "diameter")
		state := State(rapid.IntRange(0, 2).Draw(t, "state"))

		fg := color.RGBA{
			R: uint8(rapid.IntRange(1, 255).Draw(t, "fgR")),
			G: uint8(rapid.IntRange(1, 255).Draw(t, "fgG")),
			B: uint8(rapid.IntRange(1, 255).Draw(t, "fgB")),
			A: uint8(rapid.IntRange(1, 255).Draw(t, "fgA")),
		}

		bg := color.RGBA{
			R: uint8(rapid.IntRange(1, 255).Draw(t, "bgR")),
			G: uint8(rapid.IntRange(1, 255).Draw(t, "bgG")),
			B: uint8(rapid.IntRange(1, 255).Draw(t, "bgB")),
			A: uint8(rapid.IntRange(1, 255).Draw(t, "bgA")),
		}

		bounds := image.Rect(0, 0, diameter, diameter)

		// Render with invalid shape
		cfgInvalid := Config{
			Shape:      invalidShape,
			State:      state,
			Brightness: -1.0,
			Diameter:   diameter,
			Bounds:     bounds,
			Foreground: fg,
			Background: bg,
		}
		resultInvalid := Render(cfgInvalid)

		// Render with explicit Circle shape
		cfgCircle := Config{
			Shape:      Circle,
			State:      state,
			Brightness: -1.0,
			Diameter:   diameter,
			Bounds:     bounds,
			Foreground: fg,
			Background: bg,
		}
		resultCircle := Render(cfgCircle)

		if resultInvalid == nil {
			t.Fatal("expected non-nil result for valid diameter with invalid shape")
		}
		if resultCircle == nil {
			t.Fatal("expected non-nil result for valid diameter with Circle shape")
		}

		// Compare pixel-for-pixel
		imgInvalid := resultInvalid.Image
		imgCircle := resultCircle.Image

		if imgInvalid.Bounds() != imgCircle.Bounds() {
			t.Fatalf("image bounds differ: invalid=%v, circle=%v",
				imgInvalid.Bounds(), imgCircle.Bounds())
		}

		b := imgInvalid.Bounds()
		for y := b.Min.Y; y < b.Max.Y; y++ {
			for x := b.Min.X; x < b.Max.X; x++ {
				r1, g1, b1, a1 := imgInvalid.At(x, y).RGBA()
				r2, g2, b2, a2 := imgCircle.At(x, y).RGBA()
				if r1 != r2 || g1 != g2 || b1 != b2 || a1 != a2 {
					t.Fatalf("pixel mismatch at (%d,%d): invalidShape(%d)=RGBA(%d,%d,%d,%d), circle=RGBA(%d,%d,%d,%d) [diameter=%d, state=%d]",
						x, y, invalidShape, r1, g1, b1, a1, r2, g2, b2, a2, diameter, state)
				}
			}
		}

		// Also verify labels match
		if resultInvalid.Label != resultCircle.Label {
			t.Fatalf("label mismatch: invalidShape=%q, circle=%q", resultInvalid.Label, resultCircle.Label)
		}
	})
}

// TestProperty20_BrightnessClamping verifies that for any Config with Brightness
// outside the valid range (< -1.0, between -1.0 and 0.0 exclusive, > 1.0, NaN, or
// ±Infinity), the rendered output SHALL be pixel-identical to the same Config with
// Brightness clamped to the nearest valid value in [0.0, 1.0] (NaN/Inf → 0.0).

func TestProperty20_BrightnessClamping(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		diameter := rapid.IntRange(3, 64).Draw(t, "diameter")

		fg := color.RGBA{
			R: uint8(rapid.IntRange(1, 255).Draw(t, "fgR")),
			G: uint8(rapid.IntRange(1, 255).Draw(t, "fgG")),
			B: uint8(rapid.IntRange(1, 255).Draw(t, "fgB")),
			A: uint8(rapid.IntRange(1, 255).Draw(t, "fgA")),
		}

		bg := color.RGBA{
			R: uint8(rapid.IntRange(1, 255).Draw(t, "bgR")),
			G: uint8(rapid.IntRange(1, 255).Draw(t, "bgG")),
			B: uint8(rapid.IntRange(1, 255).Draw(t, "bgB")),
			A: uint8(rapid.IntRange(1, 255).Draw(t, "bgA")),
		}

		bounds := image.Rect(0, 0, diameter, diameter)

		// Generate an invalid brightness value from one of several categories
		category := rapid.IntRange(0, 4).Draw(t, "category")
		var invalidBrightness float64
		var expectedClamped float64

		switch category {
		case 0:
			// Less than -1.0 → clamp to 0.0
			invalidBrightness = rapid.Float64Range(-100.0, -1.001).Draw(t, "brightness")
			expectedClamped = 0.0
		case 1:
			// Between -1.0 and 0.0 exclusive → clamp to 0.0
			invalidBrightness = rapid.Float64Range(-0.999, -0.001).Draw(t, "brightness")
			expectedClamped = 0.0
		case 2:
			// Greater than 1.0 → clamp to 1.0
			invalidBrightness = rapid.Float64Range(1.001, 100.0).Draw(t, "brightness")
			expectedClamped = 1.0
		case 3:
			// NaN → treat as 0.0
			invalidBrightness = math.NaN()
			expectedClamped = 0.0
		case 4:
			// ±Infinity → treat as 0.0
			if rapid.Bool().Draw(t, "posInf") {
				invalidBrightness = math.Inf(1)
			} else {
				invalidBrightness = math.Inf(-1)
			}
			expectedClamped = 0.0
		}

		// Render with invalid brightness
		cfgInvalid := Config{
			Shape:      Circle,
			State:      On, // State doesn't matter when Brightness != -1.0
			Brightness: invalidBrightness,
			Diameter:   diameter,
			Bounds:     bounds,
			Foreground: fg,
			Background: bg,
		}
		resultInvalid := Render(cfgInvalid)

		// Render with expected clamped brightness
		cfgClamped := Config{
			Shape:      Circle,
			State:      On,
			Brightness: expectedClamped,
			Diameter:   diameter,
			Bounds:     bounds,
			Foreground: fg,
			Background: bg,
		}
		resultClamped := Render(cfgClamped)

		if resultInvalid == nil {
			t.Fatal("expected non-nil result for valid diameter with invalid brightness")
		}
		if resultClamped == nil {
			t.Fatal("expected non-nil result for valid diameter with clamped brightness")
		}

		// Compare pixel-for-pixel
		imgInvalid := resultInvalid.Image
		imgClamped := resultClamped.Image

		if imgInvalid.Bounds() != imgClamped.Bounds() {
			t.Fatalf("image bounds differ: invalid=%v, clamped=%v",
				imgInvalid.Bounds(), imgClamped.Bounds())
		}

		rect := imgInvalid.Bounds()
		for y := rect.Min.Y; y < rect.Max.Y; y++ {
			for x := rect.Min.X; x < rect.Max.X; x++ {
				r1, g1, b1, a1 := imgInvalid.At(x, y).RGBA()
				r2, g2, b2, a2 := imgClamped.At(x, y).RGBA()
				if r1 != r2 || g1 != g2 || b1 != b2 || a1 != a2 {
					t.Fatalf("pixel mismatch at (%d,%d): invalidBrightness(%v)=RGBA(%d,%d,%d,%d), clamped(%v)=RGBA(%d,%d,%d,%d) [diameter=%d, category=%d]",
						x, y, invalidBrightness, r1, g1, b1, a1, expectedClamped, r2, g2, b2, a2, diameter, category)
				}
			}
		}

		// Labels should also match
		if resultInvalid.Label != resultClamped.Label {
			t.Fatalf("label mismatch: invalid=%q, clamped=%q (brightness=%v → %v)",
				resultInvalid.Label, resultClamped.Label, invalidBrightness, expectedClamped)
		}
	})
}
