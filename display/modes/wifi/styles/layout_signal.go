package styles

import (
	"fmt"
	"image"
	"image/color"

	"github.com/databeast/cyberhud/display/modes/wifi/source"
	"github.com/databeast/cyberhud/display/style"
	sharedcolor "github.com/databeast/cyberhud/display/style/color"
	"github.com/databeast/cyberhud/display/style/layout"
	font "github.com/databeast/cyberhud/display/surface/fonts"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
	"github.com/databeast/cyberhud/display/widgets"
	"github.com/databeast/cyberhud/display/widgets/progressbar"
	"github.com/databeast/cyberhud/display/widgets/textbox"
)

// placeLabel renders text with fittedLabel and positions the resulting (text-
// sized) sprite centered within the band at (bandX, bandY, bandW, bandH). It
// appends to and returns the sprites slice. A nil render (empty text, no face)
// is a no-op.
func placeLabel(sprites []widgets.Sprite, ctx style.StyleContext, text string,
	bandX, bandY, bandW, bandH int, fg color.RGBA, label string) []widgets.Sprite {
	sp := fittedLabel(ctx, text, bandW, bandH, fg)
	if sp == nil {
		return sprites
	}
	b := sp.Image.Bounds()
	sp.Position = image.Point{
		X: bandX + (bandW-b.Dx())/2,
		Y: bandY + (bandH-b.Dy())/2,
	}
	sp.Label = label
	return append(sprites, *sp)
}

// buildSignalForwardWifi renders the ambient, glanceable WiFi layout.
//
// The design is signal-first: a large arc gauge, filled by the current signal
// strength and colored by quality (green strong / amber moderate / red weak),
// is the hero of the panel. A raw dBm number is meaningless to most users, so
// the GAUGE carries the meaning — mostly-full-and-green reads as "good" with no
// technical knowledge. The SSID sits below as a quiet label with a state-tinted
// WiFi icon, and the IP address is pinned to the bottom as a persistent footer
// (the one operationally critical field on a headless box). Channel, frequency,
// interface and link speed are tinkerer detail, shown only when their policy
// toggles are enabled.
//
// Everything is composited as absolute-positioned sprites; this layout emits no
// renderer text rows (ViewData.Items stays empty). That deliberately sidesteps
// the two-plane (text-rows vs sprites) alignment problem that previously caused
// overlapping lines, by keeping the whole frame on a single sprite plane.
func buildSignalForwardWifi(data source.WifiData, pol source.Policy, ctx style.StyleContext, d def) style.ViewData {
	const paddingPct = 4
	bridge := ctx.Layout(paddingPct)

	cw := bridge.AvailableContentWidth()
	chh := bridge.AvailableContentHeight()
	if cw <= 0 || chh <= 0 {
		return style.ViewData{Items: []string{"(too small)"}, Static: false, PaddingPct: paddingPct}
	}

	accent := color.RGBA{255, 255, 255, 255} // white default for non-color panels
	if d.p.Color {
		accent = resolveFGColor(pol.FGColor)
	}

	// Panels too small for a readable arc gauge degrade to a compact signal-first
	// layout (horizontal signal bar + SSID + IP). Below ~56px in either axis the
	// 270° arc and its centered readout cannot render legibly.
	const minGaugeDim = 56
	if cw < minGaugeDim || chh < minGaugeDim {
		return signalForwardCompact(data, pol, bridge, ctx, d, accent, paddingPct)
	}

	switch data.ConnectionState {
	case source.Disconnected:
		return signalForwardStatus(bridge, ctx, d, "No Network",
			color.RGBA{255, 60, 60, 255}, source.Disconnected, accent, paddingPct)
	case source.Unavailable:
		return signalForwardStatus(bridge, ctx, d, "WiFi N/A",
			color.RGBA{120, 120, 120, 255}, source.Unavailable, accent, paddingPct)
	default:
		return signalForwardConnected(data, pol, bridge, ctx, d, accent, paddingPct)
	}
}

// signalForwardStatus renders the bold disconnected / unavailable screens:
// a large state-tinted WiFi icon over a single prominent message, centered.
func signalForwardStatus(bridge layout.LayoutCalculator, ctx style.StyleContext, d def,
	msg string, msgColor color.RGBA, state source.ConnectionState, accent color.RGBA, paddingPct int) style.ViewData {

	cw := bridge.AvailableContentWidth()
	chh := bridge.AvailableContentHeight()
	ox, oy := bridge.ContentOrigin()

	var sprites []widgets.Sprite

	// Large WiFi icon: scale the 8x8 glyph up to roughly a third of the panel.
	iconSize := minInt(cw, chh) / 3
	if iconSize < wifiIconDst {
		iconSize = wifiIconDst
	}
	tint := accent
	if state == source.Disconnected {
		tint = color.RGBA{255, 60, 60, 255}
	} else if state == source.Unavailable {
		tint = color.RGBA{120, 120, 120, 255}
	}
	_ = state
	if icon := scaledWifiIcon(tint, iconSize); icon != nil {
		iconX := ox + (cw-iconSize)/2
		iconY := oy + (chh/2 - iconSize)
		if iconY < oy {
			iconY = oy
		}
		sprites = append(sprites, widgets.Sprite{Image: icon, Position: image.Point{X: iconX, Y: iconY}, Label: "wifi-icon"})
	}

	// Message: a wide, bold scaled-text band centered under the icon.
	textColor := msgColor
	if !d.p.Color {
		textColor = color.RGBA{255, 255, 255, 255}
	}
	msgW := cw
	msgH := clampInt(chh/5, 12, 48)
	msgY := oy + chh/2 + 4
	sprites = placeLabel(sprites, ctx, msg, ox, msgY, msgW, msgH, textColor, "wifi-status")

	return style.ViewData{
		Sprites:    sprites,
		Static:     false,
		Cursor:     -1,
		PaddingPct: paddingPct,
	}
}

// signalForwardConnected renders the connected hero layout: arc gauge + SSID +
// IP footer + optional detail line.
func signalForwardConnected(data source.WifiData, pol source.Policy, bridge layout.LayoutCalculator,
	ctx style.StyleContext, d def, accent color.RGBA, paddingPct int) style.ViewData {

	cw := bridge.AvailableContentWidth()
	chh := bridge.AvailableContentHeight()
	ox, oy := bridge.ContentOrigin()

	barLevel := source.SignalToBarLevel(data.SignalStrength)
	qColor := source.QualityColor(barLevel, accent)
	if !d.p.Color {
		qColor = color.RGBA{255, 255, 255, 255}
	}
	fill := float64(source.SignalPercent(data.SignalStrength)) / 100.0

	// --- Fill the panel ---
	// The gauge is the hero and must OWN the screen: size it to the smaller panel
	// dimension so it grows with the hardware instead of floating in a fixed
	// square that leaves dead black margins on larger/wider panels. Text sizes are
	// all derived as fractions of this gauge size (never fixed pixel caps), so on a
	// physically small-but-dense panel the labels scale up to stay readable.
	gaugeSize := minInt(cw, chh)

	// All text metrics scale off the gauge/panel size.
	readoutH := clampInt(gaugeSize/5, 10, 200)   // the big quality word / figure
	ssidH := clampInt(gaugeSize/9, 9, 120)       // SSID label
	footerH := clampInt(gaugeSize/11, 8, 100)    // IP / detail line
	detailLine := signalForwardDetailLine(data, pol)

	var sprites []widgets.Sprite

	// --- Hero: arc gauge, centered on the panel ---
	gaugeX := ox + (cw-gaugeSize)/2
	gaugeY := oy + (chh-gaugeSize)/2
	thickness := clampInt(gaugeSize/10, 4, 60)
	gauge := progressbar.Render(progressbar.Config{
		Style:       progressbar.Arc,
		Orientation: progressbar.OrientHorizontal,
		Value:       fill,
		Bounds:      image.Rect(gaugeX, gaugeY, gaugeX+gaugeSize, gaugeY+gaugeSize),
		Foreground:  qColor,
		Background:  color.RGBA{34, 34, 34, 255},
		Thickness:   thickness,
		SweepAngle:  270,
		RoundedCaps: true,
	})
	if gauge != nil {
		sprites = append(sprites, *gauge)
	}

	gaugeCX := gaugeX + gaugeSize/2
	gaugeCY := gaugeY + gaugeSize/2

	// --- Readout: the big quality word / figure, dead center of the gauge ---
	// This is the primary thing the user reads, so it is large (a fifth of the
	// gauge), not a tiny glyph lost in the ring.
	readout := signalForwardReadout(data, pol, barLevel)
	readoutW := int(float64(gaugeSize) * 0.7)
	sprites = placeLabel(sprites, ctx, readout, gaugeCX-readoutW/2, gaugeCY-readoutH/2,
		readoutW, readoutH, qColor, "signal-readout")

	// --- SSID + IP stacked in the lower interior of the gauge ---
	// The 270° arc opens at the bottom, leaving a clear wedge below the readout.
	// SSID and IP live there, inside the gauge footprint, so the whole composition
	// stays as one filled unit instead of scattering labels into dead margins.
	ssid := data.SSID
	if ssid == "" {
		ssid = "(no SSID)"
	}
	ssidColor := accent
	dimmed := secondaryColor(d.p.Color)

	// Stack: SSID (with icon), optional detail line, IP — below the readout.
	stackY := gaugeCY + readoutH/2 + clampInt(gaugeSize/20, 4, 40)

	// SSID band: icon + name as one centered unit. Render the name first (text-
	// sized), then place the icon immediately to its left so the pair is centered.
	ssidIconSize := ssidH
	textW := int(float64(gaugeSize) * 0.6)
	if nameSp := fittedLabel(ctx, ssid, textW, ssidH, ssidColor); nameSp != nil {
		nameW := nameSp.Image.Bounds().Dx()
		icon := scaledWifiIcon(accent, ssidIconSize)
		gap := ssidIconSize / 4
		blockW := nameW
		if icon != nil {
			blockW += ssidIconSize + gap
		}
		blockX := gaugeCX - blockW/2
		x := blockX
		if icon != nil {
			iconY := stackY + (ssidH-ssidIconSize)/2
			sprites = append(sprites, widgets.Sprite{
				Image:    icon,
				Position: image.Point{X: x, Y: iconY},
				Label:    "wifi-icon",
			})
			x += ssidIconSize + gap
		}
		nameSp.Position = image.Point{X: x, Y: stackY + (ssidH-nameSp.Image.Bounds().Dy())/2}
		nameSp.Label = "ssid"
		sprites = append(sprites, *nameSp)
	}
	stackY += ssidH + clampInt(gaugeSize/40, 2, 20)

	// Optional detail line (channel/freq/iface/speed).
	if detailLine != "" {
		sprites = placeLabel(sprites, ctx, detailLine, gaugeCX-gaugeSize/2, stackY,
			gaugeSize, footerH, dimmed, "wifi-detail")
		stackY += footerH + clampInt(gaugeSize/40, 2, 20)
	}

	// IP footer: always shown.
	if data.IPAddress != "" {
		sprites = placeLabel(sprites, ctx, data.IPAddress, gaugeCX-gaugeSize/2, stackY,
			gaugeSize, footerH, dimmed, "wifi-ip")
	}

	return style.ViewData{
		Sprites:    sprites,
		Static:     false,
		Cursor:     -1,
		PaddingPct: paddingPct,
	}
}

// signalForwardCompact is the small-panel degradation of the signal-forward
// layout: a horizontal signal bar hero (the gauge can't render below ~56px),
// with SSID and IP stacked compactly. Still signal-first, still all sprites.
func signalForwardCompact(data source.WifiData, pol source.Policy, bridge layout.LayoutCalculator,
	ctx style.StyleContext, d def, accent color.RGBA, paddingPct int) style.ViewData {

	cw := bridge.AvailableContentWidth()
	chh := bridge.AvailableContentHeight()
	ox, oy := bridge.ContentOrigin()

	var sprites []widgets.Sprite

	// Non-connected: a single centered message fills the tiny panel.
	if data.ConnectionState != source.Connected {
		msg := "No Network"
		msgColor := color.RGBA{255, 60, 60, 255}
		if data.ConnectionState == source.Unavailable {
			msg = "WiFi N/A"
			msgColor = color.RGBA{120, 120, 120, 255}
		}
		if !d.p.Color {
			msgColor = color.RGBA{255, 255, 255, 255}
		}
		sprites = placeLabel(sprites, ctx, msg, ox, oy, cw, chh, msgColor, "wifi-status")
		return style.ViewData{Sprites: sprites, Static: false, Cursor: -1, PaddingPct: paddingPct}
	}

	barLevel := source.SignalToBarLevel(data.SignalStrength)
	qColor := source.QualityColor(barLevel, accent)
	if !d.p.Color {
		qColor = color.RGBA{255, 255, 255, 255}
	}
	fill := float64(source.SignalPercent(data.SignalStrength)) / 100.0

	// Narrow-tall panels (e.g. 32x128): a horizontal bar is starved by the tiny
	// width while the generous height sits empty. Give these a VERTICAL bar that
	// climbs the long axis as the hero instead — the signal level maps to a
	// filling column, which is both more legible and uses the panel's real asset.
	if chh >= cw*2 {
		return signalForwardCompactVertical(data, bridge, ctx, d, accent, qColor, fill, paddingPct)
	}

	// Three compact bands: SSID, signal bar, IP — fit whichever rows the height allows.
	rowH := clampInt(chh/3, 7, 12)
	y := oy

	// SSID.
	ssid := data.SSID
	if ssid == "" {
		ssid = "(no SSID)"
	}
	sprites = placeLabel(sprites, ctx, ssid, ox, y, cw, rowH, accent, "ssid")
	y += rowH + 1

	// Horizontal signal bar.
	barH := clampInt(rowH, 4, 8)
	if y+barH <= oy+chh {
		bar := progressbar.Render(progressbar.Config{
			Style:       progressbar.Linear,
			Orientation: progressbar.OrientHorizontal,
			Value:       fill,
			Bounds:      image.Rect(ox, y, ox+cw, y+barH),
			Foreground:  qColor,
			Background:  color.RGBA{34, 34, 34, 255},
			RoundedCaps: true,
		})
		if bar != nil {
			sprites = append(sprites, *bar)
		}
		y += barH + 1
	}

	// IP footer, if it still fits.
	if data.IPAddress != "" && y+rowH <= oy+chh {
		dimmed := secondaryColor(d.p.Color)
		sprites = placeLabel(sprites, ctx, data.IPAddress, ox, y, cw, rowH, dimmed, "wifi-ip")
	}

	return style.ViewData{Sprites: sprites, Static: false, Cursor: -1, PaddingPct: paddingPct}
}

// signalForwardCompactVertical is the narrow-tall degradation of the signal-
// forward layout (e.g. 32x128). The panel is too narrow for the arc gauge and
// too narrow for a meaningful horizontal bar, but it is TALL — so the signal
// hero is a VERTICAL bar that climbs the long axis, filling from the bottom in
// proportion to signal strength. A short SSID label sits above it and the IP
// footer below, each only claiming the little height the bar does not need.
func signalForwardCompactVertical(data source.WifiData, bridge layout.LayoutCalculator,
	ctx style.StyleContext, d def, accent, qColor color.RGBA, fill float64, paddingPct int) style.ViewData {

	cw := bridge.AvailableContentWidth()
	chh := bridge.AvailableContentHeight()
	ox, oy := bridge.ContentOrigin()

	var sprites []widgets.Sprite

	// Reserve a compact label band at the top (SSID) and bottom (IP); the bar
	// owns everything between. Label height scales off width so it stays legible
	// on a sliver panel without eating the bar's climb.
	labelH := clampInt(cw/2, 7, 14)
	ipPresent := data.IPAddress != ""

	topBand := labelH
	botBand := 0
	if ipPresent {
		botBand = labelH
	}

	// SSID at the top.
	ssid := data.SSID
	if ssid == "" {
		ssid = "(no SSID)"
	}
	sprites = placeLabel(sprites, ctx, ssid, ox, oy, cw, topBand, accent, "ssid")

	// Vertical bar hero: spans from just below the SSID to just above the IP.
	barTop := oy + topBand + 1
	barBot := oy + chh - botBand - 1
	barH := barBot - barTop
	// A thin margin on each side so the bar does not touch the frame.
	barMargin := clampInt(cw/6, 1, 4)
	barX := ox + barMargin
	barW := cw - 2*barMargin
	if barW < 3 {
		barW = cw
		barX = ox
	}
	if barH >= 8 && barW >= 3 {
		bar := progressbar.Render(progressbar.Config{
			Style:       progressbar.Linear,
			Orientation: progressbar.OrientVertical,
			Value:       fill,
			Bounds:      image.Rect(barX, barTop, barX+barW, barTop+barH),
			Foreground:  qColor,
			Background:  color.RGBA{34, 34, 34, 255},
			RoundedCaps: true,
		})
		if bar != nil {
			sprites = append(sprites, *bar)
		}
	}

	// IP footer at the very bottom, if present.
	if ipPresent {
		dimmed := secondaryColor(d.p.Color)
		sprites = placeLabel(sprites, ctx, data.IPAddress, ox, oy+chh-botBand, cw, botBand, dimmed, "wifi-ip")
	}

	return style.ViewData{Sprites: sprites, Static: false, Cursor: -1, PaddingPct: paddingPct}
}

// signalForwardReadout returns the text shown inside the gauge.
//
// In the default "bars" mode it returns a one-word QUALITY LABEL ("STRONG",
// "GOOD", "FAIR", "WEAK") rather than a number. The word names what the arc
// means — turning an ambiguous dial into an obviously-a-signal-gauge — and it
// reinforces the fill/color in plain language a non-technical viewer reads at a
// glance, without the jargon of dBm or the redundancy of "4/4". The dbm and
// percentage modes instead surface the exact figure for the user who asked to
// see it.
func signalForwardReadout(data source.WifiData, pol source.Policy, barLevel int) string {
	switch pol.SignalDisplay {
	case "dbm":
		return source.FormatDbm(data.SignalStrength)
	case "percentage", "percent":
		return fmt.Sprintf("%d%%", source.SignalPercent(data.SignalStrength))
	default: // "bars": a qualitative word that names the arc.
		return signalQualityWord(barLevel)
	}
}

// signalQualityWord maps a 0–4 bar level to a plain-language quality word shown
// in the gauge center.
func signalQualityWord(barLevel int) string {
	switch barLevel {
	case 4:
		return "STRONG"
	case 3:
		return "GOOD"
	case 2:
		return "FAIR"
	case 1:
		return "WEAK"
	default:
		return "NONE"
	}
}

// signalForwardDetailLine assembles the optional tinkerer detail line from the
// policy toggles. Returns "" when nothing is enabled.
func signalForwardDetailLine(data source.WifiData, pol source.Policy) string {
	var parts []string
	if radio := combinedRadioLine(data, pol); radio != "" {
		parts = append(parts, radio)
	}
	if pol.ShowInterface && data.InterfaceName != "" {
		parts = append(parts, data.InterfaceName)
	}
	if data.LinkSpeed > 0 && (pol.ShowFrequency || pol.ShowChannel || pol.ShowInterface) {
		parts = append(parts, fmt.Sprintf("%dMbps", data.LinkSpeed))
	}
	return joinParts(parts, "  ")
}

// fittedLabel renders text at the largest REAL font face that fits a band of
// targetW×targetH, at that face's native pixel size — no magnification. The
// returned sprite is exactly the size of the rendered text (not the band), so
// the caller centers it within the band itself.
//
// Earlier attempts magnified a small bitmap face up to the band height with
// nearest-neighbor scaling; bitmap fonts do not magnify cleanly and the logical
// box starved the text of width, producing truncated garbage. Picking the
// biggest face that natively fits and drawing it 1:1 is crisp and correct. The
// ceiling is the largest available font asset (spleen-16x32), so on a very large
// panel the text is "only" 32px tall — modest but legible — rather than a
// blown-up smear. Growing past that needs a larger font asset, a separate change.
func fittedLabel(ctx style.StyleContext, text string, targetW, targetH int, fg color.RGBA) *widgets.Sprite {
	if text == "" || targetW <= 0 || targetH <= 0 {
		return nil
	}

	// Choose the largest tier whose native face fits BOTH the band height and the
	// text's width at that face's advance. Walk tiers largest→smallest.
	tiers := []tiercatalog.Tier{
		tiercatalog.TierColossal, tiercatalog.TierHuge, tiercatalog.TierLarge,
		tiercatalog.TierNormal, tiercatalog.TierSmall,
	}
	var face font.Face
	for _, tier := range tiers {
		f := ctx.Face("spleen", tier)
		if f == nil {
			continue
		}
		m := f.Metrics()
		if m.GlyphHeight <= 0 || m.GlyphAdvance <= 0 {
			continue
		}
		textW := len(text) * m.GlyphAdvance
		if m.GlyphHeight <= targetH && textW <= targetW {
			face = f
			break
		}
	}
	// Fallback: smallest face even if it slightly overflows (truncation handles width).
	if face == nil {
		face = ctx.Face("spleen", tiercatalog.TierSmall)
	}
	if face == nil {
		return nil
	}

	m := face.Metrics()
	// Sprite is sized to the text at native scale, capped to the band width so a
	// too-long string is truncated with an ellipsis rather than overflowing.
	spriteW := len(text) * m.GlyphAdvance
	if spriteW > targetW {
		spriteW = targetW
	}
	// Height must be the font's full RowHeight, not GlyphHeight: textbox.Render
	// produces zero ink when the bounds height equals the glyph height exactly
	// (it needs the row's descender slack), so we size to RowHeight which includes
	// it. Draw Left/Top — the sprite is text-sized and placeLabel centers the
	// whole sprite within the band, so textbox's own alignment is unused.
	spriteH := m.RowHeight
	if spriteH <= m.GlyphHeight {
		spriteH = m.GlyphHeight + 2
	}
	return textbox.Render(textbox.Config{
		Bounds:     image.Rect(0, 0, spriteW, spriteH),
		Text:       text,
		Font:       face,
		Alignment:  textbox.Left,
		VAlign:     textbox.Top,
		Overflow:   textbox.Truncate,
		Foreground: fg,
	})
}

func joinParts(parts []string, sep string) string {
	out := ""
	for i, p := range parts {
		if i > 0 {
			out += sep
		}
		out += p
	}
	return out
}

// secondaryColor returns the color for secondary text (IP footer, detail line).
// On color panels it is a dimmed white, giving visual hierarchy below the SSID
// and readout. On mono/grayscale panels it stays full white: the mono panel
// quantizes at luminance 128, and Dim(white) = {127,127,127} (luminance 127)
// falls just BELOW that threshold, so a dimmed footer renders as BLACK and
// vanishes against the background — which is why the IP disappeared on mono.
func secondaryColor(colorPanel bool) color.RGBA {
	if colorPanel {
		return sharedcolor.Dim(color.RGBA{255, 255, 255, 255})
	}
	return color.RGBA{255, 255, 255, 255}
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func clampInt(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

// combinedRadioLine builds a single consolidated radio-info fragment from the
// frequency and channel, honoring the show_frequency / show_channel policy
// toggles. Collapsing both into one fragment ("Ch 36 (5.2GHz)") keeps the detail
// line short. Returns "" when both are disabled or unavailable.
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
