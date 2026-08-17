package styles

import (
	"fmt"

	"github.com/databeast/cyberhud/display/modes/system/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/style/layout"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

var TopStyle = def{
	name: "top",
	reqs: style.SurfaceRequirements{
		MinHeight:       32,
		PreferredHeight: 64,
		MinRows:         2,
		MinCharsPerLine: 8,
	},
	p: Params{BuildFn: topBuild},
}

func topBuild(snapshot source.SystemSnapshot, pol source.Policy, ctx style.StyleContext, d def) style.ViewData {
	_ = pol
	processes := snapshot.Processes
	if processes == nil {
		return style.ViewData{Items: []string{"(no data)"}}
	}

	hints := ctx.Hints()
	bridge := layout.NewLayoutBridge(hints, layout.BridgeConfig{PaddingPct: 0})

	visibleRows := bridge.MaxVisibleRows()

	// Cap to min(len(processes), visibleRows).
	count := len(processes)
	if visibleRows > 0 && visibleRows < count {
		count = visibleRows
	}

	maxChars := 0
	if bridge.GlyphAdvance() > 0 {
		maxChars = bridge.AvailableContentWidth() / bridge.GlyphAdvance()
	}

	items := make([]string, 0, count)
	for i := 0; i < count; i++ {
		entry := processes[i]
		cpuStr := fmt.Sprintf("%d%%", entry.CPUPerc)

		name := entry.Name
		// Reserve space for " XX%" suffix (space + cpuStr).
		// Truncate name so that name + " " + cpuStr fits within maxChars.
		if maxChars > 0 {
			nameLimit := maxChars - len(cpuStr) - 1 // -1 for the space separator
			if nameLimit < 1 {
				nameLimit = 1
			}
			name = textlayout.Truncate(name, nameLimit)
		}

		items = append(items, name+" "+cpuStr)
	}

	return style.ViewData{Items: items}
}
