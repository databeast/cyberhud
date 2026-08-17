package cycle

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/databeast/cyberhud/display/modes/testsnapshot"
)

var snapshotOutputDir = filepath.Join("snapshots")

func TestCyclePlaceholderPNGSnapshot(t *testing.T) {
	if err := os.RemoveAll(snapshotOutputDir); err != nil {
		t.Fatalf("failed to clean snapshot output directory: %v", err)
	}
	if err := os.MkdirAll(snapshotOutputDir, 0o755); err != nil {
		t.Fatalf("failed to create snapshot output directory: %v", err)
	}

	pngPath := testsnapshot.RenderSnapshot(t,
		testsnapshot.WithMode("cycle"),
		testsnapshot.WithDimensions(128, 64),
		testsnapshot.WithDisplayCategory(testsnapshot.CategoryMono),
		testsnapshot.WithOutputDir(snapshotOutputDir),
		testsnapshot.WithBasename("mono-128x64-placeholder"),
		testsnapshot.WithReset(func() { SetPolicy(Policy{Interval: DefaultInterval}) }),
	)
	testsnapshot.VerifyAll(t, pngPath, 128, 64)

	if err := testsnapshot.WriteGalleryFragmentFromPaths(snapshotOutputDir, "cycle", []string{pngPath}); err != nil {
		t.Fatalf("failed to write gallery fragment: %v", err)
	}
}
