package doccheck

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestPanelPageTemplate validates Property 4: Panel Page Template Conformance.
// For each registered panel, it verifies that the documentation page contains
// required sections in the prescribed order, and that the "Input Details"
// section is present only when the panel defines InputPins.
func TestPanelPageTemplate(t *testing.T) {
	projectRoot := filepath.Join("..", "..")
	docsDir := filepath.Join(projectRoot, "ghpages", "docs", "getting-started")

	// Discover all registered panels from source.
	registeredPanels := discoverRegisteredPanels(t, projectRoot)
	t.Logf("discovered %d registered panels from source", len(registeredPanels))

	// Required sections that must appear in order.
	requiredSections := []string{
		"Quick Start",
		"Display Characteristics",
		"Troubleshooting",
		"Related Pages",
	}

	// Panels known to have InputPins defined in source.
	panelsWithInputs := map[string]bool{
		"waveshare-1.3-oled-hat":  true,
		"waveshare-1.3hat":        true,
		"waveshare-1.44":          true,
		"waveshare-triple-screen": true,
		"adafruit-2.13-ssd1680":   true,
	}

	// Track which panels have already been validated (for grouped panels
	// sharing a page, avoid testing the same file twice).
	validated := make(map[string]bool)

	for panelName := range registeredPanels {
		docFile := panelDocFilename(panelName)

		// Skip if we already validated this doc file (grouped panels).
		if validated[docFile] {
			continue
		}
		validated[docFile] = true

		fullPath := filepath.Join(docsDir, docFile)

		t.Run(panelName, func(t *testing.T) {
			// Verify the file exists.
			if _, err := os.Stat(fullPath); os.IsNotExist(err) {
				t.Fatalf("documentation page does not exist: %s", fullPath)
			}

			// Extract H2 headings.
			headings := extractHeadings(fullPath)
			if len(headings) == 0 {
				t.Fatalf("no H2 headings found in %s", docFile)
			}

			// Verify required sections are present and in order.
			lastIndex := -1
			for _, required := range requiredSections {
				idx := indexOfHeading(headings, required)
				if idx == -1 {
					t.Errorf("required section %q not found in %s", required, docFile)
					continue
				}
				if idx <= lastIndex {
					t.Errorf("section %q (index %d) is out of order in %s; expected after index %d",
						required, idx, docFile, lastIndex)
				}
				lastIndex = idx
			}

			// Verify conditional Input Details section.
			// Accept "Input Details" or "Buttons & Joystick" as the input section heading.
			hasInputDetails := indexOfHeading(headings, "Input Details") != -1 ||
				indexOfHeading(headings, "Buttons & Joystick") != -1
			expectInputs := panelsWithInputs[panelName]

			if expectInputs && !hasInputDetails {
				t.Errorf("panel %q has InputPins but page %s is missing \"Input Details\" section",
					panelName, docFile)
			}
			if !expectInputs && hasInputDetails {
				t.Errorf("panel %q has no InputPins but page %s contains \"Input Details\" section",
					panelName, docFile)
			}
		})
	}
}

// indexOfHeading returns the index of the first heading in the slice that
// matches the target (case-insensitive). Returns -1 if not found.
func indexOfHeading(headings []string, target string) int {
	for i, h := range headings {
		if strings.EqualFold(h, target) {
			return i
		}
	}
	return -1
}
