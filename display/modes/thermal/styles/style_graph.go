package styles

import (
	"github.com/databeast/cyberhud/display/modes/thermal/source"
	"github.com/databeast/cyberhud/display/style"
)

// buildGraphStyle renders stacked sparklines for all zones with abbreviated labels and current temps.
// It is the shared BuildFn used by per-resolution styles that want the graph layout.
func buildGraphStyle(snapshot source.ThermalSnapshot, pol source.Policy, ctx style.StyleContext, d def) style.ViewData {
	if len(snapshot.Zones) == 0 {
		return style.ViewData{Items: []string{"no thermal data"}}
	}
	result := buildGraph(snapshot, pol, ctx)
	// Guard: ensure Items is never empty (renderer requires at least one row).
	if len(result.Items) == 0 || allEmptyItems(result.Items) {
		result.Items = []string{"thermal"}
	}
	return result
}
