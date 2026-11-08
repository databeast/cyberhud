package styles

import (
	"image/color"

	"github.com/databeast/cyberhud/display/modes/dashboard/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/style/layout"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
	"github.com/databeast/cyberhud/runtime/host"
)

// dashRow is one candidate row for the adaptive dashboard layout, in priority
// order. Lower-priority rows are dropped first when the panel cannot fit them
// all — so identity and health survive on the smallest panels while version
// and WiFi fall away. This is what the old skeleton lacked: it emitted the
// struct's five string fields in declaration order and let FitRows trim from
// the bottom, which on a 2-row panel meant "hostname, uptime" and silently no
// health at all.
type dashRow struct {
	text string
	// sev carries the row's health severity when it is the health row, so a
	// color panel can tint it. SeverityOK for non-health rows (rendered in the
	// style's default foreground).
	sev      host.Severity
	isHealth bool
}

// prioritizedRows returns the dashboard's rows in importance order:
//
//  1. Hostname  — which box is this?
//  2. Health    — is it okay? (the fact that justifies the mode's existence)
//  3. IP        — how do I reach it?
//  4. Uptime    — the liveness heartbeat value
//  5. WiFi      — droppable context
//  6. Version   — droppable context
//
// Empty/sentinel values are kept (they still answer "no network" etc.); only
// space pressure removes rows, and it removes them from the bottom.
func prioritizedRows(data source.DashboardContent) []dashRow {
	return []dashRow{
		{text: data.Hostname},
		{text: healthText(data.Health), sev: data.Health.Overall, isHealth: true},
		{text: data.IPAddress},
		{text: data.Uptime},
		{text: data.WifiSSID},
		{text: data.Version},
	}
}

// healthText renders the aggregate verdict as a compact, panel-friendly token.
// On healthy systems it is just "OK"; when degraded it names the verdict so a
// glance distinguishes WARN from CRIT without needing color (mono panels get
// the same information color panels get, in words).
func healthText(h host.Health) string {
	return h.Overall.String()
}

// healthOKColor is the foreground for a healthy verdict on color panels. The
// house severityColor (hero.go) renders OK in the panel accent; the small
// adaptive layouts have no accent in scope, so they pass a green here — a
// healthy box reads green, warn amber, crit red (no cyan).
var healthOKColor = color.RGBA{64, 220, 96, 255}

// buildAdaptive is the shared, size- and priority-aware dashboard layout. It
// replaces the near-identical mono/grayscale/color skeletons, which all emitted
// the same five struct fields centered and let the fitter trim from the bottom.
//
// colorize selects whether rows carry explicit colors (color panels) or are
// left to the panel's monochrome default (mono/grayscale). Geometry is
// identical across capabilities — only colorization differs — so one function
// serves all six capability classes.
func buildAdaptive(data source.DashboardContent, ctx style.StyleContext, colorize bool) style.ViewData {
	lc := ctx.Layout(0)

	// Measure with the TierNormal font the renderer will actually draw with. All
	// rows are emitted at TierNormal, so ctx.Metrics(TierNormal) is the exact
	// font the renderer resolves and draws — row-count, max-chars, centering and
	// vertical fit all key off it. This is the correct-by-construction path: a
	// style cannot measure with a font it will not draw.
	m := ctx.Metrics(tiercatalog.TierNormal)

	maxRows := lc.MaxVisibleRows(m)
	if maxRows <= 0 {
		// No room for even one row: show the mode's name so the panel is not blank.
		return style.ViewData{Items: []string{"DASHBOARD"}}
	}

	rows := prioritizedRows(data)

	// Keep only as many rows as fit, taken from the top of the priority list.
	if len(rows) > maxRows {
		rows = rows[:maxRows]
	}

	cw := lc.AvailableContentWidth()

	items := make([]string, len(rows))
	tiers := make([]tiercatalog.Tier, len(rows))
	rowMetrics := make([]layout.Metrics, len(rows))
	var colors []color.Color
	if colorize {
		colors = make([]color.Color, len(rows))
	}

	defaultFG := color.RGBA{255, 255, 255, 255}

	for i, r := range rows {
		// Keep every row — just make it fit. Size the row down from TierNormal
		// until its whole text fits the width, rather than truncating it to a
		// clipped "cyberdeck-0…" that tells a glance nothing. FitTier's floor is
		// the smallest catalog tier; if even that overflows (a pathologically long
		// hostname on a tiny panel) ok is false and we keep the smallest tier and
		// let it hard-clip — a dashboard field is never dropped on width grounds.
		tier, ok := ctx.FitTier(r.text, cw, tiercatalog.TierNormal)
		if !ok {
			if all := ctx.FontCatalog().Tiers(); len(all) > 0 {
				tier = all[0] // ascending: smallest first
			}
		}
		items[i] = r.text
		tiers[i] = tier
		rowMetrics[i] = ctx.Metrics(tier)
		if colorize {
			if r.isHealth {
				colors[i] = severityColor(r.sev, healthOKColor)
			} else {
				colors[i] = defaultFG
			}
		}
	}

	// Vertically fit/center against each row's OWN resolved height — rows may now
	// sit at different tiers, so a single height would mis-center the block.
	rowHeights := make([]int, len(items))
	for i := range items {
		rowHeights[i] = rowMetrics[i].RowHeight()
	}
	_, offsetYRaw, visibleCount := lc.FitRows(rowHeights)
	offsetY := style.TTYY(offsetYRaw)
	if visibleCount < len(items) {
		items = items[:visibleCount]
		tiers = tiers[:visibleCount]
		rowMetrics = rowMetrics[:visibleCount]
		if colorize {
			colors = colors[:visibleCount]
		}
	}

	// Horizontal centering with each row's own metrics.
	offsets := make([]style.TTYX, len(items))
	for i, item := range items {
		offsets[i] = lc.CenterX(len([]rune(item)), rowMetrics[i])
	}

	rowHeightsOut := make([]int, len(items))
	for i := range items {
		rowHeightsOut[i] = rowMetrics[i].RowHeight()
	}

	return style.ViewData{
		Items:        items,
		Colors:       colors,
		Tiers:        tiers,
		LineOffsets:  offsets,
		OffsetY:      offsetY,
		RowHeights:   rowHeightsOut,
		VisibleCount: len(items),
	}
}
