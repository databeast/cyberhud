package cycle

import "github.com/databeast/cyberhud/display/coordinator"

// advanceSkippingSelf advances the mode on panelIndex, skipping "cycle" itself
// and any mode not in allowedModes (when non-empty). The loop is bounded to at
// most len(status.Modes) iterations to guarantee termination.
func advanceSkippingSelf(panelIndex int, state *coordinator.State, allowedModes []string) string {
	status, ok := state.Region(panelIndex)
	if !ok {
		return ""
	}

	maxAttempts := len(status.Modes)
	for i := 0; i < maxAttempts; i++ {
		next, err := state.Next(panelIndex)
		if err != nil {
			return ""
		}
		if next == ModeID {
			continue // skip self
		}
		if len(allowedModes) == 0 {
			return next // no filter — any non-cycle mode is acceptable
		}
		if contains(allowedModes, next) {
			return next // matches configured mode list
		}
	}

	// Exhausted all modes — no valid target found; leave current unchanged.
	return state.CurrentMode(panelIndex)
}

// contains reports whether haystack contains needle.
func contains(haystack []string, needle string) bool {
	for _, s := range haystack {
		if s == needle {
			return true
		}
	}
	return false
}
