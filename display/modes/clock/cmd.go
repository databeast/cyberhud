package clock

import (
	"strings"

	"github.com/databeast/cyberhud/display/catalog"
	"github.com/databeast/cyberhud/display/catalog/cmdutil"
	"github.com/databeast/cyberhud/display/modes/clock/source"
)

// boolStr returns "true" or "false" for a bool value.
func boolStr(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

// Handler implements action.Handler for the clock mode.
//
// Framework pattern demonstrated: command handling — input action dispatch with
// policy mutation under write lock and dirty signaling for immediate redraw.
type Handler struct{}

// cmdHandler is the shared CmdHandler for the "clock" command verb.
// Demonstrates the framework pattern: command handling with validation and atomic policy mutation.
var cmdHandler = &cmdutil.CmdHandler{
	Verb: "clock",
	Keys: []cmdutil.KeyDef{
		{Name: "style", Validate: cmdutil.AllowedValidator(registeredStyleNames())},
		{Name: "show_seconds", Validate: cmdutil.BoolValidator()},
		{Name: "time_format", Validate: cmdutil.AllowedValidator(source.AllowedTimeFormats)},
		{Name: "date_format", Validate: cmdutil.AllowedValidator(source.AllowedDateFormats)},
		{Name: "timezone", Validate: source.TimezoneValidator},
		{Name: "show_weekday", Validate: cmdutil.BoolValidator()},
		{Name: "blink_colon", Validate: cmdutil.BoolValidator()},
		{Name: "fgcolor", Validate: cmdutil.AllowedValidator(source.AllowedFGColors)},
		{Name: "show_led", Validate: cmdutil.BoolValidator()},
		{Name: "seconds_bar", Validate: cmdutil.AllowedValidator(source.AllowedSecondsBar)},
		{Name: "show_daybar", Validate: cmdutil.BoolValidator()},
		{Name: "show_border", Validate: cmdutil.BoolValidator()},
		{Name: "border_color", Validate: cmdutil.AllowedValidator(source.AllowedBorderColors)},
	},
	Get: func(key string) string {
		p := GetPolicy()
		switch key {
		case "style":
			return p.Style
		case "show_seconds":
			return boolStr(p.ShowSeconds)
		case "time_format":
			return p.TimeFormat
		case "date_format":
			return p.DateFormat
		case "timezone":
			return p.Timezone
		case "show_weekday":
			return boolStr(p.ShowWeekday)
		case "blink_colon":
			return boolStr(p.BlinkColon)
		case "fgcolor":
			return p.FGColor
		case "show_led":
			return boolStr(p.ShowLED)
		case "seconds_bar":
			return p.SecondsBar
		case "show_daybar":
			return boolStr(p.ShowDaybar)
		case "show_border":
			return boolStr(p.ShowBorder)
		case "border_color":
			return p.BorderColor
		default:
			return ""
		}
	},
	Apply: func(key, value string) {
		policyState.Lock()
		defer policyState.Unlock()
		switch key {
		case "style":
			policyState.policy.Style = strings.ToLower(strings.TrimSpace(value))
		case "show_seconds":
			if v, ok := cmdutil.ParseBool(value); ok {
				policyState.policy.ShowSeconds = v
			}
		case "time_format":
			policyState.policy.TimeFormat = strings.TrimSpace(value)
		case "date_format":
			policyState.policy.DateFormat = strings.TrimSpace(value)
		case "timezone":
			v := strings.TrimSpace(value)
			if strings.ToLower(v) == "local" {
				policyState.policy.Timezone = "local"
			} else {
				policyState.policy.Timezone = v
			}
		case "show_weekday":
			if v, ok := cmdutil.ParseBool(value); ok {
				policyState.policy.ShowWeekday = v
			}
		case "blink_colon":
			if v, ok := cmdutil.ParseBool(value); ok {
				policyState.policy.BlinkColon = v
			}
		case "fgcolor":
			policyState.policy.FGColor = strings.ToLower(strings.TrimSpace(value))
		case "show_led":
			if v, ok := cmdutil.ParseBool(value); ok {
				policyState.policy.ShowLED = v
			}
		case "seconds_bar":
			policyState.policy.SecondsBar = strings.ToLower(strings.TrimSpace(value))
		case "show_daybar":
			if v, ok := cmdutil.ParseBool(value); ok {
				policyState.policy.ShowDaybar = v
			}
		case "show_border":
			if v, ok := cmdutil.ParseBool(value); ok {
				policyState.policy.ShowBorder = v
				policyState.policy.ShowBorderExplicit = true
			}
		case "border_color":
			policyState.policy.BorderColor = strings.ToLower(strings.TrimSpace(value))
		}
	},
	PostApply: fitnessNotesPostApply,
}

// HandleCommand is the catalog command handler for the "clock" verb.
// It delegates to the shared CmdHandler which validates inputs, rejects invalid
// values with descriptive error messages, and applies accepted values atomically.
//
// Framework pattern demonstrated: command handling — verb-based CLI dispatch with
// per-key validation and atomic policy mutation via cmdutil.CmdHandler.
func HandleCommand(args []string) string {
	return cmdHandler.Handle(args)
}

func init() {
	catalog.Register(catalog.Definition{
		ID:      "clock",
		Title:   "Clock",
		Summary: "Current time, date, and weekday sized for secondary displays.",
		Order:   60,
		Options: append(source.Policy{}.Options(), []catalog.OptionDefinition{
			{Key: "style", Type: "string", Summary: "Visual layout variant for time display.", Default: "", Allowed: registeredStyleNames()},
		}...),
	})
	catalog.RegisterCommand(catalog.CommandDefinition{
		Verb:    "clock",
		Summary: "Query or set clock display options.",
		Usage:   "clock [style=<name>] [show_seconds=<true|false>] [time_format=<24h|12h>] [date_format=<YYYY-MM-DD|DD-MM-YYYY|MM-DD-YYYY|none>] [timezone=<IANA|local>] [show_weekday=<true|false>] [blink_colon=<true|false>] [show_border=<true|false>] [fgcolor=<cyan|green|amber|red|white|none>] [show_led=<true|false>] [seconds_bar=<none|horizontal|pie>] [show_daybar=<true|false>]",
		Handle:  HandleCommand,
	})
}
