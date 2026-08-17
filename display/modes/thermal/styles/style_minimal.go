package styles

import (
	"github.com/databeast/cyberhud/display/modes/thermal/source"
	"github.com/databeast/cyberhud/display/style"
)

// buildMinimalStyle renders a single large centered temperature for the hottest zone.
// It is the shared BuildFn used by per-resolution styles that want the minimal layout.
func buildMinimalStyle(snapshot source.ThermalSnapshot, pol source.Policy, ctx style.StyleContext, d def) style.ViewData {
	if len(snapshot.Zones) == 0 {
		return style.ViewData{Items: []string{"no thermal data"}, Static: true}
	}
	result := buildMinimal(snapshot, pol, ctx)
	// Guard: ensure Items is never empty (renderer requires at least one row).
	if len(result.Items) == 0 || allEmptyItems(result.Items) {
		result.Items = []string{"thermal"}
	}
	result.Static = true
	return result
}
