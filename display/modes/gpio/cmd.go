package gpio

import (
	"strings"

	"github.com/databeast/cyberhud/display/catalog"
	"github.com/databeast/cyberhud/display/catalog/cmdutil"
)

// cmdHandler is the shared CmdHandler for the "gpio" command verb.
var cmdHandler = &cmdutil.CmdHandler{
	Verb: "gpio",
	Keys: []cmdutil.KeyDef{
		{Name: "style", Validate: cmdutil.AllowedValidator(registeredStyleNames())},
		{Name: "color", Validate: cmdutil.BoolValidator()},
		{Name: "font", Validate: fontValidator},
		{Name: "fgcolor", Validate: cmdutil.AllowedValidator(allowedFGColors)},
	},
	PostApply: fitnessNotesPostApply,
	Get: func(key string) string {
		p := GetPolicy()
		switch key {
		case "style":
			return p.Style
		case "color":
			return boolStr(p.Color)
		case "font":
			return p.Font
		case "fgcolor":
			return p.FGColor
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
		case "color":
			if v, ok := cmdutil.ParseBool(value); ok {
				policyState.policy.Color = v
			}
		case "font":
			policyState.policy.Font = value
		case "fgcolor":
			policyState.policy.FGColor = strings.ToLower(strings.TrimSpace(value))
		}
	},
}

// HandleCommand is the catalog command handler for the "gpio" verb.
// It delegates to cmdHandler for per-key validation, atomic policy mutation,
// and formatted query responses.
//
// Framework pattern demonstrated: Command handling via cmdutil.CmdHandler
// with per-key validators and a PostApply fitness-notes hook.
func HandleCommand(args []string) string {
	return cmdHandler.Handle(args)
}

func init() {
	catalog.Register(catalog.Definition{
		ID:      "gpio",
		Title:   "GPIO",
		Scope:   "any",
		Summary: "GPIO pin state list on the main display or compact GPIO counters on secondary displays.",
		Order:   40,
		Options: []catalog.OptionDefinition{
			{Key: "style", Type: "string", Summary: "Visual presentation style.", Default: "", Allowed: registeredStyleNames()},
			{Key: "color", Type: "bool", Summary: "Whether to color-code pin rows by level.", Default: "true", Allowed: []string{"true", "false"}},
			{Key: "font", Type: "string", Summary: "Font selection (auto or a registered font ID).", Default: "auto"},
			{Key: "fgcolor", Type: "string", Summary: "Foreground color for GPIO display elements.", Default: "cyan", Allowed: allowedFGColors},
		},
	})
	catalog.RegisterCommand(catalog.CommandDefinition{
		Verb:    "gpio",
		Summary: "Query or set GPIO display options.",
		Usage:   "gpio [style=<list|icons|detail|dashboard|activity>] [color=<true|false>] [font=<auto|font-id>] [fgcolor=<cyan|green|amber|red|white|none>]",
		Handle:  HandleCommand,
	})
}
