package validators

// maxItemLen returns the length of the longest string in items.
// Returns 0 for an empty or nil slice.
func maxItemLen(items []string) int {
	max := 0
	for _, s := range items {
		if len(s) > max {
			max = len(s)
		}
	}
	return max
}
