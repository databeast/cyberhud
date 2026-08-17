package testpattern

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/databeast/cyberhud/display/modes/testsnapshot"
)

var snapshotOutputDir = filepath.Join("snapshots")

func TestTestPatternPNGSnapshots(t *testing.T) {
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
				testsnapshot.WithMode("testpattern"),
				testsnapshot.WithDimensions(tc.width, tc.height),
				testsnapshot.WithDisplayCategory(tc.category),
				testsnapshot.WithOutputDir(snapshotOutputDir),
				testsnapshot.WithBasename(tc.name),
			)
			testsnapshot.VerifyAll(t, pngPath, tc.width, tc.height)
			pngPaths = append(pngPaths, pngPath)
		})
	}

	if err := testsnapshot.WriteGalleryFragmentFromPaths(snapshotOutputDir, "testpattern", pngPaths); err != nil {
		t.Fatalf("failed to write gallery fragment: %v", err)
	}
}
