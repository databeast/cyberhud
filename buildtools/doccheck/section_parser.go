package doccheck

import (
	"bufio"
	"os"
	"strings"
)

// CodeBlock represents a fenced code block extracted from a markdown file.
type CodeBlock struct {
	Language   string // language tag after the opening ```
	Content    string // full text content between the fences
	LineNumber int    // line number where the code block starts (1-based)
}

// extractHeadings reads a markdown file and returns all H2 headings (lines
// starting with "## "). The returned strings contain the heading text without
// the "## " prefix.
func extractHeadings(filePath string) []string {
	f, err := os.Open(filePath)
	if err != nil {
		return nil
	}
	defer f.Close()

	var headings []string
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "## ") {
			headings = append(headings, strings.TrimPrefix(line, "## "))
		}
	}

	return headings
}

// extractTableRows finds a section by its H2 heading, then extracts all
// pipe-delimited table rows under that heading. The header separator row
// (e.g., "|---|---|") is excluded. Each row is split into cells with
// leading/trailing whitespace trimmed. Extraction stops at the next H2
// heading or EOF.
func extractTableRows(filePath, sectionHeading string) [][]string {
	f, err := os.Open(filePath)
	if err != nil {
		return nil
	}
	defer f.Close()

	var rows [][]string
	scanner := bufio.NewScanner(f)
	inSection := false

	for scanner.Scan() {
		line := scanner.Text()

		// Detect the target section heading.
		if strings.HasPrefix(line, "## ") {
			heading := strings.TrimPrefix(line, "## ")
			if heading == sectionHeading {
				inSection = true
				continue
			}
			// If we were in the section and hit another H2, stop.
			if inSection {
				break
			}
			continue
		}

		if !inSection {
			continue
		}

		// Skip non-table lines.
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, "|") {
			continue
		}

		// Skip separator rows (e.g., |---|---|---| or | --- | --- |).
		if isSeparatorRow(trimmed) {
			continue
		}

		// Parse the pipe-delimited row into cells.
		cells := parseTableRow(trimmed)
		if len(cells) > 0 {
			rows = append(rows, cells)
		}
	}

	return rows
}

// extractCodeBlocks extracts all fenced code blocks (``` delimited) from the
// file. Each CodeBlock contains the language tag, full text content between
// the fences, and the line number where the opening fence appears.
func extractCodeBlocks(filePath string) []CodeBlock {
	f, err := os.Open(filePath)
	if err != nil {
		return nil
	}
	defer f.Close()

	var blocks []CodeBlock
	scanner := bufio.NewScanner(f)
	lineNum := 0
	inBlock := false
	var currentBlock CodeBlock
	var contentLines []string

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		if !inBlock {
			// Detect opening fence.
			if strings.HasPrefix(trimmed, "```") {
				inBlock = true
				currentBlock = CodeBlock{
					Language:   strings.TrimPrefix(trimmed, "```"),
					LineNumber: lineNum,
				}
				contentLines = nil
				continue
			}
		} else {
			// Detect closing fence.
			if trimmed == "```" {
				currentBlock.Content = strings.Join(contentLines, "\n")
				blocks = append(blocks, currentBlock)
				inBlock = false
				continue
			}
			contentLines = append(contentLines, line)
		}
	}

	return blocks
}

// isSeparatorRow returns true if the row is a markdown table separator
// (contains only pipes, dashes, spaces, and colons).
func isSeparatorRow(row string) bool {
	for _, ch := range row {
		switch ch {
		case '|', '-', ' ', ':', '\t':
			continue
		default:
			return false
		}
	}
	return true
}

// parseTableRow splits a pipe-delimited row into trimmed cells.
// Leading and trailing pipes are stripped before splitting.
func parseTableRow(row string) []string {
	// Remove leading and trailing pipe characters.
	row = strings.TrimSpace(row)
	if strings.HasPrefix(row, "|") {
		row = row[1:]
	}
	if strings.HasSuffix(row, "|") {
		row = row[:len(row)-1]
	}

	parts := strings.Split(row, "|")
	cells := make([]string, 0, len(parts))
	for _, p := range parts {
		cells = append(cells, strings.TrimSpace(p))
	}
	return cells
}
