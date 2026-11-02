package styles

import (
	"image"
	"image/color"

	"github.com/databeast/cyberhud/display/modes/wifi/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/style/layout"
	"github.com/databeast/cyberhud/display/widgets"
	"github.com/databeast/cyberhud/display/widgets/borderframe"
	"github.com/databeast/cyberhud/display/widgets/progressbar"
)

// buildSlowWifi renders the signal-first WiFi layout for slow-refresh / e-paper
// panels (MonoSlow and any zero-value-Params style).
//
// It is the e-paper sibling of buildSignalForwardWifi: the same signal-first
// design — a bold signal hero, a plain-language quality word, a quiet SSID with
// a state-tinted icon, and a persistent IP footer — but rendered with widgets
// that suit a 1-bit, no-refresh display. Where the fast engine draws an
// anti-aliased 270° arc gauge, this draws a hard-edged SEGMENTED bar (crisp
// chunks that survive 1-bit quantization and never shimmer), and the whole
// frame is marked Static.
//
// Everything is composited as absolute-positioned sprites (no ViewData.Items
// text rows). That is deliberate: it keeps the frame on a single sprite plane
// (sidestepping the sprite-vs-text two-plane alignment bug), lets the SSID/IP
// auto-fit to the largest font that actually FITS the panel (so a long name is
// shown at a smaller size rather than hacked off with a truncation glyph the
// bitmap font cannot even draw), and frees the layout from counting blank rows
// as occupied height (which previously pushed the ink into the top third and
// abandoned the bottom half of short panels).
func buildSlowWifi(data source.WifiData, pol source.Policy, ctx style.StyleContext, d def) style.ViewData {
	const paddingPct = 4
	bridge := ctx.Layout(paddingPct)

	cw := bridge.AvailableContentWidth()
	chh := bridge.AvailableContentHeight()
	if cw <= 0 || chh <= 0 {
		return style.ViewData{Items: []string{"(too small)"}, Static: true, PaddingPct: paddingPct}
	}

	hints := ctx.Hints()

	// Border: e-ink safe "sharp" theme (no glow/gradient/animation).
	var sprites []widgets.Sprite
	frame := borderframe.Render(borderframe.Config{
		Bounds: image.Rect(0, 0, hints.PixelWidth, hints.PixelHeight),
		Theme:  "sharp",
	})
	if frame != nil {
		sprites = append(sprites, *frame)
	}

	// Accent: color-slow panels honor the policy FGColor; mono/grayscale slow
	// panels render monochrome (white), so quality coloring collapses to white
	// and the segmented bar carries the signal level by FILL, not hue.
	accent := color.RGBA{255, 255, 255, 255}
	if d.p.Color {
		accent = resolveFGColor(pol.FGColor)
	}

	switch data.ConnectionState {
	case source.Disconnected:
		sprites = slowStatus(sprites, bridge, ctx, d, "No Network",
			color.RGBA{255, 60, 60, 255}, source.Disconnected, accent)
	case source.Unavailable:
		sprites = slowStatus(sprites, bridge, ctx, d, "WiFi N/A",
			color.RGBA{120, 120, 120, 255}, source.Unavailable, accent)
	default:
		sprites = slowConnected(sprites, data, pol, bridge, ctx, d, accent)
	}

	return style.ViewData{
		Sprites:    sprites,
		Static:     true,
		Cursor:     -1,
		PaddingPct: paddingPct,
	}
}

// slowStatus renders the bold disconnected / unavailable screen: a large
// state-tinted WiFi icon over a single prominent message, vertically centered.
// Mirrors signalForwardStatus but stays static for e-paper.
func slowStatus(sprites []widgets.Sprite, bridge layout.LayoutCalculator, ctx style.StyleContext,
	d def, msg string, msgColor color.RGBA, state source.ConnectionState, accent color.RGBA) []widgets.Sprite {

	cw := bridge.AvailableContentWidth()
	chh := bridge.AvailableContentHeight()
	ox, oy := bridge.ContentOrigin()

	tint := accent
	textColor := msgColor
	if !d.p.Color {
		tint = color.RGBA{255, 255, 255, 255}
		textColor = color.RGBA{255, 255, 255, 255}
	} else {
		if state == source.Disconnected {
			tint = color.RGBA{255, 60, 60, 255}
		} else if state == source.Unavailable {
			tint = color.RGBA{120, 120, 120, 255}
		}
	}

	// Icon sized to ~a third of the smaller dimension, sitting above center.
	iconSize := minInt(cw, chh) / 3
	if iconSize < wifiIconDst {
		iconSize = wifiIconDst
	}
	if icon := scaledWifiIcon(tint, iconSize); icon != nil {
		iconX := ox + (cw-iconSize)/2
		iconY := oy + (chh/2 - iconSize)
		if iconY < oy {
			iconY = oy
		}
		sprites = append(sprites, widgets.Sprite{Image: icon, Position: image.Point{X: iconX, Y: iconY}, Label: "wifi-icon"})
	}

	msgH := clampInt(chh/5, 12, 48)
	msgY := oy + chh/2 + 4
	sprites = placeLabel(sprites, ctx, msg, ox, msgY, cw, msgH, textColor, "wifi-status")
	return sprites
}

// slowConnected renders the connected hero: a plain-language quality readout, a
// segmented signal bar, a quiet SSID with its icon, and an IP footer — stacked
// and vertically centered as one block so short panels breathe instead of
// abandoning their lower half.
func slowConnected(sprites []widgets.Sprite, data source.WifiData, pol source.Policy,
	bridge layout.LayoutCalculator, ctx style.StyleContext, d def, accent color.RGBA) []widgets.Sprite {

	cw := bridge.AvailableContentWidth()
	chh := bridge.AvailableContentHeight()
	ox, oy := bridge.ContentOrigin()

	barLevel := source.SignalToBarLevel(data.SignalStrength)
	qColor := source.QualityColor(barLevel, accent)
	if !d.p.Color {
		qColor = color.RGBA{255, 255, 255, 255}
	}
	fill := float64(source.SignalPercent(data.SignalStrength)) / 100.0
	dimmed := secondaryColor(d.p.Color)

	// Narrow-tall panels (e.g. 32x128): a horizontal bar is starved by the tiny
	// width while the tall axis sits empty. Render a VERTICAL segmented bar hero
	// climbing the long axis instead — same rationale as the fast engine's
	// signalForwardCompactVertical, with e-paper-safe segmented chunks.
	if chh >= cw*2 {
		return slowConnectedVertical(sprites, data, bridge, ctx, accent, qColor, fill, dimmed)
	}

	// Band heights scale off the panel's smaller dimension so labels stay legible
	// on physically-small-but-dense panels (never fixed pixel caps).
	base := minInt(cw, chh)
	readoutH := clampInt(base/4, 10, 64) // big quality word / figure
	barH := clampInt(base/8, 5, 28)      // segmented signal bar
	ssidH := clampInt(base/7, 8, 40)     // SSID label
	footerH := clampInt(base/9, 7, 32)   // IP / detail line

	detailLine := signalForwardDetailLine(data, pol)

	// --- Decide which bands fit, dropping from the bottom (detail, then IP). ---
	gap := clampInt(base/20, 2, 12)
	bands := []slowBand{
		// Readout: the plain-language quality word (bars mode) or exact figure.
		{readoutH, func(y int) {
			readout := signalForwardReadout(data, pol, barLevel)
			sprites = placeLabel(sprites, ctx, readout, ox, y, cw, readoutH, qColor, "signal-readout")
		}},
		// Segmented signal bar — the e-paper-safe hero. Hard chunks carry the
		// level by fill; color (when available) reinforces quality.
		{barH, func(y int) {
			bar := progressbar.Render(progressbar.Config{
				Style:       progressbar.Segmented,
				Orientation: progressbar.OrientHorizontal,
				Value:       fill,
				Bounds:      image.Rect(ox, y, ox+cw, y+barH),
				Foreground:  qColor,
				Background:  color.RGBA{34, 34, 34, 255},
				SegmentGap:  clampInt(cw/40, 1, 4),
			})
			if bar != nil {
				sprites = append(sprites, *bar)
			}
		}},
		// SSID + state-tinted icon as one centered unit.
		{ssidH, func(y int) {
			sprites = placeSSID(sprites, ctx, data.SSID, accent, ox, y, cw, ssidH)
		}},
	}
	if detailLine != "" {
		bands = append(bands, slowBand{footerH, func(y int) {
			sprites = placeLabel(sprites, ctx, detailLine, ox, y, cw, footerH, dimmed, "wifi-detail")
		}})
	}
	if data.IPAddress != "" {
		bands = append(bands, slowBand{footerH, func(y int) {
			sprites = placeLabel(sprites, ctx, data.IPAddress, ox, y, cw, footerH, dimmed, "wifi-ip")
		}})
	}

	// Fit from the bottom: the first three bands (readout, bar, SSID) are the
	// signal-first core and are kept; detail then IP drop when height runs out.
	const core = 3
	visible := bands
	for len(visible) > core && blockHeight(visible, gap) > chh {
		visible = visible[:len(visible)-1]
	}

	// Vertically center the surviving block.
	total := blockHeight(visible, gap)
	y := oy + (chh-total)/2
	if y < oy {
		y = oy
	}
	for _, b := range visible {
		b.render(y)
		y += b.h + gap
	}
	return sprites
}

// slowConnectedVertical is the narrow-tall slow-path layout (e.g. mono-slow
// 32x128). Like the fast engine's signalForwardCompactVertical, the hero is a
// VERTICAL bar climbing the long axis, with a compact SSID above and IP below.
// The bar is Segmented (hard chunks) rather than Linear, keeping it crisp under
// 1-bit e-paper quantization.
func slowConnectedVertical(sprites []widgets.Sprite, data source.WifiData, bridge layout.LayoutCalculator,
	ctx style.StyleContext, accent, qColor color.RGBA, fill float64, dimmed color.RGBA) []widgets.Sprite {

	cw := bridge.AvailableContentWidth()
	chh := bridge.AvailableContentHeight()
	ox, oy := bridge.ContentOrigin()

	labelH := clampInt(cw/2, 7, 14)
	ipPresent := data.IPAddress != ""
	topBand := labelH
	botBand := 0
	if ipPresent {
		botBand = labelH
	}

	// SSID at the top (name + icon as one centered unit via placeSSID).
	sprites = placeSSID(sprites, ctx, data.SSID, accent, ox, oy, cw, topBand)

	// Vertical segmented bar between the SSID and the IP footer.
	barTop := oy + topBand + 1
	barBot := oy + chh - botBand - 1
	barH := barBot - barTop
	barMargin := clampInt(cw/6, 1, 4)
	barX := ox + barMargin
	barW := cw - 2*barMargin
	if barW < 3 {
		barW = cw
		barX = ox
	}
	if barH >= 8 && barW >= 3 {
		bar := progressbar.Render(progressbar.Config{
			Style:       progressbar.Segmented,
			Orientation: progressbar.OrientVertical,
			Value:       fill,
			Bounds:      image.Rect(barX, barTop, barX+barW, barTop+barH),
			Foreground:  qColor,
			Background:  color.RGBA{34, 34, 34, 255},
			SegmentGap:  clampInt(barH/40, 1, 4),
		})
		if bar != nil {
			sprites = append(sprites, *bar)
		}
	}

	if ipPresent {
		sprites = placeLabel(sprites, ctx, data.IPAddress, ox, oy+chh-botBand, cw, botBand, dimmed, "wifi-ip")
	}
	return sprites
}

// slowBand is one vertically-stacked element of the slow-path connected layout:
// a height and a render closure that draws it at a given top Y.
type slowBand struct {
	h      int
	render func(y int)
}

// blockHeight sums band heights plus inter-band gaps.
func blockHeight(bands []slowBand, gap int) int {
	if len(bands) == 0 {
		return 0
	}
	total := 0
	for _, b := range bands {
		total += b.h
	}
	total += gap * (len(bands) - 1)
	return total
}

// placeSSID renders the SSID name (auto-fitted to the largest font that fits)
// with the state-tinted WiFi icon immediately to its left, the pair centered
// within the band. Falls back to "(no SSID)" when empty.
func placeSSID(sprites []widgets.Sprite, ctx style.StyleContext, ssid string, accent color.RGBA,
	bandX, bandY, bandW, bandH int) []widgets.Sprite {

	if ssid == "" {
		ssid = "(no SSID)"
	}
	iconSize := bandH
	// Reserve horizontal room for the icon so the name fits beside it.
	textW := bandW - iconSize - iconSize/4
	if textW < 1 {
		textW = bandW
	}
	nameSp := fittedLabel(ctx, ssid, textW, bandH, accent)
	if nameSp == nil {
		return sprites
	}
	nameW := nameSp.Image.Bounds().Dx()
	icon := scaledWifiIcon(accent, iconSize)
	iconGap := iconSize / 4
	blockW := nameW
	if icon != nil {
		blockW += iconSize + iconGap
	}
	x := bandX + (bandW-blockW)/2
	if icon != nil {
		iconY := bandY + (bandH-iconSize)/2
		sprites = append(sprites, widgets.Sprite{Image: icon, Position: image.Point{X: x, Y: iconY}, Label: "wifi-icon"})
		x += iconSize + iconGap
	}
	nameSp.Position = image.Point{X: x, Y: bandY + (bandH-nameSp.Image.Bounds().Dy())/2}
	nameSp.Label = "ssid"
	return append(sprites, *nameSp)
}
