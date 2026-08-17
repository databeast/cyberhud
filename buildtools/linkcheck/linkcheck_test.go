package linkcheck

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// TestLinkIntegrity validates Property 5: Link Integrity.
// Every relative link and fragment anchor in root documentation files
// (README.md, QUICKSTART.md, CONFIGURATION.md) resolves to an existing file
// and heading.
func TestLinkIntegrity(t *testing.T) {
	projectRoot := filepath.Join("..", "..")

	rootDocs := []string{
		filepath.Join(projectRoot, "README.md"),
		filepath.Join(projectRoot, "QUICKSTART.md"),
		filepath.Join(projectRoot, "CONFIGURATION.md"),
	}

	var allBroken []brokenLink

	for _, docPath := range rootDocs {
		links := extractMarkdownLinks(t, docPath)
		t.Logf("%s: found %d relative links", filepath.Base(docPath), len(links))

		for _, link := range links {
			if err := validateLink(projectRoot, docPath, link); err != nil {
				allBroken = append(allBroken, brokenLink{
					file:   docPath,
					link:   link,
					reason: err.Error(),
				})
			}
		}
	}

	for _, bl := range allBroken {
		t.Errorf("broken link in %s:\n  link: [%s](%s)\n  reason: %s",
			filepath.Base(bl.file), bl.link.text, bl.link.raw, bl.reason)
	}

	if len(allBroken) == 0 {
		t.Logf("all links valid across %d root documentation files", len(rootDocs))
	}
}

// TestFragmentAnchorsResolve specifically tests that links with fragment
// anchors (#heading) resolve to actual headings in target files. Links with
// non-existent fragments are flagged as invalid per requirements.
func TestFragmentAnchorsResolve(t *testing.T) {
	projectRoot := filepath.Join("..", "..")

	rootDocs := []string{
		filepath.Join(projectRoot, "README.md"),
		filepath.Join(projectRoot, "QUICKSTART.md"),
		filepath.Join(projectRoot, "CONFIGURATION.md"),
	}

	var fragmentLinks int
	var brokenFragments []brokenLink

	for _, docPath := range rootDocs {
		links := extractMarkdownLinks(t, docPath)

		for _, link := range links {
			if link.fragment == "" {
				continue
			}
			fragmentLinks++

			if err := validateLink(projectRoot, docPath, link); err != nil {
				brokenFragments = append(brokenFragments, brokenLink{
					file:   docPath,
					link:   link,
					reason: err.Error(),
				})
			}
		}
	}

	t.Logf("checked %d fragment anchor links", fragmentLinks)

	for _, bl := range brokenFragments {
		t.Errorf("invalid fragment anchor in %s:\n  link: [%s](%s)\n  fragment: #%s\n  reason: %s",
			filepath.Base(bl.file), bl.link.text, bl.link.raw, bl.link.fragment, bl.reason)
	}
}

// markdownLink represents a parsed markdown link.
type markdownLink struct {
	text     string // link text [text](...)
	raw      string // full target path#fragment
	path     string // file path part (empty for same-file anchors)
	fragment string // fragment part (without #)
}

type brokenLink struct {
	file   string
	link   markdownLink
	reason string
}

// linkPattern matches markdown links: [text](target)
// It captures the text and the target (path + optional fragment).
var linkPattern = regexp.MustCompile(`\[([^\]]*)\]\(([^)]+)\)`)

// extractMarkdownLinks reads a markdown file and extracts all relative links.
// It skips:
//   - absolute URLs (http://, https://, mailto:)
//   - image links (![...](...)
//   - links inside code blocks
func extractMarkdownLinks(t *testing.T, filePath string) []markdownLink {
	t.Helper()

	f, err := os.Open(filePath)
	if err != nil {
		t.Fatalf("cannot open %s: %v", filePath, err)
	}
	defer f.Close()

	var links []markdownLink
	scanner := bufio.NewScanner(f)
	inCodeBlock := false

	for scanner.Scan() {
		line := scanner.Text()

		// Track code block boundaries.
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "```") {
			inCodeBlock = !inCodeBlock
			continue
		}
		if inCodeBlock {
			continue
		}

		matches := linkPattern.FindAllStringSubmatch(line, -1)
		for _, match := range matches {
			text := match[1]
			target := match[2]

			// Skip absolute URLs.
			if strings.HasPrefix(target, "http://") ||
				strings.HasPrefix(target, "https://") ||
				strings.HasPrefix(target, "mailto:") {
				continue
			}

			// Skip image links (the character before [ is !).
			idx := strings.Index(line, match[0])
			if idx > 0 && line[idx-1] == '!' {
				continue
			}

			link := parseLink(text, target)
			links = append(links, link)
		}
	}

	if err := scanner.Err(); err != nil {
		t.Fatalf("error reading %s: %v", filePath, err)
	}

	return links
}

// parseLink splits a link target into path and fragment components.
func parseLink(text, target string) markdownLink {
	link := markdownLink{
		text: text,
		raw:  target,
	}

	if idx := strings.Index(target, "#"); idx >= 0 {
		link.path = target[:idx]
		link.fragment = target[idx+1:]
	} else {
		link.path = target
	}

	return link
}

// validateLink checks that a link target resolves to an existing file and,
// if a fragment is present, to an actual heading in that file.
func validateLink(projectRoot, sourceFile string, link markdownLink) error {
	// Determine the target file path.
	var targetPath string
	if link.path == "" {
		// Same-file anchor (e.g., #troubleshooting).
		targetPath = sourceFile
	} else {
		// Resolve relative to the source file's directory.
		sourceDir := filepath.Dir(sourceFile)
		targetPath = filepath.Join(sourceDir, link.path)
	}

	// Clean the path.
	targetPath = filepath.Clean(targetPath)

	// Check the file exists.
	info, err := os.Stat(targetPath)
	if err != nil {
		return fmt.Errorf("target file does not exist: %s", targetPath)
	}

	// If it's a directory, that's not a valid link target for markdown.
	if info.IsDir() {
		return fmt.Errorf("target is a directory, not a file: %s", targetPath)
	}

	// If there's a fragment, verify the heading exists.
	if link.fragment != "" {
		headings := extractHeadings(targetPath)
		if !headingExists(headings, link.fragment) {
			return fmt.Errorf("fragment #%s does not match any heading in %s (available: %s)",
				link.fragment, filepath.Base(targetPath), formatAvailableHeadings(headings, link.fragment))
		}
	}

	return nil
}

// extractHeadings reads a markdown file and returns all heading anchors
// using GitHub-flavored markdown anchor generation rules:
// - Convert to lowercase
// - Replace spaces with hyphens
// - Strip most punctuation (keep hyphens and alphanumerics)
// - Strip leading/trailing hyphens
func extractHeadings(filePath string) []string {
	f, err := os.Open(filePath)
	if err != nil {
		return nil
	}
	defer f.Close()

	var headings []string
	scanner := bufio.NewScanner(f)
	inCodeBlock := false

	for scanner.Scan() {
		line := scanner.Text()

		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "```") {
			inCodeBlock = !inCodeBlock
			continue
		}
		if inCodeBlock {
			continue
		}

		// Match ATX headings (# Heading, ## Heading, etc.).
		if strings.HasPrefix(trimmed, "#") {
			// Strip the leading #s and space.
			heading := strings.TrimLeft(trimmed, "#")
			heading = strings.TrimSpace(heading)
			if heading != "" {
				anchor := headingToAnchor(heading)
				headings = append(headings, anchor)
			}
		}
	}

	return headings
}

// headingToAnchor converts a heading text to a GitHub-flavored markdown anchor.
// Rules:
//   - Convert to lowercase
//   - Replace spaces with hyphens
//   - Remove characters that are not alphanumeric, hyphens, or spaces
//   - Collapse multiple hyphens into one
//   - Strip leading/trailing hyphens
func headingToAnchor(heading string) string {
	// Convert to lowercase.
	s := strings.ToLower(heading)

	var result strings.Builder
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z':
			result.WriteRune(r)
		case r >= '0' && r <= '9':
			result.WriteRune(r)
		case r == ' ' || r == '-':
			result.WriteRune('-')
		case r == '_':
			result.WriteRune('_')
			// All other punctuation is stripped.
		}
	}

	// Collapse multiple hyphens.
	anchor := result.String()
	for strings.Contains(anchor, "--") {
		anchor = strings.ReplaceAll(anchor, "--", "-")
	}

	// Strip leading/trailing hyphens.
	anchor = strings.Trim(anchor, "-")

	return anchor
}

// headingExists checks if a fragment matches any heading anchor in the list.
func headingExists(headings []string, fragment string) bool {
	// Normalize the fragment using the same rules.
	normalizedFragment := strings.ToLower(fragment)

	for _, h := range headings {
		if h == normalizedFragment {
			return true
		}
	}
	return false
}

// formatAvailableHeadings returns a subset of headings that are close matches
// to help debugging.
func formatAvailableHeadings(headings []string, fragment string) string {
	normalizedFragment := strings.ToLower(fragment)

	// Find headings that share a prefix or contain the fragment.
	var candidates []string
	for _, h := range headings {
		if strings.Contains(h, normalizedFragment) || strings.Contains(normalizedFragment, h) {
			candidates = append(candidates, h)
		}
	}

	if len(candidates) == 0 {
		// Return first few headings for context.
		max := 5
		if len(headings) < max {
			max = len(headings)
		}
		candidates = headings[:max]
	}

	if len(candidates) > 5 {
		candidates = candidates[:5]
	}

	return strings.Join(candidates, ", ")
}
