package styles

import (
	"image"
	"image/color"

	"github.com/databeast/cyberhud/display/modes/thermal/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/style/layout"
	"github.com/databeast/cyberhud/display/widgets/progressbar"
)

// buildPortraitThermometerStyle renders a vertical progress bar thermometer
// showing the hottest zone's temperature on portrait side panels.
// It is the shared BuildFn used by per-resolution styles that want the thermometer layout.
func buildPortraitThermometerStyle(snapshot source.ThermalSnapshot, pol source.Policy, ctx style.StyleContext, d def) style.ViewData {
	if len(snapshot.Zones) == 0 {
		return style.ViewData{Items: []string{"no thermal data"}}
	}

	hints := ctx.Hints()
	bridge := layout.NewLayoutBridge(hints, layout.BridgeConfig{PaddingPct: 5})

	ox, oy := bridge.ContentOrigin()
	width := bridge.AvailableContentWidth()
	height := bridge.AvailableContentHeight()

	// Find hottest zone: highest TempC, lowest ZoneID tiebreaker.
	hottest := snapshot.Zones[0]
	for _, z := range snapshot.Zones[1:] {
		if z.TempC > hottest.TempC || (z.TempC == hottest.TempC && z.ZoneID < hottest.ZoneID) {
			hottest = z
		}
	}

	// Compute fill proportion and severity.
	ec := effectiveCritical(hottest, float64(pol.CritThreshold))
	fill := fillProportion(hottest.TempC, ec)
	sev := severity(hottest.TempC, float64(pol.WarnThreshold), ec)
	sevColor := severityColorRGBA(sev)

	// Format temperature label.
	tempLabel := formatTemp(hottest.TempC, pol.Unit)

	vd := style.ViewData{
		Items:  []string{tempLabel},
		Colors: []color.Color{sevColor},
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

	// Progress bar occupies remaining height below the label.
	barTop := oy + labelHeight + gap
	barHeight := height - (labelHeight + gap)
	if barHeight < 1 {
		barHeight = 1
	}

	// Center the bar horizontally with some padding.
	barWidth := width
	if barWidth > 40 {
		barWidth = width * 3 / 4 // Use 75% of width for a nicer look
	}
	barX := ox + (width-barWidth)/2

	barBG := color.RGBA{R: 40, G: 40, B: 40, A: 255}

	barResult := progressbar.Render(progressbar.Config{
		Style:       progressbar.Linear,
		Orientation: progressbar.OrientVertical,
		Value:       fill,
		Bounds:      image.Rect(barX, barTop, barX+barWidth, barTop+barHeight),
		Foreground:  sevColor,
		Background:  barBG,
		RoundedCaps: true,
	})
	if barResult != nil {
		barResult.Label = "portrait-thermometer-bar"
		vd.Sprites = append(vd.Sprites, *barResult)
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
