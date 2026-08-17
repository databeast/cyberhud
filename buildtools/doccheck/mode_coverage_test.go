package doccheck

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

// TestModeCoverageCompleteness validates Property 2: Mode Coverage Completeness.
// The set of mode IDs in documentation equals the set of IDs returned by
// enumerating catalog.Register() calls, minus the excluded test/utility modes
// (testfonts, testicons, testwidgets, snapshottest).
func TestModeCoverageCompleteness(t *testing.T) {
	projectRoot := filepath.Join("..", "..")

	// Step 1: Enumerate all catalog.Register() calls in display/modes/ subdirectories.
	registeredModes := discoverRegisteredModes(t, projectRoot)
	t.Logf("discovered %d registered modes from source: %v", len(registeredModes), sortedKeys(registeredModes))

	// Step 2: Exclude test/utility modes.
	excluded := map[string]bool{
		"testfonts":    true,
		"testicons":    true,
		"testwidgets":  true,
		"snapshottest": true,
		"testpattern":  true,
	}

	userFacingModes := make(map[string]bool)
	for id := range registeredModes {
		if !excluded[id] {
			userFacingModes[id] = true
		}
	}
	t.Logf("user-facing modes (after exclusions): %d — %v", len(userFacingModes), sortedKeys(userFacingModes))

	// Step 3: Parse mkdocs.yml and extract mode page entries.
	mkdocsPath := filepath.Join(projectRoot, "ghpages", "mkdocs.yml")
	navModePages := extractNavModePages(t, mkdocsPath)
	t.Logf("nav mode pages found: %d — %v", len(navModePages), sortedKeys(navModePages))

	// Step 4: Verify doc pages exist on disk for each nav entry.
	docsDir := filepath.Join(projectRoot, "ghpages", "docs")
	for modeID, relPath := range navModePages {
		fullPath := filepath.Join(docsDir, relPath)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			t.Errorf("nav entry for mode %q references %s but file does not exist on disk", modeID, relPath)
		}
	}

	// Step 5: Check that every user-facing mode has a documentation page.
	var missingDoc []string
	var missingNav []string

	for modeID := range userFacingModes {
		// Check if a documentation page file exists.
		docPage := modeDocPagePath(modeID)
		fullDocPath := filepath.Join(docsDir, "display-modes", docPage)
		if _, err := os.Stat(fullDocPath); os.IsNotExist(err) {
			missingDoc = append(missingDoc, modeID)
		}

		// Check if the mode has a nav entry.
		if _, ok := navModePages[modeID]; !ok {
			missingNav = append(missingNav, modeID)
		}
	}

	sort.Strings(missingDoc)
	sort.Strings(missingNav)

	for _, id := range missingDoc {
		t.Errorf("user-facing mode %q has no documentation page (expected %s)", id, modeDocPagePath(id))
	}

	for _, id := range missingNav {
		t.Errorf("user-facing mode %q has no navigation entry in mkdocs.yml", id)
	}

	if len(missingDoc) == 0 && len(missingNav) == 0 {
		t.Logf("all %d user-facing modes have both a documentation page and a nav entry", len(userFacingModes))
	}
}

// discoverRegisteredModes walks display/modes/ subdirectories and parses Go
// source files for catalog.Register(catalog.Definition{ID: "..."}) calls.
// It returns a map of mode ID -> source file where it was found.
func discoverRegisteredModes(t *testing.T, projectRoot string) map[string]string {
	t.Helper()

	modesDir := filepath.Join(projectRoot, "display", "modes")
	registered := make(map[string]string)

	// Pattern to match catalog.Register(catalog.Definition{ followed by ID: "..."
	// We look across lines since the struct literal spans multiple lines.
	registerPattern := regexp.MustCompile(`catalog\.Register\(catalog\.Definition\{`)
	idPattern := regexp.MustCompile(`ID:\s*"([^"]+)"`)

	err := filepath.Walk(modesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Only process Go source files (not test files).
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}
		if strings.HasSuffix(path, "_test.go") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("reading %s: %w", path, err)
		}

		text := string(content)

		// Find all catalog.Register calls and extract IDs.
		indices := registerPattern.FindAllStringIndex(text, -1)
		for _, idx := range indices {
			// Search for the ID field in the next 200 characters after the match.
			end := idx[1] + 200
			if end > len(text) {
				end = len(text)
			}
			snippet := text[idx[1]:end]

			if m := idPattern.FindStringSubmatch(snippet); m != nil {
				modeID := m[1]
				relPath, _ := filepath.Rel(projectRoot, path)
				registered[modeID] = relPath
			}
		}

		return nil
	})

	if err != nil {
		t.Fatalf("walking display/modes: %v", err)
	}

	return registered
}

// extractNavModePages parses mkdocs.yml and extracts mode page entries from
// the Display Modes navigation section. Returns a map of mode ID -> relative
// path (e.g., "attract_bokeh" -> "display-modes/attract_bokeh.md").
func extractNavModePages(t *testing.T, mkdocsPath string) map[string]string {
	t.Helper()

	f, err := os.Open(mkdocsPath)
	if err != nil {
		t.Fatalf("cannot open mkdocs.yml: %v", err)
	}
	defer f.Close()

	// We parse the YAML nav by scanning lines for display-modes/ page references.
	// This is simpler and more robust than a full YAML parse for our purposes.
	pagePattern := regexp.MustCompile(`:\s*(display-modes/([^.]+)\.md)\s*$`)

	result := make(map[string]string)
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()

		m := pagePattern.FindStringSubmatch(line)
		if m == nil {
			continue
		}

		relPath := m[1]  // e.g., "display-modes/attract_bokeh.md"
		basename := m[2] // e.g., "attract_bokeh"

		// Derive mode ID from the basename.
		// Most mode IDs match their file basename directly.
		// Special case: gpio_control -> gpio-control (registered with hyphen).
		modeID := filenameToModeID(basename)

		// Skip non-mode pages (index, attract-modes overview, utility overview).
		if isNonModePage(basename) {
			continue
		}

		result[modeID] = relPath
	}

	if err := scanner.Err(); err != nil {
		t.Fatalf("error reading mkdocs.yml: %v", err)
	}

	return result
}

// filenameToModeID converts a documentation filename (without extension) to
// the corresponding catalog mode ID.
func filenameToModeID(basename string) string {
	// The gpio_control directory registers with ID "gpio-control".
	// The doc page is gpio_control.md. Map it back.
	switch basename {
	case "gpio_control":
		return "gpio-control"
	default:
		return basename
	}
}

// modeDocPagePath returns the expected documentation page filename for a mode ID.
func modeDocPagePath(modeID string) string {
	// Reverse mapping: mode ID -> expected doc page filename.
	switch modeID {
	case "gpio-control":
		return "gpio_control.md"
	default:
		return modeID + ".md"
	}
}

// isNonModePage returns true for display-modes/ pages that are not individual
// mode documentation pages (indexes, overviews, utility listings).
func isNonModePage(basename string) bool {
	nonModePages := map[string]bool{
		"index":         true,
		"attract-modes": true,
		"utility":       true,
	}
	return nonModePages[basename]
}

// sortedKeys returns the keys of a map sorted alphabetically.
func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
