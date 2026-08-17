package styles

import (
	"image"
	"image/color"

	"github.com/databeast/cyberhud/display/modes/thermal/source"
	"github.com/databeast/cyberhud/display/widgets"
	"github.com/databeast/cyberhud/display/widgets/led"
)

// buildLEDSprite generates the LED activity indicator sprite based on current
// policy, snapshot data, and layout parameters.
//
// Rules:
//   - If !pol.ShowLED → return nil
//   - If effectiveWidth < 6 → return nil (suppression)
//   - LED state: On for even ticks (ledTick%2 == 0), Off for odd
//   - Position: top-right corner = (contentOffsetX + effectiveWidth - 6, contentOffsetY), diameter 6
//   - Suppress if LED 6×6 rect intersects any text row rect
//   - Foreground: highest-severity zone's color (lowest ZoneID tiebreaker) on color; native FG on mono
//   - Returns Sprite with label "led/on" or "led/off"
func buildLEDSprite(
	pol source.Policy,
	snap source.ThermalSnapshot,
	ledTick int,
	effectiveWidth int,
	effectiveHeight int,
	contentOffsetX int,
	contentOffsetY int,
	fontRowHeight int,
	numTextRows int,
	isColor bool,
	nativeFG color.RGBA,
) *widgets.Sprite {
	// Suppress if show_led is disabled.
	if !pol.ShowLED {
		return nil
	}

	// Suppress if effective width is too narrow for a 6px LED.
	if effectiveWidth < 6 {
		return nil
	}

	const diameter = 6

	// Determine LED state based on tick parity.
	state := led.On
	if ledTick%2 != 0 {
		state = led.Off
	}

	// Position: top-right corner of effective content area.
	ledX := contentOffsetX + effectiveWidth - diameter
	ledY := contentOffsetY

	// Check intersection with text rows.
	// LED bounding rect: (ledX, ledY) to (ledX+6, ledY+6).
	ledRect := image.Rect(ledX, ledY, ledX+diameter, ledY+diameter)
	for i := 0; i < numTextRows; i++ {
		rowY := contentOffsetY + i*fontRowHeight
		rowRect := image.Rect(contentOffsetX, rowY, contentOffsetX+effectiveWidth, rowY+fontRowHeight)
		if ledRect.Overlaps(rowRect) {
			return nil
		}
	}

	// Determine foreground color.
	fg := determineLEDForeground(pol, snap, isColor, nativeFG)

	// Render the LED widget.
	result := led.Render(led.Config{
		State:      state,
		Brightness: -1.0,
		Diameter:   diameter,
		Bounds:     image.Rectangle{Min: image.Pt(ledX, ledY)},
		Foreground: fg,
	})
	if result == nil {
		return nil
	}

	return result
}

// determineLEDForeground determines the LED foreground color.
// On color panels: uses the highest-severity zone's severity color (lowest ZoneID tiebreaker).
// On monochrome panels: uses the native foreground color.
func determineLEDForeground(pol source.Policy, snap source.ThermalSnapshot, isColor bool, nativeFG color.RGBA) color.RGBA {
	if !isColor {
		return nativeFG
	}

	// Find the highest severity across all zones; lowest ZoneID breaks ties.
	highestSeverity := 0
	for _, z := range snap.Zones {
		ec := effectiveCritical(z, float64(pol.CritThreshold))
		sev := severity(z.TempC, float64(pol.WarnThreshold), ec)
		if sev > highestSeverity {
			highestSeverity = sev
		}
	}

	return severityColorRGBA(highestSeverity)
}
