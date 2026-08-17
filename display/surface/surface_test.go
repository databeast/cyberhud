package surface

import (
	"image"
	"image/color"
	"image/draw"
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"

	font "github.com/databeast/cyberhud/display/surface/fonts"
)

// ---------------------------------------------------------------------------
// Tests merged from scale_test.go
// ---------------------------------------------------------------------------

// For any source image.Image with positive dimensions and for any destination
// image.Rectangle with positive width and height that intersects the surface bounds,
// calling DrawImageScaled SHALL produce a framebuffer where each destination pixel
// (dx, dy) contains the color of the source pixel at
// (srcMinX + dx * srcW / dstW, srcMinY + dy * srcH / dstH) (integer division, nearest-neighbor mapping).

// scaleInput holds a randomly generated test input for the nearest-neighbor scaling property test.
type scaleInput struct {
	SrcW, SrcH int   // Source image dimensions (1..32)
	DstW, DstH int   // Destination dimensions (1..64)
	Seed       int64 // RNG seed for reproducible image generation
}

// Generate implements quick.Generator for property-based testing.
func (scaleInput) Generate(r *rand.Rand, size int) reflect.Value {
	input := scaleInput{
		SrcW: 1 + r.Intn(32),
		SrcH: 1 + r.Intn(32),
		DstW: 1 + r.Intn(64),
		DstH: 1 + r.Intn(64),
		Seed: r.Int63(),
	}
	return reflect.ValueOf(input)
}

func TestPropertyDrawImageScaledNearestNeighborCorrectness(t *testing.T) {
	config := &quick.Config{MaxCount: 200}

	prop := func(input scaleInput) bool {
		// Generate a random RGBA source image.
		rng := rand.New(rand.NewSource(input.Seed))
		src := image.NewRGBA(image.Rect(0, 0, input.SrcW, input.SrcH))
		for y := 0; y < input.SrcH; y++ {
			for x := 0; x < input.SrcW; x++ {
				src.SetRGBA(x, y, color.RGBA{
					R: uint8(rng.Intn(256)),
					G: uint8(rng.Intn(256)),
					B: uint8(rng.Intn(256)),
					A: uint8(rng.Intn(256)),
				})
			}
		}

		// Call scaleNearestNeighbor directly.
		result := scaleNearestNeighbor(src, input.DstW, input.DstH)

		// Verify the result dimensions.
		if result.Bounds().Dx() != input.DstW || result.Bounds().Dy() != input.DstH {
			t.Errorf("result dimensions mismatch: got %dx%d, want %dx%d",
				result.Bounds().Dx(), result.Bounds().Dy(), input.DstW, input.DstH)
			return false
		}

		// Verify each pixel matches the expected nearest-neighbor mapping.
		srcBounds := src.Bounds()
		srcMinX := srcBounds.Min.X
		srcMinY := srcBounds.Min.Y
		srcW := srcBounds.Dx()
		srcH := srcBounds.Dy()

		for dy := 0; dy < input.DstH; dy++ {
			for dx := 0; dx < input.DstW; dx++ {
				// Compute expected source coordinate using integer division.
				expectedSrcX := srcMinX + dx*srcW/input.DstW
				expectedSrcY := srcMinY + dy*srcH/input.DstH

				// Get the expected color from the source at that coordinate.
				expectedColor := src.RGBAAt(expectedSrcX, expectedSrcY)

				// Get the actual color from the scaled result.
				gotColor := result.RGBAAt(dx, dy)

				if gotColor != expectedColor {
					t.Errorf("pixel (%d,%d): got %v, want %v (mapped from src(%d,%d), srcSize=%dx%d, dstSize=%dx%d)",
						dx, dy, gotColor, expectedColor,
						expectedSrcX, expectedSrcY,
						input.SrcW, input.SrcH, input.DstW, input.DstH)
					return false
				}
			}
		}

		return true
	}

	if err := quick.Check(prop, config); err != nil {
		t.Errorf("Property 3 failed: %v", err)
	}
}

// ---------------------------------------------------------------------------
// Tests merged from surface_test.go
// ---------------------------------------------------------------------------

// TestSetFontWithValidFace verifies that SetFont with a non-nil face updates
// the surface's font so that FontMetrics returns the new face's metrics.

func TestSetFontWithValidFace(t *testing.T) {
	s := New(image.Rect(0, 0, 128, 64))

	// Get a known face different from the default (spleen-5x8).
	face, ok := font.Get(font.Spleen8x16ID)
	if !ok {
		t.Fatalf("font %q not registered", font.Spleen8x16ID)
	}

	s.SetFontID(font.Spleen8x16ID)

	m := s.FontMetrics()
	if m.GlyphWidth <= 0 || m.GlyphHeight <= 0 || m.GlyphAdvance <= 0 || m.RowHeight <= 0 {
		t.Fatalf("expected non-zero metrics after SetFontID, got %+v", m)
	}

	// Verify we actually got the metrics from the face we set.
	expected := face.Metrics()
	if m != expected {
		t.Errorf("FontMetrics = %+v, want %+v", m, expected)
	}
}

// TestSetFontIDValid verifies that SetFontID returns true and applies the font
// when given a known registered font ID.

func TestSetFontIDValid(t *testing.T) {
	s := New(image.Rect(0, 0, 128, 64))

	ok := s.SetFontID(font.Spleen5x8ID)
	if !ok {
		t.Fatalf("SetFontID(%q) returned false, want true", font.Spleen5x8ID)
	}

	m := s.FontMetrics()
	expected, _ := font.Get(font.Spleen5x8ID)
	if m != expected.Metrics() {
		t.Errorf("FontMetrics after SetFontID = %+v, want %+v", m, expected.Metrics())
	}
}

// TestSetFontIDInvalid verifies that SetFontID returns false for an unregistered
// font ID and does not panic.

func TestSetFontIDInvalid(t *testing.T) {
	s := New(image.Rect(0, 0, 128, 64))

	ok := s.SetFontID("nonexistent-font-id-xyz")
	if ok {
		t.Fatalf("SetFontID with unknown ID returned true, want false")
	}

	// Verify it doesn't leave the surface in a broken state — FontMetrics
	// should still return valid metrics (the default font is applied).
	m := s.FontMetrics()
	if m.GlyphWidth <= 0 || m.GlyphHeight <= 0 {
		t.Errorf("expected valid fallback metrics after invalid SetFontID, got %+v", m)
	}
}

// TestFontMetricsValid verifies that FontMetrics returns non-zero values for a
// surface with a font set.

func TestFontMetricsValid(t *testing.T) {
	s := New(image.Rect(0, 0, 128, 64))

	// Set a known font explicitly.
	ok := s.SetFontID(font.Spleen5x8ID)
	if !ok {
		t.Fatalf("SetFontID(%q) returned false", font.Spleen5x8ID)
	}

	m := s.FontMetrics()
	if m.GlyphWidth <= 0 {
		t.Errorf("GlyphWidth = %d, want > 0", m.GlyphWidth)
	}
	if m.GlyphHeight <= 0 {
		t.Errorf("GlyphHeight = %d, want > 0", m.GlyphHeight)
	}
	if m.GlyphAdvance <= 0 {
		t.Errorf("GlyphAdvance = %d, want > 0", m.GlyphAdvance)
	}
	if m.RowHeight <= 0 {
		t.Errorf("RowHeight = %d, want > 0", m.RowHeight)
	}

	// The returned metrics should match what the face itself reports.
	face, _ := font.Get(font.Spleen5x8ID)
	expected := face.Metrics()
	if m != expected {
		t.Errorf("FontMetrics = %+v, want %+v", m, expected)
	}
}

// ---------------------------------------------------------------------------
// Tests merged from drawtext_test.go
// ---------------------------------------------------------------------------

// TestDrawText_PixelsRendered verifies that DrawText actually sets foreground
// pixels on the framebuffer when drawing a non-space character. This catches
// bitmask alignment bugs where GlyphRow data uses high bits (bit 31 = leftmost)
// but the renderer checks low bits.
func TestDrawText_PixelsRendered(t *testing.T) {
	// Use spleen-5x8 which is always registered.
	face, ok := font.Get("spleen-5x8")
	if !ok {
		t.Fatal("spleen-5x8 not registered")
	}

	m := face.Metrics()
	width := m.GlyphAdvance * 3
	height := m.RowHeight

	s := New(image.Rect(0, 0, width, height))
	s.SetFontID("spleen-5x8")
	s.DrawText(0, 0, "A", color.White)

	// Count foreground pixels in the glyph area.
	fb := s.FrameBuffer()
	var fgCount int
	for y := 0; y < m.GlyphHeight; y++ {
		for x := 0; x < m.GlyphWidth; x++ {
			r, g, b, _ := fb.At(x, y).RGBA()
			if r > 0 || g > 0 || b > 0 {
				fgCount++
			}
		}
	}

	if fgCount == 0 {
		t.Fatal("DrawText produced zero foreground pixels for 'A'; bitmask alignment is broken")
	}

	// Spleen 5x8 'A' should produce a reasonable number of pixels.
	if fgCount < 5 {
		t.Errorf("DrawText produced only %d foreground pixels for 'A'; expected at least 5", fgCount)
	}
}

// TestDrawText_EmptyStringNoPixels verifies that drawing an empty string
// doesn't set any pixels.
func TestDrawText_EmptyStringNoPixels(t *testing.T) {
	s := New(image.Rect(0, 0, 50, 20))
	s.SetFontID("spleen-5x8")
	s.DrawText(0, 0, "", color.White)

	fb := s.FrameBuffer()
	bounds := fb.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := fb.At(x, y).RGBA()
			if r > 0 || g > 0 || b > 0 {
				t.Fatalf("found foreground pixel at (%d,%d) for empty string", x, y)
			}
		}
	}
}

// ---------------------------------------------------------------------------
// Tests merged from surface_draw_test.go
// ---------------------------------------------------------------------------

// For any image.Image of any concrete type (RGBA, Gray, Alpha) and for any destination
// image.Point (including negative and out-of-bounds positions), calling DrawImage SHALL
// produce a framebuffer where each destination pixel (dx, dy) within the surface bounds
// contains the composited source pixel at the corresponding source coordinate, and no pixel
// outside the intersection of the source image bounds (offset by destination) and the
// surface bounds is modified.

// drawImageInput holds a randomly generated test input for the DrawImage property test.
type drawImageInput struct {
	SurfW, SurfH int   // Surface dimensions (1..64)
	SrcW, SrcH   int   // Source image dimensions (1..32)
	DstX, DstY   int   // Destination point (-48..80)
	ImgType      int   // 0=RGBA, 1=Gray, 2=Alpha
	Seed         int64 // RNG seed for reproducible image generation
	BgR, BgG     uint8 // Background color components
	BgB, BgA     uint8
}

// Generate implements quick.Generator for property-based testing.
func (drawImageInput) Generate(r *rand.Rand, size int) reflect.Value {
	input := drawImageInput{
		SurfW:   1 + r.Intn(64),
		SurfH:   1 + r.Intn(64),
		SrcW:    1 + r.Intn(32),
		SrcH:    1 + r.Intn(32),
		DstX:    r.Intn(129) - 48, // range: -48..80
		DstY:    r.Intn(129) - 48,
		ImgType: r.Intn(3),
		Seed:    r.Int63(),
		BgR:     uint8(r.Intn(256)),
		BgG:     uint8(r.Intn(256)),
		BgB:     uint8(r.Intn(256)),
		BgA:     uint8(r.Intn(256)),
	}
	return reflect.ValueOf(input)
}

// generateImage creates a random image of the specified type with the given dimensions.
func generateImage(rng *rand.Rand, imgType, w, h int) image.Image {
	switch imgType {
	case 0: // RGBA
		img := image.NewRGBA(image.Rect(0, 0, w, h))
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				img.SetRGBA(x, y, color.RGBA{
					R: uint8(rng.Intn(256)),
					G: uint8(rng.Intn(256)),
					B: uint8(rng.Intn(256)),
					A: uint8(rng.Intn(256)),
				})
			}
		}
		return img
	case 1: // Gray
		img := image.NewGray(image.Rect(0, 0, w, h))
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				img.SetGray(x, y, color.Gray{Y: uint8(rng.Intn(256))})
			}
		}
		return img
	default: // Alpha
		img := image.NewAlpha(image.Rect(0, 0, w, h))
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				img.SetAlpha(x, y, color.Alpha{A: uint8(rng.Intn(256))})
			}
		}
		return img
	}
}

func TestDrawImagePixelCorrectnessWithClipping(t *testing.T) {
	config := &quick.Config{MaxCount: 200}

	prop := func(input drawImageInput) bool {
		// Create a surface and fill it with the background color.
		surfBounds := image.Rect(0, 0, input.SurfW, input.SurfH)
		s := New(surfBounds)
		bgColor := color.RGBA{R: input.BgR, G: input.BgG, B: input.BgB, A: input.BgA}
		s.Clear(bgColor)

		// Snapshot the background framebuffer before drawing.
		bgSnapshot := image.NewRGBA(surfBounds)
		draw.Draw(bgSnapshot, surfBounds, s.FrameBuffer(), image.Point{}, draw.Src)

		// Generate a random source image.
		rng := rand.New(rand.NewSource(input.Seed))
		src := generateImage(rng, input.ImgType, input.SrcW, input.SrcH)

		// Draw the image.
		dst := image.Pt(input.DstX, input.DstY)
		s.DrawImage(src, dst)

		// Compute the expected intersection of the source image (offset by dst) and the surface bounds.
		srcBounds := src.Bounds()
		drawnRect := image.Rectangle{
			Min: dst,
			Max: dst.Add(image.Pt(srcBounds.Dx(), srcBounds.Dy())),
		}.Intersect(surfBounds)

		// Build an expected framebuffer by compositing the source over the background
		// using the same draw.Over operation.
		expected := image.NewRGBA(surfBounds)
		draw.Draw(expected, surfBounds, bgSnapshot, image.Point{}, draw.Src)
		if !drawnRect.Empty() {
			dstRect := image.Rectangle{
				Min: dst,
				Max: dst.Add(image.Pt(srcBounds.Dx(), srcBounds.Dy())),
			}
			draw.Draw(expected, dstRect, src, srcBounds.Min, draw.Over)
		}

		fb := s.FrameBuffer()

		// Check every pixel in the surface.
		for y := surfBounds.Min.Y; y < surfBounds.Max.Y; y++ {
			for x := surfBounds.Min.X; x < surfBounds.Max.X; x++ {
				got := fb.RGBAAt(x, y)
				want := expected.RGBAAt(x, y)

				if got != want {
					pt := image.Pt(x, y)
					inDrawnRect := pt.In(drawnRect)
					t.Errorf("pixel (%d,%d) mismatch: got %v, want %v (inDrawnRect=%v, dst=%v, srcBounds=%v)",
						x, y, got, want, inDrawnRect, dst, srcBounds)
					return false
				}
			}
		}

		return true
	}

	if err := quick.Check(prop, config); err != nil {
		t.Errorf("Property 1 failed: %v", err)
	}
}

// For any background pixel color and for any source pixel with an alpha channel value,
// the resulting framebuffer pixel after DrawImage SHALL equal the result of the
// Porter-Duff source-over formula: result = src + dst * (1 - srcAlpha), computed per RGBA channel.

func TestPropertyDrawImageAlphaCompositing(t *testing.T) {
	cfg := &quick.Config{MaxCount: 200}

	f := func(dstR, dstG, dstB, dstA, srcR, srcG, srcB, srcA uint8) bool {
		// Create a 1x1 Surface and fill with the background color (premultiplied RGBA).
		surf := New(image.Rect(0, 0, 1, 1))
		bgColor := color.NRGBA{R: dstR, G: dstG, B: dstB, A: dstA}
		draw.Draw(surf.fb, surf.fb.Bounds(), image.NewUniform(bgColor), image.Point{}, draw.Src)

		// Record the actual background pixel the surface now holds (premultiplied 8-bit).
		bgPx := surf.fb.RGBAAt(0, 0)

		// Create a 1x1 NRGBA source image with the given color and alpha.
		srcImg := image.NewNRGBA(image.Rect(0, 0, 1, 1))
		srcImg.SetNRGBA(0, 0, color.NRGBA{R: srcR, G: srcG, B: srcB, A: srcA})

		// Call DrawImage to composite src over dst.
		surf.DrawImage(srcImg, image.Point{})

		// Read back the resulting pixel in premultiplied form.
		got := surf.fb.RGBAAt(0, 0)

		// Compute expected using Porter-Duff source-over, matching Go's image/draw
		// which works in 16-bit premultiplied space.
		//
		// Go's draw.Over reads both source and destination via RGBA() returning
		// premultiplied 16-bit values, applies the formula, then stores back as 8-bit.

		// Source NRGBA -> premultiplied 16-bit (matches NRGBA.RGBA())
		sr, sg, sb, sa := color.NRGBA{R: srcR, G: srcG, B: srcB, A: srcA}.RGBA()

		// Destination premultiplied RGBA -> 16-bit (matches RGBA.RGBA())
		dr, dg, db, da := bgPx.RGBA()

		// Porter-Duff source-over: out = src + dst * (0xffff - srcAlpha) / 0xffff
		a := 0xffff - sa
		outR := sr + dr*a/0xffff
		outG := sg + dg*a/0xffff
		outB := sb + db*a/0xffff
		outA := sa + da*a/0xffff

		// Convert back to 8-bit (Go uses >> 8 for the high byte).
		wantR := uint8(outR >> 8)
		wantG := uint8(outG >> 8)
		wantB := uint8(outB >> 8)
		wantA := uint8(outA >> 8)

		// Allow ±1 tolerance per channel due to integer arithmetic rounding.
		if !within1(got.R, wantR) || !within1(got.G, wantG) ||
			!within1(got.B, wantB) || !within1(got.A, wantA) {
			t.Logf("src NRGBA(%d,%d,%d,%d) over dst RGBA(%d,%d,%d,%d): got RGBA(%d,%d,%d,%d), want RGBA(%d,%d,%d,%d)",
				srcR, srcG, srcB, srcA, bgPx.R, bgPx.G, bgPx.B, bgPx.A,
				got.R, got.G, got.B, got.A, wantR, wantG, wantB, wantA)
			return false
		}
		return true
	}

	if err := quick.Check(f, cfg); err != nil {
		t.Error(err)
	}
}

// within1 checks if two uint8 values are within ±1 of each other.
func within1(a, b uint8) bool {
	diff := int(a) - int(b)
	return diff >= -1 && diff <= 1
}
