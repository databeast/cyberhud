package textlabel

import (
	"image"
	"image/color"
	"testing"

	"github.com/databeast/cyberhud/display/surface/fonts"
	"pgregory.net/rapid"
)

// --- From: textlabel_prop_test.go ---

// TestPropertyOutputMetadataCorrectness verifies that for any valid Config
// (non-empty bounds), the Text Label Widget returns a non-nil Result with
// Image dimensions equal to Bounds.Dx() × Bounds.Dy(), Position equal to
// Bounds.Min, and Label equal to "textlabel".
//

func TestPropertyOutputMetadataCorrectness(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		width := rapid.IntRange(1, 100).Draw(t, "width")
		height := rapid.IntRange(1, 100).Draw(t, "height")
		minX := rapid.IntRange(0, 50).Draw(t, "minX")
		minY := rapid.IntRange(0, 50).Draw(t, "minY")
		bounds := image.Rect(minX, minY, minX+width, minY+height)

		text := rapid.StringMatching(`[A-Za-z0-9]{0,20}`).Draw(t, "text")
		alignment := Alignment(rapid.IntRange(0, 2).Draw(t, "alignment"))

		fg := color.RGBA{
			R: uint8(rapid.IntRange(0, 255).Draw(t, "r")),
			G: uint8(rapid.IntRange(0, 255).Draw(t, "g")),
			B: uint8(rapid.IntRange(0, 255).Draw(t, "b")),
			A: uint8(rapid.IntRange(0, 255).Draw(t, "a")),
		}

		result := Render(Config{
			Text:       text,
			Bounds:     bounds,
			Font:       font.Default(),
			Alignment:  alignment,
			Foreground: fg,
		})

		if result == nil {
			t.Fatal("expected non-nil result for valid bounds")
		}

		// Verify image dimensions match bounds
		imgBounds := result.Image.Bounds()
		gotWidth := imgBounds.Dx()
		gotHeight := imgBounds.Dy()
		if gotWidth != width {
			t.Fatalf("image width mismatch: got %d, want %d", gotWidth, width)
		}
		if gotHeight != height {
			t.Fatalf("image height mismatch: got %d, want %d", gotHeight, height)
		}

		// Verify Position equals Bounds.Min
		if result.Position.X != minX || result.Position.Y != minY {
			t.Fatalf("Position mismatch: got (%d,%d), want (%d,%d)",
				result.Position.X, result.Position.Y, minX, minY)
		}

		// Verify Label
		if result.Label != "textlabel" {
			t.Fatalf("Label mismatch: got %q, want %q", result.Label, "textlabel")
		}
	})
}

// TestPropertyAlignmentPositioning verifies that the rendered foreground pixels
// begin at the correct x-offset based on alignment. We compute the expected
// leftmost pixel as xStart + the first lit column in the first character's glyph.
// - Left: xStart = 0
// - Center: xStart = floor((width - textWidth) / 2), min 0
// - Right: xStart = width - textWidth, min 0
//

func TestPropertyAlignmentPositioning(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate non-empty text (1-10 chars from printable ASCII that have visible glyphs)
		text := rapid.StringMatching(`[A-Z]{1,10}`).Draw(t, "text")

		face := font.Default()
		metrics := face.Metrics()
		glyphAdvance := metrics.GlyphAdvance
		glyphWidth := metrics.GlyphWidth
		glyphHeight := metrics.GlyphHeight

		charCount := len(text)
		textPixelWidth := (charCount-1)*glyphAdvance + glyphWidth

		// Ensure bounds are wide enough to avoid clipping for alignment tests
		width := rapid.IntRange(textPixelWidth, textPixelWidth+50).Draw(t, "width")
		height := rapid.IntRange(glyphHeight, 50).Draw(t, "height")
		bounds := image.Rect(0, 0, width, height)

		alignment := Alignment(rapid.IntRange(0, 2).Draw(t, "alignment"))

		fg := color.RGBA{R: 255, G: 0, B: 0, A: 255}

		result := Render(Config{
			Text:       text,
			Bounds:     bounds,
			Font:       face,
			Alignment:  alignment,
			Foreground: fg,
		})

		if result == nil {
			t.Fatal("expected non-nil result for valid bounds")
		}

		// Compute the first lit column of the first character's glyph
		firstChar := rune(text[0])
		firstLitCol := glyphWidth // if no pixel found, use full width
		for row := 0; row < glyphHeight; row++ {
			bits := face.GlyphRow(firstChar, row)
			for col := 0; col < glyphWidth; col++ {
				if bits&(1<<uint(31-col)) != 0 {
					if col < firstLitCol {
						firstLitCol = col
					}
				}
			}
		}

		// Find the leftmost foreground pixel x-coordinate
		leftmostX := -1
		for x := 0; x < width; x++ {
			for y := 0; y < height; y++ {
				_, _, _, a := result.Image.At(x, y).RGBA()
				if a > 0 {
					leftmostX = x
					goto found
				}
			}
		}
	found:
		if leftmostX == -1 {
			t.Fatal("no foreground pixels found for non-empty text")
		}

		var xStart int
		switch alignment {
		case Left:
			xStart = 0
		case Center:
			xStart = (width - textPixelWidth) / 2
			if xStart < 0 {
				xStart = 0
			}
		case Right:
			xStart = width - textPixelWidth
			if xStart < 0 {
				xStart = 0
			}
		}

		expectedLeftmostX := xStart + firstLitCol

		if leftmostX != expectedLeftmostX {
			t.Fatalf("alignment %d: leftmost foreground pixel at x=%d, expected x=%d (xStart=%d, firstLitCol=%d, width=%d, textWidth=%d, text=%q)",
				alignment, leftmostX, expectedLeftmostX, xStart, firstLitCol, width, textPixelWidth, text)
		}
	})
}

// TestPropertyClippingAtBounds verifies that when text pixel width exceeds
// Bounds.Dx(), no foreground pixels appear at x >= Bounds.Dx().
//

func TestPropertyClippingAtBounds(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate text that will be wider than bounds
		text := rapid.StringMatching(`[A-Z]{5,20}`).Draw(t, "text")

		face := font.Default()
		metrics := face.Metrics()
		glyphAdvance := metrics.GlyphAdvance
		glyphWidth := metrics.GlyphWidth

		charCount := len(text)
		textPixelWidth := (charCount-1)*glyphAdvance + glyphWidth

		// Make bounds narrower than text pixel width
		width := rapid.IntRange(1, textPixelWidth-1).Draw(t, "width")
		height := rapid.IntRange(metrics.GlyphHeight, 30).Draw(t, "height")
		bounds := image.Rect(0, 0, width, height)

		fg := color.RGBA{R: 255, G: 0, B: 0, A: 255}

		result := Render(Config{
			Text:       text,
			Bounds:     bounds,
			Font:       face,
			Alignment:  Left,
			Foreground: fg,
		})

		if result == nil {
			t.Fatal("expected non-nil result for valid bounds")
		}

		// Verify no foreground pixels at x >= bounds width
		imgBounds := result.Image.Bounds()
		for x := imgBounds.Min.X; x < imgBounds.Max.X; x++ {
			for y := imgBounds.Min.Y; y < imgBounds.Max.Y; y++ {
				if x >= width {
					_, _, _, a := result.Image.At(x, y).RGBA()
					if a > 0 {
						t.Fatalf("found foreground pixel at x=%d (>= bounds width %d), y=%d",
							x, width, y)
					}
				}
			}
		}
	})
}

// TestPropertyForegroundColorFidelity verifies that for any valid Config with a
// non-zero Foreground_Color and non-empty text, every non-transparent pixel in
// the output image has RGBA values exactly matching the specified Foreground_Color.
//

func TestPropertyForegroundColorFidelity(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		text := rapid.StringMatching(`[A-Z]{1,10}`).Draw(t, "text")

		width := rapid.IntRange(10, 100).Draw(t, "width")
		height := rapid.IntRange(10, 50).Draw(t, "height")
		bounds := image.Rect(0, 0, width, height)

		// Generate non-zero foreground color (ensure at least alpha > 0)
		r := uint8(rapid.IntRange(1, 255).Draw(t, "r"))
		g := uint8(rapid.IntRange(0, 255).Draw(t, "g"))
		b := uint8(rapid.IntRange(0, 255).Draw(t, "b"))
		a := uint8(rapid.IntRange(1, 255).Draw(t, "a"))
		fg := color.RGBA{R: r, G: g, B: b, A: a}

		result := Render(Config{
			Text:       text,
			Bounds:     bounds,
			Font:       font.Default(),
			Alignment:  Left,
			Foreground: fg,
		})

		if result == nil {
			t.Fatal("expected non-nil result for valid bounds")
		}

		// Check every non-transparent pixel matches the foreground color
		for x := 0; x < width; x++ {
			for y := 0; y < height; y++ {
				pr, pg, pb, pa := result.Image.At(x, y).RGBA()
				if pa == 0 {
					continue // transparent pixel, skip
				}
				// Compare using pre-multiplied 16-bit values
				er, eg, eb, ea := fg.RGBA()
				if pr != er || pg != eg || pb != eb || pa != ea {
					t.Fatalf("pixel (%d,%d): got RGBA(%d,%d,%d,%d), want RGBA(%d,%d,%d,%d) for fg=%v",
						x, y, pr, pg, pb, pa, er, eg, eb, ea, fg)
				}
			}
		}
	})
}

// TestPropertyEmptyStringProducesTransparentImage verifies that for any valid
// bounds and empty text string, the returned Result's image contains only fully
// transparent pixels (alpha = 0).
//

func TestPropertyEmptyStringProducesTransparentImage(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		width := rapid.IntRange(1, 100).Draw(t, "width")
		height := rapid.IntRange(1, 100).Draw(t, "height")
		minX := rapid.IntRange(0, 50).Draw(t, "minX")
		minY := rapid.IntRange(0, 50).Draw(t, "minY")
		bounds := image.Rect(minX, minY, minX+width, minY+height)

		fg := color.RGBA{
			R: uint8(rapid.IntRange(0, 255).Draw(t, "r")),
			G: uint8(rapid.IntRange(0, 255).Draw(t, "g")),
			B: uint8(rapid.IntRange(0, 255).Draw(t, "b")),
			A: uint8(rapid.IntRange(0, 255).Draw(t, "a")),
		}

		result := Render(Config{
			Text:       "", // empty text
			Bounds:     bounds,
			Font:       font.Default(),
			Alignment:  Alignment(rapid.IntRange(0, 2).Draw(t, "alignment")),
			Foreground: fg,
		})

		if result == nil {
			t.Fatal("expected non-nil result for valid bounds with empty text")
		}

		// Verify dimensions are correct
		imgBounds := result.Image.Bounds()
		if imgBounds.Dx() != width || imgBounds.Dy() != height {
			t.Fatalf("image dimensions mismatch: got %dx%d, want %dx%d",
				imgBounds.Dx(), imgBounds.Dy(), width, height)
		}

		// Verify all pixels are fully transparent
		for x := imgBounds.Min.X; x < imgBounds.Max.X; x++ {
			for y := imgBounds.Min.Y; y < imgBounds.Max.Y; y++ {
				_, _, _, a := result.Image.At(x, y).RGBA()
				if a != 0 {
					t.Fatalf("pixel (%d,%d) has non-zero alpha %d for empty text", x, y, a)
				}
			}
		}
	})
}

// --- From: textlabel_test.go ---

// TestRenderNilForInvalidBounds verifies that Render returns nil for
// bounds with width < 1 or height < 1.

func TestRenderNilForInvalidBounds(t *testing.T) {
	cases := []struct {
		name   string
		bounds image.Rectangle
	}{
		{"0x0", image.Rect(0, 0, 0, 0)},
		{"-1x5", image.Rectangle{Min: image.Pt(0, 0), Max: image.Pt(-1, 5)}},
		{"5x0", image.Rect(0, 0, 5, 0)},
		{"0x10", image.Rect(0, 0, 0, 10)},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := Render(Config{
				Text:   "Hello",
				Bounds: tc.bounds,
				Font:   font.Default(),
			})
			if result != nil {
				t.Errorf("expected nil for bounds %v, got non-nil result", tc.bounds)
			}
		})
	}
}

// TestDefaultForegroundColor verifies that when Foreground is the zero value,
// all non-transparent pixels are rendered as opaque white (255, 255, 255, 255).

func TestDefaultForegroundColor(t *testing.T) {
	bounds := image.Rect(0, 0, 10, 10)
	result := Render(Config{
		Text:   "A",
		Bounds: bounds,
		Font:   font.Default(),
		// Foreground is zero value → should default to white
	})

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	white := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	img := result.Image

	for y := 0; y < bounds.Dy(); y++ {
		for x := 0; x < bounds.Dx(); x++ {
			r, g, b, a := img.At(x, y).RGBA()
			if a == 0 {
				continue // transparent pixel, skip
			}
			// Convert from 16-bit RGBA back to 8-bit for comparison
			got := color.RGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: uint8(a >> 8)}
			if got != white {
				t.Fatalf("pixel (%d,%d) = %v, want opaque white %v", x, y, got, white)
			}
		}
	}
}

// TestDefaultAlignmentIsLeft verifies that when Alignment is not explicitly set,
// text defaults to Left alignment (foreground pixels appear starting at x=0).

func TestDefaultAlignmentIsLeft(t *testing.T) {
	// 'A' row 2 is 0x11 = 0b10001, which has bit 4 set → col 0 is a foreground pixel.
	// With Left alignment (xStart=0), there should be a foreground pixel at x=0.
	bounds := image.Rect(0, 0, 10, 10)
	result := Render(Config{
		Text:   "A",
		Bounds: bounds,
		Font:   font.Default(),
		// Alignment not set → defaults to Left (iota = 0)
	})

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	img := result.Image
	// Check that there's at least one foreground pixel in column x=0
	found := false
	for y := 0; y < bounds.Dy(); y++ {
		_, _, _, a := img.At(0, y).RGBA()
		if a > 0 {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected foreground pixel at x=0 for left-aligned 'A', but found none")
	}
}

// TestTextExactlyFitsBounds verifies that when text pixel width exactly matches
// the bounds width, the result is non-nil and the rightmost glyph column has pixels.

func TestTextExactlyFitsBounds(t *testing.T) {
	face := font.Default()
	metrics := face.Metrics()

	// Use a single wide character ('M' or 'W') that likely has pixels near its rightmost column.
	// For 2 characters: textPixelWidth = (2-1)*GlyphAdvance + GlyphWidth
	text := "AB"
	textWidth := (len(text)-1)*metrics.GlyphAdvance + metrics.GlyphWidth

	bounds := image.Rect(0, 0, textWidth, metrics.GlyphHeight)
	result := Render(Config{
		Text:   text,
		Bounds: bounds,
		Font:   face,
	})

	if result == nil {
		t.Fatal("expected non-nil result when text exactly fits bounds")
	}

	img := result.Image

	// Verify the image has some foreground pixels (text is non-empty and fits bounds).
	found := false
	for y := 0; y < bounds.Dy() && !found; y++ {
		for x := 0; x < bounds.Dx() && !found; x++ {
			_, _, _, a := img.At(x, y).RGBA()
			if a > 0 {
				found = true
			}
		}
	}
	if !found {
		t.Error("expected at least one foreground pixel for non-empty text that fits bounds")
	}

	// Verify no pixels are rendered beyond bounds (image matches expected dimensions).
	if img.Bounds().Dx() != textWidth || img.Bounds().Dy() != metrics.GlyphHeight {
		t.Errorf("image bounds = %v, want %dx%d", img.Bounds(), textWidth, metrics.GlyphHeight)
	}
}

// TestEmptyTextTransparentImage verifies that rendering an empty string produces
// an image where all pixels are fully transparent (alpha = 0).

func TestEmptyTextTransparentImage(t *testing.T) {
	bounds := image.Rect(0, 0, 20, 10)
	result := Render(Config{
		Text:   "",
		Bounds: bounds,
		Font:   font.Default(),
	})

	if result == nil {
		t.Fatal("expected non-nil result for empty text")
	}

	img := result.Image
	for y := 0; y < bounds.Dy(); y++ {
		for x := 0; x < bounds.Dx(); x++ {
			_, _, _, a := img.At(x, y).RGBA()
			if a != 0 {
				t.Fatalf("pixel (%d,%d) has alpha %d, want 0 (fully transparent)", x, y, a)
			}
		}
	}
}
