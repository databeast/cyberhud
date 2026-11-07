package styles

import (
	"image/color"

	"github.com/databeast/cyberhud/display/modes/dashboard/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
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

	maxRows := layout.MaxVisibleRows(ctx.Metrics(tiercatalog.TierNormal))
	if maxRows <= 0 {
		return style.ViewData{Items: []string{}}
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
	maxChars := layout.MaxChars(ctx.Metrics(tiercatalog.TierNormal))
	if maxChars > 0 {
		for i, item := range allItems {
			allItems[i] = textlayout.Truncate(item, maxChars)
		}
	}

	// Horizontal centering.
	offsets := make([]style.TTYX, len(allItems))
	for i, item := range allItems {
		offsets[i] = layout.CenterX(len([]rune(item)), ctx.Metrics(tiercatalog.TierNormal))
	}

	view.Items = allItems
	view.Colors = colors
	view.Tiers = tiers
	view.LineOffsets = offsets

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

	maxRows := layout.MaxVisibleRows(ctx.Metrics(tiercatalog.TierNormal))
	if maxRows <= 0 {
		return style.ViewData{Items: []string{}}
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
	maxChars := layout.MaxChars(ctx.Metrics(tiercatalog.TierNormal))
	if maxChars > 0 {
		for i, item := range allItems {
			allItems[i] = textlayout.Truncate(item, maxChars)
		}
	}

	// Horizontal centering.
	offsets := make([]style.TTYX, len(allItems))
	for i, item := range allItems {
		offsets[i] = layout.CenterX(len([]rune(item)), ctx.Metrics(tiercatalog.TierNormal))
	}

	view.Items = allItems
	view.Colors = colors
	view.Tiers = tiers
	view.LineOffsets = offsets

	return view
}

// The four skeleton builders below are the Build implementations the core.go
// switch dispatches to per capability. They are thin wrappers over buildAdaptive
// (see adaptive.go): mono and grayscale render without explicit colors; color
// and color-fast colorize the health row by severity. The size- and
// priority-aware layout logic lives in buildAdaptive so all six capability
// classes share one real layout instead of four near-identical copies that each
// emitted the struct's five fields in declaration order.

func buildMonoSkeleton(data source.DashboardContent, ctx style.StyleContext) style.ViewData {
	return buildAdaptive(data, ctx, false)
}

func buildGrayscaleDashboard(data source.DashboardContent, ctx style.StyleContext) style.ViewData {
	return buildAdaptive(data, ctx, false)
}

func buildColorSkeleton(data source.DashboardContent, ctx style.StyleContext) style.ViewData {
	return buildAdaptive(data, ctx, true)
}

func buildColorFastSkeleton(data source.DashboardContent, ctx style.StyleContext) style.ViewData {
	return buildAdaptive(data, ctx, true)
}
