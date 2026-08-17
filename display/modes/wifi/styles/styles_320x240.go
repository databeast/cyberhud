package styles

import (
	"fmt"
	"image"
	"image/color"

	"github.com/databeast/cyberhud/display/modes/wifi/source"
	"github.com/databeast/cyberhud/display/style"
	sharedcolor "github.com/databeast/cyberhud/display/style/color"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
	"github.com/databeast/cyberhud/display/widgets"
	"github.com/databeast/cyberhud/display/widgets/progressbar"
)

// WiFi style declarations for 320x240 panels.
//
// Each entry is a hand-tweakable declaration over the core dispatcher in
// core.go: adjust Params to select shared layouts or attach a bespoke BuildFn.

var MonoSlow320x240Style = def{
	name: "mono-slow-320x240",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 240, Capability: style.MonoSlow},
	p:    Params{MonoSlow: true},
}

var GrayscaleSlow320x240Style = def{
	name: "grayscale-slow-320x240",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 240, Capability: style.GrayscaleSlow},
}

var MonoFast320x240Style = def{
	name: "mono-fast-320x240",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 240, Capability: style.MonoFast},
	p:    Params{Fast: true},
}

var GrayscaleFast320x240Style = def{
	name: "grayscale-fast-320x240",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 240, Capability: style.GrayscaleFast},
	p:    Params{Fast: true},
}

var ColorSlow320x240Style = def{
	name: "color-slow-320x240",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 240, Capability: style.ColorSlow},
}

var Color320x240Style = def{
	name: "color-320x240",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 240, Capability: style.ColorFast},
	p:    Params{BuildFn: color320x240Build},
}

func color320x240Build(snap source.WifiData, pol source.Policy, ctx style.StyleContext, d def) style.ViewData {
	bridge := ctx.Layout(2)

	if bridge.AvailableContentWidth() == 0 || bridge.AvailableContentHeight() == 0 {
		return style.ViewData{Items: []string{"(too small)"}}
	}

	accent := resolveFGColor(pol.FGColor)
	white := color.RGBA{255, 255, 255, 255}
	dimmedWhite := sharedcolor.Dim(white)

	// Handle disconnected/unavailable states with centered message.
	if snap.ConnectionState == source.Disconnected {
		cx := bridge.CenterX(len("No Network"))
		cy := bridge.CenterBlockY([]int{bridge.RowHeight()}, 0)
		items := []string{"No Network"}
		colors := []color.Color{color.RGBA{255, 0, 0, 255}}
		offsets := []int{cx}

		var sprites []widgets.Sprite
		icon := renderWifiIcon(source.Disconnected, accent)
		if icon != nil {
			ox, oy := bridge.ContentOrigin()
			if wifiIconDst <= bridge.AvailableContentWidth() && wifiIconDst <= bridge.AvailableContentHeight() {
				sprites = append(sprites, widgets.Sprite{
					Image:    icon,
					Position: image.Point{X: ox, Y: oy},
					Label:    "wifi-icon",
				})
			}
		}

		return style.ViewData{
			Items:       items,
			Colors:      colors,
			Tiers:       []tiercatalog.Tier{tiercatalog.TierLarge},
			LineOffsets: offsets,
			OffsetY:     cy,
			Sprites:     sprites,
		}
	}

	if snap.ConnectionState == source.Unavailable {
		cx := bridge.CenterX(len("WiFi N/A"))
		cy := bridge.CenterBlockY([]int{bridge.RowHeight()}, 0)
		items := []string{"WiFi N/A"}
		colors := []color.Color{dimmedWhite}
		offsets := []int{cx}

		return style.ViewData{
			Items:       items,
			Colors:      colors,
			Tiers:       []tiercatalog.Tier{tiercatalog.TierLarge},
			LineOffsets: offsets,
			OffsetY:     cy,
		}
	}

	// --- Connected state: full zone layout ---

	// Get font metrics from tier catalog.
	smallEntry := ctx.Entry(tiercatalog.TierSmall)
	largeEntry := ctx.Entry(tiercatalog.TierLarge)

	smallRowH := smallEntry.RowHeight
	smallAdvance := smallEntry.GlyphAdvance

	largeRowH := largeEntry.RowHeight
	largeAdvance := largeEntry.GlyphAdvance

	// Build signal display text.
	barLevel := source.SignalToBarLevel(snap.SignalStrength)
	qColor := source.QualityColor(barLevel, accent)
	signalText := ""
	switch pol.SignalDisplay {
	case "dbm":
		signalText = source.FormatDbm(snap.SignalStrength)
	case "percent", "percentage":
		signalText = fmt.Sprintf("%d%%", source.SignalPercent(snap.SignalStrength))
	default: // "bars"
		signalText = fmt.Sprintf("%d/4", barLevel)
	}

	// Build detail items based on policy.
	var detailItems []string
	if pol.ShowFrequency {
		detailItems = append(detailItems, fmt.Sprintf("Freq: %.1f GHz", snap.Frequency))
	}
	if pol.ShowChannel {
		detailItems = append(detailItems, fmt.Sprintf("Ch: %d", snap.Channel))
	}
	if snap.LinkSpeed > 0 {
		detailItems = append(detailItems, fmt.Sprintf("Speed: %d Mbps", snap.LinkSpeed))
	}

	// SSID text (truncated to fit width accounting for icon + gap).
	ssid := snap.SSID
	iconGap := wifiIconDst + 2
	availSSIDWidth := bridge.AvailableContentWidth() - iconGap
	if largeAdvance > 0 && availSSIDWidth > 0 {
		maxChars := availSSIDWidth / largeAdvance
		if maxChars > 0 && len(ssid) > maxChars {
			ssid = textlayout.Truncate(ssid, maxChars)
		}
	}

	// Assemble row heights for FitRows:
	// [SSID (large), Signal (small), ProgressBar (4px), Detail1 (small), Detail2 (small), ..., Footer (small)]
	progressBarH := 4
	rowHeights := []int{largeRowH, smallRowH, progressBarH}
	for range detailItems {
		rowHeights = append(rowHeights, smallRowH)
	}
	// Footer: IP + interface
	footerText := snap.IPAddress
	if pol.ShowInterface && snap.InterfaceName != "" {
		footerText = fmt.Sprintf("%s  %s", snap.IPAddress, snap.InterfaceName)
	}
	rowHeights = append(rowHeights, smallRowH)

	spacing, offsetY, visibleCount := bridge.FitRows(rowHeights)

	// Build items, colors, offsets, and tiers arrays.
	allItems := []string{ssid, signalText}
	allColors := []color.Color{accent, qColor}
	allTiers := []tiercatalog.Tier{tiercatalog.TierLarge, tiercatalog.TierSmall}

	// Progress bar is represented as empty string (sprite-only row).
	allItems = append(allItems, "")
	allColors = append(allColors, color.RGBA{0, 0, 0, 0})
	allTiers = append(allTiers, tiercatalog.TierSmall)

	for _, detail := range detailItems {
		allItems = append(allItems, detail)
		allColors = append(allColors, dimmedWhite)
		allTiers = append(allTiers, tiercatalog.TierSmall)
	}
	allItems = append(allItems, footerText)
	allColors = append(allColors, dimmedWhite)
	allTiers = append(allTiers, tiercatalog.TierSmall)

	// Trim to visible count.
	if visibleCount < len(allItems) {
		allItems = allItems[:visibleCount]
		allColors = allColors[:visibleCount]
		allTiers = allTiers[:visibleCount]
	}

	// Compute horizontal offsets (SSID is offset for icon, rest left-aligned).
	offsets := make([]int, len(allItems))
	offsets[0] = iconGap // SSID starts after icon

	// Truncate detail/footer items to fit panel width.
	maxSmallChars := 0
	if smallAdvance > 0 {
		maxSmallChars = bridge.AvailableContentWidth() / smallAdvance
	}
	for i := 1; i < len(allItems); i++ {
		if maxSmallChars > 0 && len(allItems[i]) > maxSmallChars {
			allItems[i] = textlayout.Truncate(allItems[i], maxSmallChars)
		}
	}

	// --- Sprites ---
	var sprites []widgets.Sprite

	// WiFi icon at top-left content origin.
	ox, oy := bridge.ContentOrigin()
	icon := renderWifiIcon(source.Connected, accent)
	if icon != nil && wifiIconDst <= bridge.AvailableContentWidth() && wifiIconDst <= bridge.AvailableContentHeight() {
		sprites = append(sprites, widgets.Sprite{
			Image:    icon,
			Position: image.Point{X: ox, Y: oy + offsetY},
			Label:    "wifi-icon",
		})
	}

	// Signal bars sprite: positioned on signal row (row index 1).
	barImg := renderSignalBars(barLevel, qColor)
	if barImg != nil {
		// Place to the right of signal text.
		sigTextWidth := len(signalText) * smallAdvance
		barsX := ox + sigTextWidth + 4
		barsY := oy + offsetY + largeRowH + spacing
		// Vertically center bars in the row.
		barYOffset := (smallRowH - signalBarsHeight) / 2
		if barYOffset < 0 {
			barYOffset = 0
		}
		sprites = append(sprites, widgets.Sprite{
			Image:    barImg,
			Position: image.Point{X: barsX, Y: barsY + barYOffset},
			Label:    "signal-bars",
		})
	}

	// Progress bar: positioned on row index 2 (after SSID and Signal rows).
	progressBarY := oy + offsetY + largeRowH + spacing + smallRowH + spacing
	pbarBounds := image.Rect(
		ox,
		progressBarY,
		ox+bridge.AvailableContentWidth(),
		progressBarY+progressBarH,
	)
	pbarFg := qColor
	pbarBg := color.RGBA{30, 30, 30, 255}
	pbarValue := float64(snap.LinkQuality) / 100.0

	pbar := progressbar.Render(progressbar.Config{
		Style:      progressbar.Linear,
		Value:      pbarValue,
		Bounds:     pbarBounds,
		Foreground: pbarFg,
		Background: pbarBg,
	})
	if pbar != nil {
		sprites = append(sprites, *pbar)
	}

	return style.ViewData{
		Items:       allItems,
		Colors:      allColors,
		Tiers:       allTiers,
		LineOffsets: offsets,
		OffsetY:     offsetY + spacing,
		Sprites:     sprites,
	}

}

// Portrait variants: 240×320

var MonoSlow240x320Style = def{
	name: "mono-slow-240x320",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 320, Capability: style.MonoSlow},
}

var MonoFast240x320Style = def{
	name: "mono-fast-240x320",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 320, Capability: style.MonoFast},
	p:    Params{Fast: true},
}

var GrayscaleSlow240x320Style = def{
	name: "grayscale-slow-240x320",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 320, Capability: style.GrayscaleSlow},
}

var GrayscaleFast240x320Style = def{
	name: "grayscale-fast-240x320",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 320, Capability: style.GrayscaleFast},
	p:    Params{Fast: true},
}

var ColorSlow240x320Style = def{
	name: "color-slow-240x320",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 320, Capability: style.ColorSlow},
}

var ColorFast240x320Style = def{
	name: "color-fast-240x320",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 320, Capability: style.ColorFast},
	p:    Params{Fast: true, Color: true},
}
