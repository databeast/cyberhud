package gpio

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/databeast/cyberhud/display/modes/testsnapshot"
	"github.com/databeast/cyberhud/display/style"
	gpiomgr "github.com/databeast/cyberhud/hardware/gpio"
)

// snapshotOutputDir is the persistent directory where GPIO snapshot PNGs are written.
var snapshotOutputDir = filepath.Join("snapshots")

// testGpioPolicy returns the deterministic policy used across all GPIO snapshot tests.
func testGpioPolicy() Policy {
	return Policy{
		Style:   "list",
		Color:   true,
		Font:    "auto",
		FGColor: "cyan",
	}
}

// categoryFromStyle derives the testsnapshot.DisplayCategory from style name or requirements.
func categoryFromStyle(name string, reqs style.SurfaceRequirements) testsnapshot.DisplayCategory {
	switch {
	case strings.HasPrefix(name, "color-"):
		return testsnapshot.CategoryColor
	case strings.HasPrefix(name, "eink-"):
		return testsnapshot.CategoryEink
	case strings.HasPrefix(name, "mono-slow-"):
		return testsnapshot.CategoryEink
	case strings.HasPrefix(name, "mono-"):
		return testsnapshot.CategoryMono
	case strings.HasPrefix(name, "grayscale-"):
		return testsnapshot.CategoryGrayscale
	}
	switch reqs.Capability {
	case style.GrayscaleSlow, style.GrayscaleFast:
		return testsnapshot.CategoryGrayscale
	case style.ColorSlow, style.ColorFast:
		return testsnapshot.CategoryColor
	case style.MonoFast:
		return testsnapshot.CategoryMono
	default:
		return testsnapshot.CategoryEink
	}
}

func dimensionsFromRequirements(reqs style.SurfaceRequirements) (int, int) {
	width, height := reqs.MinWidth, reqs.MinHeight
	if width == 0 {
		width = reqs.PreferredWidth
	}
	if height == 0 {
		height = reqs.PreferredHeight
	}
	if width == 0 {
		width = 240
	}
	if height == 0 {
		height = 320
	}
	return width, height
}

func categoryString(category testsnapshot.DisplayCategory) string {
	switch category {
	case testsnapshot.CategoryMono:
		return "mono"
	case testsnapshot.CategoryEink:
		return "eink"
	case testsnapshot.CategoryGrayscale:
		return "grayscale"
	default:
		return "color"
	}
}

// TestGpioPNGSnapshots exercises all registered GPIO styles that have explicit dimensions
// and verifies that each produces a valid PNG with correct dimensions and visible content
// through the full production render pipeline via the snapshottest framework.
//
// Output PNGs are written to snapshots/ with naming:
//
//	{width}x{height}_{style}_{category}_0001.png
//

func TestGpioPNGSnapshots(t *testing.T) {
	// Ensure the output directory exists and is clean.
	if err := os.RemoveAll(snapshotOutputDir); err != nil {
		t.Fatalf("failed to clean snapshot output directory: %v", err)
	}
	if err := os.MkdirAll(snapshotOutputDir, 0o755); err != nil {
		t.Fatalf("failed to create snapshot output directory: %v", err)
	}

	styles := Registry().Enumerate()

	gm := gpiomgr.New()

	for _, s := range styles {
		s := s
		reqs := s.Requirements()
		width, height := dimensionsFromRequirements(reqs)
		category := categoryFromStyle(s.Name(), reqs)

		// Derive the category string for the basename.
		categoryStr := categoryString(category)
		basename := formatBasename(width, height, s.Name(), categoryStr)

		t.Run(s.Name(), func(t *testing.T) {
			// Set a deterministic policy for this render, using the style under test.
			policy := testGpioPolicy()
			policy.Style = s.Name()

			pngPath := testsnapshot.RenderSnapshot(t,
				testsnapshot.WithDimensions(width, height),
				testsnapshot.WithMode("gpio"),
				testsnapshot.WithOutputDir(snapshotOutputDir),
				testsnapshot.WithDisplayCategory(category),
				testsnapshot.WithGPIOManager(gm),
				testsnapshot.WithBasename(basename),
				testsnapshot.WithPreRender(func() {
					SetPolicy(policy)
				}),
			)

			testsnapshot.VerifyAll(t, pngPath, width, height)
		})
	}
}

// formatBasename produces the output filename prefix:
//
//	{width}x{height}_{style}_{category}
func formatBasename(width, height int, styleName, category string) string {
	return fmt.Sprintf("%dx%d_%s_%s", width, height, styleName, category)
}
