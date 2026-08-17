package styles

import (
	"image/color"

	"github.com/databeast/cyberhud/display/modes/dashboard/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
	"github.com/databeast/cyberhud/validators"
)

func compactLandscapeLayout(data source.DashboardContent, pol source.Policy, ctx style.StyleContext) (view style.ViewData) {
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

	view.Items = allItems
	view.Colors = colors
	view.Tiers = tiers
	view.LineOffsets = offsets
	view.Static = true

	return view
}

func compactPorttraitLayout(data source.DashboardContent, pol source.Policy, ctx style.StyleContext) (view style.ViewData) {
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

	view.Items = allItems
	view.Colors = colors
	view.Tiers = tiers
	view.LineOffsets = offsets
	view.Static = true

	return view
}

func mediumLandscapeLayout(data source.DashboardContent, pol source.Policy, ctx style.StyleContext) (view style.ViewData) {
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

	view.Items = allItems
	view.Colors = colors
	view.Tiers = tiers
	view.LineOffsets = offsets
	view.Static = true

	return view
}

func mediumPortraitLayout(data source.DashboardContent, pol source.Policy, ctx style.StyleContext) (view style.ViewData) {
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

	view.Items = allItems
	view.Colors = colors
	view.Tiers = tiers
	view.LineOffsets = offsets
	view.Static = true

	return view
}

func largeLandscapeLayout(data source.DashboardContent, pol source.Policy, ctx style.StyleContext) (view style.ViewData) {
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

	view.Items = allItems
	view.Colors = colors
	view.Tiers = tiers
	view.LineOffsets = offsets
	view.Static = true

	return view
}

func largePortraitLayout(data source.DashboardContent, pol source.Policy, ctx style.StyleContext) (view style.ViewData) {
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

	view.Items = allItems
	view.Colors = colors
	view.Tiers = tiers
	view.LineOffsets = offsets
	view.Static = true

	return view
}

// It renders dashboard data rows with no explicit colors, but with a uniform
// TierNormal declared for every row.
//
// This must declare Tiers even though it draws in the panel's monochrome
// default: without it, every row carries neither Tiers nor FontIDs, so the
// renderer's resolveTextFitFonts fallback (region_renderer.go) independently
// picks the largest font that fits each row's own text, one row at a time.
// That per-row font can differ from the single glyphAdvance this function
// used for CenterX/MaxChars above, so rows end up centred against a width
// they are not actually drawn at — visible as a rightward drift and
// right-edge clipping on longer rows. Declaring a uniform tier here forces
// resolveTierFonts to assign every row the same face this function measured
// with, keeping layout and drawing in agreement.
func buildMonoSkeleton(data source.DashboardContent, ctx style.StyleContext) style.ViewData {
	layout := ctx.Layout(0)

	maxRows := layout.MaxVisibleRows()
	if maxRows == 0 {
		return style.ViewData{Items: []string{"DASHBOARD"}, Static: true}
	}

	maxChars := layout.MaxChars()

	items := []string{
		data.Hostname,
		data.Uptime,
		data.IPAddress,
		data.WifiSSID,
		data.Version,
	}

	if len(items) > maxRows {
		items = items[:maxRows]
	}

	rowHeights := make([]int, len(items))
	for i := range items {
		rowHeights[i] = layout.RowHeight()
	}
	_, offsetY, visibleCount := layout.FitRows(rowHeights)
	if visibleCount < len(items) {
		items = items[:visibleCount]
	}

	if maxChars > 0 {
		for i, item := range items {
			items[i] = textlayout.Truncate(item, maxChars)
		}
	}

	offsets := make([]int, len(items))
	for i, item := range items {
		offsets[i] = layout.CenterX(len([]rune(item)))
	}

	tiers := make([]tiercatalog.Tier, len(items))
	for i := range items {
		tiers[i] = tiercatalog.TierNormal
	}

	if validators.AllEmptyItems(items) {
		items = []string{"DASHBOARD"}
		offsets = []int{layout.CenterX(len("DASHBOARD"))}
		tiers = []tiercatalog.Tier{tiercatalog.TierNormal}
	}

	return style.ViewData{
		Items:       items,
		Tiers:       tiers,
		LineOffsets: offsets,
		OffsetY:     offsetY,
		Static:      true,
	}
}

// buildGrayscaleDashboard is the shared Build implementation for all grayscale
// skeleton styles. It declares a uniform TierNormal per row for the same
// reason buildMonoSkeleton does: see that function's comment.
func buildGrayscaleDashboard(data source.DashboardContent, ctx style.StyleContext) style.ViewData {
	layout := ctx.Layout(0)

	maxRows := layout.MaxVisibleRows()
	if maxRows == 0 {
		return style.ViewData{Items: []string{"DASHBOARD"}, Static: true}
	}

	maxChars := layout.MaxChars()

	items := []string{
		data.Hostname,
		data.Uptime,
		data.IPAddress,
		data.WifiSSID,
		data.Version,
	}

	if len(items) > maxRows {
		items = items[:maxRows]
	}

	rowHeights := make([]int, len(items))
	for i := range items {
		rowHeights[i] = layout.RowHeight()
	}
	_, offsetY, visibleCount := layout.FitRows(rowHeights)
	if visibleCount < len(items) {
		items = items[:visibleCount]
	}

	if maxChars > 0 {
		for i, item := range items {
			items[i] = textlayout.Truncate(item, maxChars)
		}
	}

	offsets := make([]int, len(items))
	for i, item := range items {
		offsets[i] = layout.CenterX(len([]rune(item)))
	}

	tiers := make([]tiercatalog.Tier, len(items))
	for i := range items {
		tiers[i] = tiercatalog.TierNormal
	}

	if validators.AllEmptyItems(items) {
		items = []string{"DASHBOARD"}
		offsets = []int{layout.CenterX(len("DASHBOARD"))}
		tiers = []tiercatalog.Tier{tiercatalog.TierNormal}
	}

	return style.ViewData{
		Items:       items,
		Tiers:       tiers,
		LineOffsets: offsets,
		OffsetY:     offsetY,
		Static:      true,
	}
}

// buildColorSkeleton implements the SHORT panel pattern for color skeleton styles.
// All rows use TierNormal + white.
func buildColorSkeleton(data source.DashboardContent, ctx style.StyleContext) style.ViewData {
	layout := ctx.Layout(0)
	maxRows := layout.MaxVisibleRows()
	if maxRows == 0 {
		return style.ViewData{Items: []string{"DASHBOARD"}, Static: true}
	}

	white := color.RGBA{255, 255, 255, 255}

	maxChars := layout.MaxChars()

	items := []string{
		data.Hostname,
		data.Uptime,
		data.IPAddress,
		data.WifiSSID,
		data.Version,
	}

	if len(items) > maxRows {
		items = items[:maxRows]
	}

	rowHeights := make([]int, len(items))
	for i := range items {
		rowHeights[i] = layout.RowHeight()
	}
	_, offsetY, visibleCount := layout.FitRows(rowHeights)
	if visibleCount < len(items) {
		items = items[:visibleCount]
	}

	if maxChars > 0 {
		for i, item := range items {
			items[i] = textlayout.Truncate(item, maxChars)
		}
	}

	colors := make([]color.Color, len(items))
	tiers := make([]tiercatalog.Tier, len(items))
	for i := range items {
		colors[i] = white
		tiers[i] = tiercatalog.TierNormal
	}

	offsets := make([]int, len(items))
	for i, item := range items {
		offsets[i] = layout.CenterX(len([]rune(item)))
	}

	if validators.AllEmptyItems(items) {
		items = []string{"DASHBOARD"}
		colors = []color.Color{white}
		tiers = []tiercatalog.Tier{tiercatalog.TierNormal}
		offsets = []int{layout.CenterX(len("DASHBOARD"))}
	}

	return style.ViewData{
		Items:       items,
		Colors:      colors,
		Tiers:       tiers,
		LineOffsets: offsets,
		OffsetY:     offsetY,
		Static:      true,
	}
}

// buildColorFastSkeleton implements the SHORT panel pattern for color-fast skeleton styles.
// All rows use TierNormal + white.
func buildColorFastSkeleton(data source.DashboardContent, ctx style.StyleContext) style.ViewData {
	layout := ctx.Layout(0)
	maxRows := layout.MaxVisibleRows()
	if maxRows == 0 {
		return style.ViewData{Items: []string{"DASHBOARD"}, Static: true}
	}

	white := color.RGBA{255, 255, 255, 255}

	maxChars := layout.MaxChars()

	items := []string{
		data.Hostname,
		data.Uptime,
		data.IPAddress,
		data.WifiSSID,
		data.Version,
	}

	if len(items) > maxRows {
		items = items[:maxRows]
	}

	rowHeights := make([]int, len(items))
	for i := range items {
		rowHeights[i] = layout.RowHeight()
	}
	_, offsetY, visibleCount := layout.FitRows(rowHeights)
	if visibleCount < len(items) {
		items = items[:visibleCount]
	}

	if maxChars > 0 {
		for i, item := range items {
			items[i] = textlayout.Truncate(item, maxChars)
		}
	}

	colors := make([]color.Color, len(items))
	tiers := make([]tiercatalog.Tier, len(items))
	for i := range items {
		colors[i] = white
		tiers[i] = tiercatalog.TierNormal
	}

	offsets := make([]int, len(items))
	for i, item := range items {
		offsets[i] = layout.CenterX(len([]rune(item)))
	}

	if validators.AllEmptyItems(items) {
		items = []string{"DASHBOARD"}
		colors = []color.Color{white}
		tiers = []tiercatalog.Tier{tiercatalog.TierNormal}
		offsets = []int{layout.CenterX(len("DASHBOARD"))}
	}

	return style.ViewData{
		Items:       items,
		Colors:      colors,
		Tiers:       tiers,
		LineOffsets: offsets,
		OffsetY:     offsetY,
		Static:      true,
	}
}
