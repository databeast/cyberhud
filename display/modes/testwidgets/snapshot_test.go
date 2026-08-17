package testwidgets

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/databeast/cyberhud/display/modes/testsnapshot"
)

var snapshotOutputDir = filepath.Join("snapshots")

func TestTestWidgetsPNGSnapshots(t *testing.T) {
	t.Cleanup(func() { nowFunc = time.Now })

	cases := []struct {
		name     string
		width    int
		height   int
		category testsnapshot.DisplayCategory
	}{
		{"mono-128x64", 128, 64, testsnapshot.CategoryMono},
		{"color-240x240", 240, 240, testsnapshot.CategoryColor},
		{"color-320x240", 320, 240, testsnapshot.CategoryColor},
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
				testsnapshot.WithMode("testwidgets"),
				testsnapshot.WithDimensions(tc.width, tc.height),
				testsnapshot.WithDisplayCategory(tc.category),
				testsnapshot.WithOutputDir(snapshotOutputDir),
				testsnapshot.WithBasename(tc.name),
				testsnapshot.WithReset(func() { nowFunc = time.Now }),
				testsnapshot.WithPreRender(func() {
					nowFunc = func() time.Time { return time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC) }
				}),
			)
			testsnapshot.VerifyAll(t, pngPath, tc.width, tc.height)
			pngPaths = append(pngPaths, pngPath)
		})
	}

	if err := testsnapshot.WriteGalleryFragmentFromPaths(snapshotOutputDir, "testwidgets", pngPaths); err != nil {
		t.Fatalf("failed to write gallery fragment: %v", err)
	}
}
