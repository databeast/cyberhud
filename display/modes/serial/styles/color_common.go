package styles

import (
	"image/color"

	"github.com/databeast/cyberhud/display/modes/serial/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/style/layout"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
)

var (
	colorSerialHeader = color.RGBA{R: 0x50, G: 0xC8, B: 0xFF, A: 0xFF}
	colorSerialOK     = color.RGBA{R: 0x40, G: 0xE0, B: 0x80, A: 0xFF}
	colorSerialErr    = color.RGBA{R: 0xFF, G: 0x70, B: 0x70, A: 0xFF}
	colorSerialWarn   = color.RGBA{R: 0xFF, G: 0xC8, B: 0x40, A: 0xFF}
	colorSerialData   = color.RGBA{R: 0xF0, G: 0xF0, B: 0xF0, A: 0xFF}
)

func buildColorItemsOnly(snap source.Snapshot, p source.Policy, ctx style.StyleContext, _ def, defaultRows int, paddingPct int, static bool) style.ViewData {
	hints := ctx.Hints()
	bridge := layout.NewLayoutBridge(hints, layout.BridgeConfig{PaddingPct: paddingPct})

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
	var colors []color.Color

	headerText := compactHeader(snap)
	if maxChars > 0 {
		headerText = textlayout.Truncate(headerText, maxChars)
	}
	items = append(items, headerText)
	tiers = append(tiers, tiercatalog.TierNormal)
	colors = append(colors, colorSerialHeader)

	var statusText string
	statusColor := colorSerialWarn
	if snap.Connected {
		statusText = "OK"
		statusColor = colorSerialOK
	} else if snap.LastError != "" {
		statusText = "ERR: " + snap.LastError
		statusColor = colorSerialErr
	} else {
		statusText = "disconnected"
		statusColor = colorSerialWarn
	}
	if maxChars > 0 {
		statusText = textlayout.Truncate(statusText, maxChars)
	}
	items = append(items, statusText)
	tiers = append(tiers, tiercatalog.TierNormal)
	colors = append(colors, statusColor)

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
		colors = append(colors, colorSerialData)
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
		Colors:      colors,
		Tiers:       tiers,
		LineOffsets: offsets,
		OffsetY:     offsetY,
		Static:      static,
	}
}
