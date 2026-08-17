package styles

import (
	"image"
	"image/color"

	"github.com/databeast/cyberhud/display/modes/thermal/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/style/layout"
	"github.com/databeast/cyberhud/display/widgets/gradient"
	"github.com/databeast/cyberhud/display/widgets/led"
	"github.com/databeast/cyberhud/display/widgets/progressbar"
)

// buildPortraitAvgThermoStyle renders a decorative vertical thermometer gauge showing
// the arithmetic mean temperature across all zones.
// It is the shared BuildFn used by per-resolution styles that want the avg-thermo layout.
func buildPortraitAvgThermoStyle(snapshot source.ThermalSnapshot, pol source.Policy, ctx style.StyleContext, d def) style.ViewData {
	if len(snapshot.Zones) == 0 {
		return style.ViewData{Items: []string{"no thermal data"}}
	}

	hints := ctx.Hints()
	bridge := layout.NewLayoutBridge(hints, layout.BridgeConfig{PaddingPct: 5})

	ox, oy := bridge.ContentOrigin()
	width := bridge.AvailableContentWidth()
	height := bridge.AvailableContentHeight()

	// Compute arithmetic mean of all zone temperatures.
	avgTemp := 0.0
	for _, z := range snapshot.Zones {
		avgTemp += z.TempC
	}
	avgTemp /= float64(len(snapshot.Zones))

	// Use the config crit threshold directly for average display.
	ec := float64(pol.CritThreshold)
	if ec <= 0 {
		ec = 100.0
	}

	// Compute fill proportion and severity based on average temp.
	fill := fillProportion(avgTemp, ec)
	sev := severity(avgTemp, float64(pol.WarnThreshold), ec)

	// Select gradient fill stops by severity.
	gradStops := avgThermoGradientStops(sev)

	// Format temperature label.
	tempLabel := formatTemp(avgTemp, pol.Unit)

	vd := style.ViewData{
		Items:  []string{tempLabel},
		Colors: []color.Color{severityColorRGBA(sev)},
	}

	// Reserve space for the text label at the top.
	labelHeight := bridge.RowHeight()
	if labelHeight <= 0 {
		labelHeight = 12
	}

	// Gap between label and bar: proportional to height.
	gap := height / 40
	if gap < 2 {
		gap = 2
	}

	// Determine bar width (use 75% of content width for aesthetics, capped at 40px).
	barWidth := width
	if barWidth > 40 {
		barWidth = width * 3 / 4
	}
	barX := ox + (width-barWidth)/2

	// LED bulb at the bottom: diameter equals bar width.
	bulbDiameter := barWidth
	if bulbDiameter < 8 {
		bulbDiameter = 8
	}

	// Layout: label at top, then progress bar, then LED bulb at bottom.
	barTop := oy + labelHeight + gap
	bulbY := oy + height - bulbDiameter
	barHeight := bulbY - barTop
	if barHeight < 1 {
		barHeight = 1
	}

	barBG := color.RGBA{R: 40, G: 40, B: 40, A: 255}

	// Render vertical progress bar with gradient fill and rounded end caps.
	barResult := progressbar.Render(progressbar.Config{
		Style:       progressbar.Linear,
		Orientation: progressbar.OrientVertical,
		Value:       fill,
		Bounds:      image.Rect(barX, barTop, barX+barWidth, barTop+barHeight),
		Foreground:  severityColorRGBA(sev),
		Background:  barBG,
		RoundedCaps: true,
		Gradient: &progressbar.GradientFill{
			Stops: gradStops,
		},
	})
	if barResult != nil {
		barResult.Label = "portrait-avg-thermo-bar"
		vd.Sprites = append(vd.Sprites, *barResult)
	}

	// Render LED bulb at the bottom (Circle shape, severity color, glow enabled).
	bulbX := ox + (width-bulbDiameter)/2
	bulbResult := led.Render(led.Config{
		Shape:       led.Circle,
		State:       led.On,
		Brightness:  -1.0,
		Diameter:    bulbDiameter,
		Bounds:      image.Rectangle{Min: image.Pt(bulbX, bulbY)},
		Foreground:  severityColorRGBA(sev),
		GlowEnabled: true,
	})
	if bulbResult != nil {
		bulbResult.Label = "portrait-avg-thermo-bulb"
		vd.Sprites = append(vd.Sprites, *bulbResult)
	}

	// Center the temperature label text above the bar.
	glyphAdvance := bridge.GlyphAdvance()
	if glyphAdvance > 0 {
		textWidth := glyphAdvance * len(tempLabel)
		offsetX := (width - textWidth) / 2
		if offsetX < 0 {
			offsetX = 0
		}
		vd.LineOffsets = []int{offsetX}
	}

	return vd
}

// avgThermoGradientStops returns the gradient color stops for the average
// thermometer fill based on severity level.
//
//	Normal (0):   Green {0, 255, 0, 255} @ 0.0 → Cyan {0, 255, 255, 255} @ 1.0
//	Warning (1):  Yellow {255, 255, 0, 255} @ 0.0 → Orange {255, 165, 0, 255} @ 1.0
//	Critical (2): Orange {255, 165, 0, 255} @ 0.0 → Red {255, 0, 0, 255} @ 1.0
func avgThermoGradientStops(sev int) []gradient.ColorStop {
	switch sev {
	case 1: // Warning
		return []gradient.ColorStop{
			{Position: 0.0, Color: color.RGBA{R: 255, G: 255, B: 0, A: 255}},
			{Position: 1.0, Color: color.RGBA{R: 255, G: 165, B: 0, A: 255}},
		}
	case 2: // Critical
		return []gradient.ColorStop{
			{Position: 0.0, Color: color.RGBA{R: 255, G: 165, B: 0, A: 255}},
			{Position: 1.0, Color: color.RGBA{R: 255, G: 0, B: 0, A: 255}},
		}
	default: // Normal (0)
		return []gradient.ColorStop{
			{Position: 0.0, Color: color.RGBA{R: 0, G: 255, B: 0, A: 255}},
			{Position: 1.0, Color: color.RGBA{R: 0, G: 255, B: 255, A: 255}},
		}
	}
}
