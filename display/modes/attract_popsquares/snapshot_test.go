package attract_popsquares

import (
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/databeast/cyberhud/display/modes/testsnapshot"
	"github.com/databeast/cyberhud/display/style"
)

// snapshotOutputDir is the persistent directory where popsquares snapshot PNGs are written.
var snapshotOutputDir = filepath.Join("snapshots")

// categoryFromStyleName derives the testsnapshot.DisplayCategory from a
// popsquares style's name prefix.
func categoryFromStyleName(name string) testsnapshot.DisplayCategory {
	switch {
	case strings.HasPrefix(name, "color-slow-"):
		return testsnapshot.CategoryColorSlow
	case strings.HasPrefix(name, "color-"):
		return testsnapshot.CategoryColor
	case strings.HasPrefix(name, "mono-slow-"):
		return testsnapshot.CategoryEink
	case strings.HasPrefix(name, "mono-"):
		return testsnapshot.CategoryMono
	case strings.HasPrefix(name, "grayscale-slow-"):
		return testsnapshot.CategoryGrayscaleSlow
	case strings.HasPrefix(name, "grayscale-fast-"):
		return testsnapshot.CategoryGrayscale
	default:
		return testsnapshot.CategoryColor
	}
}

func TestSnapshotPanelCapabilities(t *testing.T) {
	tests := []struct {
		name     string
		category testsnapshot.DisplayCategory
		width    int
		height   int
		color    bool
		slow     bool
	}{
		{"color-slow-300x400", testsnapshot.CategoryColorSlow, 300, 400, true, true},
		{"grayscale-slow-300x400", testsnapshot.CategoryGrayscaleSlow, 300, 400, false, true},
		{"mono-slow-300x400", testsnapshot.CategoryEink, 300, 400, false, true},
		{"color-fast-240x240", testsnapshot.CategoryColor, 240, 240, true, false},
		{"grayscale-fast-240x240", testsnapshot.CategoryGrayscale, 240, 240, false, false},
		{"mono-fast-128x64", testsnapshot.CategoryMono, 128, 64, false, false},
	}

	oldPolicy := GetPolicy()
	oldFrameCounter, oldLastTick, oldTimePhase, oldEinkPanel := frameCounter, lastTick, timePhase, isEinkPanel
	t.Cleanup(func() {
		SetPolicy(oldPolicy)
		frameCounter, lastTick, timePhase, isEinkPanel = oldFrameCounter, oldLastTick, oldTimePhase, oldEinkPanel
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			category := categoryFromStyleName(tt.name)
			if category != tt.category {
				t.Fatalf("category = %v, want %v", category, tt.category)
			}
			path := testsnapshot.RenderSnapshot(t,
				testsnapshot.WithMode("attract_popsquares"),
				testsnapshot.WithDimensions(tt.width, tt.height),
				testsnapshot.WithDisplayCategory(category),
				testsnapshot.WithOutputDir(t.TempDir()),
				testsnapshot.WithBasename(tt.name),
				testsnapshot.WithFrameCount(1),
				testsnapshot.WithReset(func() {
					frameCounter = 0
					lastTick = time.Time{}
					timePhase = 0
					isEinkPanel = false
				}),
				testsnapshot.WithPreRender(func() { SetPolicy(DefaultPolicy()) }),
			)
			if isEinkPanel != tt.slow {
				t.Fatalf("slow-refresh rendering = %t, want %t", isEinkPanel, tt.slow)
			}
			f, err := os.Open(path)
			if err != nil {
				t.Fatal(err)
			}
			img, err := png.Decode(f)
			closeErr := f.Close()
			if err != nil {
				t.Fatal(err)
			}
			if closeErr != nil {
				t.Fatal(closeErr)
			}
			colored := false
			bounds := img.Bounds()
			for y := bounds.Min.Y; y < bounds.Max.Y && !colored; y++ {
				for x := bounds.Min.X; x < bounds.Max.X; x++ {
					r, g, b, _ := img.At(x, y).RGBA()
					if r != g || g != b {
						colored = true
						break
					}
				}
			}
			if colored != tt.color {
				t.Fatalf("snapshot contains color = %t, want %t", colored, tt.color)
			}
		})
	}
}

// TestPopsquaresPNGSnapshots enumerates all registered popsquares styles and
// renders each through the full production pipeline via the snapshottest framework.
func TestPopsquaresPNGSnapshots(t *testing.T) {
	styles := popsquaresRegistry.Enumerate()
	if len(styles) == 0 {
		t.Fatal("popsquaresRegistry contains zero styles")
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

			category := categoryFromStyleName(s.Name())

			frameCount := 5
			if reqs.Capability == style.MonoSlow || reqs.Capability == style.GrayscaleSlow || reqs.Capability == style.ColorSlow {
				frameCount = 1
			}

			pngPath := testsnapshot.RenderSnapshot(t,
				testsnapshot.WithMode("attract_popsquares"),
				testsnapshot.WithDimensions(reqs.MinWidth, reqs.MinHeight),
				testsnapshot.WithDisplayCategory(category),
				testsnapshot.WithOutputDir(snapshotOutputDir),
				testsnapshot.WithBasename(s.Name()),
				testsnapshot.WithFrameCount(frameCount),
				testsnapshot.WithReset(func() {
					frameCounter = 0
					lastTick = time.Time{}
					timePhase = 0
					isEinkPanel = false
				}),
				testsnapshot.WithPreRender(func() {
					SetPolicy(DefaultPolicy())
				}),
			)

			testsnapshot.VerifyAll(t, pngPath, reqs.MinWidth, reqs.MinHeight)
			pngPaths = append(pngPaths, pngPath)
		})
	}

	if err := testsnapshot.WriteGalleryFragmentFromPaths(snapshotOutputDir, "attract_popsquares", pngPaths); err != nil {
		t.Fatalf("failed to write gallery fragment: %v", err)
	}
}
