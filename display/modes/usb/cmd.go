package usb

import (
	"strconv"
	"strings"

	"github.com/databeast/cyberhud/display/catalog/cmdutil"
)

// cmdHandler is the shared CmdHandler for the "usb" command verb.
var cmdHandler = &cmdutil.CmdHandler{
	Verb: "usb",
	Keys: []cmdutil.KeyDef{
		{Name: "poll_ms", Validate: cmdutil.IntValidator(100)},
		{Name: "hold_unplugged_ms", Validate: cmdutil.IntValidator(0)},
		{Name: "hide_root_hubs", Validate: cmdutil.BoolValidator()},
		{Name: "style", Validate: cmdutil.AllowedValidator(registeredStyleNames())},
	},
	PostApply: fitnessNotesPostApply,
	Get: func(key string) string {
		p := PolicySnapshot()
		switch key {
		case "poll_ms":
			return strconv.Itoa(p.PollMS)
		case "hold_unplugged_ms":
			return strconv.Itoa(p.HoldUnpluggedMS)
		case "hide_root_hubs":
			return boolStr(p.HideRootHubs)
		case "style":
			return p.Style
		default:
			return ""
		}
	},
	Apply: func(key, value string) {
		monitorState.Lock()
		defer monitorState.Unlock()
		switch key {
		case "poll_ms":
			if n, err := strconv.Atoi(strings.TrimSpace(value)); err == nil {
				monitorState.policy.PollMS = n
			}
		case "hold_unplugged_ms":
			if n, err := strconv.Atoi(strings.TrimSpace(value)); err == nil {
				monitorState.policy.HoldUnpluggedMS = n
			}
		case "hide_root_hubs":
			if v, ok := cmdutil.ParseBool(value); ok {
				monitorState.policy.HideRootHubs = v
			}
		case "style":
			monitorState.policy.Style = strings.ToLower(strings.TrimSpace(value))
		}
	},
}

// HandleConsoleCommand handles the top-level "usb" console verb.
// It delegates to the CmdHandler for uniform key=value processing.
func HandleConsoleCommand(args []string) string {
	return cmdHandler.Handle(args)
}
