package styles

import (
	"fmt"

	"github.com/databeast/cyberhud/display/modes/wifi/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/style/layout"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
)

// WiFi style declarations for 128x64 mono panels.
//
// This panel class has room for one prominent SSID row and a few compact status
// rows, so it gets a dedicated mono-fast layout instead of falling back to the
// 128x128 family.

var MonoSlow128x64Style = def{
	name: "mono-slow-128x64",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 64, Capability: style.MonoSlow},
}

var MonoFast128x64Style = def{
	name: "mono-fast-128x64",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 64, Capability: style.MonoFast},
	p:    Params{BuildFn: monoFast128x64Build},
}

var GrayscaleSlow128x64Style = def{
	name: "grayscale-slow-128x64",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 64, Capability: style.GrayscaleSlow},
}

var GrayscaleFast128x64Style = def{
	name: "grayscale-fast-128x64",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 64, Capability: style.GrayscaleFast},
	p:    Params{Fast: true},
}

var ColorSlow128x64Style = def{
	name: "color-slow-128x64",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 64, Capability: style.ColorSlow},
}

var ColorFast128x64Style = def{
	name: "color-fast-128x64",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 64, Capability: style.ColorFast},
	p:    Params{Fast: true, Color: true},
}

func monoFast128x64Build(snap source.WifiData, pol source.Policy, ctx style.StyleContext, d def) style.ViewData {
	const paddingPct = 2

	bridge := ctx.Layout(paddingPct)
	if bridge.AvailableContentWidth() == 0 || bridge.AvailableContentHeight() == 0 {
		return style.ViewData{
			Items:      []string{"(too small)"},
			Cursor:     -1,
			PaddingPct: paddingPct,
		}
	}

	largeEntry := ctx.Entry(tiercatalog.TierLarge)
	detailEntry := ctx.Entry(tiercatalog.TierSmall)
	if snap.ConnectionState != source.Connected {
		return monoFast128x64StatusView(snap, bridge, largeEntry, paddingPct)
	}

	type row struct {
		text string
		tier tiercatalog.Tier
		rowH int
		adv  int
	}

	rows := []row{
		{
			text: compactSSID(snap.SSID),
			tier: tiercatalog.TierLarge,
			rowH: largeEntry.RowHeight,
			adv:  largeEntry.GlyphAdvance,
		},
		{
			text: source.FormatSignalText(snap, pol),
			tier: tiercatalog.TierSmall,
			rowH: detailEntry.RowHeight,
			adv:  detailEntry.GlyphAdvance,
		},
	}

	if snap.IPAddress != "" {
		rows = append(rows, row{
			text: snap.IPAddress,
			tier: tiercatalog.TierSmall,
			rowH: detailEntry.RowHeight,
			adv:  detailEntry.GlyphAdvance,
		})
	} else if snap.LinkQuality > 0 {
		rows = append(rows, row{
			text: fmt.Sprintf("Link %d%%", snap.LinkQuality),
			tier: tiercatalog.TierSmall,
			rowH: detailEntry.RowHeight,
			adv:  detailEntry.GlyphAdvance,
		})
	}

	if detail := compactRadioLine(snap, pol); detail != "" {
		rows = append(rows, row{
			text: detail,
			tier: tiercatalog.TierSmall,
			rowH: detailEntry.RowHeight,
			adv:  detailEntry.GlyphAdvance,
		})
	}

	if pol.ShowInterface && snap.InterfaceName != "" {
		rows = append(rows, row{
			text: snap.InterfaceName,
			tier: tiercatalog.TierSmall,
			rowH: detailEntry.RowHeight,
			adv:  detailEntry.GlyphAdvance,
		})
	} else if snap.LinkSpeed > 0 {
		rows = append(rows, row{
			text: fmt.Sprintf("%dMbps", snap.LinkSpeed),
			tier: tiercatalog.TierSmall,
			rowH: detailEntry.RowHeight,
			adv:  detailEntry.GlyphAdvance,
		})
	}

	// Keep the 64px panel from over-stacking detail rows. The first three rows
	// carry the most useful information and leave enough vertical room for the
	// SSID line to read cleanly.
	if len(rows) > 3 {
		rows = rows[:3]
	}

	rowHeights := make([]int, len(rows))
	for i, r := range rows {
		rowHeights[i] = r.rowH
	}
	spacing, offsetY, visibleCount := bridge.FitRows(rowHeights)
	if visibleCount <= 0 {
		visibleCount = 1
	}
	if visibleCount > len(rows) {
		visibleCount = len(rows)
	}

	items := make([]string, visibleCount)
	tiers := make([]tiercatalog.Tier, visibleCount)
	visibleHeights := make([]int, visibleCount)
	for i := 0; i < visibleCount; i++ {
		items[i] = truncateToWidth(rows[i].text, bridge.AvailableContentWidth(), rows[i].adv)
		tiers[i] = rows[i].tier
		visibleHeights[i] = rows[i].rowH
	}

	return style.ViewData{
		Items:        items,
		Tiers:        tiers,
		RowHeights:   visibleHeights,
		Spacing:      spacing,
		VisibleCount: visibleCount,
		OffsetY:      offsetY,
		Cursor:       -1,
		PaddingPct:   paddingPct,
	}
}

func monoFast128x64StatusView(
	snap source.WifiData,
	bridge layout.LayoutCalculator,
	largeEntry tiercatalog.Entry,
	paddingPct int,
) style.ViewData {
	msg := snap.StatusMessage
	if msg == "" {
		if snap.ConnectionState == source.Disconnected {
			msg = "No Network"
		} else {
			msg = "WiFi N/A"
		}
	}
	msg = truncateToWidth(msg, bridge.AvailableContentWidth(), largeEntry.GlyphAdvance)
	offsetX := bridge.CenterXWith(len([]rune(msg)), largeEntry.GlyphAdvance)
	offsetY := bridge.CenterBlockY([]int{largeEntry.RowHeight}, 0)

	return style.ViewData{
		Items:       []string{msg},
		Tiers:       []tiercatalog.Tier{tiercatalog.TierLarge},
		LineOffsets: []int{offsetX},
		RowHeights:  []int{largeEntry.RowHeight},
		OffsetY:     offsetY,
		Cursor:      -1,
		PaddingPct:  paddingPct,
	}
}

func compactSSID(ssid string) string {
	if ssid == "" {
		return "(no SSID)"
	}
	return ssid
}

func compactRadioLine(snap source.WifiData, pol source.Policy) string {
	switch {
	case pol.ShowFrequency && snap.Frequency > 0 && pol.ShowChannel && snap.Channel > 0:
		return fmt.Sprintf("%.1fG Ch%d", snap.Frequency, snap.Channel)
	case pol.ShowFrequency && snap.Frequency > 0:
		return fmt.Sprintf("%.1fGHz", snap.Frequency)
	case pol.ShowChannel && snap.Channel > 0:
		return fmt.Sprintf("Ch %d", snap.Channel)
	default:
		return ""
	}
}

func truncateToWidth(text string, width, advance int) string {
	if text == "" || width <= 0 || advance <= 0 {
		return text
	}
	maxChars := width / advance
	if maxChars <= 0 {
		return text
	}
	return textlayout.Truncate(text, maxChars)
}
