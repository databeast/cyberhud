package main

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// TestCLIExampleValidity validates every cyberhudctl command example in documentation uses a verb and argument
// structure that exists in the protocolCommand switch routing.
func TestCLIExampleValidity(t *testing.T) {
	// Find the project root (two levels up from cmd/cyberhudctl/).
	projectRoot := filepath.Join("..", "..")

	// Directories and files to scan for CLI examples.
	docPaths := []string{
		filepath.Join(projectRoot, "README.md"),
		filepath.Join(projectRoot, "QUICKSTART.md"),
		filepath.Join(projectRoot, "CONFIGURATION.md"),
		filepath.Join(projectRoot, "ghpages", "docs"),
		filepath.Join(projectRoot, "hardware"),
	}

	// Collect all cyberhudctl examples from documentation.
	examples := collectCLIExamples(t, docPaths)
	if len(examples) == 0 {
		t.Fatal("no cyberhudctl examples found in documentation — check doc paths")
	}

	t.Logf("found %d cyberhudctl examples across documentation", len(examples))

	for _, ex := range examples {
		t.Run(ex.source+"/"+strings.Join(ex.args, "_"), func(t *testing.T) {
			_, _, err := protocolCommand(ex.args)
			if err != nil {
				t.Errorf("invalid CLI example in %s:\n  command: cyberhudctl %s\n  error: %s",
					ex.file, strings.Join(ex.args, " "), err)
			}
		})
	}
}

type cliExample struct {
	file   string // relative file path
	source string // short identifier for test name
	args   []string
}

// collectCLIExamples walks the given paths and extracts cyberhudctl command
// arguments from markdown code blocks.
func collectCLIExamples(t *testing.T, paths []string) []cliExample {
	t.Helper()
	var examples []cliExample

	for _, p := range paths {
		info, err := os.Stat(p)
		if err != nil {
			t.Logf("skipping %s: %v", p, err)
			continue
		}
		if info.IsDir() {
			err = filepath.Walk(p, func(path string, fi os.FileInfo, err error) error {
				if err != nil {
					return nil
				}
				if fi.IsDir() {
					return nil
				}
				if strings.HasSuffix(fi.Name(), ".md") {
					exs := extractExamplesFromFile(t, path)
					examples = append(examples, exs...)
				}
				return nil
			})
			if err != nil {
				t.Logf("walk error for %s: %v", p, err)
			}
		} else {
			exs := extractExamplesFromFile(t, p)
			examples = append(examples, exs...)
		}
	}
	return examples
}

// cyberhudctlPattern matches lines starting with cyberhudctl (with optional $
// prompt prefix).
var cyberhudctlPattern = regexp.MustCompile(`(?:^|\$\s+)cyberhudctl\s+(.+)`)

// extractExamplesFromFile reads a markdown file and extracts cyberhudctl
// command invocations from code blocks.
func extractExamplesFromFile(t *testing.T, path string) []cliExample {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Logf("cannot open %s: %v", path, err)
		return nil
	}
	defer f.Close()

	var examples []cliExample
	scanner := bufio.NewScanner(f)
	inCodeBlock := false
	shortName := filepath.Base(filepath.Dir(path)) + "/" + filepath.Base(path)

	for scanner.Scan() {
		line := scanner.Text()

		// Track code block boundaries.
		if strings.HasPrefix(strings.TrimSpace(line), "```") {
			inCodeBlock = !inCodeBlock
			continue
		}

		// Only process lines inside code blocks to avoid matching inline
		// mentions in prose and table cells.
		if !inCodeBlock {
			continue
		}

		// Match cyberhudctl invocations.
		matches := cyberhudctlPattern.FindStringSubmatch(line)
		if len(matches) < 2 {
			continue
		}

		argStr := strings.TrimSpace(matches[1])

		// Strip trailing comments (# ...).
		if idx := strings.Index(argStr, " #"); idx >= 0 {
			argStr = strings.TrimSpace(argStr[:idx])
		}

		// Skip template/placeholder lines (e.g., "[flags] <command>").
		if strings.Contains(argStr, "[flags]") {
			continue
		}

		// Parse arguments, handling quoted strings.
		args := parseCommandArgs(argStr)
		if len(args) == 0 {
			continue
		}

		// Skip flags (-socket, -timeout) and their values.
		args = stripFlags(args)
		if len(args) == 0 {
			continue
		}
		if containsPlaceholderArg(args) {
			continue
		}

		// Skip multi-command examples containing semicolons — the
		// multi-command parser handles those differently from protocolCommand.
		if containsSemicolonArg(args) {
			continue
		}

		examples = append(examples, cliExample{
			file:   path,
			source: shortName,
			args:   args,
		})
	}
	return examples
}

// stripFlags removes -flag and -flag value pairs from an argument list.
func stripFlags(args []string) []string {
	var result []string
	for i := 0; i < len(args); i++ {
		if strings.HasPrefix(args[i], "-") {
			// Skip the flag and its value (e.g., -socket /path, -timeout 5s).
			i++ // skip value
			continue
		}
		result = append(result, args[i])
	}
	return result
}

// containsSemicolonArg checks if the args contain a semicolon token (multi-command).
func containsSemicolonArg(args []string) bool {
	for _, a := range args {
		if a == ";" || a == "';'" {
			return true
		}
	}
	return false
}

func containsPlaceholderArg(args []string) bool {
	for _, a := range args {
		if strings.Contains(a, "<") || strings.Contains(a, ">") {
			return true
		}
	}
	return false
}

// parseCommandArgs splits a command string into arguments, respecting
// single and double quotes.
func parseCommandArgs(s string) []string {
	var args []string
	var current strings.Builder
	inSingle := false
	inDouble := false

	for i := 0; i < len(s); i++ {
		ch := s[i]
		switch {
		case ch == '\'' && !inDouble:
			inSingle = !inSingle
		case ch == '"' && !inSingle:
			inDouble = !inDouble
		case ch == ' ' && !inSingle && !inDouble:
			if current.Len() > 0 {
				args = append(args, current.String())
				current.Reset()
			}
		default:
			current.WriteByte(ch)
		}
	}
	if current.Len() > 0 {
		args = append(args, current.String())
	}
	return args
}
