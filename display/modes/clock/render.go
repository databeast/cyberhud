package clock

import (
	"time"

	"github.com/databeast/cyberhud/display/modes/clock/source"
	"github.com/databeast/cyberhud/display/region"
	"github.com/databeast/cyberhud/display/style"

	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// BuildView returns the clock mode view data as a complete style.ViewData,
// including Title, Hint, and Static fields. Dispatch is performed through
// the clockRegistry: lookup the active style by name, fall back to the
// default style if unregistered, and invoke Build.
//
// All visual output (border, LED, progress bar, sparkline, text) is produced
// by the style's Build() method. BuildView only handles style selection and
// chrome metadata.
func BuildView(now time.Time, hints textlayout.TextHints) style.ViewData {
	p := GetPolicy()
	data := BuildSnapshot(now)

	// Registry-based dispatch: configured → alias → fitness.
	s, reason := style.ResolveStyle(clockRegistry, hints, "clock", p.Style)

	ctx := style.NewStyleContext(hints)
	vd := s.Build(data, p, ctx)

	vd.StyleReport = style.StyleReport{Name: s.Name(), Reason: reason}
	vd.Static = true

	return vd
}

func BuildSnapshot(now time.Time) source.ClockData {
	p := GetPolicy()
	t := source.ApplyTimezone(now, p.Timezone)
	data := source.BuildData(p)
	data.Now = now
	data.Time = source.FormatTime(now, p)
	data.Date = source.FormatDate(now, p)
	if p.ShowWeekday {
		data.Weekday = t.Weekday().String()
	} else {
		data.Weekday = ""
	}
	return data
}

// RenderCacheKey returns a change-detection string that incorporates the current
// time at second or minute precision (depending on active second-precision features)
// plus a fingerprint of all 12 policy fields. The display runtime compares
// consecutive keys to skip unnecessary re-renders.
//
// Precision-adaptive fingerprinting varies granularity with the set of active
// features (ShowSeconds, BlinkColon, ShowLED, SecondsBar) for battery efficiency.
func RenderCacheKey(now time.Time) uint32 {
	p := GetPolicy()
	t := source.ApplyTimezone(now, p.Timezone)

	needsSecondPrecision := p.ShowSeconds || p.BlinkColon || p.ShowLED || p.SecondsBar != "none"

	var timeComponent string
	if needsSecondPrecision {
		timeComponent = t.Format("20060102150405")
	} else {
		timeComponent = t.Format("200601021504")
	}

	return region.CalcRegionCacheKey(timeComponent)
}
