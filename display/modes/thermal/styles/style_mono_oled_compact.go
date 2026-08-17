package styles

import (
	"github.com/databeast/cyberhud/display/modes/thermal/source"
	"github.com/databeast/cyberhud/display/style"
)

// buildMonoOLEDCompactStyle is a compact dashboard layout designed for small
// mono OLED panels such as the Waveshare 1.3" 128×64.
//
// Layout (4 rows at typical 128×64 font metrics):
//   - Row 0: hottest zone label + current temperature
//   - Row 1: full-width sparkline history (blank text row, widget sprite)
//   - Row 2: Lo/Hi observed range for hottest zone
//   - Row 3+: compact secondary zones with pixel progress bars
//
// It is the shared BuildFn used by per-resolution styles that want the compact OLED layout.
func buildMonoOLEDCompactStyle(snapshot source.ThermalSnapshot, pol source.Policy, ctx style.StyleContext, d def) style.ViewData {
	if len(snapshot.Zones) == 0 {
		return style.ViewData{Items: []string{"no thermal data"}, Static: true}
	}
	result := buildMonoOLEDCompact(snapshot, pol, ctx)
	if len(result.Items) == 0 || allEmptyItems(result.Items) {
		result.Items = []string{"thermal"}
	}
	result.Static = true
	return result
}
