package attract_cube_test

import (
	"image/png"
	"os"
	"testing"
	"time"

	"github.com/databeast/cyberhud/display/modes/attract_cube"
	"github.com/databeast/cyberhud/display/modes/testsnapshot"
	"github.com/databeast/cyberhud/display/style"
)

func appendMatrixSnapshots(t *testing.T, outputDir string, paths []string) []string {
	t.Helper()
	for _, s := range attract_cube.SnapshotStyles() {
		req := s.Requirements()
		var category testsnapshot.DisplayCategory
		switch req.Capability {
		case style.MonoSlow:
			category = testsnapshot.CategoryEink
		case style.MonoFast:
			category = testsnapshot.CategoryMono
		case style.GrayscaleSlow:
			category = testsnapshot.CategoryGrayscaleSlow
		case style.GrayscaleFast:
			category = testsnapshot.CategoryGrayscale
		case style.ColorSlow:
			category = testsnapshot.CategoryColorSlow
		case style.ColorFast:
			category = testsnapshot.CategoryColor
		default:
			t.Fatalf("unknown capability %v", req.Capability)
		}
		t.Run(s.Name(), func(t *testing.T) {
			path := testsnapshot.RenderSnapshot(t,
				testsnapshot.WithMode("attract_cube"),
				testsnapshot.WithDimensions(req.MinWidth, req.MinHeight),
				testsnapshot.WithDisplayCategory(category),
				testsnapshot.WithOutputDir(outputDir),
				testsnapshot.WithBasename(s.Name()),
				testsnapshot.WithFrameInterval(4*time.Second),
				testsnapshot.WithPreRender(func() {
					if err := attract_cube.SetPolicy(attract_cube.DefaultPolicy()); err != nil {
						t.Fatal(err)
					}
				}),
			)
			testsnapshot.VerifyAll(t, path, req.MinWidth, req.MinHeight)
			f, err := os.Open(path)
			if err != nil {
				t.Fatal(err)
			}
			img, err := png.Decode(f)
			closeErr := f.Close()
			if err != nil {
				t.Fatal(err)
			}
			if closeErr != nil {
				t.Fatal(closeErr)
			}
			bg := img.At(0, 0)
			different := 0
			for y := 0; y < req.MinHeight; y++ {
				for x := 0; x < req.MinWidth; x++ {
					if img.At(x, y) != bg {
						different++
					}
				}
			}
			if different < 2 {
				t.Fatalf("cube has only %d non-background pixels", different)
			}
			paths = append(paths, path)
		})
	}
	return paths
}
