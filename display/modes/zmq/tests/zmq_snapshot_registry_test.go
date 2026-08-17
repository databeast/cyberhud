package tests_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/databeast/cyberhud/display/modes/testsnapshot"
	"github.com/databeast/cyberhud/display/modes/zmq"
	"github.com/databeast/cyberhud/display/modes/zmq/content"
)

// categoryFromStyleName derives the testsnapshot.DisplayCategory from a ZMQ
// style's name prefix. ZMQ styles follow the naming convention:
// color-*, mono-*, eink-*, grayscale-fast-*.
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

// TestZMQPNGSnapshots enumerates all registered ZMQ styles and generates
// a snapshot PNG for each through the full production pipeline using the
// snapshottest framework.
//

func TestZMQPNGSnapshots(t *testing.T) {
	styles := zmq.ZMQRegistryExported.Enumerate()
	if len(styles) == 0 {
		t.Fatal("zmqRegistry contains zero styles")
	}

	// Ensure the output directory exists and is clean.
	outputDir := filepath.Join("testdata", "snapshots")
	if err := os.RemoveAll(outputDir); err != nil {
		t.Fatalf("failed to clean snapshot output directory: %v", err)
	}
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		t.Fatalf("failed to create snapshot output directory: %v", err)
	}

	for _, s := range styles {
		s := s
		t.Run(s.Name(), func(t *testing.T) {
			reqs := s.Requirements()
			if reqs.MinWidth == 0 || reqs.MinHeight == 0 {
				t.Skip("skipping: style has unconstrained dimensions")
			}

			// Derive display category from the style name prefix.
			category := categoryFromStyleName(s.Name())

			// Build a policy targeting the current style.
			p := content.DefaultPolicy()
			p.Style = s.Name()

			pngPath := testsnapshot.RenderSnapshot(t,
				testsnapshot.WithMode("zmq"),
				testsnapshot.WithDimensions(reqs.MinWidth, reqs.MinHeight),
				testsnapshot.WithDisplayCategory(category),
				testsnapshot.WithOutputDir(outputDir),
				testsnapshot.WithBasename(s.Name()),
				testsnapshot.WithReset(zmq.ResetTestState),
				testsnapshot.WithPreRender(func() {
					content.SetPolicy(p)
					zmq.SeedTestMessages()
				}),
			)

			testsnapshot.VerifyAll(t, pngPath, reqs.MinWidth, reqs.MinHeight)
		})
	}
}
