package styles

import (
	"github.com/databeast/cyberhud/display/modes/serial/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/style/layout"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
)

// MonoFast128x64Style renders serial data on a 128×64 monochrome fast (OLED) panel.
//
// Profile context:
//   - MonoFast: OLED, 1-bit, fast refresh (no flicker concerns)
//   - 128×64: Small vertical space, ~4 text rows typical
//   - Target: Compact header + status + 2–3 data lines
//
// Rendering strategy:
//   - Compact header: "Port @Baud" (left-aligned)
//   - Status line: "OK" or "disconnected" or error message
//   - Data lines: Newest 2–3 lines (if space permits)
//   - All text ANSI-stripped, truncated to fit width
//   - No decoration (space tight but adequate)
var MonoFast128x64Style = def{
	name: "mono-fast-128x64",
	reqs: style.SurfaceRequirements{
		MinWidth:   128,
		MinHeight:  64,
		Capability: style.MonoFast,
	},
	p: Params{BuildFn: buildMonoFast128x64},
}

// buildMonoFast128x64 is the render function for mono-fast-128x64 style.
// Renders on a 128×64 OLED panel (typically 4 rows max).
func buildMonoFast128x64(snap source.Snapshot, p source.Policy, ctx style.StyleContext, _ def) style.ViewData {
	// Step 1: Get panel hints and layout bridge.
	hints := ctx.Hints()
	bridge := layout.NewLayoutBridge(hints, layout.BridgeConfig{PaddingPct: 0})

	// Step 2: Guard against impossible dimensions.
	if bridge.AvailableContentWidth() <= 0 || bridge.AvailableContentHeight() <= 0 {
		return style.ViewData{Items: []string{"(too small)"}, Static: true}
	}

	// Step 3: Resolve tier and font metrics.
	entry := ctx.Entry(tiercatalog.TierNormal)
	glyphAdvance := entry.GlyphAdvance
	rowHeight := entry.RowHeight

	// Step 4: Calculate layout constraints.
	maxRows := bridge.MaxVisibleRows()
	if maxRows <= 0 {
		maxRows = p.MaxLines
	}
	if maxRows <= 0 {
		maxRows = 4
	}

	maxChars := 0
	if glyphAdvance > 0 {
		maxChars = bridge.AvailableContentWidth() / glyphAdvance
	}

	ox, _ := bridge.ContentOrigin()

	// Step 5: Build items list.
	var items []string
	var tiers []tiercatalog.Tier

	// Row 0: Compact header (port @baud).
	headerText := compactHeader(snap)
	if maxChars > 0 {
		headerText = textlayout.Truncate(headerText, maxChars)
	}
	items = append(items, headerText)
	tiers = append(tiers, tiercatalog.TierNormal)

	// Row 1: Status (OK / disconnected / error).
	var statusText string
	if snap.Connected {
		statusText = "OK"
	} else if snap.LastError != "" {
		statusText = "ERR: " + snap.LastError
	} else {
		statusText = "disconnected"
	}
	if maxChars > 0 {
		statusText = textlayout.Truncate(statusText, maxChars)
	}
	items = append(items, statusText)
	tiers = append(tiers, tiercatalog.TierNormal)

	// Rows 2+: Data lines (ANSI-stripped).
	dataRowBudget := maxRows - len(items)
	if dataRowBudget < 0 {
		dataRowBudget = 0
	}
	dataLines := snap.Lines
	if len(dataLines) > dataRowBudget {
		dataLines = dataLines[len(dataLines)-dataRowBudget:]
	}
	for _, raw := range dataLines {
		text, _ := source.ParseLine(raw)
		if maxChars > 0 {
			text = textlayout.Truncate(text, maxChars)
		}
		items = append(items, text)
		tiers = append(tiers, tiercatalog.TierNormal)
	}

	// Step 6: Calculate X offsets (left-align all rows).
	offsets := make([]int, len(items))
	for i := range offsets {
		offsets[i] = ox
	}

	// Step 7: Calculate row heights and fit vertically.
	rowHeights := make([]int, len(items))
	for i := range rowHeights {
		rowHeights[i] = rowHeight
	}
	_, offsetY, _ := bridge.FitRows(rowHeights)

	// Step 8: Return ViewData.
	return style.ViewData{
		Items:       items,
		Tiers:       tiers,
		LineOffsets: offsets,
		OffsetY:     offsetY,
	}
}
