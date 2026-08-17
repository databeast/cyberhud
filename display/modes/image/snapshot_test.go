package image

import (
	stdimage "image"
	"image/color"
	"os"
	"path/filepath"
	"testing"

	"github.com/databeast/cyberhud/display/modes/testsnapshot"
)

var snapshotOutputDir = filepath.Join("snapshots")

func testImage() stdimage.Image {
	img := stdimage.NewRGBA(stdimage.Rect(0, 0, 96, 64))
	for y := 0; y < 64; y++ {
		for x := 0; x < 96; x++ {
			img.SetRGBA(x, y, color.RGBA{R: uint8(x * 255 / 95), G: uint8(y * 255 / 63), B: uint8((x + y) * 255 / 158), A: 0xff})
		}
	}
	return img
}

func TestImagePNGSnapshots(t *testing.T) {
	cases := []struct {
		name     string
		width    int
		height   int
		category testsnapshot.DisplayCategory
		policy   Policy
	}{
		{"mono-128x64-truncate", 128, 64, testsnapshot.CategoryMono, Policy{Fit: FitTruncate, Style: StyleDefault}},
		{"color-240x240-scale-bordered", 240, 240, testsnapshot.CategoryColor, Policy{Fit: FitScale, Style: StyleBordered}},
		{"color-320x240-stretch", 320, 240, testsnapshot.CategoryColor, Policy{Fit: FitStretch, Style: StyleDefault}},
	}

	if err := os.RemoveAll(snapshotOutputDir); err != nil {
		t.Fatalf("failed to clean snapshot output directory: %v", err)
	}
	if err := os.MkdirAll(snapshotOutputDir, 0o755); err != nil {
		t.Fatalf("failed to create snapshot output directory: %v", err)
	}

	var pngPaths []string
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			pngPath := testsnapshot.RenderSnapshot(t,
				testsnapshot.WithMode("image"),
				testsnapshot.WithDimensions(tc.width, tc.height),
				testsnapshot.WithDisplayCategory(tc.category),
				testsnapshot.WithOutputDir(snapshotOutputDir),
				testsnapshot.WithBasename(tc.name),
				testsnapshot.WithReset(Clear),
				testsnapshot.WithPreRender(func() {
					SetImage(testImage(), tc.policy)
				}),
			)
			testsnapshot.VerifyAll(t, pngPath, tc.width, tc.height)
			pngPaths = append(pngPaths, pngPath)
		})
	}

	if err := testsnapshot.WriteGalleryFragmentFromPaths(snapshotOutputDir, "image", pngPaths); err != nil {
		t.Fatalf("failed to write gallery fragment: %v", err)
	}
}
