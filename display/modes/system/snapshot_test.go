package system_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/databeast/cyberhud/display/modes/system"
	"github.com/databeast/cyberhud/display/modes/testsnapshot"
	"github.com/databeast/cyberhud/display/style"
)

// snapshotOutputDir is the persistent directory where system snapshot PNGs are written.
var snapshotOutputDir = filepath.Join("snapshots")

func categoryFromStyle(name string, reqs style.SurfaceRequirements) testsnapshot.DisplayCategory {
	switch {
	case strings.HasPrefix(name, "color-"):
		return testsnapshot.CategoryColor
	case strings.HasPrefix(name, "grayscale-"):
		return testsnapshot.CategoryGrayscale
	case strings.HasPrefix(name, "eink-"), strings.HasPrefix(name, "mono-slow-"):
		return testsnapshot.CategoryEink
	case strings.HasPrefix(name, "mono-"):
		return testsnapshot.CategoryMono
	}
	switch reqs.Capability {
	case style.GrayscaleSlow, style.GrayscaleFast:
		return testsnapshot.CategoryGrayscale
	case style.ColorSlow, style.ColorFast:
		return testsnapshot.CategoryColor
	case style.MonoFast:
		return testsnapshot.CategoryMono
	default:
		return testsnapshot.CategoryColor
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

// TestSystemPNGSnapshots enumerates all registered system styles and renders
// one deterministic PNG for each through the production snapshot pipeline.
func TestSystemPNGSnapshots(t *testing.T) {
	styles := system.SystemRegistryEnumerate()
	if len(styles) == 0 {
		t.Fatal("system registry contains zero styles")
	}

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
			width, height := dimensionsFromRequirements(reqs)

			pngPath := testsnapshot.RenderSnapshot(t,
				testsnapshot.WithMode("system"),
				testsnapshot.WithDimensions(width, height),
				testsnapshot.WithDisplayCategory(categoryFromStyle(s.Name(), reqs)),
				testsnapshot.WithOutputDir(snapshotOutputDir),
				testsnapshot.WithBasename(s.Name()),
				testsnapshot.WithReset(func() {
					system.SetPolicy(system.DefaultPolicy())
				}),
				testsnapshot.WithPreRender(func() {
					system.SetPolicy(system.Policy{Style: s.Name(), Font: "auto"})
				}),
			)

			testsnapshot.VerifyAll(t, pngPath, width, height)
		})
	}
}
