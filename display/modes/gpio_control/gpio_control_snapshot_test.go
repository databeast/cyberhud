package gpio_control

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/databeast/cyberhud/display/modes/gpio_control/source"
	"github.com/databeast/cyberhud/display/modes/testsnapshot"
	gpiomgr "github.com/databeast/cyberhud/hardware/gpio"
)

// mockSnapshotter implements source.Snapshotter for test injection.
type mockSnapshotter struct {
	pins []gpiomgr.PinState
}

func (m *mockSnapshotter) Snapshot() []gpiomgr.PinState {
	return m.pins
}

// snapshotOutputDir is the persistent directory where gpio_control snapshot PNGs are written.
var snapshotOutputDir = filepath.Join("snapshots")

// categoryFromStyleName derives the testsnapshot.DisplayCategory from a gpio_control
// style's name prefix. Longer prefixes are checked before shorter ones to ensure
// correct matching (e.g., "mono-slow-" before "mono-").
func categoryFromStyleName(name string) testsnapshot.DisplayCategory {
	switch {
	case strings.HasPrefix(name, "mono-slow-"):
		return testsnapshot.CategoryMono
	case strings.HasPrefix(name, "mono-fast-"):
		return testsnapshot.CategoryMono
	case strings.HasPrefix(name, "mono-"):
		return testsnapshot.CategoryMono
	case strings.HasPrefix(name, "color-slow-"):
		return testsnapshot.CategoryColor
	case strings.HasPrefix(name, "color-fast-"):
		return testsnapshot.CategoryColor
	case strings.HasPrefix(name, "color-"):
		return testsnapshot.CategoryColor
	case strings.HasPrefix(name, "grayscale-slow-"):
		return testsnapshot.CategoryGrayscale
	case strings.HasPrefix(name, "grayscale-fast-"):
		return testsnapshot.CategoryGrayscale
	case strings.HasPrefix(name, "eink-"):
		return testsnapshot.CategoryEink
	default:
		return testsnapshot.CategoryColor
	}
}

// testGpioPins returns a deterministic set of 8 GPIO pins for snapshot rendering.
// Includes: 3 ModeOutput+HIGH, 2 ModeOutput+LOW, 3 ModeInput — satisfying
// Requirement 6.2 (≥2 output HIGH, ≥2 output LOW, ≥1 input).
func testGpioPins() []gpiomgr.PinState {
	return []gpiomgr.PinState{
		{Number: 4, Name: "GPIO4", Mode: gpiomgr.ModeOutput, Level: true},
		{Number: 17, Name: "GPIO17", Mode: gpiomgr.ModeOutput, Level: false},
		{Number: 27, Name: "GPIO27", Mode: gpiomgr.ModeInput, Level: true},
		{Number: 22, Name: "GPIO22", Mode: gpiomgr.ModeOutput, Level: true},
		{Number: 5, Name: "GPIO5", Mode: gpiomgr.ModeInput, Level: false},
		{Number: 6, Name: "GPIO6", Mode: gpiomgr.ModeOutput, Level: false},
		{Number: 13, Name: "GPIO13", Mode: gpiomgr.ModeOutput, Level: true},
		{Number: 19, Name: "GPIO19", Mode: gpiomgr.ModeInput, Level: true},
	}
}

// TestGpioControlPNGSnapshots exercises all registered gpio_control styles and runs
// a subtest for each, deriving render dimensions from SurfaceRequirements and display
// category from the style name prefix.
//
// Output PNGs are written to snapshots/.
func TestGpioControlPNGSnapshots(t *testing.T) {
	// Ensure the output directory exists and is clean.
	if err := os.RemoveAll(snapshotOutputDir); err != nil {
		t.Fatalf("failed to clean snapshot output directory: %v", err)
	}
	if err := os.MkdirAll(snapshotOutputDir, 0o755); err != nil {
		t.Fatalf("failed to create snapshot output directory: %v", err)
	}

	styles := Registry().Enumerate()
	if len(styles) == 0 {
		t.Fatal("gpioControlRegistry contains zero styles")
	}

	pins := testGpioPins()
	if len(pins) < 6 {
		t.Fatalf("testGpioPins returned %d pins, need at least 6", len(pins))
	}

	for _, s := range styles {
		s := s
		t.Run(s.Name(), func(t *testing.T) {
			reqs := s.Requirements()
			if reqs.MinWidth == 0 || reqs.MinHeight == 0 {
				t.Skip("skipping: style has unconstrained dimensions")
			}

			category := categoryFromStyleName(s.Name())
			policy := DefaultPolicy()
			policy.Style = s.Name()

			// Cursor set to index 0 (GPIO4, ModeOutput) — a toggleable pin.
			pngPath := testsnapshot.RenderSnapshot(t,
				testsnapshot.WithDimensions(reqs.MinWidth, reqs.MinHeight),
				testsnapshot.WithMode("gpio-control"),
				testsnapshot.WithOutputDir(snapshotOutputDir),
				testsnapshot.WithDisplayCategory(category),
				testsnapshot.WithBasename(s.Name()),
				testsnapshot.WithReset(func() {
					source.SetGPIOControlManager(&mockSnapshotter{pins})
				}),
				testsnapshot.WithPreRender(func() {
					SetPolicy(policy)
				}),
			)

			testsnapshot.VerifyAll(t, pngPath, reqs.MinWidth, reqs.MinHeight)
		})
	}
}
