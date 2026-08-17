package systemd

import (
	"strings"

	"github.com/databeast/cyberhud/display/catalog/cmdutil"
	"github.com/databeast/cyberhud/display/modes/systemd/source"
)

// cmdHandler is the shared CmdHandler for the "systemd" command verb.
var cmdHandler = &cmdutil.CmdHandler{
	Verb: "systemd",
	Keys: []cmdutil.KeyDef{
		{Name: "style", Validate: cmdutil.AllowedValidator(registeredStyleNames())},
		{Name: "color_accent", Validate: cmdutil.AllowedValidator(source.AllowedAccents)},
	},
	PostApply: fitnessNotesPostApply,
	Get: func(key string) string {
		p := GetPolicy()
		switch key {
		case "style":
			return p.Style
		case "color_accent":
			return p.ColorAccent
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
		case "color_accent":
			policyState.policy.ColorAccent = strings.ToLower(strings.TrimSpace(value))
		}
	},
}

// HandleCommand is the catalog command handler for the "systemd" verb.
func HandleCommand(args []string) string {
	return cmdHandler.Handle(args)
}
