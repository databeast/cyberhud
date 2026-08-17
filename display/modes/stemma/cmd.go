package stemma

import (
	"strings"

	"github.com/databeast/cyberhud/display/catalog"
	"github.com/databeast/cyberhud/display/catalog/cmdutil"
)

// cmdHandler is the shared CmdHandler for the "stemma" command verb.
var cmdHandler = &cmdutil.CmdHandler{
	Verb: "stemma",
	Keys: []cmdutil.KeyDef{
		{Name: "style", Validate: cmdutil.AllowedValidator(registryStyleNames())},
	},
	Get: func(key string) string {
		p := GetPolicy()
		switch key {
		case "style":
			return p.Style
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
		}
	},
	PostApply: fitnessNotesPostApply,
}

// HandleCommand is the catalog command handler for the "stemma" verb.
func HandleCommand(args []string) string {
	return cmdHandler.Handle(args)
}

func init() {
	catalog.Register(catalog.Definition{
		ID:      "stemma",
		Title:   "STEMMA",
		Scope:   "any",
		Summary: "Detected STEMMA QT / QWIIC device list or compact device summary.",
		Order:   30,
		Options: []catalog.OptionDefinition{
			{Key: "style", Type: "string", Summary: "Visual presentation style.", Default: "", Allowed: registryStyleNames()},
		},
	})
	catalog.RegisterCommand(catalog.CommandDefinition{
		Verb:    "stemma",
		Summary: "Query or set STEMMA display options.",
		Usage:   "stemma [style=<registered-style>]",
		Handle:  HandleCommand,
	})
}
