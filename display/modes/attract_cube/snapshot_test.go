package attract_cube_test

import (
	"bytes"
	"fmt"
	"image/png"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/databeast/cyberhud/display/modes/attract_cube"
	"github.com/databeast/cyberhud/display/modes/testsnapshot"
)

func TestCubePNGSnapshots(t *testing.T) {
	t.Cleanup(func() {
		if err := attract_cube.SetPolicy(attract_cube.DefaultPolicy()); err != nil {
			t.Error(err)
		}
	})
	outputDir := filepath.Join("snapshots")
	var paths []string
	for _, size := range []struct{ width, height int }{
		{240, 240}, {320, 240}, {480, 320},
	} {
		for _, moment := range []struct {
			name     string
			interval time.Duration
		}{
			{"initial", time.Millisecond},
			{"2s", 2 * time.Second},
			{"4s", 4 * time.Second},
			{"8s", 8 * time.Second},
		} {
			name := fmt.Sprintf("color-fast-%dx%d-%s", size.width, size.height, moment.name)
			t.Run(name, func(t *testing.T) {
				path := testsnapshot.RenderSnapshot(t,
					testsnapshot.WithMode("attract_cube"),
					testsnapshot.WithDimensions(size.width, size.height),
					testsnapshot.WithDisplayCategory(testsnapshot.CategoryColor),
					testsnapshot.WithOutputDir(outputDir),
					testsnapshot.WithBasename(name),
					testsnapshot.WithFrameInterval(moment.interval),
					testsnapshot.WithPreRender(func() {
						if err := attract_cube.SetPolicy(attract_cube.DefaultPolicy()); err != nil {
							t.Fatal(err)
						}
					}),
				)
				testsnapshot.VerifyAll(t, path, size.width, size.height)
				f, err := os.Open(path)
				if err != nil {
					t.Fatal(err)
				}
				img, err := png.Decode(f)
				if closeErr := f.Close(); closeErr != nil {
					t.Fatal(closeErr)
				}
				if err != nil {
					t.Fatal(err)
				}
				background := img.At(0, 0)
				different := 0
				for y := 0; y < size.height; y++ {
					for x := 0; x < size.width; x++ {
						if img.At(x, y) != background {
							different++
						}
					}
				}
				if different < size.width*size.height/100 {
					t.Fatalf("snapshot contains only %d non-background pixels", different)
				}
				paths = append(paths, path)
			})
		}
	}
	for _, sample := range []struct{ name, color string }{
		{"amber", "amber"}, {"red", "red"}, {"custom", "#8844ff"},
	} {
		t.Run(sample.name, func(t *testing.T) {
			path := testsnapshot.RenderSnapshot(t,
				testsnapshot.WithMode("attract_cube"),
				testsnapshot.WithDimensions(320, 240),
				testsnapshot.WithDisplayCategory(testsnapshot.CategoryColor),
				testsnapshot.WithOutputDir(outputDir),
				testsnapshot.WithBasename("color-fast-320x240-"+sample.name),
				testsnapshot.WithFrameInterval(4*time.Second),
				testsnapshot.WithPreRender(func() {
					if err := attract_cube.SetPolicy(attract_cube.Policy{BaseColor: sample.color}); err != nil {
						t.Fatal(err)
					}
				}),
			)
			testsnapshot.VerifyAll(t, path, 320, 240)
			paths = append(paths, path)
		})
	}
	for _, sample := range []struct {
		name     string
		w, h     int
		category testsnapshot.DisplayCategory
		interval time.Duration
	}{
		{"grayscale-fast-128x64-initial", 128, 64, testsnapshot.CategoryGrayscale, time.Millisecond},
		{"grayscale-fast-128x64-4s", 128, 64, testsnapshot.CategoryGrayscale, 4 * time.Second},
		{"grayscale-fast-320x240-4s", 320, 240, testsnapshot.CategoryGrayscale, 4 * time.Second},
		{"grayscale-slow-128x64-initial", 128, 64, testsnapshot.CategoryGrayscaleSlow, time.Millisecond},
		{"grayscale-slow-128x64-4s", 128, 64, testsnapshot.CategoryGrayscaleSlow, 4 * time.Second},
		{"grayscale-slow-320x240-4s", 320, 240, testsnapshot.CategoryGrayscaleSlow, 4 * time.Second},
		{"mono-fast-8x16-4s", 8, 16, testsnapshot.CategoryMono, 4 * time.Second},
		{"mono-fast-128x32-4s", 128, 32, testsnapshot.CategoryMono, 4 * time.Second},
		{"mono-fast-128x64-initial", 128, 64, testsnapshot.CategoryMono, time.Millisecond},
		{"mono-fast-128x64-4s", 128, 64, testsnapshot.CategoryMono, 4 * time.Second},
		{"mono-slow-200x200-initial", 200, 200, testsnapshot.CategoryEink, time.Millisecond},
		{"mono-slow-200x200-4s", 200, 200, testsnapshot.CategoryEink, 4 * time.Second},
		{"mono-slow-250x122-4s", 250, 122, testsnapshot.CategoryEink, 4 * time.Second},
		{"color-slow-200x200-initial", 200, 200, testsnapshot.CategoryColorSlow, time.Millisecond},
		{"color-slow-200x200-4s", 200, 200, testsnapshot.CategoryColorSlow, 4 * time.Second},
	} {
		t.Run(sample.name, func(t *testing.T) {
			path := testsnapshot.RenderSnapshot(t,
				testsnapshot.WithMode("attract_cube"),
				testsnapshot.WithDimensions(sample.w, sample.h),
				testsnapshot.WithDisplayCategory(sample.category),
				testsnapshot.WithOutputDir(outputDir),
				testsnapshot.WithBasename(sample.name),
				testsnapshot.WithFrameInterval(sample.interval),
			)
			testsnapshot.VerifyAll(t, path, sample.w, sample.h)
			if sample.category == testsnapshot.CategoryMono || sample.category == testsnapshot.CategoryEink {
				f, err := os.Open(path)
				if err != nil {
					t.Fatal(err)
				}
				img, err := png.Decode(f)
				if closeErr := f.Close(); closeErr != nil {
					t.Fatal(closeErr)
				}
				if err != nil {
					t.Fatal(err)
				}
				white := 0
				for y := 0; y < sample.h; y++ {
					for x := 0; x < sample.w; x++ {
						r, g, b, _ := img.At(x, y).RGBA()
						if r != g || g != b || (r != 0 && r != 0xffff) {
							t.Fatalf("snapshot contains non-binary pixel at (%d,%d)", x, y)
						}
						if r == 0xffff {
							white++
						}
					}
				}
				if white == 0 || white >= sample.w*sample.h/2 {
					t.Fatalf("snapshot has %d white pixels, not a wireframe", white)
				}
			}
			paths = append(paths, path)
		})
	}
	slowInitial, err := os.ReadFile(filepath.Join(outputDir, "grayscale-slow-128x64-initial_0001.png"))
	if err != nil {
		t.Fatal(err)
	}
	slowLater, err := os.ReadFile(filepath.Join(outputDir, "grayscale-slow-128x64-4s_0001.png"))
	if err != nil {
		t.Fatal(err)
	}
	fastLater, err := os.ReadFile(filepath.Join(outputDir, "grayscale-fast-128x64-4s_0001.png"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(slowInitial, slowLater) || bytes.Equal(slowLater, fastLater) {
		t.Fatal("grayscale snapshots do not reflect fast and slow refresh behavior")
	}
	for _, names := range [][2]string{
		{"mono-fast-128x64-initial", "mono-fast-128x64-4s"},
		{"mono-slow-200x200-initial", "mono-slow-200x200-4s"},
		{"color-slow-200x200-initial", "color-slow-200x200-4s"},
	} {
		first, err := os.ReadFile(filepath.Join(outputDir, names[0]+"_0001.png"))
		if err != nil {
			t.Fatal(err)
		}
		second, err := os.ReadFile(filepath.Join(outputDir, names[1]+"_0001.png"))
		if err != nil {
			t.Fatal(err)
		}
		if bytes.Equal(first, second) != (names[0] != "mono-fast-128x64-initial") {
			t.Fatalf("snapshot timing incorrect for %s", names[0])
		}
	}
	paths = appendMatrixSnapshots(t, outputDir, paths)
	if err := testsnapshot.WriteGalleryFragmentFromPaths(outputDir, "attract_cube", paths); err != nil {
		t.Fatal(err)
	}
}
