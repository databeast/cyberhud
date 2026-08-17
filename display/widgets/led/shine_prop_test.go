package led

import (
	"image"
	"image/color"
	"math"
	"testing"

	"pgregory.net/rapid"
)

// TestProperty25_ShineSuppressedWhenOffOrDiameterLessThan5 verifies that for any Config
// with shine enabled (Dot or Crescent) where the LED is in Off state (or Brightness = 0.0),
// OR Diameter < 5, the output SHALL contain no shine-colored pixels (no white/highlight
// pixels at the shine position).

func TestProperty25_ShineSuppressedWhenOffOrDiameterLessThan5(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Choose a shine style (Dot or Crescent)
		shineStyle := ShineStyle(rapid.IntRange(1, 2).Draw(t, "shineStyle"))

		// Choose a suppression reason: either Off/zero-brightness OR small diameter
		suppressionReason := rapid.IntRange(0, 1).Draw(t, "suppressionReason")

		var diameter int
		var brightness float64
		var state State

		if suppressionReason == 0 {
			// Suppression due to Off state or Brightness == 0.0
			diameter = rapid.IntRange(5, 64).Draw(t, "diameter") // valid diameter ≥ 5
			// Either use discrete Off state or explicit brightness 0.0
			useDiscreteOff := rapid.Bool().Draw(t, "useDiscreteOff")
			if useDiscreteOff {
				brightness = -1.0 // sentinel: use discrete state
				state = Off
			} else {
				brightness = 0.0
				state = On // state ignored when brightness is set explicitly
			}
		} else {
			// Suppression due to Diameter < 5
			diameter = rapid.IntRange(3, 4).Draw(t, "diameter") // 3 or 4
			// LED is On so shine WOULD be drawn if diameter were large enough
			brightness = -1.0
			state = On
		}

		// Use a non-white foreground so white pixels are distinguishable
		fg := color.RGBA{
			R: uint8(rapid.IntRange(1, 200).Draw(t, "fgR")),
			G: uint8(rapid.IntRange(1, 200).Draw(t, "fgG")),
			B: uint8(rapid.IntRange(1, 50).Draw(t, "fgB")),
			A: 255,
		}

		shineOpacity := uint8(rapid.IntRange(0, 255).Draw(t, "shineOpacity"))

		cfg := Config{
			Shape:        Circle,
			State:        state,
			Brightness:   brightness,
			Diameter:     diameter,
			Bounds:       image.Rect(0, 0, diameter, diameter),
			Foreground:   fg,
			ShineStyle:   shineStyle,
			ShineOpacity: shineOpacity,
			// No border, no glow, no gradient for simplicity
			BorderWidth: 0,
			GlowEnabled: false,
		}

		result := Render(cfg)
		if result == nil {
			t.Fatalf("expected non-nil result for Diameter=%d", diameter)
		}

		img := result.Image.(*image.RGBA)
		bounds := img.Bounds()

		// Scan all pixels for any white/highlight-colored pixels.
		// Shine is rendered as white (R:255, G:255, B:255) with variable alpha.
		// No pixel should have R==255 && G==255 && B==255 with A > 0 (shine color).
		for py := bounds.Min.Y; py < bounds.Max.Y; py++ {
			for px := bounds.Min.X; px < bounds.Max.X; px++ {
				c := img.RGBAAt(px, py)
				if c.R == 255 && c.G == 255 && c.B == 255 && c.A > 0 {
					t.Fatalf("found shine pixel at (%d,%d) RGBA(%d,%d,%d,%d) but shine should be suppressed "+
						"[diameter=%d, brightness=%.1f, state=%d, shineStyle=%d, reason=%d]",
						px, py, c.R, c.G, c.B, c.A,
						diameter, brightness, state, shineStyle, suppressionReason)
				}
			}
		}
	})
}

// TestProperty26_ShineOpacityModulatedByBrightness verifies that for any Config with
// shine enabled, Diameter ≥ 5, and effective brightness B in (0.0, 1.0), the shine pixels
// SHALL have alpha equal to floor(configuredOpacity × B), where configuredOpacity defaults
// to 255 when ShineOpacity is 0.

func TestProperty26_ShineOpacityModulatedByBrightness(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Use ShineDot for predictable position
		shineStyle := ShineDot

		// Diameter ≥ 5 (need enough space for shine dot)
		diameter := rapid.IntRange(10, 64).Draw(t, "diameter")

		// Brightness in (0.0, 1.0) exclusive — use millibright to avoid float edges
		milliBright := rapid.IntRange(1, 999).Draw(t, "milliBright")
		brightness := float64(milliBright) / 1000.0

		// Random shine opacity (0 means 255)
		shineOpacity := uint8(rapid.IntRange(0, 255).Draw(t, "shineOpacity"))

		// Use a non-white foreground so white shine pixels are distinguishable
		fg := color.RGBA{
			R: uint8(rapid.IntRange(1, 100).Draw(t, "fgR")),
			G: uint8(rapid.IntRange(100, 200).Draw(t, "fgG")),
			B: uint8(rapid.IntRange(1, 50).Draw(t, "fgB")),
			A: 255,
		}

		cfg := Config{
			Shape:        Circle,
			State:        On,         // ignored when brightness is set
			Brightness:   brightness, // continuous mode
			Diameter:     diameter,
			Bounds:       image.Rect(0, 0, diameter, diameter),
			Foreground:   fg,
			ShineStyle:   shineStyle,
			ShineOpacity: shineOpacity,
			// No border, no glow, no gradient for simplicity
			BorderWidth: 0,
			GlowEnabled: false,
		}

		result := Render(cfg)
		if result == nil {
			t.Fatalf("expected non-nil result for Diameter=%d", diameter)
		}

		img := result.Image.(*image.RGBA)

		// Compute expected alpha.
		resolvedOpacity := shineOpacity
		if resolvedOpacity == 0 {
			resolvedOpacity = 255
		}
		expectedAlpha := uint8(math.Floor(float64(resolvedOpacity) * brightness))

		if expectedAlpha == 0 {
			// If expected alpha is 0, no shine pixels should exist
			// (applyShine returns early when alpha==0)
			return
		}

		// Compute expected shine dot center position.
		// Body rect with no border and no glow: full image area
		bodyRect := image.Rect(0, 0, diameter, diameter)
		bodyRadius := bodyRect.Dx() / 2
		cx := bodyRect.Min.X + bodyRect.Dx()/2
		cy := bodyRect.Min.Y + bodyRect.Dy()/2
		offsetX := int(math.Floor(float64(bodyRadius) * 0.25))
		offsetY := int(math.Floor(float64(bodyRadius) * 0.25))
		dotCx := cx - offsetX
		dotCy := cy - offsetY

		// Find a shine pixel at or near the dot center.
		// The dot center itself should be a shine pixel.
		shinePixel := img.RGBAAt(dotCx, dotCy)

		// Verify it's a shine pixel (white with the computed alpha)
		if shinePixel.R != 255 || shinePixel.G != 255 || shinePixel.B != 255 {
			t.Fatalf("expected shine pixel at dot center (%d,%d) to be white (255,255,255,*), got RGBA(%d,%d,%d,%d) "+
				"[diameter=%d, brightness=%.3f, shineOpacity=%d]",
				dotCx, dotCy, shinePixel.R, shinePixel.G, shinePixel.B, shinePixel.A,
				diameter, brightness, shineOpacity)
		}

		if shinePixel.A != expectedAlpha {
			t.Fatalf("shine pixel at (%d,%d) has alpha %d, expected %d "+
				"(resolvedOpacity=%d × brightness=%.3f = %.3f → floor=%d) "+
				"[diameter=%d, shineOpacity=%d]",
				dotCx, dotCy, shinePixel.A, expectedAlpha,
				resolvedOpacity, brightness, float64(resolvedOpacity)*brightness, expectedAlpha,
				diameter, shineOpacity)
		}
	})
}
