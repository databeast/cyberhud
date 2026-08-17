package wifi

import (
	"strings"

	"github.com/databeast/cyberhud/display/catalog/cmdutil"
	"github.com/databeast/cyberhud/display/modes/wifi/source"
)

// cmdHandler is the shared CmdHandler for the "wifi" command verb.
var cmdHandler = &cmdutil.CmdHandler{
	Verb: "wifi",
	Keys: []cmdutil.KeyDef{
		{Name: "style", Validate: cmdutil.AllowedValidator(registeredStyleNames())},
		{Name: "fgcolor", Validate: cmdutil.AllowedValidator(source.AllowedFGColors)},
		{Name: "signal_display", Validate: cmdutil.AllowedValidator(source.AllowedSignalDisplay)},
		{Name: "show_frequency", Validate: cmdutil.BoolValidator()},
		{Name: "show_interface", Validate: cmdutil.BoolValidator()},
		{Name: "show_channel", Validate: cmdutil.BoolValidator()},
	},
	Get: func(key string) string {
		p := GetPolicy()
		switch key {
		case "style":
			return p.Style
		case "fgcolor":
			return p.FGColor
		case "signal_display":
			return p.SignalDisplay
		case "show_frequency":
			return boolStr(p.ShowFrequency)
		case "show_interface":
			return boolStr(p.ShowInterface)
		case "show_channel":
			return boolStr(p.ShowChannel)
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
		case "fgcolor":
			policyState.policy.FGColor = strings.ToLower(strings.TrimSpace(value))
		case "signal_display":
			policyState.policy.SignalDisplay = strings.ToLower(strings.TrimSpace(value))
		case "show_frequency":
			if v, ok := cmdutil.ParseBool(value); ok {
				policyState.policy.ShowFrequency = v
			}
		case "show_interface":
			if v, ok := cmdutil.ParseBool(value); ok {
				policyState.policy.ShowInterface = v
			}
		case "show_channel":
			if v, ok := cmdutil.ParseBool(value); ok {
				policyState.policy.ShowChannel = v
			}
		}
	},
	PostApply: fitnessNotesPostApply,
}

// HandleCommand is the catalog command handler for the "wifi" verb.
// It delegates to the shared CmdHandler which validates inputs, rejects invalid
// values with descriptive error messages, and applies accepted values atomically.
func HandleCommand(args []string) string {
	return cmdHandler.Handle(args)
}
