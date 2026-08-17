package serial

import (
	"strconv"
	"strings"

	"github.com/databeast/cyberhud/display/catalog"
	"github.com/databeast/cyberhud/display/catalog/cmdutil"
	"github.com/databeast/cyberhud/display/modes/serial/source"
)

// handler is the package-level CmdHandler instance.
var handler = newCmdHandler()

// newCmdHandler creates the CmdHandler for the "serial" command verb.
func newCmdHandler() *cmdutil.CmdHandler {
	return &cmdutil.CmdHandler{
		Verb: "serial",
		Keys: []cmdutil.KeyDef{
			{Name: "port", Validate: portValidator},
			{Name: "baud", Validate: cmdutil.IntValidator(1)},
			{Name: "lines", Validate: cmdutil.IntValidator(1)},
			{Name: "autoselect", Validate: cmdutil.BoolValidator()},
			{Name: "scan_ms", Validate: cmdutil.IntValidator(1)},
			{Name: "style", Validate: cmdutil.AllowedValidator(registeredStyleNames())},
			{Name: "font", Validate: fontValidator},
		},
		PostApply: fitnessNotesPostApply,
		Get:       getPolicyValue,
		Apply:     applyPolicyValue,
	}
}

func init() {
	catalog.Register(catalog.Definition{
		ID:      "serial",
		Title:   "Serial Monitor",
		Scope:   "any",
		Summary: "Live serial console monitor with auto-select and configurable port/baud.",
		Order:   36,
		Options: []catalog.OptionDefinition{
			{Key: "port", Type: "string", Summary: "Manual serial device path (disables auto-select when set).", Default: ""},
			{Key: "baud", Type: "int", Summary: "Serial baud rate.", Default: "115200"},
			{Key: "lines", Type: "int", Summary: "How many recent output lines to keep.", Default: "24"},
			{Key: "autoselect", Type: "bool", Summary: "Auto-pick the best available USB serial port.", Default: "true", Allowed: []string{"true", "false"}},
			{Key: "scan_ms", Type: "int", Summary: "Milliseconds between reconnect attempts.", Default: "500"},
			{Key: "style", Type: "string", Summary: "Visual presentation style.", Default: "", Allowed: registeredStyleNames()},
			{Key: "font", Type: "string", Summary: "Font selection (auto or a registered font ID).", Default: "auto"},
		},
	})
	catalog.RegisterCommand(catalog.CommandDefinition{
		Verb:    "serial",
		Summary: "Inspect and configure the live serial monitor.",
		Usage:   "serial [key=value ...] | serial clear",
		Handle:  HandleConsoleCommand,
	})
}

// portValidator accepts empty string or any non-whitespace path.
func portValidator(value string) string {
	v := strings.TrimSpace(value)
	if v == "" {
		return "" // empty is valid (enables auto-select)
	}
	if strings.ContainsAny(v, " \t\n\r") {
		return "must be empty or a path without whitespace"
	}
	return ""
}

// fontValidator accepts "auto" or any non-empty trimmed string.
func fontValidator(value string) string {
	v := strings.TrimSpace(value)
	if v == "" {
		return "must be \"auto\" or a non-empty font ID"
	}
	return ""
}

// getPolicyValue reads the current value for a given key from policy.
func getPolicyValue(key string) string {
	p := PolicySnapshot()
	switch key {
	case "port":
		return p.Port
	case "baud":
		return strconv.Itoa(p.Baud)
	case "lines":
		return strconv.Itoa(p.MaxLines)
	case "autoselect":
		if p.AutoSelect {
			return "true"
		}
		return "false"
	case "scan_ms":
		return strconv.Itoa(p.ScanMS)
	case "style":
		return p.Style
	case "font":
		return p.Font
	default:
		return ""
	}
}

// applyPolicyValue updates the policy atomically for a single key.
func applyPolicyValue(key, value string) {
	p := PolicySnapshot()
	switch key {
	case "port":
		p.Port = strings.TrimSpace(value)
	case "baud":
		if n, err := strconv.Atoi(strings.TrimSpace(value)); err == nil {
			p.Baud = n
		}
	case "lines":
		if n, err := strconv.Atoi(strings.TrimSpace(value)); err == nil {
			p.MaxLines = n
		}
	case "autoselect":
		if v, ok := cmdutil.ParseBool(value); ok {
			p.AutoSelect = v
		}
	case "scan_ms":
		if n, err := strconv.Atoi(strings.TrimSpace(value)); err == nil {
			p.ScanMS = n
		}
	case "style":
		p.Style = strings.ToLower(strings.TrimSpace(value))
	case "font":
		p.Font = strings.TrimSpace(value)
	}
	SetPolicy(p)
}

// HandleConsoleCommand processes the serial command using CmdHandler.
// Retains the "clear" sub-command for buffer clearing.
func HandleConsoleCommand(args []string) string {
	if len(args) > 0 && strings.ToLower(args[0]) == "clear" {
		source.Clear()
		return "OK serial cleared"
	}
	return handler.Handle(args)
}
