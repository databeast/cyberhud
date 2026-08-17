package serial_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/databeast/cyberhud/display/modes/serial/tests"
	"github.com/databeast/cyberhud/display/modes/testsnapshot"
)

// snapshotOutputDir is the persistent directory where serial snapshot PNGs are written.
var snapshotOutputDir = filepath.Join("snapshots")

// TestSerialPNGSnapshots enumerates all registered serial styles and renders
// a snapshot PNG for each through the full production pipeline using the
// snapshottest framework.
//
// Styles without explicit pixel dimensions (MinWidth/MinHeight == 0) use a
// standard 128×128 default since they are text-only styles with no dimensional
// constraints.
//

func TestSerialPNGSnapshots(t *testing.T) {
	styles := tests.SerialRegistryExported.Enumerate()
	if len(styles) == 0 {
		t.Fatal("serialRegistry contains zero styles")
	}

	// Ensure the output directory exists and is clean.
	if err := os.RemoveAll(snapshotOutputDir); err != nil {
		t.Fatalf("failed to clean snapshot output directory: %v", err)
	}
	if err := os.MkdirAll(snapshotOutputDir, 0o755); err != nil {
		t.Fatalf("failed to create snapshot output directory: %v", err)
	}

	for _, s := range styles {
		s := s
		t.Run(s.Name(), func(t *testing.T) {
			reqs := s.Requirements()

			// Use explicit dimensions from requirements, clamped to a minimum of
			// 128×128 to ensure the tier catalog can build successfully. Smaller
			// dimensions (e.g., 32×32) prevent any registered font from fitting
			// the width constraint, leaving the catalog zero-valued and causing
			// panics in resolveFaceFromCatalog.
			width := reqs.MinWidth
			height := reqs.MinHeight
			if width < 128 {
				width = 128
			}
			if height < 128 {
				height = 128
			}

			// Serial styles don't use prefix-based category naming; use CategoryColor
			// as the default since serial mode renders colored elements (ANSI, LEDs).
			category := testsnapshot.CategoryColor

			pngPath := testsnapshot.RenderSnapshot(t,
				testsnapshot.WithDimensions(width, height),
				testsnapshot.WithMode("serial"),
				testsnapshot.WithOutputDir(snapshotOutputDir),
				testsnapshot.WithDisplayCategory(category),
				testsnapshot.WithBasename(s.Name()),
				testsnapshot.WithReset(tests.ResetTestState),
				testsnapshot.WithPreRender(func() {
					tests.SeedTestState(s.Name())
				}),
			)

			testsnapshot.VerifyAll(t, pngPath, width, height)
		})
	}
}
