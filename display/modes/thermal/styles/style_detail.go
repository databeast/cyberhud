package styles

import (
	"github.com/databeast/cyberhud/display/modes/thermal/source"
	"github.com/databeast/cyberhud/display/style"
)

// buildDetailStyle renders a focused view of the hottest zone with label, temp, min/max, sparkline, and trip points.
// It is the shared BuildFn used by per-resolution styles that want the detail layout.
func buildDetailStyle(snapshot source.ThermalSnapshot, pol source.Policy, ctx style.StyleContext, d def) style.ViewData {
	if len(snapshot.Zones) == 0 {
		return style.ViewData{Items: []string{"no thermal data"}, Static: true}
	}
	result := buildDetail(snapshot, pol, ctx)
	// Guard: ensure Items is never empty (renderer requires at least one row).
	if len(result.Items) == 0 || allEmptyItems(result.Items) {
		result.Items = []string{"thermal"}
	}
	result.Static = true
	return result
}
