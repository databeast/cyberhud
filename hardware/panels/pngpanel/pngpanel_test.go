package pngpanel

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// ---------------------------------------------------------------------------
// ---------------------------------------------------------------------------

// TestNew_HardwareResolutions verifies that all registered hardware resolutions
// are accepted by the constructor without error.
func TestNew_HardwareResolutions(t *testing.T) {
	resolutions := []struct {
		width, height int
	}{
		{128, 32},
		{128, 64},
		{160, 80},
		{240, 135},
		{240, 240},
		{320, 240},
		{480, 320},
		{800, 480},
	}

	tmpDir := t.TempDir()
	for _, res := range resolutions {
		t.Run(fmt.Sprintf("%dx%d", res.width, res.height), func(t *testing.T) {
			p, err := New(
				WithDimensions(res.width, res.height),
				WithOutputDir(tmpDir),
			)
			if err != nil {
				t.Fatalf("New(%dx%d) returned unexpected error: %v", res.width, res.height, err)
			}
			bounds := p.Bounds()
			if bounds.Dx() != res.width || bounds.Dy() != res.height {
				t.Fatalf("Bounds() = %dx%d, want %dx%d", bounds.Dx(), bounds.Dy(), res.width, res.height)
			}
		})
	}
}

// TestNew_DefaultColorMode verifies that when no color mode is specified,
// the panel defaults to full-color mode.
func TestNew_DefaultColorMode(t *testing.T) {
	tmpDir := t.TempDir()
	p, err := New(
		WithDimensions(10, 10),
		WithOutputDir(tmpDir),
	)
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}

	// Full-color mode means PreferEventRefresh should be false.
	hints := p.TextHints()
	if hints.PreferEventRefresh {
		t.Fatal("default color mode should be full-color (PreferEventRefresh=false), got true")
	}
}

// TestNew_DefaultThreshold verifies that when no threshold is specified,
// the default threshold of 128 is used for monochrome conversion.
func TestNew_DefaultThreshold(t *testing.T) {
	tmpDir := t.TempDir()
	p, err := New(
		WithDimensions(2, 1),
		WithColorMode(ColorModeMonochrome),
		WithOutputDir(tmpDir),
	)
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}

	// Create an image with two pixels:
	// Pixel at (0,0): R=G=B=100 → luminance = floor(100*0.299+100*0.587+100*0.114) = 100, < 128 → black
	// Pixel at (1,0): R=G=B=200 → luminance = floor(200*0.299+200*0.587+200*0.114) = 200, >= 128 → white
	img := image.NewRGBA(image.Rect(0, 0, 2, 1))
	img.Set(0, 0, color.RGBA{R: 100, G: 100, B: 100, A: 255})
	img.Set(1, 0, color.RGBA{R: 200, G: 200, B: 200, A: 255})

	err = p.DrawImage(img)
	if err != nil {
		t.Fatalf("DrawImage() unexpected error: %v", err)
	}

	// Decode the output PNG and verify pixels.
	outPath := filepath.Join(tmpDir, "frame_0001.png")
	f, err := os.Open(outPath)
	if err != nil {
		t.Fatalf("failed to open output file: %v", err)
	}
	defer f.Close()

	decoded, err := png.Decode(f)
	if err != nil {
		t.Fatalf("failed to decode PNG: %v", err)
	}

	// Pixel (0,0) should be black (luminance 100 < threshold 128).
	r0, g0, b0, _ := decoded.At(0, 0).RGBA()
	if r0 != 0 || g0 != 0 || b0 != 0 {
		t.Errorf("pixel (0,0) should be black, got R=%d G=%d B=%d", r0>>8, g0>>8, b0>>8)
	}

	// Pixel (1,0) should be white (luminance 200 >= threshold 128).
	r1, g1, b1, _ := decoded.At(1, 0).RGBA()
	if r1>>8 != 255 || g1>>8 != 255 || b1>>8 != 255 {
		t.Errorf("pixel (1,0) should be white, got R=%d G=%d B=%d", r1>>8, g1>>8, b1>>8)
	}
}

// TestDrawImage_NilImage verifies that passing nil to DrawImage returns an error.
func TestDrawImage_NilImage(t *testing.T) {
	tmpDir := t.TempDir()
	p, err := New(
		WithDimensions(10, 10),
		WithOutputDir(tmpDir),
	)
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}

	err = p.DrawImage(nil)
	if err == nil {
		t.Fatal("DrawImage(nil) should return an error")
	}
	if !strings.Contains(err.Error(), "nil") {
		t.Fatalf("error should mention nil, got: %v", err)
	}
}

// TestNew_EmptyOutputPath verifies that constructing with an empty output
// directory returns a construction error.
func TestNew_EmptyOutputPath(t *testing.T) {
	_, err := New(
		WithDimensions(10, 10),
		WithOutputDir(""),
	)
	if err == nil {
		t.Fatal("New() with empty outputDir should return an error")
	}
	if !strings.Contains(err.Error(), "output path") {
		t.Fatalf("error should mention output path, got: %v", err)
	}
}

// ---------------------------------------------------------------------------
// ---------------------------------------------------------------------------

// TestDrawImage_CreatesDirectory verifies that DrawImage creates the output
// directory (including intermediate directories) when it does not exist.
func TestDrawImage_CreatesDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	nestedDir := filepath.Join(tmpDir, "a", "b", "c")

	p, err := New(
		WithDimensions(4, 4),
		WithOutputDir(nestedDir),
	)
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}

	img := image.NewRGBA(image.Rect(0, 0, 4, 4))
	if err := p.DrawImage(img); err != nil {
		t.Fatalf("DrawImage() unexpected error: %v", err)
	}

	// Verify directory was created.
	info, err := os.Stat(nestedDir)
	if err != nil {
		t.Fatalf("output directory was not created: %v", err)
	}
	if !info.IsDir() {
		t.Fatal("output path is not a directory")
	}
}

// TestDrawImage_AtomicWrite verifies that after DrawImage completes,
// no temporary files remain in the output directory.
func TestDrawImage_AtomicWrite(t *testing.T) {
	tmpDir := t.TempDir()

	p, err := New(
		WithDimensions(4, 4),
		WithOutputDir(tmpDir),
	)
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}

	img := image.NewRGBA(image.Rect(0, 0, 4, 4))
	if err := p.DrawImage(img); err != nil {
		t.Fatalf("DrawImage() unexpected error: %v", err)
	}

	// Check no .tmp files remain.
	entries, err := os.ReadDir(tmpDir)
	if err != nil {
		t.Fatalf("failed to read directory: %v", err)
	}
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".tmp") {
			t.Fatalf("temp file left behind: %s", entry.Name())
		}
	}
}

// TestDrawImage_CounterOverflow verifies that when the counter exceeds 9999,
// the filename is formatted without zero-padding.
func TestDrawImage_CounterOverflow(t *testing.T) {
	tmpDir := t.TempDir()

	p, err := New(
		WithDimensions(4, 4),
		WithOutputDir(tmpDir),
	)
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}

	// Set counter to 9999 so next write will be frame 10000.
	p.counter = 9999

	img := image.NewRGBA(image.Rect(0, 0, 4, 4))
	if err := p.DrawImage(img); err != nil {
		t.Fatalf("DrawImage() unexpected error: %v", err)
	}

	expected := filepath.Join(tmpDir, "frame_10000.png")
	if _, err := os.Stat(expected); err != nil {
		t.Fatalf("expected file %q not found: %v", expected, err)
	}
}

// TestDrawImage_OverwriteAfterReset verifies that after a reset, the next
// frame overwrites the existing file at frame_0001.png.
func TestDrawImage_OverwriteAfterReset(t *testing.T) {
	tmpDir := t.TempDir()

	p, err := New(
		WithDimensions(4, 4),
		WithOutputDir(tmpDir),
	)
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}

	img := image.NewRGBA(image.Rect(0, 0, 4, 4))
	// Draw first frame.
	if err := p.DrawImage(img); err != nil {
		t.Fatalf("DrawImage() unexpected error: %v", err)
	}

	firstPath := filepath.Join(tmpDir, "frame_0001.png")
	info1, err := os.Stat(firstPath)
	if err != nil {
		t.Fatalf("first frame file not found: %v", err)
	}

	// Reset and draw again (should overwrite frame_0001.png).
	p.ResetCounter()

	// Draw with a different pixel color so we can verify the file changed.
	img2 := image.NewRGBA(image.Rect(0, 0, 4, 4))
	draw.Draw(img2, img2.Bounds(), &image.Uniform{color.RGBA{R: 255, A: 255}}, image.Point{}, draw.Src)

	if err := p.DrawImage(img2); err != nil {
		t.Fatalf("DrawImage() after reset unexpected error: %v", err)
	}

	info2, err := os.Stat(firstPath)
	if err != nil {
		t.Fatalf("frame_0001.png not found after reset: %v", err)
	}

	// Verify the file was overwritten (mod time should be >= first write).
	if info2.ModTime().Before(info1.ModTime()) {
		t.Fatal("file was not overwritten after reset")
	}

	// Verify it's a valid PNG.
	f, err := os.Open(firstPath)
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}
	defer f.Close()

	_, err = png.Decode(f)
	if err != nil {
		t.Fatalf("overwritten file is not a valid PNG: %v", err)
	}
}

// TestDrawImage_AllColorModes verifies that all three color modes produce
// valid PNG files when drawing the same image.
func TestDrawImage_AllColorModes(t *testing.T) {
	modes := []struct {
		name string
		mode ColorMode
	}{
		{"full-color", ColorModeFullColor},
		{"grayscale", ColorModeGrayscale},
		{"monochrome", ColorModeMonochrome},
	}

	for _, m := range modes {
		t.Run(m.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			p, err := New(
				WithDimensions(10, 10),
				WithColorMode(m.mode),
				WithOutputDir(tmpDir),
			)
			if err != nil {
				t.Fatalf("New() unexpected error: %v", err)
			}

			// Create a colorful image.
			img := image.NewRGBA(image.Rect(0, 0, 10, 10))
			for y := 0; y < 10; y++ {
				for x := 0; x < 10; x++ {
					img.Set(x, y, color.RGBA{
						R: uint8(x * 25),
						G: uint8(y * 25),
						B: uint8((x + y) * 12),
						A: 255,
					})
				}
			}

			if err := p.DrawImage(img); err != nil {
				t.Fatalf("DrawImage() unexpected error: %v", err)
			}

			// Verify output is a valid PNG.
			outPath := filepath.Join(tmpDir, "frame_0001.png")
			f, err := os.Open(outPath)
			if err != nil {
				t.Fatalf("failed to open output file: %v", err)
			}
			defer f.Close()

			decoded, err := png.Decode(f)
			if err != nil {
				t.Fatalf("output is not a valid PNG: %v", err)
			}

			if decoded.Bounds().Dx() != 10 || decoded.Bounds().Dy() != 10 {
				t.Fatalf("decoded dimensions %dx%d, want 10x10",
					decoded.Bounds().Dx(), decoded.Bounds().Dy())
			}
		})
	}
}

// ---------------------------------------------------------------------------
// ---------------------------------------------------------------------------

// TestMonochrome_BoundaryWhite verifies that an all-white image with
// threshold 0 produces all-white output. Any luminance >= 0 is white,
// so all pixels should be white.
func TestMonochrome_BoundaryWhite(t *testing.T) {
	tmpDir := t.TempDir()

	p, err := New(
		WithDimensions(10, 10),
		WithColorMode(ColorModeMonochrome),
		WithThreshold(0),
		WithOutputDir(tmpDir),
	)
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}

	// All-white image.
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

	if err := p.DrawImage(img); err != nil {
		t.Fatalf("DrawImage() unexpected error: %v", err)
	}

	// Decode and verify all pixels are white.
	outPath := filepath.Join(tmpDir, "frame_0001.png")
	f, err := os.Open(outPath)
	if err != nil {
		t.Fatalf("failed to open output file: %v", err)
	}
	defer f.Close()

	decoded, err := png.Decode(f)
	if err != nil {
		t.Fatalf("failed to decode PNG: %v", err)
	}

	bounds := decoded.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := decoded.At(x, y).RGBA()
			if r>>8 != 255 || g>>8 != 255 || b>>8 != 255 {
				t.Fatalf("pixel (%d,%d) should be white, got R=%d G=%d B=%d",
					x, y, r>>8, g>>8, b>>8)
			}
		}
	}
}

// TestMonochrome_BoundaryBlack verifies that an all-black image with
// threshold 255 produces all-black output. Luminance of 0 < 255,
// so all pixels should be black.
func TestMonochrome_BoundaryBlack(t *testing.T) {
	tmpDir := t.TempDir()

	p, err := New(
		WithDimensions(10, 10),
		WithColorMode(ColorModeMonochrome),
		WithThreshold(255),
		WithOutputDir(tmpDir),
	)
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}

	// All-black image.
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.Black}, image.Point{}, draw.Src)

	if err := p.DrawImage(img); err != nil {
		t.Fatalf("DrawImage() unexpected error: %v", err)
	}

	// Decode and verify all pixels are black.
	outPath := filepath.Join(tmpDir, "frame_0001.png")
	f, err := os.Open(outPath)
	if err != nil {
		t.Fatalf("failed to open output file: %v", err)
	}
	defer f.Close()

	decoded, err := png.Decode(f)
	if err != nil {
		t.Fatalf("failed to decode PNG: %v", err)
	}

	bounds := decoded.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := decoded.At(x, y).RGBA()
			if r != 0 || g != 0 || b != 0 {
				t.Fatalf("pixel (%d,%d) should be black, got R=%d G=%d B=%d",
					x, y, r>>8, g>>8, b>>8)
			}
		}
	}
}

// ---------------------------------------------------------------------------
// Rotation unit tests
// ---------------------------------------------------------------------------

// TestNew_DefaultRotation verifies that when no rotation is specified,
// the default is 0 degrees (no rotation).
func TestNew_DefaultRotation(t *testing.T) {
	tmpDir := t.TempDir()
	p, err := New(
		WithDimensions(10, 8),
		WithOutputDir(tmpDir),
	)
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}

	// Draw an image and verify output dimensions match input (no rotation).
	img := image.NewRGBA(image.Rect(0, 0, 10, 8))
	if err := p.DrawImage(img); err != nil {
		t.Fatalf("DrawImage() error: %v", err)
	}

	f, err := os.Open(filepath.Join(tmpDir, "frame_0001.png"))
	if err != nil {
		t.Fatalf("failed to open output: %v", err)
	}
	defer f.Close()

	decoded, err := png.Decode(f)
	if err != nil {
		t.Fatalf("failed to decode PNG: %v", err)
	}

	if decoded.Bounds().Dx() != 10 || decoded.Bounds().Dy() != 8 {
		t.Fatalf("default rotation: output %dx%d, want 10x8", decoded.Bounds().Dx(), decoded.Bounds().Dy())
	}
}

// TestDrawImage_Rotation90 verifies 90° clockwise rotation produces correct output.
func TestDrawImage_Rotation90(t *testing.T) {
	tmpDir := t.TempDir()
	p, err := New(
		WithDimensions(4, 2),
		WithColorMode(ColorModeFullColor),
		WithRotation(90),
		WithOutputDir(tmpDir),
	)
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}

	// Create a 4x2 image with a known pattern.
	// Row 0: red pixels, Row 1: blue pixels
	img := image.NewRGBA(image.Rect(0, 0, 4, 2))
	for x := 0; x < 4; x++ {
		img.SetRGBA(x, 0, color.RGBA{R: 255, A: 255})
		img.SetRGBA(x, 1, color.RGBA{B: 255, A: 255})
	}

	if err := p.DrawImage(img); err != nil {
		t.Fatalf("DrawImage() error: %v", err)
	}

	f, err := os.Open(filepath.Join(tmpDir, "frame_0001.png"))
	if err != nil {
		t.Fatalf("failed to open output: %v", err)
	}
	defer f.Close()

	decoded, err := png.Decode(f)
	if err != nil {
		t.Fatalf("failed to decode PNG: %v", err)
	}

	// 90° clockwise on 4x2 → 2x4 output
	// Original (x,y) → rotated (h-1-y, x) where h=2
	// So (0,0)→(1,0), (1,0)→(1,1), (2,0)→(1,2), (3,0)→(1,3) [red → right column]
	//    (0,1)→(0,0), (1,1)→(0,1), (2,1)→(0,2), (3,1)→(0,3) [blue → left column]
	if decoded.Bounds().Dx() != 2 || decoded.Bounds().Dy() != 4 {
		t.Fatalf("output dims = %dx%d, want 2x4", decoded.Bounds().Dx(), decoded.Bounds().Dy())
	}

	// Left column should be blue, right column should be red.
	for y := 0; y < 4; y++ {
		r, _, b, _ := decoded.At(0, y).RGBA()
		if r>>8 != 0 || b>>8 != 255 {
			t.Fatalf("left col pixel (0,%d): expected blue, got R=%d B=%d", y, r>>8, b>>8)
		}
		r, _, b, _ = decoded.At(1, y).RGBA()
		if r>>8 != 255 || b>>8 != 0 {
			t.Fatalf("right col pixel (1,%d): expected red, got R=%d B=%d", y, r>>8, b>>8)
		}
	}
}

// TestNew_InvalidRotation verifies invalid rotation values are rejected.
func TestNew_InvalidRotation(t *testing.T) {
	invalid := []int{45, 135, 360, -90, 1, 91}
	for _, rot := range invalid {
		t.Run(fmt.Sprintf("%d", rot), func(t *testing.T) {
			_, err := New(
				WithDimensions(10, 10),
				WithRotation(rot),
				WithOutputDir(t.TempDir()),
			)
			if err == nil {
				t.Fatalf("expected error for rotation %d", rot)
			}
		})
	}
}
