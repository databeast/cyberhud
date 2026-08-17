package validators

// clampList constrains cursor and topRow so the cursor stays within
// [0, itemCount) and within the visible window of maxVisible rows.
// Used by modes with scrollable lists (menu, serial, stemma, gpio, gpio-control).
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
