package styles

import (
	"github.com/databeast/cyberhud/display/modes/serial/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/style/layout"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
)

// MonoFast128x128Style renders serial data on a 128×128 monochrome fast (OLED) panel.
//
// Profile context:
//   - MonoFast: OLED, 1-bit, fast refresh
//   - 128×128: Square format, ~8 text rows typical
//   - Target: Compact header + status + 6+ data lines
//
// Rendering strategy:
//   - Compact header + status + multiple data lines
//   - All text ANSI-stripped, truncated to fit width
var MonoFast128x128Style = def{
	name: "mono-fast-128x128",
	reqs: style.SurfaceRequirements{
		MinWidth:   128,
		MinHeight:  128,
		Capability: style.MonoFast,
	},
	p: Params{BuildFn: buildMonoFast128x128},
}

func buildMonoFast128x128(snap source.Snapshot, p source.Policy, ctx style.StyleContext, _ def) style.ViewData {
	hints := ctx.Hints()
	bridge := layout.NewLayoutBridge(hints, layout.BridgeConfig{PaddingPct: 0})

	if bridge.AvailableContentWidth() <= 0 || bridge.AvailableContentHeight() <= 0 {
		return style.ViewData{Items: []string{"(too small)"}}
	}

	entry := ctx.Entry(tiercatalog.TierNormal)
	glyphAdvance := entry.GlyphAdvance
	rowHeight := entry.RowHeight

	maxRows := bridge.MaxVisibleRows()
	if maxRows <= 0 {
		maxRows = p.MaxLines
	}
	if maxRows <= 0 {
		maxRows = 8
	}

	maxChars := 0
	if glyphAdvance > 0 {
		maxChars = bridge.AvailableContentWidth() / glyphAdvance
	}

	ox, _ := bridge.ContentOrigin()

	var items []string
	var tiers []tiercatalog.Tier

	headerText := compactHeader(snap)
	if maxChars > 0 {
		headerText = textlayout.Truncate(headerText, maxChars)
	}
	items = append(items, headerText)
	tiers = append(tiers, tiercatalog.TierNormal)

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

	offsets := make([]int, len(items))
	for i := range offsets {
		offsets[i] = ox
	}

	rowHeights := make([]int, len(items))
	for i := range rowHeights {
		rowHeights[i] = rowHeight
	}
	_, offsetY, _ := bridge.FitRows(rowHeights)

	return style.ViewData{
		Items:       items,
		Tiers:       tiers,
		LineOffsets: offsets,
		OffsetY:     offsetY,
	}
}
