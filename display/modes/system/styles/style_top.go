package styles

import (
	"fmt"
	"image/color"

	"github.com/databeast/cyberhud/display/modes/system/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/style/layout"
)

// buildTopStyle is the shared "top" variant builder — the leaderboard of the
// heaviest CPU-consuming processes: name left, percentage right, worst offender
// on top. Heavy rows are accented with the green→amber→red load color so the
// troublemaker draws the eye (the same cool-to-hot logic as the cores wall).
// Names truncate gracefully; the percentage is always preserved.
func buildTopStyle(snapshot source.SystemSnapshot, pol source.Policy, ctx style.StyleContext, d def) style.ViewData {
	_ = pol
	processes := snapshot.Processes
	if processes == nil {
		return style.ViewData{Items: []string{"(no data)"}}
	}

	hints := ctx.Hints()
	bridge := layout.NewLayoutBridge(hints, layout.BridgeConfig{PaddingPct: 0})
	isColor := isColorPanel(ctx)

	m := ctx.DefaultMetrics()
	visibleRows := bridge.MaxVisibleRows(m)
	count := len(processes)
	if visibleRows > 0 && visibleRows < count {
		count = visibleRows
	}

	cw := bridge.AvailableContentWidth()

	items := make([]string, 0, count)
	colors := make([]color.Color, 0, count)
	for i := 0; i < count; i++ {
		entry := processes[i]
		cpuStr := fmt.Sprintf("%d%%", entry.CPUPerc)

		// The leaderboard keeps a UNIFORM font (one type rhythm), so a long name
		// is clipped rather than sized down — but the percentage must always
		// survive, and the clip must not append an ellipsis (the bitmap fonts have
		// no U+2026 glyph, so it renders as tofu). Budget the name by EXACT pixel
		// width against the panel, reserving the " <pct>" suffix, and clip plainly.
		name := clipNameToWidth(entry.Name, " "+cpuStr, cw, m)

		items = append(items, name+" "+cpuStr)

		// Accent by load: CPUPerc is 0-100 per core (can exceed 100 on
		// multi-core), so normalize against 100 for the heat ramp.
		load := float64(entry.CPUPerc) / 100.0
		colors = append(colors, loadColor(load, isColor))
	}

	return style.ViewData{Items: items, Colors: colors}
}

// clipNameToWidth returns name clipped to the widest rune-prefix that, together
// with the reserved suffix, fits widthPx in font m — using exact TextWidth, with
// no ellipsis (bitmap fonts render U+2026 as tofu). The suffix (" 42%") is never
// clipped; the name yields. A name that already fits is returned unchanged.
func clipNameToWidth(name, suffix string, widthPx int, m layout.Metrics) string {
	if widthPx <= 0 || !m.Usable() {
		return name
	}
	if m.Fits(name+suffix, widthPx) {
		return name
	}
	runes := []rune(name)
	for len(runes) > 0 {
		runes = runes[:len(runes)-1]
		if m.Fits(string(runes)+suffix, widthPx) {
			return string(runes)
		}
	}
	return ""
}
