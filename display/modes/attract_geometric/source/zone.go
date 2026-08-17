package source

// centralZoneOpacityCap returns the opacity cap for elements at horizontal
// position x. Elements in the central 60% of the panel (between 20% and 80%
// of panel width) are capped at 0.4 to keep overlay text readable. Elements
// outside this zone return 1.0 (no cap).
func centralZoneOpacityCap(x float64, panelWidth int) float64 {
	pw := float64(panelWidth)
	leftBound := pw * 0.2
	rightBound := pw * 0.8
	if x >= leftBound && x <= rightBound {
		return 0.4
	}
	return 1.0
}
