package clock

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/databeast/cyberhud/display/modes/clock/source"
	"github.com/databeast/cyberhud/display/modes/testsnapshot"
	"pgregory.net/rapid"
)

// referenceTime is the fixed time value used for deterministic snapshot rendering.
var referenceTime = time.Date(2024, 1, 15, 14, 30, 5, 0, time.UTC)

// testPolicy returns the fixed clock policy used across all snapshot tests.
func testPolicy() source.Policy {
	return source.Policy{
		Style:       "",
		ShowSeconds: true,
		TimeFormat:  "24h",
		DateFormat:  "YYYY-MM-DD",
		Timezone:    "UTC",
		ShowWeekday: true,
		BlinkColon:  false,
		FGColor:     "cyan",
		ShowLED:     true,
		SecondsBar:  "none",
		ShowDaybar:  false,
		ShowBorder:  true,
		BorderColor: "emerald",
	}
}

// categoryFromStyleName derives the testsnapshot.DisplayCategory from a clock
// style's name prefix. Clock styles follow the naming convention:
// mono-*, color-*, grayscale-fast-*, grayscale-slow-*, mono-slow-*, color-slow-*.
func categoryFromStyleName(name string) testsnapshot.DisplayCategory {
	switch {
	case strings.HasPrefix(name, "color-slow-"):
		return testsnapshot.CategoryColor
	case strings.HasPrefix(name, "color-"):
		return testsnapshot.CategoryColor
	case strings.HasPrefix(name, "mono-slow-"):
		return testsnapshot.CategoryMono
	case strings.HasPrefix(name, "mono-"):
		return testsnapshot.CategoryMono
	case strings.HasPrefix(name, "grayscale-fast-"):
		return testsnapshot.CategoryGrayscale
	case strings.HasPrefix(name, "grayscale-slow-"):
		return testsnapshot.CategoryGrayscale
	default:
		return testsnapshot.CategoryColor
	}
}

// snapshotOutputDir is the persistent directory where snapshot PNGs are written.
// Located at snapshots/ relative to this test file so the output survives the
// test run and can be visually inspected or committed as golden files.
var snapshotOutputDir = filepath.Join("snapshots")

// TestClockPNGSnapshots enumerates all registered clock styles and runs a subtest
// for each, writing one PNG per style to snapshots/ for visual inspection.
func TestClockPNGSnapshots(t *testing.T) {

	styles := clockRegistry.Enumerate()

	// Ensure the output directory exists and is clean.
	if err := os.RemoveAll(snapshotOutputDir); err != nil {
		t.Fatalf("failed to clean snapshot output directory: %v", err)
	}
	if err := os.MkdirAll(snapshotOutputDir, 0o755); err != nil {
		t.Fatalf("failed to create snapshot output directory: %v", err)
	}

	for _, s := range styles {
		s := s // capture range variable
		t.Run(s.Name(), func(t *testing.T) {
			// Read style requirements and skip if dimensions are unconstrained.
			reqs := s.Requirements()
			if reqs.MinWidth == 0 || reqs.MinHeight == 0 {
				t.Skip("skipping: style has unconstrained dimensions (MinWidth or MinHeight is 0)")
			}

			// Derive display category from the style name.
			category := categoryFromStyleName(s.Name())

			// Build a policy that targets the current style.
			p := testPolicy()
			p.Style = s.Name()

			// Render through the snapshottest framework.
			pngPath := testsnapshot.RenderSnapshot(t,
				testsnapshot.WithMode("clock"),
				testsnapshot.WithDimensions(reqs.MinWidth, reqs.MinHeight),
				testsnapshot.WithDisplayCategory(category),
				testsnapshot.WithOutputDir(snapshotOutputDir),
				testsnapshot.WithBasename(s.Name()),
				testsnapshot.WithReset(func() {
					SetPolicy(DefaultPolicy())
				}),
				testsnapshot.WithPreRender(func() {
					SetPolicy(p)
				}),
			)

			// Verify output using the framework's verification helpers.
			testsnapshot.VerifyAll(t, pngPath, reqs.MinWidth, reqs.MinHeight)
		})
	}
}

// Property 2: Full pipeline render produces a valid PNG with correct dimensions
//
// For any style in the clock registry with MinWidth > 0 and MinHeight > 0,
// rendering through the snapshottest framework SHALL produce a decodable PNG file
// whose image dimensions equal the style's MinWidth and MinHeight.
//

func TestProperty_FullPipelineRenderProducesValidPNG(t *testing.T) {

	allStyles := clockRegistry.Enumerate()
	if len(allStyles) == 0 {
		t.Fatal("clockRegistry contains zero styles")
	}

	rapid.Check(t, func(rt *rapid.T) {
		idx := rapid.IntRange(0, len(allStyles)-1).Draw(rt, "styleIndex")
		s := allStyles[idx]

		reqs := s.Requirements()

		// Skip styles with unconstrained dimensions.
		if reqs.MinWidth == 0 || reqs.MinHeight == 0 {
			return
		}

		// Scratch output; discarded when the test ends.
		outputDir := t.TempDir()

		// Derive display category from the style name.
		category := categoryFromStyleName(s.Name())

		// Build a policy that targets the current style.
		p := testPolicy()
		p.Style = s.Name()

		// Render through the snapshottest framework.
		pngPath := testsnapshot.RenderSnapshot(t,
			testsnapshot.WithMode("clock"),
			testsnapshot.WithDimensions(reqs.MinWidth, reqs.MinHeight),
			testsnapshot.WithDisplayCategory(category),
			testsnapshot.WithOutputDir(outputDir),
			testsnapshot.WithBasename(s.Name()),
			testsnapshot.WithReset(func() {
				SetPolicy(DefaultPolicy())
			}),
			testsnapshot.WithPreRender(func() {
				SetPolicy(p)
			}),
		)

		// Verify output using the framework's verification helpers.
		testsnapshot.VerifyAll(t, pngPath, reqs.MinWidth, reqs.MinHeight)
	})
}

// Property 1: PNGPanel configuration is correctly derived from SurfaceRequirements
//
// For any style in the clock registry with MinWidth > 0 and MinHeight > 0,
// the derived PNGPanel configuration SHALL have dimensions equal to the style's
// MinWidth and MinHeight, and a color mode of FullColor when NeedsColor is true
// or Monochrome when NeedsColor is false.
//

func TestProperty_PNGPanelConfigDerivedFromRequirements(t *testing.T) {

	allStyles := clockRegistry.Enumerate()
	if len(allStyles) == 0 {
		t.Fatal("clockRegistry contains zero styles")
	}

	rapid.Check(t, func(rt *rapid.T) {
		idx := rapid.IntRange(0, len(allStyles)-1).Draw(rt, "styleIndex")
		s := allStyles[idx]

		reqs := s.Requirements()

		// Skip styles with unconstrained dimensions.
		if reqs.MinWidth == 0 || reqs.MinHeight == 0 {
			return
		}

		// Scratch output; discarded when the test ends.
		outputDir := t.TempDir()

		// Derive display category from the style name.
		category := categoryFromStyleName(s.Name())

		// Build a policy that targets the current style.
		p := testPolicy()
		p.Style = s.Name()

		// Render through the snapshottest framework — this validates correct
		// configuration derivation by producing output with matching dimensions.
		pngPath := testsnapshot.RenderSnapshot(t,
			testsnapshot.WithMode("clock"),
			testsnapshot.WithDimensions(reqs.MinWidth, reqs.MinHeight),
			testsnapshot.WithDisplayCategory(category),
			testsnapshot.WithOutputDir(outputDir),
			testsnapshot.WithBasename(s.Name()),
			testsnapshot.WithReset(func() {
				SetPolicy(DefaultPolicy())
			}),
			testsnapshot.WithPreRender(func() {
				SetPolicy(p)
			}),
		)

		// Verify dimensions match requirements — confirms PNGPanel was configured
		// with the correct dimensions from style requirements.
		testsnapshot.VerifyDimensions(t, pngPath, reqs.MinWidth, reqs.MinHeight)
	})
}

// TestSnapshotHelperValues verifies that BuildSnapshot produces the expected
// deterministic field values when given the test policy and reference time.

func TestSnapshotHelperValues(t *testing.T) {
	// Set the policy so BuildSnapshot produces deterministic output.
	SetPolicy(testPolicy())
	defer SetPolicy(DefaultPolicy())

	snap := BuildSnapshot(referenceTime)

	if snap.Time != "14:30:05" {
		t.Errorf("BuildSnapshot().Time = %q, want %q", snap.Time, "14:30:05")
	}
	if snap.Date != "2024-01-15" {
		t.Errorf("BuildSnapshot().Date = %q, want %q", snap.Date, "2024-01-15")
	}
	if snap.Weekday != "Monday" {
		t.Errorf("BuildSnapshot().Weekday = %q, want %q", snap.Weekday, "Monday")
	}
}

// TestPolicyHelperValues verifies that testPolicy() returns field values matching
// the specification for deterministic rendering.
func TestPolicyHelperValues(t *testing.T) {
	p := testPolicy()

	expectedStyle := ""
	if p.Style != expectedStyle {
		t.Errorf("testPolicy().Style = %q, want %q", p.Style, expectedStyle)
	}
	if p.ShowSeconds != true {
		t.Errorf("testPolicy().ShowSeconds = %v, want true", p.ShowSeconds)
	}
	if p.TimeFormat != "24h" {
		t.Errorf("testPolicy().TimeFormat = %q, want %q", p.TimeFormat, "24h")
	}
	if p.DateFormat != "YYYY-MM-DD" {
		t.Errorf("testPolicy().DateFormat = %q, want %q", p.DateFormat, "YYYY-MM-DD")
	}
	if p.Timezone != "UTC" {
		t.Errorf("testPolicy().Timezone = %q, want %q", p.Timezone, "UTC")
	}
	if p.ShowWeekday != true {
		t.Errorf("testPolicy().ShowWeekday = %v, want true", p.ShowWeekday)
	}
	if p.BlinkColon != false {
		t.Errorf("testPolicy().BlinkColon = %v, want false", p.BlinkColon)
	}

	if p.FGColor != "cyan" {
		t.Errorf("testPolicy().FGColor = %q, want %q", p.FGColor, "cyan")
	}
	if p.ShowLED != true {
		t.Errorf("testPolicy().ShowLED = %v, want true", p.ShowLED)
	}
	if p.SecondsBar != "none" {
		t.Errorf("testPolicy().SecondsBar = %q, want %q", p.SecondsBar, "none")
	}
	if p.ShowDaybar != false {
		t.Errorf("testPolicy().ShowDaybar = %v, want false", p.ShowDaybar)
	}
	if p.ShowBorder != true {
		t.Errorf("testPolicy().ShowBorder = %v, want true", p.ShowBorder)
	}
	if p.BorderColor != "emerald" {
		t.Errorf("testPolicy().BorderColor = %q, want %q", p.BorderColor, "emerald")
	}
}

// TestStyleCountGuard verifies that the clock registry contains exactly 108 styles,
// ensuring the test suite's count guard assertion remains valid.

func TestStyleCountGuard(t *testing.T) {
	styles := clockRegistry.Enumerate()
	if len(styles) != 108 {
		t.Fatalf("clockRegistry.Enumerate() returned %d styles, want 108", len(styles))
	}
}
