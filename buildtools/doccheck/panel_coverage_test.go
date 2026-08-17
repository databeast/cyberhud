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

// TestPanelCoverageCompleteness validates that every registered hardware panel
// has a corresponding documentation page and navigation entry in mkdocs.yml.
func TestPanelCoverageCompleteness(t *testing.T) {
	projectRoot := filepath.Join("..", "..")

	// Step 1: Discover all panels.Register() calls in hardware/panels/.
	registeredPanels := discoverRegisteredPanels(t, projectRoot)
	t.Logf("discovered %d registered panels from source: %v", len(registeredPanels), sortedKeys(registeredPanels))

	// Step 2: Verify each panel has a documentation page.
	docsDir := filepath.Join(projectRoot, "ghpages", "docs", "getting-started")

	var missingDoc []string
	for panelName := range registeredPanels {
		expectedFile := panelDocFilename(panelName)
		fullPath := filepath.Join(docsDir, expectedFile)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			missingDoc = append(missingDoc, panelName)
		}
	}

	// Step 3: Parse mkdocs.yml nav and verify each panel has a nav entry.
	mkdocsPath := filepath.Join(projectRoot, "ghpages", "mkdocs.yml")
	navPanelPages := extractNavPanelPages(t, mkdocsPath)
	t.Logf("nav panel pages found: %d — %v", len(navPanelPages), sortedKeys(navPanelPages))

	var missingNav []string
	for panelName := range registeredPanels {
		expectedFile := panelDocFilename(panelName)
		navEntry := "getting-started/" + expectedFile
		if !navPanelPages[navEntry] {
			missingNav = append(missingNav, panelName)
		}
	}

	// Step 4: Report results.
	sort.Strings(missingDoc)
	sort.Strings(missingNav)

	for _, name := range missingDoc {
		t.Errorf("registered panel %q has no documentation page (expected getting-started/%s)", name, panelDocFilename(name))
	}

	for _, name := range missingNav {
		t.Errorf("registered panel %q has no navigation entry in mkdocs.yml (expected getting-started/%s)", name, panelDocFilename(name))
	}

	if len(missingDoc) == 0 && len(missingNav) == 0 {
		t.Logf("all %d registered panels have both a documentation page and a nav entry", len(registeredPanels))
	}
}

// discoverRegisteredPanels walks hardware/panels/ subdirectories and parses Go
// source files for panels.Register(panels.Definition{Name: "..."}) calls.
// It returns a map of panel name -> source file where it was found.
func discoverRegisteredPanels(t *testing.T, projectRoot string) map[string]string {
	t.Helper()

	panelsDir := filepath.Join(projectRoot, "hardware", "panels")
	registered := make(map[string]string)

	// Pattern to match panels.Register(panels.Definition{ followed by Name: "..."
	registerPattern := regexp.MustCompile(`panels\.Register\(panels\.Definition\{`)
	namePattern := regexp.MustCompile(`Name:\s*"([^"]+)"`)

	err := filepath.Walk(panelsDir, func(path string, info os.FileInfo, err error) error {
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

		// Find all panels.Register calls and extract Names.
		indices := registerPattern.FindAllStringIndex(text, -1)
		for _, idx := range indices {
			// Search for the Name field in the next 200 characters after the match.
			end := idx[1] + 200
			if end > len(text) {
				end = len(text)
			}
			snippet := text[idx[1]:end]

			if m := namePattern.FindStringSubmatch(snippet); m != nil {
				panelName := m[1]
				relPath, _ := filepath.Rel(projectRoot, path)
				registered[panelName] = relPath
			}
		}

		return nil
	})

	if err != nil {
		t.Fatalf("walking hardware/panels: %v", err)
	}

	return registered
}

// extractNavPanelPages parses mkdocs.yml and extracts panel page entries from
// the Getting Started navigation section. Returns a set of relative paths
// (e.g., "getting-started/waveshare-2-23-oled-hat.md").
func extractNavPanelPages(t *testing.T, mkdocsPath string) map[string]bool {
	t.Helper()

	f, err := os.Open(mkdocsPath)
	if err != nil {
		t.Fatalf("cannot open mkdocs.yml: %v", err)
	}
	defer f.Close()

	// Match lines referencing getting-started/*.md pages.
	pagePattern := regexp.MustCompile(`(getting-started/[^\s]+\.md)`)

	result := make(map[string]bool)
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()

		m := pagePattern.FindStringSubmatch(line)
		if m == nil {
			continue
		}

		result[m[1]] = true
	}

	if err := scanner.Err(); err != nil {
		t.Fatalf("error reading mkdocs.yml: %v", err)
	}

	return result
}

// TestUnclassifiedDirectoryDetection validates that every subdirectory under
// display/modes/ is either a registered mode (contains a catalog.Register call)
// or is listed in stubDirectoryExclusions. Directories that are neither are
// reported as unclassified — they need to either register a mode or be added
// to the exclusion list.
func TestUnclassifiedDirectoryDetection(t *testing.T) {
	projectRoot := filepath.Join("..", "..")
	modesDir := filepath.Join(projectRoot, "display", "modes")

	// Step 1: Discover all registered modes (directories with catalog.Register calls).
	registeredModes := discoverRegisteredModes(t, projectRoot)
	registeredDirs := make(map[string]bool)
	for _, srcPath := range registeredModes {
		// srcPath is like "display/modes/clock/mode.go" — extract the directory name.
		rel, err := filepath.Rel(filepath.Join("display", "modes"), filepath.Dir(srcPath))
		if err != nil {
			continue
		}
		// rel could be nested (e.g. "attract_bokeh"), take the top-level directory.
		parts := strings.SplitN(filepath.ToSlash(rel), "/", 2)
		registeredDirs[parts[0]] = true
	}
	t.Logf("directories with registered modes: %v", sortedKeys(registeredDirs))

	// Step 2: Read immediate subdirectories of display/modes/.
	entries, err := os.ReadDir(modesDir)
	if err != nil {
		t.Fatalf("cannot read display/modes/ directory: %v", err)
	}

	registerPattern := regexp.MustCompile(`catalog\.Register\(`)

	var unclassified []string

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		dirName := entry.Name()

		// Skip if this directory has a registered mode.
		if registeredDirs[dirName] {
			continue
		}

		// Skip if it's in the exclusion list.
		if stubDirectoryExclusions[dirName] {
			continue
		}

		// Check if any non-test .go file in this directory contains catalog.Register(.
		dirPath := filepath.Join(modesDir, dirName)
		files, err := os.ReadDir(dirPath)
		if err != nil {
			t.Errorf("cannot read directory %s: %v", dirName, err)
			continue
		}

		hasGoFiles := false
		hasRegisterCall := false

		for _, f := range files {
			if f.IsDir() {
				continue
			}
			name := f.Name()
			if !strings.HasSuffix(name, ".go") {
				continue
			}
			if strings.HasSuffix(name, "_test.go") {
				continue
			}

			hasGoFiles = true

			content, err := os.ReadFile(filepath.Join(dirPath, name))
			if err != nil {
				t.Errorf("cannot read file %s/%s: %v", dirName, name, err)
				continue
			}

			if registerPattern.Match(content) {
				hasRegisterCall = true
				break
			}
		}

		// Skip directories with no Go source files (e.g., img/ or empty dirs).
		if !hasGoFiles {
			continue
		}

		// If it has a register call, it was missed by discoverRegisteredModes
		// (perhaps a different pattern) — treat as registered.
		if hasRegisterCall {
			continue
		}

		unclassified = append(unclassified, dirName)
	}

	sort.Strings(unclassified)

	for _, dir := range unclassified {
		t.Errorf("display/modes/%s/ is unclassified: contains Go source files but no catalog.Register() call and is not in stubDirectoryExclusions. "+
			"Either add a catalog.Register() call or add %q to stubDirectoryExclusions in panel_mapping.go.", dir, dir)
	}

	if len(unclassified) == 0 {
		t.Logf("all display/modes/ subdirectories are classified (registered or excluded)")
	}
}
