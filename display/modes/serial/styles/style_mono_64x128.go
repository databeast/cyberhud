package styles

import (
	"fmt"

	"github.com/databeast/cyberhud/display/modes/serial/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/style/layout"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
)

// MonoSlow64x128Style is the preferred style for 64×128 monochrome slow panels.
// It uses a compact items-only layout: header, status, then serial lines.
var MonoSlow64x128Style = def{
	name: "mono-slow-64x128",
	reqs: style.SurfaceRequirements{
		MinWidth:   64,
		MinHeight:  128,
		Capability: style.MonoSlow,
	},
	p: Params{BuildFn: buildMonoSlow64x128},
}

func buildMonoSlow64x128(snap source.Snapshot, p source.Policy, ctx style.StyleContext, _ def) style.ViewData {
	hints := ctx.Hints()
	bridge := layout.NewLayoutBridge(hints, layout.BridgeConfig{PaddingPct: 0})

	entry := ctx.Entry(tiercatalog.TierNormal)
	glyphAdvance := entry.GlyphAdvance
	rowHeight := entry.RowHeight

	maxRows := bridge.MaxVisibleRows()
	if maxRows <= 0 {
		maxRows = p.MaxLines
	}
	if maxRows <= 0 {
		maxRows = 6
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
		statusText = fmt.Sprintf("ERR: %s", snap.LastError)
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
		Static:      true,
	}
}
