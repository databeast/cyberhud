package pngpanel

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/databeast/cyberhud/display/surface/textlayout"
	"pgregory.net/rapid"
)

func TestProperty_DimensionRoundTrip(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		width := rapid.IntRange(1, 4096).Draw(t, "width")
		height := rapid.IntRange(1, 4096).Draw(t, "height")
		mode := rapid.SampledFrom([]ColorMode{ColorModeFullColor, ColorModeGrayscale, ColorModeMonochrome}).Draw(t, "colorMode")

		panel, err := New(
			WithDimensions(width, height),
			WithColorMode(mode),
			WithOutputDir(os.TempDir()),
		)
		if err != nil {
			t.Fatalf("unexpected construction error: %v", err)
		}

		bounds := panel.Bounds()
		if bounds.Dx() != width {
			t.Fatalf("Bounds().Dx() = %d, want %d", bounds.Dx(), width)
		}
		if bounds.Dy() != height {
			t.Fatalf("Bounds().Dy() = %d, want %d", bounds.Dy(), height)
		}
		if bounds.Min.X != 0 || bounds.Min.Y != 0 {
			t.Fatalf("Bounds().Min = (%d,%d), want (0,0)", bounds.Min.X, bounds.Min.Y)
		}

		hints := panel.TextHints()
		if hints.PixelWidth != width {
			t.Fatalf("TextHints().PixelWidth = %d, want %d", hints.PixelWidth, width)
		}
		if hints.PixelHeight != height {
			t.Fatalf("TextHints().PixelHeight = %d, want %d", hints.PixelHeight, height)
		}
	})
}

func TestProperty_ThresholdValidation(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		threshold := uint8(rapid.IntRange(0, 255).Draw(t, "threshold"))

		_, err := New(
			WithDimensions(100, 100),
			WithColorMode(ColorModeMonochrome),
			WithThreshold(threshold),
			WithOutputDir(os.TempDir()),
		)
		if err != nil {
			t.Fatalf("expected success for threshold %d, got error: %v", threshold, err)
		}
	})
}

func TestProperty_TextHintsColorMode(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		width := rapid.IntRange(1, 4096).Draw(t, "width")
		height := rapid.IntRange(1, 4096).Draw(t, "height")
		mode := rapid.SampledFrom([]ColorMode{ColorModeFullColor, ColorModeGrayscale, ColorModeMonochrome}).Draw(t, "colorMode")

		panel, err := New(
			WithDimensions(width, height),
			WithColorMode(mode),
			WithOutputDir(os.TempDir()),
		)
		if err != nil {
			t.Fatalf("unexpected construction error: %v", err)
		}

		hints := panel.TextHints()

		// PNGPanel does not own a font; Region fills baseline glyph metrics later.
		if hints.GlyphWidth != 0 {
			t.Fatalf("GlyphWidth = %d, want 0", hints.GlyphWidth)
		}
		if hints.GlyphHeight != 0 {
			t.Fatalf("GlyphHeight = %d, want 0", hints.GlyphHeight)
		}
		if hints.GlyphAdvance != 0 {
			t.Fatalf("GlyphAdvance = %d, want 0", hints.GlyphAdvance)
		}
		if hints.RowHeight != 0 {
			t.Fatalf("RowHeight = %d, want 0", hints.RowHeight)
		}

		// DefaultLineMode is always truncate.
		if hints.DefaultLineMode != textlayout.LineModeTruncate {
			t.Fatalf("DefaultLineMode = %q, want %q", hints.DefaultLineMode, textlayout.LineModeTruncate)
		}

		// Mode-specific assertions.
		if mode == ColorModeMonochrome {
			if !hints.PreferEventRefresh {
				t.Fatal("monochrome: PreferEventRefresh should be true")
			}
			if hints.SupportsVerticalScroll || hints.SupportsHorizontalScroll || hints.SupportsAutoScroll {
				t.Fatal("monochrome: scrolling should be disabled")
			}
			if hints.DefaultTickerDirection != textlayout.TickerDirectionNone {
				t.Fatalf("monochrome: ticker = %q, want %q", hints.DefaultTickerDirection, textlayout.TickerDirectionNone)
			}
		} else {
			if hints.PreferEventRefresh {
				t.Fatal("color/grayscale: PreferEventRefresh should be false")
			}
			if !hints.SupportsVerticalScroll || !hints.SupportsHorizontalScroll || !hints.SupportsAutoScroll {
				t.Fatal("color/grayscale: scrolling should be enabled")
			}
			if hints.DefaultTickerDirection != textlayout.TickerDirectionVertical {
				t.Fatalf("color/grayscale: ticker = %q, want %q", hints.DefaultTickerDirection, textlayout.TickerDirectionVertical)
			}
		}
	})
}

func TestProperty_ConstructionValidation(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		width := rapid.IntRange(-100, 5000).Draw(t, "width")
		height := rapid.IntRange(-100, 5000).Draw(t, "height")

		_, err := New(
			WithDimensions(width, height),
			WithColorMode(ColorModeFullColor),
			WithOutputDir(os.TempDir()),
		)

		valid := width >= 1 && width <= 4096 && height >= 1 && height <= 4096
		if valid && err != nil {
			t.Fatalf("expected success for (%d, %d), got error: %v", width, height, err)
		}
		if !valid && err == nil {
			t.Fatalf("expected error for (%d, %d), got nil", width, height)
		}
	})
}

func TestProperty_DimensionMismatchRejected(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		panelW := rapid.IntRange(1, 4096).Draw(t, "panelW")
		panelH := rapid.IntRange(1, 4096).Draw(t, "panelH")

		// Generate image dimensions that differ from the panel.
		imgW := rapid.IntRange(1, 4096).Draw(t, "imgW")
		imgH := rapid.IntRange(1, 4096).Draw(t, "imgH")

		// Ensure at least one dimension differs.
		if imgW == panelW && imgH == panelH {
			imgW = panelW + 1
			if imgW > 4096 {
				imgW = panelW - 1
			}
		}

		dir, err := os.MkdirTemp("", "pngpanel-test-*")
		if err != nil {
			t.Fatalf("failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(dir)

		panel, err := New(
			WithDimensions(panelW, panelH),
			WithColorMode(ColorModeFullColor),
			WithOutputDir(dir),
		)
		if err != nil {
			t.Fatalf("unexpected construction error: %v", err)
		}

		img := image.NewRGBA(image.Rect(0, 0, imgW, imgH))
		drawErr := panel.DrawImage(img)

		if drawErr == nil {
			t.Fatalf("expected error for panel %dx%d with image %dx%d, got nil", panelW, panelH, imgW, imgH)
		}

		errMsg := drawErr.Error()
		expectedDims := fmt.Sprintf("%dx%d", panelW, panelH)
		receivedDims := fmt.Sprintf("%dx%d", imgW, imgH)
		if !strings.Contains(errMsg, expectedDims) {
			t.Fatalf("error %q does not contain expected dimensions %q", errMsg, expectedDims)
		}
		if !strings.Contains(errMsg, receivedDims) {
			t.Fatalf("error %q does not contain received dimensions %q", errMsg, receivedDims)
		}
	})
}

func TestProperty_PNGDimensionRoundTrip(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		width := rapid.IntRange(1, 200).Draw(t, "width")
		height := rapid.IntRange(1, 200).Draw(t, "height")
		mode := rapid.SampledFrom([]ColorMode{ColorModeFullColor, ColorModeGrayscale, ColorModeMonochrome}).Draw(t, "colorMode")

		dir, err := os.MkdirTemp("", "pngpanel-test-*")
		if err != nil {
			t.Fatalf("failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(dir)

		panel, err := New(
			WithDimensions(width, height),
			WithColorMode(mode),
			WithOutputDir(dir),
		)
		if err != nil {
			t.Fatalf("unexpected construction error: %v", err)
		}

		// Create a random RGBA image matching bounds.
		img := image.NewRGBA(image.Rect(0, 0, width, height))
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				r := uint8(rapid.IntRange(0, 255).Draw(t, fmt.Sprintf("r_%d_%d", x, y)))
				g := uint8(rapid.IntRange(0, 255).Draw(t, fmt.Sprintf("g_%d_%d", x, y)))
				b := uint8(rapid.IntRange(0, 255).Draw(t, fmt.Sprintf("b_%d_%d", x, y)))
				img.SetRGBA(x, y, color.RGBA{R: r, G: g, B: b, A: 255})
			}
		}

		if err := panel.DrawImage(img); err != nil {
			t.Fatalf("DrawImage failed: %v", err)
		}

		// Decode the output PNG and verify dimensions.
		pngPath := filepath.Join(dir, "frame_0001.png")
		f, err := os.Open(pngPath)
		if err != nil {
			t.Fatalf("failed to open output PNG: %v", err)
		}
		defer f.Close()

		decoded, err := png.Decode(f)
		if err != nil {
			t.Fatalf("failed to decode PNG: %v", err)
		}

		if decoded.Bounds().Dx() != width {
			t.Fatalf("decoded width = %d, want %d", decoded.Bounds().Dx(), width)
		}
		if decoded.Bounds().Dy() != height {
			t.Fatalf("decoded height = %d, want %d", decoded.Bounds().Dy(), height)
		}
	})
}

func TestProperty_FullColorPixelRoundTrip(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		width := rapid.IntRange(1, 50).Draw(t, "width")
		height := rapid.IntRange(1, 50).Draw(t, "height")

		dir, err := os.MkdirTemp("", "pngpanel-test-*")
		if err != nil {
			t.Fatalf("failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(dir)

		panel, err := New(
			WithDimensions(width, height),
			WithColorMode(ColorModeFullColor),
			WithOutputDir(dir),
		)
		if err != nil {
			t.Fatalf("unexpected construction error: %v", err)
		}

		// Create a random RGBA image with fully opaque pixels.
		// PNG uses non-premultiplied alpha, so transparent pixels lose color
		// information on round-trip. Use alpha=255 for lossless color fidelity.
		img := image.NewRGBA(image.Rect(0, 0, width, height))
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				r := uint8(rapid.IntRange(0, 255).Draw(t, fmt.Sprintf("r_%d_%d", x, y)))
				g := uint8(rapid.IntRange(0, 255).Draw(t, fmt.Sprintf("g_%d_%d", x, y)))
				b := uint8(rapid.IntRange(0, 255).Draw(t, fmt.Sprintf("b_%d_%d", x, y)))
				img.SetRGBA(x, y, color.RGBA{R: r, G: g, B: b, A: 255})
			}
		}

		if err := panel.DrawImage(img); err != nil {
			t.Fatalf("DrawImage failed: %v", err)
		}

		// Decode the output PNG and verify every pixel matches.
		pngPath := filepath.Join(dir, "frame_0001.png")
		f, err := os.Open(pngPath)
		if err != nil {
			t.Fatalf("failed to open output PNG: %v", err)
		}
		defer f.Close()

		decoded, err := png.Decode(f)
		if err != nil {
			t.Fatalf("failed to decode PNG: %v", err)
		}

		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				origR, origG, origB, origA := img.At(x, y).RGBA()
				decR, decG, decB, decA := decoded.At(x, y).RGBA()

				// Compare at 8-bit precision (shift 16-bit values down).
				if uint8(origR>>8) != uint8(decR>>8) ||
					uint8(origG>>8) != uint8(decG>>8) ||
					uint8(origB>>8) != uint8(decB>>8) ||
					uint8(origA>>8) != uint8(decA>>8) {
					t.Fatalf("pixel (%d,%d) mismatch: orig=(%d,%d,%d,%d) decoded=(%d,%d,%d,%d)",
						x, y,
						uint8(origR>>8), uint8(origG>>8), uint8(origB>>8), uint8(origA>>8),
						uint8(decR>>8), uint8(decG>>8), uint8(decB>>8), uint8(decA>>8))
				}
			}
		}
	})
}

func TestProperty_MonochromeOutputBinary(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		width := rapid.IntRange(1, 50).Draw(t, "width")
		height := rapid.IntRange(1, 50).Draw(t, "height")
		threshold := uint8(rapid.IntRange(0, 255).Draw(t, "threshold"))

		dir, err := os.MkdirTemp("", "pngpanel-test-*")
		if err != nil {
			t.Fatalf("failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(dir)

		panel, err := New(
			WithDimensions(width, height),
			WithColorMode(ColorModeMonochrome),
			WithThreshold(threshold),
			WithOutputDir(dir),
		)
		if err != nil {
			t.Fatalf("unexpected construction error: %v", err)
		}

		// Create a random RGBA image.
		img := image.NewRGBA(image.Rect(0, 0, width, height))
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				r := uint8(rapid.IntRange(0, 255).Draw(t, fmt.Sprintf("r_%d_%d", x, y)))
				g := uint8(rapid.IntRange(0, 255).Draw(t, fmt.Sprintf("g_%d_%d", x, y)))
				b := uint8(rapid.IntRange(0, 255).Draw(t, fmt.Sprintf("b_%d_%d", x, y)))
				img.SetRGBA(x, y, color.RGBA{R: r, G: g, B: b, A: 255})
			}
		}

		if err := panel.DrawImage(img); err != nil {
			t.Fatalf("DrawImage failed: %v", err)
		}

		// Decode the output PNG and verify all pixels are pure black or pure white.
		pngPath := filepath.Join(dir, "frame_0001.png")
		f, err := os.Open(pngPath)
		if err != nil {
			t.Fatalf("failed to open output PNG: %v", err)
		}
		defer f.Close()

		decoded, err := png.Decode(f)
		if err != nil {
			t.Fatalf("failed to decode PNG: %v", err)
		}

		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				r, g, b, _ := decoded.At(x, y).RGBA()
				r8 := uint8(r >> 8)
				g8 := uint8(g >> 8)
				b8 := uint8(b >> 8)

				isBlack := r8 == 0 && g8 == 0 && b8 == 0
				isWhite := r8 == 255 && g8 == 255 && b8 == 255
				if !isBlack && !isWhite {
					t.Fatalf("pixel (%d,%d) is neither black nor white: (%d,%d,%d)", x, y, r8, g8, b8)
				}
			}
		}
	})
}

func TestProperty_MonochromeThreshold(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		r := uint8(rapid.IntRange(0, 255).Draw(t, "r"))
		g := uint8(rapid.IntRange(0, 255).Draw(t, "g"))
		b := uint8(rapid.IntRange(0, 255).Draw(t, "b"))
		threshold := uint8(rapid.IntRange(0, 255).Draw(t, "threshold"))

		// Compute expected luminance.
		expectedLum := uint8(0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b))

		dir, err := os.MkdirTemp("", "pngpanel-test-*")
		if err != nil {
			t.Fatalf("failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(dir)

		panel, err := New(
			WithDimensions(1, 1),
			WithColorMode(ColorModeMonochrome),
			WithThreshold(threshold),
			WithOutputDir(dir),
		)
		if err != nil {
			t.Fatalf("unexpected construction error: %v", err)
		}

		// Create a 1x1 image with the test pixel.
		img := image.NewRGBA(image.Rect(0, 0, 1, 1))
		img.SetRGBA(0, 0, color.RGBA{R: r, G: g, B: b, A: 255})

		if err := panel.DrawImage(img); err != nil {
			t.Fatalf("DrawImage failed: %v", err)
		}

		// Decode the output PNG and check the pixel value.
		pngPath := filepath.Join(dir, "frame_0001.png")
		f, err := os.Open(pngPath)
		if err != nil {
			t.Fatalf("failed to open output PNG: %v", err)
		}
		defer f.Close()

		decoded, err := png.Decode(f)
		if err != nil {
			t.Fatalf("failed to decode PNG: %v", err)
		}

		outR, _, _, _ := decoded.At(0, 0).RGBA()
		outVal := uint8(outR >> 8)

		if expectedLum >= threshold {
			if outVal != 255 {
				t.Fatalf("luminance %d >= threshold %d: expected white (255), got %d", expectedLum, threshold, outVal)
			}
		} else {
			if outVal != 0 {
				t.Fatalf("luminance %d < threshold %d: expected black (0), got %d", expectedLum, threshold, outVal)
			}
		}
	})
}

func TestProperty_GrayscaleLuminance(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		r := uint8(rapid.IntRange(0, 255).Draw(t, "r"))
		g := uint8(rapid.IntRange(0, 255).Draw(t, "g"))
		b := uint8(rapid.IntRange(0, 255).Draw(t, "b"))

		// Expected grayscale value via standard luminance formula (truncated to uint8).
		expected := uint8(0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b))

		dir, err := os.MkdirTemp("", "pngpanel-test-*")
		if err != nil {
			t.Fatalf("failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(dir)

		panel, err := New(
			WithDimensions(1, 1),
			WithColorMode(ColorModeGrayscale),
			WithOutputDir(dir),
		)
		if err != nil {
			t.Fatalf("unexpected construction error: %v", err)
		}

		// Create a 1x1 image with the test pixel.
		img := image.NewRGBA(image.Rect(0, 0, 1, 1))
		img.SetRGBA(0, 0, color.RGBA{R: r, G: g, B: b, A: 255})

		if err := panel.DrawImage(img); err != nil {
			t.Fatalf("DrawImage failed: %v", err)
		}

		// Decode the output PNG and check the grayscale value.
		pngPath := filepath.Join(dir, "frame_0001.png")
		f, err := os.Open(pngPath)
		if err != nil {
			t.Fatalf("failed to open output PNG: %v", err)
		}
		defer f.Close()

		decoded, err := png.Decode(f)
		if err != nil {
			t.Fatalf("failed to decode PNG: %v", err)
		}

		// For a grayscale PNG, R=G=B=gray value.
		outR, _, _, _ := decoded.At(0, 0).RGBA()
		actual := uint8(outR >> 8)

		if actual != expected {
			t.Fatalf("grayscale(%d,%d,%d): got %d, want %d", r, g, b, actual, expected)
		}
	})
}

func TestProperty_FrameCounterSequence(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		n := rapid.IntRange(1, 20).Draw(t, "n")

		dir, err := os.MkdirTemp("", "pngpanel-test-*")
		if err != nil {
			t.Fatalf("failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(dir)

		panel, err := New(
			WithDimensions(10, 10),
			WithColorMode(ColorModeFullColor),
			WithOutputDir(dir),
		)
		if err != nil {
			t.Fatalf("unexpected construction error: %v", err)
		}

		// Draw N frames.
		img := image.NewRGBA(image.Rect(0, 0, 10, 10))
		for i := 0; i < n; i++ {
			if err := panel.DrawImage(img); err != nil {
				t.Fatalf("DrawImage call %d failed: %v", i+1, err)
			}
		}

		// Verify output files are named frame_0001.png through frame_{N:04d}.png.
		for i := 1; i <= n; i++ {
			filename := fmt.Sprintf("frame_%04d.png", i)
			path := filepath.Join(dir, filename)
			if _, err := os.Stat(path); err != nil {
				t.Fatalf("expected file %q to exist, got error: %v", filename, err)
			}
		}
	})
}

func TestProperty_FrameCounterReset(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		n := rapid.IntRange(1, 10).Draw(t, "n")

		dir, err := os.MkdirTemp("", "pngpanel-test-*")
		if err != nil {
			t.Fatalf("failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(dir)

		panel, err := New(
			WithDimensions(10, 10),
			WithColorMode(ColorModeFullColor),
			WithOutputDir(dir),
		)
		if err != nil {
			t.Fatalf("unexpected construction error: %v", err)
		}

		// Draw N frames.
		img := image.NewRGBA(image.Rect(0, 0, 10, 10))
		for i := 0; i < n; i++ {
			if err := panel.DrawImage(img); err != nil {
				t.Fatalf("DrawImage call %d failed: %v", i+1, err)
			}
		}

		// Reset counter and draw one more frame.
		panel.ResetCounter()

		if err := panel.DrawImage(img); err != nil {
			t.Fatalf("DrawImage after reset failed: %v", err)
		}

		// Verify the new file is named frame_0001.png (counter restarted).
		resetPath := filepath.Join(dir, "frame_0001.png")
		if _, err := os.Stat(resetPath); err != nil {
			t.Fatalf("expected frame_0001.png after reset, got error: %v", err)
		}
	})
}

func TestProperty_RotationDimensions(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		width := rapid.IntRange(1, 100).Draw(t, "width")
		height := rapid.IntRange(1, 100).Draw(t, "height")
		rot := rapid.SampledFrom([]int{0, 90, 180, 270}).Draw(t, "rotation")

		dir, err := os.MkdirTemp("", "pngpanel-test-*")
		if err != nil {
			t.Fatalf("failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(dir)

		panel, err := New(
			WithDimensions(width, height),
			WithColorMode(ColorModeFullColor),
			WithRotation(rot),
			WithOutputDir(dir),
		)
		if err != nil {
			t.Fatalf("unexpected construction error: %v", err)
		}

		// Bounds() should always return logical dimensions.
		bounds := panel.Bounds()
		if bounds.Dx() != width || bounds.Dy() != height {
			t.Fatalf("Bounds() = %dx%d, want %dx%d", bounds.Dx(), bounds.Dy(), width, height)
		}

		img := image.NewRGBA(image.Rect(0, 0, width, height))
		if err := panel.DrawImage(img); err != nil {
			t.Fatalf("DrawImage() error: %v", err)
		}

		// Decode and check output dimensions.
		pngPath := filepath.Join(dir, "frame_0001.png")
		f, err := os.Open(pngPath)
		if err != nil {
			t.Fatalf("failed to open output: %v", err)
		}
		defer f.Close()

		decoded, err := png.Decode(f)
		if err != nil {
			t.Fatalf("failed to decode PNG: %v", err)
		}

		var expectedW, expectedH int
		if rot == 90 || rot == 270 {
			expectedW, expectedH = height, width
		} else {
			expectedW, expectedH = width, height
		}

		if decoded.Bounds().Dx() != expectedW || decoded.Bounds().Dy() != expectedH {
			t.Fatalf("decoded %dx%d, want %dx%d (rotation=%d)",
				decoded.Bounds().Dx(), decoded.Bounds().Dy(), expectedW, expectedH, rot)
		}
	})
}

func TestProperty_Rotation180RoundTrip(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		width := rapid.IntRange(1, 30).Draw(t, "width")
		height := rapid.IntRange(1, 30).Draw(t, "height")

		// Create a random image.
		img := image.NewRGBA(image.Rect(0, 0, width, height))
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				r := uint8(rapid.IntRange(0, 255).Draw(t, fmt.Sprintf("r_%d_%d", x, y)))
				g := uint8(rapid.IntRange(0, 255).Draw(t, fmt.Sprintf("g_%d_%d", x, y)))
				b := uint8(rapid.IntRange(0, 255).Draw(t, fmt.Sprintf("b_%d_%d", x, y)))
				img.SetRGBA(x, y, color.RGBA{R: r, G: g, B: b, A: 255})
			}
		}

		// Apply 180° rotation twice using rotateImage directly (unit test of the function).
		rotated := rotateImage(img, Rotation180)
		restored := rotateImage(rotated, Rotation180)

		// Every pixel should match.
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				origR, origG, origB, _ := img.At(x, y).RGBA()
				resR, resG, resB, _ := restored.At(x, y).RGBA()
				if uint8(origR>>8) != uint8(resR>>8) || uint8(origG>>8) != uint8(resG>>8) || uint8(origB>>8) != uint8(resB>>8) {
					t.Fatalf("180° round-trip pixel (%d,%d) mismatch", x, y)
				}
			}
		}
	})
}

func TestProperty_RotationValidation(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		rot := rapid.IntRange(-360, 360).Draw(t, "rotation")
		valid := rot == 0 || rot == 90 || rot == 180 || rot == 270

		_, err := New(
			WithDimensions(10, 10),
			WithColorMode(ColorModeFullColor),
			WithRotation(rot),
			WithOutputDir(os.TempDir()),
		)

		if valid && err != nil {
			t.Fatalf("expected success for rotation %d, got error: %v", rot, err)
		}
		if !valid && err == nil {
			t.Fatalf("expected error for rotation %d, got nil", rot)
		}
	})
}
