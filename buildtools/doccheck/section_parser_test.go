package doccheck

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExtractHeadings(t *testing.T) {
	// Create a temporary markdown file.
	dir := t.TempDir()
	mdFile := filepath.Join(dir, "test.md")

	content := `# Title

Some intro text.

## Quick Start

Start instructions here.

## Display Characteristics

| Property | Value |
|----------|-------|
| Resolution | 128 × 64 |

## Troubleshooting

Something about issues.

## Related Pages

- [Hardware](hardware.md)
`
	if err := os.WriteFile(mdFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	headings := extractHeadings(mdFile)

	expected := []string{"Quick Start", "Display Characteristics", "Troubleshooting", "Related Pages"}
	if len(headings) != len(expected) {
		t.Fatalf("expected %d headings, got %d: %v", len(expected), len(headings), headings)
	}

	for i, h := range headings {
		if h != expected[i] {
			t.Errorf("heading[%d] = %q, want %q", i, h, expected[i])
		}
	}
}

func TestExtractHeadings_NonexistentFile(t *testing.T) {
	headings := extractHeadings("/nonexistent/file.md")
	if headings != nil {
		t.Errorf("expected nil for nonexistent file, got %v", headings)
	}
}

func TestExtractTableRows(t *testing.T) {
	dir := t.TempDir()
	mdFile := filepath.Join(dir, "test.md")

	content := `# Panel Name

## Display Characteristics

| Property | Value |
|----------|-------|
| Resolution | 128 × 64 |
| Color | Monochrome |
| Controller | SSD1306 |

## Troubleshooting

| Symptom | Likely Cause | Fix |
|---------|-------------|-----|
| Blank screen | Wrong bus | Check SPI config |
`
	if err := os.WriteFile(mdFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	// Extract from Display Characteristics section.
	// Includes header row + 3 data rows = 4 rows total.
	rows := extractTableRows(mdFile, "Display Characteristics")
	if len(rows) != 4 {
		t.Fatalf("expected 4 rows in Display Characteristics, got %d", len(rows))
	}

	// First row is the header.
	if rows[0][0] != "Property" || rows[0][1] != "Value" {
		t.Errorf("row[0] = %v, want [Property, Value]", rows[0])
	}
	if rows[1][0] != "Resolution" || rows[1][1] != "128 × 64" {
		t.Errorf("row[1] = %v, want [Resolution, 128 × 64]", rows[1])
	}
	if rows[3][0] != "Controller" || rows[3][1] != "SSD1306" {
		t.Errorf("row[3] = %v, want [Controller, SSD1306]", rows[3])
	}

	// Extract from Troubleshooting section.
	// Includes header row + 1 data row = 2 rows total.
	rows = extractTableRows(mdFile, "Troubleshooting")
	if len(rows) != 2 {
		t.Fatalf("expected 2 rows in Troubleshooting, got %d", len(rows))
	}
	if rows[0][0] != "Symptom" {
		t.Errorf("row[0][0] = %q, want %q", rows[0][0], "Symptom")
	}
}

func TestExtractTableRows_StopsAtNextHeading(t *testing.T) {
	dir := t.TempDir()
	mdFile := filepath.Join(dir, "test.md")

	content := `## Options

| Key | Default |
|-----|---------|
| speed | fast |

## CLI Examples

` + "```sh\ncyberhudctl display set main.0 clock\n```\n"

	if err := os.WriteFile(mdFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	// Header row + 1 data row = 2 rows.
	rows := extractTableRows(mdFile, "Options")
	if len(rows) != 2 {
		t.Fatalf("expected 2 rows, got %d: %v", len(rows), rows)
	}
	if rows[0][0] != "Key" || rows[0][1] != "Default" {
		t.Errorf("unexpected header row content: %v", rows[0])
	}
	if rows[1][0] != "speed" || rows[1][1] != "fast" {
		t.Errorf("unexpected data row content: %v", rows[1])
	}
}

func TestExtractCodeBlocks(t *testing.T) {
	dir := t.TempDir()
	mdFile := filepath.Join(dir, "test.md")

	content := "# Mode Page\n\n## Quick Start\n\n```sh\ncyberhudctl display set main.0 clock\n```\n\n## How It Works\n\nSome text.\n\n```json\n{\"mode\": \"clock\"}\n```\n"

	if err := os.WriteFile(mdFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	blocks := extractCodeBlocks(mdFile)
	if len(blocks) != 2 {
		t.Fatalf("expected 2 code blocks, got %d", len(blocks))
	}

	// First block.
	if blocks[0].Language != "sh" {
		t.Errorf("block[0].Language = %q, want %q", blocks[0].Language, "sh")
	}
	if blocks[0].Content != "cyberhudctl display set main.0 clock" {
		t.Errorf("block[0].Content = %q, want %q", blocks[0].Content, "cyberhudctl display set main.0 clock")
	}
	if blocks[0].LineNumber != 5 {
		t.Errorf("block[0].LineNumber = %d, want %d", blocks[0].LineNumber, 5)
	}

	// Second block.
	if blocks[1].Language != "json" {
		t.Errorf("block[1].Language = %q, want %q", blocks[1].Language, "json")
	}
	if blocks[1].Content != `{"mode": "clock"}` {
		t.Errorf("block[1].Content = %q, want %q", blocks[1].Content, `{"mode": "clock"}`)
	}
	if blocks[1].LineNumber != 13 {
		t.Errorf("block[1].LineNumber = %d, want %d", blocks[1].LineNumber, 13)
	}
}

func TestExtractCodeBlocks_NoLanguageTag(t *testing.T) {
	dir := t.TempDir()
	mdFile := filepath.Join(dir, "test.md")

	content := "Some text.\n\n```\nplain code\n```\n"

	if err := os.WriteFile(mdFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	blocks := extractCodeBlocks(mdFile)
	if len(blocks) != 1 {
		t.Fatalf("expected 1 code block, got %d", len(blocks))
	}
	if blocks[0].Language != "" {
		t.Errorf("block[0].Language = %q, want empty", blocks[0].Language)
	}
	if blocks[0].Content != "plain code" {
		t.Errorf("block[0].Content = %q, want %q", blocks[0].Content, "plain code")
	}
}

func TestIsSeparatorRow(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"|---|---|", true},
		{"| --- | --- | --- |", true},
		{"|:---:|:---:|", true},
		{"| Resolution | 128 |", false},
		{"|--|", true},
		{"| Key | Type | Default |", false},
	}

	for _, tt := range tests {
		got := isSeparatorRow(tt.input)
		if got != tt.want {
			t.Errorf("isSeparatorRow(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestParseTableRow(t *testing.T) {
	tests := []struct {
		input string
		want  []string
	}{
		{"| foo | bar | baz |", []string{"foo", "bar", "baz"}},
		{"|a|b|c|", []string{"a", "b", "c"}},
		{"| trimmed  |  spaces |", []string{"trimmed", "spaces"}},
	}

	for _, tt := range tests {
		got := parseTableRow(tt.input)
		if len(got) != len(tt.want) {
			t.Errorf("parseTableRow(%q) = %v, want %v", tt.input, got, tt.want)
			continue
		}
		for i := range got {
			if got[i] != tt.want[i] {
				t.Errorf("parseTableRow(%q)[%d] = %q, want %q", tt.input, i, got[i], tt.want[i])
			}
		}
	}
}
