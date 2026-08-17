package doccheck

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

// sourceOption represents one option extracted from Go source OptionDefinition entries.
type sourceOption struct {
	Key     string
	Type    string
	Default string
}

// TestOptionsTableFreshness validates Property 8: Options Table Freshness.
// For any mode whose catalog.Register() definition includes one or more
// OptionDefinition entries, the options table in that mode's documentation
// page SHALL contain exactly the same set of option keys, and for each key
// the documented Type and Default SHALL equal the corresponding fields in
// the source OptionDefinition.
func TestOptionsTableFreshness(t *testing.T) {
	projectRoot := filepath.Join("..", "..")
	modesDir := filepath.Join(projectRoot, "display", "modes")
	docsDir := filepath.Join(projectRoot, "ghpages", "docs", "display-modes")

	// Discover all registered modes from source.
	registeredModes := discoverRegisteredModes(t, projectRoot)
	t.Logf("discovered %d registered modes", len(registeredModes))

	// For each registered mode, extract options from source and compare with docs.
	for modeID := range registeredModes {
		// Skip utility/excluded modes — they may have no Options section or a
		// different template.
		if utilityModeExclusions[modeID] {
			continue
		}

		// Find the mode's source directory.
		modeDir := modeSourceDir(modesDir, modeID)
		if modeDir == "" {
			continue
		}

		// Extract options from source code.
		srcOptions := extractSourceOptions(t, modeDir)
		if len(srcOptions) == 0 {
			// Mode has no options — verify there is no Options table in docs
			// or that the doc page simply doesn't have one. Both are acceptable.
			continue
		}

		// Find the corresponding documentation page.
		docPath := filepath.Join(docsDir, modeDocPagePath(modeID))
		if _, err := os.Stat(docPath); os.IsNotExist(err) {
			// No doc page — other tests catch this.
			continue
		}

		// Extract options table from the doc page.
		docOptions := extractDocOptions(t, docPath)
		if docOptions == nil {
			t.Errorf("mode %q: source defines %d options but documentation page has no Options table", modeID, len(srcOptions))
			continue
		}

		// Compare: build key sets.
		srcKeys := make(map[string]sourceOption)
		for _, opt := range srcOptions {
			srcKeys[opt.Key] = opt
		}

		docKeys := make(map[string]sourceOption)
		for _, opt := range docOptions {
			docKeys[opt.Key] = opt
		}

		// Report keys missing from docs.
		var missingInDocs []string
		for key := range srcKeys {
			if _, ok := docKeys[key]; !ok {
				missingInDocs = append(missingInDocs, key)
			}
		}
		sort.Strings(missingInDocs)
		for _, key := range missingInDocs {
			t.Errorf("mode %q: option key %q is defined in source but missing from documentation Options table", modeID, key)
		}

		// Report extra keys in docs.
		var extraInDocs []string
		for key := range docKeys {
			if _, ok := srcKeys[key]; !ok {
				extraInDocs = append(extraInDocs, key)
			}
		}
		sort.Strings(extraInDocs)
		for _, key := range extraInDocs {
			t.Errorf("mode %q: option key %q is in documentation Options table but not defined in source", modeID, key)
		}

		// For matching keys, compare Type and Default.
		for key, srcOpt := range srcKeys {
			docOpt, ok := docKeys[key]
			if !ok {
				continue
			}
			if srcOpt.Type != "" && docOpt.Type != "" && srcOpt.Type != docOpt.Type {
				t.Errorf("mode %q: option %q type mismatch: source=%q, docs=%q", modeID, key, srcOpt.Type, docOpt.Type)
			}
			if srcOpt.Default != "" && docOpt.Default != "" && srcOpt.Default != docOpt.Default {
				t.Errorf("mode %q: option %q default mismatch: source=%q, docs=%q", modeID, key, srcOpt.Default, docOpt.Default)
			}
		}
	}
}

// modeSourceDir returns the source directory path for a given mode ID.
// It looks for a matching subdirectory under display/modes/.
func modeSourceDir(modesDir, modeID string) string {
	// Mode IDs typically match directory names directly.
	// Special case: gpio-control -> gpio_control directory.
	dirName := modeID
	switch modeID {
	case "gpio-control":
		dirName = "gpio_control"
	}

	dir := filepath.Join(modesDir, dirName)
	if info, err := os.Stat(dir); err == nil && info.IsDir() {
		return dir
	}
	return ""
}

// extractSourceOptions scans Go source files in a mode directory (and its
// source/ subdirectory) for OptionDefinition entries and extracts Key, Type,
// and Default values.
func extractSourceOptions(t *testing.T, modeDir string) []sourceOption {
	t.Helper()

	var options []sourceOption

	// Collect all Go source files in the mode directory and source/ subdirectory.
	goFiles := collectGoFiles(modeDir)
	sourceSubDir := filepath.Join(modeDir, "source")
	if info, err := os.Stat(sourceSubDir); err == nil && info.IsDir() {
		goFiles = append(goFiles, collectGoFiles(sourceSubDir)...)
	}
	contentSubDir := filepath.Join(modeDir, "content")
	if info, err := os.Stat(contentSubDir); err == nil && info.IsDir() {
		goFiles = append(goFiles, collectGoFiles(contentSubDir)...)
	}

	keyPattern := regexp.MustCompile(`Key:\s*"([^"]+)"`)
	typePattern := regexp.MustCompile(`Type:\s*"([^"]+)"`)
	defaultPattern := regexp.MustCompile(`Default:\s*"([^"]*)"`)

	seen := make(map[string]bool)

	for _, path := range goFiles {
		content, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		text := string(content)

		// Find all OptionDefinition struct literals by looking for {Key: "..."
		// entries that are part of OptionDefinition slices.
		// We look for blocks containing Key: "..." that appear after
		// []catalog.OptionDefinition{ or catalog.OptionDefinition{
		opts := extractOptionBlocks(text, keyPattern, typePattern, defaultPattern)
		for _, opt := range opts {
			if !seen[opt.Key] {
				seen[opt.Key] = true
				options = append(options, opt)
			}
		}
	}

	return options
}

// extractOptionBlocks finds OptionDefinition struct literal blocks and extracts
// Key, Type, and Default from each.
func extractOptionBlocks(text string, keyRe, typeRe, defaultRe *regexp.Regexp) []sourceOption {
	var options []sourceOption

	// Strategy: find all []catalog.OptionDefinition{ or catalog.OptionDefinition{
	// slice/struct openings, then extract all Key entries within each block.
	sliceOpener := regexp.MustCompile(`(?:\[\]catalog\.OptionDefinition\{|catalog\.OptionDefinition\{)`)
	openers := sliceOpener.FindAllStringIndex(text, -1)

	for _, opener := range openers {
		// From the opener position, scan forward for {Key: "..."} entries.
		// We scan from the opener to find all Key entries within this context.
		remaining := text[opener[1]:]

		// Extract all Key entries from this block. We scan until we hit a line
		// that starts a new non-OptionDefinition context (e.g., a function
		// definition or a closing of the outer struct). For simplicity, we
		// extract all Key entries until the next slice opener or end of the
		// enclosing expression.
		keyMatches := keyRe.FindAllStringSubmatchIndex(remaining, -1)

		for _, km := range keyMatches {
			key := remaining[km[2]:km[3]]

			// Look forward up to 300 chars for Type and Default.
			snippetStart := km[1]
			snippetEnd := snippetStart + 300
			if snippetEnd > len(remaining) {
				snippetEnd = len(remaining)
			}
			snippet := remaining[snippetStart:snippetEnd]

			// Limit snippet to the current struct literal (stop at next `Key:`).
			if nextKey := strings.Index(snippet, `Key:`); nextKey > 0 {
				snippet = snippet[:nextKey]
			}

			var optType, optDefault string
			if m := typeRe.FindStringSubmatch(snippet); m != nil {
				optType = m[1]
			}
			if m := defaultRe.FindStringSubmatch(snippet); m != nil {
				optDefault = m[1]
			}

			options = append(options, sourceOption{
				Key:     key,
				Type:    optType,
				Default: optDefault,
			})
		}
	}

	// Also handle single catalog.OptionDefinition{Key: ...} literals that appear
	// in append() calls without the slice syntax.
	singlePattern := regexp.MustCompile(`catalog\.OptionDefinition\{Key:\s*"([^"]+)"`)
	singleMatches := singlePattern.FindAllStringSubmatchIndex(text, -1)
	existingKeys := make(map[string]bool)
	for _, o := range options {
		existingKeys[o.Key] = true
	}

	for _, sm := range singleMatches {
		key := text[sm[2]:sm[3]]
		if existingKeys[key] {
			continue
		}

		// Extract Type and Default from surrounding context.
		snippetStart := sm[1]
		snippetEnd := snippetStart + 300
		if snippetEnd > len(text) {
			snippetEnd = len(text)
		}
		snippet := text[snippetStart:snippetEnd]

		if nextKey := strings.Index(snippet, `Key:`); nextKey > 0 {
			snippet = snippet[:nextKey]
		}

		var optType, optDefault string
		if m := typeRe.FindStringSubmatch(snippet); m != nil {
			optType = m[1]
		}
		if m := defaultRe.FindStringSubmatch(snippet); m != nil {
			optDefault = m[1]
		}

		options = append(options, sourceOption{
			Key:     key,
			Type:    optType,
			Default: optDefault,
		})
		existingKeys[key] = true
	}

	return options
}

// collectGoFiles returns all non-test .go files in a directory (non-recursive).
func collectGoFiles(dir string) []string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	var files []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasSuffix(name, ".go") && !strings.HasSuffix(name, "_test.go") {
			files = append(files, filepath.Join(dir, name))
		}
	}
	return files
}

// extractDocOptions parses the Options table from a mode documentation page.
// Returns nil if no Options section or table is found.
func extractDocOptions(t *testing.T, docPath string) []sourceOption {
	t.Helper()

	// Try the "Options" heading first (new template).
	rows := extractTableRows(docPath, "Options")
	if rows == nil || len(rows) == 0 {
		// Some older pages use "Configuration" as the heading.
		rows = extractTableRows(docPath, "Configuration")
	}
	if rows == nil || len(rows) == 0 {
		return nil
	}

	// The first row is the header row; extractTableRows already skips the
	// separator row. Determine column indices from the header.
	header := rows[0]
	keyCol := -1
	typeCol := -1
	defaultCol := -1

	for i, cell := range header {
		lower := strings.ToLower(strings.TrimSpace(cell))
		// Strip backticks and formatting.
		lower = strings.Trim(lower, "`")
		switch {
		case lower == "key":
			keyCol = i
		case lower == "type":
			typeCol = i
		case lower == "default":
			defaultCol = i
		}
	}

	if keyCol == -1 {
		// Cannot parse without a Key column.
		return nil
	}

	var options []sourceOption
	for _, row := range rows[1:] {
		if keyCol >= len(row) {
			continue
		}
		key := strings.TrimSpace(row[keyCol])
		// Strip backticks and formatting from key.
		key = strings.Trim(key, "`")
		if key == "" {
			continue
		}

		var optType, optDefault string
		if typeCol >= 0 && typeCol < len(row) {
			optType = strings.TrimSpace(row[typeCol])
			optType = strings.Trim(optType, "`")
		}
		if defaultCol >= 0 && defaultCol < len(row) {
			optDefault = strings.TrimSpace(row[defaultCol])
			optDefault = strings.Trim(optDefault, "`")
		}

		options = append(options, sourceOption{
			Key:     key,
			Type:    optType,
			Default: optDefault,
		})
	}

	return options
}
