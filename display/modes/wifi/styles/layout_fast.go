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

	// Build row heights array for FitRows:
	// [0] SSID header, [1] signal zone, [2] progress bar,
	// [3] frequency, [4] channel, [5] speed, [6] IP, [7] interface
	rowHeights := []int{
		ssidRowH,          // 0: header (icon + SSID)
		signalRowH,        // 1: signal zone (bars + label)
		progressBarHeight, // 2: link quality progress bar
		detailRowH,        // 3: frequency
		detailRowH,        // 4: channel
		detailRowH,        // 5: speed
		detailRowH,        // 6: IP
		detailRowH,        // 7: interface
	}

	spacing, offsetY, visibleCount := bridge.FitRows(rowHeights)

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

	// Signal label text.
	signalLabel := formatSignalLabel(data, pol, barLevel)

	// Detail row texts.
	freqText := fmt.Sprintf("%.2f GHz", data.Frequency)
	chanText := fmt.Sprintf("Ch %d", data.Channel)
	speedText := fmt.Sprintf("%d Mbps", data.LinkSpeed)
	ipText := data.IPAddress
	ifaceText := data.InterfaceName

	// All text items in order.
	allItems := []string{
		ssid,        // 0: header
		signalLabel, // 1: signal label
		"",          // 2: progress bar (no text)
		freqText,    // 3: frequency
		chanText,    // 4: channel
		speedText,   // 5: speed
		ipText,      // 6: IP
		ifaceText,   // 7: interface
	}

	// Per-row tier declarations.
	allTiers := []tiercatalog.Tier{
		primaryTier,   // 0: header
		secondaryTier, // 1: signal
		secondaryTier, // 2: progress bar
		secondaryTier, // 3: freq
		secondaryTier, // 4: channel
		secondaryTier, // 5: speed
		secondaryTier, // 6: IP
		secondaryTier, // 7: interface
	}

	// Colors: accent for SSID, quality-mapped for signal, dimmed for details.
	white := color.RGBA{255, 255, 255, 255}
	dimmedWhite := sharedcolor.Dim(white)

	var allColors []color.Color
	if d.p.Color {
		allColors = []color.Color{
			accent,      // 0: SSID in accent
			qColor,      // 1: signal in quality color
			nil,         // 2: progress bar (no text)
			dimmedWhite, // 3: freq
			dimmedWhite, // 4: channel
			dimmedWhite, // 5: speed
			dimmedWhite, // 6: IP
			dimmedWhite, // 7: interface
		}
	}

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
