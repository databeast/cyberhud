package styles

import (
	"github.com/databeast/cyberhud/display/modes/serial/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/style/layout"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
)

func buildGrayscaleItemsOnly(snap source.Snapshot, p source.Policy, ctx style.StyleContext, _ def, defaultRows int, static bool) style.ViewData {
	hints := ctx.Hints()
	bridge := layout.NewLayoutBridge(hints, layout.BridgeConfig{PaddingPct: 0})

	if bridge.AvailableContentWidth() <= 0 || bridge.AvailableContentHeight() <= 0 {
		return style.ViewData{Items: []string{"(too small)"}, Static: static}
	}

	entry := ctx.Entry(tiercatalog.TierNormal)
	glyphAdvance := entry.GlyphAdvance
	rowHeight := entry.RowHeight

	maxRows := bridge.MaxVisibleRows()
	if maxRows <= 0 {
		maxRows = p.MaxLines
	}
	if maxRows <= 0 {
		maxRows = defaultRows
	}

	maxChars := 0
	if glyphAdvance > 0 {
		maxChars = bridge.AvailableContentWidth() / glyphAdvance
	}

	ox, _ := bridge.ContentOrigin()

	var items []string
	var tiers []tiercatalog.Tier

	if snap.Port == "" && !snap.AutoSelect && !snap.Connected {
		items = append(items, "No serial port")
		tiers = append(tiers, tiercatalog.TierNormal)
	} else {
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
	}

	dataRowBudget := maxRows - len(items)
	if dataRowBudget < 0 {
		dataRowBudget = 0
	}
	lines := snap.Lines
	if len(lines) > dataRowBudget {
		lines = lines[len(lines)-dataRowBudget:]
	}
	for _, raw := range lines {
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
		Static:      static,
	}
}
