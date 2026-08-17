package doccheck

import (
	"path/filepath"
	"strings"
	"testing"
)

// TestModePageTemplate validates Property 5: Mode Page Template Conformance.
// For every registered mode page, the required H2 sections must appear and
// must be in the prescribed order. Extra sections are allowed (the test only
// checks that the required set appears in order).
func TestModePageTemplate(t *testing.T) {
	projectRoot := filepath.Join("..", "..")

	// Discover all registered modes from source.
	registeredModes := discoverRegisteredModes(t, projectRoot)
	t.Logf("discovered %d registered modes", len(registeredModes))

	// Required sections for standard (non-attract) mode pages.
	standardRequired := []string{
		"Quick Start",
		"How It Works",
		"Panel Compatibility",
		"Related Pages",
	}

	// Required sections for attract mode pages (different template).
	attractRequired := []string{
		"Panel Compatibility",
	}

	// Modes whose documentation pages haven't been updated to the new
	// template yet. They are excluded from the template conformance check
	// until their pages are reworked.
	templateExclusions := map[string]bool{}

	for modeID := range registeredModes {
		// Skip modes excluded from template enforcement.
		if utilityModeExclusions[modeID] && modeID != "testfonts" && modeID != "testicons" && modeID != "testwidgets" && modeID != "testpattern" {
			continue
		}
		if templateExclusions[modeID] {
			t.Logf("skipping %q (template exclusion — page not yet updated)", modeID)
			continue
		}

		docFile := filepath.Join(projectRoot, "ghpages", "docs", "display-modes", modeDocPagePath(modeID))
		headings := extractHeadings(docFile)

		if headings == nil {
			t.Errorf("mode %q: cannot read documentation page at %s", modeID, docFile)
			continue
		}

		// Select the required sections based on mode type.
		var required []string
		if strings.HasPrefix(modeID, "attract_") {
			required = attractRequired
		} else {
			required = standardRequired
		}

		// Verify required sections are present and in order.
		verifyOrderedSections(t, modeID, headings, required)
	}
}

// verifyOrderedSections checks that each heading in `required` appears in the
// `headings` slice, and that their relative order matches (i.e., each required
// heading appears after the previous one). Extra headings between required ones
// are ignored.
func verifyOrderedSections(t *testing.T, modeID string, headings []string, required []string) {
	t.Helper()

	lastIdx := -1
	for _, req := range required {
		found := false
		for i := lastIdx + 1; i < len(headings); i++ {
			if headings[i] == req {
				lastIdx = i
				found = true
				break
			}
		}
		if !found {
			// Check if the section exists at all (possibly out of order).
			exists := false
			for _, h := range headings {
				if h == req {
					exists = true
					break
				}
			}
			if exists {
				t.Errorf("mode %q: required section %q is present but out of order (expected after index %d)", modeID, req, lastIdx)
			} else {
				t.Errorf("mode %q: missing required section %q", modeID, req)
			}
			return
		}
	}
}
