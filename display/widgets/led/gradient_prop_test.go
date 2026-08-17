package led

import (
	"image"
	"image/color"
	"math"
	"testing"

	"github.com/databeast/cyberhud/display/widgets/gradient"
	"pgregory.net/rapid"
)

// TestProperty5_GradientRadialFillCorrectness verifies that for any valid Config
// with a GradientConfig containing ≥ 2 ColorStops and LED in On state, the pixel at
// the body center SHALL match the interpolated color at gradient position 0.0, and
// pixels at the body perimeter SHALL match the interpolated color at position 1.0,
// within ±1 per channel for rounding.

func TestProperty5_GradientRadialFillCorrectness(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Use odd diameters so that the geometric center aligns with a pixel center,
		// giving us a true position-0.0 pixel at the center.
		halfDiameter := rapid.IntRange(4, 32).Draw(t, "halfDiameter")
		diameter := halfDiameter*2 + 1 // Always odd: 9, 11, 13, ... 65

		// Generate two color stops: one at position 0.0 (center) and one at 1.0 (edge).
		centerColor := color.RGBA{
			R: uint8(rapid.IntRange(1, 255).Draw(t, "centerR")),
			G: uint8(rapid.IntRange(1, 255).Draw(t, "centerG")),
			B: uint8(rapid.IntRange(1, 255).Draw(t, "centerB")),
			A: 255,
		}
		edgeColor := color.RGBA{
			R: uint8(rapid.IntRange(1, 255).Draw(t, "edgeR")),
			G: uint8(rapid.IntRange(1, 255).Draw(t, "edgeG")),
			B: uint8(rapid.IntRange(1, 255).Draw(t, "edgeB")),
			A: 255,
		}

		stops := []gradient.ColorStop{
			{Position: 0.0, Color: centerColor},
			{Position: 1.0, Color: edgeColor},
		}

		cfg := Config{
			Shape:      Circle,
			State:      On,
			Brightness: -1.0, // Discrete On → effectiveBrightness = 1.0
			Diameter:   diameter,
			Bounds:     image.Rect(0, 0, diameter, diameter),
			Foreground: centerColor, // Non-zero to avoid default resolution
			Gradient:   &GradientConfig{Stops: stops},
			// No border, no glow for simplicity
			BorderWidth: 0,
			GlowEnabled: false,
			ShineStyle:  ShineNone,
		}

		result := Render(cfg)
		if result == nil {
			t.Fatalf("expected non-nil result for Diameter=%d with gradient", diameter)
		}

		img := result.Image.(*image.RGBA)

		// The body area is the full image (no border, no glow).
		// For odd diameter, the geometric center is at (diameter/2.0, diameter/2.0).
		// The pixel at index (diameter/2, diameter/2) has its center at
		// (diameter/2 + 0.5, diameter/2 + 0.5). For odd diameter d, d/2 (int) = (d-1)/2,
		// so pixel center = (d-1)/2 + 0.5 = d/2.0 exactly. This pixel IS the center.
		cx := diameter / 2
		cy := diameter / 2

		// Verify the center pixel's normalized distance is effectively 0.
		// For odd diameter: geometric center = d/2.0, pixel at index d/2 has center
		// at d/2 + 0.5. E.g., d=9: center = 4.5, pixel[4] center = 4.5. Perfect.
		centerPx := img.RGBAAt(cx, cy)
		if !colorWithinTolerance(centerPx, centerColor, 1) {
			t.Fatalf("center pixel (%d,%d) = RGBA(%d,%d,%d,%d), expected ~RGBA(%d,%d,%d,%d) within ±1 [diameter=%d]",
				cx, cy, centerPx.R, centerPx.G, centerPx.B, centerPx.A,
				centerColor.R, centerColor.G, centerColor.B, centerColor.A, diameter)
		}

		// Check a pixel at the body perimeter.
		// Find the outermost filled pixel on the center row (from the right).
		radius := float64(diameter) / 2.0
		var edgePx color.RGBA
		var edgeX int
		foundEdge := false
		for x := diameter - 1; x >= cx; x-- {
			px := img.RGBAAt(x, cy)
			if px.A > 0 { // First non-transparent pixel from the right edge
				edgePx = px
				edgeX = x
				foundEdge = true
				break
			}
		}

		if !foundEdge {
			t.Fatalf("no non-transparent pixel found on center row from right edge [diameter=%d]", diameter)
		}

		// Compute the expected color at this pixel's normalized distance.
		pcx := float64(edgeX) + 0.5
		geomCenter := float64(diameter) / 2.0
		dist := math.Abs(pcx - geomCenter)
		normDist := dist / radius
		if normDist > 1.0 {
			normDist = 1.0
		}

		// Interpolate expected color at normDist between centerColor (0.0) and edgeColor (1.0).
		expectedR := math.Round(float64(centerColor.R) + (float64(edgeColor.R)-float64(centerColor.R))*normDist)
		expectedG := math.Round(float64(centerColor.G) + (float64(edgeColor.G)-float64(centerColor.G))*normDist)
		expectedB := math.Round(float64(centerColor.B) + (float64(edgeColor.B)-float64(centerColor.B))*normDist)
		expectedEdge := color.RGBA{R: uint8(expectedR), G: uint8(expectedG), B: uint8(expectedB), A: 255}

		if !colorWithinTolerance(edgePx, expectedEdge, 1) {
			t.Fatalf("edge pixel (%d,%d) normDist=%.3f = RGBA(%d,%d,%d,%d), expected ~RGBA(%d,%d,%d,%d) within ±1 [diameter=%d]",
				edgeX, cy, normDist, edgePx.R, edgePx.G, edgePx.B, edgePx.A,
				expectedEdge.R, expectedEdge.G, expectedEdge.B, expectedEdge.A, diameter)
		}
	})
}

// TestProperty6_GradientFallbackToSolidFill verifies that for any Config where
// GradientConfig is nil, has an empty Stops slice, or has fewer than 2 valid
// ColorStops (after discarding NaN/Inf positions), the rendered output SHALL be
// pixel-identical to the same Config with no gradient configured.

func TestProperty6_GradientFallbackToSolidFill(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		diameter := rapid.IntRange(3, 64).Draw(t, "diameter")
		shape := Shape(rapid.IntRange(0, 3).Draw(t, "shape"))

		fg := color.RGBA{
			R: uint8(rapid.IntRange(1, 255).Draw(t, "fgR")),
			G: uint8(rapid.IntRange(1, 255).Draw(t, "fgG")),
			B: uint8(rapid.IntRange(1, 255).Draw(t, "fgB")),
			A: 255,
		}

		bg := color.RGBA{
			R: uint8(rapid.IntRange(1, 255).Draw(t, "bgR")),
			G: uint8(rapid.IntRange(1, 255).Draw(t, "bgG")),
			B: uint8(rapid.IntRange(1, 255).Draw(t, "bgB")),
			A: 255,
		}

		bounds := image.Rect(0, 0, diameter, diameter)

		// Generate an invalid gradient config from several categories
		category := rapid.IntRange(0, 3).Draw(t, "category")
		var gradCfg *GradientConfig

		switch category {
		case 0:
			// nil gradient
			gradCfg = nil
		case 1:
			// Empty stops slice
			gradCfg = &GradientConfig{Stops: []gradient.ColorStop{}}
		case 2:
			// Only 1 valid stop
			gradCfg = &GradientConfig{Stops: []gradient.ColorStop{
				{Position: 0.5, Color: color.RGBA{R: 255, G: 0, B: 0, A: 255}},
			}}
		case 3:
			// Stops with NaN/Inf positions leaving < 2 valid
			gradCfg = &GradientConfig{Stops: []gradient.ColorStop{
				{Position: math.NaN(), Color: color.RGBA{R: 255, G: 0, B: 0, A: 255}},
				{Position: math.Inf(1), Color: color.RGBA{R: 0, G: 255, B: 0, A: 255}},
				{Position: 0.5, Color: color.RGBA{R: 0, G: 0, B: 255, A: 255}},
			}}
		}

		// Render with invalid gradient
		cfgWithGradient := Config{
			Shape:       shape,
			State:       On,
			Brightness:  -1.0,
			Diameter:    diameter,
			Bounds:      bounds,
			Foreground:  fg,
			Background:  bg,
			Gradient:    gradCfg,
			BorderWidth: 0,
			GlowEnabled: false,
			ShineStyle:  ShineNone,
		}
		resultWithGradient := Render(cfgWithGradient)

		// Render without gradient (nil)
		cfgNoGradient := Config{
			Shape:       shape,
			State:       On,
			Brightness:  -1.0,
			Diameter:    diameter,
			Bounds:      bounds,
			Foreground:  fg,
			Background:  bg,
			Gradient:    nil,
			BorderWidth: 0,
			GlowEnabled: false,
			ShineStyle:  ShineNone,
		}
		resultNoGradient := Render(cfgNoGradient)

		if resultWithGradient == nil {
			t.Fatal("expected non-nil result for valid diameter with invalid gradient")
		}
		if resultNoGradient == nil {
			t.Fatal("expected non-nil result for valid diameter without gradient")
		}

		// Compare pixel-for-pixel
		assertImagesIdentical(t, resultWithGradient.Image.(*image.RGBA), resultNoGradient.Image.(*image.RGBA),
			"gradient-fallback", "no-gradient", diameter, int(shape))
	})
}

// TestProperty7_GradientTruncationTo16Stops verifies that for any Config with a
// GradientConfig containing more than 16 ColorStops, the rendered output SHALL be
// pixel-identical to the same Config with only the first 16 stops (in original order).

func TestProperty7_GradientTruncationTo16Stops(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		diameter := rapid.IntRange(7, 64).Draw(t, "diameter")

		fg := color.RGBA{
			R: uint8(rapid.IntRange(1, 255).Draw(t, "fgR")),
			G: uint8(rapid.IntRange(1, 255).Draw(t, "fgG")),
			B: uint8(rapid.IntRange(1, 255).Draw(t, "fgB")),
			A: 255,
		}

		bounds := image.Rect(0, 0, diameter, diameter)

		// Generate more than 16 stops (17–20)
		numStops := rapid.IntRange(17, 20).Draw(t, "numStops")
		stops := make([]gradient.ColorStop, numStops)
		for i := 0; i < numStops; i++ {
			pos := float64(i) / float64(numStops-1) // Spread evenly in [0, 1]
			stops[i] = gradient.ColorStop{
				Position: pos,
				Color: color.RGBA{
					R: uint8(rapid.IntRange(0, 255).Draw(t, "stopR")),
					G: uint8(rapid.IntRange(0, 255).Draw(t, "stopG")),
					B: uint8(rapid.IntRange(0, 255).Draw(t, "stopB")),
					A: 255,
				},
			}
		}

		// Render with all stops (>16)
		cfgAll := Config{
			Shape:       Circle,
			State:       On,
			Brightness:  -1.0,
			Diameter:    diameter,
			Bounds:      bounds,
			Foreground:  fg,
			Gradient:    &GradientConfig{Stops: stops},
			BorderWidth: 0,
			GlowEnabled: false,
			ShineStyle:  ShineNone,
		}
		resultAll := Render(cfgAll)

		// Render with only the first 16 stops
		first16 := make([]gradient.ColorStop, 16)
		copy(first16, stops[:16])

		cfgTruncated := Config{
			Shape:       Circle,
			State:       On,
			Brightness:  -1.0,
			Diameter:    diameter,
			Bounds:      bounds,
			Foreground:  fg,
			Gradient:    &GradientConfig{Stops: first16},
			BorderWidth: 0,
			GlowEnabled: false,
			ShineStyle:  ShineNone,
		}
		resultTruncated := Render(cfgTruncated)

		if resultAll == nil {
			t.Fatal("expected non-nil result for valid config with >16 stops")
		}
		if resultTruncated == nil {
			t.Fatal("expected non-nil result for valid config with 16 stops")
		}

		// Compare pixel-for-pixel
		assertImagesIdentical(t, resultAll.Image.(*image.RGBA), resultTruncated.Image.(*image.RGBA),
			"all-stops", "first-16-stops", diameter, 0)
	})
}

// TestProperty8_GradientIgnoredWhenLEDIsOff verifies that for any Config with a
// valid GradientConfig and LED in Off state (or Brightness = 0.0), the rendered
// output SHALL be pixel-identical to the same Config without a gradient in Off state.

func TestProperty8_GradientIgnoredWhenLEDIsOff(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		diameter := rapid.IntRange(3, 64).Draw(t, "diameter")
		shape := Shape(rapid.IntRange(0, 3).Draw(t, "shape"))

		fg := color.RGBA{
			R: uint8(rapid.IntRange(1, 255).Draw(t, "fgR")),
			G: uint8(rapid.IntRange(1, 255).Draw(t, "fgG")),
			B: uint8(rapid.IntRange(1, 255).Draw(t, "fgB")),
			A: 255,
		}

		bg := color.RGBA{
			R: uint8(rapid.IntRange(1, 255).Draw(t, "bgR")),
			G: uint8(rapid.IntRange(1, 255).Draw(t, "bgG")),
			B: uint8(rapid.IntRange(1, 255).Draw(t, "bgB")),
			A: 255,
		}

		bounds := image.Rect(0, 0, diameter, diameter)

		// Generate a valid gradient (≥ 2 stops)
		numStops := rapid.IntRange(2, 8).Draw(t, "numStops")
		stops := make([]gradient.ColorStop, numStops)
		for i := 0; i < numStops; i++ {
			pos := float64(i) / float64(numStops-1)
			stops[i] = gradient.ColorStop{
				Position: pos,
				Color: color.RGBA{
					R: uint8(rapid.IntRange(0, 255).Draw(t, "stopR")),
					G: uint8(rapid.IntRange(0, 255).Draw(t, "stopG")),
					B: uint8(rapid.IntRange(0, 255).Draw(t, "stopB")),
					A: 255,
				},
			}
		}

		// Choose off-state representation: either discrete Off or Brightness = 0.0
		useDiscrete := rapid.Bool().Draw(t, "useDiscrete")
		var state State
		var brightness float64
		if useDiscrete {
			state = Off
			brightness = -1.0
		} else {
			state = On // State doesn't matter when Brightness != -1.0
			brightness = 0.0
		}

		// Render with gradient in Off state
		cfgWithGradient := Config{
			Shape:       shape,
			State:       state,
			Brightness:  brightness,
			Diameter:    diameter,
			Bounds:      bounds,
			Foreground:  fg,
			Background:  bg,
			Gradient:    &GradientConfig{Stops: stops},
			BorderWidth: 0,
			GlowEnabled: false,
			ShineStyle:  ShineNone,
		}
		resultWithGradient := Render(cfgWithGradient)

		// Render without gradient in Off state
		cfgNoGradient := Config{
			Shape:       shape,
			State:       state,
			Brightness:  brightness,
			Diameter:    diameter,
			Bounds:      bounds,
			Foreground:  fg,
			Background:  bg,
			Gradient:    nil,
			BorderWidth: 0,
			GlowEnabled: false,
			ShineStyle:  ShineNone,
		}
		resultNoGradient := Render(cfgNoGradient)

		if resultWithGradient == nil {
			t.Fatal("expected non-nil result for valid diameter in Off state with gradient")
		}
		if resultNoGradient == nil {
			t.Fatal("expected non-nil result for valid diameter in Off state without gradient")
		}

		// Compare pixel-for-pixel
		assertImagesIdentical(t, resultWithGradient.Image.(*image.RGBA), resultNoGradient.Image.(*image.RGBA),
			"off-with-gradient", "off-no-gradient", diameter, int(shape))
	})
}

// --- Helper functions ---

// colorWithinTolerance checks if two colors match within ±tolerance per channel.
func colorWithinTolerance(actual, expected color.RGBA, tolerance uint8) bool {
	return absDiffU8(actual.R, expected.R) <= tolerance &&
		absDiffU8(actual.G, expected.G) <= tolerance &&
		absDiffU8(actual.B, expected.B) <= tolerance &&
		absDiffU8(actual.A, expected.A) <= tolerance
}

// absDiffU8 returns the absolute difference between two uint8 values.
func absDiffU8(a, b uint8) uint8 {
	if a > b {
		return a - b
	}
	return b - a
}

// assertImagesIdentical compares two images pixel-by-pixel and fails the test
// if any pixel differs.
func assertImagesIdentical(t *rapid.T, img1, img2 *image.RGBA, name1, name2 string, diameter, shape int) {
	t.Helper()

	if img1.Bounds() != img2.Bounds() {
		t.Fatalf("image bounds differ: %s=%v, %s=%v [diameter=%d, shape=%d]",
			name1, img1.Bounds(), name2, img2.Bounds(), diameter, shape)
	}

	b := img1.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			c1 := img1.RGBAAt(x, y)
			c2 := img2.RGBAAt(x, y)
			if c1 != c2 {
				t.Fatalf("pixel mismatch at (%d,%d): %s=RGBA(%d,%d,%d,%d), %s=RGBA(%d,%d,%d,%d) [diameter=%d, shape=%d]",
					x, y, name1, c1.R, c1.G, c1.B, c1.A, name2, c2.R, c2.G, c2.B, c2.A, diameter, shape)
			}
		}
	}
}
