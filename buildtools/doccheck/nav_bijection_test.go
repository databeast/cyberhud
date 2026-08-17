package doccheck

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// TestNavFileBijection validates Property 1: Navigation–File Bijection.
// For every entry in mkdocs.yml nav, a corresponding markdown file exists at
// the referenced path under docs/. Conversely, no orphan documentation page
// exists without a nav entry.
func TestNavFileBijection(t *testing.T) {
	ghpagesRoot := filepath.Join("..", "..", "ghpages")

	mkdocsPath := filepath.Join(ghpagesRoot, "mkdocs.yml")
	docsRoot := filepath.Join(ghpagesRoot, "docs")

	// Step 1: Parse mkdocs.yml to extract all nav file paths.
	navPaths := extractNavPaths(t, mkdocsPath)
	t.Logf("found %d nav entries in mkdocs.yml", len(navPaths))

	// Step 2: For each nav entry, verify the file exists under docs/.
	var missingFiles []string
	for _, relPath := range navPaths {
		fullPath := filepath.Join(docsRoot, filepath.FromSlash(relPath))
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			missingFiles = append(missingFiles, relPath)
		}
	}

	if len(missingFiles) > 0 {
		t.Errorf("nav entries referencing non-existent files (%d):", len(missingFiles))
		for _, f := range missingFiles {
			t.Errorf("  missing: docs/%s", f)
		}
	}

	// Step 3: Walk docs/ to find ALL markdown files.
	allDocFiles := walkMarkdownFiles(t, docsRoot)
	t.Logf("found %d markdown files under docs/", len(allDocFiles))

	// Step 4: Build a lookup set from nav paths for orphan detection.
	navSet := make(map[string]bool, len(navPaths))
	for _, p := range navPaths {
		// Normalize to forward slashes for comparison.
		navSet[filepath.ToSlash(p)] = true
	}

	// Step 5: Identify orphan pages (files in docs/ with no nav entry).
	// Exclude: files in img/ directories, .gitkeep files, and other non-content files.
	var orphans []string
	for _, docFile := range allDocFiles {
		relPath := filepath.ToSlash(docFile)

		// Skip files in img/ directories.
		if strings.Contains(relPath, "/img/") || strings.HasPrefix(relPath, "img/") {
			continue
		}

		if !navSet[relPath] {
			orphans = append(orphans, relPath)
		}
	}

	// Step 6: Report orphan pages.
	if len(orphans) > 0 {
		t.Errorf("orphan documentation pages with no nav entry (%d):", len(orphans))
		for _, o := range orphans {
			t.Errorf("  orphan: docs/%s", o)
		}
	}

	if len(missingFiles) == 0 && len(orphans) == 0 {
		t.Logf("navigation-file bijection validated: %d nav entries, %d doc files, 0 missing, 0 orphans",
			len(navPaths), len(allDocFiles))
	}
}

// TestNavEntriesHaveFiles verifies the forward direction of the bijection:
// every mkdocs.yml nav entry references an existing file.
func TestNavEntriesHaveFiles(t *testing.T) {
	ghpagesRoot := filepath.Join("..", "..", "ghpages")

	mkdocsPath := filepath.Join(ghpagesRoot, "mkdocs.yml")
	docsRoot := filepath.Join(ghpagesRoot, "docs")

	navPaths := extractNavPaths(t, mkdocsPath)

	for _, relPath := range navPaths {
		fullPath := filepath.Join(docsRoot, filepath.FromSlash(relPath))
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			t.Errorf("nav entry %q references non-existent file: %s", relPath, fullPath)
		}
	}
}

// TestNoOrphanPages verifies the reverse direction of the bijection:
// no markdown file under docs/ exists without a corresponding nav entry.
func TestNoOrphanPages(t *testing.T) {
	ghpagesRoot := filepath.Join("..", "..", "ghpages")

	mkdocsPath := filepath.Join(ghpagesRoot, "mkdocs.yml")
	docsRoot := filepath.Join(ghpagesRoot, "docs")

	navPaths := extractNavPaths(t, mkdocsPath)
	navSet := make(map[string]bool, len(navPaths))
	for _, p := range navPaths {
		navSet[filepath.ToSlash(p)] = true
	}

	allDocFiles := walkMarkdownFiles(t, docsRoot)

	for _, docFile := range allDocFiles {
		relPath := filepath.ToSlash(docFile)

		// Skip non-content files.
		if strings.Contains(relPath, "/img/") || strings.HasPrefix(relPath, "img/") {
			continue
		}

		if !navSet[relPath] {
			t.Errorf("orphan page (no nav entry): docs/%s", relPath)
		}
	}
}

// extractNavPaths parses mkdocs.yml and extracts all .md file paths from the
// nav section. It uses a regex approach to find lines containing markdown file
// references (pattern: some-file.md or path/to/file.md).
func extractNavPaths(t *testing.T, mkdocsPath string) []string {
	t.Helper()

	f, err := os.Open(mkdocsPath)
	if err != nil {
		t.Fatalf("cannot open mkdocs.yml: %v", err)
	}
	defer f.Close()

	// The nav section starts with "nav:" at the beginning of a line.
	// Each nav entry eventually references a .md file.
	// We extract all .md file paths from lines within the nav section.
	//
	// Nav entries can look like:
	//   - Home: index.md
	//   - getting-started/index.md
	//   - 'Waveshare Zero LCD HAT (A)': getting-started/waveshare-zero-lcd-hat-a.md
	//   - Setup: getting-started/setup.md

	// Pattern matches .md file paths (with optional leading path components).
	mdPathRe := regexp.MustCompile(`([\w][\w./-]*\.md)`)

	var paths []string
	scanner := bufio.NewScanner(f)
	inNav := false

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		// Detect start of nav section.
		if trimmed == "nav:" {
			inNav = true
			continue
		}

		// Detect end of nav section (a new top-level key that isn't indented).
		if inNav && len(line) > 0 && line[0] != ' ' && line[0] != '\t' && !strings.HasPrefix(trimmed, "-") {
			break
		}

		if !inNav {
			continue
		}

		// Extract .md file path from this line.
		matches := mdPathRe.FindStringSubmatch(line)
		if len(matches) > 1 {
			paths = append(paths, matches[1])
		}
	}

	if err := scanner.Err(); err != nil {
		t.Fatalf("error reading mkdocs.yml: %v", err)
	}

	if len(paths) == 0 {
		t.Fatal("no nav entries found in mkdocs.yml - parsing may have failed")
	}

	return paths
}

// walkMarkdownFiles recursively walks the docs directory and returns all
// .md file paths relative to docsRoot.
func walkMarkdownFiles(t *testing.T, docsRoot string) []string {
	t.Helper()

	var files []string
	err := filepath.Walk(docsRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories.
		if info.IsDir() {
			return nil
		}

		// Only include .md files.
		if !strings.HasSuffix(strings.ToLower(info.Name()), ".md") {
			return nil
		}

		// Get path relative to docsRoot.
		rel, err := filepath.Rel(docsRoot, path)
		if err != nil {
			return err
		}

		files = append(files, rel)
		return nil
	})

	if err != nil {
		t.Fatalf("error walking docs directory: %v", err)
	}

	return files
}
