package styles

import (
	"github.com/databeast/cyberhud/display/modes/serial/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/style/layout"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
)

// MonoFast128x32Style renders serial data on a 128×32 monochrome fast (OLED) panel.
//
// Profile context:
//   - MonoFast: OLED, 1-bit, fast refresh (no flicker concerns)
//   - 128×32: Tiny vertical space, ~2–3 text rows max
//   - Target: Tightly packed header + status, minimal data lines
//
// Rendering strategy:
//   - Compact header: "Port @Baud" (left-aligned)
//   - Status line: "OK" or "disconnected" or error message
//   - Data lines: Only newest 1 line (if space permits)
//   - All text ANSI-stripped, truncated to fit width
//   - No decoration or widgets (space too tight)
var MonoFast128x32Style = def{
	name: "mono-fast-128x32",
	reqs: style.SurfaceRequirements{
		MinWidth:   128,
		MinHeight:  32,
		Capability: style.MonoFast,
	},
	p: Params{BuildFn: buildMonoFast128x32},
}

// buildMonoFast128x32 is the render function for mono-fast-128x32 style.
// Renders on a 128×32 OLED panel (typically 2–3 rows max).
func buildMonoFast128x32(snap source.Snapshot, p source.Policy, ctx style.StyleContext, _ def) style.ViewData {
	// Step 1: Get panel hints and layout bridge.
	hints := ctx.Hints()
	bridge := layout.NewLayoutBridge(hints, layout.BridgeConfig{PaddingPct: 0})

	// Step 2: Guard against impossible dimensions.
	if bridge.AvailableContentWidth() <= 0 || bridge.AvailableContentHeight() <= 0 {
		return style.ViewData{Items: []string{"(too small)"}, Static: true}
	}

	// Step 3: Resolve tier and font metrics.
	// 128×32 is tiny; we get only TierNormal from the catalog.
	entry := ctx.Entry(tiercatalog.TierNormal)
	glyphAdvance := entry.GlyphAdvance
	rowHeight := entry.RowHeight

	// Step 4: Calculate layout constraints.
	maxRows := bridge.MaxVisibleRows()
	if maxRows <= 0 {
		maxRows = p.MaxLines
	}
	if maxRows <= 0 {
		maxRows = 2 // Default for 128×32: ~2 rows
	}

	maxChars := 0
	if glyphAdvance > 0 {
		maxChars = bridge.AvailableContentWidth() / glyphAdvance
	}

	ox, _ := bridge.ContentOrigin() // Left-align X offset

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
	// For 128×32 (only 2–3 rows total), we have minimal budget.
	dataRowBudget := maxRows - len(items)
	if dataRowBudget < 0 {
		dataRowBudget = 0
	}
	dataLines := snap.Lines
	if len(dataLines) > dataRowBudget {
		dataLines = dataLines[len(dataLines)-dataRowBudget:]
	}
	for _, raw := range dataLines {
		text, _ := source.ParseLine(raw) // Strip ANSI codes
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
	// bridge.FitRows() centers the block vertically within available space.
	rowHeights := make([]int, len(items))
	for i := range rowHeights {
		rowHeights[i] = rowHeight
	}
	_, offsetY, _ := bridge.FitRows(rowHeights)

	// Step 8: Return ViewData.
	// Note: No Sprites or Static flag needed for MonoFast (OLED, fast refresh).
	return style.ViewData{
		Items:       items,
		Tiers:       tiers,
		LineOffsets: offsets,
		OffsetY:     offsetY,
	}
}
