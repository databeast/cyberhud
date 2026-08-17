package systemd_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/databeast/cyberhud/display/modes/systemd"
	"github.com/databeast/cyberhud/display/modes/testsnapshot"
	"github.com/databeast/cyberhud/display/style"
)

// snapshotOutputDir is the persistent directory where systemd snapshot PNGs are written.
var snapshotOutputDir = filepath.Join("snapshots")

// categoryFromCapability maps a style's declared hardware Capability to the
// appropriate testsnapshot.DisplayCategory. Derived from the style's own
// SurfaceRequirements rather than name heuristics, so new styles registered
// in the registry are handled automatically.
func categoryFromCapability(cap style.Capability) testsnapshot.DisplayCategory {
	switch cap {
	case style.MonoSlow:
		return testsnapshot.CategoryEink
	case style.MonoFast:
		return testsnapshot.CategoryMono
	case style.GrayscaleSlow, style.GrayscaleFast:
		return testsnapshot.CategoryGrayscale
	case style.ColorSlow, style.ColorFast:
		return testsnapshot.CategoryColor
	default:
		return testsnapshot.CategoryColor
	}
}

// TestSystemdPNGSnapshots enumerates all registered systemd styles and generates
// a snapshot PNG for each via the full production render pipeline using the
// snapshottest framework.
//
// Output PNGs are written to snapshots/ for visual inspection.
//

func TestSystemdPNGSnapshots(t *testing.T) {
	styles := systemd.SystemdRegistryExported.Enumerate()
	if len(styles) == 0 {
		t.Fatal("systemdRegistry contains zero styles")
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
			if reqs.MinWidth == 0 || reqs.MinHeight == 0 {
				t.Skip("skipping: style has unconstrained dimensions")
			}

			// Derive display category from the style's declared capability.
			category := categoryFromCapability(reqs.Capability)

			// Build a policy that targets the current style.
			p := systemd.DefaultPolicy()
			p.Style = s.Name()

			// Render through the snapshottest framework.
			pngPath := testsnapshot.RenderSnapshot(t,
				testsnapshot.WithMode("systemd"),
				testsnapshot.WithDimensions(reqs.MinWidth, reqs.MinHeight),
				testsnapshot.WithDisplayCategory(category),
				testsnapshot.WithOutputDir(snapshotOutputDir),
				testsnapshot.WithBasename(s.Name()),
				testsnapshot.WithReset(func() {
					systemd.SetPolicy(systemd.DefaultPolicy())
				}),
				testsnapshot.WithPreRender(func() {
					systemd.SetPolicy(p)
				}),
			)

			// Verify output using the framework's verification helpers.
			testsnapshot.VerifyAll(t, pngPath, reqs.MinWidth, reqs.MinHeight)
		})
	}
}
