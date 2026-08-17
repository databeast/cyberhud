package source

// MenuSnapshot captures the state needed by all menu styles for rendering.
// Items are the menu entry strings, Cursor is the highlighted index,
// TopRow is the first visible row for scroll windowing.
type MenuSnapshot struct {
	Items  []string
	Cursor int
	TopRow int
}

// BuildData returns the menu snapshot consumed by styles.
func BuildData(cursor, topRow int) MenuSnapshot {
	return MenuSnapshot{
		Items:  Items(),
		Cursor: cursor,
		TopRow: topRow,
	}
}

// AllEmptyItems returns true if items is nil, empty, or contains only empty strings.
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
