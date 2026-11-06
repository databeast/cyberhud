package storage

import (
	"strconv"
	"strings"

	"github.com/databeast/cyberhud/display/catalog/cmdutil"
	"github.com/databeast/cyberhud/display/modes/storage/source"
)

// cmdHandler is the declarative CmdHandler for the "storage" command verb.
var cmdHandler = &cmdutil.CmdHandler{
	Verb: "storage",
	Keys: []cmdutil.KeyDef{
		{Name: "style", Validate: cmdutil.AllowedValidator(allowedStyleNames())},
		{Name: "font", Validate: cmdutil.FontValidator()},
		{Name: "refresh_ms", Validate: cmdutil.IntValidator(500)},
		{Name: "unit", Validate: cmdutil.AllowedValidator(source.AllowedUnits)},
		{Name: "show_all", Validate: cmdutil.BoolValidator()},
	},
	Get:       getPolicyValue,
	Apply:     applyPolicyValue,
	PostApply: fitnessNotesPostApply,
}

// getPolicyValue returns the current value for a given policy key.
func getPolicyValue(key string) string {
	p := GetPolicy()
	switch key {
	case "style":
		return p.Style
	case "font":
		return p.Font
	case "refresh_ms":
		return strconv.Itoa(p.RefreshMS)
	case "unit":
		return p.Unit
	case "show_all":
		return boolStr(p.ShowAll)
	default:
		return ""
	}
}

// applyPolicyValue updates a single policy key with its validated value, then
// re-normalizes the whole policy so clamping and registry checks apply.
func applyPolicyValue(key, value string) {
	policyState.Lock()
	defer policyState.Unlock()
	switch key {
	case "style":
		policyState.policy.Style = strings.ToLower(strings.TrimSpace(value))
	case "font":
		policyState.policy.Font = strings.TrimSpace(value)
	case "refresh_ms":
		if n, err := strconv.Atoi(strings.TrimSpace(value)); err == nil {
			policyState.policy.RefreshMS = n
		}
	case "unit":
		policyState.policy.Unit = strings.ToLower(strings.TrimSpace(value))
	case "show_all":
		if v, ok := cmdutil.ParseBool(value); ok {
			policyState.policy.ShowAll = v
		}
	}
	policyState.policy = normalizePolicy(policyState.policy)
}

// HandleCommand is the catalog command handler for the "storage" verb.
func HandleCommand(args []string) string {
	return cmdHandler.Handle(args)
}

func boolStr(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
