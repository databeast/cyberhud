package pngpanel

import (
	"image"
	"os"
	"path/filepath"
	"testing"

	"github.com/databeast/cyberhud/display/region"
)

// ---------------------------------------------------------------------------
// ---------------------------------------------------------------------------

// TestPNGPanel_ActivatePanel verifies that a PNGPanel can be used as a
// ScreenPosition.Target through the ActivatePanel infrastructure.
// It constructs a PNGPanel, passes it as a screen target, activates the panel,
// then flushes and confirms DrawImage is invoked (producing a PNG file).
func TestPNGPanel_ActivatePanel(t *testing.T) {
	tmpDir := t.TempDir()

	panel, err := New(
		WithDimensions(240, 135),
		WithColorMode(ColorModeFullColor),
		WithOutputDir(tmpDir),
	)
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}

	screen := region.ScreenPosition{
		Index:        0,
		Name:         "test-png-panel",
		Bounds:       image.Rect(0, 0, 240, 135),
		Target:       panel,
		HintProvider: panel.TextHints,
	}

	activation, err := region.ActivatePanel(region.PanelActivationConfig{
		Screens:     []region.ScreenPosition{screen},
		DefaultMode: "clock",
		AvailModes:  []string{"clock"},
		ModeValidator: func(mode string) bool {
			return mode == "clock"
		},
	})
	if err != nil {
		t.Fatalf("ActivatePanel() unexpected error: %v", err)
	}

	// FlushPath should call DrawImage on the PNGPanel target.
	err = activation.FlushPath.Flush()
	if err != nil {
		t.Fatalf("Flush() unexpected error: %v", err)
	}

	// Verify that a PNG file was written (DrawImage was called).
	outPath := filepath.Join(tmpDir, "frame_0001.png")
	info, err := os.Stat(outPath)
	if err != nil {
		t.Fatalf("expected output PNG file at %q, got error: %v", outPath, err)
	}
	if info.Size() == 0 {
		t.Fatalf("output PNG file at %q is empty", outPath)
	}
}
