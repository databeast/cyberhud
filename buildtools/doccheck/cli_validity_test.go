package doccheck

import (
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// knownVerbs is the set of top-level command verbs accepted by the
// protocolCommand routing in cmd/cyberhudctl/main.go. This list must be kept
// in sync with the switch cases in that function.
var knownVerbs = map[string]bool{
	"display": true,
	"gpio":    true,
	"help":    true,
	"status":  true,
	"stemma":  true,
	"freeze":  true,
	"policy":  true,
	"region":  true,
	"raw":     true,
}

// cyberhudctlVerbRe captures the first argument (verb) after "cyberhudctl".
// It's applied only after isCyberhudctlInvocation confirms the line is a
// direct invocation and not a build target or argument to another command.
var cyberhudctlVerbRe = regexp.MustCompile(`cyberhudctl\s+(\S+)`)

// isCyberhudctlInvocation returns true if the line represents a direct
// cyberhudctl command invocation (as opposed to cyberhudctl appearing as an
// argument to another command like "go build -o cyberhudctl ...").
func isCyberhudctlInvocation(line string) bool {
	trimmed := strings.TrimSpace(line)

	// Strip common shell prompt prefixes.
	trimmed = strings.TrimPrefix(trimmed, "$ ")
	trimmed = strings.TrimPrefix(trimmed, "# ")
	trimmed = strings.TrimPrefix(trimmed, "sudo ")

	return strings.HasPrefix(trimmed, "cyberhudctl")
}

// TestCLIExampleValidity validates Property 3 (doccheck variant): CLI Example Validity.
// Every cyberhudctl command example in documentation uses a verb that exists
// in the protocolCommand routing switch. Unlike the test in cmd/cyberhudctl/
// which calls protocolCommand() directly, this version validates against a
// hardcoded set of known verbs since it lives in a separate package.
func TestCLIExampleValidity(t *testing.T) {
	docsDir := filepath.Join("..", "..", "ghpages", "docs")

	// Walk all markdown files under ghpages/docs/.
	mdFiles := walkMarkdownFiles(t, docsDir)
	if len(mdFiles) == 0 {
		t.Fatal("no markdown files found under ghpages/docs/")
	}

	type cliInvocation struct {
		file string
		line int
		verb string
		raw  string
	}

	var invocations []cliInvocation

	for _, relPath := range mdFiles {
		fullPath := filepath.Join(docsDir, relPath)
		blocks := extractCodeBlocks(fullPath)

		for _, block := range blocks {
			// Only examine sh/bash code blocks.
			lang := strings.ToLower(strings.TrimSpace(block.Language))
			if lang != "sh" && lang != "bash" && lang != "shell" {
				continue
			}

			lines := strings.Split(block.Content, "\n")
			for i, line := range lines {
				// Skip lines with placeholders like <command>, [flags].
				if strings.Contains(line, "<command>") || strings.Contains(line, "[flags]") {
					continue
				}

				// Only consider lines that are direct cyberhudctl invocations.
				if !isCyberhudctlInvocation(line) {
					continue
				}

				// Find cyberhudctl invocations and extract the verb.
				match := cyberhudctlVerbRe.FindStringSubmatch(line)
				if match == nil {
					continue
				}

				verb := match[1]

				// Skip flags (arguments starting with -).
				if strings.HasPrefix(verb, "-") {
					continue
				}

				// Skip placeholders in the verb position.
				if strings.Contains(verb, "<") || strings.Contains(verb, ">") ||
					strings.Contains(verb, "[") || strings.Contains(verb, "]") {
					continue
				}

				invocations = append(invocations, cliInvocation{
					file: relPath,
					line: block.LineNumber + i + 1, // +1 for the opening fence line
					verb: verb,
					raw:  strings.TrimSpace(line),
				})
			}
		}
	}

	if len(invocations) == 0 {
		t.Fatal("no cyberhudctl invocations found in documentation code blocks")
	}

	t.Logf("found %d cyberhudctl invocations across documentation", len(invocations))

	var failures int
	for _, inv := range invocations {
		if !knownVerbs[inv.verb] {
			t.Errorf("unknown verb %q in %s (line %d):\n  %s",
				inv.verb, inv.file, inv.line, inv.raw)
			failures++
		}
	}

	if failures == 0 {
		t.Logf("all %d CLI examples use valid verbs", len(invocations))
	}
}
