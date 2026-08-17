package attract_shapes

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/databeast/cyberhud/display/modes/testsnapshot"
	"github.com/databeast/cyberhud/display/style"
)

// snapshotOutputDir is the persistent directory where shapes snapshot PNGs are written.
var snapshotOutputDir = filepath.Join("snapshots")

// categoryFromStyleName derives the testsnapshot.DisplayCategory from a shapes
// style's name prefix.
func categoryFromStyleName(name string) testsnapshot.DisplayCategory {
	switch {
	case strings.HasPrefix(name, "color-slow-"):
		return testsnapshot.CategoryColor
	case strings.HasPrefix(name, "color-fast-"):
		return testsnapshot.CategoryColor
	case strings.HasPrefix(name, "mono-slow-"):
		return testsnapshot.CategoryEink
	case strings.HasPrefix(name, "mono-fast-"):
		return testsnapshot.CategoryMono
	case strings.HasPrefix(name, "grayscale-slow-"):
		return testsnapshot.CategoryGrayscale
	case strings.HasPrefix(name, "grayscale-fast-"):
		return testsnapshot.CategoryGrayscale
	default:
		return testsnapshot.CategoryColor
	}
}

// TestShapesPNGSnapshots enumerates all registered shapes styles and renders each
// through the full production pipeline via the snapshottest framework.
func TestShapesPNGSnapshots(t *testing.T) {
	styles := SnapshotStyles()
	if len(styles) == 0 {
		t.Fatal("shapesRegistry contains zero styles")
	}

	if err := os.RemoveAll(snapshotOutputDir); err != nil {
		t.Fatalf("failed to clean snapshot output directory: %v", err)
	}
	if err := os.MkdirAll(snapshotOutputDir, 0o755); err != nil {
		t.Fatalf("failed to create snapshot output directory: %v", err)
	}

	var pngPaths []string

	for _, s := range styles {
		s := s
		t.Run(s.Name(), func(t *testing.T) {
			reqs := s.Requirements()
			if reqs.MinWidth == 0 || reqs.MinHeight == 0 {
				t.Skip("skipping: style has unconstrained dimensions")
			}

			category := categoryFromStyleName(s.Name())

			frameCount := 5
			if reqs.Capability == style.MonoSlow || reqs.Capability == style.GrayscaleSlow || reqs.Capability == style.ColorSlow {
				frameCount = 1
			}

			pngPath := testsnapshot.RenderSnapshot(t,
				testsnapshot.WithMode("attract_shapes"),
				testsnapshot.WithDimensions(reqs.MinWidth, reqs.MinHeight),
				testsnapshot.WithDisplayCategory(category),
				testsnapshot.WithOutputDir(snapshotOutputDir),
				testsnapshot.WithBasename(s.Name()),
				testsnapshot.WithFrameCount(frameCount),
				testsnapshot.WithReset(func() {
					ResetSnapshotState()
				}),
				testsnapshot.WithPreRender(func() {
					SetPolicy(DefaultPolicy())
				}),
			)

			testsnapshot.VerifyAll(t, pngPath, reqs.MinWidth, reqs.MinHeight)
			pngPaths = append(pngPaths, pngPath)
		})
	}

	// Write gallery fragment for this mode's snapshots.
	if err := testsnapshot.WriteGalleryFragmentFromPaths(snapshotOutputDir, "attract_shapes", pngPaths); err != nil {
		t.Fatalf("failed to write gallery fragment: %v", err)
	}
}
