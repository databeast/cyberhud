package styles

import (
	"image"
	"image/color"

	"github.com/databeast/cyberhud/display/modes/thermal/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/style/layout"
	"github.com/databeast/cyberhud/display/widgets/gradient"
)

// buildPortraitHeatmapStyle renders a full-panel gradient background colored by
// severity, with the hottest temperature as a centered text overlay.
// It is the shared BuildFn used by per-resolution styles that want the heatmap layout.
func buildPortraitHeatmapStyle(snapshot source.ThermalSnapshot, pol source.Policy, ctx style.StyleContext, d def) style.ViewData {
	if len(snapshot.Zones) == 0 {
		return style.ViewData{Items: []string{"no thermal data"}}
	}

	hints := ctx.Hints()
	bridge := layout.NewLayoutBridge(hints, layout.BridgeConfig{PaddingPct: 5})

	// Find hottest zone: highest TempC, lowest ZoneID as tiebreaker.
	hottest := snapshot.Zones[0]
	for _, z := range snapshot.Zones[1:] {
		if z.TempC > hottest.TempC || (z.TempC == hottest.TempC && z.ZoneID < hottest.ZoneID) {
			hottest = z
		}
	}

	// Determine severity level and select gradient color stops.
	ec := effectiveCritical(hottest, float64(pol.CritThreshold))
	sev := severity(hottest.TempC, float64(pol.WarnThreshold), ec)

	stops := heatmapGradientStops(sev)

	vd := style.ViewData{}

	// Render gradient covering the full panel bounds.
	if hints.PixelWidth > 0 && hints.PixelHeight > 0 {
		gradSprite := gradient.Render(gradient.Config{
			Style:  gradient.Linear,
			Angle:  180,
			Bounds: image.Rect(0, 0, hints.PixelWidth, hints.PixelHeight),
			Stops:  stops,
		})
		if gradSprite != nil {
			gradSprite.Label = "portrait-heatmap-gradient"
			vd.Sprites = append(vd.Sprites, *gradSprite)
		}
	}

	// Render hottest temperature as centered text overlay.
	tempText := formatTemp(hottest.TempC, pol.Unit)
	vd.Items = append(vd.Items, tempText)
	vd.Colors = append(vd.Colors, color.RGBA{R: 255, G: 255, B: 255, A: 255})

	// Center the text vertically and horizontally using bridge methods.
	offsetX := bridge.CenterX(len(tempText))
	height := bridge.AvailableContentHeight()
	offsetY := (height - bridge.RowHeight()) / 2
	if offsetY < 0 {
		offsetY = 0
	}
	vd.LineOffsets = []int{offsetX}
	vd.OffsetY = offsetY

	return vd
}

// heatmapGradientStops returns the gradient color stops for a given severity level.
//
//	Normal (0):   Dark blue @ 0.0 → Green @ 1.0
//	Warning (1):  Dark blue @ 0.0 → Yellow @ 0.5 → Orange @ 1.0
//	Critical (2): Orange @ 0.0 → Red @ 1.0
func heatmapGradientStops(sev int) []gradient.ColorStop {
	switch sev {
	case 1: // Warning
		return []gradient.ColorStop{
			{Position: 0.0, Color: color.RGBA{R: 0, G: 0, B: 80, A: 255}},
			{Position: 0.5, Color: color.RGBA{R: 255, G: 255, B: 0, A: 255}},
			{Position: 1.0, Color: color.RGBA{R: 255, G: 165, B: 0, A: 255}},
		}
	case 2: // Critical
		return []gradient.ColorStop{
			{Position: 0.0, Color: color.RGBA{R: 255, G: 165, B: 0, A: 255}},
			{Position: 1.0, Color: color.RGBA{R: 255, G: 0, B: 0, A: 255}},
		}
	default: // Normal (0)
		return []gradient.ColorStop{
			{Position: 0.0, Color: color.RGBA{R: 0, G: 0, B: 80, A: 255}},
			{Position: 1.0, Color: color.RGBA{R: 0, G: 255, B: 0, A: 255}},
		}
	}
}
