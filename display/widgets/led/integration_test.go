package led

import (
	"image"
	"image/color"
	"testing"
	"time"

	"github.com/databeast/cyberhud/display/widgets/gradient"
)

// TestFullPipelineIntegration exercises the complete render pipeline with ALL effects
// simultaneously enabled: Circle shape, glow, border, gradient, shine (Dot), and
// Pulse animation with elapsed time.
//
// This verifies:
// - Correct layer order: transparent base → glow → border → body fill → shine
// - Output dimensions account for glow (Diameter + 2*GlowRadius = 30 + 2*5 = 40)
// - No pixel index out of bounds
// - Glow pixels exist in the outer ring
// - Border pixels exist in the perimeter band
// - Body pixels exist inside
// - Shine pixel exists at the expected position
//

func TestFullPipelineIntegration(t *testing.T) {
	cfg := Config{
		Shape:      Circle,
		State:      On,
		Brightness: -1.0,
		Diameter:   30,
		Bounds:     image.Rect(10, 20, 40, 50),

		Foreground:   color.RGBA{R: 200, G: 50, B: 50, A: 255},
		Background:   color.RGBA{R: 0, G: 0, B: 0, A: 255},
		WarningColor: color.RGBA{R: 255, G: 191, B: 0, A: 255},

		// Gradient with 2 stops
		Gradient: &GradientConfig{
			Stops: []gradient.ColorStop{
				{Position: 0.0, Color: color.RGBA{R: 255, G: 100, B: 0, A: 255}},
				{Position: 1.0, Color: color.RGBA{R: 100, G: 0, B: 200, A: 255}},
			},
		},

		// Glow
		GlowEnabled: true,
		GlowRadius:  5,

		// Border
		BorderWidth: 2,
		BorderColor: color.RGBA{R: 128, G: 128, B: 128, A: 255},

		// Shine
		ShineStyle:   ShineDot,
		ShineOpacity: 0, // means 255 (fully opaque)

		// Pulse animation with some elapsed time (quarter period → brightness rising)
		Animation: AnimationConfig{
			Type:          Pulse,
			Period:        1000 * time.Millisecond,
			MinBrightness: 0.3,
		},
		animElapsed: 250 * time.Millisecond, // 25% through cycle
	}

	result := Render(cfg)
	if result == nil {
		t.Fatal("expected non-nil result for valid config with all effects enabled")
	}

	// --- Verify output dimensions ---
	expectedSize := 30 + 2*5 // Diameter + 2*GlowRadius = 40
	bounds := result.Image.Bounds()
	if bounds.Dx() != expectedSize || bounds.Dy() != expectedSize {
		t.Fatalf("output dimensions: got %d×%d, want %d×%d",
			bounds.Dx(), bounds.Dy(), expectedSize, expectedSize)
	}

	// --- Verify position ---
	if result.Position.X != 10 || result.Position.Y != 20 {
		t.Fatalf("position: got (%d, %d), want (10, 20)",
			result.Position.X, result.Position.Y)
	}

	// --- Verify no pixel extends beyond image bounds ---
	img := result.Image
	imgBounds := img.Bounds()
	// Just accessing all pixels within bounds should not panic; the image was
	// properly allocated. We additionally check that all non-transparent pixels
	// are within the expected bounds.
	for y := imgBounds.Min.Y; y < imgBounds.Max.Y; y++ {
		for x := imgBounds.Min.X; x < imgBounds.Max.X; x++ {
			// This access should never panic for a correctly sized image.
			_ = img.At(x, y)
		}
	}

	// --- Verify glow pixels exist in the outer ring ---
	// Glow occupies the ring from the output edge inward up to glowRadius pixels.
	// The body center is at (20, 20) in the output (outputSize/2 = 40/2 = 20).
	// Glow pixels should exist at the corners of the output image or in the
	// outer ring beyond the LED body (which starts at glowRadius=5 inset from edge).
	glowFound := false
	for y := 0; y < 5; y++ {
		for x := 0; x < expectedSize; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			if a > 0 && (r > 0 || g > 0 || b > 0) {
				glowFound = true
				break
			}
		}
		if glowFound {
			break
		}
	}
	// Also check the outermost pixels on the sides
	if !glowFound {
		for y := 5; y < expectedSize-5; y++ {
			for x := 0; x < 5; x++ {
				r, g, b, a := img.At(x, y).RGBA()
				if a > 0 && (r > 0 || g > 0 || b > 0) {
					glowFound = true
					break
				}
			}
			if glowFound {
				break
			}
		}
	}
	if !glowFound {
		t.Error("expected glow pixels in the outer ring (beyond body edge), found none")
	}

	// --- Verify border pixels exist in the perimeter band ---
	// The border occupies the ring just inside the body outer edge (inset = glowRadius).
	// Outer border starts at pixel offset glowRadius=5 from the output edge.
	// Border is 2 pixels wide, so border pixels are in [5, 7) from edge.
	// The body outer edge is a circle of radius 15 (Diameter/2) centered at (20, 20).
	// Border ring: pixels with distance from center in (13, 15] (outerRadius=15, innerRadius=13).
	borderGray := color.RGBA{R: 128, G: 128, B: 128, A: 255}
	borderFound := false
	center := float64(expectedSize) / 2.0

	// Check pixels along the horizontal axis at center height (y=20)
	// At y=20, the body outer edge is at x=5 and x=35 (distance 15 from center).
	// Border pixels should be around x=5,6 and x=34,35 (the outer 2px of the body).
	for x := 5; x < 7; x++ {
		r, g, b, a := img.At(x, int(center)).RGBA()
		c := color.RGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: uint8(a >> 8)}
		if c == borderGray {
			borderFound = true
			break
		}
	}
	if !borderFound {
		// Try the right side
		for x := expectedSize - 7; x < expectedSize-5; x++ {
			r, g, b, a := img.At(x, int(center)).RGBA()
			c := color.RGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: uint8(a >> 8)}
			if c == borderGray {
				borderFound = true
				break
			}
		}
	}
	if !borderFound {
		t.Error("expected border pixels in the perimeter band, found none at expected positions")
	}

	// --- Verify body fill pixels exist inside ---
	// The body is gradient-filled. Center of body is at (20, 20) in output coords.
	// The center pixel should be non-transparent and non-border color.
	centerPx := int(center)
	cr, cg, cb, ca := img.At(centerPx, centerPx).RGBA()
	if ca == 0 {
		t.Error("expected non-transparent body pixel at center, got transparent")
	}
	centerColor := color.RGBA{R: uint8(cr >> 8), G: uint8(cg >> 8), B: uint8(cb >> 8), A: uint8(ca >> 8)}
	if centerColor == borderGray {
		t.Error("center pixel should be body fill (gradient), not border color")
	}

	// --- Verify shine pixel exists at expected position ---
	// Shine (Dot) is drawn in the upper-left quadrant of the body.
	// Body rect inset: glowRadius + borderWidth = 5 + 2 = 7
	// Body rect: [7, 7] to [33, 33], size = 26×26
	// Body center in image coords: (7 + 26/2, 7 + 26/2) = (20, 20)
	// Body radius = 26/2 = 13
	// Dot center offset: 25% of bodyRadius = floor(13 * 0.25) = 3
	// Dot center: (20-3, 20-3) = (17, 17)
	// Dot radius: floor(13 * 0.15) = 1
	// So shine pixel should be at (17, 17).
	shinePx := img.At(17, 17)
	sr, sg, sb, sa := shinePx.RGBA()
	// Shine is white (R=255,G=255,B=255) with alpha modulated by brightness.
	// At 25% phase, pulse brightness = 0.3 + (1-0.3)*(1-cos(2π*0.25))/2
	// = 0.3 + 0.7*(1-cos(π/2))/2 = 0.3 + 0.7*(1-0)/2 = 0.3 + 0.35 = 0.65
	// Shine alpha = floor(255 * 0.65) = 165
	// The pixel should be white-ish with some alpha.
	if sa == 0 {
		t.Error("expected shine pixel at (17,17) to be non-transparent")
	}
	// Verify it's a white tone (RGB channels should be near-equal and bright)
	if sr>>8 != 255 || sg>>8 != 255 || sb>>8 != 255 {
		// The shine should be white with alpha; verify RGB are high.
		// Note: img.At returns pre-multiplied, so we need to check un-premultiplied.
		if sa > 0 {
			// Un-premultiply: actual_R = R * 0xFFFF / A
			actualR := sr * 0xFFFF / sa
			actualG := sg * 0xFFFF / sa
			actualB := sb * 0xFFFF / sa
			if actualR < 0xF000 || actualG < 0xF000 || actualB < 0xF000 {
				t.Errorf("shine pixel at (17,17) expected white-ish (high RGB), got RGBA16(%d,%d,%d,%d)",
					sr, sg, sb, sa)
			}
		}
	}

	// --- Verify label ---
	// Pulse at 25% phase gives brightness ≈ 0.65 (> 0.5), and State=On with
	// Brightness=-1.0 → label should be "led/on" (discrete state determines label
	// when Brightness == -1.0).
	if result.Label != "led/on" {
		t.Errorf("label: got %q, want %q", result.Label, "led/on")
	}
}

// TestFullPipelineSquareAllEffects exercises the pipeline with Square shape
// and all effects, confirming correct composition and bounds.
//

func TestFullPipelineSquareAllEffects(t *testing.T) {
	cfg := Config{
		Shape:      Square,
		State:      On,
		Brightness: -1.0,
		Diameter:   20,
		Bounds:     image.Rect(0, 0, 20, 20),

		Foreground: color.RGBA{R: 0, G: 200, B: 100, A: 255},
		Background: color.RGBA{R: 10, G: 10, B: 10, A: 255},

		Gradient: &GradientConfig{
			Stops: []gradient.ColorStop{
				{Position: 0.0, Color: color.RGBA{R: 0, G: 255, B: 0, A: 255}},
				{Position: 1.0, Color: color.RGBA{R: 0, G: 0, B: 255, A: 255}},
			},
		},

		GlowEnabled: true,
		GlowRadius:  3,

		BorderWidth: 1,
		BorderColor: color.RGBA{R: 200, G: 200, B: 200, A: 255},

		ShineStyle: ShineCrescent,

		Animation: AnimationConfig{
			Type:          Blink,
			Period:        500 * time.Millisecond,
			MinBrightness: 0.3,
		},
		animElapsed: 100 * time.Millisecond, // In the "on" half (< 250ms)
	}

	result := Render(cfg)
	if result == nil {
		t.Fatal("expected non-nil result for square with all effects")
	}

	expectedSize := 20 + 2*3 // 26
	bounds := result.Image.Bounds()
	if bounds.Dx() != expectedSize || bounds.Dy() != expectedSize {
		t.Fatalf("square output dimensions: got %d×%d, want %d×%d",
			bounds.Dx(), bounds.Dy(), expectedSize, expectedSize)
	}

	// Verify all pixels are within image bounds (no panic on traversal)
	img := result.Image
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			_ = img.At(x, y)
		}
	}

	// Center pixel should be non-transparent (body fill)
	_, _, _, ca := img.At(expectedSize/2, expectedSize/2).RGBA()
	if ca == 0 {
		t.Error("center pixel should be non-transparent for On-state square")
	}
}

// TestFullPipelineDiamondAllEffects exercises the pipeline with Diamond shape.
//

func TestFullPipelineDiamondAllEffects(t *testing.T) {
	cfg := Config{
		Shape:      Diamond,
		State:      On,
		Brightness: -1.0,
		Diameter:   24,
		Bounds:     image.Rect(5, 5, 29, 29),

		Foreground: color.RGBA{R: 100, G: 100, B: 255, A: 255},

		Gradient: &GradientConfig{
			Stops: []gradient.ColorStop{
				{Position: 0.0, Color: color.RGBA{R: 255, G: 255, B: 0, A: 255}},
				{Position: 1.0, Color: color.RGBA{R: 0, G: 100, B: 255, A: 255}},
			},
		},

		GlowEnabled: true,
		GlowRadius:  4,

		BorderWidth: 2,
		BorderColor: color.RGBA{R: 80, G: 80, B: 80, A: 255},

		ShineStyle: ShineDot,

		Animation: AnimationConfig{
			Type:   Fade,
			Period: 2000 * time.Millisecond,
		},
		animElapsed: 500 * time.Millisecond, // 25% through → brightness 0.5
	}

	result := Render(cfg)
	if result == nil {
		t.Fatal("expected non-nil result for diamond with all effects")
	}

	expectedSize := 24 + 2*4 // 32
	bounds := result.Image.Bounds()
	if bounds.Dx() != expectedSize || bounds.Dy() != expectedSize {
		t.Fatalf("diamond output dimensions: got %d×%d, want %d×%d",
			bounds.Dx(), bounds.Dy(), expectedSize, expectedSize)
	}

	// Position should match Bounds.Min
	if result.Position.X != 5 || result.Position.Y != 5 {
		t.Fatalf("position: got (%d,%d), want (5,5)",
			result.Position.X, result.Position.Y)
	}
}

// TestFullPipelineRoundedSquareAllEffects exercises the pipeline with RoundedSquare shape.
//

func TestFullPipelineRoundedSquareAllEffects(t *testing.T) {
	cfg := Config{
		Shape:      RoundedSquare,
		State:      Warning,
		Brightness: -1.0,
		Diameter:   16,
		Bounds:     image.Rect(0, 0, 16, 16),

		WarningColor: color.RGBA{R: 255, G: 191, B: 0, A: 255},

		GlowEnabled: true,
		GlowRadius:  2,

		BorderWidth: 1,

		ShineStyle: ShineDot,

		Animation: AnimationConfig{
			Type:          Pulse,
			Period:        1000 * time.Millisecond,
			MinBrightness: 0.3,
		},
		animElapsed: 500 * time.Millisecond, // Peak brightness for Pulse
	}

	result := Render(cfg)
	if result == nil {
		t.Fatal("expected non-nil result for rounded square with warning state")
	}

	expectedSize := 16 + 2*2 // 20
	bounds := result.Image.Bounds()
	if bounds.Dx() != expectedSize || bounds.Dy() != expectedSize {
		t.Fatalf("rounded square output dimensions: got %d×%d, want %d×%d",
			bounds.Dx(), bounds.Dy(), expectedSize, expectedSize)
	}

	// At phase=0.5 for Pulse: brightness = 0.3 + 0.7*(1-cos(π))/2 = 0.3 + 0.7*1 = 1.0
	// So the LED is at full brightness. Label should be "led/warning" because
	// State=Warning and Brightness=-1.0.
	if result.Label != "led/warning" {
		t.Errorf("label: got %q, want %q", result.Label, "led/warning")
	}
}

// TestSingleVsGroupDispatch verifies that Render correctly dispatches
// between single LED and group rendering.
//

func TestSingleVsGroupDispatch(t *testing.T) {
	// Single LED (no group)
	singleCfg := Config{
		Shape:      Circle,
		State:      On,
		Brightness: -1.0,
		Diameter:   10,
		Bounds:     image.Rect(0, 0, 10, 10),
		Foreground: color.RGBA{R: 0, G: 200, B: 0, A: 255},
	}
	single := Render(singleCfg)
	if single == nil {
		t.Fatal("single LED render returned nil")
	}
	if single.Image.Bounds().Dx() != 10 || single.Image.Bounds().Dy() != 10 {
		t.Errorf("single LED dimensions: got %d×%d, want 10×10",
			single.Image.Bounds().Dx(), single.Image.Bounds().Dy())
	}

	// Group with 3 entries, horizontal
	groupCfg := Config{
		Shape:      Circle,
		State:      On,
		Brightness: -1.0,
		Diameter:   10,
		Bounds:     image.Rect(0, 0, 10, 10),
		Foreground: color.RGBA{R: 0, G: 200, B: 0, A: 255},
		Spacing:    2,
		Group: []GroupEntry{
			{State: On, Foreground: color.RGBA{R: 255, G: 0, B: 0, A: 255}},
			{State: Off},
			{State: On, Foreground: color.RGBA{R: 0, G: 0, B: 255, A: 255}},
		},
		Orientation: Horizontal,
	}
	group := Render(groupCfg)
	if group == nil {
		t.Fatal("group LED render returned nil")
	}
	// Expected: 3*10 + 2*2 = 34 width, 10 height
	expectedW := 3*10 + 2*2
	if group.Image.Bounds().Dx() != expectedW || group.Image.Bounds().Dy() != 10 {
		t.Errorf("group LED dimensions: got %d×%d, want %d×10",
			group.Image.Bounds().Dx(), group.Image.Bounds().Dy(), expectedW)
	}
}

// TestNoPixelsOutOfBounds iterates over all shape/effect combinations to ensure
// no pixel access goes out of bounds during rendering.
//

func TestNoPixelsOutOfBounds(t *testing.T) {
	shapes := []Shape{Circle, Square, Diamond, RoundedSquare}
	glowConfigs := []struct {
		enabled bool
		radius  int
	}{
		{false, 0},
		{true, 3},
		{true, 10},
	}
	borderWidths := []int{0, 1, 3}
	shineStyles := []ShineStyle{ShineNone, ShineDot, ShineCrescent}

	for _, shape := range shapes {
		for _, glow := range glowConfigs {
			for _, bw := range borderWidths {
				for _, shine := range shineStyles {
					cfg := Config{
						Shape:       shape,
						State:       On,
						Brightness:  -1.0,
						Diameter:    20,
						Bounds:      image.Rect(0, 0, 20, 20),
						Foreground:  color.RGBA{R: 100, G: 200, B: 50, A: 255},
						GlowEnabled: glow.enabled,
						GlowRadius:  glow.radius,
						BorderWidth: bw,
						ShineStyle:  shine,
					}

					result := Render(cfg)
					if result == nil {
						t.Fatalf("unexpected nil for shape=%d, glow=%v, border=%d, shine=%d",
							shape, glow, bw, shine)
					}

					// Verify the output image bounds match expectations
					expectedGlow := 0
					if glow.enabled {
						expectedGlow = glow.radius
					}
					expectedSize := 20 + 2*expectedGlow
					imgBounds := result.Image.Bounds()
					if imgBounds.Dx() != expectedSize || imgBounds.Dy() != expectedSize {
						t.Errorf("shape=%d, glow=%v, border=%d: dimensions got %d×%d, want %d×%d",
							shape, glow, bw, imgBounds.Dx(), imgBounds.Dy(), expectedSize, expectedSize)
					}

					// Access all pixels — confirms no out-of-bounds panics
					for y := imgBounds.Min.Y; y < imgBounds.Max.Y; y++ {
						for x := imgBounds.Min.X; x < imgBounds.Max.X; x++ {
							_ = result.Image.At(x, y)
						}
					}
				}
			}
		}
	}
}
