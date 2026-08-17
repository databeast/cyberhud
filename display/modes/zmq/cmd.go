package zmq

import (
	"strings"

	"github.com/databeast/cyberhud/display/catalog/cmdutil"
	"github.com/databeast/cyberhud/display/modes/zmq/content"
)

// cmdHandler is the package-level CmdHandler instance.
var cmdHandler = newCmdHandler()

// newCmdHandler creates the CmdHandler for the "zmq" command verb.
func newCmdHandler() *cmdutil.CmdHandler {
	return &cmdutil.CmdHandler{
		Verb: "zmq",
		Keys: []cmdutil.KeyDef{
			{Name: "endpoint", Validate: endpointValidator},
			{Name: "socket_type", Validate: cmdutil.AllowedValidator([]string{"sub", "pull"})},
			{Name: "topic", Validate: nil}, // accepts any string
			{Name: "max_lines", Validate: cmdutil.IntValidator(1)},
			{Name: "json_fields", Validate: nil}, // accepts any string
			{Name: "style", Validate: cmdutil.AllowedValidator(registeredStyleNames())},
			{Name: "font", Validate: fontValidator},
		},
		Get:       getZmqPolicyValue,
		Apply:     applyZmqPolicyValue,
		PostApply: nil, // SetPolicy already handles reconnection detection
	}
}

// HandleCommand processes the zmq command using CmdHandler.
// Handles the "clear" sub-command for buffer clearing before delegating
// to the standard CmdHandler for key=value operations.
func HandleCommand(args []string) string {
	if len(args) > 0 && strings.ToLower(args[0]) == "clear" {
		content.Clear()
		return "OK zmq cleared"
	}
	return cmdHandler.Handle(args)
}
