package styles

import (
	"fmt"
	"image/color"
	"strings"

	"github.com/databeast/cyberhud/display/modes/thermal/source"
)

func allEmptyItems(items []string) bool {
	if len(items) == 0 {
		return true
	}
	for _, item := range items {
		if item != "" {
			return false
		}
	}
	return true
}

func toFahrenheit(celsius float64) float64 {
	return celsius*9.0/5.0 + 32.0
}

func formatTemp(celsius float64, unit string) string {
	if unit == "F" {
		return fmt.Sprintf("%.1f°F", toFahrenheit(celsius))
	}
	return fmt.Sprintf("%.1f°C", celsius)
}

func severity(tempC, warnThreshold, effectiveCrit float64) int {
	if tempC >= effectiveCrit {
		return 2
	}
	if tempC >= warnThreshold {
		return 1
	}
	return 0
}

func effectiveCritical(zone source.ZoneReading, configCrit float64) float64 {
	kernelCrit := 0.0
	found := false
	for _, tp := range zone.TripPoints {
		if strings.ToLower(tp.Type) == "critical" {
			if !found || tp.TempC < kernelCrit {
				kernelCrit = tp.TempC
				found = true
			}
		}
	}
	if found && kernelCrit < configCrit {
		if kernelCrit == 0 && configCrit == 0 {
			return 100.0
		}
		return kernelCrit
	}
	if configCrit == 0 && !found {
		return 100.0
	}
	if configCrit == 0 && found {
		return kernelCrit
	}
	return configCrit
}

func severityColor(level int) color.Color {
	return severityColorRGBA(level)
}

func severityColorRGBA(level int) color.RGBA {
	switch level {
	case 1:
		return color.RGBA{R: 255, G: 255, B: 0, A: 255}
	case 2:
		return color.RGBA{R: 255, G: 0, B: 0, A: 255}
	default:
		return color.RGBA{R: 0, G: 255, B: 0, A: 255}
	}
}

func accentColor(name string) color.RGBA {
	switch name {
	case "cyan":
		return color.RGBA{R: 0, G: 255, B: 255, A: 255}
	case "green":
		return color.RGBA{R: 0, G: 200, B: 0, A: 255}
	case "amber":
		return color.RGBA{R: 255, G: 191, B: 0, A: 255}
	case "red":
		return color.RGBA{R: 255, G: 0, B: 0, A: 255}
	case "white":
		return color.RGBA{R: 255, G: 255, B: 255, A: 255}
	default:
		return color.RGBA{R: 255, G: 255, B: 255, A: 255}
	}
}

func resolveAccentForLabel(accent string, isColor bool, severityLevel int, nativeFG color.RGBA) color.RGBA {
	if !isColor {
		return nativeFG
	}
	switch accent {
	case "thermal":
		return severityColorRGBA(severityLevel)
	case "none":
		return color.RGBA{R: 255, G: 255, B: 255, A: 255}
	default:
		return accentColor(accent)
	}
}

func resolveAccentForTemp(accent string, isColor bool, severityLevel int, nativeFG color.RGBA) color.RGBA {
	if !isColor {
		return nativeFG
	}
	if accent == "none" {
		return color.RGBA{R: 255, G: 255, B: 255, A: 255}
	}
	return severityColorRGBA(severityLevel)
}

func resolveAccentForBar(accent string, isColor bool, severityLevel int, nativeFG color.RGBA) color.RGBA {
	if !isColor {
		return nativeFG
	}
	return severityColorRGBA(severityLevel)
}
