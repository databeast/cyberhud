package source

import (
	"image/color"
	"strings"
	"time"

	sharedcolor "github.com/databeast/cyberhud/display/style/color"
)

// applyTimezone converts now to the specified IANA timezone.
// Returns now unchanged for "", "local", or if the timezone cannot be loaded.
func ApplyTimezone(now time.Time, tz string) time.Time {
	if tz == "" || tz == "local" {
		return now
	}
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return now
	}
	return now.In(loc)
}

// formatDate applies the configured timezone and returns the date in the
// configured format. Returns "" when DateFormat is "none".
func FormatDate(now time.Time, p Policy) string {
	t := ApplyTimezone(now, p.Timezone)
	switch p.DateFormat {
	case "DD-MM-YYYY":
		return t.Format("02-01-2006")
	case "MM-DD-YYYY":
		return t.Format("01-02-2006")
	case "none":
		return ""
	default: // "YYYY-MM-DD"
		return t.Format("2006-01-02")
	}
}

// formatTime formats the current time according to the policy's TimeFormat,
// ShowSeconds, Timezone, and BlinkColon settings.
func FormatTime(now time.Time, p Policy) string {
	t := ApplyTimezone(now, p.Timezone)
	sec := t.Second()

	var timeStr string
	switch p.TimeFormat {
	case "12h":
		if p.ShowSeconds {
			timeStr = t.Format("3:04:05 PM")
		} else {
			timeStr = t.Format("3:04 PM")
		}
	default: // "24h"
		if p.ShowSeconds {
			timeStr = t.Format("15:04:05")
		} else {
			timeStr = t.Format("15:04")
		}
	}

	if p.BlinkColon && sec%2 == 1 {
		timeStr = replaceColons(timeStr, ' ')
	}

	return timeStr
}

// replaceColons replaces all colon characters in s with the given rune.
// In Go's time format output, colons only appear as HH:MM and MM:SS separators,
// so replacing all ':' characters is safe.
func replaceColons(s string, replacement rune) string {
	return strings.ReplaceAll(s, ":", string(replacement))
}

// allowedFGColors lists all valid fgcolor policy values.
var AllowedFGColors = []string{"cyan", "green", "amber", "red", "white", "none"}

// resolveFGColor returns the primary RGBA for a named accent.
// Unknown or "none" values resolve to opaque white.
func resolveFGColor(name string) color.RGBA {
	return sharedcolor.ResolveAccent(name)
}

// TimezoneValidator validates a timezone string. Accepts "local" as a special
// case; otherwise requires a valid IANA timezone identifier.
func TimezoneValidator(value string) string {
	v := strings.TrimSpace(value)
	if strings.ToLower(v) == "local" {
		return ""
	}
	_, err := time.LoadLocation(v)
	if err != nil {
		return "must be a valid IANA timezone identifier or \"local\""
	}
	return ""
}

// BuildClockData formats current time fields for panel, respecting the current policy.
//
// Framework pattern demonstrated: data-signature change detection — produces the
// formatted time/date/weekday strings that BuildView consumes, with timezone
// application and conditional weekday inclusion driven by the active policy.
func BuildData(p Policy) ClockData {
	now := time.Now()
	t := ApplyTimezone(now, p.Timezone)

	var weekday string
	if p.ShowWeekday {
		weekday = t.Weekday().String()
	}

	return ClockData{
		Time:    FormatTime(now, p),
		Date:    FormatDate(now, p),
		Weekday: weekday,
		Now:     now,
	}
}
