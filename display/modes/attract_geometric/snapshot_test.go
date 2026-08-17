package attract_geometric

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/databeast/cyberhud/display/modes/testsnapshot"
	"github.com/databeast/cyberhud/display/style"
)

var snapshotOutputDir = filepath.Join("snapshots")

func categoryFromStyleName(name string) testsnapshot.DisplayCategory {
	switch {
	case strings.HasPrefix(name, "color-slow-"), strings.HasPrefix(name, "color-"):
		return testsnapshot.CategoryColor
	case strings.HasPrefix(name, "mono-slow-"):
		return testsnapshot.CategoryEink
	case strings.HasPrefix(name, "mono-"):
		return testsnapshot.CategoryMono
	case strings.HasPrefix(name, "grayscale-slow-"), strings.HasPrefix(name, "grayscale-fast-"):
		return testsnapshot.CategoryGrayscale
	default:
		return testsnapshot.CategoryColor
	}
}

func TestGeometricPNGSnapshots(t *testing.T) {
	t.Cleanup(ResetSnapshotState)

	styles := SnapshotStyles()
	if len(styles) == 0 {
		t.Fatal("geometric registry contains zero styles")
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

			frameCount := 5
			if reqs.Capability == style.MonoSlow || reqs.Capability == style.GrayscaleSlow || reqs.Capability == style.ColorSlow {
				frameCount = 1
			}

			pngPath := testsnapshot.RenderSnapshot(t,
				testsnapshot.WithMode("attract_geometric"),
				testsnapshot.WithDimensions(reqs.MinWidth, reqs.MinHeight),
				testsnapshot.WithDisplayCategory(categoryFromStyleName(s.Name())),
				testsnapshot.WithOutputDir(snapshotOutputDir),
				testsnapshot.WithBasename(s.Name()),
				testsnapshot.WithFrameCount(frameCount),
				testsnapshot.WithReset(ResetSnapshotState),
				testsnapshot.WithPreRender(func() {
					SetSnapshotDeterminism(424242, time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC), 16*time.Millisecond)
					SetPolicy(DefaultPolicy())
				}),
			)

			testsnapshot.VerifyAll(t, pngPath, reqs.MinWidth, reqs.MinHeight)
			pngPaths = append(pngPaths, pngPath)
		})
	}

	if err := testsnapshot.WriteGalleryFragmentFromPaths(snapshotOutputDir, "attract_geometric", pngPaths); err != nil {
		t.Fatalf("failed to write gallery fragment: %v", err)
	}
}
