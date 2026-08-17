package tests_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/databeast/cyberhud/display/modes/testsnapshot"
	"github.com/databeast/cyberhud/display/modes/wifi"
	"github.com/databeast/cyberhud/display/modes/wifi/source"
)

// snapshotOutputDir is the persistent directory where wifi snapshot PNGs are written.
var snapshotOutputDir = filepath.Join("..", "snapshots")

// categoryFromStyleName derives the testsnapshot.DisplayCategory from a wifi
// style's name prefix. Wifi styles follow the naming convention:
// color-*, grayscale-fast-*.
func categoryFromStyleName(name string) testsnapshot.DisplayCategory {
	switch {
	case strings.HasPrefix(name, "color-"):
		return testsnapshot.CategoryColor
	case strings.HasPrefix(name, "grayscale-fast-"):
		return testsnapshot.CategoryGrayscale
	case strings.HasPrefix(name, "mono-"):
		return testsnapshot.CategoryMono
	case strings.HasPrefix(name, "eink-"):
		return testsnapshot.CategoryEink
	default:
		return testsnapshot.CategoryColor
	}
}

// TestWifiPNGSnapshots enumerates all registered wifi styles and generates
// a snapshot PNG for each via the full production render pipeline using the
// snapshottest framework.
// Note: On non-Linux platforms, gatherWifiState returns Unavailable state,
// which exercises the disconnected/unavailable rendering paths.
//

func TestWifiPNGSnapshots(t *testing.T) {
	styles := wifi.WifiRegistryExported.Enumerate()
	if len(styles) == 0 {
		t.Fatal("wifiRegistry contains zero styles")
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

			// Derive display category from the style name.
			category := categoryFromStyleName(s.Name())

			// Build a policy that targets the current style.
			p := source.DefaultPolicy()
			p.Style = s.Name()

			// Render through the snapshottest framework.
			pngPath := testsnapshot.RenderSnapshot(t,
				testsnapshot.WithMode("wifi"),
				testsnapshot.WithDimensions(reqs.MinWidth, reqs.MinHeight),
				testsnapshot.WithDisplayCategory(category),
				testsnapshot.WithOutputDir(snapshotOutputDir),
				testsnapshot.WithBasename(s.Name()),
				testsnapshot.WithReset(func() {
					wifi.SetPolicy(source.DefaultPolicy())
					source.ResetTestWifiState()
				}),
				testsnapshot.WithPreRender(func() {
					source.SetTestWifiState(source.MockWifiData())
					wifi.SetPolicy(p)
				}),
			)

			// Verify output using the framework's verification helpers.
			testsnapshot.VerifyAll(t, pngPath, reqs.MinWidth, reqs.MinHeight)
		})
	}
}
