package tests_test

import (
	"testing"

	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/fonts"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// --- From: evaluate_fitness_test.go ---

// testFontFace is a minimal font.Face implementation for unit testing.
type testFontFace struct {
	id      string
	metrics font.Metrics
}

func (f testFontFace) ID() string                { return f.id }
func (f testFontFace) Metrics() font.Metrics     { return f.metrics }
func (f testFontFace) GlyphRow(rune, int) uint32 { return 0 }

// --- Requirement 1.2: Zero-dimension rejection ---

func TestEvaluateFitness_ZeroPixelWidth(t *testing.T) {
	reqs := style.SurfaceRequirements{}
	hints := textlayout.TextHints{PixelWidth: 0, PixelHeight: 100}
	got := style.EvaluateFitness(reqs, hints)
	if got != style.Unsupported {
		t.Fatalf("expected Unsupported for zero PixelWidth, got %d", got)
	}
}

func TestEvaluateFitness_ZeroPixelHeight(t *testing.T) {
	reqs := style.SurfaceRequirements{}
	hints := textlayout.TextHints{PixelWidth: 100, PixelHeight: 0}
	got := style.EvaluateFitness(reqs, hints)
	if got != style.Unsupported {
		t.Fatalf("expected Unsupported for zero PixelHeight, got %d", got)
	}
}

func TestEvaluateFitness_BothDimensionsZero(t *testing.T) {
	reqs := style.SurfaceRequirements{}
	hints := textlayout.TextHints{PixelWidth: 0, PixelHeight: 0}
	got := style.EvaluateFitness(reqs, hints)
	if got != style.Unsupported {
		t.Fatalf("expected Unsupported for both dimensions zero, got %d", got)
	}
}

// --- Requirement 1.3: MinWidth boundary cases ---

func TestEvaluateFitness_BelowMinWidth(t *testing.T) {
	reqs := style.SurfaceRequirements{MinWidth: 100}
	hints := textlayout.TextHints{PixelWidth: 99, PixelHeight: 200}
	got := style.EvaluateFitness(reqs, hints)
	if got != style.Unsupported {
		t.Fatalf("expected Unsupported when PixelWidth < MinWidth, got %d", got)
	}
}

func TestEvaluateFitness_ExactMinWidth(t *testing.T) {
	reqs := style.SurfaceRequirements{MinWidth: 100}
	hints := textlayout.TextHints{PixelWidth: 100, PixelHeight: 200}
	got := style.EvaluateFitness(reqs, hints)
	if got == style.Unsupported {
		t.Fatalf("expected supported (Full or Optimal) when PixelWidth == MinWidth, got Unsupported")
	}
}

func TestEvaluateFitness_AboveMinWidth(t *testing.T) {
	reqs := style.SurfaceRequirements{MinWidth: 100}
	hints := textlayout.TextHints{PixelWidth: 101, PixelHeight: 200}
	got := style.EvaluateFitness(reqs, hints)
	if got == style.Unsupported {
		t.Fatalf("expected supported when PixelWidth > MinWidth, got Unsupported")
	}
}

func TestEvaluateFitness_MinWidthZeroIgnored(t *testing.T) {
	// MinWidth of 0 means unconstrained; any positive PixelWidth should pass.
	reqs := style.SurfaceRequirements{MinWidth: 0}
	hints := textlayout.TextHints{PixelWidth: 1, PixelHeight: 200}
	got := style.EvaluateFitness(reqs, hints)
	if got == style.Unsupported {
		t.Fatalf("expected supported when MinWidth is 0 (unconstrained), got Unsupported")
	}
}

// --- Requirement 1.4: MinHeight boundary cases ---

func TestEvaluateFitness_BelowMinHeight(t *testing.T) {
	reqs := style.SurfaceRequirements{MinHeight: 50}
	hints := textlayout.TextHints{PixelWidth: 200, PixelHeight: 49}
	got := style.EvaluateFitness(reqs, hints)
	if got != style.Unsupported {
		t.Fatalf("expected Unsupported when PixelHeight < MinHeight, got %d", got)
	}
}

func TestEvaluateFitness_ExactMinHeight(t *testing.T) {
	reqs := style.SurfaceRequirements{MinHeight: 50}
	hints := textlayout.TextHints{PixelWidth: 200, PixelHeight: 50}
	got := style.EvaluateFitness(reqs, hints)
	if got == style.Unsupported {
		t.Fatalf("expected supported when PixelHeight == MinHeight, got Unsupported")
	}
}

func TestEvaluateFitness_MinHeightZeroIgnored(t *testing.T) {
	reqs := style.SurfaceRequirements{MinHeight: 0}
	hints := textlayout.TextHints{PixelWidth: 200, PixelHeight: 1}
	got := style.EvaluateFitness(reqs, hints)
	if got == style.Unsupported {
		t.Fatalf("expected supported when MinHeight is 0 (unconstrained), got Unsupported")
	}
}

// --- Requirement 1.5: Capability ordering gating ---

func TestEvaluateFitness_CapabilityRequired_NotAvailable(t *testing.T) {
	reqs := style.SurfaceRequirements{Capability: style.ColorFast}
	hints := textlayout.TextHints{PixelWidth: 200, PixelHeight: 200, Capability: int(style.MonoFast)}
	got := style.EvaluateFitness(reqs, hints)
	if got != style.Unsupported {
		t.Fatalf("expected Unsupported when reqs.Capability > hints.Capability, got %d", got)
	}
}

func TestEvaluateFitness_CapabilityRequired_Available(t *testing.T) {
	reqs := style.SurfaceRequirements{Capability: style.ColorFast}
	hints := textlayout.TextHints{PixelWidth: 200, PixelHeight: 200, Capability: int(style.ColorFast)}
	got := style.EvaluateFitness(reqs, hints)
	if got == style.Unsupported {
		t.Fatalf("expected supported when hints.Capability >= reqs.Capability, got Unsupported")
	}
}

func TestEvaluateFitness_CapabilityNotRequired_LowPanel(t *testing.T) {
	// When Capability is MonoSlow (zero value), any panel satisfies.
	reqs := style.SurfaceRequirements{Capability: style.MonoSlow}
	hints := textlayout.TextHints{PixelWidth: 200, PixelHeight: 200, Capability: int(style.MonoSlow)}
	got := style.EvaluateFitness(reqs, hints)
	if got == style.Unsupported {
		t.Fatalf("expected supported when Capability is MonoSlow (unconstrained), got Unsupported")
	}
}

// --- Requirement 1.6: TextFitness delegation ---

func TestEvaluateFitness_TextFitnessDelegation_MinRowsPositive(t *testing.T) {
	// Register a font that can satisfy the text constraints.
	restore := font.SnapshotAndClear()
	defer restore()

	font.Register(testFontFace{
		id:      "test-5x8",
		metrics: font.Metrics{GlyphWidth: 5, GlyphHeight: 8, GlyphAdvance: 6, RowHeight: 10},
	})

	// MinRows=2, RowHeight=10 → needs PixelHeight >= 20.
	reqs := style.SurfaceRequirements{MinRows: 2}
	hints := textlayout.TextHints{PixelWidth: 200, PixelHeight: 200}
	got := style.EvaluateFitness(reqs, hints)
	if got == style.Unsupported {
		t.Fatalf("expected supported when a font satisfies MinRows constraint, got Unsupported")
	}
}

func TestEvaluateFitness_TextFitnessDelegation_MinRowsUnsatisfiable(t *testing.T) {
	// Register a font whose RowHeight makes it impossible to fit MinRows.
	restore := font.SnapshotAndClear()
	defer restore()

	font.Register(testFontFace{
		id:      "test-big",
		metrics: font.Metrics{GlyphWidth: 5, GlyphHeight: 8, GlyphAdvance: 6, RowHeight: 100},
	})

	// MinRows=3, RowHeight=100 → needs PixelHeight >= 300, but we only have 200.
	reqs := style.SurfaceRequirements{MinRows: 3}
	hints := textlayout.TextHints{PixelWidth: 200, PixelHeight: 200}
	got := style.EvaluateFitness(reqs, hints)
	if got != style.Unsupported {
		t.Fatalf("expected Unsupported when no font can satisfy MinRows, got %d", got)
	}
}

func TestEvaluateFitness_TextFitnessDelegation_MinCharsPerLinePositive(t *testing.T) {
	restore := font.SnapshotAndClear()
	defer restore()

	font.Register(testFontFace{
		id:      "test-5x8",
		metrics: font.Metrics{GlyphWidth: 5, GlyphHeight: 8, GlyphAdvance: 6, RowHeight: 10},
	})

	// MinCharsPerLine=10, GlyphAdvance=6 → needs PixelWidth >= 60.
	reqs := style.SurfaceRequirements{MinCharsPerLine: 10}
	hints := textlayout.TextHints{PixelWidth: 200, PixelHeight: 200}
	got := style.EvaluateFitness(reqs, hints)
	if got == style.Unsupported {
		t.Fatalf("expected supported when a font satisfies MinCharsPerLine constraint, got Unsupported")
	}
}

func TestEvaluateFitness_TextFitnessDelegation_MinCharsUnsatisfiable(t *testing.T) {
	restore := font.SnapshotAndClear()
	defer restore()

	font.Register(testFontFace{
		id:      "test-wide",
		metrics: font.Metrics{GlyphWidth: 10, GlyphHeight: 8, GlyphAdvance: 50, RowHeight: 10},
	})

	// MinCharsPerLine=10, GlyphAdvance=50 → needs PixelWidth >= 500, but only 200.
	reqs := style.SurfaceRequirements{MinCharsPerLine: 10}
	hints := textlayout.TextHints{PixelWidth: 200, PixelHeight: 200}
	got := style.EvaluateFitness(reqs, hints)
	if got != style.Unsupported {
		t.Fatalf("expected Unsupported when no font can satisfy MinCharsPerLine, got %d", got)
	}
}

func TestEvaluateFitness_TextFitnessDelegation_EmptyRegistry(t *testing.T) {
	restore := font.SnapshotAndClear()
	defer restore()

	// No fonts registered; TextFitness should return Unsupported.
	reqs := style.SurfaceRequirements{MinRows: 1, MinCharsPerLine: 1}
	hints := textlayout.TextHints{PixelWidth: 200, PixelHeight: 200}
	got := style.EvaluateFitness(reqs, hints)
	if got != style.Unsupported {
		t.Fatalf("expected Unsupported when font registry is empty and text constraints are set, got %d", got)
	}
}

// --- Requirement 1.7: Optimal when preferred dimensions met ---

func TestEvaluateFitness_Optimal_BothPreferredMet(t *testing.T) {
	reqs := style.SurfaceRequirements{
		MinWidth:        50,
		MinHeight:       50,
		PreferredWidth:  100,
		PreferredHeight: 100,
	}
	hints := textlayout.TextHints{PixelWidth: 150, PixelHeight: 150}
	got := style.EvaluateFitness(reqs, hints)
	if got < style.Optimal {
		t.Fatalf("expected Optimal tier when both preferred dimensions are met, got %d", got)
	}
}

func TestEvaluateFitness_Optimal_ExactPreferred(t *testing.T) {
	reqs := style.SurfaceRequirements{
		PreferredWidth:  100,
		PreferredHeight: 100,
	}
	hints := textlayout.TextHints{PixelWidth: 100, PixelHeight: 100}
	got := style.EvaluateFitness(reqs, hints)
	if got < style.Optimal {
		t.Fatalf("expected Optimal tier when exactly at preferred dimensions, got %d", got)
	}
}

func TestEvaluateFitness_Optimal_PreferredZeroMeansUnconstrained(t *testing.T) {
	// PreferredWidth=0 and PreferredHeight=0 means "no preference" → always Optimal (if min met).
	reqs := style.SurfaceRequirements{
		PreferredWidth:  0,
		PreferredHeight: 0,
	}
	hints := textlayout.TextHints{PixelWidth: 10, PixelHeight: 10}
	got := style.EvaluateFitness(reqs, hints)
	if got < style.Optimal {
		t.Fatalf("expected Optimal tier when both preferred are 0 (unconstrained), got %d", got)
	}
}

// --- Requirement 1.8: Full when min met but preferred not ---

func TestEvaluateFitness_Full_PreferredWidthNotMet(t *testing.T) {
	reqs := style.SurfaceRequirements{
		MinWidth:       50,
		PreferredWidth: 200,
	}
	hints := textlayout.TextHints{PixelWidth: 100, PixelHeight: 200}
	got := style.EvaluateFitness(reqs, hints)
	if got < style.Full || got >= style.Optimal {
		t.Fatalf("expected Full tier when min met but preferred width not met, got %d", got)
	}
}

func TestEvaluateFitness_Full_PreferredHeightNotMet(t *testing.T) {
	reqs := style.SurfaceRequirements{
		MinHeight:       50,
		PreferredHeight: 200,
	}
	hints := textlayout.TextHints{PixelWidth: 200, PixelHeight: 100}
	got := style.EvaluateFitness(reqs, hints)
	if got < style.Full || got >= style.Optimal {
		t.Fatalf("expected Full tier when min met but preferred height not met, got %d", got)
	}
}

func TestEvaluateFitness_Full_BothPreferredNotMet(t *testing.T) {
	reqs := style.SurfaceRequirements{
		MinWidth:        50,
		MinHeight:       50,
		PreferredWidth:  300,
		PreferredHeight: 300,
	}
	hints := textlayout.TextHints{PixelWidth: 100, PixelHeight: 100}
	got := style.EvaluateFitness(reqs, hints)
	if got < style.Full || got >= style.Optimal {
		t.Fatalf("expected Full tier when min met but neither preferred is met, got %d", got)
	}
}

// --- Combined scenarios ---

func TestEvaluateFitness_AllChecksPass_NoPreferred(t *testing.T) {
	// When there's no preferred and all minimums are met, result is Optimal tier.
	reqs := style.SurfaceRequirements{
		MinWidth:  50,
		MinHeight: 50,
	}
	hints := textlayout.TextHints{PixelWidth: 100, PixelHeight: 100}
	got := style.EvaluateFitness(reqs, hints)
	if got < style.Optimal {
		t.Fatalf("expected Optimal tier when all mins met and no preferred set, got %d", got)
	}
}

func TestEvaluateFitness_MinWidthRejectsBeforePreferred(t *testing.T) {
	// MinWidth check happens before preferred dimension check.
	reqs := style.SurfaceRequirements{
		MinWidth:       200,
		PreferredWidth: 300,
	}
	hints := textlayout.TextHints{PixelWidth: 100, PixelHeight: 200}
	got := style.EvaluateFitness(reqs, hints)
	if got != style.Unsupported {
		t.Fatalf("expected Unsupported when below MinWidth, regardless of PreferredWidth, got %d", got)
	}
}

func TestEvaluateFitness_CapabilityRejectsBeforePreferred(t *testing.T) {
	// Capability check happens before preferred dimension check.
	reqs := style.SurfaceRequirements{
		Capability:     style.ColorFast,
		PreferredWidth: 100,
	}
	hints := textlayout.TextHints{PixelWidth: 200, PixelHeight: 200, Capability: int(style.MonoSlow)}
	got := style.EvaluateFitness(reqs, hints)
	if got != style.Unsupported {
		t.Fatalf("expected Unsupported when capability unmet, got %d", got)
	}
}
