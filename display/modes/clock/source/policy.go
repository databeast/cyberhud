package source

import (
	"fmt"

	"github.com/databeast/cyberhud/display/catalog"
)

// Allowed time and date format values.
var AllowedTimeFormats = []string{"24h", "12h"}
var AllowedDateFormats = []string{"YYYY-MM-DD", "DD-MM-YYYY", "MM-DD-YYYY", "none"}

// allowedSecondsBar lists the valid seconds_bar policy values.
var AllowedSecondsBar = []string{"none", "horizontal", "pie"}

// allowedBorderColors lists the valid border_color policy values.
var AllowedBorderColors = []string{"cyan", "green", "emerald", "amber", "red", "white", "none", "auto"}

// Policy captures all runtime-configurable parameters for the clock mode.
// All 13 fields are exposed through catalog registration and CLI command handling.
//
// Framework pattern demonstrated: policy definition — a typed struct that serves
// as the single source of truth for mode behavior, referenced by catalog defaults,
// normalization, CLI handlers, and the rendering pipeline.
type Policy struct {
	Style              string // one of 24 registered style names
	ShowSeconds        bool
	TimeFormat         string // "12h" or "24h"
	DateFormat         string // "YYYY-MM-DD", "DD-MM-YYYY", "MM-DD-YYYY", or "none"
	Timezone           string // IANA timezone string or "local"
	ShowWeekday        bool
	BlinkColon         bool
	FGColor            string // "cyan", "green", "amber", "red", "white", "none"
	ShowLED            bool   // LED seconds indicator
	SecondsBar         string // "none", "horizontal", "pie"
	ShowDaybar         bool   // Sparkline day-progress bar
	ShowBorder         bool   // Enable rounded border frame
	ShowBorderExplicit bool   // True when user explicitly set show_border (prevents auto-enablement)
	BorderColor        string // "cyan", "green", "amber", "red", "white", "none", "auto"
}

func (p Policy) Options() []catalog.OptionDefinition {
	return []catalog.OptionDefinition{
		{Key: "show_seconds", Type: "bool", Summary: "Include seconds in the time string.", Default: "true", Allowed: []string{"true", "false"}},
		{Key: "time_format", Type: "string", Summary: "12-hour or 24-hour time display.", Default: "24h", Allowed: AllowedTimeFormats},
		{Key: "date_format", Type: "string", Summary: "Date layout or none to hide date row.", Default: "YYYY-MM-DD", Allowed: AllowedDateFormats},
		{Key: "timezone", Type: "string", Summary: "IANA timezone identifier or local.", Default: "local", Allowed: []string{}},
		{Key: "show_weekday", Type: "bool", Summary: "Display the weekday name row.", Default: "true", Allowed: []string{"true", "false"}},
		{Key: "blink_colon", Type: "bool", Summary: "Animate colon on/off each second.", Default: "false", Allowed: []string{"true", "false"}},
		{Key: "fgcolor", Type: "string", Summary: "Foreground color for time text on color panels.", Default: "cyan", Allowed: AllowedFGColors},
		{Key: "show_led", Type: "bool", Summary: "Show LED seconds indicator when seconds digits hidden.", Default: "true", Allowed: []string{"true", "false"}},
		{Key: "seconds_bar", Type: "string", Summary: "Progress bar style showing seconds within current minute.", Default: "none", Allowed: AllowedSecondsBar},
		{Key: "show_daybar", Type: "bool", Summary: "Show sparkline bar indicating day progress.", Default: "false", Allowed: []string{"true", "false"}},
		{Key: "show_border", Type: "bool", Summary: "Show rounded border frame.", Default: "false", Allowed: []string{"true", "false"}},
		{Key: "border_color", Type: "string", Summary: "Border frame color or auto to inherit fgcolor.", Default: "auto", Allowed: AllowedBorderColors},
	}
}

// DefaultPolicy returns the default clock policy with all 13 fields initialized.
//
// Framework pattern demonstrated: policy definition — provides the canonical
// default state that catalog registration and normalization reference as the baseline.
func DefaultPolicy() Policy {
	return Policy{
		Style:              "",
		ShowSeconds:        true,
		TimeFormat:         "24h",
		DateFormat:         "YYYY-MM-DD",
		Timezone:           "local",
		ShowWeekday:        true,
		BlinkColon:         false,
		FGColor:            "cyan",
		ShowLED:            true,
		SecondsBar:         "none",
		ShowDaybar:         false,
		ShowBorder:         false,
		ShowBorderExplicit: false,
		BorderColor:        "auto",
	}
}

func (p Policy) Fingerprint() string {
	return policyFingerprint(p)
}

func (p Policy) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"style":        p.Style,
		"show_seconds": p.ShowSeconds,
		"time_format":  p.TimeFormat,
		"date_format":  p.DateFormat,
		"timezone":     p.Timezone,
		"show_weekday": p.ShowWeekday,
		"blink_colon":  p.BlinkColon,
		"fgcolor":      p.FGColor,
		"show_led":     p.ShowLED,
		"seconds_bar":  p.SecondsBar,
		"show_daybar":  p.ShowDaybar,
		"show_border":  p.ShowBorder,
		"border_color": p.BorderColor,
	}
}

// policyFingerprint encodes all policy fields into a short deterministic string.
// All 13 fields are included so that any policy change triggers a redraw.
func policyFingerprint(p Policy) string {
	return fmt.Sprintf("%s|%v|%s|%s|%s|%v|%v|%s|%v|%s|%v|%v|%v|%s",
		p.Style, p.ShowSeconds, p.TimeFormat, p.DateFormat,
		p.Timezone, p.ShowWeekday, p.BlinkColon,
		p.FGColor, p.ShowLED, p.SecondsBar, p.ShowDaybar,
		p.ShowBorder, p.ShowBorderExplicit, p.BorderColor)
}
