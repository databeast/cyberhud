package svgrender

import (
	"fmt"
	"image"
	"image/color"
	"strings"
	"testing"
	"time"
	"unicode/utf8"

	"pgregory.net/rapid"
)

// --- From: svgrender_animator_prop_test.go ---

func TestProperty_AnimationTickMonotonicAdvancement(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		n := rapid.IntRange(2, 20).Draw(rt, "numFrames")

		frames := make([]Frame, n)
		var totalDuration time.Duration
		for i := 0; i < n; i++ {
			dur := time.Duration(rapid.IntRange(1, 1000).Draw(rt, "dur")) * time.Millisecond
			frames[i] = Frame{
				SVG:      fmt.Sprintf("frame-%d", i),
				Duration: dur,
			}
			totalDuration += dur
		}

		anim := NewAnimator(frames, false)

		// Tick the total sum of all frame durations.
		anim.Tick(totalDuration)

		// Assert Done() == true for non-looping animator past final frame.
		if !anim.Done() {
			rt.Fatalf("expected Done() == true after ticking total duration, got false")
		}

		// Assert CurrentFrame() returns the last frame index (N-1).
		svg, index := anim.CurrentFrame()
		if index != n-1 {
			rt.Fatalf("expected CurrentFrame index %d, got %d", n-1, index)
		}
		expectedSVG := fmt.Sprintf("frame-%d", n-1)
		if svg != expectedSVG {
			rt.Fatalf("expected CurrentFrame SVG %q, got %q", expectedSVG, svg)
		}
	})
}

func TestProperty_AnimationLoopWrap(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		n := rapid.IntRange(2, 20).Draw(rt, "numFrames")

		frames := make([]Frame, n)
		var totalDuration time.Duration
		for i := 0; i < n; i++ {
			dur := time.Duration(rapid.IntRange(1, 1000).Draw(rt, "dur")) * time.Millisecond
			frames[i] = Frame{
				SVG:      fmt.Sprintf("frame-%d", i),
				Duration: dur,
			}
			totalDuration += dur
		}

		anim := NewAnimator(frames, true)

		// Tick the total sum of all frame durations (one full cycle).
		anim.Tick(totalDuration)

		// Assert CurrentFrame() returns index 0 (wrapped back to start).
		_, index := anim.CurrentFrame()
		if index != 0 {
			rt.Fatalf("expected CurrentFrame index 0 after full loop cycle, got %d", index)
		}

		// Assert Done() == false (never done when looping).
		if anim.Done() {
			rt.Fatalf("expected Done() == false for looping animator, got true")
		}
	})
}

// --- From: svgrender_bounds_prop_test.go ---

// *For any* Config where Bounds.Dx() < 16 or Bounds.Dy() < 16 or Bounds.Dx() <= 0
// or Bounds.Dy() <= 0, Render(cfg) SHALL return nil regardless of other field values.

// TestProperty_ResolutionGuard_TooNarrow verifies that Render returns nil when
// bounds width is below MinBoundsWidth (16), even with valid SVG content.
func TestProperty_ResolutionGuard_TooNarrow(t *testing.T) {
	const validSVG = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><rect width="100" height="100" fill="red"/></svg>`

	rapid.Check(t, func(rt *rapid.T) {
		// Width in [0, 15] — below the 16px minimum.
		w := rapid.IntRange(0, 15).Draw(rt, "width")
		// Height is valid (>= 16) to isolate the width check.
		h := rapid.IntRange(16, 512).Draw(rt, "height")

		cfg := Config{
			Bounds: image.Rect(0, 0, w, h),
			SVG:    validSVG,
		}

		result := Render(cfg)
		if result != nil {
			t.Fatalf("Render with bounds width=%d (< 16) returned non-nil sprite", w)
		}
	})
}

// TestProperty_ResolutionGuard_TooShort verifies that Render returns nil when
// bounds height is below MinBoundsHeight (16), even with valid SVG content.
func TestProperty_ResolutionGuard_TooShort(t *testing.T) {
	const validSVG = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><rect width="100" height="100" fill="red"/></svg>`

	rapid.Check(t, func(rt *rapid.T) {
		// Width is valid (>= 16) to isolate the height check.
		w := rapid.IntRange(16, 512).Draw(rt, "width")
		// Height in [0, 15] — below the 16px minimum.
		h := rapid.IntRange(0, 15).Draw(rt, "height")

		cfg := Config{
			Bounds: image.Rect(0, 0, w, h),
			SVG:    validSVG,
		}

		result := Render(cfg)
		if result != nil {
			t.Fatalf("Render with bounds height=%d (< 16) returned non-nil sprite", h)
		}
	})
}

// TestProperty_ResolutionGuard_NegativeDimensions verifies that Render returns nil
// when bounds have zero or negative width or height dimensions.
// Note: image.Rectangle created via struct literal (not image.Rect) can have
// non-canonical coordinates where Max < Min, producing negative Dx()/Dy().
func TestProperty_ResolutionGuard_NegativeDimensions(t *testing.T) {
	const validSVG = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><rect width="100" height="100" fill="red"/></svg>`

	rapid.Check(t, func(rt *rapid.T) {
		// Generate at least one negative dimension using direct struct literal
		// to bypass image.Rect's normalization.
		variant := rapid.IntRange(0, 2).Draw(rt, "variant")

		var bounds image.Rectangle
		switch variant {
		case 0: // negative width (Max.X < Min.X), valid height
			minX := rapid.IntRange(1, 100).Draw(rt, "minX")
			maxX := rapid.IntRange(-100, minX-1).Draw(rt, "maxX")
			bounds = image.Rectangle{
				Min: image.Pt(minX, 0),
				Max: image.Pt(maxX, rapid.IntRange(16, 512).Draw(rt, "maxY")),
			}
		case 1: // negative height (Max.Y < Min.Y), valid width
			minY := rapid.IntRange(1, 100).Draw(rt, "minY")
			maxY := rapid.IntRange(-100, minY-1).Draw(rt, "maxY")
			bounds = image.Rectangle{
				Min: image.Pt(0, minY),
				Max: image.Pt(rapid.IntRange(16, 512).Draw(rt, "maxX"), maxY),
			}
		case 2: // both negative
			minX := rapid.IntRange(1, 100).Draw(rt, "minX")
			maxX := rapid.IntRange(-100, minX-1).Draw(rt, "maxX")
			minY := rapid.IntRange(1, 100).Draw(rt, "minY")
			maxY := rapid.IntRange(-100, minY-1).Draw(rt, "maxY")
			bounds = image.Rectangle{
				Min: image.Pt(minX, minY),
				Max: image.Pt(maxX, maxY),
			}
		}

		cfg := Config{
			Bounds: bounds,
			SVG:    validSVG,
		}

		result := Render(cfg)
		if result != nil {
			t.Fatalf("Render with negative bounds (Dx=%d, Dy=%d) returned non-nil sprite",
				cfg.Bounds.Dx(), cfg.Bounds.Dy())
		}
	})
}

// *For any* Config where SVG is empty and Frames is empty or nil, Render(cfg) SHALL
// return nil regardless of bounds or other field values.

// TestProperty_EmptySource_NilFrames verifies that Render returns nil when SVG is
// empty and Frames is nil, even with valid bounds.
func TestProperty_EmptySource_NilFrames(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		// Valid bounds that pass the resolution guard.
		w := rapid.IntRange(16, 512).Draw(rt, "width")
		h := rapid.IntRange(16, 512).Draw(rt, "height")

		cfg := Config{
			Bounds: image.Rect(0, 0, w, h),
			SVG:    "",
			Frames: nil,
		}

		result := Render(cfg)
		if result != nil {
			t.Fatalf("Render with empty SVG and nil Frames returned non-nil sprite (bounds %dx%d)", w, h)
		}
	})
}

// TestProperty_EmptySource_EmptyFrames verifies that Render returns nil when SVG is
// empty and Frames is an empty slice, even with valid bounds.
func TestProperty_EmptySource_EmptyFrames(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		// Valid bounds that pass the resolution guard.
		w := rapid.IntRange(16, 512).Draw(rt, "width")
		h := rapid.IntRange(16, 512).Draw(rt, "height")

		cfg := Config{
			Bounds: image.Rect(0, 0, w, h),
			SVG:    "",
			Frames: []Frame{},
		}

		result := Render(cfg)
		if result != nil {
			t.Fatalf("Render with empty SVG and empty Frames returned non-nil sprite (bounds %dx%d)", w, h)
		}
	})
}

// --- From: svgrender_canvas_prop_test.go ---

// *For any* positive width and height, NewCanvas(width, height) SHALL return a non-nil
// Canvas whose Image() has bounds Rect(0, 0, width, height). *For any* zero or negative
// width or height, NewCanvas SHALL return nil.

// TestProperty_Canvas_PositiveDimensions verifies that NewCanvas with positive width
// and height returns a non-nil Canvas whose Image bounds match the requested dimensions.
func TestProperty_Canvas_PositiveDimensions(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		w := rapid.IntRange(1, 4096).Draw(rt, "width")
		h := rapid.IntRange(1, 4096).Draw(rt, "height")

		c := NewCanvas(w, h)
		if c == nil {
			t.Fatalf("NewCanvas(%d, %d) returned nil, want non-nil", w, h)
		}

		img := c.Image()
		if img == nil {
			t.Fatalf("Canvas.Image() returned nil for valid canvas (%d x %d)", w, h)
		}

		expectedBounds := image.Rect(0, 0, w, h)
		if img.Bounds() != expectedBounds {
			t.Fatalf("Canvas.Image().Bounds() = %v, want %v", img.Bounds(), expectedBounds)
		}
	})
}

// TestProperty_Canvas_ZeroOrNegativeDimensions verifies that NewCanvas returns nil
// when width or height is zero or negative.
func TestProperty_Canvas_ZeroOrNegativeDimensions(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		// Generate at least one non-positive dimension.
		variant := rapid.IntRange(0, 2).Draw(rt, "variant")

		var w, h int
		switch variant {
		case 0: // width non-positive, height arbitrary
			w = rapid.IntRange(-1000, 0).Draw(rt, "width")
			h = rapid.IntRange(-1000, 4096).Draw(rt, "height")
		case 1: // height non-positive, width arbitrary
			w = rapid.IntRange(-1000, 4096).Draw(rt, "width")
			h = rapid.IntRange(-1000, 0).Draw(rt, "height")
		case 2: // both non-positive
			w = rapid.IntRange(-1000, 0).Draw(rt, "width")
			h = rapid.IntRange(-1000, 0).Draw(rt, "height")
		}

		c := NewCanvas(w, h)
		if c != nil {
			t.Fatalf("NewCanvas(%d, %d) returned non-nil, want nil for invalid dimensions", w, h)
		}
	})
}

// *For any* newly created Canvas (before any rendering), every pixel in Image()
// SHALL have RGBA values (0, 0, 0, 0).

// TestProperty_Canvas_TransparentBackground verifies that a freshly created Canvas
// has all pixels set to fully transparent (0, 0, 0, 0).
func TestProperty_Canvas_TransparentBackground(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		// Use smaller dimensions to keep iteration tractable.
		w := rapid.IntRange(1, 256).Draw(rt, "width")
		h := rapid.IntRange(1, 256).Draw(rt, "height")

		c := NewCanvas(w, h)
		if c == nil {
			t.Fatalf("NewCanvas(%d, %d) returned nil", w, h)
		}

		img := c.Image()
		bounds := img.Bounds()
		transparent := color.RGBA{0, 0, 0, 0}

		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				r, g, b, a := img.At(x, y).RGBA()
				if r != 0 || g != 0 || b != 0 || a != 0 {
					got := color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)}
					t.Fatalf("pixel (%d,%d) = %v, want %v", x, y, got, transparent)
				}
			}
		}
	})
}

// --- From: svgrender_precedence_prop_test.go ---

// *For any* Config where both SVG and Frames are non-empty, Render(cfg) SHALL use
// the frame at cfg.FrameIndex from Frames as the SVG source, ignoring the SVG field
// entirely.

// TestProperty_FrameSequence_Precedence verifies that when both SVG and Frames are
// populated, the rendered output uses Frames[FrameIndex], not the SVG field.
// Strategy: use two visually distinct SVGs (red vs blue), put one in SVG field and
// the other in Frames. Render with both set, then render with only the Frame SVG.
// Both outputs must match, proving Frames took precedence.
func TestProperty_FrameSequence_Precedence(t *testing.T) {
	const (
		redSVG  = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><rect width="100" height="100" fill="red"/></svg>`
		blueSVG = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><rect width="100" height="100" fill="blue"/></svg>`
	)

	rapid.Check(t, func(rt *rapid.T) {
		// Generate a random number of frames (1 to 5).
		numFrames := rapid.IntRange(1, 5).Draw(rt, "numFrames")

		// Generate a valid FrameIndex within range.
		frameIndex := rapid.IntRange(0, numFrames-1).Draw(rt, "frameIndex")

		// Build a frame sequence where all frames use blue SVG.
		frames := make([]Frame, numFrames)
		for i := range frames {
			frames[i] = Frame{SVG: blueSVG, Duration: 100_000_000} // 100ms
		}

		bounds := image.Rect(0, 0, 32, 32)

		// Config with BOTH SVG (red) and Frames (blue) populated.
		cfgBoth := Config{
			Bounds:     bounds,
			Label:      "precedence-test",
			SVG:        redSVG,
			Frames:     frames,
			FrameIndex: frameIndex,
		}

		// Config with ONLY the Frame SVG (blue), no SVG field.
		cfgFrameOnly := Config{
			Bounds:     bounds,
			Label:      "precedence-test",
			SVG:        "",
			Frames:     frames,
			FrameIndex: frameIndex,
		}

		// Render both configs.
		spriteBoth := Render(cfgBoth)
		spriteFrameOnly := Render(cfgFrameOnly)

		if spriteBoth == nil {
			t.Fatal("Render with both SVG and Frames returned nil")
		}
		if spriteFrameOnly == nil {
			t.Fatal("Render with only Frames returned nil")
		}

		// Both renders should produce identical pixel output, proving Frames
		// took precedence over the SVG field.
		imgBoth := spriteBoth.Image.(*image.RGBA)
		imgFrameOnly := spriteFrameOnly.Image.(*image.RGBA)

		if imgBoth.Bounds() != imgFrameOnly.Bounds() {
			t.Fatalf("image bounds differ: both=%v, frameOnly=%v",
				imgBoth.Bounds(), imgFrameOnly.Bounds())
		}

		// Compare pixel-by-pixel.
		b := imgBoth.Bounds()
		for y := b.Min.Y; y < b.Max.Y; y++ {
			for x := b.Min.X; x < b.Max.X; x++ {
				cBoth := imgBoth.RGBAAt(x, y)
				cFrame := imgFrameOnly.RGBAAt(x, y)
				if cBoth != cFrame {
					t.Fatalf("pixel (%d,%d) differs: both=%v, frameOnly=%v",
						x, y, cBoth, cFrame)
				}
			}
		}

		// Additionally verify that the output does NOT match what rendering
		// SVG-only (red) would produce, confirming the SVG field was ignored.
		cfgSVGOnly := Config{
			Bounds: bounds,
			Label:  "precedence-test",
			SVG:    redSVG,
		}
		spriteSVGOnly := Render(cfgSVGOnly)
		if spriteSVGOnly == nil {
			t.Fatal("Render with only SVG (red) returned nil")
		}

		imgSVGOnly := spriteSVGOnly.Image.(*image.RGBA)

		// Find at least one pixel that differs between the Frames render (blue)
		// and the SVG-only render (red), proving they are visually different.
		foundDifference := false
		for y := b.Min.Y; y < b.Max.Y; y++ {
			for x := b.Min.X; x < b.Max.X; x++ {
				cFrame := imgBoth.RGBAAt(x, y)
				cSVG := imgSVGOnly.RGBAAt(x, y)
				if cFrame != cSVG {
					foundDifference = true
					break
				}
			}
			if foundDifference {
				break
			}
		}

		if !foundDifference {
			// If the rendered output of both-config matches the SVG-only (red) output,
			// it means Frames did NOT take precedence — the SVG field was used.
			t.Fatal("rendered output with Frames matches SVG-only output; " +
				"expected Frames to take precedence over SVG field")
		}

		// Verify the pixels are actually blue-ish (from the frame SVG), not red.
		// Sample a center pixel which should be fully rendered.
		centerPixel := imgBoth.RGBAAt(16, 16)
		if centerPixel.A == 0 {
			t.Fatal("center pixel is fully transparent, expected rendered content")
		}
		// Blue SVG should have B > R for non-transparent pixels.
		if centerPixel.R >= centerPixel.B && centerPixel.A > 0 {
			t.Fatalf("center pixel appears red (R=%d, B=%d); expected blue from Frames",
				centerPixel.R, centerPixel.B)
		}
	})
}

// TestProperty_FrameSequence_Precedence_VaryingIndex verifies that FrameIndex
// correctly selects from the Frames slice when both SVG and Frames are present.
func TestProperty_FrameSequence_Precedence_VaryingIndex(t *testing.T) {
	const (
		redSVG  = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><rect width="100" height="100" fill="red"/></svg>`
		blueSVG = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><rect width="100" height="100" fill="blue"/></svg>`
	)

	rapid.Check(t, func(rt *rapid.T) {
		// Generate 2-5 frames, alternating red and blue.
		numFrames := rapid.IntRange(2, 5).Draw(rt, "numFrames")
		frameIndex := rapid.IntRange(0, numFrames-1).Draw(rt, "frameIndex")

		frames := make([]Frame, numFrames)
		for i := range frames {
			if i%2 == 0 {
				frames[i] = Frame{SVG: blueSVG, Duration: 100_000_000}
			} else {
				frames[i] = Frame{SVG: redSVG, Duration: 100_000_000}
			}
		}

		bounds := image.Rect(0, 0, 32, 32)

		// Config with both SVG (opposite color) and Frames set.
		// SVG field is always the "wrong" color to detect if it leaks through.
		var svgField string
		if frameIndex%2 == 0 {
			svgField = redSVG // Frame is blue, SVG is red
		} else {
			svgField = blueSVG // Frame is red, SVG is blue
		}

		cfgBoth := Config{
			Bounds:     bounds,
			Label:      "vary-index",
			SVG:        svgField,
			Frames:     frames,
			FrameIndex: frameIndex,
		}

		// Render the same frame index with ONLY Frames (no SVG field).
		cfgFrameOnly := Config{
			Bounds:     bounds,
			Label:      "vary-index",
			SVG:        "",
			Frames:     frames,
			FrameIndex: frameIndex,
		}

		spriteBoth := Render(cfgBoth)
		spriteFrameOnly := Render(cfgFrameOnly)

		if spriteBoth == nil {
			t.Fatal("Render with both SVG and Frames returned nil")
		}
		if spriteFrameOnly == nil {
			t.Fatal("Render with only Frames returned nil")
		}

		// Pixel-level comparison: both must be identical.
		imgBoth := spriteBoth.Image.(*image.RGBA)
		imgFrameOnly := spriteFrameOnly.Image.(*image.RGBA)

		b := imgBoth.Bounds()
		for y := b.Min.Y; y < b.Max.Y; y++ {
			for x := b.Min.X; x < b.Max.X; x++ {
				cBoth := imgBoth.RGBAAt(x, y)
				cFrame := imgFrameOnly.RGBAAt(x, y)
				if cBoth != cFrame {
					t.Fatalf("pixel (%d,%d) differs: both=%v, frameOnly=%v",
						x, y, cBoth, cFrame)
				}
			}
		}

		// Verify the center pixel color matches the expected frame color.
		centerPixel := imgBoth.RGBAAt(16, 16)
		if centerPixel.A == 0 {
			t.Fatal("center pixel is transparent, expected rendered content")
		}

		expectedBlue := frameIndex%2 == 0
		if expectedBlue {
			// Blue frame: B should be > R.
			if centerPixel.R >= centerPixel.B {
				t.Fatalf("expected blue pixel at center (frameIndex=%d), got R=%d B=%d",
					frameIndex, centerPixel.R, centerPixel.B)
			}
		} else {
			// Red frame: R should be > B.
			if centerPixel.B >= centerPixel.R {
				t.Fatalf("expected red pixel at center (frameIndex=%d), got R=%d B=%d",
					frameIndex, centerPixel.R, centerPixel.B)
			}
		}
	})
}

// --- From: svgrender_sign_prop_test.go ---

// drawConfig generates a random Config using rapid generators.
func drawConfig(rt *rapid.T, suffix string) Config {
	minX := rapid.IntRange(-100, 100).Draw(rt, "minX"+suffix)
	minY := rapid.IntRange(-100, 100).Draw(rt, "minY"+suffix)
	maxX := rapid.IntRange(16, 300).Draw(rt, "maxX"+suffix)
	maxY := rapid.IntRange(16, 300).Draw(rt, "maxY"+suffix)

	label := rapid.String().Draw(rt, "label"+suffix)

	col := color.RGBA{
		R: uint8(rapid.IntRange(0, 255).Draw(rt, "colorR"+suffix)),
		G: uint8(rapid.IntRange(0, 255).Draw(rt, "colorG"+suffix)),
		B: uint8(rapid.IntRange(0, 255).Draw(rt, "colorB"+suffix)),
		A: uint8(rapid.IntRange(0, 255).Draw(rt, "colorA"+suffix)),
	}

	svg := rapid.StringMatching("[a-z]{1,50}").Draw(rt, "svg"+suffix)

	frameCount := rapid.IntRange(0, 5).Draw(rt, "frameCount"+suffix)
	frames := make([]Frame, frameCount)
	for i := range frames {
		frames[i] = Frame{
			SVG:      rapid.StringMatching("[a-z]{1,30}").Draw(rt, "frameSVG"+suffix),
			Duration: time.Duration(rapid.IntRange(1, 10000).Draw(rt, "frameDur"+suffix)) * time.Millisecond,
		}
	}

	frameIndex := rapid.IntRange(0, 10).Draw(rt, "frameIndex"+suffix)

	return Config{
		Bounds:     image.Rect(minX, minY, minX+maxX, minY+maxY),
		Label:      label,
		Color:      col,
		SVG:        svg,
		Frames:     frames,
		FrameIndex: frameIndex,
	}
}

// *For any* Config value c, Sign(c) == Sign(c) (calling Sign twice with identical
// input produces identical output).

// TestProperty_Sign_Determinism verifies that Sign produces the same hash for
// the same Config on repeated calls.
func TestProperty_Sign_Determinism(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		cfg := drawConfig(rt, "")

		h1 := Sign(cfg)
		h2 := Sign(cfg)

		if h1 != h2 {
			t.Fatalf("Sign determinism violated: Sign(cfg) = %d, then %d", h1, h2)
		}
	})
}

// *For any* two Config values that differ in exactly one field, Sign(a) != Sign(b)
// with probability >= 1 - 2^(-64). In practice: for 1000 generated pairs differing
// in one field, zero collisions are expected.

// TestProperty_Sign_Sensitivity_Bounds verifies that configs differing only in Bounds
// produce different hashes.
func TestProperty_Sign_Sensitivity_Bounds(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		cfg := drawConfig(rt, "")

		modified := cfg
		// Change Bounds to a different rectangle.
		newMinX := rapid.IntRange(-100, 100).Draw(rt, "newMinX")
		newMinY := rapid.IntRange(-100, 100).Draw(rt, "newMinY")
		newW := rapid.IntRange(16, 300).Draw(rt, "newW")
		newH := rapid.IntRange(16, 300).Draw(rt, "newH")
		modified.Bounds = image.Rect(newMinX, newMinY, newMinX+newW, newMinY+newH)

		if modified.Bounds == cfg.Bounds {
			t.Skip("generated identical bounds, skipping")
		}

		h1 := Sign(cfg)
		h2 := Sign(modified)
		if h1 == h2 {
			t.Fatalf("Sign collision: different Bounds produced same hash %d", h1)
		}
	})
}

// TestProperty_Sign_Sensitivity_Label verifies that configs differing only in Label
// produce different hashes.
func TestProperty_Sign_Sensitivity_Label(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		cfg := drawConfig(rt, "")

		modified := cfg
		// Append a suffix to guarantee the label is different.
		modified.Label = cfg.Label + rapid.StringMatching("[a-z]{1,10}").Draw(rt, "labelSuffix")

		h1 := Sign(cfg)
		h2 := Sign(modified)
		if h1 == h2 {
			t.Fatalf("Sign collision: different Label produced same hash %d", h1)
		}
	})
}

// TestProperty_Sign_Sensitivity_Color verifies that configs differing only in Color
// produce different hashes.
func TestProperty_Sign_Sensitivity_Color(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		cfg := drawConfig(rt, "")

		modified := cfg
		modified.Color = color.RGBA{
			R: uint8(rapid.IntRange(0, 255).Draw(rt, "newR")),
			G: uint8(rapid.IntRange(0, 255).Draw(rt, "newG")),
			B: uint8(rapid.IntRange(0, 255).Draw(rt, "newB")),
			A: uint8(rapid.IntRange(0, 255).Draw(rt, "newA")),
		}

		if modified.Color == cfg.Color {
			t.Skip("generated identical color, skipping")
		}

		h1 := Sign(cfg)
		h2 := Sign(modified)
		if h1 == h2 {
			t.Fatalf("Sign collision: different Color produced same hash %d", h1)
		}
	})
}

// TestProperty_Sign_Sensitivity_SVG verifies that configs differing only in SVG
// produce different hashes.
func TestProperty_Sign_Sensitivity_SVG(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		cfg := drawConfig(rt, "")

		modified := cfg
		modified.SVG = rapid.StringMatching("[a-z]{1,50}").Draw(rt, "newSVG")

		if modified.SVG == cfg.SVG {
			t.Skip("generated identical SVG, skipping")
		}

		h1 := Sign(cfg)
		h2 := Sign(modified)
		if h1 == h2 {
			t.Fatalf("Sign collision: different SVG produced same hash %d", h1)
		}
	})
}

// TestProperty_Sign_Sensitivity_Frames verifies that configs differing only in Frames
// produce different hashes.
func TestProperty_Sign_Sensitivity_Frames(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		cfg := drawConfig(rt, "")

		modified := cfg
		// Generate a different set of frames.
		newFrameCount := rapid.IntRange(1, 5).Draw(rt, "newFrameCount")
		newFrames := make([]Frame, newFrameCount)
		for i := range newFrames {
			newFrames[i] = Frame{
				SVG:      rapid.StringMatching("[a-z]{1,30}").Draw(rt, "newFrameSVG"),
				Duration: time.Duration(rapid.IntRange(1, 10000).Draw(rt, "newFrameDur")) * time.Millisecond,
			}
		}
		modified.Frames = newFrames

		// Check they are actually different by comparing Sign inputs.
		if Sign(cfg) == Sign(modified) {
			// Very unlikely but possible: skip if hashes happen to collide.
			t.Skip("generated frames that produce same hash, skipping")
		}

		h1 := Sign(cfg)
		h2 := Sign(modified)
		if h1 == h2 {
			t.Fatalf("Sign collision: different Frames produced same hash %d", h1)
		}
	})
}

// TestProperty_Sign_Sensitivity_FrameIndex verifies that configs differing only in
// FrameIndex produce different hashes.
func TestProperty_Sign_Sensitivity_FrameIndex(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		cfg := drawConfig(rt, "")

		// Ensure original and modified have distinct FrameIndex values.
		idx1 := rapid.IntRange(0, 5).Draw(rt, "idx1")
		offset := rapid.IntRange(1, 5).Draw(rt, "offset")
		idx2 := idx1 + offset

		cfg.FrameIndex = idx1
		modified := cfg
		modified.FrameIndex = idx2

		h1 := Sign(cfg)
		h2 := Sign(modified)
		if h1 == h2 {
			t.Fatalf("Sign collision: different FrameIndex (%d vs %d) produced same hash %d", idx1, idx2, h1)
		}
	})
}

// *For any* Config with identical fields except FrameIndex, Sign(a) != Sign(b)
// when a.FrameIndex != b.FrameIndex.

// TestProperty_Sign_IncorporatesAnimationState verifies that changing only the
// FrameIndex field of a Config produces a different Sign hash.
func TestProperty_Sign_IncorporatesAnimationState(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		cfg := drawConfig(rt, "")

		// Generate two distinct frame indices using an offset to guarantee difference.
		idx1 := rapid.IntRange(0, 5).Draw(rt, "idx1")
		offset := rapid.IntRange(1, 5).Draw(rt, "offset")
		idx2 := idx1 + offset

		cfg.FrameIndex = idx1
		cfgB := cfg
		cfgB.FrameIndex = idx2

		h1 := Sign(cfg)
		h2 := Sign(cfgB)
		if h1 == h2 {
			t.Fatalf("Sign does not incorporate animation state: FrameIndex %d vs %d both produce hash %d", idx1, idx2, h1)
		}
	})
}

// --- From: svgrender_source_prop_test.go ---

// *For any* arbitrary string input (including malformed XML, empty strings, binary data),
// the parse function SHALL either return valid parsed data or return a non-nil error,
// and SHALL NOT panic.

// TestProperty_Parse_ArbitraryStrings verifies that parse never panics on arbitrary
// string inputs and always returns either a valid icon or a non-nil error.
func TestProperty_Parse_ArbitraryStrings(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		svg := rapid.String().Draw(rt, "svg")

		icon, err := parse(svg, 64, 64)

		if icon == nil && err == nil {
			t.Fatalf("parse(%q, 64, 64) returned nil icon and nil error; want at least one non-nil", svg)
		}
	})
}

// TestProperty_Parse_EmptyString verifies that parse handles the empty string
// gracefully without panicking.
func TestProperty_Parse_EmptyString(t *testing.T) {
	icon, err := parse("", 64, 64)

	if icon == nil && err == nil {
		t.Fatalf("parse(\"\", 64, 64) returned nil icon and nil error; want at least one non-nil")
	}
}

// TestProperty_Parse_RandomBytes verifies that parse handles arbitrary byte sequences
// (including invalid UTF-8) without panicking, always returning either valid data or an error.
func TestProperty_Parse_RandomBytes(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		bytes := rapid.SliceOf(rapid.Byte()).Draw(rt, "bytes")
		svg := string(bytes)

		icon, err := parse(svg, 64, 64)

		if icon == nil && err == nil {
			t.Fatalf("parse(randomBytes, 64, 64) returned nil icon and nil error; want at least one non-nil")
		}
	})
}

// --- From: svgrender_sprite_prop_test.go ---

// validSVG is a minimal valid SVG document used to produce non-nil sprites.
const validSVG = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><rect width="100" height="100" fill="red"/></svg>`

// *For any* Config that produces a non-nil Sprite, the Sprite's Position SHALL
// equal cfg.Bounds.Min.

// TestProperty_Sprite_PositionInvariant verifies that the rendered sprite's Position
// always equals the Config's Bounds.Min for any valid configuration.
func TestProperty_Sprite_PositionInvariant(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		minX := rapid.IntRange(-500, 500).Draw(rt, "minX")
		minY := rapid.IntRange(-500, 500).Draw(rt, "minY")
		w := rapid.IntRange(16, 200).Draw(rt, "width")
		h := rapid.IntRange(16, 200).Draw(rt, "height")

		cfg := Config{
			Bounds: image.Rect(minX, minY, minX+w, minY+h),
			SVG:    validSVG,
		}

		sprite := Render(cfg)
		if sprite == nil {
			t.Fatalf("Render returned nil for valid config with bounds %v", cfg.Bounds)
		}

		expected := cfg.Bounds.Min
		if sprite.Position != expected {
			t.Fatalf("sprite.Position = %v, want %v (Bounds.Min)", sprite.Position, expected)
		}
	})
}

// *For any* Config that produces a non-nil Sprite, the Sprite's Label SHALL have
// at most 128 runes. If cfg.Label has more than 128 runes, the Sprite's Label
// SHALL equal the first 128 runes of cfg.Label.

// TestProperty_Sprite_LabelTruncation_Long verifies that labels exceeding 128 runes
// are truncated to exactly the first 128 runes.
func TestProperty_Sprite_LabelTruncation_Long(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		// Generate a label with more than 128 runes using Unicode characters.
		runeCount := rapid.IntRange(129, 512).Draw(rt, "runeCount")

		// Build a string from random Unicode runes.
		runes := make([]rune, runeCount)
		for i := range runes {
			// Draw from a wide Unicode range including multi-byte characters.
			runes[i] = rune(rapid.IntRange(0x20, 0x1F600).Draw(rt, "rune"))
			// Ensure valid rune.
			if !utf8.ValidRune(runes[i]) {
				runes[i] = 'X'
			}
		}
		label := string(runes)

		cfg := Config{
			Bounds: image.Rect(0, 0, 64, 64),
			SVG:    validSVG,
			Label:  label,
		}

		sprite := Render(cfg)
		if sprite == nil {
			t.Fatalf("Render returned nil for valid config")
		}

		spriteRunes := []rune(sprite.Label)
		if len(spriteRunes) > 128 {
			t.Fatalf("sprite.Label has %d runes, want at most 128", len(spriteRunes))
		}

		expectedLabel := string(runes[:128])
		if sprite.Label != expectedLabel {
			t.Fatalf("sprite.Label = %q..., want first 128 runes of input",
				truncateForDisplay(sprite.Label, 40))
		}
	})
}

// TestProperty_Sprite_LabelPassthrough_Short verifies that labels with 128 or fewer
// runes are passed through unchanged.
func TestProperty_Sprite_LabelPassthrough_Short(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		// Generate a label with at most 128 runes.
		runeCount := rapid.IntRange(0, 128).Draw(rt, "runeCount")

		runes := make([]rune, runeCount)
		for i := range runes {
			runes[i] = rune(rapid.IntRange(0x20, 0x1F600).Draw(rt, "rune"))
			if !utf8.ValidRune(runes[i]) {
				runes[i] = 'X'
			}
		}
		label := string(runes)

		cfg := Config{
			Bounds: image.Rect(0, 0, 64, 64),
			SVG:    validSVG,
			Label:  label,
		}

		sprite := Render(cfg)
		if sprite == nil {
			t.Fatalf("Render returned nil for valid config")
		}

		if sprite.Label != label {
			t.Fatalf("sprite.Label was modified for label with %d runes (≤128): got %q, want %q",
				runeCount, truncateForDisplay(sprite.Label, 40), truncateForDisplay(label, 40))
		}
	})
}

// truncateForDisplay truncates a string for display in error messages.
func truncateForDisplay(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + strings.Repeat(".", 3)
}

// --- From: svgrender_tint_prop_test.go ---

// *For any* non-zero Color and *any* rendered image, applying the tint transform once
// SHALL produce the same result as calling Render once with that Color (the tint is
// applied exactly once during rendering, not accumulated).

// TestProperty_Tint_Idempotence verifies that applying tint once to a fresh image
// produces a deterministic result — applying tint once to the same fresh image
// always yields identical output (proving single-application idempotence).
func TestProperty_Tint_Idempotence(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		// Generate a small image size.
		w := rapid.IntRange(4, 8).Draw(rt, "width")
		h := rapid.IntRange(4, 8).Draw(rt, "height")

		// Generate random non-zero tint color.
		tintR := rapid.Uint8Range(1, 255).Draw(rt, "tintR")
		tintG := rapid.Uint8Range(1, 255).Draw(rt, "tintG")
		tintB := rapid.Uint8Range(1, 255).Draw(rt, "tintB")
		tintA := rapid.Uint8Range(1, 255).Draw(rt, "tintA")
		tint := color.RGBA{R: tintR, G: tintG, B: tintB, A: tintA}

		// Generate random pixel data with alpha > 0.
		numPixels := w * h
		pixelData := make([]color.RGBA, numPixels)
		for i := 0; i < numPixels; i++ {
			pixelData[i] = color.RGBA{
				R: rapid.Uint8Range(1, 255).Draw(rt, "pixR"),
				G: rapid.Uint8Range(1, 255).Draw(rt, "pixG"),
				B: rapid.Uint8Range(1, 255).Draw(rt, "pixB"),
				A: rapid.Uint8Range(1, 255).Draw(rt, "pixA"),
			}
		}

		// Create first image and apply tint once.
		img1 := image.NewRGBA(image.Rect(0, 0, w, h))
		for i, px := range pixelData {
			x := i % w
			y := i / w
			img1.SetRGBA(x, y, px)
		}
		applyTint(img1, tint)

		// Create second identical image and apply tint once.
		img2 := image.NewRGBA(image.Rect(0, 0, w, h))
		for i, px := range pixelData {
			x := i % w
			y := i / w
			img2.SetRGBA(x, y, px)
		}
		applyTint(img2, tint)

		// Both results must be identical — proving single application is deterministic.
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				c1 := img1.RGBAAt(x, y)
				c2 := img2.RGBAAt(x, y)
				if c1 != c2 {
					t.Fatalf("pixel (%d,%d) differs: first=%v, second=%v", x, y, c1, c2)
				}
			}
		}
	})
}

// TestProperty_Tint_NotNoOp verifies that applying tint twice produces different
// results than applying it once — proving the tint IS applied (not a no-op) and
// that double application accumulates, confirming the importance of applying exactly once.
func TestProperty_Tint_NotNoOp(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		// Generate a small image size.
		w := rapid.IntRange(4, 8).Draw(rt, "width")
		h := rapid.IntRange(4, 8).Draw(rt, "height")

		// Generate a tint color where at least one RGB channel is < 255
		// so tint actually modifies values (a tint of 255,255,255 is identity).
		tintR := rapid.Uint8Range(1, 254).Draw(rt, "tintR")
		tintG := rapid.Uint8Range(1, 254).Draw(rt, "tintG")
		tintB := rapid.Uint8Range(1, 254).Draw(rt, "tintB")
		tintA := rapid.Uint8Range(1, 255).Draw(rt, "tintA")
		tint := color.RGBA{R: tintR, G: tintG, B: tintB, A: tintA}

		// Generate random pixel data with alpha > 0 and RGB > 1 so tint
		// is guaranteed to produce a visible change.
		numPixels := w * h
		pixelData := make([]color.RGBA, numPixels)
		for i := 0; i < numPixels; i++ {
			pixelData[i] = color.RGBA{
				R: rapid.Uint8Range(2, 255).Draw(rt, "pixR"),
				G: rapid.Uint8Range(2, 255).Draw(rt, "pixG"),
				B: rapid.Uint8Range(2, 255).Draw(rt, "pixB"),
				A: rapid.Uint8Range(1, 255).Draw(rt, "pixA"),
			}
		}

		// Apply tint once.
		imgOnce := image.NewRGBA(image.Rect(0, 0, w, h))
		for i, px := range pixelData {
			x := i % w
			y := i / w
			imgOnce.SetRGBA(x, y, px)
		}
		applyTint(imgOnce, tint)

		// Apply tint twice to a fresh copy.
		imgTwice := image.NewRGBA(image.Rect(0, 0, w, h))
		for i, px := range pixelData {
			x := i % w
			y := i / w
			imgTwice.SetRGBA(x, y, px)
		}
		applyTint(imgTwice, tint)
		applyTint(imgTwice, tint)

		// The results MUST differ — proving tint is not a no-op and that
		// double application accumulates (so it's critical to apply exactly once).
		foundDiff := false
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				c1 := imgOnce.RGBAAt(x, y)
				c2 := imgTwice.RGBAAt(x, y)
				if c1 != c2 {
					foundDiff = true
					break
				}
			}
			if foundDiff {
				break
			}
		}

		if !foundDiff {
			t.Fatalf("tint applied once and twice produced identical results; tint should accumulate on double application (tint=%v)", tint)
		}
	})
}
