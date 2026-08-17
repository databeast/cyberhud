package thermal

import (
	"image/color"

	"github.com/databeast/cyberhud/display/modes/thermal/source"
	"github.com/databeast/cyberhud/display/modes/thermal/styles"
	"github.com/databeast/cyberhud/display/widgets"
	"github.com/databeast/cyberhud/runtime/action"
)

type ZoneHistory = source.ZoneHistory

func NormalizePolicy(p Policy) (Policy, error) { return normalizePolicy(p) }
func NewZoneHistory() *ZoneHistory    { return &source.ZoneHistory{} }
func ToFahrenheit(c float64) float64  { return styles.ToFahrenheit(c) }
func Severity(tempC, warnThreshold, effectiveCrit float64) int {
	return styles.Severity(tempC, warnThreshold, effectiveCrit)
}
func AccentColor(name string) color.RGBA { return styles.AccentColor(name) }
func ResolveAccentForLabel(accent string, isColor bool, severityLevel int, nativeFG color.RGBA) color.RGBA {
	return styles.ResolveAccentForLabel(accent, isColor, severityLevel, nativeFG)
}
func ResolveAccentForTemp(accent string, isColor bool, severityLevel int, nativeFG color.RGBA) color.RGBA {
	return styles.ResolveAccentForTemp(accent, isColor, severityLevel, nativeFG)
}
func ResolveAccentForBar(accent string, isColor bool, severityLevel int, nativeFG color.RGBA) color.RGBA {
	return styles.ResolveAccentForBar(accent, isColor, severityLevel, nativeFG)
}
func SeverityColorRGBA(level int) color.RGBA { return styles.SeverityColorRGBA(level) }
func EffectiveCritical(zone source.ZoneReading, configCrit float64) float64 {
	return styles.EffectiveCritical(zone, configCrit)
}
func ApplyPanelDefaults(p *Policy, isColor bool, width, height int, explicit map[string]bool) {
	styles.ApplyPanelDefaults((*source.Policy)(p), isColor, width, height, explicit)
}
func AllowedStyles() []string { return allowedStyleNames() }

var ThermalRegistry = thermalRegistry

func ThermalRegistryEnumerate() []interface{} {
	out := make([]interface{}, 0, len(thermalRegistry.Enumerate()))
	for _, s := range thermalRegistry.Enumerate() {
		out = append(out, s)
	}
	return out
}

var AllowedFGColors = source.AllowedFGColors
var AllowedUnits = source.AllowedUnits

func CycleStyle(delta int) action.Result { return cycleStyle(delta) }
func BuildLEDSprite(pol Policy, snap ThermalSnapshot, ledTick, effectiveWidth, effectiveHeight, contentOffsetX, contentOffsetY, fontRowHeight, numTextRows int, isColor bool, nativeFG color.RGBA) *widgets.Sprite {
	return styles.BuildLEDSprite(pol, snap, ledTick, effectiveWidth, effectiveHeight, contentOffsetX, contentOffsetY, fontRowHeight, numTextRows, isColor, nativeFG)
}
func DetermineLEDForeground(pol Policy, snap ThermalSnapshot, isColor bool, nativeFG color.RGBA) color.RGBA {
	return styles.DetermineLEDForeground(pol, snap, isColor, nativeFG)
}
func GetHistory(zoneID int) []float64         { return source.GetHistory(zoneID) }
func RecordHistory(zoneID int, tempC float64) { source.RecordHistory(zoneID, tempC) }
func UpdateSnapshot(snap ThermalSnapshot)     { source.UpdateSnapshot(snap) }
func ResetHistoryState()                      { source.ResetHistoryStateForTest() }
func ResetSnapshotState()                     { source.UpdateSnapshot(source.ThermalSnapshot{}) }
