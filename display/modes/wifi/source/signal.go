package source

import (
	"fmt"
	"image/color"
)

// Signal bar level thresholds (dBm boundaries).
// signalToBarLevel maps a dBm value to a discrete bar level (0–4).
//
// Mapping:
//
//	dBm < -80  → 0 bars
//	-80 ≤ dBm < -70 → 1 bar
//	-70 ≤ dBm < -60 → 2 bars
//	-60 ≤ dBm < -50 → 3 bars
//	dBm ≥ -50 → 4 bars
func SignalToBarLevel(dBm int) int {
	switch {
	case dBm >= -50:
		return 4
	case dBm >= -60:
		return 3
	case dBm >= -70:
		return 2
	case dBm >= -80:
		return 1
	default:
		return 0
	}
}

// formatDbm formats a dBm value as a string in the form "-NNdBm".
// No decimal places are included.
func FormatDbm(dBm int) string {
	return fmt.Sprintf("%ddBm", dBm)
}

// signalPercent converts a dBm value to a percentage (0–100) using linear
// interpolation from -100 dBm (0%) to -30 dBm (100%), clamped at both bounds.
//
// Formula: clamp((dBm + 100) * 100 / 70, 0, 100)
func SignalPercent(dBm int) int {
	pct := (dBm + 100) * 100 / 70
	if pct < 0 {
		return 0
	}
	if pct > 100 {
		return 100
	}
	return pct
}

// qualityColor returns the quality-mapped color for the given signal bar level.
//
// Mapping:
//
//	levels 3–4 → accent color (passed in)
//	level 2    → amber (255, 191, 0, 255)
//	levels 0–1 → red (255, 0, 0, 255)
func QualityColor(barLevel int, accent color.RGBA) color.RGBA {
	switch {
	case barLevel >= 3:
		return accent
	case barLevel == 2:
		return color.RGBA{255, 191, 0, 255}
	default:
		return color.RGBA{255, 0, 0, 255}
	}
}

// formatSignalText returns a text-only representation of signal strength
// for grayscale-fast styles that cannot render bar sprites or color-coded indicators.
// Output varies based on the snapshot's SignalDisplay policy setting.
func FormatSignalText(snap WifiData, pol Policy) string {
	switch pol.SignalDisplay {
	case "dbm":
		return FormatDbm(snap.SignalStrength)
	case "percentage", "percent":
		return fmt.Sprintf("%d%%", SignalPercent(snap.SignalStrength))
	default: // "bars" — show as text level since bar sprites aren't available
		level := SignalToBarLevel(snap.SignalStrength)
		return fmt.Sprintf("Signal: %d/4", level)
	}
}
