package menu

import (
	"strings"

	"github.com/databeast/cyberhud/display/catalog/cmdutil"
)

// cmdHandler is the shared command handler for the "menu" verb.
var cmdHandler = &cmdutil.CmdHandler{
	Verb: "menu",
	Keys: []cmdutil.KeyDef{
		{Name: "style", Validate: cmdutil.AllowedValidator(registeredStyleNames())},
	},
	PostApply: fitnessNotesPostApply,
	Get: func(key string) string {
		p := GetPolicy()
		switch strings.ToLower(key) {
		case "style":
			return p.Style
		}
		return ""
	},
	Apply: func(key, value string) {
		policyState.Lock()
		defer policyState.Unlock()
		switch strings.ToLower(key) {
		case "style":
			policyState.policy.Style = strings.ToLower(value)
		}
	},
}

// HandleCommand is the catalog command handler for the "menu" verb.
func HandleCommand(args []string) string {
	return cmdHandler.Handle(args)
}
