package styles

import (
	"fmt"
	"image/color"

	"github.com/databeast/cyberhud/display/modes/system/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/style/layout"
	"github.com/databeast/cyberhud/display/surface/textlayout"
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

	visibleRows := bridge.MaxVisibleRows()
	count := len(processes)
	if visibleRows > 0 && visibleRows < count {
		count = visibleRows
	}

	maxChars := 0
	if bridge.GlyphAdvance() > 0 {
		maxChars = bridge.AvailableContentWidth() / bridge.GlyphAdvance()
	}

	items := make([]string, 0, count)
	colors := make([]color.Color, 0, count)
	for i := 0; i < count; i++ {
		entry := processes[i]
		cpuStr := fmt.Sprintf("%d%%", entry.CPUPerc)

		name := entry.Name
		// Truncate the NAME so that name + " " + cpuStr fits within maxChars;
		// the percentage must always survive. This prevents the renderer from
		// shrinking an overflowing row to a smaller font (which breaks the
		// leaderboard's type rhythm).
		if maxChars > 0 {
			nameLimit := maxChars - len(cpuStr) - 1 // -1 for the space separator
			if nameLimit < 1 {
				nameLimit = 1
			}
			name = textlayout.Truncate(name, nameLimit)
		}

		items = append(items, name+" "+cpuStr)

		// Accent by load: CPUPerc is 0-100 per core (can exceed 100 on
		// multi-core), so normalize against 100 for the heat ramp.
		load := float64(entry.CPUPerc) / 100.0
		colors = append(colors, loadColor(load, isColor))
	}

	return style.ViewData{Items: items, Colors: colors}
}
