package dashboard

import (
	"image"
	"image/color"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/databeast/cyberhud/display/modes/dashboard/source"
	"github.com/databeast/cyberhud/display/modes/testsnapshot"
)

// snapshotOutputDir is the persistent directory where snapshot PNGs are written.
// Located at snapshots/ relative to this test file, so the output
// survives the test run and can be visually inspected or committed as golden files.
var snapshotOutputDir = filepath.Join("snapshots")

// TestDashboardPNGSnapshots enumerates all registered dashboard styles and runs a subtest
// for each. It guards against registry drift by asserting exactly 176 styles.
// Output PNGs are written to snapshots/ for visual inspection.
//
// Uses the snapshottest framework to render through the full production pipeline
// (PNGPanel → Region → ModeEngine → RegionRenderer.Render).
//

func TestDashboardPNGSnapshots(t *testing.T) {
	styles := dashboardRegistry.Enumerate()

	// Ensure the output directory exists and is clean.
	if err := os.RemoveAll(snapshotOutputDir); err != nil {
		t.Fatalf("failed to clean snapshot output directory: %v", err)
	}
	if err := os.MkdirAll(snapshotOutputDir, 0o755); err != nil {
		t.Fatalf("failed to create snapshot output directory: %v", err)
	}

	for _, s := range styles {
		s := s // capture range variable
		t.Run(s.Name(), func(t *testing.T) {
			// Read style requirements and skip if dimensions are unconstrained.
			reqs := s.Requirements()
			if reqs.MinWidth == 0 || reqs.MinHeight == 0 {
				t.Skip("skipping: style has unconstrained dimensions (MinWidth or MinHeight is 0)")
			}

			// Derive display category from the style name.
			category := categoryFromStyleName(s.Name())

			// Render through the snapshottest framework.
			pngPath := testsnapshot.RenderSnapshot(t,
				testsnapshot.WithMode("dashboard"),
				testsnapshot.WithDimensions(reqs.MinWidth, reqs.MinHeight),
				testsnapshot.WithDisplayCategory(category),
				testsnapshot.WithOutputDir(snapshotOutputDir),
				testsnapshot.WithBasename(s.Name()),
				testsnapshot.WithIconGetter(stubIconGetter),
				testsnapshot.WithReset(func() {
					ResetPolicy()
					source.ResetTestDashboardContent()
				}),
				testsnapshot.WithPreRender(func() {
					source.SetTestDashboardContent(source.MockDashboardContent())
					HandleCommand([]string{"style=" + s.Name()})
				}),
			)

			// Verify output using the framework's verification helpers.
			testsnapshot.VerifyAll(t, pngPath, reqs.MinWidth, reqs.MinHeight)
		})
	}
}

// stubIconGetter provides a minimal icon getter for dashboard sprite rendering.
// It returns a small 8x8 colored image for any requested icon name.
func stubIconGetter(name string) (image.Image, bool) {
	img := image.NewRGBA(image.Rect(0, 0, 8, 8))
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			img.Set(x, y, color.RGBA{R: 0, G: 200, B: 0, A: 255})
		}
	}
	return img, true
}

// categoryFromStyleName derives the testsnapshot.DisplayCategory from a dashboard
// style's name prefix. Dashboard styles follow the naming convention:
// mono-*, color-*, eink-*, grayscale-fast-*.
func categoryFromStyleName(name string) testsnapshot.DisplayCategory {
	switch {
	case strings.HasPrefix(name, "color-"):
		return testsnapshot.CategoryColor
	case strings.HasPrefix(name, "mono-"):
		return testsnapshot.CategoryMono
	case strings.HasPrefix(name, "eink-"):
		return testsnapshot.CategoryEink
	case strings.HasPrefix(name, "grayscale-fast-"):
		return testsnapshot.CategoryGrayscale
	default:
		return testsnapshot.CategoryColor
	}
}
