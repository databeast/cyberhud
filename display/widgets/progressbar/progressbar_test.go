package progressbar

import (
	"image"
	"image/color"
	"math"
	"testing"
	"time"

	"github.com/databeast/cyberhud/display/widgets"
	"github.com/databeast/cyberhud/display/widgets/gradient"
	"pgregory.net/rapid"
)

// --- From: progressbar_prop_test.go ---

// TestPropertyPieFillAreaProportionality verifies that for any valid bounds and value,
// the ratio of foreground pixels within the inscribed circle approximates the value.

func TestPropertyPieFillAreaProportionality(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		value := rapid.Float64Range(0.0, 1.0).Draw(t, "value")
		width := rapid.IntRange(3, 100).Draw(t, "width")
		height := rapid.IntRange(3, 100).Draw(t, "height")

		bounds := image.Rect(0, 0, width, height)

		fg := color.RGBA{R: 255, G: 0, B: 0, A: 255}
		bg := color.RGBA{R: 0, G: 0, B: 255, A: 255}

		result := Render(Config{
			Style:      Pie,
			Value:      value,
			Bounds:     bounds,
			Foreground: fg,
			Background: bg,
		})

		if result == nil {
			t.Fatal("expected non-nil result for valid pie bounds")
		}

		// Compute inscribed circle parameters
		cx := float64(width) / 2.0
		cy := float64(height) / 2.0
		r := math.Min(float64(width), float64(height)) / 2.0

		// Count foreground and total pixels within the inscribed circle
		fgCount := 0
		totalCount := 0

		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				// Use pixel center for distance calculation
				px := float64(x) + 0.5
				py := float64(y) + 0.5
				dx := px - cx
				dy := py - cy
				dist := math.Sqrt(dx*dx + dy*dy)

				if dist <= r {
					totalCount++
					pixel := result.Image.At(x, y)
					pr, pg, pb, pa := pixel.RGBA()
					// Compare with foreground color (red)
					if pr == 0xFFFF && pg == 0 && pb == 0 && pa == 0xFFFF {
						fgCount++
					}
				}
			}
		}

		if totalCount == 0 {
			t.Fatal("no pixels found within inscribed circle")
		}

		// Compute actual ratio of foreground pixels
		actualRatio := float64(fgCount) / float64(totalCount)

		// Tolerance accounts for rasterization error in small circles
		tolerance := 2.0/r + 0.02

		if math.Abs(actualRatio-value) > tolerance {
			t.Fatalf("pie fill ratio mismatch: value=%f, actualRatio=%f, tolerance=%f, width=%d, height=%d, fgCount=%d, totalCount=%d",
				value, actualRatio, tolerance, width, height, fgCount, totalCount)
		}
	})
}

// TestPropertyVerticalFillProportionality verifies that for any valid bounds and value,
// the correct number of rows from the bottom are filled with foreground.
//
// For any valid bounds (width ≥ 1, height ≥ 1) and any value v in [0.0, 1.0],
// rendering a vertical progress bar SHALL produce an image where exactly
// floor(height * v) rows from the bottom are filled with the foreground color
// and the remaining rows are filled with the background color.

func TestPropertyVerticalFillProportionality(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		value := rapid.Float64Range(0.0, 1.0).Draw(t, "value")
		width := rapid.IntRange(1, 100).Draw(t, "width")
		height := rapid.IntRange(1, 200).Draw(t, "height")

		bounds := image.Rect(0, 0, width, height)
		fg := color.RGBA{R: 255, G: 0, B: 0, A: 255}
		bg := color.RGBA{R: 0, G: 0, B: 255, A: 255}

		result := Render(Config{
			Style:       Linear,
			Orientation: OrientVertical,
			Value:       value,
			Bounds:      bounds,
			Foreground:  fg,
			Background:  bg,
		})

		if result == nil {
			t.Fatal("Render returned nil for valid bounds")
		}

		expectedFillRows := int(math.Floor(float64(height) * value))

		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				got := result.Image.At(x, y)
				r, g, b, a := got.RGBA()
				if y >= height-expectedFillRows {
					// Should be foreground (red) - bottom rows
					if r != 0xFFFF || g != 0 || b != 0 || a != 0xFFFF {
						t.Fatalf("pixel (%d,%d): expected fg %v, got (%d,%d,%d,%d) [value=%.6f, height=%d, fillRows=%d]",
							x, y, fg, r, g, b, a, value, height, expectedFillRows)
					}
				} else {
					// Should be background (blue) - top rows
					if r != 0 || g != 0 || b != 0xFFFF || a != 0xFFFF {
						t.Fatalf("pixel (%d,%d): expected bg %v, got (%d,%d,%d,%d) [value=%.6f, height=%d, fillRows=%d]",
							x, y, bg, r, g, b, a, value, height, expectedFillRows)
					}
				}
			}
		}
	})
}

// For any valid bounds (width ≥ 1, height ≥ 1) and any value v in [0.0, 1.0],
// rendering a horizontal progress bar SHALL produce an image where exactly
// floor(width * v) columns from the left are filled with the foreground color
// and the remaining columns are filled with the background color.

func TestPropertyHorizontalFillProportionality(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		value := rapid.Float64Range(0.0, 1.0).Draw(t, "value")
		width := rapid.IntRange(1, 200).Draw(t, "width")
		height := rapid.IntRange(1, 100).Draw(t, "height")

		bounds := image.Rect(0, 0, width, height)
		fg := color.RGBA{R: 255, G: 0, B: 0, A: 255}
		bg := color.RGBA{R: 0, G: 0, B: 255, A: 255}

		result := Render(Config{
			Style:       Linear,
			Orientation: OrientHorizontal,
			Value:       value,
			Bounds:      bounds,
			Foreground:  fg,
			Background:  bg,
		})

		if result == nil {
			t.Fatal("Render returned nil for valid bounds")
		}

		expectedFillCols := int(math.Floor(float64(width) * value))

		for x := 0; x < width; x++ {
			for y := 0; y < height; y++ {
				got := result.Image.At(x, y)
				r, g, b, a := got.RGBA()
				if x < expectedFillCols {
					// Should be foreground (red)
					if r != 0xFFFF || g != 0 || b != 0 || a != 0xFFFF {
						t.Fatalf("pixel (%d,%d): expected fg %v, got (%d,%d,%d,%d) [value=%.6f, width=%d, fillCols=%d]",
							x, y, fg, r, g, b, a, value, width, expectedFillCols)
					}
				} else {
					// Should be background (blue)
					if r != 0 || g != 0 || b != 0xFFFF || a != 0xFFFF {
						t.Fatalf("pixel (%d,%d): expected bg %v, got (%d,%d,%d,%d) [value=%.6f, width=%d, fillCols=%d]",
							x, y, bg, r, g, b, a, value, width, expectedFillCols)
					}
				}
			}
		}
	})
}

// TestPropertyOutputMetadataCorrectness verifies that for any style, valid bounds, and value,
// the returned Result has Position equal to bounds.Min and Label equal to the expected
// style identifier string.

func TestPropertyOutputMetadataCorrectness(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		style := Style(rapid.IntRange(0, 4).Draw(t, "style"))

		minX := rapid.IntRange(0, 100).Draw(t, "minX")
		minY := rapid.IntRange(0, 100).Draw(t, "minY")

		var width, height int
		switch style {
		case Pie, Ring, Arc:
			width = rapid.IntRange(3, 100).Draw(t, "width")
			height = rapid.IntRange(3, 100).Draw(t, "height")
		default:
			width = rapid.IntRange(1, 100).Draw(t, "width")
			height = rapid.IntRange(1, 100).Draw(t, "height")
		}

		bounds := image.Rect(minX, minY, minX+width, minY+height)
		value := rapid.Float64Range(0.0, 1.0).Draw(t, "value")

		fg := color.RGBA{R: 200, G: 100, B: 50, A: 255}
		bg := color.RGBA{R: 10, G: 20, B: 30, A: 255}

		result := Render(Config{
			Style:      style,
			Value:      value,
			Bounds:     bounds,
			Foreground: fg,
			Background: bg,
		})

		if result == nil {
			t.Fatal("expected non-nil result for valid bounds")
		}

		// Assert Position == bounds.Min
		if result.Position.X != minX || result.Position.Y != minY {
			t.Fatalf("Position mismatch: got (%d, %d), want (%d, %d)",
				result.Position.X, result.Position.Y, minX, minY)
		}

		// Assert Label matches expected style identifier
		var expectedLabel string
		switch style {
		case Linear:
			expectedLabel = "progressbar/linear"
		case Pie:
			expectedLabel = "progressbar/pie"
		case Segmented:
			expectedLabel = "progressbar/segmented"
		case Ring:
			expectedLabel = "progressbar/ring"
		case Arc:
			expectedLabel = "progressbar/arc"
		}

		if result.Label != expectedLabel {
			t.Fatalf("Label mismatch: got %q, want %q (style=%d)",
				result.Label, expectedLabel, style)
		}
	})
}

// TestPropertyOutputDimensionsMatchBounds verifies that for any style and any valid bounds
// that do not trigger nil output, the rendered image's width equals bounds.Dx() and its
// height equals bounds.Dy().

func TestPropertyOutputDimensionsMatchBounds(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		style := Style(rapid.IntRange(0, 4).Draw(t, "style"))

		var width, height int
		switch style {
		case Pie, Ring, Arc:
			width = rapid.IntRange(3, 150).Draw(t, "width")
			height = rapid.IntRange(3, 150).Draw(t, "height")
		default:
			width = rapid.IntRange(1, 200).Draw(t, "width")
			height = rapid.IntRange(1, 200).Draw(t, "height")
		}

		bounds := image.Rect(0, 0, width, height)
		fg := color.RGBA{R: 200, G: 100, B: 50, A: 255}
		bg := color.RGBA{R: 10, G: 20, B: 30, A: 255}

		value := rapid.Float64Range(0.0, 1.0).Draw(t, "value")

		result := Render(Config{
			Style:      style,
			Value:      value,
			Bounds:     bounds,
			Foreground: fg,
			Background: bg,
		})

		if result == nil {
			t.Fatal("expected non-nil result for valid bounds")
		}

		gotWidth := result.Image.Bounds().Dx()
		gotHeight := result.Image.Bounds().Dy()

		if gotWidth != width {
			t.Fatalf("image width mismatch: got %d, want %d (style=%d, bounds=%v)",
				gotWidth, width, style, bounds)
		}
		if gotHeight != height {
			t.Fatalf("image height mismatch: got %d, want %d (style=%d, bounds=%v)",
				gotHeight, height, style, bounds)
		}
	})
}

// TestPropertyValueClampingIdempotence verifies that for any style, bounds, and
// out-of-range value, the rendered output is pixel-identical to rendering with
// the value clamped to [0.0, 1.0].

func TestPropertyValueClampingIdempotence(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate random style
		style := Style(rapid.IntRange(0, 4).Draw(t, "style"))

		// Generate bounds appropriate for style
		var w, h int
		switch style {
		case Pie, Ring, Arc:
			w = rapid.IntRange(3, 50).Draw(t, "width")
			h = rapid.IntRange(3, 50).Draw(t, "height")
		default:
			w = rapid.IntRange(1, 50).Draw(t, "width")
			h = rapid.IntRange(1, 50).Draw(t, "height")
		}
		bounds := image.Rect(0, 0, w, h)

		// Generate out-of-range value from [-100, -0.001] ∪ [1.001, 100]
		negative := rapid.SampledFrom([]bool{true, false}).Draw(t, "negative")
		var value float64
		if negative {
			value = rapid.Float64Range(-100.0, -0.001).Draw(t, "value")
		} else {
			value = rapid.Float64Range(1.001, 100.0).Draw(t, "value")
		}

		// Compute clamped value manually
		clampedValue := value
		if clampedValue < 0.0 {
			clampedValue = 0.0
		}
		if clampedValue > 1.0 {
			clampedValue = 1.0
		}

		// Use non-zero colors to avoid default resolution
		fg := color.RGBA{R: 200, G: 100, B: 50, A: 255}
		bg := color.RGBA{R: 10, G: 20, B: 30, A: 255}

		// Render with out-of-range value
		result1 := Render(Config{
			Style:      style,
			Value:      value,
			Bounds:     bounds,
			Foreground: fg,
			Background: bg,
		})

		// Render with clamped value
		result2 := Render(Config{
			Style:      style,
			Value:      clampedValue,
			Bounds:     bounds,
			Foreground: fg,
			Background: bg,
		})

		// Assert both are non-nil
		if result1 == nil {
			t.Fatal("expected non-nil result for out-of-range value render")
		}
		if result2 == nil {
			t.Fatal("expected non-nil result for clamped value render")
		}

		// Assert pixel-identical: iterate all pixels and compare RGBA values
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				r1, g1, b1, a1 := result1.Image.At(x, y).RGBA()
				r2, g2, b2, a2 := result2.Image.At(x, y).RGBA()
				if r1 != r2 || g1 != g2 || b1 != b2 || a1 != a2 {
					t.Fatalf("pixel mismatch at (%d,%d): out-of-range=(%d,%d,%d,%d), clamped=(%d,%d,%d,%d) [style=%d, value=%f, clamped=%f, bounds=%dx%d]",
						x, y, r1, g1, b1, a1, r2, g2, b2, a2, style, value, clampedValue, w, h)
				}
			}
		}
	})
}

// TestProperty16_ValueClampingIdempotence verifies Property 16: For any style and any
// value outside [0.0, 1.0] (including NaN and ±Infinity), the rendered output SHALL be
// pixel-identical to rendering with the value clamped (NaN/Inf → 0.0, negatives → 0.0,
// >1.0 → 1.0). Extended to cover all 5 styles and both orientations.

func TestProperty16_ValueClampingIdempotence(t *testing.T) {
	t.Run("out_of_range_values", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			style := Style(rapid.IntRange(0, 4).Draw(t, "style"))
			orientation := Orientation(rapid.IntRange(0, 1).Draw(t, "orientation"))

			var w, h int
			switch style {
			case Pie, Ring, Arc:
				w = rapid.IntRange(3, 50).Draw(t, "width")
				h = rapid.IntRange(3, 50).Draw(t, "height")
			default:
				w = rapid.IntRange(1, 50).Draw(t, "width")
				h = rapid.IntRange(1, 50).Draw(t, "height")
			}
			bounds := image.Rect(0, 0, w, h)

			// Generate out-of-range value: negative or > 1
			negative := rapid.SampledFrom([]bool{true, false}).Draw(t, "negative")
			var value float64
			if negative {
				value = rapid.Float64Range(-1000.0, -0.001).Draw(t, "value")
			} else {
				value = rapid.Float64Range(1.001, 1000.0).Draw(t, "value")
			}

			// Compute expected clamped value
			clampedValue := value
			if clampedValue < 0.0 {
				clampedValue = 0.0
			}
			if clampedValue > 1.0 {
				clampedValue = 1.0
			}

			fg := color.RGBA{R: 200, G: 100, B: 50, A: 255}
			bg := color.RGBA{R: 10, G: 20, B: 30, A: 255}

			baseCfg := Config{
				Style:       style,
				Orientation: orientation,
				Bounds:      bounds,
				Foreground:  fg,
				Background:  bg,
			}

			cfg1 := baseCfg
			cfg1.Value = value

			cfg2 := baseCfg
			cfg2.Value = clampedValue

			result1 := Render(cfg1)
			result2 := Render(cfg2)

			if result1 == nil {
				t.Fatal("expected non-nil result for out-of-range value render")
			}
			if result2 == nil {
				t.Fatal("expected non-nil result for clamped value render")
			}

			for y := 0; y < h; y++ {
				for x := 0; x < w; x++ {
					r1, g1, b1, a1 := result1.Image.At(x, y).RGBA()
					r2, g2, b2, a2 := result2.Image.At(x, y).RGBA()
					if r1 != r2 || g1 != g2 || b1 != b2 || a1 != a2 {
						t.Fatalf("pixel mismatch at (%d,%d): out-of-range=(%d,%d,%d,%d), clamped=(%d,%d,%d,%d) [style=%d, orient=%d, value=%f, clamped=%f]",
							x, y, r1, g1, b1, a1, r2, g2, b2, a2, style, orientation, value, clampedValue)
					}
				}
			}
		})
	})

	t.Run("nan_and_infinity", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			style := Style(rapid.IntRange(0, 4).Draw(t, "style"))
			orientation := Orientation(rapid.IntRange(0, 1).Draw(t, "orientation"))

			var w, h int
			switch style {
			case Pie, Ring, Arc:
				w = rapid.IntRange(3, 50).Draw(t, "width")
				h = rapid.IntRange(3, 50).Draw(t, "height")
			default:
				w = rapid.IntRange(1, 50).Draw(t, "width")
				h = rapid.IntRange(1, 50).Draw(t, "height")
			}
			bounds := image.Rect(0, 0, w, h)

			// Pick a special value: NaN, +Inf, or -Inf
			specialIdx := rapid.IntRange(0, 2).Draw(t, "specialIdx")
			var value float64
			switch specialIdx {
			case 0:
				value = math.NaN()
			case 1:
				value = math.Inf(1)
			case 2:
				value = math.Inf(-1)
			}

			// All special values clamp to 0.0
			clampedValue := 0.0

			fg := color.RGBA{R: 200, G: 100, B: 50, A: 255}
			bg := color.RGBA{R: 10, G: 20, B: 30, A: 255}

			baseCfg := Config{
				Style:       style,
				Orientation: orientation,
				Bounds:      bounds,
				Foreground:  fg,
				Background:  bg,
			}

			cfg1 := baseCfg
			cfg1.Value = value

			cfg2 := baseCfg
			cfg2.Value = clampedValue

			result1 := Render(cfg1)
			result2 := Render(cfg2)

			if result1 == nil {
				t.Fatal("expected non-nil result for NaN/Inf value render")
			}
			if result2 == nil {
				t.Fatal("expected non-nil result for clamped (0.0) value render")
			}

			for y := 0; y < h; y++ {
				for x := 0; x < w; x++ {
					r1, g1, b1, a1 := result1.Image.At(x, y).RGBA()
					r2, g2, b2, a2 := result2.Image.At(x, y).RGBA()
					if r1 != r2 || g1 != g2 || b1 != b2 || a1 != a2 {
						t.Fatalf("pixel mismatch at (%d,%d): special=(%d,%d,%d,%d), clamped=(%d,%d,%d,%d) [style=%d, orient=%d, specialIdx=%d]",
							x, y, r1, g1, b1, a1, r2, g2, b2, a2, style, orientation, specialIdx)
					}
				}
			}
		})
	})
}

// TestPropertyInvalidBoundsProduceNil verifies that Render returns nil when bounds
// are too small for meaningful rendering.

func TestPropertyInvalidBoundsProduceNil(t *testing.T) {
	t.Run("zero_or_negative_bounds", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			style := Style(rapid.IntRange(0, 4).Draw(t, "style"))

			// Decide which dimension(s) to make invalid: 0=width, 1=height, 2=both
			invalidDim := rapid.IntRange(0, 2).Draw(t, "invalidDim")

			var bounds image.Rectangle
			switch invalidDim {
			case 0: // width invalid (Dx() ≤ 0), height valid
				w := rapid.IntRange(-10, 0).Draw(t, "width")
				h := rapid.IntRange(1, 50).Draw(t, "height")
				bounds = image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: w, Y: h}}
			case 1: // height invalid (Dy() ≤ 0), width valid
				w := rapid.IntRange(1, 50).Draw(t, "width")
				h := rapid.IntRange(-10, 0).Draw(t, "height")
				bounds = image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: w, Y: h}}
			case 2: // both invalid
				w := rapid.IntRange(-10, 0).Draw(t, "width")
				h := rapid.IntRange(-10, 0).Draw(t, "height")
				bounds = image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: w, Y: h}}
			}

			result := Render(Config{
				Style:      style,
				Value:      rapid.Float64Range(0.0, 1.0).Draw(t, "value"),
				Bounds:     bounds,
				Foreground: color.RGBA{R: 255, G: 0, B: 0, A: 255},
				Background: color.RGBA{R: 0, G: 0, B: 255, A: 255},
			})

			if result != nil {
				t.Fatalf("expected nil for invalid bounds (Dx=%d, Dy=%d, style=%d), got non-nil result",
					bounds.Dx(), bounds.Dy(), style)
			}
		})
	})

	t.Run("circular_too_small", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			// Test Pie, Ring, and Arc styles with bounds < 3×3
			style := Style(rapid.SampledFrom([]Style{Pie, Ring, Arc}).Draw(t, "style"))
			width := rapid.IntRange(1, 2).Draw(t, "width")
			height := rapid.IntRange(1, 2).Draw(t, "height")

			bounds := image.Rect(0, 0, width, height)
			result := Render(Config{
				Style:      style,
				Value:      rapid.Float64Range(0.0, 1.0).Draw(t, "value"),
				Bounds:     bounds,
				Foreground: color.RGBA{R: 255, G: 0, B: 0, A: 255},
				Background: color.RGBA{R: 0, G: 0, B: 255, A: 255},
			})

			if result != nil {
				t.Fatalf("expected nil for %d style with small bounds (w=%d, h=%d), got non-nil result",
					style, width, height)
			}
		})
	})
}

// TestProperty17_NilForInvalidBoundsOrStyle verifies Property 17: For any configuration
// where Bounds.Dx() ≤ 0 OR Bounds.Dy() ≤ 0, OR Style is Ring/Arc with bounds < 3×3,
// OR Style is an unrecognized value, Render SHALL return nil.

func TestProperty17_NilForInvalidBoundsOrStyle(t *testing.T) {
	t.Run("zero_or_negative_bounds_all_styles", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			style := Style(rapid.IntRange(0, 4).Draw(t, "style"))
			orientation := Orientation(rapid.IntRange(0, 1).Draw(t, "orientation"))

			// Decide which dimension(s) to make invalid: 0=width, 1=height, 2=both
			invalidDim := rapid.IntRange(0, 2).Draw(t, "invalidDim")

			var bounds image.Rectangle
			switch invalidDim {
			case 0: // width invalid (Dx() ≤ 0), height valid
				w := rapid.IntRange(-10, 0).Draw(t, "width")
				h := rapid.IntRange(1, 50).Draw(t, "height")
				bounds = image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: w, Y: h}}
			case 1: // height invalid (Dy() ≤ 0), width valid
				w := rapid.IntRange(1, 50).Draw(t, "width")
				h := rapid.IntRange(-10, 0).Draw(t, "height")
				bounds = image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: w, Y: h}}
			case 2: // both invalid
				w := rapid.IntRange(-10, 0).Draw(t, "width")
				h := rapid.IntRange(-10, 0).Draw(t, "height")
				bounds = image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: w, Y: h}}
			}

			result := Render(Config{
				Style:       style,
				Orientation: orientation,
				Value:       rapid.Float64Range(0.0, 1.0).Draw(t, "value"),
				Bounds:      bounds,
				Foreground:  color.RGBA{R: 255, G: 0, B: 0, A: 255},
				Background:  color.RGBA{R: 0, G: 0, B: 255, A: 255},
			})

			if result != nil {
				t.Fatalf("expected nil for invalid bounds (Dx=%d, Dy=%d, style=%d, orient=%d), got non-nil",
					bounds.Dx(), bounds.Dy(), style, orientation)
			}
		})
	})

	t.Run("ring_arc_bounds_less_than_3x3", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			style := Style(rapid.SampledFrom([]Style{Ring, Arc}).Draw(t, "style"))
			orientation := Orientation(rapid.IntRange(0, 1).Draw(t, "orientation"))

			// At least one dimension < 3 (but both positive so we don't hit the zero/negative case)
			dimCase := rapid.IntRange(0, 2).Draw(t, "dimCase")
			var width, height int
			switch dimCase {
			case 0: // width < 3, height valid
				width = rapid.IntRange(1, 2).Draw(t, "width")
				height = rapid.IntRange(3, 50).Draw(t, "height")
			case 1: // height < 3, width valid
				width = rapid.IntRange(3, 50).Draw(t, "width")
				height = rapid.IntRange(1, 2).Draw(t, "height")
			case 2: // both < 3
				width = rapid.IntRange(1, 2).Draw(t, "width")
				height = rapid.IntRange(1, 2).Draw(t, "height")
			}

			bounds := image.Rect(0, 0, width, height)
			result := Render(Config{
				Style:       style,
				Orientation: orientation,
				Value:       rapid.Float64Range(0.0, 1.0).Draw(t, "value"),
				Bounds:      bounds,
				Foreground:  color.RGBA{R: 255, G: 0, B: 0, A: 255},
				Background:  color.RGBA{R: 0, G: 0, B: 255, A: 255},
			})

			if result != nil {
				t.Fatalf("expected nil for %d style with bounds < 3×3 (w=%d, h=%d, orient=%d), got non-nil",
					style, width, height, orientation)
			}
		})
	})

	t.Run("unrecognized_style", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			// Generate a style value outside the recognized range [0, 4]
			negative := rapid.SampledFrom([]bool{true, false}).Draw(t, "negative")
			var styleVal int
			if negative {
				styleVal = rapid.IntRange(-100, -1).Draw(t, "style")
			} else {
				styleVal = rapid.IntRange(5, 100).Draw(t, "style")
			}

			orientation := Orientation(rapid.IntRange(0, 1).Draw(t, "orientation"))
			width := rapid.IntRange(1, 50).Draw(t, "width")
			height := rapid.IntRange(1, 50).Draw(t, "height")
			bounds := image.Rect(0, 0, width, height)

			result := Render(Config{
				Style:       Style(styleVal),
				Orientation: orientation,
				Value:       rapid.Float64Range(0.0, 1.0).Draw(t, "value"),
				Bounds:      bounds,
				Foreground:  color.RGBA{R: 255, G: 0, B: 0, A: 255},
				Background:  color.RGBA{R: 0, G: 0, B: 255, A: 255},
			})

			if result != nil {
				t.Fatalf("expected nil for unrecognized style=%d (w=%d, h=%d, orient=%d), got non-nil",
					styleVal, width, height, orientation)
			}
		})
	})
}

// --- From: progressbar_test.go ---

// TestDefaultColors verifies that zero-value Foreground/Background in Config
// produce opaque white foreground and opaque black background.

func TestDefaultColors(t *testing.T) {
	bounds := image.Rect(0, 0, 10, 5)

	t.Run("full_value_all_foreground_white", func(t *testing.T) {
		cfg := Config{
			Style:  Linear,
			Value:  1.0,
			Bounds: bounds,
			// Zero-value Foreground and Background → defaults apply
		}
		result := Render(cfg)
		if result == nil {
			t.Fatal("expected non-nil result")
		}

		white := color.RGBA{R: 255, G: 255, B: 255, A: 255}
		for y := 0; y < 5; y++ {
			for x := 0; x < 10; x++ {
				got := result.Image.At(x, y)
				r, g, b, a := got.RGBA()
				wr, wg, wb, wa := white.RGBA()
				if r != wr || g != wg || b != wb || a != wa {
					t.Fatalf("pixel (%d,%d): expected white %v, got %v", x, y, white, got)
				}
			}
		}
	})

	t.Run("zero_value_all_background_black", func(t *testing.T) {
		cfg := Config{
			Style:  Linear,
			Value:  0.0,
			Bounds: bounds,
		}
		result := Render(cfg)
		if result == nil {
			t.Fatal("expected non-nil result")
		}

		black := color.RGBA{R: 0, G: 0, B: 0, A: 255}
		for y := 0; y < 5; y++ {
			for x := 0; x < 10; x++ {
				got := result.Image.At(x, y)
				r, g, b, a := got.RGBA()
				br, bg, bb, ba := black.RGBA()
				if r != br || g != bg || b != bb || a != ba {
					t.Fatalf("pixel (%d,%d): expected black %v, got %v", x, y, black, got)
				}
			}
		}
	})
}

// TestCustomColors verifies that specified RGBA foreground/background colors
// appear in the correct regions of the rendered output.

func TestCustomColors(t *testing.T) {
	fg := color.RGBA{R: 100, G: 150, B: 200, A: 255}
	bg := color.RGBA{R: 50, G: 60, B: 70, A: 255}

	cfg := Config{
		Style:      Linear,
		Value:      0.5,
		Bounds:     image.Rect(0, 0, 10, 5),
		Foreground: fg,
		Background: bg,
	}
	result := Render(cfg)
	if result == nil {
		t.Fatal("expected non-nil result")
	}

	// With value 0.5 and width 10: floor(10*0.5) = 5 columns filled
	fillCols := 5

	for y := 0; y < 5; y++ {
		for x := 0; x < 10; x++ {
			got := result.Image.At(x, y)
			r, g, b, a := got.RGBA()

			if x < fillCols {
				// Should be foreground
				er, eg, eb, ea := fg.RGBA()
				if r != er || g != eg || b != eb || a != ea {
					t.Fatalf("pixel (%d,%d): expected fg %v, got %v", x, y, fg, got)
				}
			} else {
				// Should be background
				er, eg, eb, ea := bg.RGBA()
				if r != er || g != eg || b != eb || a != ea {
					t.Fatalf("pixel (%d,%d): expected bg %v, got %v", x, y, bg, got)
				}
			}
		}
	}
}

func TestRoundedCapsThickBorderLeavesVisibleInterior(t *testing.T) {
	fg := color.RGBA{R: 255, G: 0, B: 0, A: 255}
	bg := color.RGBA{R: 0, G: 0, B: 255, A: 255}
	border := color.RGBA{R: 255, G: 255, B: 255, A: 255}

	cfg := Config{
		Style:       Linear,
		Orientation: OrientHorizontal,
		Value:       1.0,
		Bounds:      image.Rect(0, 0, 60, 20),
		Foreground:  fg,
		Background:  bg,
		RoundedCaps: true,
		BorderWidth: 2,
		BorderWall:  6,
		BorderColor: border,
	}

	result := Render(cfg)
	if result == nil {
		t.Fatal("expected non-nil result for thick rounded border")
	}

	outsideBorderFound := false
	centerPixelIsFill := false
	for y := 0; y < 20; y++ {
		for x := 0; x < 60; x++ {
			pr, pg, pb, pa := result.Image.At(x, y).RGBA()
			if pr == 0xFFFF && pg == 0xFFFF && pb == 0xFFFF && pa == 0xFFFF {
				outsideBorderFound = true
			}
			if x == 30 && y == 10 {
				if pr == 0xFFFF && pg == 0 && pb == 0 && pa == 0xFFFF {
					centerPixelIsFill = true
				}
			}
		}
	}

	if !outsideBorderFound {
		t.Fatal("expected the border wall to draw visible border pixels")
	}
	if !centerPixelIsFill {
		t.Fatal("expected the thick wall to leave a visible red interior region in the center")
	}
}

// TestPieCentering verifies that non-square bounds produce a centered inscribed circle.

func TestPieCentering(t *testing.T) {
	fg := color.RGBA{R: 200, G: 100, B: 50, A: 255}
	bg := color.RGBA{R: 10, G: 20, B: 30, A: 255}

	// Non-square bounds: width=20, height=10
	// Inscribed circle: radius = min(20,10)/2 = 5, centered at (10, 5)
	cfg := Config{
		Style:      Pie,
		Value:      1.0, // Full circle
		Bounds:     image.Rect(0, 0, 20, 10),
		Foreground: fg,
		Background: bg,
	}
	result := Render(cfg)
	if result == nil {
		t.Fatal("expected non-nil result")
	}

	// Helper to check pixel color
	isFg := func(x, y int) bool {
		got := result.Image.At(x, y)
		r, g, b, a := got.RGBA()
		er, eg, eb, ea := fg.RGBA()
		return r == er && g == eg && b == eb && a == ea
	}
	isBg := func(x, y int) bool {
		got := result.Image.At(x, y)
		r, g, b, a := got.RGBA()
		er, eg, eb, ea := bg.RGBA()
		return r == er && g == eg && b == eb && a == ea
	}

	// Center pixel (10, 5) should be within the circle → foreground
	if !isFg(10, 5) {
		t.Errorf("center pixel (10,5): expected foreground color")
	}

	// Corners should be outside the circle → background
	corners := [][2]int{{0, 0}, {19, 0}, {0, 9}, {19, 9}}
	for _, c := range corners {
		if !isBg(c[0], c[1]) {
			t.Errorf("corner pixel (%d,%d): expected background color", c[0], c[1])
		}
	}

	// Pixel at (10, 0) is at the top center of the circle boundary
	// Distance from center (10, 5): pixel center is (10.5, 0.5), center is (10, 5)
	// dx = 10.5-10 = 0.5, dy = 0.5-5 = -4.5, dist = sqrt(0.25+20.25) ≈ 4.53 < 5
	// So it should be inside the circle → foreground
	if !isFg(10, 0) {
		t.Errorf("top-center pixel (10,0): expected foreground color (inside circle)")
	}
}

// TestLabelStrings verifies that each style produces the correct label string.

func TestLabelStrings(t *testing.T) {
	tests := []struct {
		style       Style
		orientation Orientation
		label       string
	}{
		{Linear, OrientHorizontal, "progressbar/linear"},
		{Linear, OrientVertical, "progressbar/linear"},
		{Pie, OrientHorizontal, "progressbar/pie"},
		{Segmented, OrientHorizontal, "progressbar/segmented"},
		{Ring, OrientHorizontal, "progressbar/ring"},
		{Arc, OrientHorizontal, "progressbar/arc"},
	}

	for _, tt := range tests {
		t.Run(tt.label, func(t *testing.T) {
			cfg := Config{
				Style:       tt.style,
				Orientation: tt.orientation,
				Value:       0.5,
				Bounds:      image.Rect(0, 0, 10, 10),
			}
			result := Render(cfg)
			if result == nil {
				t.Fatal("expected non-nil result")
			}
			if result.Label != tt.label {
				t.Errorf("expected label %q, got %q", tt.label, result.Label)
			}
		})
	}
}

// TestProperty24_DefaultOrientationEquivalence verifies Property 24: For any valid Config
// where the Orientation field is the zero value (OrientHorizontal), the rendered output
// SHALL be pixel-identical to the same Config with Orientation explicitly set to
// OrientHorizontal. Equivalently: omitting the Orientation field always produces
// horizontal-oriented output.

func TestProperty24_DefaultOrientationEquivalence(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		style := Style(rapid.IntRange(0, 4).Draw(t, "style"))

		var w, h int
		switch style {
		case Pie, Ring, Arc:
			w = rapid.IntRange(3, 50).Draw(t, "width")
			h = rapid.IntRange(3, 50).Draw(t, "height")
		default:
			w = rapid.IntRange(1, 50).Draw(t, "width")
			h = rapid.IntRange(1, 50).Draw(t, "height")
		}
		bounds := image.Rect(0, 0, w, h)

		value := rapid.Float64Range(0.0, 1.0).Draw(t, "value")
		fg := color.RGBA{
			R: uint8(rapid.IntRange(0, 255).Draw(t, "fgR")),
			G: uint8(rapid.IntRange(0, 255).Draw(t, "fgG")),
			B: uint8(rapid.IntRange(0, 255).Draw(t, "fgB")),
			A: 255,
		}
		bg := color.RGBA{
			R: uint8(rapid.IntRange(0, 255).Draw(t, "bgR")),
			G: uint8(rapid.IntRange(0, 255).Draw(t, "bgG")),
			B: uint8(rapid.IntRange(0, 255).Draw(t, "bgB")),
			A: 255,
		}

		// Config with zero-value Orientation (omitted = OrientHorizontal)
		cfgZero := Config{
			Style:      style,
			Value:      value,
			Bounds:     bounds,
			Foreground: fg,
			Background: bg,
			// Orientation is zero value (OrientHorizontal)
		}

		// Config with Orientation explicitly set to OrientHorizontal
		cfgExplicit := Config{
			Style:       style,
			Orientation: OrientHorizontal,
			Value:       value,
			Bounds:      bounds,
			Foreground:  fg,
			Background:  bg,
		}

		result1 := Render(cfgZero)
		result2 := Render(cfgExplicit)

		if result1 == nil {
			t.Fatal("expected non-nil result for zero-orientation config")
		}
		if result2 == nil {
			t.Fatal("expected non-nil result for explicit OrientHorizontal config")
		}

		// Assert pixel-identical output
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				r1, g1, b1, a1 := result1.Image.At(x, y).RGBA()
				r2, g2, b2, a2 := result2.Image.At(x, y).RGBA()
				if r1 != r2 || g1 != g2 || b1 != b2 || a1 != a2 {
					t.Fatalf("pixel mismatch at (%d,%d): zero-orient=(%d,%d,%d,%d), explicit=(%d,%d,%d,%d) [style=%d, value=%f]",
						x, y, r1, g1, b1, a1, r2, g2, b2, a2, style, value)
				}
			}
		}
	})
}

// TestProperty1_HorizontalGradientFillMapsStopsToFillRegion verifies Property 1:
// For any Linear bar with Horizontal orientation and a valid GradientFill (≥2 stops),
// any value in [0.0, 1.0], and any valid bounds, each pixel in the fill region SHALL
// have a color equal to the linearly interpolated gradient color at position
// (pixel_x / fill_width), within ±1 per channel for rounding.

func TestProperty1_HorizontalGradientFillMapsStopsToFillRegion(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate valid bounds (width ≥ 2 to have a meaningful fill region with gradient)
		width := rapid.IntRange(2, 100).Draw(t, "width")
		height := rapid.IntRange(1, 50).Draw(t, "height")
		bounds := image.Rect(0, 0, width, height)

		// Generate a value that produces at least 2 fill columns for gradient interpolation
		// fillCols = floor(width * value), need fillCols >= 2
		// value must be > 1/width to get at least 2 fill cols
		minValue := 2.0 / float64(width)
		if minValue > 1.0 {
			minValue = 1.0
		}
		value := rapid.Float64Range(minValue, 1.0).Draw(t, "value")

		// Generate gradient stops (2 to 8 stops)
		numStops := rapid.IntRange(2, 8).Draw(t, "numStops")
		stops := make([]gradient.ColorStop, numStops)
		for i := 0; i < numStops; i++ {
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

		cfg := Config{
			Style:       Linear,
			Orientation: OrientHorizontal,
			Value:       value,
			Bounds:      bounds,
			Foreground:  color.RGBA{R: 255, G: 255, B: 255, A: 255},
			Background:  color.RGBA{R: 0, G: 0, B: 0, A: 255},
			Gradient:    &GradientFill{Stops: stops},
		}

		result := Render(cfg)
		if result == nil {
			t.Fatal("expected non-nil result for valid config")
		}

		fillCols := int(float64(width) * value)
		if fillCols < 2 {
			// Should not happen given our value constraint, but guard
			return
		}

		// For each pixel in the fill region, verify it matches the expected
		// interpolated gradient color within ±1 per channel.
		for x := 0; x < fillCols; x++ {
			// Compute the gradient t at this pixel position
			var gradT float64
			if fillCols > 1 {
				gradT = float64(x) / float64(fillCols-1)
			}

			// Compute expected color using interpolateGradient (same function used by render)
			expected := interpolateGradient(stops, gradT)

			// Check all rows (gradient is uniform across rows for horizontal bars)
			for y := 0; y < height; y++ {
				r, g, b, a := result.Image.At(x, y).RGBA()
				// Convert from 16-bit back to 8-bit
				gotR := uint8(r >> 8)
				gotG := uint8(g >> 8)
				gotB := uint8(b >> 8)
				gotA := uint8(a >> 8)

				if !withinTolerance(gotR, expected.R, 1) ||
					!withinTolerance(gotG, expected.G, 1) ||
					!withinTolerance(gotB, expected.B, 1) ||
					!withinTolerance(gotA, expected.A, 1) {
					t.Fatalf("pixel (%d,%d): got RGBA(%d,%d,%d,%d), expected RGBA(%d,%d,%d,%d) ±1 "+
						"[value=%.6f, width=%d, fillCols=%d, gradT=%.6f, numStops=%d]",
						x, y, gotR, gotG, gotB, gotA,
						expected.R, expected.G, expected.B, expected.A,
						value, width, fillCols, gradT, numStops)
				}
			}
		}
	})
}

// TestProperty2_VerticalGradientFillInterpolatesBottomToTop verifies Property 2:
// For any Linear bar with Vertical orientation and a valid GradientFill (≥2 stops),
// any value in [0.0, 1.0], and any valid bounds, each pixel in the fill region SHALL
// have a color equal to the linearly interpolated gradient color at position
// ((fill_bottom - pixel_y) / fill_height), within ±1 per channel for rounding.
//
// The fill region for vertical bar: pixels where y >= h - fillRows, where fillRows = floor(height * value).
// Gradient t at pixel y in fill region: t = float64((h-1) - y) / float64(fillRows - 1) when fillRows > 1.
// This means t=0 at bottom (y=h-1), t=1 at top of fill region (y=h-fillRows).

func TestProperty2_VerticalGradientFillInterpolatesBottomToTop(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate valid bounds (height ≥ 2 to have a meaningful fill region with gradient)
		width := rapid.IntRange(1, 50).Draw(t, "width")
		height := rapid.IntRange(2, 100).Draw(t, "height")
		bounds := image.Rect(0, 0, width, height)

		// Generate a value that produces at least 2 fill rows for gradient interpolation
		// fillRows = floor(height * value), need fillRows >= 2
		// value must be > 1/height to get at least 2 fill rows
		minValue := 2.0 / float64(height)
		if minValue > 1.0 {
			minValue = 1.0
		}
		value := rapid.Float64Range(minValue, 1.0).Draw(t, "value")

		// Generate gradient stops (2 to 8 stops)
		numStops := rapid.IntRange(2, 8).Draw(t, "numStops")
		stops := make([]gradient.ColorStop, numStops)
		for i := 0; i < numStops; i++ {
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

		cfg := Config{
			Style:       Linear,
			Orientation: OrientVertical,
			Value:       value,
			Bounds:      bounds,
			Foreground:  color.RGBA{R: 255, G: 255, B: 255, A: 255},
			Background:  color.RGBA{R: 0, G: 0, B: 0, A: 255},
			Gradient:    &GradientFill{Stops: stops},
		}

		result := Render(cfg)
		if result == nil {
			t.Fatal("expected non-nil result for valid config")
		}

		fillRows := int(float64(height) * value)
		if fillRows < 2 {
			// Should not happen given our value constraint, but guard
			return
		}

		// For each pixel in the fill region, verify it matches the expected
		// interpolated gradient color within ±1 per channel.
		// Fill region: y >= h - fillRows (bottom portion of the image)
		for y := height - fillRows; y < height; y++ {
			// Compute the gradient t at this pixel position
			// t = 0.0 at bottom (y = h-1), t = 1.0 at top of fill region (y = h-fillRows)
			var gradT float64
			if fillRows > 1 {
				gradT = float64((height-1)-y) / float64(fillRows-1)
			}

			// Compute expected color using interpolateGradient (same function used by render)
			expected := interpolateGradient(stops, gradT)

			// Check all columns (gradient is uniform across columns for vertical bars)
			for x := 0; x < width; x++ {
				r, g, b, a := result.Image.At(x, y).RGBA()
				// Convert from 16-bit back to 8-bit
				gotR := uint8(r >> 8)
				gotG := uint8(g >> 8)
				gotB := uint8(b >> 8)
				gotA := uint8(a >> 8)

				if !withinTolerance(gotR, expected.R, 1) ||
					!withinTolerance(gotG, expected.G, 1) ||
					!withinTolerance(gotB, expected.B, 1) ||
					!withinTolerance(gotA, expected.A, 1) {
					t.Fatalf("pixel (%d,%d): got RGBA(%d,%d,%d,%d), expected RGBA(%d,%d,%d,%d) ±1 "+
						"[value=%.6f, height=%d, fillRows=%d, gradT=%.6f, numStops=%d]",
						x, y, gotR, gotG, gotB, gotA,
						expected.R, expected.G, expected.B, expected.A,
						value, height, fillRows, gradT, numStops)
				}
			}
		}
	})
}

// TestProperty3_GradientFallbackToSolidFill verifies Property 3: For any bar
// configuration where the GradientFill has fewer than 2 stops OR the style is Pie,
// the rendered output SHALL be pixel-identical to the same configuration with no
// gradient (nil GradientFill).

func TestProperty3_GradientFallbackToSolidFill(t *testing.T) {
	t.Run("fewer_than_2_stops", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			// Generate a random style (all 5 styles).
			style := Style(rapid.IntRange(0, 4).Draw(t, "style"))
			orientation := Orientation(rapid.IntRange(0, 1).Draw(t, "orientation"))

			// Generate bounds appropriate for the style.
			var w, h int
			switch style {
			case Pie, Ring, Arc:
				w = rapid.IntRange(3, 50).Draw(t, "width")
				h = rapid.IntRange(3, 50).Draw(t, "height")
			default:
				w = rapid.IntRange(1, 50).Draw(t, "width")
				h = rapid.IntRange(1, 50).Draw(t, "height")
			}
			bounds := image.Rect(0, 0, w, h)

			value := rapid.Float64Range(0.0, 1.0).Draw(t, "value")

			fg := color.RGBA{
				R: uint8(rapid.IntRange(0, 255).Draw(t, "fgR")),
				G: uint8(rapid.IntRange(0, 255).Draw(t, "fgG")),
				B: uint8(rapid.IntRange(0, 255).Draw(t, "fgB")),
				A: 255,
			}
			bg := color.RGBA{
				R: uint8(rapid.IntRange(0, 255).Draw(t, "bgR")),
				G: uint8(rapid.IntRange(0, 255).Draw(t, "bgG")),
				B: uint8(rapid.IntRange(0, 255).Draw(t, "bgB")),
				A: 255,
			}

			// Generate a gradient with 0 or 1 stops (fewer than 2).
			numStops := rapid.IntRange(0, 1).Draw(t, "numStops")
			var grad *GradientFill
			if numStops == 0 {
				grad = &GradientFill{Stops: []gradient.ColorStop{}}
			} else {
				grad = &GradientFill{Stops: []gradient.ColorStop{
					{
						Position: rapid.Float64Range(0.0, 1.0).Draw(t, "stopPos"),
						Color: color.RGBA{
							R: uint8(rapid.IntRange(0, 255).Draw(t, "stopR")),
							G: uint8(rapid.IntRange(0, 255).Draw(t, "stopG")),
							B: uint8(rapid.IntRange(0, 255).Draw(t, "stopB")),
							A: uint8(rapid.IntRange(0, 255).Draw(t, "stopA")),
						},
					},
				}}
			}

			// Render with the invalid gradient.
			cfgWithGrad := Config{
				Style:       style,
				Orientation: orientation,
				Value:       value,
				Bounds:      bounds,
				Foreground:  fg,
				Background:  bg,
				Gradient:    grad,
			}

			// Render without gradient (nil).
			cfgNoGrad := Config{
				Style:       style,
				Orientation: orientation,
				Value:       value,
				Bounds:      bounds,
				Foreground:  fg,
				Background:  bg,
				Gradient:    nil,
			}

			result1 := Render(cfgWithGrad)
			result2 := Render(cfgNoGrad)

			if result1 == nil {
				t.Fatal("expected non-nil result for config with invalid gradient")
			}
			if result2 == nil {
				t.Fatal("expected non-nil result for config without gradient")
			}

			// Assert pixel-identical output.
			for y := 0; y < h; y++ {
				for x := 0; x < w; x++ {
					r1, g1, b1, a1 := result1.Image.At(x, y).RGBA()
					r2, g2, b2, a2 := result2.Image.At(x, y).RGBA()
					if r1 != r2 || g1 != g2 || b1 != b2 || a1 != a2 {
						t.Fatalf("pixel mismatch at (%d,%d): withGrad=(%d,%d,%d,%d), noGrad=(%d,%d,%d,%d) "+
							"[style=%d, orient=%d, value=%.4f, numStops=%d, bounds=%dx%d]",
							x, y, r1, g1, b1, a1, r2, g2, b2, a2,
							style, orientation, value, numStops, w, h)
					}
				}
			}
		})
	})

	t.Run("pie_ignores_gradient", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			orientation := Orientation(rapid.IntRange(0, 1).Draw(t, "orientation"))

			// Pie style requires bounds ≥ 3×3.
			w := rapid.IntRange(3, 50).Draw(t, "width")
			h := rapid.IntRange(3, 50).Draw(t, "height")
			bounds := image.Rect(0, 0, w, h)

			value := rapid.Float64Range(0.0, 1.0).Draw(t, "value")

			fg := color.RGBA{
				R: uint8(rapid.IntRange(0, 255).Draw(t, "fgR")),
				G: uint8(rapid.IntRange(0, 255).Draw(t, "fgG")),
				B: uint8(rapid.IntRange(0, 255).Draw(t, "fgB")),
				A: 255,
			}
			bg := color.RGBA{
				R: uint8(rapid.IntRange(0, 255).Draw(t, "bgR")),
				G: uint8(rapid.IntRange(0, 255).Draw(t, "bgG")),
				B: uint8(rapid.IntRange(0, 255).Draw(t, "bgB")),
				A: 255,
			}

			// Generate a valid gradient (≥2 stops) that Pie should ignore.
			numStops := rapid.IntRange(2, 8).Draw(t, "numStops")
			stops := make([]gradient.ColorStop, numStops)
			for i := 0; i < numStops; i++ {
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

			// Render Pie with a valid gradient.
			cfgWithGrad := Config{
				Style:       Pie,
				Orientation: orientation,
				Value:       value,
				Bounds:      bounds,
				Foreground:  fg,
				Background:  bg,
				Gradient:    &GradientFill{Stops: stops},
			}

			// Render Pie without gradient (nil).
			cfgNoGrad := Config{
				Style:       Pie,
				Orientation: orientation,
				Value:       value,
				Bounds:      bounds,
				Foreground:  fg,
				Background:  bg,
				Gradient:    nil,
			}

			result1 := Render(cfgWithGrad)
			result2 := Render(cfgNoGrad)

			if result1 == nil {
				t.Fatal("expected non-nil result for Pie with gradient")
			}
			if result2 == nil {
				t.Fatal("expected non-nil result for Pie without gradient")
			}

			// Assert pixel-identical output.
			for y := 0; y < h; y++ {
				for x := 0; x < w; x++ {
					r1, g1, b1, a1 := result1.Image.At(x, y).RGBA()
					r2, g2, b2, a2 := result2.Image.At(x, y).RGBA()
					if r1 != r2 || g1 != g2 || b1 != b2 || a1 != a2 {
						t.Fatalf("pixel mismatch at (%d,%d): pieWithGrad=(%d,%d,%d,%d), pieNoGrad=(%d,%d,%d,%d) "+
							"[orient=%d, value=%.4f, numStops=%d, bounds=%dx%d]",
							x, y, r1, g1, b1, a1, r2, g2, b2, a2,
							orientation, value, numStops, w, h)
					}
				}
			}
		})
	})
}

// withinTolerance checks if two uint8 values are within ±tol of each other.
func withinTolerance(got, expected uint8, tol uint8) bool {
	if got > expected {
		return (got - expected) <= tol
	}
	return (expected - got) <= tol
}

// TestProperty4_SegmentedBarGapCorrectness verifies Property 4:
// For any segmented bar configuration with valid bounds, the pixels at computed gap
// positions between adjacent cells SHALL all be the background color, and each gap
// SHALL be exactly the configured gap width (clamped to [1, 4]) in pixels.

func TestProperty4_SegmentedBarGapCorrectness(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		orientation := Orientation(rapid.IntRange(0, 1).Draw(t, "orientation"))

		// Ensure bounds are large enough for at least 2 segments.
		// For horizontal: width is primary axis; for vertical: height is primary axis.
		var primaryLen, minorLen int
		primaryLen = rapid.IntRange(20, 100).Draw(t, "primaryLen")
		minorLen = rapid.IntRange(2, 30).Draw(t, "minorLen")

		var bounds image.Rectangle
		if orientation == OrientHorizontal {
			bounds = image.Rect(0, 0, primaryLen, minorLen)
		} else {
			bounds = image.Rect(0, 0, minorLen, primaryLen)
		}

		// Configure gap (will be clamped to [1, 4])
		gapRaw := rapid.IntRange(0, 6).Draw(t, "gapRaw")
		gap := gapRaw
		if gap < 1 {
			gap = 1
		}
		if gap > 4 {
			gap = 4
		}

		// Generate segment count that yields at least 2 cells with valid cellWidth
		segCount := rapid.IntRange(2, 10).Draw(t, "segCount")

		// Compute expected cell width
		cellWidth := (primaryLen - (segCount-1)*gap) / segCount
		if cellWidth < 1 {
			// Skip if can't fit segments properly
			return
		}

		// Verify the bar won't fall back to unsegmented
		if primaryLen < 2*cellWidth+gap {
			return
		}

		value := rapid.Float64Range(0.0, 1.0).Draw(t, "value")

		fg := color.RGBA{R: 255, G: 0, B: 0, A: 255}
		bg := color.RGBA{R: 0, G: 0, B: 255, A: 255}

		cfg := Config{
			Style:        Segmented,
			Orientation:  orientation,
			Value:        value,
			Bounds:       bounds,
			Foreground:   fg,
			Background:   bg,
			SegmentCount: segCount,
			SegmentGap:   gapRaw,
		}

		result := Render(cfg)
		if result == nil {
			t.Fatal("expected non-nil result for valid segmented bar config")
		}

		// Verify gap pixels between adjacent cells.
		// Cell i starts at: i * (cellWidth + gap)
		// Cell i ends at: i * (cellWidth + gap) + cellWidth
		// Gap between cell i and cell i+1: [cellEnd, cellEnd + gap)
		for i := 0; i < segCount-1; i++ {
			cellEnd := i*(cellWidth+gap) + cellWidth
			for gPos := cellEnd; gPos < cellEnd+gap; gPos++ {
				// Check all pixels along the minor axis at this gap position
				for m := 0; m < minorLen; m++ {
					var x, y int
					if orientation == OrientHorizontal {
						x = gPos
						y = m
					} else {
						// Vertical: primary axis is Y (bottom-to-top)
						// Position 0 along primary axis maps to pixel y = height-1
						x = m
						y = bounds.Dy() - 1 - gPos
					}

					if x < 0 || x >= bounds.Dx() || y < 0 || y >= bounds.Dy() {
						continue
					}

					r, g, b, a := result.Image.At(x, y).RGBA()
					bgR, bgG, bgB, bgA := bg.RGBA()
					if r != bgR || g != bgG || b != bgB || a != bgA {
						t.Fatalf("gap pixel at primary=%d, minor=%d (image x=%d, y=%d): "+
							"expected bg color, got (%d,%d,%d,%d) "+
							"[orient=%d, segCount=%d, gap=%d, cellWidth=%d, cellEnd=%d]",
							gPos, m, x, y, r>>8, g>>8, b>>8, a>>8,
							orientation, segCount, gap, cellWidth, cellEnd)
					}
				}
			}
		}
	})
}

// TestProperty5_SegmentedCellFillAllOrNothing verifies Property 5:
// Each cell is either fully colored or fully background-colored. A cell is filled
// iff its center position along the primary axis falls within the fill region.
// Pixels beyond the last segment are background-colored.

func TestProperty5_SegmentedCellFillAllOrNothing(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		orientation := Orientation(rapid.IntRange(0, 1).Draw(t, "orientation"))

		primaryLen := rapid.IntRange(20, 100).Draw(t, "primaryLen")
		minorLen := rapid.IntRange(2, 30).Draw(t, "minorLen")

		var bounds image.Rectangle
		if orientation == OrientHorizontal {
			bounds = image.Rect(0, 0, primaryLen, minorLen)
		} else {
			bounds = image.Rect(0, 0, minorLen, primaryLen)
		}

		gap := rapid.IntRange(1, 4).Draw(t, "gap")
		segCount := rapid.IntRange(2, 10).Draw(t, "segCount")

		cellWidth := (primaryLen - (segCount-1)*gap) / segCount
		if cellWidth < 1 {
			return
		}
		if primaryLen < 2*cellWidth+gap {
			return
		}

		value := rapid.Float64Range(0.0, 1.0).Draw(t, "value")

		fg := color.RGBA{R: 255, G: 0, B: 0, A: 255}
		bg := color.RGBA{R: 0, G: 0, B: 255, A: 255}

		cfg := Config{
			Style:        Segmented,
			Orientation:  orientation,
			Value:        value,
			Bounds:       bounds,
			Foreground:   fg,
			Background:   bg,
			SegmentCount: segCount,
			SegmentGap:   gap,
		}

		result := Render(cfg)
		if result == nil {
			t.Fatal("expected non-nil result")
		}

		fillExtent := int(float64(primaryLen) * value)

		for i := 0; i < segCount; i++ {
			cellStart := i * (cellWidth + gap)
			cellEnd := cellStart + cellWidth
			cellCenter := cellStart + cellWidth/2

			isFilled := cellCenter < fillExtent

			// Check all pixels within this cell
			for pos := cellStart; pos < cellEnd; pos++ {
				for m := 0; m < minorLen; m++ {
					var x, y int
					if orientation == OrientHorizontal {
						x = pos
						y = m
					} else {
						x = m
						y = bounds.Dy() - 1 - pos
					}

					if x < 0 || x >= bounds.Dx() || y < 0 || y >= bounds.Dy() {
						continue
					}

					r, g, b, a := result.Image.At(x, y).RGBA()
					bgR, bgG, bgB, bgA := bg.RGBA()
					fgR, fgG, fgB, fgA := fg.RGBA()

					if isFilled {
						if r != fgR || g != fgG || b != fgB || a != fgA {
							t.Fatalf("cell %d (filled) pixel at primary=%d, minor=%d (x=%d, y=%d): "+
								"expected fg, got (%d,%d,%d,%d) "+
								"[orient=%d, value=%.4f, fillExtent=%d, cellCenter=%d]",
								i, pos, m, x, y, r>>8, g>>8, b>>8, a>>8,
								orientation, value, fillExtent, cellCenter)
						}
					} else {
						if r != bgR || g != bgG || b != bgB || a != bgA {
							t.Fatalf("cell %d (unfilled) pixel at primary=%d, minor=%d (x=%d, y=%d): "+
								"expected bg, got (%d,%d,%d,%d) "+
								"[orient=%d, value=%.4f, fillExtent=%d, cellCenter=%d]",
								i, pos, m, x, y, r>>8, g>>8, b>>8, a>>8,
								orientation, value, fillExtent, cellCenter)
						}
					}
				}
			}
		}

		// Verify pixels beyond the last segment are background-colored
		lastCellEnd := (segCount-1)*(cellWidth+gap) + cellWidth
		for pos := lastCellEnd; pos < primaryLen; pos++ {
			for m := 0; m < minorLen; m++ {
				var x, y int
				if orientation == OrientHorizontal {
					x = pos
					y = m
				} else {
					x = m
					y = bounds.Dy() - 1 - pos
				}

				if x < 0 || x >= bounds.Dx() || y < 0 || y >= bounds.Dy() {
					continue
				}

				r, g, b, a := result.Image.At(x, y).RGBA()
				bgR, bgG, bgB, bgA := bg.RGBA()
				if r != bgR || g != bgG || b != bgB || a != bgA {
					t.Fatalf("remainder pixel at primary=%d, minor=%d (x=%d, y=%d): "+
						"expected bg, got (%d,%d,%d,%d) [orient=%d]",
						pos, m, x, y, r>>8, g>>8, b>>8, a>>8, orientation)
				}
			}
		}
	})
}

// TestProperty6_SegmentedGradientAtCellCenter verifies Property 6:
// For any segmented bar combined with a valid GradientFill, each filled cell SHALL
// have pixels colored according to the gradient interpolated at the cell's center
// position along the bar axis, within ±1 per channel for rounding.

func TestProperty6_SegmentedGradientAtCellCenter(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		orientation := Orientation(rapid.IntRange(0, 1).Draw(t, "orientation"))

		primaryLen := rapid.IntRange(20, 100).Draw(t, "primaryLen")
		minorLen := rapid.IntRange(2, 20).Draw(t, "minorLen")

		var bounds image.Rectangle
		if orientation == OrientHorizontal {
			bounds = image.Rect(0, 0, primaryLen, minorLen)
		} else {
			bounds = image.Rect(0, 0, minorLen, primaryLen)
		}

		gap := rapid.IntRange(1, 4).Draw(t, "gap")
		segCount := rapid.IntRange(2, 8).Draw(t, "segCount")

		cellWidth := (primaryLen - (segCount-1)*gap) / segCount
		if cellWidth < 1 {
			return
		}
		if primaryLen < 2*cellWidth+gap {
			return
		}

		// Use a value large enough to fill at least one cell
		value := rapid.Float64Range(0.2, 1.0).Draw(t, "value")

		// Generate gradient stops (2 to 6 stops)
		numStops := rapid.IntRange(2, 6).Draw(t, "numStops")
		stops := make([]gradient.ColorStop, numStops)
		for i := 0; i < numStops; i++ {
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

		fg := color.RGBA{R: 255, G: 255, B: 255, A: 255}
		bg := color.RGBA{R: 0, G: 0, B: 0, A: 255}

		cfg := Config{
			Style:        Segmented,
			Orientation:  orientation,
			Value:        value,
			Bounds:       bounds,
			Foreground:   fg,
			Background:   bg,
			SegmentCount: segCount,
			SegmentGap:   gap,
			Gradient:     &GradientFill{Stops: stops},
		}

		result := Render(cfg)
		if result == nil {
			t.Fatal("expected non-nil result")
		}

		fillExtent := int(float64(primaryLen) * value)

		for i := 0; i < segCount; i++ {
			cellStart := i * (cellWidth + gap)
			cellEnd := cellStart + cellWidth
			cellCenter := cellStart + cellWidth/2

			if cellCenter >= fillExtent {
				continue // Cell is not filled
			}

			// Compute the expected gradient color at cell center
			var gradT float64
			if primaryLen > 1 {
				gradT = float64(cellCenter) / float64(primaryLen-1)
			}
			expected := interpolateGradient(stops, gradT)

			// Check all pixels within this filled cell
			for pos := cellStart; pos < cellEnd; pos++ {
				for m := 0; m < minorLen; m++ {
					var x, y int
					if orientation == OrientHorizontal {
						x = pos
						y = m
					} else {
						x = m
						y = bounds.Dy() - 1 - pos
					}

					if x < 0 || x >= bounds.Dx() || y < 0 || y >= bounds.Dy() {
						continue
					}

					r, g, b, a := result.Image.At(x, y).RGBA()
					gotR := uint8(r >> 8)
					gotG := uint8(g >> 8)
					gotB := uint8(b >> 8)
					gotA := uint8(a >> 8)

					if !withinTolerance(gotR, expected.R, 1) ||
						!withinTolerance(gotG, expected.G, 1) ||
						!withinTolerance(gotB, expected.B, 1) ||
						!withinTolerance(gotA, expected.A, 1) {
						t.Fatalf("cell %d pixel at primary=%d, minor=%d (x=%d, y=%d): "+
							"got RGBA(%d,%d,%d,%d), expected RGBA(%d,%d,%d,%d) ±1 "+
							"[orient=%d, value=%.4f, cellCenter=%d, gradT=%.4f]",
							i, pos, m, x, y,
							gotR, gotG, gotB, gotA,
							expected.R, expected.G, expected.B, expected.A,
							orientation, value, cellCenter, gradT)
					}
				}
			}
		}
	})
}

// TestProperty21_SegmentedOrientationFillDirection verifies Property 21 (segmented portion):
// For Horizontal, filled cells are contiguous from the left. For Vertical, filled cells
// are contiguous from the bottom. No filled cell shall appear beyond the fill boundary
// while an unfilled cell exists before it along the configured direction.

func TestProperty21_SegmentedOrientationFillDirection(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		orientation := Orientation(rapid.IntRange(0, 1).Draw(t, "orientation"))

		primaryLen := rapid.IntRange(20, 100).Draw(t, "primaryLen")
		minorLen := rapid.IntRange(2, 30).Draw(t, "minorLen")

		var bounds image.Rectangle
		if orientation == OrientHorizontal {
			bounds = image.Rect(0, 0, primaryLen, minorLen)
		} else {
			bounds = image.Rect(0, 0, minorLen, primaryLen)
		}

		gap := rapid.IntRange(1, 4).Draw(t, "gap")
		segCount := rapid.IntRange(2, 10).Draw(t, "segCount")

		cellWidth := (primaryLen - (segCount-1)*gap) / segCount
		if cellWidth < 1 {
			return
		}
		if primaryLen < 2*cellWidth+gap {
			return
		}

		// Value between 0 and 1 exclusive to have a mix of filled/unfilled cells
		value := rapid.Float64Range(0.01, 0.99).Draw(t, "value")

		fg := color.RGBA{R: 255, G: 0, B: 0, A: 255}
		bg := color.RGBA{R: 0, G: 0, B: 255, A: 255}

		cfg := Config{
			Style:        Segmented,
			Orientation:  orientation,
			Value:        value,
			Bounds:       bounds,
			Foreground:   fg,
			Background:   bg,
			SegmentCount: segCount,
			SegmentGap:   gap,
		}

		result := Render(cfg)
		if result == nil {
			t.Fatal("expected non-nil result")
		}

		fillExtent := int(float64(primaryLen) * value)

		// Determine which cells are filled based on cell center
		cellFilled := make([]bool, segCount)
		for i := 0; i < segCount; i++ {
			cellStart := i * (cellWidth + gap)
			cellCenter := cellStart + cellWidth/2
			cellFilled[i] = cellCenter < fillExtent
		}

		// Verify contiguity: once we see an unfilled cell, all subsequent must be unfilled.
		// This ensures filled cells are contiguous from the start of the direction.
		foundUnfilled := false
		for i := 0; i < segCount; i++ {
			if !cellFilled[i] {
				foundUnfilled = true
			} else if foundUnfilled {
				t.Fatalf("filled cell %d found after unfilled cell "+
					"[orient=%d, value=%.4f, fillExtent=%d, segCount=%d]",
					i, orientation, value, fillExtent, segCount)
			}
		}

		// Also verify the pixel colors match expectations:
		// For Horizontal: filled cells should have fg from the left
		// For Vertical: filled cells should have fg from the bottom (low primary pos = bottom)
		for i := 0; i < segCount; i++ {
			cellStart := i * (cellWidth + gap)
			// Sample the center pixel of the cell
			samplePos := cellStart + cellWidth/2
			sampleMinor := minorLen / 2

			var x, y int
			if orientation == OrientHorizontal {
				x = samplePos
				y = sampleMinor
			} else {
				x = sampleMinor
				y = bounds.Dy() - 1 - samplePos
			}

			if x < 0 || x >= bounds.Dx() || y < 0 || y >= bounds.Dy() {
				continue
			}

			r, g, b, a := result.Image.At(x, y).RGBA()
			fgR, fgG, fgB, fgA := fg.RGBA()
			bgR, bgG, bgB, bgA := bg.RGBA()

			if cellFilled[i] {
				if r != fgR || g != fgG || b != fgB || a != fgA {
					t.Fatalf("cell %d expected filled (fg) but got (%d,%d,%d,%d) "+
						"[orient=%d, value=%.4f, x=%d, y=%d]",
						i, r>>8, g>>8, b>>8, a>>8, orientation, value, x, y)
				}
			} else {
				if r != bgR || g != bgG || b != bgB || a != bgA {
					t.Fatalf("cell %d expected unfilled (bg) but got (%d,%d,%d,%d) "+
						"[orient=%d, value=%.4f, x=%d, y=%d]",
						i, r>>8, g>>8, b>>8, a>>8, orientation, value, x, y)
				}
			}
		}
	})
}

// TestProperty9_RingFillProportionality verifies Property 9: For any Ring-style bar with
// valid bounds (≥3×3), any thickness, and any value in [0.0, 1.0], the ratio of filled
// pixels within the ring annulus to total annulus pixels SHALL approximate the value
// within a tolerance of (2/radius + 0.02).

func TestProperty9_RingFillProportionality(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		width := rapid.IntRange(3, 80).Draw(t, "width")
		height := rapid.IntRange(3, 80).Draw(t, "height")
		value := rapid.Float64Range(0.0, 1.0).Draw(t, "value")
		orientation := Orientation(rapid.IntRange(0, 1).Draw(t, "orientation"))

		// Generate a thickness or 0 (auto-compute)
		minDim := width
		if height < minDim {
			minDim = height
		}
		maxThickness := minDim / 2
		if maxThickness < 1 {
			maxThickness = 1
		}
		thickness := rapid.IntRange(0, maxThickness).Draw(t, "thickness")

		bounds := image.Rect(0, 0, width, height)
		fg := color.RGBA{R: 255, G: 0, B: 0, A: 255}
		bg := color.RGBA{R: 0, G: 0, B: 255, A: 255}

		cfg := Config{
			Style:       Ring,
			Orientation: orientation,
			Value:       value,
			Bounds:      bounds,
			Foreground:  fg,
			Background:  bg,
			Thickness:   thickness,
		}

		result := Render(cfg)
		if result == nil {
			t.Fatal("expected non-nil result for valid ring bounds")
		}

		// After validate, determine actual thickness
		cfgCopy := cfg
		validate(&cfgCopy)
		actualThickness := cfgCopy.Thickness

		// Ring geometry
		cx := float64(width) / 2.0
		cy := float64(height) / 2.0
		outerR := math.Min(float64(width), float64(height)) / 2.0
		innerR := outerR - float64(actualThickness)
		if innerR < 0 {
			innerR = 0
		}

		// Count foreground and total pixels within the annulus
		fgCount := 0
		totalCount := 0

		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				px := float64(x) + 0.5
				py := float64(y) + 0.5
				dx := px - cx
				dy := py - cy
				dist := math.Sqrt(dx*dx + dy*dy)

				if dist >= innerR && dist <= outerR {
					totalCount++
					r, g, b, a := result.Image.At(x, y).RGBA()
					// Check if pixel is foreground (red)
					if r == 0xFFFF && g == 0 && b == 0 && a == 0xFFFF {
						fgCount++
					}
				}
			}
		}

		if totalCount == 0 {
			// Very small ring, skip
			return
		}

		actualRatio := float64(fgCount) / float64(totalCount)
		radius := outerR
		tolerance := 2.0/radius + 0.02

		if math.Abs(actualRatio-value) > tolerance {
			t.Fatalf("ring fill ratio mismatch: value=%f, actualRatio=%f, tolerance=%f, "+
				"width=%d, height=%d, thickness=%d, orientation=%d, fgCount=%d, totalCount=%d",
				value, actualRatio, tolerance, width, height, actualThickness, orientation,
				fgCount, totalCount)
		}
	})
}

// TestProperty10_ArcSweepAngleClamping verifies Property 10: For any Arc-style
// configuration with a sweep angle outside [90°, 350°], the rendered output SHALL be
// pixel-identical to the same configuration with the sweep angle clamped to [90°, 350°].

func TestProperty10_ArcSweepAngleClamping(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		width := rapid.IntRange(3, 60).Draw(t, "width")
		height := rapid.IntRange(3, 60).Draw(t, "height")
		value := rapid.Float64Range(0.0, 1.0).Draw(t, "value")
		orientation := Orientation(rapid.IntRange(0, 1).Draw(t, "orientation"))

		// Generate a sweep angle outside [90, 350]
		outOfRange := rapid.SampledFrom([]bool{true, false}).Draw(t, "below90")
		var sweepAngle float64
		if outOfRange {
			// Below 90
			sweepAngle = rapid.Float64Range(1.0, 89.9).Draw(t, "sweepAngle")
		} else {
			// Above 350
			sweepAngle = rapid.Float64Range(350.1, 720.0).Draw(t, "sweepAngle")
		}

		// Compute clamped value
		clampedSweep := sweepAngle
		if clampedSweep < 90 {
			clampedSweep = 90
		}
		if clampedSweep > 350 {
			clampedSweep = 350
		}

		bounds := image.Rect(0, 0, width, height)
		fg := color.RGBA{R: 200, G: 100, B: 50, A: 255}
		bg := color.RGBA{R: 10, G: 20, B: 30, A: 255}

		// Render with out-of-range sweep angle
		result1 := Render(Config{
			Style:       Arc,
			Orientation: orientation,
			Value:       value,
			Bounds:      bounds,
			Foreground:  fg,
			Background:  bg,
			SweepAngle:  sweepAngle,
		})

		// Render with clamped sweep angle
		result2 := Render(Config{
			Style:       Arc,
			Orientation: orientation,
			Value:       value,
			Bounds:      bounds,
			Foreground:  fg,
			Background:  bg,
			SweepAngle:  clampedSweep,
		})

		if result1 == nil {
			t.Fatal("expected non-nil result for out-of-range sweep angle render")
		}
		if result2 == nil {
			t.Fatal("expected non-nil result for clamped sweep angle render")
		}

		// Assert pixel-identical
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				r1, g1, b1, a1 := result1.Image.At(x, y).RGBA()
				r2, g2, b2, a2 := result2.Image.At(x, y).RGBA()
				if r1 != r2 || g1 != g2 || b1 != b2 || a1 != a2 {
					t.Fatalf("pixel mismatch at (%d,%d): unclamped=(%d,%d,%d,%d), clamped=(%d,%d,%d,%d) "+
						"[sweep=%f, clamped=%f, value=%f, orient=%d, bounds=%dx%d]",
						x, y, r1, g1, b1, a1, r2, g2, b2, a2,
						sweepAngle, clampedSweep, value, orientation, width, height)
				}
			}
		}
	})
}

// TestProperty22_RingStartAngle verifies Property 22: For any Ring-style bar with valid
// bounds (≥3×3) and any value in (0.0, 0.5], when Orientation is Horizontal the first
// filled arc pixel SHALL be adjacent to 12-o'clock (top center), and when Orientation is
// Vertical the first filled arc pixel SHALL be adjacent to 9-o'clock (left center).

func TestProperty22_RingStartAngle(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Use larger bounds to ensure clearly distinguishable pixels at cardinal points.
		// With radius ≥ 15 and value ≥ 0.05, the fill arc spans at least ~18° which
		// guarantees multiple pixels are filled near the start angle.
		width := rapid.IntRange(30, 80).Draw(t, "width")
		height := rapid.IntRange(30, 80).Draw(t, "height")
		// Value large enough that at least a few pixels are filled at start position
		value := rapid.Float64Range(0.05, 0.5).Draw(t, "value")
		orientation := Orientation(rapid.IntRange(0, 1).Draw(t, "orientation"))

		bounds := image.Rect(0, 0, width, height)
		fg := color.RGBA{R: 255, G: 0, B: 0, A: 255}
		bg := color.RGBA{R: 0, G: 0, B: 255, A: 255}

		cfg := Config{
			Style:       Ring,
			Orientation: orientation,
			Value:       value,
			Bounds:      bounds,
			Foreground:  fg,
			Background:  bg,
		}

		result := Render(cfg)
		if result == nil {
			t.Fatal("expected non-nil result for valid ring bounds")
		}

		// After validate, determine actual thickness
		cfgCopy := cfg
		validate(&cfgCopy)
		actualThickness := cfgCopy.Thickness

		cx := float64(width) / 2.0
		cy := float64(height) / 2.0
		outerR := math.Min(float64(width), float64(height)) / 2.0
		innerR := outerR - float64(actualThickness)
		if innerR < 0 {
			innerR = 0
		}

		isFg := func(x, y int) bool {
			if x < 0 || x >= width || y < 0 || y >= height {
				return false
			}
			r, g, b, a := result.Image.At(x, y).RGBA()
			return r == 0xFFFF && g == 0 && b == 0 && a == 0xFFFF
		}

		// Check that at least one foreground pixel exists within a small angular wedge
		// near the expected start angle, within the annulus.
		// The ring renderer uses atan2(dx, -dy) for clockwise-from-north angles.
		// Horizontal start: angle 0 (12-o'clock / top center)
		// Vertical start: angle 3π/2 (9-o'clock / left center)
		var startAngle float64
		if orientation == OrientHorizontal {
			startAngle = 0 // 12-o'clock in atan2(dx,-dy) space
		} else {
			startAngle = 3.0 * math.Pi / 2.0 // 9-o'clock in atan2(dx,-dy) space
		}

		// Search for foreground pixels within a small angular wedge (±15°) around startAngle
		// and within the annulus radii
		wedge := 15.0 * math.Pi / 180.0 // ±15° tolerance
		found := false
		for y := 0; y < height && !found; y++ {
			for x := 0; x < width && !found; x++ {
				px := float64(x) + 0.5
				py := float64(y) + 0.5
				dx := px - cx
				dy := py - cy
				dist := math.Sqrt(dx*dx + dy*dy)

				// Must be within annulus
				if dist < innerR || dist > outerR {
					continue
				}

				// Compute clockwise-from-north angle (same as renderer)
				angle := math.Atan2(dx, -dy)
				if angle < 0 {
					angle += 2.0 * math.Pi
				}

				// Check if within wedge around startAngle
				diff := angle - startAngle
				// Normalize to [-π, π]
				for diff > math.Pi {
					diff -= 2.0 * math.Pi
				}
				for diff < -math.Pi {
					diff += 2.0 * math.Pi
				}

				if math.Abs(diff) <= wedge && isFg(x, y) {
					found = true
				}
			}
		}

		if !found {
			orientStr := "horizontal (12-o'clock)"
			if orientation == OrientVertical {
				orientStr = "vertical (9-o'clock)"
			}
			t.Fatalf("ring %s: no foreground pixel found near start angle, "+
				"value=%f, width=%d, height=%d, thickness=%d",
				orientStr, value, width, height, actualThickness)
		}
	})
}

// TestProperty23_ArcFillEndpoint verifies Property 23: For any Arc-style bar with valid
// bounds (≥3×3) and any value in (0.0, 1.0), when Orientation is Horizontal the fill
// SHALL begin from the left endpoint of the sweep, and when Orientation is Vertical the
// fill SHALL begin from the bottom endpoint of the sweep.

func TestProperty23_ArcFillEndpoint(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Use larger bounds for sub-pixel accuracy
		width := rapid.IntRange(20, 80).Draw(t, "width")
		height := rapid.IntRange(20, 80).Draw(t, "height")
		value := rapid.Float64Range(0.05, 0.95).Draw(t, "value")
		orientation := Orientation(rapid.IntRange(0, 1).Draw(t, "orientation"))

		bounds := image.Rect(0, 0, width, height)
		fg := color.RGBA{R: 255, G: 0, B: 0, A: 255}
		bg := color.RGBA{R: 0, G: 0, B: 255, A: 255}

		cfg := Config{
			Style:       Arc,
			Orientation: orientation,
			Value:       value,
			Bounds:      bounds,
			Foreground:  fg,
			Background:  bg,
			// SweepAngle defaults to 270 after validation
		}

		result := Render(cfg)
		if result == nil {
			t.Fatal("expected non-nil result for valid arc bounds")
		}

		// After validate, get resolved parameters
		cfgCopy := cfg
		validate(&cfgCopy)
		actualThickness := cfgCopy.Thickness
		sweepAngle := cfgCopy.SweepAngle // Should be 270 degrees

		cx := float64(width) / 2.0
		cy := float64(height) / 2.0
		outerR := math.Min(float64(width), float64(height)) / 2.0
		innerR := outerR - float64(actualThickness)
		if innerR < 0 {
			innerR = 0
		}
		midR := (outerR + innerR) / 2.0

		// Compute the start angle of the arc:
		// The arc is centered symmetrically about the orientation reference angle.
		// geom.StartAngle: -π/2 for horizontal (12-o'clock), π for vertical (9-o'clock)
		sweepRad := sweepAngle * math.Pi / 180.0
		var refAngle float64
		if orientation == OrientHorizontal {
			refAngle = -math.Pi / 2.0
		} else {
			refAngle = math.Pi
		}
		arcStart := refAngle - sweepRad/2.0

		// The first filled pixel should be at arcStart (the left/bottom endpoint).
		// Convert arcStart (clockwise from north) to screen coordinates:
		// In our atan2(dx, -dy) convention: angle clockwise from north.
		// To convert back to x,y: x = cx + midR*sin(angle), y = cy - midR*cos(angle)
		startX := int(cx + midR*math.Sin(arcStart))
		startY := int(cy - midR*math.Cos(arcStart))

		isFg := func(x, y int) bool {
			if x < 0 || x >= width || y < 0 || y >= height {
				return false
			}
			r, g, b, a := result.Image.At(x, y).RGBA()
			return r == 0xFFFF && g == 0 && b == 0 && a == 0xFFFF
		}

		// Check that a foreground pixel exists near the expected start position
		found := isFg(startX, startY)
		if !found {
			// Check 3x3 neighborhood for rounding tolerance
			for dy := -1; dy <= 1; dy++ {
				for dx := -1; dx <= 1; dx++ {
					if isFg(startX+dx, startY+dy) {
						found = true
						break
					}
				}
				if found {
					break
				}
			}
		}

		if !found {
			orientStr := "horizontal"
			if orientation == OrientVertical {
				orientStr = "vertical"
			}
			t.Fatalf("arc %s: expected foreground pixel near start endpoint (%d,%d), "+
				"value=%f, width=%d, height=%d, thickness=%d, sweepAngle=%f, arcStart=%f rad",
				orientStr, startX, startY, value, width, height, actualThickness, sweepAngle, arcStart)
		}
	})
}

// TestProperty25_VerticalGradientDirection verifies Property 25: For any Ring or Arc style
// bar with a valid GradientFill (≥2 stops) and Vertical orientation, the gradient color
// interpolation SHALL proceed along the vertical fill axis.

func TestProperty25_VerticalGradientDirection(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		style := Style(rapid.SampledFrom([]Style{Ring, Arc}).Draw(t, "style"))

		// Use reasonable sizes for clear gradient visibility
		width := rapid.IntRange(30, 60).Draw(t, "width")
		height := rapid.IntRange(30, 60).Draw(t, "height")

		bounds := image.Rect(0, 0, width, height)

		// Use a simple two-stop gradient: red at position 0.0, green at position 1.0
		// This makes it easy to verify gradient direction by checking color progression.
		stops := []gradient.ColorStop{
			{Position: 0.0, Color: color.RGBA{R: 255, G: 0, B: 0, A: 255}},
			{Position: 1.0, Color: color.RGBA{R: 0, G: 255, B: 0, A: 255}},
		}

		cfg := Config{
			Style:       style,
			Orientation: OrientVertical,
			Value:       1.0, // Full fill so gradient covers entire arc
			Bounds:      bounds,
			Foreground:  color.RGBA{R: 255, G: 255, B: 255, A: 255},
			Background:  color.RGBA{R: 0, G: 0, B: 0, A: 255},
			Gradient:    &GradientFill{Stops: stops},
		}

		result := Render(cfg)
		if result == nil {
			t.Fatal("expected non-nil result for valid config")
		}

		// After validate, get resolved parameters
		cfgCopy := cfg
		validate(&cfgCopy)
		actualThickness := cfgCopy.Thickness

		cx := float64(width) / 2.0
		cy := float64(height) / 2.0
		outerR := math.Min(float64(width), float64(height)) / 2.0
		innerR := outerR - float64(actualThickness)
		if innerR < 0 {
			innerR = 0
		}
		midR := (outerR + innerR) / 2.0

		// The ring/arc renderer uses atan2(dx, -dy) to compute angles clockwise from north.
		// For vertical orientation:
		//   Ring: startOffset = 3π/2 (9-o'clock), fills clockwise. Gradient t=0 at 9-o'clock.
		//   Arc: arcStart = π - sweepRad/2, fills from arcStart clockwise.
		//
		// To convert atan2(dx,-dy) angle θ to screen coords:
		//   x = cx + midR * sin(θ)
		//   y = cy - midR * cos(θ)
		//
		// For Ring vertical:
		//   Start (t=0) at atan2 angle 3π/2: x = cx + midR*sin(3π/2) = cx - midR, y = cy - midR*cos(3π/2) = cy
		//   That's the left center (9-o'clock). ✓
		//   Halfway (t=0.5) at atan2 angle 3π/2 + π = 5π/2 → normalized = π/2:
		//   x = cx + midR*sin(π/2) = cx + midR, y = cy - midR*cos(π/2) = cy
		//   That's the right center (3-o'clock).

		// Sample pixels at t≈0.1 (near start) and t≈0.5 (halfway)
		var earlyAngle, lateAngle float64

		if style == Ring {
			// Ring: start at atan2 angle 3π/2, progresses clockwise (increasing atan2 angle mod 2π)
			// t fraction maps to angle: startOffset + t * 2π
			earlyAngle = 3.0*math.Pi/2.0 + 0.1*2.0*math.Pi // t=0.1
			lateAngle = 3.0*math.Pi/2.0 + 0.5*2.0*math.Pi  // t=0.5
		} else {
			// Arc: arcStart = geom.StartAngle - sweepRad/2
			// geom.StartAngle for vertical Arc = π (9-o'clock in atan2(dx,-dy) system??)
			// Actually: for Arc, the renderer uses arcStart = geom.StartAngle - sweepRad/2
			// and geom.StartAngle comes from resolveOrientation which returns π for vertical Arc.
			//
			// But the Arc renderer computes relAngle = angle - arcStart in the atan2(dx,-dy) space.
			// So the start angle in atan2 space is indeed: the angle where relAngle = 0.
			//
			// geom.StartAngle for vertical = π. But that's the "reference" angle.
			// arcStart = π - sweepRad/2 (this is a clockwise-from-north angle).
			//
			// Wait, looking at renderArc: it uses math.Atan2(dx, -dy) just like renderRing,
			// and arcStart = geom.StartAngle - sweepRad/2.
			// For vertical: geom.StartAngle = π (from resolveOrientation for Arc vertical)
			//
			// But geom.StartAngle is in what coordinate system? Looking at resolveOrientation:
			// For Arc vertical: startAngle = math.Pi → same -π/2 vs π as Ring.
			//
			// In renderArc: angle = math.Atan2(dx, -dy) which gives clockwise-from-north.
			// Then relAngle = angle - arcStart.
			// arcStart = geom.StartAngle - sweepRad/2.
			//
			// Hmm, but geom.StartAngle = π for vertical. atan2 gives values in [-π, π].
			// And arcStart = π - sweepRad/2.
			// For default sweep 270° = 3π/2 rad: arcStart = π - 3π/4 = π/4.
			//
			// So the arc starts at atan2 angle π/4 (which is clockwise 45° from north = ~1:30 position)
			// That seems correct for a 270° arc centered on 9-o'clock (π in atan2 space... wait)
			//
			// Let me reconsider. atan2(dx,-dy) gives clockwise-from-north in [-π,π].
			// π in this system is... sin(π)=0, cos(π)=-1 → x=cx, y=cy+midR → bottom center!
			// That's 6-o'clock, not 9-o'clock.
			//
			// Wait: x = cx + midR*sin(θ), y = cy - midR*cos(θ)
			// For θ=π: x = cx + 0 = cx, y = cy - midR*(-1) = cy + midR → bottom center (6 o'clock)
			// For θ=-π/2 (or 3π/2): x = cx + midR*sin(-π/2) = cx - midR, y = cy - midR*cos(-π/2) = cy → left center (9 o'clock)
			//
			// So in resolveOrientation: vertical Arc gets startAngle = π (which is 6-o'clock, not 9).
			// But the design says vertical should start at bottom endpoint...
			// "Vertical orientation: fill begins from the bottom endpoint of the sweep"
			// The arc is centered symmetrically about the reference angle. For vertical,
			// geom.StartAngle = π (6-o'clock/bottom center). arcStart = π - sweepRad/2.
			// That means the arc spans from (π - sweepRad/2) to (π + sweepRad/2), centered on bottom.
			// Fill starts from arcStart = the left endpoint of this bottom-centered sweep.
			// This is the "bottom endpoint" per the design.

			sweepRad := cfgCopy.SweepAngle * math.Pi / 180.0
			arcStartAngle := math.Pi - sweepRad/2.0
			earlyAngle = arcStartAngle + 0.1*sweepRad // t=0.1 along the arc
			lateAngle = arcStartAngle + 0.5*sweepRad  // t=0.5 along the arc
		}

		// Normalize angles to [0, 2π)
		normalizeAngle := func(a float64) float64 {
			a = math.Mod(a, 2*math.Pi)
			if a < 0 {
				a += 2 * math.Pi
			}
			return a
		}
		earlyAngle = normalizeAngle(earlyAngle)
		lateAngle = normalizeAngle(lateAngle)

		// Convert to pixel coordinates
		earlyX := int(cx + midR*math.Sin(earlyAngle))
		earlyY := int(cy - midR*math.Cos(earlyAngle))
		lateX := int(cx + midR*math.Sin(lateAngle))
		lateY := int(cy - midR*math.Cos(lateAngle))

		// Clamp to image bounds
		clamp := func(v, max int) int {
			if v < 0 {
				return 0
			}
			if v >= max {
				return max - 1
			}
			return v
		}
		earlyX = clamp(earlyX, width)
		earlyY = clamp(earlyY, height)
		lateX = clamp(lateX, width)
		lateY = clamp(lateY, height)

		// Get red channel at early and late positions
		sr, _, _, _ := result.Image.At(earlyX, earlyY).RGBA()
		lr, _, _, _ := result.Image.At(lateX, lateY).RGBA()

		earlyR := int(sr >> 8)
		lateR := int(lr >> 8)

		// Gradient goes red(255,0,0) → green(0,255,0).
		// At t=0.1 the red channel should be ~230 (high).
		// At t=0.5 the red channel should be ~128 (medium).
		// So earlyR should be > lateR (more red near the start).
		// Allow generous tolerance for sub-pixel rounding.
		if earlyR+40 < lateR {
			t.Fatalf("gradient direction incorrect: early pixel (%d,%d) R=%d should be > "+
				"late pixel (%d,%d) R=%d for vertical %d style [width=%d, height=%d, thickness=%d]",
				earlyX, earlyY, earlyR,
				lateX, lateY, lateR,
				style, width, height, actualThickness)
		}
	})
}

// TestProperty20_SignDeterminismAndAnimationSensitivity verifies Property 20:
// For any two Config values that are identical in all fields including animElapsed,
// sign SHALL produce identical uint64 values. For any Config where only animElapsed
// differs, sign SHALL produce different uint64 values.

func TestProperty20_SignDeterminismAndAnimationSensitivity(t *testing.T) {
	t.Run("identical_configs_identical_hashes", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			cfg := drawRandomConfig(t)

			hash1 := sign(cfg)
			hash2 := sign(cfg)

			if hash1 != hash2 {
				t.Fatalf("sign not deterministic: same config produced %d and %d", hash1, hash2)
			}
		})
	})

	t.Run("different_animElapsed_different_hashes", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			cfg := drawRandomConfig(t)

			// Ensure two distinct animElapsed values
			elapsed1 := rapid.Int64Range(0, 1_000_000_000).Draw(t, "elapsed1")
			elapsed2 := rapid.Int64Range(0, 1_000_000_000).Filter(func(v int64) bool {
				return v != elapsed1
			}).Draw(t, "elapsed2")

			cfg1 := cfg
			cfg1.animElapsed = time.Duration(elapsed1)

			cfg2 := cfg
			cfg2.animElapsed = time.Duration(elapsed2)

			hash1 := sign(cfg1)
			hash2 := sign(cfg2)

			if hash1 == hash2 {
				t.Fatalf("sign collision: different animElapsed (%v vs %v) produced same hash %d",
					cfg1.animElapsed, cfg2.animElapsed, hash1)
			}
		})
	})
}

// drawRandomConfig generates a random Config with all fields populated for property testing.
func drawRandomConfig(t *rapid.T) Config {
	style := Style(rapid.IntRange(0, 4).Draw(t, "style"))
	orientation := Orientation(rapid.IntRange(0, 1).Draw(t, "orientation"))
	value := rapid.Float64Range(0.0, 1.0).Draw(t, "value")

	minX := rapid.IntRange(0, 100).Draw(t, "minX")
	minY := rapid.IntRange(0, 100).Draw(t, "minY")
	width := rapid.IntRange(3, 100).Draw(t, "width")
	height := rapid.IntRange(3, 100).Draw(t, "height")
	bounds := image.Rect(minX, minY, minX+width, minY+height)

	fg := color.RGBA{
		R: uint8(rapid.IntRange(0, 255).Draw(t, "fgR")),
		G: uint8(rapid.IntRange(0, 255).Draw(t, "fgG")),
		B: uint8(rapid.IntRange(0, 255).Draw(t, "fgB")),
		A: uint8(rapid.IntRange(0, 255).Draw(t, "fgA")),
	}
	bg := color.RGBA{
		R: uint8(rapid.IntRange(0, 255).Draw(t, "bgR")),
		G: uint8(rapid.IntRange(0, 255).Draw(t, "bgG")),
		B: uint8(rapid.IntRange(0, 255).Draw(t, "bgB")),
		A: uint8(rapid.IntRange(0, 255).Draw(t, "bgA")),
	}

	// Gradient: sometimes nil, sometimes with valid stops
	var grad *GradientFill
	hasGradient := rapid.SampledFrom([]bool{true, false}).Draw(t, "hasGradient")
	if hasGradient {
		numStops := rapid.IntRange(2, 8).Draw(t, "numStops")
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
		grad = &GradientFill{Stops: stops}
	}

	segmentCount := rapid.IntRange(0, 20).Draw(t, "segmentCount")
	segmentGap := rapid.IntRange(1, 4).Draw(t, "segmentGap")
	thickness := rapid.IntRange(0, 20).Draw(t, "thickness")
	sweepAngle := rapid.Float64Range(90.0, 350.0).Draw(t, "sweepAngle")
	roundedCaps := rapid.SampledFrom([]bool{true, false}).Draw(t, "roundedCaps")
	borderWidth := rapid.IntRange(0, 16).Draw(t, "borderWidth")
	borderColor := color.RGBA{
		R: uint8(rapid.IntRange(0, 255).Draw(t, "borderR")),
		G: uint8(rapid.IntRange(0, 255).Draw(t, "borderG")),
		B: uint8(rapid.IntRange(0, 255).Draw(t, "borderB")),
		A: uint8(rapid.IntRange(0, 255).Draw(t, "borderA")),
	}

	// Markers: 0–8 random markers
	numMarkers := rapid.IntRange(0, 8).Draw(t, "numMarkers")
	markers := make([]ThresholdMarker, numMarkers)
	for i := range markers {
		markers[i] = ThresholdMarker{
			Value: rapid.Float64Range(0.0, 1.0).Draw(t, "markerVal"),
			Color: color.RGBA{
				R: uint8(rapid.IntRange(0, 255).Draw(t, "markerR")),
				G: uint8(rapid.IntRange(0, 255).Draw(t, "markerG")),
				B: uint8(rapid.IntRange(0, 255).Draw(t, "markerB")),
				A: uint8(rapid.IntRange(0, 255).Draw(t, "markerA")),
			},
		}
	}

	// Animation config
	animType := Animation(rapid.IntRange(0, 3).Draw(t, "animType"))
	animPeriod := time.Duration(rapid.Int64Range(100_000_000, 10_000_000_000).Draw(t, "animPeriod"))
	animSpeed := rapid.IntRange(10, 500).Draw(t, "animSpeed")

	// animElapsed
	animElapsed := time.Duration(rapid.Int64Range(0, 5_000_000_000).Draw(t, "animElapsed"))

	return Config{
		Style:        style,
		Orientation:  orientation,
		Value:        value,
		Bounds:       bounds,
		Foreground:   fg,
		Background:   bg,
		Gradient:     grad,
		SegmentCount: segmentCount,
		SegmentGap:   segmentGap,
		Thickness:    thickness,
		SweepAngle:   sweepAngle,
		RoundedCaps:  roundedCaps,
		BorderWidth:  borderWidth,
		BorderColor:  borderColor,
		Markers:      markers,
		Animation: AnimationConfig{
			Type:   animType,
			Period: animPeriod,
			Speed:  animSpeed,
		},
		animElapsed: animElapsed,
	}
}

// TestProperty11_BorderOccupiesPerimeterBand verifies Property 11: For any bar
// configuration with a border width in [1, 16] and minor axis ≥ (2×border_width + 2),
// the outermost border_width pixels on each edge of the bar SHALL be the configured
// border color, and no fill or track pixels SHALL appear within the border band.
//
// Tests Linear and Segmented styles without RoundedCaps (rectangular border).

func TestProperty11_BorderOccupiesPerimeterBand(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Choose Linear or Segmented (rectangular border cases)
		style := Style(rapid.SampledFrom([]Style{Linear, Segmented}).Draw(t, "style"))
		orientation := Orientation(rapid.IntRange(0, 1).Draw(t, "orientation"))

		// Border width in [1, 16]
		bw := rapid.IntRange(1, 16).Draw(t, "borderWidth")

		// Minor axis must be ≥ 2*bw + 2 so border is drawn.
		// For Horizontal orientation: minor axis = height; for Vertical: minor axis = width.
		minMinor := 2*bw + 2
		minorAxis := rapid.IntRange(minMinor, minMinor+50).Draw(t, "minorAxis")
		// Primary axis just needs to be valid (≥ 2*bw + 2 so border is visible there too)
		primaryAxis := rapid.IntRange(minMinor, minMinor+80).Draw(t, "primaryAxis")

		var w, h int
		if orientation == OrientHorizontal {
			w = primaryAxis
			h = minorAxis
		} else {
			w = minorAxis
			h = primaryAxis
		}

		bounds := image.Rect(0, 0, w, h)
		value := rapid.Float64Range(0.0, 1.0).Draw(t, "value")

		// Use distinct non-zero colors that won't be confused with each other
		fg := color.RGBA{R: 255, G: 0, B: 0, A: 255}
		bg := color.RGBA{R: 0, G: 0, B: 255, A: 255}
		bc := color.RGBA{R: 0, G: 255, B: 0, A: 255} // border = green, distinct from fg and bg

		cfg := Config{
			Style:       style,
			Orientation: orientation,
			Value:       value,
			Bounds:      bounds,
			Foreground:  fg,
			Background:  bg,
			BorderWidth: bw,
			BorderColor: bc,
			RoundedCaps: false, // Rectangular border only
		}

		result := Render(cfg)
		if result == nil {
			t.Fatal("expected non-nil result for valid bordered bar config")
		}

		img := result.Image

		// Helper to check if a pixel matches the border color
		isBorderColor := func(x, y int) bool {
			r, g, b, a := img.At(x, y).RGBA()
			er, eg, eb, ea := bc.RGBA()
			return r == er && g == eg && b == eb && a == ea
		}

		// Check the rectangular border band:
		// Top band: rows [0, bw)
		for y := 0; y < bw; y++ {
			for x := 0; x < w; x++ {
				if !isBorderColor(x, y) {
					t.Fatalf("top border pixel (%d,%d): expected border color %v, got %v [bw=%d, w=%d, h=%d, style=%d, orient=%d]",
						x, y, bc, img.At(x, y), bw, w, h, style, orientation)
				}
			}
		}
		// Bottom band: rows [h-bw, h)
		for y := h - bw; y < h; y++ {
			for x := 0; x < w; x++ {
				if !isBorderColor(x, y) {
					t.Fatalf("bottom border pixel (%d,%d): expected border color %v, got %v [bw=%d, w=%d, h=%d, style=%d, orient=%d]",
						x, y, bc, img.At(x, y), bw, w, h, style, orientation)
				}
			}
		}
		// Left band: rows [bw, h-bw), cols [0, bw)
		for y := bw; y < h-bw; y++ {
			for x := 0; x < bw; x++ {
				if !isBorderColor(x, y) {
					t.Fatalf("left border pixel (%d,%d): expected border color %v, got %v [bw=%d, w=%d, h=%d, style=%d, orient=%d]",
						x, y, bc, img.At(x, y), bw, w, h, style, orientation)
				}
			}
		}
		// Right band: rows [bw, h-bw), cols [w-bw, w)
		for y := bw; y < h-bw; y++ {
			for x := w - bw; x < w; x++ {
				if !isBorderColor(x, y) {
					t.Fatalf("right border pixel (%d,%d): expected border color %v, got %v [bw=%d, w=%d, h=%d, style=%d, orient=%d]",
						x, y, bc, img.At(x, y), bw, w, h, style, orientation)
				}
			}
		}

		// Additionally verify no fill or track pixels exist within border band.
		// We already confirmed they're all border color above, so this is implicitly satisfied.
		// But also check that the interior does NOT contain border-colored pixels in the
		// fill/track region (interior is [bw, w-bw) x [bw, h-bw)).
		// NOTE: This is NOT required by the property — the property only asserts border band content.
		// Interior pixels may legitimately include border color if markers use it.
		// The critical property is that the border band pixels ARE the border color.
	})
}

// TestProperty12_BorderSkippedWhenTooThin verifies Property 12: For any bar
// configuration where the minor axis is less than (2×border_width + 2), the output
// SHALL contain no border-colored pixels (border rendering is skipped).

func TestProperty12_BorderSkippedWhenTooThin(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Choose Linear or Segmented (rectangular border cases)
		style := Style(rapid.SampledFrom([]Style{Linear, Segmented}).Draw(t, "style"))
		orientation := Orientation(rapid.IntRange(0, 1).Draw(t, "orientation"))

		// Border width in [1, 16]
		bw := rapid.IntRange(1, 16).Draw(t, "borderWidth")

		// Minor axis must be < 2*bw + 2 so border is skipped.
		// Minor axis must be at least 1 for valid rendering.
		maxMinor := 2*bw + 1
		if maxMinor < 1 {
			maxMinor = 1
		}
		minorAxis := rapid.IntRange(1, maxMinor).Draw(t, "minorAxis")

		// Primary axis should be valid (at least 1)
		primaryAxis := rapid.IntRange(1, 50).Draw(t, "primaryAxis")

		var w, h int
		if orientation == OrientHorizontal {
			w = primaryAxis
			h = minorAxis
		} else {
			w = minorAxis
			h = primaryAxis
		}

		bounds := image.Rect(0, 0, w, h)
		value := rapid.Float64Range(0.0, 1.0).Draw(t, "value")

		// Use a distinct border color that won't appear as fg or bg
		fg := color.RGBA{R: 255, G: 0, B: 0, A: 255}
		bg := color.RGBA{R: 0, G: 0, B: 255, A: 255}
		bc := color.RGBA{R: 0, G: 255, B: 0, A: 255} // border = green, distinct from fg and bg

		cfg := Config{
			Style:       style,
			Orientation: orientation,
			Value:       value,
			Bounds:      bounds,
			Foreground:  fg,
			Background:  bg,
			BorderWidth: bw,
			BorderColor: bc,
			RoundedCaps: false,
		}

		result := Render(cfg)
		if result == nil {
			t.Fatal("expected non-nil result for valid bar config (even without border)")
		}

		img := result.Image

		// Pre-compute border color RGBA for comparison
		bcR, bcG, bcB, bcA := bc.RGBA()

		// Verify NO pixel in the output matches the border color
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				r, g, b, a := img.At(x, y).RGBA()
				if r == bcR && g == bcG && b == bcB && a == bcA {
					t.Fatalf("found border-colored pixel at (%d,%d) but border should be skipped [minor=%d < 2*%d+2=%d, w=%d, h=%d, style=%d, orient=%d]",
						x, y, minorAxis, bw, 2*bw+2, w, h, style, orientation)
				}
			}
		}
	})
}

// TestProperty13_ThresholdMarkersAtCorrectPositions verifies Property 13:
// For any bar configuration with 1–8 threshold markers, for each marker a 1-pixel-wide
// line perpendicular to the bar axis SHALL exist at the pixel position corresponding to
// clamp(marker.Value, 0, 1) × axis_length. The line SHALL span the full cross-axis extent.
// Marker pixels SHALL be the configured marker color (or opaque white if zero-value).
// Markers SHALL be visible regardless of the current Value (drawn on top of fill and track).

func TestProperty13_ThresholdMarkersAtCorrectPositions(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Use Linear or Segmented style for simpler position calculation.
		style := Style(rapid.SampledFrom([]Style{Linear, Segmented}).Draw(t, "style"))
		orientation := Orientation(rapid.IntRange(0, 1).Draw(t, "orientation"))

		// Generate bounds large enough to have meaningful marker positions.
		var w, h int
		w = rapid.IntRange(5, 80).Draw(t, "width")
		h = rapid.IntRange(5, 80).Draw(t, "height")
		bounds := image.Rect(0, 0, w, h)

		// Generate 1–8 markers with distinct colors different from fg and bg.
		numMarkers := rapid.IntRange(1, 8).Draw(t, "numMarkers")
		markers := make([]ThresholdMarker, numMarkers)
		for i := 0; i < numMarkers; i++ {
			// Generate marker values in [-0.5, 1.5] to exercise clamping.
			markerVal := rapid.Float64Range(-0.5, 1.5).Draw(t, "markerVal")
			// Use a marker color that is distinguishable from fg and bg.
			// Use non-zero marker color (R=0, G=255, B=0, A=255) — green.
			markers[i] = ThresholdMarker{
				Value: markerVal,
				Color: color.RGBA{R: 0, G: 255, B: 0, A: 255},
			}
		}

		// Use fg=red, bg=blue so marker (green) is unambiguous.
		fg := color.RGBA{R: 255, G: 0, B: 0, A: 255}
		bg := color.RGBA{R: 0, G: 0, B: 255, A: 255}

		// Generate a random value to test markers are visible in both fill and track regions.
		value := rapid.Float64Range(0.0, 1.0).Draw(t, "value")

		cfg := Config{
			Style:       style,
			Orientation: orientation,
			Value:       value,
			Bounds:      bounds,
			Foreground:  fg,
			Background:  bg,
			Markers:     markers,
		}

		// For segmented style, use enough segments to avoid unsegmented fallback.
		if style == Segmented {
			cfg.SegmentCount = 4
			cfg.SegmentGap = 1
		}

		result := Render(cfg)
		if result == nil {
			t.Fatal("expected non-nil result for valid config with markers")
		}

		markerColor := color.RGBA{R: 0, G: 255, B: 0, A: 255}
		markerR, markerG, markerB, markerA := markerColor.RGBA()

		for _, m := range markers {
			val := m.Value
			// Clamp marker value to [0, 1]
			if val < 0.0 {
				val = 0.0
			}
			if val > 1.0 {
				val = 1.0
			}

			if orientation == OrientHorizontal {
				// Horizontal: marker at x = int(val * float64(w-1))
				expectedX := int(val * float64(w-1))
				// Check that all pixels in that column are the marker color.
				for y := 0; y < h; y++ {
					r, g, b, a := result.Image.At(expectedX, y).RGBA()
					if r != markerR || g != markerG || b != markerB || a != markerA {
						t.Fatalf("marker at value=%.4f: pixel (%d,%d) expected marker color "+
							"RGBA(%d,%d,%d,%d), got (%d,%d,%d,%d) "+
							"[orient=Horizontal, style=%d, barValue=%.4f, w=%d, h=%d]",
							m.Value, expectedX, y,
							markerColor.R, markerColor.G, markerColor.B, markerColor.A,
							uint8(r>>8), uint8(g>>8), uint8(b>>8), uint8(a>>8),
							style, value, w, h)
					}
				}
			} else {
				// Vertical: marker at y = h-1-int(val*float64(h-1))
				expectedY := h - 1 - int(val*float64(h-1))
				// Check that all pixels in that row are the marker color.
				for x := 0; x < w; x++ {
					r, g, b, a := result.Image.At(x, expectedY).RGBA()
					if r != markerR || g != markerG || b != markerB || a != markerA {
						t.Fatalf("marker at value=%.4f: pixel (%d,%d) expected marker color "+
							"RGBA(%d,%d,%d,%d), got (%d,%d,%d,%d) "+
							"[orient=Vertical, style=%d, barValue=%.4f, w=%d, h=%d]",
							m.Value, x, expectedY,
							markerColor.R, markerColor.G, markerColor.B, markerColor.A,
							uint8(r>>8), uint8(g>>8), uint8(b>>8), uint8(a>>8),
							style, value, w, h)
					}
				}
			}
		}
	})
}

// TestProperty14_MarkerTruncationTo8 verifies Property 14:
// For any configuration with more than 8 threshold markers, the rendered output SHALL be
// pixel-identical to the same configuration with only the first 8 markers.

func TestProperty14_MarkerTruncationTo8(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Use Linear or Segmented style (consistent with Property 13).
		style := Style(rapid.SampledFrom([]Style{Linear, Segmented}).Draw(t, "style"))
		orientation := Orientation(rapid.IntRange(0, 1).Draw(t, "orientation"))

		var w, h int
		w = rapid.IntRange(5, 80).Draw(t, "width")
		h = rapid.IntRange(5, 80).Draw(t, "height")
		bounds := image.Rect(0, 0, w, h)

		// Generate 9–16 markers.
		numMarkers := rapid.IntRange(9, 16).Draw(t, "numMarkers")
		markers := make([]ThresholdMarker, numMarkers)
		for i := 0; i < numMarkers; i++ {
			markerVal := rapid.Float64Range(0.0, 1.0).Draw(t, "markerVal")
			markerColor := color.RGBA{
				R: uint8(rapid.IntRange(1, 255).Draw(t, "mR")),
				G: uint8(rapid.IntRange(1, 255).Draw(t, "mG")),
				B: uint8(rapid.IntRange(1, 255).Draw(t, "mB")),
				A: 255,
			}
			markers[i] = ThresholdMarker{
				Value: markerVal,
				Color: markerColor,
			}
		}

		fg := color.RGBA{R: 200, G: 50, B: 50, A: 255}
		bg := color.RGBA{R: 20, G: 20, B: 80, A: 255}
		value := rapid.Float64Range(0.0, 1.0).Draw(t, "value")

		// Config with all markers (>8).
		cfgAll := Config{
			Style:       style,
			Orientation: orientation,
			Value:       value,
			Bounds:      bounds,
			Foreground:  fg,
			Background:  bg,
			Markers:     markers,
		}

		// Config with only the first 8 markers.
		cfgFirst8 := Config{
			Style:       style,
			Orientation: orientation,
			Value:       value,
			Bounds:      bounds,
			Foreground:  fg,
			Background:  bg,
			Markers:     markers[:8],
		}

		// For segmented style, use consistent segment config.
		if style == Segmented {
			cfgAll.SegmentCount = 4
			cfgAll.SegmentGap = 1
			cfgFirst8.SegmentCount = 4
			cfgFirst8.SegmentGap = 1
		}

		result1 := Render(cfgAll)
		result2 := Render(cfgFirst8)

		if result1 == nil {
			t.Fatal("expected non-nil result for config with >8 markers")
		}
		if result2 == nil {
			t.Fatal("expected non-nil result for config with first 8 markers")
		}

		// Assert pixel-identical output.
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				r1, g1, b1, a1 := result1.Image.At(x, y).RGBA()
				r2, g2, b2, a2 := result2.Image.At(x, y).RGBA()
				if r1 != r2 || g1 != g2 || b1 != b2 || a1 != a2 {
					t.Fatalf("pixel mismatch at (%d,%d): allMarkers=(%d,%d,%d,%d), first8=(%d,%d,%d,%d) "+
						"[style=%d, orient=%d, value=%.4f, numMarkers=%d, bounds=%dx%d]",
						x, y, r1>>8, g1>>8, b1>>8, a1>>8, r2>>8, g2>>8, b2>>8, a2>>8,
						style, orientation, value, numMarkers, w, h)
				}
			}
		}
	})
}

// TestProperty7_PulseAnimationBrightnessBounds verifies Property 7: For any bar with
// Pulse animation, any valid period, and any elapsed time, every pixel in the fill
// region SHALL have brightness between 30% and 100% (inclusive) of the configured
// foreground color intensity.

func TestProperty7_PulseAnimationBrightnessBounds(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate a Linear bar with solid foreground (no gradient) and Pulse animation.
		width := rapid.IntRange(5, 80).Draw(t, "width")
		height := rapid.IntRange(3, 40).Draw(t, "height")

		// Ensure at least one fill column: value must be >= 1/width so floor(width*value) >= 1.
		minValue := 1.0 / float64(width)
		if minValue > 1.0 {
			minValue = 1.0
		}
		value := rapid.Float64Range(minValue+0.001, 1.0).Draw(t, "value")

		// Use a foreground color with channels high enough that even at 30% brightness
		// the result won't collapse to background (0,0,0). Min channel = 4 ensures
		// floor(4 * 0.30) = 1 ≠ 0 (background).
		fgR := uint8(rapid.IntRange(4, 255).Draw(t, "fgR"))
		fgG := uint8(rapid.IntRange(4, 255).Draw(t, "fgG"))
		fgB := uint8(rapid.IntRange(4, 255).Draw(t, "fgB"))
		fg := color.RGBA{R: fgR, G: fgG, B: fgB, A: 255}

		// Background must differ from foreground to identify fill pixels.
		bg := color.RGBA{R: 0, G: 0, B: 0, A: 255}

		// Random valid period [100ms, 10s] and random elapsed time.
		period := time.Duration(rapid.Int64Range(100_000_000, 10_000_000_000).Draw(t, "period"))
		elapsed := time.Duration(rapid.Int64Range(0, 20_000_000_000).Draw(t, "elapsed"))

		bounds := image.Rect(0, 0, width, height)

		cfg := Config{
			Style:       Linear,
			Orientation: OrientHorizontal,
			Value:       value,
			Bounds:      bounds,
			Foreground:  fg,
			Background:  bg,
			Animation: AnimationConfig{
				Type:   Pulse,
				Period: period,
			},
			animElapsed: elapsed,
		}

		result := Render(cfg)
		if result == nil {
			t.Fatal("expected non-nil result for valid config")
		}

		// Check every pixel that is NOT background (i.e., fill pixels).
		// Each channel must be between floor(0.30 * fg_channel) - 1 and fg_channel.
		bgArr := [4]uint8{bg.R, bg.G, bg.B, bg.A}

		foundFillPixel := false
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				c := result.Image.(*image.RGBA).RGBAAt(x, y)

				// Skip background pixels.
				if c.R == bgArr[0] && c.G == bgArr[1] && c.B == bgArr[2] && c.A == bgArr[3] {
					continue
				}
				// Skip transparent pixels.
				if c.A == 0 {
					continue
				}

				foundFillPixel = true

				// Check R channel bounds.
				minR := int(float64(fgR)*0.30) - 1
				if minR < 0 {
					minR = 0
				}
				if int(c.R) < minR || int(c.R) > int(fgR) {
					t.Fatalf("pixel (%d,%d) R channel %d out of bounds [%d, %d] (fg.R=%d, period=%v, elapsed=%v)",
						x, y, c.R, minR, fgR, fgR, period, elapsed)
				}

				// Check G channel bounds.
				minG := int(float64(fgG)*0.30) - 1
				if minG < 0 {
					minG = 0
				}
				if int(c.G) < minG || int(c.G) > int(fgG) {
					t.Fatalf("pixel (%d,%d) G channel %d out of bounds [%d, %d] (fg.G=%d, period=%v, elapsed=%v)",
						x, y, c.G, minG, fgG, fgG, period, elapsed)
				}

				// Check B channel bounds.
				minB := int(float64(fgB)*0.30) - 1
				if minB < 0 {
					minB = 0
				}
				if int(c.B) < minB || int(c.B) > int(fgB) {
					t.Fatalf("pixel (%d,%d) B channel %d out of bounds [%d, %d] (fg.B=%d, period=%v, elapsed=%v)",
						x, y, c.B, minB, fgB, fgB, period, elapsed)
				}
			}
		}

		if !foundFillPixel {
			t.Fatal("expected at least one fill pixel but found none")
		}
	})
}

// TestProperty8_NoAnimationEquivalence verifies Property 8: For any bar configuration
// where the animation type is NoAnimation, OR the animation period/speed is ≤ 0,
// the rendered output SHALL be pixel-identical to the same configuration rendered
// without any animation effects.

func TestProperty8_NoAnimationEquivalence(t *testing.T) {
	t.Run("NoAnimation_type_identical_to_disabled", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			// Generate a random valid config with NoAnimation type but arbitrary period/speed/elapsed.
			style := Style(rapid.IntRange(0, 4).Draw(t, "style"))
			orientation := Orientation(rapid.IntRange(0, 1).Draw(t, "orientation"))

			var w, h int
			switch style {
			case Pie, Ring, Arc:
				w = rapid.IntRange(3, 50).Draw(t, "width")
				h = rapid.IntRange(3, 50).Draw(t, "height")
			default:
				w = rapid.IntRange(1, 50).Draw(t, "width")
				h = rapid.IntRange(1, 50).Draw(t, "height")
			}
			bounds := image.Rect(0, 0, w, h)
			value := rapid.Float64Range(0.0, 1.0).Draw(t, "value")

			fg := color.RGBA{
				R: uint8(rapid.IntRange(1, 255).Draw(t, "fgR")),
				G: uint8(rapid.IntRange(1, 255).Draw(t, "fgG")),
				B: uint8(rapid.IntRange(1, 255).Draw(t, "fgB")),
				A: 255,
			}
			bg := color.RGBA{
				R: uint8(rapid.IntRange(0, 254).Draw(t, "bgR")),
				G: uint8(rapid.IntRange(0, 254).Draw(t, "bgG")),
				B: uint8(rapid.IntRange(0, 254).Draw(t, "bgB")),
				A: 255,
			}

			// Random period and elapsed (should not matter since NoAnimation).
			period := time.Duration(rapid.Int64Range(100_000_000, 10_000_000_000).Draw(t, "period"))
			speed := rapid.IntRange(10, 500).Draw(t, "speed")
			elapsed := time.Duration(rapid.Int64Range(0, 5_000_000_000).Draw(t, "elapsed"))

			// Config WITH NoAnimation type but with period/speed/elapsed set.
			cfgWithAnim := Config{
				Style:       style,
				Orientation: orientation,
				Value:       value,
				Bounds:      bounds,
				Foreground:  fg,
				Background:  bg,
				Animation: AnimationConfig{
					Type:   NoAnimation,
					Period: period,
					Speed:  speed,
				},
				animElapsed: elapsed,
			}

			// Config WITHOUT any animation fields set.
			cfgNoAnim := Config{
				Style:       style,
				Orientation: orientation,
				Value:       value,
				Bounds:      bounds,
				Foreground:  fg,
				Background:  bg,
			}

			result1 := Render(cfgWithAnim)
			result2 := Render(cfgNoAnim)

			if result1 == nil || result2 == nil {
				// Both should be non-nil for valid configs; skip if nil (invalid).
				if result1 == nil && result2 == nil {
					return
				}
				t.Fatalf("nil mismatch: withAnim=%v, noAnim=%v", result1 == nil, result2 == nil)
			}

			// Compare pixel by pixel.
			for y := 0; y < h; y++ {
				for x := 0; x < w; x++ {
					r1, g1, b1, a1 := result1.Image.At(x, y).RGBA()
					r2, g2, b2, a2 := result2.Image.At(x, y).RGBA()
					if r1 != r2 || g1 != g2 || b1 != b2 || a1 != a2 {
						t.Fatalf("pixel mismatch at (%d,%d): NoAnim-config=(%d,%d,%d,%d), no-fields=(%d,%d,%d,%d) [style=%d, orient=%d, value=%f]",
							x, y, r1, g1, b1, a1, r2, g2, b2, a2, style, orientation, value)
					}
				}
			}
		})
	})

	t.Run("Pulse_with_period_le_zero_identical_to_no_animation", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			// Generate a config with Pulse animation but period ≤ 0 (disabled).
			style := Style(rapid.IntRange(0, 4).Draw(t, "style"))
			orientation := Orientation(rapid.IntRange(0, 1).Draw(t, "orientation"))

			var w, h int
			switch style {
			case Pie, Ring, Arc:
				w = rapid.IntRange(3, 50).Draw(t, "width")
				h = rapid.IntRange(3, 50).Draw(t, "height")
			default:
				w = rapid.IntRange(1, 50).Draw(t, "width")
				h = rapid.IntRange(1, 50).Draw(t, "height")
			}
			bounds := image.Rect(0, 0, w, h)
			value := rapid.Float64Range(0.0, 1.0).Draw(t, "value")

			fg := color.RGBA{
				R: uint8(rapid.IntRange(1, 255).Draw(t, "fgR")),
				G: uint8(rapid.IntRange(1, 255).Draw(t, "fgG")),
				B: uint8(rapid.IntRange(1, 255).Draw(t, "fgB")),
				A: 255,
			}
			bg := color.RGBA{
				R: uint8(rapid.IntRange(0, 254).Draw(t, "bgR")),
				G: uint8(rapid.IntRange(0, 254).Draw(t, "bgG")),
				B: uint8(rapid.IntRange(0, 254).Draw(t, "bgB")),
				A: 255,
			}

			// Period is ≤ 0 (negative or zero), which disables Pulse.
			period := time.Duration(rapid.Int64Range(-5_000_000_000, 0).Draw(t, "period"))
			elapsed := time.Duration(rapid.Int64Range(0, 5_000_000_000).Draw(t, "elapsed"))

			// Config with Pulse animation but disabled via period ≤ 0.
			cfgPulseDisabled := Config{
				Style:       style,
				Orientation: orientation,
				Value:       value,
				Bounds:      bounds,
				Foreground:  fg,
				Background:  bg,
				Animation: AnimationConfig{
					Type:   Pulse,
					Period: period,
				},
				animElapsed: elapsed,
			}

			// Config with no animation.
			cfgNoAnim := Config{
				Style:       style,
				Orientation: orientation,
				Value:       value,
				Bounds:      bounds,
				Foreground:  fg,
				Background:  bg,
			}

			result1 := Render(cfgPulseDisabled)
			result2 := Render(cfgNoAnim)

			if result1 == nil || result2 == nil {
				if result1 == nil && result2 == nil {
					return
				}
				t.Fatalf("nil mismatch: pulseDisabled=%v, noAnim=%v", result1 == nil, result2 == nil)
			}

			for y := 0; y < h; y++ {
				for x := 0; x < w; x++ {
					r1, g1, b1, a1 := result1.Image.At(x, y).RGBA()
					r2, g2, b2, a2 := result2.Image.At(x, y).RGBA()
					if r1 != r2 || g1 != g2 || b1 != b2 || a1 != a2 {
						t.Fatalf("pixel mismatch at (%d,%d): pulse-disabled=(%d,%d,%d,%d), no-anim=(%d,%d,%d,%d) [style=%d, orient=%d, value=%f, period=%v]",
							x, y, r1, g1, b1, a1, r2, g2, b2, a2, style, orientation, value, period)
					}
				}
			}
		})
	})

	t.Run("Shimmer_with_speed_le_zero_identical_to_no_animation", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			// Generate a config with Shimmer animation but speed < 0 (disabled).
			style := Style(rapid.IntRange(0, 4).Draw(t, "style"))
			orientation := Orientation(rapid.IntRange(0, 1).Draw(t, "orientation"))

			var w, h int
			switch style {
			case Pie, Ring, Arc:
				w = rapid.IntRange(3, 50).Draw(t, "width")
				h = rapid.IntRange(3, 50).Draw(t, "height")
			default:
				w = rapid.IntRange(1, 50).Draw(t, "width")
				h = rapid.IntRange(1, 50).Draw(t, "height")
			}
			bounds := image.Rect(0, 0, w, h)
			value := rapid.Float64Range(0.0, 1.0).Draw(t, "value")

			fg := color.RGBA{
				R: uint8(rapid.IntRange(1, 255).Draw(t, "fgR")),
				G: uint8(rapid.IntRange(1, 255).Draw(t, "fgG")),
				B: uint8(rapid.IntRange(1, 255).Draw(t, "fgB")),
				A: 255,
			}
			bg := color.RGBA{
				R: uint8(rapid.IntRange(0, 254).Draw(t, "bgR")),
				G: uint8(rapid.IntRange(0, 254).Draw(t, "bgG")),
				B: uint8(rapid.IntRange(0, 254).Draw(t, "bgB")),
				A: 255,
			}

			// Negative speed disables Shimmer.
			speed := rapid.IntRange(-500, -1).Draw(t, "speed")
			elapsed := time.Duration(rapid.Int64Range(0, 5_000_000_000).Draw(t, "elapsed"))

			cfgShimmerDisabled := Config{
				Style:       style,
				Orientation: orientation,
				Value:       value,
				Bounds:      bounds,
				Foreground:  fg,
				Background:  bg,
				Animation: AnimationConfig{
					Type:  Shimmer,
					Speed: speed,
				},
				animElapsed: elapsed,
			}

			cfgNoAnim := Config{
				Style:       style,
				Orientation: orientation,
				Value:       value,
				Bounds:      bounds,
				Foreground:  fg,
				Background:  bg,
			}

			result1 := Render(cfgShimmerDisabled)
			result2 := Render(cfgNoAnim)

			if result1 == nil || result2 == nil {
				if result1 == nil && result2 == nil {
					return
				}
				t.Fatalf("nil mismatch: shimmerDisabled=%v, noAnim=%v", result1 == nil, result2 == nil)
			}

			for y := 0; y < h; y++ {
				for x := 0; x < w; x++ {
					r1, g1, b1, a1 := result1.Image.At(x, y).RGBA()
					r2, g2, b2, a2 := result2.Image.At(x, y).RGBA()
					if r1 != r2 || g1 != g2 || b1 != b2 || a1 != a2 {
						t.Fatalf("pixel mismatch at (%d,%d): shimmer-disabled=(%d,%d,%d,%d), no-anim=(%d,%d,%d,%d) [style=%d, orient=%d, value=%f, speed=%d]",
							x, y, r1, g1, b1, a1, r2, g2, b2, a2, style, orientation, value, speed)
					}
				}
			}
		})
	})

	t.Run("MarchingStripes_with_speed_le_zero_identical_to_no_animation", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			// Generate a config with MarchingStripes animation but speed < 0 (disabled).
			// Note: MarchingStripes on Ring/Arc falls back to Shimmer, which also checks speed.
			style := Style(rapid.IntRange(0, 4).Draw(t, "style"))
			orientation := Orientation(rapid.IntRange(0, 1).Draw(t, "orientation"))

			var w, h int
			switch style {
			case Pie, Ring, Arc:
				w = rapid.IntRange(3, 50).Draw(t, "width")
				h = rapid.IntRange(3, 50).Draw(t, "height")
			default:
				w = rapid.IntRange(1, 50).Draw(t, "width")
				h = rapid.IntRange(1, 50).Draw(t, "height")
			}
			bounds := image.Rect(0, 0, w, h)
			value := rapid.Float64Range(0.0, 1.0).Draw(t, "value")

			fg := color.RGBA{
				R: uint8(rapid.IntRange(1, 255).Draw(t, "fgR")),
				G: uint8(rapid.IntRange(1, 255).Draw(t, "fgG")),
				B: uint8(rapid.IntRange(1, 255).Draw(t, "fgB")),
				A: 255,
			}
			bg := color.RGBA{
				R: uint8(rapid.IntRange(0, 254).Draw(t, "bgR")),
				G: uint8(rapid.IntRange(0, 254).Draw(t, "bgG")),
				B: uint8(rapid.IntRange(0, 254).Draw(t, "bgB")),
				A: 255,
			}

			// Negative speed disables MarchingStripes.
			speed := rapid.IntRange(-500, -1).Draw(t, "speed")
			elapsed := time.Duration(rapid.Int64Range(0, 5_000_000_000).Draw(t, "elapsed"))

			cfgStripesDisabled := Config{
				Style:       style,
				Orientation: orientation,
				Value:       value,
				Bounds:      bounds,
				Foreground:  fg,
				Background:  bg,
				Animation: AnimationConfig{
					Type:  MarchingStripes,
					Speed: speed,
				},
				animElapsed: elapsed,
			}

			cfgNoAnim := Config{
				Style:       style,
				Orientation: orientation,
				Value:       value,
				Bounds:      bounds,
				Foreground:  fg,
				Background:  bg,
			}

			result1 := Render(cfgStripesDisabled)
			result2 := Render(cfgNoAnim)

			if result1 == nil || result2 == nil {
				if result1 == nil && result2 == nil {
					return
				}
				t.Fatalf("nil mismatch: stripesDisabled=%v, noAnim=%v", result1 == nil, result2 == nil)
			}

			for y := 0; y < h; y++ {
				for x := 0; x < w; x++ {
					r1, g1, b1, a1 := result1.Image.At(x, y).RGBA()
					r2, g2, b2, a2 := result2.Image.At(x, y).RGBA()
					if r1 != r2 || g1 != g2 || b1 != b2 || a1 != a2 {
						t.Fatalf("pixel mismatch at (%d,%d): stripes-disabled=(%d,%d,%d,%d), no-anim=(%d,%d,%d,%d) [style=%d, orient=%d, value=%f, speed=%d]",
							x, y, r1, g1, b1, a1, r2, g2, b2, a2, style, orientation, value, speed)
					}
				}
			}
		})
	})
}

// TestProperty15_OutputDimensionsAndPositionMatchBounds verifies Property 15:
// For any valid configuration (positive bounds, recognized style, Ring/Arc ≥ 3×3),
// the returned Sprite SHALL have Image dimensions equal to Bounds.Dx() × Bounds.Dy()
// and Position equal to Bounds.Min.

func TestProperty15_OutputDimensionsAndPositionMatchBounds(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		style := Style(rapid.IntRange(0, 4).Draw(t, "style"))
		orientation := Orientation(rapid.IntRange(0, 1).Draw(t, "orientation"))

		// Generate non-zero Min point to exercise Position correctness
		minX := rapid.IntRange(-50, 200).Draw(t, "minX")
		minY := rapid.IntRange(-50, 200).Draw(t, "minY")

		var width, height int
		switch style {
		case Pie, Ring, Arc:
			width = rapid.IntRange(3, 100).Draw(t, "width")
			height = rapid.IntRange(3, 100).Draw(t, "height")
		default:
			width = rapid.IntRange(1, 100).Draw(t, "width")
			height = rapid.IntRange(1, 100).Draw(t, "height")
		}

		bounds := image.Rect(minX, minY, minX+width, minY+height)
		value := rapid.Float64Range(0.0, 1.0).Draw(t, "value")

		fg := color.RGBA{R: 200, G: 50, B: 100, A: 255}
		bg := color.RGBA{R: 10, G: 30, B: 60, A: 255}

		result := Render(Config{
			Style:       style,
			Orientation: orientation,
			Value:       value,
			Bounds:      bounds,
			Foreground:  fg,
			Background:  bg,
		})

		if result == nil {
			t.Fatal("expected non-nil result for valid configuration")
		}

		// Assert Image dimensions equal Bounds.Dx() × Bounds.Dy()
		gotWidth := result.Image.Bounds().Dx()
		gotHeight := result.Image.Bounds().Dy()
		if gotWidth != width {
			t.Fatalf("image width mismatch: got %d, want %d (style=%d, orient=%d, bounds=%v)",
				gotWidth, width, style, orientation, bounds)
		}
		if gotHeight != height {
			t.Fatalf("image height mismatch: got %d, want %d (style=%d, orient=%d, bounds=%v)",
				gotHeight, height, style, orientation, bounds)
		}

		// Assert Position equals Bounds.Min
		if result.Position.X != minX || result.Position.Y != minY {
			t.Fatalf("position mismatch: got (%d, %d), want (%d, %d) (style=%d, orient=%d, bounds=%v)",
				result.Position.X, result.Position.Y, minX, minY, style, orientation, bounds)
		}
	})
}

// TestProperty18_Value0ProducesNoFillPixels verifies Property 18:
// For any valid configuration with Value = 0.0 (and no animation, no gradient),
// the output SHALL contain zero pixels matching the foreground fill color within
// the interior (excluding border and marker pixels).
//
// We test Linear and Segmented styles (no border, no markers, no animation, no gradient)
// to ensure no foreground pixels exist when Value = 0.0.

func TestProperty18_Value0ProducesNoFillPixels(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Only test Linear and Segmented — circular styles have background
		// outside the circle that complicates interior detection.
		style := Style(rapid.SampledFrom([]Style{Linear, Segmented}).Draw(t, "style"))
		orientation := Orientation(rapid.IntRange(0, 1).Draw(t, "orientation"))

		width := rapid.IntRange(1, 100).Draw(t, "width")
		height := rapid.IntRange(1, 100).Draw(t, "height")
		bounds := image.Rect(0, 0, width, height)

		// Use distinct, non-zero colors so we can identify fg vs bg
		fg := color.RGBA{R: 200, G: 50, B: 100, A: 255}
		bg := color.RGBA{R: 10, G: 30, B: 60, A: 255}

		cfg := Config{
			Style:       style,
			Orientation: orientation,
			Value:       0.0,
			Bounds:      bounds,
			Foreground:  fg,
			Background:  bg,
			// No gradient, no animation, no border, no markers
		}

		result := Render(cfg)
		if result == nil {
			t.Fatal("expected non-nil result for valid configuration")
		}

		// Scan all pixels: none should match the foreground color
		fgR, fgG, fgB, fgA := fg.RGBA()
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				r, g, b, a := result.Image.At(x, y).RGBA()
				if r == fgR && g == fgG && b == fgB && a == fgA {
					t.Fatalf("found foreground pixel at (%d,%d) with Value=0.0 "+
						"[style=%d, orient=%d, bounds=%dx%d]",
						x, y, style, orientation, width, height)
				}
			}
		}
	})
}

// TestProperty19_Value1ProducesNoTrackPixels verifies Property 19:
// For any valid configuration with Value = 1.0 (and no gradient), the output SHALL
// contain zero pixels matching the background/track color within the interior
// (excluding border and marker pixels).
//
// For Linear style: entire image should be foreground (no background).
// For Segmented style: all cell pixels should be foreground; only gap pixels
// may be background-colored (gaps are structural, not track).
// Circular styles are excluded because pixels outside the inscribed circle are
// background by design.

func TestProperty19_Value1ProducesNoTrackPixels(t *testing.T) {
	t.Run("linear", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			orientation := Orientation(rapid.IntRange(0, 1).Draw(t, "orientation"))

			width := rapid.IntRange(1, 100).Draw(t, "width")
			height := rapid.IntRange(1, 100).Draw(t, "height")
			bounds := image.Rect(0, 0, width, height)

			// Use distinct, non-zero colors so we can identify fg vs bg
			fg := color.RGBA{R: 200, G: 50, B: 100, A: 255}
			bg := color.RGBA{R: 10, G: 30, B: 60, A: 255}

			cfg := Config{
				Style:       Linear,
				Orientation: orientation,
				Value:       1.0,
				Bounds:      bounds,
				Foreground:  fg,
				Background:  bg,
				// No gradient, no border, no markers
			}

			result := Render(cfg)
			if result == nil {
				t.Fatal("expected non-nil result for valid configuration")
			}

			// Scan all pixels: none should match the background/track color
			bgR, bgG, bgB, bgA := bg.RGBA()
			for y := 0; y < height; y++ {
				for x := 0; x < width; x++ {
					r, g, b, a := result.Image.At(x, y).RGBA()
					if r == bgR && g == bgG && b == bgB && a == bgA {
						t.Fatalf("found background/track pixel at (%d,%d) with Value=1.0 "+
							"[style=Linear, orient=%d, bounds=%dx%d]",
							x, y, orientation, width, height)
					}
				}
			}
		})
	})

	t.Run("segmented_cells_all_filled", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			orientation := Orientation(rapid.IntRange(0, 1).Draw(t, "orientation"))

			// Use large enough bounds to accommodate segments
			primaryLen := rapid.IntRange(20, 100).Draw(t, "primaryLen")
			minorLen := rapid.IntRange(2, 30).Draw(t, "minorLen")

			var bounds image.Rectangle
			if orientation == OrientHorizontal {
				bounds = image.Rect(0, 0, primaryLen, minorLen)
			} else {
				bounds = image.Rect(0, 0, minorLen, primaryLen)
			}

			gap := rapid.IntRange(1, 4).Draw(t, "gap")
			segCount := rapid.IntRange(2, 10).Draw(t, "segCount")

			cellWidth := (primaryLen - (segCount-1)*gap) / segCount
			if cellWidth < 1 {
				return
			}
			if primaryLen < 2*cellWidth+gap {
				return
			}

			// Use distinct, non-zero colors
			fg := color.RGBA{R: 200, G: 50, B: 100, A: 255}
			bg := color.RGBA{R: 10, G: 30, B: 60, A: 255}

			cfg := Config{
				Style:        Segmented,
				Orientation:  orientation,
				Value:        1.0,
				Bounds:       bounds,
				Foreground:   fg,
				Background:   bg,
				SegmentCount: segCount,
				SegmentGap:   gap,
				// No gradient, no border, no markers
			}

			result := Render(cfg)
			if result == nil {
				t.Fatal("expected non-nil result for valid configuration")
			}

			// For Segmented at Value=1.0: all cell pixels should be foreground.
			// Gap pixels and remainder pixels are allowed to be background.
			// Check only within cell regions for absence of track color.
			fgR, fgG, fgB, fgA := fg.RGBA()
			for i := 0; i < segCount; i++ {
				cellStart := i * (cellWidth + gap)
				cellEnd := cellStart + cellWidth

				for pos := cellStart; pos < cellEnd; pos++ {
					for m := 0; m < minorLen; m++ {
						var x, y int
						if orientation == OrientHorizontal {
							x = pos
							y = m
						} else {
							x = m
							y = bounds.Dy() - 1 - pos
						}

						if x < 0 || x >= bounds.Dx() || y < 0 || y >= bounds.Dy() {
							continue
						}

						r, g, b, a := result.Image.At(x, y).RGBA()
						if r != fgR || g != fgG || b != fgB || a != fgA {
							t.Fatalf("cell %d pixel at primary=%d, minor=%d (x=%d, y=%d): "+
								"expected fg color with Value=1.0, got (%d,%d,%d,%d) "+
								"[orient=%d, segCount=%d, gap=%d, cellWidth=%d]",
								i, pos, m, x, y, r>>8, g>>8, b>>8, a>>8,
								orientation, segCount, gap, cellWidth)
						}
					}
				}
			}
		})
	})
}

// =============================================================================
// Widget Wrapper Integration Tests (Task 11.3)
// =============================================================================

// TestWidgetNew verifies that New returns a non-nil widgets.Renderable
// and that the returned value also satisfies Configurable, Animated, and Described.
//
// _Requirements: 3.4, 8.4, 10.1_
func TestWidgetNew(t *testing.T) {
	cfg := Config{
		Style:      Linear,
		Value:      0.5,
		Bounds:     image.Rect(0, 0, 20, 10),
		Foreground: color.RGBA{R: 255, G: 0, B: 0, A: 255},
		Background: color.RGBA{R: 0, G: 0, B: 255, A: 255},
	}

	w := New(cfg)
	if w == nil {
		t.Fatal("New returned nil")
	}

	// Verify that it implements widgets.Renderable (RenderFrame)
	sprite := w.RenderFrame()
	if sprite == nil {
		t.Fatal("RenderFrame returned nil for valid config")
	}

	// Verify Described interface
	pw := w.(*progressbarWidget)
	desc := pw.Describe()
	if desc.Name != "progressbar" {
		t.Errorf("Describe().Name = %q, want %q", desc.Name, "progressbar")
	}
}

// TestWidgetConfigure verifies that Configure updates the widget's internal config
// so that subsequent RenderFrame calls use the new parameters.
//
// _Requirements: 8.4_
func TestWidgetConfigure(t *testing.T) {
	cfg1 := Config{
		Style:      Linear,
		Value:      0.5,
		Bounds:     image.Rect(0, 0, 20, 10),
		Foreground: color.RGBA{R: 255, G: 0, B: 0, A: 255},
		Background: color.RGBA{R: 0, G: 0, B: 255, A: 255},
	}
	cfg2 := Config{
		Style:      Pie,
		Value:      0.8,
		Bounds:     image.Rect(0, 0, 30, 30),
		Foreground: color.RGBA{R: 0, G: 255, B: 0, A: 255},
		Background: color.RGBA{R: 0, G: 0, B: 0, A: 255},
	}

	w := New(cfg1)
	sprite1 := w.RenderFrame()
	if sprite1 == nil {
		t.Fatal("RenderFrame returned nil before Configure")
	}

	// Configure with new config
	pw := w.(*progressbarWidget)
	pw.Configure(cfg2)

	sprite2 := w.RenderFrame()
	if sprite2 == nil {
		t.Fatal("RenderFrame returned nil after Configure")
	}

	// Verify the output changed (Pie style renders differently than Linear)
	if sprite2.Label != "progressbar/pie" {
		t.Errorf("after Configure to Pie style, label = %q, want %q", sprite2.Label, "progressbar/pie")
	}

	// Verify dimensions changed (30×30 vs 20×10)
	if sprite2.Image.Bounds().Dx() != 30 || sprite2.Image.Bounds().Dy() != 30 {
		t.Errorf("after Configure, image bounds = %v, want 30×30", sprite2.Image.Bounds())
	}
}

// TestWidgetRenderFrame verifies that RenderFrame produces correct output
// consistent with calling Render directly on the same config.
//
// _Requirements: 8.4_
func TestWidgetRenderFrame(t *testing.T) {
	cfg := Config{
		Style:      Linear,
		Value:      0.7,
		Bounds:     image.Rect(5, 10, 25, 20),
		Foreground: color.RGBA{R: 200, G: 100, B: 50, A: 255},
		Background: color.RGBA{R: 10, G: 20, B: 30, A: 255},
	}

	w := New(cfg)
	sprite := w.RenderFrame()
	if sprite == nil {
		t.Fatal("RenderFrame returned nil")
	}

	// Verify the sprite position matches bounds.Min
	if sprite.Position.X != 5 || sprite.Position.Y != 10 {
		t.Errorf("Position = %v, want (5, 10)", sprite.Position)
	}

	// Verify sprite dimensions match bounds
	if sprite.Image.Bounds().Dx() != 20 || sprite.Image.Bounds().Dy() != 10 {
		t.Errorf("Image dimensions = %d×%d, want 20×10",
			sprite.Image.Bounds().Dx(), sprite.Image.Bounds().Dy())
	}

	// Verify label
	if sprite.Label != "progressbar/linear" {
		t.Errorf("Label = %q, want %q", sprite.Label, "progressbar/linear")
	}
}

// TestWidgetTick verifies the animation state machine:
// - Pulse animation: Tick advances animElapsed
// - NoAnimation: Tick does not change animElapsed
//
// _Requirements: 3.4_
func TestWidgetTick(t *testing.T) {
	t.Run("Pulse_advances_animElapsed", func(t *testing.T) {
		cfg := Config{
			Style:  Linear,
			Value:  0.5,
			Bounds: image.Rect(0, 0, 20, 10),
			Animation: AnimationConfig{
				Type:   Pulse,
				Period: time.Second,
			},
		}

		w := New(cfg)
		pw := w.(*progressbarWidget)

		// Initially animElapsed should be zero
		if pw.cfg.animElapsed != 0 {
			t.Fatalf("initial animElapsed = %v, want 0", pw.cfg.animElapsed)
		}

		// Tick with 100ms
		pw.Tick(100 * time.Millisecond)
		if pw.cfg.animElapsed != 100*time.Millisecond {
			t.Fatalf("after Tick(100ms), animElapsed = %v, want 100ms", pw.cfg.animElapsed)
		}

		// Tick again with 50ms
		pw.Tick(50 * time.Millisecond)
		if pw.cfg.animElapsed != 150*time.Millisecond {
			t.Fatalf("after Tick(50ms), animElapsed = %v, want 150ms", pw.cfg.animElapsed)
		}
	})

	t.Run("NoAnimation_skips_Tick", func(t *testing.T) {
		cfg := Config{
			Style:  Linear,
			Value:  0.5,
			Bounds: image.Rect(0, 0, 20, 10),
			Animation: AnimationConfig{
				Type: NoAnimation,
			},
		}

		w := New(cfg)
		pw := w.(*progressbarWidget)

		// Tick should be a no-op
		pw.Tick(100 * time.Millisecond)
		if pw.cfg.animElapsed != 0 {
			t.Fatalf("NoAnimation: after Tick(100ms), animElapsed = %v, want 0", pw.cfg.animElapsed)
		}

		pw.Tick(500 * time.Millisecond)
		if pw.cfg.animElapsed != 0 {
			t.Fatalf("NoAnimation: after Tick(500ms), animElapsed = %v, want 0", pw.cfg.animElapsed)
		}
	})
}

// TestWidgetCaching verifies the render cache integration:
// - Same config → cached sprite (same pointer)
// - Different config via Configure → re-render (different result)
// - Orientation change → cache invalidation
// - Style change → cache invalidation
//
// _Requirements: 8.4, 9.10_
func TestWidgetCaching(t *testing.T) {
	t.Run("same_config_returns_cached_sprite", func(t *testing.T) {
		cfg := Config{
			Style:      Linear,
			Value:      0.5,
			Bounds:     image.Rect(0, 0, 20, 10),
			Foreground: color.RGBA{R: 255, G: 0, B: 0, A: 255},
			Background: color.RGBA{R: 0, G: 0, B: 255, A: 255},
		}

		w := New(cfg, widgets.WithCaching())
		sprite1 := w.RenderFrame()
		sprite2 := w.RenderFrame()

		if sprite1 == nil || sprite2 == nil {
			t.Fatal("RenderFrame returned nil")
		}

		// Same config should return the same cached pointer
		if sprite1 != sprite2 {
			t.Error("expected same pointer for identical config with caching enabled")
		}
	})

	t.Run("different_config_returns_new_sprite", func(t *testing.T) {
		cfg := Config{
			Style:      Linear,
			Value:      0.5,
			Bounds:     image.Rect(0, 0, 20, 10),
			Foreground: color.RGBA{R: 255, G: 0, B: 0, A: 255},
			Background: color.RGBA{R: 0, G: 0, B: 255, A: 255},
		}

		w := New(cfg, widgets.WithCaching())
		sprite1 := w.RenderFrame()

		// Change config via Configure
		pw := w.(*progressbarWidget)
		pw.Configure(Config{
			Style:      Linear,
			Value:      0.8,
			Bounds:     image.Rect(0, 0, 20, 10),
			Foreground: color.RGBA{R: 255, G: 0, B: 0, A: 255},
			Background: color.RGBA{R: 0, G: 0, B: 255, A: 255},
		})
		sprite2 := w.RenderFrame()

		if sprite1 == nil || sprite2 == nil {
			t.Fatal("RenderFrame returned nil")
		}

		// Different config should produce a different sprite
		if sprite1 == sprite2 {
			t.Error("expected different pointer after Configure with different value")
		}
	})

	t.Run("orientation_change_invalidates_cache", func(t *testing.T) {
		cfg := Config{
			Style:       Linear,
			Orientation: OrientHorizontal,
			Value:       0.5,
			Bounds:      image.Rect(0, 0, 20, 20),
			Foreground:  color.RGBA{R: 255, G: 0, B: 0, A: 255},
			Background:  color.RGBA{R: 0, G: 0, B: 255, A: 255},
		}

		w := New(cfg, widgets.WithCaching())
		sprite1 := w.RenderFrame()

		// Change only the orientation
		pw := w.(*progressbarWidget)
		newCfg := cfg
		newCfg.Orientation = OrientVertical
		pw.Configure(newCfg)
		sprite2 := w.RenderFrame()

		if sprite1 == nil || sprite2 == nil {
			t.Fatal("RenderFrame returned nil")
		}

		// Orientation change should invalidate the cache
		if sprite1 == sprite2 {
			t.Error("expected different pointer after orientation change")
		}
	})

	t.Run("style_change_invalidates_cache", func(t *testing.T) {
		cfg := Config{
			Style:      Linear,
			Value:      0.5,
			Bounds:     image.Rect(0, 0, 20, 20),
			Foreground: color.RGBA{R: 255, G: 0, B: 0, A: 255},
			Background: color.RGBA{R: 0, G: 0, B: 255, A: 255},
		}

		w := New(cfg, widgets.WithCaching())
		sprite1 := w.RenderFrame()

		// Change only the style
		pw := w.(*progressbarWidget)
		newCfg := cfg
		newCfg.Style = Pie
		pw.Configure(newCfg)
		sprite2 := w.RenderFrame()

		if sprite1 == nil || sprite2 == nil {
			t.Fatal("RenderFrame returned nil")
		}

		// Style change should invalidate the cache
		if sprite1 == sprite2 {
			t.Error("expected different pointer after style change")
		}
	})
}

// TestWidgetLabels verifies that each of the five styles produces the correct
// sprite label string when rendered through the widget wrapper.
//
// _Requirements: 10.1_
func TestWidgetLabels(t *testing.T) {
	tests := []struct {
		style Style
		label string
	}{
		{Linear, "progressbar/linear"},
		{Pie, "progressbar/pie"},
		{Segmented, "progressbar/segmented"},
		{Ring, "progressbar/ring"},
		{Arc, "progressbar/arc"},
	}

	for _, tt := range tests {
		t.Run(tt.label, func(t *testing.T) {
			cfg := Config{
				Style:      tt.style,
				Value:      0.5,
				Bounds:     image.Rect(0, 0, 20, 20),
				Foreground: color.RGBA{R: 255, G: 0, B: 0, A: 255},
				Background: color.RGBA{R: 0, G: 0, B: 255, A: 255},
			}

			w := New(cfg)
			sprite := w.RenderFrame()
			if sprite == nil {
				t.Fatal("RenderFrame returned nil")
			}

			if sprite.Label != tt.label {
				t.Errorf("Label = %q, want %q", sprite.Label, tt.label)
			}
		})
	}
}

// TestWidgetDefaultOrientation verifies that a zero-value Orientation field
// in the widget config produces horizontal output (identical to explicitly
// setting OrientHorizontal).
//
// _Requirements: 9.10_
func TestWidgetDefaultOrientation(t *testing.T) {
	// Config with zero-value Orientation (default = horizontal)
	cfgDefault := Config{
		Style:      Linear,
		Value:      0.6,
		Bounds:     image.Rect(0, 0, 20, 10),
		Foreground: color.RGBA{R: 255, G: 0, B: 0, A: 255},
		Background: color.RGBA{R: 0, G: 0, B: 255, A: 255},
	}

	// Config with explicit OrientHorizontal
	cfgExplicit := Config{
		Style:       Linear,
		Orientation: OrientHorizontal,
		Value:       0.6,
		Bounds:      image.Rect(0, 0, 20, 10),
		Foreground:  color.RGBA{R: 255, G: 0, B: 0, A: 255},
		Background:  color.RGBA{R: 0, G: 0, B: 255, A: 255},
	}

	w1 := New(cfgDefault)
	w2 := New(cfgExplicit)

	sprite1 := w1.RenderFrame()
	sprite2 := w2.RenderFrame()

	if sprite1 == nil || sprite2 == nil {
		t.Fatal("RenderFrame returned nil")
	}

	// Verify both produce the same image dimensions
	if sprite1.Image.Bounds().Dx() != sprite2.Image.Bounds().Dx() ||
		sprite1.Image.Bounds().Dy() != sprite2.Image.Bounds().Dy() {
		t.Fatal("dimension mismatch between default and explicit horizontal orientation")
	}

	// Verify pixel-identical output
	w := sprite1.Image.Bounds().Dx()
	h := sprite1.Image.Bounds().Dy()
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r1, g1, b1, a1 := sprite1.Image.At(x, y).RGBA()
			r2, g2, b2, a2 := sprite2.Image.At(x, y).RGBA()
			if r1 != r2 || g1 != g2 || b1 != b2 || a1 != a2 {
				t.Fatalf("pixel mismatch at (%d,%d): default=(%d,%d,%d,%d), explicit=(%d,%d,%d,%d)",
					x, y, r1, g1, b1, a1, r2, g2, b2, a2)
			}
		}
	}
}
