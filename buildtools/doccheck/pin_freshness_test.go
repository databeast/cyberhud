package doccheck

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

// panelEntry pairs a panel name with its source file path.
type panelEntry struct {
	name    string
	srcPath string
}

// TestPinTableFreshness validates that GPIO pin assignments declared in panel
// source code are accurately reflected in the corresponding documentation page's
// Pin Assignments table. Any source-defined pin not found in the docs is reported.
func TestPinTableFreshness(t *testing.T) {
	projectRoot := filepath.Join("..", "..")

	// Step 1: Discover all registered panels.
	registeredPanels := discoverRegisteredPanels(t, projectRoot)
	t.Logf("discovered %d registered panels", len(registeredPanels))

	docsDir := filepath.Join(projectRoot, "ghpages", "docs", "getting-started")

	// Group panels by their documentation file. Multiple panel variants (e.g.,
	// SPI and I2C) may share a single doc page and source file.
	docGroups := make(map[string][]panelEntry)

	for panelName, srcRelPath := range registeredPanels {
		docFilename := panelDocFilename(panelName)
		srcPath := filepath.Join(projectRoot, srcRelPath)
		docGroups[docFilename] = append(docGroups[docFilename], panelEntry{
			name:    panelName,
			srcPath: srcPath,
		})
	}

	docFilenames := sortedKeys(docGroups)
	for _, docFilename := range docFilenames {
		entries := docGroups[docFilename]
		docPath := filepath.Join(docsDir, docFilename)

		// Collect the union of source pins across all variants sharing this doc.
		sourcePinSet := make(map[string]bool)
		for _, entry := range entries {
			for _, pin := range extractSourcePins(t, entry.srcPath) {
				sourcePinSet[pin] = true
			}
		}

		if len(sourcePinSet) == 0 {
			t.Logf("skipping %s (no pin assignments in source for %v)", docFilename, panelEntryNames(entries))
			continue
		}

		// Verify the doc page exists.
		if _, err := os.Stat(docPath); os.IsNotExist(err) {
			t.Errorf("doc page %s does not exist (cannot validate pins for %v)", docFilename, panelEntryNames(entries))
			continue
		}

		// Extract documented GPIO pins from the Pin Assignments section.
		docPins := extractDocPins(t, docPath)
		if len(docPins) == 0 {
			t.Errorf("doc page %s has no Pin Assignments table or no GPIO entries (panels: %v)", docFilename, panelEntryNames(entries))
			continue
		}

		// Compare: every source pin should appear in docs.
		var missing []string
		for pin := range sourcePinSet {
			if !docPins[pin] {
				missing = append(missing, pin)
			}
		}

		sort.Strings(missing)
		for _, pin := range missing {
			t.Errorf("doc %s: source pin %s not found in Pin Assignments table (panels: %v)", docFilename, pin, panelEntryNames(entries))
		}

		if len(missing) == 0 {
			t.Logf("doc %s: all %d source pins found (panels: %v)", docFilename, len(sourcePinSet), panelEntryNames(entries))
		}
	}
}

// extractSourcePins reads a panel source file and extracts all GPIO pin
// identifiers from display pin fields (DCPin, RSTPin, BLPin, BusyPin) and
// input pin fields (Key1, Key2, Key3, JoyUp, JoyDown, JoyLeft, JoyRight,
// JoyPressed). Returns a deduplicated sorted slice of GPIO identifiers.
func extractSourcePins(t *testing.T, srcPath string) []string {
	t.Helper()

	content, err := os.ReadFile(srcPath)
	if err != nil {
		t.Fatalf("cannot read source file %s: %v", srcPath, err)
	}

	text := string(content)

	// Patterns for display pins and input pins referencing panels.GPIOxx constants.
	displayPinRe := regexp.MustCompile(`(?:DCPin|RSTPin|BLPin|BusyPin):\s*panels\.(GPIO\d+)`)
	inputPinRe := regexp.MustCompile(`(?:Key1|Key2|Key3|JoyUp|JoyDown|JoyLeft|JoyRight|JoyPressed):\s*panels\.(GPIO\d+)`)

	pinSet := make(map[string]bool)

	for _, m := range displayPinRe.FindAllStringSubmatch(text, -1) {
		pinSet[m[1]] = true
	}
	for _, m := range inputPinRe.FindAllStringSubmatch(text, -1) {
		pinSet[m[1]] = true
	}

	pins := make([]string, 0, len(pinSet))
	for pin := range pinSet {
		pins = append(pins, pin)
	}
	sort.Strings(pins)
	return pins
}

// extractDocPins reads a panel documentation page and extracts all GPIO pin
// references from the Pin Assignments section (and subsections under it).
// Returns a set of GPIO identifiers found in the table(s).
func extractDocPins(t *testing.T, docPath string) map[string]bool {
	t.Helper()

	// First try extracting from a flat "Pin Assignments" H2 section.
	rows := extractTableRows(docPath, "Pin Assignments")

	// Also scan for sub-headed pin tables (H3 sections under Pin Assignments).
	// For multi-section pin tables, we parse the entire file looking for GPIO
	// references in any table row within the Pin Assignments section.
	pinSet := make(map[string]bool)
	gpioRe := regexp.MustCompile(`GPIO\d+`)

	// Extract from standard H2 table rows.
	for _, row := range rows {
		for _, cell := range row {
			if m := gpioRe.FindString(cell); m != "" {
				pinSet[m] = true
			}
		}
	}

	// For docs with H3 subsections (e.g., waveshare-triple-screen with separate
	// screen sections, or waveshare-2-23-oled-hat with SPI/I2C modes), we need
	// to also scan tables under H3 headings within the Pin Assignments area.
	// Re-read the file and parse all tables between "## Pin Assignments" and
	// the next H2 heading.
	pinSet = extractAllPinSectionGPIOs(docPath, pinSet)

	return pinSet
}

// extractAllPinSectionGPIOs reads the doc file and extracts all GPIO references
// from table rows between "## Pin Assignments" and the next "## " heading.
// This handles H3-subdivided pin sections.
func extractAllPinSectionGPIOs(docPath string, existing map[string]bool) map[string]bool {
	content, err := os.ReadFile(docPath)
	if err != nil {
		return existing
	}

	gpioRe := regexp.MustCompile(`GPIO\d+`)
	lines := strings.Split(string(content), "\n")
	inPinSection := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Detect the Pin Assignments H2 heading.
		if strings.HasPrefix(line, "## ") {
			heading := strings.TrimPrefix(line, "## ")
			if strings.TrimSpace(heading) == "Pin Assignments" {
				inPinSection = true
				continue
			}
			// Hit the next H2 — stop.
			if inPinSection {
				break
			}
			continue
		}

		if !inPinSection {
			continue
		}

		// Process table rows within the pin section.
		if !strings.HasPrefix(trimmed, "|") {
			continue
		}
		if isSeparatorRow(trimmed) {
			continue
		}

		// Extract GPIO references from this row.
		for _, m := range gpioRe.FindAllString(trimmed, -1) {
			existing[m] = true
		}
	}

	return existing
}

// panelEntryNames extracts panel names from a slice of panelEntry for logging.
func panelEntryNames(entries []panelEntry) []string {
	names := make([]string, len(entries))
	for i, e := range entries {
		names[i] = e.name
	}
	sort.Strings(names)
	return names
}
