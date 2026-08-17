package thermal_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/databeast/cyberhud/display/modes/testsnapshot"
	"github.com/databeast/cyberhud/display/modes/thermal"
	"github.com/databeast/cyberhud/display/modes/thermal/tests"
)

// snapshotOutputDir is the persistent directory where thermal snapshot PNGs are written.
var snapshotOutputDir = filepath.Join("snapshots")

// testThermalSnapshot returns a deterministic ThermalSnapshot for snapshot rendering.
func testThermalSnapshot() thermal.ThermalSnapshot {
	return thermal.ThermalSnapshot{
		Zones: []thermal.ZoneReading{
			{ZoneID: 0, Label: "cpu-thermal", TempC: 52.3, TripPoints: []thermal.TripPoint{{Type: "critical", TempC: 100.0}}},
			{ZoneID: 1, Label: "gpu-thermal", TempC: 47.8, TripPoints: []thermal.TripPoint{{Type: "critical", TempC: 95.0}}},
			{ZoneID: 2, Label: "pmic-thermal", TempC: 38.1, TripPoints: []thermal.TripPoint{{Type: "critical", TempC: 120.0}}},
		},
		Timestamp: time.Date(2024, 1, 15, 14, 30, 5, 0, time.UTC),
	}
}

// displayCategory represents a display category name, mirroring the gallery
// pipeline's category.Category type for prefix classification.
type displayCategory string

const (
	catColor     displayCategory = "Color"
	catEInk      displayCategory = "E-Ink"
	catGrayscale displayCategory = "Grayscale"
	catMono      displayCategory = "Mono"
)

// prefixMapping associates a filename prefix with a display category.
type prefixMapping struct {
	prefix   string
	category displayCategory
}

// recognizedPrefixes mirrors category.Prefixes — ordered longest-first for
// correct longest-prefix matching. This must stay in sync with
// buildtools/docsnap/internal/category.Prefixes.
var recognizedPrefixes = []prefixMapping{
	{prefix: "color-slow-", category: catColor},
	{prefix: "color-fast-", category: catColor},
	{prefix: "grayscale-fast-", category: catGrayscale},
	{prefix: "grayscale-slow-", category: catGrayscale},
	{prefix: "mono-slow-", category: catMono},
	{prefix: "mono-fast-", category: catMono},
	{prefix: "color-", category: catColor},
	{prefix: "mono-", category: catMono},
	{prefix: "eink-", category: catEInk},
}

// matchCategory returns the display category for a style name using
// longest-prefix matching, mirroring category.Match from the gallery pipeline.
func matchCategory(name string) (displayCategory, bool) {
	var bestCat displayCategory
	bestLen := 0
	for _, pm := range recognizedPrefixes {
		if strings.HasPrefix(name, pm.prefix) && len(pm.prefix) > bestLen {
			bestCat = pm.category
			bestLen = len(pm.prefix)
		}
	}
	if bestLen == 0 {
		return "", false
	}
	return bestCat, true
}

// categoryToDisplayCategory maps a local displayCategory to testsnapshot.DisplayCategory.
func categoryToDisplayCategory(cat displayCategory) testsnapshot.DisplayCategory {
	switch cat {
	case catColor:
		return testsnapshot.CategoryColor
	case catMono:
		return testsnapshot.CategoryMono
	case catEInk:
		return testsnapshot.CategoryEink
	case catGrayscale:
		return testsnapshot.CategoryGrayscale
	default:
		return testsnapshot.CategoryColor
	}
}

// TestThermalPNGSnapshots enumerates all registered thermal styles and generates
// a snapshot PNG for each via the full production render pipeline using the
// snapshottest framework.
func TestThermalPNGSnapshots(t *testing.T) {
	styles := tests.ThermalRegistry.Enumerate()
	if len(styles) == 0 {
		t.Fatal("thermalRegistry contains zero styles")
	}

	// Pre-flight check: verify every style name is recognized by matchCategory
	// (mirrors category.Match from the gallery pipeline).
	var prefixList []string
	for _, pm := range recognizedPrefixes {
		prefixList = append(prefixList, pm.prefix)
	}
	knownPrefixes := strings.Join(prefixList, ", ")
	for _, s := range styles {
		if _, ok := matchCategory(s.Name()); !ok {
			t.Fatalf("style %q does not match any recognized category prefix; recognized prefixes: [%s]",
				s.Name(), knownPrefixes)
		}
	}

	// Ensure the output directory exists and is clean.
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

			// Use preferred dimensions for a representative screenshot; fall back
			// to minimum if preferred is unset.
			width, height := reqs.PreferredWidth, reqs.PreferredHeight
			if width == 0 {
				width = reqs.MinWidth
			}
			if height == 0 {
				height = reqs.MinHeight
			}
			if width == 0 || height == 0 {
				t.Skip("skipping: style has unconstrained dimensions")
			}

			// Derive display category from style name via matchCategory
			// (mirrors category.Match from the gallery pipeline).
			cat, ok := matchCategory(s.Name())
			if !ok {
				t.Fatalf("style %q: matchCategory returned no match", s.Name())
			}
			displayCat := categoryToDisplayCategory(cat)

			// Build a policy that targets the current style.
			p := thermal.DefaultPolicy()
			p.Style = s.Name()

			// Seed deterministic snapshot and history data.
			snap := testThermalSnapshot()

			// Render through the snapshottest framework.
			pngPath := testsnapshot.RenderSnapshot(t,
				testsnapshot.WithMode("thermal"),
				testsnapshot.WithDimensions(width, height),
				testsnapshot.WithDisplayCategory(displayCat),
				testsnapshot.WithOutputDir(snapshotOutputDir),
				testsnapshot.WithBasename(s.Name()),
				testsnapshot.WithReset(func() {
					_ = thermal.SetPolicy(thermal.DefaultPolicy())
					tests.ResetHistoryState()
					tests.ResetSnapshotState()
				}),
				testsnapshot.WithPreRender(func() {
					_ = thermal.SetPolicy(p)
					thermal.UpdateSnapshot(snap)
					for _, z := range snap.Zones {
						for i := 0; i < 64; i++ {
							thermal.RecordHistory(z.ZoneID, z.TempC+float64(i%5)*0.1)
						}
					}
				}),
			)

			// Verify output using the framework's verification helpers.
			testsnapshot.VerifyAll(t, pngPath, width, height)
		})
	}
}
