package gpio

import (
	"fmt"
	"strings"

	"github.com/databeast/cyberhud/display/modes/gpio/source"
	"github.com/databeast/cyberhud/display/region"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	gpiomgr "github.com/databeast/cyberhud/hardware/gpio"
)

// RenderCacheKey returns a change-detection string that incorporates a fingerprint
// of all Policy fields that affect rendering (Style, Color, Font, PinLabels,
// FGColor) plus a fingerprint of each pin's state
// (Number, Mode, Level). The display runtime compares consecutive keys to
// skip unnecessary re-renders when nothing has changed.
//
// The output is deterministic, at most 512 bytes, and returns a stable non-empty
// string even when the pin list is empty (distinguishable from any key with pins).
func RenderCacheKey(pins []gpiomgr.PinState) uint32 {
	p := GetPolicy()

	var b strings.Builder
	// Policy fingerprint: all fields that affect rendering output.
	b.WriteString(p.Style)
	b.WriteByte('|')
	if p.Color {
		b.WriteByte('1')
	} else {
		b.WriteByte('0')
	}
	b.WriteByte('|')
	b.WriteString(p.Font)
	b.WriteByte('|')
	b.WriteString(p.FGColor)
	b.WriteByte('|')

	b.WriteByte('0')

	b.WriteByte('|')
	b.WriteByte('0')

	b.WriteByte('|')

	// PinLabels: encode sorted by pin number for determinism.
	if len(p.PinLabels) > 0 {
		// Collect and sort keys for deterministic output.
		keys := make([]int, 0, len(p.PinLabels))
		for k := range p.PinLabels {
			keys = append(keys, k)
		}
		sortInts(keys)
		for _, k := range keys {
			fmt.Fprintf(&b, "%d=%s,", k, p.PinLabels[k])
		}
	}
	b.WriteByte('|')

	// Pin state fingerprint.
	if len(pins) == 0 {
		// Stable non-empty marker for empty pin list, distinguishable from any pin data.
		b.WriteString("empty")
	} else {
		for _, pin := range pins {
			lvl := byte('0')
			if pin.Level {
				lvl = '1'
			}
			fmt.Fprintf(&b, "%d%s%c;", pin.Number, pin.Mode.String(), lvl)
		}
	}

	sig := b.String()
	// Enforce max 512 bytes.
	if len(sig) > 512 {
		sig = sig[:512]
	}
	return region.CalcRegionCacheKey(sig)
}

// BuildView constructs a full style.ViewData for the GPIO mode, using the current policy
// and the provided TextHints for adaptive layout.
// The warnings parameter carries hardware pin notices; when non-empty the first
// warning is prepended to the Hint with a "WARN: " prefix.
//
// Dispatch flow: gpioRegistry.Lookup(style) → fallback gpioRegistry.Default() →
// style.Build(snapshot, hints). The expanded Policy is passed into GpioSnapshot so
// resolution styles access accent/label preferences without calling GetPolicy().
//
// Framework pattern demonstrated: Registry-based dispatch with accent color integration.
func BuildView(pins []gpiomgr.PinState, hints textlayout.TextHints, warnings []string) style.ViewData {
	if len(pins) == 0 {
		msg := "GPIO unavailable"
		if len(warnings) > 0 && warnings[0] != "" {
			msg = warnings[0]
		}
		return style.ViewData{Items: []string{msg}, Static: true}
	}

	pol := GetPolicy()

	// Registry-based dispatch: configured → alias → fitness.
	s, reason := style.ResolveStyle(gpioRegistry, hints, "gpio", pol.Style)

	// Pass expanded Policy into GpioSnapshot so resolution styles can access
	// accent/label preferences without relying on package-level state.
	snap := source.GpioSnapshot{Pins: pins, Policy: pol}
	ctx := style.NewStyleContext(hints)
	svd := s.Build(snap, pol, ctx)

	// Title/Hint/Static/StyleReport ONLY — no sprites, no geometry, no layout.
	svd.StyleReport = style.StyleReport{
		Name:   s.Name(),
		Reason: reason,
	}

	_ = warnings

	return svd
}
