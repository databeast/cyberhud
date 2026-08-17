package tests_test

import (
	"testing"

	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// TestScoredFitness_ColorFastPanel_PrefersColorStyle is a regression test for
// the bug where a MonoFast 128x32 style was incorrectly selected on a
// ColorFast 128x128 panel (Waveshare 1.44" LCD HAT) because both scored
// identically as Optimal with no tie-breaking.
func TestScoredFitness_ColorFastPanel_PrefersColorStyle(t *testing.T) {
	// Panel: Waveshare 1.44" — 128x128, ColorFast (capability=5)
	hints := textlayout.TextHints{
		PixelWidth:  128,
		PixelHeight: 128,
		Capability:  int(style.ColorFast),
	}

	monoFitness := style.EvaluateFitness(style.SurfaceRequirements{
		MinWidth: 128, MinHeight: 32, Capability: style.MonoFast,
	}, hints)

	colorFitness := style.EvaluateFitness(style.SurfaceRequirements{
		MinWidth: 128, MinHeight: 128, Capability: style.ColorFast,
	}, hints)

	if monoFitness < style.Optimal {
		t.Fatalf("mono style should be Optimal tier, got %d", monoFitness)
	}
	if colorFitness < style.Optimal {
		t.Fatalf("color style should be Optimal tier, got %d", colorFitness)
	}
	if colorFitness <= monoFitness {
		t.Fatalf("ColorFast 128x128 (fitness=%d) must beat MonoFast 128x32 (fitness=%d) on ColorFast 128x128 panel",
			colorFitness, monoFitness)
	}
}

// TestScoredFitness_CapabilityProximityTiebreak verifies that styles with
// closer capability to the panel score higher than distant capability matches.
func TestScoredFitness_CapabilityProximityTiebreak(t *testing.T) {
	hints := textlayout.TextHints{
		PixelWidth:  128,
		PixelHeight: 128,
		Capability:  int(style.ColorFast),
	}

	// Both fit the panel, but ColorFast should score higher than MonoFast.
	monoFitness := style.EvaluateFitness(style.SurfaceRequirements{
		MinWidth: 128, MinHeight: 128, Capability: style.MonoFast,
	}, hints)

	colorFitness := style.EvaluateFitness(style.SurfaceRequirements{
		MinWidth: 128, MinHeight: 128, Capability: style.ColorFast,
	}, hints)

	if colorFitness <= monoFitness {
		t.Fatalf("exact capability match (fitness=%d) should beat distant match (fitness=%d)",
			colorFitness, monoFitness)
	}
}

// TestScoredFitness_ResolutionProximityTiebreak verifies that styles whose
// min resolution is closer to the panel score higher (less wasted area).
func TestScoredFitness_ResolutionProximityTiebreak(t *testing.T) {
	hints := textlayout.TextHints{
		PixelWidth:  128,
		PixelHeight: 128,
		Capability:  int(style.ColorFast),
	}

	exactFitness := style.EvaluateFitness(style.SurfaceRequirements{
		MinWidth: 128, MinHeight: 128, Capability: style.ColorFast,
	}, hints)

	smallFitness := style.EvaluateFitness(style.SurfaceRequirements{
		MinWidth: 80, MinHeight: 80, Capability: style.ColorFast,
	}, hints)

	if exactFitness <= smallFitness {
		t.Fatalf("exact resolution match (fitness=%d) should beat smaller (fitness=%d)",
			exactFitness, smallFitness)
	}
}

// TestScoredFitness_MonoPanel_BlocksColorStyles confirms the capability gate
// still correctly blocks styles that require higher capability than the panel.
func TestScoredFitness_MonoPanel_BlocksColorStyles(t *testing.T) {
	hints := textlayout.TextHints{
		PixelWidth:  128,
		PixelHeight: 64,
		Capability:  int(style.MonoFast),
	}

	colorFitness := style.EvaluateFitness(style.SurfaceRequirements{
		MinWidth: 128, MinHeight: 64, Capability: style.ColorFast,
	}, hints)

	if colorFitness != style.Unsupported {
		t.Fatalf("ColorFast style on MonoFast panel should be Unsupported, got %d", colorFitness)
	}
}
