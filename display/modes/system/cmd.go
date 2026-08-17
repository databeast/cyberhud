package system

import (
	"fmt"
	"strings"
)

// KeyDef and CmdHandler types are defined in the parent displaymodes package.
// To avoid an import cycle (system → displaymodes → registry → system),
// the command handler is registered via catalog directly with a local implementation.

// HandleCommand is the catalog command handler for the "system" verb.
// It follows the uniform command pattern: zero args returns current state,
// key=value args validate and apply atomically.
func HandleCommand(args []string) string {
	if len(args) == 0 {
		return queryResponse()
	}

	// Parse key=value pairs.
	type kv struct {
		key string
		val string
	}
	pairs := make([]kv, 0, len(args))
	for _, arg := range args {
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) != 2 {
			return fmt.Sprintf("ERR unknown key %s", arg)
		}
		pairs = append(pairs, kv{key: strings.TrimSpace(parts[0]), val: strings.TrimSpace(parts[1])})
	}

	// Validate all pairs before applying any.
	for _, p := range pairs {
		switch strings.ToLower(p.key) {
		case "style":
			lower := strings.ToLower(p.val)
			if systemRegistry.Lookup(lower) == nil {
				return fmt.Sprintf("ERR style: must be one of [%s]", strings.Join(registeredStyleNames(), ", "))
			}
		case "font":
			if reason := fontValidator(p.val); reason != "" {
				return fmt.Sprintf("ERR font: %s", reason)
			}
		default:
			return fmt.Sprintf("ERR unknown key %s", p.key)
		}
	}

	// All valid — apply atomically, tracking which keys were applied.
	appliedKeys := make([]string, 0, len(pairs))
	policyState.Lock()
	for _, p := range pairs {
		switch strings.ToLower(p.key) {
		case "style":
			policyState.policy.Style = strings.ToLower(p.val)
			appliedKeys = append(appliedKeys, "style")
		case "font":
			policyState.policy.Font = p.val
			appliedKeys = append(appliedKeys, "font")
		}
	}
	policyState.Unlock()

	resp := queryResponse()
	if notes := fitnessNotesPostApply(appliedKeys); len(notes) > 0 {
		resp += "\n" + strings.Join(notes, "\n")
	}
	return resp
}
