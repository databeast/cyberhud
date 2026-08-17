package styles

import (
	"github.com/databeast/cyberhud/display/modes/serial/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/style/layout"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
)

// MonoFast480x320Style renders serial data on a 480×320 monochrome fast panel.
//
// Profile context:
//   - MonoFast: OLED, 1-bit, fast refresh
//   - 480×320: Large landscape resolution (~20 text rows typical)
//   - Target: Header + status + abundant data lines with wide width
//
// Rendering strategy:
//   - Excellent character budget (~60 chars per line)
//   - Many rows available for scrolling data
var MonoFast480x320Style = def{
	name: "mono-fast-480x320",
	reqs: style.SurfaceRequirements{
		MinWidth:   480,
		MinHeight:  320,
		Capability: style.MonoFast,
	},
	p: Params{BuildFn: buildMonoFast480x320},
}

func buildMonoFast480x320(snap source.Snapshot, p source.Policy, ctx style.StyleContext, _ def) style.ViewData {
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
		maxRows = 20
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
