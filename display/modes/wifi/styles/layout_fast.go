package styles

import (
	"fmt"
	"image"
	"image/color"

	"github.com/databeast/cyberhud/display/modes/wifi/source"
	"github.com/databeast/cyberhud/display/style"
	sharedcolor "github.com/databeast/cyberhud/display/style/color"
	"github.com/databeast/cyberhud/display/style/layout"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
	"github.com/databeast/cyberhud/display/widgets"
	"github.com/databeast/cyberhud/display/widgets/progressbar"
)

// buildFastWifi implements the fast-refresh layout engine for WiFi styles.
// It renders WiFi state with graphical elements (icons, signal bars, progress
// bars) suitable for continuously-refreshed panels (TFT, OLED, fast grayscale).
//
// Connected state: WiFi icon sprite (16×16) + SSID header + signal bars sprite
// + signal label + link quality progress bar + detail rows (frequency, channel,
// speed, IP, interface).
//
// Disconnected state: centered "No Network" in warning color + dimmed WiFi icon.
//
// Unavailable state: centered "WiFi N/A" in dimmed color.
//
// Uses style.TierForHeight for adaptive tier selection and LayoutBridge.FitRows
// to drop detail rows from the bottom when height is insufficient.
// Produces ViewData.Static = false.
func buildFastWifi(data source.WifiData, pol source.Policy, ctx style.StyleContext, d def) style.ViewData {
	const paddingPct = 2

	bridge := ctx.Layout(paddingPct)
	if bridge.AvailableContentWidth() <= 0 || bridge.AvailableContentHeight() <= 0 {
		return style.ViewData{Items: []string{"(too small)"}, Static: false, PaddingPct: paddingPct}
	}

	// Resolve accent color based on Color flag.
	accent := color.RGBA{255, 255, 255, 255} // white default for non-color
	if d.p.Color {
		accent = resolveFGColor(pol.FGColor)
	}

	switch data.ConnectionState {
	case source.Disconnected:
		return buildFastWifiDisconnected(bridge, ctx, d, accent, paddingPct)
	case source.Unavailable:
		return buildFastWifiUnavailable(bridge, ctx, d, paddingPct)
	default:
		return buildFastWifiConnected(data, pol, bridge, ctx, d, accent, paddingPct)
	}
}

// buildFastWifiDisconnected renders the Disconnected state:
// centered "No Network" in warning color + dimmed WiFi icon.
func buildFastWifiDisconnected(bridge layout.LayoutCalculator, ctx style.StyleContext, d def, accent color.RGBA, paddingPct int) style.ViewData {
	msg := "No Network"

	primaryTier := style.TierForHeight(d.reqs.MinHeight)
	entry := ctx.Entry(primaryTier)

	// Truncate message to fit.
	if entry.GlyphAdvance > 0 {
		maxChars := bridge.AvailableContentWidth() / entry.GlyphAdvance
		if maxChars > 0 {
			msg = textlayout.Truncate(msg, maxChars)
		}
	}

	offsetX := bridge.CenterXWith(len([]rune(msg)), entry.GlyphAdvance)
	offsetY := bridge.CenterBlockY([]int{entry.RowHeight}, 0)

	// Warning color: red for disconnected.
	warningColor := color.RGBA{255, 60, 60, 255}

	// WiFi icon in dimmed state.
	var sprites []widgets.Sprite
	iconImg := renderWifiIcon(source.Disconnected, accent)
	if iconImg != nil {
		ox, oy := bridge.ContentOrigin()
		sprites = append(sprites, widgets.Sprite{
			Image:    iconImg,
			Position: image.Point{X: ox, Y: oy},
			Label:    "wifi-icon",
		})
	}

	var colors []color.Color
	if d.p.Color {
		colors = []color.Color{warningColor}
	}

	return style.ViewData{
		Items:       []string{msg},
		Tiers:       []tiercatalog.Tier{primaryTier},
		LineOffsets: []int{offsetX},
		OffsetY:     offsetY,
		RowHeights:  []int{entry.RowHeight},
		Sprites:     sprites,
		Colors:      colors,
		Static:      false,
		PaddingPct:  paddingPct,
	}
}

// buildFastWifiUnavailable renders the Unavailable state:
// centered "WiFi N/A" in dimmed color.
func buildFastWifiUnavailable(bridge layout.LayoutCalculator, ctx style.StyleContext, d def, paddingPct int) style.ViewData {
	msg := "WiFi N/A"

	primaryTier := style.TierForHeight(d.reqs.MinHeight)
	entry := ctx.Entry(primaryTier)

	// Truncate message to fit.
	if entry.GlyphAdvance > 0 {
		maxChars := bridge.AvailableContentWidth() / entry.GlyphAdvance
		if maxChars > 0 {
			msg = textlayout.Truncate(msg, maxChars)
		}
	}

	offsetX := bridge.CenterXWith(len([]rune(msg)), entry.GlyphAdvance)
	offsetY := bridge.CenterBlockY([]int{entry.RowHeight}, 0)

	// Dimmed color for unavailable state.
	dimmedColor := color.RGBA{80, 80, 80, 255}

	var colors []color.Color
	if d.p.Color {
		colors = []color.Color{dimmedColor}
	}

	return style.ViewData{
		Items:       []string{msg},
		Tiers:       []tiercatalog.Tier{primaryTier},
		LineOffsets: []int{offsetX},
		OffsetY:     offsetY,
		RowHeights:  []int{entry.RowHeight},
		Colors:      colors,
		Static:      false,
		PaddingPct:  paddingPct,
	}
}

// buildFastWifiConnected renders the Connected state with full graphical elements:
// WiFi icon + SSID header + signal bars + signal label + link quality progress bar
// + detail rows (frequency, channel, speed, IP, interface).
func buildFastWifiConnected(data source.WifiData, pol source.Policy, bridge layout.LayoutCalculator, ctx style.StyleContext, d def, accent color.RGBA, paddingPct int) style.ViewData {
	// Tier selection based on panel height.
	primaryTier := style.TierForHeight(d.reqs.MinHeight)
	secondaryTier := style.SecondaryTier(primaryTier)

	primaryEntry := ctx.Entry(primaryTier)
	secondaryEntry := ctx.Entry(secondaryTier)

	// Signal computation.
	barLevel := source.SignalToBarLevel(data.SignalStrength)
	qColor := source.QualityColor(barLevel, accent)

	// Row height definitions.
	ssidRowH := primaryEntry.RowHeight
	signalRowH := signalBarsHeight // 16px for signal bars sprite
	progressBarHeight := 4
	detailRowH := secondaryEntry.RowHeight

	// Row heights for the three fixed leading rows. Detail rows (all detailRowH)
	// are appended below once the policy-driven detail list is known, then
	// FitRows is run against the final array.
	// [0] SSID header, [1] signal zone, [2] progress bar
	rowHeights := []int{
		ssidRowH,          // 0: header (icon + SSID)
		signalRowH,        // 1: signal zone (bars + label)
		progressBarHeight, // 2: link quality progress bar
	}

	var spacing, offsetY, visibleCount int

	// --- Build SSID text (truncated to fit width minus icon) ---
	iconWidth := wifiIconDst // 16
	ssidGap := 2
	ssidAvailableWidth := bridge.AvailableContentWidth() - iconWidth - ssidGap
	ssidMaxChars := 0
	if primaryEntry.GlyphAdvance > 0 && ssidAvailableWidth > 0 {
		ssidMaxChars = ssidAvailableWidth / primaryEntry.GlyphAdvance
	}
	ssid := data.SSID
	if ssid == "" {
		ssid = "(no SSID)"
	}
	if ssidMaxChars > 0 {
		ssid = textlayout.Truncate(ssid, ssidMaxChars)
	}

	// Signal label text. In "bars" mode the signal-bars sprite already conveys
	// the level, so the redundant "Signal: N/4" string is dropped and the sprite
	// stands alone. For percentage/dbm modes the sprite has no numeric readout,
	// so the label is kept to give an exact value.
	signalLabel := ""
	if pol.SignalDisplay != "" && pol.SignalDisplay != "bars" {
		signalLabel = formatSignalLabel(data, pol, barLevel)
	}

	// Colors palette for the detail rows.
	white := color.RGBA{255, 255, 255, 255}
	dimmedWhite := sharedcolor.Dim(white)

	// --- Detail rows (policy-aware, consolidated) ---
	// Frequency and channel collapse into a single row when both are enabled,
	// so a height-constrained panel spends one row on radio info instead of two.
	// The show_frequency / show_channel / show_interface toggles are honored here
	// (the shared fast engine previously ignored them). IP and link speed are
	// always shown when present, matching the bespoke panel layouts.
	detailItems := make([]string, 0, 5)

	if radio := combinedRadioLine(data, pol); radio != "" {
		detailItems = append(detailItems, radio)
	}
	if data.LinkSpeed > 0 {
		detailItems = append(detailItems, fmt.Sprintf("%d Mbps", data.LinkSpeed))
	}
	if data.IPAddress != "" {
		detailItems = append(detailItems, data.IPAddress)
	}
	if pol.ShowInterface && data.InterfaceName != "" {
		detailItems = append(detailItems, data.InterfaceName)
	}

	// Assemble the full row list: fixed header/signal/progress rows followed by
	// the dynamic detail rows. Row heights, tiers, and colors are kept in lockstep.
	allItems := []string{
		ssid,        // 0: header
		signalLabel, // 1: signal label ("" in bars mode → sprite only)
		"",          // 2: progress bar (no text)
	}
	allTiers := []tiercatalog.Tier{primaryTier, secondaryTier, secondaryTier}
	var allColors []color.Color
	if d.p.Color {
		allColors = []color.Color{accent, qColor, nil}
	}

	for _, item := range detailItems {
		allItems = append(allItems, item)
		allTiers = append(allTiers, secondaryTier)
		rowHeights = append(rowHeights, detailRowH)
		if d.p.Color {
			allColors = append(allColors, dimmedWhite)
		}
	}

	// Trim the row-height scratch array to the rows we actually declared, then
	// recompute the fit so visibleCount reflects the real (consolidated) layout.
	rowHeights = rowHeights[:len(allItems)]
	spacing, offsetY, visibleCount = bridge.FitRows(rowHeights)

	// Trim to visible count.
	if visibleCount > len(allItems) {
		visibleCount = len(allItems)
	}
	if visibleCount < 1 {
		visibleCount = 1
	}
	items := allItems[:visibleCount]
	tiers := allTiers[:visibleCount]
	var colors []color.Color
	if allColors != nil {
		colors = allColors[:visibleCount]
	}

	// Truncate detail rows to fit width.
	for i := 3; i < len(items); i++ {
		if secondaryEntry.GlyphAdvance > 0 {
			maxChars := bridge.AvailableContentWidth() / secondaryEntry.GlyphAdvance
			if maxChars > 0 {
				items[i] = textlayout.Truncate(items[i], maxChars)
			}
		}
	}

	// --- Compute horizontal offsets ---
	offsets := make([]int, len(items))
	ox, _ := bridge.ContentOrigin()
	for i := range items {
		if i == 0 {
			// SSID is offset by icon + gap.
			offsets[i] = ox + iconWidth + ssidGap
		} else {
			offsets[i] = ox
		}
	}

	// --- Build sprites ---
	var sprites []widgets.Sprite

	// WiFi icon sprite (positioned in header zone).
	iconImg := renderWifiIcon(source.Connected, accent)
	if iconImg != nil {
		iconX, contentOriginY := bridge.ContentOrigin()
		iconY := contentOriginY + offsetY + (ssidRowH-wifiIconDst)/2
		if iconY < contentOriginY {
			iconY = contentOriginY
		}
		sprites = append(sprites, widgets.Sprite{
			Image:    iconImg,
			Position: image.Point{X: iconX, Y: iconY},
			Label:    "wifi-icon",
		})
	}

	// Signal bars sprite (positioned in signal zone).
	signalBarsImg := renderSignalBars(barLevel, qColor)
	if signalBarsImg != nil {
		_, contentOriginY := bridge.ContentOrigin()
		signalY := contentOriginY + offsetY + ssidRowH + spacing
		sprites = append(sprites, widgets.Sprite{
			Image:    signalBarsImg,
			Position: image.Point{X: ox, Y: signalY},
			Label:    "signal-bars",
		})
	}

	// Progress bar sprite (link quality) — only if row 2 is visible.
	if visibleCount > 2 {
		_, contentOriginY := bridge.ContentOrigin()
		progressY := contentOriginY + offsetY + ssidRowH + spacing + signalRowH + spacing
		progressWidth := bridge.AvailableContentWidth()
		progressValue := float64(data.LinkQuality) / 100.0

		progressBg := color.RGBA{30, 30, 30, 255}
		progressFg := qColor
		if !d.p.Color {
			progressFg = white
		}
		progressCfg := progressbar.Config{
			Style:      progressbar.Linear,
			Value:      progressValue,
			Bounds:     image.Rect(ox, progressY, ox+progressWidth, progressY+progressBarHeight),
			Foreground: progressFg,
			Background: progressBg,
		}
		pbSprite := progressbar.Render(progressCfg)
		if pbSprite != nil {
			sprites = append(sprites, widgets.Sprite{
				Image:    pbSprite.Image,
				Position: pbSprite.Position,
				Label:    "link-quality",
			})
		}
	}

	return style.ViewData{
		Items:        items,
		Colors:       colors,
		Tiers:        tiers,
		LineOffsets:  offsets,
		OffsetY:      offsetY,
		RowHeights:   rowHeights[:visibleCount],
		Spacing:      spacing,
		VisibleCount: visibleCount,
		Sprites:      sprites,
		Static:       false,
		PaddingPct:   paddingPct,
	}
}

// combinedRadioLine builds a single consolidated radio-info row from the
// frequency and channel, honoring the show_frequency / show_channel policy
// toggles. Collapsing both into one row ("Ch 36 (5.2GHz)") keeps the detail
// stack short on height-constrained panels. Returns "" when both are disabled
// or unavailable.
func combinedRadioLine(data source.WifiData, pol source.Policy) string {
	showFreq := pol.ShowFrequency && data.Frequency > 0
	showChan := pol.ShowChannel && data.Channel > 0

	switch {
	case showFreq && showChan:
		return fmt.Sprintf("Ch %d (%.1fGHz)", data.Channel, data.Frequency)
	case showChan:
		return fmt.Sprintf("Ch %d", data.Channel)
	case showFreq:
		return fmt.Sprintf("%.1f GHz", data.Frequency)
	default:
		return ""
	}
}
