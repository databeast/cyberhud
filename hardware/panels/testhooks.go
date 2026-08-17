package panels

// TestFilterExcluded exposes filterExcluded to the external tests package.
func TestFilterExcluded(modes []string, excluded []string) []string {
	return filterExcluded(modes, excluded)
}
