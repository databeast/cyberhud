package style

// Shared view-data and list-geometry helpers used across display modes.
//
// These were previously scattered as byte-identical private copies in a dozen
// modes (allEmptyItems, clampList, longestItemLen, guardEmptyItems) plus a
// single-importer `validators` package. They live here because they operate on
// ViewData and list geometry that this package already owns, and every mode
// already imports style — so no mode grows a new dependency to use them.

// AllEmptyItems reports whether items is nil, empty, or contains only empty
// strings. Styles use it to decide whether to substitute placeholder content.
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

// GuardEmptyItems ensures vd.Items is never empty: when it has length 0, it is
// replaced with a single-element slice containing one empty string, giving the
// renderer a row to work with. A non-empty Items slice is left unchanged.
func GuardEmptyItems(vd *ViewData) {
	if len(vd.Items) == 0 {
		vd.Items = []string{""}
	}
}

// MaxItemLen returns the length (in bytes) of the longest string in items,
// or 0 for a nil or empty slice. Used for width-based layout decisions.
func MaxItemLen(items []string) int {
	max := 0
	for _, s := range items {
		if len(s) > max {
			max = len(s)
		}
	}
	return max
}

// ClampList constrains cursor and topRow so the cursor stays within
// [0, itemCount) and within the visible window of maxVisible rows. Used by
// modes with scrollable lists (menu, serial, stemma, gpio, gpio-control).
func ClampList(cursor, topRow, maxVisible, itemCount int) (int, int) {
	if itemCount <= 0 {
		return 0, 0
	}
	if cursor < 0 {
		cursor = 0
	}
	if cursor >= itemCount {
		cursor = itemCount - 1
	}
	if cursor < topRow {
		topRow = cursor
	}
	if maxVisible > 0 && cursor >= topRow+maxVisible {
		topRow = cursor - maxVisible + 1
	}
	if topRow < 0 {
		topRow = 0
	}
	return cursor, topRow
}
