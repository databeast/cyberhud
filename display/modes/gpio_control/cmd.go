package gpio_control

import (
	"github.com/databeast/cyberhud/display/catalog"
	"github.com/databeast/cyberhud/display/catalog/cmdutil"
)

var cmdHandler = &cmdutil.CmdHandler{
	Verb: "gpio-control",
	Keys: []cmdutil.KeyDef{
		{Name: "style", Validate: cmdutil.AllowedValidator(registeredStyleNames())},
		{Name: "font", Validate: fontValidator},
	},
	PostApply: fitnessNotesPostApply,
	Get: func(key string) string {
		p := PolicySnapshot()
		switch key {
		case "style":
			return p.Style
		case "font":
			return p.Font
		}
		return ""
	},
	Apply: func(key, value string) {
		policyMu.Lock()
		defer policyMu.Unlock()
		switch key {
		case "style":
			policy.Style = value
		case "font":
			policy.Font = value
		}
	},
}

// HandleCommand is the catalog command handler for the "gpio-control" verb.
func HandleCommand(args []string) string {
	return cmdHandler.Handle(args)
}

func init() {
	catalog.Register(catalog.Definition{
		ID:      "gpio-control",
		Title:   "GPIO Control",
		Summary: "Interactive mode for toggling GPIO pins. Navigate with up/down, toggle with primary button.",
		Order:   42,
		Options: []catalog.OptionDefinition{
			{Key: "style", Type: "string", Summary: "Visual presentation style", Default: "", Allowed: registeredStyleNames()},
			{Key: "font", Type: "string", Summary: "Font selection (auto or a registered font ID).", Default: "auto"},
		},
	})

	catalog.RegisterCommand(catalog.CommandDefinition{
		Verb:    "gpio-control",
		Summary: "Query or set GPIO control mode configuration",
		Usage:   "gpio-control [style=<list|compact|grid>] [font=<auto|font-id>]",
		Handle:  HandleCommand,
	})
}
