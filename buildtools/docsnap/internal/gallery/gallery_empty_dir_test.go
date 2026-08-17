package gallery_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/databeast/cyberhud/buildtools/docsnap/internal/gallery"
)

// TestEmptyDirectoryPreservesGalleryMarkers verifies Requirement 5.5:
// IF the img/thermal/ directory exists but is empty, THEN the gallery tool
// SHALL inject empty gallery markers to maintain consistent document structure.
//
// The gallery tool's Run function skips empty directories (no images found),
// which means it does not modify the page file. Pre-existing empty markers
// in the page are preserved, maintaining consistent document structure.
func TestEmptyDirectoryPreservesGalleryMarkers(t *testing.T) {
	// Create a temporary directory structure simulating the docs layout.
	tmpDir := t.TempDir()

	imgRoot := filepath.Join(tmpDir, "img")
	pagesRoot := filepath.Join(tmpDir, "pages")

	// Create an empty thermal image directory.
	thermalImgDir := filepath.Join(imgRoot, "thermal")
	if err := os.MkdirAll(thermalImgDir, 0755); err != nil {
		t.Fatalf("failed to create thermal img dir: %v", err)
	}

	// Create the pages directory and a thermal.md page with pre-existing empty markers.
	if err := os.MkdirAll(pagesRoot, 0755); err != nil {
		t.Fatalf("failed to create pages dir: %v", err)
	}

	pageContent := `# Thermal

## Styles

Some style descriptions here.

<!-- snapshot-gallery:start -->
<!-- snapshot-gallery:end -->

## Options

| Key | Type |
`

	pagePath := filepath.Join(pagesRoot, "thermal.md")
	if err := os.WriteFile(pagePath, []byte(pageContent), 0644); err != nil {
		t.Fatalf("failed to write thermal.md: %v", err)
	}

	// Run the gallery tool.
	if err := gallery.Run(imgRoot, pagesRoot); err != nil {
		t.Fatalf("gallery.Run failed: %v", err)
	}

	// Read the page back and verify markers are preserved.
	got, err := os.ReadFile(pagePath)
	if err != nil {
		t.Fatalf("failed to read thermal.md after gallery run: %v", err)
	}

	gotStr := string(got)

	// Verify the empty markers are still present.
	if !strings.Contains(gotStr, "<!-- snapshot-gallery:start -->") {
		t.Error("expected <!-- snapshot-gallery:start --> marker to be preserved, but it was removed")
	}
	if !strings.Contains(gotStr, "<!-- snapshot-gallery:end -->") {
		t.Error("expected <!-- snapshot-gallery:end --> marker to be preserved, but it was removed")
	}

	// Verify the page content is unchanged (gallery tool should not modify it).
	if gotStr != pageContent {
		t.Errorf("page content was modified when directory is empty.\nWant:\n%s\nGot:\n%s", pageContent, gotStr)
	}
}

// TestEmptyDirectoryNoMarkersDoesNotInject verifies that when a page has NO
// pre-existing markers and the image directory is empty, the gallery tool
// does not inject anything (it skips the mode entirely).
func TestEmptyDirectoryNoMarkersDoesNotInject(t *testing.T) {
	tmpDir := t.TempDir()

	imgRoot := filepath.Join(tmpDir, "img")
	pagesRoot := filepath.Join(tmpDir, "pages")

	// Create an empty thermal image directory.
	thermalImgDir := filepath.Join(imgRoot, "thermal")
	if err := os.MkdirAll(thermalImgDir, 0755); err != nil {
		t.Fatalf("failed to create thermal img dir: %v", err)
	}

	// Create a page WITHOUT markers.
	if err := os.MkdirAll(pagesRoot, 0755); err != nil {
		t.Fatalf("failed to create pages dir: %v", err)
	}

	pageContent := `# Thermal

## Styles

Some style descriptions here.

## Options
`

	pagePath := filepath.Join(pagesRoot, "thermal.md")
	if err := os.WriteFile(pagePath, []byte(pageContent), 0644); err != nil {
		t.Fatalf("failed to write thermal.md: %v", err)
	}

	// Run the gallery tool.
	if err := gallery.Run(imgRoot, pagesRoot); err != nil {
		t.Fatalf("gallery.Run failed: %v", err)
	}

	// Read the page back — it should be completely unchanged.
	got, err := os.ReadFile(pagePath)
	if err != nil {
		t.Fatalf("failed to read thermal.md after gallery run: %v", err)
	}

	if string(got) != pageContent {
		t.Errorf("page was modified when directory is empty and no markers exist.\nWant:\n%s\nGot:\n%s", pageContent, string(got))
	}
}
