package gradient_test

import (
	"bytes"
	"image"
	"image/color"
	"math"
	"reflect"
	"sort"
	"testing"

	"github.com/databeast/cyberhud/display/widgets/gradient"
	"pgregory.net/rapid"
)

// --- From: gradient_bench_test.go ---

// BenchmarkRenderLinear measures allocation count for a 100×100, 3-stop linear gradient.
// Target: 1 alloc/op (the image.NewRGBA buffer).

func BenchmarkRenderLinear(b *testing.B) {
	cfg := gradient.Config{
		Style:  gradient.Linear,
		Angle:  45,
		Bounds: image.Rect(0, 0, 100, 100),
		Stops: []gradient.ColorStop{
			{Position: 0.0, Color: color.RGBA{255, 0, 0, 255}},
			{Position: 0.5, Color: color.RGBA{0, 255, 0, 255}},
			{Position: 1.0, Color: color.RGBA{0, 0, 255, 255}},
		},
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gradient.Render(cfg)
	}
}

// BenchmarkRenderRadial measures allocation count for a 100×100, 3-stop radial gradient.
// Target: 1 alloc/op (the image.NewRGBA buffer).

func BenchmarkRenderRadial(b *testing.B) {
	cfg := gradient.Config{
		Style:  gradient.Radial,
		Angle:  0,
		Bounds: image.Rect(0, 0, 100, 100),
		Stops: []gradient.ColorStop{
			{Position: 0.0, Color: color.RGBA{255, 0, 0, 255}},
			{Position: 0.5, Color: color.RGBA{0, 255, 0, 255}},
			{Position: 1.0, Color: color.RGBA{0, 0, 255, 255}},
		},
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gradient.Render(cfg)
	}
}

// --- From: gradient_prop_test.go ---

// For any valid Config (style ∈ {Linear, Radial}, ≥ 2 stops, bounds with Dx() ≥ 1
// and Dy() ≥ 1, no NaN/Inf), the rendered image width SHALL equal Bounds.Dx() and
// height SHALL equal Bounds.Dy().

func TestPropertyOutputDimensionsMatchBounds(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate random style: Linear (0) or Radial (1)
		style := gradient.Style(rapid.IntRange(0, 1).Draw(t, "style"))

		// Generate random bounds with Dx() >= 1 and Dy() >= 1
		minX := rapid.IntRange(0, 100).Draw(t, "minX")
		minY := rapid.IntRange(0, 100).Draw(t, "minY")
		width := rapid.IntRange(1, 200).Draw(t, "width")
		height := rapid.IntRange(1, 200).Draw(t, "height")
		bounds := image.Rectangle{
			Min: image.Point{X: minX, Y: minY},
			Max: image.Point{X: minX + width, Y: minY + height},
		}

		// Generate a valid angle (finite, no NaN/Inf)
		angle := rapid.Float64Range(-720.0, 720.0).Draw(t, "angle")

		// Generate 2–10 random color stops with valid positions in [0, 1]
		numStops := rapid.IntRange(2, 10).Draw(t, "numStops")
		stops := make([]gradient.ColorStop, numStops)
		for i := range stops {
			stops[i] = gradient.ColorStop{
				Position: rapid.Float64Range(0.0, 1.0).Draw(t, "stopPos"),
				Color: color.RGBA{
					R: uint8(rapid.IntRange(0, 255).Draw(t, "stopR")),
					G: uint8(rapid.IntRange(0, 255).Draw(t, "stopG")),
					B: uint8(rapid.IntRange(0, 255).Draw(t, "stopB")),
					A: uint8(rapid.IntRange(0, 255).Draw(t, "stopA")),
				},
			}
		}

		cfg := gradient.Config{
			Style:  style,
			Angle:  angle,
			Bounds: bounds,
			Stops:  stops,
		}

		result := gradient.Render(cfg)

		// Result must be non-nil for a valid config
		if result == nil {
			t.Fatalf("Render returned nil for valid config: %+v", cfg)
		}

		// Check output image dimensions match bounds
		img := result.Image
		if img == nil {
			t.Fatalf("Result.Image is nil for valid config: %+v", cfg)
		}

		imgBounds := img.Bounds()
		gotWidth := imgBounds.Dx()
		gotHeight := imgBounds.Dy()

		if gotWidth != width {
			t.Fatalf("image width = %d, want %d (bounds: %v)", gotWidth, width, bounds)
		}
		if gotHeight != height {
			t.Fatalf("image height = %d, want %d (bounds: %v)", gotHeight, height, bounds)
		}
	})
}

// For any valid Config with exactly two color stops at positions 0.0 and 1.0
// (style Linear, angle 0, height ≥ 2), the pixel at the vertical midpoint row
// SHALL have each RGBA channel within ±1 of the arithmetic mean of the two stop
// colors' corresponding channels.

func TestPropertyLinearInterpolationMidpoint(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate random bounds with odd height >= 5 so that height/2 maps
		// to exactly t=0.5 on the gradient axis. For odd h, midY = h/2 = (h-1)/2,
		// and t = midY/(h-1) = 0.5 exactly.
		halfH := rapid.IntRange(2, 50).Draw(t, "halfH")
		height := halfH*2 + 1 // odd, range [5, 101]
		width := rapid.IntRange(1, 100).Draw(t, "width")

		// Generate two random colors for stops at 0.0 and 1.0
		c1 := color.RGBA{
			R: uint8(rapid.IntRange(0, 255).Draw(t, "c1R")),
			G: uint8(rapid.IntRange(0, 255).Draw(t, "c1G")),
			B: uint8(rapid.IntRange(0, 255).Draw(t, "c1B")),
			A: uint8(rapid.IntRange(0, 255).Draw(t, "c1A")),
		}
		c2 := color.RGBA{
			R: uint8(rapid.IntRange(0, 255).Draw(t, "c2R")),
			G: uint8(rapid.IntRange(0, 255).Draw(t, "c2G")),
			B: uint8(rapid.IntRange(0, 255).Draw(t, "c2B")),
			A: uint8(rapid.IntRange(0, 255).Draw(t, "c2A")),
		}

		cfg := gradient.Config{
			Style:  gradient.Linear,
			Angle:  0, // top-to-bottom
			Bounds: image.Rect(0, 0, width, height),
			Stops: []gradient.ColorStop{
				{Position: 0.0, Color: c1},
				{Position: 1.0, Color: c2},
			},
		}

		result := gradient.Render(cfg)
		if result == nil {
			t.Fatal("Render returned nil for valid config")
		}

		img := result.Image

		// Check the pixel at (0, height/2) — the exact midpoint row for odd heights
		midY := height / 2
		r, g, b, a := img.At(0, midY).RGBA()
		// img.At returns pre-multiplied 16-bit values; convert to 8-bit
		pixR := uint8(r >> 8)
		pixG := uint8(g >> 8)
		pixB := uint8(b >> 8)
		pixA := uint8(a >> 8)

		// Expected: arithmetic mean of corresponding channels
		expectedR := (float64(c1.R) + float64(c2.R)) / 2.0
		expectedG := (float64(c1.G) + float64(c2.G)) / 2.0
		expectedB := (float64(c1.B) + float64(c2.B)) / 2.0
		expectedA := (float64(c1.A) + float64(c2.A)) / 2.0

		// Verify each channel is within ±1 of the arithmetic mean
		// The ±1 tolerance accounts for rounding in interpolation
		diff := func(got uint8, exp float64) float64 {
			if float64(got) > exp {
				return float64(got) - exp
			}
			return exp - float64(got)
		}

		if diff(pixR, expectedR) > 1 {
			t.Fatalf("R channel: pixel=%d, expected≈%.1f (c1.R=%d, c2.R=%d, midY=%d, height=%d)",
				pixR, expectedR, c1.R, c2.R, midY, height)
		}
		if diff(pixG, expectedG) > 1 {
			t.Fatalf("G channel: pixel=%d, expected≈%.1f (c1.G=%d, c2.G=%d, midY=%d, height=%d)",
				pixG, expectedG, c1.G, c2.G, midY, height)
		}
		if diff(pixB, expectedB) > 1 {
			t.Fatalf("B channel: pixel=%d, expected≈%.1f (c1.B=%d, c2.B=%d, midY=%d, height=%d)",
				pixB, expectedB, c1.B, c2.B, midY, height)
		}
		if diff(pixA, expectedA) > 1 {
			t.Fatalf("A channel: pixel=%d, expected≈%.1f (c1.A=%d, c2.A=%d, midY=%d, height=%d)",
				pixA, expectedA, c1.A, c2.A, midY, height)
		}
	})
}

// For any valid Config that produces a non-nil Result, Result.Position SHALL equal
// Bounds.Min, and Result.Label SHALL equal "gradient/linear" for Linear style or
// "gradient/radial" for Radial style.

func TestPropertyOutputMetadataCorrectness(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate random style: Linear (0) or Radial (1)
		style := gradient.Style(rapid.IntRange(0, 1).Draw(t, "style"))

		// Generate random bounds with non-zero Min and Dx() >= 1, Dy() >= 1
		minX := rapid.IntRange(1, 500).Draw(t, "minX")
		minY := rapid.IntRange(1, 500).Draw(t, "minY")
		width := rapid.IntRange(1, 200).Draw(t, "width")
		height := rapid.IntRange(1, 200).Draw(t, "height")
		bounds := image.Rectangle{
			Min: image.Point{X: minX, Y: minY},
			Max: image.Point{X: minX + width, Y: minY + height},
		}

		// Generate a valid angle
		angle := rapid.Float64Range(-720.0, 720.0).Draw(t, "angle")

		// Generate 2–5 random color stops with valid positions
		numStops := rapid.IntRange(2, 5).Draw(t, "numStops")
		stops := make([]gradient.ColorStop, numStops)
		for i := range stops {
			stops[i] = gradient.ColorStop{
				Position: rapid.Float64Range(0.0, 1.0).Draw(t, "stopPos"),
				Color: color.RGBA{
					R: uint8(rapid.IntRange(0, 255).Draw(t, "stopR")),
					G: uint8(rapid.IntRange(0, 255).Draw(t, "stopG")),
					B: uint8(rapid.IntRange(0, 255).Draw(t, "stopB")),
					A: uint8(rapid.IntRange(0, 255).Draw(t, "stopA")),
				},
			}
		}

		cfg := gradient.Config{
			Style:  style,
			Angle:  angle,
			Bounds: bounds,
			Stops:  stops,
		}

		result := gradient.Render(cfg)

		// Result must be non-nil for a valid config
		if result == nil {
			t.Fatalf("Render returned nil for valid config: %+v", cfg)
		}

		// Check Result.Position == Bounds.Min
		if result.Position != cfg.Bounds.Min {
			t.Fatalf("Result.Position = %v, want %v (Bounds.Min)", result.Position, cfg.Bounds.Min)
		}

		// Check Result.Label based on style
		var expectedLabel string
		if style == gradient.Linear {
			expectedLabel = "gradient/linear"
		} else {
			expectedLabel = "gradient/radial"
		}

		if result.Label != expectedLabel {
			t.Fatalf("Result.Label = %q, want %q (style=%d)", result.Label, expectedLabel, style)
		}
	})
}

// For any valid Config, calling Render twice with the same Config SHALL produce
// pixel-identical Results, and the Config struct SHALL not be mutated by the call.

func TestPropertyRenderPurity(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate random style: Linear (0) or Radial (1)
		style := gradient.Style(rapid.IntRange(0, 1).Draw(t, "style"))

		// Generate small bounds (1-50)
		minX := rapid.IntRange(0, 50).Draw(t, "minX")
		minY := rapid.IntRange(0, 50).Draw(t, "minY")
		width := rapid.IntRange(1, 50).Draw(t, "width")
		height := rapid.IntRange(1, 50).Draw(t, "height")
		bounds := image.Rectangle{
			Min: image.Point{X: minX, Y: minY},
			Max: image.Point{X: minX + width, Y: minY + height},
		}

		// Generate valid angle
		angle := rapid.Float64Range(-720.0, 720.0).Draw(t, "angle")

		// Generate 2-5 random color stops with valid positions
		numStops := rapid.IntRange(2, 5).Draw(t, "numStops")
		stops := make([]gradient.ColorStop, numStops)
		for i := range stops {
			stops[i] = gradient.ColorStop{
				Position: rapid.Float64Range(0.0, 1.0).Draw(t, "stopPos"),
				Color: color.RGBA{
					R: uint8(rapid.IntRange(0, 255).Draw(t, "stopR")),
					G: uint8(rapid.IntRange(0, 255).Draw(t, "stopG")),
					B: uint8(rapid.IntRange(0, 255).Draw(t, "stopB")),
					A: uint8(rapid.IntRange(0, 255).Draw(t, "stopA")),
				},
			}
		}

		cfg := gradient.Config{
			Style:  style,
			Angle:  angle,
			Bounds: bounds,
			Stops:  stops,
		}

		// Deep-copy the Config before calling Render to verify no mutation
		cfgCopy := gradient.Config{
			Style:  cfg.Style,
			Angle:  cfg.Angle,
			Bounds: cfg.Bounds,
			Stops:  make([]gradient.ColorStop, len(cfg.Stops)),
		}
		copy(cfgCopy.Stops, cfg.Stops)

		// Call Render twice with the same config
		result1 := gradient.Render(cfg)
		result2 := gradient.Render(cfg)

		// Verify both results are non-nil
		if result1 == nil {
			t.Fatal("First Render returned nil for valid config")
		}
		if result2 == nil {
			t.Fatal("Second Render returned nil for valid config")
		}

		// Compare pixel data of both rendered images
		img1, ok1 := result1.Image.(*image.RGBA)
		if !ok1 {
			t.Fatal("First result Image is not *image.RGBA")
		}
		img2, ok2 := result2.Image.(*image.RGBA)
		if !ok2 {
			t.Fatal("Second result Image is not *image.RGBA")
		}

		if !bytes.Equal(img1.Pix, img2.Pix) {
			t.Fatal("Render produced different pixel data on two calls with same Config")
		}

		// Verify Config was not mutated by Render calls
		if !reflect.DeepEqual(cfg, cfgCopy) {
			t.Fatalf("Config was mutated by Render: got %+v, original %+v", cfg, cfgCopy)
		}
	})
}

// For any Config where bounds have width < 1 or height < 1, OR fewer than 2 color
// stops are provided, OR style is not Linear or Radial, OR any stop position is
// NaN/Inf, OR angle is NaN/Inf, Render SHALL return nil.

func TestPropertyInvalidInputProducesNil(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Pick which invalid condition to inject (0-6)
		invalidKind := rapid.IntRange(0, 6).Draw(t, "invalidKind")

		// Start with an otherwise-valid config
		minX := rapid.IntRange(0, 50).Draw(t, "minX")
		minY := rapid.IntRange(0, 50).Draw(t, "minY")
		width := rapid.IntRange(1, 100).Draw(t, "width")
		height := rapid.IntRange(1, 100).Draw(t, "height")
		bounds := image.Rectangle{
			Min: image.Point{X: minX, Y: minY},
			Max: image.Point{X: minX + width, Y: minY + height},
		}
		angle := rapid.Float64Range(-720.0, 720.0).Draw(t, "angle")
		style := gradient.Style(rapid.IntRange(0, 1).Draw(t, "style"))

		numStops := rapid.IntRange(2, 10).Draw(t, "numStops")
		stops := make([]gradient.ColorStop, numStops)
		for i := range stops {
			stops[i] = gradient.ColorStop{
				Position: rapid.Float64Range(0.0, 1.0).Draw(t, "stopPos"),
				Color: color.RGBA{
					R: uint8(rapid.IntRange(0, 255).Draw(t, "stopR")),
					G: uint8(rapid.IntRange(0, 255).Draw(t, "stopG")),
					B: uint8(rapid.IntRange(0, 255).Draw(t, "stopB")),
					A: uint8(rapid.IntRange(0, 255).Draw(t, "stopA")),
				},
			}
		}

		// Now inject the chosen invalid condition
		switch invalidKind {
		case 0:
			// Invalid bounds: width < 1
			bounds.Max.X = bounds.Min.X + rapid.IntRange(-10, 0).Draw(t, "badWidth")
		case 1:
			// Invalid bounds: height < 1
			bounds.Max.Y = bounds.Min.Y + rapid.IntRange(-10, 0).Draw(t, "badHeight")
		case 2:
			// Fewer than 2 stops
			stopCount := rapid.IntRange(0, 1).Draw(t, "fewStops")
			stops = stops[:stopCount]
		case 3:
			// Invalid style (not Linear=0 or Radial=1)
			badStyle := rapid.IntRange(2, 100).Draw(t, "badStylePos")
			// Also allow negative style values
			if rapid.IntRange(0, 1).Draw(t, "negStyle") == 1 {
				badStyle = -rapid.IntRange(1, 100).Draw(t, "badStyleNeg")
			}
			style = gradient.Style(badStyle)
		case 4:
			// NaN stop position
			nanIdx := rapid.IntRange(0, len(stops)-1).Draw(t, "nanIdx")
			stops[nanIdx].Position = math.NaN()
		case 5:
			// Inf stop position
			infIdx := rapid.IntRange(0, len(stops)-1).Draw(t, "infIdx")
			sign := rapid.IntRange(0, 1).Draw(t, "infSign")
			if sign == 0 {
				stops[infIdx].Position = math.Inf(1)
			} else {
				stops[infIdx].Position = math.Inf(-1)
			}
		case 6:
			// NaN or Inf angle
			if rapid.IntRange(0, 1).Draw(t, "angleNanOrInf") == 0 {
				angle = math.NaN()
			} else {
				if rapid.IntRange(0, 1).Draw(t, "angleInfSign") == 0 {
					angle = math.Inf(1)
				} else {
					angle = math.Inf(-1)
				}
			}
		}

		cfg := gradient.Config{
			Style:  style,
			Angle:  angle,
			Bounds: bounds,
			Stops:  stops,
		}

		result := gradient.Render(cfg)
		if result != nil {
			t.Fatalf("Render returned non-nil for invalid config (invalidKind=%d): %+v", invalidKind, cfg)
		}
	})
}

// For any valid radial Config with ≥ 2 stops and bounds where min(Dx, Dy) ≥ 3,
// all pixels whose Euclidean distance from the center exceeds
// min(Bounds.Dx(), Bounds.Dy()) / 2 SHALL have the color of the last sorted stop.

func TestPropertyPixelsOutsideInscribedCircleGetLastStopColor(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate non-square bounds where one dimension is significantly larger
		// so there are guaranteed pixels outside the inscribed circle.
		orientation := rapid.IntRange(0, 1).Draw(t, "orientation")
		var width, height int
		if orientation == 0 {
			// Wide rectangle: width >> height
			width = rapid.IntRange(20, 50).Draw(t, "width")
			height = rapid.IntRange(3, 10).Draw(t, "height")
		} else {
			// Tall rectangle: height >> width
			width = rapid.IntRange(3, 10).Draw(t, "width")
			height = rapid.IntRange(20, 50).Draw(t, "height")
		}

		bounds := image.Rect(0, 0, width, height)

		// Generate 2-4 stops; ensure position 0.0 and 1.0 are present
		// with a distinct last stop color at position 1.0
		numMiddleStops := rapid.IntRange(0, 2).Draw(t, "numMiddleStops")

		firstColor := color.RGBA{
			R: uint8(rapid.IntRange(0, 255).Draw(t, "firstR")),
			G: uint8(rapid.IntRange(0, 255).Draw(t, "firstG")),
			B: uint8(rapid.IntRange(0, 255).Draw(t, "firstB")),
			A: uint8(rapid.IntRange(0, 255).Draw(t, "firstA")),
		}

		lastColor := color.RGBA{
			R: uint8(rapid.IntRange(0, 255).Draw(t, "lastR")),
			G: uint8(rapid.IntRange(0, 255).Draw(t, "lastG")),
			B: uint8(rapid.IntRange(0, 255).Draw(t, "lastB")),
			A: uint8(rapid.IntRange(0, 255).Draw(t, "lastA")),
		}

		stops := make([]gradient.ColorStop, 0, 2+numMiddleStops)
		stops = append(stops, gradient.ColorStop{Position: 0.0, Color: firstColor})

		for i := 0; i < numMiddleStops; i++ {
			pos := rapid.Float64Range(0.01, 0.99).Draw(t, "midPos")
			c := color.RGBA{
				R: uint8(rapid.IntRange(0, 255).Draw(t, "midR")),
				G: uint8(rapid.IntRange(0, 255).Draw(t, "midG")),
				B: uint8(rapid.IntRange(0, 255).Draw(t, "midB")),
				A: uint8(rapid.IntRange(0, 255).Draw(t, "midA")),
			}
			stops = append(stops, gradient.ColorStop{Position: pos, Color: c})
		}

		stops = append(stops, gradient.ColorStop{Position: 1.0, Color: lastColor})

		cfg := gradient.Config{
			Style:  gradient.Radial,
			Angle:  0,
			Bounds: bounds,
			Stops:  stops,
		}

		result := gradient.Render(cfg)
		if result == nil {
			t.Fatal("Render returned nil for valid radial config")
		}

		img := result.Image

		// Compute center and radius using integer division (matching implementation)
		cx := float64(width / 2)
		cy := float64(height / 2)
		r := min(width, height) / 2

		// Iterate all pixels and check those outside the inscribed circle
		outsideCount := 0
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				dx := float64(x) - cx
				dy := float64(y) - cy
				dist := math.Sqrt(dx*dx + dy*dy)

				if dist > float64(r) {
					outsideCount++
					pr, pg, pb, pa := img.At(x, y).RGBA()
					// Convert from pre-multiplied 16-bit to 8-bit
					gotR := uint8(pr >> 8)
					gotG := uint8(pg >> 8)
					gotB := uint8(pb >> 8)
					gotA := uint8(pa >> 8)

					if gotR != lastColor.R || gotG != lastColor.G || gotB != lastColor.B || gotA != lastColor.A {
						t.Fatalf("pixel (%d,%d) outside inscribed circle (dist=%.2f, r=%d) "+
							"has color RGBA(%d,%d,%d,%d), want RGBA(%d,%d,%d,%d)",
							x, y, dist, r,
							gotR, gotG, gotB, gotA,
							lastColor.R, lastColor.G, lastColor.B, lastColor.A)
					}
				}
			}
		}

		// Sanity check: with non-square bounds, there must be pixels outside the circle
		if outsideCount == 0 {
			t.Fatal("no pixels found outside inscribed circle — test setup is invalid")
		}
	})
}

// For any valid Config, rendering with color stops in arbitrary order SHALL produce
// pixel-identical output to rendering with the same stops pre-sorted by position in
// ascending order (preserving original order among same-position stops).

func TestPropertyStopSortingEquivalence(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate random style: Linear (0) or Radial (1)
		style := gradient.Style(rapid.IntRange(0, 1).Draw(t, "style"))

		// Generate small bounds (5-20) for efficiency
		width := rapid.IntRange(5, 20).Draw(t, "width")
		height := rapid.IntRange(5, 20).Draw(t, "height")
		bounds := image.Rect(0, 0, width, height)

		// Generate a valid angle
		angle := rapid.Float64Range(-720.0, 720.0).Draw(t, "angle")

		// Generate 2-10 color stops at distinct random positions [0, 1]
		// Positions must be distinct to avoid same-position tie-breaking ambiguity,
		// which is covered by Property 9 (duplicate-position stop uses last in input order).
		numStops := rapid.IntRange(2, 10).Draw(t, "numStops")
		stops := make([]gradient.ColorStop, numStops)
		// Generate distinct positions by using i-based offsets within sub-ranges
		for i := range stops {
			// Each stop gets a position in a unique sub-range [i/numStops, (i+1)/numStops)
			base := float64(i) / float64(numStops)
			offset := rapid.Float64Range(0.0, 0.9/float64(numStops)).Draw(t, "stopOffset")
			stops[i] = gradient.ColorStop{
				Position: base + offset,
				Color: color.RGBA{
					R: uint8(rapid.IntRange(0, 255).Draw(t, "stopR")),
					G: uint8(rapid.IntRange(0, 255).Draw(t, "stopG")),
					B: uint8(rapid.IntRange(0, 255).Draw(t, "stopB")),
					A: uint8(rapid.IntRange(0, 255).Draw(t, "stopA")),
				},
			}
		}

		// Create a shuffled (random permutation) copy of the stops
		shuffled := make([]gradient.ColorStop, numStops)
		copy(shuffled, stops)
		// Fisher-Yates shuffle using rapid for reproducibility
		for i := numStops - 1; i > 0; i-- {
			j := rapid.IntRange(0, i).Draw(t, "shuffleIdx")
			shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
		}

		// Create a stable-sorted copy of the stops by position (ascending)
		sorted := make([]gradient.ColorStop, numStops)
		copy(sorted, stops)
		sort.SliceStable(sorted, func(i, j int) bool {
			return sorted[i].Position < sorted[j].Position
		})

		// Render with the shuffled stops
		cfgShuffled := gradient.Config{
			Style:  style,
			Angle:  angle,
			Bounds: bounds,
			Stops:  shuffled,
		}
		resultShuffled := gradient.Render(cfgShuffled)
		if resultShuffled == nil {
			t.Fatal("Render returned nil for valid config with shuffled stops")
		}

		// Render with the sorted stops
		cfgSorted := gradient.Config{
			Style:  style,
			Angle:  angle,
			Bounds: bounds,
			Stops:  sorted,
		}
		resultSorted := gradient.Render(cfgSorted)
		if resultSorted == nil {
			t.Fatal("Render returned nil for valid config with sorted stops")
		}

		// Compare pixel data — both should produce identical images
		imgShuffled, ok := resultShuffled.Image.(*image.RGBA)
		if !ok {
			t.Fatal("Shuffled result Image is not *image.RGBA")
		}
		imgSorted, ok := resultSorted.Image.(*image.RGBA)
		if !ok {
			t.Fatal("Sorted result Image is not *image.RGBA")
		}

		if !bytes.Equal(imgShuffled.Pix, imgSorted.Pix) {
			t.Fatalf("Rendering with shuffled stops produced different pixels than sorted stops (style=%d, angle=%.1f, bounds=%v, numStops=%d)",
				style, angle, bounds, numStops)
		}
	})
}

// For any valid Config where color stop positions are outside [0.0, 1.0],
// rendering with those positions SHALL produce pixel-identical output to
// rendering with all positions clamped to [0.0, 1.0].

func TestPropertyPositionClampingEquivalence(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate random style: Linear (0) or Radial (1)
		style := gradient.Style(rapid.IntRange(0, 1).Draw(t, "style"))

		// Generate small bounds (5-20) for efficiency
		width := rapid.IntRange(5, 20).Draw(t, "width")
		height := rapid.IntRange(5, 20).Draw(t, "height")
		bounds := image.Rect(0, 0, width, height)

		// Generate a valid angle
		angle := rapid.Float64Range(-720.0, 720.0).Draw(t, "angle")

		// Generate 2-5 color stops with positions in [-1.0, 2.0] (extending beyond [0,1])
		numStops := rapid.IntRange(2, 5).Draw(t, "numStops")
		stops := make([]gradient.ColorStop, numStops)
		for i := range stops {
			stops[i] = gradient.ColorStop{
				Position: rapid.Float64Range(-1.0, 2.0).Draw(t, "stopPos"),
				Color: color.RGBA{
					R: uint8(rapid.IntRange(0, 255).Draw(t, "stopR")),
					G: uint8(rapid.IntRange(0, 255).Draw(t, "stopG")),
					B: uint8(rapid.IntRange(0, 255).Draw(t, "stopB")),
					A: uint8(rapid.IntRange(0, 255).Draw(t, "stopA")),
				},
			}
		}

		// Create the original config with out-of-range positions
		cfgOriginal := gradient.Config{
			Style:  style,
			Angle:  angle,
			Bounds: bounds,
			Stops:  stops,
		}

		// Create a clamped version where each position is clamped to [0, 1]
		clampedStops := make([]gradient.ColorStop, numStops)
		for i, s := range stops {
			clampedStops[i] = gradient.ColorStop{
				Position: math.Max(0, math.Min(1, s.Position)),
				Color:    s.Color,
			}
		}

		cfgClamped := gradient.Config{
			Style:  style,
			Angle:  angle,
			Bounds: bounds,
			Stops:  clampedStops,
		}

		// Render both configs
		resultOriginal := gradient.Render(cfgOriginal)
		resultClamped := gradient.Render(cfgClamped)

		// Both must be non-nil (configs are valid: bounds >= 1, >= 2 stops, no NaN/Inf)
		if resultOriginal == nil {
			t.Fatal("Render returned nil for original config with out-of-range positions")
		}
		if resultClamped == nil {
			t.Fatal("Render returned nil for clamped config")
		}

		// Compare pixel data — both should produce identical images
		imgOriginal, ok := resultOriginal.Image.(*image.RGBA)
		if !ok {
			t.Fatal("Original result Image is not *image.RGBA")
		}
		imgClamped, ok := resultClamped.Image.(*image.RGBA)
		if !ok {
			t.Fatal("Clamped result Image is not *image.RGBA")
		}

		if !bytes.Equal(imgOriginal.Pix, imgClamped.Pix) {
			t.Fatalf("Rendering with out-of-range positions produced different pixels "+
				"than rendering with clamped positions.\nStyle=%d, Angle=%.2f, Bounds=%v\n"+
				"Original stops: %+v\nClamped stops: %+v",
				style, angle, bounds, stops, clampedStops)
		}
	})
}

// For any valid Config containing two or more stops at the same position value,
// the rendered color at that normalized position SHALL equal the color of the last
// such stop in original input order.

func TestPropertyDuplicatePositionStopUsesLastInInputOrder(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate two different random colors for the duplicate stops at position 0.0
		color1 := color.RGBA{
			R: uint8(rapid.IntRange(0, 255).Draw(t, "dup1R")),
			G: uint8(rapid.IntRange(0, 255).Draw(t, "dup1G")),
			B: uint8(rapid.IntRange(0, 255).Draw(t, "dup1B")),
			A: uint8(rapid.IntRange(0, 255).Draw(t, "dup1A")),
		}
		color2 := color.RGBA{
			R: uint8(rapid.IntRange(0, 255).Draw(t, "dup2R")),
			G: uint8(rapid.IntRange(0, 255).Draw(t, "dup2G")),
			B: uint8(rapid.IntRange(0, 255).Draw(t, "dup2B")),
			A: uint8(rapid.IntRange(0, 255).Draw(t, "dup2A")),
		}

		// Generate a color for the end stop at position 1.0
		endColor := color.RGBA{
			R: uint8(rapid.IntRange(0, 255).Draw(t, "endR")),
			G: uint8(rapid.IntRange(0, 255).Draw(t, "endG")),
			B: uint8(rapid.IntRange(0, 255).Draw(t, "endB")),
			A: uint8(rapid.IntRange(0, 255).Draw(t, "endA")),
		}

		// Build config: two stops at position 0.0 (color1 first, color2 second),
		// and one stop at position 1.0. Linear style, angle 0 (top-to-bottom).
		// The top row should render using color2 (last in input order at position 0.0).
		width := rapid.IntRange(1, 50).Draw(t, "width")
		height := rapid.IntRange(2, 50).Draw(t, "height")

		cfg := gradient.Config{
			Style:  gradient.Linear,
			Angle:  0, // top-to-bottom
			Bounds: image.Rect(0, 0, width, height),
			Stops: []gradient.ColorStop{
				{Position: 0.0, Color: color1}, // first duplicate at 0.0
				{Position: 0.0, Color: color2}, // second duplicate at 0.0 (last in input order)
				{Position: 1.0, Color: endColor},
			},
		}

		result := gradient.Render(cfg)
		if result == nil {
			t.Fatal("Render returned nil for valid config")
		}

		img := result.Image

		// The top-left pixel (0, 0) corresponds to gradient position 0.0.
		// It should have color2 since that's the last stop at position 0.0.
		r, g, b, a := img.At(0, 0).RGBA()
		gotR := uint8(r >> 8)
		gotG := uint8(g >> 8)
		gotB := uint8(b >> 8)
		gotA := uint8(a >> 8)

		if gotR != color2.R || gotG != color2.G || gotB != color2.B || gotA != color2.A {
			t.Fatalf("top-left pixel (0,0) has color RGBA(%d,%d,%d,%d), want RGBA(%d,%d,%d,%d) "+
				"(last duplicate stop at position 0.0); first duplicate was RGBA(%d,%d,%d,%d)",
				gotR, gotG, gotB, gotA,
				color2.R, color2.G, color2.B, color2.A,
				color1.R, color1.G, color1.B, color1.A)
		}
	})
}

// --- From: gradient_test.go ---

// TestLinear_TwoStops_Angle0 verifies a red-to-blue vertical gradient at 10×10.
// Angle 0 means top-to-bottom: top row = red, bottom row = blue, middle ≈ purple.

func TestLinear_TwoStops_Angle0(t *testing.T) {
	red := color.RGBA{R: 255, G: 0, B: 0, A: 255}
	blue := color.RGBA{R: 0, G: 0, B: 255, A: 255}

	cfg := gradient.Config{
		Style:  gradient.Linear,
		Angle:  0,
		Bounds: image.Rect(0, 0, 10, 10),
		Stops: []gradient.ColorStop{
			{Position: 0.0, Color: red},
			{Position: 1.0, Color: blue},
		},
	}

	result := gradient.Render(cfg)
	if result == nil {
		t.Fatal("Render returned nil for valid config")
	}

	img := result.Image

	// Top-left should be red (position 0.0)
	r, g, b, a := img.At(0, 0).RGBA()
	assertColorEqual(t, "top-left", uint8(r>>8), uint8(g>>8), uint8(b>>8), uint8(a>>8), red)

	// Bottom-left should be blue (position 1.0)
	r, g, b, a = img.At(0, 9).RGBA()
	assertColorEqual(t, "bottom-left", uint8(r>>8), uint8(g>>8), uint8(b>>8), uint8(a>>8), blue)

	// Middle row (y=4 or y=5) should be approximately purple (≈128, 0, ≈128, 255)
	r, g, b, a = img.At(0, 5).RGBA()
	midR, midG, midB, midA := uint8(r>>8), uint8(g>>8), uint8(b>>8), uint8(a>>8)
	// Allow ±2 tolerance for rounding
	if absDiff(midR, 128) > 15 {
		t.Errorf("middle row R=%d, want ≈128", midR)
	}
	if midG > 2 {
		t.Errorf("middle row G=%d, want ≈0", midG)
	}
	if absDiff(midB, 128) > 15 {
		t.Errorf("middle row B=%d, want ≈128", midB)
	}
	if midA != 255 {
		t.Errorf("middle row A=%d, want 255", midA)
	}
}

// TestLinear_TwoStops_Angle90 verifies a left-to-right horizontal gradient at 10×10.
// Angle 90 means left-to-right: left column = red, right column = blue.

func TestLinear_TwoStops_Angle90(t *testing.T) {
	red := color.RGBA{R: 255, G: 0, B: 0, A: 255}
	blue := color.RGBA{R: 0, G: 0, B: 255, A: 255}

	cfg := gradient.Config{
		Style:  gradient.Linear,
		Angle:  90,
		Bounds: image.Rect(0, 0, 10, 10),
		Stops: []gradient.ColorStop{
			{Position: 0.0, Color: red},
			{Position: 1.0, Color: blue},
		},
	}

	result := gradient.Render(cfg)
	if result == nil {
		t.Fatal("Render returned nil for valid config")
	}

	img := result.Image

	// Top-left should be red (left column at angle 90 = position 0.0)
	r, g, b, a := img.At(0, 0).RGBA()
	assertColorEqual(t, "top-left", uint8(r>>8), uint8(g>>8), uint8(b>>8), uint8(a>>8), red)

	// Top-right should be blue (right column at angle 90 = position 1.0)
	r, g, b, a = img.At(9, 0).RGBA()
	assertColorEqual(t, "top-right", uint8(r>>8), uint8(g>>8), uint8(b>>8), uint8(a>>8), blue)
}

// TestLinear_Angle45_Diagonal verifies that a 45° angle produces a diagonal gradient
// where the top-left corner differs from the bottom-right corner.

func TestLinear_Angle45_Diagonal(t *testing.T) {
	white := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	black := color.RGBA{R: 0, G: 0, B: 0, A: 255}

	cfg := gradient.Config{
		Style:  gradient.Linear,
		Angle:  45,
		Bounds: image.Rect(0, 0, 10, 10),
		Stops: []gradient.ColorStop{
			{Position: 0.0, Color: white},
			{Position: 1.0, Color: black},
		},
	}

	result := gradient.Render(cfg)
	if result == nil {
		t.Fatal("Render returned nil for valid config")
	}

	img := result.Image

	// Top-left corner (0,0): at 45° clockwise from north, direction is (sin45, cos45) ≈ (0.707, 0.707)
	// projection at (0,0) = 0 → should be position 0.0 (white)
	r0, g0, b0, _ := img.At(0, 0).RGBA()
	// Bottom-right corner (9,9): max projection → should be position 1.0 (black)
	r9, g9, b9, _ := img.At(9, 9).RGBA()

	// Top-left should be lighter than bottom-right
	tlLum := uint8(r0>>8) + uint8(g0>>8) + uint8(b0>>8)
	brLum := uint8(r9>>8) + uint8(g9>>8) + uint8(b9>>8)

	if tlLum <= brLum {
		t.Errorf("top-left luminance (%d) should be greater than bottom-right (%d) at 45°", tlLum, brLum)
	}

	// Top-left should be white and bottom-right should be black
	assertColorEqual(t, "top-left (45°)", uint8(r0>>8), uint8(g0>>8), uint8(b0>>8), 255, white)
	assertColorEqual(t, "bottom-right (45°)", uint8(r9>>8), uint8(g9>>8), uint8(b9>>8), 255, black)
}

// TestRadial_TwoStops_Square verifies radial symmetry in a square (20×20) bounds.
// The center pixel should be the first stop color.

func TestRadial_TwoStops_Square(t *testing.T) {
	yellow := color.RGBA{R: 255, G: 255, B: 0, A: 255}
	purple := color.RGBA{R: 128, G: 0, B: 128, A: 255}

	cfg := gradient.Config{
		Style:  gradient.Radial,
		Angle:  0,
		Bounds: image.Rect(0, 0, 20, 20),
		Stops: []gradient.ColorStop{
			{Position: 0.0, Color: yellow},
			{Position: 1.0, Color: purple},
		},
	}

	result := gradient.Render(cfg)
	if result == nil {
		t.Fatal("Render returned nil for valid radial config")
	}

	img := result.Image

	// Center pixel: (20/2, 20/2) = (10, 10)
	r, g, b, a := img.At(10, 10).RGBA()
	assertColorEqual(t, "center", uint8(r>>8), uint8(g>>8), uint8(b>>8), uint8(a>>8), yellow)
}

// TestRadial_TwoStops_Rectangle verifies a radial gradient in non-square bounds (40×10).
// Center should be first stop, corners should be last stop (outside inscribed circle).

func TestRadial_TwoStops_Rectangle(t *testing.T) {
	green := color.RGBA{R: 0, G: 255, B: 0, A: 255}
	red := color.RGBA{R: 255, G: 0, B: 0, A: 255}

	cfg := gradient.Config{
		Style:  gradient.Radial,
		Angle:  0,
		Bounds: image.Rect(0, 0, 40, 10),
		Stops: []gradient.ColorStop{
			{Position: 0.0, Color: green},
			{Position: 1.0, Color: red},
		},
	}

	result := gradient.Render(cfg)
	if result == nil {
		t.Fatal("Render returned nil for valid radial config")
	}

	img := result.Image

	// Center pixel: (40/2, 10/2) = (20, 5)
	r, g, b, a := img.At(20, 5).RGBA()
	assertColorEqual(t, "center", uint8(r>>8), uint8(g>>8), uint8(b>>8), uint8(a>>8), green)

	// Corners are outside inscribed circle (radius = min(40,10)/2 = 5).
	// The corner (0,0) distance from center (20,5): sqrt(20²+5²) ≈ 20.6 >> 5
	r, g, b, a = img.At(0, 0).RGBA()
	assertColorEqual(t, "corner (0,0)", uint8(r>>8), uint8(g>>8), uint8(b>>8), uint8(a>>8), red)

	r, g, b, a = img.At(39, 9).RGBA()
	assertColorEqual(t, "corner (39,9)", uint8(r>>8), uint8(g>>8), uint8(b>>8), uint8(a>>8), red)
}

// TestThreeStops_Linear verifies a multi-stop gradient with red/green/blue at 0.0/0.5/1.0.
// The midpoint row should be approximately green.

func TestThreeStops_Linear(t *testing.T) {
	red := color.RGBA{R: 255, G: 0, B: 0, A: 255}
	green := color.RGBA{R: 0, G: 255, B: 0, A: 255}
	blue := color.RGBA{R: 0, G: 0, B: 255, A: 255}

	// Use 11 pixel height so midpoint row (y=5) is exactly at t=0.5
	cfg := gradient.Config{
		Style:  gradient.Linear,
		Angle:  0,
		Bounds: image.Rect(0, 0, 10, 11),
		Stops: []gradient.ColorStop{
			{Position: 0.0, Color: red},
			{Position: 0.5, Color: green},
			{Position: 1.0, Color: blue},
		},
	}

	result := gradient.Render(cfg)
	if result == nil {
		t.Fatal("Render returned nil for valid config")
	}

	img := result.Image

	// Midpoint row (y=5): t = 5/10 = 0.5, should be green
	r, g, b, a := img.At(0, 5).RGBA()
	gotR, gotG, gotB, gotA := uint8(r>>8), uint8(g>>8), uint8(b>>8), uint8(a>>8)
	assertColorEqual(t, "midpoint row", gotR, gotG, gotB, gotA, green)
}

// TestMaxStops_64 verifies that a gradient with 64 stops renders successfully.

func TestMaxStops_64(t *testing.T) {
	stops := make([]gradient.ColorStop, 64)
	for i := range stops {
		pos := float64(i) / 63.0
		v := uint8(float64(i) * 255.0 / 63.0)
		stops[i] = gradient.ColorStop{
			Position: pos,
			Color:    color.RGBA{R: v, G: v, B: v, A: 255},
		}
	}

	cfg := gradient.Config{
		Style:  gradient.Linear,
		Angle:  0,
		Bounds: image.Rect(0, 0, 50, 50),
		Stops:  stops,
	}

	result := gradient.Render(cfg)
	if result == nil {
		t.Fatal("Render returned nil for valid 64-stop config")
	}

	// Verify image dimensions
	imgBounds := result.Image.Bounds()
	if imgBounds.Dx() != 50 || imgBounds.Dy() != 50 {
		t.Errorf("image dimensions = %dx%d, want 50x50", imgBounds.Dx(), imgBounds.Dy())
	}
}

// TestNilCases verifies that each invalid condition returns nil.

func TestNilCases(t *testing.T) {
	validStops := []gradient.ColorStop{
		{Position: 0.0, Color: color.RGBA{R: 255, A: 255}},
		{Position: 1.0, Color: color.RGBA{B: 255, A: 255}},
	}

	tests := []struct {
		name string
		cfg  gradient.Config
	}{
		{
			name: "zero width bounds",
			cfg: gradient.Config{
				Style:  gradient.Linear,
				Angle:  0,
				Bounds: image.Rect(0, 0, 0, 10),
				Stops:  validStops,
			},
		},
		{
			name: "zero height bounds",
			cfg: gradient.Config{
				Style:  gradient.Linear,
				Angle:  0,
				Bounds: image.Rect(0, 0, 10, 0),
				Stops:  validStops,
			},
		},
		{
			name: "negative width bounds",
			cfg: gradient.Config{
				Style:  gradient.Linear,
				Angle:  0,
				Bounds: image.Rectangle{Min: image.Point{X: 10, Y: 0}, Max: image.Point{X: 5, Y: 10}},
				Stops:  validStops,
			},
		},
		{
			name: "fewer than 2 stops (empty)",
			cfg: gradient.Config{
				Style:  gradient.Linear,
				Angle:  0,
				Bounds: image.Rect(0, 0, 10, 10),
				Stops:  []gradient.ColorStop{},
			},
		},
		{
			name: "fewer than 2 stops (one)",
			cfg: gradient.Config{
				Style:  gradient.Linear,
				Angle:  0,
				Bounds: image.Rect(0, 0, 10, 10),
				Stops:  validStops[:1],
			},
		},
		{
			name: "invalid style",
			cfg: gradient.Config{
				Style:  gradient.Style(99),
				Angle:  0,
				Bounds: image.Rect(0, 0, 10, 10),
				Stops:  validStops,
			},
		},
		{
			name: "NaN stop position",
			cfg: gradient.Config{
				Style:  gradient.Linear,
				Angle:  0,
				Bounds: image.Rect(0, 0, 10, 10),
				Stops: []gradient.ColorStop{
					{Position: 0.0, Color: color.RGBA{R: 255, A: 255}},
					{Position: math.NaN(), Color: color.RGBA{B: 255, A: 255}},
				},
			},
		},
		{
			name: "+Inf stop position",
			cfg: gradient.Config{
				Style:  gradient.Linear,
				Angle:  0,
				Bounds: image.Rect(0, 0, 10, 10),
				Stops: []gradient.ColorStop{
					{Position: 0.0, Color: color.RGBA{R: 255, A: 255}},
					{Position: math.Inf(1), Color: color.RGBA{B: 255, A: 255}},
				},
			},
		},
		{
			name: "-Inf stop position",
			cfg: gradient.Config{
				Style:  gradient.Linear,
				Angle:  0,
				Bounds: image.Rect(0, 0, 10, 10),
				Stops: []gradient.ColorStop{
					{Position: math.Inf(-1), Color: color.RGBA{R: 255, A: 255}},
					{Position: 1.0, Color: color.RGBA{B: 255, A: 255}},
				},
			},
		},
		{
			name: "NaN angle",
			cfg: gradient.Config{
				Style:  gradient.Linear,
				Angle:  math.NaN(),
				Bounds: image.Rect(0, 0, 10, 10),
				Stops:  validStops,
			},
		},
		{
			name: "+Inf angle",
			cfg: gradient.Config{
				Style:  gradient.Linear,
				Angle:  math.Inf(1),
				Bounds: image.Rect(0, 0, 10, 10),
				Stops:  validStops,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := gradient.Render(tc.cfg)
			if result != nil {
				t.Errorf("expected nil for %s, got non-nil result", tc.name)
			}
		})
	}
}

// --- Helpers ---

func assertColorEqual(t *testing.T, label string, gotR, gotG, gotB, gotA uint8, want color.RGBA) {
	t.Helper()
	if gotR != want.R || gotG != want.G || gotB != want.B || gotA != want.A {
		t.Errorf("%s: got RGBA(%d,%d,%d,%d), want RGBA(%d,%d,%d,%d)",
			label, gotR, gotG, gotB, gotA, want.R, want.G, want.B, want.A)
	}
}

func absDiff(a, b uint8) uint8 {
	if a > b {
		return a - b
	}
	return b - a
}
