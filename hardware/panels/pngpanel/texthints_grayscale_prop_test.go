package pngpanel

import (
	"os"
	"testing"

	"github.com/databeast/cyberhud/display/surface/textlayout"
	"pgregory.net/rapid"
)

// This test encodes the EXPECTED (correct) behavior: a PNGPanel constructed with
// ColorModeGrayscale should return CapGrayscaleFast (3) from TextHints().Capability.
//
// On UNFIXED code this test will FAIL because the else branch returns CapColorFast (5)
// for both full-color and grayscale modes. Failure confirms the bug exists.

func TestProperty_GrayscaleCapability_BugCondition(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		width := rapid.IntRange(1, 4096).Draw(t, "width")
		height := rapid.IntRange(1, 4096).Draw(t, "height")

		panel, err := New(
			WithDimensions(width, height),
			WithColorMode(ColorModeGrayscale),
			WithOutputDir(os.TempDir()),
		)
		if err != nil {
			t.Fatalf("unexpected construction error: %v", err)
		}

		hints := panel.TextHints()

		// Assert expected behavior: grayscale panels should advertise CapGrayscaleFast (3).
		// Bug condition: on unfixed code, this returns CapColorFast (5) instead.
		if hints.Capability != textlayout.CapGrayscaleFast {
			t.Fatalf("TextHints().Capability = %d (CapColorFast), want %d (CapGrayscaleFast) for grayscale panel with dimensions (%d, %d)",
				hints.Capability, textlayout.CapGrayscaleFast, width, height)
		}
	})
}
