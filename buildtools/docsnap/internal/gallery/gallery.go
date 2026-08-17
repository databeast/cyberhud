// Package gallery implements the core logic for the Gallery_Generator tool.
// It reads collected snapshot images and produces markdown gallery sections
// for display mode documentation pages.
package gallery

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/databeast/cyberhud/buildtools/docsnap/internal/category"
)

// dimensionPattern matches a leading "<width>x<height>" token, used to locate
// the resolution embedded in a snapshot filename immediately after its
// category prefix (e.g. "320x240" in "color-320x240-overview_0001.png").
var dimensionPattern = regexp.MustCompile(`^(\d+)x(\d+)`)

// CategoryFromFilename derives the Display_Category from a snapshot filename
// using longest-prefix matching. Returns empty string if no prefix matches.
func CategoryFromFilename(filename string) string {
	cat, ok := category.Match(filename)
	if !ok {
		return ""
	}
	return string(cat)
}

// DimensionsFromFilename extracts width and height from a filename like
// "color-240x320_0001.png" or "color-320x240-overview_0001.png" (where an
// optional style-name suffix follows the dimensions). Returns (0, 0, false)
// if dimensions cannot be parsed.
func DimensionsFromFilename(filename string) (width, height int, ok bool) {
	// Find the longest matching prefix to determine where dimensions start.
	prefixLen := 0
	for _, pm := range category.Prefixes {
		if strings.HasPrefix(filename, pm.Prefix) && len(pm.Prefix) > prefixLen {
			prefixLen = len(pm.Prefix)
		}
	}
	if prefixLen == 0 {
		return 0, 0, false
	}

	// Require the expected snapshot suffix.
	const suffix = "_0001.png"
	if !strings.HasSuffix(filename, suffix) {
		return 0, 0, false
	}

	// Strip the prefix and match the leading <W>x<H> token; anything after
	// it (e.g. "-overview") is an optional style-name suffix and is ignored.
	remainder := filename[prefixLen:]
	m := dimensionPattern.FindStringSubmatch(remainder)
	if m == nil {
		return 0, 0, false
	}

	w, err := strconv.Atoi(m[1])
	if err != nil {
		return 0, 0, false
	}
	h, err := strconv.Atoi(m[2])
	if err != nil {
		return 0, 0, false
	}

	return w, h, true
}

// StyleNameFromFilename extracts the style-name portion (filename without _0001.png suffix).
func StyleNameFromFilename(filename string) string {
	const suffix = "_0001.png"
	if strings.HasSuffix(filename, suffix) {
		return filename[:len(filename)-len(suffix)]
	}
	return filename
}

// imageEntry holds parsed data for a single snapshot image.
type imageEntry struct {
	filename  string
	styleName string
	width     int
	height    int
}

// GenerateGallery produces the markdown gallery section for a mode given its
// list of collected image filenames and the relative path from the page to images.
func GenerateGallery(imageFiles []string, imgRelPath string) string {
	// Group images by category, skipping unrecognized files and unparseable dimensions.
	groups := make(map[string][]imageEntry)
	for _, f := range imageFiles {
		cat := CategoryFromFilename(f)
		if cat == "" {
			continue
		}
		w, h, ok := DimensionsFromFilename(f)
		if !ok {
			continue
		}
		groups[cat] = append(groups[cat], imageEntry{
			filename:  f,
			styleName: StyleNameFromFilename(f),
			width:     w,
			height:    h,
		})
	}

	// If no images were recognized, return empty string.
	if len(groups) == 0 {
		return ""
	}

	// Sort category names alphabetically.
	cats := make([]string, 0, len(groups))
	for c := range groups {
		cats = append(cats, c)
	}
	sort.Strings(cats)

	// Sort images within each category by width ascending, then height ascending.
	for _, c := range cats {
		sort.Slice(groups[c], func(i, j int) bool {
			if groups[c][i].width != groups[c][j].width {
				return groups[c][i].width < groups[c][j].width
			}
			return groups[c][i].height < groups[c][j].height
		})
	}

	// Build the gallery markdown.
	var sb strings.Builder
	sb.WriteString("<!-- snapshot-gallery:start -->\n")
	sb.WriteString("## Snapshots\n")

	for _, cat := range cats {
		entries := groups[cat]
		sb.WriteString("\n")
		sb.WriteString(fmt.Sprintf("### %s\n", cat))

		useGrid := len(entries) > 6
		if useGrid {
			sb.WriteString("\n")
			sb.WriteString("<div class=\"grid\" style=\"display:grid;grid-template-columns:repeat(3,1fr);gap:1rem;\">\n")
		}

		for _, img := range entries {
			sb.WriteString("\n")
			sb.WriteString("<figure>\n")
			src := imgRelPath + img.filename
			alt := fmt.Sprintf("%s %dx%d", img.styleName, img.width, img.height)
			dimCaption := strconv.Itoa(img.width) + "x" + strconv.Itoa(img.height)
			sb.WriteString(fmt.Sprintf("  <img src=\"%s\" alt=\"%s\" style=\"max-width:320px;width:100%%;\">", src, alt))
			sb.WriteString("\n")
			sb.WriteString(fmt.Sprintf("  <figcaption>%s</figcaption>", dimCaption))
			sb.WriteString("\n")
			sb.WriteString("</figure>\n")
		}

		if useGrid {
			sb.WriteString("\n")
			sb.WriteString("</div>\n")
		}
	}

	sb.WriteString("\n")
	sb.WriteString("<!-- snapshot-gallery:end -->\n")

	return sb.String()
}

// InjectGallery inserts or replaces the gallery section in a mode page's content.
// Uses <!-- snapshot-gallery:start --> and <!-- snapshot-gallery:end --> markers.
func InjectGallery(pageContent, gallerySection string) string {
	const startMarker = "<!-- snapshot-gallery:start -->"
	const endMarker = "<!-- snapshot-gallery:end -->"

	startIdx := strings.Index(pageContent, startMarker)
	endIdx := strings.Index(pageContent, endMarker)
	hasMarkers := startIdx >= 0 && endIdx >= 0 && endIdx >= startIdx

	if gallerySection == "" {
		// Empty gallery: remove existing marker block if present.
		if hasMarkers {
			before := pageContent[:startIdx]
			after := pageContent[endIdx+len(endMarker):]
			// Trim a trailing newline after the end marker.
			if strings.HasPrefix(after, "\n") {
				after = after[1:]
			}
			// If before ends with a double newline and after starts with a newline,
			// collapse to avoid triple newlines at the seam.
			if strings.HasSuffix(before, "\n\n") && strings.HasPrefix(after, "\n") {
				before = before[:len(before)-1]
			}
			return before + after
		}
		return pageContent
	}

	if hasMarkers {
		// Replace existing marker block with new gallery section.
		before := pageContent[:startIdx]
		after := pageContent[endIdx+len(endMarker):]
		return before + gallerySection + after
	}

	// No existing markers: append gallery section at end.
	if len(pageContent) > 0 && !strings.HasSuffix(pageContent, "\n") {
		return pageContent + "\n" + gallerySection
	}
	return pageContent + gallerySection
}

// Run executes the gallery generation pipeline: walks the image root for mode
// directories, generates gallery sections, and injects them into mode pages.
func Run(imgRoot, pagesRoot string) error {
	// If the image root doesn't exist, there are no collected images — nothing to do.
	if _, err := os.Stat(imgRoot); err != nil {
		return nil
	}

	// Read immediate subdirectories of imgRoot — these are mode names.
	entries, err := os.ReadDir(imgRoot)
	if err != nil {
		return fmt.Errorf("reading image root %q: %w", imgRoot, err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		mode := entry.Name()

		// Read .png files in imgRoot/<mode>/.
		modeDir := filepath.Join(imgRoot, mode)
		files, err := os.ReadDir(modeDir)
		if err != nil {
			return fmt.Errorf("reading mode directory %q: %w", modeDir, err)
		}

		var imageFiles []string
		for _, f := range files {
			if f.IsDir() {
				continue
			}
			if strings.HasSuffix(f.Name(), ".png") {
				imageFiles = append(imageFiles, f.Name())
			}
		}

		// Skip modes with no collected images for per-page injection.
		if len(imageFiles) == 0 {
			continue
		}

		// Compute relative path from the page to the mode's image directory.
		// MkDocs with directory URLs turns clock.md into clock/index.html,
		// so the relative path needs to go up one level.
		imgRelPath := "../img/" + mode + "/"

		// Generate gallery markdown.
		gallerySection := GenerateGallery(imageFiles, imgRelPath)

		// Determine the mode page path.
		pagePath := filepath.Join(pagesRoot, mode+".md")

		// If the page file doesn't exist, log warning and continue.
		if _, err := os.Stat(pagePath); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "gallery: warning: mode page %q does not exist, skipping\n", pagePath)
			continue
		}

		// Read the page file content.
		pageBytes, err := os.ReadFile(pagePath)
		if err != nil {
			return fmt.Errorf("reading mode page %q: %w", pagePath, err)
		}

		// Inject gallery section into page content.
		updated := InjectGallery(string(pageBytes), gallerySection)

		// Write updated content back to the page file.
		if err := os.WriteFile(pagePath, []byte(updated), 0644); err != nil {
			return fmt.Errorf("writing mode page %q: %w", pagePath, err)
		}
	}

	return nil
}
