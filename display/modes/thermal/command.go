package thermal

import (
	"strconv"
	"strings"

	"github.com/databeast/cyberhud/display/catalog/cmdutil"
	"github.com/databeast/cyberhud/display/modes/thermal/source"
)

// fontValidator accepts "auto" or any non-empty trimmed string.
// Whitespace-only values are rejected.
func fontValidator(value string) string {
	v := strings.TrimSpace(value)
	if v == "" {
		return "must be \"auto\" or a non-empty font ID"
	}
	return ""
}

// cmdHandler is the declarative CmdHandler for the "thermal" command verb.
var cmdHandler = &cmdutil.CmdHandler{
	Verb: "thermal",
	Keys: []cmdutil.KeyDef{
		{Name: "style", Validate: cmdutil.AllowedValidator(allowedStyleNames())},
		{Name: "font", Validate: fontValidator},
		{Name: "refresh_ms", Validate: cmdutil.IntValidator(500)},
		{Name: "warn_threshold", Validate: cmdutil.IntValidator(0)},
		{Name: "crit_threshold", Validate: cmdutil.IntValidator(0)},
		{Name: "unit", Validate: cmdutil.AllowedValidator(source.AllowedUnits)},
		{Name: "fgcolor", Validate: cmdutil.AllowedValidator(source.AllowedFGColors)},
		{Name: "show_led", Validate: cmdutil.BoolValidator()},
		{Name: "show_refresh_bar", Validate: cmdutil.BoolValidator()},
		{Name: "show_border", Validate: cmdutil.BoolValidator()},
	},
	Get:       getPolicy,
	Apply:     applyPolicy,
	PostApply: fitnessNotesPostApply,
}

// getPolicy returns the current value for a given policy key.
func getPolicy(key string) string {
	p := GetPolicy()
	switch key {
	case "style":
		return p.Style
	case "font":
		return p.Font
	case "refresh_ms":
		return strconv.Itoa(p.RefreshMS)
	case "warn_threshold":
		return strconv.Itoa(p.WarnThreshold)
	case "crit_threshold":
		return strconv.Itoa(p.CritThreshold)
	case "unit":
		return p.Unit
	case "fgcolor":
		return p.FGColor
	case "show_led":
		if p.ShowLED {
			return "true"
		}
		return "false"
	case "show_refresh_bar":
		if p.ShowRefreshBar {
			return "true"
		}
		return "false"
	case "show_border":
		if p.ShowBorder {
			return "true"
		}
		return "false"
	default:
		return ""
	}
}

// applyPolicy updates a single policy key with its validated value.
func applyPolicy(key, value string) {
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
	case "warn_threshold":
		if n, err := strconv.Atoi(strings.TrimSpace(value)); err == nil {
			policyState.policy.WarnThreshold = n
		}
	case "crit_threshold":
		if n, err := strconv.Atoi(strings.TrimSpace(value)); err == nil {
			policyState.policy.CritThreshold = n
		}
	case "unit":
		policyState.policy.Unit = strings.ToUpper(strings.TrimSpace(value))
	case "fgcolor":
		policyState.policy.FGColor = strings.ToLower(strings.TrimSpace(value))
	case "show_led":
		if v, ok := cmdutil.ParseBool(value); ok {
			policyState.policy.ShowLED = v
		}
	case "show_refresh_bar":
		if v, ok := cmdutil.ParseBool(value); ok {
			policyState.policy.ShowRefreshBar = v
		}
	case "show_border":
		if v, ok := cmdutil.ParseBool(value); ok {
			policyState.policy.ShowBorder = v
		}
	}
}

// HandleCommand is the catalog command handler for the "thermal" verb.
// It delegates to the declarative CmdHandler but adds cross-field validation
// via SetPolicy: if the resulting policy fails normalization (e.g. threshold
// rejection), the change is rolled back and the error is surfaced to the caller.
func HandleCommand(args []string) string {
	if len(args) == 0 {
		return cmdHandler.Handle(args)
	}

	// Snapshot the current policy before applying changes.
	policyState.RLock()
	snapshot := policyState.policy
	policyState.RUnlock()

	result := cmdHandler.Handle(args)

	// If the handler returned an error, no changes were applied — return as-is.
	if strings.HasPrefix(result, "ERR") {
		return result
	}

	// Run the mutated policy through SetPolicy for full normalization and
	// threshold validation. If SetPolicy rejects it, rollback and surface the error.
	policyState.RLock()
	p := policyState.policy
	policyState.RUnlock()

	if err := SetPolicy(p); err != nil {
		// Rollback to pre-change state.
		policyState.Lock()
		policyState.policy = snapshot
		policyState.Unlock()
		return "ERR " + err.Error()
	}

	return result
}
