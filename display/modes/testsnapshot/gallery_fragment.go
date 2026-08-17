package testsnapshot

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// GalleryEntry represents a single snapshot image for gallery fragment generation.
type GalleryEntry struct {
	Filename string // e.g. "color-240x240_0001.png"
	Width    int
	Height   int
}

// WriteGalleryFragment writes a _gallery.md file into outputDir containing
// HTML <img> tags for each PNG produced by a snapshot test. This fragment is
// later collected and assembled into the documentation gallery page.
//
// The imgRelPath parameter specifies the relative path from the gallery page
// to the image directory (e.g., "display-modes/img/attract_bokeh/").
//
// Entries are sorted by width ascending for consistent output.
func WriteGalleryFragment(outputDir, modeID string, entries []GalleryEntry) error {
	if len(entries) == 0 {
		return nil
	}

	// Sort by width, then height for deterministic output.
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Width != entries[j].Width {
			return entries[i].Width < entries[j].Width
		}
		return entries[i].Height < entries[j].Height
	})

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("<!-- gallery-fragment: %s -->\n", modeID))

	for _, e := range entries {
		alt := fmt.Sprintf("%s %dx%d", modeID, e.Width, e.Height)
		sb.WriteString("<figure>\n")
		sb.WriteString(fmt.Sprintf("  <img src=\"display-modes/img/%s/%s\" alt=\"%s\" style=\"max-width:320px;width:100%%;\">\n", modeID, e.Filename, alt))
		sb.WriteString(fmt.Sprintf("  <figcaption>%s — %dx%d</figcaption>\n", modeID, e.Width, e.Height))
		sb.WriteString("</figure>\n")
	}

	sb.WriteString(fmt.Sprintf("<!-- /gallery-fragment: %s -->\n", modeID))

	fragmentPath := filepath.Join(outputDir, "_gallery.md")
	return os.WriteFile(fragmentPath, []byte(sb.String()), 0o644)
}

// WriteGalleryFragmentFromPaths is a convenience that builds GalleryEntry
// values from a list of PNG paths produced by RenderSnapshot. It parses
// dimensions from filenames using the standard "{basename}_0001.png" pattern.
func WriteGalleryFragmentFromPaths(outputDir, modeID string, pngPaths []string) error {
	var entries []GalleryEntry
	for _, p := range pngPaths {
		filename := filepath.Base(p)
		w, h, ok := parseDimensionsFromFilename(filename)
		if !ok {
			// Include file even without parsed dimensions — use 0x0.
			entries = append(entries, GalleryEntry{Filename: filename})
			continue
		}
		entries = append(entries, GalleryEntry{
			Filename: filename,
			Width:    w,
			Height:   h,
		})
	}
	return WriteGalleryFragment(outputDir, modeID, entries)
}

// parseDimensionsFromFilename extracts width and height from filenames like
// "color-240x240_0001.png" or "compact_128x128_0001.png".
// It scans for the last occurrence of a NxM pattern before "_0001.png".
func parseDimensionsFromFilename(filename string) (w, h int, ok bool) {
	// Strip _NNNN.png suffix.
	base := strings.TrimSuffix(filename, filepath.Ext(filename))
	if idx := strings.LastIndex(base, "_"); idx > 0 {
		base = base[:idx]
	}

	// Find the last segment containing "NxN" by splitting on common separators.
	// Try to find a WxH pattern anywhere in the remaining string.
	parts := strings.FieldsFunc(base, func(r rune) bool {
		return r == '-' || r == '_'
	})

	for i := len(parts) - 1; i >= 0; i-- {
		dims := strings.SplitN(parts[i], "x", 2)
		if len(dims) == 2 {
			var width, height int
			if _, err := fmt.Sscanf(dims[0], "%d", &width); err != nil {
				continue
			}
			if _, err := fmt.Sscanf(dims[1], "%d", &height); err != nil {
				continue
			}
			if width > 0 && height > 0 {
				return width, height, true
			}
		}
	}

	return 0, 0, false
}
