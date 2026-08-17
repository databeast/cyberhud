package thermal

// isAllowed checks if value is present in the allowed list.
func isAllowed(value string, allowed []string) bool {
	for _, a := range allowed {
		if value == a {
			return true
		}
	}
	return false
}

// clampRefreshMS clamps a RefreshMS value to [500, 120000], defaulting to 2000 if zero.
func clampRefreshMS(ms int) int {
	if ms <= 0 {
		ms = 2000
	}
	if ms < 500 {
		ms = 500
	}
	if ms > 120000 {
		ms = 120000
	}
	return ms
}
