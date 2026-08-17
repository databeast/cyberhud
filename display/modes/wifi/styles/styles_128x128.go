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

// WiFi style declarations for 128x128 panels.
//
// Each entry is a hand-tweakable declaration over the core dispatcher in
// core.go: adjust Params to select shared layouts or attach a bespoke BuildFn.

var ColorSmall128x128Style = def{
	name: "color-128x128",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 128, Capability: style.ColorFast},
	p:    Params{BuildFn: colorSmall128x128Build},
}

var MonoSlow128x128Style = def{
	name: "mono-slow-128x128",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 128, Capability: style.MonoSlow},
	p:    Params{MonoSlow: true},
}

var GrayscaleSlow128x128Style = def{
	name: "grayscale-slow-128x128",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 128, Capability: style.GrayscaleSlow},
}

var GrayscaleFast128x128Style = def{
	name: "grayscale-fast-128x128",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 128, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: grayscaleFast128x128Build},
}

var ColorSlow128x128Style = def{
	name: "color-slow-128x128",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 128, Capability: style.ColorSlow},
}

var MonoFast128x128Style = def{
	name: "mono-fast-128x128",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 128, Capability: style.MonoFast},
	p:    Params{Fast: true},
}

func colorSmall128x128Build(snapshot source.WifiData, pol source.Policy, ctx style.StyleContext, d def) style.ViewData {
	bridge := ctx.Layout(2)

	if bridge.AvailableContentWidth() == 0 || bridge.AvailableContentHeight() == 0 {
		return style.ViewData{Items: []string{"(too small)"}}
	}

	accent := resolveFGColor(pol.FGColor)
	white := color.RGBA{255, 255, 255, 255}
	dimmedWhite := sharedcolor.Dim(white)
	red := color.RGBA{255, 0, 0, 255}

	// Handle disconnected state: centered "No Network" in red + dimmed icon.
	if snapshot.ConnectionState == source.Disconnected {
		return colorSmall128x128Disconnected(bridge, ctx, accent, red, dimmedWhite)
	}

	// Handle unavailable state: centered "WiFi N/A" in dimmed white.
	if snapshot.ConnectionState == source.Unavailable {
		return colorSmall128x128Unavailable(bridge, ctx, dimmedWhite)
	}

	// Connected state: full multi-zone layout.
	return colorSmall128x128Connected(snapshot, pol, bridge, ctx, accent, white, dimmedWhite)

}

func colorSmall128x128Disconnected(bridge layout.LayoutCalculator, ctx style.StyleContext, accent, red, dimmedWhite color.RGBA) style.ViewData {
	msg := "No Network"

	// Use catalog metrics for layout (TierLarge for single prominent message).
	largeEntry := ctx.Entry(tiercatalog.TierLarge)
	glyphAdvance := largeEntry.GlyphAdvance
	rowH := largeEntry.RowHeight

	offsetX := bridge.CenterXWith(len(msg), glyphAdvance)

	// Vertically center.
	offsetY := bridge.CenterBlockY([]int{rowH}, 0)

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

	return style.ViewData{
		Items:       []string{msg},
		Colors:      []color.Color{red},
		Tiers:       []tiercatalog.Tier{tiercatalog.TierLarge},
		LineOffsets: []int{offsetX},
		OffsetY:     offsetY,
		Sprites:     sprites,
	}

}

func colorSmall128x128Unavailable(bridge layout.LayoutCalculator, ctx style.StyleContext, dimmedWhite color.RGBA) style.ViewData {
	msg := "WiFi N/A"

	// Use catalog metrics for layout (TierLarge for single prominent message).
	largeEntry := ctx.Entry(tiercatalog.TierLarge)
	glyphAdvance := largeEntry.GlyphAdvance
	rowH := largeEntry.RowHeight

	offsetX := bridge.CenterXWith(len(msg), glyphAdvance)
	offsetY := bridge.CenterBlockY([]int{rowH}, 0)

	// Render WiFi icon in unavailable (dimmed gray) tint per req 6.7.
	var sprites []widgets.Sprite
	iconImg := renderWifiIcon(source.Unavailable, dimmedWhite)
	if iconImg != nil {
		ox, oy := bridge.ContentOrigin()
		sprites = append(sprites, widgets.Sprite{
			Image:    iconImg,
			Position: image.Point{X: ox, Y: oy},
			Label:    "wifi-icon",
		})
	}
	if sprites == nil {
		sprites = []widgets.Sprite{}
	}

	return style.ViewData{
		Items:       []string{msg},
		Colors:      []color.Color{dimmedWhite},
		Tiers:       []tiercatalog.Tier{tiercatalog.TierLarge},
		LineOffsets: []int{offsetX},
		OffsetY:     offsetY,
		Sprites:     sprites,
	}

}

func colorSmall128x128Connected(snapshot source.WifiData, pol source.Policy, bridge layout.LayoutCalculator, ctx style.StyleContext, accent, white, dimmedWhite color.RGBA) style.ViewData {
	// --- Font metrics from tier catalog ---
	// Small tier: detail/footer text rows.
	smallEntry := ctx.Entry(tiercatalog.TierSmall)
	smallRowH := smallEntry.RowHeight

	// Large tier: SSID header row (slightly larger than small).
	largeEntry := ctx.Entry(tiercatalog.TierLarge)
	ssidRowH := largeEntry.RowHeight
	ssidAdvance := largeEntry.GlyphAdvance

	// --- Row height definitions for FitRows ---
	// header (SSID row height), signal (signal bars sprite height=16),
	// progress (4px), details row, details row, details row, footer row, footer row
	progressBarHeight := 4
	signalRowH := signalBarsHeight // 16px

	rowHeights := []int{
		ssidRowH,          // 0: header (icon + SSID)
		signalRowH,        // 1: signal zone (bars + label)
		progressBarHeight, // 2: progress bar
		smallRowH,         // 3: detail - frequency
		smallRowH,         // 4: detail - channel
		smallRowH,         // 5: detail - speed
		smallRowH,         // 6: footer - IP
		smallRowH,         // 7: footer - interface
	}

	spacing, offsetY, visibleCount := bridge.FitRows(rowHeights)

	// --- Compute signal level and colors ---
	barLevel := source.SignalToBarLevel(snapshot.SignalStrength)
	qColor := source.QualityColor(barLevel, accent)

	// --- Build text items and colors ---
	// Truncate SSID to fit available width accounting for icon (16px) + 2px gap.
	iconWidth := wifiIconDst // 16
	ssidGap := 2
	ssidAvailableWidth := bridge.AvailableContentWidth() - iconWidth - ssidGap
	ssidMaxChars := 0
	if ssidAdvance > 0 && ssidAvailableWidth > 0 {
		ssidMaxChars = ssidAvailableWidth / ssidAdvance
	}
	ssid := snapshot.SSID
	if ssidMaxChars > 0 {
		ssid = textlayout.Truncate(ssid, ssidMaxChars)
	}

	// Signal label text (depends on SignalDisplay policy).
	signalLabel := formatSignalLabel(snapshot, pol, barLevel)

	// Detail rows.
	freqText := fmt.Sprintf("Freq: %.1f GHz", snapshot.Frequency)
	chanText := fmt.Sprintf("Ch: %d", snapshot.Channel)
	speedText := fmt.Sprintf("Speed: %d Mbps", snapshot.LinkSpeed)
	ipText := snapshot.IPAddress
	ifaceText := snapshot.InterfaceName

	allItems := []string{
		ssid,
		signalLabel,
		"", // progress bar row (no text)
		freqText,
		chanText,
		speedText,
		ipText,
		ifaceText,
	}

	allColors := []color.Color{
		accent,      // SSID in accent
		qColor,      // signal label in quality color
		nil,         // progress bar (no text)
		dimmedWhite, // freq value
		dimmedWhite, // channel value
		dimmedWhite, // speed value
		dimmedWhite, // IP
		dimmedWhite, // interface
	}

	// Per-row tier declarations: SSID row uses TierLarge, rest use TierSmall.
	allTiers := []tiercatalog.Tier{
		tiercatalog.TierLarge, // header
		tiercatalog.TierSmall, // signal
		tiercatalog.TierSmall, // progress bar
		tiercatalog.TierSmall, // freq
		tiercatalog.TierSmall, // channel
		tiercatalog.TierSmall, // speed
		tiercatalog.TierSmall, // IP
		tiercatalog.TierSmall, // interface
	}

	// Trim to visible count.
	items := allItems[:visibleCount]
	colors := allColors[:visibleCount]
	tiers := allTiers[:visibleCount]

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

	// WiFi icon sprite (positioned in header zone, vertically centered to SSID baseline).
	iconImg := renderWifiIcon(source.Connected, accent)
	if iconImg != nil {
		iconX, _ := bridge.ContentOrigin()
		// Vertically align icon center with SSID text center at row 0.
		_, contentOriginY := bridge.ContentOrigin()
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

	// Progress bar sprite (positioned below signal zone).
	if visibleCount > 2 {
		_, contentOriginY := bridge.ContentOrigin()
		progressY := contentOriginY + offsetY + ssidRowH + spacing + signalRowH + spacing
		progressWidth := bridge.AvailableContentWidth()
		progressValue := float64(snapshot.LinkQuality) / 100.0

		progressBg := color.RGBA{30, 30, 30, 255}
		progressCfg := progressbar.Config{
			Style:      progressbar.Linear,
			Value:      progressValue,
			Bounds:     image.Rect(ox, progressY, ox+progressWidth, progressY+progressBarHeight),
			Foreground: qColor,
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
		Items:       items,
		Colors:      colors,
		Tiers:       tiers,
		LineOffsets: offsets,
		OffsetY:     offsetY,
		Sprites:     sprites,
	}
}

// formatSignalLabel returns the signal quality label text based on display policy.
func formatSignalLabel(snapshot source.WifiData, pol source.Policy, barLevel int) string {
	switch pol.SignalDisplay {
	case "dbm":
		return source.FormatDbm(snapshot.SignalStrength)
	case "percentage", "percent":
		return fmt.Sprintf("%d%%", source.SignalPercent(snapshot.SignalStrength))
	default: // "bars"
		return fmt.Sprintf("Signal: %d/4", barLevel)
	}
}

func grayscaleFast128x128Build(snap source.WifiData, pol source.Policy, ctx style.StyleContext, d def) style.ViewData {
	bridge := ctx.Layout(2)

	if bridge.AvailableContentWidth() == 0 || bridge.AvailableContentHeight() == 0 {
		return style.ViewData{Items: []string{"(too small)"}}
	}

	white := color.RGBA{255, 255, 255, 255}
	lightGray := color.RGBA{180, 180, 180, 255}

	// Handle disconnected/unavailable states with centered message.
	if snap.ConnectionState == source.Disconnected {
		msg := snap.StatusMessage
		if msg == "" {
			msg = "No Network"
		}
		return style.ViewData{
			Items:  []string{msg},
			Colors: []color.Color{white},
			Tiers:  []tiercatalog.Tier{tiercatalog.TierLarge},
		}
	}
	if snap.ConnectionState == source.Unavailable {
		msg := snap.StatusMessage
		if msg == "" {
			msg = "WiFi N/A"
		}
		return style.ViewData{
			Items:  []string{msg},
			Colors: []color.Color{lightGray},
			Tiers:  []tiercatalog.Tier{tiercatalog.TierLarge},
		}
	}

	// Connected state: build text-only rows.
	items := []string{}
	colors := []color.Color{}
	tiers := []tiercatalog.Tier{}

	// Row 0: SSID (truncated to fit) — use large tier for header.
	ssid := snap.SSID
	if ssid == "" {
		ssid = "(no SSID)"
	}
	maxChars := bridge.AvailableContentWidth() / bridge.GlyphAdvance()
	if maxChars > 0 && len(ssid) > maxChars {
		ssid = ssid[:maxChars]
	}
	items = append(items, ssid)
	colors = append(colors, white)
	tiers = append(tiers, tiercatalog.TierLarge)

	// Row 1: Signal level as text (uses shared formatSignalText).
	items = append(items, source.FormatSignalText(snap, pol))
	colors = append(colors, lightGray)
	tiers = append(tiers, tiercatalog.TierSmall)

	// Row 2: Link quality.
	items = append(items, fmt.Sprintf("Link: %d%%", snap.LinkQuality))
	colors = append(colors, lightGray)
	tiers = append(tiers, tiercatalog.TierSmall)

	// Row 3: Frequency + Channel (if enabled).
	if pol.ShowFrequency {
		freqCh := fmt.Sprintf("%.1fGHz", snap.Frequency)
		if pol.ShowChannel && snap.Channel > 0 {
			freqCh += fmt.Sprintf(" Ch%d", snap.Channel)
		}
		items = append(items, freqCh)
		colors = append(colors, lightGray)
		tiers = append(tiers, tiercatalog.TierSmall)
	} else if pol.ShowChannel && snap.Channel > 0 {
		items = append(items, fmt.Sprintf("Ch %d", snap.Channel))
		colors = append(colors, lightGray)
		tiers = append(tiers, tiercatalog.TierSmall)
	}

	// Row 4: IP address.
	if snap.IPAddress != "" {
		items = append(items, snap.IPAddress)
		colors = append(colors, lightGray)
		tiers = append(tiers, tiercatalog.TierSmall)
	}

	// Row 5: Interface name (if enabled).
	if pol.ShowInterface && snap.InterfaceName != "" {
		items = append(items, snap.InterfaceName)
		colors = append(colors, lightGray)
		tiers = append(tiers, tiercatalog.TierSmall)
	}

	// Use FitRows to determine spacing and vertical centering.
	rowH := bridge.RowHeight()
	rowHeights := make([]int, len(items))
	for i := range rowHeights {
		rowHeights[i] = rowH
	}
	_, offsetY, _ := bridge.FitRows(rowHeights)

	return style.ViewData{
		Items:   items,
		Colors:  colors,
		Tiers:   tiers,
		OffsetY: offsetY,
	}

}
