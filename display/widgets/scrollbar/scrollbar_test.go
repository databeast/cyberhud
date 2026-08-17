package scrollbar

import (
	"image"
	"image/color"
	"testing"

	"pgregory.net/rapid"
)

// --- From: scrollbar_prop_test.go ---

// TestPropertyScrollbarOutputMetadataCorrectness verifies that for any valid scrollbar Config
// (bounds width ≥ 1, height ≥ 1, total ≥ 1), the Scrollbar_Widget returns a non-nil Result
// where Image dimensions equal Bounds.Dx() × Bounds.Dy(), Position equals Bounds.Min,
// and Label equals "scrollbar".
//

func TestPropertyScrollbarOutputMetadataCorrectness(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		minX := rapid.IntRange(0, 100).Draw(t, "minX")
		minY := rapid.IntRange(0, 100).Draw(t, "minY")
		width := rapid.IntRange(1, 100).Draw(t, "width")
		height := rapid.IntRange(1, 100).Draw(t, "height")

		bounds := image.Rect(minX, minY, minX+width, minY+height)

		totalItems := rapid.IntRange(1, 500).Draw(t, "totalItems")
		visibleItems := rapid.IntRange(1, totalItems).Draw(t, "visibleItems")
		scrollOffset := rapid.IntRange(0, totalItems-1).Draw(t, "scrollOffset")

		fg := color.RGBA{R: 255, G: 0, B: 0, A: 255}
		bg := color.RGBA{R: 0, G: 0, B: 255, A: 255}

		result := Render(Config{
			TotalItems:   totalItems,
			VisibleItems: visibleItems,
			ScrollOffset: scrollOffset,
			Bounds:       bounds,
			Foreground:   fg,
			Background:   bg,
		})

		if result == nil {
			t.Fatal("expected non-nil result for valid scrollbar config")
		}

		// Verify image dimensions match bounds
		gotWidth := result.Image.Bounds().Dx()
		gotHeight := result.Image.Bounds().Dy()
		if gotWidth != width {
			t.Fatalf("image width mismatch: got %d, want %d", gotWidth, width)
		}
		if gotHeight != height {
			t.Fatalf("image height mismatch: got %d, want %d", gotHeight, height)
		}

		// Verify Position equals Bounds.Min
		if result.Position.X != minX || result.Position.Y != minY {
			t.Fatalf("Position mismatch: got (%d, %d), want (%d, %d)",
				result.Position.X, result.Position.Y, minX, minY)
		}

		// Verify Label
		if result.Label != "scrollbar" {
			t.Fatalf("Label mismatch: got %q, want %q", result.Label, "scrollbar")
		}
	})
}

// TestPropertyScrollbarThumbSizingAndPositioning verifies that for any valid scrollbar Config
// where total > visible, the number of contiguous foreground-colored pixel rows equals
// max(1, floor(Bounds.Dy() × visible / total)), and the first foreground row is at
// min(floor(Bounds.Dy() × clampedOffset / total), Bounds.Dy() - thumbHeight).
//

func TestPropertyScrollbarThumbSizingAndPositioning(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		width := rapid.IntRange(1, 50).Draw(t, "width")
		height := rapid.IntRange(1, 200).Draw(t, "height")

		bounds := image.Rect(0, 0, width, height)

		totalItems := rapid.IntRange(2, 200).Draw(t, "totalItems")
		visibleItems := rapid.IntRange(1, totalItems-1).Draw(t, "visibleItems")
		scrollOffset := rapid.IntRange(0, totalItems-1).Draw(t, "scrollOffset")

		fg := color.RGBA{R: 255, G: 0, B: 0, A: 255}
		bg := color.RGBA{R: 0, G: 0, B: 255, A: 255}

		result := Render(Config{
			TotalItems:   totalItems,
			VisibleItems: visibleItems,
			ScrollOffset: scrollOffset,
			Bounds:       bounds,
			Foreground:   fg,
			Background:   bg,
		})

		if result == nil {
			t.Fatal("expected non-nil result for valid scrollbar config")
		}

		// Expected thumb height = max(1, floor(h × visible / total))
		expectedThumbHeight := (height * visibleItems) / totalItems
		if expectedThumbHeight < 1 {
			expectedThumbHeight = 1
		}

		// Expected thumb top = min(floor(h × offset / total), h - thumbHeight)
		expectedThumbTop := (height * scrollOffset) / totalItems
		if expectedThumbTop > height-expectedThumbHeight {
			expectedThumbTop = height - expectedThumbHeight
		}

		// Count foreground rows and find the first foreground row
		// A row is "foreground" if all pixels in that row match the foreground color
		firstFgRow := -1
		fgRowCount := 0

		for y := 0; y < height; y++ {
			isFgRow := true
			for x := 0; x < width; x++ {
				r, g, b, a := result.Image.At(x, y).RGBA()
				// foreground is red: 0xFFFF, 0, 0, 0xFFFF
				if r != 0xFFFF || g != 0 || b != 0 || a != 0xFFFF {
					isFgRow = false
					break
				}
			}
			if isFgRow {
				if firstFgRow == -1 {
					firstFgRow = y
				}
				fgRowCount++
			}
		}

		// Verify thumb height
		if fgRowCount != expectedThumbHeight {
			t.Fatalf("thumb height mismatch: got %d, want %d (height=%d, visible=%d, total=%d, offset=%d)",
				fgRowCount, expectedThumbHeight, height, visibleItems, totalItems, scrollOffset)
		}

		// Verify thumb top position
		if firstFgRow != expectedThumbTop {
			t.Fatalf("thumb top mismatch: got %d, want %d (height=%d, visible=%d, total=%d, offset=%d)",
				firstFgRow, expectedThumbTop, height, visibleItems, totalItems, scrollOffset)
		}

		// Verify foreground rows are contiguous (single thumb block)
		if fgRowCount > 0 {
			for y := firstFgRow; y < firstFgRow+fgRowCount; y++ {
				for x := 0; x < width; x++ {
					r, g, b, a := result.Image.At(x, y).RGBA()
					if r != 0xFFFF || g != 0 || b != 0 || a != 0xFFFF {
						t.Fatalf("non-foreground pixel at (%d,%d) within expected thumb region [%d, %d)",
							x, y, firstFgRow, firstFgRow+fgRowCount)
					}
				}
			}
		}
	})
}

// TestPropertyScrollbarFullFillWhenNoScrollingNeeded verifies that for any valid scrollbar Config
// where total ≤ visible, every pixel in the output image is the Foreground_Color.
//

func TestPropertyScrollbarFullFillWhenNoScrollingNeeded(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		width := rapid.IntRange(1, 100).Draw(t, "width")
		height := rapid.IntRange(1, 100).Draw(t, "height")

		bounds := image.Rect(0, 0, width, height)

		totalItems := rapid.IntRange(1, 100).Draw(t, "totalItems")
		// visible >= total (no scrolling needed)
		visibleItems := rapid.IntRange(totalItems, totalItems+100).Draw(t, "visibleItems")
		scrollOffset := rapid.IntRange(0, totalItems-1).Draw(t, "scrollOffset")

		fg := color.RGBA{R: 255, G: 0, B: 0, A: 255}
		bg := color.RGBA{R: 0, G: 0, B: 255, A: 255}

		result := Render(Config{
			TotalItems:   totalItems,
			VisibleItems: visibleItems,
			ScrollOffset: scrollOffset,
			Bounds:       bounds,
			Foreground:   fg,
			Background:   bg,
		})

		if result == nil {
			t.Fatal("expected non-nil result for valid scrollbar config")
		}

		// Every pixel should be the foreground color (red)
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				r, g, b, a := result.Image.At(x, y).RGBA()
				if r != 0xFFFF || g != 0 || b != 0 || a != 0xFFFF {
					t.Fatalf("pixel (%d,%d): expected foreground (0xFFFF,0,0,0xFFFF), got (%d,%d,%d,%d) [total=%d, visible=%d]",
						x, y, r, g, b, a, totalItems, visibleItems)
				}
			}
		}
	})
}

// TestPropertyScrollbarOffsetClampingIdempotence verifies that for any scrollbar Config with
// scroll_offset < 0 or scroll_offset >= total_count, the rendered output is pixel-identical
// to rendering with the offset clamped to [0, total_count - 1].
//

func TestPropertyScrollbarOffsetClampingIdempotence(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		width := rapid.IntRange(1, 50).Draw(t, "width")
		height := rapid.IntRange(1, 100).Draw(t, "height")

		bounds := image.Rect(0, 0, width, height)

		totalItems := rapid.IntRange(1, 200).Draw(t, "totalItems")
		visibleItems := rapid.IntRange(1, totalItems).Draw(t, "visibleItems")

		// Generate an out-of-range offset: either negative or >= totalItems
		negative := rapid.SampledFrom([]bool{true, false}).Draw(t, "negative")
		var outOfRangeOffset int
		if negative {
			outOfRangeOffset = rapid.IntRange(-100, -1).Draw(t, "offset")
		} else {
			outOfRangeOffset = rapid.IntRange(totalItems, totalItems+100).Draw(t, "offset")
		}

		// Compute clamped offset
		clampedOffset := outOfRangeOffset
		if clampedOffset < 0 {
			clampedOffset = 0
		}
		if clampedOffset >= totalItems {
			clampedOffset = totalItems - 1
		}

		fg := color.RGBA{R: 255, G: 0, B: 0, A: 255}
		bg := color.RGBA{R: 0, G: 0, B: 255, A: 255}

		// Render with out-of-range offset
		result1 := Render(Config{
			TotalItems:   totalItems,
			VisibleItems: visibleItems,
			ScrollOffset: outOfRangeOffset,
			Bounds:       bounds,
			Foreground:   fg,
			Background:   bg,
		})

		// Render with clamped offset
		result2 := Render(Config{
			TotalItems:   totalItems,
			VisibleItems: visibleItems,
			ScrollOffset: clampedOffset,
			Bounds:       bounds,
			Foreground:   fg,
			Background:   bg,
		})

		if result1 == nil {
			t.Fatal("expected non-nil result for out-of-range offset render")
		}
		if result2 == nil {
			t.Fatal("expected non-nil result for clamped offset render")
		}

		// Compare pixel-by-pixel
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				r1, g1, b1, a1 := result1.Image.At(x, y).RGBA()
				r2, g2, b2, a2 := result2.Image.At(x, y).RGBA()
				if r1 != r2 || g1 != g2 || b1 != b2 || a1 != a2 {
					t.Fatalf("pixel mismatch at (%d,%d): out-of-range=(%d,%d,%d,%d), clamped=(%d,%d,%d,%d) [total=%d, visible=%d, rawOffset=%d, clampedOffset=%d]",
						x, y, r1, g1, b1, a1, r2, g2, b2, a2, totalItems, visibleItems, outOfRangeOffset, clampedOffset)
				}
			}
		}
	})
}

// --- From: scrollbar_test.go ---
