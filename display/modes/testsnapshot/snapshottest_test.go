package testsnapshot

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	gpiomgr "github.com/databeast/cyberhud/hardware/gpio"
	"github.com/databeast/cyberhud/hardware/panels/pngpanel"

	// Blank imports to trigger mode init() self-registration.
	_ "github.com/databeast/cyberhud/display/modes/attract_matrix"
	_ "github.com/databeast/cyberhud/display/modes/clock"
	_ "github.com/databeast/cyberhud/display/modes/cycle"
	_ "github.com/databeast/cyberhud/display/modes/dashboard"
	_ "github.com/databeast/cyberhud/display/modes/gpio"
	_ "github.com/databeast/cyberhud/display/modes/gpio_control"
	_ "github.com/databeast/cyberhud/display/modes/image"
	_ "github.com/databeast/cyberhud/display/modes/menu"
	_ "github.com/databeast/cyberhud/display/modes/serial"
	_ "github.com/databeast/cyberhud/display/modes/stemma"
	_ "github.com/databeast/cyberhud/display/modes/system"
	_ "github.com/databeast/cyberhud/display/modes/systemd"
	_ "github.com/databeast/cyberhud/display/modes/testfonts"
	_ "github.com/databeast/cyberhud/display/modes/testpattern"
	_ "github.com/databeast/cyberhud/display/modes/thermal"
	_ "github.com/databeast/cyberhud/display/modes/ticker"
	_ "github.com/databeast/cyberhud/display/modes/usb"
	_ "github.com/databeast/cyberhud/display/modes/wifi"
	_ "github.com/databeast/cyberhud/display/modes/zmq"
)

// TestPreRenderCallbackInvoked verifies that the PreRender callback runs before render.
// We use a flag variable set inside the callback and check it was set after RenderSnapshot.
func TestPreRenderCallbackInvoked(t *testing.T) {
	outputDir := t.TempDir()
	called := false

	pngPath := RenderSnapshot(t,
		WithDimensions(240, 320),
		WithMode("clock"),
		WithOutputDir(outputDir),
		WithDisplayCategory(CategoryColor),
		WithPreRender(func() {
			called = true
		}),
	)

	if !called {
		t.Fatal("PreRender callback was not invoked")
	}
	VerifyPNG(t, pngPath)
}

// TestResetCallbackInvokedBeforePreRender verifies that the Reset callback runs
// before the PreRender callback, establishing the correct ordering.
func TestResetCallbackInvokedBeforePreRender(t *testing.T) {
	outputDir := t.TempDir()
	var order []string

	pngPath := RenderSnapshot(t,
		WithDimensions(240, 320),
		WithMode("clock"),
		WithOutputDir(outputDir),
		WithDisplayCategory(CategoryColor),
		WithReset(func() {
			order = append(order, "reset")
		}),
		WithPreRender(func() {
			order = append(order, "preRender")
		}),
	)

	if len(order) != 2 {
		t.Fatalf("expected 2 callbacks, got %d: %v", len(order), order)
	}
	if order[0] != "reset" {
		t.Fatalf("expected reset first, got %q", order[0])
	}
	if order[1] != "preRender" {
		t.Fatalf("expected preRender second, got %q", order[1])
	}
	VerifyPNG(t, pngPath)
}

// TestGPIOManagerPassthrough verifies that passing a real GPIO Manager via
// WithGPIOManager allows the gpio mode to render successfully.
func TestGPIOManagerPassthrough(t *testing.T) {
	outputDir := t.TempDir()

	gm := gpiomgr.New()

	pngPath := RenderSnapshot(t,
		WithDimensions(240, 135),
		WithMode("gpio"),
		WithOutputDir(outputDir),
		WithDisplayCategory(CategoryColor),
		WithGPIOManager(gm),
	)

	VerifyAll(t, pngPath, 240, 135)
}

// TestPassiveModeNormalizesMenuToDashboard verifies that setting mode to "menu"
// in passive mode (inputMapper=nil) still renders successfully because the
// ModeEngine normalizes menu→dashboard.
func TestPassiveModeNormalizesMenuToDashboard(t *testing.T) {
	outputDir := t.TempDir()

	// "menu" is a registered mode but in passive mode it normalizes to dashboard.
	// The key assertion is that this does NOT fail — it renders dashboard output.
	pngPath := RenderSnapshot(t,
		WithDimensions(240, 320),
		WithMode("menu"),
		WithOutputDir(outputDir),
		WithDisplayCategory(CategoryColor),
	)

	VerifyAll(t, pngPath, 240, 320)
}

// --- Fatal/error case tests ---
//
// These tests verify that validate() calls t.Fatal for invalid configurations.
// Since t.Fatal terminates the test goroutine, we use the subprocess pattern:
// the test re-invokes itself as a subprocess targeting a helper test function,
// and checks that the subprocess exits with a non-zero exit code.

// TestMissingDimensionsFatals verifies that omitting dimensions causes a fatal.
func TestMissingDimensionsFatals(t *testing.T) {
	if os.Getenv("SNAPSHOTTEST_FATAL_HELPER") == "1" {
		// This runs in the subprocess. It should fatal.
		outputDir := t.TempDir()
		RenderSnapshot(t,
			// No WithDimensions — should fatal
			WithMode("clock"),
			WithOutputDir(outputDir),
			WithDisplayCategory(CategoryColor),
		)
		return
	}

	// Run the test as a subprocess.
	cmd := exec.Command(os.Args[0], "-test.run=^TestMissingDimensionsFatals$", "-test.v")
	cmd.Env = append(os.Environ(), "SNAPSHOTTEST_FATAL_HELPER=1")
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected subprocess to fail (missing dimensions), but it succeeded")
	}
	if !strings.Contains(string(output), "missing required option") {
		t.Fatalf("expected 'missing required option' in output, got:\n%s", output)
	}
}

// TestInvalidModeFatals verifies that specifying an unregistered mode causes a fatal.
func TestInvalidModeFatals(t *testing.T) {
	if os.Getenv("SNAPSHOTTEST_FATAL_HELPER") == "1" {
		outputDir := t.TempDir()
		RenderSnapshot(t,
			WithDimensions(240, 320),
			WithMode("nonexistent_mode_xyz"),
			WithOutputDir(outputDir),
			WithDisplayCategory(CategoryColor),
		)
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=^TestInvalidModeFatals$", "-test.v")
	cmd.Env = append(os.Environ(), "SNAPSHOTTEST_FATAL_HELPER=1")
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected subprocess to fail (invalid mode), but it succeeded")
	}
	if !strings.Contains(string(output), "not found in registry") {
		t.Fatalf("expected 'not found in registry' in output, got:\n%s", output)
	}
}

// TestBothCategoryAndColorModeFatals verifies that specifying both
// WithDisplayCategory and WithColorMode causes a fatal.
func TestBothCategoryAndColorModeFatals(t *testing.T) {
	if os.Getenv("SNAPSHOTTEST_FATAL_HELPER") == "1" {
		outputDir := t.TempDir()
		RenderSnapshot(t,
			WithDimensions(240, 320),
			WithMode("clock"),
			WithOutputDir(outputDir),
			WithDisplayCategory(CategoryColor),
			WithColorMode(pngpanel.ColorModeMonochrome),
		)
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=^TestBothCategoryAndColorModeFatals$", "-test.v")
	cmd.Env = append(os.Environ(), "SNAPSHOTTEST_FATAL_HELPER=1")
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected subprocess to fail (both category and colorMode), but it succeeded")
	}
	if !strings.Contains(string(output), "mutually exclusive") {
		t.Fatalf("expected 'mutually exclusive' in output, got:\n%s", output)
	}
}

// TestOutputPathFormat verifies that the output path follows the expected format.
func TestOutputPathFormat(t *testing.T) {
	outputDir := t.TempDir()

	pngPath := RenderSnapshot(t,
		WithDimensions(240, 320),
		WithMode("clock"),
		WithOutputDir(outputDir),
		WithDisplayCategory(CategoryColor),
		WithBasename("mytest"),
	)

	expected := fmt.Sprintf("%s_0001.png", "mytest")
	if !strings.HasSuffix(pngPath, expected) {
		t.Fatalf("expected path ending with %q, got %q", expected, pngPath)
	}
}
