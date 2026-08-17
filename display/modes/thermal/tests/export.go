package tests

import "github.com/databeast/cyberhud/display/modes/thermal"

var NormalizePolicy = thermal.NormalizePolicy

type ZoneHistoryExported = thermal.ZoneHistory

func NewZoneHistory() *ZoneHistoryExported { return thermal.NewZoneHistory() }

var ToFahrenheit = thermal.ToFahrenheit
var Severity = thermal.Severity
var AccentColor = thermal.AccentColor
var ResolveAccentForLabel = thermal.ResolveAccentForLabel
var ResolveAccentForTemp = thermal.ResolveAccentForTemp
var ResolveAccentForBar = thermal.ResolveAccentForBar
var SeverityColorRGBA = thermal.SeverityColorRGBA
var EffectiveCritical = thermal.EffectiveCritical
var ApplyPanelDefaults = thermal.ApplyPanelDefaults
var AllowedStyles = thermal.AllowedStyles()
var ThermalRegistry = thermal.ThermalRegistry
var AllowedFGColors = thermal.AllowedFGColors
var AllowedUnits = thermal.AllowedUnits
var CycleStyle = thermal.CycleStyle
var BuildLEDSprite = thermal.BuildLEDSprite
var DetermineLEDForeground = thermal.DetermineLEDForeground
var GetHistory = thermal.GetHistory

func ResetHistoryState()  { thermal.ResetHistoryState() }
func ResetSnapshotState() { thermal.ResetSnapshotState() }
