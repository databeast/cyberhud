package source

import (
	"github.com/databeast/cyberhud/display/style/layout"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// computeOffsetY determines the vertical centering offset for ticker content.
// Returns 0 when centering should not be applied:
//   - Auto-scroll active and content overflows available rows
//   - Horizontal direction with marquee strips active
//   - Content block height >= available content height (CenterBlockY returns 0 naturally)
func ComputeOffsetY(formatted []FormattedLine, hints textlayout.TextHints, effective Policy) int {
	bridge := layout.NewLayoutBridge(hints, layout.BridgeConfig{})
	maxRows := bridge.MaxVisibleRows()

	// Disable centering when scrolling is active and content overflows.
	if effective.AutoScrollMS > 0 && effective.Direction != textlayout.TickerDirectionNone {
		if len(formatted) > maxRows {
			return 0
		}
	}

	// Disable centering when horizontal direction with marquee strips.
	if effective.Direction == "horizontal" && effective.AutoScrollMS > 0 {
		return 0
	}

	rowHeights := computeRowHeights(formatted, hints)
	spacing := 0 // ticker uses zero inter-row spacing (lines are contiguous)
	return bridge.CenterBlockY(rowHeights, spacing)
}

// computeRowHeights returns per-row pixel heights from formatted lines.
// Each row height is derived from the tier catalog's RowHeight for that line's tier.
func computeRowHeights(formatted []FormattedLine, hints textlayout.TextHints) []int {
	heights := make([]int, len(formatted))
	for i, fl := range formatted {
		face := ResolveFace(hints, "spleen", fl.Tier)
		heights[i] = face.Metrics().RowHeight
	}
	return heights
}
