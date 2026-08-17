package styles

import (
	"fmt"
	"image"
	"image/color"

	"github.com/databeast/cyberhud/display/modes/wifi/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
	"github.com/databeast/cyberhud/display/widgets"
	"github.com/databeast/cyberhud/display/widgets/progressbar"
)

// WiFi style declarations for 240x240 panels.
//
// Each entry is a hand-tweakable declaration over the core dispatcher in
// core.go: adjust Params to select shared layouts or attach a bespoke BuildFn.

var ColorSlow240x240Style = def{
	name: "color-slow-240x240",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 240, Capability: style.ColorSlow},
}

var GrayscaleSlow240x240Style = def{
	name: "grayscale-slow-240x240",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 240, Capability: style.GrayscaleSlow},
}

var MonoSlow240x240Style = def{
	name: "mono-slow-240x240",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 240, Capability: style.MonoSlow},
	p:    Params{MonoSlow: true},
}

var MonoFast240x240Style = def{
	name: "mono-fast-240x240",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 240, Capability: style.MonoFast},
	p:    Params{Fast: true},
}

var GrayscaleFast240x240Style = def{
	name: "grayscale-fast-240x240",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 240, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: grayscaleFast240x240Build},
}

var Color240x240Style = def{
	name: "color-240x240",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 240, Capability: style.ColorFast},
	p:    Params{BuildFn: color240x240Build},
}

func grayscaleFast240x240Build(snap source.WifiData, pol source.Policy, ctx style.StyleContext, d def) style.ViewData {
	bridge := ctx.Layout(2)

	if bridge.AvailableContentWidth() == 0 || bridge.AvailableContentHeight() == 0 {
		return style.ViewData{Items: []string{"(too small)"}}
	}

	white := color.RGBA{255, 255, 255, 255}
	dimmed := color.RGBA{180, 180, 180, 255}

	// Disconnected: centered status message.
	if snap.ConnectionState == source.Disconnected {
		msg := snap.StatusMessage
		if msg == "" {
			msg = "No Network"
		}
		return style.ViewData{
			Items:  []string{msg},
			Colors: []color.Color{color.RGBA{255, 80, 80, 255}},
			Tiers:  []tiercatalog.Tier{tiercatalog.TierLarge},
		}
	}

	// Unavailable: centered status message.
	if snap.ConnectionState == source.Unavailable {
		msg := snap.StatusMessage
		if msg == "" {
			msg = "WiFi N/A"
		}
		return style.ViewData{
			Items:  []string{msg},
			Colors: []color.Color{dimmed},
			Tiers:  []tiercatalog.Tier{tiercatalog.TierLarge},
		}
	}

	// Connected: text-only layout with signal as text.
	items := []string{snap.SSID}
	colors := []color.Color{white}
	tiers := []tiercatalog.Tier{tiercatalog.TierLarge}

	// Signal line: display as text based on policy signal_display preference.
	sigLine := source.FormatSignalText(snap, pol)
	items = append(items, sigLine)
	colors = append(colors, dimmed)
	tiers = append(tiers, tiercatalog.TierSmall)

	// IP address.
	if snap.IPAddress != "" {
		items = append(items, snap.IPAddress)
		colors = append(colors, dimmed)
		tiers = append(tiers, tiercatalog.TierSmall)
	}

	// Channel + Frequency (if policy allows).
	if pol.ShowChannel && snap.Channel > 0 {
		chLine := fmt.Sprintf("Ch %d", snap.Channel)
		if pol.ShowFrequency && snap.Frequency > 0 {
			chLine = fmt.Sprintf("Ch %d (%.1fGHz)", snap.Channel, snap.Frequency)
		}
		items = append(items, chLine)
		colors = append(colors, dimmed)
		tiers = append(tiers, tiercatalog.TierSmall)
	} else if pol.ShowFrequency && snap.Frequency > 0 {
		items = append(items, fmt.Sprintf("%.1f GHz", snap.Frequency))
		colors = append(colors, dimmed)
		tiers = append(tiers, tiercatalog.TierSmall)
	}

	// Interface name (if policy allows).
	if pol.ShowInterface && snap.InterfaceName != "" {
		items = append(items, snap.InterfaceName)
		colors = append(colors, dimmed)
		tiers = append(tiers, tiercatalog.TierSmall)
	}

	// Link speed (if available).
	if snap.LinkSpeed > 0 {
		items = append(items, fmt.Sprintf("%d Mbps", snap.LinkSpeed))
		colors = append(colors, dimmed)
		tiers = append(tiers, tiercatalog.TierSmall)
	}

	// Use FitRows for vertical distribution.
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

func color240x240Build(snap source.WifiData, pol source.Policy, ctx style.StyleContext, d def) style.ViewData {
	bridge := ctx.Layout(2)

	if bridge.AvailableContentWidth() == 0 || bridge.AvailableContentHeight() == 0 {
		return style.ViewData{Items: []string{"(too small)"}}
	}

	accent := resolveFGColor(pol.FGColor)
	white := color.RGBA{255, 255, 255, 255}
	dimmedWhite := color.RGBA{R: 180, G: 180, B: 180, A: 255}
	red := color.RGBA{255, 0, 0, 255}

	// --- Disconnected state: centered message + dimmed icon ---
	if snap.ConnectionState == source.Disconnected {
		msg := snap.StatusMessage
		if msg == "" {
			msg = "No Network"
		}
		offsetX := bridge.CenterX(len(msg))

		vd := style.ViewData{
			Items:       []string{msg},
			Colors:      []color.Color{red},
			Tiers:       []tiercatalog.Tier{tiercatalog.TierLarge},
			LineOffsets: []int{offsetX},
		}

		// Add WiFi icon in disconnected state.
		iconImg := renderWifiIcon(source.Disconnected, accent)
		if iconImg != nil {
			ox, oy := bridge.ContentOrigin()
			if wifiIconDst <= bridge.AvailableContentWidth() && wifiIconDst <= bridge.AvailableContentHeight() {
				vd.Sprites = append(vd.Sprites, widgets.Sprite{
					Image:    iconImg,
					Position: image.Point{X: ox, Y: oy},
					Label:    "wifi-icon",
				})
			}
		}
		return vd
	}

	// --- Unavailable state: centered message, no sprites ---
	if snap.ConnectionState == source.Unavailable {
		msg := snap.StatusMessage
		if msg == "" {
			msg = "WiFi N/A"
		}
		offsetX := bridge.CenterX(len(msg))
		return style.ViewData{
			Items:       []string{msg},
			Colors:      []color.Color{dimmedWhite},
			Tiers:       []tiercatalog.Tier{tiercatalog.TierLarge},
			LineOffsets: []int{offsetX},
		}
	}

	// --- Connected state: full layout ---

	// Get font metrics from tier catalog: large for SSID header, small for details.
	smallEntry := ctx.Entry(tiercatalog.TierSmall)
	largeEntry := ctx.Entry(tiercatalog.TierLarge)

	smallRowH := smallEntry.RowHeight
	smallAdvance := smallEntry.GlyphAdvance

	largeRowH := largeEntry.RowHeight
	largeAdvance := largeEntry.GlyphAdvance

	// Build content rows and compute row heights.
	type row struct {
		text  string
		color color.Color
		tier  tiercatalog.Tier
		rowH  int
		adv   int
	}

	var rows []row

	// Row 0 — Header: SSID (large font, accent color).
	ssid := snap.SSID
	if ssid == "" {
		ssid = "(no SSID)"
	}
	// Truncate SSID to fit: available width minus icon(16) minus gap(2).
	ssidMaxW := bridge.AvailableContentWidth() - wifiIconDst - 2
	if largeAdvance > 0 && ssidMaxW > 0 {
		maxChars := ssidMaxW / largeAdvance
		if maxChars > 0 {
			ssid = textlayout.Truncate(ssid, maxChars)
		}
	}
	rows = append(rows, row{text: ssid, color: accent, tier: tiercatalog.TierLarge, rowH: largeRowH, adv: largeAdvance})

	// Row 1 — Signal label (small font, quality-mapped color).
	barLevel := source.SignalToBarLevel(snap.SignalStrength)
	qColor := source.QualityColor(barLevel, accent)
	signalLabel := signalLabelForPolicy(pol.SignalDisplay, snap.SignalStrength, barLevel)
	rows = append(rows, row{text: signalLabel, color: qColor, tier: tiercatalog.TierSmall, rowH: smallRowH, adv: smallAdvance})

	// Row 2 — Link quality line placeholder (progress bar rendered as sprite).
	// Use a small-height row to represent the progress bar zone.
	const progressBarH = 4
	rows = append(rows, row{text: "", color: nil, tier: tiercatalog.TierSmall, rowH: progressBarH + 4, adv: smallAdvance}) // 4px bar + 4px spacing

	// Row 3+ — Detail rows.
	if pol.ShowFrequency && snap.Frequency > 0 {
		freqText := fmt.Sprintf("Freq: %.1f GHz", snap.Frequency)
		if pol.ShowChannel && snap.Channel > 0 {
			freqText = fmt.Sprintf("%.1fGHz Ch%d", snap.Frequency, snap.Channel)
		}
		rows = append(rows, row{text: freqText, color: dimmedWhite, tier: tiercatalog.TierSmall, rowH: smallRowH, adv: smallAdvance})
	} else if pol.ShowChannel && snap.Channel > 0 {
		rows = append(rows, row{text: fmt.Sprintf("Ch %d", snap.Channel), color: dimmedWhite, tier: tiercatalog.TierSmall, rowH: smallRowH, adv: smallAdvance})
	}

	if snap.LinkSpeed > 0 {
		rows = append(rows, row{text: fmt.Sprintf("Speed: %dMbps", snap.LinkSpeed), color: dimmedWhite, tier: tiercatalog.TierSmall, rowH: smallRowH, adv: smallAdvance})
	}

	// Footer rows.
	if snap.IPAddress != "" {
		rows = append(rows, row{text: snap.IPAddress, color: white, tier: tiercatalog.TierSmall, rowH: smallRowH, adv: smallAdvance})
	}
	if pol.ShowInterface && snap.InterfaceName != "" {
		rows = append(rows, row{text: snap.InterfaceName, color: dimmedWhite, tier: tiercatalog.TierSmall, rowH: smallRowH, adv: smallAdvance})
	}

	// Compute FitRows layout.
	rowHeights := make([]int, len(rows))
	for i, r := range rows {
		rowHeights[i] = r.rowH
	}
	spacing, offsetY, visibleCount := bridge.FitRows(rowHeights)

	// Build ViewData items from visible rows.
	items := make([]string, visibleCount)
	colors := make([]color.Color, visibleCount)
	offsets := make([]int, visibleCount)
	tiers := make([]tiercatalog.Tier, visibleCount)

	for i := 0; i < visibleCount; i++ {
		items[i] = rows[i].text
		colors[i] = rows[i].color
		tiers[i] = rows[i].tier
		// Center each text row horizontally.
		if rows[i].adv > 0 {
			offsets[i] = bridge.CenterXWith(len(rows[i].text), rows[i].adv)
		}
	}

	vd := style.ViewData{
		Items:       items,
		Colors:      colors,
		Tiers:       tiers,
		LineOffsets: offsets,
		OffsetY:     offsetY + spacing,
		Static:      false,
	}

	// --- Sprites ---
	ox, oy := bridge.ContentOrigin()

	// WiFi icon sprite (16×16) at top-left of content area.
	iconImg := renderWifiIcon(source.Connected, accent)
	if iconImg != nil {
		if wifiIconDst <= bridge.AvailableContentWidth() && wifiIconDst <= bridge.AvailableContentHeight() {
			vd.Sprites = append(vd.Sprites, widgets.Sprite{
				Image:    iconImg,
				Position: image.Point{X: ox, Y: oy + offsetY},
				Label:    "wifi-icon",
			})
		}
	}

	// Signal bar sprite positioned in the signal zone row.
	signalBarImg := renderSignalBars(barLevel, qColor)
	if signalBarImg != nil {
		// Position signal bars after SSID row (row 0 height + spacing).
		signalY := oy + offsetY + rows[0].rowH + spacing
		vd.Sprites = append(vd.Sprites, widgets.Sprite{
			Image:    signalBarImg,
			Position: image.Point{X: ox, Y: signalY},
			Label:    "signal-bars",
		})
	}

	// Progress bar sprite below signal zone.
	progressY := oy + offsetY + rows[0].rowH + spacing + rows[1].rowH + spacing
	barWidth := bridge.AvailableContentWidth()
	if barWidth > 0 {
		// Determine progress bar foreground color.
		var pbFg color.RGBA
		if snap.ConnectionState == source.Connected {
			pbFg = qColor
		} else {
			pbFg = color.RGBA{30, 30, 30, 255}
		}
		pbBg := color.RGBA{30, 30, 30, 255}
		linkValue := float64(snap.LinkQuality) / 100.0
		if linkValue < 0 {
			linkValue = 0
		}
		if linkValue > 1 {
			linkValue = 1
		}

		barBounds := image.Rect(ox, progressY, ox+barWidth, progressY+progressBarH)
		pbSprite := progressbar.Render(progressbar.Config{
			Style:      progressbar.Linear,
			Value:      linkValue,
			Bounds:     barBounds,
			Foreground: pbFg,
			Background: pbBg,
		})
		if pbSprite != nil {
			vd.Sprites = append(vd.Sprites, *pbSprite)
		}
	}

	return vd

}

// signalLabelForPolicy returns the signal quality label text based on the
// policy's SignalDisplay preference.
func signalLabelForPolicy(displayMode string, dBm int, barLevel int) string {
	switch displayMode {
	case "dbm":
		return source.FormatDbm(dBm)
	case "percentage", "percent":
		return fmt.Sprintf("%d%%", source.SignalPercent(dBm))
	default: // "bars"
		return fmt.Sprintf("%d/4", barLevel)
	}
}
