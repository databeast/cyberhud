package testsnapshot

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// modeFromPath extracts a human-readable identifier from a PNG path for use
// in error messages. It strips the directory and the _NNNN.png suffix, leaving
// the basename (which typically matches the mode ID).
func modeFromPath(pngPath string) string {
	base := filepath.Base(pngPath)
	base = strings.TrimSuffix(base, filepath.Ext(base))
	// Strip trailing _NNNN frame number if present.
	if idx := strings.LastIndex(base, "_"); idx > 0 {
		suffix := base[idx+1:]
		allDigits := true
		for _, ch := range suffix {
			if ch < '0' || ch > '9' {
				allDigits = false
				break
			}
		}
		if allDigits && len(suffix) > 0 {
			base = base[:idx]
		}
	}
	return base
}

// decodePNG is a shared helper that opens and decodes a PNG file.
// It calls t.Fatal on any error.
func decodePNG(t *testing.T, pngPath string) image.Image {
	t.Helper()

	f, err := os.Open(pngPath)
	if err != nil {
		if os.IsNotExist(err) {
			t.Fatalf("snapshottest: PNG not found at %q", pngPath)
		}
		t.Fatalf("snapshottest: PNG decode failed for %q: %v", pngPath, err)
	}
	defer f.Close()

	img, err := png.Decode(f)
	if err != nil {
		t.Fatalf("snapshottest: PNG decode failed for %q: %v", pngPath, err)
	}
	return img
}

// VerifyPNG asserts the output file exists and is a decodable PNG.
func VerifyPNG(t *testing.T, pngPath string) {
	t.Helper()
	decodePNG(t, pngPath)
}

// VerifyDimensions asserts the decoded PNG has the expected width and height.
func VerifyDimensions(t *testing.T, pngPath string, wantW, wantH int) {
	t.Helper()

	img := decodePNG(t, pngPath)
	bounds := img.Bounds()
	gotW := bounds.Dx()
	gotH := bounds.Dy()

	if gotW != wantW || gotH != wantH {
		mode := modeFromPath(pngPath)
		t.Fatalf("snapshottest: dimensions mismatch for mode %q: want %dx%d, got %dx%d", mode, wantW, wantH, gotW, gotH)
	}
}

// VerifyNonBlank asserts the image contains at least one non-zero pixel.
func VerifyNonBlank(t *testing.T, pngPath string) {
	t.Helper()

	img := decodePNG(t, pngPath)
	bounds := img.Bounds()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			if r != 0 || g != 0 || b != 0 || a != 0 {
				return
			}
		}
	}

	mode := modeFromPath(pngPath)
	t.Fatalf("snapshottest: mode %q produced blank output (all-zero pixels)", mode)
}

// VerifyGolden compares the rendered PNG pixel-by-pixel against a golden reference.
// Fails if any pixel differs.
func VerifyGolden(t *testing.T, pngPath, goldenPath string) {
	t.Helper()

	got := decodePNG(t, pngPath)
	want := decodePNG(t, goldenPath)

	gotBounds := got.Bounds()
	wantBounds := want.Bounds()

	if gotBounds.Dx() != wantBounds.Dx() || gotBounds.Dy() != wantBounds.Dy() {
		t.Fatalf("snapshottest: golden mismatch: image dimensions differ: got %dx%d, want %dx%d",
			gotBounds.Dx(), gotBounds.Dy(), wantBounds.Dx(), wantBounds.Dy())
	}

	for y := gotBounds.Min.Y; y < gotBounds.Max.Y; y++ {
		for x := gotBounds.Min.X; x < gotBounds.Max.X; x++ {
			gotPx := got.At(x, y)
			wantPx := want.At(x+wantBounds.Min.X-gotBounds.Min.X, y+wantBounds.Min.Y-gotBounds.Min.Y)

			gr, gg, gb, ga := gotPx.RGBA()
			wr, wg, wb, wa := wantPx.RGBA()

			if gr != wr || gg != wg || gb != wb || ga != wa {
				t.Fatalf("snapshottest: golden mismatch at (%d,%d): got %v, want %v",
					x, y, color.NRGBA{
						R: uint8(gr >> 8), G: uint8(gg >> 8),
						B: uint8(gb >> 8), A: uint8(ga >> 8),
					}, color.NRGBA{
						R: uint8(wr >> 8), G: uint8(wg >> 8),
						B: uint8(wb >> 8), A: uint8(wa >> 8),
					})
			}
		}
	}
}

// VerifyAll runs VerifyPNG + VerifyDimensions + VerifyNonBlank as a convenience.
func VerifyAll(t *testing.T, pngPath string, wantW, wantH int) {
	t.Helper()
	VerifyPNG(t, pngPath)
	VerifyDimensions(t, pngPath, wantW, wantH)
	VerifyNonBlank(t, pngPath)
}
