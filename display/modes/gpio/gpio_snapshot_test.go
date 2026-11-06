package gpio

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/databeast/cyberhud/display/modes/testsnapshot"
	"github.com/databeast/cyberhud/display/style"
	gpiomgr "github.com/databeast/cyberhud/hardware/gpio"
)

// snapshotOutputDir is the persistent directory where GPIO snapshot PNGs are written.
var snapshotOutputDir = filepath.Join("snapshots")

// fakeGPIOManager is a deterministic, hardware-free GPIO source for snapshots.
//
// The real *gpiomgr.Manager sources pin state from the periph gpioreg registry,
// which is empty on a dev/CI host with no Pi hardware: every pin comes back
// ModeUnknown/Low, so BuildView takes its "GPIO unavailable"/"GPIO idle" path
// and every snapshot renders the same empty message instead of a real layout.
// This fake returns a fixed, representative roster so each style renders actual
// content and the snapshot corpus is stable across runs and machines.
type fakeGPIOManager struct{ pins []gpiomgr.PinState }

func (f fakeGPIOManager) Snapshot() []gpiomgr.PinState { return f.pins }

// snapshotPins is the deterministic pin roster used across all GPIO snapshots.
// It models a realistic headless deployment so the live-functions view is fully
// exercised by a single render: a few STANDALONE digital pins (outputs and
// inputs, HIGH and LOW) that render individually, plus two active BUSES — I2C1
// (2 pins) and SPI0 (5 pins) — that must COLLAPSE to one row each, plus DEAD
// unconfigured pins that must be dropped entirely. Pin numbers and bus
// assignments follow the real BCM pinmux (I2C1 on 2/3, SPI0 on 7–11).
var snapshotPins = []gpiomgr.PinState{
	// Active I2C1 bus (two pins → one "I2C1 up" row).
	{Number: 2, Name: "GPIO2", Mode: gpiomgr.ModeAlt, Func: "I2C1_SDA"},
	{Number: 3, Name: "GPIO3", Mode: gpiomgr.ModeAlt, Func: "I2C1_SCL"},
	// Active SPI0 bus (five pins → one "SPI0 up" row).
	{Number: 7, Name: "GPIO7", Mode: gpiomgr.ModeAlt, Func: "SPI0_CE1"},
	{Number: 8, Name: "GPIO8", Mode: gpiomgr.ModeAlt, Func: "SPI0_CE0"},
	{Number: 9, Name: "GPIO9", Mode: gpiomgr.ModeAlt, Func: "SPI0_MISO"},
	{Number: 10, Name: "GPIO10", Mode: gpiomgr.ModeAlt, Func: "SPI0_MOSI"},
	{Number: 11, Name: "GPIO11", Mode: gpiomgr.ModeAlt, Func: "SPI0_SCLK"},
	// Standalone digital pins (each its own row, level is the signal).
	{Number: 17, Name: "GPIO17", Mode: gpiomgr.ModeOutput, Level: true, Func: "Out"},
	{Number: 22, Name: "GPIO22", Mode: gpiomgr.ModeOutput, Level: false, Func: "Out"},
	{Number: 27, Name: "GPIO27", Mode: gpiomgr.ModeInput, Level: true, Func: "In"},
	// Dead/unconfigured pins (must be dropped — no row).
	{Number: 4, Name: "GPIO4", Mode: gpiomgr.ModeUnknown},
	{Number: 23, Name: "GPIO23", Mode: gpiomgr.ModeUnknown},
}

// testGpioPolicy returns the deterministic policy used across all GPIO snapshot tests.
func testGpioPolicy() Policy {
	return Policy{
		Style:   "list",
		Color:   true,
		Font:    "auto",
		FGColor: "cyan",
	}
}

// categoryFromStyle derives the testsnapshot.DisplayCategory from style name or requirements.
func categoryFromStyle(name string, reqs style.SurfaceRequirements) testsnapshot.DisplayCategory {
	switch {
	case strings.HasPrefix(name, "color-"):
		return testsnapshot.CategoryColor
	case strings.HasPrefix(name, "mono-slow-"):
		return testsnapshot.CategoryEink
	case strings.HasPrefix(name, "mono-"):
		return testsnapshot.CategoryMono
	case strings.HasPrefix(name, "grayscale-"):
		return testsnapshot.CategoryGrayscale
	}
	switch reqs.Capability {
	case style.GrayscaleSlow, style.GrayscaleFast:
		return testsnapshot.CategoryGrayscale
	case style.ColorSlow, style.ColorFast:
		return testsnapshot.CategoryColor
	case style.MonoFast:
		return testsnapshot.CategoryMono
	default:
		return testsnapshot.CategoryEink
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

func categoryString(category testsnapshot.DisplayCategory) string {
	switch category {
	case testsnapshot.CategoryMono:
		return "mono"
	case testsnapshot.CategoryEink:
		return "eink"
	case testsnapshot.CategoryGrayscale:
		return "grayscale"
	default:
		return "color"
	}
}

// TestGpioPNGSnapshots exercises all registered GPIO styles that have explicit dimensions
// and verifies that each produces a valid PNG with correct dimensions and visible content
// through the full production render pipeline via the snapshottest framework.
//
// Output PNGs are written to snapshots/ with naming:
//
//	{width}x{height}_{style}_{category}_0001.png
//

func TestGpioPNGSnapshots(t *testing.T) {
	// Ensure the output directory exists and is clean.
	if err := os.RemoveAll(snapshotOutputDir); err != nil {
		t.Fatalf("failed to clean snapshot output directory: %v", err)
	}
	if err := os.MkdirAll(snapshotOutputDir, 0o755); err != nil {
		t.Fatalf("failed to create snapshot output directory: %v", err)
	}

	styles := Registry().Enumerate()

	for _, s := range styles {
		s := s
		reqs := s.Requirements()
		width, height := dimensionsFromRequirements(reqs)
		category := categoryFromStyle(s.Name(), reqs)

		// Derive the category string for the basename.
		categoryStr := categoryString(category)
		basename := formatBasename(width, height, s.Name(), categoryStr)

		t.Run(s.Name(), func(t *testing.T) {
			// Set a deterministic policy for this render, using the style under test.
			policy := testGpioPolicy()
			policy.Style = s.Name()

			pngPath := testsnapshot.RenderSnapshot(t,
				testsnapshot.WithDimensions(width, height),
				testsnapshot.WithMode("gpio"),
				testsnapshot.WithOutputDir(snapshotOutputDir),
				testsnapshot.WithDisplayCategory(category),
				testsnapshot.WithBasename(basename),
				// Register the deterministic fake BEFORE the instance is
				// constructed. RenderSnapshot runs WithReset prior to SetMode,
				// and the gpio instance captures source.GPIOManager() at
				// construction into i.mgr; a fake set later (in WithPreRender)
				// would be too late to affect that captured reference.
				testsnapshot.WithReset(func() {
					SetGPIOManager(fakeGPIOManager{pins: snapshotPins})
				}),
				testsnapshot.WithPreRender(func() {
					SetPolicy(policy)
				}),
			)

			testsnapshot.VerifyAll(t, pngPath, width, height)
		})
	}
}

// formatBasename produces the output filename prefix:
//
//	{width}x{height}_{style}_{category}
func formatBasename(width, height int, styleName, category string) string {
	return fmt.Sprintf("%dx%d_%s_%s", width, height, styleName, category)
}
