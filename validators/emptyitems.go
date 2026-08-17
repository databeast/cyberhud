package validators

import "github.com/databeast/cyberhud/display/style"

// allEmptyItems returns true if items is nil, empty, or contains only empty strings.
func AllEmptyItems(items []string) bool {
	if len(items) == 0 {
		return true
	}
	for _, item := range items {
		if item != "" {
			return false
		}
	}
	return true
}

// allEmptyItems returns true if items is nil, empty, or contains only empty strings.
func allEmptyItems(items []string) bool {
	if len(items) == 0 {
		return true
	}
	for _, item := range items {
		if item != "" {
			return false
		}
	}
	return true
}

// guardEmptyItems ensures a ViewData always has at least one non-empty item.
// When all snapshot fields are empty, some styles produce Items with only
// empty strings; substitute a clock placeholder so the renderer has content.
func guardEmptyItems(vd *style.ViewData) {
	if allEmptyItems(vd.Items) {
		vd.Items = []string{"--:--"}
	}
}
