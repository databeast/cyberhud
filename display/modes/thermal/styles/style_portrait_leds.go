package styles

import (
	"fmt"
	"image"
	"image/color"

	"github.com/databeast/cyberhud/display/modes/thermal/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/style/layout"
	"github.com/databeast/cyberhud/display/widgets/led"
)

// buildPortraitLEDsStyle renders a vertical column of LED indicators, one per thermal
// zone, on portrait side panels.
// It is the shared BuildFn used by per-resolution styles that want the LEDs layout.
func buildPortraitLEDsStyle(snapshot source.ThermalSnapshot, pol source.Policy, ctx style.StyleContext, d def) style.ViewData {
	if len(snapshot.Zones) == 0 {
		return style.ViewData{Items: []string{"no thermal data"}}
	}

	hints := ctx.Hints()
	bridge := layout.NewLayoutBridge(hints, layout.BridgeConfig{PaddingPct: 5})

	ox, oy := bridge.ContentOrigin()
	width := bridge.AvailableContentWidth()
	height := bridge.AvailableContentHeight()

	// Zone truncation: each LED needs at least 16px of vertical space.
	maxZones := height / 16
	if maxZones < 1 {
		maxZones = 1
	}

	zones := snapshot.Zones
	if len(zones) > maxZones {
		zones = zones[:maxZones]
	}

	// Even spacing: distribute LEDs vertically across the content height.
	perLEDSpacing := height / len(zones)
	if perLEDSpacing < 1 {
		perLEDSpacing = 1
	}

	vd := style.ViewData{
		Items:  []string{"thermal"},
		Colors: []color.Color{color.RGBA{R: 255, G: 255, B: 255, A: 255}},
	}

	for i, z := range zones {
		ec := effectiveCritical(z, float64(pol.CritThreshold))
		sev := severity(z.TempC, float64(pol.WarnThreshold), ec)
		sevColor := severityColorRGBA(sev)

		// LED diameter: min(width, perLEDSpacing) - 4, minimum 8.
		// Additionally cap to ~77% of width so that the glow halo
		// (which adds ~30% of body radius on each side) stays within bounds.
		maxForGlow := width * 10 / 13
		diameter := width
		if perLEDSpacing < diameter {
			diameter = perLEDSpacing
		}
		diameter -= 4
		if diameter > maxForGlow {
			diameter = maxForGlow
		}
		if diameter < 8 {
			diameter = 8
		}

		// Center LED horizontally within the content area.
		ledX := ox + (width-diameter)/2
		if ledX < ox {
			ledX = ox
		}

		// Position LED vertically with even spacing, centered within its slot.
		slotY := oy + i*perLEDSpacing
		ledY := slotY + (perLEDSpacing-diameter)/2
		if ledY < slotY {
			ledY = slotY
		}

		ledResult := led.Render(led.Config{
			Shape:       led.Circle,
			State:       led.On,
			Brightness:  -1.0,
			Diameter:    diameter,
			Bounds:      image.Rectangle{Min: image.Pt(ledX, ledY)},
			Foreground:  sevColor,
			GlowEnabled: true,
		})
		if ledResult != nil {
			ledResult.Label = fmt.Sprintf("portrait-led-zone-%d", i)
			vd.Sprites = append(vd.Sprites, *ledResult)
		}
	}

	return vd
}
