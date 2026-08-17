package attract_matrix

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/databeast/cyberhud/display/modes/testsnapshot"
	"github.com/databeast/cyberhud/display/style"
)

// snapshotOutputDir is the persistent directory where matrix snapshot PNGs are written.
var snapshotOutputDir = filepath.Join("snapshots")

// categoryFromStyleName derives the testsnapshot.DisplayCategory from a matrix
// style's name prefix.
func categoryFromStyleName(name string) testsnapshot.DisplayCategory {
	switch {
	case strings.HasPrefix(name, "color-slow-"):
		return testsnapshot.CategoryColor
	case strings.HasPrefix(name, "color-"):
		return testsnapshot.CategoryColor
	case strings.HasPrefix(name, "mono-slow-"):
		return testsnapshot.CategoryEink
	case strings.HasPrefix(name, "mono-"):
		return testsnapshot.CategoryMono
	case strings.HasPrefix(name, "eink-"):
		return testsnapshot.CategoryEink
	case strings.HasPrefix(name, "grayscale-slow-"):
		return testsnapshot.CategoryGrayscale
	case strings.HasPrefix(name, "grayscale-fast-"):
		return testsnapshot.CategoryGrayscale
	default:
		return testsnapshot.CategoryColor
	}
}

// TestMatrixPNGSnapshots enumerates all registered matrix styles and renders each
// through the full production pipeline via the snapshottest framework.
func TestMatrixPNGSnapshots(t *testing.T) {
	styles := SnapshotStyles()
	if len(styles) == 0 {
		t.Fatal("matrixRegistry contains zero styles")
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
			if reqs.MinWidth == 0 || reqs.MinHeight == 0 {
				t.Skip("skipping: style has unconstrained dimensions")
			}

			category := categoryFromStyleName(s.Name())

			frameCount := 7
			if reqs.Capability == style.MonoSlow || reqs.Capability == style.GrayscaleSlow || reqs.Capability == style.ColorSlow {
				frameCount = 1
			}

			pngPath := testsnapshot.RenderSnapshot(t,
				testsnapshot.WithMode("attract_matrix"),
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
					// Position the cycle at peak (midpoint of 45s cycle) so the
					// "CYBERHUD" splash is visible in snapshots.
					SetCycleStartForTest(time.Now().Add(-22500 * time.Millisecond))
				}),
			)

			testsnapshot.VerifyAll(t, pngPath, reqs.MinWidth, reqs.MinHeight)
		})
	}
}

// TestMatrixShowcaseSnapshots produces a comprehensive matrix of PNG images
// combining policy variations with panel resolutions.
func TestMatrixShowcaseSnapshots(t *testing.T) {
	outputDir := filepath.Join("snapshots", "showcase")

	if err := os.RemoveAll(outputDir); err != nil {
		t.Fatalf("failed to clean showcase output directory: %v", err)
	}
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		t.Fatalf("failed to create showcase output directory: %v", err)
	}

	type policyConfig struct {
		name   string
		policy Policy
	}

	policies := []policyConfig{
		{"default", Policy{MinSpeed: 3, MaxSpeed: 12, TrailLength: 16, Density: 1.0, ShowBackground: false}},
		{"sparse", Policy{MinSpeed: 3, MaxSpeed: 12, TrailLength: 8, Density: 0.3, ShowBackground: false}},
		{"background", Policy{MinSpeed: 3, MaxSpeed: 12, TrailLength: 16, Density: 1.0, ShowBackground: true}},
		{"short-trails", Policy{MinSpeed: 3, MaxSpeed: 12, TrailLength: 4, Density: 1.0, ShowBackground: false}},
	}

	type panelSpec struct {
		name     string
		width    int
		height   int
		category testsnapshot.DisplayCategory
	}

	panels := []panelSpec{
		{"240x135_color", 240, 135, testsnapshot.CategoryColor},
		{"320x240_color", 320, 240, testsnapshot.CategoryColor},
		{"128x64_mono", 128, 64, testsnapshot.CategoryMono},
		{"296x128_eink", 296, 128, testsnapshot.CategoryEink},
	}

	for _, pc := range policies {
		for _, pan := range panels {
			pc := pc
			pan := pan
			name := pan.name + "_" + pc.name
			t.Run(name, func(t *testing.T) {
				frameCount := 7
				if pan.category == testsnapshot.CategoryEink {
					frameCount = 1
				}

				pngPath := testsnapshot.RenderSnapshot(t,
					testsnapshot.WithMode("attract_matrix"),
					testsnapshot.WithDimensions(pan.width, pan.height),
					testsnapshot.WithDisplayCategory(pan.category),
					testsnapshot.WithOutputDir(outputDir),
					testsnapshot.WithBasename(name),
					testsnapshot.WithFrameCount(frameCount),
					testsnapshot.WithReset(func() {
						ResetSnapshotState()
					}),
					testsnapshot.WithPreRender(func() {
						SetPolicy(pc.policy)
						SetCycleStartForTest(time.Now().Add(-22500 * time.Millisecond))
					}),
				)

				testsnapshot.VerifyAll(t, pngPath, pan.width, pan.height)
			})
		}
	}
}
