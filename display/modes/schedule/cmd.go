package schedule

import (
	"strconv"
	"strings"

	"github.com/databeast/cyberhud/display/catalog"
	"github.com/databeast/cyberhud/display/catalog/cmdutil"
	"github.com/databeast/cyberhud/display/modes/schedule/source"
)

// handler is the package-level CmdHandler for the "schedule" verb.
var handler = newCmdHandler()

func newCmdHandler() *cmdutil.CmdHandler {
	return &cmdutil.CmdHandler{
		Verb: "schedule",
		Keys: []cmdutil.KeyDef{
			{Name: "source", Validate: sourceValidator},
			{Name: "refresh_sec", Validate: cmdutil.IntValidator(source.MinRefreshSec)},
			{Name: "time_format", Validate: cmdutil.AllowedValidator(source.AllowedTimeFormats)},
			{Name: "show_upcoming", Validate: cmdutil.BoolValidator()},
			{Name: "style", Validate: cmdutil.AllowedValidator(scheduleRegistry.RegisteredNames())},
			{Name: "font", Validate: cmdutil.FontValidator()},
		},
		Get:   getPolicyValue,
		Apply: applyPolicyValue,
	}
}

// HandleCommand is the catalog command handler for the "schedule" verb.
func HandleCommand(args []string) string {
	return handler.Handle(args)
}

func init() {
	catalog.Register(catalog.Definition{
		ID:       "schedule",
		Title:    "Schedule",
		Category: "Data",
		Scope:    "any",
		Summary:  "Shows the calendar event active right now from a local or remote iCalendar (.ics) source.",
		Order:    37,
		Options: append(source.Policy{}.Options(), []catalog.OptionDefinition{
			{Key: "style", Type: "string", Summary: "Visual presentation style.", Default: "", Allowed: scheduleRegistry.RegisteredNames()},
		}...),
	})
	catalog.RegisterCommand(catalog.CommandDefinition{
		Verb:    "schedule",
		Summary: "Query or set the schedule calendar source and display options.",
		Usage:   "schedule [source=<url|path>] [refresh_sec=<seconds>] [time_format=<24h|12h>] [show_upcoming=<true|false>] [style=<name>] [font=<auto|id>]",
		Handle:  HandleCommand,
	})
}

// sourceValidator accepts an empty string (disables fetching) or any
// non-whitespace value (a URL or path).
func sourceValidator(value string) string {
	v := strings.TrimSpace(value)
	if v == "" {
		return ""
	}
	if strings.ContainsAny(v, "\t\n\r") {
		return "must be a URL or path without whitespace control characters"
	}
	return ""
}

func getPolicyValue(key string) string {
	p := PolicySnapshot()
	switch key {
	case "source":
		return p.Source
	case "refresh_sec":
		return strconv.Itoa(p.RefreshSec)
	case "time_format":
		return p.TimeFormat
	case "show_upcoming":
		if p.ShowUpcoming {
			return "true"
		}
		return "false"
	case "style":
		return p.Style
	case "font":
		return p.Font
	default:
		return ""
	}
}

func applyPolicyValue(key, value string) {
	p := PolicySnapshot()
	switch key {
	case "source":
		p.Source = strings.TrimSpace(value)
	case "refresh_sec":
		if n, err := strconv.Atoi(strings.TrimSpace(value)); err == nil {
			p.RefreshSec = n
		}
	case "time_format":
		p.TimeFormat = strings.ToLower(strings.TrimSpace(value))
	case "show_upcoming":
		if v, ok := cmdutil.ParseBool(value); ok {
			p.ShowUpcoming = v
		}
	case "style":
		p.Style = strings.ToLower(strings.TrimSpace(value))
	case "font":
		p.Font = strings.TrimSpace(value)
	}
	SetPolicy(p)
}
