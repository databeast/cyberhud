package styles

import (
	"fmt"
	"image"

	"github.com/databeast/cyberhud/display/modes/wifi/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/style/layout"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
	"github.com/databeast/cyberhud/display/widgets"
	"github.com/databeast/cyberhud/display/widgets/borderframe"
)

// buildSlowWifi implements the 9-step polish pipeline for all WiFi slow-refresh
// style variants. The panel's MinHeight from reqs determines tier selection.
// This is the default dispatch target for both MonoSlow=true and zero-value Params.
func buildSlowWifi(snap source.WifiData, ctx style.StyleContext, reqs style.SurfaceRequirements) style.ViewData {
	hints := ctx.Hints()

	// Step 1 (Guard): LayoutBridge with PaddingPct:2, early return for zero dims.
	bridge := layout.NewLayoutBridge(hints, layout.BridgeConfig{PaddingPct: 2})
	if bridge.AvailableContentWidth() <= 0 || bridge.AvailableContentHeight() <= 0 {
		return style.ViewData{Items: []string{"(too small)"}, Static: true}
	}

	// Step 2 (Borderframe): Theme "sharp", e-ink safe (all animation/glow/gradient zero).
	frameCfg := borderframe.Config{
		Bounds: image.Rect(0, 0, hints.PixelWidth, hints.PixelHeight),
		Theme:  "sharp",
	}
	frameSprite := borderframe.Render(frameCfg)
	var sprites []widgets.Sprite
	if frameSprite != nil {
		sprites = append(sprites, *frameSprite)
	}

	// Step 3 (Tier resolution): based on ConnectionState and panel height.
	type rowSpec struct {
		text  string
		tier  tiercatalog.Tier
		entry tiercatalog.Entry
	}

	var header []rowSpec
	var secondary []rowSpec

	if snap.ConnectionState == source.Connected {
		// Header tier: TierLarge for panels >= 240, TierHuge for smaller panels.
		var headerTier tiercatalog.Tier
		if reqs.MinHeight >= 240 {
			headerTier = tiercatalog.TierLarge
		} else {
			headerTier = tiercatalog.TierHuge
		}

		headerEntry := ctx.Entry(headerTier)
		smallEntry := ctx.Entry(tiercatalog.TierSmall)

		// Step 4 (Row assembly — Connected): header=SSID, secondary=[signal/quality, IP, channel]
		header = []rowSpec{
			{text: snap.SSID, tier: headerTier, entry: headerEntry},
		}
		secondary = []rowSpec{
			{text: fmt.Sprintf("%d dBm (%d%%)", snap.SignalStrength, snap.LinkQuality), tier: tiercatalog.TierSmall, entry: smallEntry},
			{text: snap.IPAddress, tier: tiercatalog.TierSmall, entry: smallEntry},
			{text: fmt.Sprintf("Ch %d (%.1f GHz)", snap.Channel, snap.Frequency), tier: tiercatalog.TierSmall, entry: smallEntry},
		}
	} else {
		// Step 5 (Row assembly — Disconnected/Unavailable): single row = StatusMessage at TierLarge.
		largeEntry := ctx.Entry(tiercatalog.TierLarge)
		header = []rowSpec{
			{text: snap.StatusMessage, tier: tiercatalog.TierLarge, entry: largeEntry},
		}
		// No secondary rows for disconnected/unavailable state.
	}

	// GroupGap: half header tier row height, minimum 8px.
	groupGap := 0
	if len(header) > 0 {
		groupGap = header[0].entry.RowHeight / 2
		if groupGap < 8 {
			groupGap = 8
		}
	}

	// Step 6 (Adaptive fitting): drop secondary rows from end when height insufficient.
	allRows := append(header, secondary...)
	includeSecondary := len(secondary) > 0

	totalHeight := 0
	for _, r := range header {
		totalHeight += r.entry.RowHeight
	}
	for _, r := range secondary {
		totalHeight += r.entry.RowHeight
	}
	if includeSecondary {
		totalHeight += groupGap
	}

	visibleRows := allRows
	if totalHeight > bridge.AvailableContentHeight() {
		// Drop secondary rows from end until fit.
		includeSecondary = false
		visibleRows = header
		for i := len(secondary); i > 0; i-- {
			candidate := append(header, secondary[:i]...)
			h := 0
			for _, r := range candidate {
				h += r.entry.RowHeight
			}
			h += groupGap
			if h <= bridge.AvailableContentHeight() {
				visibleRows = candidate
				includeSecondary = true
				break
			}
		}
	}

	// Step 7 (Truncation): per-tier maxChars.
	for i := range visibleRows {
		ga := visibleRows[i].entry.GlyphAdvance
		if ga > 0 {
			maxChars := bridge.AvailableContentWidth() / ga
			if maxChars > 0 {
				visibleRows[i].text = textlayout.Truncate(visibleRows[i].text, maxChars)
			}
		}
	}

	// Step 8 (Centering): CenterXWith per-tier, vertical centering.
	offsets := make([]int, len(visibleRows))
	for i, r := range visibleRows {
		offsets[i] = bridge.CenterXWith(len([]rune(r.text)), r.entry.GlyphAdvance)
	}

	var blockHeight int
	for _, r := range visibleRows {
		blockHeight += r.entry.RowHeight
	}
	if includeSecondary && len(visibleRows) > len(header) {
		blockHeight += groupGap
	}
	offsetY := (bridge.AvailableContentHeight() - blockHeight) / 2
	if offsetY < 0 {
		offsetY = 0
	}

	// Step 9 (Emit): ViewData with Items, Tiers, LineOffsets, OffsetY, Static:true, Sprites.
	items := make([]string, len(visibleRows))
	tiers := make([]tiercatalog.Tier, len(visibleRows))
	for i, r := range visibleRows {
		items[i] = r.text
		tiers[i] = r.tier
	}

	return style.ViewData{
		Items:       items,
		Tiers:       tiers,
		LineOffsets: offsets,
		OffsetY:     offsetY,
		Static:      true,
		Sprites:     sprites,
	}
}
