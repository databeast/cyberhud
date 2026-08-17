package styles

import (
	"image/color"

	"github.com/databeast/cyberhud/display/modes/thermal/source"
	"github.com/databeast/cyberhud/display/widgets"
)

func ToFahrenheit(c float64) float64 { return toFahrenheit(c) }
func Severity(tempC, warnThreshold, effectiveCrit float64) int {
	return severity(tempC, warnThreshold, effectiveCrit)
}
func AccentColor(name string) color.RGBA { return accentColor(name) }
func ResolveAccentForLabel(accent string, isColor bool, severityLevel int, nativeFG color.RGBA) color.RGBA {
	return resolveAccentForLabel(accent, isColor, severityLevel, nativeFG)
}
func ResolveAccentForTemp(accent string, isColor bool, severityLevel int, nativeFG color.RGBA) color.RGBA {
	return resolveAccentForTemp(accent, isColor, severityLevel, nativeFG)
}
func ResolveAccentForBar(accent string, isColor bool, severityLevel int, nativeFG color.RGBA) color.RGBA {
	return resolveAccentForBar(accent, isColor, severityLevel, nativeFG)
}
func SeverityColorRGBA(level int) color.RGBA { return severityColorRGBA(level) }
func EffectiveCritical(zone source.ZoneReading, configCrit float64) float64 {
	return effectiveCritical(zone, configCrit)
}
func ApplyPanelDefaults(p *source.Policy, isColor bool, width, height int, explicit map[string]bool) {
	applyPanelDefaults(p, isColor, width, height, explicit)
}
func BuildLEDSprite(pol source.Policy, snap source.ThermalSnapshot, ledTick, effectiveWidth, effectiveHeight, contentOffsetX, contentOffsetY, fontRowHeight, numTextRows int, isColor bool, nativeFG color.RGBA) *widgets.Sprite {
	return buildLEDSprite(pol, snap, ledTick, effectiveWidth, effectiveHeight, contentOffsetX, contentOffsetY, fontRowHeight, numTextRows, isColor, nativeFG)
}
func DetermineLEDForeground(pol source.Policy, snap source.ThermalSnapshot, isColor bool, nativeFG color.RGBA) color.RGBA {
	return determineLEDForeground(pol, snap, isColor, nativeFG)
}
