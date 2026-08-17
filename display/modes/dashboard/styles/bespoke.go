package styles

import (
	"image"
	"image/color"

	"github.com/databeast/cyberhud/display/modes/dashboard/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
	"github.com/databeast/cyberhud/display/widgets"
	"github.com/databeast/cyberhud/display/widgets/borderframe"
)

func buildGrayscaleFast128x128(data source.DashboardContent, pol source.Policy, ctx style.StyleContext, d styleDef) style.ViewData {
	layout := ctx.Layout(0)
	accent := source.ResolveUptimeAccent(pol.ColorAccent)
	items := []string{data.Hostname, data.Uptime, data.IPAddress, data.WifiSSID}
	colors := []color.Color{
		color.RGBA{255, 255, 255, 255},
		accent,
		color.RGBA{255, 255, 255, 255},
		color.RGBA{180, 180, 180, 255},
	}

	// Horizontal centering.
	offsets := make([]int, len(items))
	for i, item := range items {
		offsets[i] = layout.CenterX(len([]rune(item)))
	}

	// Declare a uniform TierNormal per row. Without this, resolveTierFonts
	// (region_renderer.go) leaves Tiers/FontIDs empty and the renderer's
	// resolveTextFitFonts fallback picks the largest font that independently
	// fits each row's own text — a different face per row than the single
	// glyphAdvance used above for CenterX, so rows drift off-centre and can
	// clip past the panel edge.
	tiers := uniformTiers(len(items))

	return style.ViewData{Items: items, Colors: colors, Tiers: tiers, LineOffsets: offsets, Static: true}
}

func buildMonoStyle128x128(data source.DashboardContent, pol source.Policy, ctx style.StyleContext, d styleDef) style.ViewData {
	layout := ctx.Layout(0)

	items := []string{data.Hostname, data.Uptime, data.IPAddress}

	maxRows := layout.MaxVisibleRows()
	if maxRows > 0 && len(items) > maxRows {
		items = items[:maxRows]
	}

	maxChars := 0
	if layout.GlyphAdvance() > 0 {
		maxChars = layout.AvailableContentWidth() / layout.GlyphAdvance()
	}
	if maxChars > 0 {
		for i, item := range items {
			items[i] = textlayout.Truncate(item, maxChars)
		}
	}

	// Horizontal centering.
	offsets := make([]int, len(items))
	for i, item := range items {
		offsets[i] = layout.CenterX(len([]rune(item)))
	}

	return style.ViewData{Items: items, Tiers: uniformTiers(len(items)), LineOffsets: offsets, Static: true}
}

func buildMonoStyle32x128(data source.DashboardContent, pol source.Policy, ctx style.StyleContext, d styleDef) style.ViewData {
	layout := ctx.Layout(0)

	// Single-char abbreviations.
	h := abbreviate(data.Hostname)
	u := abbreviate(data.Uptime)
	ip := abbreviate(data.IPAddress)
	items := []string{h, u, ip}

	maxRows := layout.MaxVisibleRows()
	if maxRows > 0 && len(items) > maxRows {
		items = items[:maxRows]
	}

	// Horizontal centering.
	offsets := make([]int, len(items))
	for i, item := range items {
		offsets[i] = layout.CenterX(len([]rune(item)))
	}

	return style.ViewData{Items: items, Tiers: uniformTiers(len(items)), LineOffsets: offsets, Static: true}
}

func buildMonoStyle128x32(data source.DashboardContent, pol source.Policy, ctx style.StyleContext, d styleDef) style.ViewData {
	layout := ctx.Layout(0)
	items := []string{data.Hostname, data.Uptime, data.IPAddress}

	maxRows := layout.MaxVisibleRows()
	if maxRows > 0 && len(items) > maxRows {
		items = items[:maxRows]
	}

	maxChars := layout.MaxChars()
	if maxChars > 0 {
		for i, item := range items {
			items[i] = textlayout.Truncate(item, maxChars)
		}
	}

	// Horizontal centering.
	offsets := make([]int, len(items))
	for i, item := range items {
		offsets[i] = layout.CenterX(len([]rune(item)))
	}

	return style.ViewData{
		Items:       items,
		Tiers:       uniformTiers(len(items)),
		LineOffsets: offsets,
		Static:      true,
	}
}

func buildMonoStyle64x128(data source.DashboardContent, pol source.Policy, ctx style.StyleContext, d styleDef) style.ViewData {
	layout := ctx.Layout(0)
	items := []string{data.Hostname, data.Uptime, data.IPAddress}

	maxRows := layout.MaxVisibleRows()
	if maxRows > 0 && len(items) > maxRows {
		items = items[:maxRows]
	}

	maxChars := layout.MaxChars()
	if maxChars > 0 {
		for i, item := range items {
			items[i] = textlayout.Truncate(item, maxChars)
		}
	}

	// Horizontal centering.
	offsets := make([]int, len(items))
	for i, item := range items {
		offsets[i] = layout.CenterX(len([]rune(item)))
	}

	return style.ViewData{Items: items, Tiers: uniformTiers(len(items)), LineOffsets: offsets, Static: true}
}

func buildGrayscaleFast128x160(data source.DashboardContent, pol source.Policy, ctx style.StyleContext, d styleDef) style.ViewData {
	_ = ctx.Layout(0)
	accent := source.ResolveUptimeAccent(pol.ColorAccent)
	items := []string{data.Hostname, data.Uptime, data.IPAddress, data.WifiSSID}
	colors := []color.Color{
		color.RGBA{255, 255, 255, 255},
		accent,
		color.RGBA{255, 255, 255, 255},
		color.RGBA{180, 180, 180, 255},
	}
	return style.ViewData{Items: items, Colors: colors, Static: true}
}

func buildColorStyle160x128(data source.DashboardContent, pol source.Policy, ctx style.StyleContext, d styleDef) style.ViewData {
	layout := ctx.Layout(0)

	white := color.RGBA{255, 255, 255, 255}

	// Short panel (h ≤ 160): Hostname + Uptime always first, then secondary
	// fields in order (IP, WiFi, Version, ActivePanelName). Omit from the
	// bottom until rows fit.
	allItems := []string{
		data.Hostname,
		data.Uptime,
		data.IPAddress,
		data.WifiSSID,
		data.Version,
	}

	maxRows := layout.MaxVisibleRows()
	if maxRows <= 0 {
		return style.ViewData{Items: []string{}, Static: true}
	}
	if len(allItems) > maxRows {
		allItems = allItems[:maxRows]
	}

	// All rows: TierNormal + white.
	colors := make([]color.Color, len(allItems))
	tiers := make([]tiercatalog.Tier, len(allItems))
	for i := range allItems {
		colors[i] = white
		tiers[i] = tiercatalog.TierNormal
	}

	// Use FitRows for vertical centering and potential further row reduction.
	rowHeights := make([]int, len(allItems))
	for i := range allItems {
		rowHeights[i] = layout.RowHeight()
	}
	_, offsetY, visibleCount := layout.FitRows(rowHeights)
	if visibleCount < len(allItems) {
		allItems = allItems[:visibleCount]
		colors = colors[:visibleCount]
		tiers = tiers[:visibleCount]
	}

	// Truncate text per row.
	maxChars := 0
	if layout.GlyphAdvance() > 0 {
		maxChars = layout.AvailableContentWidth() / layout.GlyphAdvance()
	}
	if maxChars > 0 {
		for i, item := range allItems {
			allItems[i] = textlayout.Truncate(item, maxChars)
		}
	}

	// Horizontal centering.
	offsets := make([]int, len(allItems))
	for i, item := range allItems {
		offsets[i] = layout.CenterX(len([]rune(item)))
	}

	return style.ViewData{
		Items:       allItems,
		Colors:      colors,
		Tiers:       tiers,
		LineOffsets: offsets,
		OffsetY:     offsetY,
		Static:      true,
	}
}

func buildGrayscaleFast80x160(data source.DashboardContent, pol source.Policy, ctx style.StyleContext, d styleDef) style.ViewData {
	_ = ctx.Layout(0)
	accent := source.ResolveUptimeAccent(pol.ColorAccent)
	// Single combined line: hostname + uptime
	line := data.Hostname + " " + data.Uptime
	return style.ViewData{Items: []string{line}, Colors: []color.Color{accent}, Static: true}
}

func buildColorStyle160x80(data source.DashboardContent, pol source.Policy, ctx style.StyleContext, d styleDef) style.ViewData {
	layout := ctx.Layout(0)
	white := color.RGBA{255, 255, 255, 255}

	// Short panel (h ≤ 160): Hostname + Uptime always first, then secondary
	// fields in order (IP, WiFi, Version, ActivePanelName). Omit from the
	// bottom until rows fit.
	allItems := []string{
		data.Hostname,
		data.Uptime,
		data.IPAddress,
		data.WifiSSID,
		data.Version,
	}

	maxRows := layout.MaxVisibleRows()
	if maxRows <= 0 {
		return style.ViewData{Items: []string{}, Static: true}
	}
	if len(allItems) > maxRows {
		allItems = allItems[:maxRows]
	}

	// All rows: TierNormal + white.
	colors := make([]color.Color, len(allItems))
	tiers := make([]tiercatalog.Tier, len(allItems))
	for i := range allItems {
		colors[i] = white
		tiers[i] = tiercatalog.TierNormal
	}

	// Use FitRows for vertical centering and potential further row reduction.
	rowHeights := make([]int, len(allItems))
	for i := range allItems {
		rowHeights[i] = layout.RowHeight()
	}
	_, offsetY, visibleCount := layout.FitRows(rowHeights)
	if visibleCount < len(allItems) {
		allItems = allItems[:visibleCount]
		colors = colors[:visibleCount]
		tiers = tiers[:visibleCount]
	}

	// Truncate text per row.
	maxChars := layout.MaxChars()
	if maxChars > 0 {
		for i, item := range allItems {
			allItems[i] = textlayout.Truncate(item, maxChars)
		}
	}

	// Horizontal centering.
	offsets := make([]int, len(allItems))
	for i, item := range allItems {
		offsets[i] = layout.CenterX(len([]rune(item)))
	}

	return style.ViewData{
		Items:       allItems,
		Colors:      colors,
		Tiers:       tiers,
		LineOffsets: offsets,
		OffsetY:     offsetY,
		Static:      true,
	}
}

func buildEinkStyle200x200(data source.DashboardContent, pol source.Policy, ctx style.StyleContext, d styleDef) style.ViewData {
	layout := ctx.Layout(0)

	wctx := widgets.SuppressionContext{
		AvailableWidth:  layout.AvailableContentWidth(),
		AvailableHeight: layout.AvailableContentHeight(),
		IsEink:          true,
	}
	comp := widgets.NewCompositor(wctx, widgets.SuppressOnEink())

	items := []string{data.Hostname, data.Uptime, data.IPAddress}

	// Truncate to visible rows
	maxRows := layout.MaxVisibleRows()
	if maxRows > 0 && len(items) > maxRows {
		items = items[:maxRows]
	}

	// Truncate each row
	maxChars := 0
	if layout.GlyphAdvance() > 0 {
		maxChars = layout.AvailableContentWidth() / layout.GlyphAdvance()
	}
	if maxChars > 0 {
		for i, item := range items {
			items[i] = textlayout.Truncate(item, maxChars)
		}
	}

	// Horizontal centering.
	offsets := make([]int, len(items))
	for i, item := range items {
		offsets[i] = layout.CenterX(len([]rune(item)))
	}

	return style.ViewData{
		Items:       items,
		Tiers:       uniformTiers(len(items)),
		LineOffsets: offsets,
		Static:      true,
		Sprites:     comp.Sprites(),
	}
}

func buildColorStyle135x240(data source.DashboardContent, pol source.Policy, ctx style.StyleContext, d styleDef) style.ViewData {
	layout := ctx.Layout(0)

	accent := source.ResolveUptimeAccent(pol.ColorAccent)
	white := color.RGBA{255, 255, 255, 255}
	dimmed := color.RGBA{180, 180, 180, 255}

	items := []string{
		data.Hostname,
		data.Uptime,
		data.IPAddress,
		data.WifiSSID,
		data.Version,
	}
	colors := []color.Color{
		accent, // hostname: accent
		white,  // uptime: white
		dimmed, // IP: dimmed
		dimmed, // WiFi SSID: dimmed
		dimmed, // version: dimmed
		dimmed, // active panel name: dimmed
	}
	tiers := []tiercatalog.Tier{
		tiercatalog.TierLarge, // hostname
		tiercatalog.TierLarge, // uptime
		tiercatalog.TierSmall, // IP
		tiercatalog.TierSmall, // WiFi SSID
		tiercatalog.TierSmall, // version
		tiercatalog.TierSmall, // active panel name
	}

	// Truncate to visible rows.
	maxRows := layout.MaxVisibleRows()
	if maxRows > 0 && len(items) > maxRows {
		items = items[:maxRows]
		colors = colors[:maxRows]
		tiers = tiers[:maxRows]
	}

	// Truncate each row's text to fit available width.
	maxChars := 0
	if layout.GlyphAdvance() > 0 {
		maxChars = layout.AvailableContentWidth() / layout.GlyphAdvance()
	}
	if maxChars > 0 {
		for i, item := range items {
			items[i] = textlayout.Truncate(item, maxChars)
		}
	}

	// Horizontal centering.
	offsets := make([]int, len(items))
	for i, item := range items {
		offsets[i] = layout.CenterX(len([]rune(item)))
	}

	return style.ViewData{Items: items, Colors: colors, Tiers: tiers, LineOffsets: offsets, Static: true}
}

func buildGrayscaleFast240x135(data source.DashboardContent, pol source.Policy, ctx style.StyleContext, d styleDef) style.ViewData {
	_ = ctx.Layout(0)
	accent := source.ResolveUptimeAccent(pol.ColorAccent)

	items := []string{
		data.Hostname,
		data.Uptime,
		data.IPAddress,
		data.WifiSSID,
	}
	colors := []color.Color{
		color.RGBA{255, 255, 255, 255}, // hostname: white
		accent,                         // uptime: accent
		color.RGBA{255, 255, 255, 255}, // IP: white
		color.RGBA{180, 180, 180, 255}, // clock: dimmed
	}

	return style.ViewData{Items: items, Colors: colors, Static: true}
}

func buildColorStyle240x240(data source.DashboardContent, pol source.Policy, ctx style.StyleContext, d styleDef) style.ViewData {
	layout := ctx.Layout(0)

	accent := source.ResolveUptimeAccent(pol.ColorAccent)
	white := color.RGBA{255, 255, 255, 255}
	dimmed := color.RGBA{180, 180, 180, 255}

	items := []string{
		data.Hostname,
		data.Uptime,
		data.IPAddress,
		data.WifiSSID,
		data.Version,
	}
	colors := []color.Color{
		accent, // hostname: accent
		white,  // uptime: white
		dimmed, // IP: dimmed
		dimmed, // WiFi SSID: dimmed
		dimmed, // version: dimmed
		dimmed, // active panel name: dimmed
	}
	tiers := []tiercatalog.Tier{
		tiercatalog.TierLarge, // hostname
		tiercatalog.TierLarge, // uptime
		tiercatalog.TierSmall, // IP
		tiercatalog.TierSmall, // WiFi SSID
		tiercatalog.TierSmall, // version
		tiercatalog.TierSmall, // active panel name
	}

	// Truncate to visible rows.
	maxRows := layout.MaxVisibleRows()
	if maxRows > 0 && len(items) > maxRows {
		items = items[:maxRows]
		colors = colors[:maxRows]
		tiers = tiers[:maxRows]
	}

	// Use FitRows for vertical centering and potential further row reduction.
	rowHeights := make([]int, len(items))
	for i := range items {
		rowHeights[i] = layout.RowHeight()
	}
	_, offsetY, visibleCount := layout.FitRows(rowHeights)
	if visibleCount < len(items) {
		items = items[:visibleCount]
		colors = colors[:visibleCount]
		tiers = tiers[:visibleCount]
	}

	// Truncate each row's text to fit available width.
	maxChars := layout.MaxChars()
	if maxChars > 0 {
		for i, item := range items {
			items[i] = textlayout.Truncate(item, maxChars)
		}
	}

	// Horizontal centering.
	offsets := make([]int, len(items))
	for i, item := range items {
		offsets[i] = layout.CenterX(len([]rune(item)))
	}

	return style.ViewData{Items: items, Colors: colors, Tiers: tiers, LineOffsets: offsets, OffsetY: offsetY, Static: true}
}

func buildColorStyle320x240(data source.DashboardContent, pol source.Policy, ctx style.StyleContext, d styleDef) style.ViewData {
	layout := ctx.Layout(0)

	accent := source.ResolveUptimeAccent(pol.ColorAccent)
	white := color.RGBA{255, 255, 255, 255}
	dimmed := color.RGBA{180, 180, 180, 255}

	items := []string{
		data.Hostname,
		data.Uptime,
		data.IPAddress,
		data.WifiSSID,
		data.Version,
	}
	colors := []color.Color{
		accent, // hostname: accent
		white,  // uptime: white
		dimmed, // IP: dimmed
		dimmed, // WiFi SSID: dimmed
		dimmed, // version: dimmed
		dimmed, // active panel name: dimmed
	}
	tiers := []tiercatalog.Tier{
		tiercatalog.TierLarge, // hostname
		tiercatalog.TierLarge, // uptime
		tiercatalog.TierSmall, // IP
		tiercatalog.TierSmall, // WiFi SSID
		tiercatalog.TierSmall, // version
		tiercatalog.TierSmall, // active panel name
	}

	// Truncate to visible rows.
	maxRows := layout.MaxVisibleRows()
	if maxRows > 0 && len(items) > maxRows {
		items = items[:maxRows]
		colors = colors[:maxRows]
		tiers = tiers[:maxRows]
	}

	// Use FitRows for vertical centering and potential further row reduction.
	rowHeights := make([]int, len(items))
	for i := range items {
		rowHeights[i] = layout.RowHeight()
	}
	_, offsetY, visibleCount := layout.FitRows(rowHeights)
	if visibleCount < len(items) {
		items = items[:visibleCount]
		colors = colors[:visibleCount]
		tiers = tiers[:visibleCount]
	}

	// Truncate each row's text to fit available width.
	maxChars := 0
	if layout.GlyphAdvance() > 0 {
		maxChars = layout.AvailableContentWidth() / layout.GlyphAdvance()
	}
	if maxChars > 0 {
		for i, item := range items {
			items[i] = textlayout.Truncate(item, maxChars)
		}
	}

	// Horizontal centering.
	offsets := make([]int, len(items))
	for i, item := range items {
		offsets[i] = layout.CenterX(len([]rune(item)))
	}

	return style.ViewData{Items: items, Colors: colors, Tiers: tiers, LineOffsets: offsets, OffsetY: offsetY, Static: true}
}

func buildEinkStyle800x480(data source.DashboardContent, pol source.Policy, ctx style.StyleContext, d styleDef) style.ViewData {
	layout := ctx.Layout(0)

	// Step 2 (Borderframe): render decorative border sprite.
	frameCfg := borderframe.Config{
		Bounds: image.Rect(0, 0, ctx.Hints().PixelWidth, ctx.Hints().PixelHeight),
		Theme:  "sharp",
		// All animation/glow/gradient fields remain at zero value (e-ink safe).
	}
	frameSprite := borderframe.Render(frameCfg)

	// If borderframe returns nil (bounds < 16×16), omit sprite but proceed.
	var sprites []widgets.Sprite
	if frameSprite != nil {
		sprites = append(sprites, *frameSprite)
	}

	// Step 4 (Tier resolution): query catalog for per-tier metrics.
	hugeEntry := ctx.Entry(tiercatalog.TierHuge)
	largeEntry := ctx.Entry(tiercatalog.TierLarge)

	// Step 5 (Row assembly): define row groups with tier metadata.
	type rowSpec struct {
		text  string
		tier  tiercatalog.Tier
		entry tiercatalog.Entry
	}

	primary := []rowSpec{
		{text: data.Hostname, tier: tiercatalog.TierHuge, entry: hugeEntry},
		{text: data.Uptime, tier: tiercatalog.TierLarge, entry: largeEntry},
		{text: data.IPAddress, tier: tiercatalog.TierLarge, entry: largeEntry},
		{text: data.WifiSSID, tier: tiercatalog.TierLarge, entry: largeEntry},
	}

	groupGap := largeEntry.RowHeight / 2
	if groupGap < 8 {
		groupGap = 8
	}

	// Step 6 (Adaptive fitting): determine visible rows.
	visibleRows := append([]rowSpec(nil), primary...)
	includeSecondary := true

	// Compute total height with all rows + groupGap.
	totalHeight := 0
	for _, r := range primary {
		totalHeight += r.entry.RowHeight
	}

	totalHeight += groupGap // inter-group gap

	if totalHeight > layout.AvailableContentHeight() {
		// Drop secondary group.
		visibleRows = primary
		includeSecondary = false
	}

	// If still insufficient with primary only, use FitRows to drop from end.
	if !includeSecondary {
		rowHeights := make([]int, len(primary))
		for i, r := range primary {
			rowHeights[i] = r.entry.RowHeight
		}
		_, _, visibleCount := layout.FitRows(rowHeights)
		if visibleCount < len(primary) {
			visibleRows = primary[:visibleCount]
		}
	}

	// Step 7 (Truncation): per-tier text truncation.
	for i := range visibleRows {
		ga := visibleRows[i].entry.GlyphAdvance
		if ga > 0 {
			maxChars := layout.AvailableContentWidth() / ga
			if maxChars > 0 {
				visibleRows[i].text = textlayout.Truncate(visibleRows[i].text, maxChars)
			}
		}
	}

	// Step 8 (Centering): horizontal per-tier centering.
	offsets := make([]int, len(visibleRows))
	for i, r := range visibleRows {
		offsets[i] = layout.CenterXWith(len([]rune(r.text)), r.entry.GlyphAdvance)
	}

	// Vertical centering: compute total block height and center it.
	rowHeightsForCenter := make([]int, len(visibleRows))
	for i, r := range visibleRows {
		rowHeightsForCenter[i] = r.entry.RowHeight
	}

	// Determine intra-row spacing and group gap for vertical centering.
	var blockHeight int
	if includeSecondary && len(visibleRows) > len(primary) {
		// Total block = sum of row heights + groupGap between primary and secondary
		for _, h := range rowHeightsForCenter {
			blockHeight += h
		}
		blockHeight += groupGap
	} else {
		for _, h := range rowHeightsForCenter {
			blockHeight += h
		}
	}

	offsetY := (layout.AvailableContentHeight() - blockHeight) / 2
	if offsetY < 0 {
		offsetY = 0
	}

	// Build Items and Tiers from visible rows.
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

func abbreviate(s string) string {
	for _, r := range s {
		return string(r)
	}
	return ""
}

// uniformTiers returns a slice of n TierNormal entries. Used by bespoke
// builders that lay out and centre text against the region's single baseline
// glyph advance (via layout.CenterX/GlyphAdvance), so every row must resolve
// to that same TierNormal face rather than being independently re-picked by
// the renderer's resolveTextFitFonts fallback.
func uniformTiers(n int) []tiercatalog.Tier {
	tiers := make([]tiercatalog.Tier, n)
	for i := range tiers {
		tiers[i] = tiercatalog.TierNormal
	}
	return tiers
}
