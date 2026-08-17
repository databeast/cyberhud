package led

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"testing"

	"github.com/databeast/cyberhud/display/widgets/gradient"
)

// ============================================================================
// ============================================================================

// TestOffStateOutlineDimmedColor verifies that Off state draws a 1px outline
// at dimmed color (floor(R×0.3), floor(G×0.3), floor(B×0.3), A) for all shapes.

func TestOffStateOutlineDimmedColor(t *testing.T) {
	shapes := []struct {
		name  string
		shape Shape
	}{
		{"Circle", Circle},
		{"Square", Square},
		{"Diamond", Diamond},
		{"RoundedSquare", RoundedSquare},
	}

	fg := color.RGBA{R: 200, G: 100, B: 50, A: 255}
	bg := color.RGBA{R: 10, G: 10, B: 10, A: 255}
	expectedDimmed := color.RGBA{
		R: uint8(math.Floor(float64(fg.R) * 0.3)),
		G: uint8(math.Floor(float64(fg.G) * 0.3)),
		B: uint8(math.Floor(float64(fg.B) * 0.3)),
		A: fg.A,
	}

	for _, tc := range shapes {
		t.Run(tc.name, func(t *testing.T) {
			cfg := Config{
				Shape:      tc.shape,
				State:      Off,
				Brightness: -1.0,
				Diameter:   20,
				Bounds:     image.Rect(0, 0, 20, 20),
				Foreground: fg,
				Background: bg,
			}
			result := Render(cfg)
			if result == nil {
				t.Fatal("expected non-nil result")
			}

			img := result.Image
			bounds := img.Bounds()

			// Scan for outline pixels: pixels that are the dimmed color
			foundDimmed := false
			foundInterior := false
			for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
				for x := bounds.Min.X; x < bounds.Max.X; x++ {
					r, g, b, a := img.At(x, y).RGBA()
					px := color.RGBA{
						R: uint8(r >> 8), G: uint8(g >> 8),
						B: uint8(b >> 8), A: uint8(a >> 8),
					}
					if px == expectedDimmed {
						foundDimmed = true
					}
					if px == bg {
						foundInterior = true
					}
				}
			}
			if !foundDimmed {
				t.Errorf("shape %s: no dimmed outline pixels found (expected %v)",
					tc.name, expectedDimmed)
			}
			if !foundInterior {
				t.Errorf("shape %s: no background interior pixels found (expected %v)",
					tc.name, bg)
			}

			// Verify NO full-brightness foreground pixels exist
			for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
				for x := bounds.Min.X; x < bounds.Max.X; x++ {
					r, g, b, a := img.At(x, y).RGBA()
					px := color.RGBA{
						R: uint8(r >> 8), G: uint8(g >> 8),
						B: uint8(b >> 8), A: uint8(a >> 8),
					}
					if px == fg {
						t.Fatalf("shape %s: found full-brightness foreground pixel at (%d,%d)",
							tc.name, x, y)
					}
				}
			}
		})
	}
}

// TestBrightness0EqualsStateOff verifies that Brightness=0.0 renders
// pixel-identically to State=Off (Brightness=-1.0) for all shapes.

func TestBrightness0EqualsStateOff(t *testing.T) {
	shapes := []struct {
		name  string
		shape Shape
	}{
		{"Circle", Circle},
		{"Square", Square},
		{"Diamond", Diamond},
		{"RoundedSquare", RoundedSquare},
	}

	fg := color.RGBA{R: 150, G: 80, B: 200, A: 255}
	bg := color.RGBA{R: 20, G: 20, B: 20, A: 255}

	for _, tc := range shapes {
		t.Run(tc.name, func(t *testing.T) {
			// Render with State=Off, Brightness=-1.0 (discrete)
			cfgOff := Config{
				Shape:      tc.shape,
				State:      Off,
				Brightness: -1.0,
				Diameter:   15,
				Bounds:     image.Rect(0, 0, 15, 15),
				Foreground: fg,
				Background: bg,
			}
			resultOff := Render(cfgOff)

			// Render with Brightness=0.0 (continuous)
			cfgB0 := Config{
				Shape:      tc.shape,
				State:      On, // State should be ignored
				Brightness: 0.0,
				Diameter:   15,
				Bounds:     image.Rect(0, 0, 15, 15),
				Foreground: fg,
				Background: bg,
			}
			resultB0 := Render(cfgB0)

			if resultOff == nil || resultB0 == nil {
				t.Fatal("expected non-nil results")
			}

			// Compare pixel-by-pixel
			boundsOff := resultOff.Image.Bounds()
			boundsB0 := resultB0.Image.Bounds()
			if boundsOff != boundsB0 {
				t.Fatalf("image bounds differ: Off=%v, B0=%v",
					boundsOff, boundsB0)
			}

			for y := boundsOff.Min.Y; y < boundsOff.Max.Y; y++ {
				for x := boundsOff.Min.X; x < boundsOff.Max.X; x++ {
					r1, g1, b1, a1 := resultOff.Image.At(x, y).RGBA()
					r2, g2, b2, a2 := resultB0.Image.At(x, y).RGBA()
					if r1 != r2 || g1 != g2 || b1 != b2 || a1 != a2 {
						t.Fatalf("shape %s: pixel mismatch at (%d,%d): Off=(%d,%d,%d,%d) B0=(%d,%d,%d,%d)",
							tc.name, x, y, r1, g1, b1, a1, r2, g2, b2, a2)
					}
				}
			}
		})
	}
}

// TestBrightness1EqualsStateOn verifies that Brightness=1.0 renders
// pixel-identically to State=On (Brightness=-1.0) for all shapes.

func TestBrightness1EqualsStateOn(t *testing.T) {
	shapes := []struct {
		name  string
		shape Shape
	}{
		{"Circle", Circle},
		{"Square", Square},
		{"Diamond", Diamond},
		{"RoundedSquare", RoundedSquare},
	}

	fg := color.RGBA{R: 150, G: 80, B: 200, A: 255}
	bg := color.RGBA{R: 20, G: 20, B: 20, A: 255}

	for _, tc := range shapes {
		t.Run(tc.name, func(t *testing.T) {
			// Render with State=On, Brightness=-1.0 (discrete)
			cfgOn := Config{
				Shape:      tc.shape,
				State:      On,
				Brightness: -1.0,
				Diameter:   15,
				Bounds:     image.Rect(0, 0, 15, 15),
				Foreground: fg,
				Background: bg,
			}
			resultOn := Render(cfgOn)

			// Render with Brightness=1.0 (continuous)
			cfgB1 := Config{
				Shape:      tc.shape,
				State:      Off, // State should be ignored
				Brightness: 1.0,
				Diameter:   15,
				Bounds:     image.Rect(0, 0, 15, 15),
				Foreground: fg,
				Background: bg,
			}
			resultB1 := Render(cfgB1)

			if resultOn == nil || resultB1 == nil {
				t.Fatal("expected non-nil results")
			}

			// Compare pixel-by-pixel
			boundsOn := resultOn.Image.Bounds()
			boundsB1 := resultB1.Image.Bounds()
			if boundsOn != boundsB1 {
				t.Fatalf("image bounds differ: On=%v, B1=%v",
					boundsOn, boundsB1)
			}

			for y := boundsOn.Min.Y; y < boundsOn.Max.Y; y++ {
				for x := boundsOn.Min.X; x < boundsOn.Max.X; x++ {
					r1, g1, b1, a1 := resultOn.Image.At(x, y).RGBA()
					r2, g2, b2, a2 := resultB1.Image.At(x, y).RGBA()
					if r1 != r2 || g1 != g2 || b1 != b2 || a1 != a2 {
						t.Fatalf("shape %s: pixel mismatch at (%d,%d): On=(%d,%d,%d,%d) B1=(%d,%d,%d,%d)",
							tc.name, x, y, r1, g1, b1, a1, r2, g2, b2, a2)
					}
				}
			}
		})
	}
}

// TestContinuousBrightnessScalesRGB verifies that continuous brightness B
// scales each RGB channel: floor(channel × B), alpha unchanged.

func TestContinuousBrightnessScalesRGB(t *testing.T) {
	fg := color.RGBA{R: 200, G: 100, B: 50, A: 200}
	bg := color.RGBA{R: 10, G: 10, B: 10, A: 255}
	brightness := 0.5

	expectedScaled := color.RGBA{
		R: uint8(math.Floor(float64(fg.R) * brightness)),
		G: uint8(math.Floor(float64(fg.G) * brightness)),
		B: uint8(math.Floor(float64(fg.B) * brightness)),
		A: fg.A,
	}

	cfg := Config{
		Shape:      Circle,
		State:      On,
		Brightness: brightness,
		Diameter:   20,
		Bounds:     image.Rect(0, 0, 20, 20),
		Foreground: fg,
		Background: bg,
	}
	result := Render(cfg)
	if result == nil {
		t.Fatal("expected non-nil result")
	}

	// Center pixel should be brightness-scaled foreground
	cx, cy := 10, 10
	r, g, b, a := result.Image.At(cx, cy).RGBA()
	centerPx := color.RGBA{
		R: uint8(r >> 8), G: uint8(g >> 8),
		B: uint8(b >> 8), A: uint8(a >> 8),
	}
	if centerPx != expectedScaled {
		t.Errorf("center pixel: got %v, want %v", centerPx, expectedScaled)
	}

	// Alpha must be preserved (unchanged)
	if uint8(a>>8) != fg.A {
		t.Errorf("alpha changed: got %d, want %d", uint8(a>>8), fg.A)
	}
}

// TestOffStateNoGlowPixels verifies that glow is suppressed when LED is Off,
// even when GlowEnabled=true.

func TestOffStateNoGlowPixels(t *testing.T) {
	fg := color.RGBA{R: 0, G: 200, B: 0, A: 255}
	bg := color.RGBA{R: 10, G: 10, B: 10, A: 255}

	cfg := Config{
		Shape:       Circle,
		State:       Off,
		Brightness:  -1.0,
		Diameter:    15,
		Bounds:      image.Rect(0, 0, 15, 15),
		Foreground:  fg,
		Background:  bg,
		GlowEnabled: true,
		GlowRadius:  5,
	}
	result := Render(cfg)
	if result == nil {
		t.Fatal("expected non-nil result")
	}

	img := result.Image
	bounds := img.Bounds()
	// Output size should include glow padding: 15 + 2*5 = 25
	if bounds.Dx() != 25 || bounds.Dy() != 25 {
		t.Fatalf("expected 25×25 image (with glow padding), got %d×%d",
			bounds.Dx(), bounds.Dy())
	}

	// Check that the glow region (outside the body) has no glow pixels
	dimmed := dimColor(fg)
	center := float64(bounds.Dx()) / 2.0
	bodyRadius := float64(cfg.Diameter) / 2.0

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			px := color.RGBA{
				R: uint8(r >> 8), G: uint8(g >> 8),
				B: uint8(b >> 8), A: uint8(a >> 8),
			}

			// Skip transparent, dimmed, bg, and border pixels
			if a == 0 {
				continue // transparent - ok
			}
			if px == dimmed || px == bg {
				continue // expected off-state pixels
			}

			// Check if this pixel is in the glow region (outside body)
			dx := float64(x) + 0.5 - center
			dy := float64(y) + 0.5 - center
			dist := math.Sqrt(dx*dx + dy*dy)
			if dist > bodyRadius {
				// This pixel is outside the body - should be transparent
				if a != 0 {
					t.Fatalf("glow pixel found at (%d,%d) dist=%.1f: color=%v (should be transparent in Off state)",
						x, y, dist, px)
				}
			}
		}
	}
}

// TestOffStateNoGradientPixels verifies that gradient is ignored when LED is Off.

func TestOffStateNoGradientPixels(t *testing.T) {
	fg := color.RGBA{R: 200, G: 0, B: 0, A: 255}
	bg := color.RGBA{R: 10, G: 10, B: 10, A: 255}

	cfg := Config{
		Shape:      Circle,
		State:      Off,
		Brightness: -1.0,
		Diameter:   20,
		Bounds:     image.Rect(0, 0, 20, 20),
		Foreground: fg,
		Background: bg,
		Gradient: &GradientConfig{
			Stops: []gradient.ColorStop{
				{Position: 0.0, Color: color.RGBA{R: 255, G: 0, B: 0, A: 255}},
				{Position: 1.0, Color: color.RGBA{R: 0, G: 0, B: 255, A: 255}},
			},
		},
	}

	// Render with gradient (Off state)
	resultGrad := Render(cfg)
	if resultGrad == nil {
		t.Fatal("expected non-nil result")
	}

	// Render without gradient (Off state)
	cfgNoGrad := cfg
	cfgNoGrad.Gradient = nil
	resultNoGrad := Render(cfgNoGrad)
	if resultNoGrad == nil {
		t.Fatal("expected non-nil result")
	}

	// They should be pixel-identical
	bounds := resultGrad.Image.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r1, g1, b1, a1 := resultGrad.Image.At(x, y).RGBA()
			r2, g2, b2, a2 := resultNoGrad.Image.At(x, y).RGBA()
			if r1 != r2 || g1 != g2 || b1 != b2 || a1 != a2 {
				t.Fatalf("pixel mismatch at (%d,%d): withGrad=(%d,%d,%d,%d) noGrad=(%d,%d,%d,%d)",
					x, y, r1, g1, b1, a1, r2, g2, b2, a2)
			}
		}
	}
}

// TestOffStateNoShinePixels verifies that shine is suppressed when LED is Off.

func TestOffStateNoShinePixels(t *testing.T) {
	fg := color.RGBA{R: 0, G: 200, B: 0, A: 255}
	bg := color.RGBA{R: 10, G: 10, B: 10, A: 255}

	for _, style := range []ShineStyle{ShineDot, ShineCrescent} {
		name := "Dot"
		if style == ShineCrescent {
			name = "Crescent"
		}
		t.Run(name, func(t *testing.T) {
			cfg := Config{
				Shape:      Circle,
				State:      Off,
				Brightness: -1.0,
				Diameter:   20,
				Bounds:     image.Rect(0, 0, 20, 20),
				Foreground: fg,
				Background: bg,
				ShineStyle: style,
			}
			result := Render(cfg)
			if result == nil {
				t.Fatal("expected non-nil result")
			}

			img := result.Image
			bounds := img.Bounds()

			// Scan for white/shine pixels (shine is white with variable alpha)
			for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
				for x := bounds.Min.X; x < bounds.Max.X; x++ {
					r, g, b, a := img.At(x, y).RGBA()
					// Shine pixels are white (R=G=B=255) with non-zero alpha
					if uint8(r>>8) == 255 && uint8(g>>8) == 255 &&
						uint8(b>>8) == 255 && a > 0 {
						t.Fatalf("shine pixel found at (%d,%d) in Off state with style %s",
							x, y, name)
					}
				}
			}
		})
	}
}

// TestOffStateBorderFullColor verifies that border is rendered at full configured
// color (no dimming) when LED is Off.

func TestOffStateBorderFullColor(t *testing.T) {
	fg := color.RGBA{R: 0, G: 200, B: 0, A: 255}
	bg := color.RGBA{R: 10, G: 10, B: 10, A: 255}
	borderColor := color.RGBA{R: 200, G: 50, B: 50, A: 255}

	shapes := []struct {
		name  string
		shape Shape
	}{
		{"Circle", Circle},
		{"Square", Square},
		{"Diamond", Diamond},
		{"RoundedSquare", RoundedSquare},
	}

	for _, tc := range shapes {
		t.Run(tc.name, func(t *testing.T) {
			cfg := Config{
				Shape:       tc.shape,
				State:       Off,
				Brightness:  -1.0,
				Diameter:    20,
				Bounds:      image.Rect(0, 0, 20, 20),
				Foreground:  fg,
				Background:  bg,
				BorderWidth: 2,
				BorderColor: borderColor,
			}
			result := Render(cfg)
			if result == nil {
				t.Fatal("expected non-nil result")
			}

			img := result.Image
			bounds := img.Bounds()

			// Verify border pixels exist at full color
			foundBorder := false
			for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
				for x := bounds.Min.X; x < bounds.Max.X; x++ {
					r, g, b, a := img.At(x, y).RGBA()
					px := color.RGBA{
						R: uint8(r >> 8), G: uint8(g >> 8),
						B: uint8(b >> 8), A: uint8(a >> 8),
					}
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
				t.Errorf("shape %s: no full-color border pixels found (expected %v)",
					tc.name, borderColor)
			}

			// Verify border is NOT dimmed
			dimmedBorder := dimColor(borderColor)
			for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
				for x := bounds.Min.X; x < bounds.Max.X; x++ {
					r, g, b, a := img.At(x, y).RGBA()
					px := color.RGBA{
						R: uint8(r >> 8), G: uint8(g >> 8),
						B: uint8(b >> 8), A: uint8(a >> 8),
					}
					if px == dimmedBorder {
						// dimmedBorder is (60, 15, 15, 255).
						// Make sure this isn't also the dimmed fg.
						dimmedFg := dimColor(fg)
						if px != dimmedFg {
							t.Errorf("shape %s: found dimmed border pixel at (%d,%d): %v (border should NOT be dimmed)",
								tc.name, x, y, px)
						}
					}
				}
			}
		})
	}
}

// TestContinuousBrightnessMultipleValues verifies brightness scaling at several
// values for all shapes.

func TestContinuousBrightnessMultipleValues(t *testing.T) {
	fg := color.RGBA{R: 100, G: 200, B: 50, A: 180}
	bg := color.RGBA{R: 10, G: 10, B: 10, A: 255}

	brightnesses := []float64{0.1, 0.25, 0.5, 0.75, 0.9}

	for _, b := range brightnesses {
		t.Run("brightness_"+formatFloat(b), func(t *testing.T) {
			cfg := Config{
				Shape:      Circle,
				State:      On,
				Brightness: b,
				Diameter:   20,
				Bounds:     image.Rect(0, 0, 20, 20),
				Foreground: fg,
				Background: bg,
			}
			result := Render(cfg)
			if result == nil {
				t.Fatal("expected non-nil result")
			}

			// Check center pixel has correct scaling
			expectedR := uint8(math.Floor(float64(fg.R) * b))
			expectedG := uint8(math.Floor(float64(fg.G) * b))
			expectedB := uint8(math.Floor(float64(fg.B) * b))
			expectedA := fg.A

			r, g, bb, a := result.Image.At(10, 10).RGBA()
			px := color.RGBA{
				R: uint8(r >> 8), G: uint8(g >> 8),
				B: uint8(bb >> 8), A: uint8(a >> 8),
			}
			expected := color.RGBA{R: expectedR, G: expectedG, B: expectedB, A: expectedA}
			if px != expected {
				t.Errorf("brightness %.2f: center pixel got %v, want %v",
					b, px, expected)
			}
		})
	}
}

func formatFloat(f float64) string {
	return fmt.Sprintf("%.2f", f)
}

// TestBrightness0EqualsStateOffWithEffects verifies pixel-identical output
// between Brightness=0.0 and State=Off when glow, gradient, and shine are configured.
// These effects should all be suppressed identically in both cases.

func TestBrightness0EqualsStateOffWithEffects(t *testing.T) {
	fg := color.RGBA{R: 0, G: 200, B: 0, A: 255}
	bg := color.RGBA{R: 10, G: 10, B: 10, A: 255}

	baseConfig := Config{
		Shape:       Circle,
		Diameter:    20,
		Bounds:      image.Rect(0, 0, 20, 20),
		Foreground:  fg,
		Background:  bg,
		GlowEnabled: true,
		GlowRadius:  4,
		ShineStyle:  ShineDot,
		BorderWidth: 2,
		BorderColor: color.RGBA{R: 128, G: 128, B: 128, A: 255},
		Gradient: &GradientConfig{
			Stops: []gradient.ColorStop{
				{Position: 0.0, Color: color.RGBA{R: 255, G: 0, B: 0, A: 255}},
				{Position: 1.0, Color: color.RGBA{R: 0, G: 0, B: 255, A: 255}},
			},
		},
	}

	// State=Off, Brightness=-1.0 (discrete)
	cfgOff := baseConfig
	cfgOff.State = Off
	cfgOff.Brightness = -1.0
	resultOff := Render(cfgOff)

	// Brightness=0.0 (continuous, state ignored)
	cfgB0 := baseConfig
	cfgB0.State = On
	cfgB0.Brightness = 0.0
	resultB0 := Render(cfgB0)

	if resultOff == nil || resultB0 == nil {
		t.Fatal("expected non-nil results")
	}

	boundsOff := resultOff.Image.Bounds()
	boundsB0 := resultB0.Image.Bounds()
	if boundsOff != boundsB0 {
		t.Fatalf("bounds differ: Off=%v, B0=%v", boundsOff, boundsB0)
	}

	for y := boundsOff.Min.Y; y < boundsOff.Max.Y; y++ {
		for x := boundsOff.Min.X; x < boundsOff.Max.X; x++ {
			r1, g1, b1, a1 := resultOff.Image.At(x, y).RGBA()
			r2, g2, b2, a2 := resultB0.Image.At(x, y).RGBA()
			if r1 != r2 || g1 != g2 || b1 != b2 || a1 != a2 {
				t.Fatalf("pixel mismatch at (%d,%d): Off=(%d,%d,%d,%d) B0=(%d,%d,%d,%d)",
					x, y, r1, g1, b1, a1, r2, g2, b2, a2)
			}
		}
	}
}

// TestBrightness1EqualsStateOnWithEffects verifies pixel-identical output
// between Brightness=1.0 and State=On when glow, gradient, and shine are configured.

func TestBrightness1EqualsStateOnWithEffects(t *testing.T) {
	fg := color.RGBA{R: 0, G: 200, B: 0, A: 255}
	bg := color.RGBA{R: 10, G: 10, B: 10, A: 255}

	baseConfig := Config{
		Shape:       Circle,
		Diameter:    20,
		Bounds:      image.Rect(0, 0, 20, 20),
		Foreground:  fg,
		Background:  bg,
		GlowEnabled: true,
		GlowRadius:  4,
		ShineStyle:  ShineDot,
		BorderWidth: 2,
		BorderColor: color.RGBA{R: 128, G: 128, B: 128, A: 255},
		Gradient: &GradientConfig{
			Stops: []gradient.ColorStop{
				{Position: 0.0, Color: color.RGBA{R: 255, G: 0, B: 0, A: 255}},
				{Position: 1.0, Color: color.RGBA{R: 0, G: 0, B: 255, A: 255}},
			},
		},
	}

	// State=On, Brightness=-1.0 (discrete)
	cfgOn := baseConfig
	cfgOn.State = On
	cfgOn.Brightness = -1.0
	resultOn := Render(cfgOn)

	// Brightness=1.0 (continuous, state ignored)
	cfgB1 := baseConfig
	cfgB1.State = Off
	cfgB1.Brightness = 1.0
	resultB1 := Render(cfgB1)

	if resultOn == nil || resultB1 == nil {
		t.Fatal("expected non-nil results")
	}

	boundsOn := resultOn.Image.Bounds()
	boundsB1 := resultB1.Image.Bounds()
	if boundsOn != boundsB1 {
		t.Fatalf("bounds differ: On=%v, B1=%v", boundsOn, boundsB1)
	}

	for y := boundsOn.Min.Y; y < boundsOn.Max.Y; y++ {
		for x := boundsOn.Min.X; x < boundsOn.Max.X; x++ {
			r1, g1, b1, a1 := resultOn.Image.At(x, y).RGBA()
			r2, g2, b2, a2 := resultB1.Image.At(x, y).RGBA()
			if r1 != r2 || g1 != g2 || b1 != b2 || a1 != a2 {
				t.Fatalf("pixel mismatch at (%d,%d): On=(%d,%d,%d,%d) B1=(%d,%d,%d,%d)",
					x, y, r1, g1, b1, a1, r2, g2, b2, a2)
			}
		}
	}
}

// TestOffStateNoFullBrightnessForeground verifies that no pixel in the output
// matches the full (undimmed) foreground color when LED is Off.

func TestOffStateNoFullBrightnessForeground(t *testing.T) {
	fg := color.RGBA{R: 0, G: 200, B: 0, A: 255}
	bg := color.RGBA{R: 10, G: 10, B: 10, A: 255}

	shapes := []struct {
		name  string
		shape Shape
	}{
		{"Circle", Circle},
		{"Square", Square},
		{"Diamond", Diamond},
		{"RoundedSquare", RoundedSquare},
	}

	for _, tc := range shapes {
		t.Run(tc.name, func(t *testing.T) {
			cfg := Config{
				Shape:      tc.shape,
				State:      Off,
				Brightness: -1.0,
				Diameter:   20,
				Bounds:     image.Rect(0, 0, 20, 20),
				Foreground: fg,
				Background: bg,
			}
			result := Render(cfg)
			if result == nil {
				t.Fatal("expected non-nil result")
			}

			img := result.Image
			bounds := img.Bounds()
			for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
				for x := bounds.Min.X; x < bounds.Max.X; x++ {
					r, g, b, a := img.At(x, y).RGBA()
					px := color.RGBA{
						R: uint8(r >> 8), G: uint8(g >> 8),
						B: uint8(b >> 8), A: uint8(a >> 8),
					}
					if px == fg {
						t.Fatalf("shape %s: full foreground pixel at (%d,%d) in Off state",
							tc.name, x, y)
					}
				}
			}
		})
	}
}
