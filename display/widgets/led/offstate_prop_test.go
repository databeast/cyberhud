package led

import (
	"image"
	"image/color"
	"math"
	"testing"

	"pgregory.net/rapid"
)

// TestProperty17_Brightness0RendersIdenticallyToOffState verifies that for any Config
// with Brightness = 0.0 (not -1.0), the rendered output SHALL be pixel-identical to the
// same Config with Brightness = -1.0 and State = Off.

func TestProperty17_Brightness0RendersIdenticallyToOffState(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		shape := Shape(rapid.IntRange(0, 3).Draw(t, "shape"))
		diameter := rapid.IntRange(3, 40).Draw(t, "diameter")

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

		// Optional border
		borderWidth := rapid.IntRange(0, 4).Draw(t, "borderWidth")
		borderColor := color.RGBA{
			R: uint8(rapid.IntRange(1, 255).Draw(t, "borderR")),
			G: uint8(rapid.IntRange(1, 255).Draw(t, "borderG")),
			B: uint8(rapid.IntRange(1, 255).Draw(t, "borderB")),
			A: 255,
		}

		// Render with Brightness=0.0 (continuous, state field ignored)
		cfgB0 := Config{
			Shape:       shape,
			State:       On, // should be ignored
			Brightness:  0.0,
			Diameter:    diameter,
			Bounds:      image.Rect(0, 0, diameter, diameter),
			Foreground:  fg,
			Background:  bg,
			BorderWidth: borderWidth,
			BorderColor: borderColor,
			// No glow, no gradient, no shine, no animation — isolate brightness behavior
			GlowEnabled: false,
		}

		// Render with State=Off, Brightness=-1.0 (discrete)
		cfgOff := Config{
			Shape:       shape,
			State:       Off,
			Brightness:  -1.0,
			Diameter:    diameter,
			Bounds:      image.Rect(0, 0, diameter, diameter),
			Foreground:  fg,
			Background:  bg,
			BorderWidth: borderWidth,
			BorderColor: borderColor,
			GlowEnabled: false,
		}

		resultB0 := Render(cfgB0)
		resultOff := Render(cfgOff)

		if resultB0 == nil || resultOff == nil {
			t.Fatalf("expected non-nil results for diameter=%d", diameter)
		}

		imgB0 := resultB0.Image.(*image.RGBA)
		imgOff := resultOff.Image.(*image.RGBA)

		boundsB0 := imgB0.Bounds()
		boundsOff := imgOff.Bounds()

		if boundsB0 != boundsOff {
			t.Fatalf("image bounds differ: B0=%v, Off=%v", boundsB0, boundsOff)
		}

		for y := boundsB0.Min.Y; y < boundsB0.Max.Y; y++ {
			for x := boundsB0.Min.X; x < boundsB0.Max.X; x++ {
				pxB0 := imgB0.RGBAAt(x, y)
				pxOff := imgOff.RGBAAt(x, y)
				if pxB0 != pxOff {
					t.Fatalf("pixel mismatch at (%d,%d): Brightness=0.0 → %v, State=Off → %v "+
						"[shape=%d, diameter=%d, fg=%v]",
						x, y, pxB0, pxOff, shape, diameter, fg)
				}
			}
		}
	})
}

// TestProperty18_Brightness1RendersIdenticallyToOnState verifies that for any Config
// with Brightness = 1.0 (not -1.0), the rendered output SHALL be pixel-identical to the
// same Config with Brightness = -1.0 and State = On at full foreground color.

func TestProperty18_Brightness1RendersIdenticallyToOnState(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		shape := Shape(rapid.IntRange(0, 3).Draw(t, "shape"))
		diameter := rapid.IntRange(3, 40).Draw(t, "diameter")

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

		// Optional border
		borderWidth := rapid.IntRange(0, 4).Draw(t, "borderWidth")
		borderColor := color.RGBA{
			R: uint8(rapid.IntRange(1, 255).Draw(t, "borderR")),
			G: uint8(rapid.IntRange(1, 255).Draw(t, "borderG")),
			B: uint8(rapid.IntRange(1, 255).Draw(t, "borderB")),
			A: 255,
		}

		// Render with Brightness=1.0 (continuous, state field ignored)
		cfgB1 := Config{
			Shape:       shape,
			State:       Off, // should be ignored
			Brightness:  1.0,
			Diameter:    diameter,
			Bounds:      image.Rect(0, 0, diameter, diameter),
			Foreground:  fg,
			Background:  bg,
			BorderWidth: borderWidth,
			BorderColor: borderColor,
			GlowEnabled: false,
		}

		// Render with State=On, Brightness=-1.0 (discrete)
		cfgOn := Config{
			Shape:       shape,
			State:       On,
			Brightness:  -1.0,
			Diameter:    diameter,
			Bounds:      image.Rect(0, 0, diameter, diameter),
			Foreground:  fg,
			Background:  bg,
			BorderWidth: borderWidth,
			BorderColor: borderColor,
			GlowEnabled: false,
		}

		resultB1 := Render(cfgB1)
		resultOn := Render(cfgOn)

		if resultB1 == nil || resultOn == nil {
			t.Fatalf("expected non-nil results for diameter=%d", diameter)
		}

		imgB1 := resultB1.Image.(*image.RGBA)
		imgOn := resultOn.Image.(*image.RGBA)

		boundsB1 := imgB1.Bounds()
		boundsOn := imgOn.Bounds()

		if boundsB1 != boundsOn {
			t.Fatalf("image bounds differ: B1=%v, On=%v", boundsB1, boundsOn)
		}

		for y := boundsB1.Min.Y; y < boundsB1.Max.Y; y++ {
			for x := boundsB1.Min.X; x < boundsB1.Max.X; x++ {
				pxB1 := imgB1.RGBAAt(x, y)
				pxOn := imgOn.RGBAAt(x, y)
				if pxB1 != pxOn {
					t.Fatalf("pixel mismatch at (%d,%d): Brightness=1.0 → %v, State=On → %v "+
						"[shape=%d, diameter=%d, fg=%v]",
						x, y, pxB1, pxOn, shape, diameter, fg)
				}
			}
		}
	})
}

// TestProperty19_ContinuousBrightnessScalesRGBPreservingAlpha verifies that for any
// Config with Brightness B in (0.0, 1.0) exclusive, every body-fill pixel SHALL have
// RGB channels equal to floor(foreground.R × B), floor(foreground.G × B),
// floor(foreground.B × B) and alpha equal to the original foreground alpha unchanged.

func TestProperty19_ContinuousBrightnessScalesRGBPreservingAlpha(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		shape := Shape(rapid.IntRange(0, 3).Draw(t, "shape"))
		diameter := rapid.IntRange(5, 40).Draw(t, "diameter")

		// Use millibright to get values strictly in (0.0, 1.0) exclusive
		milliBright := rapid.IntRange(1, 999).Draw(t, "milliBright")
		brightness := float64(milliBright) / 1000.0

		fg := color.RGBA{
			R: uint8(rapid.IntRange(1, 255).Draw(t, "fgR")),
			G: uint8(rapid.IntRange(1, 255).Draw(t, "fgG")),
			B: uint8(rapid.IntRange(1, 255).Draw(t, "fgB")),
			A: uint8(rapid.IntRange(1, 255).Draw(t, "fgA")),
		}
		// Use a distinct background so we can differentiate body pixels
		bg := color.RGBA{R: 5, G: 5, B: 5, A: 255}

		cfg := Config{
			Shape:       shape,
			State:       On, // ignored when brightness is set
			Brightness:  brightness,
			Diameter:    diameter,
			Bounds:      image.Rect(0, 0, diameter, diameter),
			Foreground:  fg,
			Background:  bg,
			BorderWidth: 0,
			GlowEnabled: false,
		}

		result := Render(cfg)
		if result == nil {
			t.Fatalf("expected non-nil result for diameter=%d", diameter)
		}

		img := result.Image.(*image.RGBA)
		bounds := img.Bounds()

		expectedR := uint8(math.Floor(float64(fg.R) * brightness))
		expectedG := uint8(math.Floor(float64(fg.G) * brightness))
		expectedB := uint8(math.Floor(float64(fg.B) * brightness))
		expectedA := fg.A
		expectedFill := color.RGBA{R: expectedR, G: expectedG, B: expectedB, A: expectedA}

		// Scan all pixels: any non-transparent, non-background pixel should be the expected fill
		foundBodyPixel := false
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				px := img.RGBAAt(x, y)
				if px.A == 0 {
					continue // transparent exterior
				}
				if px == bg {
					continue // background (shouldn't appear in On state but skip if present)
				}
				// This is a body fill pixel
				foundBodyPixel = true
				if px != expectedFill {
					t.Fatalf("body pixel at (%d,%d): got %v, want %v "+
						"[shape=%d, diameter=%d, brightness=%.3f, fg=%v]",
						x, y, px, expectedFill, shape, diameter, brightness, fg)
				}
			}
		}

		if !foundBodyPixel {
			t.Fatalf("no body fill pixels found [shape=%d, diameter=%d, brightness=%.3f]",
				shape, diameter, brightness)
		}
	})
}

// TestProperty28_OffStateContainsNoFullBrightnessForegroundPixels verifies that for any
// Config in Off state (State = Off or Brightness = 0.0), no pixel shall match the full
// foreground color.

func TestProperty28_OffStateContainsNoFullBrightnessForegroundPixels(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		shape := Shape(rapid.IntRange(0, 3).Draw(t, "shape"))
		diameter := rapid.IntRange(3, 40).Draw(t, "diameter")

		fg := color.RGBA{
			R: uint8(rapid.IntRange(10, 255).Draw(t, "fgR")),
			G: uint8(rapid.IntRange(10, 255).Draw(t, "fgG")),
			B: uint8(rapid.IntRange(10, 255).Draw(t, "fgB")),
			A: uint8(rapid.IntRange(1, 255).Draw(t, "fgA")),
		}
		bg := color.RGBA{
			R: uint8(rapid.IntRange(0, 9).Draw(t, "bgR")),
			G: uint8(rapid.IntRange(0, 9).Draw(t, "bgG")),
			B: uint8(rapid.IntRange(0, 9).Draw(t, "bgB")),
			A: 255,
		}

		// Optional border
		borderWidth := rapid.IntRange(0, 3).Draw(t, "borderWidth")
		borderColor := color.RGBA{
			R: uint8(rapid.IntRange(50, 200).Draw(t, "borderR")),
			G: uint8(rapid.IntRange(50, 200).Draw(t, "borderG")),
			B: uint8(rapid.IntRange(50, 200).Draw(t, "borderB")),
			A: 255,
		}

		// Choose between discrete Off or continuous Brightness=0.0
		useDiscreteOff := rapid.Bool().Draw(t, "useDiscreteOff")
		var brightness float64
		var state State
		if useDiscreteOff {
			brightness = -1.0
			state = Off
		} else {
			brightness = 0.0
			state = On // ignored when brightness is set
		}

		cfg := Config{
			Shape:       shape,
			State:       state,
			Brightness:  brightness,
			Diameter:    diameter,
			Bounds:      image.Rect(0, 0, diameter, diameter),
			Foreground:  fg,
			Background:  bg,
			BorderWidth: borderWidth,
			BorderColor: borderColor,
			GlowEnabled: false,
		}

		result := Render(cfg)
		if result == nil {
			t.Fatalf("expected non-nil result for diameter=%d", diameter)
		}

		img := result.Image.(*image.RGBA)
		bounds := img.Bounds()

		// Resolve foreground in case it was zero-value (but we generate non-zero)
		resolvedFg := fg

		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				px := img.RGBAAt(x, y)
				if px == resolvedFg {
					t.Fatalf("found full-brightness foreground pixel at (%d,%d): %v "+
						"[shape=%d, diameter=%d, off=%v, fg=%v]",
						x, y, px, shape, diameter, useDiscreteOff, fg)
				}
			}
		}
	})
}

// TestProperty29_OffStateBorderAtFullColor verifies that for any Config in Off state
// with a border configured, border pixels SHALL be at full border color (no dimming).

func TestProperty29_OffStateBorderAtFullColor(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		shape := Shape(rapid.IntRange(0, 3).Draw(t, "shape"))
		diameter := rapid.IntRange(10, 40).Draw(t, "diameter")

		fg := color.RGBA{
			R: uint8(rapid.IntRange(10, 255).Draw(t, "fgR")),
			G: uint8(rapid.IntRange(10, 255).Draw(t, "fgG")),
			B: uint8(rapid.IntRange(10, 255).Draw(t, "fgB")),
			A: 255,
		}
		bg := color.RGBA{
			R: uint8(rapid.IntRange(0, 9).Draw(t, "bgR")),
			G: uint8(rapid.IntRange(0, 9).Draw(t, "bgG")),
			B: uint8(rapid.IntRange(0, 9).Draw(t, "bgB")),
			A: 255,
		}

		// Border width must be > 0 for this property
		borderWidth := rapid.IntRange(1, 3).Draw(t, "borderWidth")
		// Use a distinctive border color that can't be confused with other colors
		borderColor := color.RGBA{
			R: uint8(rapid.IntRange(100, 255).Draw(t, "borderR")),
			G: uint8(rapid.IntRange(100, 255).Draw(t, "borderG")),
			B: uint8(rapid.IntRange(100, 255).Draw(t, "borderB")),
			A: 255,
		}

		cfg := Config{
			Shape:       shape,
			State:       Off,
			Brightness:  -1.0,
			Diameter:    diameter,
			Bounds:      image.Rect(0, 0, diameter, diameter),
			Foreground:  fg,
			Background:  bg,
			BorderWidth: borderWidth,
			BorderColor: borderColor,
			GlowEnabled: false,
		}

		result := Render(cfg)
		if result == nil {
			t.Fatalf("expected non-nil result for diameter=%d", diameter)
		}

		img := result.Image.(*image.RGBA)
		bounds := img.Bounds()

		// Verify that border pixels exist at full border color
		foundBorder := false
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				px := img.RGBAAt(x, y)
				if px == borderColor {
					foundBorder = true
					break
				}
			}
			if foundBorder {
				break
			}
		}

		if !foundBorder {
			t.Fatalf("no full-color border pixels found "+
				"[shape=%d, diameter=%d, borderWidth=%d, borderColor=%v]",
				shape, diameter, borderWidth, borderColor)
		}

		// Verify that no dimmed border pixels exist (border should NOT be dimmed)
		dimmedBorder := dimColor(borderColor)
		// Only check if dimmedBorder is distinct from other expected colors
		dimmedFg := dimColor(fg)
		if dimmedBorder != dimmedFg && dimmedBorder != bg && dimmedBorder != (color.RGBA{}) {
			for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
				for x := bounds.Min.X; x < bounds.Max.X; x++ {
					px := img.RGBAAt(x, y)
					if px == dimmedBorder {
						t.Fatalf("found dimmed border pixel at (%d,%d): %v "+
							"(border should be at full color, not dimmed) "+
							"[shape=%d, diameter=%d, borderColor=%v]",
							x, y, px, shape, diameter, borderColor)
					}
				}
			}
		}
	})
}

// TestProperty30_DefaultColorResolutionAppliedBeforeRenderingLogic verifies that for any
// Config with zero-value Foreground, the resolved foreground SHALL be (0, 200, 0, 255)
// and subsequent dimming operates on this resolved value.

func TestProperty30_DefaultColorResolutionAppliedBeforeRenderingLogic(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		shape := Shape(rapid.IntRange(0, 3).Draw(t, "shape"))
		diameter := rapid.IntRange(5, 40).Draw(t, "diameter")

		// Zero-value foreground — should resolve to (0, 200, 0, 255)
		zeroFg := color.RGBA{R: 0, G: 0, B: 0, A: 0}

		// Zero-value background — should resolve to (0, 0, 0, 255)
		zeroBg := color.RGBA{R: 0, G: 0, B: 0, A: 0}

		cfg := Config{
			Shape:       shape,
			State:       Off,
			Brightness:  -1.0,
			Diameter:    diameter,
			Bounds:      image.Rect(0, 0, diameter, diameter),
			Foreground:  zeroFg,
			Background:  zeroBg,
			BorderWidth: 0,
			GlowEnabled: false,
		}

		result := Render(cfg)
		if result == nil {
			t.Fatalf("expected non-nil result for diameter=%d", diameter)
		}

		img := result.Image.(*image.RGBA)
		bounds := img.Bounds()

		// Resolved foreground: (0, 200, 0, 255)
		resolvedFg := color.RGBA{R: 0, G: 200, B: 0, A: 255}
		// Dimmed resolved foreground: floor(0*0.3)=0, floor(200*0.3)=60, floor(0*0.3)=0
		expectedDimmed := color.RGBA{
			R: uint8(math.Floor(float64(resolvedFg.R) * 0.3)),
			G: uint8(math.Floor(float64(resolvedFg.G) * 0.3)),
			B: uint8(math.Floor(float64(resolvedFg.B) * 0.3)),
			A: resolvedFg.A,
		}
		// Resolved background: (0, 0, 0, 255)
		resolvedBg := color.RGBA{R: 0, G: 0, B: 0, A: 255}

		// Verify dimmed outline pixels exist (from resolved foreground, not zero)
		foundDimmed := false
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				px := img.RGBAAt(x, y)
				if px == expectedDimmed {
					foundDimmed = true
					break
				}
			}
			if foundDimmed {
				break
			}
		}

		if !foundDimmed {
			t.Fatalf("no dimmed outline pixels found from resolved foreground "+
				"(expected dimmed green %v) [shape=%d, diameter=%d]",
				expectedDimmed, shape, diameter)
		}

		// Verify all non-transparent pixels are either dimmed outline, resolved background, or transparent
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				px := img.RGBAAt(x, y)
				if px.A == 0 {
					continue // transparent exterior
				}
				if px == expectedDimmed || px == resolvedBg {
					continue // expected off-state pixels
				}
				// No pixel should be the unresolve zero-value foreground
				if px == zeroFg {
					t.Fatalf("found zero-value foreground pixel at (%d,%d): %v "+
						"(default color resolution should have been applied before rendering) "+
						"[shape=%d, diameter=%d]",
						x, y, px, shape, diameter)
				}
				// Any other pixel is unexpected
				t.Fatalf("unexpected pixel at (%d,%d): %v "+
					"(expected only dimmed=%v or bg=%v or transparent) "+
					"[shape=%d, diameter=%d]",
					x, y, px, expectedDimmed, resolvedBg, shape, diameter)
			}
		}
	})
}
