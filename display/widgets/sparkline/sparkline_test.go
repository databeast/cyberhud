package sparkline

import (
	"image"
	"image/color"
	"math"
	"testing"

	"pgregory.net/rapid"
)

// --- From: sparkline_prop_test.go ---

// TestPropertySparklineOutputMetadataCorrectness verifies that for any valid sparkline Config
// (bounds width ≥ 1, height ≥ 1), the Sparkline_Widget returns a non-nil Result where Image
// dimensions equal Bounds.Dx() × Bounds.Dy(), Position equals Bounds.Min, and Label matches style.
//

func TestPropertySparklineOutputMetadataCorrectness(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate valid bounds with non-zero offset
		minX := rapid.IntRange(0, 50).Draw(t, "minX")
		minY := rapid.IntRange(0, 50).Draw(t, "minY")
		width := rapid.IntRange(1, 100).Draw(t, "width")
		height := rapid.IntRange(1, 100).Draw(t, "height")
		bounds := image.Rect(minX, minY, minX+width, minY+height)

		// Generate random data (0-20 points)
		dataLen := rapid.IntRange(0, 20).Draw(t, "dataLen")
		data := make([]float64, dataLen)
		for i := range data {
			data[i] = rapid.Float64Range(0.0, 1.0).Draw(t, "dataPoint")
		}

		// Generate random style
		style := Style(rapid.IntRange(0, 1).Draw(t, "style"))

		fg := color.RGBA{R: 200, G: 100, B: 50, A: 255}
		bg := color.RGBA{R: 10, G: 20, B: 30, A: 255}

		result := Render(Config{
			Data:       data,
			Style:      style,
			Bounds:     bounds,
			Foreground: fg,
			Background: bg,
		})

		if result == nil {
			t.Fatal("expected non-nil result for valid bounds")
		}

		// Verify image dimensions
		gotWidth := result.Image.Bounds().Dx()
		gotHeight := result.Image.Bounds().Dy()
		if gotWidth != width {
			t.Fatalf("image width mismatch: got %d, want %d", gotWidth, width)
		}
		if gotHeight != height {
			t.Fatalf("image height mismatch: got %d, want %d", gotHeight, height)
		}

		// Verify Position == Bounds.Min
		if result.Position.X != minX || result.Position.Y != minY {
			t.Fatalf("Position mismatch: got (%d, %d), want (%d, %d)",
				result.Position.X, result.Position.Y, minX, minY)
		}

		// Verify Label matches style
		var expectedLabel string
		switch style {
		case Line:
			expectedLabel = "sparkline/line"
		case Bar:
			expectedLabel = "sparkline/bar"
		}
		if result.Label != expectedLabel {
			t.Fatalf("Label mismatch: got %q, want %q", result.Label, expectedLabel)
		}
	})
}

// TestPropertySparklineBarFillHeight verifies that for Bar-style sparkline, each bar column's
// foreground pixel count from the bottom equals floor(Bounds.Dy() × clamp(value, 0, 1)).
//

func TestPropertySparklineBarFillHeight(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate 1-10 data points in [0,1]
		count := rapid.IntRange(1, 10).Draw(t, "count")
		data := make([]float64, count)
		for i := range data {
			data[i] = rapid.Float64Range(0.0, 1.0).Draw(t, "dataPoint")
		}

		// Generate bounds wider than data count to ensure barWidth >= 1
		height := rapid.IntRange(2, 50).Draw(t, "height")
		width := rapid.IntRange(count, count*10).Draw(t, "width")
		bounds := image.Rect(0, 0, width, height)

		fg := color.RGBA{R: 200, G: 100, B: 50, A: 255}
		bg := color.RGBA{R: 10, G: 20, B: 30, A: 255}

		result := Render(Config{
			Data:       data,
			Style:      Bar,
			Bounds:     bounds,
			Foreground: fg,
			Background: bg,
		})

		if result == nil {
			t.Fatal("expected non-nil result for valid bounds")
		}

		h := height
		w := width
		barWidth := w / count
		if barWidth < 1 {
			barWidth = 1
		}

		// For each bar, pick a representative column and count foreground pixels
		for barIdx := 0; barIdx < count; barIdx++ {
			// Pick the first column of this bar as representative
			col := barIdx * barWidth
			if col >= w {
				break
			}

			expectedBarHeight := int(float64(h) * data[barIdx])

			// Count foreground pixels from bottom in this column
			fgCount := 0
			for y := 0; y < h; y++ {
				pixel := result.Image.At(col, y)
				r, g, b, a := pixel.RGBA()
				// Compare with fg (200,100,50,255) → pre-multiplied 16-bit
				fgR, fgG, fgB, fgA := fg.RGBA()
				if r == fgR && g == fgG && b == fgB && a == fgA {
					fgCount++
				}
			}

			if fgCount != expectedBarHeight {
				t.Fatalf("bar %d (col %d): foreground pixel count = %d, want %d (value=%f, h=%d, w=%d, barWidth=%d)",
					barIdx, col, fgCount, expectedBarHeight, data[barIdx], h, w, barWidth)
			}
		}
	})
}

// TestPropertySparklineLineFillHeight verifies that for Line-style sparkline, each column's
// foreground pixel count matches the interpolated value (±1 pixel tolerance).
//

func TestPropertySparklineLineFillHeight(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate 2-10 data points in [0,1]
		count := rapid.IntRange(2, 10).Draw(t, "count")
		data := make([]float64, count)
		for i := range data {
			data[i] = rapid.Float64Range(0.0, 1.0).Draw(t, "dataPoint")
		}

		height := rapid.IntRange(3, 50).Draw(t, "height")
		width := rapid.IntRange(count, 100).Draw(t, "width")
		bounds := image.Rect(0, 0, width, height)

		fg := color.RGBA{R: 200, G: 100, B: 50, A: 255}
		bg := color.RGBA{R: 10, G: 20, B: 30, A: 255}

		result := Render(Config{
			Data:       data,
			Style:      Line,
			Bounds:     bounds,
			Foreground: fg,
			Background: bg,
		})

		if result == nil {
			t.Fatal("expected non-nil result for valid bounds")
		}

		h := height
		w := width

		// Compute x positions for each data point (same formula as implementation)
		xPositions := make([]int, count)
		for i := 0; i < count; i++ {
			xPositions[i] = i * (w - 1) / (count - 1)
		}

		fgR, fgG, fgB, fgA := fg.RGBA()

		// For each pixel column, verify foreground count matches interpolated value
		for x := 0; x < w; x++ {
			// Compute expected interpolated value at this column
			interpValue := computeInterpolation(x, xPositions, data, count)
			expectedFillHeight := int(float64(h) * interpValue)

			// Count foreground pixels in this column
			fgCount := 0
			for y := 0; y < h; y++ {
				pixel := result.Image.At(x, y)
				r, g, b, a := pixel.RGBA()
				if r == fgR && g == fgG && b == fgB && a == fgA {
					fgCount++
				}
			}

			// Allow ±1 pixel tolerance for rounding
			if abs(fgCount-expectedFillHeight) > 1 {
				t.Fatalf("column %d: foreground pixel count = %d, expected %d (±1 tolerance), interpValue=%f, h=%d",
					x, fgCount, expectedFillHeight, interpValue, h)
			}
		}
	})
}

// TestPropertySparklineTruncationOfExcessDataPoints verifies that when len(Data) > Bounds.Dx(),
// the output is pixel-identical to rendering with only the last Bounds.Dx() data points.
//

func TestPropertySparklineTruncationOfExcessDataPoints(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate bounds
		width := rapid.IntRange(2, 20).Draw(t, "width")
		height := rapid.IntRange(2, 30).Draw(t, "height")
		bounds := image.Rect(0, 0, width, height)

		// Generate data with more points than width
		excess := rapid.IntRange(1, 20).Draw(t, "excess")
		totalPoints := width + excess
		data := make([]float64, totalPoints)
		for i := range data {
			data[i] = rapid.Float64Range(0.0, 1.0).Draw(t, "dataPoint")
		}

		// Only the last `width` points should be used
		truncatedData := data[len(data)-width:]

		style := Style(rapid.IntRange(0, 1).Draw(t, "style"))

		fg := color.RGBA{R: 200, G: 100, B: 50, A: 255}
		bg := color.RGBA{R: 10, G: 20, B: 30, A: 255}

		// Render with full data
		result1 := Render(Config{
			Data:       data,
			Style:      style,
			Bounds:     bounds,
			Foreground: fg,
			Background: bg,
		})

		// Render with truncated data
		result2 := Render(Config{
			Data:       truncatedData,
			Style:      style,
			Bounds:     bounds,
			Foreground: fg,
			Background: bg,
		})

		if result1 == nil || result2 == nil {
			t.Fatal("expected non-nil results for valid bounds")
		}

		// Compare pixel-by-pixel
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				r1, g1, b1, a1 := result1.Image.At(x, y).RGBA()
				r2, g2, b2, a2 := result2.Image.At(x, y).RGBA()
				if r1 != r2 || g1 != g2 || b1 != b2 || a1 != a2 {
					t.Fatalf("pixel mismatch at (%d,%d): full=(%d,%d,%d,%d), truncated=(%d,%d,%d,%d) [style=%d, totalPoints=%d, width=%d]",
						x, y, r1, g1, b1, a1, r2, g2, b2, a2, style, totalPoints, width)
				}
			}
		}
	})
}

// TestPropertySparklineDataClampingIdempotence verifies that out-of-range and NaN data values
// produce identical output to manually clamped values.
//

func TestPropertySparklineDataClampingIdempotence(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		width := rapid.IntRange(2, 30).Draw(t, "width")
		height := rapid.IntRange(2, 30).Draw(t, "height")
		bounds := image.Rect(0, 0, width, height)

		// Generate data with values outside [0,1] and NaN
		count := rapid.IntRange(1, 10).Draw(t, "count")
		data := make([]float64, count)
		clampedData := make([]float64, count)

		for i := range data {
			// Generate a mix: some in-range, some out-of-range, some NaN
			kind := rapid.IntRange(0, 4).Draw(t, "kind")
			switch kind {
			case 0: // negative
				data[i] = rapid.Float64Range(-100.0, -0.001).Draw(t, "negValue")
				clampedData[i] = 0.0
			case 1: // above 1
				data[i] = rapid.Float64Range(1.001, 100.0).Draw(t, "highValue")
				clampedData[i] = 1.0
			case 2: // NaN
				data[i] = math.NaN()
				clampedData[i] = 0.0
			default: // valid [0,1]
				v := rapid.Float64Range(0.0, 1.0).Draw(t, "validValue")
				data[i] = v
				clampedData[i] = v
			}
		}

		style := Style(rapid.IntRange(0, 1).Draw(t, "style"))

		fg := color.RGBA{R: 200, G: 100, B: 50, A: 255}
		bg := color.RGBA{R: 10, G: 20, B: 30, A: 255}

		// Render with raw data
		result1 := Render(Config{
			Data:       data,
			Style:      style,
			Bounds:     bounds,
			Foreground: fg,
			Background: bg,
		})

		// Render with manually clamped data
		result2 := Render(Config{
			Data:       clampedData,
			Style:      style,
			Bounds:     bounds,
			Foreground: fg,
			Background: bg,
		})

		if result1 == nil || result2 == nil {
			t.Fatal("expected non-nil results for valid bounds")
		}

		// Compare pixel-by-pixel
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				r1, g1, b1, a1 := result1.Image.At(x, y).RGBA()
				r2, g2, b2, a2 := result2.Image.At(x, y).RGBA()
				if r1 != r2 || g1 != g2 || b1 != b2 || a1 != a2 {
					t.Fatalf("pixel mismatch at (%d,%d): raw=(%d,%d,%d,%d), clamped=(%d,%d,%d,%d) [style=%d]",
						x, y, r1, g1, b1, a1, r2, g2, b2, a2, style)
				}
			}
		}
	})
}

// TestPropertySparklineColorExclusivity verifies that every pixel in the output image is either
// the specified Foreground_Color or the specified Background_Color (no other colors present).
//

func TestPropertySparklineColorExclusivity(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		width := rapid.IntRange(1, 50).Draw(t, "width")
		height := rapid.IntRange(1, 50).Draw(t, "height")
		bounds := image.Rect(0, 0, width, height)

		// Generate random data
		dataLen := rapid.IntRange(0, 15).Draw(t, "dataLen")
		data := make([]float64, dataLen)
		for i := range data {
			data[i] = rapid.Float64Range(0.0, 1.0).Draw(t, "dataPoint")
		}

		style := Style(rapid.IntRange(0, 1).Draw(t, "style"))

		// Use explicit non-zero colors that are different from each other
		fg := color.RGBA{R: 200, G: 100, B: 50, A: 255}
		bg := color.RGBA{R: 10, G: 20, B: 30, A: 255}

		result := Render(Config{
			Data:       data,
			Style:      style,
			Bounds:     bounds,
			Foreground: fg,
			Background: bg,
		})

		if result == nil {
			t.Fatal("expected non-nil result for valid bounds")
		}

		fgR, fgG, fgB, fgA := fg.RGBA()
		bgR, bgG, bgB, bgA := bg.RGBA()

		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				r, g, b, a := result.Image.At(x, y).RGBA()
				isFg := r == fgR && g == fgG && b == fgB && a == fgA
				isBg := r == bgR && g == bgG && b == bgB && a == bgA
				if !isFg && !isBg {
					t.Fatalf("pixel (%d,%d) is neither fg nor bg: got (%d,%d,%d,%d), fg=(%d,%d,%d,%d), bg=(%d,%d,%d,%d)",
						x, y, r, g, b, a, fgR, fgG, fgB, fgA, bgR, bgG, bgB, bgA)
				}
			}
		}
	})
}

// --- Helper functions for property tests ---

// computeInterpolation computes the linearly interpolated data value at pixel column x.
// This mirrors the implementation's interpolateAtX logic.
func computeInterpolation(x int, xPositions []int, data []float64, count int) float64 {
	for i := 0; i < count-1; i++ {
		x0 := xPositions[i]
		x1 := xPositions[i+1]
		if x >= x0 && x <= x1 {
			if x0 == x1 {
				return data[i]
			}
			t := float64(x-x0) / float64(x1-x0)
			return data[i] + t*(data[i+1]-data[i])
		}
	}
	return data[count-1]
}

// abs returns the absolute value of an integer.
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// --- From: sparkline_test.go ---

// TestRenderNilForInvalidBounds verifies that Render returns nil when
// given bounds with zero or negative width/height.
func TestRenderNilForInvalidBounds(t *testing.T) {
	tests := []struct {
		name   string
		bounds image.Rectangle
	}{
		{"zero width", image.Rect(0, 0, 0, 10)},
		{"zero height", image.Rect(0, 0, 5, 0)},
		{"both zero", image.Rect(0, 0, 0, 0)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Config{
				Data:   []float64{0.5},
				Style:  Line,
				Bounds: tt.bounds,
			}
			result := Render(cfg)
			if result != nil {
				t.Errorf("expected nil for bounds %v, got non-nil result", tt.bounds)
			}
		})
	}
}

// TestEmptyDataProducesBackgroundOnly verifies that an empty data slice
// produces an image filled entirely with the background color.

func TestEmptyDataProducesBackgroundOnly(t *testing.T) {
	bounds := image.Rect(0, 0, 8, 6)
	bg := color.RGBA{R: 0, G: 0, B: 0, A: 255}

	styles := []struct {
		name  string
		style Style
	}{
		{"Line", Line},
		{"Bar", Bar},
	}

	for _, s := range styles {
		t.Run(s.name, func(t *testing.T) {
			cfg := Config{
				Data:   []float64{},
				Style:  s.style,
				Bounds: bounds,
			}
			result := Render(cfg)
			if result == nil {
				t.Fatal("expected non-nil result for empty data")
			}

			img := result.Image
			w := bounds.Dx()
			h := bounds.Dy()

			for y := 0; y < h; y++ {
				for x := 0; x < w; x++ {
					r, g, b, a := img.At(x, y).RGBA()
					er, eg, eb, ea := bg.RGBA()
					if r != er || g != eg || b != eb || a != ea {
						t.Fatalf("pixel (%d,%d): got (%d,%d,%d,%d), want background (%d,%d,%d,%d)",
							x, y, r, g, b, a, er, eg, eb, ea)
					}
				}
			}
		})
	}
}

// TestSingleDataPointLine verifies that a single data point in Line style
// renders only column x=0 with foreground pixels.
func TestSingleDataPointLine(t *testing.T) {
	bounds := image.Rect(0, 0, 10, 10)
	cfg := Config{
		Data:   []float64{0.5},
		Style:  Line,
		Bounds: bounds,
	}
	result := Render(cfg)
	if result == nil {
		t.Fatal("expected non-nil result")
	}

	img := result.Image
	fg := color.RGBA{R: 255, G: 255, B: 255, A: 255} // default white
	bg := color.RGBA{R: 0, G: 0, B: 0, A: 255}       // default black

	// Column 0 should have floor(10 * 0.5) = 5 foreground pixels from the bottom
	fgCount := 0
	for y := 0; y < 10; y++ {
		r, g, b, a := img.At(0, y).RGBA()
		er, eg, eb, ea := fg.RGBA()
		if r == er && g == eg && b == eb && a == ea {
			fgCount++
		}
	}
	if fgCount != 5 {
		t.Errorf("column 0: expected 5 foreground pixels, got %d", fgCount)
	}

	// All other columns should be entirely background
	for x := 1; x < 10; x++ {
		for y := 0; y < 10; y++ {
			r, g, b, a := img.At(x, y).RGBA()
			er, eg, eb, ea := bg.RGBA()
			if r != er || g != eg || b != eb || a != ea {
				t.Fatalf("pixel (%d,%d): expected background, got foreground", x, y)
			}
		}
	}
}

// TestSingleDataPointBar verifies that a single data point in Bar style
// fills all columns with the bar height.
func TestSingleDataPointBar(t *testing.T) {
	bounds := image.Rect(0, 0, 10, 10)
	cfg := Config{
		Data:   []float64{0.5},
		Style:  Bar,
		Bounds: bounds,
	}
	result := Render(cfg)
	if result == nil {
		t.Fatal("expected non-nil result")
	}

	img := result.Image
	fg := color.RGBA{R: 255, G: 255, B: 255, A: 255}

	// barWidth = 10/1 = 10, so all columns 0-9 should have 5 foreground pixels from bottom
	for x := 0; x < 10; x++ {
		fgCount := 0
		for y := 0; y < 10; y++ {
			r, g, b, a := img.At(x, y).RGBA()
			er, eg, eb, ea := fg.RGBA()
			if r == er && g == eg && b == eb && a == ea {
				fgCount++
			}
		}
		if fgCount != 5 {
			t.Errorf("column %d: expected 5 foreground pixels, got %d", x, fgCount)
		}
	}
}

// TestDefaultColors verifies that zero-value colors resolve to white foreground
// and black background.

func TestDefaultColors(t *testing.T) {
	bounds := image.Rect(0, 0, 5, 5)
	cfg := Config{
		Data:       []float64{1.0},
		Style:      Bar,
		Bounds:     bounds,
		Foreground: color.RGBA{}, // zero value → white
		Background: color.RGBA{}, // zero value → black
	}
	result := Render(cfg)
	if result == nil {
		t.Fatal("expected non-nil result")
	}

	img := result.Image
	white := color.RGBA{R: 255, G: 255, B: 255, A: 255}

	// With data=[1.0] and Bar style, barWidth=5/1=5, all pixels should be foreground (white)
	for y := 0; y < 5; y++ {
		for x := 0; x < 5; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			er, eg, eb, ea := white.RGBA()
			if r != er || g != eg || b != eb || a != ea {
				t.Fatalf("pixel (%d,%d): expected white (255,255,255,255), got different color", x, y)
			}
		}
	}

	// Also verify background color with data=[0.0]
	cfg2 := Config{
		Data:       []float64{0.0},
		Style:      Bar,
		Bounds:     bounds,
		Foreground: color.RGBA{},
		Background: color.RGBA{},
	}
	result2 := Render(cfg2)
	if result2 == nil {
		t.Fatal("expected non-nil result for data=[0.0]")
	}

	img2 := result2.Image
	black := color.RGBA{R: 0, G: 0, B: 0, A: 255}

	// With data=[0.0], height=floor(5*0.0)=0, all pixels should be background (black)
	for y := 0; y < 5; y++ {
		for x := 0; x < 5; x++ {
			r, g, b, a := img2.At(x, y).RGBA()
			er, eg, eb, ea := black.RGBA()
			if r != er || g != eg || b != eb || a != ea {
				t.Fatalf("pixel (%d,%d): expected black (0,0,0,255), got different color", x, y)
			}
		}
	}
}

// TestFewerDataPointsThanWidth verifies that bars are spaced correctly when
// there are fewer data points than pixel columns.
func TestFewerDataPointsThanWidth(t *testing.T) {
	bounds := image.Rect(0, 0, 10, 10)
	cfg := Config{
		Data:   []float64{0.5, 1.0},
		Style:  Bar,
		Bounds: bounds,
	}
	result := Render(cfg)
	if result == nil {
		t.Fatal("expected non-nil result")
	}

	img := result.Image
	fg := color.RGBA{R: 255, G: 255, B: 255, A: 255}

	// barWidth = 10/2 = 5
	// First bar (cols 0-4): height = floor(10 * 0.5) = 5
	for x := 0; x < 5; x++ {
		fgCount := 0
		for y := 0; y < 10; y++ {
			r, g, b, a := img.At(x, y).RGBA()
			er, eg, eb, ea := fg.RGBA()
			if r == er && g == eg && b == eb && a == ea {
				fgCount++
			}
		}
		if fgCount != 5 {
			t.Errorf("first bar column %d: expected 5 foreground pixels, got %d", x, fgCount)
		}
	}

	// Second bar (cols 5-9): height = floor(10 * 1.0) = 10
	for x := 5; x < 10; x++ {
		fgCount := 0
		for y := 0; y < 10; y++ {
			r, g, b, a := img.At(x, y).RGBA()
			er, eg, eb, ea := fg.RGBA()
			if r == er && g == eg && b == eb && a == ea {
				fgCount++
			}
		}
		if fgCount != 10 {
			t.Errorf("second bar column %d: expected 10 foreground pixels, got %d", x, fgCount)
		}
	}
}
