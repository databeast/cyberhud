package styles

import (
	"fmt"
	"image"
	"image/color"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/databeast/cyberhud/display/modes/thermal/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/style/layout"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
	"github.com/databeast/cyberhud/display/widgets/borderframe"
	"github.com/databeast/cyberhud/display/widgets/progressbar"
	"github.com/databeast/cyberhud/display/widgets/scaledtextbox"
	"github.com/databeast/cyberhud/display/widgets/sparkline"
	"github.com/databeast/cyberhud/display/widgets/textbox"
)

// ──────────────────────────────────────────────────────────────────────────────
// Per-resolution style declarations.
//
// Each var is a def following the clock mode pattern:
//   - name: "<category_prefix><WxH>[-<variant>]"
//   - reqs: SurfaceRequirements with MinWidth=W, MinHeight=H, Capability matching prefix
//   - p:    Params{BuildFn: shared core layout function}
//
// Category prefix → Capability mapping:
//   color-          → ColorFast
//   color-slow-     → ColorSlow
//   mono-           → MonoFast
//   mono-slow-      → MonoSlow
//   grayscale-fast- → GrayscaleFast
//   grayscale-slow- → GrayscaleSlow
// ──────────────────────────────────────────────────────────────────────────────

// --- ColorFast 320×240 (landscape) ---

var Color320x240OverviewStyle = def{
	name: "color-320x240-overview",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 240, Capability: style.ColorFast},
	p:    Params{BuildFn: buildOverviewStyle},
}

var Color320x240TimegraphStyle = def{
	name: "color-320x240-timegraph",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 240, Capability: style.ColorFast},
	p:    Params{BuildFn: buildTimegraphStyle},
}

// --- ColorFast 240×320 (portrait) ---

var Color240x320ThermometerStyle = def{
	name: "color-240x320-thermometer",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 320, Capability: style.ColorFast},
	p:    Params{BuildFn: buildPortraitThermometerStyle},
}

var Color240x320SparkStyle = def{
	name: "color-240x320-spark",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 320, Capability: style.ColorFast},
	p:    Params{BuildFn: buildPortraitSparkStyle},
}

var Color240x320HeatmapStyle = def{
	name: "color-240x320-heatmap",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 320, Capability: style.ColorFast},
	p:    Params{BuildFn: buildPortraitHeatmapStyle},
}

var Color240x320LEDsStyle = def{
	name: "color-240x320-leds",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 320, Capability: style.ColorFast},
	p:    Params{BuildFn: buildPortraitLEDsStyle},
}

var Color240x320AvgThermoStyle = def{
	name: "color-240x320-avg-thermo",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 320, Capability: style.ColorFast},
	p:    Params{BuildFn: buildPortraitAvgThermoStyle},
}

// --- MonoSlow 128×64 ---

var MonoSlow128x64Style = def{
	name: "mono-slow-128x64",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 64, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildMonoOLEDCompactStyle},
}

// --- MonoFast 128×128 ---

var Mono128x128Style = def{
	name: "mono-128x128",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 128, Capability: style.MonoFast},
	p:    Params{BuildFn: buildDetailStyle},
}

// --- GrayscaleSlow 296×128 (e-ink) ---

var GrayscaleSlow296x128Style = def{
	name: "grayscale-slow-296x128",
	reqs: style.SurfaceRequirements{MinWidth: 296, MinHeight: 128, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildOverviewStyle},
}

// --- GrayscaleFast 400×300 ---

var GrayscaleFast400x300Style = def{
	name: "grayscale-fast-400x300",
	reqs: style.SurfaceRequirements{MinWidth: 400, MinHeight: 300, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildOverviewStyle},
}


// ──────────────────────────────────────────────────────────────────────────────
// New layout build functions referenced in the per-resolution declarations.
// ──────────────────────────────────────────────────────────────────────────────

// buildTimegraphStyle renders a time-series graph layout. It is functionally
// identical to the graph layout but produces non-static output suitable for
// fast-refresh panels with continuous sparkline animation.
func buildTimegraphStyle(snapshot source.ThermalSnapshot, pol source.Policy, ctx style.StyleContext, d def) style.ViewData {
	if len(snapshot.Zones) == 0 {
		return style.ViewData{Items: []string{"no thermal data"}}
	}
	result := buildGraph(snapshot, pol, ctx)
	if len(result.Items) == 0 || allEmptyItems(result.Items) {
		result.Items = []string{"thermal"}
	}
	// Time-graph is non-static: sparklines animate on fast panels.
	result.Static = false
	return result
}

// buildThermometerStyle renders a vertical thermometer bar layout suitable for
// portrait panels. It delegates to the portrait thermometer core layout.
func buildThermometerStyle(snapshot source.ThermalSnapshot, pol source.Policy, ctx style.StyleContext, d def) style.ViewData {
	return buildPortraitThermometerStyle(snapshot, pol, ctx, d)
}

// ──────────────────────────────────────────────────────────────────────────────
// Shared core layout functions.
//
// These implement the actual rendering logic and are invoked by per-resolution
// styles via Params.BuildFn.
// ──────────────────────────────────────────────────────────────────────────────

// blinkIndicator is the filled square glyph appended to critical temperature text.
const blinkIndicator = " \u25A0"

// buildOverview renders the "overview" style: one row per zone with label,
// formatted temperature, and progress bar fill proportion.
// Items: "label temp [fill_bar]" with blinking indicator for critical zones.
// Colors: per-row severity color.
// Truncated to visible rows ordered by zone ID ascending.
func buildOverview(snap source.ThermalSnapshot, pol source.Policy, ctx style.StyleContext) style.ViewData {
	hints := ctx.Hints()
	bridge := layout.NewLayoutBridge(hints, layout.BridgeConfig{PaddingPct: 5})

	ox, oy := bridge.ContentOrigin()
	width := bridge.AvailableContentWidth()
	hasBorder := pol.ShowBorder && hints.PixelWidth >= 16 && hints.PixelHeight >= 16

	maxRows := bridge.MaxVisibleRows()
	zones := snap.Zones
	if maxRows > 0 && len(zones) > maxRows {
		zones = zones[:maxRows]
	}

	vd := style.ViewData{
		OffsetY: 0,
	}

	// Add borderframe sprites if applicable (panel-covering).
	if hasBorder {
		cfg := borderframe.Config{Bounds: image.Rect(0, 0, hints.PixelWidth, hints.PixelHeight)}
		if pol.FGColor == "none" {
			cfg.ColorTint = color.RGBA{R: 255, G: 255, B: 255, A: 255}
		}
		borderSprite := borderframe.Render(cfg)
		if borderSprite != nil {
			vd.Sprites = append(vd.Sprites, *borderSprite)
		}
	}

	rowHeight := bridge.RowHeight()

	// Derive isColor from policy accent.
	isColor := pol.FGColor != "none"
	nativeFG := color.RGBA{R: 255, G: 255, B: 255, A: 255}

	// Determine progress bar background based on panel type.
	var barBG color.RGBA
	if isColor {
		barBG = color.RGBA{R: 40, G: 40, B: 40, A: 255}
	} else {
		barBG = color.RGBA{R: 0, G: 0, B: 0, A: 0}
	}

	for i, z := range zones {
		ec := effectiveCritical(z, float64(pol.CritThreshold))
		sev := severity(z.TempC, float64(pol.WarnThreshold), ec)
		fill := fillProportion(z.TempC, ec)

		// Resolve accent-driven colors for labels and bar fills.
		labelColor := resolveAccentForLabel(pol.FGColor, isColor, sev, nativeFG)
		barFG := resolveAccentForBar(pol.FGColor, isColor, sev, nativeFG)

		// Build the text row: label + temp + progress bar representation.
		tempStr := formatTemp(z.TempC, pol.Unit)
		if sev == 2 {
			tempStr += blinkIndicator
		}

		// Build a text-based progress bar.
		barWidth := computeBarWidthFromBridge(bridge, z.Label, tempStr)
		bar := renderTextBar(fill, barWidth)

		row := fmt.Sprintf("%s %s %s", z.Label, tempStr, bar)
		vd.Items = append(vd.Items, row)
		vd.Colors = append(vd.Colors, labelColor)

		// Compute Y position for this row using content origin.
		yPos := oy + i*rowHeight

		// Render progressbar widget for this zone using severity-driven bar fill color.
		// Position: to the right of text area, leaving space for label+temp text.
		textChars := utf8.RuneCountInString(z.Label) + 1 + utf8.RuneCountInString(tempStr) + 1
		glyphAdvance := bridge.GlyphAdvance()
		barX := ox + textChars*glyphAdvance
		barPixelWidth := width - textChars*glyphAdvance
		if barPixelWidth < 1 {
			barPixelWidth = 1
		}
		barHeight := rowHeight - 2 // 1px vertical margin within row
		if barHeight < 1 {
			barHeight = 1
		}
		barY := yPos + 1

		barResult := progressbar.Render(progressbar.Config{
			Style:      progressbar.Linear,
			Value:      fill,
			Bounds:     image.Rect(barX, barY, barX+barPixelWidth, barY+barHeight),
			Foreground: barFG,
			Background: barBG,
		})
		if barResult != nil {
			barResult.Label = fmt.Sprintf("thermal-overview-bar-%d", i)
			vd.Sprites = append(vd.Sprites, *barResult)
		}

	}
	return vd
}

// computeBarWidthFromBridge determines how many characters are available for the progress bar
// after accounting for label and temperature text, using the LayoutBridge for geometry.
func computeBarWidthFromBridge(bridge layout.LayoutCalculator, label, tempStr string) int {
	glyphAdvance := bridge.GlyphAdvance()
	if glyphAdvance <= 0 {
		return 8 // default fallback
	}
	maxChars := bridge.AvailableContentWidth() / glyphAdvance
	if maxChars <= 0 {
		return 8 // default fallback
	}
	// Account for: label + space + temp + space + bar
	used := utf8.RuneCountInString(label) + 1 + utf8.RuneCountInString(tempStr) + 1
	remaining := maxChars - used
	if remaining < 3 {
		return 3
	}
	return remaining
}

// renderTextBar returns a text-based progress bar of the given width.
// Uses block characters: "█" for filled, "░" for empty, wrapped in brackets.
func renderTextBar(fill float64, width int) string {
	if width <= 2 {
		return "[]"
	}
	inner := width - 2 // subtract brackets
	filled := int(fill*float64(inner) + 0.5)
	if filled > inner {
		filled = inner
	}
	if filled < 0 {
		filled = 0
	}
	var b strings.Builder
	b.WriteByte('[')
	for i := 0; i < filled; i++ {
		b.WriteRune('\u2588') // █
	}
	for i := filled; i < inner; i++ {
		b.WriteRune('\u2591') // ░
	}
	b.WriteByte(']')
	return b.String()
}

// buildDetail renders the "detail" style: focused view of the hottest zone.
// Shows: label, current temp, min/max observed, sparkline history data, trip points.
// Primary zone: highest temp; lowest zone ID on tie.
//
// Color accent logic:
//   - Zone label uses accent color via resolveAccentForLabel
//   - Temperature and min/max use severity color via resolveAccentForTemp
//   - Sparkline foreground uses severity color via resolveAccentForBar
//
// Sparkline height (scales with panel size):
//   - Panel height ≥ 128px → 2 × rowHeight
//   - Panel height 64–127px → 1 × rowHeight
//   - Panel height < 64px → suppress sparkline and trip point rows
func buildDetail(snap source.ThermalSnapshot, pol source.Policy, ctx style.StyleContext) style.ViewData {
	hints := ctx.Hints()
	bridge := layout.NewLayoutBridge(hints, layout.BridgeConfig{PaddingPct: 5})

	// Derive isColor from the panel's color capability.
	isColor := pol.FGColor != "none"
	nativeFG := color.RGBA{R: 255, G: 255, B: 255, A: 255}

	ox, oy := bridge.ContentOrigin()
	width := bridge.AvailableContentWidth()
	height := bridge.AvailableContentHeight()
	hasBorder := pol.ShowBorder && hints.PixelWidth >= 16 && hints.PixelHeight >= 16

	// Find primary zone: highest temp, lowest zone ID on tie.
	primary := snap.Zones[0]
	for _, z := range snap.Zones[1:] {
		if z.TempC > primary.TempC || (z.TempC == primary.TempC && z.ZoneID < primary.ZoneID) {
			primary = z
		}
	}

	ec := effectiveCritical(primary, float64(pol.CritThreshold))
	sev := severity(primary.TempC, float64(pol.WarnThreshold), ec)

	// Resolve per-element colors using the color accent system.
	labelColor := resolveAccentForLabel(pol.FGColor, isColor, sev, nativeFG)
	tempColor := resolveAccentForTemp(pol.FGColor, isColor, sev, nativeFG)
	sparkFG := resolveAccentForBar(pol.FGColor, isColor, sev, nativeFG)

	// Get history for sparkline and min/max computation.
	history := source.GetHistory(primary.ZoneID)

	// Compute min/max from observed history (non-zero entries).
	minTemp, maxTemp := computeMinMax(history)

	// Build Items.
	vd := style.ViewData{
		OffsetY: 0,
	}

	// Add borderframe sprites if applicable (panel-covering).
	if hasBorder {
		cfg := borderframe.Config{Bounds: image.Rect(0, 0, hints.PixelWidth, hints.PixelHeight)}
		if pol.FGColor == "none" {
			cfg.ColorTint = color.RGBA{R: 255, G: 255, B: 255, A: 255}
		}
		borderSprite := borderframe.Render(cfg)
		if borderSprite != nil {
			vd.Sprites = append(vd.Sprites, *borderSprite)
		}
	}

	rowHeight := bridge.RowHeight()

	// Determine effective panel height for sparkline/trip point suppression.
	suppressSparkAndTrips := height < 64

	// Row 1: label — uses accent color.
	labelRow := primary.Label
	vd.Items = append(vd.Items, labelRow)
	vd.Colors = append(vd.Colors, labelColor)

	// Row 2: current temperature with blinking indicator if critical.
	tempStr := formatTemp(primary.TempC, pol.Unit)
	if sev == 2 {
		tempStr += blinkIndicator
	}
	vd.Items = append(vd.Items, tempStr)
	vd.Colors = append(vd.Colors, tempColor)

	// Row 3: min/max observed.
	minStr := formatTemp(minTemp, pol.Unit)
	maxStr := formatTemp(maxTemp, pol.Unit)
	vd.Items = append(vd.Items, fmt.Sprintf("Lo:%s Hi:%s", minStr, maxStr))
	vd.Colors = append(vd.Colors, tempColor)

	// Sparkline and trip points: suppressed when panel is too short.
	var sparkData []float64
	if !suppressSparkAndTrips {
		// Row 4: sparkline placeholder (rendered as widget below)
		sparkData = normalizeSparkline(history, ec)
		vd.Items = append(vd.Items, "───────") // sparkline placeholder row
		vd.Colors = append(vd.Colors, tempColor)

		// Row 5+: trip points
		for _, tp := range primary.TripPoints {
			tripStr := fmt.Sprintf("%s: %s", tp.Type, formatTemp(tp.TempC, pol.Unit))
			vd.Items = append(vd.Items, tripStr)
			vd.Colors = append(vd.Colors, tempColor)
		}
	}

	// Truncate to visible rows.
	maxRows := bridge.MaxVisibleRows()
	if maxRows > 0 && len(vd.Items) > maxRows {
		vd.Items = vd.Items[:maxRows]
		vd.Colors = vd.Colors[:maxRows]
	}

	// Render sparkline widget for temperature history.
	if !suppressSparkAndTrips {
		sparkY := oy + len(vd.Items)*rowHeight
		sparkWidth := width
		if sparkWidth < 1 {
			sparkWidth = hints.PixelWidth
		}

		// Sparkline height: 2×rowHeight for panels ≥128px tall, 1×rowHeight for 64–127px.
		sparkHeight := rowHeight
		if height >= 128 {
			sparkHeight = 2 * rowHeight
		}

		sparkResult := sparkline.Render(sparkline.Config{
			Data:       sparkData,
			Style:      sparkline.Line,
			Bounds:     image.Rect(ox, sparkY, ox+sparkWidth, sparkY+sparkHeight),
			Foreground: sparkFG,
		})
		if sparkResult != nil {
			sparkResult.Label = "thermal-detail-sparkline"
			vd.Sprites = append(vd.Sprites, *sparkResult)
		}
	}

	return vd
}

// fillProportion computes current_temp / effective_critical clamped to [0.0, 1.0].
// Returns 0.0 when effectiveCrit is zero or negative (avoids division by zero).
func fillProportion(tempC, effectiveCrit float64) float64 {
	if effectiveCrit <= 0 {
		return 0.0
	}
	ratio := tempC / effectiveCrit
	if ratio < 0 {
		return 0.0
	}
	if ratio > 1.0 {
		return 1.0
	}
	return ratio
}

// normalizeSparkline normalizes history samples relative to effective critical,
// clamping each value to [0.0, 1.0]. If effectiveCrit is 0, all values are 0.0.
func normalizeSparkline(history []float64, effectiveCrit float64) []float64 {
	result := make([]float64, len(history))
	if effectiveCrit <= 0 {
		return result
	}
	for i, v := range history {
		ratio := v / effectiveCrit
		if ratio < 0 {
			ratio = 0.0
		}
		if ratio > 1.0 {
			ratio = 1.0
		}
		result[i] = ratio
	}
	return result
}

// computeMinMax returns the minimum and maximum non-zero values from a history slice.
// If all values are zero, returns (0.0, 0.0).
func computeMinMax(history []float64) (min, max float64) {
	first := true
	for _, v := range history {
		if v == 0 {
			continue
		}
		if first {
			min = v
			max = v
			first = false
			continue
		}
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	return min, max
}

// buildGraph renders the "graph" style: stacked sparklines for all zones.
// Each row: first char of label (lowercased) + current temp + sparkline text representation.
// Sparkline width = pixel width minus abbreviation and temp text space.
//
// Color accent logic:
//   - Sparkline foreground = zone severity color via resolveAccentForBar
//   - Text prefix foreground = resolveAccentForLabel (accent-driven or severity)
//   - Monochrome panels: native foreground for all elements
//
// Sparkline widgets are suppressed when effective width < 64px (text-only fallback).
func buildGraph(snap source.ThermalSnapshot, pol source.Policy, ctx style.StyleContext) style.ViewData {
	vd := style.ViewData{
		OffsetY: 0,
	}

	// Determine color panel status and native foreground for monochrome.
	isColor := pol.FGColor != "none"
	nativeFG := color.RGBA{R: 255, G: 255, B: 255, A: 255}

	hints := ctx.Hints()
	bridge := layout.NewLayoutBridge(hints, layout.BridgeConfig{PaddingPct: 5})

	ox, oy := bridge.ContentOrigin()
	width := bridge.AvailableContentWidth()

	rowHeight := bridge.RowHeight()
	glyphAdvance := bridge.GlyphAdvance()

	hasBorder := pol.ShowBorder && hints.PixelWidth >= 16 && hints.PixelHeight >= 16

	// Add borderframe sprites if applicable (panel-covering).
	if hasBorder {
		cfg := borderframe.Config{Bounds: image.Rect(0, 0, hints.PixelWidth, hints.PixelHeight)}
		if pol.FGColor == "none" {
			cfg.ColorTint = color.RGBA{R: 255, G: 255, B: 255, A: 255}
		}
		borderSprite := borderframe.Render(cfg)
		if borderSprite != nil {
			vd.Sprites = append(vd.Sprites, *borderSprite)
		}
	}

	for i, z := range snap.Zones {
		ec := effectiveCritical(z, float64(pol.CritThreshold))
		sev := severity(z.TempC, float64(pol.WarnThreshold), ec)
		sevColor := severityColor(sev)

		// First character of label lowercased (unicode-aware).
		abbr := "?"
		if z.Label != "" {
			r, _ := utf8.DecodeRuneInString(z.Label)
			if r != utf8.RuneError {
				abbr = string(unicode.ToLower(r))
			}
		}

		tempStr := formatTemp(z.TempC, pol.Unit)

		// Compute sparkline width in characters.
		maxChars := 0
		if glyphAdvance > 0 {
			maxChars = width / glyphAdvance
		}
		prefixLen := utf8.RuneCountInString(abbr) + 1 + utf8.RuneCountInString(tempStr) + 1
		sparkWidth := maxChars - prefixLen
		if sparkWidth < 1 {
			sparkWidth = 1
		}

		// Get history and normalize for sparkline.
		history := source.GetHistory(z.ZoneID)
		sparkData := normalizeSparkline(history, ec)

		// Render text sparkline from normalized data.
		sparkText := renderTextSparkline(sparkData, sparkWidth)

		row := fmt.Sprintf("%s %s %s", abbr, tempStr, sparkText)
		vd.Items = append(vd.Items, row)
		vd.Colors = append(vd.Colors, sevColor)

		// Compute Y position and prefix width for sparkline widget placement.
		yPos := oy + i*rowHeight
		prefixPixelWidth := prefixLen * glyphAdvance

		// Suppress sparkline widgets when panel is too narrow (show text only).
		if width >= 64 {
			// Render sparkline widget for this zone stacked vertically.
			sparkX := ox + prefixPixelWidth
			sparkPixelWidth := width - prefixPixelWidth
			if sparkPixelWidth < 1 {
				sparkPixelWidth = 1
			}

			// Resolve sparkline foreground via accent system.
			sparkFG := resolveAccentForBar(pol.FGColor, isColor, sev, nativeFG)

			sparkResult := sparkline.Render(sparkline.Config{
				Data:       sparkData,
				Style:      sparkline.Line,
				Bounds:     image.Rect(sparkX, yPos, sparkX+sparkPixelWidth, yPos+rowHeight),
				Foreground: sparkFG,
			})
			if sparkResult != nil {
				sparkResult.Label = fmt.Sprintf("thermal-graph-spark-%d", i)
				vd.Sprites = append(vd.Sprites, *sparkResult)
			}
		}
	}

	// Truncate to visible rows.
	maxRows := bridge.MaxVisibleRows()
	if maxRows > 0 && len(vd.Items) > maxRows {
		vd.Items = vd.Items[:maxRows]
		vd.Colors = vd.Colors[:maxRows]
	}

	return vd
}

// renderTextSparkline renders normalized [0.0, 1.0] data as text using block characters.
// Uses the last `width` samples from the 64-sample data array.
func renderTextSparkline(data []float64, width int) string {
	// Use the most recent `width` samples.
	start := len(data) - width
	if start < 0 {
		start = 0
	}
	samples := data[start:]

	// Sparkline characters from lowest to highest.
	bars := []rune{'▁', '▂', '▃', '▄', '▅', '▆', '▇', '█'}

	var b strings.Builder
	for _, v := range samples {
		idx := int(v * float64(len(bars)-1))
		if idx < 0 {
			idx = 0
		}
		if idx >= len(bars) {
			idx = len(bars) - 1
		}
		b.WriteRune(bars[idx])
	}
	return b.String()
}

// buildMinimal renders the "minimal" style: single highest temperature centered.
// Font selected with MinVisibleRows=1 for maximum size.
// One decimal place + degree symbol.
// Uses ScaledTextBox when panel ≥ 240×240; falls back to textlabel on nil return.
func buildMinimal(snap source.ThermalSnapshot, pol source.Policy, ctx style.StyleContext) style.ViewData {
	hints := ctx.Hints()
	bridge := layout.NewLayoutBridge(hints, layout.BridgeConfig{PaddingPct: 5})

	// Find highest temperature across all zones.
	highest := snap.Zones[0]
	for _, z := range snap.Zones[1:] {
		if z.TempC > highest.TempC {
			highest = z
		}
	}

	ec := effectiveCritical(highest, float64(pol.CritThreshold))
	sev := severity(highest.TempC, float64(pol.WarnThreshold), ec)

	// Derive isColor from policy accent.
	isColor := pol.FGColor != "none"
	nativeFG := color.RGBA{R: 255, G: 255, B: 255, A: 255}

	// Resolve temperature text color using accent resolution.
	tempColor := resolveAccentForTemp(pol.FGColor, isColor, sev, nativeFG)

	vd := style.ViewData{
		OffsetY: 0,
	}

	tempText := formatTemp(highest.TempC, pol.Unit)
	vd.Items = append(vd.Items, tempText)
	vd.Colors = append(vd.Colors, tempColor)

	// Resolve font face for rendering via tier catalog.
	face := ctx.Face("spleen", tiercatalog.TierNormal)

	// Use face metrics when available; fall back to bridge metrics otherwise.
	var glyphAdvance, rowHeight int
	if face != nil {
		m := face.Metrics()
		glyphAdvance = m.GlyphAdvance
		rowHeight = m.RowHeight
	}
	if glyphAdvance <= 0 {
		glyphAdvance = bridge.GlyphAdvance()
	}
	if rowHeight <= 0 {
		rowHeight = bridge.RowHeight()
	}

	ox, oy := bridge.ContentOrigin()
	width := bridge.AvailableContentWidth()
	height := bridge.AvailableContentHeight()
	hasBorder := pol.ShowBorder && hints.PixelWidth >= 16 && hints.PixelHeight >= 16

	// Add borderframe sprites if applicable (panel-covering).
	if hasBorder {
		cfg := borderframe.Config{Bounds: image.Rect(0, 0, hints.PixelWidth, hints.PixelHeight)}
		if pol.FGColor == "none" {
			cfg.ColorTint = color.RGBA{R: 255, G: 255, B: 255, A: 255}
		}
		borderSprite := borderframe.Render(cfg)
		if borderSprite != nil {
			vd.Sprites = append(vd.Sprites, *borderSprite)
		}
	}

	// Compute centering offsets.
	offsetX := bridge.CenterXWith(utf8.RuneCountInString(tempText), glyphAdvance)
	offsetY := (height - rowHeight) / 2
	if offsetY < 0 {
		offsetY = 0
	}
	vd.Tiers = []tiercatalog.Tier{tiercatalog.TierNormal}
	vd.LineOffsets = []int{offsetX}
	vd.OffsetY = offsetY

	// Use ScaledTextBox for large panels (≥ 240×240) to render temperature prominently.
	if hints.PixelWidth >= 240 && hints.PixelHeight >= 240 && face != nil {
		logicalWidth := glyphAdvance * utf8.RuneCountInString(tempText)
		logicalHeight := rowHeight
		if logicalWidth <= 0 {
			logicalWidth = 1
		}
		if logicalHeight <= 0 {
			logicalHeight = 1
		}

		stbResult := scaledtextbox.Render(scaledtextbox.Config{
			LogicalSize: image.Point{X: logicalWidth, Y: logicalHeight},
			TargetSize:  image.Point{X: width, Y: height},
			Position:    image.Point{X: ox, Y: oy},
			Text:        tempText,
			Font:        face,
			Alignment:   textbox.Center,
			VAlign:      textbox.Middle,
			Foreground:  tempColor,
			Label:       "thermal-minimal-scaled",
		})
		if stbResult != nil {
			vd.Sprites = append(vd.Sprites, *stbResult)
			return vd
		}
		// ScaledTextBox returned nil — fall through to standard text rendering via Items.
	}

	return vd
}

// buildMonoOLEDCompact renders a compact dashboard for small mono OLED panels
// (optimised for 128×64). Uses the available rows as:
//
//   - Row 0:  hottest zone label + current temperature (with severity indicator)
//   - Row 1:  blank — a full-width sparkline history widget occupies this slot
//   - Row 2:  Lo/Hi observed range for the hottest zone
//   - Row 3+: compact secondary zones: abbrev + temp + pixel progress bar
func buildMonoOLEDCompact(snap source.ThermalSnapshot, pol source.Policy, ctx style.StyleContext) style.ViewData {
	hints := ctx.Hints()
	bridge := layout.NewLayoutBridge(hints, layout.BridgeConfig{PaddingPct: 5})

	isColor := pol.FGColor != "none"
	nativeFG := color.RGBA{R: 255, G: 255, B: 255, A: 255}

	ox, oy := bridge.ContentOrigin()
	width := bridge.AvailableContentWidth()
	rowHeight := bridge.RowHeight()
	glyphAdvance := bridge.GlyphAdvance()
	maxRows := bridge.MaxVisibleRows()

	hasBorder := pol.ShowBorder && hints.PixelWidth >= 16 && hints.PixelHeight >= 16

	// Find primary zone: highest temp, lowest ID on tie.
	primary := snap.Zones[0]
	for _, z := range snap.Zones[1:] {
		if z.TempC > primary.TempC || (z.TempC == primary.TempC && z.ZoneID < primary.ZoneID) {
			primary = z
		}
	}

	ec := effectiveCritical(primary, float64(pol.CritThreshold))
	sev := severity(primary.TempC, float64(pol.WarnThreshold), ec)

	tempColor := resolveAccentForTemp(pol.FGColor, isColor, sev, nativeFG)
	labelColor := resolveAccentForLabel(pol.FGColor, isColor, sev, nativeFG)
	sparkFG := resolveAccentForBar(pol.FGColor, isColor, sev, nativeFG)

	history := source.GetHistory(primary.ZoneID)
	minTemp, maxTemp := computeMinMax(history)
	sparkData := normalizeSparkline(history, ec)

	vd := style.ViewData{OffsetY: 0}

	if hasBorder {
		cfg := borderframe.Config{Bounds: image.Rect(0, 0, hints.PixelWidth, hints.PixelHeight)}
		if !isColor {
			cfg.ColorTint = nativeFG
		}
		if s := borderframe.Render(cfg); s != nil {
			vd.Sprites = append(vd.Sprites, *s)
		}
	}

	// Row 0: hottest zone label + current temperature.
	tempStr := formatTemp(primary.TempC, pol.Unit)
	if sev == 2 {
		tempStr += blinkIndicator
	}
	vd.Items = append(vd.Items, fmt.Sprintf("%s %s", primary.Label, tempStr))
	vd.Colors = append(vd.Colors, tempColor)

	// Row 1: blank placeholder — sparkline widget renders here.
	vd.Items = append(vd.Items, "")
	vd.Colors = append(vd.Colors, tempColor)

	sparkY := oy + rowHeight
	sparkH := rowHeight - 2
	if sparkH < 1 {
		sparkH = 1
	}
	if sparkResult := sparkline.Render(sparkline.Config{
		Data:       sparkData,
		Style:      sparkline.Line,
		Bounds:     image.Rect(ox, sparkY, ox+width, sparkY+sparkH),
		Foreground: sparkFG,
	}); sparkResult != nil {
		sparkResult.Label = "thermal-compact-sparkline"
		vd.Sprites = append(vd.Sprites, *sparkResult)
	}

	// Row 2: lo/hi observed range for the hottest zone.
	if maxRows > 2 {
		minStr := formatTemp(minTemp, pol.Unit)
		maxStr := formatTemp(maxTemp, pol.Unit)
		vd.Items = append(vd.Items, fmt.Sprintf("Lo:%s Hi:%s", minStr, maxStr))
		vd.Colors = append(vd.Colors, labelColor)
	}

	// Rows 3+: remaining zones compact (one-char abbrev + temp + progress bar).
	var barBG color.RGBA
	if isColor {
		barBG = color.RGBA{R: 40, G: 40, B: 40, A: 255}
	}

	for i, z := range snap.Zones {
		if z.ZoneID == primary.ZoneID {
			continue
		}
		if maxRows > 0 && len(vd.Items) >= maxRows {
			break
		}

		zec := effectiveCritical(z, float64(pol.CritThreshold))
		zsev := severity(z.TempC, float64(pol.WarnThreshold), zec)
		zfill := fillProportion(z.TempC, zec)
		zColor := resolveAccentForLabel(pol.FGColor, isColor, zsev, nativeFG)
		zBarFG := resolveAccentForBar(pol.FGColor, isColor, zsev, nativeFG)

		abbr := "?"
		if z.Label != "" {
			r, _ := utf8.DecodeRuneInString(z.Label)
			if r != utf8.RuneError {
				abbr = string(unicode.ToLower(r))
			}
		}
		zTempStr := formatTemp(z.TempC, pol.Unit)
		prefixLen := utf8.RuneCountInString(abbr) + 1 + utf8.RuneCountInString(zTempStr) + 1

		maxChars := 0
		if glyphAdvance > 0 {
			maxChars = width / glyphAdvance
		}
		barWidth := maxChars - prefixLen
		if barWidth < 3 {
			barWidth = 3
		}
		bar := renderTextBar(zfill, barWidth)

		rowIdx := len(vd.Items)
		yPos := bridge.RowY(rowIdx, 0)

		vd.Items = append(vd.Items, fmt.Sprintf("%s %s %s", abbr, zTempStr, bar))
		vd.Colors = append(vd.Colors, zColor)

		// Pixel progress bar overlaid on the text bar region.
		prefixPx := prefixLen * glyphAdvance
		barX := ox + prefixPx
		barPxWidth := width - prefixPx
		if barPxWidth < 1 {
			barPxWidth = 1
		}
		barH := rowHeight - 2
		if barH < 1 {
			barH = 1
		}
		if barResult := progressbar.Render(progressbar.Config{
			Style:      progressbar.Linear,
			Value:      zfill,
			Bounds:     image.Rect(barX, yPos+1, barX+barPxWidth, yPos+barH),
			Foreground: zBarFG,
			Background: barBG,
		}); barResult != nil {
			barResult.Label = fmt.Sprintf("thermal-compact-zone-%d", i)
			vd.Sprites = append(vd.Sprites, *barResult)
		}
	}

	return vd
}

// colorToRGBA converts a color.Color to color.RGBA for use with widget configs.
func colorToRGBA(c color.Color) color.RGBA {
	if c == nil {
		return color.RGBA{R: 255, G: 255, B: 255, A: 255}
	}
	if rgba, ok := c.(color.RGBA); ok {
		return rgba
	}
	r, g, b, a := c.RGBA()
	return color.RGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: uint8(a >> 8)}
}

// applyPanelDefaults sets panel-adaptive default values on the given Policy
// based on the panel's color capability and dimensions. Fields whose CLI key
// is present in the explicit map are left unchanged, allowing user-specified
// values to take precedence over automatic defaults.
//
// Panel classification rules:
//  1. Color panel with width ≥ 240 AND height ≥ 240 → detail + border + thermal accent
//  2. Color panel with width < 240 OR height < 240  → overview + no border + thermal accent
//  3. Mono panel with width ≥ 16 AND height ≥ 16    → style by height + border + none accent
//  4. Mono panel with width < 16 OR height < 16     → style by height + no border + none accent
//  5. Zero dimensions (width=0 OR height=0)         → DefaultPolicy values (overview, no border, "none" accent)
//
// The explicit map keys match CLI option keys: "style",
// "fgcolor", "show_led", "show_refresh_bar", "show_border".
func applyPanelDefaults(p *source.Policy, isColor bool, width, height int, explicit map[string]bool) {
	// Determine target defaults based on panel classification.
	var targetStyle string
	var targetAccent string
	var targetShowBorder bool

	if width == 0 || height == 0 {
		// Zero dimensions: use DefaultPolicy values.
		targetStyle = "overview"
		targetAccent = "none"
		targetShowBorder = false
	} else if isColor {
		if width >= 240 && height >= 240 {
			// Color ≥240×240: detail + border + thermal
			targetStyle = "detail"
			targetAccent = "thermal"
			targetShowBorder = true
		} else {
			// Color <240 (either dimension): overview + no border + thermal
			targetStyle = "overview"
			targetAccent = "thermal"
			targetShowBorder = false
		}
	} else {
		// Monochrome panel: accent is always "none".
		targetAccent = "none"
		if height <= 32 {
			targetStyle = "minimal"
		} else if height <= 64 {
			targetStyle = "overview"
		} else {
			targetStyle = "detail"
		}
		// Mono panels ≥ 16×16 get border by default.
		if width >= 16 && height >= 16 {
			targetShowBorder = true
		} else {
			targetShowBorder = false
		}
	}

	// Apply defaults only for fields not explicitly set by the user.
	if !explicit["style"] {
		p.Style = targetStyle
	}
	if !explicit["fgcolor"] {
		p.FGColor = targetAccent
	}
	if !explicit["show_border"] {
		p.ShowBorder = targetShowBorder
	}
}
