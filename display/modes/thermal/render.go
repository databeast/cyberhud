package thermal

import (
	"fmt"
	"image/color"
	"strings"

	"github.com/databeast/cyberhud/display/modes/thermal/source"
	"github.com/databeast/cyberhud/display/region"
	"github.com/databeast/cyberhud/display/style"

	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// buildView constructs a full style.ViewData for the thermal mode, using the current
// policy and the provided TextHints for adaptive layout.
func buildView(hints textlayout.TextHints) style.ViewData {
	pol := GetPolicy()
	snap := source.CurrentSnapshot()
	// Handle empty snapshot: if the sampling loop is active, show a brief
	// "Sampling..." message instead of the misleading "No thermal sensors found".
	// On a system that genuinely has no thermal zones, the loop will also produce
	// an empty snapshot, but we distinguish via loopState: if the loop is running,
	// we just haven't gotten data yet.
	if len(snap.Zones) == 0 {
		// If the loop is active, this is a transient startup state — sample inline.
		if fresh, ok := source.SampleActive(); ok {
			snap = fresh
			// Fall through to normal rendering below.
		} else {
			return style.ViewData{
				Items:  []string{"Sampling..."},
				Static: true,
			}
		}
	}

	// Dispatch to style-specific renderer.
	// Font selection is handled by the registry's tier-based resolution
	// (configured as spleen/TierNormal for thermal). We pass hints as-is.
	//
	// Adaptive layout: on space-constrained panels, auto-select minimal style
	// and suppress chrome (title/hint) to maximize useful content rows.
	totalRows := textlayout.MaxVisibleRows(hints, 0)

	if totalRows <= 4 {
		// constrained panels omit non-data chrome in style output
	}
	if totalRows <= 2 {
		// constrained panels omit non-data chrome in style output
	}
	if totalRows <= 3 && pol.Style != "minimal" {
		pol.Style = "minimal"
	}

	// Registry-based dispatch: configured → alias → fitness.
	s, reason := style.ResolveStyle(thermalRegistry, hints, "thermal", pol.Style)

	// Construct StyleContext for the style boundary.
	// Each style constructs its own bridge internally with its own PaddingPct.
	ctx := style.NewStyleContext(hints)
	vd := s.Build(snap, pol, ctx)

	// Report style resolution metadata to the registry layer.
	vd.StyleReport = style.StyleReport{
		Name:   s.Name(),
		Reason: reason,
	}

	// Populate constant fields — suppress title/hint on constrained panels.
	vd.Static = true

	if vd.OffsetY < 0 {
		vd.OffsetY = 0
	}

	return vd
}

// RenderCacheKey returns a string composed of formatted temperature values from
// the latest snapshot concatenated with a policy fingerprint encoding all 10
// policy fields, so the display registry can skip redundant redraws when neither
// data nor config changed. When any zone is critical, a blink-state token is
// appended that alternates each blinkTick. When show_led is true, an LED tick
// parity token is appended that alternates each ledTick.
//
// pixelWidth and pixelHeight encode the instance-level panel dimensions so that
// two instances on panels with different pixel dimensions produce distinct cache
// keys. Pass 0, 0 when instance hints are unavailable.
func RenderCacheKey(snap source.ThermalSnapshot, blinkTick int, ledTick int, elapsedMS int, pixelWidth, pixelHeight int) uint32 {
	pol := GetPolicy()

	var b strings.Builder

	// Zone data section: "label:temp;label:temp;..."
	for i, z := range snap.Zones {
		if i > 0 {
			b.WriteByte(';')
		}
		b.WriteString(fmt.Sprintf("%s:%.1f", z.Label, z.TempC))
	}

	// Policy fingerprint section: all 10 fields
	b.WriteByte('|')
	b.WriteString(fmt.Sprintf("style=%s,font=%s,refresh_ms=%d,warn=%d,crit=%d,unit=%s,accent=%s,led=%v,refresh_bar=%v,border=%v",
		pol.Style, pol.Font, pol.RefreshMS, pol.WarnThreshold, pol.CritThreshold, pol.Unit,
		pol.FGColor, pol.ShowLED, pol.ShowRefreshBar, pol.ShowBorder))

	// Blink token: appended when any zone is critical
	anyCritical := false
	for _, z := range snap.Zones {
		ec := effectiveCritical(z, float64(pol.CritThreshold))
		if z.TempC >= ec {
			anyCritical = true
			break
		}
	}
	if anyCritical {
		if blinkTick%2 == 0 {
			b.WriteString("|B0")
		} else {
			b.WriteString("|B1")
		}
	}

	// LED tick parity token: appended when show_led is true
	if pol.ShowLED {
		if ledTick%2 == 0 {
			b.WriteString("|L0")
		} else {
			b.WriteString("|L1")
		}
	}

	// Pixel dimensions section: encodes instance-level panel size so distinct
	// panels produce distinct cache keys.
	b.WriteString(fmt.Sprintf("|%dx%d", pixelWidth, pixelHeight))

	return region.CalcRegionCacheKey(b.String())
}

// toFahrenheit converts a Celsius temperature to Fahrenheit.
func toFahrenheit(celsius float64) float64 {
	return celsius*9.0/5.0 + 32.0
}

// formatTemp formats a temperature value with one decimal place and unit suffix.
// The celsius parameter is the temperature in Celsius; unit is "C" or "F".
func formatTemp(celsius float64, unit string) string {
	if unit == "F" {
		f := toFahrenheit(celsius)
		return fmt.Sprintf("%.1f°F", f)
	}
	return fmt.Sprintf("%.1f°C", celsius)
}

// severity returns the severity level for a temperature reading:
//
//	0 = normal (below warn threshold)
//	1 = warning (at or above warn, below effective critical)
//	2 = critical (at or above effective critical)
func severity(tempC, warnThreshold, effectiveCrit float64) int {
	if tempC >= effectiveCrit {
		return 2
	}
	if tempC >= warnThreshold {
		return 1
	}
	return 0
}

// effectiveCritical returns the effective critical threshold for a zone.
// It finds the lowest kernel "critical" trip point and returns min(configCrit, kernelCrit).
// If both are 0, returns 100.0 as a safe fallback.
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

// severityColor returns the color corresponding to a severity level:
//
//	0 → green (0, 255, 0)
//	1 → yellow (255, 255, 0)
//	2 → red (255, 0, 0)
func severityColor(level int) color.Color {
	switch level {
	case 1:
		return color.RGBA{R: 255, G: 255, B: 0, A: 255}
	case 2:
		return color.RGBA{R: 255, G: 0, B: 0, A: 255}
	default:
		return color.RGBA{R: 0, G: 255, B: 0, A: 255}
	}
}
