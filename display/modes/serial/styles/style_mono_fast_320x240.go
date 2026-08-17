package styles

import (
	"github.com/databeast/cyberhud/display/modes/serial/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/style/layout"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
)

// MonoFast320x240Style renders serial data on a 320×240 monochrome fast panel.
//
// Profile context:
//   - MonoFast: OLED, 1-bit, fast refresh
//   - 320×240: Common VGA-like resolution (~15 text rows typical)
//   - Target: Header + status + many data lines with good width
//
// Rendering strategy:
//   - Excellent character budget (~40 chars per line)
//   - Ample row count for scrolling data
var MonoFast320x240Style = def{
	name: "mono-fast-320x240",
	reqs: style.SurfaceRequirements{
		MinWidth:   320,
		MinHeight:  240,
		Capability: style.MonoFast,
	},
	p: Params{BuildFn: buildMonoFast320x240},
}

func buildMonoFast320x240(snap source.Snapshot, p source.Policy, ctx style.StyleContext, _ def) style.ViewData {
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
		maxRows = 15
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
