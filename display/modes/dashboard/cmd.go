package dashboard

import (
	"strings"

	"github.com/databeast/cyberhud/display/catalog"
	"github.com/databeast/cyberhud/display/catalog/cmdutil"
	"github.com/databeast/cyberhud/display/modes/dashboard/source"
)

// HandleCommand is the catalog command handler for the "dashboard" verb.
func HandleCommand(args []string) string {
	return cmdHandler.Handle(args)
}

var cmdHandler *cmdutil.CmdHandler

func init() {
	cmdHandler = &cmdutil.CmdHandler{
		Verb: "dashboard",
		Keys: []cmdutil.KeyDef{
			{
				Name:     "style",
				Validate: cmdutil.AllowedValidator(registeredStyleNames()),
			},
			{
				Name:     "color_accent",
				Validate: cmdutil.AllowedValidator(source.AllowedAccents),
			},
		},
		PostApply: fitnessNotesPostApply,
		Get: func(key string) string {
			policyMu.RLock()
			defer policyMu.RUnlock()
			switch key {
			case "style":
				return policy.Style
			case "color_accent":
				return policy.ColorAccent
			}
			return ""
		},
		Apply: func(key, value string) {
			policyMu.Lock()
			defer policyMu.Unlock()
			switch key {
			case "style":
				policy.Style = strings.ToLower(strings.TrimSpace(value))
			case "color_accent":
				policy.ColorAccent = strings.ToLower(strings.TrimSpace(value))
			}
		},
	}

	catalog.RegisterCommand(catalog.CommandDefinition{
		Verb:    "dashboard",
		Summary: "Query or set dashboard mode options",
		Usage:   "dashboard [style=...] [show_border=true|false] [color_accent=...]",
		Handle:  func(args []string) string { return cmdHandler.Handle(args) },
	})
}
