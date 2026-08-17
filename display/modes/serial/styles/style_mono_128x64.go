package styles

import (
	"fmt"

	"github.com/databeast/cyberhud/display/modes/serial/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/style/layout"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
)

// MonoSlow128x64Style is the preferred style for 128×64 monochrome slow panels.
// It uses an items-only layout (no sprite path) for clean, predictable rendering
// on the tiny screen: a compact header row, a status row, then serial data lines.
var MonoSlow128x64Style = def{
	name: "mono-slow-128x64",
	reqs: style.SurfaceRequirements{
		MinWidth:   128,
		MinHeight:  64,
		Capability: style.MonoSlow,
	},
	p: Params{BuildFn: buildMonoSlow128x64},
}

func buildMonoSlow128x64(snap source.Snapshot, p source.Policy, ctx style.StyleContext, _ def) style.ViewData {
	hints := ctx.Hints()
	bridge := layout.NewLayoutBridge(hints, layout.BridgeConfig{PaddingPct: 0})

	entry := ctx.Entry(tiercatalog.TierNormal)
	glyphAdvance := entry.GlyphAdvance
	rowHeight := entry.RowHeight

	maxRows := bridge.MaxVisibleRows()
	maxChars := 0
	if glyphAdvance > 0 {
		maxChars = bridge.AvailableContentWidth() / glyphAdvance
	}

	if maxRows <= 0 {
		maxRows = p.MaxLines
	}
	if maxRows <= 0 {
		maxRows = 6
	}

	// Left-align helper: content origin X.
	ox, _ := bridge.ContentOrigin()

	// No port configured at all — single idle row.
	if snap.Port == "" && !snap.AutoSelect && !snap.Connected {
		return style.ViewData{
			Items:       []string{"No serial port"},
			Tiers:       []tiercatalog.Tier{tiercatalog.TierNormal},
			LineOffsets: []int{ox},
			OffsetY:     (bridge.AvailableContentHeight() - rowHeight) / 2,
		}
	}

	// Row 0: compact port+baud header.
	headerText := compactHeader(snap)
	if maxChars > 0 {
		headerText = textlayout.Truncate(headerText, maxChars)
	}

	// Row 1: status.
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

	allItems := []string{headerText, statusText}
	allTiers := []tiercatalog.Tier{tiercatalog.TierNormal, tiercatalog.TierNormal}

	// Data rows: ANSI-stripped serial lines, newest last, up to available rows.
	dataRowBudget := maxRows - 2
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
		allItems = append(allItems, text)
		allTiers = append(allTiers, tiercatalog.TierNormal)
	}

	// Left-align all rows.
	offsets := make([]int, len(allItems))
	for i := range offsets {
		offsets[i] = ox
	}

	rowHeights := make([]int, len(allItems))
	for i := range rowHeights {
		rowHeights[i] = rowHeight
	}
	_, offsetY, _ := bridge.FitRows(rowHeights)

	return style.ViewData{
		Items:       allItems,
		Tiers:       allTiers,
		LineOffsets: offsets,
		OffsetY:     offsetY,
		Static:      true,
	}
}

// compactHeader returns a short port+baud string for the 128×64 header row.
func compactHeader(snap source.Snapshot) string {
	port := snap.Port
	if port == "" {
		if snap.AutoSelect {
			port = "auto"
		} else {
			port = "(none)"
		}
	}
	return fmt.Sprintf("%s @%d", port, snap.Baud)
}
